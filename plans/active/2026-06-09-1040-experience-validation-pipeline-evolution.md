---
id: 2026-06-09-1040-experience-validation-pipeline-evolution
plan_kind: main
status: completed
owner: linyihong
created: 2026-06-09
priority: P2
required_for_completion: false
---

# Experience Validation Pipeline Evolution

**Status**: `completed`（2026-06-15；Phase 4 Graduation Adjudication = Branch A：typed context taxonomy `rejected_for_now`，保留 watch-list reopen triggers）
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

Follow-up execution in [`2026-06-10-0908-user-journey-validation-integration.md`](2026-06-10-0908-user-journey-validation-integration.md) adds evidence-backed Outcome Validation pressure: BDD owns Journey Specification, software-delivery validation owns Journey Execution, and the first landing models Journey as `validation_scope` rather than `validation_domain` or typed Context Taxonomy.

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

### ADR Promotion Criteria（completed 時驗證 — evaluated 2026-06-15）

評估結果：**6 條中只有 2 條成立，promotion 不通過 → Branch A**（見 §Phase 4 Graduation Adjudication）。

- [x] 至少 3 類 typed context 產生真實 scenario 壓力 — **NOT met**：0/6 family 通過 3 題 gate。
- [x] 有至少 2 個例子顯示 domain × context 比獨立 domain 更清楚 — **NOT met**：render/interaction 各僅單一 domain 消費。
- [x] Validation Matrix 的 context selection 不再能靠 prose gate 管理，需要 formal artifact shape — **NOT met**：required/optional/high-risk prose gate 仍足夠（Phase 2）。
- [x] `viewport_metadata` / capture metadata 被至少一個 review artifact 實際消費 — **MET**：H5 fixed-bottom review 消費了 viewport / visual-viewport / frame metadata（Phase 1）。
- [x] Capture Envelope vs per-artifact metadata ownership 有清楚決策 — **MET**：shared Capture Envelope 擁有共通 metadata、artifact-specific 僅在差異時 local（Phase 1）。
- [x] 沒有更輕量 promotion target（workflow note / validation scenario / template field）即可解決 — **NOT met**：更輕的 workflow note + coverage model + evidence selection 已能解，**不需** ADR-level taxonomy。

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
| P2 | Coverage Model | State Coverage + Context Coverage + Evidence Coverage 已有真實 workflow 壓力；Journey follow-up adds Outcome Validation pressure while keeping Journey as `validation_scope` for the first landing. |
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

