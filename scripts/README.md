# Scripts

| 檔案 | 用途 |
| --- | --- |
| [`init-new-project.sh`](init-new-project.sh) | **新專案初始化**：在目標專案目錄中一次建立 Roo Code（`.roomodes`）、Cursor（`.cursor/rules/`）、Claude Code（`CLAUDE.md`）的設定檔，全部指向 Ai-skill 知識庫；bootstrap 內含 **專案 durable Markdown 預設**（寫入 `docs/`、`README.md` 等前先讀 `workflow/documentation/`）。開新專案時跑一次就好。 |
| [`sync-cursor-bundle.sh`](sync-cursor-bundle.sh) | 可選的 Cursor symlink / bundle bridge：把本庫同步到 **`~/.cursor/bundles/enforcement`**（共用規則）與 **`~/.cursor/bundles/ai-skill/`**（workflow / analysis / intelligence source）。Reference-only 不需要執行。 |
| [`ai-skill-close-loop.sh`](ai-skill-close-loop.sh) | 保守收尾工具：偵測 active close-loop lock、列出 dirty 檔案 owner group；預設 dry-run，`--commit` 才分組提交，`--push` 才推遠端。 |
| [`ai-skill-cli/`](ai-skill-cli/README.md) | 跨平台 Go CLI / runtime toolchain 的開發根目錄；`docs/` 放文件先行產物與舊腳本 parity 盤點，未來 `cmd/`、`internal/`、`testdata/` 放程式碼與 fixtures。 |
| [`agent-goals.sh`](agent-goals.sh) | 工具中立的專案暫存 goal ledger helper：在 `<PROJECT_ROOT>/.agent-goals/` 建立、更新、拆解、暫停、完成刪除對話目標；不提交 goal 檔。 |
| [`validate-knowledge-runtime.rb`](validate-knowledge-runtime.rb) | 驗證 knowledge runtime generated surfaces：routing registry、refresh policy、summaries 與 graph records 的 YAML / Markdown 格式、必要欄位與 canonical path。 |
| [`generate-knowledge-runtime-report.rb`](generate-knowledge-runtime-report.rb) | 從 routing registry、summaries、graphs 與 refresh policy 產生 deterministic runtime report。 |
| [`generate-model-context-report.rb`](generate-model-context-report.rb) | 從 routing registry 的 model 欄位產生 model-aware context loading report。 |
| [`generate-model-checklists.rb`](generate-model-checklists.rb) | 從 routing registry 產生 per-model context-loading checklist artifact。 |
| [`generate-runtime-sqlite-index.rb`](generate-runtime-sqlite-index.rb) | 從 summaries、routing registry、graphs 與 feedback lessons 產生本機 SQLite / FTS lookup cache。 |
| [`query-runtime-index.rb`](query-runtime-index.rb) | 用 keyword 查詢本機 SQLite runtime index，依 rank / priority / confidence / context cost 回傳少量 candidate source paths。 |
| [`validate-runtime-sqlite-index.rb`](validate-runtime-sqlite-index.rb) | 驗證 SQLite runtime index 的 integrity、row counts、source paths、FTS、source checksum 與 git ignore 邊界。 |
| [`query-knowledge-graph.rb`](query-knowledge-graph.rb) | 查詢 graph edges，支援 source / target / type / keyword filters。 |
| [`refresh-knowledge-runtime.rb`](refresh-knowledge-runtime.rb) | 一鍵重建 model/runtime reports、SQLite index，並執行 runtime validators。 |
| [`git-hooks/post-commit`](git-hooks/post-commit) | **可選。**在本 repo 設定 `git config core.hooksPath scripts/git-hooks` 且 `AI_SKILL_SYNC_CURSOR_BUNDLE=1` 時，**`git commit`** 後會執行 `sync-cursor-bundle.sh`。 |

## Go CLI migration map

跨平台 Go CLI 的 source 在 [`ai-skill-cli/`](ai-skill-cli/README.md)。完整 parity source-of-truth 是 [`ai-skill-cli/docs/script-parity-inventory.md`](ai-skill-cli/docs/script-parity-inventory.md)；本表只列 scripts 入口的使用者導覽。

