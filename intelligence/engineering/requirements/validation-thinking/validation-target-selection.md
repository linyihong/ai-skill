# Validation Target Selection

**Status**: `candidate-intelligence`

## 判斷原則

Validation target 應由風險決定，而不是由工具方便性決定。

| 風險 | 合適 target |
| --- | --- |
| simple pure logic | unit / table-driven test |
| API / schema compatibility | contract test / generated client check |
| behavior flow | BDD / scenario / integration test |
| domain invariant | invariant / property-based / mutation test |
| external service | fixture-backed or gated live evidence |
| hardware / manual workflow | manual evidence + checklist |

## 測試有效性升級

當 coverage 高但 correctness confidence 低，或變更涉及 AI-generated logic、critical branch、permission/security、domain invariant、refactor-no-behavior-change claim 時，使用 [`mutation-testing-effectiveness.md`](mutation-testing-effectiveness.md) 做 targeted check。

Mutation testing 的目標不是追求 mutant 數量，而是確認測試至少能殺掉有價值的錯誤版本；若 mutant survived，應補 validation target 或縮小完成宣告。
