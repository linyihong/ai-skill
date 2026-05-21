# 舊腳本 parity 盤點：Ai-skill CLI Runtime

> **上游計畫**：[`2026-05-21-0834-cross-platform-go-script-runtime.md`](../../../plans/archived/2026-05-21-0834-cross-platform-go-script-runtime.md)

本文件補上 Phase 0 的舊功能覆蓋率盤點。任何 Go CLI 實作若宣稱替代既有腳本，必須能從本文件反查：舊腳本做了什麼、未來命令接到哪裡、目前 parity 狀態、需要哪個 fixture 或測試證明。

## 狀態分類

| 狀態 | 意義 |
| --- | --- |
| `native target` | Phase 2 或 Phase 3 目標是 Go native 實作，舊腳本保留到 parity 測試通過。 |
| `wrapper first` | 先由 CLI 呼叫既有 Ruby / shell 行為，等 parity fixture 完成後再移植。 |
| `deferred` | 不在第一批 CLI 取代範圍內，但必須保留可執行或明確列為後續工作。 |
| `tool-specific` | 僅服務特定工具或本機設定，不應變成通用 CLI 預設行為。 |
| `hook adapter` | Git hook 仍是 Git adapter surface；CLI 可安裝或觸發，但不取代 hook 本身。 |
| `deleted` | 新 CLI parity、fixtures、文件與 release gate 通過後，舊入口已刪除。 |
| `legacy rollback` | 新 CLI 已是預設 path，但舊入口短期保留作 rollback / parity reference。 |

## 高風險覆蓋規則

- 每個現有腳本至少要有一列 parity disposition；不得以「未盤點」進入 Go implementation。
- 會寫檔、commit、push、同步 tool mirror、更新 runtime.db、更新 generated reports 或寫入使用者設定的腳本，必須有 fixture 或 BDD scenario。
- `close-loop`、`runtime refresh`、`runtime compile`、`init-project`、`sync-cursor-bundle` 是高風險路徑；缺少 parity 測試時不得宣稱 replacement 完成。
- 手機與桌面支援決策不能只看新 CLI；必須確認舊腳本在該平台的不可攜假設已被新命令處理或明確排除。
- Phase 6 closure 預設刪除 replacement 範圍內的舊 shell / Ruby / Python script。只有 `hook adapter`、`tool-specific` adapter 或明確短期 thin wrapper 可保留；保留時必須寫明 owner、移除條件與期限。

## 舊腳本功能盤點

