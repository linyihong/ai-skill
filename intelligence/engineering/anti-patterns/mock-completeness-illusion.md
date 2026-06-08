# Mock Completeness Illusion

**Status**: `candidate-intelligence`

## 定義

Mock completeness illusion 是一種工程 anti-pattern：unit、domain service、adapter、API mock 或 fixture 測試看起來完整，但它們只證明局部邏輯，沒有證明產品層 claim 成立。

問題不是測試數量不足，而是 evidence depth 不足。

## 常見症狀

- 所有 unit tests 通過，但 live flow 失敗。
- Domain service test 證明規則正確，但身份、cookie、tenant、feature flag 或 entitlement 沒被真實套用。
- Adapter test 證明 wrapper 會回 success，但外部 side effect 沒發生。
- API mock test 證明 response shape 正確，但 DB、queue、email、storage、SSR/UI 沒有一致性證據。

## 為什麼危險

Mock tests 常常固定了理想世界：

```text
mocked precondition
→ local function success
→ asserted local result
```

真實產品通常是：

```text
identity / state source
→ boundary
→ side effect
→ persisted or external state
→ readback
→ user observable result
```

兩者驗證的 claim 不同。

## 修正方式

- 依 [`../execution/validation-reasoning/evidence-depth.md`](../execution/validation-reasoning/evidence-depth.md) 判斷最低 evidence depth。
- 命中 state visibility gap 時，依 [`../execution/validation-reasoning/evidence-chain-validation.md`](../execution/validation-reasoning/evidence-chain-validation.md) 補完整鏈路。
- 保留 mock tests 作為低層回歸，但不要用它們支持 live system claim。

## 相關知識

- [`../execution/validation-reasoning/state-visibility-gap.md`](../execution/validation-reasoning/state-visibility-gap.md)
- [`../execution/validation-reasoning/evidence-model.md`](../execution/validation-reasoning/evidence-model.md)
- [`validation-proxy-trap.md`](validation-proxy-trap.md)
