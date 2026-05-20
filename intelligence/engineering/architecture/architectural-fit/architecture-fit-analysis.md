# Architecture Fit Analysis

**Status**: `candidate-intelligence`

## 判斷原則

Architecture fit analysis 先評估問題形狀，再選架構。它禁止 `project => full DDD` 這種預設推論。

## 評估軸

| Axis | 問題 |
| --- | --- |
| `domain_complexity` | 業務規則是否超過 CRUD？ |
| `invariant_density` | 不變量是否關鍵且常被破壞？ |
| `business_language_instability` | 詞彙是否常被不同角色誤用？ |
| `workflow_complexity` | 流程是否跨多步、多角色、多狀態？ |
| `integration_pressure` | 外部模型是否會污染內部模型？ |
| `lifecycle_length` | 系統是否長期維護？ |
| `team_scale` | 是否多團隊或多人長期協作？ |
| `bounded_context_count` | 是否有多個 subdomain？ |
| `event_coordination_need` | 是否需要非同步協調？ |
| `delivery_speed_priority` | MVP 速度是否高於長期模型純度？ |

## 輸出格式

```text
Recommended architecture: <strategy>
Why this fits: <evidence>
Rejected heavier option: <why>
Rejected lighter option: <why>
Upgrade trigger: <future signal>
Validation: <how to confirm>
```
