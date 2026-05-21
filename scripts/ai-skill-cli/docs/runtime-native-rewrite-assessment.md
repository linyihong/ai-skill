# Runtime Native Rewrite 評估

> 上游計畫：[`2026-05-21-0834-cross-platform-go-script-runtime.md`](../../../plans/archived/2026-05-21-0834-cross-platform-go-script-runtime.md)

本文件決定 Phase 3 中哪些 Ruby runtime tooling 適合先改寫為 Go native，哪些應維持 wrapper-first，避免在沒有 parity fixture 前直接替換 production compiler。

## 決策原則

- 無寫入、輸出穩定、資料來源明確的 validator / query 優先 native。
- 會寫 generated Markdown、SQLite index 或 `runtime.db` 的 generator / compiler 必須先有 golden fixture。
- 涉及 compiler source-of-truth、schema 建立、prose extraction 或 aggregation 的工具，必須先建立 Ruby vs Go parity test。
- wrapper mode 可保留，但必須固定 UTF-8 env、清楚回報 Ruby / `sqlite3` / Git 依賴，且不得 partial success。

## 建議排序

| 舊入口 | 目前 CLI | Native 優先度 | 建議處置 | 必要 parity |
| --- | --- | --- | --- | --- |
| `scripts/validate-runtime-db.rb` | `ai-skill runtime validate` | 高 | 已刪除。Go native 覆蓋 integrity、required tables、row count、JSON、compiler metadata 與 stale metadata warning。 | invalid DB、missing table、invalid JSON、stale metadata fixture 已覆蓋。 |
| `scripts/validate-runtime-sqlite-index.rb` | `ai-skill runtime validate` | 高 | 已刪除。SQLite integrity、tables、row counts、source references、checksum、FTS count、ranked query 與 git-ignore boundary 已 native。 | missing DB / table、stale checksum、FTS count mismatch、git-ignore boundary fixture 已覆蓋。 |
| `scripts/query-runtime-index.rb` | `ai-skill runtime query` | 高 | 已刪除。查詢 SQLite / FTS、filter、limit 與 empty result 行為已 native；不依賴外部 `sqlite3` CLI。 | ranking / filter / empty result / missing DB fixture 已覆蓋。 |
| `scripts/query-knowledge-graph.rb` | `ai-skill runtime query --graph` | 中 | 已刪除。YAML parsing、source / target / type / keyword filters、empty result 與 missing filter 已 fixed。 | source / target / type / keyword filter、empty result fixture 已覆蓋。 |
| `scripts/generate-runtime-sqlite-index.rb` | `ai-skill runtime refresh` | 中 | 已刪除。Runtime refresh 預設寫入 Go-native SQLite index。 | golden DB row counts、FTS fixture、source checksum、recursive feedback fixture 已覆蓋。 |
| `scripts/generate-knowledge-runtime-report.rb` | `ai-skill runtime refresh` | 中 | 已刪除。Runtime refresh 預設寫入 Go-native runtime report。 | routing registry / summaries / graphs golden anchors 與 native report write 已覆蓋。 |
| `scripts/generate-model-context-report.rb` | `ai-skill runtime refresh` | 中 | 已刪除。Runtime refresh 預設寫入 Go-native model context report。 | profile / compression grouping anchors 與 native report write 已覆蓋。 |
| `scripts/generate-model-checklists.rb` | `ai-skill runtime refresh` | 中 | 已刪除。Runtime refresh 預設寫入 Go-native model checklists。 | per-model checklist anchors 與 native report write 已覆蓋。 |
| `scripts/refresh-knowledge-runtime.rb` | `ai-skill runtime refresh` | 低 | 已刪除。`--legacy-wrapper` 對 refresh 回 `legacy_runtime_refresh_removed`；預設 refresh 由 Go 寫 reports / index 並跑 native checks。 | native no-Ruby fixture、generated surface fixture、no partial success fixture 已覆蓋。 |
| `runtime/compiler/compiler-engine.rb` | `ai-skill runtime compile` | 低 | 預設 Go native snapshot mode；`--legacy-wrapper` 回 Ruby compiler。 | runtime source keyword、`--check` no-op、`runtime.db` generated surface assertion、Ruby vs native generated-surfaces / metadata snapshot parity；完整 runtime config row-count parity 待真正 port compiler。 |
| `scripts/migrate-runtime-config-to-sqlite.rb` | `ai-skill runtime migrate` / `compile` | Deferred | 暫不改寫，現有 compiler path 已吸收大部分需求。 | idempotent migration fixture。 |
| `scripts/init-runtime-state-db.rb` | `ai-skill runtime state init` | Deferred | 等 mutable runtime-state scope 明確後再處理。 | custom DB path、idempotent schema fixture。 |
| `scripts/sync-runtime-yaml-from-embedded.rb` | `ai-skill runtime sync-yaml` | Deferred | 屬於 source sync / recovery 工具，需先確認 lifecycle owner。 | embedded-to-yaml golden fixture。 |

## 下一步

1. CI / release gate 需改跑 native default smoke；runtime refresh / validate legacy wrapper 已移除。
2. Compiler snapshot mode 已移除 Ruby dependency，但真正 source-to-DB compiler port 仍需完整 runtime config row-count parity。
