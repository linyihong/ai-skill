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
- **Go-first automation**：新增 repository automation 必須先設計成 `ai-skill` Go CLI command。`.sh` 只可作為 Git hook adapter、thin bootstrap wrapper，或保留到 Go write-mode parity 完成的 legacy surface。

## 初始命令範圍

| 命令 | 目的 | 會修改檔案 | 需要 Git | 主要 Phase |
| --- | --- | --- | --- | --- |
| `ai-skill version` | 輸出 binary build metadata | 否 | 否 | Phase 4 |
| `ai-skill doctor` | 檢查 runtime、Git、repo root、PATH、write permission、hooksPath 與平台支援 | 否 | 可選檢查 | Phase 1 |
| `ai-skill init-project` | 建立新專案 AI tool bootstrap 設定 | 是 | 否 | Phase 2 |
| `ai-skill goals` | 管理 `.agent-goals/` 暫存目標 | 是 | 否 | Phase 2 |
| `ai-skill hooks install` | 安裝本 repo git hooks（4 hooks：pre-commit / commit-msg / post-commit / pre-push） | 是 | 是 | Phase 2 |
| `ai-skill hooks run pre-commit` | 執行 Git pre-commit hook logic（runtime source 變更時自動 recompile + stage runtime.db） | 是 | 是 | Shell To Go |
| `ai-skill hooks run commit-msg` | 執行 Git commit-msg hook logic：commit body 必須含 Cognitive Contract v2 報告，並通過 behavioral validators（execution floors、governance consistency、memory subdir、cognitive cost、activation signals、capability snippet、plan sync、token budget、adaptive triggers、bootstrap thinness、CLI doc sync、runtime YAML projection、Markdown/YAML sync）。Compact form 只適用於 all-default 且 staged diff 未命中 runtime / routing / workflow executable contract / active plan 等會要求 DEEP+STRICT 的 surfaces；compact path 也必須跑 staged-file validators。 | 否 | 是 | ADR-008 / runtime-cognitive-contract-v2 |
| `ai-skill hooks run post-commit` | 執行 Git post-commit hook logic（cursor bundle sync 等） | 否 | 是 | Shell To Go |
| `ai-skill hooks run pre-push` | 執行 Git pre-push hook logic：CLI source 變動時 go test ./... preflight | 否 | 是 | Shell To Go |
| `ai-skill hooks run session-start` | 執行 Claude Code SessionStart hook：query runtime.db、讀 4 個 bootstrap 文件、輸出 hookSpecificOutput JSON、寫 SessionStart flag（TTL 120s） | 是（`/tmp/ai-skill-sessionstart-<hash>.flag`） | 否 | Cross-platform Go script runtime |
| `ai-skill hooks run pre-tool-use` | 執行 Claude Code PreToolUse hook：scan transcript for Bootstrap Receipt；Read tool 一律 allow，其他 tool 在 Receipt 出現前 block（exit 2） | 是（`/tmp/ai-skill-bootstrap-<hash>.done`） | 否 | Cross-platform Go script runtime |
| `ai-skill hooks run post-tool-use` | 執行 Claude Code PostToolUse hook：Bootstrap Receipt 不在 transcript 時注入 reminder via hookSpecificOutput；always exit 0 | 是（cache file） | 否 | Cross-platform Go script runtime |
| `ai-skill hooks run user-prompt-submit` | 執行 Claude Code UserPromptSubmit hook：每次 user turn 注入 final close-out reminder + CORE_BOOTSTRAP.md as additionalContext；若 project root 底下有 dirty nested Git repos，注入合併 `### Project Git Report` 要求 | 否 | 否 | Cross-platform Go script runtime |
| `ai-skill hooks run stop` | 執行 Stop / final-response hook：從 transcript 或 hook payload 檢查對話含 Bootstrap Receipt、last assistant message 含 Cognitive Mode block（compact 或 full table）；若有 dirty root / nested Git repos 也要求 `### Project Git Report`；Claude-style stop 缺少時 block（exit 2）；Cursor stop 缺少時一次彙整缺項輸出 `followup_message` 並 exit 0 以 loop back；Cursor plan / todo / mode-switch tool-generated 非 final 訊息不觸發 close-out loop | 否 | 否 | Cross-platform Go script runtime |
| `ai-skill sync-cursor-bundle` | 同步 Cursor bundle / mirror | 是 | 否 | Phase 2 |
| `ai-skill close-loop` | 檢查 dirty owner group、commit、push、readback | 是 | 是 | Phase 2 |
| `ai-skill runtime refresh` | 重建 knowledge runtime reports / SQLite index | 是 | 否 | Phase 3 |
| `ai-skill runtime compile` | 編譯 `runtime/runtime.db` | 是 | 否 | Phase 3 |
| `ai-skill runtime validate` | 驗證 runtime.db、knowledge runtime、SQLite assertions | 否 | 否 | Phase 3 |
| `ai-skill runtime query` | 查詢 runtime index / generated surfaces | 否 | 否 | Phase 3 |
| `ai-skill runtime obligations` | 列出目前 active bootstrap obligations（per_session / per_turn / per_commit）並附 Bootstrap Receipt line，從 `generated_surfaces[runtime.core_bootstrap.contract]` 與 runtime phase/gate tables 讀取 | 否 | 否 | bootstrap-yaml-migration Phase 3 |
| `ai-skill runtime receipt` | 輸出 canonical Bootstrap Receipt 與 active per-turn obligation IDs，供 hooks / stop repair 使用；避免 agent 臨時拼 SQLite 查詢 | 否 | 否 | bootstrap receipt repair hardening |
| `ai-skill runtime audit` | 4-way 分類 routes / generated_surfaces / scenarios（auto-detected / consumed / intentionally-manual / orphan）。預設 markdown 報告；`--json` 切換 JSON。`runtime validate` 自動以 warning-only check 引用其 orphan 統計 | 否 | 否 | gen3-runtime-trigger-audit Phase 2 |
| `ai-skill hooks run commit-msg` validator `validatePlanCheckboxSync` | commit-msg hook 第 16 個 validator：當 commit body 引用 `plans/active/*.md` 且 stage 真工作（Go / scenarios / runtime / governance / enforcement），plan 必須同 stage 且 staged diff 含 `[ ]` → `[x]` transition。block default；opt-out `[skip-plan-checkbox-sync]` | 否 | 是 | gen3-runtime-trigger-audit Phase 5 |
| `ai-skill hooks run commit-msg` validator `validateRuntimeTriggerWiring` | commit-msg hook 第 17 個 validator：staged diff 新增 `route.*` 或 `target_key:` 但無 discovery signal / Go consumer / `manual_activation` annotation 則 block。enforces governance §`define_runtime_trigger_flow`；opt-out `[skip-runtime-trigger-wiring]` | 否 | 是 | gen3-runtime-trigger-audit Phase 5 |
| `ai-skill hooks run commit-msg` validator `validatePlanArchivalAudit` | commit-msg hook 第 19 個 validator：當 commit 把 `plans/active/<name>.md` 移到 `plans/archived/<name>.md`（staged 同時含刪除與新增同 basename），archived 版本若仍有 `- [ ]` 行且 commit body 無 justification keyword（deferred / non-goal / scope reduced / handover / 延後 / 拆分）則 block。純 Go 掃描（呼叫 `ScanCheckboxesInFile`），零 shell 依賴；opt-out `[skip-plan-archival-audit]` | 否 | 是 | plan-archival-audit-validator Phase 2 |
| `ai-skill scan-checkboxes <file>` | 掃描任意 Markdown 檔案的 task-list checkboxes（`- [ ]` / `- [x]` / `- [X]`），純 Go 實作，零 shell 依賴，可跨專案使用（release checklist、規格確認單等）。flags：`--format plain\|json`、`--exit-code`（有未完成項時 exit 1，方便 CI / pre-push hooks 使用） | 否 | 否 | plan-archival-audit-validator Phase 2 |
| `ai-skill glossary validate` | 驗證 `knowledge/glossary/*.md` 的 entry schema、status / owner / relation enum、naming convention、alias 規則、`introduced-by` / `deprecated-by` 形狀、symmetric relation 對稱性與 `excludes` 引用 | 否 | 否 | context-language-glossary-system Phase 2 |
| `ai-skill enforcement lint` | 對 `enforcement/enforcement-registry.yaml` 跑 13 條 Phase 3 lint check（orphan_rule / missing_executor_symbol / behavioral_only_* / deprecated_* / upstream_chain / class_size / baseline_snapshot / pending_implementation_child_plan_validity）。Thin CLI wrapper 直接複用 `LintEnforcementRegistry` 引擎；支援 `--check <substr>` 過濾、`--registry <path>` 覆寫、`--expect-finding <substr>` 跨平台 assertion mode（scenario 用） | 否 | 否 | mechanical-enforcement-registry Phase 4 |
| `ai-skill enforcement coverage` | 把 registry 聚合成 6-bucket coverage report（mechanical / behavioral_only / not_mechanizable / pending_implementation / research_required / deprecated）+ verification level 分類 + runtime_observed gap（Phase 5 wire 前回 `null` + alert）。`--format text|json|markdown`、`--diff <ref>` 對比 git revision、`--detail` 開 per-class 表、`--self-check` 跨平台驗 3 個 format schema | 否 | 否 | mechanical-enforcement-registry Phase 4 |
| `ai-skill roo set-global-custom-instructions` | guarded 寫入 Roo Code 全域 Custom Instructions | 是 | 否 | Tool adapter |
| `ai-skill copilot start` | 產生 GitHub Copilot 新 session 第一則 bootstrap prompt | 否 | 否 | Tool adapter |

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
- 寫入模式：可能建立 `.roomodes`、`.cursor/rules/`、`.cursor/hooks.json`、`CLAUDE.md`、`.claude/settings.json`、`.agent-goals/`、`.ai-skill/.gitignore`、`.ai-skill/local.env` 或等效 project-local 設定。
- 寫入內容不得包含本機 Ai-skill repository 的絕對路徑；以 `<AI_SKILL_REPO>` placeholder 與 `AI_SKILL_REPO` 環境變數連回 canonical repo。
- `.ai-skill/local.env` 是唯一例外：它是 machine-local bridge file，初始化時寫入目前 Ai-skill repo path，權限為 `0600`，且必須由 `.ai-skill/.gitignore` 排除，避免一般 `git add .` commit 本機路徑。
- Phase 2 初始切片只開放 dry-run planner；write mode 在 template parity、fixture 與覆蓋策略完成前必須回傳 `partial_close_loop_blocked`。
- Shell To Go migration 後，write mode 必須實作並成為預設；舊 `scripts/init-new-project.sh` 必須刪除。

