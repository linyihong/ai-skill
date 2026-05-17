#!/usr/bin/env ruby
# frozen_string_literal: true

# Runtime Compiler Engine
# 將 canonical prose source 與 structured YAML 編譯為 runtime/runtime.db (SQLite)。
#
# 使用方式：
#   ruby runtime/compiler/compiler-engine.rb          # 編譯所有 source → runtime.db
#   ruby runtime/compiler/compiler-engine.rb --check  # 只檢查是否需要編譯
#   ruby runtime/compiler/compiler-engine.rb --diff   # 顯示預期變更
#   ruby runtime/compiler/compiler-engine.rb --db PATH  # 指定輸出 SQLite 路徑
#
# 設計原則：
# - Deterministic：相同輸入 → 相同輸出
# - Idempotent：重複執行不改變結果
# - 只編譯 deterministic state，不編譯 judgment/heuristics
# - SQLite 是 execution surface，Markdown/YAML 是 source-of-truth

require 'yaml'
require 'fileutils'
require 'time'
require 'optparse'
require 'open3'
require 'json'

COMPILER_VERSION = '1.1.0'
GENERATED_DIR = File.join(File.dirname(__FILE__), '..', 'generated')
COMPILER_RULES = File.join(File.dirname(__FILE__), 'compiler-rules.yaml')
DEFAULT_DB_PATH = File.join(File.dirname(__FILE__), '..', 'runtime.db')

@mapping = nil

def load_mapping
  rules = YAML.safe_load(File.read(COMPILER_RULES), permitted_classes: [Date])
  @mapping = rules['source_target_mapping']
end

def generated_header(source_path)
  { 'generated_from' => source_path, 'generated_at' => Time.now.utc.strftime('%Y-%m-%dT%H:%M:%SZ'), 'compiler_version' => COMPILER_VERSION, 'status' => 'synced' }
end

def target_path_for(source_path, mapping_entry)
  domain = extract_domain(source_path)
  target = mapping_entry['target'].gsub('{domain}', domain)
  target.include?('/') ? target : File.join(GENERATED_DIR, File.basename(target))
end

def extract_domain(source_path)
  m = source_path.match(%r{workflow/([^/]+)/}); return m[1] if m
  m = source_path.match(%r{enforcement/(.+)\.md$}); return m[1].tr('-', '_') if m
  return 'plans-index' if source_path.match?(%r{plans/})
  return 'apk-workflow' if source_path.match?(%r{analysis/apk/workflows/})
  'unknown'
end

