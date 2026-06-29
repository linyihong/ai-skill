# Implementation Execution Modes（Structure Preparation / Direct Change）

> **Cognitive Slice**：`sd-implementation`（retained in [`execution-flow.md`](../execution-flow.md) lifecycle；canonical execution-mode contract 在本檔。SDK 缺陷閉環與同工作階段閉環仍留在 `execution-flow.md` §3–§4）。

| slice 欄位 | 值 |
|---|---|
| `id` | `sd-implementation` |
| `purpose` | 在 implementation 階段選擇 **direct_change** 或 **preparatory_refactoring** execution mode；以 Change Intent Lock、Observable Equivalence Checkpoint、stop condition 安全改動既有 code |
| `type` | `execution` |
| `tags` | implementation, execution-mode, structure-preparation, change-intent-lock |
| `load_when` | 實際進行程式碼變更；`change_kind: feature` 且 `blocked_by_structure`；implementation plan 需宣告 `execution_mode` 或 step `intent` |
| `do_not_load_when` | 純 intake / 純 contract / 純 validation / evidence-only；`change_kind: replacement` 或 `migration`（用 `sd-intake` parity）；greenfield 以新檔為主且無 structure blocker |
| `owner_layer` | workflow |
| `canonical_source` | 本檔 + [`execution-flow.md`](../execution-flow.md) §3 SDK 缺陷閉環、§4 同工作階段閉環 |
| `dependencies` | `sd-intake`（change_kind / blocked_by_structure）、`sd-test-strategy`（intent → validation）、`sd-surgical-caveats`（feature intent scope）、`sd-validation`（完成證據） |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |

Plan source：[`plans/active/2026-06-29-1430-preparatory-refactoring-workflow.md`](../../../plans/active/2026-06-29-1430-preparatory-refactoring-workflow.md) Phase 1。

## 1. Execution modes（非新 lifecycle stage）

Implementation 仍是單一 lifecycle stage。差異在 **execution mode**，不是多一個 intake → preparatory → implementation 階段。

```yaml
implementation:
  execution_modes:
    - id: direct_change
      when: structure sufficient; feature can land locally without seam work

    - id: preparatory_refactoring
      when: change_kind is feature and blocked_by_structure is true
      sub_modes:
        - structure   # behavior_change.allowed: false
        - feature     # new acceptance / observable behavior
```

| mode | 何時 | 做什麼 |
|------|------|--------|
| `direct_change` | 結構已足夠；下一步 feature 是局部改動 | 直接實作 feature；仍遵守 [`surgical-changes.md`](../surgical-changes.md) |
| `preparatory_refactoring` | feature 被既有結構 block（untestable、god method、缺 seam） | 先 structure steps（行為不變）→ stop → feature steps |

**雙軸路由**（intake 細節 Phase 2 擴充；此處只定 implementation 消費面）：

```yaml
change_kind: feature          # 軸 1 — 本質仍是 feature，不是 refactor kind
blocked_by_structure: true    # 路由條件
execution_mode: preparatory_refactoring   # 軸 2
```

`preparatory_refactoring` 是 **execution mode**，不是 `change_kind`。勿與 `replacement` / `migration` parity 混用。

## 2. Change Intent Lock

每個 implementation plan step（`preparatory_refactoring` mode 或 plan 明確選 mode 時）**必須**宣告 `intent`：

| intent | `behavior_change.allowed` | validation |
|--------|---------------------------|------------|
| `structure` | `false` | observable equivalence / parity — 舊行為仍受保護 |
| `feature` | `true` | new acceptance / BDD / contract proof — 見 [`test-strategy.md`](../test-strategy.md) |

範例：

```yaml
steps:
  - id: prep-01
    intent: structure
    behavior_change:
      allowed: false
    action: extract apply_ranges()
    checkpoint:
      observable_equivalence:
        required: true

  - id: feat-01
    intent: feature
    behavior_change:
      allowed: true
    action: add multi-range support
    validation: new acceptance criteria
```

Commit / PR 分離是副產物，不是主要治理機制。

### Intent Transition Rule（intent state machine）

```yaml
transition:
  - from: structure
    to: feature
    require:
      - observable_equivalence_passed

  - from: feature
    to: structure
    require:
      - explicit_reopen_reason
```

- 預設禁止無理由 `feature → structure → feature` oscillation。
- `explicit_reopen_reason` 必須可 review（seam 不足、acceptance 無法 expressible）；不可是「再整理一下」。
- 未來 advisory validator 優先驗 **`illegal transition`**（guard 未滿足），而非「有沒有做 structure step」。

