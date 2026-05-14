#!/usr/bin/env ruby
# frozen_string_literal: true

# Activation Engine
# 讀取 runtime/router/activation-rules.yaml，根據輸入的 task intent、
# file change、user signal 等條件，輸出應該 activate 的 rule 列表。
#
# Usage:
#   # 顯示所有規則的 activation 狀態（無輸入）
#   ./runtime/router/activation-engine.rb
#
#   # 指定 task intent
#   ./runtime/router/activation-engine.rb --intent migration
#   ./runtime/router/activation-engine.rb --intent debug --intent refactor
#
#   # 指定 file changes
#   ./runtime/router/activation-engine.rb --file-changed enforcement/rule-weight.md
#   ./runtime/router/activation-engine.rb --file-changed README.md --file-changed enforcement/linked-updates.md
#
#   # 指定 user signal
#   ./runtime/router/activation-engine.rb --signal 連動
#   ./runtime/router/activation-engine.rb --signal 錯誤 --signal 失誤
#
#   # 複合條件
#   ./runtime/router/activation-engine.rb --intent migration --file-changed "**/*.md" --file-changed enforcement/linked-updates.md
#
#   # Dry-run 模式（顯示判斷邏輯）
#   ./runtime/router/activation-engine.rb --intent migration --dry-run
#
#   # 列出所有已知的 intent/signal/file pattern
#   ./runtime/router/activation-engine.rb --list-known

require 'yaml'
require 'optparse'
require 'fileutils'

ENGINE_VERSION = 'v1'
ACTIVATION_RULES_PATH = File.join(__dir__, 'activation-rules.yaml')

# ──────────────────────────────────────────────
# 資料結構
# ──────────────────────────────────────────────

Rule = Struct.new(
  :rule_id,
  :description,
  :activation_conditions,  # Array of condition hashes
  :load_strategy,
  :priority,
  :estimated_tokens,
  :source,
  keyword_init: true
)

ActivationResult = Struct.new(
  :rule_id,
  :description,
  :activated,
  :matched_conditions,  # Array of condition descriptions that matched
  :priority,
  :estimated_tokens,
  :source,
  keyword_init: true
)

# ──────────────────────────────────────────────
# 載入規則
# ──────────────────────────────────────────────

def load_activation_rules(path = ACTIVATION_RULES_PATH)
  data = YAML.safe_load(File.read(path))
  rules_data = data['rules'] || []
  core_bootstrap = data['core_bootstrap'] || []

  rules = rules_data.map do |r|
    Rule.new(
      rule_id: r['rule_id'],
      description: r['description'],
      activation_conditions: r['activation']['when'] || [],
      load_strategy: r['load']['strategy'] || 'lazy',
      priority: r['load']['priority'] || 'P3',
      estimated_tokens: r['load']['estimated_tokens'] || 0,
      source: r['load']['source'] || ''
    )
  end

  { rules: rules, core_bootstrap: core_bootstrap, version: data['activation_rules_version'], status: data['status'] }
end

# ──────────────────────────────────────────────
# 條件比對引擎
# ──────────────────────────────────────────────

# 比對單一 condition 是否滿足
def condition_matches?(condition, inputs)
  type = condition['type']
  case type
  when 'task_intent'
    intents = inputs[:intents] || []
    matches = condition['matches'] || []
    # 任一 intent 符合任一 match pattern（substring match）
    intents.any? { |i| matches.any? { |m| i.downcase.include?(m.downcase) } }

  when 'user_signal'
    signals = inputs[:signals] || []
    matches = condition['matches'] || []
    signals.any? { |s| matches.any? { |m| s.downcase.include?(m.downcase) } }

  when 'file_change'
    changed_files = inputs[:changed_files] || []
    pattern = condition['pattern'] || ''
    min_count = (condition['count'] || '>=1').sub('>=', '').to_i
    # 將 glob pattern 轉為 regex（簡化版：只處理 **/ 前綴和後綴）
    regex = glob_to_regex(pattern)
    matched = changed_files.count { |f| f.match?(regex) }
    matched >= min_count

  when 'file_has_todo'
    # 此條件需要掃描文件內容，這裡假設外部已提供 has_todo_files 列表
    has_todo_files = inputs[:has_todo_files] || []
    pattern = condition['pattern'] || '**/*.md'
    regex = glob_to_regex(pattern)
    has_todo_files.any? { |f| f.match?(regex) }

  when 'task_complexity'
    routes = inputs[:routes] || 1
    threshold_str = condition['routes'] || '>=1'
    threshold = threshold_str.sub('>=', '').to_i
    routes >= threshold

  when 'validation_gap'
    signals = inputs[:signals] || []
    matches = condition['matches'] || []
    signals.any? { |s| matches.any? { |m| s.downcase.include?(m.downcase) } }

  when 'file_size'
    # 此條件需要實際檢查文件行數，這裡假設外部已提供 oversized_files 列表
    oversized = inputs[:oversized_files] || []
    oversized.any?

  when 'tool_active'
    active_tools = inputs[:active_tools] || []
    tool_name = condition['name'] || ''
    active_tools.include?(tool_name)

  else
    false
  end
