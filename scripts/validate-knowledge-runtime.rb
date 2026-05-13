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

SKIP_SUBDIRS = %w[feedback_history].freeze

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

def markdown_links_from_text(text)
  # Extract both standard markdown links [...](...) and table cell links [...](...)
  links = markdown_links(text)
  # Also extract links from table cells: | [text](path) |
  text.scan(/\|?\s*\[[^\]]+\]\(([^)]+)\)\s*\|?/).flatten.each { |l| links << l }
  links.uniq
end

def files_listed_in_readme(readme_path)
  return [] unless readme_path.exist?

  readme_text = read_text(readme_path)
  links = markdown_links_from_text(readme_text)
  links.map { |t| t.split("#", 2).first }
       .reject(&:empty?)
       .map { |t| (readme_path.dirname + t).cleanpath }
end

def validate_directory_structure
  # Check that every subdirectory under intelligence/engineering/<domain>/ has a README.md
  domains_dir = ROOT + "intelligence/engineering"
  if domains_dir.exist?
    domains_dir.each_child do |domain_dir|
      next unless domain_dir.directory?
      next if domain_dir.basename.to_s.start_with?(".")

      domain_rel = rel(domain_dir)
      readme = domain_dir + "README.md"
      add_error("#{domain_rel}: missing README.md") unless readme.exist?

      # Check sub-categories (heuristics, anti-patterns, failure, signals) have README.md
      domain_dir.each_child do |sub_dir|
        next unless sub_dir.directory?
        next if sub_dir.basename.to_s.start_with?(".")
        next if SKIP_SUBDIRS.include?(sub_dir.basename.to_s)

        sub_rel = rel(sub_dir)
        sub_readme = sub_dir + "README.md"
        add_error("#{sub_rel}: missing README.md") unless sub_readme.exist?

        # Check that every .md file in sub-category is listed in its README.md
        if sub_readme.exist?
          listed = files_listed_in_readme(sub_readme)

          sub_dir.each_child do |file|
            next unless file.extname == ".md"
            next if file.basename.to_s == "README.md"

            add_error("#{rel(file)}: not listed in #{sub_rel}/README.md") unless listed.include?(file.cleanpath)
          end
        end
      end
    end
  end

  # Check that every subdirectory under analysis/<domain>/workflows/ has its .md files listed in README.md
  analysis_dir = ROOT + "analysis"
  if analysis_dir.exist?
    analysis_dir.each_child do |domain_dir|
      next unless domain_dir.directory?
      next if domain_dir.basename.to_s.start_with?(".")

      workflows_dir = domain_dir + "workflows"
      next unless workflows_dir.exist?

      workflows_rel = rel(workflows_dir)
      workflows_readme = workflows_dir + "README.md"
      add_error("#{workflows_rel}: missing README.md") unless workflows_readme.exist?

      if workflows_readme.exist?
        listed = files_listed_in_readme(workflows_readme)

        workflows_dir.each_child do |file|
          next unless file.extname == ".md"
          next if file.basename.to_s == "README.md"

          add_error("#{rel(file)}: not listed in #{workflows_rel}/README.md") unless listed.include?(file.cleanpath)
        end
      end
    end
  end

  # Check that every subdirectory under workflow/<domain>/ has its .md files listed in README.md
  workflow_dir = ROOT + "workflow"
  if workflow_dir.exist?
    workflow_dir.each_child do |domain_dir|
      next unless domain_dir.directory?
      next if domain_dir.basename.to_s.start_with?(".")

      domain_rel = rel(domain_dir)
      domain_readme = domain_dir + "README.md"
      add_error("#{domain_rel}: missing README.md") unless domain_readme.exist?

      if domain_readme.exist?
        listed = files_listed_in_readme(domain_readme)

        domain_dir.each_child do |file|
          next unless file.extname == ".md"
          next if file.basename.to_s == "README.md"

          add_error("#{rel(file)}: not listed in #{domain_rel}/README.md") unless listed.include?(file.cleanpath)
        end
      end
    end
  end

  # Count engineering subdomains + top-level intelligence subdirectories (ide, business, travel, etc.)
  eng_domains = Dir.glob((ROOT + "intelligence/engineering/*").to_s).count { |p| File.directory?(p) }
  top_domains = Dir.glob((ROOT + "intelligence/*").to_s).count { |p| File.directory?(p) && File.basename(p) != "engineering" }
  COUNTS[:intelligence_domains] = eng_domains + top_domains
  COUNTS[:analysis_domains] = Dir.glob((ROOT + "analysis/*").to_s).count { |p| File.directory?(p) }
  COUNTS[:workflow_domains] = Dir.glob((ROOT + "workflow/*").to_s).count { |p| File.directory?(p) }
