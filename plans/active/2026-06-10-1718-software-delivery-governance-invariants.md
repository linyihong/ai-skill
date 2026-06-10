---
id: 2026-06-10-1718-software-delivery-governance-invariants
plan_kind: main
status: draft
owner: linyihong
created: 2026-06-10
priority: P1
required_for_completion: false
---

# Software Delivery Governance Invariants

**Status**: `draft`
Owner: framework maintainer (linyihong)
**建立日期**：2026-06-10
**Priority**：**P1**（software-delivery workflow hardening；把 incident feedback 提煉成 governance invariants）

## Why this plan exists

Recent software-delivery feedback listed concrete incidents such as browser API fallback, share counters, localhost/base URL deployment drift, nested repo dirty states, interrupted deploys, and browser smoke tests. These are real issues, but most of them are symptoms.

The higher-level gap is that `workflow/software-delivery` needs reusable governance invariants for:

- activating the right workflow before implementation
- escalating workflow scope when task complexity changes
- making root-cause hypotheses explicit before patching
- assumptions about runtime capabilities
- business-truth ownership for side effects
- deployed-state readback instead of config-input validation
- change ownership before committing across repos or ownership groups
- operational transactions that can be partial, interrupted, or unverified
- post-deploy validation at user-journey level

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

These are not separate categories of bugs. They are governance invariant gaps around activation, escalation, hypothesis discipline, authority, readback, ownership, and operational closure.

### Decision

Draft reusable invariants for `workflow/software-delivery`, with `authority-coupled-side-effects`, `configuration-readback-validation`, and `operational-transaction-closure` treated as the strongest validation-family candidates, and workflow activation / escalation treated as routing-family pressure that must cross-link rather than duplicate the Workflow Activation Engine.

| Invariant | Prevents | Core question |
|---|---|---|
| `workflow-activation` | workflow exists but is not loaded | Which workflow slices must activate before implementation? |
| `workflow-escalation` | task starts narrow but needs broader governance | When does a task escalate from UI/CSS to runtime capability, side effect, deployment, or journey validation? |
| `explicit-root-cause-hypothesis` | silent assumptions and patch-first debugging | What root cause is claimed, what evidence supports it, what patch follows, and how will it be validated? |
| `runtime-capability-validation` | browser / platform supported-assumption bugs | Does the runtime actually provide the capability, and is fallback validated? |
| `authority-coupled-side-effects` | fake success metrics / premature counters | Which event owns the business truth? |
| `configuration-readback-validation` | validated config input but wrong deployed state | Did we read back actual deployed state? |
| `change-boundary-detection` / `change-ownership-validation` | mixed-task commits across repos or owner groups | What belongs to this task, and who owns each changed boundary? |
| `operational-transaction-closure` | interrupted / partial deployment ambiguity | Did the operational transaction finish and get verified? |
| `journey-level-post-deploy-validation` | component pass but user flow fail | Did the critical journey pass after deploy? |

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
- D. Treat nested repo dirty state as Git-specific: reject. The reusable concept is `change-ownership-validation`, which also applies to monorepos and generated artifacts.
- E. Add browser smoke tests as a separate category: reject. Most post-deploy smoke tests should map to `journey-level-post-deploy-validation`.
- F. Add UI fast path: reject as unsafe. Many incidents looked like small UI work but required runtime capability, side-effect authority, route/history, or journey validation. Prefer `workflow-escalation` and `complexity-reassessment`.
- G. Add change brief / checklist layer: reject as primary solution. Briefs can be useful output, but the deeper issue is activating the right workflow at the right time.
- H. Extract governance invariants, then decide their workflow placement, routing relationship, and validation evidence: accept.

### Why Not an ADR Yet

This is not ready for ADR promotion. The invariant family is plausible, but needs:

- mapping to current `workflow/software-delivery` surfaces without duplicating existing State Visibility / Journey Validation content
- mapping activation / escalation concerns to the Workflow Activation Engine without duplicating `workflow-activation-discovery-bridge`
- scenario evidence showing the invariants catch more than the original incident examples
- ownership decisions: workflow vs shared validation reasoning vs future intelligence atom
- proof that the plan does not create a catch-all taxonomy or over-broad runtime schema

### ADR Promotion Criteria（completed 時驗證）

- [ ] At least 3 invariants land with clear workflow placement and no duplicated source-of-truth.
- [ ] At least 3 validation scenarios prove the invariants catch distinct failure families.
- [ ] Workflow activation / escalation concerns are cross-linked to the Workflow Activation Engine or explicitly deferred.
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

#### 風險

