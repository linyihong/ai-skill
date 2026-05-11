# Feedback

`feedback/` 負責「系統如何持續演化」。本層保存 feedback replay、lesson extraction、refinement 與 promotion 的系統設計，讓一次性觀察能被驗證後推進到 workflow、intelligence、memory 或 shared rules。

## 放什麼

- Feedback extraction 與 replay 的流程設計。
- Lesson promotion、refinement 與退回條件。
- 如何把 agent failure、review comment 或實作經驗轉成可重用知識。
- Feedback loop 與 `analysis/`、`intelligence/`、`workflow/` 的關係。

## 不放什麼

- 單一 skill lesson 正文；仍放在 `skills/<skill>/feedback_history/`，直到遷移策略穩定。
- 全庫可執行 feedback lesson 檔名與模板規則；放到 `shared-rules/feedback-lessons.md`。
- Failure pattern 的可執行 prevention gate；放到 `shared-rules/failure-patterns/` 或相關 shared rule。
- 專案 raw evidence 或未去敏資料；留在業務專案。

## 與既有層的關係

- `skills/*/feedback_history/` 仍是目前 lesson storage 的相容層。
- `shared-rules/failure-learning-system.md` 與 `shared-rules/feedback-lessons.md` 仍是可執行規則。
- `intelligence/` 承接成熟的工程智慧。
- `governance/` 定義 promotion、deprecation 與 validation lifecycle。

## 第一批候選遷移來源

- `shared-rules/failure-learning-system.md`
- `shared-rules/feedback-lessons.md`
- `skills/*/feedback_history/README.md`
- `architecture/next-stage-upgrade-plan.md` 的 Intelligence Feedback Loop
