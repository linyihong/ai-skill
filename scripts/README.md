# Scripts

| 檔案 | 用途 |
| --- | --- |
| [`ai-skill-cli/`](ai-skill-cli/README.md) | 跨平台 Go CLI / runtime toolchain 的開發根目錄；`docs/` 放文件先行產物與舊腳本 parity 盤點，未來 `cmd/`、`internal/`、`testdata/` 放程式碼與 fixtures。 |
| [`git-hooks/pre-commit`](git-hooks/pre-commit) / [`git-hooks/post-commit`](git-hooks/post-commit) | 最小 Git hook adapter；實際邏輯在 `ai-skill hooks run pre-commit/post-commit`。 |

## Go-first script policy

新增 automation 預設必須進 [`ai-skill-cli/`](ai-skill-cli/README.md) 的 Go CLI；不要再新增長期 `.sh`、`.rb` 或 `.py` 入口。Shell 只允許作為 Git hook adapter 或轉呼叫 repo-local `ai-skill` binary 的 thin bootstrap wrapper。

已完成 parity 的 shell 入口必須刪除；新功能先做成 Go command，再補 fixture、文件與 release gate。

## Go CLI migration map

跨平台 Go CLI 的 source 在 [`ai-skill-cli/`](ai-skill-cli/README.md)。開發者 handoff map 是 [`ai-skill-cli/docs/legacy-to-go-migration-map.md`](ai-skill-cli/docs/legacy-to-go-migration-map.md)，完整 parity source-of-truth 是 [`ai-skill-cli/docs/script-parity-inventory.md`](ai-skill-cli/docs/script-parity-inventory.md)；本表只列 scripts 入口的使用者導覽。

| 現有入口 | 目標 CLI | 目前狀態 | 收尾政策 |
| --- | --- | --- | --- |
| deleted `init-new-project.sh` | `ai-skill init-project` | write mode 已實作；支援 `--dry-run`、`--force`、`--tools` | 舊 shell 已刪；新專案初始化只能走 Go CLI。 |
| deleted `agent-goals.sh` | `ai-skill goals` | `init/status/start/update/split/pause/complete/cleanup` 已由 Go 實作 | 舊 shell 已刪；goal ledger lifecycle 只能走 Go CLI。 |
| deleted legacy hook installer | `ai-skill hooks install` / `ai-skill hooks run ...` | install 會安裝最小 adapter；pre/post-commit 業務邏輯在 Go runner | Git hook files 只保留 adapter，不承載 reusable automation logic。 |
| deleted `sync-cursor-bundle.sh` | `ai-skill sync-cursor-bundle` | 舊 shell 已刪；Go command 已有 explicit-target dry-run planner，但 write mode 仍回 `write_mode_not_implemented` | 未來若要恢復 Cursor mirror 寫入，只能補 `ai-skill sync-cursor-bundle` Go write mode。 |
| deleted `ai-skill-close-loop.sh` | `ai-skill close-loop` | `--dry-run`、`--commit`、`--push` 已由 Go 實作 | 舊 shell 已刪；close-loop 寫入流程只走 Go CLI。 |
| runtime Ruby helpers | `ai-skill runtime ...` | 已刪除 runtime report/index/query/validation/migration/state/sync Ruby entrypoints；`runtime validate`、`runtime refresh`、`runtime query`、`runtime compile` 的 desktop path 預設都走 Go-native，不依賴 Ruby、Python 或 `sqlite3` CLI。 | 已完成 native 覆蓋或易誤用的 scripts 直接刪除；runtime compiler source 已恢復為 YAML，Go compiler 是唯一 active compile path。 |
| deleted Roo Python helper | `ai-skill roo set-global-custom-instructions` | 已刪除 `scripts/set-roo-global-custom-instructions.py`；guarded tool adapter 由 Go CLI 實作。 | fake VS Code SQLite DB tests 通過後刪除，避免未來錯誤引用 Python helper。 |

Legacy script closure policy：[`ai-skill-cli/docs/legacy-script-disposition.md`](ai-skill-cli/docs/legacy-script-disposition.md) 是舊 shell / Ruby / Python entrypoints 的最終 disposition source；runtime desktop CLI 已是 primary，已覆蓋的 scripts 已刪除，Git hook shell 只保留最小 adapter。

## New project initialization

