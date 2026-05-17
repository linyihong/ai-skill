#!/usr/bin/env ruby
# frozen_string_literal: true

# validate-runtime-db.rb
#
# Validates runtime.db integrity and data consistency.
#
# Checks:
#   1. File exists and is a valid SQLite database
#   2. All required tables exist
#   3. Required tables have minimum row counts
#   4. Foreign key-like consistency (e.g., gates reference valid phases)
#   5. JSON columns contain valid JSON
#   6. Compiler metadata is present and recent
#
# Usage:
#   ruby scripts/validate-runtime-db.rb
#   ruby scripts/validate-runtime-db.rb --db /path/to/runtime.db

require 'open3'
require 'optparse'
require 'json'
require 'time'

DEFAULT_DB_PATH = File.join(File.dirname(__FILE__), '..', 'runtime', 'runtime.db')

REQUIRED_TABLES = %w[
  phases phase_transitions obligations gates
  transaction_states transaction_transitions transaction_rules transaction_templates
  activation_rules core_bootstrap_rules
  discovery_checkpoints discovery_search_strategy
  generated_surfaces compiler_metadata
  runtime_budget context_ttl_policy circuit_breaker context_pollution
  context_health_score intelligence_routing obligation_ledger
  language_policy output_rules governance_gates blocking_gates
  phase_machine pipeline_context_flow guard_chain relevance_engine
  session_lifecycle prompt_artifact_templates prompt_composition_rules
  recovery_strategies state_repair obligation_rebuild phase_reconciliation
  execution_queue priority_scheduler activation_rules_mirror
  transaction_templates_ext distributed_locks multi_agent_coordination
  async_job_lifecycle capability_checkpoints
].freeze

MINIMUM_ROWS = {
  'phases' => 8,
  'obligations' => 15,
  'gates' => 15,
  'activation_rules' => 10,
  'core_bootstrap_rules' => 2,
  'discovery_checkpoints' => 3,
  'compiler_metadata' => 2,
  # Runtime Config Tables (v2 migration) — each should have at least 1 row (__config__)
  'runtime_budget' => 1,
  'context_ttl_policy' => 1,
  'circuit_breaker' => 1,
  'context_pollution' => 1,
  'context_health_score' => 1,
  'intelligence_routing' => 1,
  'obligation_ledger' => 1,
  'language_policy' => 1,
  'output_rules' => 1,
  'governance_gates' => 1,
  'blocking_gates' => 1,
  'phase_machine' => 1,
  'pipeline_context_flow' => 1,
  'guard_chain' => 1,
  'relevance_engine' => 1,
  'session_lifecycle' => 1,
  'prompt_artifact_templates' => 1,
  'prompt_composition_rules' => 1,
  'recovery_strategies' => 1,
  'state_repair' => 1,
  'obligation_rebuild' => 1,
  'phase_reconciliation' => 1,
  'execution_queue' => 1,
  'priority_scheduler' => 1,
  'activation_rules_mirror' => 1,
  'transaction_templates_ext' => 1,
  'distributed_locks' => 1,
  'multi_agent_coordination' => 1,
  'async_job_lifecycle' => 1,
  'capability_checkpoints' => 1
}.freeze

JSON_COLUMNS = {
  'phases' => %w[entry_conditions allowed_actions forbidden_actions blocking_gates obligations phase_transition_triggers],
  'obligations' => %w[verification depends_on linked_gates],
  'transaction_states' => %w[entry_conditions allowed_actions forbidden_actions blocking_gates],
  'activation_rules' => %w[activation_when],
  'discovery_checkpoints' => %w[discovery_targets]
}.freeze

ERRORS = []
WARNINGS = []

def sqlite_query(db_path, sql)
  o, e, s = Open3.capture3('sqlite3', '-json', db_path, sql)
  unless s.success?
    ERRORS << "SQL error: #{e.strip} (SQL: #{sql[0..80]})"
    return []
  end
  return [] if o.strip.empty?
  JSON.parse(o)
rescue JSON::ParserError => e
  ERRORS << "JSON parse error from sqlite3 output: #{e.message}"
  []
end

def sqlite_raw(db_path, sql)
  o, e, s = Open3.capture3('sqlite3', db_path, sql)
  ERRORS << "SQL error: #{e.strip}" unless s.success?
  o.strip
end

def check_file_exists(db_path)
  unless File.exist?(db_path)
    ERRORS << "File not found: #{db_path}"
    return false
  end
  if File.size(db_path) < 1024
    ERRORS << "File too small (#{File.size(db_path)} bytes): #{db_path}"
  end
  true
