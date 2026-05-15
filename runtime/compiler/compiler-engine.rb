#!/usr/bin/env ruby
# frozen_string_literal: true

# Runtime Compiler Engine
# 將 canonical prose source 編譯為 runtime/generated/*.yaml。
#
# 使用方式：
#   ruby runtime/compiler/compiler-engine.rb          # 編譯所有 source
#   ruby runtime/compiler/compiler-engine.rb --check  # 只檢查是否需要編譯
#   ruby runtime/compiler/compiler-engine.rb --diff   # 顯示預期變更
#
# 設計原則：
# - Deterministic：相同輸入 → 相同輸出
# - Idempotent：重複執行不改變結果
# - 只編譯 deterministic state，不編譯 judgment/heuristics

require 'yaml'
require 'fileutils'
require 'time'
require 'optparse'

COMPILER_VERSION = '1.0.0'
GENERATED_DIR = File.join(File.dirname(__FILE__), '..', 'generated')
COMPILER_RULES = File.join(File.dirname(__FILE__), 'compiler-rules.yaml')

# Source-target mapping loaded from compiler-rules.yaml
@mapping = nil

def load_mapping
  rules = YAML.safe_load(File.read(COMPILER_RULES), permitted_classes: [Date])
  @mapping = rules['source_target_mapping']
end

def generated_header(source_path)
  {
    'generated_from' => source_path,
    'generated_at' => Time.now.utc.strftime('%Y-%m-%dT%H:%M:%SZ'),
    'compiler_version' => COMPILER_VERSION,
    'status' => 'synced'
  }
end

def target_path_for(source_path, mapping_entry)
  domain = extract_domain(source_path)
  target = mapping_entry['target']
    .gsub('{domain}', domain)
  # If target contains a directory path (e.g., analysis/apk/workflows/generated-phases.yaml),
  # use it directly. Otherwise, join with GENERATED_DIR.
  if target.include?('/')
    target
  else
    File.join(GENERATED_DIR, File.basename(target))
  end
end

def extract_domain(source_path)
  # Extract domain from path like workflow/apk-analysis/execution-flow.md → apk-analysis
  match = source_path.match(%r{workflow/([^/]+)/})
  return match[1] if match

  # enforcement/dependency-reading.md → transaction-machine
  match = source_path.match(%r{enforcement/(.+)\.md$})
  return match[1].tr('-', '_') if match

  # plans/active/*.md → plans-index (single target for all plans)
  match = source_path.match(%r{plans/})
  return 'plans-index' if match

  # analysis/apk/workflows/*.md → apk-workflow (single target for all APK workflows)
  match = source_path.match(%r{analysis/apk/workflows/})
  return 'apk-workflow' if match

  'unknown'
end

