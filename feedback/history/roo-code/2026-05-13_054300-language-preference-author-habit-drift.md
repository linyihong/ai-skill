# 語言偏好漂移（作者習慣型）

## 發生時間

2026-05-13 13:42 (JST)

## 情境

使用者在確認 `vscode-extension-global-state.md` 昇華到 `intelligence/ide/` 的實作結果後，指出 `attempt_completion` 的內容混用了中英文（表格欄位用英文、分析段落中關鍵詞用英文）。

## 當下使用的工具

Roo Code（Code mode），deepseek-chat model

## 分析

這次的語言飄移**不是** Roo Code system prompt 的問題，也不是模型的問題。原因是：

1. 使用者的提問是中文，agent 也用中文回應
2. 但在撰寫 `attempt_completion` 時，agent 不自覺地在表格欄位使用英文（Change、Description）
3. 在分析「為什麼當初沒建議」的段落中，關鍵概念用英文表達（cognitive bias、reusable knowledge vs tool config）
4. 這是 agent 的**寫作慣性**：訓練資料中技術分析類內容大量使用英文，導致在「寫作模式」下回歸到統計慣性

## 與前一次語言漂移的差異

| 面向 | 前一次（2026-05-13 12:48） | 本次（2026-05-13 13:42） |
|------|---------------------------|--------------------------|
| 類型 | system prompt 強制（Type A） | 作者習慣漂移（Type B） |
| 原因 | Roo Code 的 `language: "en"` 欄位 + Custom Instructions 固定英文偏好 | Agent 寫作慣性，訓練資料的統計特性 |
| 表現 | 整個回應都用英文 | 對話是中文，但結構化輸出混入英文 |
| 解決 | 修改 SQLite language 欄位 + `.roomodes` 軟性偏好 | 需要在輸出前增加語言一致性檢查 |

## 教訓

1. 即使對話語言正確，agent 在產生結構化輸出（表格、分析報告）時仍可能漂移
2. 需要在 Custom Instructions 中加入「所有輸出（包含表格、標題、分析報告）都必須與使用者當前語言一致」
3. 使用者在發現漂移時應立即糾正，避免慣性持續

## 相關文件

- Failure pattern: [`enforcement/failure-patterns/language-preference-drift.md`](../../enforcement/failure-patterns/language-preference-drift.md)（已新增子類型 B）
- 前一次 lesson: [`feedback/history/roo-code/2026-05-13_124800-language-preference-drift.md`](./2026-05-13_124800-language-preference-drift.md)