開新專案時，用 repo-local Go CLI 一次設定所有 AI 工具：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 init-project --project ~/projects/my-new-app
```

這會在目標專案中建立：

| 工具 | 產出檔案 | 內容 |
|------|---------|------|
| Roo Code | `.roomodes` | 5 個 mode（code/architect/ask/debug/orchestrator），含語言規則 + 知識更新 checkpoint |
| Cursor | `.cursor/rules/ai-skill-bootstrap.mdc` | alwaysApply 規則，含啟動流程 + **專案 durable Markdown 預設**（`workflow/documentation/`）+ 語言規則 + checkpoint |
| Cursor | `.cursor/hooks.json` | sessionStart hook 提醒 |
| Claude Code | `CLAUDE.md` | 自動載入 Core Bootstrap + 語言規則 + checkpoint |
| 通用 | `.agent-goals/` | 對話目標帳本目錄 |

只設定特定工具：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 init-project --project ~/projects/my-new-app --tools roo,cursor
```

預覽模式（不實際寫入）：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 init-project --project ~/projects/my-new-app --dry-run
```

覆蓋已有檔案：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 init-project --project ~/projects/my-new-app --force
```

完整說明見 [`ai-tools/new-project-onboarding.md`](../ai-tools/new-project-onboarding.md)。

**規則：**reference-only 是預設，不需要跑 bundle sync。舊 `sync-cursor-bundle.sh` 已刪；新的同步功能必須加到 `ai-skill sync-cursor-bundle` Go command，不能加回 shell。

## Close-loop automation

Close-loop 寫入入口已移到 Go CLI；舊 `ai-skill-close-loop.sh` 已刪除。

先檢查，不提交：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 close-loop --dry-run
```

沒有人正在操作、所有 dirty path 都可歸屬時，依 owner 分組提交。預設不跑 Cursor bundle sync：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 close-loop --commit
```

提交後也推送目前 branch：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 close-loop --commit --push
```

安全條件：

- 若 `.git/ai-skill-agent.lock` 顯示其他 agent / user 仍活躍，腳本會停止，不 commit、不 push。
- 若存在 merge / rebase / cherry-pick 狀態，腳本會停止。
- 若 dirty path 無法歸到 `architecture`（含下一階段 top-level 分層）、`enforcement`、`workflow`、`analysis`、`intelligence`、`scripts`、`ai-tools` 或 `.cursor/rules` owner，腳本會停止。
- 預設只 dry-run；必須明確加 `--commit` / `--push` 才會寫入 git。
- 根目錄 `CONTRIBUTING.md` 與 `README.md`、`enforcement/`、`.gitignore` 同屬 `shared` owner group（見 `closeLoopGroupForPath` in `scripts/ai-skill-cli/internal/app/close_loop.go`）。

## Conversation goal ledger helper

規則正文見 [`enforcement/conversation-goal-ledger.md`](../enforcement/conversation-goal-ledger.md)。Goal ledger 是專案本地暫存狀態，放在 `<PROJECT_ROOT>/.agent-goals/`，用來在 agent 中斷、轉移目標、多 agent 操作、context compact、已有 TodoWrite 或看到 dirty files 後回溯目前目標，不取代 git / issue tracker / Ai-skill writeback transaction。

進入多步驟工作、使用者要求「繼續」前一個任務、或已看到 active project 有 modified / staged / untracked files 時，先檢查 active goal、priority、owner、lock、parallelization mode、plan/todo links、open missing/decision/strengthen：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 goals status --project <PROJECT_ROOT>
```

若尚未建立 ledger 且任務不是單一回覆即可完成，先初始化：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 goals init --project <PROJECT_ROOT>
```

建立目前主要目標：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 goals start --project <PROJECT_ROOT> \
  --id P1-example-goal \
  --title "Example goal" \
  --source "User request summary" \
  --parallelization single-owner \
  --next "Next concrete action" \
  --criteria "Observable completion condition" \
  --plan "docs/implementation-plan.md#example" \
  --todo "implement-example"
```

