# Feedback History

所有 feedback lesson 統一集中在此目錄，不再分散到各層（`workflow/`、`analysis/`、`intelligence/`）底下的 `feedback_history/`。

## 結構

```text
feedback/history/
  README.md                     # 本檔
  <domain>/                     # 對應 skill 領域名稱
    README.md                   # 該領域 lesson 索引
    common/                     # 跨分類或全域 lesson
    <category>/                 # 分類 lesson（依 skill 既有分類）
      YYYY-MM-DD_HHMMSS-<slug>.md
```

## 目前 domains

| Domain | 對應 skill | 狀態 |
|--------|-----------|------|
| `apk-analysis/` | APK 分析 | lesson 已就位 |
| `app-development-guidance/` | 應用開發指引 | lesson 已就位 |
| `travel-planning/` | 旅遊規劃 | 尚無 lesson |

## 規則

- 新 lesson 一律寫入 `feedback/history/<domain>/` 對應分類，**不得**再寫入 `skills/<name>/feedback_history/` 或 `workflow/<name>/feedback_history/` 等舊路徑。
- 檔名規則、模板、agent 行為見 [`shared-rules/feedback-lessons.md`](../../shared-rules/feedback-lessons.md)。
- 若 domain 下尚無對應分類目錄，**應主動建立**，而非退回舊路徑。
- 舊 `skills/<name>/feedback_history/` 中的 lesson 將分批搬遷至此，搬遷完成後刪除舊目錄。

## 與其他層的關係

- `feedback/extraction/`：存放從舊 skill 提取的 lesson 索引，指向 `feedback/history/` 的新位置。
- `shared-rules/feedback-lessons.md`：定義 lesson 的檔名規則、模板、agent 行為。
