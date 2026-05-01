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

## 疑義時

以本庫 **`shared-rules/`** 與 **`skills/`** 為準；`.cursor` 內皆為同步產物。

← [回到共用規則索引](README.md)
