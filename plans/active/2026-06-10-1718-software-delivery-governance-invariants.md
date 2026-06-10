---
id: 2026-06-10-1718-software-delivery-governance-invariants
plan_kind: main
status: in-progress
owner: linyihong
created: 2026-06-10
priority: P1
required_for_completion: false
---

# Software Delivery Governance Invariants

**Status**: `in-progress`
Owner: framework maintainer (linyihong)
**建立日期**：2026-06-10
**Priority**：**P1**（software-delivery workflow hardening；把 incident feedback 提煉成 governance invariants）

## Why this plan exists

Recent software-delivery feedback listed concrete incidents such as browser API fallback, share counters, localhost/base URL deployment drift, nested repo dirty states, interrupted deploys, and browser smoke tests. These are real issues, but most of them are symptoms.

The higher-level gap is that `workflow/software-delivery` needs reusable validation-family invariants for:

- assumptions about runtime capabilities
- business-truth ownership for side effects
- deployed-state readback instead of config-input validation
- ownership boundaries before committing across repos or ownership groups
- operational transactions that can be partial, interrupted, or unverified
- post-deploy validation at user-journey level

Activation and escalation are related, but not the same layer: they are routing pressure. They explain why the Workflow Activation Engine must route to these invariants at the right time, but they should not be treated as validation-family invariants themselves.

This plan avoids adding one rule per incident (`navigator.share`, clipboard, localhost, deploy log, share count). Instead it extracts a small set of reusable invariants that align with existing framework direction:

- Evidence-Oriented Validation
- State Visibility Gap
- Failure Authority
- Evidence Chain Validation
- Journey Validation
- Contract First
- Workflow Activation / Routing / Economics

## Decision Rationale

### Problem & Why Now

The workflow already has strong coverage for contract-first development, UI governance, responsive evidence, journey validation, and state visibility gaps. However, multiple recent incidents show recurring workflow misses:

```text
workflow exists but is not activated
  -> agent patches from surface symptom

task starts as small UI fix
  -> runtime capability / side-effect / journey scope appears
  -> workflow does not escalate

root cause assumed silently
  -> patch precedes evidence

assumed capability exists
  -> no fallback validated

side effect counted too early
  -> fake success metric

configuration input checked
  -> deployed state not read back

mixed working tree committed
  -> unrelated ownership included

deploy command started
  -> completion state ambiguous

component/API smoke passes
  -> critical user flow still fails
```

These are not separate categories of bugs. They are governance invariant gaps around authority, readback, ownership, and operational closure, plus routing pressure around activation/escalation and evidence discipline around root-cause hypotheses.

### Decision

Draft reusable validation-family invariants for `workflow/software-delivery`, with `authority-coupled-side-effects`, `configuration-readback-validation`, and `operational-transaction-closure` treated as the strongest current candidates. `runtime-capability-validation` is also a high-probability candidate because it can cover browser, platform, filesystem, container, and orchestration capabilities, but Phase 0 must validate breadth before placement. Treat workflow activation / escalation as routing pressure that must cross-link rather than duplicate the Workflow Activation Engine. Treat `explicit-root-cause-hypothesis` as evidence acquisition discipline, not as a software-delivery invariant unless Phase 0 proves it belongs here.

| Candidate | Layer | Prevents | Core question |
|---|---|---|---|
| `runtime-capability-validation` | validation invariant | browser / platform supported-assumption bugs | Does the runtime actually provide the capability, and is fallback validated? |
| `authority-coupled-side-effects` | validation invariant | fake success metrics / premature counters | Which event owns the business truth? |
| `configuration-readback-validation` | validation invariant | validated config input but wrong deployed state | Did we read back actual deployed state? |
| `ownership-boundary-validation` | validation / closure candidate | mixed-task commits across repos or owner groups | Is the failure about task scope, ownership awareness, or both? |
| `operational-transaction-closure` | validation / closure invariant | interrupted / partial deployment ambiguity | Did the operational transaction finish and get verified? |
| `journey-level-post-deploy-validation` | validation invariant | component pass but user flow fail | Did the critical journey pass after deploy? |
| `workflow-activation` | routing pressure | workflow exists but is not loaded | Which workflow slices must activate before implementation? |
| `workflow-escalation` | routing pressure | task starts narrow but needs broader governance | When does a task escalate from UI/CSS to runtime capability, side effect, deployment, or journey validation? |
| `explicit-root-cause-hypothesis` | evidence acquisition discipline | silent assumptions and patch-first debugging | What root cause is claimed, what evidence supports it, what patch follows, and how will it be validated? |

`ownership-boundary-validation` is a temporary working name. Phase 0 may split it into `task-scope-validation` and `ownership-awareness` if evidence shows task boundary and ownership boundary are separate models.

The intended abstraction:

```text
incident example
  -> symptom family
  -> governance invariant
  -> workflow gate / artifact shape / validation scenario
```

This is a planning step. It should not immediately add broad runtime gates, generic context schemas, or per-API rules.

