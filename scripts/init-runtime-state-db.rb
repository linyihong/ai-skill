#!/usr/bin/env ruby
# frozen_string_literal: true

# init-runtime-state-db.rb
#
# Creates/initializes runtime-state.db — the mutable runtime state database
# that agents write to during execution.
#
# Usage:
#   ruby scripts/init-runtime-state-db.rb
#   ruby scripts/init-runtime-state-db.rb --db /path/to/runtime-state.db
#
# This is idempotent: re-running will not destroy existing data (CREATE IF NOT EXISTS).

require 'open3'
require 'optparse'
require 'fileutils'

DEFAULT_DB_PATH = File.join(File.dirname(__FILE__), '..', 'runtime', 'runtime-state.db')

def sqlite_exec(db_path, sql)
  o, e, s = Open3.capture3('sqlite3', db_path, stdin_data: sql)
  warn "sqlite3 error: #{e.strip}" unless s.success?
  s.success?
end

def create_schema(db_path)
  schema = <<~SQL
    -- ============================================================
    -- execution_state: tracks the current phase and execution context
    -- ============================================================
    CREATE TABLE IF NOT EXISTS execution_state (
      id            INTEGER PRIMARY KEY AUTOINCREMENT,
      phase         TEXT NOT NULL DEFAULT 'bootstrap',
      sub_phase     TEXT,
      status        TEXT NOT NULL DEFAULT 'idle'
                    CHECK (status IN ('idle','running','paused','completed','failed')),
      started_at    TEXT,
      updated_at    TEXT NOT NULL DEFAULT (datetime('now')),
      metadata      TEXT  -- JSON: arbitrary key-value context
    );

    -- ============================================================
    -- obligation_status: tracks completion status of each obligation
    -- ============================================================
    CREATE TABLE IF NOT EXISTS obligation_status (
      id              INTEGER PRIMARY KEY AUTOINCREMENT,
      obligation_id   TEXT NOT NULL,
      phase           TEXT NOT NULL,
      status          TEXT NOT NULL DEFAULT 'pending'
                      CHECK (status IN ('pending','in_progress','completed','blocked','skipped')),
      verified_at     TEXT,
      verified_by     TEXT,  -- e.g. 'agent', 'human', 'validator'
      notes           TEXT,
      metadata        TEXT,  -- JSON
      created_at      TEXT NOT NULL DEFAULT (datetime('now')),
      updated_at      TEXT NOT NULL DEFAULT (datetime('now')),
      UNIQUE(obligation_id, phase)
    );

    -- ============================================================
    -- transaction_state: tracks the current transaction lifecycle
    -- ============================================================
    CREATE TABLE IF NOT EXISTS transaction_state (
      id              INTEGER PRIMARY KEY AUTOINCREMENT,
      transaction_id  TEXT NOT NULL UNIQUE,
      state           TEXT NOT NULL DEFAULT 'closed'
                      CHECK (state IN ('closed','open','staging','commit_ready','committed','pushed','verified')),
      phase           TEXT,
      description     TEXT,
      started_at      TEXT,
      committed_at    TEXT,
      pushed_at       TEXT,
      verified_at     TEXT,
      metadata        TEXT,  -- JSON
      created_at      TEXT NOT NULL DEFAULT (datetime('now')),
      updated_at      TEXT NOT NULL DEFAULT (datetime('now'))
    );

    -- ============================================================
    -- execution_log: append-only log of execution events
    -- ============================================================
    CREATE TABLE IF NOT EXISTS execution_log (
      id          INTEGER PRIMARY KEY AUTOINCREMENT,
      phase       TEXT NOT NULL,
      event       TEXT NOT NULL,
      detail      TEXT,
      metadata    TEXT,  -- JSON
      created_at  TEXT NOT NULL DEFAULT (datetime('now'))
    );

    -- ============================================================
    -- Indexes for common queries
    -- ============================================================
    CREATE INDEX IF NOT EXISTS idx_obligation_status_phase
      ON obligation_status(phase, status);
    CREATE INDEX IF NOT EXISTS idx_transaction_state_status
      ON transaction_state(state);
    CREATE INDEX IF NOT EXISTS idx_execution_log_phase
      ON execution_log(phase, created_at);
  SQL

  sqlite_exec(db_path, schema)
end

# ---- CLI ----
db_path = DEFAULT_DB_PATH

OptionParser.new do |opts|
  opts.banner = "Usage: #{$PROGRAM_NAME} [options]"
  opts.on('--db PATH', 'Path to runtime-state.db (default: runtime/runtime-state.db)') { |v| db_path = v }
  opts.on('-h', '--help', 'Show help') { puts opts; exit }
end.parse!

FileUtils.mkdir_p(File.dirname(db_path))

if create_schema(db_path)
  puts "✓ runtime-state.db initialized at #{db_path}"
  sqlite_exec(db_path, "INSERT OR IGNORE INTO execution_state (phase, status, started_at) VALUES ('bootstrap', 'idle', datetime('now'));")
  puts "  - execution_state: bootstrap/idle"
  puts "  - obligation_status: ready"
  puts "  - transaction_state: ready"
  puts "  - execution_log: ready"
else
  puts "✗ Failed to initialize runtime-state.db"
  exit 1
end