end

# 簡化版 glob → regex 轉換（支援 **/ 前綴和後綴）
def glob_to_regex(pattern)
  regex_str = pattern
    .gsub('.', '\\.')
    .gsub('**/', '(.+/)?')
    .gsub('*', '[^/]*')
  /\A#{regex_str}\z/
end

# ──────────────────────────────────────────────
# 判斷引擎
# ──────────────────────────────────────────────

def evaluate_rules(rules, inputs, dry_run: false)
  results = []

  rules.each do |rule|
    matched_conditions = []

    rule.activation_conditions.each do |condition|
      if condition_matches?(condition, inputs)
        matched_conditions << describe_condition(condition)
      end
    end

    activated = !matched_conditions.empty?

    if dry_run
      puts "  [#{activated ? 'ACTIVATE' : 'SKIP'}] #{rule.rule_id}"
      puts "        描述: #{rule.description}"
      if matched_conditions.any?
        puts "        匹配條件:"
        matched_conditions.each { |mc| puts "          - #{mc}" }
      else
        puts "        未匹配任何條件"
      end
      puts "        載入策略: #{rule.load_strategy} | 優先權: #{rule.priority} | Tokens: #{rule.estimated_tokens}"
      puts "        來源: #{rule.source}"
      puts
    end

    results << ActivationResult.new(
      rule_id: rule.rule_id,
      description: rule.description,
      activated: activated,
      matched_conditions: matched_conditions,
      priority: rule.priority,
      estimated_tokens: rule.estimated_tokens,
      source: rule.source
    )
  end

  results
end

def describe_condition(condition)
  type = condition['type']
  case type
  when 'task_intent'
    "task_intent 匹配: #{condition['matches'].join(', ')}"
  when 'user_signal'
    "user_signal 匹配: #{condition['matches'].join(', ')}"
  when 'file_change'
    "file_change 匹配: #{condition['pattern']} (count>=#{condition['count']&.sub('>=', '') || 1})"
  when 'file_has_todo'
    "file_has_todo 匹配: #{condition['pattern']}"
  when 'task_complexity'
    "task_complexity: routes #{condition['routes']}"
  when 'validation_gap'
    "validation_gap 匹配: #{condition['matches'].join(', ')}"
  when 'file_size'
    "file_size: #{condition['threshold']}"
  when 'tool_active'
    "tool_active: #{condition['name']}"
  else
    "unknown type: #{type}"
  end
end

# ──────────────────────────────────────────────
# 輸出格式化
# ──────────────────────────────────────────────

def format_results(results, core_bootstrap)
  activated = results.select(&:activated)
  skipped = results.reject(&:activated)

  puts "=" * 60
  puts "Activation Engine #{ENGINE_VERSION}"
  puts "=" * 60
  puts

  # Core Bootstrap（永遠 active）
  puts "▶ Core Bootstrap（永遠 preload）:"
  core_bootstrap.each { |r| puts "   #{r}" }
  puts

  # Activated lazy rules
  if activated.any?
    puts "▶ Activated Lazy-load Rules (#{activated.size}):"
    activated.sort_by { |r| priority_weight(r.priority) }.each do |r|
      puts "   [#{r.priority}] #{r.rule_id}"
      puts "          描述: #{r.description}"
      puts "          Tokens: #{r.estimated_tokens}"
      puts "          來源: #{r.source}"
      r.matched_conditions.each { |mc| puts "          ✓ #{mc}" }
      puts
    end
  else
    puts "▶ 無 Lazy-load Rules 被 activate"
    puts
  end

  # Summary
  total_tokens = activated.sum(&:estimated_tokens)
  puts "-" * 40
  puts "Core Bootstrap: #{core_bootstrap.size} rules (~800 tokens)"
  puts "Activated Lazy: #{activated.size} rules (#{total_tokens} tokens)"
  puts "Skipped:        #{skipped.size} rules"
  puts "Total (est.):   #{800 + total_tokens} tokens"
  puts "=" * 60