end

def validate_no_outdated_active_entrypoint
  # Check that no files under workflow/, analysis/, runtime/onboarding/ still reference
  # skills/ as "active entrypoint". These should have been migrated to "遷移狀態" headers.
  scan_dirs = %w[workflow analysis runtime/onboarding]
  pattern = /仍為 active( skill)? (entrypoint|source of truth)/i

  scan_dirs.each do |dir|
    base = ROOT + dir
    next unless base.exist?

    Dir.glob((base + "**/*.md").to_s).sort.each do |file|
      next if file.include?("feedback_history/")
      next if file.include?("/archived/")

      text = read_text(file)
      if text.match?(pattern)
        add_error("#{rel(file)}: contains outdated '仍為 active entrypoint' reference (should be '遷移狀態')")
      end
    end
  end
end

def validate_intelligence_ide_knowledge
  # Check that intelligence/ide/ exists and has proper structure.
  # This knowledge was promoted from ai-tools/ide/ to intelligence/ide/.
  ide_dir = ROOT + "intelligence/ide"
  unless ide_dir.exist?
    add_error("intelligence/ide/ does not exist (expected after promotion from ai-tools/ide/)")
    return
  end

  # Must have README.md
  readme = ide_dir + "README.md"
  unless readme.exist?
    add_error("intelligence/ide/README.md is missing")
    return
  end

  # README.md must list all .md files in the directory (excluding itself)
  listed = files_listed_in_readme(readme).map(&:to_s)
  actual = ide_dir.each_child
                  .select { |f| f.extname == ".md" && f.basename.to_s != "README.md" }
                  .map(&:to_s)

  actual.each do |f|
    unless listed.include?(f)
      basename = File.basename(f)
      add_error("intelligence/ide/README.md does not list #{basename}")
    end
  end

  # The vscode-extension-global-state.md should NOT exist in ai-tools/ide/ anymore
  old_path = ROOT + "ai-tools/ide/vscode-extension-global-state.md"
  if old_path.exist?
    add_error("ai-tools/ide/vscode-extension-global-state.md still exists (should have been promoted to intelligence/ide/)")
  end
end

def validate_language_consistency
  # Check for Author Habit Drift (Type B): Chinese documents with English table headers.
  # This is a simple heuristic — if a .md file is mostly Chinese but has common English
  # table header patterns like "| Change |" or "| Description |", flag it.
  # Scans files under shared-rules/failure-patterns/ and intelligence/ide/ as these
  # are most prone to this drift.
  scan_dirs = %w[shared-rules/failure-patterns intelligence/ide]

  # Common English table headers that should be in Chinese if the document is Chinese
  english_headers = /\|\s*(Change|Description|Status|Notes|Example|Type|Name|Value|Key|Field|Method|Source|Target|Action|Result|Category|Priority|Risk|Impact|Scope|Trigger|Failure Mode|Root Cause|Prevention Gate|Validation Method)\s*\|/

  scan_dirs.each do |dir|
    base = ROOT + dir
    next unless base.exist?

    Dir.glob((base + "**/*.md").to_s).sort.each do |file|
      next if file.include?("feedback_history/")
      next if file.include?("/archived/")

      text = read_text(file)
      # Only check files that have Chinese content (to avoid false positives on English-only files)
      has_chinese = text.match?(/[\u4e00-\u9fff]/)
      next unless has_chinese

      # Skip content inside code blocks (```...```) — these are intentional examples
      lines = text.lines
      in_code_block = false
      lines.each_with_index do |line, idx|
        line_num = idx + 1
        if line.strip.start_with?("```")
          in_code_block = !in_code_block
          next
        end
        next if in_code_block

        if line.match?(english_headers) && line.include?("|")
          # Check if this is a table separator line (---|---|---)
          next if line.strip.match?(/^[\|\s\-:]+$/)
          add_error("#{rel(file)}:#{line_num}: possible Author Habit Drift — English table header '#{line.strip}' in Chinese document")
        end
      end
    end
  end
