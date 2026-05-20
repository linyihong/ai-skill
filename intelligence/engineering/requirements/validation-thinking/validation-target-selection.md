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