### Alternatives Considered

- A. Add `browser-capability-fallback` as a specific rule: reject. `navigator.share`, clipboard, fullscreen, autoplay, WebAuthn, Payment Request, camera, geolocation, Bluetooth, and notifications are all examples of `runtime-capability-validation`.
- B. Add `side-effect-counting-semantics` as a low-level counter rule: reject as too narrow. The real question is authority: button click, API call, API 200, DB commit, external provider acceptance, and consumer-observed result have different authority levels.
- C. Add deploy-env checks for localhost/base URL: reject as too narrow. The stronger invariant is `configuration-readback-validation`: desired configuration must be read back from actual runtime/deployed state.
- D. Treat nested repo dirty state as Git-specific: reject. The reusable question is whether the change exceeds task scope, crosses ownership boundaries, or both; Phase 0 decides whether this remains one model.
- E. Add browser smoke tests as a separate category: reject. Most post-deploy smoke tests should map to `journey-level-post-deploy-validation`.
- F. Add UI fast path: reject as unsafe. Many incidents looked like small UI work but required runtime capability, side-effect authority, route/history, or journey validation. Prefer `workflow-escalation` and `complexity-reassessment`.
- G. Add change brief / checklist layer: reject as primary solution. Briefs can be useful output, but the deeper issue is activating the right workflow at the right time.
- H. Extract governance invariants, then decide their workflow placement, routing relationship, and validation evidence: accept.

### Why Not an ADR Yet

This is not ready for ADR promotion. The invariant family is plausible, but needs:

- mapping to current `workflow/software-delivery` surfaces without duplicating existing State Visibility / Journey Validation content
- mapping activation / escalation concerns to the Workflow Activation Engine without duplicating `workflow-activation-discovery-bridge`
- deciding whether `explicit-root-cause-hypothesis` belongs in software-delivery, or should move to shared evidence-acquisition / validation reasoning
- keeping `claim_validation` as a hypothesis only; it must not be introduced unless multiple invariants prove they share identical validation structure
- scenario evidence showing the invariants catch more than the original incident examples
- ownership decisions: workflow vs shared validation reasoning vs future intelligence atom
- proof that the plan does not create a catch-all taxonomy or over-broad runtime schema

### ADR Promotion Criteria（completed 時驗證）

- [ ] At least 3 invariants land with clear workflow placement and no duplicated source-of-truth.
- [ ] At least 3 validation scenarios prove the invariants catch distinct failure families.
- [ ] Workflow activation / escalation concerns are cross-linked to the Workflow Activation Engine or explicitly deferred.
- [ ] `explicit-root-cause-hypothesis` is either placed outside software-delivery or justified as a software-delivery checkpoint.
- [ ] Task boundary and ownership boundary are evaluated separately; `ownership-boundary-validation` is kept only if Phase 0 proves it is one stable model.
- [ ] Claim Validation is evaluated but not introduced unless at least three invariants share identical validation structure; similar concepts are not enough.
- [ ] `authority-coupled-side-effects`, `configuration-readback-validation`, and `operational-transaction-closure` are either implemented or explicitly deferred with evidence.
- [ ] Existing Journey Validation, State Visibility Gap, Evidence Chain Validation, and Failure Authority references remain consistent.
- [ ] No new generic runtime schema is introduced without named consumers and validation scenarios.
- [ ] If promoted beyond workflow, shared Validation Reasoning is evaluated before ADR.

### Consequences（預期）

#### 正面

- Turns incident feedback into reusable workflow invariants.
- Reduces one-off rule growth for every browser API, deployment setting, or side-effect counter.
- Makes workflow activation and escalation explicit instead of adding more task-specific checklists.
- Forces root-cause hypothesis to be stated before patching when the task is diagnostic or ambiguous.
- Strengthens completion claims around authority, readback, and operational closure.
- Aligns post-deploy checks with Journey Validation rather than component-only smoke tests.

#### 負面

- More validation evidence may be required before declaring done.
- Some operational workflows will need explicit transaction state instead of a single deploy pass/fail flag.
- Side-effect claims may need project-specific authority maps that cannot be fully generic.
- More tasks may require an activation / escalation checkpoint before implementation, increasing upfront reasoning.
- Some candidates may move out of software-delivery into shared validation reasoning, so this plan may shrink after Phase 0 rather than grow.

#### 風險

- Over-expanding the invariants could create another catch-all taxonomy.
- Authority maps can become too abstract if not tied to concrete evidence classes.
- Workflow activation concerns could duplicate the existing activation engine unless this plan treats them as cross-links and evidence pressure.
- `explicit-root-cause-hypothesis` could become a workflow-local debugging ritual even though it is cross-domain evidence acquisition discipline.
- `claim_validation` could become a premature parent taxonomy before authority, readback, closure, and journey scenarios prove identical structure.
- Operational transaction closure can overlap with CI/CD tooling and may belong under execution reasoning rather than software-delivery workflow if migrations, backfills, cache rebuilds, and deploys share the same transaction-state model.

