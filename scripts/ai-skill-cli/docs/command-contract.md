# 命令契約：`ai-skill` CLI

> **上游計畫**：[`2026-05-21-0834-cross-platform-go-script-runtime.md`](../../../plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md)

## 契約原則

- **Dry-run 優先**：任何寫檔、安裝 hook、同步 bundle、commit、push、migration、repair 的命令都必須支援 `--dry-run`。
- **機器可讀輸出**：所有命令必須支援 `--json`，且 schema 穩定。
- **人類可讀輸出**：所有命令必須支援 `--plain` 或預設人類可讀輸出，不依賴顏色。
- **穩定 exit code**：相同失敗類型在所有 OS 上回傳相同 exit code。
- **明確 side effect**：每個命令必須列出讀取路徑、寫入路徑、外部依賴與 git 操作。
- **禁止 partial close-loop**：linked-update、writeback、commit、push、runtime sync 若缺 Git 或 repo 狀態不安全，必須阻斷。
- **Native default**：Runtime core commands 預設走 Go-native path；Ruby、Python 與外部 `sqlite3` CLI 不屬於 active runtime dependency。

## 初始命令範圍

| 命令 | 目的 | 會修改檔案 | 需要 Git | 主要 Phase |
| --- | --- | --- | --- | --- |
| `ai-skill version` | 輸出 binary build metadata | 否 | 否 | Phase 4 |
| `ai-skill doctor` | 檢查 runtime、Git、repo root、PATH、write permission、hooksPath 與平台支援 | 否 | 可選檢查 | Phase 1 |
| `ai-skill init-project` | 建立新專案 AI tool bootstrap 設定 | 是 | 否 | Phase 2 |
| `ai-skill goals` | 管理 `.agent-goals/` 暫存目標 | 是 | 否 | Phase 2 |
| `ai-skill hooks install` | 安裝本 repo git hooks | 是 | 是 | Phase 2 |
| `ai-skill sync-cursor-bundle` | 同步 Cursor bundle / mirror | 是 | 否 | Phase 2 |
| `ai-skill close-loop` | 檢查 dirty owner group、commit、push、readback | 是 | 是 | Phase 2 |
| `ai-skill runtime refresh` | 重建 knowledge runtime reports / SQLite index|`ai-skill runtime refresh` | 重建 knowledge runtime reports / SQLite index | 是 | 否 | Phase 3 |
| `ai-skill runtime compile` | 編譯 `runtime/runtime.db` | 是 | 否 | Phase 3 |
| `ai-skill runtime validate` | 驗證 runtime.db、knowledge runtime、SQLite assertions | 否 | 否 | Phase 3 |
| `ai-skill runtime query` | 查詢 runtime index / generated surfaces | 否 | 否 | Phase 3 |
| `ai-skill roo set-global-custom-instructions` | guarded 寫入 Roo Code 全域 Custom Instructions | 是 | 否 | Tool adapter |

## 共通輸出契約

所有命令的 `--json` 輸出必須至少包含：