## 3. Observable Equivalence Checkpoint

`intent: structure` 的 canonical 要求：

> 允許內部結構改變，但 **observable contract 必須等價**。

不要求固定 `insert no-op` 手法。Project-local 範例：noop adapter、delegation、wrapper、branch disabled、dead-path extraction、fixture-backed parity assert、targeted mutation check（見 [`test-strategy.md`](../test-strategy.md) mutation 節）。

每個 structure step 標記：

```yaml
checkpoint:
  observable_equivalence:
    required: true
    evidence: <fixture | regression test | mutation check | manual parity note>
```

## 4. Stop condition

### 正向 exit（`exit_when`）

Structure mode 在任一成立時 exit，進入 feature intent 或結束 structure 批次：

```yaml
exit_when:
  - target_change_becomes_local
  - target_test_becomes_expressible
  - new_abstraction_created
```

原則：**下一步 feature 已容易加入就停；不追求漂亮。**

### Force exit（`force_exit_when`）

正向未觸發但 structure 無進展時，**必須**停止 structure mode：

```yaml
force_exit_when:
  - no_structure_progress_detected
  - abstraction_not_used_by_next_feature
```

`no_structure_progress_detected`（第一版 **不數值化** N-step）：以下皆不成立即視為無進展 —

- 下一步 feature 仍未變成 local change
- 目標 acceptance 仍不可 expressible 為 failing test
- 新 seam 未被下一步 feature step 引用

`abstraction_not_used_by_next_feature`：anti-architecture-as-output — 架構整理了但 feature 踩不上去。

觸發後不得再加 structure step；必須 (a) 進 feature mode、(b) `direct_change` 縮小 scope、或 (c) 回 intake 重判雙軸 / 拆 plan。

### 禁止（`avoid`）

即使 `intent: structure` 也不允許：

```yaml
avoid:
  - broad_cleanup
  - style_only
  - debt_harvesting
  - opportunistic_refactor
```

## 5. Compatibility（opt-in enrichment）

```yaml
compatibility:
  existing_plans:
    default_execution_mode: direct_change
  existing_implementation_plans:
    missing_intent_on_steps: allowed
    missing_execution_mode: direct_change
```

未宣告 `execution_mode` 的 plan 視為 `direct_change`。僅當走 `preparatory_refactoring` 或 plan 明確選 mode 時才強制 Change Intent Lock。

## 6. Structure mode `recommended_recipe`（advisory）

Fowler「預備式重構」八步作 **參考 checklist**，非 mandatory order、非新 lifecycle step：

1. 寫最終 acceptance / E2E 目標（可先 skip / pending，project-local）
2. Extract / 結構隔離 — 把 feature 將落腳的邏輯抽出
3. Insert seam — delegation / adapter / wrapper（project-local）
4. Observable equivalence checkpoint — 證明輸出未變
5. 漸進拆解 — start point / segment / end point
6. 滿足 `exit_when` 或觸發 `force_exit_when` 處置
7. Transition `structure → feature`（需 `observable_equivalence_passed`）
8. Feature intent steps — 實作新行為

可跳步；不可跳 checkpoint 與 transition guard。

## 7. Surgical-changes reconciliation

[`surgical-changes.md`](../surgical-changes.md) 適用於 `direct_change` 與 `intent: feature` steps。

**例外**：plan 內 `intent: structure` steps，具 Observable Equivalence Checkpoint、符合 stop condition、且不在 `avoid` 清單 — 不算 scope creep。

## 8. Failure modes

| failure | 徵象 | recovery |
|---------|------|----------|
| infinite refactor | structure steps 無 `exit_when` / `force_exit_when` 收斂 | 觸發 force exit；縮 scope 或回 intake |
| intent oscillation | feature → structure 無 `explicit_reopen_reason` | 停止；補 reason 或改 plan |
| intent 混用 | 同 step 混 structure 與 feature 語意 | 拆 step；每步單一 intent |
| replacement parity 誤用 | preparatory structure work 走 parity inventory | 改標 `change_kind: feature` + `execution_mode: preparatory_refactoring` |
| skip checkpoint | structure step 無 equivalence 證據 | 補 checkpoint 或縮小 correctness claim |
| illegal transition | structure → feature 無 equivalence pass | 補證據或維持 structure |
