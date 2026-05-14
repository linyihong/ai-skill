# 可重用規則與專案證據邊界

本檔是 **Ai-skill repository 全庫鐵則**：任何 reusable skill、shared rule、template、feedback lesson、checklist 或 workflow，都只能沉澱**可重用原因、決策規則、驗證方法與適用邊界**。專案 incident 的具體證據留在專案文件、issue、runbook 或 integration notes，不進 reusable skill 正文。

## 核心規則

| 內容 | 放置位置 |
| --- | --- |
| 可重用失敗模式、原因分析、決策規則、驗證流程、適用 / 不適用條件 | `enforcement/` 或對應 skill 文件 |
| 具體 App / 專案名稱、module 名稱、class/test 名稱、endpoint、host、payload fragment、sample ID、帳號/裝置/環境細節、一次 live run 結果 | 業務專案文件、issue、runbook、integration notes |
| 從 incident 推導出的 lesson | `feedback_history/`，但正文只寫 generalized lesson；incident 原文與可識別資訊不可複製 |

## 必須先分析原因

當使用者指出 skill / rule 更新「閉環不完整」或「寫進了專案特例」時，agent 不可只刪掉字串或只補一個檔案。必須在同一輪完成：

1. **原因分析**：說明為什麼漏掉閉環，例如只更新局部 skill、沒有查 promotion target、沒有依 `linked-updates.md` 檢查、沒有搜尋殘留專案字串。
2. **規則強化**：把防止重犯的規則補到正確層級。若是全庫行為，放 `enforcement/`；若是單一 skill 行為，才放 skill 文件。
3. **失效學習**：若該問題可能重複發生，依 [`failure-learning-system.md`](failure-learning-system.md) 分類失效模式，必要時新增或更新 `enforcement/failure-patterns/` 或對應 skill 的 `feedback_history/`。
4. **連動更新**：依 [`linked-updates.md`](linked-updates.md) 同步 README、入口、checklist、templates、feedback index 或受影響 skill。
5. **驗證**：搜尋 reusable docs 是否仍含專案特例；反查新規則是否可從 shared rule index 與相關 skill 入口找到。

## 寫作檢查

寫入 reusable 文件前，先問：

- 這句話是否換到另一個專案仍然成立？
- 是否包含某個專案的 name/path/endpoint/class/sample/live result？
- 是否把「這次看到的現象」誤寫成全域規則？
- 是否只描述修正結果，卻沒有描述可重用原因與驗證方法？
- 是否已把具體證據留在專案文件，而非 skill/shared rule？

任一答案顯示內容仍綁定專案，就不要寫進 reusable skill；改寫成 generalized rule，並把具體證據移到專案文件。

## 與其他規則的關係

- 去敏與 placeholder 依 [`sanitization.md`](sanitization.md)。
- 內容放置位置依 [`content-layering.md`](content-layering.md)。
- 新增或推廣 lesson 依 [`feedback-lessons.md`](feedback-lessons.md)。
- 重複 agent 失效模式的分類與 prevention gate 依 [`failure-learning-system.md`](failure-learning-system.md)。
- 任何跨文件影響依 [`linked-updates.md`](linked-updates.md)。
- 每個重要結論需有目標、執行與驗證，依 [`goal-action-validation.md`](goal-action-validation.md)。

← [回到共用規則索引](README.md)