def extract_phase_sections(content)
  # Split content into phase sections by ## N. heading
  # Returns array of { heading:, line_number:, body: }
  sections = []
  content.scan(/^(##\s+\d+\.\s+.+)$/) do
    heading_text = $1.strip
    heading_line = $`.lines.count + 1

    # Find the start of this section
    section_start = $`.size
    # Find the next ## heading or end of content
    remaining = content[section_start + $&.size..]
    next_section_match = remaining.match(/^##\s+\d+\.\s+/)
    section_body = if next_section_match
                     remaining[0...next_section_match.begin(0)]
                   else
                     remaining
                   end

    sections << {
      'heading' => heading_text.sub(/^##\s+/, '').strip,
      'line_number' => heading_line,
      'body' => section_body.strip
    }
  end
  sections
end

def extract_allowed_actions_from_section(body)
  # Extract allowed actions from:
  # 1. Markdown tables with "行動" or "Action" column
  # 2. Bullet lists under "記錄" or "必要行動" or similar
  # 3. Numbered step lists
  actions = []

  # Pattern 1: Table rows with action descriptions (| N | action |)
  body.scan(/^\|\s*\d+\s*\|\s*(.+?)\s*\|$/) do |match|
    actions << { 'action' => match[0].strip, 'source' => 'table' }
  end

  # Pattern 2: Table rows with "問題 | 必要行動" format (second column is action)
  # Only capture rows where the second column is a meaningful action description
  body.scan(/^\|\s*(?:[^|]+)\s*\|\s*(.+?)\s*\|$/) do |match|
    col2 = match[0].strip
    next if col2 == '---' || col2 == '必要行動' || col2 == 'Action' || col2 == '行動'
    next if col2 == '緩解措施' || col2 == '回填要求' || col2 == '證明'
    next if col2.start_with?('`<')
    next if col2 =~ /^(步驟|Step|Reset|測試類型|文件|問題|根本原因|用途)/
    next if col2.length < 10  # skip short/non-action content
    next if col2 =~ /^\|/  # skip nested table artifacts
    actions << { 'action' => col2, 'source' => 'table' }
  end

  # Pattern 3: Bullet list items under "記錄" or action-oriented paragraphs
  body.scan(/^-\s+(.+)$/) do |match|
    item = match[0].strip
    next if item.start_with?('`<')
    next if item =~ /^(不要|禁止|Don't|Do not|Never|避免)/
    actions << { 'action' => item, 'source' => 'bullet' }
  end

  # Pattern 4: Numbered step items (1. **Action**)
  body.scan(/^\d+\.\s+\*\*([^*]+)\*\*[：:]\s*(.+)$/) do |match|
    actions << { 'action' => "#{match[0].strip}：#{match[1].strip}", 'source' => 'step' }
  end

  actions.uniq { |a| a['action'] }
end

def extract_blocking_gates_from_section(body)
  # Extract blocking gates from:
  # 1. "禁止" or "Don't" bullet lists
  # 2. "阻擋項" or "blocking" mentions
  # 3. Validation/check conditions
  gates = []

  # Pattern 1: "禁止" / "Don't" / "不要" bullet lists
  body.scan(/^-\s+(不要|禁止|Don't|Do not|Never|避免)\s*(.+)$/i) do |prefix, rest|
    gates << {
      'id' => "gate.#{rest.strip.downcase.gsub(/[^a-z0-9]+/, '_').gsub(/^_|_$/, '')}",
      'description' => "#{prefix}#{rest}",
      'severity' => 'critical',
      'source' => 'prohibition'
    }
  end

  # Pattern 2: "阻擋項" or "blocking" mentions
  body.scan(/(.+?)(?:是|為|屬於)\s*阻擋項/) do |match|
    gates << {
      'id' => "gate.blocking.#{match[0].strip.downcase.gsub(/[^a-z0-9]+/, '_').gsub(/^_|_$/, '')}",
      'description' => "#{match[0].strip}是阻擋項",
      'severity' => 'critical',
      'source' => 'blocking_condition'
    }
  end

  # Pattern 3: "必須" / "must" conditions that act as gates
  body.scan(/(?:在|於|在).+?(?:之前|前).+?(?:必須|需要|應).+?。/) do |match|
    gates << {
      'id' => "gate.prerequisite.#{match[0..30].strip.downcase.gsub(/[^a-z0-9]+/, '_').gsub(/^_|_$/, '')}",
      'description' => match[0].strip,
      'severity' => 'high',
      'source' => 'prerequisite'
    }
  end

  gates.uniq { |g| g['id'] }
end

def extract_tables_from_section(body)
  # Extract markdown tables as structured data
  tables = []
  lines = body.split("\n")
  current_table = []
  in_table = false

  lines.each do |line|
    if line.match?(/^\|.+\|$/)
      current_table << line.strip
      in_table = true
    else
      if in_table && current_table.length >= 3 # header + separator + at least 1 row
        # Parse the table
        header = current_table[0].split('|').map(&:strip).reject(&:empty?)
        rows = current_table[2..].map { |r| r.split('|').map(&:strip).reject(&:empty?) }
        tables << { 'header' => header, 'rows' => rows }
      end
      current_table = []
      in_table = false
    end
  end

  # Don't forget last table
  if in_table && current_table.length >= 3
    header = current_table[0].split('|').map(&:strip).reject(&:empty?)
    rows = current_table[2..].map { |r| r.split('|').map(&:strip).reject(&:empty?) }
    tables << { 'header' => header, 'rows' => rows }
  end

  tables
end

def compile_workflow_phases(source_path, mapping_entry)
  content = File.read(source_path)
  sections = extract_phase_sections(content)

  phases = sections.map do |sec|
    allowed = extract_allowed_actions_from_section(sec['body'])
    gates = extract_blocking_gates_from_section(sec['body'])
    tables = extract_tables_from_section(sec['body'])

    phase_entry = {
      'name' => sec['heading'],
      'source_line' => sec['line_number']
    }
    phase_entry['allowed_actions'] = allowed unless allowed.empty?
    phase_entry['blocking_gates'] = gates unless gates.empty?
    phase_entry['tables'] = tables unless tables.empty?
    phase_entry
  end

  # Also extract global gates (not tied to a specific phase)
  all_gates = phases.flat_map { |p| p['blocking_gates'] || [] }.uniq { |g| g['id'] }

  target = target_path_for(source_path, mapping_entry)
  header = generated_header(source_path)

  yaml_content = {
    'header' => header,
    'phases' => phases
  }
  yaml_content['gates'] = all_gates unless all_gates.empty?

  FileUtils.mkdir_p(File.dirname(target))
  File.write(target, YAML.dump(yaml_content))
  puts "  ✓ #{target}"
end

def extract_apk_step_sections(content)
  # Split content into step sections by ## 步驟 N： heading
  # Returns array of { heading:, line_number:, body: }
  sections = []
  content.scan(/^(##\s+步驟\s+\d+[：:]\s*.+)$/) do
    heading_text = $1.strip
    heading_line = $`.lines.count + 1

    # Find the start of this section
    section_start = $`.size
    # Find the next ## 步驟 heading or end of content
    remaining = content[section_start + $&.size..]
    next_section_match = remaining.match(/^##\s+步驟\s+\d+[：:]/)
    section_body = if next_section_match
                     remaining[0...next_section_match.begin(0)]
                   else
                     remaining
                   end

    sections << {
      'heading' => heading_text.sub(/^##\s+/, '').strip,
      'line_number' => heading_line,
      'body' => section_body.strip
    }
  end
  sections
end

def extract_prerequisites(content)
  # Extract prerequisites from ## 前置準備 section
  prereq_match = content.match(/^##\s+前置準備\n(.+?)(?=\n##\s+|\z)/m)
  return [] unless prereq_match

  prereq_body = prereq_match[1].strip
  items = []

  # Extract bullet list items
  prereq_body.scan(/^-\s+(.+)$/) do |match|
    items << match[0].strip
  end

  # Extract ### sub-headings as prerequisite categories
  prereq_body.scan(/^###\s+(.+)$/) do |match|
    items << match[0].strip
  end

  items
end

def extract_output_format(content)
  # Extract output format from ## 成功產出格式 section
  output_match = content.match(/^##\s+成功產出格式\n(.+?)(?=\n##\s+|\z)/m)
  return nil unless output_match

  output_match[1].strip
end

def extract_prerequisites_clean(content)
  # Extract prerequisites from ## 前置準備 section, excluding ### sub-headings
  prereq_match = content.match(/^##\s+前置準備\n(.+?)(?=\n##\s+|\z)/m)
  return [] unless prereq_match

  prereq_body = prereq_match[1].strip
  items = []

  # Extract bullet list items only (skip ### sub-headings like "必要條件", "工具")
  prereq_body.scan(/^-\s+(.+)$/) do |match|
    items << match[0].strip
  end

  items
end

def compile_apk_workflow_phases(source_path, _mapping_entry)
  # Skip README — not a workflow
  filename = File.basename(source_path, '.md')
  if filename == 'README'
    puts "  - #{source_path} (skipped — README)"
    return nil
  end

  content = File.read(source_path)

  # Extract step sections using ## 步驟 N： heading format
  sections = extract_apk_step_sections(content)

  steps = sections.map do |sec|
    {
      'title' => sec['heading'],
      'source_line' => sec['line_number']
    }
  end

  # Extract prerequisites (bullet items only, skip ### sub-headings)
  prerequisites = extract_prerequisites_clean(content)

  # Extract output format
  output_format = extract_output_format(content)

  # Return entry for aggregation (run method handles writing)
  entry = {
    'workflow_name' => filename,
    'source_path' => source_path,
    'total_steps' => steps.length,
    'steps' => steps
  }
  entry['prerequisites'] = prerequisites unless prerequisites.empty?
  entry['output_format'] = output_format unless output_format.nil?

  puts "  ✓ #{source_path} (#{steps.length} steps)"
  entry
end

def compile_enforcement_transactions(source_path, mapping_entry)
  content = File.read(source_path)

  # Extract transaction states from markdown sections
  states = []
  content.scan(/^###?\s+(.+?)(?:\s*\((.+?)\))?$/) do |match|
    states << { 'name' => match[0].strip, 'context' => match[1]&.strip }
  end

  target = target_path_for(source_path, mapping_entry)
  header = generated_header(source_path)

  yaml_content = {
    'header' => header,
    'compiled_from' => source_path,
    'states' => states
  }

  FileUtils.mkdir_p(File.dirname(target))
  File.write(target, YAML.dump(yaml_content))
  puts "  ✓ #{target}"
end

def compile_output_governance(source_path, mapping_entry)
  content = File.read(source_path)

  # Extract governance rules from markdown content
  rules = []
  content.scan(/^###?\s+(.+?)$/) do |match|
    rules << { 'section' => match[0].strip }
  end

  # Extract validation gates
  gates = []
  content.scan(/^\*\*([^*]+)\*\*：(.+)$/) do |match|
    gates << { 'name' => match[0].strip, 'description' => match[1].strip }
  end

  target = target_path_for(source_path, mapping_entry)
  header = generated_header(source_path)

  yaml_content = {
    'header' => header,
    'compiled_from' => source_path,
    'rules' => rules,
    'gates' => gates
  }

  FileUtils.mkdir_p(File.dirname(target))
  File.write(target, YAML.dump(yaml_content))
  puts "  ✓ #{target}"
end

def compile_knowledge_update_flow(source_path, _mapping_entry)
  content = File.read(source_path)

  # Extract 11 steps from ## level headings with step numbers
  steps = []
  content.scan(/^##\s+Step\s+(\d+)[：:]\s*(.+?)$/) do |match|
    step_num = match[0].strip.to_i
    step_name = match[1].strip

    # Find the step's content block (from this ## to next ## or end)
    step_start = content.index("## Step #{step_num}：")
    next unless step_start

    remaining = content[step_start..]
    next_step_match = remaining.index(/^## Step #{step_num + 1}[：:]/)
    step_content = if next_step_match
                     remaining[0...next_step_match]
                   else
                     remaining
                   end

    # Extract entry conditions (判斷結果 table)
    entry_conditions = []
    step_content.scan(/^\|\s*(\w[\w\s]+?)\s*\|\s*(.+?)\s*\|$/) do |row|
      entry_conditions << { 'condition' => row[0].strip, 'next_step' => row[1].strip }
    end

    # Extract reference sources
    references = []
    step_content.scan(/\[`([^`]+)`\]\(([^)]+)\)/) do |ref|
      references << { 'name' => ref[0], 'path' => ref[1] }
    end

    steps << {
      'step' => step_num,
      'name' => step_name,
      'entry_conditions' => entry_conditions,
      'references' => references
    }
  end

  target = File.join(GENERATED_DIR, 'knowledge-update-phases.yaml')
  header = generated_header(source_path)

  yaml_content = {
    'header' => header,
    'compiled_from' => source_path,
    'total_steps' => 11,
    'steps' => steps
  }

  FileUtils.mkdir_p(GENERATED_DIR)
  File.write(target, YAML.dump(yaml_content))
  puts "  ✓ #{target}"
end

def compile_workflow_artifacts(source_path, mapping_entry)
  content = File.read(source_path)

  # Extract artifact sections from ## level headings
  artifacts = []
  content.scan(/^##\s+\d+\.\s+(.+)$/) do |match|
    artifacts << { 'name' => match[0].strip }
  end

  # Extract verification gates (completion gates, quality gates)
  gates = []
  content.scan(/^###\s+(.+?)(?:\s*Gate|gate)?$/) do |match|
    gates << { 'name' => match[0].strip, 'type' => 'verification_gate' }
  end

  # Extract table-based artifact definitions (markdown tables with | Artifact | ... |)
  tables = []
  content.scan(/^\|.+\|.+\|$/) do |line|
    next if line.match?(/^\|[\s-]+\|[\s-]+\|$/) # skip separator rows
    next if line.match?(/^\|.*#.*\|$/) # skip non-artifact tables
    tables << line.strip
  end

  # Extract required items from numbered lists under artifact sections
  required_items = []
  content.scan(/^\d+\.\s+\*\*([^*]+)\*\*(.*)$/) do |match|
    required_items << { 'name' => match[0].strip, 'description' => match[1]&.strip }
  end

  target = target_path_for(source_path, mapping_entry)
  header = generated_header(source_path)

  yaml_content = {
    'header' => header,
    'compiled_from' => source_path,
    'artifacts' => artifacts,
    'verification_gates' => gates,
    'required_items' => required_items
  }

  FileUtils.mkdir_p(File.dirname(target))
  File.write(target, YAML.dump(yaml_content))
  puts "  ✓ #{target}"
end

def compile_goal_action_gates(source_path, mapping_entry)
  content = File.read(source_path)

  # Extract the core goal/execution/validation table (Chinese headers: 欄位 | 必填內容)
  core_fields = []
  # Match table rows with Chinese or English content in first column
  content.scan(/^\|\s*(.+?)\s*\|\s*(.+?)\s*\|$/) do |match|
    col1 = match[0].strip
    col2 = match[1].strip
    next if col1 == '---' || col1 == '欄位' || col1 == '必填內容'
    next if col1 == '情境' || col1 == '要求'
    next if col1 == '工作類型' || col1 == '驗證方式'
    next if col1 == '工作單元' || col1 == '目標' || col1 == '執行' || col1 == '驗證 / 參考來源'
    next if col1.start_with?('`<') # skip template placeholders
    # Only capture the 3 core fields: 目標, 執行, 驗證
    if %w[目標 執行 驗證].include?(col1)
      core_fields << { 'field' => col1, 'description' => col2 }
    end
  end

  # Extract usage scenarios table (情境 | 要求)
  scenarios = []
  content.scan(/^\|\s*(.+?)\s*\|\s*(.+?)\s*\|$/) do |match|
    col1 = match[0].strip
    col2 = match[1].strip
    next if col1 == '---' || col1 == '情境' || col1 == '要求'
    next if col1 == '欄位' || col1 == '必填內容'
    next if col1 == '工作類型' || col1 == '驗證方式'
    next if col1 == '工作單元' || col1 == '目標' || col1 == '執行' || col1 == '驗證 / 參考來源'
    next if col1.start_with?('`<')
    # These are the scenario rows (modify files, produce analysis, etc.)
    scenarios << { 'scenario' => col1, 'requirement' => col2 }
  end

  # Extract validation gates from "驗證 Gate 參考" section
  gate_section = content[/### 驗證 Gate 參考\n(.+?)(?=\n## |\n### |\z)/m]
  gates = []
  if gate_section
    gate_section.scan(/^\d+\.\s+(.+)$/) do |match|
      gates << { 'rule' => match[0].strip }
    end
  end

  # Extract verification examples table (工作類型 | 驗證方式)
  verification_examples = []
  content.scan(/^\|\s*(.+?)\s*\|\s*(.+?)\s*\|$/) do |match|
    col1 = match[0].strip
    col2 = match[1].strip
    next if col1 == '---' || col1 == '工作類型' || col1 == '驗證方式'
    next if col1 == '欄位' || col1 == '必填內容'
    next if col1 == '情境' || col1 == '要求'
    next if col1 == '工作單元' || col1 == '目標' || col1 == '執行' || col1 == '驗證 / 參考來源'
    next if col1.start_with?('`<')
    verification_examples << { 'work_type' => col1, 'verification_method' => col2 }
  end

  target = target_path_for(source_path, mapping_entry)
  header = generated_header(source_path)

  yaml_content = {
    'header' => header,
    'compiled_from' => source_path,
    'core_fields' => core_fields,
    'scenarios' => scenarios,
    'gates' => gates,
    'verification_examples' => verification_examples
  }

  FileUtils.mkdir_p(File.dirname(target))
  File.write(target, YAML.dump(yaml_content))
  puts "  ✓ #{target}"
end

def compile_failure_recovery(source_path, mapping_entry)
  content = File.read(source_path)

  # Extract failure taxonomy classes
  taxonomy = []
  content.scan(/^\|\s*`([^`]+)`\s*\|\s*(.+?)\s*\|\s*(.+?)\s*\|$/) do |match|
    taxonomy << {
      'class' => match[0].strip,
      'meaning' => match[1].strip,
      'common_prevention' => match[2].strip
    }
  end

  # Extract storage rules
  storage_rules = []
  content.scan(/^\|\s*(.+?)\s*\|\s*(.+?)\s*\|$/) do |match|
    next if match[0] == '---' || match[0] == '內容' || match[0] == 'Durable location'
    storage_rules << { 'content_type' => match[0].strip, 'location' => match[1].strip }
  end

  # Extract promotion decisions
  promotion_decisions = []
  content.scan(/^\|\s*(.+?)\s*\|\s*(.+?)\s*\|$/) do |match|
    next if match[0] == '---' || match[0] == 'Failure scope' || match[0] == 'Promotion target'
    promotion_decisions << { 'failure_scope' => match[0].strip, 'promotion_target' => match[1].strip }
  end

  # Extract core loop steps (capture → classify → contain → promote → strengthen → validate)
  loop_steps = []
  content.scan(/^\d+\.\s+\*\*([^*]+)\*\*[：:]\s*(.+)$/) do |match|
    loop_steps << { 'step' => match[0].strip, 'description' => match[1].strip }
  end

  target = target_path_for(source_path, mapping_entry)
  header = generated_header(source_path)

  yaml_content = {
    'header' => header,
    'compiled_from' => source_path,
    'failure_taxonomy' => taxonomy,
    'storage_rules' => storage_rules,
    'promotion_decisions' => promotion_decisions,
    'loop_steps' => loop_steps
  }

  FileUtils.mkdir_p(File.dirname(target))
  File.write(target, YAML.dump(yaml_content))
  puts "  ✓ #{target}"
end

def compile_plans_index(source_path, _mapping_entry)
  content = File.read(source_path)

  # Extract plan_id from filename (YYYY-MM-DD-HHMM-<slug>.md)
  filename = File.basename(source_path, '.md')
  plan_id = filename.sub(/^\d{4}-\d{2}-\d{2}-\d{4}-/, '')

  # Extract status from front matter (> **狀態**: draft/in-progress/completed)
  status = 'unknown'
  status_match = content.match(/>\s*\*\*狀態\*\*[：:]\s*(.+)$/)
  status = status_match[1].strip if status_match

  # Extract title from first # heading
  title = ''
  title_match = content.match(/^#\s+(.+)$/)
  title = title_match[1].strip if title_match

  # Extract purpose from front matter (> **目的**: ...)
  purpose = ''
  purpose_match = content.match(/>\s*\*\*目的\*\*[：:]\s*(.+)$/)
  purpose = purpose_match[1].strip if purpose_match

  # Fallback: if no front matter status, check if file is in archived/
  if status == 'unknown'
    status = source_path.include?('/archived/') ? 'completed' : 'draft'
  end

  # Extract phases from ### Phase N: Title (P0/P1/P2) headings
  # Supports both half-width (P0) and full-width （P0）parentheses
  phases = []
  content.scan(/^###\s+Phase\s+(\d+)[：:]\s*(.+?)$/) do |match|
    phase_num = match[0].strip.to_i
    phase_title_line = match[1].strip

    # Extract priority from parenthetical (P0/P1/P2/P3) — both half and full width
    # Matches formats: (P0), （P0）, (P0 — 立即), （P0 — 立即）
    priority = 'P?'
    priority_match = phase_title_line.match(/[（(]P(\d)/)
    priority = "P#{priority_match[1]}" if priority_match

    # Clean title (remove priority annotation like (P0) or （P0 — 立即）)
    phase_title = phase_title_line.gsub(/[（(]P\d[^)）]*[)）]/, '').strip

    phases << {
      'phase' => phase_num,
      'title' => phase_title,
      'priority' => priority
    }
  end

  # Extract affected files from tables with "檔案 | 變更類型 | Phase" header
  affected_files = []
  in_affected_table = false
  content.each_line do |line|
    if line.match?(/^\|\s*檔案\s*\|\s*變更類型\s*\|\s*Phase\s*\|$/)
      in_affected_table = true
      next
    end
    if in_affected_table
      break unless line.match?(/^\|.+\|.+\|.+\|$/)
      next if line.match?(/^\|[\s-]+\|[\s-]+\|[\s-]+\|$/)
      cols = line.split('|').map(&:strip).reject(&:empty?)
      next if cols.length < 3
      affected_files << {
        'path' => cols[0],
        'change_type' => cols[1],
        'phase' => cols[2]
      }
    end
  end

  {
    'plan_id' => plan_id,
    'filename' => filename,
    'title' => title,
    'status' => status,
    'purpose' => purpose,
    'phases' => phases,
    'affected_files' => affected_files
  }
end

def compile_classification_rules(_source_path, _mapping_entry)
  # ── Step 1: Read intelligence/engineering/README.md for known dimensions ──
  eng_readme_path = 'intelligence/engineering/README.md'
  dimensions = []
  if File.exist?(eng_readme_path)
    content = File.read(eng_readme_path)
    # Extract dimension rows from the subdirectory table
    content.scan(/^\|\s*\[`([^`]+)`\]\(([^)]+)\)\s*\|\s*(.+?)\s*\|$/) do |match|
      dir_name = match[0].strip
      rel_path = match[1].strip
      description = match[2].strip
      # Skip non-dimension entries (like the table header)
      next if dir_name == '子目錄' || dir_name == '---'
      # Determine subdirectories by scanning the actual directory
      subdirs = []
      full_path = File.join(File.dirname(eng_readme_path), dir_name)
      if Dir.exist?(full_path)
        Dir.entries(full_path).each do |entry|
          next if entry.start_with?('.')
          next unless File.directory?(File.join(full_path, entry))
          subdirs << entry
        end
      end
      dimensions << {
        'name' => dir_name,
        'description' => description,
        'path' => rel_path,
        'subdirectories' => subdirs.sort
      }
    end
  end

  # ── Step 2: Read language-specific/README.md for known languages ──
  lang_readme_path = 'intelligence/engineering/language-specific/README.md'
  known_languages = []
  if File.exist?(lang_readme_path)
    content = File.read(lang_readme_path)
    # Extract language rows from the "Current Languages" table
    in_lang_table = false
    content.each_line do |line|
      if line.match?(/^\|\s*Language\s*\|\s*Directory\s*\|\s*Atoms\s*\|$/)
        in_lang_table = true
        next
      end
      if in_lang_table
        break unless line.match?(/^\|.+\|.+\|.+\|$/)
        next if line.match?(/^\|[\s-]+\|[\s-]+\|[\s-]+\|$/)
        cols = line.split('|').map(&:strip).reject(&:empty?)
        next if cols.length < 3
        known_languages << {
          'name' => cols[0],
          'path' => "intelligence/engineering/language-specific/#{cols[0].downcase}/",
          'atoms' => cols[2].to_s
        }
      end
    end
  end

  # ── Step 3: Scan for framework-specific/ and platform-specific/ directories ──
  known_frameworks = []
  known_platforms = []
  eng_dir = 'intelligence/engineering'
  if Dir.exist?(eng_dir)
    Dir.entries(eng_dir).each do |entry|
      next if entry.start_with?('.')
      full_path = File.join(eng_dir, entry)
      next unless File.directory?(full_path)

      if entry == 'framework-specific'
        Dir.entries(full_path).each do |fw|
          next if fw.start_with?('.')
          next unless File.directory?(File.join(full_path, fw))
          known_frameworks << { 'name' => fw, 'path' => "#{eng_dir}/#{entry}/#{fw}/" }
        end
      elsif entry == 'platform-specific'
        Dir.entries(full_path).each do |pl|
          next if pl.start_with?('.')
          next unless File.directory?(File.join(full_path, pl))
          known_platforms << { 'name' => pl, 'path' => "#{eng_dir}/#{entry}/#{pl}/" }
        end
      end
    end
  end

  # ── Step 4: Extract decision tree from knowledge-update-flow.md Step 2.4 ──
  decision_tree = []
  kuf_path = 'governance/lifecycle/knowledge-update-flow.md'
  if File.exist?(kuf_path)
    content = File.read(kuf_path)
    # Find Step 2.4 section
    step_section = content[/^###\s+2\.4\s.*?\n(.*?)(?=^###\s+2\.5|\z)/m]
    if step_section
      # Extract decision branches from code blocks
      step_section.scan(/├─\s*(.+?)$/) do |match|
        branch = match[0].strip
        next if branch.start_with?('─') || branch.empty?
        decision_tree << { 'branch' => branch }
      end
      # Also extract the "→ 考慮建立" patterns
      step_section.scan(/→\s*(.+?)$/) do |match|
        action = match[0].strip
        decision_tree << { 'action' => action } unless action.empty?
      end
    end
  end

  # ── Step 5: Write classification-rules.yaml ──
  target = File.join(GENERATED_DIR, 'classification-rules.yaml')
  header = generated_header('governance/lifecycle/knowledge-update-flow.md')

  yaml_content = {
    'header' => header,
    'compiled_from' => 'governance/lifecycle/knowledge-update-flow.md + intelligence/engineering/README.md',
    'classification_dimensions' => dimensions,
    'known_languages' => known_languages,
    'known_frameworks' => known_frameworks,
    'known_platforms' => known_platforms,
    'decision_tree' => decision_tree
  }

  FileUtils.mkdir_p(GENERATED_DIR)
  File.write(target, YAML.dump(yaml_content))
  puts "  ✓ #{target}"
end

def compile_source(source_path, mapping_entry)
  compile_rule = mapping_entry['compile_rule']

  case compile_rule
  when /從 execution-flow 的章節標題提取 phase 定義/
    compile_workflow_phases(source_path, mapping_entry)
  when /從 writeback transaction 章節提取 state machine 定義/
    compile_enforcement_transactions(source_path, mapping_entry)
  when /提取 language policy 定義|提取 sanitization 定義|提取 tool neutrality 定義/
    compile_output_governance(source_path, mapping_entry)
  when /從 knowledge-update-flow\.md 的 11 個步驟標題與判斷表格提取 phase 定義/
    compile_knowledge_update_flow(source_path, mapping_entry)
  when /從 artifact gates 的檢查清單提取 required artifacts 與 verification criteria/
    compile_workflow_artifacts(source_path, mapping_entry)
  when /從 validation gate 描述提取 gate 定義/
    compile_goal_action_gates(source_path, mapping_entry)
  when /從 failure taxonomy 與 recovery 描述提取 pattern 與 strategy/
    compile_failure_recovery(source_path, mapping_entry)
  when /從 plans\/active\/\*\.md 的 front matter、phase 標題、受影響檔案表格提取 plan index/
    compile_plans_index(source_path, mapping_entry)
  when /從 knowledge-update-flow\.md Step 2\.4 的決策樹與 intelligence\/engineering\/ 的 README 提取分類維度定義/
    compile_classification_rules(source_path, mapping_entry)
  when /從 analysis\/apk\/workflows\/\*\.md 的「步驟 N：」標題提取 step 定義/
    compile_apk_workflow_phases(source_path, mapping_entry)
  else
    puts "  ⚠  Unknown compile rule: #{compile_rule}"
  end
end

def check_modified_sources
  modified = []
  plans_modified = false
  apk_workflow_modified = false
  plans_target = File.join(GENERATED_DIR, 'plans-index.yaml')
  classification_target = File.join(GENERATED_DIR, 'classification-rules.yaml')
  apk_mapping = @mapping.find { |e| e['compile_rule']&.include?('analysis/apk/workflows/*.md') }
  apk_workflow_target = apk_mapping ? apk_mapping['target'] : File.join(GENERATED_DIR, 'apk-workflow-phases.yaml')

  @mapping.each do |entry|
    source_glob = entry['source']
    Dir.glob(source_glob).each do |source_path|
      # Plans index: check if any plan file is newer than the aggregated target
      if entry['compile_rule']&.include?('plans/active/*.md')
        if !File.exist?(plans_target) || File.mtime(source_path) > File.mtime(plans_target)
          plans_modified = true
        end
      # APK workflow: check if any APK workflow file is newer than the aggregated target
      elsif entry['compile_rule']&.include?('analysis/apk/workflows/*.md')
        # Skip README — not a workflow
        next if File.basename(source_path) == 'README.md'
        if !File.exist?(apk_workflow_target) || File.mtime(source_path) > File.mtime(apk_workflow_target)
          apk_workflow_modified = true
        end
      # Classification rules: also depends on intelligence/engineering/ README files
      elsif entry['compile_rule']&.include?('分類維度定義')
        target = classification_target
        # Check all dependency files
        deps = [
          source_path,
          'intelligence/engineering/README.md',
          'intelligence/engineering/language-specific/README.md'
        ]
        deps.each do |dep|
          if File.exist?(dep) && (!File.exist?(target) || File.mtime(dep) > File.mtime(target))
            modified << { source: source_path, mapping: entry }
            break
          end
        end
      else
        target = target_path_for(source_path, entry)
        if !File.exist?(target) || File.mtime(source_path) > File.mtime(target)
          modified << { source: source_path, mapping: entry }
        end
      end
    end
  end

  if plans_modified
    modified << { source: 'plans/active/*.md', mapping: @mapping.find { |e| e['compile_rule']&.include?('plans/active/*.md') } }
  end

  if apk_workflow_modified
    modified << { source: 'analysis/apk/workflows/*.md', mapping: @mapping.find { |e| e['compile_rule']&.include?('analysis/apk/workflows/*.md') } }
  end

  modified
end

def run(options)
  load_mapping

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
  puts "Compiling prose sources to generated YAML..."
  puts

  # Track plans index entries and APK workflow entries for aggregation
  plans_entries = []
  apk_workflow_entries = []

  @mapping.each do |entry|
    source_glob = entry['source']
    Dir.glob(source_glob).each do |source_path|
      # Special handling for plans index: aggregate all plans into one file
      if entry['compile_rule']&.include?('plans/active/*.md')
        plans_entries << compile_plans_index(source_path, entry)
      # Special handling for APK workflows: aggregate all workflows into one file
      elsif entry['compile_rule']&.include?('analysis/apk/workflows/*.md')
        apk_workflow_entries << compile_apk_workflow_phases(source_path, entry)
      else
        compile_source(source_path, entry)
      end
    end
  end

  # Write aggregated plans index
  unless plans_entries.empty?
    target = File.join(GENERATED_DIR, 'plans-index.yaml')
    header = generated_header('plans/active/*.md')
    yaml_content = {
      'header' => header,
      'compiled_from' => 'plans/active/*.md',
      'total_plans' => plans_entries.length,
      'plans' => plans_entries
    }
    File.write(target, YAML.dump(yaml_content))
    puts "  ✓ #{target}"
  end

  # Write aggregated APK workflow phases
  apk_workflow_entries = apk_workflow_entries.compact
  unless apk_workflow_entries.empty?
    apk_mapping = @mapping.find { |e| e['compile_rule']&.include?('analysis/apk/workflows/*.md') }
    target = apk_mapping ? apk_mapping['target'] : File.join(GENERATED_DIR, 'apk-workflow-phases.yaml')
    header = generated_header('analysis/apk/workflows/*.md')
    yaml_content = {
      'header' => header,
      'compiled_from' => 'analysis/apk/workflows/*.md',
      'total_workflows' => apk_workflow_entries.length,
      'workflows' => apk_workflow_entries
    }
    File.write(target, YAML.dump(yaml_content))
    puts "  ✓ #{target}"
  end

  puts
  puts "Compilation complete."
end

# CLI entry point
options = {}
OptionParser.new do |opts|
  opts.banner = "Usage: compiler-engine.rb [options]"
  opts.on('--check', 'Check if compilation is needed') { |v| options[:check] = v }
  opts.on('--diff', 'Show what would change') { |v| options[:diff] = v }
end.parse!

run(options)
