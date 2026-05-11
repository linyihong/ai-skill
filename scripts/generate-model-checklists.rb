#!/usr/bin/env ruby
# encoding: UTF-8
# Generate per-model context loading checklists from the routing registry.

require "date"
require "pathname"
require "yaml"

ROOT = Pathname.new(__dir__).parent.realpath
OUTPUT_PATH = ROOT + "knowledge/runtime/model-checklists.md"

def read_text(path)
  File.read(path.to_s, encoding: "UTF-8")
end

def yaml_file(relative_path)
  YAML.safe_load(read_text(ROOT + relative_path), permitted_classes: [Date], aliases: false)
end

def md_escape(value)
  value.to_s.gsub("|", "\\|").gsub("\n", " ")
end

def registry_records
  Array(yaml_file("knowledge/runtime/routing-registry.yaml")["records"])
end

def profile_guardrails(profile)
  case profile
  when "small"
    [
      "先讀 index、registry、summary 或 generated lookup。",
      "不可跳過 required bootstrap、source-of-truth gate 或 validation signal。",
      "需要修改 canonical source、遇到 conflict、缺 validation signal 時升級。"
    ]
  when "large"
    [
      "讀 primary source、required dependencies 與 task-relevant related sources。",
      "回報 deferred sources 與 validation signal。",
      "涉及 writeback、promotion、shared rules 或 migration 時保持 source-backed。"
    ]
  when "specialized"
    [
      "先讀 routing registry 與 primary source，再讀 domain workflow / technique / adapter。",
      "不得讓工具能力覆蓋 shared rules、authorization 或 source-of-truth。",
      "保留 domain-specific validation 與 project evidence boundary。"
    ]
  else
    [
      "先確認 registry record 的 model profile。",
      "依 `models/profiles/README.md` 與 `models/compression/README.md` 選讀取深度。"
    ]
  end
end

def build_report
  records = registry_records
  grouped = records.group_by { |record| record.dig("model", "profile").to_s }.sort_by { |profile, _| profile }

  lines = []
  lines << "# Model Checklists"
  lines << ""
  lines << "本檔由 `ruby scripts/generate-model-checklists.rb --write` 產生，將 routing registry 中的 model profile / compression level 轉成 agent 可直接使用的 context-loading checklist。"
  lines << ""
  lines << "## Source Surfaces"
  lines << ""
  lines << "| Surface | Path | Purpose |"
  lines << "| --- | --- | --- |"
  lines << "| Routing registry | [`routing-registry.yaml`](routing-registry.yaml) | 提供 route、primary source、dependencies、model profile 與 compression level。 |"
  lines << "| Model profiles | [`../../models/profiles/README.md`](../../models/profiles/README.md) | 定義 profile guardrails。 |"
  lines << "| Compression strategy | [`../../models/compression/README.md`](../../models/compression/README.md) | 定義 escalation rules。 |"
  lines << ""
  lines << "## Profile Checklists"
  lines << ""
  grouped.each do |profile, profile_records|
    display_profile = profile.empty? ? "unspecified" : profile
    lines << "### `#{md_escape(display_profile)}`"
    lines << ""
    lines << "Guardrails:"
    lines << ""
    profile_guardrails(profile).each { |item| lines << "- #{item}" }
    lines << ""
    lines << "| Route | Checklist |"
    lines << "| --- | --- |"
    profile_records.each do |record|
      model = record["model"] || {}
      dependencies = Array(record["required_dependencies"]).map { |path| "`#{path}`" }.join("<br>")
      checklist = [
        "Primary: `#{record["primary_source"]}`",
        "Compression: `#{model["compression_level"]}`",
        "Required: #{dependencies}",
        "Validation: #{md_escape(record["validation_signal"])}"
      ].join("<br>")
      lines << "| `#{md_escape(record["id"])}` | #{checklist} |"
    end
    lines << ""
  end
  lines << "## Escalation Checklist"
  lines << ""
  lines << "- Summary / registry 與 source-of-truth 可能不一致時，讀全文。"
  lines << "- 任務需要修改、commit、push、readback 或 promotion 時，升級到 `source-backed`。"
  lines << "- 涉及 safety、secrets、authorization、source/mirror 或 destructive actions 時，升級到 full source 和 shared rules。"
  lines << "- Routing registry 指向 candidate path，但 old entrypoint 仍 active 時，保留 old entrypoint gate。"
  lines << "- Validation signal 不足以支持結論時，停止並讀 required dependencies。"
  lines << ""
  lines << "## Validation"
  lines << ""
  lines << "- 產生前應先確認 `routing-registry.yaml` 可通過 `ruby scripts/validate-knowledge-runtime.rb`。"
  lines << "- 產生後應重新執行 `ruby scripts/validate-knowledge-runtime.rb`，檢查本 report links。"
  lines << "- 本檔是 generated view，不取代 model source docs 或 routing registry。"
  lines << ""
  lines.join("\n")
end

if ARGV.include?("--write")
  OUTPUT_PATH.write(build_report, mode: "w", encoding: "UTF-8")
  puts "Wrote #{OUTPUT_PATH.relative_path_from(ROOT)}"
else
  puts build_report
end