必要行為：

- 不寫入使用者真實 home 或 tool mirror，除非使用者提供明確目標路徑。
- 若目標檔存在且未傳 `--force`，必須阻斷並列出衝突檔案。
- 所有 template 中的 Ai-skill reference 必須指向 canonical repo 或可攜 reference；不得寫入一次性本機私有路徑到可重用文件。

### `ai-skill copilot start`

目的：替 GitHub Copilot 新 session 產生第一則 guided bootstrap prompt，降低 Copilot 忽略 project instructions 或 scoped instructions 的機率。

輸入：

- `--project <path>`（預設目前目錄）
- `--json`
- `--plain`

副作用：無。此命令只讀取 project path 與 `.ai-skill/local.env` 是否存在，不修改檔案、不啟動外部 editor、不執行 Git。

必要行為：

- prompt 必須指向 `<AI_SKILL_REPO>/CORE_BOOTSTRAP.md`、`<AI_SKILL_REPO>/runtime/core-bootstrap.yaml` 與 `<AI_SKILL_REPO>/ai-tools/agent/copilot.md`。
- prompt 必須要求 Copilot 在回覆任何使用者請求前完成 bootstrap，不得把簡單列檔、read-only query、說明原因或 resumed context 視為豁免。
- prompt 必須明確說明 Copilot instructions 是 guided startup，不是 hard enforcement；runtime gates 仍由 hooks、CI、`ai-skill runtime validate` 負責。
- `init-project --tools copilot` 可產生 `.copilot/start-copilot.sh` 作為 temporary thin wrapper；wrapper 只能呼叫 repo-local `ai-skill copilot start`，並必須寫明 deletion/removal condition。

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
- Shell To Go migration 後，`init/status/start/update/split/pause/complete/cleanup` 都必須由 Go 原生實作；舊 `scripts/agent-goals.sh` 必須刪除。

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
- Shell To Go migration 後，hook installation 必須安裝 Go-owned hook runner；hook file 若保留，只能作最小 binary adapter。

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
- Shell To Go migration 後，`--commit` / `--push` 必須由 Go 原生執行；舊 `scripts/ai-skill-close-loop.sh` 必須刪除。

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