更新、拆解、暫停與完成：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 goals update --project <PROJECT_ROOT> --id P1-example-goal --note "Read dependencies" --next "Implement the change"
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 goals update --project <PROJECT_ROOT> --id P1-example-goal --missing "Validation examples are not written" --decision "Choose whether this remains P1" --strengthen "Add stronger completion criteria"
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 goals update --project <PROJECT_ROOT> --id P1-example-goal --parallelization non-parallelizable --note "Live capture must stay single-owner"
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 goals split --project <PROJECT_ROOT> --parent P1-example-goal --id P2-child-goal --title "Child goal"
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 goals pause --project <PROJECT_ROOT> --id P1-example-goal --reason "User changed priority"
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 goals complete --project <PROJECT_ROOT> --id P1-example-goal --validated --note "Validation passed"
```

安全條件：

- `complete` 只有在傳入 `--validated` 時才會刪除 goal 檔；否則會保留並標成 `needs-validation`。
- 完成條件、驗證與 final/handoff 都成立時，同一輪必須執行 `complete --validated` 或手動刪除 goal 並刷新主表；不要把 `completed` row 長期留在 `.agent-goals/README.md`。
- 若 goal 完成後仍代表長期 roadmap、phase、未完成能力、migration / promotion / deprecation 狀態或治理決策，先把該狀態回寫到 `architecture/`、layer README、`governance/`、`knowledge/`、`metadata/`、正式 project docs 或 issue，再執行 `complete --validated`。`.agent-goals/` 只保存 active conversation goal，不保存長期目標 archive。
- `.agent-goals/README.md` 會自動刷新成主目標表，連到 `goals/*.md`，並顯示 mode、owner、lock、open work / decisions、plan/todo links、下一步與更新時間。
- `start`、`update`、`split` 可重複使用 `--plan` 與 `--todo`，把 planning 文件章節、TodoWrite ID、checklist item 或 issue ID 連到 goal。
- `start`、`update`、`split` 可用 `--parallelization parallelizable|single-owner|non-parallelizable`，讓主表顯示目前是否能分工。
- `update` 可用 `--missing`、`--decision`、`--strengthen` 把未完成、待決策與待補強項目放進主表。
- 每個 goal 更新時會使用 `.agent-goals/locks/<goal-id>.lock/` 防止多 agent 同時寫入。
- 若主表或 `status` 顯示重疊 goal 已被其他 owner/lock 處理，停止修改並提示使用者決定：等待、接手、拆子目標或另開非重疊 goal。
- 對 git 合併/發版、Ai-skill writeback transaction、資料遷移、credential rotation、破壞性操作等不可分工流程，將 goal 標成 `non-parallelizable`。
- Stale lock 可用 `cleanup` 清理；TTL 預設 30 分鐘，可用 `AGENT_GOALS_LOCK_TTL_SECONDS` 覆寫。

## Knowledge runtime validation

> **重要**：修改 `knowledge/` 或 `validation/` 下的檔案後，**必須**執行對應平台的 repo-local binary（例如 macOS Apple Silicon：`scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 runtime refresh`）確認 native validators 通過，再提交。Pre-commit hook（`scripts/git-hooks/pre-commit`）會在 commit 時自動檢查，但建議在修改後立即執行以加速迭代。

產生並檢查 generated knowledge surfaces：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 runtime refresh
```

或逐步執行：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 runtime refresh
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 runtime validate
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 runtime query feedback --limit 5
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 runtime query feedback --layer feedback --limit 5
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 runtime query --graph --type depends_on --limit 5
```

`runtime refresh` 會更新 `knowledge/runtime/runtime-report.md`、`knowledge/runtime/model-context-report.md`、`knowledge/runtime/model-checklists.md`，並產生被 git ignore 的本機 `knowledge/runtime/sqlite/runtime-index.sqlite`。
`runtime query` 只輸出候選來源，不取代 canonical Markdown / YAML。

此 helper 目前驗證：

- `knowledge/runtime/routing-registry.yaml` 的 records、required dependencies、candidate sources、model profile 與 metadata 欄位。
- `knowledge/runtime/refresh-policy.yaml` 的 surfaces、decision values 與 validation / close-loop 欄位。
- `knowledge/summaries/*.md` 的必要 summary table 欄位、source links，以及 `knowledge/summaries/README.md` 是否列出 summary。
- `knowledge/graphs/*.yaml` 的 source、edge types、edge targets 與 metadata 欄位。
- `knowledge/runtime/runtime-report.md` 與 `knowledge/runtime/model-context-report.md` 的 Markdown links。

這個 helper 只做 deterministic validation；它不自動修改 summaries、graphs 或 registry。若檢查失敗，先修 source / generated surface，再執行 lints、Markdown link check、close-loop dry run 與 commit / push / readback。

## Script CI / GitHub workflow guard

修改 `scripts/ai-skill-cli/**`、`scripts/git-hooks/**`、repo-local binaries 或 script CI 行為時，必須同步檢查 `.github/workflows/ai-skill-cli.yml`：

- workflow `paths` 是否涵蓋被修改的 script surface。
- `go test ./...` 是否仍覆蓋 repo-local binary checksum / `BUILDINFO` parity。
- artifact build step 是否仍產出五個平台 binary 與 `SHA256SUMS`。
- 若 CLI source 變更，先 commit source，再用該 source commit 重建 `scripts/ai-skill-cli/bin/`，最後用第二個 commit 提交 binary refresh。

Push 前可跑：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 hooks run pre-push --repo . --json
```

這會在本分支相對 upstream 修改 CLI、git hooks 或 CLI workflow 時，先執行 `cd scripts/ai-skill-cli && go test ./...`，用來在 push 前攔截 GitHub Actions 會報的 stale binary / checksum / workflow regression。
