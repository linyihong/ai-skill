# Intelligence

`intelligence/` 負責「沉澱工程智慧與領域知識」。本層保存從 analysis、workflow 執行與 feedback 中抽出的可重用判斷，讓 agent 能引用工程 decision、trade-off、anti-pattern 與 domain knowledge。

## 放什麼

- 工程決策原則、trade-off 與架構 lesson。
- 可跨專案重用的 domain knowledge。
- 失效模式、anti-pattern 與改善策略的抽象結論。
- 從分析證據萃取出的穩定判斷。

## 不放什麼

- 觀察與拆解的原始方法；放到 `analysis/`。
- 逐步執行流程、review flow 或 task orchestration；放到 `workflow/`。
- 對話暫存 goal、目前 owner 或 next action；放到 `.agent-goals/`。
- 可執行 policy 與 close-loop gate；放到 `shared-rules/`。

## 與既有層的關係

- `skills/` 目前仍提供能力入口；成熟的工程智慧可逐步抽到本層。
- `workflow/` 應 reference 本層，而不是內嵌大量知識。
- `feedback/` 可把新 lesson promotion 到本層。
- `governance/` 定義本層知識的 lifecycle、清理與 validation。

## 第一批候選遷移來源

- `skills/app-development-guidance/implementation/`
- `skills/app-development-guidance/controls/`
- `skills/*/feedback_history/` 中已成熟且跨專案可重用的 lesson
- `shared-rules/failure-patterns/` 中偏工程判斷的 pattern 摘要
