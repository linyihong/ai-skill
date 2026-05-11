#!/usr/bin/env ruby
# encoding: UTF-8
# Validate the generated SQLite runtime lookup cache.

require "open3"
require "pathname"
require "digest"

ROOT = Pathname.new(__dir__).parent.realpath
DB = ROOT + "knowledge/runtime/sqlite/runtime-index.sqlite"

def read_text(path)
  File.read(path.to_s, encoding: "UTF-8")
end

def utf8(value)
  text = value.to_s
  text = text.dup.force_encoding("UTF-8") if text.encoding == Encoding::ASCII_8BIT || text.encoding == Encoding::US_ASCII
  text.encode("UTF-8", invalid: :replace, undef: :replace)
end

def run_sql(sql)
  stdout, stderr, status = Open3.capture3("sqlite3", "-batch", DB.to_s, sql)
  unless status.success?
    warn stderr
    exit status.exitstatus || 1
  end
  utf8(stdout).strip
end

unless DB.exist?
  warn "Missing SQLite index: #{DB.relative_path_from(ROOT)}"
  warn "Run: ruby scripts/generate-runtime-sqlite-index.rb"
  exit 1
end

integrity = run_sql("PRAGMA integrity_check;")
abort "SQLite integrity failed: #{integrity}" unless integrity == "ok"

tables = run_sql("SELECT name FROM sqlite_master WHERE type IN ('table', 'virtual') ORDER BY name;").lines.map(&:strip)
%w[atoms sources edges fts].each do |table|
  abort "Missing table: #{table}" unless tables.include?(table)
end

atom_count = run_sql("SELECT COUNT(*) FROM atoms;").to_i
source_count = run_sql("SELECT COUNT(*) FROM sources;").to_i
edge_count = run_sql("SELECT COUNT(*) FROM edges;").to_i
fts_count = run_sql("SELECT COUNT(*) FROM fts;").to_i

abort "atoms table is empty" if atom_count.zero?
abort "sources table is empty" if source_count.zero?
abort "fts count does not match atoms count" unless fts_count == atom_count

missing_sources = run_sql("SELECT source_path FROM atoms WHERE source_path NOT IN (SELECT source_path FROM sources) LIMIT 10;")
abort "Atoms reference missing sources:\n#{missing_sources}" unless missing_sources.empty?

source_rows = run_sql("SELECT source_path || char(9) || checksum FROM sources WHERE checksum IS NOT NULL AND checksum != '';")
source_rows.each_line do |line|
  source_path, checksum = line.chomp.split("\t", 2)
  path = ROOT + source_path
  abort "Source path missing on disk: #{source_path}" unless path.exist?
  next unless path.file?

  current = Digest::SHA256.hexdigest(read_text(path))
  abort "Stale checksum for #{source_path}" unless current == checksum
end

feedback_hits = run_sql(%q{SELECT COUNT(*) FROM fts WHERE fts MATCH '"feedback"';}).to_i
route_hits = run_sql(%q{SELECT COUNT(*) FROM fts WHERE fts MATCH '"route"';}).to_i
ranked_route = run_sql(%q{SELECT atoms.id FROM fts JOIN atoms ON atoms.id = fts.id WHERE fts MATCH '"feedback"' ORDER BY bm25(fts), CASE atoms.priority WHEN 'P0' THEN 0 WHEN 'P1' THEN 1 WHEN 'P2' THEN 2 ELSE 3 END LIMIT 1;})
abort "Expected feedback FTS hits" if feedback_hits.zero?
abort "Expected route FTS hits" if route_hits.zero?
abort "Expected ranked query result" if ranked_route.empty?

ignored_stdout, = Open3.capture2("git", "check-ignore", DB.relative_path_from(ROOT).to_s, chdir: ROOT.to_s)
abort "Generated DB is not ignored by git" if ignored_stdout.strip.empty?

puts "SQLite runtime index validation OK"
puts "atoms=#{atom_count}"
puts "sources=#{source_count}"
puts "edges=#{edge_count}"
puts "fts=#{fts_count}"