- [x] Should `Responsive` remain a governance domain, or become a cross-cutting `context.render` dimension composed with domains such as Accessibility, Behavior, Design System, and Contract? — **resolved (Phase 3 + Phase 4)**: keep `Responsive` as workflow-local domain; Phase 4 graduation inventory shows `render` does not meet the cross-domain typed-context gate. Reopen per Phase 4 Branch A triggers.
- [x] Should the future context model be a typed `context` taxonomy with render / interaction / accessibility / environment / appearance / locale groups? — **resolved (Phase 4 Branch A)**: `rejected_for_now` — 0/6 families pass the 3-question gate; pressure was evidence-depth, not missing context families.
- [x] What is the minimum capture metadata required for browser evidence to remain reviewable: viewport width/height, user agent, emulation/device profile, orientation, DPR, safe-area, render_context, or all of the above? — resolved (Phase 1 baseline): required = viewport width/height, orientation, render context; optional = DPR, user agent, emulation profile, safe area. Reopen only if future evidence needs safe-area, DPR, browser-specific, or similar metadata as required.
- [x] Who owns evidence metadata: a shared Capture Envelope containing metadata plus artifacts, or per-artifact metadata under screenshot / DOM snapshot / accessibility scan / interaction trace? — resolved (Phase 1 baseline): shared Capture Envelope owns common capture metadata; artifact-specific metadata remains local only when it differs. Reopen only if future artifact examples show shared capture metadata causes ambiguity.
- [x] Should validation coverage be modeled explicitly as State Coverage + Context Coverage + Evidence Coverage? — **resolved (Phase 2)**: yes, as a non-executable watch-list note in `validation.md`; no executable gate added. Promotion to shared Validation Reasoning deferred (next OQ).
- [x] Does Validation Coverage Model belong to UI workflow, or should it graduate to shared Validation Reasoning once API / runtime / workflow validation show the same coverage pressure? — **deferred with evidence (Phase 2/4)**: graduate to shared Validation Reasoning only after ≥3 non-UI domains reuse `state/context/evidence coverage + validation_scope`; the membership_purchase journey added outcome-depth pressure in that direction but is one domain.
- [x] How should matrix explosion be controlled: `required_contexts`, `optional_contexts`, `high_risk_contexts`, risk-triggered expansion, or another shape? — **resolved (Phase 2)**: required / optional / high-risk contexts + risk-triggered expansion, not a Cartesian product.
- [x] Should Capture Envelope remain browser-specific, or evolve into a generic Evidence Envelope shared across UI, API, runtime, and workflow validation? — **resolved (Phase 4 spike)**: keep browser-specific; generic Evidence Envelope deferred until ≥3 non-browser examples converge on a minimal schema (promote under validation reasoning, not UI workflow).

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
| capture metadata minimum | resolved (Phase 1 baseline) | Phase 1 defines required vs optional fields in `workflow/software-delivery/validation.md`; reopen only for concrete safe-area / DPR / browser-specific evidence pressure |
| evidence metadata ownership | resolved (Phase 1 baseline) | Phase 1 chooses evidence-scoped Capture Envelope for shared metadata; reopen only if artifact examples show ambiguity |
| coverage dimensions model | still-open | Needs validation scenario pressure |
| coverage model owner | still-open | Could remain UI-local or graduate to shared validation-reasoning |
| matrix explosion control | still-open | Needs at least one expanded matrix example |
| Capture Envelope vs Evidence Envelope | deferred after initial spike | Non-browser examples exist, but common shape is not stable enough; keep browser Capture Envelope local for now |

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
- Metadata minimum and ownership are considered Phase 1 baseline decisions, not open design debates, unless future evidence introduces counterexamples.
- Project feedback from an H5 fixed-bottom tab issue: required capture metadata stayed stable (`viewport_width`, `viewport_height`, `orientation`, `render_context`), but the claim needed optional visual viewport / bottom-frame evidence because `document` dimensions alone did not explain a user-visible iPhone browser frame problem. This supports keeping Capture Envelope evidence-scoped and adding claim-specific optional metadata instead of expanding a typed Context Taxonomy.
- Project feedback from an H5 player-to-drama navigation issue: static route/component markers proved the contract existed, but did not prove browser history behavior. The useful evidence was a real browser interaction trace: open player, open episode sheet, click drama title, reach `/drama/:id`, click detail back, and verify return to the original `/player/:episodeId`. This supports treating route/history trace as a browser evidence artifact, not as screenshot/RWD-only evidence.

## Phase 2 — Coverage Model Watch-List

- [x] Add a short coverage model note: UI validation is trending toward State Coverage + Context Coverage + Evidence Coverage.
- [x] Record that Coverage Model likely belongs to shared Validation Reasoning if API / runtime / workflow validation show the same state-context-evidence shape.
- [x] Inventory at least one non-UI example before proposing any shared model:
  - API: success / error / timeout states; authenticated / anonymous contexts; response / logs / traces evidence.
  - Runtime: startup / steady-state / shutdown states; single-node / multi-node contexts; metrics / logs / traces evidence.
- [x] Document candidate coverage dimensions without adding executable gates:
  - state: loading / empty / success / error
  - context: desktop / mobile / dark_mode / offline / keyboard_only / screen_reader
  - evidence: screenshot / DOM snapshot / accessibility scan / interaction trace
