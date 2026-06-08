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

- [ ] N ≥ 5 confirmed samples (one more genuine sample, not retrofitted)
- [ ] At least one sample that DOES NOT fit perfectly — analyze why (essential signal: 6-step is correct, or 6-step has variants, or 6-step is wrong)
- [ ] Cross-check: the 6 steps each have ≥3 samples (each step's universality verified, not just the chain's)

If the gate is met → extract to `governance/lifecycle/governance-pattern-template.md` as a positive template (sibling of failure patterns; failure patterns capture anti-patterns, this captures positive recurring shape).

If not met → keep observation in plan body / lessons file; revisit when 5th + 6th samples emerge.

## Phase 0 — Sample inventory + draft analysis

Place draft analysis at: [`governance/lifecycle/governance-pattern-library-draft.md`](../../governance/lifecycle/governance-pattern-library-draft.md)

- [ ] List the 4 current samples (column per 6-step phase, row per sample) — table form
- [ ] Identify counter-samples or partial samples on main (likely candidates: `runtime-trigger-wiring` validator, `plan-tree-hierarchy` plan, `bootstrap-contract-yaml-migration`, `commit-time registry reference consistency` (spawn chip task_a068faa6))
- [ ] For each sample, fill the 6 columns. Empty cell = step missing. Document which cells are empty for which samples.
- [ ] Check whether the pattern is sequential (always Observation → ... → Validation) or has parallel/optional branches

## Phase 1 — Gate decision

- [ ] If 5+ samples found AND ≥3 samples per step AND 6-step shape consistent → proceed to Phase 2
- [ ] If samples found but step coverage uneven → revise template shape (e.g. some steps are optional, document variants)
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
