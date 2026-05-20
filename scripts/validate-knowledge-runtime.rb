#!/usr/bin/env ruby
# encoding: UTF-8
# Validate Ai-skill knowledge runtime generated surfaces.

require "date"
require "pathname"
require "set"
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

  # Check that every top-level directory (excluding hidden dirs and files) has a README.md
  top_level_dirs = ROOT.each_child
                       .select { |d| d.directory? && !d.basename.to_s.start_with?(".") }
                       .sort
  top_level_dirs.each do |dir|
    dir_rel = rel(dir)
    readme = dir + "README.md"
    add_error("#{dir_rel}: missing README.md — every top-level directory must have an entry point") unless readme.exist?
  end

  # Count engineering subdomains + top-level intelligence subdirectories (ide, business, travel, etc.)
  eng_domains = Dir.glob((ROOT + "intelligence/engineering/*").to_s).count { |p| File.directory?(p) }
  top_domains = Dir.glob((ROOT + "intelligence/*").to_s).count { |p| File.directory?(p) && File.basename(p) != "engineering" }
  COUNTS[:intelligence_domains] = eng_domains + top_domains
  COUNTS[:analysis_domains] = Dir.glob((ROOT + "analysis/*").to_s).count { |p| File.directory?(p) }
  COUNTS[:workflow_domains] = Dir.glob((ROOT + "workflow/*").to_s).count { |p| File.directory?(p) }
end

