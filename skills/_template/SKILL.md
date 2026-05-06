---
name: <skill-kebab-name>
description: <英文：Agent 如何辨識何時套用此 skill。寫觸發情境、允許範圍、不要寫特定產品名或機密。>
---

# <Skill 標題（可中英並列）>

**共用政策：**先讀 [`shared-rules` 索引](../../shared-rules/README.md)、[`feedback-lessons.md`](../../shared-rules/feedback-lessons.md)、[`neutral-language.md`](../../shared-rules/neutral-language.md)、[`goal-action-validation.md`](../../shared-rules/goal-action-validation.md)、[`document-sizing.md`](../../shared-rules/document-sizing.md) 與 [`cross-skill-references.md`](../../shared-rules/cross-skill-references.md)。`feedback_history/` 內各 lesson 只引用、不重複貼全文；文件標題、檔名、slug、索引與摘要使用中性低爭議用語；每個重要工作單元要能反查目標、執行、驗證或參考來源；文件變大時改成資料夾與 `README.md` 目錄拆分。若本 skill 需要引用另一個 skill 的規範、模板或交接產物，只寫觸發條件、交接 artifact 與 ownership boundary，不複製對方全文。

## 何時使用本 skill

- 情境一（例：使用者提到 …）
- 情境二

## 不在範圍內

- 例：未授權目標、個資、特定產品機密結論（應放在業務專案文件）

## Quick Start

1. …
2. …

## 預設流程

- 詳見 [WORKFLOW.md](WORKFLOW.md)（若尚未建立，可先在本節寫精簡步驟，日後再拆檔）
- 若某段流程開始變長或分出多個分類，依 `shared-rules/document-sizing.md` 拆成子資料夾與目錄索引

## Cross-Skill References（可選）

- 若需要引用其他 skill，依 [`shared-rules/cross-skill-references.md`](../../shared-rules/cross-skill-references.md) 寫清楚：target skill、觸發條件、交接 artifact、ownership boundary、去敏邊界與 linked updates。

## 工具與環境

- 詳見 [TOOLS.md](TOOLS.md)（可選）

## 文件化與產出

- 詳見 [DOCUMENTATION.md](DOCUMENTATION.md)（可選）

## 回饋與技巧沉澱

- 新 lesson：未分類 skill 用 `feedback_history/YYYY-MM-DD_HHMMSS-<slug>.md`；已分類 skill 用 `feedback_history/<category>/YYYY-MM-DD_HHMMSS-<slug>.md`，跨分類用 `feedback_history/common/`（規則見 [`shared-rules/feedback-lessons.md`](../../shared-rules/feedback-lessons.md)）
- 可選入口：[FEEDBACK.md](FEEDBACK.md)（僅連到上述檔）

## 新專案／Runbook（可選）

- [RUNBOOK.md](RUNBOOK.md) — 第一天怎麼套用、如何回饋中央庫