end

def priority_weight(priority)
  case priority
  when 'P0' then 0
  when 'P1' then 1
  when 'P2' then 2
  when 'P3' then 3
  else 99
  end
end

# ──────────────────────────────────────────────
# 已知條件列表
# ──────────────────────────────────────────────

def list_known_conditions(rules)
  intents = Set.new
  signals = Set.new
  file_patterns = Set.new

  rules.each do |rule|
    rule.activation_conditions.each do |cond|
      (cond['matches'] || []).each { |m| intents.add(m) } if cond['type'] == 'task_intent'
      (cond['matches'] || []).each { |m| signals.add(m) } if cond['type'] == 'user_signal'
      (cond['matches'] || []).each { |m| signals.add(m) } if cond['type'] == 'validation_gap'
      file_patterns.add("#{cond['pattern']} (#{cond['count'] || '>=1'})") if cond['type'] == 'file_change'
      file_patterns.add(cond['pattern']) if cond['type'] == 'file_has_todo'
    end
  end

  puts "已知 Task Intents:"
  intents.sort.each { |i| puts "  - #{i}" }
  puts
  puts "已知 User Signals / Validation Gaps:"
  signals.sort.each { |s| puts "  - #{s}" }
  puts
  puts "已知 File Patterns:"
  file_patterns.sort.each { |p| puts "  - #{p}" }
end

# ──────────────────────────────────────────────
# CLI Entrypoint
# ──────────────────────────────────────────────

if __FILE__ == $PROGRAM_NAME
  options = {
    intents: [],
    signals: [],
    changed_files: [],
    has_todo_files: [],
    routes: 1,
    oversized_files: [],
    active_tools: [],
    dry_run: false,
    list_known: false
  }

  OptionParser.new do |opts|
    opts.banner = "Usage: activation-engine.rb [options]"

    opts.on('-i', '--intent INTENT', 'Task intent (can be specified multiple times)') do |v|
      options[:intents] << v
    end

    opts.on('-s', '--signal SIGNAL', 'User signal (can be specified multiple times)') do |v|
      options[:signals] << v
    end

    opts.on('-f', '--file-changed FILE', 'Changed file path (can be specified multiple times)') do |v|
      options[:changed_files] << v
    end

    opts.on('--has-todo FILE', 'File containing TODO markers (can be specified multiple times)') do |v|
      options[:has_todo_files] << v
    end

    opts.on('-r', '--routes COUNT', Integer, 'Number of available routes (for task_complexity)') do |v|
      options[:routes] = v
    end

    opts.on('--oversized FILE', 'Oversized file (can be specified multiple times)') do |v|
      options[:oversized_files] << v
    end

    opts.on('-t', '--tool NAME', 'Active tool name (e.g., cursor)') do |v|
      options[:active_tools] << v
    end

    opts.on('-n', '--dry-run', 'Dry-run mode: show matching logic') do |v|
      options[:dry_run] = v
    end

    opts.on('-l', '--list-known', 'List all known intents, signals, and file patterns') do |v|
      options[:list_known] = v
    end

    opts.on('-h', '--help', 'Show help') do
      puts opts
      exit
    end
  end.parse!

  # 載入規則
  config = load_activation_rules
  rules = config[:rules]
  core_bootstrap = config[:core_bootstrap]

  puts "Activation Rules: #{config[:version]} (#{config[:status]})"
  puts "Loaded #{rules.size} lazy-load rules, #{core_bootstrap.size} core bootstrap rules"
  puts

  if options[:list_known]
    list_known_conditions(rules)
    exit 0
  end

  # 準備輸入
  inputs = {
    intents: options[:intents],
    signals: options[:signals],
    changed_files: options[:changed_files],
    has_todo_files: options[:has_todo_files],
    routes: options[:routes],
    oversized_files: options[:oversized_files],
    active_tools: options[:active_tools]
  }

  has_input = options.values_at(:intents, :signals, :changed_files, :has_todo_files, :oversized_files, :active_tools).any? { |v| v.respond_to?(:any?) && v.any? } || options[:routes] > 1

  unless has_input
    puts "⚠ 未提供任何輸入條件。使用 --dry-run 查看所有規則的判斷邏輯，或提供 --intent/--signal/--file-changed 等參數。"
    puts
    options[:dry_run] = true if options[:dry_run].nil?
  end

  # 執行判斷
  results = evaluate_rules(rules, inputs, dry_run: options[:dry_run])

  # 輸出結果
  format_results(results, core_bootstrap) unless options[:dry_run]
end