end

def validate_intelligence_classification_boundary
  # Check that every top-level subdirectory under intelligence/ is listed in
  # intelligence/README.md's structure diagram. This prevents "framework dependency
  # bias" — putting a directory under engineering/ when it should be a sibling.
  #
  # The structure diagram in intelligence/README.md defines the canonical classification.
  # Any new directory must be added there first, with a clear description of its boundary.
  readme_path = ROOT + "intelligence/README.md"
  unless readme_path.exist?
    add_error("intelligence/README.md is missing (cannot validate classification boundary)")
    return
  end

  readme_text = read_text(readme_path)

  # Parse the structure diagram: find lines between ```text and ``` that contain "  dirname/"
  in_diagram = false
  listed_dirs = []
  readme_text.each_line do |line|
    if line.strip == "```text"
      in_diagram = true
      next
    end
    if line.strip == "```" && in_diagram
      in_diagram = false
      next
    end
    next unless in_diagram

    # Match lines like "  engineering/" or "  ide/" (indented directory names)
    if line.match?(/^\s{2}[a-z][a-z0-9_-]+\/\s*(#.*)?$/)
      dir_name = line.strip.split("/").first
      listed_dirs << dir_name
    end
  end

  # Get actual top-level directories under intelligence/
  actual_dirs = []
  intel_dir = ROOT + "intelligence"
  if intel_dir.exist?
    intel_dir.each_child do |child|
      next unless child.directory?
      next if child.basename.to_s.start_with?(".")
      actual_dirs << child.basename.to_s
    end
  end

  # Check for directories not listed in the structure diagram
  actual_dirs.each do |dir|
    unless listed_dirs.include?(dir)
      add_error("intelligence/#{dir}/ exists but is not listed in intelligence/README.md structure diagram (possible classification boundary violation)")
    end
  end

  # Check for directories listed in diagram but not existing (optional warning)
  listed_dirs.each do |dir|
    unless actual_dirs.include?(dir)
      add_error("intelligence/README.md lists '#{dir}/' in structure diagram but directory does not exist")
    end
  end
end

# ──────────────────────────────────────────────
# Failure-to-Validator Closure Test
# ──────────────────────────────────────────────
def validate_failure_pattern_validator_coverage
  # Check that every failure pattern in shared-rules/failure-patterns/ has at least
  # one entry in its "Linked Validation Scenarios" section. This prevents the
  # "failure-to-validator closure" pattern where an error is fixed but no reusable
  # test is added to prevent recurrence.
  #
  # Rationale: If a failure pattern's Linked Validation Scenarios is empty, it means
  # the agent who created the pattern didn't add a corresponding validator test.
  # This is itself a failure-to-validator closure.
  patterns_dir = ROOT + "shared-rules/failure-patterns"
  return unless patterns_dir.exist?

  patterns_dir.each_child do |file|
    next unless file.extname == ".md"
    next if file.basename.to_s == "README.md"

    text = read_text(file)
    lines = text.each_line.to_a

    # Find the "Linked Validation Scenarios" section
    in_section = false
    scenario_lines = []
    lines.each_with_index do |line, idx|
      if line.match?(/^## Linked Validation Scenarios/)
        in_section = true
        next
      end
      if in_section && line.start_with?("## ")
        break
      end
      if in_section
        scenario_lines << line
      end
    end

    # Check if there's at least one non-empty scenario entry
    has_scenario = scenario_lines.any? { |l| l.match?(/^- `[^`]+`/) }

    unless has_scenario
      add_error("failure pattern #{file.basename} has empty Linked Validation Scenarios — add at least one validator test reference (failure-to-validator closure)")
    end
  end
end

validate_registry
validate_refresh_policy
validate_summaries
validate_graphs
validate_directory_structure
validate_intelligence_classification_boundary
validate_failure_pattern_validator_coverage
validate_no_outdated_active_entrypoint
validate_intelligence_ide_knowledge
validate_language_consistency
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
  puts "intelligence_domains=#{COUNTS[:intelligence_domains]}"
  puts "analysis_domains=#{COUNTS[:analysis_domains]}"
  puts "workflow_domains=#{COUNTS[:workflow_domains]}"
else
  warn "Knowledge runtime validation failed:"
  ERRORS.each { |error| warn "- #{error}" }
  exit 1
end
