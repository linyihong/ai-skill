# ADR-012: `route_type` is an Activation Behavior Family (Interim Single-Axis Classifier)

## Status

**Accepted**

## Framework Generation

- **世代分類**：Gen 3 runtime hardening — Workflow Activation Engine 的 route classification 基礎
- **當前世代文件**：[`knowledge/runtime/routing-registry.yaml`](../knowledge/runtime/routing-registry.yaml) §`route_type_spec` / §`activation_mode_spec`；operational spec [`governance/workflow-activation-engine.md`](../governance/workflow-activation-engine.md)
- **適用狀態**：本 ADR 鎖定 `route_type` 的語意邊界（activation behavior family only），並明確標記 16-type enum 為過渡相容層，非長期 ontology。

## Date

2026-06-04

## Source Plan

[`plans/archived/2026-05-31-1900-workflow-activation-engine.md`](../plans/archived/2026-05-31-1900-workflow-activation-engine.md) §Phase 0.2（Route Type + Activation Mode 分類）+ §Open Questions Q10（route_type ontology collapse）

## Context

Workflow Activation Engine 需要每條 route 自宣告分類，讓 detector 由分類推導預設 `activation_mode`。Phase 0.2 落地時發現 plan 草擬的 12-value `route_type` enum 未涵蓋實際 57 routes（20 個 prefix）。補齊時面臨設計岔路：

```
route_type = 12  (plan draft)
  ↓ 補實際 prefix
route_type = 16  (本 ADR)
  ↓ 若無邊界
route_type = 25 → 30+  (每個 prefix 一個 type)
  ↓
route_type == prefix（失去抽象，複製 namespace）
```

這正是 plan Q10 指出的 **ontology collapse**：`route_type` 同時被當成 capability、activation behavior、knowledge domain 三個正交軸使用，導致同一主題（如 apk）散落在 `analysis.apk.*` / `workflow.apk-analysis` / `intelligence.apk-analysis.*` 三條 path。

### 觀察到的失效模式

- 分類欄位一旦缺乏「描述哪一個軸」的明確契約，每出現一個新 prefix 就有理由新增一個 type
- 沒有 explicit invariant → review 時無拒絕標準 → enum 漸進式爆炸
- 最終 `route_type` 退化成 prefix 的複製品，無法支撐 detector 的 mode 推導

## Decision

**`route_type` 的唯一語意是 Activation Behavior Family（一條 route 何時／如何變成 active）。它不得編碼 knowledge domain（apk / travel / architecture）或 capability。**

### Scope Boundary — IS for

- 描述 activation 行為族群：`always-on` / `auto-detect` / `on-demand` / `advisory` / `must-declare`（mixed）。
- 由 type 推導預設 `activation_mode`（單一語意明確的 type 給預設；mixed-layer type 標 `must-declare`，強制顯式宣告）。

### Scope Boundary — NOT for

- 不表達 knowledge domain（apk / travel / web / architecture）。
- 不表達 capability identity。
- 不作為 routing namespace 的鏡像（`route_type` ≠ prefix）。

### 16-type enum 是過渡相容層

本次落地的 16 個 type（plan 12 + 新增 `runtime_doc` / `memory` / `validation` / `anti_patterns`，並 fold `skill`→metadata、`evaluations`→validation、`traces`→metadata）是 **interim compatibility layer**，不是長期 ontology 模型。

長期方向（Q10 ontology-split future plan）將 `route_type` 拆為三個正交軸：

```
capability_domain:  runtime / workflow / intelligence / governance / memory / ...
activation_family:  always-on / auto-detect / on-demand / advisory
knowledge_domain:   apk / travel / architecture / software-delivery / documentation / ...
```

到那時 `route_type` 可能完全消失。任何未來重構不得被這 16 個 enum 綁死。

### Classification Heuristic（機械可判定）

採用 **session-entry-point test** 取代主觀 primary/secondary 判斷：

- route 可單獨作為使用者主任務入口 → primary candidate → 預設 `auto-detect`
- route 只有依附其他 route 才有意義 → secondary candidate → `advisory` 或 `on-demand`

此規則在 route 數量擴張到 100~300 時比 prose 描述更穩定。

### Pending TODO（type split triggers）

- `traces.*` 目前 fold 進 `metadata` 為暫時措施。當 traces / observability 類 route 超過 3~5 條，拆出獨立 `observability` type。
- 任何新出現的 mixed-layer prefix 應標 `must-declare`，不得給單一預設 mode。

## Consequences

### 正面

- `route_type` 有明確單軸契約，review 時對 enum 擴張有拒絕標準。
- detector 的 mode 推導有穩定來源；`must-declare` 防止 mixed-layer route 被靜默誤分類。
- 過渡層標記讓未來 ontology-split 重構不被既有 enum 綁死。

### 代價 / 已知限制

- 16-type enum 仍是單軸，短期內 `runtime_doc` 這類 type 實際混了 capability(runtime) + activation(reference) 兩個概念；這是過渡層的已知妥協，由 Q10 future plan 處理。

## Related

- [`knowledge/runtime/routing-registry.yaml`](../knowledge/runtime/routing-registry.yaml) §`route_type_spec`、§`activation_mode_spec`
- [`governance/workflow-activation-engine.md`](../governance/workflow-activation-engine.md)（operational spec，Phase 1 交付）
- [`constitution/ADR-006-registry-first-workflow-activation.md`](ADR-006-registry-first-workflow-activation.md)
- [`constitution/ADR-010-registry-upstream-classes-scope-freeze.md`](ADR-010-registry-upstream-classes-scope-freeze.md)（同類 scope-freeze 先例）
- plan §Open Question Q10（route ontology split — future plan）
