# State Visibility Gap

**Status**: `candidate-intelligence`

## 定義

State visibility gap 是一種 validation failure pattern：測試或產品表面觀察到的狀態，與系統真實狀態不一致。

核心風險不是「沒有測試」，而是測試看到的成功訊號只覆蓋局部狀態，沒有證明狀態已沿著完整鏈路傳播到最終可觀察結果。

## 常見形狀

| 觀察到的狀態 | 可能缺失的真實狀態 |
| --- | --- |
| API 回 `200` | 資料庫沒有寫入、事件沒有送出、權限沒有更新 |
| Adapter 回 success | 外部系統沒有接受、使用者沒有收到、business effect 沒有發生 |
| UI 顯示成功 | persisted state、SSR state、entitlement state 不一致 |
| 測試環境可見 | production audience、tenant、feature flag 或 geo scope 不可見 |

## 風險維度

不要把每個案例命名成單獨 subtype。先拆成兩個維度：

```yaml
state_visibility_gap:
  state_sources:
    - identity
    - entitlement
    - tenant
    - ownership
    - feature_flag
    - payment_state

  propagation_modes:
    - synchronous
    - asynchronous
    - external_system
    - persisted_state
    - rendered_ui
```

`state_sources` 表示狀態從哪裡來；`propagation_modes` 表示狀態如何跨邊界傳播。

## 驗證要求

命中 state visibility gap 時，agent 必須：

1. 描述觀察到的成功訊號。
2. 描述該訊號覆蓋的 scope。
3. 找出尚未驗證的真實狀態。
4. 依 [`evidence-chain-validation.md`](evidence-chain-validation.md) 補完整狀態傳播鏈證據。
5. 依 [`evidence-depth.md`](evidence-depth.md) 選擇最低 evidence depth。

## 相關知識

- [`evidence-model.md`](evidence-model.md)
- [`evidence-depth.md`](evidence-depth.md)
- [`evidence-chain-validation.md`](evidence-chain-validation.md)
- [`../../../anti-patterns/mock-completeness-illusion.md`](../../../anti-patterns/mock-completeness-illusion.md)
- [`../../../anti-patterns/validation-proxy-trap.md`](../../../anti-patterns/validation-proxy-trap.md)
