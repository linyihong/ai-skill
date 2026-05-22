# ADR-006: Registry-First Workflow Activation

## Status

**Accepted**（2026-05-18）

## Framework Generation

- **世代分類**：Gen 3
- **當前世代文件**：[`architecture/ai-native-cognitive-execution-system.md`](../architecture/ai-native-cognitive-execution-system.md)
- **適用狀態**：registry-first activation 仍為 Gen 3 workflow 進入點預設；activation-table #27 + `route.workflow.*.activation_triggers` + `workflow/workflow-routing.md` 歧義表的三層分工不變。

## Context

Workflow 進入點曾出現三種維護方式並存：

1. `activation-table` 為每個 workflow 開專向列（#27 開發、#28 APK…），新增 workflow 需持續加列。
2. `routing-registry.yaml` 僅部分 route 有 `activation_triggers`，且曾反向引用 `activation_table_ref`。
3. Agent 可只載入單一 intelligence 規則（如 docs-first）而跳過 `routing-philosophy` → `execution-flow`。

使用者要求：察覺開發時**強制 routing discovery**；且 workflow 變多時**不要**無限擴充 activation 編號。

## Decision

採用 **registry-first** 三層分工：

| 層 | 職責 |
| --- | --- |
| **activation-table #27** | 通用「Workflow 編排閘門」+ §Workflow Discovery SOP（固定五步） |
| **`route.workflow.*.activation_triggers`** | 各 workflow 的觸發條件與 `required_dependencies`（機器可讀） |
| **`workflow/workflow-routing.md` §歧義** | 多 route 同時命中時的裁決（如 software-delivery vs apk-analysis） |

**禁止**：

- 為每個 workflow 新增 activation-table 專向列（#28、#29…）。
- 另立 `route.project.*` 作為 workflow 進入點；專案 overlay（如 `workflow-activation.yaml`）僅在選定 **software-delivery** 之後附加。

**強制**：

- 命中 #27 或任一 `route.workflow.*` 時，須先完成 [`governance/lifecycle/routing-philosophy.md`](../governance/lifecycle/routing-philosophy.md) Pipeline，再載入選定 route 的 `primary_source`。

## Consequences

### 正面

- 新增 workflow 只改 registry 一筆 + 必要時補歧義表。
- 觸發條件與 lazy rules 單一來源，利於 `activation-engine` 未來比對 registry。
- 與既有 `knowledge/indexes`、`dependency-reading` 對齊。

### 負面

- Agent 須學會「先掃 registry 再選路」，不能只記 #27=開發。
- 多 route 命中時仍需讀歧義表（人類維護成本轉移而非消失）。

## Alternatives Considered

- **每 workflow 一列 activation（#27/#28/…）**：已拒絕；不可擴展。
- **僅 workflow-routing.md 無 registry triggers**：已拒絕；無法機器化 dry-run。
- **`route.project.apk-analysis-sdk`**：已拒絕；與 software-delivery 平行入口混淆。

## Related

- [`workflow/workflow-routing.md`](../workflow/workflow-routing.md)
- [`runtime/router/activation-table.md`](../runtime/router/activation-table.md) §#27
- [`knowledge/runtime/routing-registry.yaml`](../knowledge/runtime/routing-registry.yaml)
- [`feedback/history/development-guidance/common/2026-05-18_160000-registry-first-workflow-activation.md`](../feedback/history/development-guidance/common/2026-05-18_160000-registry-first-workflow-activation.md)
- [ADR-003](ADR-003-three-layer-architecture.md) — workflow 層定位