| 現有入口 | 目標 CLI | 目前狀態 | 收尾政策 |
| --- | --- | --- | --- |
| `init-new-project.sh` | `ai-skill init-project` | `--dry-run` planner 已實作；write mode 等 template parity | CLI parity、fixtures、文件通過後刪除舊 shell 入口。 |
| `agent-goals.sh` | `ai-skill goals` | `status` read-only 與 `init --dry-run` planner 已實作；寫入命令待 parity | 完整 goal lifecycle parity 通過後刪除舊 shell 入口。 |
| `install-hooks.sh` | `ai-skill hooks install` | dry-run planner 已實作；copy / chmod write mode 待 parity | hook install parity 通過後刪除舊 shell 入口；Git hook files 本身可作為 hook adapter 保留。 |
| `sync-cursor-bundle.sh` | `ai-skill sync-cursor-bundle` | explicit-target dry-run planner 已實作；managed mirror write mode 待 parity | 保留條件需寫明 owner、期限與移除條件；不得成為通用 CLI 預設行為。 |
| `ai-skill-close-loop.sh` | `ai-skill close-loop` | dry-run inspection 已實作；commit / push 待 parity | close-loop lock、dirty group、merge/rebase、dry-run、commit/push parity 通過後刪除或降為短期 thin wrapper。 |
| runtime Ruby helpers | `ai-skill runtime ...` | Phase 3 wrapper-first / native split；`runtime refresh` 已覆蓋 ordered steps / first-failure block，`generate-model-context-report.rb` 已有第一個 Ruby vs Go exact parity guard，`runtime validate` 已開始 native DB / SQLite index / git-ignore checks，`runtime query` 已開始 native SQLite index / knowledge graph query | runtime refresh / compile / validate / query parity 通過後刪除被取代的 Ruby entrypoints；hook adapter 例外需文件化。 |

## New project initialization

開新專案時，用 [`init-new-project.sh`](init-new-project.sh) 一次設定所有 AI 工具：

```bash
./scripts/init-new-project.sh ~/projects/my-new-app
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
./scripts/init-new-project.sh ~/projects/my-new-app --tools roo,cursor
```

預覽模式（不實際寫入）：

```bash
./scripts/init-new-project.sh ~/projects/my-new-app --dry-run
```

覆蓋已有檔案：

```bash
./scripts/init-new-project.sh ~/projects/my-new-app --force
```

完整說明見 [`ai-tools/new-project-onboarding.md`](../ai-tools/new-project-onboarding.md)。

**規則：**reference-only 是預設，不需要跑 bundle sync。只有本機明確用 Cursor symlink / bundle / copy mirror 佈署，且希望 mirror 立刻跟上時，才跑 `sync-cursor-bundle.sh`（或以 `AI_SKILL_SYNC_CURSOR_BUNDLE=1` 啟用上述 hook / close-loop helper 同步）。

在本庫根目錄執行：

```bash
chmod +x scripts/sync-cursor-bundle.sh   # 只需做一次
./scripts/sync-cursor-bundle.sh
```

啟用 commit 後自動同步（選用）：

```bash
git config core.hooksPath scripts/git-hooks
export AI_SKILL_SYNC_CURSOR_BUNDLE=1
```

## Close-loop automation

先檢查，不提交：

```bash
./scripts/ai-skill-close-loop.sh
```

沒有人正在操作、所有 dirty path 都可歸屬時，依 owner 分組提交。預設不跑 Cursor bundle sync：

```bash
./scripts/ai-skill-close-loop.sh --commit
```

若本機 Cursor bundle / mirror 需要跟上，明確啟用同步：

```bash
AI_SKILL_SYNC_CURSOR_BUNDLE=1 ./scripts/ai-skill-close-loop.sh --commit
```

提交後也推送目前 branch：

```bash
./scripts/ai-skill-close-loop.sh --commit --push
```

安全條件：

- 若 `.git/ai-skill-agent.lock` 顯示其他 agent / user 仍活躍，腳本會停止，不 commit、不 push。
- 若存在 merge / rebase / cherry-pick 狀態，腳本會停止。
- 若 dirty path 無法歸到 `architecture`（含下一階段 top-level 分層）、`enforcement`、`workflow`、`analysis`、`intelligence`、`scripts`、`ai-tools` 或 `.cursor/rules` owner，腳本會停止。
- 預設只 dry-run；必須明確加 `--commit` / `--push` 才會寫入 git。
- 根目錄 `CONTRIBUTING.md` 與 `README.md`、`enforcement/`、`.gitignore` 同屬 `shared` owner group（見 `group_for_path` in `ai-skill-close-loop.sh`）。

## Conversation goal ledger helper