Glossary Impact: yes — candidate terms `workflow_activation`, `workflow_escalation`, `explicit_root_cause_hypothesis`, `evidence_sufficiency_validation`, `runtime_capability_validation`, `authority_coupled_side_effects`, `configuration_readback_validation`, `ownership_boundary_validation`, `operational_transaction_closure`, `journey_level_post_deploy_validation`, `business_truth`, `observable_proxy`; `claim_validation` remains a hypothesis-only watch-list term unless Phase 0 proves at least three invariants share identical validation structure. Phase 2 decides whether accepted candidate terms should register in `knowledge/glossary/ai-skill.md` or remain plan-local.

Watch-Out List citation: Gen 4 forward scope must avoid autonomous taxonomy expansion and avoid turning incidents into broad runtime schemas without evidence. See [`architecture/ai-native-cognitive-ecosystem-system.md`](../../architecture/ai-native-cognitive-ecosystem-system.md) §Watch-Out List.

## Runtime Execution Path

Runtime owner: `workflow/software-delivery/` first; possible future promotion target is `intelligence/engineering/execution/validation-reasoning/`.

Planned trigger flow:

```text
software-delivery incident or feedback points to validation miss
  -> classify symptom family
  -> map to governance invariant
  -> decide workflow placement / artifact gate / validation scenario
  -> update software-delivery workflow docs
  -> runtime refresh indexes updated workflow/scenario sources
  -> validation confirms scenario coverage and no taxonomy overreach
```

Potential runtime surfaces:

| Surface | Intended consumer | Notes |
|---|---|---|
| `workflow/software-delivery/validation.md` | validation workflow | likely home for readback, authority, and post-deploy evidence expectations |
| `workflow/software-delivery/test-strategy.md` | test planning | likely home for capability/fallback and post-deploy journey test selection |
| `workflow/software-delivery/intake.md` | intake / routing checkpoint | possible home for root-cause hypothesis and change-boundary detection before implementation |
| `workflow/software-delivery/closure.md` | close-out workflow | likely home for operational transaction closure and change ownership close-loop |
| `workflow/software-delivery/artifact-gates.md` | artifact completeness | likely home for evidence shape of authority/readback/operational closure |
| `workflow/software-delivery/execution-flow.yaml` | executable workflow gate loader | only if Phase 1 proves a gate should become executable |
| `plans/active/2026-06-06-1700-workflow-activation-discovery-bridge.md` | workflow activation hardening plan | activation / escalation pressure should cross-link here rather than duplicate runtime routing design |
| `validation/scenarios/software-delivery/*.yaml` | runtime scenario validation | needed before implementation phases that change workflow behavior |

This plan does not add a new `runtime/*.yaml` source and does not use deferred runtime projection. If a later phase adds a runtime surface, it must add a named consumer and validation scenario in the same phase.

## Open Questions

- [ ] Which invariants belong in `workflow/software-delivery`, and which should remain in shared Validation Reasoning?
- [x] Which activation/escalation concerns belong in this plan, and which should be deferred to Workflow Activation Discovery Bridge? Phase 0: defer routing mechanics to Workflow Activation Discovery Bridge; this plan only records routing pressure.
- [x] Should activation/escalation be downgraded from invariant candidates to routing pressure only? Phase 0: yes, unless future scenario evidence proves a separate software-delivery invariant.
- [x] Should `explicit-root-cause-hypothesis` be moved out of software-delivery into shared evidence acquisition / validation reasoning? Phase 0: yes for invariant placement; keep only as evidence-acquisition discipline / diagnostic checkpoint candidate.
- [x] Is the boundary problem one concept, or should it split into `task-scope-validation` plus `ownership-awareness`? Phase 1: split for now; scenarios should test them separately.
- [x] Do authority, readback, closure, and journey validation share identical validation structure, or are they only similar concepts? Phase 0: not proven; `claim_validation` remains hypothesis-only.
- [ ] What evidence is sufficient to decide workflow escalation from UI -> runtime capability -> side effect -> journey validation?
- [x] Should `authority-coupled-side-effects` become a workflow gate, an artifact evidence shape, or an intelligence atom first? Phase 1: artifact / validation evidence shape first; no executable gate before scenarios.
- [ ] What is the minimum authority ladder for side effects without hardcoding domain-specific business events?
- [x] How should `configuration-readback-validation` relate to State Visibility Gap without duplicating it? Phase 0: treat as workflow evidence shape backed by State Visibility Gap / Evidence Chain, not a new shared reasoning atom yet.
- [ ] Does `operational-transaction-closure` belong in software-delivery workflow, or should it graduate to shared execution reasoning for deploys, migrations, backfills, cache rebuilds, and similar operations? Phase 0: likely shared execution reasoning, but needs Phase 1 placement matrix and scenarios.
- [ ] What validation scenarios are needed before any of these become executable gates?
- [ ] Should runtime capability validation be limited to user-visible surfaces first, or include backend/platform runtime capabilities too?