### `ai-skill hooks run`

目的：讓 Git hook adapter 將業務邏輯交給 Go，避免 hook 依賴 shell 實作。

輸入：

- Git hooks: `pre-commit` / `commit-msg` / `post-commit` / `pre-push`
- Claude Code hooks: `session-start` / `user-prompt-submit` / `pre-tool-use` / `post-tool-use` / `stop`
- `--repo <path>`
- `--json` / `--plain`
- positional：commit-msg 收 `<commit-msg-file-path>`（由 git 透過 `"$@"` 從 adapter 傳入）

副作用：

- `pre-commit`：當 staged `runtime/runtime.db` 或 knowledge / validation surface 改動時執行 `runtime validate`；不再因 committed runtime YAML mirror 觸發 compile。新增 shell script gate：`validateNoNewShellScripts` 偵測 staged Added（`--diff-filter=A`）`.sh` 檔案，發現新建 `.sh` 時以 `new_shell_script_forbidden` block；修改既有 `.sh` 不觸發。Opt-out：commit message 含獨立一行 `[skip-go-migration]`（僅用於短暫過渡 wrapper，必須在 commit message 說明 Go migration plan）。
- `commit-msg`：讀 commit message file，依 `runtime/cognitive-modes-*.yaml` + `runtime/cognitive-modes-cost-class.yaml` + `runtime/plan-status-sync-enforcement.yaml` + `runtime/cognitive-modes-token-budget.yaml` + `runtime/cognitive-modes-adaptive.yaml` + `runtime/bootstrap-entry-points.yaml` + `runtime/cli-modification-policy.yaml` + `runtime/core-bootstrap.yaml` contracts 執行 15 個 validators（cognitive contract block / executionFloors / governanceConsistency / memorySubdir / cognitiveCost / activationSignals / capabilitySnippet / planStatusSync / tokenBudget / adaptiveTriggers / bootstrapEntryThinness / **cliDocSync** / **runtimeYamlProjects** / **markdownYamlSync** / **glossaryRetroOwn**）。Validator block 時 exit 30。Opt-out trailers：`[skip-cognitive-mode]` / `[skip-plan-status-sync]` / `[skip-token-budget]` / `[skip-adaptive]` / `[skip-bootstrap-thinness]` / `[skip-cli-doc-sync]` / `[skip-runtime-yaml-projection]` / `[skip-markdown-yaml-sync]` / `[skip-glossary-retro-own]`。

  **`glossaryRetroOwn` validator**（Phase 6 of context-language-glossary-system）：staged diff 動到 framework cognitive vocabulary surface（`runtime/cognitive-modes*.yaml`、`runtime/economics/**`、`ecosystem/**`）時，`knowledge/glossary/ai-skill.md` 必須同時 stage，避免新 framework term 漂移為 subsystem-local vocabulary。Opt-out `[skip-glossary-retro-own]` 適用於 typo / refactor / comment-only 變更。詳見 [`runtime/cli-modification-policy.yaml`](../../../runtime/cli-modification-policy.yaml) §`gate.glossary.retro_own_required` 與 [`validation/scenarios/failure-derived/glossary-retro-own-missing-v1.yaml`](../../../validation/scenarios/failure-derived/glossary-retro-own-missing-v1.yaml)。