- [x] Add guidance that validation matrix growth should use required / optional / high-risk contexts instead of a Cartesian product.

Phase 2 evidence:

- `workflow/software-delivery/validation.md` §Validation Coverage Model Watch-List now documents state / context / evidence coverage.
- The note explicitly remains non-executable and points shared promotion toward validation reasoning if API / runtime / workflow examples prove reuse.
- The Vidoe-Test player-to-drama navigation pilot adds interaction-context pressure: the governed state is route/history/back-stack behavior, the high-risk context is a specific entrypoint (`player episode sheet -> drama detail`), and the evidence is a browser route/history interaction trace. This strengthens the Coverage Model watch-list without requiring a new runtime gate.

## Phase 3 — Responsive Domain Decision

- [x] Add a compact open question to `workflow/software-delivery/ui-governance.md`: Responsive domain vs cross-cutting context.render.
- [x] Inventory responsive scenarios and determine whether each is better modeled as `domain: Responsive` or another domain plus `context.render`.
- [x] When most Responsive failures are better explained as another governance domain evaluated under `context.render`, draft a downgrade path for `Responsive`.
- [x] If `Responsive` remains clearer for routing, keep it as a workflow-local domain and record why.

Phase 3 evidence:

- `workflow/software-delivery/ui-governance.md` §Render Contexts now includes Responsive domain downgrade watch.
- Current decision: keep `Responsive` as workflow-local domain while it helps route viewport-specific failures clearly.
- Downgrade condition: when most responsive failures are better explained as another governance domain evaluated under `context.render`, draft the downgrade path.
- Existing scenario pressure is mixed: mobile-only layout failures justify `Responsive`; mobile accessibility failure already suggests domain × context may become clearer later.

## Phase 4 — Typed Context Taxonomy Graduation Decision

- [x] Inventory future scenarios involving dark mode, high zoom, keyboard-only, screen reader, low bandwidth, offline, touch-only, route/history back-stack behavior, or localization. — done; see §Phase 4 Graduation Adjudication (these are incidental or single-scenario in Vidoe-Test, no cross-domain family).
- [x] Inventory non-browser evidence envelopes, such as API captures with endpoint/auth context, runtime captures with node/cluster context, or workflow captures with route/execution context.
  - Evidence: `intelligence/engineering/execution/validation-reasoning/evidence-model.md` defines evidence type / confidence / scope / proves for non-browser API/log/database/user-observable evidence.
  - Evidence: `intelligence/engineering/execution/validation-reasoning/evidence-chain-validation.md` defines claim → chain → segment evidence → gap/depth reasoning.
  - Evidence: `validation/scenarios/software-delivery/ui-governance-static-runtime-acquisition.yaml` and `ui-governance-contract-readback-acquisition.yaml` show static/runtime/contract readback acquisition shape.
  - Evidence: `validation/scenarios/runtime/workflow-detector-deterministic-match-v1.yaml` and `workflow-detector-conflict-resolution-v1.yaml` show runtime scenario transcript / constraints / trace / verification shape.
- [x] Decide whether each scenario is best represented as domain, typed context, coverage dimension, or evidence coverage. — see §Phase 4 Graduation Adjudication (render→`Responsive` domain; interaction→evidence-coverage; appearance→multi-theme contract; a11y/offline→evidence artifact / ungoverned).
- [x] Decide whether browser-specific Capture Envelope should remain local or graduate into generic Evidence Envelope after non-browser examples exist.
  - Initial decision: keep browser-specific Capture Envelope local for now.
  - Reason: browser Capture Envelope is artifact-bundle oriented; validation reasoning examples are claim/chain/scope/confidence oriented; runtime scenarios are replay/trace/verification oriented.
  - Promotion risk: a generic envelope now would flatten artifact vs evaluation vs finding, or become a second runtime context schema.