## Phase 0 — Open Questions Check + Architecture Compatibility Preflight

### Phase 0.0 — Open Questions 核對（公版，必填）

逐條核對本 plan §Open Questions，標記處置並回寫：

- [x] 已讀本 plan §Open Questions 全部條目
- [x] 對每條標記 `resolved`（附 Phase 0 證據）/ `still-open` / `deferred`（附原因）
- [x] `resolved` 的條目已同步勾選 / 附註於 §Open Questions
- [x] 若盤點新發現問題，已加入 §Open Questions

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| invariant placement | still-open | Current workflow slices and validation reasoning sources support candidate placement, but final owner decisions require Phase 1 matrix |
| activation / escalation placement | resolved | Workflow Activation Discovery Bridge owns detector miss -> advisory; software-delivery plan only records routing pressure |
| activation / escalation layer | resolved | Treat as routing pressure, not validation-family invariants |
| explicit root-cause hypothesis | resolved | Evidence Collapse Point says collapse point is not root cause; root-cause hypothesis belongs to evidence acquisition / diagnostic reasoning, not software-delivery invariant |
| task scope vs ownership | resolved | Phase 1 splits the working name into `task-scope-validation` and `ownership-awareness`; Phase 2 should test both failure families separately |
| Claim Validation hypothesis | resolved | Evidence Hierarchy already owns claim scope; State Visibility / Evidence Chain / Evidence Depth own proof sufficiency; no parent taxonomy introduced |
| authority side effects shape | resolved | Phase 1 places it as artifact / validation evidence shape first; no executable gate or shared intelligence atom before scenarios |
| configuration readback boundary | resolved | Treat as State Visibility Gap / Evidence Chain application in `validation.md`, not a duplicate shared concept |
| operational transaction placement | still-open | `closure.md` has DoD/close-loop, but no deploy/migration/backfill/cache rebuild transaction-state model exists under execution reasoning yet |
| validation scenarios | still-open | Phase 2 must define scenarios before executable gates |

### Phase 0.0a — Focused Decision Questions

Phase 0 should answer these four questions before expanding vocabulary:

- [ ] Are activation and escalation only routing pressure for the Workflow Activation Engine?
- [ ] Does root-cause hypothesis belong to evidence acquisition rather than software-delivery workflow?
- [ ] Is boundary validation one concept, or should it split into task scope and ownership awareness?
- [ ] Do authority, readback, closure, and journey validation share identical validation structure? If not, keep `claim_validation` as a hypothesis and do not introduce a parent abstraction.

### Phase 0.0b — Validation Targets

Phase 0 is approved to proceed when validation focuses on evidence instead of abstract fit:

- [x] Each proposed invariant should block at least two different incident families before it is treated as stable. Example: `authority-coupled-side-effects` should cover both share-count authority and payment-success authority, not just one incident.
- [x] Placement remains undecided until evidence distinguishes workflow invariant, shared validation reasoning, and shared execution reasoning. Example: `operational-transaction-closure` may belong under execution reasoning if deploys, migrations, backfills, and cache rebuilds share one transaction-state model.
- [x] `claim_validation` remains a possible parent abstraction only. Record structural similarities, but do not create the abstraction unless the three-invariant identical-structure threshold is met.

### Phase 0.0c — Pre-build Interrogation Record

- Goal: execute Phase 0 validation so this plan can enter Phase 1 placement without adding premature workflow gates, runtime projections, or parent taxonomy.
- Scope: this active plan, `plans/README.md`, current `workflow/software-delivery` surfaces, Workflow Activation Discovery Bridge, and shared validation / execution reasoning sources.
- Non-goals: no workflow doc implementation yet, no executable gate, no runtime YAML, no `claim_validation` document, no glossary registration, and no new generic runtime schema.
- Acceptance / validation target: Open Questions receive Phase 0 dispositions; unresolved items become Phase 1 or Phase 2 inputs; duplicate source-of-truth risks are documented.
- Framework discovery: canonical source is this active plan. Software-delivery owner surfaces are `intake.md`, `validation.md`, `test-strategy.md`, `closure.md`, `artifact-gates.md`, `execution-flow.yaml`, and `review-checklist.md`. Routing mechanics belong to `plans/active/2026-06-06-1700-workflow-activation-discovery-bridge.md`. Proof sufficiency already exists in `state-visibility-gap.md`, `evidence-model.md`, `evidence-chain-validation.md`, `evidence-depth.md`, and `evidence-hierarchy.md`.
- Duplication risk: high if this plan redefines State Visibility Gap, Evidence Chain, Journey Validation, Workflow Activation, or claim scope. Phase 0 decision is to cross-link or defer instead of duplicating these sources.
- Open questions: no blocker question for proceeding to Phase 1 placement matrix; unresolved questions require evidence/scenario work before Phase 3 workflow documentation changes.
- Assumptions: current workflow docs are the source of truth for software-delivery behavior; no hidden runtime projection should be added during Phase 0.
- Decision: proceed to Phase 1 placement matrix after Phase 0 writeback and validation.

