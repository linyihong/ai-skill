---
id: 2026-06-08-1544-evidence-acquisition-layer
status: completed (auto-detected)
owner: linyihong
created: 2026-06-08
priority: P2
---

# Evidence Acquisition Layer

## Phase 1 of Validation Reasoning Taxonomy

**Status**: `completed (auto-detected)`
Owner: framework maintainer (linyihong)
**世代**：Gen 3 software-delivery workflow hardening
**建立日期**：2026-06-08
**最後更新**：2026-06-08（Plan Completion Closure complete）
**Priority**：**P2**

This plan extends the `sd-ui-governance` / `sd-validation` taxonomy with an explicit **Evidence Acquisition Layer** as Phase 1 of a broader **Validation Reasoning Taxonomy**. The goal is to model how evidence is collected before it is evaluated, without promoting Browser Review into a governance domain, validation mechanism, standalone slice, runtime YAML, or enforcement rule.

---

## Decision Rationale

### Problem & Why Now

The completed [`UI Governance Workflow Integration`](../archived/2026-06-08-1408-ui-governance-workflow.md) plan introduced the useful distinction between:

```yaml
domain: what is governed
mechanism: how evidence is evaluated
evidence_class: what evidence/result is produced or trusted
```

That model correctly prevents tools such as screenshot diff, AI review, or manual review from becoming governance domains. However, it leaves one layer implicit: **how evidence is acquired**.

This creates a taxonomy gap around Browser Review and similar runtime inspection activities:

```text
Open Browser
  -> Navigate
  -> Observe actual state
  -> Collect screenshot / DOM / accessibility tree / interaction trace
  -> Feed deterministic, screenshot_diff, ai_review, or manual_review mechanisms
```

Browser Review is not the core missing abstraction. Browser Review is one acquisition method inside a broader Evidence Acquisition Layer. The same gap also applies to contract readback, static analysis, runtime trace, telemetry, production signals, and feedback records.

This plan is no longer only a UI governance extension. The examples already include `contract_readback`, `runtime_trace`, `feedback_record`, `telemetry`, and `production_signal`, which are not UI-specific. UI governance is the first consumer because it exposed the Browser Review ambiguity, but the abstraction points toward shared validation reasoning.

This also exposes a second taxonomy boundary: the current `evidence_class` list is not symmetric. It mixes evidence bodies, acquisition artifacts, evaluation results, evaluation mechanisms, and finding outputs:

| Current item | More precise role |
|---|---|
| `contract` / `runtime` | evidence body / state source |
| `screenshot` | acquisition artifact |
| `accessibility_scan` | evaluation result / finding source |
| `visual_diff` | comparison result |
| `ai_review` | evaluation mechanism or review output depending on usage |
| `human_review` | review output / finding source |

Phase 1 should not solve this whole split. It should name the problem and avoid making it worse while validating whether the acquisition layer adds value.

The deeper owner question is where this taxonomy belongs. If `collection_method` is part of the classification tuple:

```yaml
domain: behavior
collection_method: runtime_trace
mechanism: deterministic
evidence_class: runtime
```

then it should not be owned only by `sd-validation`. It belongs to the governance classification model, while validation owns execution of acquisition and evaluation. Since the same tuple already works for architecture, runtime, documentation, and UI examples, the long-term owner should be shared `validation-reasoning`, not `sd-ui-governance`.

### Decision

Add `collection_method` as a first-class taxonomy field between `domain` and `mechanism`:

```yaml
domain: what is governed
collection_method: how evidence is acquired
mechanism: how evidence is evaluated
evidence_class: what evidence/result is produced
```

Phase 1 ownership split:

```text
sd-ui-governance
  = local UI classification usage
    domain + collection_method + mechanism + evidence_class

sd-validation
  = execution model
    acquire evidence + run evaluation + record result
```

Long-term ownership target:

```text
validation-reasoning
  = owns shared classification taxonomy
    governance domain + acquisition + artifact + evaluation + finding

UI Governance / Architecture Governance / Runtime Governance / Documentation Governance
  = consumers of the shared taxonomy
```

In other words, `sd-ui-governance` owns local UI usage in Phase 1 only. `sd-validation` owns how acquisition and evaluation are actually executed during validation. The intended mature state is: `validation-reasoning` owns the shared taxonomy, governance slices consume it, and validation executes it.

This owner split is a Phase 1 decision, not an open question:

```text
UI governance owns local classification usage.
Validation owns execution.
Shared validation-reasoning is the expected long-term taxonomy owner if Phase 1 proves reusable.
```

