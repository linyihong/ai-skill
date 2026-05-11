#!/usr/bin/env ruby
# encoding: UTF-8
# Generate model-aware context loading views from the routing registry.

require "date"
require "pathname"
require "yaml"

ROOT = Pathname.new(__dir__).parent.realpath
OUTPUT_PATH = ROOT + "knowledge/runtime/model-context-report.md"

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

def group_by_profile(records)
  records.group_by { |record| record.dig("model", "profile").to_s }
         .sort_by { |profile, _| profile }
end

def group_by_compression(records)
  records.group_by { |record| record.dig("model", "compression_level").to_s }
         .sort_by { |level, _| level }
end

def build_report
  records = registry_records

  lines = []
  lines << "# Model Context Report"
  lines << ""
  lines << "本檔由 `ruby scripts/generate-model-context-report.rb --write` 產生，依 `knowledge/runtime/routing-registry.yaml` 的 model 欄位整理 model-aware context loading view。"
  lines << ""
  lines << "## Source Surfaces"
  lines << ""
  lines << "| Surface | Path | Purpose |"
  lines << "| --- | --- | --- |"
  lines << "| Routing registry | [`routing-registry.yaml`](routing-registry.yaml) | 提供每條 route 的 model profile 與 compression level。 |"
  lines << "| Model profiles | [`../../models/profiles/README.md`](../../models/profiles/README.md) | 定義 `small`、`large`、`specialized` 的讀取深度與 guardrails。 |"
  lines << "| Compression strategy | [`../../models/compression/README.md`](../../models/compression/README.md) | 定義 `summary-first`、`source-backed`、`graph-assisted` 等壓縮層級。 |"
  lines << ""
  lines << "## Profile View"
  lines << ""
  group_by_profile(records).each do |profile, profile_records|
    lines << "### `#{md_escape(profile.empty? ? "unspecified" : profile)}`"
    lines << ""
    lines << "| Route | Primary source | Compression | Reason |"
    lines << "| --- | --- | --- | --- |"
    profile_records.each do |record|
      model = record["model"] || {}
      lines << "| `#{md_escape(record["id"])}` | `#{md_escape(record["primary_source"])}` | `#{md_escape(model["compression_level"])}` | #{md_escape(model["reason"])} |"
    end
    lines << ""
  end
  lines << "## Compression View"
  lines << ""
  lines << "| Compression level | Routes | Escalation note |"
  lines << "| --- | --- | --- |"
  group_by_compression(records).each do |level, level_records|
    route_ids = level_records.map { |record| "`#{md_escape(record["id"])}`" }.join(", ")
    note = case level
           when "summary-first"
             "適合先用 registry / summary 判斷 relevance；修改 source 時升級。"
           when "source-backed"
             "需要 primary source 與 required dependencies；適合 writeback、migration 或 domain work。"
           when "graph-assisted"
             "需要 graph records 輔助 dependency / conflict / promotion reasoning。"
           else
             "依 `models/compression/README.md` 的 escalation rules 判斷。"
           end
    lines << "| `#{md_escape(level)}` | #{route_ids} | #{note} |"
  end
  lines << ""
  lines << "## Agent Output Shape"
  lines << ""
  lines << "使用本 report 決定 model-aware loading 時，回報："
  lines << ""
  lines << "```text"
  lines << "Profile:"
  lines << "Compression level:"
  lines << "Primary source:"
  lines << "Summaries used:"
  lines << "Required full sources:"
  lines << "Deferred sources:"
  lines << "Escalation trigger:"
  lines << "Validation signal:"
  lines << "```"
  lines << ""
  lines << "## Validation"
  lines << ""
  lines << "- 產生前應先確認 `routing-registry.yaml` 可通過 `ruby scripts/validate-knowledge-runtime.rb`。"
  lines << "- 產生後應重新執行 `ruby scripts/validate-knowledge-runtime.rb`，檢查本 report links。"
  lines << "- 本報告是 generated view，不取代 `models/profiles/README.md`、`models/compression/README.md` 或 routing registry。"
  lines << ""
  lines.join("\n")
end

if ARGV.include?("--write")
  OUTPUT_PATH.write(build_report, mode: "w", encoding: "UTF-8")
  puts "Wrote #{OUTPUT_PATH.relative_path_from(ROOT)}"
else
  puts build_report
end
