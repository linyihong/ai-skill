# Command Contract：`ai-skill` CLI

> **上游計畫**：[`2026-05-21-0834-cross-platform-go-script-runtime.md`](../../../plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md)

## Contract Principles

- **Dry-run first**：任何寫檔、安裝 hook、同步 bundle、commit、push、migration、repair 的命令都必須支援 `--dry-run`。
- **Machine-readable output**：所有命令必須支援 `--json`，且 schema 穩定。
- **Human-readable output**：所有命令必須支援 `--plain` 或預設人類可讀輸出，不依賴顏色。
- **Stable exit codes**：相同失敗類型在所有 OS 上回傳相同 exit code。
- **Explicit side effects**：每個命令必須列出讀取路徑、寫入路徑、外部依賴與 git 操作。
- **No partial close-loop**：linked-update、writeback、commit、push、runtime sync 若缺 Git 或 repo 狀態不安全，必須阻斷。

## Initial Command Surface

| Command | Purpose | Mutates Files | Requires Git | Primary Phase |
| --- | --- | --- | --- | --- |
| `ai-skill doctor` | 檢查 runtime、Git、repo root、PATH、write permission、hooksPath 與平台支援 | no | optional check | Phase 1 |
| `ai-skill init-project` | 建立新專案 AI tool bootstrap 設定 | yes | no | Phase 2 |
| `ai-skill goals` | 管理 `.agent-goals/` 暫存目標 | yes | no | Phase 2 |
| `ai-skill hooks install` | 安裝本 repo git hooks | yes | yes | Phase 2 |
| `ai-skill sync-cursor-bundle` | 同步 Cursor bundle / mirror | yes | no | Phase 2 |
| `ai-skill close-loop` | 檢查 dirty owner group、commit、push、readback | yes | yes | Phase 2 |
| `ai-skill runtime refresh` | 重建 knowledge runtime reports / SQLite index | yes | no | Phase 3 |
| `ai-skill runtime compile` | 編譯 `runtime/runtime.db` | yes | no | Phase 3 |
| `ai-skill runtime validate` | 驗證 runtime.db、knowledge runtime、SQLite assertions | no | no | Phase 3 |
| `ai-skill runtime query` | 查詢 runtime index / generated surfaces | no | no | Phase 3 |

## Command Details

### `ai-skill doctor`

Purpose：在執行其他命令前檢查本機與 repository 是否可安全運行。

Inputs：

- `--json`
- `--plain`
- `--require-git`
- `--require-write`
- `--check-runtime`

Side effects：無。

Required checks：

- Repo root 是否存在且可讀。
- Git 是否存在、版本是否符合最低需求、是否可執行 `git rev-parse` / `git status`。
- SQLite native path 是否可用；Phase 1 後應使用 pure Go SQLite，不依賴 `sqlite3` CLI。
- Ruby / Python 僅在 wrapper mode 需要時提示，不應是長期核心依賴。
- PATH、write permission、hooksPath、平台支援狀態。

Missing Git behavior：

- 若只是一般 `doctor`：回報 `git.status = missing`，exit code 依 severity 設定。
- 若傳入 `--require-git`：必須阻斷，提示安裝 Git。

### `ai-skill close-loop`

Purpose：替代 `scripts/ai-skill-close-loop.sh`，保守處理 writeback transaction、dirty owner group、commit、push。

Inputs：

- `--dry-run`
- `--commit`
- `--push`
- `--json`
- `--plain`

Side effects：

- dry-run：無寫入。
- `--commit`：可能執行 `git add`、`git commit`。
- `--push`：可能執行 `git push`。

Required Git checks：

- `git` binary exists。
- repo root 可由 `git rev-parse --show-toplevel` 確認。
- 不在 merge / rebase / cherry-pick 狀態。
- dirty files 可歸屬 owner group。

Missing Git behavior：

- 不得 fallback 成手動檔案掃描後繼續。
- 必須阻斷並提示安裝 Git。
- JSON output 必須包含 `error.code = "missing_git"`。

### `ai-skill runtime compile`

Purpose：編譯 `runtime/runtime.db`，未來可從 Ruby wrapper 過渡到 Go native compiler。

Inputs：

- `--dry-run`
- `--assert-source <path>`
- `--assert-keyword <keyword>`
- `--json`

Side effects：

- 可能更新 `runtime/runtime.db`。
- wrapper mode 可能呼叫 Ruby compiler；native mode 不應依賴 Ruby。

Validation：

- `runtime/runtime.db` integrity check。
- `generated_surfaces` content assertion。
- compiler version / schema version 存在。

### `ai-skill runtime validate`

Purpose：執行 runtime.db、knowledge runtime、SQLite index 與 assertion 驗證。

Inputs：

- `--json`
- `--plain`
- `--source <path>`
- `--keyword <keyword>`

Side effects：無。

Required behavior：

- 驗證失敗時不得修改檔案。
- 必須區分 missing dependency、schema invalid、assertion failed、dirty generated surface。

## Exit Code Table

| Code | Name | Meaning |
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
| `40` | partial_close_loop_blocked | 會造成半套 close-loop，因此被阻斷 |

## Side-Effect Registry

| Command | Read Paths | Write Paths | External Dependencies |
| --- | --- | --- | --- |
| `doctor` | repo root、PATH、git config | none | optional Git |
| `init-project` | `CORE_BOOTSTRAP.md`、`workflow/documentation/`、tool templates | target project config、`.agent-goals/` | none |
| `goals` | `<PROJECT_ROOT>/.agent-goals/` | `<PROJECT_ROOT>/.agent-goals/` | none |
| `hooks install` | `scripts/git-hooks/`、git config | `.git/config` or hooks path | Git |
| `sync-cursor-bundle` | Ai-skill source | Cursor bundle / mirror path | filesystem permissions |
| `close-loop` | git status、repo files、rules | git index、commits、remote branch | Git |
| `runtime refresh` | `knowledge/`、`feedback/`、runtime sources | generated reports、SQLite index | wrapper mode may need Ruby |
| `runtime compile` | runtime compiler sources、prose sources | `runtime/runtime.db` | wrapper mode may need Ruby |
| `runtime validate` | generated reports、runtime.db | none | none after native migration |