def validate_directory_naming
  # ── Allowlist ──────────────────────────────────────────────────────────────
  # Pre-existing directories with clear boundary documentation in their README.md.
  # These are intentionally named to match a top-level category but have distinct
  # content scope (e.g., engineering anti-patterns vs agent operation anti-patterns).
  # New directories should NOT be added here — they should use distinct names.
  allowed_same_name_dirs = %w[
    anti-patterns
    architecture
  ].to_set

  # Check 1: Same-name cross-layer conflict detection
  # Detect directories under intelligence/engineering/ that share a name with
  # top-level directories (e.g., intelligence/engineering/anti-patterns/ vs anti-patterns/).
  # This is a naming governance violation — the name should reflect content essence,
  # not just reuse an existing category name.
  top_level_dirs = ROOT.each_child
                       .select { |d| d.directory? && !d.basename.to_s.start_with?(".") }
                       .map { |d| d.basename.to_s }
                       .to_set

  eng_dir = ROOT + "intelligence/engineering"
  if eng_dir.exist?
    eng_dir.each_child do |domain_dir|
      next unless domain_dir.directory?
      dirname = domain_dir.basename.to_s
      next if dirname.start_with?(".")
      next if allowed_same_name_dirs.include?(dirname)

      if top_level_dirs.include?(dirname)
        add_error("#{rel(domain_dir)}: same-name cross-layer conflict — '#{dirname}' also exists as top-level directory '#{dirname}/'. " \
                  "Rename to reflect content essence (see governance/lifecycle/directory-structure-governance.md Step 1)")
      end
    end
  end

  # Check 2: Inertial naming detection
  # Detect directory names that appear to be shortened versions of old skill names
  # rather than reflecting content essence. This is a heuristic check based on
  # known old skill names.
  old_skill_names = %w[apk-analysis travel-planning repo-analysis]
  if eng_dir.exist?
    eng_dir.each_child do |domain_dir|
      next unless domain_dir.directory?
      dirname = domain_dir.basename.to_s

      # Check if dirname is a prefix/shortened form of an old skill name
      old_skill_names.each do |skill|
        # e.g., "analysis" is a prefix of "apk-analysis" or "repo-analysis"
        if skill.start_with?(dirname) && skill != dirname
          add_error("#{rel(domain_dir)}: possible inertial naming — '#{dirname}' is a prefix of old skill '#{skill}'. " \
                    "Rename to reflect content essence, not source skill name (see governance/lifecycle/directory-structure-governance.md Step 3)")
        end
      end
    end
  end

  # Check 3: Path depth warning
  # Detect directories deeper than 4 levels from root (excluding root itself).
  # Deep nesting makes navigation harder and increases cognitive load.
  max_depth = 4
  Dir.glob((ROOT + "**/*/").to_s).sort.each do |dir_path|
    next if dir_path.include?(".git/")
    next if dir_path.include?("node_modules/")

    dir = Pathname.new(dir_path)
    relative = dir.relative_path_from(ROOT).to_s
    depth = relative.split.length

    if depth > max_depth
      add_error("#{relative}: path depth #{depth} exceeds maximum #{max_depth}. " \
                "Consider flattening the directory structure (see governance/lifecycle/directory-structure-governance.md Step 4)")
    end
  end
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
  # Check for Author Habit Drift (Type B): Chinese documents with English headings
  # or table headers. This detects when an agent writes Chinese content but uses
  # English section titles (e.g., "## Token Impact") or English table headers.
  #
  # Scans intelligence/, workflow/, analysis/, enforcement/ — any directory where
  # Chinese documents may have English headings.
  scan_dirs = %w[intelligence workflow analysis enforcement]

  # English section headings that commonly appear in Chinese documents (author habit drift)
  # These are heading patterns like "## Token Impact" or "### Risk Assessment"
  english_headings = /\A#{'#{1,6}'}\s+(Token Impact|Token 影響|Risk|Scope|Trigger|Failure Mode|Root Cause|Prevention Gate|Validation Method|Decision Rule|Preferred Pattern|Validation Signal|Boundaries|Intelligence Status|Examples Of|Example|Change|Description|Status|Notes|Type|Name|Value|Key|Field|Method|Source|Target|Action|Result|Category|Priority|Impact|Summary|Overview|Background|Purpose|Goal|Prerequisites|Steps|Output|Next Steps|Related|See Also|References)\s*\z/

  # English table headers that should be in Chinese if the document is Chinese
  english_table_headers = /\|\s*(Change|Description|Status|Notes|Example|Type|Name|Value|Key|Field|Method|Source|Target|Action|Result|Category|Priority|Risk|Impact|Scope|Trigger|Failure Mode|Root Cause|Prevention Gate|Validation Method)\s*\|/

  # ── Allowlist ──────────────────────────────────────────────────────────────
  # These are intentional structural conventions, NOT author habit drift.
  # Failure patterns use standard English template headings across all patterns.
  # Intelligence documents use "## Token Impact" as a standard section.
  # These should NOT be flagged.
  allowed_english_headings = [
    # Failure pattern template headings (intentional structural convention)
    "Trigger", "Failure Mode", "Risk", "Root Cause",
    "Required Agent Action", "Prevention Gate", "Validation Method",
    "Linked Rules", "Linked Failure Patterns",
    "Linked Feedback Lessons", "Linked Validation Scenarios",
    # README standard structural headings
    "Scope",
    # Intelligence standard sections
    "Token Impact", "Token 影響",
    # Intelligence document template (highest-leverage-analysis-path)
    "Intelligence Status", "Decision Rule", "Preferred Pattern",
    "Validation Signal", "Boundaries",
  ].freeze

  # ── Table header allowlist ─────────────────────────────────────────────────
  # These are intentional conventions where English field names serve as
  # identifiers in Chinese documents (e.g., goal ledger field names like
  # "Priority", "Status", "Source", "Scope" used as row labels).
  # The table header row itself should be Chinese, but content rows with
  # English field identifiers are intentional.
  allowed_table_field_names = %w[
    Priority Status Source Scope Goal Trigger
    Field Value Atom
    Parallelization mode Owner Owner/lock decision
    Subgoals Planning/todo links Open work/decisions
    Dependencies Next action Completion criteria
    Target skill Expected input/output Ownership boundary
    Sanitization boundary Linked updates
    Required set Read Not applicable Deferred/blocked Validation
  ].freeze

  # Check if a heading is in the allowlist (case-insensitive match on the heading text after ##)
  def allowed_heading?(heading_text, allowed_list)
    # Extract the heading text after the ## markers
    clean = heading_text.sub(/\A#+\s+/, "").strip
    allowed_list.any? { |a| clean.casecmp?(a) }
  end

  scan_dirs.each do |dir|
    base = ROOT + dir
    next unless base.exist?

    Dir.glob((base + "**/*.md").to_s).sort.each do |file|
      next if file.include?("feedback_history/")
      next if file.include?("/archived/")
      next if file.include?("node_modules/")

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

        # Check for English section headings (## Token Impact)
        stripped = line.strip
        if stripped.match?(english_headings)
          # Skip if this heading is in the allowlist (intentional structural convention)
          next if allowed_heading?(stripped, allowed_english_headings)
          add_error("#{rel(file)}:#{line_num}: possible Author Habit Drift — English heading '#{stripped}' in Chinese document")
        end

        # Check for English table headers — only flag the actual header row
        # (the first row of a table, before the separator line).
        # Content rows with English field names (e.g., "| Priority | P0 |")
        # are intentional identifiers and should NOT be flagged.
        if line.match?(english_table_headers) && line.include?("|")
          # Check if this is a table separator line (---|---|---)
          next if line.strip.match?(/^[\|\s\-:]+$/)
          # Check if the NEXT line is a table separator (---|---|---)
          # If so, this is the header row — flag it if it has English headers
          next_line = lines[idx + 1]
          if next_line && next_line.strip.match?(/^[\|\s\-:]+$/)
            # This is a table header row — check if all cells are English
            # (mixed Chinese/English headers like "| Scope | 放置位置 |" are borderline)
            cells = line.split("|").map(&:strip).reject(&:empty?)
            english_cells = cells.count { |c| c.match?(/\A[a-zA-Z\s\/]+\z/) }
            total_cells = cells.length
            if english_cells == total_cells
              add_error("#{rel(file)}:#{line_num}: possible Author Habit Drift — English table header '#{line.strip}' in Chinese document")
            end
          end
          # Content rows (not header rows) are NOT flagged — English field names
          # like "Priority", "Status", "Source" are intentional identifiers
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
  # Check that every failure pattern in enforcement/failure-patterns/ has at least
  # one entry in its "Linked Validation Scenarios" section. This prevents the
  # "failure-to-validator closure" pattern where an error is fixed but no reusable
  # test is added to prevent recurrence.
  #
  # Rationale: If a failure pattern's Linked Validation Scenarios is empty, it means
  # the agent who created the pattern didn't add a corresponding validator test.
  # This is itself a failure-to-validator closure.
  patterns_dir = ROOT + "enforcement/failure-patterns"
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

# ──────────────────────────────────────────────
# Runtime Recovery Scenario Test
# ──────────────────────────────────────────────
def validate_runtime_recovery_scenarios
  scenario_files = Dir.glob((ROOT + "validation/scenarios/failure-derived/runtime-recovery-*.yaml").to_s).sort
  required_then_fields = %w[expected_route required_reload_set forbidden_routes expected_final_route validation].freeze

  scenario_files.each do |path|
    relative = rel(path)
    data = require_mapping(yaml_file(relative), relative)
    next if data.empty?

    if data["type"] != "failure-recovery"
      add_error("#{relative}: runtime recovery scenario must use type=failure-recovery")
    end

    then_block = require_mapping(data["then"], "#{relative}.then")
    check_fields(then_block, required_then_fields, "#{relative}.then")

    reload_set = require_array(then_block["required_reload_set"], "#{relative}.then.required_reload_set")
    add_error("#{relative}: required_reload_set must not be empty") if reload_set.empty?
    reload_set.each { |item| check_repo_path(item, "#{relative}.then.required_reload_set") if item.is_a?(String) }

    forbidden_routes = require_array(then_block["forbidden_routes"], "#{relative}.then.forbidden_routes")
    add_error("#{relative}: forbidden_routes must not be empty") if forbidden_routes.empty?

    final_route = then_block["expected_final_route"].to_s
    unless final_route.include?("recovery")
      add_error("#{relative}: expected_final_route must point to a recovery route")
    end

    expected_route = require_mapping(then_block["expected_route"], "#{relative}.then.expected_route")
    steps = require_array(expected_route["steps"], "#{relative}.then.expected_route.steps")
    unless steps.any? { |step| step.to_s.include?("reload") || step.to_s.include?("rebuild") }
      add_error("#{relative}: expected_route.steps must include reload/rebuild behavior")
    end

    COUNTS[:runtime_recovery_scenarios] += 1
  end
end

# ──────────────────────────────────────────────
# Intelligence Entry/Solution Cross-Reference Test
# ──────────────────────────────────────────────
def validate_intelligence_entry_solution_crossref
  # Check that intelligence/<domain>/ subdirectories follow the Entry/Solution layering
  # rules defined in enforcement/content-layering.md:
  #
  # 1. If a file in failure/ says "解法見 heuristics/X.md", then heuristics/X.md must exist
  # 2. If a file in heuristics/ has a corresponding file in failure/ with the same basename,
  #    the failure file should reference it (cross-reference completeness)
  #
  # Entry directories: failure, signals, anti-patterns
  # Solution directories: heuristics, patterns
  entry_dirs = %w[failure signals anti-patterns].freeze
  solution_dirs = %w[heuristics patterns].freeze

  intel_dir = ROOT + "intelligence"
  return unless intel_dir.exist?

  intel_dir.each_child do |domain_dir|
    next unless domain_dir.directory?
    next if domain_dir.basename.to_s.start_with?(".")

    # Check each entry/solution subdirectory pair
    entry_dirs.each do |entry_name|
      entry_dir = domain_dir + entry_name
      next unless entry_dir.exist? && entry_dir.directory?

      entry_dir.each_child do |entry_file|
        next unless entry_file.extname == ".md"
        next if entry_file.basename.to_s == "README.md"

        text = read_text(entry_file)
        # Look for "解法見" or "→ see" patterns pointing to solution files
        text.each_line do |line|
          solution_dirs.each do |sol_name|
            # Match patterns like: "解法見 heuristics/X.md" or "→ see heuristics/X.md"
            if line.match?(/#{Regexp.escape(sol_name)}\/([a-z][a-z0-9_-]+)\.md/i)
              target_basename = $1
              target_path = domain_dir + sol_name + "#{target_basename}.md"
              unless target_path.exist?
                add_error("#{rel(entry_file)}: references #{sol_name}/#{target_basename}.md but file does not exist (broken entry/solution cross-reference)")
              end
            end
          end
        end
      end
    end

    # Check solution files have corresponding entry files (optional warning)
    solution_dirs.each do |sol_name|
      sol_dir = domain_dir + sol_name
      next unless sol_dir.exist? && sol_dir.directory?

      sol_dir.each_child do |sol_file|
        next unless sol_file.extname == ".md"
        next if sol_file.basename.to_s == "README.md"

        basename = sol_file.basename.to_s
        # Check if any entry directory has a file with the same basename
        has_entry = entry_dirs.any? do |entry_name|
          entry_dir = domain_dir + entry_name
          (entry_dir + basename).exist?
        end

        # If no entry file exists, that's fine — not all solutions need entries.
        # But if an entry file exists, the solution file should reference it.
        if has_entry
          text = read_text(sol_file)
          # Check if solution file references back to entry (top 10 lines)
          first_lines = text.each_line.first(10).join
          has_backref = entry_dirs.any? { |en| first_lines.match?(/#{Regexp.escape(en)}\//) }
          unless has_backref
            add_error("#{rel(sol_file)}: has corresponding entry file in failure/ but does not reference it in top 10 lines (missing back-reference per content-layering.md rule 5)")
          end
        end
      end
    end
  end
end

def validate_tool_config_standards
  # Validate that all agent tool config files in ai-tools/agent/ follow the
  # architecture defined in agent-onboarding.md:
  #
  # Required sections (from agent-onboarding.md 必要設定 table):
  # 1. 自動載入入口 — references CORE_BOOTSTRAP.md for startup flow
  # 2. 語言偏好 — soft language preference + 語言一致性強制規則
  # 3. 對話目標閉環 — goal ledger integration
  # 4. 全域設定 vs 專案設定 — global vs project config guide
  #
  # Each tool may have its own customization (hooks, modes, etc.),
  # but all must follow this architecture.
  agent_dir = ROOT + "ai-tools/agent"
  return unless agent_dir.exist? && agent_dir.directory?

  # Architecture requirements from agent-onboarding.md
  required_architecture = [
    {
      name: "自動載入入口",
      patterns: [/CORE_BOOTSTRAP\.md/],
      hint: "必須引用 CORE_BOOTSTRAP.md 作為啟動流程入口"
    },
    {
      name: "語言偏好設定（含語言一致性強制規則）",
      patterns: [
        /語言偏好設定/,
        /語言一致性強制規則/,
        /所有輸出.*attempt_completion.*commit message.*都必須與使用者當前語言一致/
      ],
      hint: "必須包含軟性語言偏好 + 語言一致性強制規則（所有輸出包含 attempt_completion、表格、commit message 都必須與使用者語言一致）"
    },
    {
      name: "對話目標閉環",
      patterns: [/conversation-goal-ledger\.md/, /goal ledger/, /\.agent-goals/],
      hint: "必須引用 conversation-goal-ledger.md 或實作 goal ledger 整合"
    },
    {
      name: "全域設定 vs 專案設定",
      patterns: [
        /##\s+全域設定\s*vs\s+專案設定/,
        /設定層級說明/,
        /層級 A[：:]/,
        /層級 B[：:]/,
        /建議策略/
      ],
      hint: "必須包含全域 vs 專案設定的層級說明（層級 A/層級 B）與建議策略"
    },
  ]

  agent_dir.each_child do |file|
    next unless file.extname == ".md"
    next if file.basename.to_s == "README.md"

    text = read_text(file)
    basename = file.basename.to_s

    required_architecture.each do |arch|
      matched = arch[:patterns].any? { |pattern| text.match?(pattern) }
      unless matched
        add_error("#{rel(file)}: 缺少必要架構項目「#{arch[:name]}」— #{arch[:hint]}（依據 agent-onboarding.md 必要設定表）")
      end
    end

    COUNTS[:tool_config_files] += 1
  end
end

validate_registry
validate_refresh_policy
validate_summaries
validate_graphs
validate_directory_structure
validate_directory_naming
validate_intelligence_classification_boundary
validate_intelligence_entry_solution_crossref
validate_failure_pattern_validator_coverage
validate_runtime_recovery_scenarios
validate_no_outdated_active_entrypoint
validate_intelligence_ide_knowledge
validate_language_consistency
validate_tool_config_standards
check_markdown_links("knowledge/runtime/README.md")
check_markdown_links("knowledge/runtime/runtime-report.md") if (ROOT + "knowledge/runtime/runtime-report.md").exist?
check_markdown_links("knowledge/runtime/model-context-report.md") if (ROOT + "knowledge/runtime/model-context-report.md").exist?
check_markdown_links("knowledge/runtime/model-checklists.md") if (ROOT + "knowledge/runtime/model-checklists.md").exist?
check_markdown_links("knowledge/indexes/README.md")

if ERRORS.empty?
  puts "Knowledge runtime validation OK"
  puts "registry_records=#{COUNTS[:registry_records]}"
  puts "summaries=#{COUNTS[:summaries]}"
  puts "graphs=#{COUNTS[:graphs]}"
  puts "intelligence_domains=#{COUNTS[:intelligence_domains]}"
  puts "analysis_domains=#{COUNTS[:analysis_domains]}"
  puts "workflow_domains=#{COUNTS[:workflow_domains]}"
  puts "tool_config_files=#{COUNTS[:tool_config_files]}"
  puts "runtime_recovery_scenarios=#{COUNTS[:runtime_recovery_scenarios]}"
else
  warn "Knowledge runtime validation failed:"
  ERRORS.each { |error| warn "- #{error}" }
  exit 1
end
