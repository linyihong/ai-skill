# Scripts

| 檔案 | 用途 |
| --- | --- |
| [`sync-cursor-bundle.sh`](sync-cursor-bundle.sh) | 把本庫同步到 **`~/.cursor/bundles/shared-rules`**（共用規則）與 **`~/.cursor/bundles/ai-skill/`**（各 skill），再讓 `~/.cursor/shared-rules`、`~/.cursor/skills/<name>` 指向 bundle（與其他 `.cursor` 內容分流）。 |
| [`ai-skill-close-loop.sh`](ai-skill-close-loop.sh) | 保守收尾工具：偵測 active close-loop lock、列出 dirty 檔案 owner group；預設 dry-run，`--commit` 才分組提交，`--push` 才推遠端。 |
| [`agent-goals.sh`](agent-goals.sh) | 工具中立的專案暫存 goal ledger helper：在 `<PROJECT_ROOT>/.agent-goals/` 建立、更新、拆解、暫停、完成刪除對話目標；不提交 goal 檔。 |
| [`git-hooks/post-commit`](git-hooks/post-commit) | **可選。**在本 repo 設定 `git config core.hooksPath scripts/git-hooks` 後，每次 **`git commit`** 會自動執行 `sync-cursor-bundle.sh`。 |

**規則：**只要改過 **`shared-rules/`** 或 **`skills/`** 且本機用 bundles 佈署，就應跑 `sync-cursor-bundle.sh`（或依賴上述 hook）。

在本庫根目錄執行：

```bash
chmod +x scripts/sync-cursor-bundle.sh   # 只需做一次
./scripts/sync-cursor-bundle.sh
```

啟用 commit 後自動同步（選用）：

```bash
git config core.hooksPath scripts/git-hooks
```

## Close-loop automation

先檢查，不提交：

```bash
./scripts/ai-skill-close-loop.sh
```

沒有人正在操作、所有 dirty path 都可歸屬時，依 owner 分組提交：

```bash
./scripts/ai-skill-close-loop.sh --commit
```

提交後也推送目前 branch：

```bash
./scripts/ai-skill-close-loop.sh --commit --push
```

安全條件：

- 若 `.git/ai-skill-agent.lock` 顯示其他 agent / user 仍活躍，腳本會停止，不 commit、不 push。
- 若存在 merge / rebase / cherry-pick 狀態，腳本會停止。
- 若 dirty path 無法歸到 `shared-rules`、`scripts` 或某個 `skills/<name>` owner，腳本會停止。
- 預設只 dry-run；必須明確加 `--commit` / `--push` 才會寫入 git。

## Conversation goal ledger helper

規則正文見 [`shared-rules/conversation-goal-ledger.md`](../shared-rules/conversation-goal-ledger.md)。Goal ledger 是專案本地暫存狀態，放在 `<PROJECT_ROOT>/.agent-goals/`，用來在 agent 中斷、轉移目標、多 agent 操作或 context compact 後回溯目前目標，不取代 git / issue tracker / Ai-skill writeback transaction。

初始化目前專案的 goal ledger，並把 `.agent-goals/` 寫入 `.git/info/exclude`：

```bash
./scripts/agent-goals.sh --project <PROJECT_ROOT> init
```

查看目前 active goals 與 locks。`.agent-goals/README.md` 是給人類與 AI 一開始判斷「還沒做什麼、要先做哪個、哪裡還要補強」的主表：

```bash
./scripts/agent-goals.sh --project <PROJECT_ROOT> status
```

建立目前主要目標：

```bash
./scripts/agent-goals.sh --project <PROJECT_ROOT> start \
  --id P1-example-goal \
  --title "Example goal" \
  --source "User asked for the example outcome" \
  --next "Read the relevant files" \
  --criteria "User-visible outcome is complete and validated" \
  --plan "docs/implementation-plan.md#example" \
  --todo "implement-example"
```

更新、拆解、暫停與完成：

```bash
./scripts/agent-goals.sh --project <PROJECT_ROOT> update --id P1-example-goal --note "Read dependencies" --next "Implement the change"
./scripts/agent-goals.sh --project <PROJECT_ROOT> update --id P1-example-goal --missing "Validation examples are not written" --decision "Choose whether this remains P1" --strengthen "Add stronger completion criteria"
./scripts/agent-goals.sh --project <PROJECT_ROOT> split --parent P1-example-goal --id P2-child-goal --title "Child goal"
./scripts/agent-goals.sh --project <PROJECT_ROOT> pause --id P1-example-goal --reason "User changed priority"
./scripts/agent-goals.sh --project <PROJECT_ROOT> complete --id P1-example-goal --validated --note "Validation passed"
```

安全條件：

- `complete` 只有在傳入 `--validated` 時才會刪除 goal 檔；否則會保留並標成 `needs-validation`。
- `.agent-goals/README.md` 會自動刷新成主目標表，連到 `goals/*.md`，並顯示 open work / decisions、plan/todo links、下一步與更新時間。
- `start`、`update`、`split` 可重複使用 `--plan` 與 `--todo`，把 planning 文件章節、TodoWrite ID、checklist item 或 issue ID 連到 goal。
- `update` 可用 `--missing`、`--decision`、`--strengthen` 把未完成、待決策與待補強項目放進主表。
- 每個 goal 更新時會使用 `.agent-goals/locks/<goal-id>.lock/` 防止多 agent 同時寫入。
- Stale lock 可用 `cleanup` 清理；TTL 預設 30 分鐘，可用 `AGENT_GOALS_LOCK_TTL_SECONDS` 覆寫。
