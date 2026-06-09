---
id: 2026-06-09-1040-experience-validation-pipeline-evolution
plan_kind: main
status: in-progress
owner: linyihong
created: 2026-06-09
priority: P2
required_for_completion: false
---

# Experience Validation Pipeline Evolution

**Status**: `in-progress`
Owner: framework maintainer (linyihong)
**建立日期**：2026-06-09
**Priority**：**P2**（taxonomy evolution / watch-list；不阻塞既有 responsive gate；Phase 1 metadata work 可獨立升 P1）

## Why this plan exists

`responsive_validation_complete` 已把 UI responsive 檢查從「多裝置測試」提升成 workflow artifact gate：

```text
UI Contract
  -> Render Context
  -> Browser Evidence Collection
  -> Validation Matrix
  -> Evidence-backed Completion
```

這次使用者 review 指出兩個後續方向：

1. `Responsive` 目前作為 governance domain 合理，但長期可能只是 cross-cutting `render_context` dimension。
2. Browser Review / Validation Matrix 會自然往更廣的 Experience Validation Pipeline 演化：State Coverage + Context Coverage + Evidence Coverage。
3. `experience_context` 這個名稱太泛，未來容易變成垃圾桶欄位；更可控的方向是 typed Context Taxonomy。

本 plan 的目的不是立即擴張 taxonomy，而是把後續設計壓力、open questions 與安全演化順序記錄下來，避免剛落地的 responsive gate 直接膨脹成 device matrix 或 Cartesian product。

Although this plan originates from UI responsive validation, the emerging Coverage Model may eventually become a shared Validation Reasoning concept rather than a UI-specific taxonomy. API validation, runtime validation, workflow validation, and UI validation can all face the same question: which states, contexts, and evidence classes are covered well enough to support the completion claim?

## Decision Rationale

### Problem & Why Now

目前 `sd-ui-governance` 已有：

```yaml
domain:
  - Responsive
render_context:
  - desktop
  - mobile
```

這解決了「desktop 正常、mobile 壞掉」不應被單純視為一般 contract failure 的問題。但若後續加入 dark mode、offline、high zoom、screen reader、touch only、keyboard only，taxonomy 可能膨脹成多個類似 domain：

```yaml
domain:
  - responsive
  - accessibility
  - dark_mode
  - offline
```

這會模糊 domain 與 context 的責任邊界，也可能讓 validation matrix 變成所有 state × context × evidence 的笛卡兒積。

更好的 long-term shape 可能是 typed context taxonomy：

```yaml
context:
  render:
    - desktop
    - mobile
    - tablet
  interaction:
    - mouse
    - touch
    - keyboard_only
  accessibility:
    - screen_reader
    - high_zoom
  environment:
    - offline
    - slow_network
  appearance:
    - light
    - dark
  locale:
    - en
    - ja
    - ar
```

這保留 context 的結構，而不是把 mobile、dark mode、offline、keyboard only、localization 全丟進一個 generic `experience_context`。

### Decision

短期保留 `Responsive` as governance domain，因為它是目前最容易理解的 failure class，也能讓 `responsive_validation_complete` gate 立即產生治理效果。

同時記錄未來演化方向：

- Browser evidence artifacts 需要 capture metadata，避免 screenshot、DOM snapshot、interaction trace 或 accessibility scan 失去證據語境。這是最接近現有 workflow 壓力的第一優先。
- Metadata ownership 應優先評估 Capture Envelope，而不是讓每個 artifact 重複貼 viewport metadata。
- Capture Envelope is evidence-scoped. It records only metadata required to interpret collected evidence; it is not intended to become a full runtime environment descriptor.
- Validation coverage model 應朝 State Coverage + Context Coverage + Evidence Coverage 觀察，不急著新增 executable gate。
- `Responsive` 可能降為 cross-cutting `context.render` dimension。
- `experience_context` 不採用為 active 名稱；未來若需要，優先設計 typed `context` taxonomy。
- Validation Matrix 未來若擴張，應先引入 `required_contexts` / `optional_contexts` / `high_risk_contexts`，不要直接展開 Cartesian product。

### Alternatives Considered

- A. 立即把 `Responsive` 改成 pure `render_context` dimension：reject for now，因為剛落地的 gate 需要一個清楚可理解的 failure class，直接抽象化會降低可操作性。
- B. 立即新增 generic `experience_contexts` taxonomy：reject，因為名稱過泛，會把 render、interaction、accessibility、environment、appearance、locale 混成垃圾桶欄位。
- C. 保留 Responsive domain，但把 domain-vs-context 問題寫入 plan：accept，保留當前治理效果，同時避免未來把第一版 taxonomy 誤認成最終模型。
- D. 先推 Browser Evidence Metadata：accept as Phase 1 priority，因為 screenshot / DOM / interaction evidence 已經存在，缺 metadata 會直接削弱 reviewability 和 trust。
- E. 每個 artifact 自帶 metadata：tentatively reject，因為 screenshot、DOM snapshot、accessibility scan、interaction trace 很可能共享同一次 capture context；重複欄位會造成 drift。Phase 1 應先評估 Capture Envelope。

