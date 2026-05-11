#!/usr/bin/env ruby
# encoding: UTF-8
# Query the local SQLite runtime lookup cache.

require "open3"
require "pathname"

ROOT = Pathname.new(__dir__).parent.realpath
DEFAULT_DB = ROOT + "knowledge/runtime/sqlite/runtime-index.sqlite"

query = nil
limit = 8
db = DEFAULT_DB
layer = nil
type = nil
status = nil

ARGV.each_with_index do |arg, index|
  case arg
  when "--limit"
    limit = ARGV[index + 1].to_i if ARGV[index + 1]
  when "--db"
    db = ROOT + ARGV[index + 1] if ARGV[index + 1]
  when "--layer"
    layer = ARGV[index + 1] if ARGV[index + 1]
  when "--type"
    type = ARGV[index + 1] if ARGV[index + 1]
  when "--status"
    status = ARGV[index + 1] if ARGV[index + 1]
  else
    query ||= arg unless arg.start_with?("--")
  end
end

unless query && !query.strip.empty?
  warn "Usage: ruby scripts/query-runtime-index.rb <query> [--limit N] [--layer L] [--type T] [--status S] [--db path]"
  exit 2
end

unless db.exist?
  warn "Missing SQLite index: #{db.relative_path_from(ROOT)}"
  warn "Run: ruby scripts/generate-runtime-sqlite-index.rb"
  exit 1
end

phrase = query.gsub('"', '""').gsub("'", "''")
match_literal = "'\"#{phrase}\"'"
filters = []
filters << "atoms.layer = '#{layer.gsub("'", "''")}'" if layer
filters << "atoms.type = '#{type.gsub("'", "''")}'" if type
filters << "atoms.status = '#{status.gsub("'", "''")}'" if status
filter_sql = filters.empty? ? "" : "AND #{filters.join(" AND ")}"
sql = <<~SQL
  SELECT bm25(fts) AS rank,
         atoms.id,
         atoms.source_path,
         atoms.layer,
         atoms.type,
         atoms.status,
         atoms.priority,
         atoms.confidence,
         atoms.context_cost,
         atoms.summary
  FROM fts
  JOIN atoms ON atoms.id = fts.id
  WHERE fts MATCH #{match_literal}
  #{filter_sql}
  ORDER BY rank,
           CASE atoms.priority WHEN 'P0' THEN 0 WHEN 'P1' THEN 1 WHEN 'P2' THEN 2 ELSE 3 END,
           CASE atoms.confidence WHEN 'high' THEN 0 WHEN 'medium' THEN 1 ELSE 2 END,
           CASE atoms.context_cost WHEN 'low' THEN 0 WHEN 'medium' THEN 1 ELSE 2 END,
           atoms.id
  LIMIT #{limit};
SQL

stdout, stderr, status = Open3.capture3(
  "sqlite3",
  "-header",
  "-separator",
  "\t",
  db.to_s,
  sql
)

unless status.success?
  warn stderr
  exit status.exitstatus || 1
end

puts stdout
