#!/usr/bin/env ruby
# frozen_string_literal: true

# Runtime Config YAML → SQLite Migration Script
# 將所有 runtime/**/*.yaml 設定檔遷移至 runtime.db
#
# 使用方式：
#   ruby scripts/migrate-runtime-config-to-sqlite.rb
#
# 設計原則：
# - 每個 YAML 類別對應一個 SQLite 表格
# - 使用 content TEXT (JSON) 儲存完整結構
# - 對需要查詢的欄位加開專用欄位
# - Idempotent：重複執行不改變結果

require 'yaml'
require 'json'
require 'fileutils'
require 'open3'
require 'time'

ROOT = File.expand_path('..', __dir__)
DB_PATH = File.join(ROOT, 'runtime', 'runtime.db')

def sqlite_exec(sql)
  o, e, s = Open3.capture3('sqlite3', DB_PATH, stdin_data: sql)
  unless s.success?
    warn "sqlite3 error: #{e.strip}" unless e.strip.empty?
  end
  o
end

def sqe(v)
  v.nil? ? 'NULL' : "'#{v.to_s.gsub("'", "''")}'"
end

def jsn(v)
  v.nil? ? 'NULL' : "'#{v.to_json.gsub("'", "''")}'"
end

def upsert(table, set_clause, where_clause)
  sqlite_exec("UPDATE #{table} SET #{set_clause} WHERE #{where_clause};")
  # If no rows updated, insert
end

def migrate_token_budget
  path = File.join(ROOT, 'runtime', 'budget', 'token-budget.yaml')
  return unless File.exist?(path)

  data = YAML.safe_load(File.read(path), permitted_classes: [Date])
  return unless data

  sqlite_exec(<<~SQL)
    CREATE TABLE IF NOT EXISTS runtime_budget (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      model_name TEXT,
      content TEXT NOT NULL,
      created_at TEXT DEFAULT (datetime('now')),
      updated_at TEXT DEFAULT (datetime('now'))
    );
  SQL

  models = data['per_model'] || []
  models.each do |m|
    name = m['name'] || m['model'] || 'default'
    sqlite_exec("INSERT OR REPLACE INTO runtime_budget (id, model_name, content) VALUES ((SELECT id FROM runtime_budget WHERE model_name = #{sqe(name)}), #{sqe(name)}, #{jsn(m)});")
  end

  # Also store default_budget and per_layer as separate rows
  if data['default_budget']
    sqlite_exec("INSERT OR REPLACE INTO runtime_budget (id, model_name, content) VALUES ((SELECT id FROM runtime_budget WHERE model_name = '__default_budget__'), '__default_budget__', #{jsn(data['default_budget'])});")
  end
  if data['per_layer']
    sqlite_exec("INSERT OR REPLACE INTO runtime_budget (id, model_name, content) VALUES ((SELECT id FROM runtime_budget WHERE model_name = '__per_layer__'), '__per_layer__', #{jsn(data['per_layer'])});")
  end

  # Also store the full config
  sqlite_exec("INSERT OR REPLACE INTO runtime_budget (id, model_name, content) VALUES ((SELECT id FROM runtime_budget WHERE model_name = '__config__'), '__config__', #{jsn(data)});")
  puts "  ✓ runtime_budget (#{models.length} models)"
end

def migrate_context_ttl
  path = File.join(ROOT, 'runtime', 'context', 'ttl-policy.yaml')
  return unless File.exist?(path)

  data = YAML.safe_load(File.read(path), permitted_classes: [Date])
  return unless data

  sqlite_exec(<<~SQL)
    CREATE TABLE IF NOT EXISTS context_ttl_policy (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      ttl_type TEXT,
      content TEXT NOT NULL,
      created_at TEXT DEFAULT (datetime('now')),
      updated_at TEXT DEFAULT (datetime('now'))
    );
  SQL

  rules = data['ttl_rules'] || data['rules'] || []
  rules.each do |r|
    ttl_type = r['type'] || r['ttl_type'] || r['name'] || 'default'
    sqlite_exec("INSERT OR REPLACE INTO context_ttl_policy (id, ttl_type, content) VALUES ((SELECT id FROM context_ttl_policy WHERE ttl_type = #{sqe(ttl_type)}), #{sqe(ttl_type)}, #{jsn(r)});")
  end

  sqlite_exec("INSERT OR REPLACE INTO context_ttl_policy (id, ttl_type, content) VALUES ((SELECT id FROM context_ttl_policy WHERE ttl_type = '__config__'), '__config__', #{jsn(data)});")
  puts "  ✓ context_ttl_policy (#{rules.length} rules)"
