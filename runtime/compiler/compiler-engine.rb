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
  rules = YAML.safe_load(File.read(COMPILER_RULES))
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
  File.join(GENERATED_DIR, File.basename(target))
end

def extract_domain(source_path)
  # Extract domain from path like workflow/apk-analysis/execution-flow.md → apk-analysis
  match = source_path.match(%r{workflow/([^/]+)/})
  return match[1] if match

  # enforcement/dependency-reading.md → transaction-machine
  match = source_path.match(%r{enforcement/(.+)\.md$})
  return match[1].tr('-', '_') if match

  'unknown'
end

def extract_phases_from_flow(content)
  # Extract phase definitions from markdown headings (## level)
  phases = []
  content.scan(/^##\s+\d+\.\s+(.+)$/) do |match|
    phases << { 'name' => match[0].strip, 'source_line' => $`.lines.count + 1 }
  end
  phases
end

def extract_gates_from_content(content)
  # Extract blocking gates from markdown lists or gate descriptions
  gates = []
  content.scan(/^\*\*([^*]+)\*\*：(.+)$/) do |match|
    gates << { 'name' => match[0].strip, 'description' => match[1].strip }
  end
  gates
end

def compile_workflow_phases(source_path, mapping_entry)
  content = File.read(source_path)
  phases = extract_phases_from_flow(content)
  gates = extract_gates_from_content(content)

  target = target_path_for(source_path, mapping_entry)
  header = generated_header(source_path)

  yaml_content = {
    'header' => header,
    'phases' => phases,
    'gates' => gates
  }

  FileUtils.mkdir_p(File.dirname(target))
  File.write(target, YAML.dump(yaml_content))
  puts "  ✓ #{target}"
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

def compile_source(source_path, mapping_entry)
  compile_rule = mapping_entry['compile_rule']

  case compile_rule
  when /從 execution-flow 的章節標題提取 phase 定義/
    compile_workflow_phases(source_path, mapping_entry)
  when /從 writeback transaction 章節提取 state machine 定義/
    compile_enforcement_transactions(source_path, mapping_entry)
  else
    puts "  ⚠  Unknown compile rule: #{compile_rule}"
  end
end

def check_modified_sources
  modified = []
  @mapping.each do |entry|
    source_glob = entry['source']
    Dir.glob(source_glob).each do |source_path|
      target = target_path_for(source_path, entry)
      if !File.exist?(target) || File.mtime(source_path) > File.mtime(target)
        modified << { source: source_path, mapping: entry }
      end
    end
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

  @mapping.each do |entry|
    source_glob = entry['source']
    Dir.glob(source_glob).each do |source_path|
      compile_source(source_path, entry)
    end
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
