#!/usr/bin/env ruby
# encoding: UTF-8
# Query knowledge graph records without loading every graph file into context.

require "date"
require "pathname"
require "yaml"

ROOT = Pathname.new(__dir__).parent.realpath

options = {
  source: nil,
  target: nil,
  type: nil,
  query: nil,
  limit: 20
}

ARGV.each_with_index do |arg, index|
  case arg
  when "--source"
    options[:source] = ARGV[index + 1]
  when "--target"
    options[:target] = ARGV[index + 1]
  when "--type"
    options[:type] = ARGV[index + 1]
  when "--query"
    options[:query] = ARGV[index + 1]
  when "--limit"
    options[:limit] = ARGV[index + 1].to_i if ARGV[index + 1]
  end
end

if options.values_at(:source, :target, :type, :query).all? { |value| value.nil? || value.empty? }
  warn "Usage: ruby scripts/query-knowledge-graph.rb [--source path] [--target path] [--type edge] [--query text] [--limit N]"
  exit 2
end

def read_text(path)
  File.read(path.to_s, encoding: "UTF-8")
end

def yaml_file(path)
  YAML.safe_load(read_text(path), permitted_classes: [Date], aliases: false)
end

def rel(path)
  Pathname.new(path).relative_path_from(ROOT).to_s
end

def includes?(value, needle)
  return true if needle.nil? || needle.empty?

  value.to_s.downcase.include?(needle.downcase)
end

records = Dir.glob((ROOT + "knowledge/graphs/*.yaml").to_s).sort.flat_map do |path|
  graph_file = Pathname.new(path)
  data = yaml_file(graph_file)
  Array(data["edges"]).map do |edge|
    {
      graph_file: rel(graph_file),
      graph_id: data["id"],
      graph_source: data["source"],
      graph_status: data["status"],
      edge_type: edge["type"],
      target: edge["target"],
      reason: edge["reason"],
      validation: edge["validation"]
    }
  end
end

filtered = records.select do |record|
  next false unless includes?(record[:graph_source], options[:source])
  next false unless includes?(record[:target], options[:target])
  next false unless includes?(record[:edge_type], options[:type])

  if options[:query] && !options[:query].empty?
    haystack = [
      record[:graph_id],
      record[:graph_source],
      record[:edge_type],
      record[:target],
      record[:reason],
      record[:validation]
    ].join(" ")
    next false unless includes?(haystack, options[:query])
  end

  true
end.first(options[:limit])

puts [
  "graph_id",
  "graph_source",
  "edge_type",
  "target",
  "reason",
  "validation",
  "graph_file"
].join("\t")

filtered.each do |record|
  puts [
    record[:graph_id],
    record[:graph_source],
    record[:edge_type],
    record[:target],
    record[:reason],
    record[:validation],
    record[:graph_file]
  ].map { |value| value.to_s.gsub("\t", " ").gsub("\n", " ") }.join("\t")
end