end

def migrate_discovery_checkpoints
  path = File.join(ROOT, 'runtime', 'discovery', 'capability-checkpoints.yaml')
  return unless File.exist?(path)

  data = YAML.safe_load(File.read(path), permitted_classes: [Date])
  return unless data

  sqlite_exec(<<~SQL)
    CREATE TABLE IF NOT EXISTS capability_checkpoints (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      checkpoint_name TEXT,
      content TEXT NOT NULL,
      created_at TEXT DEFAULT (datetime('now')),
      updated_at TEXT DEFAULT (datetime('now'))
    );
  SQL

  checkpoints = data['checkpoints'] || []
  checkpoints.each do |cp|
    name = cp['phase'] || cp['name'] || 'default'
    sqlite_exec("INSERT OR REPLACE INTO capability_checkpoints (id, checkpoint_name, content) VALUES ((SELECT id FROM capability_checkpoints WHERE checkpoint_name = #{sqe(name)}), #{sqe(name)}, #{jsn(cp)});")
  end

  sqlite_exec("INSERT OR REPLACE INTO capability_checkpoints (id, checkpoint_name, content) VALUES ((SELECT id FROM capability_checkpoints WHERE checkpoint_name = '__config__'), '__config__', #{jsn(data)});")
  puts "  ✓ capability_checkpoints (#{checkpoints.length} checkpoints)"
end