- `post-commit`：reference-only 預設 no-op；若 `AI_SKILL_SYNC_CURSOR_BUNDLE=1`，只回報 Go mirror write mode 狀態。
- `pre-push`：CLI source（`scripts/ai-skill-cli/...`、GitHub workflows、Git hooks）變動時跑 `go test ./...` preflight；其他情況跳過。
- `user-prompt-submit`：Claude Code hook；注入 final close-out Cognitive Mode reminder，並掃描 project root 底下的 dirty nested Git repositories。若一個 repo dirty，final response 應回報該 repo；若多個 repo dirty，合併成單一 `### Project Git Report` section。
- `stop`：Claude Code hook；final response 必須包含 Cognitive Mode block。若 project root 或 nested Git repositories dirty，final response 也必須包含 `### Project Git Report`，避免 root 非 git repo 的多 repo workspace 漏報。

必要行為：

- 不使用 shell grep / uname 判斷業務邏輯。
- Staged-file decision、runtime DB validation、knowledge validation、commit-msg parsing 都在 Go 中完成。
- Hook adapter 只解析 repo root / binary path，並把 git 傳的 positional args（如 commit-msg file path）透過 `"$@"` 轉發給 Go runner。

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

目的：用 Go-native compiler 從 SQLite canonical runtime documents 與 deterministic prose mappings refresh `runtime/runtime.db` projections。

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
- 預設 native SQLite-canonical projection mode 不依賴 Ruby、Python 或外部 `sqlite3` CLI。

