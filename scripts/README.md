# Scripts

| 檔案 | 用途 |
| --- | --- |
| [`sync-cursor-bundle.sh`](sync-cursor-bundle.sh) | 把本庫同步到 **`~/.cursor/bundles/shared-rules`**（共用規則）與 **`~/.cursor/bundles/ai-skill/`**（各 skill），再讓 `~/.cursor/shared-rules`、`~/.cursor/skills/<name>` 指向 bundle（與其他 `.cursor` 內容分流）。 |

在本庫根目錄執行：

```bash
chmod +x scripts/sync-cursor-bundle.sh   # 只需做一次
./scripts/sync-cursor-bundle.sh
```
