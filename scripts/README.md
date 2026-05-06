# Scripts

| 檔案 | 用途 |
| --- | --- |
| [`sync-cursor-bundle.sh`](sync-cursor-bundle.sh) | 把本庫同步到 **`~/.cursor/bundles/shared-rules`**（共用規則）與 **`~/.cursor/bundles/ai-skill/`**（各 skill），再讓 `~/.cursor/shared-rules`、`~/.cursor/skills/<name>` 指向 bundle（與其他 `.cursor` 內容分流）。 |
| [`ai-skill-close-loop.sh`](ai-skill-close-loop.sh) | 保守收尾工具：偵測 active close-loop lock、列出 dirty 檔案 owner group；預設 dry-run，`--commit` 才分組提交，`--push` 才推遠端。 |
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