### Phase 0.0d — Source Review Findings

| Candidate / concern | Phase 0 disposition | Evidence |
|---|---|---|
| `workflow-activation` / `workflow-escalation` | resolved as routing pressure | Workflow Activation Discovery Bridge explicitly keeps Discovery advisory-only and forbids proposal confidence from satisfying activation triggers. `execution-flow.yaml` already declares software-delivery activation events, so this plan should not create another routing path. |
| `explicit-root-cause-hypothesis` | evidence acquisition / diagnostic discipline, not software-delivery invariant | `evidence-collapse-point.md` distinguishes collapse point from root cause; the useful workflow check is to state hypothesis/evidence before patching, not to create a delivery invariant. |
| `ownership-boundary-validation` | temporary working name remains open | `intake.md` has change scope / parity / blocker classification; `review-checklist.md` separately reviews ownership/security boundaries. Phase 1 must decide whether one model is stable or split into `task-scope-validation` plus `ownership-awareness`. |
| `configuration-readback-validation` | resolved as workflow evidence shape backed by existing reasoning | `validation.md` already escalates state visibility gaps to State Visibility / Evidence Chain / Evidence Depth; this candidate should not duplicate those reasoning docs. |
| `authority-coupled-side-effects` | still open; likely workflow evidence shape or narrow Evidence Chain specialization | `evidence-chain-validation.md` covers payment/event/email chains, and `identity-coupled-side-effect-validation.md` is a specialization. Authority may be broader than identity, so Phase 1 needs placement. |
| `runtime-capability-validation` | high-probability candidate; placement open | `test-strategy.md` and `validation.md` cover runtime/browser proof patterns, but capability breadth spans browser, platform, filesystem, container, and orchestration surfaces. Phase 1 must test whether this lands in workflow first or shared reasoning later. |
| `operational-transaction-closure` | placement open; likely shared execution reasoning candidate | `closure.md` has DoD / close-loop, but current `intelligence/engineering/execution/` has no transaction-state reasoning doc for deploys, migrations, backfills, cache rebuilds, or external operations. |
| `journey-level-post-deploy-validation` | do not duplicate Journey Validation | `test-strategy.md`, `validation.md`, `artifact-gates.md`, and `review-checklist.md` already define BDD-owned Journey Specification and validation-owned execution/evidence. This plan may reference the existing shape. |
| `claim_validation` | hypothesis-only | `evidence-hierarchy.md` already owns claim scope, evidence quality, and confidence integrity; no three-invariant identical structure has been proven. |

### Phase 0.1 — Preflight

- [x] Read current `workflow/software-delivery/README.md`, `validation.md`, `test-strategy.md`, `closure.md`, `artifact-gates.md`, `execution-flow.yaml`, and `review-checklist.md`.
- [x] Read `plans/active/2026-06-06-1700-workflow-activation-discovery-bridge.md` and decide how activation / escalation pressure should cross-link.
- [x] Read shared validation reasoning sources: `state-visibility-gap.md`, `evidence-model.md`, `evidence-chain-validation.md`, `evidence-depth.md`, and any Failure Authority source if present.
- [x] Check shared execution reasoning sources before keeping `operational-transaction-closure` inside software-delivery.
- [x] Search for existing Claim Validation, Failure Authority, authority classification, or evidence sufficiency vocabulary before adding new terms.
- [x] Confirm no duplicate vocabulary already exists in glossary or current workflow.
- [x] Decide whether this is a sibling plan to Experience Validation Pipeline or should cross-link as a downstream hardening plan.
- [x] Confirm no new runtime projection is needed for Phase 0.
- [x] Record not-applicable dependencies explicitly if a source does not exist.

Not applicable / missing source notes:

- No dedicated `claim_validation` source exists and Phase 0 does not create one.
- No dedicated shared execution reasoning document currently owns operational transaction state; Phase 1 must decide whether to create one later.
- No separate Failure Authority source was found under `intelligence/engineering/execution/validation-reasoning/`; existing authority vocabulary appears in `evidence-hierarchy.md` and recent active/archived governance plans.

## Phase 1 — Invariant Placement Draft

- [x] Create a placement matrix that separates validation-family invariants, routing pressure, and evidence acquisition discipline:
  - workflow-activation
  - workflow-escalation
  - explicit-root-cause-hypothesis
  - evidence-sufficiency-validation
  - runtime-capability-validation
  - authority-coupled-side-effects
  - configuration-readback-validation
  - ownership-boundary-validation
  - operational-transaction-closure
  - journey-level-post-deploy-validation
- [x] Decide whether `ownership-boundary-validation` remains one model or splits into `task-scope-validation` plus `ownership-awareness`.
- [x] Evaluate whether authority, readback, operational closure, and journey validation share identical validation structure. Keep `claim_validation` hypothesis-only unless at least three invariants satisfy that threshold.
- [x] For each invariant, decide primary owner surface: `intake.md`, `validation.md`, `test-strategy.md`, `closure.md`, `artifact-gates.md`, `execution-flow.yaml`, Workflow Activation plan, shared validation reasoning, or defer.
- [x] Define non-goals so incident-specific examples do not become hardcoded framework rules.
- [x] Update this plan with the placement decision before editing workflow docs.

