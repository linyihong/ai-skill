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

ARGV.each_with_index do |arg, index|
  case arg
  when "--limit"
    limit = ARGV[index + 1].to_i if ARGV[index + 1]
  when "--db"
    db = ROOT + ARGV[index + 1] if ARGV[index + 1]
  else
    query ||= arg unless arg.start_with?("--")
  end
end

unless query && !query.strip.empty?
  warn "Usage: ruby scripts/query-runtime-index.rb <query> [--limit N] [--db path]"
  exit 2
end

unless db.exist?
  warn "Missing SQLite index: #{db.relative_path_from(ROOT)}"
  warn "Run: ruby scripts/generate-runtime-sqlite-index.rb"
  exit 1
end

phrase = query.gsub('"', '""').gsub("'", "''")
match_literal = "'\"#{phrase}\"'"
sql = <<~SQL
  SELECT atoms.id,
         atoms.source_path,
         atoms.layer,
         atoms.type,
         atoms.status,
         atoms.confidence,
         atoms.context_cost,
         atoms.summary
  FROM fts
  JOIN atoms ON atoms.id = fts.id
  WHERE fts MATCH #{match_literal}
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