| 舊入口 | 目前功能 | 輸入 | 輸出 / 副作用 | 外部依賴 | 目標 CLI | 狀態 | 必要驗證 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `scripts/init-new-project.sh` | 建立新專案 tool bootstrap：`.roomodes`、`.cursor/rules/`、`.cursor/hooks.json`、`CLAUDE.md`、`.agent-goals/` | 目標專案路徑、`--dry-run`、`--force`、`--tools` | 寫入目標專案設定；dry-run 只列出計畫 | POSIX shell、檔案系統權限 | `ai-skill init-project` | `native target` | fake home / target project fixture；既有檔案衝突 fixture；dry-run 無寫入 assertion（Phase 2 dry-run planner 已覆蓋；write mode 待 template parity） |
| `scripts/agent-goals.sh` | 管理 project-local `.agent-goals/`、locks、goal index、split / pause / complete | `--project`、`init/status/start/update/split/pause/complete/cleanup` | 寫入 `.agent-goals/`，更新 `.git/info/exclude`，建立 / 清理 lock | POSIX shell、Python 3、可選 Git | `ai-skill goals` | `native target` | fake project fixture；lock fixture；`complete` 缺 `--validated` 不刪除 assertion（Phase 2 已覆蓋 `status` read-only 與 `init --dry-run` planner；write/start/update/split/pause/complete/cleanup 待 parity） |
| `scripts/install-hooks.sh` | 從 `.githooks/` 複製 hooks 到 `.git/hooks/` 並加執行權限 | 無正式 flags | 寫入 `.git/hooks/`，設定可執行權限 | POSIX shell、Git repo、檔案權限 | `ai-skill hooks install` using `scripts/git-hooks/` source | `deleted` | Deleted with `.githooks/`; Go CLI dry-run planner now points at `scripts/git-hooks/`. Write copy / chmod still blocked until parity fixtures are complete. |
| `scripts/sync-cursor-bundle.sh` | 同步 Cursor bundle / mirror symlink，清理會造成循環的 symlink | 環境中的 `$HOME` 與 repo root | 寫入 `~/.cursor/bundles/`、`~/.cursor/shared-rules`、`~/.cursor/skills/`；可能移動既有非 symlink 目標 | POSIX shell、find、ln、檔案權限、Cursor 目錄慣例 | `ai-skill sync-cursor-bundle` | `tool-specific` | fake home fixture；managed / unmanaged mirror safety fixture；reference-only 不自動同步 assertion（Phase 2 已覆蓋 explicit target dry-run、target outside repo、copy-fallback strategy、無寫入；managed mirror write mode 待 parity） |
| `scripts/ai-skill-close-loop.sh` | 檢查 lock、dirty owner group、plan closure、knowledge runtime validation，依 group commit，可選 push | `--commit`、`--push`、`AI_SKILL_SYNC_CURSOR_BUNDLE`、`AI_SKILL_LOCK_TTL_SECONDS` | dry-run 列群組；`--commit` 執行 `git add/commit`；`--push` 執行 `git push`；可選同步 Cursor bundle | POSIX shell、Git、repo-local `ai-skill` binary、可選 Cursor sync | `ai-skill close-loop` | `wrapper first` | clean / dirty owner group fixture；merge / rebase / lock fixture；missing Git fixture；private path scan fixture（Phase 2 已覆蓋 native dry-run inspection；commit/push、plan closure、private path scan 待 parity；runtime validation 已改呼叫 repo-local CLI） |
| `scripts/refresh-knowledge-runtime.rb` | 串接 model report、checklist、runtime report、SQLite index 生成與 validation | 無 flags | 更新 generated Markdown reports 與本機 SQLite index；執行 validators | Ruby、sqlite3 CLI、canonical knowledge files | `ai-skill runtime refresh` | `deleted` | Native refresh default 寫 reports / index 並跑 native checks；legacy wrapper 已刪除。 |
| `scripts/generate-knowledge-runtime-report.rb` | 從 routing registry、summaries、graphs、refresh policy 產生 deterministic runtime report | `--write` | 可輸出或寫入 `knowledge/runtime/runtime-report.md` | Ruby、YAML parser | `ai-skill runtime refresh` | `deleted` | Native builder 與 golden anchors 覆蓋；runtime refresh 寫入 generated report。 |
| `scripts/generate-model-context-report.rb` | 從 routing registry model 欄位產生 model-aware loading report | `--write` | 可輸出或寫入 `knowledge/runtime/model-context-report.md` | Ruby、YAML parser | `ai-skill runtime refresh` | `deleted` | Native builder 與 golden anchors 覆蓋；runtime refresh 寫入 generated report。 |
| `scripts/generate-model-checklists.rb` | 從 routing registry 產生 per-model context-loading checklist | `--write` | 可輸出或寫入 `knowledge/runtime/model-checklists.md` | Ruby、YAML parser | `ai-skill runtime refresh` | `deleted` | Native builder 與 golden anchors 覆蓋；runtime refresh 寫入 generated report。 |
| `scripts/generate-runtime-sqlite-index.rb` | 產生本機 SQLite / FTS runtime lookup cache | `--output` 或預設路徑 | 寫入 git-ignored `knowledge/runtime/sqlite/runtime-index.sqlite` 或指定 DB | Ruby、sqlite3 CLI、canonical knowledge files | `ai-skill runtime refresh` | `deleted` | Native index builder 覆蓋 row counts、FTS hits、checksums 與 recursive feedback fixture；runtime refresh 寫入 SQLite index。 |
| `scripts/validate-runtime-sqlite-index.rb` | 驗證本機 SQLite runtime index integrity、row count、FTS、checksum、git ignore | 無 flags | 無寫入；失敗時 non-zero exit | Ruby、sqlite3 CLI、Git | `ai-skill runtime validate` | `deleted` | Native SQLite index validation 覆蓋 missing DB / table、integrity、row counts、atom source references、checksum、FTS count、ranked query 與 git-ignore boundary。 |
| `scripts/query-runtime-index.rb` | 用 keyword / layer / type / status 查詢 SQLite runtime index | query、`--limit`、`--db`、`--layer`、`--type`、`--status` | 讀取 SQLite，輸出 ranked candidates；無寫入 | Ruby、sqlite3 CLI | `ai-skill runtime query` | `deleted` | Native query fixture 覆蓋 keyword / positional query、filter、limit、empty result、missing DB；不依賴外部 `sqlite3` CLI。 |
| `scripts/query-knowledge-graph.rb` | 查詢 graph edge，支援 source / target / type / query / limit | `--source`、`--target`、`--type`、`--query`、`--limit` | 讀取 graph YAML，輸出 edge records；無寫入 | Ruby、YAML parser | `ai-skill runtime query --graph` | `deleted` | Native graph query fixture 覆蓋 source / target / type / keyword filter、limit、empty result 與 missing filter。 |
| `runtime/compiler/compiler-engine.rb` | 編譯 canonical prose source / structured YAML 到 `runtime/runtime.db` | `--check`、`--diff`、`--db` | 更新或檢查 `runtime/runtime.db`，可輸出 diff | Ruby、YAML parser、sqlite3 CLI、runtime source | `ai-skill runtime compile` | `deleted` | Replaced by Go-native source-to-DB compiler reading runtime YAML plus `runtime/compiler/compiler-rules.yaml`; custom DB, generated surface, row-count, and stable snapshot tests covered. |
| `scripts/validate-runtime-db.rb` | 驗證 `runtime/runtime.db` integrity、required tables、row counts、JSON columns、compiler metadata | `--db` | 無寫入；失敗時 non-zero exit | Ruby、sqlite3 CLI | `ai-skill runtime validate` | `deleted` | Native runtime DB validation 覆蓋 integrity、required tables、row counts、JSON、compiler metadata 與 stale metadata warning。 |
| `scripts/migrate-runtime-config-to-sqlite.rb` | 將 runtime YAML config 遷移到 `runtime/runtime.db` tables | 無 flags | 寫入 `runtime/runtime.db`；idempotent migration | Ruby、sqlite3 CLI、runtime YAML | `ai-skill runtime compile` future source-to-DB compiler | `deleted` | Deleted as obsolete migration helper; current supported path is `ai-skill runtime compile` / pre-commit compiler integration. |
| `scripts/init-runtime-state-db.rb` | 建立 mutable `runtime/runtime-state.db` schema | `--db` | 建立 / 更新 runtime-state SQLite schema；idempotent | Ruby、sqlite3 CLI | future `ai-skill runtime state init` | `deleted` | Deleted until mutable runtime-state scope has a Go-native command and validation contract. |
| `scripts/sync-runtime-yaml-from-embedded.rb` | 從 embedded runtime data 回寫 runtime YAML 與缺失 README | 無 flags | 寫入多個 `runtime/**/*.yaml` 與部分 README；一方向同步 | Ruby、runtime compiler embedded data | source restoration completed | `deleted` | Runtime YAML sources restored; no embedded runtime source remains. |
| `scripts/set-roo-global-custom-instructions.py` | 修改 VS Code globalStorage SQLite 中 Roo Code customInstructions | 無 flags；依固定 DB path | 寫入使用者 VS Code `state.vscdb`；要求 VS Code 關閉 | Python 3、sqlite3、macOS VS Code path、Roo Code key | `ai-skill roo set-global-custom-instructions` | `deleted` | Go command preserves guarded write behavior; fake VS Code DB tests cover success, missing DB, and missing Roo row. |
| `scripts/git-hooks/pre-commit` | staged runtime source 變更時編譯 `runtime.db`；staged knowledge / validation 變更時跑 repo-local CLI validator | Git staged files、`SKIP_*` 環境變數目前只在提示中出現 | 可更新並 stage `runtime/runtime.db`；可阻斷 commit | Git hook、repo-local `ai-skill` binary | `ai-skill hooks install` + hook adapter | `hook adapter` | staged runtime fixture；compiler failure blocks commit；runtime.db staged assertion；knowledge validation path uses repo-local binary |
| `scripts/git-hooks/post-commit` | 當 `AI_SKILL_SYNC_CURSOR_BUNDLE=1` 時執行 Cursor bundle sync | `AI_SKILL_SYNC_CURSOR_BUNDLE` | 可寫入 Cursor bundle / mirror；failure 不阻斷 commit | Git hook、POSIX shell、Cursor sync script | `ai-skill hooks install` + hook adapter | `hook adapter` | env-off no-op fixture；env-on fake home sync fixture |

