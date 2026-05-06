# 同步到 Cursor

## 本庫結構（記憶用）

- **`shared-rules/`**：分類後的共用規則（正本）。
- **`skills/<name>/`**：各 skill（例如 `skills/apk-analysis/` 內含 `SKILL.md`、`feedback_history/` 等）。

## 建議佈署順序

部署到 `<PROJECT_ROOT>/.cursor/` 時：

1. **複製整個** `shared-rules/` → `<PROJECT_ROOT>/.cursor/shared-rules/`（或你與團隊約定的固定路徑，與 skill 並列即可）。
2. 再將 **`skills/apk-analysis/`** 複製或 symbolic link 到 **`<PROJECT_ROOT>/.cursor/skills/apk-analysis/`**（Cursor 慣用 skill 掃描路徑）。

這樣 Agent 同時讀得到「分類後共用規則」與「apk-analysis 技巧包」，無須把共用條文拆進每一則技巧檔。

## Symbolic link 注意

若 `.cursor/skills/apk-analysis` 連結到本庫的 `skills/apk-analysis`，**仍須**另行同步 **`shared-rules/`**（連結不會自動包含上一層的共用目錄）。

## 本機共用：建議用 `bundles/` 並列放（與其他規則隔離）

若 **`shared-rules/`** 與各 **skill** 算「同一套共用資產」，又不希望和別的工具或手動建立的 `~/.cursor/*` 混在一起，建議在本機 **`~/.cursor/bundles/`** 下**分兩條**（都用 symlink 指回 `<AI_SKILL_REPO>`）：

```text
~/.cursor/bundles/shared-rules  -> <AI_SKILL_REPO>/shared-rules

~/.cursor/bundles/ai-skill/     # 只放各 skill，不放 shared-rules
  apk-analysis   -> <AI_SKILL_REPO>/skills/apk-analysis
  <其他 skill>/  -> <AI_SKILL_REPO>/skills/<name>/
```

再讓 Cursor 慣用路徑**經 bundle 轉接**（若要一次還原，只要認 `bundles/` 底下這兩塊）：

```text
~/.cursor/shared-rules          -> ~/.cursor/bundles/shared-rules
~/.cursor/skills/apk-analysis  -> ~/.cursor/bundles/ai-skill/apk-analysis
```

這樣 **共用規則**與 **skill 包**在 `bundles/` 裡**並列、語意分開**，比較不容易被「同一個資料夾裡混了別種連結」搞混；`~/.cursor/` 其他設定仍獨立在外。

本庫提供可重複執行的腳本：**[`scripts/sync-cursor-bundle.sh`](../scripts/sync-cursor-bundle.sh)**（見 [`scripts/README.md`](../scripts/README.md)）。

## 改動 `shared-rules/` 或 `skills/` 之後

凡編輯過這兩處且希望 **本機 Cursor** 立刻經 `~/.cursor/bundles` 跟上時，請在 `<AI_SKILL_REPO>` 執行 `./scripts/sync-cursor-bundle.sh`（可重複執行、無害）。可選：於 repo 根目錄 `git config core.hooksPath scripts/git-hooks`，則每次 **`git commit`** 後會自動跑該腳本。若 skill 清單仍未更新，再 **Developer: Reload Window**。

同步與 `git push` 完成後，agent 還必須重新讀取本次更新過的 skill/shared-rule 入口與主要依賴文件。`sync-cursor-bundle.sh` 只更新檔案路徑，不會自動更新 agent 已載入的上下文；讀回 gate 依 [`dependency-reading.md`](dependency-reading.md) 執行。

## 為什麼會看到「兩個」shared-rules？

**沒有兩份內容。**真實檔案只在 **`<AI_SKILL_REPO>/shared-rules/`** 這一處。

本機會出現兩個路徑名稱，是因為 **symlink 鏈**（接力指向同一資料夾）：

1. **`~/.cursor/bundles/shared-rules`** — 直接連到 repo（bundle 裡的「共用規則」條目，和 `ai-skill/` 並列）。
2. **`~/.cursor/shared-rules`** — 連到 **`bundles/shared-rules`**，給習慣從 `~/.cursor/shared-rules` 讀取的流程用。

兩者 **`realpath` 相同**，編輯任一邊看到的都是同一套檔案。若你希望畫面上少一層目錄概念，也可以自行改成 **只有** `~/.cursor/shared-rules` 直接指到 repo（就不經 bundles）；缺點是 bundle 底下少了獨立的 `shared-rules` 條目，與「規則與 skill 都放在 bundles 並列」的整理方式不一致。

**警告：**請勿在 **本庫**裡再建立「指回 `~/.cursor/bundles/...`」的同名連結，會與 bundle **繞圈**，IDE 會像無限巢狀：

- 不要在 **`shared-rules/`** 底下再放 **`shared-rules`**。
- 不要在 **`skills/apk-analysis/`** 底下再放 **`apk-analysis`**（其他 skill 同理：**`skills/<name>/` 內不要同名 `<name>`**）。

正確做法：**只在 `~/.cursor` 側**用腳本或手動建立 symlink；repo 內只放正文與 `feedback_history/`。本庫根目錄 **`.gitignore`** 已忽略常見誤建路徑。

## 疑義時

以本庫 **`shared-rules/`** 與 **`skills/`** 為準；`.cursor` 內皆為同步產物。

← [回到共用規則索引](README.md)
