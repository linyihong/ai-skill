#!/usr/bin/env ruby
# encoding: UTF-8
# Validate Ai-skill knowledge runtime generated surfaces.

require "date"
require "pathname"
require "yaml"

ROOT = Pathname.new(__dir__).parent.realpath

ERRORS = []
COUNTS = Hash.new(0)

EDGE_TYPES = %w[
  depends_on
  related_to
  conflicts_with
  replaces
  preserves_entrypoint
  promotes_from
  routes_to
].freeze

ROUTING_RECORD_FIELDS = %w[
  id
  task_intent
  primary_source
  required_dependencies
  candidate_sources
  source_of_truth_gate
  ranking_reason
  validation_signal
  metadata
  model
].freeze

ROUTING_METADATA_FIELDS = %w[
  priority
  confidence
  context_cost
  compatibility_state
].freeze

ROUTING_MODEL_FIELDS = %w[
  profile
  compression_level
  reason
].freeze

SUMMARY_FIELDS = [
  "Atom ID",
  "Source path",
  "Lifecycle",
  "Summary",
  "When to read",
  "Do not use for",
  "Validation signal"
].freeze

def add_error(message)
  ERRORS << message
end

def read_text(path)
  File.read(path.to_s, encoding: "UTF-8")
end

def rel(path)
  Pathname.new(path).relative_path_from(ROOT).to_s
rescue ArgumentError
  path.to_s
end

def path_exists?(relative_path)
  return false if relative_path.nil? || relative_path.strip.empty?
  return true if relative_path.start_with?("#")
  return true if relative_path.match?(/\A[a-z][a-z0-9+.-]*:\/\//i)

  target = relative_path.split("#", 2).first
  return true if target.nil? || target.empty?

  (ROOT + target).exist?
end

def check_repo_path(relative_path, context)
  return if path_exists?(relative_path)

  add_error("#{context}: missing path #{relative_path}")
end

def yaml_file(relative_path)
  path = ROOT + relative_path
  YAML.safe_load(read_text(path), permitted_classes: [Date], aliases: false)
rescue Psych::SyntaxError => e
  add_error("#{relative_path}: YAML parse failed: #{e.message}")
  nil
rescue Errno::ENOENT
  add_error("#{relative_path}: file does not exist")
  nil
end

def require_mapping(value, context)
  return value if value.is_a?(Hash)

  add_error("#{context}: expected mapping")
  {}
end

def require_array(value, context)
  return value if value.is_a?(Array)

  add_error("#{context}: expected array")
  []
end

def check_fields(mapping, fields, context)
  fields.each do |field|
    next if mapping.key?(field) && !mapping[field].nil? && mapping[field] != ""

    add_error("#{context}: missing #{field}")
  end
end

def markdown_links(text)
  text.scan(/\[[^\]]+\]\(([^)]+)\)/).flatten
end

