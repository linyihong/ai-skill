#!/usr/bin/env ruby
# encoding: UTF-8
# Generate a deterministic report from knowledge runtime surfaces.

require "date"
require "pathname"
require "yaml"

ROOT = Pathname.new(__dir__).parent.realpath
OUTPUT_PATH = ROOT + "knowledge/runtime/runtime-report.md"

def read_text(path)
  File.read(path.to_s, encoding: "UTF-8")
end

def yaml_file(relative_path)
  YAML.safe_load(read_text(ROOT + relative_path), permitted_classes: [Date], aliases: false)
end

def markdown_links(text)
  text.scan(/\[[^\]]+\]\(([^)]+)\)/).flatten
end

def parse_summary_table(path)
  fields = {}
  read_text(path).each_line do |line|
    next unless line.start_with?("|")

    cells = line.strip.split("|").map(&:strip)
    next unless cells.length >= 3
    next if cells[1] == "欄位" || cells[1].match?(/\A-+\z/)

    fields[cells[1]] = cells[2]
  end
  fields
end

def summary_records
  Dir.glob((ROOT + "knowledge/summaries/*.md").to_s)
     .reject { |path| File.basename(path) == "README.md" }
     .sort
     .map do |path|
       fields = parse_summary_table(Pathname.new(path))
       {
         file: Pathname.new(path).relative_path_from(ROOT).to_s,
         atom_id: fields["Atom ID"].to_s.gsub("`", ""),
         lifecycle: fields["Lifecycle"].to_s.gsub("`", ""),
         summary: fields["Summary"].to_s
       }
     end
end

def graph_records
  Dir.glob((ROOT + "knowledge/graphs/*.yaml").to_s).sort.map do |path|
    relative = Pathname.new(path).relative_path_from(ROOT).to_s
    data = yaml_file(relative)
    {
      file: relative,
      id: data["id"],
      source: data["source"],
      status: data["status"],
      edge_count: Array(data["edges"]).length
    }
  end
end

def registry_records
  data = yaml_file("knowledge/runtime/routing-registry.yaml")
  Array(data["records"])
end

def refresh_policy
  yaml_file("knowledge/runtime/refresh-policy.yaml")
end

def md_escape(value)
  value.to_s.gsub("|", "\\|").gsub("\n", " ")
end

def build_report
  routes = registry_records
  summaries = summary_records
  graphs = graph_records
  policy = refresh_policy

  lines = []
  lines << "# Knowledge Runtime Report"
  lines << ""
  lines << "本檔由 `ruby scripts/generate-knowledge-runtime-report.rb --write` 產生，彙整 runtime registry、summaries、graphs 與 refresh policy 的目前狀態。"
  lines << ""
  lines << "## Source Surfaces"
  lines << ""
  lines << "| Surface | Path | Count / Status |"
  lines << "| --- | --- | --- |"
  lines << "| Routing registry | [`routing-registry.yaml`](routing-registry.yaml) | #{routes.length} records |"
  lines << "| Refresh policy | [`refresh-policy.yaml`](refresh-policy.yaml) | #{policy["status"] || "unknown"} |"
  lines << "| Model context report | [`model-context-report.md`](model-context-report.md) | generated view |"
  lines << "| SQLite runtime index | [`sqlite/`](sqlite/) | planned lookup cache |"
  lines << "| Summaries | [`../summaries/`](../summaries/) | #{summaries.length} files |"
  lines << "| Graph records | [`../graphs/`](../graphs/) | #{graphs.length} files |"
  lines << ""
  lines << "## Routing Records"
  lines << ""
  lines << "| ID | Primary source | Model | Compression | Validation signal |"
  lines << "| --- | --- | --- | --- | --- |"
  routes.each do |record|
    model = record["model"] || {}
    lines << "| `#{md_escape(record["id"])}` | `#{md_escape(record["primary_source"])}` | `#{md_escape(model["profile"])}` | `#{md_escape(model["compression_level"])}` | #{md_escape(record["validation_signal"])} |"
  end
  lines << ""
  lines << "## Summary Records"
  lines << ""
  lines << "| Atom ID | Lifecycle | File | Summary |"
  lines << "| --- | --- | --- | --- |"
  summaries.each do |summary|
    lines << "| `#{md_escape(summary[:atom_id])}` | `#{md_escape(summary[:lifecycle])}` | [`#{md_escape(File.basename(summary[:file]))}`](../summaries/#{File.basename(summary[:file])}) | #{md_escape(summary[:summary])} |"
  end
  lines << ""
  lines << "## Graph Records"
  lines << ""
  lines << "| ID | Source | Status | Edges | File |"
  lines << "| --- | --- | --- | --- | --- |"
  graphs.each do |graph|
    lines << "| `#{md_escape(graph[:id])}` | `#{md_escape(graph[:source])}` | `#{md_escape(graph[:status])}` | #{graph[:edge_count]} | [`#{md_escape(File.basename(graph[:file]))}`](../graphs/#{File.basename(graph[:file])}) |"
  end
  lines << ""
  lines << "## Refresh Decisions"
  lines << ""
  lines << "| Decision value | Meaning |"
  lines << "| --- | --- |"
  Array(policy["decision_values"]).each do |decision|
    lines << "| `#{md_escape(decision)}` | 由 `refresh-policy.yaml` 定義的 generated surface decision。 |"
  end
  lines << ""
  lines << "## Validation"
  lines << ""
  lines << "- 產生前應先執行 `ruby scripts/validate-knowledge-runtime.rb`。"
  lines << "- 產生後應執行 Markdown link check、lints、close-loop dry run、commit / push / readback。"
  lines << "- 本報告是 generated view，不取代 `routing-registry.yaml`、`refresh-policy.yaml`、summary 或 graph source files。"
  lines << ""
  lines.join("\n")
end

if ARGV.include?("--write")
  OUTPUT_PATH.write(build_report, mode: "w", encoding: "UTF-8")
  puts "Wrote #{OUTPUT_PATH.relative_path_from(ROOT)}"
else
  puts build_report
end