驗證：

- `runtime/runtime.db` integrity check。
- `generated_surfaces` content assertion。
- compiler version / schema version 存在。
- Go compiler 讀取 `runtime/runtime.db` 的 `runtime_config_documents`、compiler mapping document、以及 deterministic prose sources，產生指定 `--db` output，並通過 native runtime DB validation。

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

### `ai-skill runtime obligations`

目的：列出目前 active bootstrap obligations，從 `runtime/runtime.db` 的 `generated_surfaces[runtime.core_bootstrap.contract]` JSON 讀取，並附上 canonical Bootstrap Receipt line。Source-of-truth 是 [`runtime/core-bootstrap.yaml`](../../../runtime/core-bootstrap.yaml) 與 runtime phase / gate tables；本 command read-only。

輸入：

- `--repo <path>`
- `--json` / `--plain`

副作用：無。

必要行為：

- 回報 `bootstrap_receipt` check，格式為 `Bootstrap: rules=✓ phase=<phase-id> obligations=<n> gates=<n>`。
- 列出 `per_session_obligations` / `per_turn_obligations` / `per_commit_obligations` 各自的 obligation IDs。
- 若 `runtime.core_bootstrap.contract` projection 不存在 → exit 30 並提示執行 `ai-skill runtime compile + refresh`。
- 不可修改 runtime.db 或 generated surfaces。
- 用途：debug "為什麼 commit-msg hook 擋我" + Phase 6 per-obligation dispatcher 對齊 hook 與 contract。

### `ai-skill runtime receipt`

目的：提供 Bootstrap Receipt 的正式 CLI surface，避免 agent 或 hook 直接拼 `sqlite3` 查詢或猜 runtime schema。

輸入：

- `--repo <path>`
- `--json` / `--plain`

副作用：無。

必要行為：

- 回報 `bootstrap_receipt` check，格式為 `Bootstrap: rules=✓ phase=<phase-id> obligations=<n> gates=<n>`。
- 回報 `per_turn_obligations` check，內容為 active per-turn obligation IDs。
- phase 取 runtime phase machine 的第一個非 `__config__` phase；obligations / gates 取 runtime tables row counts。
- 若 runtime projection 不存在或 schema 不可讀，阻斷並提示執行 `ai-skill runtime compile && ai-skill runtime refresh`。

### `ai-skill runtime audit`