- Over-expanding the invariants could create another catch-all taxonomy.
- Authority maps can become too abstract if not tied to concrete evidence classes.
- Workflow activation concerns could duplicate the existing activation engine unless this plan treats them as cross-links and evidence pressure.
- Operational transaction closure can overlap with CI/CD tooling unless scoped to workflow evidence and completion claims.

Glossary Impact: yes — candidate terms `workflow_activation`, `workflow_escalation`, `explicit_root_cause_hypothesis`, `change_boundary_detection`, `evidence_sufficiency_validation`, `runtime_capability_validation`, `authority_coupled_side_effects`, `configuration_readback_validation`, `change_ownership_validation`, `operational_transaction_closure`, `journey_level_post_deploy_validation`, `business_truth`, `observable_proxy`; Phase 2 decides whether to register them in `knowledge/glossary/ai-skill.md` or keep them plan-local.

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
- [ ] Which activation/escalation concerns belong in this plan, and which should be deferred to Workflow Activation Discovery Bridge?
- [ ] Should `explicit-root-cause-hypothesis` be an intake checkpoint, a debug/review artifact shape, or a general execution invariant?
- [ ] Should `change-boundary-detection` be distinct from `change-ownership-validation`, or should one be the pre-commit form of the other?
- [ ] What evidence is sufficient to decide workflow escalation from UI -> runtime capability -> side effect -> journey validation?
- [ ] Should `authority-coupled-side-effects` become a workflow gate, an artifact evidence shape, or an intelligence atom first?
- [ ] What is the minimum authority ladder for side effects without hardcoding domain-specific business events?
- [ ] How should `configuration-readback-validation` relate to State Visibility Gap without duplicating it?
- [ ] Does `operational-transaction-closure` belong in `closure.md`, `validation.md`, `execution-flow.yaml`, or a separate operational slice?
- [ ] What validation scenarios are needed before any of these become executable gates?
- [ ] Should runtime capability validation be limited to user-visible surfaces first, or include backend/platform runtime capabilities too?

## Phase 0 — Open Questions Check + Architecture Compatibility Preflight

### Phase 0.0 — Open Questions 核對（公版，必填）

逐條核對本 plan §Open Questions，標記處置並回寫：

- [ ] 已讀本 plan §Open Questions 全部條目
- [ ] 對每條標記 `resolved`（附 Phase 0 證據）/ `still-open` / `deferred`（附原因）
- [ ] `resolved` 的條目已同步勾選 / 附註於 §Open Questions
- [ ] 若盤點新發現問題，已加入 §Open Questions

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| invariant placement | pending | Phase 0 must read current workflow slices and validation reasoning sources |
| activation / escalation placement | pending | Must compare with Workflow Activation Discovery Bridge and avoid duplicate routing design |
| explicit root-cause hypothesis | pending | Needs intake / validation / review placement decision |
| change boundary vs ownership | pending | Needs commit/closure boundary decision |
| authority side effects shape | pending | Needs evidence-chain and failure-authority comparison |
| configuration readback boundary | pending | Needs State Visibility Gap comparison |
| operational transaction placement | pending | Needs closure / execution-flow comparison |
| validation scenarios | pending | Must be designed before implementation |

### Phase 0.1 — Preflight

- [ ] Read current `workflow/software-delivery/README.md`, `validation.md`, `test-strategy.md`, `closure.md`, `artifact-gates.md`, `execution-flow.yaml`, and `review-checklist.md`.
- [ ] Read `plans/active/2026-06-06-1700-workflow-activation-discovery-bridge.md` and decide how activation / escalation pressure should cross-link.
- [ ] Read shared validation reasoning sources: `state-visibility-gap.md`, `evidence-model.md`, `evidence-chain-validation.md`, `evidence-depth.md`, and any Failure Authority source if present.
- [ ] Confirm no duplicate vocabulary already exists in glossary or current workflow.
- [ ] Decide whether this is a sibling plan to Experience Validation Pipeline or should cross-link as a downstream hardening plan.
- [ ] Confirm no new runtime projection is needed for Phase 0.
- [ ] Record not-applicable dependencies explicitly if a source does not exist.

## Phase 1 — Invariant Placement Draft

- [ ] Create a placement matrix for the six invariants:
  - workflow-activation
  - workflow-escalation
  - explicit-root-cause-hypothesis
  - change-boundary-detection
  - evidence-sufficiency-validation
  - runtime-capability-validation
  - authority-coupled-side-effects
  - configuration-readback-validation
  - change-ownership-validation
  - operational-transaction-closure
  - journey-level-post-deploy-validation