| 欄位 | 意義 |
| --- | --- |
| `command` | 實際執行的 command path，例如 `runtime compile` |
| `mode` | `dry_run`、`check`、`write`、`commit` 或 `push` |
| `status` | `success`、`blocked`、`failed` |
| `exit_code` | 對應 [Exit Code 表](#exit-code-表) 的穩定代碼 |
| `checks` | 已執行的 dependency / repo / platform / permission checks |
| `planned_actions` | dry-run 或 blocked 時原本會做的動作 |
| `mutations` | 實際寫入、git index、commit、push 或 generated artifact 變更；無則為空陣列 |
| `error` | 失敗時包含 `code`、`message`、`remediation` |

`--plain` 或預設輸出必須使用相同語意，只是以人類可讀段落呈現；不得只靠顏色或 terminal control sequence 表示 pass / fail。

## 命令細節

### `ai-skill version`

目的：輸出 release artifact 的 version、commit 與 build date。

輸入：

- `--json`
- `--plain`

副作用：無。

必要行為：

- 支援透過 Go `-ldflags` 注入 `Version`、`Commit`、`Date`。
- JSON output 必須包含同一組 build metadata checks，供 CI / release verification 解析。
- Repo-local binaries under `scripts/ai-skill-cli/bin/` must report the commit they were built from; `go test ./...` verifies `BUILDINFO`, `SHA256SUMS`, and the current host binary smoke. Rebuild them only when CLI source changes.

### `ai-skill doctor`

目的：在執行其他命令前檢查本機與 repository 是否可安全運行。

輸入：

- `--json`
- `--plain`
- `--require-git`
- `--require-write`
- `--check-runtime`

副作用：無。

必要檢查：

- Repo root 是否存在且可讀。
- Git 是否存在、版本是否符合最低需求、是否可執行 `git rev-parse` / `git status`。
- SQLite native path 是否可用；Phase 1 後應使用 pure Go SQLite，不依賴 `sqlite3` CLI。
- PATH、write permission、hooksPath、平台支援狀態。
- `--check-runtime` 必須使用 native SQLite driver 做 smoke query；若找到 `runtime.db`，執行 integrity check。

缺 Git 行為：

- 若只是一般 `doctor`：回報 `git.status = missing`，exit code 依嚴重度設定。
- 若傳入 `--require-git`：必須阻斷，提示安裝 Git。

### `ai-skill init-project`

目的：替代 `scripts/init-new-project.sh`，在目標專案建立 AI tool bootstrap 設定與 `.agent-goals/` 結構。

輸入：

- `--project <path>`
- `--tools <list>`
- `--dry-run`
- `--force`
- `--json`
- `--plain`

副作用：

- dry-run：只列出將寫入的檔案。
- 寫入模式：可能建立 `.roomodes`、`.cursor/rules/`、`.cursor/hooks.json`、`CLAUDE.md`、`.agent-goals/` 或等效 project-local 設定。
- Phase 2 初始切片只開放 dry-run planner；write mode 在 template parity、fixture 與覆蓋策略完成前必須回傳 `partial_close_loop_blocked`。

必要行為：

- 不寫入使用者真實 home 或 tool mirror，除非使用者提供明確目標路徑。
- 若目標檔存在且未傳 `--force`，必須阻斷並列出衝突檔案。
- 所有 template 中的 Ai-skill reference 必須指向 canonical repo 或可攜 reference；不得寫入一次性本機私有路徑到可重用文件。

### `ai-skill goals`

目的：替代 `scripts/agent-goals.sh`，管理 project-local `.agent-goals/` 暫存目標。

輸入：

- `status`
- `init`
- `start`
- `update`
- `split`
- `pause`
- `complete --validated`
- `cleanup`
- `--project <path>`
- `--json`
- `--plain`

副作用：

- 讀寫 `<PROJECT_ROOT>/.agent-goals/`。
- 建立或清理 lock directory。
- Phase 2 初始切片只開放 `status` read-only 與 `init --dry-run` planner；write mode 在 fixture parity 完成前必須回傳 `partial_close_loop_blocked`。

必要行為：

- `.agent-goals/` 預設視為 temporary project state，不應加入 git index。
- `complete` 未傳 `--validated` 時不得刪除 goal file，只能標記 needs-validation。
- 若 active lock 屬於其他 owner 且未判定 stale，不得覆寫。

### `ai-skill hooks install`

目的：替代已刪除的 `scripts/install-hooks.sh` / `.githooks/` 舊面或手動 hook 設定，安裝本 repo git hooks。

輸入：

- `--dry-run`
- `--repo <path>`
- `--force`
- `--json`
- `--plain`

副作用：

- dry-run：列出會設定的 hooks path。
- 寫入模式：可能修改 repo-local git config 或 hooks path。
- Phase 2 初始切片只開放 dry-run planner；write mode 在 hook copy / chmod parity 完成前必須回傳 `partial_close_loop_blocked`。

必要 Git 檢查：

- `git` binary 存在。
- 目標 repo 可由 `git rev-parse --show-toplevel` 確認。
- hook source 是 `scripts/git-hooks/`；target 是 repo-local `.git/hooks/`。
- 若 repo 正在 merge / rebase / cherry-pick，安裝 hook 仍可被允許，但必須明確報告目前不安全狀態，且不得觸發 commit / push。

### `ai-skill sync-cursor-bundle`

目的：替代 `scripts/sync-cursor-bundle.sh`，在使用者明確選擇 Cursor bundle / mirror 策略時同步本庫到工具 mirror。

輸入：

- `--dry-run`
- `--repo <path>`
- `--target <path>`
- `--json`
- `--plain`

副作用：

- dry-run：列出會同步的來源與目標。
- 寫入模式：可能建立、更新或刪除目標 mirror 內的受管理檔案。
- Phase 2 初始切片只開放 dry-run planner；write mode 在 managed mirror 與 symlink/copy fallback parity 完成前必須回傳 `partial_close_loop_blocked`。

必要行為：

- Reference-only 工作流不應自動執行此命令。
- `--target` 必須明確傳入 Cursor root / fake Cursor root，不得預設寫入真實 `$HOME`。
- 必須區分 managed mirror content 與使用者自訂檔案，避免刪除非本工具管理的內容。
- 權限不足時回傳 `permission_denied`，不得半同步後宣稱完成。
- Windows / 權限受限環境的預設策略是 copy fallback；symlink 只能作為未來明確 opt-in 並通過權限 fixture 後啟用。

### `ai-skill close-loop`

目的：替代 `scripts/ai-skill-close-loop.sh`，保守處理 writeback transaction、dirty owner group、commit、push。

輸入：

- `--dry-run`
- `--repo <path>`
- `--commit`
- `--push`
- `--json`
- `--plain`

副作用：

- dry-run：無寫入。
- `--commit`：可能執行 `git add`、`git commit`。
- `--push`：可能執行 `git push`。
- Phase 2 初始切片只開放 dry-run inspection；`--commit` / `--push` 在 commit parity 完成前必須回傳 `partial_close_loop_blocked`。

必要 Git 檢查：

- `git` binary 存在。
- repo root 可由 `git rev-parse --show-toplevel` 確認。
- 不在 merge / rebase / cherry-pick 狀態。
- 不存在 active `.git/ai-skill-agent.lock`。
- dirty files 可歸屬 owner group。

缺 Git 行為：

- 不得 fallback 成手動檔案掃描後繼續。
- 必須阻斷並提示安裝 Git。
- JSON 輸出必須包含 `error.code = "missing_git"`。

### `ai-skill runtime refresh`

目的：重建 knowledge runtime reports、model reports 與 SQLite index；已取代並刪除舊 `scripts/refresh-knowledge-runtime.rb`。

輸入：

- `--dry-run`
- `--repo <path>`
- `--json`
- `--plain`
- `--native-reports`
- `--assert-source <path>`

副作用：

- 可能更新 `knowledge/runtime/runtime-report.md`、`knowledge/runtime/model-context-report.md`、`knowledge/runtime/model-checklists.md` 與本機 SQLite index。
- 預設 native mode 不依賴 Ruby、Python 或外部 `sqlite3` CLI。

必要行為：

- 必須回報哪些 generated surfaces 被更新、哪些 validator 被執行。
- 不得只更新部分 generated surface 後回傳 success。
- 預設 native mode 寫入 Go-generated Markdown reports 與 SQLite index，並執行 native runtime DB / index / knowledge runtime checks；dry-run 只列出將執行的 native actions，不寫入 generated surfaces。
- 任一 native refresh step 失敗時，CLI 必須停止後續 steps、回傳 `runtime_refresh_failed`，且 JSON checks 必須保留已執行 steps 與 failing step。
- `--native-reports` / `--native-index` 已是預設 native refresh 行為；保留 flags 供舊 automation 顯式表示意圖。

### `ai-skill runtime compile`

目的：用 Go-native compiler 從 runtime YAML 與 deterministic prose mappings 編譯 `runtime/runtime.db`。

輸入：

- `--dry-run`
- `--repo <path>`
- `--db <path>`
- `--native-compiler`（deprecated no-op；compile 已是 native）
- `--assert-source <path>`
- `--assert-keyword <keyword>`
- `--json`

副作用：

- 可能更新 `runtime/runtime.db`。
- 預設 native source-to-DB mode 不依賴 Ruby、Python 或外部 `sqlite3` CLI。

驗證：

- `runtime/runtime.db` integrity check。
- `generated_surfaces` content assertion。
- compiler version / schema version 存在。
- Go compiler 讀取 `runtime/compiler/compiler-rules.yaml`、runtime YAML source、以及 deterministic prose sources，產生指定 `--db` output，並通過 native runtime DB validation。

### `ai-skill runtime validate`

目的：執行 runtime.db、knowledge runtime、SQLite index 與 assertion 驗證。

輸入：

- `--dry-run`
- `--repo <path>`
- `--json`
- `--plain`
- `--source <path>`
- `--keyword <keyword>`

副作用：無。

必要行為：

- 驗證失敗時不得修改檔案。
- 必須區分 missing dependency、schema invalid、assertion failed、dirty generated surface。
- 預設 native mode 執行 Go runtime DB、SQLite index 與 knowledge runtime checks，不依賴 Ruby、Python 或外部 `sqlite3` CLI。
- `runtime.db` native slice 已用 Go / `modernc.org/sqlite` 檢查 integrity、required tables、minimum row counts、JSON columns、compiler metadata 與 stale metadata warning；stale warning 不阻斷成功狀態。
- SQLite runtime index native slice 已用 Go / `modernc.org/sqlite` 檢查 missing DB、integrity、required tables、row counts、atom source references、source checksums、FTS count 與 basic ranked FTS query；git-ignore boundary 以 Go 呼叫 Git 檢查，缺 Git 時回 `missing_dependency`。

### `ai-skill runtime query`

目的：查詢 runtime SQLite index、knowledge graph 與 runtime DB；已取代並刪除舊 runtime query Ruby helpers。

輸入：

- `<query>` positional term 或 `--keyword <term>`
- `--graph`
- `--db <path>`
- `--layer <name>`
- `--source <path>`
- `--target <path>`
- `--type <name>`
- `--status <name>`
- `--limit <n>`
- `--json`
- `--plain`

副作用：無。

必要行為：

- Query 命令不得修改 SQLite DB 或 generated reports。
- 查不到資料時回傳 success 並提供空 results，除非 requested table / DB schema 不存在。
- `--json` results 必須包含 source path、rank / priority（若有）、match reason 與 validation signal（若有）。
- Native query 覆蓋舊 runtime index query：keyword / positional query、`--db`、`--layer`、`--type`、`--status`、`--limit`、empty result 與 missing DB。
- `--graph` native query 覆蓋舊 knowledge graph query：`--source`、`--target`、`--type`、`--keyword` / positional query、`--limit`、empty result 與 missing filter。

## 舊 Script 對應

完整功能盤點、side effects、外部依賴與測試證據見 [`script-parity-inventory.md`](script-parity-inventory.md)。本節只保留命令契約層的摘要映射；不得用本表取代 parity 驗收。

| 現有 script / 入口 | CLI 命令 | 遷移定位 |
| --- | --- | --- |
| `scripts/init-new-project.sh` | `ai-skill init-project` | Phase 2 native 候選 |
| `scripts/agent-goals.sh` | `ai-skill goals` | Phase 2 native 候選 |
| deleted `scripts/install-hooks.sh` / `.githooks/` | `ai-skill hooks install` | dry-run planner uses `scripts/git-hooks/`; write mode still blocked until fixture-backed |
| `scripts/sync-cursor-bundle.sh` | `ai-skill sync-cursor-bundle` | Phase 2 native 候選，需 mirror safety gate |
| `scripts/ai-skill-close-loop.sh` | `ai-skill close-loop` | Phase 2 先 wrapper，owner-group parity 後 native |
| Runtime report / SQLite generators | `ai-skill runtime refresh` | Native completed; old Ruby entrypoints deleted |
| Runtime validators | `ai-skill runtime validate` | Native completed; old Ruby entrypoints deleted |
| Runtime query helpers | `ai-skill runtime query` | Native completed; old Ruby entrypoints deleted |
| Runtime compiler | `ai-skill runtime compile` | Go-native source-to-DB compiler completed; old Ruby compiler deleted |
| Runtime migration / state helpers | future Go-native runtime commands | old Ruby helpers deleted; reintroduce only with command contract and fixtures |
| Roo global setting helper | `ai-skill roo set-global-custom-instructions` | guarded tool-specific adapter with fake DB tests |

## Exit Code 表

| 代碼 | 名稱 | 意義 |
| --- | --- | --- |
| `0` | success | 命令成功 |
| `1` | general_failure | 未分類錯誤，應逐步消除 |
| `2` | invalid_usage | 參數錯誤或互斥 flag |
| `10` | missing_dependency | 缺必要外部依賴，例如 Git |
| `11` | unsupported_platform | 平台不支援該命令 |
| `20` | unsafe_repo_state | merge / rebase / dirty owner group 不安全 |
| `21` | permission_denied | 權限不足 |
| `30` | validation_failed | validator 或 assertion 失敗 |
| `31` | generated_surface_stale | generated artifact 未同步 |
| `40` | partial_close_loop_blocked | 會造成 partial close-loop，因此被阻斷 |

## 舊 Script 刪除政策

本 CLI 是 replacement 入口，不是永久新增第二套工具。當某個舊 script 的功能已由 `ai-skill` 命令覆蓋，且 parity fixture、BDD scenario、文件與 release gate 通過後，舊 script 必須刪除；不得只標記 deprecated 後長期保留。

例外：

- Git hook adapter 可保留 hook surface，但安裝、檢查與文件應由 `ai-skill hooks install` 管理。
- Tool-specific adapter 可保留在工具層，但不得成為通用 CLI 預設。
- 短期 thin wrapper 只能轉呼叫 `ai-skill`，並必須記錄 owner、刪除條件與期限。

## 副作用登錄表

| 命令 | 讀取路徑 | 寫入路徑 | 外部依賴 |
| --- | --- | --- | --- |
| `doctor` | repo root、PATH、git config | 無 | 可選 Git |
| `init-project` | `CORE_BOOTSTRAP.md`、`workflow/documentation/`、tool templates | 目標專案設定、`.agent-goals/` | 無 |
| `goals` | `<PROJECT_ROOT>/.agent-goals/` | `<PROJECT_ROOT>/.agent-goals/` | 無 |
| `hooks install` | `scripts/git-hooks/`、git config | `.git/config` 或 hooks path | Git |
| `sync-cursor-bundle` | Ai-skill source | Cursor bundle / mirror path | filesystem permissions |
| `close-loop` | git status、repo files、rules | git index、commits、remote branch | Git |
| `runtime refresh` | `knowledge/`、`feedback/`、runtime sources | generated reports、SQLite index | 無 |
| `runtime compile` | runtime YAML、compiler rules、prose sources | `runtime/runtime.db` | 無 |
| `roo set-global-custom-instructions` | VS Code `state.vscdb`、optional instructions file | VS Code `state.vscdb` | 無 |
| `runtime validate` | generated reports、runtime.db | 無 | SQLite index git-ignore boundary 需 Git |
| `runtime query` | `knowledge/runtime/sqlite/runtime-index.sqlite` 或 `--db` 指定 SQLite index | 無 | 無 |
| `runtime query --graph` | `knowledge/graphs/*.yaml` | 無 | 無 |
