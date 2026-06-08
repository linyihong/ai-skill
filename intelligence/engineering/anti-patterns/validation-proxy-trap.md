# Validation Proxy Trap

**Status**: `candidate-intelligence`

## 定義

Validation proxy trap 是把代理訊號誤認成最終證據的 anti-pattern。

代理訊號可以是真實的，但它只證明某一段成功，不證明 business effect 或 user observable result 已發生。

## 常見錯誤

| Proxy signal | 不等於 |
| --- | --- |
| SMTP adapter success | 使用者收到 email |
| Upload API success | 檔案可被正確讀取、權限與 metadata 正確 |
| Queue publish success | Consumer 已處理且 business state 已改變 |
| Payment API accepted | 金流 settled、webhook processed、entitlement granted |
| Login API success | cookie/session 可存取 protected resource |
| API `200` | DB row、SSR state、UI CTA 與權限狀態一致 |

## 判斷規則

如果證據只來自同一層的 success response，且 claim 涉及外部世界、持久化、非同步處理、身份權限或 UI readback，就不能宣稱完成。

需要補：

- persisted state evidence
- external confirmation
- consumer/readback evidence
- user observable state
- independent observation for critical paths

## 修正方式

1. 用 [`../execution/validation-reasoning/evidence-model.md`](../execution/validation-reasoning/evidence-model.md) 標出 proxy signal 的 scope。
2. 用 [`../execution/validation-reasoning/evidence-chain-validation.md`](../execution/validation-reasoning/evidence-chain-validation.md) 列出缺少的下游 segment。
3. 用 [`../execution/validation-reasoning/evidence-depth.md`](../execution/validation-reasoning/evidence-depth.md) 決定是否需要 `live_system` 或 `independent_observation`。

## 相關知識

- [`mock-completeness-illusion.md`](mock-completeness-illusion.md)
- [`../execution/validation-reasoning/state-visibility-gap.md`](../execution/validation-reasoning/state-visibility-gap.md)
