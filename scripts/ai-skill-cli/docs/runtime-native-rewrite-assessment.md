# Runtime Native Rewrite 評估

> 上游計畫：[`2026-05-21-0834-cross-platform-go-script-runtime.md`](../../../plans/archived/2026-05-21-0834-cross-platform-go-script-runtime.md)

本文件記錄 Phase 3 中哪些 Ruby runtime tooling 已改寫為 Go native，以及刪除舊 runtime surfaces 的驗收依據。

## 決策原則

- 無寫入、輸出穩定、資料來源明確的 validator / query 優先 native。
- 會寫 generated Markdown、SQLite index 或 `runtime.db` 的 generator / compiler 必須先有 golden fixture。
- 涉及 compiler source-of-truth、schema 建立、prose extraction 或 aggregation 的工具，必須有 source-to-DB fixture 與 stable snapshot test。
- Runtime path 不保留 Ruby / `sqlite3` CLI wrapper；rollback 應走 git revert 或 release artifact，而不是雙 compiler。

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
| `scripts/refresh-knowledge-runtime.rb` | `ai-skill runtime refresh` | 低 | 已刪除。預設 refresh 由 Go 寫 reports / index 並跑 native checks。 | native no-Ruby fixture、generated surface fixture、no partial success fixture 已覆蓋。 |
| `runtime/compiler/compiler-engine.rb` | `ai-skill runtime compile` | 完成 | 已刪除。Go native source-to-DB compiler 讀 runtime YAML、compiler-rules mapping 與 deterministic prose sources。 | custom DB fixture、generated surface assertion、row-count validation、stable snapshot test。 |
| `scripts/migrate-runtime-config-to-sqlite.rb` | `ai-skill runtime compile` future source-to-DB compiler | 低 | 已刪除。舊 migration helper 會誤導 agent 走 Ruby / sqlite3 CLI path；目前 supported path 是 `ai-skill runtime compile` 與 pre-commit compiler integration。 | runtime DB validation、compiler snapshot fixture。 |
| `scripts/init-runtime-state-db.rb` | future `ai-skill runtime state init` | 低 | 已刪除。Mutable runtime-state scope 尚未啟用，保留 Ruby initializer 會造成誤用。 | 等 scope 明確後新增 Go-native fixture。 |
| `scripts/sync-runtime-yaml-from-embedded.rb` | future source restoration migration | 低 | 已刪除。避免從 embedded data 回寫 stale YAML；若要恢復 standalone YAML，需獨立 source restoration plan。 | source restoration plan 再補 fixture。 |

## 下一步

1. CI / release gate 跑 native default smoke；runtime refresh / validate / compile legacy wrapper 已移除。
2. 若 runtime schema 未來新增 table，先更新 YAML source、Go compiler fixture、runtime DB validation，再刪除或重建 generated output。