### Why Not an ADR Yet

這仍是 taxonomy evolution watch-list，不是 accepted architecture decision。尚未滿足 ADR promotion 條件：

- typed Context Taxonomy 尚未實際落地。
- Validation coverage model 尚未有跨 scenario / cross-domain 使用證據。
- `Responsive` 是否保留 domain 還有未解 open question。
- Evidence metadata ownership 尚未決定：Capture Envelope vs per-artifact metadata。

### ADR Promotion Criteria（completed 時驗證）

- [ ] 至少 3 類 typed context（例如 render、appearance、environment、accessibility 或 interaction）產生真實 scenario 壓力。
- [ ] 有至少 2 個例子顯示 domain × context 比獨立 domain 更清楚。
- [ ] Validation Matrix 的 context selection 不再能靠 prose gate 管理，需要 formal artifact shape。
- [ ] `viewport_metadata` / capture metadata 被至少一個 review artifact 實際消費。
- [ ] Capture Envelope vs per-artifact metadata ownership 有清楚決策。
- [ ] 沒有更輕量 promotion target（workflow note / validation scenario / template field）即可解決。

### Consequences（預期）

#### 正面

- 避免 Responsive gate 退化成 device SKU matrix。
- 保留 `responsive_validation_complete` 的高價值 gate。
- 為未來 dark mode、offline、keyboard-only、screen reader、localization 等 context 留下 typed taxonomy 演化路線。

#### 負面

- 第一版仍同時有 `domain: Responsive` 與 `render_context`，taxonomy 看起來可能有些重疊。
- 在沒有更多 scenarios 前，coverage model 只能是 watch-list，不可過早機械化。

#### 風險

- 若過早新增 generic `experience_contexts`，workflow 可能變成過重且無結構的 matrix checklist。
- 若 Capture Envelope 未限制為 evidence-scoped，可能膨脹成第二個 runtime context schema（browser version、OS、locale、timezone、network、CPU、memory 等都被塞進 capture metadata）。
- 若不記錄 open question，未來可能把 `Responsive` 當永久 domain，導致 dark mode / offline 等也被錯誤 domain 化。

## Priority Order

| Priority | Focus | Rationale |
|---|---|---|
| P1 | Browser Evidence Metadata | 已有 screenshots / DOM snapshots / interaction traces；缺 viewport/capture metadata 會讓 evidence 無法 review。 |
| P2 | Coverage Model | State Coverage + Context Coverage + Evidence Coverage 已有真實 workflow 壓力，但仍需先保持 doc/watch-list。 |
| P3 | Responsive domain downgrade decision | 需要更多 usage evidence 判斷 `Responsive` 是否應降為 `context.render`。 |
| P4 | Typed Context Taxonomy | render / interaction / accessibility / environment / appearance / locale 需要更多 scenarios，不能先做成 generic `experience_context`。 |

P1 metadata working split:

```yaml
capture:
  metadata:
    required:
      - viewport_width
      - viewport_height
      - orientation
      - render_context
    optional:
      - dpr
      - user_agent
      - emulation_profile
      - safe_area
  artifacts:
    - screenshot
    - dom_snapshot
    - interaction_trace
    - accessibility_scan
```

Reasoning: without `viewport_width` / `viewport_height` / `orientation` / `render_context`, reviewers cannot tell whether evidence actually represents mobile, desktop, landscape, or another governed context. `dpr`, `user_agent`, `emulation_profile`, and `safe_area` improve trust for specific issues, but many reviews can still proceed without them.

Capture Envelope boundary:

- Evidence-scoped only: include metadata required to interpret the collected artifacts.
- Not a runtime environment descriptor: do not grow it into browser / OS / locale / timezone / network / CPU / memory inventory unless a specific evidence claim requires that field.
- Shared ownership: metadata that applies to the whole browser capture belongs in `capture.metadata`; artifact-specific metadata should stay under the artifact only when it differs from the shared capture context.

Glossary Impact: yes — candidate terms: `context_taxonomy`, `validation_matrix`, `coverage_dimensions`, `viewport_metadata`; do not register until Phase 1 proves they are stable framework vocabulary rather than local plan terms.

