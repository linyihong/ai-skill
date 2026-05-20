# Transactional Consistency

**Status**: `candidate-intelligence`

## 判斷原則

Consistency model 應由 invariant 決定。不是每個流程都需要 ACID，也不是每個跨 context 流程都能接受 eventual consistency。

## 檢查

| 問題 | 決策 |
| --- | --- |
| 不變量是否必須即時成立？ | 若是，放在同一 consistency boundary。 |
| 延遲一致是否對使用者或業務可接受？ | 若是，可使用 event / saga / compensation。 |
| 失敗後是否能補償？ | 若不能，避免過早非同步化。 |
| 事件順序是否影響 correctness？ | 若是，需明確 ordering / idempotency。 |

## 與 BDD 的關係

BDD 描述 observable behavior；consistency model 解釋該 behavior 需要哪種 domain correctness 保證。
