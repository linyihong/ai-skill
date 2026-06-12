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
  - **Evidence**: F19 (Validation Scenario Governance Executor, archived 2026-06-12) added as sample #5 in [`governance-pattern-library-draft.md`](../../governance/lifecycle/governance-pattern-library-draft.md) §"Sample inventory". Genuine, not retrofitted — F19 landed independently for its own promotion (`pending_implementation`→`mechanical`), not to fill this gate.
- [x] At least one sample that DOES NOT fit perfectly — analyze why (essential signal: 6-step is correct, or 6-step has variants, or 6-step is wrong)
  - **Evidence**: now **two** non-fitting samples on **two different steps**: #4 (Runtime Index) missing *Rule* (structural-invariant condition), #5 (F19) missing *Projection* (compile-time-direct-read condition). Analysed in draft §"Per-step counts" — Projection ⇔ hot-path hypothesis. Essential signal = **"6-step has variants"** (middle of the three).
- [x] Cross-check: the 6 steps each have ≥3 samples (each step's universality verified, not just the chain's)
  - **Evidence**: per-step counts (N=5) — Observation 5, Rule 4, Registry 5, Projection 4, Executor 5, Validation 5; all ≥4.

**Gate met 2026-06-12.** Outcome is NOT "promote 6-step as-is" — the uneven (but principled) step coverage routes to Phase 1 middle branch: extract as a **4-step invariant core + 2 conditional steps**, not a rigid 6-step checklist.

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

## Phase 1 — Gate decision

- [ ] If 5+ samples found AND ≥3 samples per step AND 6-step shape consistent → proceed to Phase 2
- [x] If samples found but step coverage uneven → revise template shape (e.g. some steps are optional, document variants)
  - **Indicated branch (2026-06-12, pending owner confirmation)**: coverage is uneven in a principled way — Rule (4/5) and Projection (4/5) each conditionally optional under *distinct* falsifiable conditions (structural-invariant; compile-time-direct-read). Template should be **4-step invariant core + 2 conditional steps**, not rigid 6. Recommend confirming the two conditional rules against a 6th sample before freezing variant wording, then proceed to Phase 2.
- [ ] If samples diverge significantly → abandon promotion; keep observation in draft only

## Phase 2 — Template extraction (only if gate passes)

- [ ] Create `governance/lifecycle/governance-pattern-template.md` as **positive template**
- [ ] 6-step checklist with canonical example link per step
- [ ] Cross-link to relevant failure patterns:
  - Step 5 Executor → `enforcement/failure-patterns/validation-coverage-gap-executor-placement.md` (placement variants)
  - Step 3 Registry → `enforcement/failure-patterns/rule-without-executor.md` (executor missing variant)
- [ ] R1/R2/R3 self-governance lint rules (e.g. new mechanical rule_class must declare each of the 6 steps OR explicitly justify the missing one)
- [ ] Update `governance/lifecycle/README.md` to index this template
- [ ] Archive `governance-pattern-library-draft.md` once template lands

## Out of scope

- 不立刻寫 template（這正是本 plan 的重點 — 等樣本足夠）
- 不 retrofit 既有 plans 套入 template 格式
- 不寫 commit-msg validator 強制 future plans 套用模板（那是 Phase 3+）

## Reference / context

- Reviewer observation captured in 2026-06-08 conversation with linyihong: "你可能把目前的成功模式誤認成永遠的模式 ... 等真的出現第五次、第六次、第七次都長同樣形狀，再升級成 governance-pattern-template.md 比較穩"
- Sibling spawned plan: task_a068faa6 "Commit-Time Registry Reference Consistency" — may itself become sample #5 once landed (executor placement variant)
- Failure pattern doc landed 2026-06-08: [`enforcement/failure-patterns/validation-coverage-gap-executor-placement.md`](../../enforcement/failure-patterns/validation-coverage-gap-executor-placement.md) — already captures the "executor placement" anti-pattern; this plan captures the positive template the anti-pattern violates against