Watch-Out List citation: [`architecture/ai-native-cognitive-ecosystem-system.md`](../../architecture/ai-native-cognitive-ecosystem-system.md) §Watch-Out List — avoid over-engineering new runtime abstractions before evidence shows the taxonomy is reusable.

## Runtime Execution Path

This plan is **doc-first / watch-list only**.

No new runtime projection, route, commit-msg validator, or enforcement registry entry is planned in Phase 0. Existing runtime surfaces remain:

| Surface | Current consumer | Planned change |
|---|---|---|
| `workflow/software-delivery/ui-governance.md` | software-delivery workflow loading surface | Add Responsive-domain open question and typed context taxonomy watch-list only |
| `workflow/software-delivery/validation.md` | validation slice | Add browser evidence metadata, Capture Envelope decision, and coverage model note only |
| `workflow/software-delivery/execution-flow.yaml` | executable workflow contract | No semantic gate change unless Phase 1 proves wording clarification is needed |
| `validation/scenarios/software-delivery/ui-governance-responsive-*.yaml` | runtime scenario inventory | No new scenario unless a new active rule is introduced |

Runtime validation remains `ai-skill runtime compile`, `ai-skill runtime refresh`, and `ai-skill runtime validate` if any runtime-indexed source changes.

## Open Questions

- [ ] Should `Responsive` remain a governance domain, or become a cross-cutting `context.render` dimension composed with domains such as Accessibility, Behavior, Design System, and Contract?
- [ ] Should the future context model be a typed `context` taxonomy with render / interaction / accessibility / environment / appearance / locale groups?
- [x] What is the minimum capture metadata required for browser evidence to remain reviewable: viewport width/height, user agent, emulation/device profile, orientation, DPR, safe-area, render_context, or all of the above? — resolved provisionally in Phase 1: required = viewport width/height, orientation, render context; optional = DPR, user agent, emulation profile, safe area.
- [x] Who owns evidence metadata: a shared Capture Envelope containing metadata plus artifacts, or per-artifact metadata under screenshot / DOM snapshot / accessibility scan / interaction trace? — resolved provisionally in Phase 1: shared Capture Envelope owns common capture metadata; artifact-specific metadata remains local only when it differs.
- [ ] Should validation coverage be modeled explicitly as State Coverage + Context Coverage + Evidence Coverage?
- [ ] Does Validation Coverage Model belong to UI workflow, or should it graduate to shared Validation Reasoning once API / runtime / workflow validation show the same coverage pressure?
- [ ] How should matrix explosion be controlled: `required_contexts`, `optional_contexts`, `high_risk_contexts`, risk-triggered expansion, or another shape?

## Phase 0 — Plan and Current-State Alignment

### Phase 0.0 — Open Questions 核對（公版，必填）

逐條核對本 plan §Open Questions，標記處置並回寫：

- [x] 已讀本 plan §Open Questions 全部條目
- [x] 對每條標記 `resolved`（附 Phase 0 證據）/ `still-open` / `deferred`（附原因）
- [x] `resolved` 的條目已同步勾選 / 附註於 §Open Questions
- [x] 若盤點新發現問題，已加入 §Open Questions

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| Responsive domain vs context.render dimension | still-open | Needs post-gate usage evidence |
| typed context taxonomy graduation | still-open | Explicitly deferred until Phase 2+ examples exist |
| capture metadata minimum | resolved | Phase 1 defines required vs optional fields in `workflow/software-delivery/validation.md` |
| evidence metadata ownership | resolved | Phase 1 chooses evidence-scoped Capture Envelope for shared metadata |
| coverage dimensions model | still-open | Needs validation scenario pressure |
| coverage model owner | still-open | Could remain UI-local or graduate to shared validation-reasoning |
| matrix explosion control | still-open | Needs at least one expanded matrix example |

- [x] Read current `sd-ui-governance`, `sd-ui-contracts`, `sd-validation`, and `execution-flow.yaml`.
  - Evidence: prior responsive gate landing read/updated all four surfaces; this execution reread `sd-validation` before editing.
- [x] Confirm existing `responsive_validation_complete` gate remains correct and should not be weakened.
  - Evidence: Phase 1 adds evidence metadata only; no semantic change to `workflow/software-delivery/execution-flow.yaml`.
- [x] Confirm no new runtime surface is needed for Phase 0.
  - Evidence: Phase 1 is doc-only on an existing workflow slice; no new runtime projection, route, validator, or scenario.
- [x] Update this plan if current sources already answer one open question.
  - Evidence: capture metadata minimum and metadata ownership are marked resolved provisionally after Phase 1.

## Phase 1 — Browser Evidence Metadata