Phase 1 landing starts with these collection methods:

```yaml
collection_method:
  - contract_readback
  - static_analysis
  - runtime_trace
  - browser_review
  - human_observation
```

This list intentionally mixes system acquisition methods and interaction channels in Phase 1. The plan records that as a deferred refinement instead of adding a fifth taxonomy layer now:

```yaml
future_refinement:
  collection_method:
    - contract_readback
    - static_analysis
    - runtime_trace
  collection_channel:
    - browser
    - human
    - telemetry
```

Example mappings:

```yaml
- domain: accessibility
  collection_method: browser_review
  mechanism: deterministic
  evidence_class: accessibility_scan

- domain: design_system
  collection_method: browser_review
  mechanism: screenshot_diff
  evidence_class: visual_diff

- domain: contract
  collection_method: contract_readback
  mechanism: deterministic
  evidence_class: contract
```

Add the following invariants to the workflow taxonomy:

- Browser Review is not a governance domain.
- Browser Review is not a validation mechanism.
- Browser Review is an evidence collection method.
- A collection method may support multiple mechanisms.
- A mechanism may consume evidence from multiple collection methods.
- Evidence acquisition must not imply validation success by itself.
- `sd-validation` may execute acquisition, but it must not redefine the UI governance taxonomy.

### Alternatives Considered

- **A. Add `sd-browser-review` as a standalone slice** — reject for first landing because Browser Review is a method, not an owner-layer workflow stage. Promoting it to a slice would make one acquisition method look more canonical than contract readback, static analysis, or runtime trace.
- **B. Add `browser_review` as a validation mechanism** — reject because Browser Review collects artifacts; it does not evaluate them. It can feed `deterministic`, `screenshot_diff`, `ai_review`, or `manual_review`.
- **C. Add `browser_review` as a governance domain** — reject because domains answer what is governed, while Browser Review answers how evidence is obtained.
- **D. Add Evidence Acquisition Layer as workflow taxonomy field** — accept. This preserves the existing domain/mechanism separation and gives Browser Review, contract readback, static analysis, runtime traces, and human observation a shared layer.
- **E. Split `collection_method` and `collection_channel` immediately** — reject for first landing because it would add a fifth dimension before the four-layer model has been validated. Keep this as a named future refinement.

### Why Not an ADR Yet

This is a workflow taxonomy refinement, not a cross-system irreversible architecture decision. The initial landing should prove that `collection_method` improves UI governance and validation evidence classification without creating taxonomy overhead.

ADR promotion is deferred until the taxonomy is used across multiple workflow domains or becomes a stable runtime/enforcement contract.

### ADR Promotion Criteria（completed 時驗證）

ADR promotion is deferred. This plan lands a workflow taxonomy refinement and validation scenarios, not an accepted cross-system validation theory.

| Criterion | Status | Evidence / revisit condition |
|---|---|---|
| `collection_method` is used in real UI governance or validation evidence | pass | Phase 1 updates `ui-governance.md`, `validation.md`, executable gates, artifact gates, template, checklist, and README references |
| Browser Review remains modeled as acquisition, not domain or mechanism | pass | UI governance, validation execution guidance, artifact gates, checklist, and scenarios all classify Browser Review as acquisition |
| At least two non-browser acquisition methods are exercised (`contract_readback`, `static_analysis`, or `runtime_trace`) | pass | Phase 2 scenarios cover `contract_readback`, `static_analysis`, and `runtime_trace` |
| Evidence Source Taxonomy refinement need is either resolved or explicitly deferred with examples | pass | Taxonomy Evolution Risk and future Artifact / Evaluation / Finding split are recorded |
| Ownership graduation from UI-local usage to shared validation-reasoning is either deferred or promoted to a follow-up | pass | Deferred to closure/future validation-reasoning follow-up; no shared owner extraction in Phase 1 |
| Finding standardization is either deferred or promoted to a shared validation-reasoning follow-up | pass | Deferred as likely future shared validation-reasoning work |
| No standalone `sd-browser-review`, `rule_class: browser_review`, or `runtime/browser-review.yaml` is introduced by accident | pass | No new slice, runtime YAML, or enforcement rule_class added |

### Consequences

#### Positive

- Gives Browser Review an explicit place without over-promoting it.
- Makes evidence chains clearer: collection -> evaluation -> evidence/result.
- Creates a path for future tools such as Playwright, Chrome MCP, Browser Use, OpenAI Computer Use, Claude Browser, telemetry, and production signals without forcing each tool into its own slice.
- Reduces confusion between acquisition artifacts, validation outputs, and review outputs.