### Phase 1 Result — Placement Matrix

| Candidate / pressure | Layer decision | Primary owner surface | Phase 2 scenario? | Notes |
|---|---|---|---|---|
| `workflow-activation` | routing pressure | Workflow Activation Discovery Bridge | yes | Do not create software-delivery-local routing logic. Scenario should prove a workflow exists but was not activated. |
| `workflow-escalation` | routing pressure | Workflow Activation Discovery Bridge + `workflow/software-delivery/execution-flow.yaml` as consumer | yes | Software-delivery may declare escalation evidence, but routing mechanics stay outside this plan. |
| `explicit-root-cause-hypothesis` | evidence acquisition / diagnostic discipline | shared evidence acquisition / diagnostic reasoning; no software-delivery invariant yet | yes | Keep as a checkpoint candidate for ambiguous debugging, not a validation-family invariant. |
| `evidence-sufficiency-validation` | existing shared reasoning | `enforcement/evidence-hierarchy.md` + validation reasoning sources | no new concept | Existing Claim Scope Gate and Evidence Chain already own sufficiency; this plan should reference, not duplicate. |
| `runtime-capability-validation` | workflow validation evidence shape first | `workflow/software-delivery/test-strategy.md` + `workflow/software-delivery/validation.md` | yes | High-probability candidate. Start as workflow evidence shape; shared reasoning only if browser/platform/filesystem/container/orchestration scenarios prove common structure. |
| `authority-coupled-side-effects` | workflow evidence shape backed by Evidence Chain | `workflow/software-delivery/validation.md` + `workflow/software-delivery/artifact-gates.md` | yes | Broader than identity-coupled side effects because business truth may be owned by payment gateway, DB settlement, provider acceptance, or user-observable result. |
| `configuration-readback-validation` | workflow evidence shape backed by State Visibility Gap | `workflow/software-delivery/validation.md` + `workflow/software-delivery/artifact-gates.md` | yes | Do not create a separate shared atom unless later scenarios show structure beyond State Visibility / Evidence Chain. |
| `ownership-boundary-validation` | split working name | `task-scope-validation` -> `workflow/software-delivery/intake.md`; `ownership-awareness` -> `intake.md`, `closure.md`, `review-checklist.md` | yes, split into two scenarios | Phase 1 rejects one combined invariant for now. Task boundary and ownership boundary can diverge. |
| `operational-transaction-closure` | defer workflow placement; likely shared execution reasoning candidate | Phase 2/4 decision; possible future `intelligence/engineering/execution/` doc | yes | `closure.md` has DoD, but deploy/migration/backfill/cache rebuild transaction states need cross-operation evidence before workflow implementation. |
| `journey-level-post-deploy-validation` | existing Journey Validation reference | `test-strategy.md`, `validation.md`, `artifact-gates.md`, `review-checklist.md` | yes only for post-deploy smoke claim | Do not create a new journey taxonomy; use existing BDD-owned Journey Specification and validation-owned execution. |
| `claim_validation` | hypothesis-only | none | no | Similarity across authority/readback/closure/journey is not enough. Existing evidence hierarchy owns claim scope. |

### Phase 1 Result — Non-Goals Before Workflow Docs

- Do not add workflow gates until Phase 2 scenarios exist and fail-by-absence, unless a scenario is explicitly marked doc-only spike.
- Do not introduce `claim_validation` as a document, glossary entry, route, runtime surface, or parent taxonomy in this plan phase.
- Do not merge task scope and ownership awareness unless Phase 2 proves one model catches both failure families without ambiguity.
- Do not duplicate State Visibility Gap, Evidence Chain, Evidence Depth, Evidence Hierarchy, Journey Validation, or Workflow Activation mechanics.
- Do not hardcode incident examples such as `navigator.share`, localhost deploy config, nested repos, share count, or browser smoke tests as framework-canonical rules.

## Phase 2 — Scenario-First Validation

- [x] Add scenario: runtime capability assumed supported but unavailable fallback is missing -> expected `runtime-capability-validation` failure.
- [x] Add scenario: workflow exists but is not activated for a surface-level UI task -> expected activation pressure / routing failure routed to Workflow Activation Discovery Bridge.
- [x] Add scenario: task starts as UI but evidence reveals browser capability / side effect / journey scope and workflow does not escalate -> expected escalation pressure / routing failure.
- [x] Add scenario: patch is made without explicit root-cause hypothesis despite ambiguous evidence -> expected evidence acquisition / diagnostic checkpoint failure.
- [x] Add scenario: change set exceeds declared task scope before commit -> expected `task-scope-validation` failure.
- [x] Add scenario: change set crosses repo or owner boundary without explicit ownership awareness -> expected `ownership-awareness` failure.
- [x] Add scenario: side-effect counter increments on low-authority event while business truth is not confirmed -> expected authority evidence-shape failure.
- [x] Add scenario: deployment config input is correct but runtime readback shows stale/wrong value -> expected configuration readback evidence-shape failure.
- [x] Add scenario: deploy, migration, backfill, or cache rebuild is started/interrupted/partial without verified final state -> expected operational transaction closure placement evidence.
- [x] Add scenario: component/API smoke passes but post-deploy journey fails -> expected existing Journey Validation failure, not new journey taxonomy.
- [x] Verify scenarios fail-by-absence before workflow implementation, unless Phase 2 is explicitly marked doc-only spike.

