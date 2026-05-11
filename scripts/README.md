# Scripts

| 檔案 | 用途 |
| --- | --- |
| [`sync-cursor-bundle.sh`](sync-cursor-bundle.sh) | 可選的 Cursor symlink / bundle bridge：把本庫同步到 **`~/.cursor/bundles/shared-rules`**（共用規則）與 **`~/.cursor/bundles/ai-skill/`**（各 skill），再讓 `~/.cursor/shared-rules`、`~/.cursor/skills/<name>` 指向 bundle（與其他 `.cursor` 內容分流）。Reference-only 不需要執行。 |
| [`ai-skill-close-loop.sh`](ai-skill-close-loop.sh) | 保守收尾工具：偵測 active close-loop lock、列出 dirty 檔案 owner group；預設 dry-run，`--commit` 才分組提交，`--push` 才推遠端。 |
| [`agent-goals.sh`](agent-goals.sh) | 工具中立的專案暫存 goal ledger helper：在 `<PROJECT_ROOT>/.agent-goals/` 建立、更新、拆解、暫停、完成刪除對話目標；不提交 goal 檔。 |
| [`git-hooks/post-commit`](git-hooks/post-commit) | **可選。**在本 repo 設定 `git config core.hooksPath scripts/git-hooks` 且 `AI_SKILL_SYNC_CURSOR_BUNDLE=1` 時，**`git commit`** 後會執行 `sync-cursor-bundle.sh`。 |

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
- 若 dirty path 無法歸到 `architecture`、`shared-rules`、`scripts`、`ai-tools`、`.cursor/rules` 或某個 `skills/<name>` owner，腳本會停止。
- 預設只 dry-run；必須明確加 `--commit` / `--push` 才會寫入 git。

## Conversation goal ledger helper

規則正文見 [`shared-rules/conversation-goal-ledger.md`](../shared-rules/conversation-goal-ledger.md)。Goal ledger 是專案本地暫存狀態，放在 `<PROJECT_ROOT>/.agent-goals/`，用來在 agent 中斷、轉移目標、多 agent 操作、context compact、已有 TodoWrite 或看到 dirty files 後回溯目前目標，不取代 git / issue tracker / Ai-skill writeback transaction。

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
- `.agent-goals/README.md` 會自動刷新成主目標表，連到 `goals/*.md`，並顯示 mode、owner、lock、open work / decisions、plan/todo links、下一步與更新時間。
- `start`、`update`、`split` 可重複使用 `--plan` 與 `--todo`，把 planning 文件章節、TodoWrite ID、checklist item 或 issue ID 連到 goal。
- `start`、`update`、`split` 可用 `--parallelization parallelizable|single-owner|non-parallelizable`，讓主表顯示目前是否能分工。
- `update` 可用 `--missing`、`--decision`、`--strengthen` 把未完成、待決策與待補強項目放進主表。
- 每個 goal 更新時會使用 `.agent-goals/locks/<goal-id>.lock/` 防止多 agent 同時寫入。
- 若主表或 `status` 顯示重疊 goal 已被其他 owner/lock 處理，停止修改並提示使用者決定：等待、接手、拆子目標或另開非重疊 goal。
- 對 git 合併/發版、Ai-skill writeback transaction、資料遷移、credential rotation、破壞性操作等不可分工流程，將 goal 標成 `non-parallelizable`。
- Stale lock 可用 `cleanup` 清理；TTL 預設 30 分鐘，可用 `AGENT_GOALS_LOCK_TTL_SECONDS` 覆寫。
