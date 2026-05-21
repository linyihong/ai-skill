# Runtime Native Rewrite 評估

> 上游計畫：[`2026-05-21-0834-cross-platform-go-script-runtime.md`](../../../plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md)

本文件決定 Phase 3 中哪些 Ruby runtime tooling 適合先改寫為 Go native，哪些應維持 wrapper-first，避免在沒有 parity fixture 前直接替換 production compiler。

## 決策原則

- 無寫入、輸出穩定、資料來源明確的 validator / query 優先 native。
- 會寫 generated Markdown、SQLite index 或 `runtime.db` 的 generator / compiler 必須先有 golden fixture。
- 涉及 compiler source-of-truth、schema 建立、prose extraction 或 aggregation 的工具，必須先建立 Ruby vs Go parity test。
- wrapper mode 可保留，但必須固定 UTF-8 env、清楚回報 Ruby / `sqlite3` / Git 依賴，且不得 partial success。

## 建議排序

| 舊入口 | 目前 CLI | Native 優先度 | 建議處置 | 必要 parity |
| --- | --- | --- | --- | --- |
| `scripts/validate-runtime-db.rb` | `ai-skill runtime validate` | 高 | 已開始 Go native slice。`modernc.org/sqlite` 已覆蓋 integrity、required tables、row count、JSON、compiler metadata 與 stale metadata warning；Ruby validator 暫保留作 parity guard。 | invalid DB、missing table、invalid JSON、stale metadata fixture 已覆蓋。 |
| `scripts/validate-runtime-sqlite-index.rb` | `ai-skill runtime validate` | 高 | 已開始 Go native slice。SQLite integrity、tables、row counts、source references、checksum、FTS count 與 basic ranked query 已 native；Git ignore 邊界仍需 external Git / Ruby parity guard。 | missing DB / table、stale checksum、FTS count mismatch 已覆蓋；git-ignore boundary fixture 待補。 |
| `scripts/query-runtime-index.rb` | `ai-skill runtime query` | 高 | 已開始 Go native slice。查詢 SQLite / FTS、filter、limit 與 empty result 行為已可直接測；不依賴外部 `sqlite3` CLI。 | ranking / filter / empty result / missing DB fixture 已覆蓋。 |
| `scripts/query-knowledge-graph.rb` | `ai-skill runtime query` | 中 | 可 Go native，但 YAML parsing 與 graph schema 需先固定 golden cases。 | source / target / type / keyword filter fixture。 |
| `scripts/generate-runtime-sqlite-index.rb` | `ai-skill runtime refresh` | 中 | 先 wrapper，建立 golden SQLite index 後再 native。寫入 git-ignored DB，需驗證 deterministic rows / checksums / FTS。 | golden DB row counts、checksum、FTS fixture。 |
| `scripts/generate-knowledge-runtime-report.rb` | `ai-skill runtime refresh` | 中 | 先 wrapper，建立 golden Markdown output 後再 native。 | routing registry / summaries / graphs golden report。 |
| `scripts/generate-model-context-report.rb` | `ai-skill runtime refresh` | 中 | 先 wrapper，建立 model grouping golden output 後再 native。 | profile / compression grouping fixture。 |
| `scripts/generate-model-checklists.rb` | `ai-skill runtime refresh` | 中 | 先 wrapper，建立 checklist golden output 後再 native。 | per-model checklist fixture。 |
| `scripts/refresh-knowledge-runtime.rb` | `ai-skill runtime refresh` | 低 | 保持 wrapper orchestrator，直到各 generator / validator native 後再改成 Go orchestrator。 | partial failure blocks success、ordered step summary、no partial success fixture。 |
| `runtime/compiler/compiler-engine.rb` | `ai-skill runtime compile` | 低 | 保持 wrapper。不得直接替換 production compiler；先建立 Ruby vs Go parity test。 | runtime source keyword、`--check` no-op、`runtime.db` generated surface assertion、schema parity。 |
| `scripts/migrate-runtime-config-to-sqlite.rb` | `ai-skill runtime migrate` / `compile` | Deferred | 暫不改寫，現有 compiler path 已吸收大部分需求。 | idempotent migration fixture。 |
| `scripts/init-runtime-state-db.rb` | `ai-skill runtime state init` | Deferred | 等 mutable runtime-state scope 明確後再處理。 | custom DB path、idempotent schema fixture。 |
| `scripts/sync-runtime-yaml-from-embedded.rb` | `ai-skill runtime sync-yaml` | Deferred | 屬於 source sync / recovery 工具，需先確認 lifecycle owner。 | embedded-to-yaml golden fixture。 |

## 下一步

1. 補 `runtime validate` 的 SQLite index git-ignore boundary fixture，決定保留 external Git check 或以 Go 呼叫 Git。
2. 補 `runtime query` 的 `query-knowledge-graph.rb` native slice，先固定 graph filter / empty result fixture。
3. Generator / compiler 只能在 golden fixture 與 Ruby vs Go parity test 完成後替換。
