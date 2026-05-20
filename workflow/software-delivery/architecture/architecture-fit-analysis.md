# Architecture Fit Analysis Workflow

## 目的

在 software-delivery 任務中，先用 evidence 判斷架構 fit，再提出 architecture recommendation。

## 步驟

1. 讀取 change brief、product brief、BDD、domain contract 或 design note。
2. 評估 `domain_complexity`、`invariant_density`、`business_language_instability`、`workflow_complexity`、`integration_pressure`、`lifecycle_length`、`team_scale`、`bounded_context_count`、`event_coordination_need` 與 `delivery_speed_priority`。
3. 對照 `metadata/architecture/architecture-fit-matrix.yaml` 選出最小可行架構。
4. 若使用者提到 DDD，先判斷 DDD Lite 或 Full DDD，而不是預設 full DDD。
5. 輸出 chosen strategy、rejected lighter option、rejected heavier option、fit evidence、validation plan 與 upgrade/downgrade trigger。

## Gate

若建議 CQRS、event sourcing、Full DDD 或 microservices，必須先完成 overengineering review。
