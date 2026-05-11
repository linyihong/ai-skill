#!/usr/bin/env ruby
# encoding: UTF-8
# Regenerate and validate all knowledge runtime generated surfaces.

require "open3"
require "pathname"

ROOT = Pathname.new(__dir__).parent.realpath

COMMANDS = [
  ["ruby", "scripts/generate-model-context-report.rb", "--write"],
  ["ruby", "scripts/generate-model-checklists.rb", "--write"],
  ["ruby", "scripts/generate-knowledge-runtime-report.rb", "--write"],
  ["ruby", "scripts/generate-runtime-sqlite-index.rb"],
  ["ruby", "scripts/validate-runtime-sqlite-index.rb"],
  ["ruby", "scripts/validate-knowledge-runtime.rb"]
].freeze

def run_command(command)
  puts "$ #{command.join(" ")}"
  stdout, stderr, status = Open3.capture3(*command, chdir: ROOT.to_s)
  puts stdout unless stdout.empty?
  warn stderr unless stderr.empty?
  abort "Command failed: #{command.join(" ")}" unless status.success?
end

COMMANDS.each { |command| run_command(command) }

puts "Knowledge runtime refresh OK"
