---
id: 2026-06-08-2100-governance-pattern-library-extraction
plan_kind: main
status: in-progress
owner: linyihong
created: 2026-06-08
priority: P2
required_for_completion: false
---

# Governance Pattern Library Extraction

**Status**: `in-progress`
Owner: framework maintainer (linyihong)
**建立日期**：2026-06-08
**Priority**：**P2**（observation-stage；尚未通過 acceptance gate，不阻擋其他工作）

## Why this plan exists, not an immediate write-up

在 sanitization plan Phase 1 review 對話（2026-06-08）中，觀察到近 6 週 land 的 4 個 governance subsystem 似乎都長同樣形狀：

```
Observation → Rule → Registry → Projection → Executor → Validation
```

樣本：

1. **Workflow Activation Engine** (parent plan 2026-05-31-1900, archived) — detector miss observation → routing-registry.yaml rule → enforcement-registry binding → runtime.db routes projection → detector.go + PreToolUse gate → regression scenarios
2. **Discovery Bridge** (plan 2026-06-06-1700, Phase A landed) — detector-miss-no-fallback observation → discovery-bridge.yaml rule → capability_discovery rule_class → runtime.discovery.config projection → discovery.go + advisory injector → workflow-discovery-bridge-light-v1.yaml scenario
3. **Sanitization Mechanical Enforcement** (plan 2026-06-06-1800, in design) — leak-on-canonical-write observation → sanitization-patterns.yaml + topology → enforcement-registry rule_class → derived_match_tokens projection → sanitization_scan.go → metadata-derived-fail/pass scenarios
4. **Runtime Index Freshness** (commit c5874a8, landed) — source-tree checksum drift observation → (implicit) rule → runtime-index.sqlite registry → sources table projection → nativeRuntimeIndexChecksumsCheck + commit-msg validator → runtime-index-freshness tests

The 6-step shape fits all 4. **But 4 samples is not enough to confirm a framework law.** Two failure modes if we promote now:

- The shape might be transient — next governance need could be `Observation → Rule → Executor → Validation` (no registry / no projection), and the template would over-constrain
- Some samples may have been forced into the shape because the author was reading the others — confirmation bias, not pattern discovery

## Acceptance gate (do NOT extract until met)

- [x] N ≥ 5 confirmed samples (one more genuine sample, not retrofitted)
  - **Evidence**: now **N=7** in [`governance-pattern-library-draft.md`](../../governance/lifecycle/governance-pattern-library-draft.md) §"Sample inventory". #5 = F19 (Validation Scenario Governance Executor). #6 = Plan-Tree Hierarchy Governance. #7 = Runtime Trigger Wiring (`runtime_trigger_wiring`), added 2026-06-12 as a deliberate predicate-1 challenge. None retrofitted.
- [x] At least one sample that DOES NOT fit perfectly — analyze why (essential signal: 6-step is correct, or 6-step has variants, or 6-step is wrong)
  - **Evidence**: **four** non-fitting samples: #4 (Runtime Index) missing *Rule* (structural invariant); #5/#6/#7 missing *Projection* (executor consumes the authoritative source directly). Essential signal = **"6-step is not an invariant — it is a 4-step invariant core + 2 conditional stages"** (variants branch).
- [x] Cross-check: the 6 steps each have ≥3 samples (each step's universality verified, not just the chain's)
  - **Evidence**: per-step counts (N=7) — Observation 7, Rule 6, Registry 7, Projection 4, Executor 7, Validation 7; all ≥4.

