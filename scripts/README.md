# Scripts

| 檔案 | 用途 |
| --- | --- |
| [`sync-cursor-bundle.sh`](sync-cursor-bundle.sh) | 把本庫同步到 **`~/.cursor/bundles/shared-rules`**（共用規則）與 **`~/.cursor/bundles/ai-skill/`**（各 skill），再讓 `~/.cursor/shared-rules`、`~/.cursor/skills/<name>` 指向 bundle（與其他 `.cursor` 內容分流）。 |
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