def extract_phase_sections(content)
  sections = []
  content.scan(/^(##\s+\d+\.\s+.+)$/) do
    heading_text = $1.strip
    heading_line = $`.lines.count + 1
    section_start = $`.size
    remaining = content[section_start + $&.size..]
    next_section_match = remaining.match(/^##\s+\d+\.\s+/)
    section_body = next_section_match ? remaining[0...next_section_match.begin(0)] : remaining
    sections << { 'heading' => heading_text.sub(/^##\s+/, '').strip, 'line_number' => heading_line, 'body' => section_body.strip }
  end
  sections
end

def extract_allowed_actions_from_section(body)
  actions = []
  body.scan(/^\|\s*\d+\s*\|\s*(.+?)\s*\|$/) { |m| actions << { 'action' => m[0].strip, 'source' => 'table' } }
  body.scan(/^\|\s*(?:[^|]+)\s*\|\s*(.+?)\s*\|$/) do |m|
    c2 = m[0].strip; next if %w[--- 必要行動 Action 行動 緩解措施 回填要求 證明].include?(c2) || c2.start_with?('`<') || c2 =~ /^(步驟|Step|Reset|測試類型|文件|問題|根本原因|用途)/ || c2.length < 10 || c2 =~ /^\|/
    actions << { 'action' => c2, 'source' => 'table' }
  end
  body.scan(/^-\s+(.+)$/) { |m| actions << { 'action' => m[0].strip, 'source' => 'bullet' } unless m[0].start_with?('`<') || m[0] =~ /^(不要|禁止|Don't|Do not|Never|避免)/ }
  body.scan(/^\d+\.\s+\*\*([^*]+)\*\*[：:]\s*(.+)$/) { |m| actions << { 'action' => "#{m[0].strip}：#{m[1].strip}", 'source' => 'step' } }
  actions.uniq { |a| a['action'] }
end

def extract_blocking_gates_from_section(body)
  gates = []
  body.scan(/^-\s+(不要|禁止|Don't|Do not|Never|避免)\s*(.+)$/i) { |p, r| gates << { 'id' => "gate.#{r.strip.downcase.gsub(/[^a-z0-9]+/, '_').gsub(/^_|_$/, '')}", 'description' => "#{p}#{r}", 'severity' => 'critical', 'source' => 'prohibition' } }
  body.scan(/(.+?)(?:是|為|屬於)\s*阻擋項/) { |m| gates << { 'id' => "gate.blocking.#{m[0].strip.downcase.gsub(/[^a-z0-9]+/, '_').gsub(/^_|_$/, '')}", 'description' => "#{m[0].strip}是阻擋項", 'severity' => 'critical', 'source' => 'blocking_condition' } }
  body.scan(/(?:在|於|在).+?(?:之前|前).+?(?:必須|需要|應).+?。/) { |m| gates << { 'id' => "gate.prerequisite.#{m[0..30].strip.downcase.gsub(/[^a-z0-9]+/, '_').gsub(/^_|_$/, '')}", 'description' => m[0].strip, 'severity' => 'high', 'source' => 'prerequisite' } }
  gates.uniq { |g| g['id'] }
end

def extract_tables_from_section(body)
  tables = []; lines = body.split("\n"); ct = []; in_t = false
  lines.each do |line|
    if line.match?(/^\|.+\|$/)
      ct << line.strip; in_t = true
    else
      if in_t && ct.length >= 3
        h = ct[0].split('|').map(&:strip).reject(&:empty?)
        r = ct[2..].map { |x| x.split('|').map(&:strip).reject(&:empty?) }
        tables << { 'header' => h, 'rows' => r }
      end
      ct = []; in_t = false
    end
  end
  if in_t && ct.length >= 3
    h = ct[0].split('|').map(&:strip).reject(&:empty?)
    r = ct[2..].map { |x| x.split('|').map(&:strip).reject(&:empty?) }
    tables << { 'header' => h, 'rows' => r }
  end
  tables
end

# ═══════════════════════════════════════════════════════════════
# SQLite Helpers
# ═══════════════════════════════════════════════════════════════

def sqlite_exec(db_path, sql)
  o, e, s = Open3.capture3('sqlite3', db_path, stdin_data: sql)
  warn "sqlite3 error: #{e.strip}" unless s.success?
  s.success? ? o : nil
end

def sqe(v); v.nil? ? 'NULL' : "'#{v.to_s.gsub("'", "''")}'"; end
def jsn(v); v.nil? ? 'NULL' : "'#{((v.is_a?(String) ? v : v.to_json)).gsub("'", "''")}'"; end

def insert_gs(db_path, source_path, target_key, compile_rule, data)
  sqlite_exec(db_path, "INSERT OR REPLACE INTO generated_surfaces (source_path, target_key, compile_rule, compiled_at, compiler_version, status, data) VALUES (#{sqe(source_path)}, #{sqe(target_key)}, #{sqe(compile_rule)}, datetime('now'), #{sqe(COMPILER_VERSION)}, 'synced', #{jsn(data)});")
end

# ═══════════════════════════════════════════════════════════════
# Schema Creation
# ═══════════════════════════════════════════════════════════════

def create_runtime_db_schema(db_path)
  sqlite_exec(db_path, <<~SQL)
    CREATE TABLE IF NOT EXISTS phases (
      id TEXT PRIMARY KEY, name TEXT NOT NULL, description TEXT,
      entry_conditions TEXT, allowed_actions TEXT, forbidden_actions TEXT,
      blocking_gates TEXT, obligations TEXT, next_phase TEXT,
      phase_transition_triggers TEXT, metadata TEXT,
      created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now'))
    );
    CREATE TABLE IF NOT EXISTS phase_transitions (
      id INTEGER PRIMARY KEY AUTOINCREMENT, from_phase TEXT NOT NULL, to_phase TEXT NOT NULL,
      trigger TEXT, condition TEXT, rule_type TEXT DEFAULT 'normal', description TEXT,
      FOREIGN KEY (from_phase) REFERENCES phases(id), FOREIGN KEY (to_phase) REFERENCES phases(id)
    );
    CREATE TABLE IF NOT EXISTS obligations (
      id TEXT PRIMARY KEY, phase TEXT NOT NULL, name TEXT NOT NULL, description TEXT,
      verification TEXT, severity TEXT DEFAULT 'high', depends_on TEXT, linked_gates TEXT, metadata TEXT,
      created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')),
      FOREIGN KEY (phase) REFERENCES phases(id)
    );
    CREATE TABLE IF NOT EXISTS gates (
      id TEXT PRIMARY KEY, phase TEXT NOT NULL, name TEXT NOT NULL, description TEXT,
      severity TEXT DEFAULT 'high', check_type TEXT, check_target TEXT, check_verification TEXT,
      failure_action TEXT, failure_message TEXT, metadata TEXT,
      created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')),
      FOREIGN KEY (phase) REFERENCES phases(id)
    );
    CREATE TABLE IF NOT EXISTS transaction_states (
      state TEXT PRIMARY KEY, description TEXT, entry_conditions TEXT, allowed_actions TEXT,
      forbidden_actions TEXT, blocking_gates TEXT, metadata TEXT,
      created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now'))
    );
    CREATE TABLE IF NOT EXISTS transaction_transitions (
      id INTEGER PRIMARY KEY AUTOINCREMENT, from_state TEXT NOT NULL, to_state TEXT NOT NULL,
      trigger TEXT, condition TEXT, description TEXT,
      FOREIGN KEY (from_state) REFERENCES transaction_states(state), FOREIGN KEY (to_state) REFERENCES transaction_states(state)
    );
    CREATE TABLE IF NOT EXISTS transaction_rules (
      id TEXT PRIMARY KEY, name TEXT NOT NULL, description TEXT, rule TEXT, severity TEXT DEFAULT 'high'
    );
    CREATE TABLE IF NOT EXISTS transaction_templates (
      id TEXT PRIMARY KEY, name TEXT NOT NULL, description TEXT, typical_steps TEXT
    );
    CREATE TABLE IF NOT EXISTS activation_rules (
      rule_id TEXT PRIMARY KEY, description TEXT, activation_when TEXT,
      load_strategy TEXT DEFAULT 'lazy', load_priority TEXT DEFAULT 'P2',
      load_estimated_tokens INTEGER DEFAULT 0, load_source TEXT, metadata TEXT,
      created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now'))
    );
    CREATE TABLE IF NOT EXISTS core_bootstrap_rules (
      rule_id TEXT PRIMARY KEY, ordinal INTEGER DEFAULT 0
    );
    CREATE TABLE IF NOT EXISTS discovery_checkpoints (
      id INTEGER PRIMARY KEY AUTOINCREMENT, phase TEXT NOT NULL, trigger TEXT NOT NULL,
      description TEXT, discovery_targets TEXT, metadata TEXT,
      FOREIGN KEY (phase) REFERENCES phases(id)
    );
    CREATE TABLE IF NOT EXISTS discovery_search_strategy (
      id INTEGER PRIMARY KEY AUTOINCREMENT, priority_order TEXT, fallback TEXT, min_confidence_threshold TEXT
    );
    CREATE TABLE IF NOT EXISTS compiler_metadata (key TEXT PRIMARY KEY, value TEXT NOT NULL);
    CREATE TABLE IF NOT EXISTS generated_surfaces (
      id INTEGER PRIMARY KEY AUTOINCREMENT, source_path TEXT NOT NULL, target_key TEXT NOT NULL,
      compile_rule TEXT, compiled_at TEXT DEFAULT (datetime('now')), compiler_version TEXT,
      status TEXT DEFAULT 'synced', data TEXT, UNIQUE(source_path, target_key)
    );
    CREATE INDEX IF NOT EXISTS idx_obligations_phase ON obligations(phase);
    CREATE INDEX IF NOT EXISTS idx_gates_phase ON gates(phase);
    CREATE INDEX IF NOT EXISTS idx_discovery_checkpoints_phase ON discovery_checkpoints(phase);
    CREATE INDEX IF NOT EXISTS idx_phase_transitions_from ON phase_transitions(from_phase);
    CREATE INDEX IF NOT EXISTS idx_transaction_transitions_from ON transaction_transitions(from_state);
  SQL
end

# ═══════════════════════════════════════════════════════════════
# Build Runtime DB from Structured YAML
# ═══════════════════════════════════════════════════════════════

def build_runtime_db(db_path)
  puts "  Building runtime.db from structured YAML..."
  %w[phases phase_transitions obligations gates transaction_states transaction_transitions transaction_rules transaction_templates activation_rules core_bootstrap_rules discovery_checkpoints discovery_search_strategy generated_surfaces compiler_metadata].each { |t| sqlite_exec(db_path, "DELETE FROM #{t};") }

  rd = File.dirname(File.dirname(__FILE__))

  # 1. Phases
  pm_path = File.join(rd, 'phases', 'phase-machine.yaml')
  if File.exist?(pm_path)
    pm = YAML.safe_load(File.read(pm_path), permitted_classes: [Date])
    (pm['phases'] || []).each do |ph|
      sqlite_exec(db_path, "INSERT INTO phases (id, name, description, entry_conditions, allowed_actions, forbidden_actions, blocking_gates, obligations, next_phase, phase_transition_triggers, metadata) VALUES (#{sqe(ph['id'])}, #{sqe(ph['name'])}, #{sqe(ph['description'])}, #{jsn(ph['entry_conditions'])}, #{jsn(ph['allowed_actions'])}, #{jsn(ph['forbidden_actions'])}, #{jsn(ph['blocking_gates'])}, #{jsn(ph['obligations'])}, #{sqe(ph['next_phase'])}, #{jsn(ph['phase_transition_triggers'])}, #{jsn(ph['metadata'] || {})});")
    end
    puts "    ✓ #{pm['phases']&.length || 0} phases"
    (pm['phase_transition_rules'] || []).each do |r|
      sqlite_exec(db_path, "INSERT INTO phase_transitions (from_phase, to_phase, trigger, condition, rule_type, description) VALUES (#{sqe(r['from'])}, #{sqe(r['to'])}, #{sqe(r['trigger'])}, #{sqe(r['condition'])}, #{sqe(r['type'] || 'normal')}, #{sqe(r['description'])});")
    end
    puts "    ✓ #{pm['phase_transition_rules']&.length || 0} phase transitions"
  end

  # 2. Obligations
  ol_path = File.join(rd, 'obligations', 'obligation-ledger.yaml')
  if File.exist?(ol_path)
    ol = YAML.safe_load(File.read(ol_path), permitted_classes: [Date])
    (ol['obligations'] || []).each do |ob|
      sqlite_exec(db_path, "INSERT INTO obligations (id, phase, name, description, verification, severity, depends_on, linked_gates, metadata) VALUES (#{sqe(ob['id'])}, #{sqe(ob['phase'])}, #{sqe(ob['name'])}, #{sqe(ob['description'])}, #{jsn(ob['verification'])}, #{sqe(ob['severity'] || 'high')}, #{jsn(ob['depends_on'])}, #{jsn(ob['linked_gates'])}, #{jsn(ob['metadata'] || {})});")
    end
    puts "    ✓ #{ol['obligations']&.length || 0} obligations"
  end

  # 3. Gates
  bg_path = File.join(rd, 'gates', 'blocking-gates.yaml')
  if File.exist?(bg_path)
    bg = YAML.safe_load(File.read(bg_path), permitted_classes: [Date])
    (bg['gates'] || []).each do |g|
      ck = g['check'] || {}
      sqlite_exec(db_path, "INSERT INTO gates (id, phase, name, description, severity, check_type, check_target, check_verification, failure_action, failure_message, metadata) VALUES (#{sqe(g['id'])}, #{sqe(g['phase'])}, #{sqe(g['name'])}, #{sqe(g['description'])}, #{sqe(g['severity'] || 'high')}, #{sqe(ck['type'])}, #{sqe(ck['target'])}, #{sqe(ck['verification'])}, #{sqe(g['failure_action'])}, #{sqe(g['failure_message'])}, #{jsn(g['metadata'] || {})});")
    end
    puts "    ✓ #{bg['gates']&.length || 0} gates"
  end

  # 4. Transaction Machine
  tx_path = File.join(rd, 'transactions', 'transaction-machine.yaml')
  if File.exist?(tx_path)
    tx = YAML.safe_load(File.read(tx_path), permitted_classes: [Date])
    (tx['transaction_states'] || []).each do |st|
      sqlite_exec(db_path, "INSERT INTO transaction_states (state, description, entry_conditions, allowed_actions, forbidden_actions, blocking_gates, metadata) VALUES (#{sqe(st['state'])}, #{sqe(st['description'])}, #{jsn(st['entry_conditions'])}, #{jsn(st['allowed_actions'])}, #{jsn(st['forbidden_actions'])}, #{jsn(st['blocking_gates'])}, #{jsn(st['metadata'] || {})});")
      (st['transitions'] || []).each do |tr|
        sqlite_exec(db_path, "INSERT INTO transaction_transitions (from_state, to_state, trigger, condition, description) VALUES (#{sqe(st['state'])}, #{sqe(tr['to'])}, #{sqe(tr['trigger'])}, #{sqe(tr['condition'])}, #{sqe(tr['description'])});")
      end
    end
    puts "    ✓ #{tx['transaction_states']&.length || 0} transaction states"
    (tx['transaction_rules'] || []).each { |r| sqlite_exec(db_path, "INSERT INTO transaction_rules (id, name, description, rule, severity) VALUES (#{sqe(r['id'])}, #{sqe(r['name'])}, #{sqe(r['description'])}, #{sqe(r['rule'])}, #{sqe(r['severity'] || 'high')});") }
    (tx['transaction_templates'] || []).each { |t| sqlite_exec(db_path, "INSERT INTO transaction_templates (id, name, description, typical_steps) VALUES (#{sqe(t['id'])}, #{sqe(t['name'])}, #{sqe(t['description'])}, #{jsn(t['typical_steps'])});") }
    puts "    ✓ #{tx['transaction_rules']&.length || 0} rules, #{tx['transaction_templates']&.length || 0} templates"
  end

  # 5. Activation Rules
  ar_path = File.join(rd, 'router', 'activation-rules.yaml')
  if File.exist?(ar_path)
    ar = YAML.safe_load(File.read(ar_path), permitted_classes: [Date])
    (ar['core_bootstrap'] || []).each_with_index { |rid, i| sqlite_exec(db_path, "INSERT INTO core_bootstrap_rules (rule_id, ordinal) VALUES (#{sqe(rid)}, #{i});") }
    (ar['rules'] || []).each do |r|
      ac = r['activation'] || {}; ld = r['load'] || {}
      sqlite_exec(db_path, "INSERT INTO activation_rules (rule_id, description, activation_when, load_strategy, load_priority, load_estimated_tokens, load_source, metadata) VALUES (#{sqe(r['rule_id'])}, #{sqe(r['description'])}, #{jsn(ac['when'])}, #{sqe(ld['strategy'] || 'lazy')}, #{sqe(ld['priority'] || 'P2')}, #{ld['estimated_tokens'] || 0}, #{sqe(ld['source'])}, #{jsn(r['metadata'] || {})});")
    end
    puts "    ✓ #{ar['core_bootstrap']&.length || 0} core, #{ar['rules']&.length || 0} activation rules"
  end

  # 6. Discovery Checkpoints
  dc_path = File.join(rd, 'discovery', 'capability-checkpoints.yaml')
  if File.exist?(dc_path)
    dc = YAML.safe_load(File.read(dc_path), permitted_classes: [Date])
    (dc['checkpoints'] || []).each do |cp|
      sqlite_exec(db_path, "INSERT INTO discovery_checkpoints (phase, trigger, description, discovery_targets, metadata) VALUES (#{sqe(cp['phase'])}, #{sqe(cp['trigger'])}, #{sqe(cp['description'])}, #{jsn(cp['discovery_targets'])}, #{jsn(cp['metadata'] || {})});")
    end
    puts "    ✓ #{dc['checkpoints']&.length || 0} discovery checkpoints"
    st = dc['search_strategy'] || {}
    sqlite_exec(db_path, "INSERT INTO discovery_search_strategy (priority_order, fallback, min_confidence_threshold) VALUES (#{jsn(st['priority_order'])}, #{jsn(st['fallback'])}, #{sqe(st['min_confidence_threshold'])});") unless st.empty?
  end

  # 7. Compiler Metadata
  sqlite_exec(db_path, "INSERT INTO compiler_metadata (key, value) VALUES ('compiler_version', #{sqe(COMPILER_VERSION)});")
  sqlite_exec(db_path, "INSERT INTO compiler_metadata (key, value) VALUES ('compiled_at', #{sqe(Time.now.utc.strftime('%Y-%m-%dT%H:%M:%SZ'))});")
  sqlite_exec(db_path, "INSERT INTO compiler_metadata (key, value) VALUES ('schema_version', '1.0');")
  puts "  ✓ compiler metadata written"
end

# ═══════════════════════════════════════════════════════════════
# Prose → SQLite Compilation Functions
# ═══════════════════════════════════════════════════════════════

def compile_workflow_phases_to_sqlite(db_path, source_path, mapping_entry)
  content = File.read(source_path)
  sections = extract_phase_sections(content)
  phases = sections.map do |sec|
    allowed = extract_allowed_actions_from_section(sec['body'])
    gates = extract_blocking_gates_from_section(sec['body'])
    tables = extract_tables_from_section(sec['body'])
    pe = { 'name' => sec['heading'], 'source_line' => sec['line_number'] }
    pe['allowed_actions'] = allowed unless allowed.empty?
    pe['blocking_gates'] = gates unless gates.empty?
    pe['tables'] = tables unless tables.empty?
    pe
  end
  ag = phases.flat_map { |p| p['blocking_gates'] || [] }.uniq { |g| g['id'] }
  data = { 'phases' => phases }; data['gates'] = ag unless ag.empty?
  insert_gs(db_path, source_path, "workflow.#{extract_domain(source_path)}.phases", mapping_entry['compile_rule'], data)
  puts "  ✓ #{source_path}"
end

def compile_enforcement_transactions_to_sqlite(db_path, source_path, mapping_entry)
  content = File.read(source_path)
  states = []
  content.scan(/^###?\s+(.+?)(?:\s*\((.+?)\))?$/) { |m| states << { 'name' => m[0].strip, 'context' => m[1]&.strip } }
  insert_gs(db_path, source_path, "enforcement.#{extract_domain(source_path)}.transactions", mapping_entry['compile_rule'], { 'compiled_from' => source_path, 'states' => states })
  puts "  ✓ #{source_path}"
end

def compile_output_governance_to_sqlite(db_path, source_path, mapping_entry)
  content = File.read(source_path)
  rules = []; content.scan(/^###?\s+(.+?)$/) { |m| rules << { 'section' => m[0].strip } }
  gates = []; content.scan(/^\*\*([^*]+)\*\*：(.+)$/) { |m| gates << { 'name' => m[0].strip, 'description' => m[1].strip } }
  insert_gs(db_path, source_path, "governance.#{extract_domain(source_path)}", mapping_entry['compile_rule'], { 'compiled_from' => source_path, 'rules' => rules, 'gates' => gates })
  puts "  ✓ #{source_path}"
end

def compile_knowledge_update_flow_to_sqlite(db_path, source_path, _mapping_entry)
  content = File.read(source_path)
  steps = []
  content.scan(/^##\s+Step\s+(\d+)[：:]\s*(.+?)$/) do |m|
    sn = m[0].strip.to_i; sname = m[1].strip
    ss = content.index("## Step #{sn}："); next unless ss
    rem = content[ss..]; nsm = rem.index(/^## Step #{sn + 1}[：:]/)
    sc = nsm ? rem[0...nsm] : rem
    ec = []; sc.scan(/^\|\s*(\w[\w\s]+?)\s*\|\s*(.+?)\s*\|$/) { |r| ec << { 'condition' => r[0].strip, 'next_step' => r[1].strip } }
    refs = []; sc.scan(/\[`([^`]+)`\]\(([^)]+)\)/) { |r| refs << { 'name' => r[0], 'path' => r[1] } }
    steps << { 'step' => sn, 'name' => sname, 'entry_conditions' => ec, 'references' => refs }
  end
  insert_gs(db_path, source_path, 'governance.knowledge_update_flow', '從 knowledge-update-flow.md 的 11 個步驟標題與判斷表格提取 phase 定義', { 'compiled_from' => source_path, 'total_steps' => 11, 'steps' => steps })
  puts "  ✓ #{source_path}"
end

def compile_workflow_artifacts_to_sqlite(db_path, source_path, mapping_entry)
  content = File.read(source_path)
  artifacts = []; content.scan(/^##\s+\d+\.\s+(.+)$/) { |m| artifacts << { 'name' => m[0].strip } }
  gates = []; content.scan(/^###\s+(.+?)(?:\s*Gate|gate)?$/) { |m| gates << { 'name' => m[0].strip, 'type' => 'verification_gate' } }
  tables = []; content.scan(/^\|.+\|.+\|$/) { |l| next if l.match?(/^\|[\s-]+\|[\s-]+\|$/) || l.match?(/^\|.*#.*\|$/); tables << l.strip }
  ri = []; content.scan(/^\d+\.\s+\*\*([^*]+)\*\*(.*)$/) { |m| ri << { 'name' => m[0].strip, 'description' => m[1]&.strip } }
  insert_gs(db_path, source_path, "workflow.#{extract_domain(source_path)}.artifacts", mapping_entry['compile_rule'], { 'compiled_from' => source_path, 'artifacts' => artifacts, 'verification_gates' => gates, 'required_items' => ri })
  puts "  ✓ #{source_path}"
end

def compile_goal_action_gates_to_sqlite(db_path, source_path, mapping_entry)
  content = File.read(source_path)
  cf = []; sc = []; ve = []
  content.scan(/^\|\s*(.+?)\s*\|\s*(.+?)\s*\|$/) do |m|
    c1 = m[0].strip; c2 = m[1].strip
    next if %w[--- 欄位 必填內容 情境 要求 工作類型 驗證方式 工作單元].include?(c1) || c1 == '目標' || c1 == '執行' || c1 == '驗證' || c1 == '驗證 / 參考來源' || c1.start_with?('`<')
    if %w[目標 執行 驗證].include?(c1); cf << { 'field' => c1, 'description' => c2 }
    elsif c1.length > 2 && c2.length > 2; sc << { 'scenario' => c1, 'requirement' => c2 } end
  end
  gs = []; gs_section = content[/### 驗證 Gate 參考\n(.+?)(?=\n## |\n### |\z)/m]
  gs_section&.scan(/^\d+\.\s+(.+)$/) { |m| gs << { 'rule' => m[0].strip } }
  insert_gs(db_path, source_path, "enforcement.#{extract_domain(source_path)}.gates", mapping_entry['compile_rule'], { 'compiled_from' => source_path, 'core_fields' => cf, 'scenarios' => sc, 'gates' => gs })
  puts "  ✓ #{source_path}"
end

def compile_failure_recovery_to_sqlite(db_path, source_path, mapping_entry)
  content = File.read(source_path)
  tx = []; content.scan(/^\|\s*`([^`]+)`\s*\|\s*(.+?)\s*\|\s*(.+?)\s*\|$/) { |m| tx << { 'class' => m[0].strip, 'meaning' => m[1].strip, 'common_prevention' => m[2].strip } }
  sr = []; content.scan(/^\|\s*(.+?)\s*\|\s*(.+?)\s*\|$/) { |m| sr << { 'content_type' => m[0].strip, 'location' => m[1].strip } unless m[0] == '---' || m[0] == '內容' || m[0] == 'Durable location' }
  pd = []; content.scan(/^\|\s*(.+?)\s*\|\s*(.+?)\s*\|$/) { |m| pd << { 'failure_scope' => m[0].strip, 'promotion_target' => m[1].strip } unless m[0] == '---' || m[0] == 'Failure scope' || m[0] == 'Promotion target' }
  ls = []; content.scan(/^\d+\.\s+\*\*([^*]+)\*\*[：:]\s*(.+)$/) { |m| ls << { 'step' => m[0].strip, 'description' => m[1].strip } }
  insert_gs(db_path, source_path, "enforcement.#{extract_domain(source_path)}.recovery", mapping_entry['compile_rule'], { 'compiled_from' => source_path, 'failure_taxonomy' => tx, 'storage_rules' => sr, 'promotion_decisions' => pd, 'loop_steps' => ls })
  puts "  ✓ #{source_path}"
end

def compile_plans_index_to_sqlite(db_path, source_path, _mapping_entry)
  content = File.read(source_path)
  filename = File.basename(source_path, '.md')
  plan_id = filename.sub(/^\d{4}-\d{2}-\d{2}-\d{4}-/, '')
  status = 'unknown'
  sm = content.match(/>\s*\*\*狀態\*\*[：:]\s*(.+)$/); status = sm[1].strip if sm
  tm = content.match(/^#\s+(.+)$/); title = tm ? tm[1].strip : ''
  pm = content.match(/>\s*\*\*目的\*\*[：:]\s*(.+)$/); purpose = pm ? pm[1].strip : ''
  status = source_path.include?('/archived/') ? 'completed' : 'draft' if status == 'unknown'
  phases = []
  content.scan(/^###\s+Phase\s+(\d+)[：:]\s*(.+?)$/) do |m|
    pn = m[0].strip.to_i; pt = m[1].strip
    pr = 'P?'; prm = pt.match(/[（(]P(\d)/); pr = "P#{prm[1]}" if prm
    phases << { 'phase' => pn, 'title' => pt.gsub(/[（(]P\d[^)）]*[)）]/, '').strip, 'priority' => pr }
  end
  af = []; in_af = false
  content.each_line do |line|
    if line.match?(/^\|\s*檔案\s*\|\s*變更類型\s*\|\s*Phase\s*\|$/); in_af = true; next; end
    if in_af
      break unless line.match?(/^\|.+\|.+\|.+\|$/)
      next if line.match?(/^\|[\s-]+\|[\s-]+\|[\s-]+\|$/)
      cols = line.split('|').map(&:strip).reject(&:empty?)
      af << { 'path' => cols[0], 'change_type' => cols[1], 'phase' => cols[2] } if cols.length >= 3
    end
  end
  { 'plan_id' => plan_id, 'filename' => filename, 'title' => title, 'status' => status, 'purpose' => purpose, 'phases' => phases, 'affected_files' => af }
end

def compile_classification_rules_to_sqlite(db_path, _source_path, _mapping_entry)
  eng_readme_path = 'intelligence/engineering/README.md'
  dimensions = []
  if File.exist?(eng_readme_path)
    content = File.read(eng_readme_path)
    content.scan(/^\|\s*\[`([^`]+)`\]\(([^)]+)\)\s*\|\s*(.+?)\s*\|$/) do |m|
      dn = m[0].strip; rp = m[1].strip; desc = m[2].strip
      next if %w[子目錄 ---].include?(dn)
      subdirs = []
      full_path = File.join(File.dirname(eng_readme_path), dn)
      if Dir.exist?(full_path)
        Dir.entries(full_path).each do |e|
          next if e.start_with?('.') || !File.directory?(File.join(full_path, e))
          subdirs << e
        end
      end
      dimensions << { 'name' => dn, 'description' => desc, 'path' => rp, 'subdirectories' => subdirs.sort }
    end
  end

  lang_readme_path = 'intelligence/engineering/language-specific/README.md'
  known_languages = []
  if File.exist?(lang_readme_path)
    content = File.read(lang_readme_path)
    in_lt = false
    content.each_line do |line|
      if line.match?(/^\|\s*Language\s*\|\s*Directory\s*\|\s*Atoms\s*\|$/); in_lt = true; next; end
      if in_lt
        break unless line.match?(/^\|.+\|.+\|.+\|$/)
        next if line.match?(/^\|[\s-]+\|[\s-]+\|[\s-]+\|$/)
        cols = line.split('|').map(&:strip).reject(&:empty?)
        known_languages << { 'name' => cols[0], 'path' => "intelligence/engineering/language-specific/#{cols[0].downcase}/", 'atoms' => cols[2].to_s } if cols.length >= 3
      end
    end
  end

  known_frameworks = []; known_platforms = []
  eng_dir = 'intelligence/engineering'
  if Dir.exist?(eng_dir)
    Dir.entries(eng_dir).each do |e|
      next if e.start_with?('.'); fp = File.join(eng_dir, e); next unless File.directory?(fp)
      if e == 'framework-specific'
        Dir.entries(fp).each { |fw| next if fw.start_with?('.'); known_frameworks << { 'name' => fw, 'path' => "#{eng_dir}/#{e}/#{fw}/" } if File.directory?(File.join(fp, fw)) }
      elsif e == 'platform-specific'
        Dir.entries(fp).each { |pl| next if pl.start_with?('.'); known_platforms << { 'name' => pl, 'path' => "#{eng_dir}/#{e}/#{pl}/" } if File.directory?(File.join(fp, pl)) }
      end
    end
  end

  decision_tree = []
  kuf_path = 'governance/lifecycle/knowledge-update-flow.md'
  if File.exist?(kuf_path)
    content = File.read(kuf_path)
    ss = content[/^###\s+2\.4\s.*?\n(.*?)(?=^###\s+2\.5|\z)/m]
    if ss
      ss.scan(/├─\s*(.+?)$/) { |m| decision_tree << { 'branch' => m[0].strip } unless m[0].strip.start_with?('─') || m[0].strip.empty? }
      ss.scan(/→\s*(.+?)$/) { |m| decision_tree << { 'action' => m[0].strip } unless m[0].strip.empty? }
    end
  end

  data = {
    'compiled_from' => 'governance/lifecycle/knowledge-update-flow.md + intelligence/engineering/README.md',
    'classification_dimensions' => dimensions, 'known_languages' => known_languages,
    'known_frameworks' => known_frameworks, 'known_platforms' => known_platforms, 'decision_tree' => decision_tree
  }
  insert_gs(db_path, 'governance/lifecycle/knowledge-update-flow.md', 'governance.classification_rules', '從 knowledge-update-flow.md Step 2.4 的決策樹與 intelligence/engineering/ 的 README 提取分類維度定義', data)
  puts "  ✓ classification-rules"
end

def compile_system_upgrade_governance_to_sqlite(db_path, source_path, _mapping_entry)
  content = File.read(source_path)

  upgrade_conditions = []
  if content =~ /## 1\. 什麼是「大型系統升級」\n\n(.*?)(?:\n\n##|\z)/m
    $1.scan(/^\|\s*(🏷️|🏛️|🔄|📄|🗑️)\s*\*\*(.+?)\*\*\s*\|\s*(.+?)\s*\|\s*(.+?)\s*\|$/) do |emoji, cond, desc, ex|
      upgrade_conditions << { 'emoji' => emoji, 'condition' => cond.strip, 'description' => desc.strip, 'example' => ex.strip }
    end
  end

  checklist_categories = []
  if content =~ /## 2\. 升級計畫書必須包含的檢查清單\n(.*?)(?:\n\n##\s+\d+\.|\z)/m
    $1.split(/^###\s+/).each do |sub|
      next if sub.strip.empty?
      if sub =~ /^(\d+\.\d+)\s+(.+?)\n(.*)/m
        items = []
        $3.scan(/^-\s+\[( |x)\]\s+\*\*(.+?)\*\*(?::|：)\s*(.+?)?$/) { |chk, tit, desc| items << { 'id' => "#{$1.tr('.', '-')}-#{items.length + 1}", 'title' => tit.strip, 'description' => (desc || '').strip, 'required' => true } }
        checklist_categories << { 'id' => $1.strip, 'name' => $2.strip, 'items' => items }
      end
    end
  end

  forced_rules = []
  if content =~ /## 3\. 從三次升級提煉的強制規則\n(.*?)(?:\n\n##\s+\d+\.|\z)/m
    $1.scan(/^###\s+規則\s+(\d+)[：:]\s*(.+?)\n\n\*\*教訓\*\*：(.+?)\n\n\*\*強制\*\*：(.+?)(?=\n\n###|\n\n---|\z)/m) { |num, tit, les, enf| forced_rules << { 'id' => "rule.#{num}", 'title' => tit.strip, 'lesson' => les.strip, 'enforcement' => enf.strip } }
  end

  template_categories = []
  if content =~ /## 5\. 檢查清單範本（可直接複製到計畫書）\n\n```markdown\n(.*?)```/m
    tmpl_content = $1
    tmpl_content.split(/^###\s+/).each do |sub|
      next if sub.strip.empty?
      m = sub.match(/^(.+?)\n(.*)/m)
      next unless m
      name = m[1].strip
      items = []
      m[2].scan(/^-\s+\[( |x)\]\s+(.+?)$/) { |_chk, tit| items << tit.strip }
      template_categories << { 'name' => name, 'items' => items }
    end
  end

  data = { 'compiled_from' => source_path, 'upgrade_conditions' => upgrade_conditions, 'checklist_categories' => checklist_categories, 'forced_rules' => forced_rules, 'checklist_template' => template_categories }
  insert_gs(db_path, source_path, 'governance.system_upgrade', '從系統升級治理要則的 §1 條件表格、§2 檢查清單分類與項目、§3 強制規則提取升級治理定義', data)
  puts "  ✓ #{source_path}"
end

def compile_apk_workflow_phases_to_sqlite(db_path, source_path, _mapping_entry)
  filename = File.basename(source_path, '.md')
  if filename == 'README'
    puts "  - #{source_path} (skipped — README)"
    return nil
  end
  content = File.read(source_path)
  sections = []
  content.scan(/^(##\s+步驟\s+\d+[：:]\s*.+)$/) do
    ht = $1.strip; hl = $`.lines.count + 1
    ss = $`.size; rem = content[ss + $&.size..]
    nsm = rem.match(/^##\s+步驟\s+\d+[：:]/)
    sb = nsm ? rem[0...nsm.begin(0)] : rem
    sections << { 'heading' => ht.sub(/^##\s+/, '').strip, 'line_number' => hl, 'body' => sb.strip }
  end
  steps = sections.map { |s| { 'title' => s['heading'], 'source_line' => s['line_number'] } }
  prereq = []; pm = content.match(/^##\s+前置準備\n(.+?)(?=\n##\s+|\z)/m)
  if pm; pm[1].scan(/^-\s+(.+)$/) { |m| prereq << m[0].strip }; end
  of = nil; om = content.match(/^##\s+成功產出格式\n(.+?)(?=\n##\s+|\z)/m); of = om[1].strip if om
  entry = { 'workflow_name' => filename, 'source_path' => source_path, 'total_steps' => steps.length, 'steps' => steps }
  entry['prerequisites'] = prereq unless prereq.empty?
  entry['output_format'] = of unless of.nil?
  puts "  ✓ #{source_path} (#{steps.length} steps)"
  entry
end

# ═══════════════════════════════════════════════════════════════
# Compile Source Dispatcher
# ═══════════════════════════════════════════════════════════════

def compile_source_to_sqlite(db_path, source_path, mapping_entry)
  rule = mapping_entry['compile_rule']
  case rule
  when /從 execution-flow 的章節標題提取 phase 定義/; compile_workflow_phases_to_sqlite(db_path, source_path, mapping_entry)
  when /從 writeback transaction 章節提取 state machine 定義/; compile_enforcement_transactions_to_sqlite(db_path, source_path, mapping_entry)
  when /提取 language policy 定義|提取 sanitization 定義|提取 tool neutrality 定義/; compile_output_governance_to_sqlite(db_path, source_path, mapping_entry)
  when /從 knowledge-update-flow\.md 的 11 個步驟標題與判斷表格提取 phase 定義/; compile_knowledge_update_flow_to_sqlite(db_path, source_path, mapping_entry)
  when /從 artifact gates 的檢查清單提取 required artifacts 與 verification criteria/; compile_workflow_artifacts_to_sqlite(db_path, source_path, mapping_entry)
  when /從 validation gate 描述提取 gate 定義/; compile_goal_action_gates_to_sqlite(db_path, source_path, mapping_entry)
  when /從 failure taxonomy 與 recovery 描述提取 pattern 與 strategy/; compile_failure_recovery_to_sqlite(db_path, source_path, mapping_entry)
  when /從 plans\/active\/\*\.md 的 front matter、phase 標題、受影響檔案表格提取 plan index/; compile_plans_index_to_sqlite(db_path, source_path, mapping_entry)
  when /從 knowledge-update-flow\.md Step 2\.4 的決策樹與 intelligence\/engineering\/ 的 README 提取分類維度定義/; compile_classification_rules_to_sqlite(db_path, source_path, mapping_entry)
  when /從系統升級治理要則的 §1 條件表格、§2 檢查清單分類與項目、§3 強制規則提取升級治理定義/; compile_system_upgrade_governance_to_sqlite(db_path, source_path, mapping_entry)
  when /從 analysis\/apk\/workflows\/\*\.md 的「步驟 N：」標題提取 step 定義/; compile_apk_workflow_phases_to_sqlite(db_path, source_path, mapping_entry)
  else; puts "  ⚠  Unknown compile rule: #{rule}"
  end
end

# ═══════════════════════════════════════════════════════════════
# Check Modified Sources
# ═══════════════════════════════════════════════════════════════

def check_modified_sources
  modified = []
  plans_modified = false; apk_workflow_modified = false
  plans_target = File.join(GENERATED_DIR, 'plans-index.yaml')
  classification_target = File.join(GENERATED_DIR, 'classification-rules.yaml')
  apk_mapping = @mapping.find { |e| e['compile_rule']&.include?('analysis/apk/workflows/*.md') }
  apk_workflow_target = apk_mapping ? apk_mapping['target'] : File.join(GENERATED_DIR, 'apk-workflow-phases.yaml')

  @mapping.each do |entry|
    Dir.glob(entry['source']).each do |source_path|
      if entry['compile_rule']&.include?('plans/active/*.md')
        plans_modified = true if !File.exist?(plans_target) || File.mtime(source_path) > File.mtime(plans_target)
      elsif entry['compile_rule']&.include?('analysis/apk/workflows/*.md')
        next if File.basename(source_path) == 'README.md'
        apk_workflow_modified = true if !File.exist?(apk_workflow_target) || File.mtime(source_path) > File.mtime(apk_workflow_target)
      elsif entry['compile_rule']&.include?('分類維度定義')
        target = classification_target
        deps = [source_path, 'intelligence/engineering/README.md', 'intelligence/engineering/language-specific/README.md']
        deps.each { |dep| if File.exist?(dep) && (!File.exist?(target) || File.mtime(dep) > File.mtime(target)); modified << { source: source_path, mapping: entry }; break; end }
      else
        target = target_path_for(source_path, entry)
        modified << { source: source_path, mapping: entry } if !File.exist?(target) || File.mtime(source_path) > File.mtime(target)
      end
    end
  end

  if plans_modified; modified << { source: 'plans/active/*.md', mapping: @mapping.find { |e| e['compile_rule']&.include?('plans/active/*.md') } }; end
  if apk_workflow_modified; modified << { source: 'analysis/apk/workflows/*.md', mapping: @mapping.find { |e| e['compile_rule']&.include?('analysis/apk/workflows/*.md') } }; end
  modified
end

# ═══════════════════════════════════════════════════════════════
# Main Run
# ═══════════════════════════════════════════════════════════════

def run(options)
  load_mapping
  db_path = options[:db] || DEFAULT_DB_PATH

  if options[:check]
    modified = check_modified_sources
    if modified.empty?
      puts "All generated surfaces are up to date."
      exit 0
    else
      puts "Stale generated surfaces:"
      modified.each { |m| puts "  #{m[:source]} → #{target_path_for(m[:source], m[:mapping])}" }
      exit 1
    end
  end

  if options[:diff]
    modified = check_modified_sources
    if modified.empty?
      puts "No changes needed."
    else
      puts "Would compile:"
      modified.each { |m| puts "  #{m[:source]} → #{target_path_for(m[:source], m[:mapping])}" }
    end
    exit 0
  end

  # Default: compile all
  puts "Runtime Compiler v#{COMPILER_VERSION}"
  puts "Compiling to SQLite: #{db_path}"
  puts

  # Create schema
  FileUtils.mkdir_p(File.dirname(db_path))
  create_runtime_db_schema(db_path)
  puts "  ✓ Schema created"

  # Build runtime.db from structured YAML
  build_runtime_db(db_path)

  # Compile prose sources to SQLite
  puts
  puts "  Compiling prose sources..."
  plans_entries = []; apk_workflow_entries = []

  @mapping.each do |entry|
    Dir.glob(entry['source']).each do |source_path|
      if entry['compile_rule']&.include?('plans/active/*.md')
        plans_entries << compile_plans_index_to_sqlite(db_path, source_path, entry)
      elsif entry['compile_rule']&.include?('analysis/apk/workflows/*.md')
        apk_workflow_entries << compile_apk_workflow_phases_to_sqlite(db_path, source_path, entry)
      else
        compile_source_to_sqlite(db_path, source_path, entry)
      end
    end
  end

  # Write aggregated plans index
  unless plans_entries.empty?
    insert_gs(db_path, 'plans/active/*.md', 'plans.index', '從 plans/active/*.md 的 front matter、phase 標題、受影響檔案表格提取 plan index', { 'compiled_from' => 'plans/active/*.md', 'total_plans' => plans_entries.length, 'plans' => plans_entries })
    puts "  ✓ plans-index (#{plans_entries.length} plans)"
  end

  # Write aggregated APK workflow phases
  apk_workflow_entries = apk_workflow_entries.compact
  unless apk_workflow_entries.empty?
    insert_gs(db_path, 'analysis/apk/workflows/*.md', 'apk.workflows', '從 analysis/apk/workflows/*.md 的「步驟 N：」標題提取 step 定義', { 'compiled_from' => 'analysis/apk/workflows/*.md', 'total_workflows' => apk_workflow_entries.length, 'workflows' => apk_workflow_entries })
    puts "  ✓ apk-workflows (#{apk_workflow_entries.length} workflows)"
  end

  puts
  puts "Compilation complete. DB: #{db_path}"
end

# CLI entry point
options = {}
OptionParser.new do |opts|
  opts.banner = "Usage: compiler-engine.rb [options]"
  opts.on('--check', 'Check if compilation is needed') { |v| options[:check] = v }
  opts.on('--diff', 'Show what would change') { |v| options[:diff] = v }
  opts.on('--db PATH', 'SQLite output path (default: runtime/runtime.db)') { |v| options[:db] = v }
end.parse!

run(options)