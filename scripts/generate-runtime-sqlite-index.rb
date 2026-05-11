#!/usr/bin/env ruby
# encoding: UTF-8
# Generate a local SQLite / FTS runtime lookup cache from canonical sources.

require "date"
require "digest"
require "fileutils"
require "open3"
require "pathname"
require "shellwords"
require "yaml"

ROOT = Pathname.new(__dir__).parent.realpath
DEFAULT_DB = ROOT + "knowledge/runtime/sqlite/runtime-index.sqlite"

def utf8(value)
  text = value.to_s
  text = text.dup.force_encoding("UTF-8") if text.encoding == Encoding::ASCII_8BIT || text.encoding == Encoding::US_ASCII
  text.encode("UTF-8", invalid: :replace, undef: :replace)
end

def read_text(path)
  File.read(path.to_s, encoding: "UTF-8")
end

def yaml_file(relative_path)
  YAML.safe_load(read_text(ROOT + relative_path), permitted_classes: [Date], aliases: false)
end

def sql(value)
  return "NULL" if value.nil?

  text = utf8(value)
  "'#{text.gsub("'", "''")}'"
end

def rel(path)
  utf8(Pathname.new(path).relative_path_from(ROOT).to_s)
end

def layer_for(path)
  path.to_s.split("/", 2).first
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

def strip_markup(value)
  utf8(value).gsub(/\[([^\]]+)\]\([^)]+\)/, "\\1").gsub("`", "")
end

def links_from_cell(cell)
  cell.to_s.scan(/\[[^\]]+\]\(([^)]+)\)/).flatten.map { |target| target.split("#", 2).first }.reject(&:empty?)
end

def repo_relative_link(base_path, target)
  return nil if target.nil? || target.empty?
  return target if target.start_with?("http://", "https://")

  resolved = (base_path.dirname + target).cleanpath
  rel(resolved)
end