### Phase 2 Result — Scenario Files

All Phase 2 scenarios are scenario-first inputs. They intentionally record `current_expected_result: FAIL_BY_ABSENCE` because Phase 3 workflow docs and any future executable gates have not been implemented yet.

| Scenario | Covers | Expected failure |
|---|---|---|
| `validation/scenarios/software-delivery/governance-runtime-capability-missing-fallback.yaml` | runtime capability breadth / fallback | `runtime_capability_validation_fail` |
| `validation/scenarios/software-delivery/governance-workflow-activation-pressure.yaml` | activation as routing pressure | `activation_pressure_routing_fail` |
| `validation/scenarios/software-delivery/governance-workflow-escalation-pressure.yaml` | escalation as routing pressure | `escalation_pressure_routing_fail` |
| `validation/scenarios/software-delivery/governance-root-cause-hypothesis-missing.yaml` | evidence acquisition / diagnostic discipline | `root_cause_hypothesis_missing` |
| `validation/scenarios/software-delivery/governance-task-scope-overreach.yaml` | task scope overreach | `task_scope_validation_fail` |
| `validation/scenarios/software-delivery/governance-ownership-awareness-cross-boundary.yaml` | ownership boundary awareness | `ownership_awareness_fail` |
| `validation/scenarios/software-delivery/governance-authority-side-effect-low-authority.yaml` | authority-coupled side-effect evidence | `authority_coupled_side_effect_validation_fail` |
| `validation/scenarios/software-delivery/governance-configuration-readback-stale-runtime.yaml` | configuration readback / State Visibility Gap | `configuration_readback_validation_fail` |
| `validation/scenarios/software-delivery/governance-operational-transaction-partial.yaml` | operational transaction closure placement | `operational_transaction_closure_fail` |
| `validation/scenarios/software-delivery/governance-post-deploy-journey-smoke-fail.yaml` | post-deploy journey validation via existing Journey Validation | `post_deploy_journey_validation_fail` |

## Phase 3 — Workflow Documentation Update

- [x] Update `workflow/software-delivery/validation.md` with authority/readback validation guidance if Phase 1 places it there.
- [x] Update `workflow/software-delivery/intake.md` with explicit root-cause hypothesis, task-scope validation, or ownership-boundary awareness only if Phase 1 places them there.
- [x] Update `workflow/software-delivery/test-strategy.md` with runtime capability and post-deploy journey test selection if Phase 1 places it there.
- [x] Update `workflow/software-delivery/closure.md` with change ownership and operational transaction closure if Phase 1 places them there.
- [x] Update `workflow/software-delivery/artifact-gates.md` with minimum evidence shapes for authority, configuration readback, and operational transaction closure if needed.
- [x] Update `workflow/software-delivery/execution-flow.yaml` only if Phase 2 scenarios prove an executable gate is warranted.
- [x] Update `workflow/software-delivery/README.md` and `plans/README.md` with concise status/index changes.

### Phase 3 Result — Documentation / Evidence Shape Landing

Phase 3 intentionally did not add runtime gates, `execution-flow.yaml` gates, routes, or automatic enforcement. The Phase 2 scenarios are still fail-by-absence with respect to executable enforcement; Phase 3 only gives agents and reviewers evidence shapes to use.

| Surface | Documentation update |
|---|---|
| `workflow/software-delivery/validation.md` | Added evidence shapes for runtime capability validation, authority-coupled side effects, configuration readback validation, and operational transaction closure. |
| `workflow/software-delivery/artifact-gates.md` | Added minimum artifact fields for runtime capability, authority-coupled side effect, configuration readback, and operational transaction closure. |
| `workflow/software-delivery/intake.md` | Added Task Scope / Ownership Awareness Check and a structured change-boundary review shape. |
| `workflow/software-delivery/test-strategy.md` | Added runtime capability test selection guidance for supported and absent/denied/unavailable cases. |
| `workflow/software-delivery/closure.md` | Added operational transaction close-out evidence shape and retained shared execution reasoning as the promotion path if scenarios converge. |
| `workflow/software-delivery/execution-flow.yaml` | No change; Phase 2 did not prove an executable gate is warranted yet. |

## Phase 4 — Glossary / Shared Reasoning Decision

