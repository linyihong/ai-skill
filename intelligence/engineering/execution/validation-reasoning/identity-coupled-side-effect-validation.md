# Identity-Coupled Side Effect Validation

**Status**: `candidate-intelligence`

## 定義

Identity-coupled side effect validation 是 [`evidence-chain-validation.md`](evidence-chain-validation.md) 的特化案例：產品行為同時依賴呼叫者身份狀態，並產生可觀察 side effect。

它不應成為一級 workflow rule。先用通用 pattern 判斷，再在身份與副作用交會時套用此特化。

## 觸發條件

同時出現：

```yaml
validation_risk:
  - state_visibility_gap
state_sources:
  - identity
  - entitlement
  - tenant
  - ownership
propagation_modes:
  - persisted_state
  - rendered_ui
  - external_system
```

## 最低證據鏈

典型鏈路：

```text
Identity Material
→ Authenticated Product API
→ Permission / Entitlement Decision
→ Side Effect
→ Persisted or External State
→ SSR/API Readback
→ User Observable UI
```

最低要求：

- 使用真實身份材料，例如產品 API 取得的 session、cookie、token 或等價測試身份。
- 走產品正式 API 或 UI/H5 path，不只呼叫 service 或 adapter。
- 驗證 side effect 的 durable state 或 independent observation。
- 驗證 readback path 與 UI/SSR/API state 一致。
- 覆蓋至少 unauthenticated、authenticated without entitlement、authenticated with entitlement；若功能有空資料狀態，也覆蓋 missing/empty data。

## 不足證據

- Login API success alone.
- Service/domain test alone.
- Adapter success alone.
- API response without persisted or user-observable readback.
- UI text without backend state consistency.

## 相關知識

- [`state-visibility-gap.md`](state-visibility-gap.md)
- [`evidence-model.md`](evidence-model.md)
- [`evidence-depth.md`](evidence-depth.md)
- [`evidence-chain-validation.md`](evidence-chain-validation.md)
- [`../../../anti-patterns/validation-proxy-trap.md`](../../../anti-patterns/validation-proxy-trap.md)
