# Architecture Overengineering Detection

## 目的

防止 AI architecture inflation，確保架構複雜度與 business complexity 對齊。

## 檢查清單

| Signal | Required review |
| --- | --- |
| Aggregate count 快速增加 | 檢查每個 aggregate 的 invariant。 |
| CQRS 無 scale / read-write divergence | 合併 command/query model。 |
| Event sourcing 無 audit / replay need | 改用 state transition log 或普通 persistence。 |
| Repository 單一 implementation | 確認是否真的保護 aggregate boundary。 |
| Microservice 無 ownership / deployment boundary | 回到 modular monolith。 |
| Abstraction count > domain concept count | 移除 speculative abstraction。 |

## 輸出

Overengineering review 必須提出：keep / simplify / defer 三種 verdict 之一。