- [ ] Decide whether candidate terms should be registered in `knowledge/glossary/ai-skill.md`.
- [ ] Keep `claim_validation` as a hypothesis unless Phase 0 + scenarios prove at least three invariants share identical validation structure.
- [ ] Decide whether `authority-coupled-side-effects`, `configuration-readback-validation`, or `operational-transaction-closure` should become shared Validation Reasoning documents instead of workflow-only guidance.
- [ ] Decide whether `operational-transaction-closure` should graduate to shared execution reasoning rather than software-delivery workflow.
- [ ] Decide whether `explicit-root-cause-hypothesis` belongs under shared evidence acquisition / validation reasoning instead of software-delivery.
- [ ] If shared reasoning is chosen, create focused documents under `intelligence/engineering/execution/validation-reasoning/` and update indexes.
- [ ] If workflow-only is chosen, explicitly state why the concept should not yet graduate.
- [ ] If activation / escalation is deferred to Workflow Activation Discovery Bridge, update this plan with a cross-link instead of duplicating implementation.

## Phase 5 — Runtime Refresh + Closure

- [ ] Run ReadLints for changed docs/scenarios.
- [ ] Run `git diff --check`.
- [ ] Run `ai-skill runtime refresh`.
- [ ] Run `ai-skill runtime validate`.
- [ ] Run `ai-skill plans tree --state active --format text`.
- [ ] Confirm `git status --short --branch` is clean after commit/push if implementation proceeds.

## 完成條件

- [x] Incident feedback is transformed into reusable governance invariants rather than one-off rules.
- [x] Activation / escalation feedback is routed to the correct plan or represented only as evidence pressure here.
- [x] Explicit root-cause hypothesis has a placement decision and is not forced into software-delivery if shared evidence acquisition is the better owner.
- [x] Task scope and ownership boundary are merged or deliberately split with evidence.
- [x] `claim_validation` remains hypothesis-only unless structural identity is proven across at least three invariants.
- [ ] `operational-transaction-closure` has a placement decision: software-delivery workflow, shared execution reasoning, or deferred.
- [x] Each invariant has an owner surface, a deferred reason, or a shared-reasoning promotion decision.
- [x] At least 3 validation scenarios cover distinct failure families before executable workflow changes.
- [x] Workflow docs do not duplicate existing State Visibility Gap, Journey Validation, or Evidence Chain guidance.
- [x] No generic runtime schema or context taxonomy is added without named consumers.
- [x] `plans/README.md` status row is updated.

## Stakeholder 同意項目

- [x] Do not add per-browser-API fallback rules.
- [x] Do not add UI fast path as a bypass around workflow activation.
- [x] Treat activation / escalation as routing pressure unless Phase 0 proves they are software-delivery invariants.
- [x] Prefer workflow escalation / complexity reassessment over checklist sprawl.
- [x] Treat explicit root-cause hypothesis as evidence acquisition discipline unless Phase 0 places it elsewhere.
- [x] Do not assume task boundary and ownership boundary are the same concept; Phase 0 must decide whether to split them.
- [x] Do not introduce `claim_validation` as a parent abstraction unless at least three invariants share identical validation structure.
- [x] Treat evidence sufficiency as the check for whether validation proof is enough, not as another smoke-test checklist.
- [x] Extract `runtime-capability-validation` instead of `browser-capability-fallback`.
- [x] Extract `authority-coupled-side-effects` instead of low-level side-effect counter semantics.
- [x] Extract `configuration-readback-validation` instead of deploy-env input checks.
- [x] Extract task-scope / ownership-boundary guidance instead of nested-repo-only Git checks.
- [x] Keep `operational-transaction-closure` as a first-class candidate.
- [x] Treat browser smoke tests as journey-level post-deploy validation when a user flow is at stake.

## Per-surface consumer 表

| Generated surface key | Named consumer(s) | Consumer 類型 |
|---|---|---|
| n/a | n/a | Draft plan only; no runtime surface, route, validator, or generated projection added yet |

## 與其他 plans 的關係

- Builds on [`archived/2026-06-08-1544-evidence-acquisition-layer`](../archived/2026-06-08-1544-evidence-acquisition-layer.md): these invariants rely on evidence acquisition/evaluation separation.
- Builds on [`archived/2026-06-10-0908-user-journey-validation-integration`](../archived/2026-06-10-0908-user-journey-validation-integration.md): post-deploy smoke should become journey-level validation when the claim is a user flow.
- Related to [`active/2026-06-09-1040-experience-validation-pipeline-evolution`](2026-06-09-1040-experience-validation-pipeline-evolution.md): this plan adds governance-invariant pressure but should not promote Typed Context Taxonomy or generic Evidence Envelope by itself.
- Related to [`active/2026-06-06-1700-workflow-activation-discovery-bridge`](2026-06-06-1700-workflow-activation-discovery-bridge.md): workflow activation / escalation failures should feed that route/routing plan rather than becoming software-delivery-local routing logic.
- Related to existing validation reasoning sources under `intelligence/engineering/execution/validation-reasoning/`; Phase 0 decides whether new shared reasoning documents are needed.