**Gate met 2026-06-12 — and the result is stronger than "gate passed".** Invariant core = Observation / Registry / Executor / Validation (7/7). The two conditionals are now **asymmetric in attestation**:
- **Projection-optional — solid (3 confirmations #5/#6/#7)**: optional for *direct-consumption* executors.
- **Rule-optional — rare carve-out (1 clean case #4)**: optional only for *pure structural (in)equality* invariants. The attempt to find a 2nd case via #7 (`runtime_trigger_wiring`) **failed informatively** — its reference-resolvability invariant looks structural but carries a policy Rule (`system-upgrade-governance.yaml §define_runtime_trigger_flow` enumerates valid wiring forms). Lesson: most "structural-looking" invariants smuggle in a policy definition of the valid relation, so Rule is **near-universal**.

The promotable artifact is a **4-step invariant core + Projection (conditional) + Rule (near-universal, narrow structural exemption)**.

If the gate is met → extract to `governance/lifecycle/governance-pattern-template.md` as a positive template (sibling of failure patterns; failure patterns capture anti-patterns, this captures positive recurring shape).

If not met → keep observation in plan body / lessons file; revisit when 5th + 6th samples emerge.

## Phase 0 — Sample inventory + draft analysis

Place draft analysis at: [`governance/lifecycle/governance-pattern-library-draft.md`](../../governance/lifecycle/governance-pattern-library-draft.md)

- [x] List the 4 current samples (column per 6-step phase, row per sample) — table form
  - **Evidence**: draft §"Sample inventory" — 4-sample table filled with Observation / Rule / Registry / Projection / Executor / Validation cells per sample
- [x] Identify counter-samples or partial samples on main (likely candidates: `runtime-trigger-wiring` validator, `plan-tree-hierarchy` plan, `bootstrap-contract-yaml-migration`, `commit-time registry reference consistency` (spawn chip task_a068faa6))
  - **Evidence**: draft §"Counter-sample candidates (to inventory next)" — 4 candidates listed with the analysis question each one should answer. Analysis itself deferred to next Phase 0 pass when a 5th genuine sample lands (likely via Commit-Time Registry Reference Consistency chip).
- [x] For each sample, fill the 6 columns. Empty cell = step missing. Document which cells are empty for which samples.
  - **Evidence**: draft §"Per-step counts (N=4)" — sample #4 (Runtime Index Freshness) explicitly flagged as having `(implicit)` Rule cell. First soft signal that *Rule* step may be conditionally optional for structural invariants.
- [x] Check whether the pattern is sequential (always Observation → ... → Validation) or has parallel/optional branches
  - **Evidence**: draft §"Sequential vs. branching shape" — answered with three nuances: (1) Rule is conditionally optional (sample #4), (2) Executor has within-step branching (samples #2, #3 have multi-entry executors), (3) no feedback loop observed yet but Discovery Bridge Phase D may produce one. Working interpretation: "mostly sequential with a single optional step + within-step branching at Executor".

**Phase 0 additional output (not in original plan, captured for completeness)**:

- [x] **Parallel Reference Integrity family observation** captured in same draft as separate section
  - 4 Reference Integrity samples inventoried (sanitization rule↔executor, workflow detector↔fallback, runtime-index content↔checksum, registry path↔pointer)
  - Distinguished from 6-step family (build-time vs maintenance-time; rule promotion vs surface consistency)
  - Independent acceptance gate (also at N≥5; same falsification discipline)
  - 3 pre-emptive falsification candidates listed (conversation goal ledger, cognitive mode reporting, discovery advisory output) — each tests whether Reference Integrity is universal or narrower than 6-step
  - This is **not** a scope expansion of this plan; it is a sibling observation sharing the same incubator surface. If either family passes its gate, they will be siblings in `governance/lifecycle/governance-pattern-library/` folder, not parent/child.

- [x] **Parallel Failure Authority family observation** captured in same draft as separate section (added 2026-06-10, from sanitization Phase 1D review)
  - Invariant: only compile-authoritative sources (shared-layer / tracked / runtime-index `sources` row) may block runtime compilation; non-authoritative sources (`shared_layer:false`, `owner:project-local`, untracked) may warn only
  - Names the conflation being fixed: `Metadata Presence ⇒ Compile Authority` (file existence + invalidity silently becomes pipeline halt regardless of source class)
  - 3 samples inventoried (Runtime Index source-row scope ✅, Workflow gate fail-open ✅, Sanitization Phase 1D `.agent-goals/` hard-fail ❌ — the live violation = sanitization plan Finding A)
  - Distinguished from the other two families: this is **enforcement-time / blast-radius authority**, not build-time onboarding (6-step) nor maintenance-time surface agreement (Reference Integrity)
  - Rides on topology v2 path-classification — promoting it protects the governance value of the owner/purpose/shared_layer work
  - Independent acceptance gate (N≥5; already has 1 analysed non-fitting sample); 3 falsification candidates listed
  - **Resolution binding**: sanitization Finding A is deferred to a Failure-Authority-governed fix at Phase 4, NOT a one-off `.agent-goals/` special-case in the sanitization scanner
  - **Conceptual core = Standing** (2026-06-10): the family's deepest framing is legal *standing* — `Validity ≠ Authority`. Many things can be wrong; not every wrong thing has the right to halt the process. Finding A proves it (100%-invalid metadata ≠ standing to halt compile). Recorded as the most valuable distinction of this family.
  - **Authority Classification Contract layer inserted** (2026-06-10): chain is `invariant → Classification Contract → Classifier → Executor`. The contract is **subject-based** (`AuthoritySubject{kind: discovery-provider | runtime-index-row | metadata-file | generated-surface}`, `ClassifyFailureAuthority(subject)`), NOT path-based (`isCompileAuthoritative(path)` would weld authority:=path and break for non-file subjects). Contract specified docs-first, no Go, so the first classifier implementation conforms to a shared definition of standing reusable by ≥5 executors. Build order: family ✅ → contract ✅ → classifier impl ✅ → Finding A executor ⏭.
  - **Classifier landed** (2026-06-10): `scripts/ai-skill-cli/internal/app/failure_authority.go` — `ClassifyFailureAuthority(AuthoritySubject)` with tri-state `SharedLayer`, demotion-wins precedence, and unknown-kind-must-earn-standing default. Unit-tested across metadata-file / runtime-index-row / discovery-provider / unknown kinds (proving multi-kind, not path-welded). No caller wired yet; Finding A wiring is the next step and contributes the family's 4th sample.

## Phase 1 — Gate decision（completed 2026-06-12）

**Status**: `completed`. The gate's decision tree resolves deterministically to the middle branch, and the shape revision is done, so Phase 1 is a closed decision — what remains is Phase 2 wording refinement, not a branch choice.

- [ ] If 5+ samples found AND ≥3 samples per step AND 6-step shape consistent → proceed to Phase 2
- [x] If samples found but step coverage uneven → revise template shape (e.g. some steps are optional, document variants)
  - **Decision (2026-06-12)**: coverage is uneven but **principled** — Rule (4/5) optional for *policy-derived* governance vs structural invariants; Projection (4/5) optional when the executor *consumes the authoritative source directly* (predicate stated at contract level, not implementation). Revised shape = **4-step invariant core (Observation → Registry → Executor → Validation) + 2 conditional steps (Rule, Projection)**. This closes Phase 1.
- [ ] If samples diverge significantly → abandon promotion; keep observation in draft only

**Phase 2 precondition — discharged (2026-06-12)** across samples #6 and #7:
- **Projection-optional — confirmed three times (#5/#6/#7).** Two executor kinds (compile-time lint + commit-msg validator), three source shapes (registry yaml, plan frontmatter, routing-registry + tree walk). Solid.
- **Rule-optional — tested and reframed, NOT independently re-confirmed.** #6's structural sub-invariant (`validatePlanTreeUniqueID`) rides on a policy rule_class (rule_class granularity). #7 (`runtime_trigger_wiring`) was the deliberate predicate-1 candidate and **did not close it**: its invariant is policy-derived (`§define_runtime_trigger_flow` enumerates valid wiring forms), not structural. Conclusion: **Rule is near-universal; the structural exemption is a rare carve-out, attested by exactly #4.**
- **Validation sub-forms** confirmed additive: #7 has both Go test + scenario yaml.

**Predicate-1 status**: closing it with a *second* clean structural-no-rule case is now a low-priority open item — the practical finding is that such cases are rare (reference-resolvability-style invariants keep proving policy-backed). For Phase 2 the template should present Rule as **near-universal with a criterion-gated structural exemption** (criterion: pure structural (in)equality, no human-authored relation — checksum / raw uniqueness), which #4 already justifies. Not gate-blocking.

**Cross-family side-effect**: #7 is also Reference Integrity sample #5; that family now has N=5 but its gate is still **not** met (no non-fitting sample analysed). See draft §Reference Integrity.

## Phase 2 — Template extraction (only if gate passes)

- [ ] Create `governance/lifecycle/governance-pattern-template.md` as **positive template**
- [ ] 6-step checklist with canonical example link per step
- [ ] Cross-link to relevant failure patterns:
  - Step 5 Executor → `enforcement/failure-patterns/validation-coverage-gap-executor-placement.md` (placement variants)
  - Step 3 Registry → `enforcement/failure-patterns/rule-without-executor.md` (executor missing variant)
- [ ] R1/R2/R3 self-governance lint rules (e.g. new mechanical rule_class must declare each of the 6 steps OR explicitly justify the missing one)
- [ ] Update `governance/lifecycle/README.md` to index this template
- [ ] Archive `governance-pattern-library-draft.md` once template lands

## Evidence Rule

> Machine-readable evidence-rule（schema `evidence-rule-v1`），索引於
> [`governance/evidence-candidates/evidence-rules/governance-pattern.pointer.yaml`](../../governance/evidence-candidates/evidence-rules/governance-pattern.pointer.yaml)。
> **Phase 1A Step 2（consumer attach）**：本 section 成立 = consumer hook 建立；criterion 內容是
> **Step 3（criteria authoring）**，下方刻意留 placeholder。rule 定義 owner = 本 plan。acceptance-gate
> 形狀候選 `sample_count >= 5 + falsification >= 1`（notify 達 gate = 「可 review」≠「已 promotion」），
> 屬 acceptance-gate，不在 evidence_rule。設計來源見
> [`evidence-candidate-system`](2026-06-16-1131-evidence-candidate-system.md)。

```yaml
evidence_rule:
  collect: true
  match:
    artifact_types: []   # Step 3
    criteria: []         # Step 3（候選方向：6-step / Reference Integrity / Failure Authority sample inventory）
  exclusions: []         # Step 3
```

## Out of scope

- 不立刻寫 template（這正是本 plan 的重點 — 等樣本足夠）
- 不 retrofit 既有 plans 套入 template 格式
- 不寫 commit-msg validator 強制 future plans 套用模板（那是 Phase 3+）

## Reference / context

- Reviewer observation captured in 2026-06-08 conversation with linyihong: "你可能把目前的成功模式誤認成永遠的模式 ... 等真的出現第五次、第六次、第七次都長同樣形狀，再升級成 governance-pattern-template.md 比較穩"
- Sibling spawned plan: task_a068faa6 "Commit-Time Registry Reference Consistency" — may itself become sample #5 once landed (executor placement variant)
- Failure pattern doc landed 2026-06-08: [`enforcement/failure-patterns/validation-coverage-gap-executor-placement.md`](../../enforcement/failure-patterns/validation-coverage-gap-executor-placement.md) — already captures the "executor placement" anti-pattern; this plan captures the positive template the anti-pattern violates against