目的：對 `knowledge/runtime/routing-registry.yaml` routes、`runtime/runtime.db` `generated_surfaces` 與 `validation/scenarios/**/*.yaml` 進行 4-way 分類（auto-detected / consumed / intentionally-manual / orphan），以揭露違反 `governance/lifecycle/system-upgrade-governance.yaml` §`define_runtime_trigger_flow` forbidden rules 的 orphan 條目。Source-of-truth 是上述三個 surface；本 command read-only。

輸入：

- `--repo <path>`
- `--json`（預設輸出 markdown 報告；`--json` 切換 JSON inventory）
- `--plain`

副作用：無。

必要行為：

- 預設輸出三表 markdown 報告（routes / generated surfaces / validation scenarios）+ summary 計數表 + orphan_total。
- `--json` 輸出 `Inventory` struct（含 `routes` / `surfaces` / `scenarios` / `summary` / `warnings` 欄位）。
- 分類規則：
  - **manual_activation**：routing-registry entry 含 `manual_activation: { reason: ... }` annotation。
  - **auto-detected via signal**：`runtime/cognitive-modes-discovery.yaml` 任一 signal description / pattern 提及該 route id。
  - **consumed**：`scripts/ai-skill-cli/**/*.go` 任一 Go 檔案內容包含該 route id / target_key / scenario id。
  - **orphan**：以上皆否。
- `ai-skill runtime validate` 自動以 warning-only check `runtime_audit_warning` 引用其 orphan 統計；audit 自身失敗不阻斷 validate。
- 額外 glossary coverage warning pass：掃 `plans/active/`、`architecture/`、`workflow/`、`analysis/`、`intelligence/`、`runtime/`、`ecosystem/` 內的 backtick-wrapped identifiers + snake_case ≥ 2 segments terms；若 `knowledge/runtime/sqlite/runtime-index.sqlite` 的 `glossary_terms.term` / `aliases` 找不到，emit `inventory.warnings`（依出現頻次排序，最多 50 條 + 截斷提示）。Heuristic 排除 path references（含 `/`）、單一英文短詞與 < 3 char terms 以降 false positive。
- 不可修改 runtime.db、routing-registry 或 generated surfaces。

### `ai-skill enforcement lint`

目的：對 `enforcement/enforcement-registry.yaml` 跑 Phase 3 lint 引擎並輸出 severity-grouped 結果。Source-of-truth 是 `enforcement/enforcement-registry.yaml`；本 command read-only。Phase 3 已 wire 進 `ai-skill runtime compile`（FAIL block / WARNING print），本 CLI 是 standalone 入口，方便 scenario / CI / local debug 直接呼叫。

輸入：

- `--repo <path>`（預設 `.`；用於解析 rule yaml 與 executor 檔案位置）
- `--registry <path>`（覆寫 registry 路徑，預設 `<repo>/enforcement/enforcement-registry.yaml`；scenario 用 synthetic 副本時必填）
- `--check <substr>`（依 finding type 過濾，substring match；如 `--check orphan_rule`）
- `--expect-finding <substr>`（assertion mode：若任一 finding 的 type / message / field key/value 含此 substring 則 exit 0，否則 exit 30；cross-platform，無需 shell `grep`）
- `--expect-severity <FAIL|WARNING>`（窄化 `--expect-finding` 到指定 severity）
- `--json` / `--plain`

副作用：當提供 `--registry` 時，在 `os.TempDir()` 建立 shadow repo（複製 enforcement / runtime / governance yaml + hooks.go），結束後刪除。

必要行為：

- 預設模式列每個 finding 為 `lint.<type>` Check，按 severity 分組統計 `findings_fail` / `findings_warn`。
- 任一 FAIL → exit 30 (`validation_failed`)；只有 WARNING → exit 0；零 finding → exit 0。
- Assertion mode (`--expect-finding`) 改 `mode=assert`，純看是否命中；命中 → exit 0，未命中 → exit 30。
- 不寫入 `enforcement/enforcement-registry.yaml`、不寫入 `runtime.db`、不執行 git 操作。

### `ai-skill enforcement coverage`