def title_from_markdown(text)
  text.each_line do |line|
    return line.sub(/^#+\s*/, "").strip if line.start_with?("#")
  end
  ""
end

def first_heading_after(text, marker)
  lines = text.lines
  index = lines.find_index { |line| line.strip == marker }
  return "" unless index

  lines[(index + 1)..].to_a.each do |line|
    value = line.strip
    next if value.empty?
    return value unless value.start_with?("#")
  end
  ""
end

def feedback_records
  Dir.glob((ROOT + "skills/*/feedback_history/**/*.md").to_s).sort.map do |path|
    pathname = Pathname.new(path)
    next if pathname.basename.to_s == "README.md"

    text = read_text(pathname)
    relative = rel(pathname)
    skill = relative.split("/")[1]
    category = relative.split("/")[3..-2]&.join("/")
    title = text.lines.find { |line| line.start_with?("### ") }.to_s.sub(/^###\s*/, "").strip
    status = text[/^Status:\s*([^\n]+)/, 1].to_s.strip
    summary = first_heading_after(text, "#### One-line Summary")
    {
      id: "feedback.#{skill}.#{File.basename(relative, ".md")}",
      source_path: relative,
      layer: "skills",
      type: "feedback-pattern",
      status: status.empty? ? "candidate" : status,
      priority: "P2",
      confidence: status == "promoted" || status == "validated" ? "high" : "medium",
      context_cost: "medium",
      tags: [skill, "feedback", category].compact.join(","),
      domains: skill,
      title: title.empty? ? File.basename(relative, ".md") : title,
      summary: summary,
      when_to_read: "Feedback lesson lookup for #{skill}.",
      validation_signal: "Open canonical feedback lesson at #{relative}."
    }
  end.compact
end

def summary_records
  Dir.glob((ROOT + "knowledge/summaries/*.md").to_s).sort.map do |path|
    pathname = Pathname.new(path)
    next if pathname.basename.to_s == "README.md"

    fields = parse_summary_table(pathname)
    source_path = repo_relative_link(pathname, links_from_cell(fields["Source path"]).first)
    {
      id: strip_markup(fields["Atom ID"]),
      source_path: source_path || rel(pathname),
      layer: layer_for(source_path || rel(pathname)),
      type: "summary",
      status: strip_markup(fields["Lifecycle"]),
      priority: "P2",
      confidence: strip_markup(fields["Lifecycle"]) == "validated" ? "high" : "medium",
      context_cost: "low",
      tags: "summary,atom",
      domains: "",
      title: strip_markup(fields["Atom ID"]),
      summary: fields["Summary"].to_s,
      when_to_read: fields["When to read"].to_s,
      validation_signal: fields["Validation signal"].to_s
    }
  end.compact
end

def route_records
  records = Array(yaml_file("knowledge/runtime/routing-registry.yaml")["records"])
  records.map do |record|
    metadata = record["metadata"] || {}
    model = record["model"] || {}
    {
      id: record["id"],
      source_path: record["primary_source"],
      layer: layer_for(record["primary_source"]),
      type: "route",
      status: metadata["compatibility_state"],
      priority: metadata["priority"],
      confidence: metadata["confidence"],
      context_cost: metadata["context_cost"],
      tags: ["route", model["profile"], model["compression_level"]].compact.join(","),
      domains: "",
      title: record["task_intent"],
      summary: record["ranking_reason"],
      when_to_read: record["task_intent"],
      validation_signal: record["validation_signal"]
    }
  end
end

def graph_edges
  Dir.glob((ROOT + "knowledge/graphs/*.yaml").to_s).sort.flat_map do |path|
    relative = rel(Pathname.new(path))
    data = yaml_file(relative)
    Array(data["edges"]).map do |edge|
      {
        graph_id: data["id"],
        source_path: data["source"],
        edge_type: edge["type"],
        target_path: edge["target"],
        reason: edge["reason"],
        validation: edge["validation"]
      }
    end
  end
end

def source_rows(records)
  records.map { |record| record[:source_path] }.compact.uniq.sort.map do |source_path|
    path = ROOT + source_path
    next unless path.exist?

    text = path.file? ? read_text(path) : ""
    {
      source_path: source_path,
      layer: layer_for(source_path),
      title: path.file? ? title_from_markdown(text) : File.basename(source_path),
      checksum: path.file? ? Digest::SHA256.hexdigest(text) : "",
      bytes: path.file? ? text.bytesize : 0
    }
  end.compact
end

def build_sql(records)
  sources = source_rows(records)
  edges = graph_edges
  statements = []
  statements << "PRAGMA journal_mode=OFF;"
  statements << "PRAGMA synchronous=OFF;"
  statements << "CREATE TABLE atoms (id TEXT PRIMARY KEY, source_path TEXT, layer TEXT, type TEXT, status TEXT, priority TEXT, confidence TEXT, context_cost TEXT, tags TEXT, domains TEXT, title TEXT, summary TEXT, when_to_read TEXT, validation_signal TEXT);"
  statements << "CREATE TABLE sources (source_path TEXT PRIMARY KEY, layer TEXT, title TEXT, checksum TEXT, bytes INTEGER);"
  statements << "CREATE TABLE edges (graph_id TEXT, source_path TEXT, edge_type TEXT, target_path TEXT, reason TEXT, validation TEXT);"
  statements << "CREATE VIRTUAL TABLE fts USING fts5(id, source_path, title, summary, tags, when_to_read, validation_signal);"
  records.each do |record|
    statements << "INSERT INTO atoms VALUES (#{[
      record[:id], record[:source_path], record[:layer], record[:type], record[:status],
      record[:priority], record[:confidence], record[:context_cost], record[:tags],
      record[:domains], record[:title], record[:summary], record[:when_to_read],
      record[:validation_signal]
    ].map { |value| sql(value) }.join(", ")});"
    statements << "INSERT INTO fts VALUES (#{[
      record[:id], record[:source_path], record[:title], record[:summary], record[:tags],
      record[:when_to_read], record[:validation_signal]
    ].map { |value| sql(value) }.join(", ")});"
  end
  sources.each do |source|
    statements << "INSERT INTO sources VALUES (#{[
      source[:source_path], source[:layer], source[:title], source[:checksum], source[:bytes]
    ].map { |value| sql(value) }.join(", ")});"
  end
  edges.each do |edge|
    statements << "INSERT INTO edges VALUES (#{[
      edge[:graph_id], edge[:source_path], edge[:edge_type], edge[:target_path], edge[:reason], edge[:validation]
    ].map { |value| sql(value) }.join(", ")});"
  end
  statements.join("\n")
end

output = DEFAULT_DB
ARGV.each_with_index do |arg, index|
  output = ROOT + ARGV[index + 1] if arg == "--output" && ARGV[index + 1]
end

records = (summary_records + route_records + feedback_records)
FileUtils.mkdir_p(output.dirname)
FileUtils.rm_f(output)
stdout, stderr, status = Open3.capture3("sqlite3", output.to_s, stdin_data: build_sql(records))
unless status.success?
  warn stderr
  warn stdout
  exit status.exitstatus || 1
end

puts "Wrote #{output.relative_path_from(ROOT)}"
puts "atoms=#{records.length}"
puts "sources=#{source_rows(records).length}"
puts "edges=#{graph_edges.length}"