- [x] Only draft an active typed Context Taxonomy update if promotion evidence meets all gates: — gate NOT met (0/6 families pass); no draft. See §Phase 4 Graduation Adjudication.
  - at least 3 context families are represented
  - each represented family has multiple scenarios
  - at least 2 workflow domains consume the taxonomy
- [x] Candidate context families must demonstrate cross-domain use before promotion, for example: — not demonstrated (render→UI governance only; interaction→UI validation only). See §Phase 4 Graduation Adjudication.
  - render used by responsive and accessibility
  - appearance used by accessibility and design-system
  - environment used by behavior and runtime validation
  - interaction used by UI behavior navigation and accessibility / input-modality validation
- [x] If promotion evidence passes, draft the active taxonomy with typed families: — N/A, promotion did not pass (Branch A). See §Phase 4 Graduation Adjudication.
  - render
  - interaction
  - accessibility
  - environment
  - appearance
  - locale
- [x] If examples remain one-off, keep the current Responsive domain model and close the plan with no taxonomy expansion. — **adopted (Branch A)**. See §Phase 4 Graduation Adjudication.

### Phase 4 Initial Evidence Envelope Spike

Initial result: **do not promote generic Evidence Envelope yet**.

Observed partial common shape:

```text
claim / task context
  -> acquisition or source
  -> artifacts, observations, traces, or chain segments
  -> evaluation / verification
  -> scope / confidence / result
```

Why this is not enough for promotion:

- Browser Capture Envelope is centered on a single observation with shared capture metadata and multiple artifacts.
- Fixed-bottom UI evidence showed a useful claim-specific extension: visual viewport and bottom-frame measurements can belong to the browser capture artifact/optional metadata without changing the required metadata set.
- Player-to-drama navigation evidence showed another useful claim-specific extension: route/history trace can belong to the browser capture artifacts without turning Capture Envelope into a generic runtime context or environment schema.
- API / runtime / workflow examples need claim, propagation chain, trace, verification, scope, and confidence fields more than viewport-like metadata.
- A generic Evidence Envelope would risk becoming too broad before non-browser examples prove a stable minimum schema.

Future promotion should happen in shared validation reasoning, not UI workflow, and only after at least 3 non-browser examples converge on a minimal envelope shape.

### Phase 4 Observation — Journey Validation Pilot

Vidoe-Test `membership_purchase` pilot added useful pressure, but it does **not** justify promoting Typed Context Taxonomy or a generic Evidence Envelope yet.

Observed result:

- Journey / Outcome Validation should be treated as `validation_scope`, not as a new context family and not as a UI governance domain.
- The pilot reinforced the Coverage Model direction: the journey consumed state coverage (`membership_active`, `playback_allowed`), context coverage (`authenticated_user`), and evidence coverage (DB readback, payment event, authenticated record/profile evidence, protected resource readback).
- Evidence selection mattered more than context taxonomy expansion. Client Component HTML was not stable proof for `membership_active`; DB readback was the correct outcome proof. Protected HLS key / protected resource readback was stronger proof for `playback_allowed` than page HTML.
- The browser Capture Envelope remains local to browser evidence. The journey pilot needed claim / side-effect chain / outcome proof / protected resource readback, which matches shared validation reasoning pressure more than browser capture metadata.
- The generic Evidence Envelope remains deferred. The common shape is still claim -> evidence chain -> verification -> scope/result, but the concrete minimum schema is not stable across browser captures, API/DB readback, runtime traces, and journey validation.
- Typed Context Taxonomy remains deferred. The new pressure is outcome/evidence depth, not new context families such as appearance, environment, interaction, or locale.

Consequence:

- Keep this plan as an evolution watch-list.
- Do not add `experience_context`, generic `context`, or generic `evidence_bundle` runtime schema from this pilot alone.
- If multiple non-UI workflow domains later reuse `state/context/evidence coverage + validation_scope`, consider promoting that model under shared Validation Reasoning before considering UI-specific taxonomy changes.

### Phase 4 — Graduation Adjudication（2026-06-15）