規則正文見 [`enforcement/conversation-goal-ledger.md`](../enforcement/conversation-goal-ledger.md)。Goal ledger 是專案本地暫存狀態，放在 `<PROJECT_ROOT>/.agent-goals/`，用來在 agent 中斷、轉移目標、多 agent 操作、context compact、已有 TodoWrite 或看到 dirty files 後回溯目前目標，不取代 git / issue tracker / Ai-skill writeback transaction。

進入多步驟工作、使用者要求「繼續」前一個任務、或已看到 active project 有 modified / staged / untracked files 時，先檢查 active goal、priority、owner、lock、parallelization mode、plan/todo links、open missing/decision/strengthen：

```bash
./scripts/agent-goals.sh --project <PROJECT_ROOT> status
```

若尚未建立 ledger 且任務不是單一回覆即可完成，先初始化：

```bash
./scripts/agent-goals.sh --project <PROJECT_ROOT> init
```

建立目前主要目標：

```bash
./scripts/agent-goals.sh --project <PROJECT_ROOT> start \
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
./scripts/agent-goals.sh --project <PROJECT_ROOT> update --id P1-example-goal --note "Read dependencies" --next "Implement the change"
./scripts/agent-goals.sh --project <PROJECT_ROOT> update --id P1-example-goal --missing "Validation examples are not written" --decision "Choose whether this remains P1" --strengthen "Add stronger completion criteria"
./scripts/agent-goals.sh --project <PROJECT_ROOT> update --id P1-example-goal --parallelization non-parallelizable --note "Live capture must stay single-owner"
./scripts/agent-goals.sh --project <PROJECT_ROOT> split --parent P1-example-goal --id P2-child-goal --title "Child goal"
./scripts/agent-goals.sh --project <PROJECT_ROOT> pause --id P1-example-goal --reason "User changed priority"
./scripts/agent-goals.sh --project <PROJECT_ROOT> complete --id P1-example-goal --validated --note "Validation passed"
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

> **重要**：修改 `knowledge/` 或 `validation/` 下的檔案後，**必須**執行 `ruby scripts/refresh-knowledge-runtime.rb` 確認所有 validator 通過，再提交。Pre-commit hook（`scripts/git-hooks/pre-commit`）會在 commit 時自動檢查，但建議在修改後立即執行以加速迭代。

產生並檢查 generated knowledge surfaces：

```bash
ruby scripts/refresh-knowledge-runtime.rb
```

或逐步執行：

```bash
ruby scripts/generate-knowledge-runtime-report.rb --write
ruby scripts/generate-model-context-report.rb --write
ruby scripts/generate-model-checklists.rb --write
ruby scripts/validate-knowledge-runtime.rb
ruby scripts/generate-runtime-sqlite-index.rb
ruby scripts/validate-runtime-sqlite-index.rb
ruby scripts/query-runtime-index.rb feedback --limit 5
ruby scripts/query-runtime-index.rb feedback --layer feedback --limit 5
ruby scripts/query-knowledge-graph.rb --type depends_on --limit 5
```

`generate-knowledge-runtime-report.rb --write` 會更新 `knowledge/runtime/runtime-report.md`，讓 agent 可快速檢視目前 routes、summaries、graphs 與 refresh decisions。
`generate-model-context-report.rb --write` 會更新 `knowledge/runtime/model-context-report.md`，依 profile 與 compression level 整理 model-aware loading view。
`generate-model-checklists.rb --write` 會更新 `knowledge/runtime/model-checklists.md`，依 profile 產生可執行的 context-loading checklist。
`generate-runtime-sqlite-index.rb` 會產生被 git ignore 的本機 `knowledge/runtime/sqlite/runtime-index.sqlite`；query helper 只輸出候選來源，不取代 canonical Markdown / YAML。

此 helper 目前驗證：

- `knowledge/runtime/routing-registry.yaml` 的 records、required dependencies、candidate sources、model profile 與 metadata 欄位。
- `knowledge/runtime/refresh-policy.yaml` 的 surfaces、decision values 與 validation / close-loop 欄位。
- `knowledge/summaries/*.md` 的必要 summary table 欄位、source links，以及 `knowledge/summaries/README.md` 是否列出 summary。
- `knowledge/graphs/*.yaml` 的 source、edge types、edge targets 與 metadata 欄位。
- `knowledge/runtime/runtime-report.md` 與 `knowledge/runtime/model-context-report.md` 的 Markdown links。

這個 helper 只做 deterministic validation；它不自動修改 summaries、graphs 或 registry。若檢查失敗，先修 source / generated surface，再執行 lints、Markdown link check、close-loop dry run 與 commit / push / readback。
