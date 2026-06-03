# ADR-010: Registry `upstream_classes` Scope Freeze (Promotion Traceability Only)

## Status

**Accepted**

## Framework Generation

- **世代分類**：Gen 3 Layer 2.5 (Meta Governance) 子系統治理約束（非世代升級）
- **當前世代文件**：[`enforcement/enforcement-registry.yaml`](../enforcement/enforcement-registry.yaml) §self_governance + `rule_classes[].upstream_classes`
- **適用狀態**：本 ADR 鎖定 `upstream_classes` 欄位的 scope；不擴充 registry schema，只定義邊界。

## Date

2026-06-01

## Source Plan

[`plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md`](../plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md) §Phase 3 Round-4/5 Review (T5 / U3)

## Context

Mechanical Enforcement Registry round-2 引入 `rule_classes[].upstream_classes: []` 處理 cross-class promotion governance（F20 `decision_promotion_pipeline.upstream_classes = [failure_learning_system]`）。Round-3 加 cycle detection；round-4 user 評審指出此欄位演化軌跡將觸發連鎖擴張：

```
upstream_classes  (round-2)
  ↓
+ cycle detect    (round-3)
  ↓
? downstream_classes
? promotion_role: {source, intermediate, sink}
? artifact_type
? dependency_type
  ↓
Governance DAG (Neo4j-tier complexity)
```

Round-4 T5 決定 **freeze 設計**；round-5 U3 要求 ADR 必須明確化 scope，避免半年後欄位被 silent 擴張回 DAG 方向。

### 觀察到的失效模式

- 治理 metadata 一旦長出第一個 graph 欄位 (`upstream_classes`)，後續每個 phase 都有理由補一個新欄位
- 沒有 explicit boundary 文件 → review 時無拒絕標準 → 漸進式擴張
- 最終 registry 變成 graph database mirror，違背「輕量 binding table」原始定位

## Decision

**`upstream_classes` 維持單一向上引用語意 + cycle detection，scope frozen 至 promotion traceability。**

### Scope Boundary — IS for

`upstream_classes` **僅用於**下列場景：

1. **Promotion traceability**：標示某 rule_class 接收哪些上游 class 的 promotion artifact（如 `decision_promotion_pipeline` 接收 `failure_learning_system` 升級上來的 pattern）
2. **Cycle prevention**：確保 promotion 鏈不會形成循環（lint `upstream_chain_resolution` 機械強制）
3. **Coverage report visualization**（future Phase 4）：在 governance dashboard 顯示 promotion chain

### Scope Boundary — IS NOT for

`upstream_classes` **明確禁止**用於下列場景：

1. **Execution ordering**：不是 DAG scheduler，executor 順序由 hook phase / dispatcher 定義
2. **Dependency injection**：不是 IoC container，executor 之間的程式依賴在 Go code 中表達
3. **Runtime orchestration**：不是 workflow engine，runtime 行為由 `runtime/phases/phase-machine.yaml` 定義
4. **DAG-based scheduling**：不是 task scheduler，scheduling 由 `runtime/scheduler/` 子系統管理
5. **Generic graph metadata**：不接受 `downstream_classes` / `promotion_role` / `artifact_type` / `dependency_type` 等補強欄位

### Supersession Clause

任何跨越上述 IS NOT boundary 的新欄位（包括但不限於 `downstream_classes`、`promotion_role`、`artifact_type`、`dependency_type`、`scheduling_hint`、`execution_order`），**必須先寫新 ADR 顯式 supersede 本 ADR**，不得 silent 加進 registry.yaml schema。

寫新 ADR 時必須回答：

1. 為什麼 `upstream_classes` + 既有子系統不足？
2. 為什麼新欄位不屬於既有子系統（runtime/scheduler / phase-machine / hook dispatcher / Go code）？
3. 半年觀察證據：是否有 ≥3 個實際 use case 需要此擴充？
4. 升級為 DAG 後的治理成本估計

### Revision Policy

```yaml
adr_revision_policy: supersede_required
```

不允許 in-place 修改本 ADR 來放寬 boundary；必須建立新 ADR 引用本 ADR 為 superseded，舊 ADR 保留為歷史紀錄。

## Consequences

### Positive

- Registry 維持輕量 binding table 定位，治理成本可控
- 未來欄位增加有 explicit gate（必須過新 ADR）→ silent 漂移不可能發生
- 既有 `upstream_classes` 用例（F20）不受影響

### Negative / 風險

- 如果未來真實需要 DAG-style 治理，需先過 ADR review → 增加 latency
- ADR review 過程可能拒絕合理的擴充需求（false negative）

### Mitigations

- 半年內若出現 ≥3 個明確 use case 仍被本 ADR 阻擋，視為「scope 過嚴」signal，啟動 ADR review
- Registry self-governance lint 未來可加 `adr_supersede_check`：偵測 schema diff 是否新增 graph 欄位但無對應 ADR

## Acceptance

- [x] 本 ADR 明寫 IS / IS NOT scope boundary
- [x] 本 ADR 明寫 supersession clause
- [x] `enforcement/enforcement-registry.yaml` companion `.md` 引用本 ADR 為 upstream_classes 設計依據（在 round-5 Step 1 schema patch 階段補上）

## Related

- [`enforcement/enforcement-registry.yaml`](../enforcement/enforcement-registry.yaml) §rule_classes[].upstream_classes
- [`plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md`](../plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md) §Phase 3 Round-4 T5 / Round-5 U3
- [`ADR-006-registry-first-workflow-activation.md`](ADR-006-registry-first-workflow-activation.md) — registry-first principle，本 ADR 延伸其治理邊界規範