> **Reframe**：Phase 4 的工作不是發明 taxonomy，而是**嘗試推翻 promotion gate**。下表對每個
> candidate family 問三題；**只有三題全 Yes 才升**（gate：≥3 families × 每 family ≥2 scenarios ×
> ≥2 consuming domains × cross-domain）。Evidence base = Vidoe-Test pilots（fixed-bottom /
> player→drama nav / membership_purchase journey）+ 該專案 frontend-contracts。

| Family | ≥2 governed scenarios? | 被 ≥2 domains 消費? | 不做 typed context 會失去可判定性? | 升? |
|---|---|---|---|---|
| **render** | Yes（H5 fixed-bottom、mobile-only layout） | No（僅 UI governance） | No — 已有 `Responsive` domain + Capture Envelope metadata 判定 | ❌ |
| **interaction** | No（1：player→drama→back route/history trace） | No | No — interaction-trace evidence 即可判定 | ❌ |
| **appearance** | No（1：theme switch；`alt` 為 placeholder token，品牌未確認） | No（單一 multi-theme 前端契約） | No — `multi-theme.md` 契約直接治理 | ❌ |
| **accessibility** | No（aria- 散見於 player tests；a11y-scan 是 **evidence artifact**，非 context family） | — | No | ❌ |
| **locale** | No（i18n 是基礎設施，無 locale-**validation** scenarios；非 RTL/翻譯完整性治理） | — | No | ❌ |
| **environment** | No（offline/network 散見、未治理） | — | No | ❌ |

**裁決：0 個 family 通過三題。** 所有觀察到的壓力都被 **evidence depth / coverage / 既有 domain**
解掉，不是被「缺 context family」造成。收斂結論：`Experience Validation Pipeline ≠ Context
Explosion；≈ Coverage + Evidence Selection + Outcome Proof`。

```yaml
# Branch A（採用）— typed context taxonomy 不升
decision:
  typed_context_taxonomy: rejected_for_now
reason:
  - evidence pressure resolved by evidence depth (DB readback / protected-resource readback)
  - coverage model (state/context/evidence) sufficient
  - context families not converged (0/6 pass the 3-question gate)
retain:
  - Responsive (workflow-local domain)
  - Capture Envelope (browser-evidence-scoped)
promote_only:
  - validation reasoning (state/context/evidence coverage), 待 ≥3 non-UI domains 復用後再議
reopen_triggers:
  - ≥3 context families each with ≥2 scenarios consumed by ≥2 domains
  - a governed scenario that loses decidability without a typed context

# Branch B（未觸發）— 保留以備將來
decision:
  typed_context_taxonomy: accepted
requirements:
  - ">=3 families"
  - ">=2 domains"
  - reusable artifact contract（render used by responsive+accessibility, etc.）
```

此裁決收束 §Phase 4 開頭的 6 個 open items（inventory / decide / gate / cross-domain / draft / close）為 **Branch A**；那 6 個 checkbox 已於原處勾選並指回本表。

## Completion Criteria

- [x] Open questions are either resolved or explicitly deferred with evidence. — all 8 OQs resolved/deferred (see §Open Questions).
- [x] Any docs changed by Phase 1 are validated with lints and `git diff --check`. — Phase 1 docs (`validation.md`) landed in prior commits; `git diff --check` clean (2026-06-15).
- [x] If runtime-indexed docs or scenarios change, runtime compile / refresh / validate pass. — this closure touches only this plan md (not a runtime-indexed source); `ai-skill runtime validate` success.
- [x] `plans/README.md` status row is updated. — status → completed; summary reflects Branch A closure.
- [x] Commit / push completed if implementation phases are executed. — closure commit pushed (no runtime/code change; adjudication + bookkeeping only).

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
- Complements [`archived/2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime.md`](2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime.md): that plan owns L3 `evidence_type` catalog + gate vocabulary; this plan owns journey / experience-validation pipeline execution — do not collapse `collection_method` into `evidence_type`.