- [x] Add browser evidence capture metadata note to `workflow/software-delivery/validation.md`.
- [x] Define provisional required metadata fields: `viewport_width`, `viewport_height`, `orientation`, `render_context`.
- [x] Define provisional optional metadata fields: `dpr`, `user_agent`, `emulation_profile`, `safe_area`.
- [x] Clarify evidence chain: Evidence → Metadata → Reviewability → Trust.
- [x] Decide whether metadata belongs in a shared Capture Envelope or under each artifact.
- [x] If Capture Envelope wins, define the provisional shape as `capture.metadata` + `capture.artifacts`.

Phase 1 evidence:

- `workflow/software-delivery/validation.md` §Browser Evidence Collection now defines evidence-scoped Capture Envelope.
- Shared metadata belongs under `capture.metadata`; artifact-specific metadata is local only when it differs.
- Capture Envelope is explicitly not a full runtime environment descriptor.

## Phase 2 — Coverage Model Watch-List

- [ ] Add a short coverage model note: UI validation is trending toward State Coverage + Context Coverage + Evidence Coverage.
- [ ] Record that Coverage Model likely belongs to shared Validation Reasoning if API / runtime / workflow validation show the same state-context-evidence shape.
- [ ] Inventory at least one non-UI example before proposing any shared model:
  - API: success / error / timeout states; authenticated / anonymous contexts; response / logs / traces evidence.
  - Runtime: startup / steady-state / shutdown states; single-node / multi-node contexts; metrics / logs / traces evidence.
- [ ] Document candidate coverage dimensions without adding executable gates:
  - state: loading / empty / success / error
  - context: desktop / mobile / dark_mode / offline / keyboard_only / screen_reader
  - evidence: screenshot / DOM snapshot / accessibility scan / interaction trace
- [ ] Add guidance that validation matrix growth should use required / optional / high-risk contexts instead of a Cartesian product.

## Phase 3 — Responsive Domain Decision

- [ ] Add a compact open question to `workflow/software-delivery/ui-governance.md`: Responsive domain vs cross-cutting context.render.
- [ ] Inventory responsive scenarios and determine whether each is better modeled as `domain: Responsive` or another domain plus `context.render`.
- [ ] When most Responsive failures are better explained as another governance domain evaluated under `context.render`, draft a downgrade path for `Responsive`.
- [ ] If `Responsive` remains clearer for routing, keep it as a workflow-local domain and record why.

## Phase 4 — Typed Context Taxonomy Graduation Decision

- [ ] Inventory future scenarios involving dark mode, high zoom, keyboard-only, screen reader, low bandwidth, offline, touch-only, or localization.
- [ ] Decide whether each scenario is best represented as domain, typed context, coverage dimension, or evidence coverage.
- [ ] Only draft an active typed Context Taxonomy update if promotion evidence meets all gates:
  - at least 3 context families are represented
  - each represented family has multiple scenarios
  - at least 2 workflow domains consume the taxonomy
- [ ] Candidate context families must demonstrate cross-domain use before promotion, for example:
  - render used by responsive and accessibility
  - appearance used by accessibility and design-system
  - environment used by behavior and runtime validation
- [ ] If promotion evidence passes, draft the active taxonomy with typed families:
  - render
  - interaction
  - accessibility
  - environment
  - appearance
  - locale
- [ ] If examples remain one-off, keep the current Responsive domain model and close the plan with no taxonomy expansion.

## Completion Criteria

- [ ] Open questions are either resolved or explicitly deferred with evidence.
- [ ] Any docs changed by Phase 1 are validated with lints and `git diff --check`.
- [ ] If runtime-indexed docs or scenarios change, runtime compile / refresh / validate pass.
- [ ] `plans/README.md` status row is updated.
- [ ] Commit / push completed if implementation phases are executed.

## Stakeholder 同意項目

- Keep `responsive_validation_complete` as the valuable first gate.
- Do not add device-SKU testing taxonomy.
- Do not activate generic `experience_contexts`; if scenario pressure proves the need, prefer typed Context Taxonomy.
- Treat the current plan as watch-list / evolution planning, not an immediate runtime expansion.

## 與其他 plans 的關係

- Follows the landed responsive render-context governance commit `bba1020`.
- Builds on [`archived/2026-06-08-1408-ui-governance-workflow.md`](../archived/2026-06-08-1408-ui-governance-workflow.md).
- Builds on [`archived/2026-06-08-1544-evidence-acquisition-layer.md`](../archived/2026-06-08-1544-evidence-acquisition-layer.md).
- Related to [`active/2026-06-08-2100-governance-pattern-library-extraction.md`](2026-06-08-2100-governance-pattern-library-extraction.md) as another observation-stage taxonomy evolution, but it is not a governance-pattern-library sample unless future implementation adds rule / registry / projection / executor / validation shape.