#### Negative

- Adds one more taxonomy field to UI governance classification.
- Existing templates, gates, and scenarios need a small but broad linked update.

#### Risks

- `collection_method` may be misread as proof quality. Mitigation: define that acquisition alone never proves compliance.
- Evidence class names currently mix acquisition artifacts (`screenshot`), validation outputs (`visual_diff`, `accessibility_scan`), and review outputs (`ai_review`, `human_review`). Mitigation: record this as a future taxonomy refinement, not part of Phase 1.
- Collection methods currently mix acquisition mode and interaction channel. Mitigation: defer `collection_channel` until the four-layer model produces enough examples.
- Taxonomy Evolution Risk: Phase 1 may accidentally ossify `evidence_class` before Artifact / Evaluation / Finding boundaries are clarified. If teams start relying on mixed `evidence_class` values for enforcement, analytics, dashboards, or trend reporting, later migration becomes expensive. Mitigation: keep Phase 1 wording explicit that `evidence_class` is provisional, and require completion closure to decide whether finding standardization or shared validation-reasoning follow-up is needed.

### Future Refinement: Evidence Pipeline Taxonomy

This plan intentionally stops at a four-layer taxonomy for Phase 1:

```text
Governance Domain
  -> Collection Method
  -> Mechanism
  -> Evidence Class
```

The longer-term model may become a fuller evidence pipeline:

```text
Governance Domain
  -> Acquisition
  -> Artifact
  -> Evaluation
  -> Finding
```

Example mappings:

```yaml
- domain: accessibility
  acquisition: browser_review
  artifact: accessibility_tree
  evaluation: deterministic
  finding: accessibility_scan

- domain: design_system
  acquisition: browser_review
  artifact: screenshot
  evaluation: screenshot_diff
  finding: visual_regression
```

Potential split:

```yaml
artifact:
  - screenshot
  - dom_snapshot
  - interaction_trace

evaluation_result:
  - visual_diff
  - accessibility_scan

review_output:
  - ai_review
  - human_review
```

Potential finding taxonomy:

```yaml
finding:
  - missing_loading_state
  - missing_error_state
  - raw_design_token
  - visual_regression
  - accessibility_violation
```

Findings are the unit most likely to support warning/block severity, enforcement, analytics, and trend reporting. Screenshots, browser review, telemetry, and deterministic evaluation explain how a finding was acquired or produced; they are not the durable thing being counted.

This is likely reusable beyond UI governance and may later move into `intelligence/engineering/execution/validation-reasoning/`. Do not implement this split in Phase 1 unless the workflow edits become ambiguous without it.

Potential long-term shared model:

```text
validation-reasoning/
  ├── governance domain
  ├── acquisition / collection method
  ├── artifact taxonomy
  ├── evaluation mechanism
  └── finding taxonomy
```

UI Governance, Architecture Governance, Runtime Governance, and Documentation Governance would then reference this shared evidence theory instead of each inventing local acquisition/evaluation/finding terminology.

Example future consumers:

```yaml
- domain: architecture
  collection_method: contract_readback
  mechanism: deterministic
  evidence_class: architecture_contract

- domain: runtime
  collection_method: telemetry
  mechanism: deterministic
  evidence_class: runtime_signal

- domain: documentation
  collection_method: static_analysis
  mechanism: deterministic
  evidence_class: coverage_report
```

Glossary Impact: yes — candidate framework terms: `evidence_acquisition_layer`, `collection_method`, `browser_review`. Register only if Phase 1 confirms the terms should be reusable beyond this plan; otherwise keep them local to `workflow/software-delivery/`.

---

## Runtime Execution Path

Runtime owner: `workflow/software-delivery/`.

Trigger flow:

```text
UI compliance claim / validation evidence claim
  -> workflow/software-delivery execution flow loads `sd-ui-governance` and `sd-validation`
  -> `sd-ui-governance` applies Phase 1 local classification: domain + collection_method + mechanism + evidence_class
  -> `sd-validation` executes acquisition/evaluation and records how evidence was obtained and evaluated
  -> future shared validation-reasoning may own the cross-governance taxonomy if reuse is proven
  -> closure can distinguish "evidence collected" from "claim validated"
```

Planned runtime-facing changes:

- Update `workflow/software-delivery/ui-governance.md` taxonomy and minimum review fields.
- Update `workflow/software-delivery/validation.md` with Browser Evidence Collection as an execution subsection under validation.
- Update `workflow/software-delivery/execution-flow.yaml` so `classify_ui_governance` and relevant evidence gates mention `collection_method`.
- Update artifact/template/checklist surfaces if they mention UI evidence fields.
- Add validation scenarios proving Browser Review is acquisition, not mechanism/domain.

Not planned:

- No new `sd-browser-review` slice.
- No `rule_class: browser_review`.
- No `runtime/browser-review.yaml`.
- No Playwright / Chrome MCP / Browser Use / Computer Use integration.
- No mechanical validator executor in this landing.

### Deferred Runtime Projection

No standalone `runtime/*.yaml` is planned. This remains a workflow taxonomy and executable workflow contract update. If a future phase introduces a shared `runtime/evidence-acquisition*.yaml`, it must name a runtime consumer, projection target, evidence threshold, and graduation condition.

### Per-surface Consumer 表

| Generated surface key | Named consumer(s) | Consumer 類型 |
|---|---|---|
| `workflow.software_delivery.execution_flow.contract` update | software-delivery workflow route / executable contract validation | existing routable workflow contract |
| validation scenarios under `validation/scenarios/software-delivery/` | runtime validation scenario inventory | validation scenario |

---

## Open Questions

- [x] Should taxonomy ownership graduate from UI-local usage to shared `validation-reasoning`? **Deferred**: Phase 1 remains UI-local; closure must decide whether to open a shared validation-reasoning follow-up.
- [x] Should `contract_readback` include generated surface readback and runtime refresh validation, or should those become separate methods later? **Resolved for Phase 1**: keep `contract_readback` broad in UI governance examples; do not split generated/runtime readback subtypes yet.
- [x] Should future taxonomy split `evidence_artifact` from `evidence_class` to separate screenshots / DOM snapshots from validation outputs and review results? **Deferred**: record Taxonomy Evolution Risk; do not add a fifth layer in Phase 1.
- [x] Should findings be standardized for severity, enforcement, analytics, and trend reporting? **Deferred**: finding taxonomy is a likely follow-up, not part of Phase 1.
- [x] Which collection methods should be first-class in Phase 1 vs. listed as future candidates (`telemetry`, `production_signal`, `feedback_record`)? **Resolved for Phase 1**: first-class methods are `contract_readback`, `static_analysis`, `runtime_trace`, `browser_review`, and `human_observation`; telemetry/production/feedback stay future candidates.
- [x] Should `collection_method` later split from `collection_channel`, or is the added precision not worth a fifth taxonomy layer? **Deferred**: keep one `collection_method` field in Phase 1.
- [x] Should the long-term reusable model become `Governance Domain -> Acquisition -> Artifact -> Evaluation -> Finding` under `intelligence/engineering/execution/validation-reasoning/`? **Deferred**: record as expected long-term direction; do not extract in Phase 1.

---

## Architecture Compatibility Preflight

### Phase 0.0 — Open Questions 核對（公版，必填）

逐條核對本 plan §Open Questions，標記處置並回寫：

- [x] 已讀本 plan §Open Questions 全部條目
- [x] 對每條標記 `resolved`（附 Phase 0 證據）/ `still-open` / `deferred`（附原因）
- [x] `resolved` 的條目已同步勾選 / 附註於 §Open Questions
- [x] 若盤點新發現問題，已加入 §Open Questions

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| `collection_method` owner layer | resolved | Phase 1: `sd-ui-governance` owns local classification usage; `sd-validation` owns execution |
| Ownership graduation to validation-reasoning | deferred | Closure must decide whether to open shared validation-reasoning follow-up |
| `contract_readback` scope | resolved | Keep broad in Phase 1; do not split generated/runtime readback subtypes yet |
| `evidence_artifact` split | deferred | Record Taxonomy Evolution Risk; no fifth layer in Phase 1 |
| Finding standardization | deferred | Candidate follow-up for enforcement/analytics/trend reporting |
| Phase 1 first-class methods | resolved | `contract_readback`, `static_analysis`, `runtime_trace`, `browser_review`, `human_observation` |
| `collection_channel` split | deferred | Avoid fifth taxonomy layer in Phase 1 |
| Evidence pipeline taxonomy | deferred | Expected long-term direction under validation-reasoning; no extraction in Phase 1 |

### Candidate Files