## Parity 驗收矩陣

| 舊能力類別 | 必須證明的 parity | 對應命令 | 最低測試證據 |
| --- | --- | --- | --- |
| 新專案 bootstrap | dry-run、既有檔案衝突、指定工具、force 覆蓋、`.agent-goals/` 初始化都與舊腳本等價或有明確差異 | `ai-skill init-project` | fake project fixture、golden planned actions |
| Goal ledger | status / init / start / update / split / pause / complete / cleanup 的狀態轉換與 lock 行為可重現 | `ai-skill goals` | fake project fixture、lock fixture、validated completion fixture |
| Git hook 安裝 | hook 來源、目標、權限與 unsafe repo 報告可驗證 | `ai-skill hooks install` | temporary git repo fixture |
| Cursor mirror sync | reference-only 不自動同步；啟用時只管理預期 bundle path，不碰真實 home | `ai-skill sync-cursor-bundle` | fake home fixture、managed mirror fixture |
| Close-loop | lock、merge/rebase/cherry-pick、dirty owner group、private path scan、commit / push gating 全部可驗證 | `ai-skill close-loop` | temp repo dirty / unsafe / missing Git fixtures |
| Runtime refresh | 每個 generated surface 有完整更新或整體失敗，不允許 partial success | `ai-skill runtime refresh` | generated report golden fixture、partial failure fixture |
| Runtime compile | source change 能進 `runtime/runtime.db`，custom DB path 可隔離測試 | `ai-skill runtime compile` | runtime source-to-DB compiler fixture、custom DB fixture、stable snapshot fixture |
| Runtime validate | runtime.db 與 runtime index 的 schema、row count、checksum、FTS 與 generated surface freshness 可驗證 | `ai-skill runtime validate` | invalid DB / stale checksum fixtures |
| Runtime query | keyword、filter、limit、empty result、missing DB 行為穩定 | `ai-skill runtime query` | query fixture、empty result fixture（SQLite runtime index 與 knowledge graph query native slices 已覆蓋） |
| Tool-specific global settings | 不作為通用 CLI 預設；需 guarded adapter 並使用 fake DB 測試 | `ai-skill roo set-global-custom-instructions` | fake VS Code DB fixture |