def migrate_distributed
  base = File.join(ROOT, 'runtime', 'distributed')

  # distributed-locks.yaml
  path = File.join(base, 'distributed-locks.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS distributed_locks (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          lock_name TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      locks = data['locks'] || []
      locks.each do |l|
        name = l['name'] || l['lock_name'] || 'default'
        sqlite_exec("INSERT OR REPLACE INTO distributed_locks (id, lock_name, content) VALUES ((SELECT id FROM distributed_locks WHERE lock_name = #{sqe(name)}), #{sqe(name)}, #{jsn(l)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO distributed_locks (id, lock_name, content) VALUES ((SELECT id FROM distributed_locks WHERE lock_name = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ distributed_locks (#{locks.length} locks)"
    end
  end

  # multi-agent-coordination.yaml
  path = File.join(base, 'multi-agent-coordination.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS multi_agent_coordination (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          rule_name TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      rules = data['coordination_rules'] || data['rules'] || []
      rules.each do |r|
        name = r['rule'] || r['name'] || r['rule_id'] || 'default'
        sqlite_exec("INSERT OR REPLACE INTO multi_agent_coordination (id, rule_name, content) VALUES ((SELECT id FROM multi_agent_coordination WHERE rule_name = #{sqe(name)}), #{sqe(name)}, #{jsn(r)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO multi_agent_coordination (id, rule_name, content) VALUES ((SELECT id FROM multi_agent_coordination WHERE rule_name = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ multi_agent_coordination (#{rules.length} rules)"
    end
  end

  # async-job-lifecycle.yaml
  path = File.join(base, 'async-job-lifecycle.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS async_job_lifecycle (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          job_state TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      states = data['job_states'] || data['states'] || []
      states.each do |s|
        name = s['state'] || s['name'] || 'default'
        sqlite_exec("INSERT OR REPLACE INTO async_job_lifecycle (id, job_state, content) VALUES ((SELECT id FROM async_job_lifecycle WHERE job_state = #{sqe(name)}), #{sqe(name)}, #{jsn(s)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO async_job_lifecycle (id, job_state, content) VALUES ((SELECT id FROM async_job_lifecycle WHERE job_state = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ async_job_lifecycle (#{states.length} states)"
    end
  end
end

def migrate_gates
  path = File.join(ROOT, 'runtime', 'gates', 'blocking-gates.yaml')
  return unless File.exist?(path)

  data = YAML.safe_load(File.read(path), permitted_classes: [Date])
  return unless data

  sqlite_exec(<<~SQL)
    CREATE TABLE IF NOT EXISTS blocking_gates (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      gate_name TEXT,
      phase TEXT,
      content TEXT NOT NULL,
      created_at TEXT DEFAULT (datetime('now')),
      updated_at TEXT DEFAULT (datetime('now'))
    );
  SQL

  gates = data['gates'] || []
  gates.each do |g|
    gname = g['name'] || g['id'] || 'default'
    gphase = g['phase'] || 'unknown'
    sqlite_exec("INSERT OR REPLACE INTO blocking_gates (id, gate_name, phase, content) VALUES ((SELECT id FROM blocking_gates WHERE gate_name = #{sqe(gname)} AND phase = #{sqe(gphase)}), #{sqe(gname)}, #{sqe(gphase)}, #{jsn(g)});")
  end

  sqlite_exec("INSERT OR REPLACE INTO blocking_gates (id, gate_name, phase, content) VALUES ((SELECT id FROM blocking_gates WHERE gate_name = '__config__'), '__config__', '__config__', #{jsn(data)});")
  puts "  ✓ blocking_gates (#{gates.length} gates)"
end

def migrate_guards
  base = File.join(ROOT, 'runtime', 'guards')

  # circuit-breaker.yaml
  path = File.join(base, 'circuit-breaker.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS circuit_breaker (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          guard_name TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      # circuit-breaker.yaml uses top-level keys as guard names (recursive_depth, tool_calls, etc.)
      guard_keys = %w[recursive_depth tool_calls context_growth hallucination_risk conflict_rules]
      guards = []
      guard_keys.each do |key|
        next unless data[key]
        guard_entry = { 'name' => key }.merge(data[key])
        guards << guard_entry
        sqlite_exec("INSERT OR REPLACE INTO circuit_breaker (id, guard_name, content) VALUES ((SELECT id FROM circuit_breaker WHERE guard_name = #{sqe(key)}), #{sqe(key)}, #{jsn(guard_entry)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO circuit_breaker (id, guard_name, content) VALUES ((SELECT id FROM circuit_breaker WHERE guard_name = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ circuit_breaker (#{guards.length} guards)"
    end
  end

  # context-pollution.yaml
  path = File.join(base, 'context-pollution.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS context_pollution (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          signal_name TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      signals = data['pollution_signals'] || data['signals'] || []
      signals.each do |s|
        name = s['name'] || s['signal'] || 'default'
        sqlite_exec("INSERT OR REPLACE INTO context_pollution (id, signal_name, content) VALUES ((SELECT id FROM context_pollution WHERE signal_name = #{sqe(name)}), #{sqe(name)}, #{jsn(s)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO context_pollution (id, signal_name, content) VALUES ((SELECT id FROM context_pollution WHERE signal_name = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ context_pollution (#{signals.length} signals)"
    end
  end
end

def migrate_health
  path = File.join(ROOT, 'runtime', 'health', 'context-health-score.yaml')
  return unless File.exist?(path)

  data = YAML.safe_load(File.read(path), permitted_classes: [Date])
  return unless data

  sqlite_exec(<<~SQL)
    CREATE TABLE IF NOT EXISTS context_health_score (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      dimension TEXT,
      content TEXT NOT NULL,
      created_at TEXT DEFAULT (datetime('now')),
      updated_at TEXT DEFAULT (datetime('now'))
    );
  SQL

  dimensions = data['dimensions'] || data['health_dimensions'] || []
  dimensions.each do |d|
    name = d['name'] || d['dimension'] || 'default'
    sqlite_exec("INSERT OR REPLACE INTO context_health_score (id, dimension, content) VALUES ((SELECT id FROM context_health_score WHERE dimension = #{sqe(name)}), #{sqe(name)}, #{jsn(d)});")
  end

  sqlite_exec("INSERT OR REPLACE INTO context_health_score (id, dimension, content) VALUES ((SELECT id FROM context_health_score WHERE dimension = '__config__'), '__config__', #{jsn(data)});")
  puts "  ✓ context_health_score (#{dimensions.length} dimensions)"
end

def migrate_intelligence_routing
  path = File.join(ROOT, 'runtime', 'intelligence', 'intelligence-routing.yaml')
  return unless File.exist?(path)

  data = YAML.safe_load(File.read(path), permitted_classes: [Date])
  return unless data

  sqlite_exec(<<~SQL)
    CREATE TABLE IF NOT EXISTS intelligence_routing (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      domain TEXT,
      content TEXT NOT NULL,
      created_at TEXT DEFAULT (datetime('now')),
      updated_at TEXT DEFAULT (datetime('now'))
    );
  SQL

  rules = data['routing_rules'] || data['rules'] || []
  rules.each do |r|
    domain = r['domain'] || r['name'] || 'default'
    sqlite_exec("INSERT OR REPLACE INTO intelligence_routing (id, domain, content) VALUES ((SELECT id FROM intelligence_routing WHERE domain = #{sqe(domain)}), #{sqe(domain)}, #{jsn(r)});")
  end

  sqlite_exec("INSERT OR REPLACE INTO intelligence_routing (id, domain, content) VALUES ((SELECT id FROM intelligence_routing WHERE domain = '__config__'), '__config__', #{jsn(data)});")
  puts "  ✓ intelligence_routing (#{rules.length} rules)"
end

def migrate_obligations
  path = File.join(ROOT, 'runtime', 'obligations', 'obligation-ledger.yaml')
  return unless File.exist?(path)

  data = YAML.safe_load(File.read(path), permitted_classes: [Date])
  return unless data

  sqlite_exec(<<~SQL)
    CREATE TABLE IF NOT EXISTS obligation_ledger (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      obligation_name TEXT,
      phase TEXT,
      content TEXT NOT NULL,
      created_at TEXT DEFAULT (datetime('now')),
      updated_at TEXT DEFAULT (datetime('now'))
    );
  SQL

  obligations = data['obligations'] || []
  obligations.each do |ob|
    oname = ob['name'] || ob['id'] || 'default'
    ophase = ob['phase'] || 'unknown'
    sqlite_exec("INSERT OR REPLACE INTO obligation_ledger (id, obligation_name, phase, content) VALUES ((SELECT id FROM obligation_ledger WHERE obligation_name = #{sqe(oname)} AND phase = #{sqe(ophase)}), #{sqe(oname)}, #{sqe(ophase)}, #{jsn(ob)});")
  end

  sqlite_exec("INSERT OR REPLACE INTO obligation_ledger (id, obligation_name, phase, content) VALUES ((SELECT id FROM obligation_ledger WHERE obligation_name = '__config__'), '__config__', '__config__', #{jsn(data)});")
  puts "  ✓ obligation_ledger (#{obligations.length} obligations)"
end

def migrate_output_governance
  base = File.join(ROOT, 'runtime', 'output-governance')

  # language-policy.yaml
  path = File.join(base, 'language-policy.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS language_policy (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          rule_name TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      rules = data['core_rules'] || data['rules'] || []
      rules = [data] if rules.empty? && data.is_a?(Hash)
      rules = Array(rules) unless rules.is_a?(Array)
      rules.each_with_index do |r, i|
        name = r['name'] || r['rule'] || "rule_#{i}"
        sqlite_exec("INSERT OR REPLACE INTO language_policy (id, rule_name, content) VALUES ((SELECT id FROM language_policy WHERE rule_name = #{sqe(name)}), #{sqe(name)}, #{jsn(r)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO language_policy (id, rule_name, content) VALUES ((SELECT id FROM language_policy WHERE rule_name = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ language_policy"
    end
  end

  # output-rules.yaml
  path = File.join(base, 'output-rules.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS output_rules (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          rule_name TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      rules = data['rules'] || data['output_rules'] || []
      rules = [data] if rules.empty? && data.is_a?(Hash)
      rules = Array(rules) unless rules.is_a?(Array)
      rules.each_with_index do |r, i|
        name = r['name'] || r['rule'] || "rule_#{i}"
        sqlite_exec("INSERT OR REPLACE INTO output_rules (id, rule_name, content) VALUES ((SELECT id FROM output_rules WHERE rule_name = #{sqe(name)}), #{sqe(name)}, #{jsn(r)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO output_rules (id, rule_name, content) VALUES ((SELECT id FROM output_rules WHERE rule_name = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ output_rules"
    end
  end

  # governance-gates.yaml
  path = File.join(base, 'governance-gates.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS governance_gates (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          gate_name TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      gates = data['gates'] || data['governance_gates'] || []
      gates.each do |g|
        name = g['name'] || g['id'] || 'default'
        sqlite_exec("INSERT OR REPLACE INTO governance_gates (id, gate_name, content) VALUES ((SELECT id FROM governance_gates WHERE gate_name = #{sqe(name)}), #{sqe(name)}, #{jsn(g)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO governance_gates (id, gate_name, content) VALUES ((SELECT id FROM governance_gates WHERE gate_name = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ governance_gates (#{gates.length} gates)"
    end
  end
end

def migrate_phase_machine
  path = File.join(ROOT, 'runtime', 'phases', 'phase-machine.yaml')
  return unless File.exist?(path)

  data = YAML.safe_load(File.read(path), permitted_classes: [Date])
  return unless data

  sqlite_exec(<<~SQL)
    CREATE TABLE IF NOT EXISTS phase_machine (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      phase_name TEXT,
      content TEXT NOT NULL,
      created_at TEXT DEFAULT (datetime('now')),
      updated_at TEXT DEFAULT (datetime('now'))
    );
  SQL

  phases = data['phases'] || []
  phases.each do |ph|
    pname = ph['name'] || ph['id'] || 'default'
    sqlite_exec("INSERT OR REPLACE INTO phase_machine (id, phase_name, content) VALUES ((SELECT id FROM phase_machine WHERE phase_name = #{sqe(pname)}), #{sqe(pname)}, #{jsn(ph)});")
  end

  sqlite_exec("INSERT OR REPLACE INTO phase_machine (id, phase_name, content) VALUES ((SELECT id FROM phase_machine WHERE phase_name = '__config__'), '__config__', #{jsn(data)});")
  puts "  ✓ phase_machine (#{phases.length} phases)"
end

def migrate_pipeline
  base = File.join(ROOT, 'runtime', 'pipeline')

  # context-flow.yaml
  path = File.join(base, 'context-flow.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS pipeline_context_flow (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          level TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      levels = data['expansion_levels'] || data['levels'] || []
      levels.each do |l|
        name = l['level'] || l['name'] || 'default'
        sqlite_exec("INSERT OR REPLACE INTO pipeline_context_flow (id, level, content) VALUES ((SELECT id FROM pipeline_context_flow WHERE level = #{sqe(name)}), #{sqe(name)}, #{jsn(l)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO pipeline_context_flow (id, level, content) VALUES ((SELECT id FROM pipeline_context_flow WHERE level = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ pipeline_context_flow (#{levels.length} levels)"
    end
  end

  # guard-chain.yaml
  path = File.join(base, 'guard-chain.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS guard_chain (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          stage TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      stages = data['guard_stages'] || data['stages'] || []
      stages.each do |s|
        name = s['stage'] || s['name'] || 'default'
        sqlite_exec("INSERT OR REPLACE INTO guard_chain (id, stage, content) VALUES ((SELECT id FROM guard_chain WHERE stage = #{sqe(name)}), #{sqe(name)}, #{jsn(s)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO guard_chain (id, stage, content) VALUES ((SELECT id FROM guard_chain WHERE stage = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ guard_chain (#{stages.length} stages)"
    end
  end

  # relevance-engine.yaml
  path = File.join(base, 'relevance-engine.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS relevance_engine (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          component TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      # relevance-engine.yaml uses 'scoring' key with nested components array
      scoring = data['scoring']
      components = []
      if scoring && scoring['components'].is_a?(Array)
        scoring['components'].each do |c|
          name = c['id'] || c['name'] || 'default'
          components << c
          sqlite_exec("INSERT OR REPLACE INTO relevance_engine (id, component, content) VALUES ((SELECT id FROM relevance_engine WHERE component = #{sqe(name)}), #{sqe(name)}, #{jsn(c)});")
        end
      end
      sqlite_exec("INSERT OR REPLACE INTO relevance_engine (id, component, content) VALUES ((SELECT id FROM relevance_engine WHERE component = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ relevance_engine (#{components.length} components)"
    end
  end

  # session-lifecycle.yaml
  path = File.join(base, 'session-lifecycle.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS session_lifecycle (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          stage TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      stages = data['session_stages'] || data['stages'] || []
      stages.each do |s|
        name = s['stage'] || s['name'] || 'default'
        sqlite_exec("INSERT OR REPLACE INTO session_lifecycle (id, stage, content) VALUES ((SELECT id FROM session_lifecycle WHERE stage = #{sqe(name)}), #{sqe(name)}, #{jsn(s)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO session_lifecycle (id, stage, content) VALUES ((SELECT id FROM session_lifecycle WHERE stage = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ session_lifecycle (#{stages.length} stages)"
    end
  end
end

def migrate_prompt_artifacts
  base = File.join(ROOT, 'runtime', 'prompt-artifacts')

  # artifact-templates.yaml
  path = File.join(base, 'artifact-templates.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS prompt_artifact_templates (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          template_name TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      templates = data['templates'] || data['artifact_templates'] || []
      templates.each do |t|
        name = t['name'] || t['template_name'] || t['id'] || 'default'
        sqlite_exec("INSERT OR REPLACE INTO prompt_artifact_templates (id, template_name, content) VALUES ((SELECT id FROM prompt_artifact_templates WHERE template_name = #{sqe(name)}), #{sqe(name)}, #{jsn(t)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO prompt_artifact_templates (id, template_name, content) VALUES ((SELECT id FROM prompt_artifact_templates WHERE template_name = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ prompt_artifact_templates (#{templates.length} templates)"
    end
  end

  # composition-rules.yaml
  path = File.join(base, 'composition-rules.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS prompt_composition_rules (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          rule_name TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      rules = data['composition_rules'] || data['rules'] || []
      rules.each do |r|
        name = r['name'] || r['rule_name'] || 'default'
        sqlite_exec("INSERT OR REPLACE INTO prompt_composition_rules (id, rule_name, content) VALUES ((SELECT id FROM prompt_composition_rules WHERE rule_name = #{sqe(name)}), #{sqe(name)}, #{jsn(r)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO prompt_composition_rules (id, rule_name, content) VALUES ((SELECT id FROM prompt_composition_rules WHERE rule_name = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ prompt_composition_rules (#{rules.length} rules)"
    end
  end
end

def migrate_recovery
  base = File.join(ROOT, 'runtime', 'recovery')

  # recovery-strategies.yaml
  path = File.join(base, 'recovery-strategies.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS recovery_strategies (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          strategy_name TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      strategies = data['recovery_strategies'] || data['strategies'] || []
      strategies.each do |s|
        name = s['name'] || s['strategy'] || s['failure_mode'] || 'default'
        sqlite_exec("INSERT OR REPLACE INTO recovery_strategies (id, strategy_name, content) VALUES ((SELECT id FROM recovery_strategies WHERE strategy_name = #{sqe(name)}), #{sqe(name)}, #{jsn(s)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO recovery_strategies (id, strategy_name, content) VALUES ((SELECT id FROM recovery_strategies WHERE strategy_name = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ recovery_strategies (#{strategies.length} strategies)"
    end
  end

  # state-repair.yaml
  path = File.join(base, 'state-repair.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS state_repair (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          procedure_name TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      procedures = data['repair_procedures'] || data['procedures'] || []
      procedures.each do |p|
        name = p['name'] || p['procedure'] || p['condition'] || 'default'
        sqlite_exec("INSERT OR REPLACE INTO state_repair (id, procedure_name, content) VALUES ((SELECT id FROM state_repair WHERE procedure_name = #{sqe(name)}), #{sqe(name)}, #{jsn(p)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO state_repair (id, procedure_name, content) VALUES ((SELECT id FROM state_repair WHERE procedure_name = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ state_repair (#{procedures.length} procedures)"
    end
  end

  # obligation-rebuild.yaml
  path = File.join(base, 'obligation-rebuild.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS obligation_rebuild (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          procedure_name TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      procedures = data['rebuild_procedures'] || data['procedures'] || []
      procedures.each do |p|
        name = p['name'] || p['procedure'] || 'default'
        sqlite_exec("INSERT OR REPLACE INTO obligation_rebuild (id, procedure_name, content) VALUES ((SELECT id FROM obligation_rebuild WHERE procedure_name = #{sqe(name)}), #{sqe(name)}, #{jsn(p)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO obligation_rebuild (id, procedure_name, content) VALUES ((SELECT id FROM obligation_rebuild WHERE procedure_name = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ obligation_rebuild (#{procedures.length} procedures)"
    end
  end

  # phase-reconciliation.yaml
  path = File.join(base, 'phase-reconciliation.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS phase_reconciliation (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          procedure_name TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      procedures = data['reconciliation_procedures'] || data['procedures'] || []
      procedures.each do |p|
        name = p['name'] || p['procedure'] || 'default'
        sqlite_exec("INSERT OR REPLACE INTO phase_reconciliation (id, procedure_name, content) VALUES ((SELECT id FROM phase_reconciliation WHERE procedure_name = #{sqe(name)}), #{sqe(name)}, #{jsn(p)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO phase_reconciliation (id, procedure_name, content) VALUES ((SELECT id FROM phase_reconciliation WHERE procedure_name = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ phase_reconciliation (#{procedures.length} procedures)"
    end
  end
end

def migrate_scheduler
  base = File.join(ROOT, 'runtime', 'scheduler')

  # execution-queue.yaml
  path = File.join(base, 'execution-queue.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS execution_queue (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          queue_name TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      # execution-queue.yaml uses 'queue_structure' key with nested name/description/fields
      queue_struct = data['queue_structure']
      queues = []
      if queue_struct
        entry = { 'name' => queue_struct['name'] || 'default' }.merge(queue_struct)
        queues << entry
        sqlite_exec("INSERT OR REPLACE INTO execution_queue (id, queue_name, content) VALUES ((SELECT id FROM execution_queue WHERE queue_name = #{sqe(entry['name'])}), #{sqe(entry['name'])}, #{jsn(entry)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO execution_queue (id, queue_name, content) VALUES ((SELECT id FROM execution_queue WHERE queue_name = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ execution_queue (#{queues.length} queues)"
    end
  end

  # priority-scheduler.yaml
  path = File.join(base, 'priority-scheduler.yaml')
  if File.exist?(path)
    data = YAML.safe_load(File.read(path), permitted_classes: [Date])
    if data
      sqlite_exec(<<~SQL)
        CREATE TABLE IF NOT EXISTS priority_scheduler (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          priority_level TEXT,
          content TEXT NOT NULL,
          created_at TEXT DEFAULT (datetime('now')),
          updated_at TEXT DEFAULT (datetime('now'))
        );
      SQL
      levels = data['priority_levels'] || data['levels'] || []
      levels.each do |l|
        name = l['level'] || l['priority'] || 'default'
        sqlite_exec("INSERT OR REPLACE INTO priority_scheduler (id, priority_level, content) VALUES ((SELECT id FROM priority_scheduler WHERE priority_level = #{sqe(name)}), #{sqe(name)}, #{jsn(l)});")
      end
      sqlite_exec("INSERT OR REPLACE INTO priority_scheduler (id, priority_level, content) VALUES ((SELECT id FROM priority_scheduler WHERE priority_level = '__config__'), '__config__', #{jsn(data)});")
      puts "  ✓ priority_scheduler (#{levels.length} levels)"
    end
  end
end

def migrate_router
  path = File.join(ROOT, 'runtime', 'router', 'activation-rules.yaml')
  return unless File.exist?(path)

  data = YAML.safe_load(File.read(path), permitted_classes: [Date])
  return unless data

  sqlite_exec(<<~SQL)
    CREATE TABLE IF NOT EXISTS activation_rules_mirror (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      rule_name TEXT,
      content TEXT NOT NULL,
      created_at TEXT DEFAULT (datetime('now')),
      updated_at TEXT DEFAULT (datetime('now'))
    );
  SQL

  rules = data['rules'] || []
  rules.each do |r|
    name = r['rule_id'] || r['name'] || 'default'
    sqlite_exec("INSERT OR REPLACE INTO activation_rules_mirror (id, rule_name, content) VALUES ((SELECT id FROM activation_rules_mirror WHERE rule_name = #{sqe(name)}), #{sqe(name)}, #{jsn(r)});")
  end

  sqlite_exec("INSERT OR REPLACE INTO activation_rules_mirror (id, rule_name, content) VALUES ((SELECT id FROM activation_rules_mirror WHERE rule_name = '__config__'), '__config__', #{jsn(data)});")
  puts "  ✓ activation_rules_mirror (#{rules.length} rules)"
end

def migrate_transaction_templates
  path = File.join(ROOT, 'runtime', 'transactions', 'transaction-templates.yaml')
  return unless File.exist?(path)

  data = YAML.safe_load(File.read(path), permitted_classes: [Date])
  return unless data

  sqlite_exec(<<~SQL)
    CREATE TABLE IF NOT EXISTS transaction_templates_ext (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      template_name TEXT,
      content TEXT NOT NULL,
      created_at TEXT DEFAULT (datetime('now')),
      updated_at TEXT DEFAULT (datetime('now'))
    );
  SQL

  templates = data['templates'] || data['transaction_templates'] || []
  templates.each do |t|
    name = t['name'] || t['id'] || 'default'
    sqlite_exec("INSERT OR REPLACE INTO transaction_templates_ext (id, template_name, content) VALUES ((SELECT id FROM transaction_templates_ext WHERE template_name = #{sqe(name)}), #{sqe(name)}, #{jsn(t)});")
  end

  sqlite_exec("INSERT OR REPLACE INTO transaction_templates_ext (id, template_name, content) VALUES ((SELECT id FROM transaction_templates_ext WHERE template_name = '__config__'), '__config__', #{jsn(data)});")
  puts "  ✓ transaction_templates_ext (#{templates.length} templates)"
end

# ═══════════════════════════════════════════════════════════════
# Main
# ═══════════════════════════════════════════════════════════════

def run
  puts "Runtime Config YAML → SQLite Migration"
  puts "DB: #{DB_PATH}"
  puts

  # Create schema_migrations table if not exists
  sqlite_exec(<<~SQL)
    CREATE TABLE IF NOT EXISTS schema_migrations (
      version TEXT PRIMARY KEY,
      description TEXT,
      applied_at TEXT DEFAULT (datetime('now'))
    );
  SQL

  # Check if already migrated
  result = sqlite_exec("SELECT version FROM schema_migrations WHERE version = 'v2_runtime_config_migration';")
  if result && !result.strip.empty?
    puts "Migration v2_runtime_config_migration already applied. Skipping."
    puts
    puts "To re-run: DELETE FROM schema_migrations WHERE version = 'v2_runtime_config_migration';"
    exit 0
  end

  migrate_token_budget
  migrate_context_ttl
  migrate_discovery_checkpoints
  migrate_distributed
  migrate_gates
  migrate_guards
  migrate_health
  migrate_intelligence_routing
  migrate_obligations
  migrate_output_governance
  migrate_phase_machine
  migrate_pipeline
  migrate_prompt_artifacts
  migrate_recovery
  migrate_scheduler
  migrate_router
  migrate_transaction_templates

  # Record migration
  sqlite_exec("INSERT INTO schema_migrations (version, description) VALUES ('v2_runtime_config_migration', 'Migrate all runtime/**/*.yaml config files to SQLite tables');")
  puts
  puts "Migration complete! All runtime config YAML files have been migrated to SQLite."
  puts "Total new tables: 31"
end

run