| Path | Role |
|---|---|
| `workflow/software-delivery/ui-governance.md` | Add `collection_method` taxonomy and invariants |
| `workflow/software-delivery/validation.md` | Add Evidence Acquisition / Browser Evidence Collection subsection |
| `workflow/software-delivery/execution-flow.yaml` | Add `collection_method` to UI governance classification step and gates |
| `workflow/software-delivery/artifact-gates.md` | Add collection method to UI governance evidence shape if needed |
| `workflow/software-delivery/artifact-gates.yaml` | Add executable artifact gate field if needed |
| `workflow/software-delivery/templates/ui-governance-evidence-template.md` | Add collection method field |
| `workflow/software-delivery/review-checklist.md` | Add review check that Browser Review is acquisition, not validator |
| `validation/scenarios/software-delivery/` | Add scenarios for evidence acquisition taxonomy |
| `plans/README.md` | Track this plan while active |

### Compatibility Checklist

| Check | Result | Notes |
|---|---|---|
| Candidate files exist | pass | Read `ui-governance.md`, `validation.md`, `execution-flow.yaml`, artifact gates, template, checklist, and scenario schema |
| Source-of-truth consistency | pass | Workflow markdown + executable YAML remain canonical for this landing |
| Layer responsibility | pass | Phase 1 updates workflow taxonomy/execution only; no runtime/enforcement rule_class |
| Compiler / generated surface | pass | Runtime refresh required after executable YAML or validation scenario changes |
| Linked updates | pass | Artifact gates, template, checklist, execution flow, and scenarios are in scope |
| Execution decision | proceed | Proceed with Phase 1 taxonomy landing and Phase 2 scenarios |

---

## Phase 0 — Pre-Build Interrogation / Compatibility

- [x] Complete Open Questions disposition.
- [x] Confirm `sd-ui-governance` vs. `sd-validation` ownership split:
  - `sd-ui-governance`: owns Phase 1 local UI compliance classification usage (`domain`, `collection_method`, `mechanism`, `evidence_class`).
  - `sd-validation`: owns acquisition/evaluation execution and result recording during validation.
  - future `validation-reasoning`: expected long-term owner if Phase 1 validates reuse beyond UI governance.
- [x] Confirm no new standalone slice is needed.
- [x] Confirm no runtime YAML or enforcement rule_class is needed.
- [x] Confirm linked update list before implementation.

Acceptance:

- [x] Phase 0 records decision to proceed / revise / stop.
- [x] Open Questions are updated with disposition.
- [x] Scope remains Phase 1 taxonomy landing only.

---

## Phase 1 — Workflow Taxonomy Landing

- [x] Update `workflow/software-delivery/ui-governance.md`:
  - [x] Add `collection_method` between domain and mechanism.
  - [x] Add supported methods: `contract_readback`, `static_analysis`, `runtime_trace`, `browser_review`, `human_observation`.
  - [x] Add invariants for Browser Review and acquisition/evaluation separation.
  - [x] Update Minimum UI Governance Review field list.
- [x] Update `workflow/software-delivery/validation.md`:
  - [x] Add Evidence Acquisition execution subsection.
  - [x] Add Browser Evidence Collection examples and outputs:
    - `screenshot`
    - `dom_snapshot`
    - `accessibility_scan`
    - `interaction_trace`
    - `responsive_capture`
  - [x] Name consumers: `ui_governance`, `screenshot_diff`, `ai_review`, `accessibility_validator`, `behavior_validation`.
- [x] Update `workflow/software-delivery/execution-flow.yaml`:
  - [x] Add `collection_method` to `classify_ui_governance`.
  - [x] Add gate language requiring acquisition method for UI evidence claims.
- [x] Update artifact gates / template / review checklist surfaces as linked updates if current fields omit acquisition method.

Acceptance:

- [x] Browser Review is explicitly not a domain or validation mechanism.
- [x] Evidence acquisition and evidence evaluation are separate.
- [x] `sd-ui-governance` owns Phase 1 local classification usage; `sd-validation` owns execution.
- [x] Completion evidence can state how UI evidence was acquired.
- [x] No new `sd-browser-review`, runtime YAML, or rule_class is introduced.

---

## Phase 2 — Validation Scenarios

- [x] Add scenario: Browser Review used as `collection_method` feeding `ai_review`, not as mechanism.
- [x] Add scenario: Contract readback used as `collection_method` feeding deterministic contract validation.
- [x] Add scenario: Static analysis or runtime trace feeding deterministic accessibility / behavior validation.
- [x] Add scenario: A completion claim with mechanism and evidence class but missing collection method receives a warning or classification failure, depending on existing scenario schema capabilities.
- [x] Run runtime refresh / validate.
- [x] Check lints for updated markdown and YAML files.