目的：將 `enforcement/enforcement-registry.yaml` 聚合成 6-bucket coverage report + verification level + runtime observation gap，作為 mechanical-enforcement-registry plan §Phase 4 主要交付。Source-of-truth 是 registry yaml；本 command read-only。

輸入：

- `--repo <path>`（預設 `.`）
- `--registry <path>`（覆寫 registry 路徑）
- `--format text|json|markdown`（預設 `text`）
- `--diff <ref>`（與 git revision 比較；`git show <ref>:enforcement/enforcement-registry.yaml` 失敗時轉 alert 不阻斷）
- `--detail`（text/markdown 模式開 per-class 表；預設只列 summary）
- `--self-check`（跨平台驗 3 format schema 的內建模式；列每項 check pass/fail，全 pass exit 0，否則 exit 30）

副作用：無。

必要行為：

- **6-bucket enum 鎖死**：`mechanical` / `behavioral_only` / `not_mechanizable` / `pending_implementation` / `research_required` / `deprecated`。即使 count=0 也輸出（JSON schema 穩定）。
- **Verification level 分類**：`full`（mechanical + scenario file exists）/ `symbol_only`（mechanical 無 scenario）/ `planned`（pending_implementation 或 research_required）/ `behavioral`（behavioral_only）/ `not_applicable`（not_mechanizable 或 deprecated）。
- **Runtime observation gap**：偵測 `runtime/runtime.db` 是否含 `executor_observations` table（Phase 5 wire 前不存在）。未 wire → `runtime_observed_pct` 一律 `null` + 全域 alert `runtime_observations_not_wired`；wire 後依 `rule_class_id` 查最近 `observation_window_days` 觸發次數。
- **Text format**：第一行固定 `Enforcement Coverage Report (YYYY-MM-DD)`，第二行為 `═` 重複；之後是 bucket summary + 條件 section（Pending impl / Research required / Behavioral_only / Not_mechanizable / Deprecated）；非 TTY 不得輸出 ANSI escape。
- **JSON format**：top-level keys `schema_version`（=1，數字）/ `generated_at`（ISO-8601）/ `total_rule_classes` / `observation_window_days` / `buckets`（6 keys）/ `per_class`（array of `{id, coverage, verification, runtime_observed_pct, alerts, ...}`）/ `alerts`（array，可空）/ `diff`（optional）。所有 key snake_case，禁混 camelCase。
- **Markdown format**：`# Enforcement Coverage Report` h1 + `## Summary` h2 表格 + 條件 `## Per-class detail` / `## Alerts` / `## Changes vs <ref>`，不得含 raw HTML。
- **Self-check**：內建 9 個 check（text 三項、markdown 兩項、json 四項），驗 first-line regex / schema_version 數字 / per_class array / snake_case keys 等；任一失敗 exit 30。

### `ai-skill glossary validate`

目的：驗證 `knowledge/glossary/*.md` 內的 glossary entries 是否符合 [`knowledge/glossary/README.md`](../../../knowledge/glossary/README.md) 定義的 entry schema、symmetry classification、命名規則與 `excludes` 引用合法性。Source-of-truth 是 README.md 內的 schema spec；本 command read-only validator。

輸入：

- `--repo <path>`（預設 `.`，glossary 根目錄 = `<repo>/knowledge/glossary/`）
- `--glossary <path>`（覆寫 glossary 根目錄，預設 `<repo>/knowledge/glossary/`）
- `--json` / `--plain`

副作用：無。

必要行為：

- 解析每個 `knowledge/glossary/*.md`（除 `README.md` 外）內的 H2 heading + 緊接 YAML code block；H2 文字必須與 YAML block 的 `term:` 完全相同。
- 檢查 required fields：`term`、`status`、`meaning`、`affects`、`owner-layer`。
- 檢查 optional fields shape：`aliases`、`anti-meaning`、`excludes`、`related-terms`、`introduced-by`、`deprecated-by`。
- 檢查 enum：
  - `status` ∈ `canonical` / `candidate` / `deprecated` / `superseded` / `alias-only` / `experimental` / `project-local`
  - `owner-layer` ∈ README.md `semantic_owner_domains:` 列出的 domains
  - `relation_type` (in `related-terms`) ∈ `alias_of` / `related_to` / `conflicts_with` / `owned_by` / `used_by` / `deprecated_by` / `replaced_by` / `derived_from` / `aggregates`
