# Evidence Chain Validation

**Status**: `candidate-intelligence`

## 定義

Evidence chain validation 是一種 validation reasoning pattern：驗證必須沿著狀態傳播鏈一路追到最終可觀察結果，並在每個重要 state transition 上附上足夠證據。

它回答的問題不是「有沒有測試」，而是「證據是否覆蓋了 claim 所依賴的整條狀態鏈」。

## 基本形狀

```text
state source
→ propagation step
→ persisted or external state
→ readback path
→ user/business observable result
```

範例鏈路：

```text
Login → Session Created → Cookie Issued → Protected Resource → UI State
API → Domain Logic → SMTP Provider → Inbox → User Observable Email
Order → Gateway → Webhook → DB → Entitlement → UI
Producer → Topic → Consumer → DB → Business Effect
Flag Rule → Audience Match → Frontend Fetch → Rendered Feature
```

## 操作步驟

1. **列 claim**：要證明的產品行為是什麼。
2. **列 chain**：從狀態來源到最終可觀察結果，逐段列出 propagation path。
3. **標 evidence**：每段現有證據的類型、scope 和 confidence，依 [`evidence-model.md`](evidence-model.md) 判斷。
4. **找缺口**：任何沒有 readback、persisted state、external confirmation 或 user-observable state 的段落都不能被視為完整。
5. **選 depth**：依 [`evidence-depth.md`](evidence-depth.md) 決定是否需要 `live_system` 或 `independent_observation`。

## 判斷規則

- API success 只覆蓋 API segment，不自動覆蓋 DB、queue、external system 或 UI。
- Adapter success 只覆蓋 adapter segment，不自動覆蓋外部世界。
- UI success 若沒有 persisted/readback evidence，可能只是 presentation state。
- DB success 若沒有 user-observable readback，可能仍不支持完整 product claim。
- 高風險鏈路需要至少一個 readback path，critical 鏈路需要 independent observation。

## 與特化案例的關係

Identity-coupled side effect validation 是本 pattern 的特化案例，而不是獨立的一級規則：

```text
Evidence Chain Validation
├── Identity-Coupled Flow
├── Payment Flow
├── Event-Driven Flow
├── Email Flow
└── Feature Flag Flow
```

## 相關知識

- [`state-visibility-gap.md`](state-visibility-gap.md)
- [`evidence-model.md`](evidence-model.md)
- [`evidence-depth.md`](evidence-depth.md)
- [`evidence-collapse-point.md`](evidence-collapse-point.md)
