# SQLite Runtime Index

`knowledge/runtime/sqlite/` 規劃以 SQLite 作為本 repository 的 generated lookup cache，支援快速搜尋、FTS、tag filtering 與 task-intent routing。它的目的在於降低 agent 讀取大量 Markdown / YAML 的 context cost，不取代任何 canonical source。

## Source-Of-Truth Rule

| Surface | Role |
| --- | --- |
| Markdown / YAML source files | Canonical source-of-truth，例如 `skills/*/feedback_history/`、`shared-rules/`、`knowledge/summaries/`、`knowledge/graphs/`、`knowledge/runtime/routing-registry.yaml`。 |
| SQLite database | Generated lookup cache；可刪除、可重建、不可作為唯一資料來源。 |
| Generator script | 從 canonical sources deterministic 產生 SQLite tables / FTS index。 |
| Query helper | 從 SQLite 回傳少量 candidate `source_path`、summary、tags、score 與 validation signal。 |

SQLite row 可以索引 feedback lessons，但 feedback lesson 全文仍留在 `skills/<skill>/feedback_history/`。Agent 查詢 SQLite 後，只能把結果當作 candidate list；需要執行、修改或高信心判斷時，仍需讀 canonical source。

## Why It Can Save Tokens

SQLite index 可讓 agent 先做低成本查詢：

1. 使用 task intent、keyword、tag、layer、status、priority 或 confidence 篩選候選知識。
2. 只讀前幾筆候選的 `summary`、`source_path`、`when_to_read` 與 `validation_signal`。
3. 依 `source_of_truth_gate` 決定是否讀完整 Markdown / YAML。
4. 對 low-risk lookup 使用 index / summary；對 writeback、promotion、shared-rule 或 skill execution 升級到 source-backed 讀取。

省 token 的關鍵不是把內容搬進 SQLite，而是用 SQLite 避免一開始讀完整資料夾。

## Planned Tables

第一版 schema 應從 `metadata/schema.md`、summary table、graph records 與 routing registry 產生。

| Table | Purpose | Canonical source |
| --- | --- | --- |
| `atoms` | Knowledge Atom / summary / route 的主要索引欄位。 | `knowledge/summaries/*.md`、`metadata/schema.md` |
| `sources` | Canonical `source_path`、layer、title、mtime / checksum。 | Markdown / YAML source files |
| `routes` | Task intent 到 primary source / dependencies / model profile。 | `knowledge/runtime/routing-registry.yaml` |
| `edges` | depends / related / routes_to / promotes_from 關係。 | `knowledge/graphs/*.yaml` |
| `feedback_lessons` | skill feedback lesson 的標題、status、category、summary 與 source path。 | `skills/*/feedback_history/**` |
| `fts` | FTS5 full-text index，索引 title、summary、tags、trigger、when_to_read。 | Generated from tables above |

## Minimum Fields

SQLite records 至少保留：

| Field | Reason |
| --- | --- |
| `id` | Stable lookup key。 |
| `source_path` | 回到 canonical source。 |
| `layer` | 對應 `knowledge`、`feedback`、`intelligence`、`skills`、`shared-rules` 等分層。 |
| `type` | `rule`、`workflow`、`intelligence`、`index`、`schema`、`reference` 等。 |
| `status` | `candidate`、`validated`、`stable`、`deprecated`。 |
| `priority` / `confidence` / `context_cost` | 支援 runtime ranking。 |
| `tags` / `domains` | 支援快速篩選。 |
| `summary` | 低 token preview。 |
| `when_to_read` | 何時升級讀全文。 |
| `validation_signal` | 查詢結果是否足以支持下一步。 |

## Repository Policy

- 不 commit 生成出的 `.sqlite` / `.db` binary，除非未來另有明確 governance 決策。
- Commit schema、generator、query helper、validation tests 與 deterministic fixtures。
- SQLite generator 必須能在 clean checkout 重新產生相同 lookup content。
- SQLite query helper 預設只輸出 top candidates，不輸出大量全文。
- 若 SQLite index stale，agent 必須回到 canonical sources，並依 `refresh-policy.yaml` 判斷是否 regenerate。

## Cold Data Archive Role

SQLite / FTS 是 cold feedback lesson 與 runtime navigation 的查找層。當 lesson 數量增加、或 agent 只需要候選列表而不需要讀全文時，先查 SQLite：

```bash
ruby scripts/query-runtime-index.rb "<keyword>" --limit 5
```

查詢結果仍只是 candidate list。需要 promotion、debug、failure learning、修改 lesson 或高信心判斷時，必須讀回 `source_path` 指向的 Markdown / YAML。

## Tooling

| Tool | Role |
| --- | --- |
| `scripts/generate-runtime-sqlite-index.rb` | 從 summaries、graphs、registry、feedback lessons 產生本機 SQLite DB。 |
| `scripts/query-runtime-index.rb` | 用 keyword 查候選 source，輸出少量 rows。 |
| `scripts/validate-runtime-sqlite-index.rb` | 檢查 DB integrity、row counts、source path existence、FTS availability 與 git ignore 邊界。 |

## Usage

```bash
ruby scripts/generate-runtime-sqlite-index.rb
ruby scripts/validate-runtime-sqlite-index.rb
ruby scripts/query-runtime-index.rb feedback --limit 5
```

預設 DB 路徑：

```text
knowledge/runtime/sqlite/runtime-index.sqlite
```

此 DB 由 `.gitignore` 排除，應可在 clean checkout 由 generator 重建。

## Validation

第一版導入完成前，至少要驗證：

- SQLite schema 可建立。
- Generated DB 不進 git status。
- Query 能用 keyword / tag 找到 feedback lesson、summary、graph route。
- Query result 只作候選，不跳過 source-of-truth gate。
- Runtime reports、routing registry、summary / graph counts 與 SQLite counts 可交叉檢查。
