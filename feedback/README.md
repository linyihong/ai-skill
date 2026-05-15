# Feedback

`feedback/` 負責「系統如何持續演化」。本層保存 feedback replay、lesson extraction、refinement 與 promotion 的系統設計，讓一次性觀察能被驗證後推進到 workflow、intelligence、memory 或 shared rules。

## 放什麼

- Feedback extraction 與 replay 的流程設計。
- Lesson promotion、refinement 與退回條件。
- 如何把 agent failure、review comment 或實作經驗轉成可重用知識。
- Feedback loop 與 `analysis/`、`intelligence/`、`workflow/` 的關係。

## 不放什麼

- 單一 lesson 正文；統一放在 `feedback/history/<domain>/`。
- 全庫可執行 feedback lesson 檔名與模板規則；放到 `enforcement/feedback-lessons.md`。
- Failure pattern 的可執行 prevention gate；放到 `enforcement/failure-patterns/` 或相關 shared rule。
- 專案 raw evidence 或未去敏資料；留在業務專案。

## 與既有層的關係

- `feedback/history/` 是 lesson 的統一目標路徑。
- `skills/*/feedback_history/` 已於 2026-05-13 刪除（apk-analysis、app-development-guidance、travel-planning），所有 lesson 已搬遷至 `feedback/history/`。
- `enforcement/failure-learning-system.md` 與 `enforcement/feedback-lessons.md` 仍是可執行規則。
- `intelligence/` 承接成熟的工程智慧。
- `governance/` 定義 promotion、deprecation 與 validation lifecycle。
- `knowledge/runtime/sqlite/` 可索引冷 feedback lessons，協助低 token 查找；lesson 全文仍留在 `feedback/history/`。

## 目前入口

- [`history/`](history/README.md)：所有 lesson 的統一存放路徑（依 domain 分類）。
- [`replay/`](replay/README.md)：定義經驗重播的系統設計（從過往 session、failure 與成功經驗中提取可重複使用的教訓）。
- [`extraction/`](extraction/README.md)：定義智慧抽取的系統設計（從 feedback lesson、replay 結果中提煉 intelligence atom）。
- [`refinement/`](refinement/README.md)：定義流程精煉的系統設計（從實作經驗中持續改進 workflow）。
- [`promotion/`](promotion/README.md)：定義 feedback lesson 如何從 history 推進到 workflow、intelligence、shared rules、memory 或 runtime surfaces。

## 誰會參考這裡（Inbound References）

- [`route.feedback.promotion-pipeline`](../knowledge/runtime/routing-registry.yaml:290) — candidate_sources 引用 `feedback/README.md`
- [`route.feedback.history`](../knowledge/runtime/routing-registry.yaml:664) — candidate_sources 引用 `feedback/README.md`
- [`enforcement/feedback-lessons.md`](../enforcement/feedback-lessons.md) — 定義 feedback lesson 的寫作規則與分類
- [`governance/lifecycle/intelligence-extraction-pipeline.md`](../governance/lifecycle/intelligence-extraction-pipeline.md) — Step 4 使用 feedback history 做 extraction

## 第一批候選遷移來源

- `enforcement/failure-learning-system.md`
- `enforcement/feedback-lessons.md`
- ✅ `skills/*/feedback_history/`（已於 2026-05-13 全部搬遷至 `feedback/history/`，舊目錄已刪除）
- `plans/archived/2026-05-11-next-stage-upgrade-plan.md` 的 Intelligence Feedback Loop