end

def check_sqlite_integrity(db_path)
  result = sqlite_raw(db_path, 'PRAGMA integrity_check;')
  if result == 'ok'
    puts "  ✓ integrity_check"
  else
    ERRORS << "Integrity check failed: #{result}"
  end
end

def check_tables_exist(db_path)
  tables = sqlite_raw(db_path, ".tables").split
  REQUIRED_TABLES.each do |t|
    if tables.include?(t)
      puts "  ✓ table #{t}"
    else
      ERRORS << "Missing required table: #{t}"
    end
  end
end

def check_row_counts(db_path)
  MINIMUM_ROWS.each do |table, min|
    count = sqlite_raw(db_path, "SELECT COUNT(*) FROM #{table};").to_i
    if count >= min
      puts "  ✓ #{table}: #{count} rows (min #{min})"
    else
      ERRORS << "#{table}: #{count} rows, expected at least #{min}"
    end
  end
end

def check_json_columns(db_path)
  JSON_COLUMNS.each do |table, columns|
    rows = sqlite_query(db_path, "SELECT #{columns.join(', ')} FROM #{table} LIMIT 5;")
    columns.each do |col|
      rows.each_with_index do |row, i|
        val = row[col]
        next if val.nil? || val == '[]' || val == '{}'
        begin
          JSON.parse(val) if val.is_a?(String)
        rescue JSON::ParserError
          ERRORS << "#{table}.#{col} (row #{i}): invalid JSON"
        end
      end
    end
  end
  puts "  ✓ JSON columns validated" unless ERRORS.any? { |e| e.include?('invalid JSON') }
end

def check_phase_references(db_path)
  # Gates reference valid phases
  gates = sqlite_query(db_path, "SELECT DISTINCT phase FROM gates;")
  phases = sqlite_raw(db_path, "SELECT GROUP_CONCAT(id) FROM phases;").split(',')
  gates.each do |g|
    unless phases.include?(g['phase'])
      WARNINGS << "Gate references unknown phase: #{g['phase']}"
    end
  end

  # Obligations reference valid phases
  obligations = sqlite_query(db_path, "SELECT DISTINCT phase FROM obligations;")
  obligations.each do |o|
    unless phases.include?(o['phase'])
      WARNINGS << "Obligation references unknown phase: #{o['phase']}"
    end
  end

  puts "  ✓ phase references checked"
end

def check_compiler_metadata(db_path)
  meta = sqlite_query(db_path, "SELECT key, value FROM compiler_metadata;")
  meta_map = meta.each_with_object({}) { |r, h| h[r['key']] = r['value'] }

  unless meta_map['compiler_version']
    ERRORS << "compiler_metadata missing 'compiler_version'"
    return
  end

  unless meta_map['compiled_at']
    ERRORS << "compiler_metadata missing 'compiled_at'"
    return
  end

  # Check compilation is recent (within 24 hours)
  compiled = Time.parse(meta_map['compiled_at'])
  age_hours = (Time.now - compiled) / 3600
  if age_hours > 24
    WARNINGS << "runtime.db is #{age_hours.round(1)} hours old (compiled at #{meta_map['compiled_at']})"
  end

  puts "  ✓ compiler: v#{meta_map['compiler_version']} at #{meta_map['compiled_at']}"
end

# ---- CLI ----
db_path = DEFAULT_DB_PATH

OptionParser.new do |opts|
  opts.banner = "Usage: #{$PROGRAM_NAME} [options]"
  opts.on('--db PATH', 'Path to runtime.db (default: runtime/runtime.db)') { |v| db_path = v }
  opts.on('-h', '--help', 'Show help') { puts opts; exit }
end.parse!

puts "Validating runtime.db: #{db_path}"
puts ""

exit 1 unless check_file_exists(db_path)

check_sqlite_integrity(db_path)
check_tables_exist(db_path)
check_row_counts(db_path)
check_json_columns(db_path)
check_phase_references(db_path)
check_compiler_metadata(db_path)

puts ""
if ERRORS.empty?
  puts "✓ All checks passed"
else
  puts "✗ #{ERRORS.size} error(s):"
  ERRORS.each { |e| puts "  ✗ #{e}" }
end

unless WARNINGS.empty?
  puts "  ⚠ #{WARNINGS.size} warning(s):"
  WARNINGS.each { |w| puts "  ⚠ #{w}" }
end

exit ERRORS.empty? ? 0 : 1