Acceptance:

- [x] Scenario coverage distinguishes domain, collection method, mechanism, and evidence class.
- [x] Scenario wording does not imply `browser_review` is a validator or that acquisition alone proves a finding.
- [x] Scenarios avoid introducing tool-specific runtime integrations.
- [x] Runtime validation passes.

---

## Phase 3 — Closure / Archive

- [x] Review linked updates against `enforcement/linked-updates.md`.
- [x] Update this plan status to `completed (auto-detected)` after Phase 1–2 completion.
- [x] Update `plans/README.md`.
- [x] Execute Plan Completion Closure.
- [x] Archive plan to `plans/archived/`.
- [x] Commit and push when phase completion is valid.

Acceptance:

- [x] Plan closure evidence records runtime validation and linked update review.
- [x] Working tree is clean and pushed after archive commit.

---

## 完成條件

- [x] `collection_method` is documented as Evidence Acquisition Layer.
- [x] `browser_review` is modeled only as an acquisition method.
- [x] `sd-ui-governance` owns Phase 1 local usage of the four-layer UI compliance classification taxonomy.
- [x] `sd-validation` documents Browser Evidence Collection execution outputs and consumers.
- [x] UI governance minimum evidence classification includes acquisition method.
- [x] Executable workflow contract reflects the new field.
- [x] Validation scenarios prove taxonomy separation.
- [x] `collection_channel` and `Artifact -> Evaluation -> Finding` split are explicitly deferred or promoted with justification.
- [x] Ownership graduation to shared validation-reasoning is explicitly deferred or captured as a follow-up.
- [x] Finding standardization is explicitly deferred or captured as a shared validation-reasoning follow-up.
- [x] No standalone Browser Review slice, runtime YAML, or enforcement rule is added.
- [x] Plan Completion Closure executed when all phases are done.

---

## Stakeholder 同意項目

- [x] Accept Evidence Acquisition Layer as the missing abstraction, not Browser Review as a standalone slice.
- [x] Accept `collection_method` between `domain` and `mechanism`.
- [x] Accept `sd-ui-governance` as Phase 1 local taxonomy consumer and `sd-validation` as owner of execution.
- [x] Accept Phase 1 as workflow taxonomy landing only.
- [x] Accept `evidence_artifact` split as future refinement, not part of first landing.
- [x] Accept `collection_channel` split as future refinement, not part of first landing.
- [x] Accept shared Validation / Evidence Taxonomy as the expected long-term owner, with graduation timing deferred beyond Phase 1.
- [x] Accept standardized finding taxonomy as a likely future direction, not part of Phase 1 landing.
- [x] Accept no tool integration in this plan.

---

## Plan Completion Closure

| Check | Result | Evidence |
| --- | --- | --- |
| All phases complete | pass | Phase 0–3 checklists and acceptance items are complete. |
| Runtime validation | pass | `ai-skill runtime refresh`, `ai-skill runtime validate`, YAML parse, lints, and `git diff --check` passed. |
| Linked updates reviewed | pass | `ui-governance.md`, `validation.md`, `execution-flow.yaml`, artifact gates, UI evidence template, template index, review checklist, README, and validation scenarios were synchronized. |
| Runtime / enforcement boundary | pass | No standalone `sd-browser-review`, `runtime/browser-review.yaml`, or enforcement rule_class added. |
| Shared taxonomy boundary | deferred | Ownership graduation to shared `validation-reasoning`, finding taxonomy, and Artifact / Evaluation / Finding split are recorded as future follow-up decisions. |
| Archive action | pass | Plan moved from `plans/active/` to `plans/archived/` after closure. |

---

## 與其他 plans 的關係

- Related: [`2026-06-08-1408-ui-governance-workflow.md`](2026-06-08-1408-ui-governance-workflow.md) — parent context for `sd-ui-governance` taxonomy and validation scenarios.
- Related: [`../active/2026-06-06-1700-workflow-activation-discovery-bridge.md`](../active/2026-06-06-1700-workflow-activation-discovery-bridge.md) — workflow discovery behavior may later affect when validation/evidence acquisition surfaces are loaded.
- Related: [`../active/2026-06-06-1800-sanitization-mechanical-enforcement.md`](../active/2026-06-06-1800-sanitization-mechanical-enforcement.md) — same theme of keeping method/tool evidence in the correct layer before mechanical enforcement.