## Phase Gate

在新增 `scripts/ai-skill-cli/go.mod`、`scripts/ai-skill-cli/cmd/ai-skill/` 或 production Go 實作前，必須滿足：

- [x] 本文件每個舊入口都有狀態與目標命令。
- [x] 每個 `native target` 或 `wrapper first` 舊入口都有最低測試證據或 fixture 名稱。
- [x] 高風險路徑已在 [`bdd-scenarios.md`](bdd-scenarios.md) 或 [`test-fixture-plan.md`](test-fixture-plan.md) 中出現。
- [x] 若某舊能力被 `deferred` 或 `tool-specific`，已說明為何不阻擋 Phase 1。
- [x] `command-contract.md` 只保留摘要表；完整 parity 以本文件為 source-of-truth。

## Phase 6 刪除 Gate

功能完成後，每個舊入口必須更新最終 disposition。當前 closure source-of-truth 是 [`legacy-script-disposition.md`](legacy-script-disposition.md)：

- `native target` / `wrapper first` 完成 parity 並由新 CLI 覆蓋後，舊 script 必須刪除並改標 `deleted`。
- `hook adapter` 可保留 hook surface，但安裝、檢查與文件應由 `ai-skill hooks install` 管理；保留原因必須寫清楚。
- `tool-specific` adapter 若仍需要，必須留在工具 adapter 層，不得成為通用 CLI 預設；保留時寫明移除條件。
- 若短期保留 thin wrapper，wrapper 只能轉呼叫 `ai-skill`，且必須在同一 closure 記錄刪除條件與期限。