def check_markdown_links(relative_path)
  path = ROOT + relative_path
  text = read_text(path)
  markdown_links(text).each do |target|
    next if target.start_with?("#") || target.match?(/\A[a-z][a-z0-9+.-]*:\/\//i)

    target_path = target.split("#", 2).first
    next if target_path.nil? || target_path.empty?

    resolved = (path.dirname + target_path).cleanpath
    add_error("#{relative_path}: missing markdown link #{target}") unless resolved.exist?
  end
end

def parse_summary_table(text)
  fields = {}
  text.each_line do |line|
    next unless line.start_with?("|")

    cells = line.strip.split("|").map(&:strip)
    next unless cells.length >= 3
    next if cells[1] == "欄位" || cells[1].match?(/\A-+\z/)

    fields[cells[1]] = cells[2]
  end
  fields
end

def links_from_cell(cell)
  markdown_links(cell).map { |target| target.split("#", 2).first }.reject(&:empty?)
end

def validate_registry
  data = require_mapping(yaml_file("knowledge/runtime/routing-registry.yaml"), "routing-registry")
  check_fields(data, %w[registry_version status owner_layer source_of_truth_policy records], "routing-registry")
  check_repo_path(data["refresh_policy"], "routing-registry refresh_policy") if data["refresh_policy"]

  records = require_array(data["records"], "routing-registry records")
  seen = {}
  records.each_with_index do |record, index|
    context = "routing-registry records[#{index}]"
    record = require_mapping(record, context)
    check_fields(record, ROUTING_RECORD_FIELDS, context)

    id = record["id"]
    add_error("#{context}: duplicate id #{id}") if id && seen[id]
    seen[id] = true if id

    check_repo_path(record["primary_source"], "#{context} primary_source") if record["primary_source"]
    require_array(record["required_dependencies"], "#{context} required_dependencies").each do |path|
      check_repo_path(path, "#{context} required_dependencies")
    end
    require_array(record["candidate_sources"], "#{context} candidate_sources").each do |path|
      check_repo_path(path, "#{context} candidate_sources")
    end

    metadata = require_mapping(record["metadata"], "#{context} metadata")
    check_fields(metadata, ROUTING_METADATA_FIELDS, "#{context} metadata")

    model = require_mapping(record["model"], "#{context} model")
    check_fields(model, ROUTING_MODEL_FIELDS, "#{context} model")
  end
  COUNTS[:registry_records] = records.length
end

def validate_refresh_policy
  data = require_mapping(yaml_file("knowledge/runtime/refresh-policy.yaml"), "refresh-policy")
  check_fields(data, %w[policy_version status owner_layer surfaces decision_values validation output_shape], "refresh-policy")

  surfaces = require_mapping(data["surfaces"], "refresh-policy surfaces")
  %w[summaries graphs routing_registry].each do |surface|
    value = require_mapping(surfaces[surface], "refresh-policy surfaces.#{surface}")
    check_fields(value, %w[path refresh_when revalidate_when downgrade_when], "refresh-policy surfaces.#{surface}")
    check_repo_path(value["path"], "refresh-policy surfaces.#{surface}.path") if value["path"]
  end

  validation = require_mapping(data["validation"], "refresh-policy validation")
  check_fields(validation, %w[required_checks close_loop], "refresh-policy validation")
  require_array(validation["required_checks"], "refresh-policy validation.required_checks")
  require_array(validation["close_loop"], "refresh-policy validation.close_loop")
end

def validate_summaries
  summary_files = Dir.glob((ROOT + "knowledge/summaries/*.md").to_s)
                     .map { |path| Pathname.new(path) }
                     .reject { |path| path.basename.to_s == "README.md" }
                     .sort

  summary_files.each do |path|
    relative = rel(path)
    text = read_text(path)
    fields = parse_summary_table(text)
    SUMMARY_FIELDS.each do |field|
      add_error("#{relative}: missing summary field #{field}") unless fields.key?(field) && fields[field] != ""
    end

    links_from_cell(fields["Source path"].to_s).each do |target|
      resolved = (path.dirname + target).cleanpath
      add_error("#{relative}: missing source link #{target}") unless resolved.exist?
    end

    check_markdown_links(relative)
  end

  readme = ROOT + "knowledge/summaries/README.md"
  listed = markdown_links(read_text(readme))
             .map { |target| (readme.dirname + target.split("#", 2).first).cleanpath }
             .select { |path| path.to_s.end_with?(".md") && path.basename.to_s != "README.md" }
  summary_files.each do |path|
    next if listed.include?(path.cleanpath)

    add_error("knowledge/summaries/README.md: summary not listed #{rel(path)}")
  end

  check_markdown_links("knowledge/summaries/README.md")
  COUNTS[:summaries] = summary_files.length
end

def validate_graphs
  graph_files = Dir.glob((ROOT + "knowledge/graphs/*.yaml").to_s).sort
  graph_files.each do |graph_file|
    relative = rel(graph_file)
    data = require_mapping(yaml_file(relative), relative)
    check_fields(data, %w[id source status summary edges metadata], relative)
    check_repo_path(data["source"], "#{relative} source") if data["source"]

    edges = require_array(data["edges"], "#{relative} edges")
    edges.each_with_index do |edge, index|
      context = "#{relative} edges[#{index}]"
      edge = require_mapping(edge, context)
      check_fields(edge, %w[type target reason validation], context)
      add_error("#{context}: unknown edge type #{edge["type"]}") if edge["type"] && !EDGE_TYPES.include?(edge["type"])
      check_repo_path(edge["target"], "#{context} target") if edge["target"]
    end

    metadata = require_mapping(data["metadata"], "#{relative} metadata")
    check_fields(metadata, %w[priority confidence compatibility_state last_checked], "#{relative} metadata")
  end

  check_markdown_links("knowledge/graphs/README.md")
  COUNTS[:graphs] = graph_files.length
end

validate_registry
validate_refresh_policy
validate_summaries
validate_graphs
check_markdown_links("knowledge/runtime/README.md")
check_markdown_links("knowledge/runtime/runtime-report.md") if (ROOT + "knowledge/runtime/runtime-report.md").exist?
check_markdown_links("knowledge/runtime/model-context-report.md") if (ROOT + "knowledge/runtime/model-context-report.md").exist?
check_markdown_links("knowledge/runtime/model-checklists.md") if (ROOT + "knowledge/runtime/model-checklists.md").exist?
check_markdown_links("runtime/routing/README.md")
check_markdown_links("knowledge/indexes/README.md")

if ERRORS.empty?
  puts "Knowledge runtime validation OK"
  puts "registry_records=#{COUNTS[:registry_records]}"
  puts "summaries=#{COUNTS[:summaries]}"
  puts "graphs=#{COUNTS[:graphs]}"
else
  warn "Knowledge runtime validation failed:"
  ERRORS.each { |error| warn "- #{error}" }
  exit 1
end