- [ ] For each invariant, decide primary owner surface: `intake.md`, `validation.md`, `test-strategy.md`, `closure.md`, `artifact-gates.md`, `execution-flow.yaml`, Workflow Activation plan, shared validation reasoning, or defer.
- [ ] Define non-goals so incident-specific examples do not become hardcoded framework rules.
- [ ] Update this plan with the placement decision before editing workflow docs.

## Phase 2 — Scenario-First Validation

- [ ] Add scenario: runtime capability assumed supported but unavailable fallback is missing -> expected validation failure.
- [ ] Add scenario: workflow exists but is not activated for a surface-level UI task -> expected activation failure.
- [ ] Add scenario: task starts as UI but evidence reveals browser capability / side effect / journey scope and workflow does not escalate -> expected escalation failure.
- [ ] Add scenario: patch is made without explicit root-cause hypothesis despite ambiguous evidence -> expected hypothesis failure.
- [ ] Add scenario: change set includes unrelated owner boundary before commit -> expected change-boundary failure.
- [ ] Add scenario: side-effect counter increments on low-authority event while business truth is not confirmed -> expected authority failure.
- [ ] Add scenario: deployment config input is correct but runtime readback shows stale/wrong value -> expected configuration readback failure.
- [ ] Add scenario: deploy transaction is started/interrupted/partial without verified final state -> expected operational closure failure.
- [ ] Add scenario: component/API smoke passes but post-deploy journey fails -> expected journey-level post-deploy validation failure.
- [ ] Verify scenarios fail-by-absence before workflow implementation, unless Phase 2 is explicitly marked doc-only spike.

## Phase 3 — Workflow Documentation Update

- [ ] Update `workflow/software-delivery/validation.md` with authority/readback validation guidance if Phase 1 places it there.
- [ ] Update `workflow/software-delivery/intake.md` with explicit root-cause hypothesis or change-boundary detection if Phase 1 places them there.
- [ ] Update `workflow/software-delivery/test-strategy.md` with runtime capability and post-deploy journey test selection if Phase 1 places it there.
- [ ] Update `workflow/software-delivery/closure.md` with change ownership and operational transaction closure if Phase 1 places them there.
- [ ] Update `workflow/software-delivery/artifact-gates.md` with minimum evidence shapes for authority, configuration readback, and operational transaction closure if needed.
- [ ] Update `workflow/software-delivery/execution-flow.yaml` only if Phase 2 scenarios prove an executable gate is warranted.
- [ ] Update `workflow/software-delivery/README.md` and `plans/README.md` with concise status/index changes.

## Phase 4 — Glossary / Shared Reasoning Decision

- [ ] Decide whether candidate terms should be registered in `knowledge/glossary/ai-skill.md`.
- [ ] Decide whether `authority-coupled-side-effects`, `configuration-readback-validation`, or `operational-transaction-closure` should become shared Validation Reasoning documents instead of workflow-only guidance.
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

- [ ] Incident feedback is transformed into reusable governance invariants rather than one-off rules.
- [ ] Activation / escalation feedback is routed to the correct plan or represented only as evidence pressure here.
- [ ] Explicit root-cause hypothesis and change-boundary detection have placement decisions.
- [ ] Each invariant has an owner surface, a deferred reason, or a shared-reasoning promotion decision.
- [ ] At least 3 validation scenarios cover distinct failure families before executable workflow changes.
- [ ] Workflow docs do not duplicate existing State Visibility Gap, Journey Validation, or Evidence Chain guidance.
- [ ] No generic runtime schema or context taxonomy is added without named consumers.
- [ ] `plans/README.md` status row is updated.

## Stakeholder 同意項目

- [ ] Do not add per-browser-API fallback rules.
- [ ] Do not add UI fast path as a bypass around workflow activation.
- [ ] Prefer workflow escalation / complexity reassessment over checklist sprawl.
- [ ] Add explicit root-cause hypothesis as a candidate invariant before patching ambiguous defects.
- [ ] Treat change-boundary detection as higher-level than nested-repo checklist.
- [ ] Treat evidence sufficiency as the check for whether validation proof is enough, not as another smoke-test checklist.
- [ ] Extract `runtime-capability-validation` instead of `browser-capability-fallback`.
- [ ] Extract `authority-coupled-side-effects` instead of low-level side-effect counter semantics.
- [ ] Extract `configuration-readback-validation` instead of deploy-env input checks.
- [ ] Extract `change-ownership-validation` instead of nested-repo-only Git checks.
- [ ] Keep `operational-transaction-closure` as a first-class candidate.
- [ ] Treat browser smoke tests as journey-level post-deploy validation when a user flow is at stake.

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