- 檢查 naming：`term:` 必須為 snake_case。
- 檢查 alias 規則：(a) `aliases:` 中字串不得出現為任何 entry 的 `term:`；(b) alias chain 不得形成 cycle；(c) 新 entry 禁用 `status: alias-only`（僅供 deprecated 過渡）。
- 檢查 `introduced-by` / `deprecated-by` 形狀：必須為 `plans/<path>` 或 `constitution/ADR-XXX.md`，禁止 commit SHA、issue 編號、PR URL。
- 檢查 `excludes:` 引用合法性：所有列出的字串必須為現有 entry 的 `term:`。
- 檢查 symmetry：symmetric relation（`related_to` / `conflicts_with`，依 README.md symmetry 分類表）必須雙向出現於兩端 entry；asymmetric 僅單向。
- 檢查 forbidden patterns：term body 不得包含 project-specific hosts / paths / class names / sample IDs / incident evidence / executable contract semantics / raw memory replay。
- 失敗時 exit 30 (`validation_failed`)；JSON `checks` 列出每個 violation 的 entry path、term、rule id、remediation。
- 通過時 exit 0，並回報 entries / aliases / relations 計數。
- 不修改 `knowledge/glossary/` 內任何檔案。

驗證：

- Go unit tests in `scripts/ai-skill-cli/internal/glossary/validator_test.go`（或等效路徑），fixture 見 [`test-fixture-plan.md`](test-fixture-plan.md) `fixture/glossary-*`。
- BDD scenarios 見 [`bdd-scenarios.md`](bdd-scenarios.md) §Glossary。
- 若未來接入 `ai-skill runtime validate`，須同時更新本檔與 runtime validate 行為說明。

## 舊 Script 對應

完整功能盤點、side effects、外部依賴與測試證據見 [`script-parity-inventory.md`](script-parity-inventory.md)。開發者 handoff map 見 [`legacy-to-go-migration-map.md`](legacy-to-go-migration-map.md)。本節只保留命令契約層的摘要映射；不得用本表取代 parity 驗收。

| 現有 script / 入口 | CLI 命令 | 遷移定位 |
| --- | --- | --- |
| deleted `scripts/init-new-project.sh` | `ai-skill init-project` | Go write mode complete |
| deleted `scripts/agent-goals.sh` | `ai-skill goals` | Go lifecycle write mode complete |
| deleted `scripts/install-hooks.sh` / `.githooks/` | `ai-skill hooks install` | dry-run planner uses `scripts/git-hooks/`; write mode still blocked until fixture-backed |
| deleted `scripts/sync-cursor-bundle.sh` | `ai-skill sync-cursor-bundle` | shell 已刪；Go dry-run 已有，write mode 待 Go parity |
| deleted `scripts/ai-skill-close-loop.sh` | `ai-skill close-loop` | Go commit / push parity complete |
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
- 不得新增長期 shell script 作為新功能入口；若需要 automation，先擴充本 command contract 與 Go CLI。

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
| `glossary validate` | `knowledge/glossary/*.md`（H2 + YAML block entries） | 無 | 無 |
| `enforcement lint` | `enforcement/enforcement-registry.yaml` + `enforcement/runtime/governance/**/*.yaml` + `scripts/ai-skill-cli/internal/app/hooks.go` | 僅 `os.TempDir()` shadow repo（使用 `--registry` 時，結束即刪） | 無 |
| `enforcement coverage` | `enforcement/enforcement-registry.yaml` + `runtime/runtime.db`（檢 `executor_observations`） + `validation/scenarios/`（heuristic scenario 命名比對） + git `show <ref>:enforcement/enforcement-registry.yaml`（`--diff` 模式） | 無 | `--diff` 模式需 Git |
