# Dogfood Evidence — Vidoe-Test project-layer (Change Intent Lock pilot)

**Plan**: [`_plan.md`](_plan.md) Phase 4 (second evidence path)  
**Date**: 2026-06-25 (pilot design) · **observed / partial-verified** 2026-06-29  
**Repo**: `Vidoe-Test` (project-layer; **not** Ai-skill `plans/active/`)  
**Evidence class**: **partial-verified** (happy path / structure-transition only)

Complements [`01-dogfood-evidence.md`](01-dogfood-evidence.md) (**verified** — `force_exit` path). This file records **cross-repo** evidence that `preparatory_refactoring` + Change Intent Lock can constrain real H5 work without ai-skill commit-msg runtime projection.

> **Not** a success-story completion claim. Phase 1 feature work remains in progress; this evidence validates **execution contract classification** (structure → transition → feature started), not full feature delivery.

---

## Evidence maturity (this file)

```text
Observed → Partial Verified → Verified (behavior proven) → Promoted (independently auditable)
```

| Level | Meaning | This file |
|-------|---------|-----------|
| Observed | Narrative + artifacts named | yes |
| Partial Verified | Structure transition + guard observed; canonical `exit_when` mapped | **current** |
| Verified | Observable equivalence / behavior proven (Gate A) | **not yet** |
| Promoted | Independently auditable pointer + ai-skill validator/scenario wiring | **explicitly no** (observation period) |

**Upgrade gates**: Gate A → **Verified**; pointer/SHA → **Promoted** (not blocking collection or Verified).

---

## Intake routing (mapped to dual-axis)

```yaml
change_kind: feature
blocked_by_structure: true
# Rationale: ImmersivePlayerFrame is preview-gated, swipe-mutex, tab-shell z-index,
# and HLS metadata coupling — feature cannot land safely without contracts first.
execution_mode: preparatory_refactoring   # project analog; native plan uses Phase 0/1 wording
```

**Not** `replacement` / `migration` — no parity inventory; orthogonal per `intake.md`.

---

## Implementation plan (retrospective mapping)

> **Retro-map declaration**: Implementation YAML shown below is **retrospective mapping for evidence review**; it is **not** native Vidoe-Test plan schema. The authoritative project plan is prose Phase 0 / Phase 1 in `docs/plans/2026-06-25-landscape-horizontal-player-mode.md`.

Project plan: `Vidoe-Test/docs/plans/2026-06-25-landscape-horizontal-player-mode.md`

```yaml
execution_mode: preparatory_refactoring   # implied; plan labels Phase 0 / Phase 1
steps:
  - id: phase-0-contracts
    intent: structure
    behavior_change:
      allowed: false
    action: screen-mapping + player-spec §Landscape + hazard notes + BDD trace refs
    checkpoint:
      observable_equivalence:
        required: true
        # Intent (not yet proven): portrait + vertical-snap + preview-gate unchanged

  - id: phase-1-mvp
    intent: feature
    behavior_change:
      allowed: true
    action: ImmersivePlayerFrame landscapeMode + page.module.css variants
    validation: player-landscape-mode.integration.mjs + deploy smoke

  - id: phase-3-closure
    intent: feature
    behavior_change:
      allowed: true
    action: integration green + claim_registry + plan status completed
```

**Intent transition**: `structure → feature` allowed after Phase 0 artifacts exist (project file_exists + human gate). **Not** the same as `observable_equivalence_passed` until regression evidence is attached (see Gate A).

---

## Path classification

| Segment | Classification | Status (2026-06-29) |
|---------|----------------|---------------------|
| Phase 0 contracts | structure intent | **observed** — artifacts exist |
| Structure → feature transition | transition gate | **partial-verified** — guard + mapping/spec gate |
| Phase 1 implementation | feature | **in progress** — WIP; not independently auditable from Ai-skill |
| force_exit | N/A | Not triggered — contracts reduced local feature cost |

**Valid evidence path**: **Happy path (partial)** — structure phase + transition observed → feature work started. Differs from [`01-dogfood-evidence.md`](01-dogfood-evidence.md) (`force_exit` teaching case).

---

## Gate A — observable_equivalence（blocking for Verified）

Per [`execution-modes.md`](../../../workflow/software-delivery/implementation/execution-modes.md): **checkpoint ≠ observable equivalence**.

| Claim level | What we have | Status |
|-------------|--------------|--------|
| `checkpoint_exists` | mapping file + spec §Landscape + integration scaffold | **yes** (2026-06-29) |
| `checkpoint_valid` | portrait / snap / preview **regression executed and green** | **no** — not recorded in this evidence file |

**Do not conflate** project guard `file_exists` checks with `observable_equivalence_passed`.

**Verified 2026-06-29 (artifact gate only)**:

```text
screen-mapping/player-landscape-mode.md -> ok
player-spec §Landscape -> ok
```

**Pending for Verified upgrade**: command + date + pass/fail for portrait-player regression / integration owned by Phase 0 checkpoint intent.

**Failure lens**: aligns with §8 `fake equivalence` if we claimed Valid without regression proof.

---

## Gate B — exit_when（recorded at Partial Verified)

Canonical vocabulary from [`execution-modes.md`](../../../workflow/software-delivery/implementation/execution-modes.md) §4. Satisfied for **partial-verified**; does not alone satisfy **Verified** (see Gate A).

Canonical vocabulary from [`execution-modes.md`](../../../workflow/software-delivery/implementation/execution-modes.md) §4:

| `exit_when` candidate | Applies? | Evidence |
|-----------------------|----------|----------|
| `target_change_becomes_local` | **partial** | Phase 0 contracts + guard unblocked **starting** ImmersivePlayerFrame work; feature not closed |
| `target_test_becomes_expressible` | **yes (primary)** | `player-landscape-mode.integration.mjs` + fixture scaffold exist — landscape acceptance can be expressed before/alongside implementation |
| `new_abstraction_created` | no | No new seam/abstraction retained for feature to consume |

**Recorded exit_when (partial-verified)**: **`target_test_becomes_expressible`** — integration scaffold + fixture make landscape behavior testable; structure phase exited into feature intent on that basis **plus** artifact gate, **not** on equivalence proof (Gate A gap).

`force_exit_when`: **none** (happy partial path).

---

## Mechanical enforcement (project overlay)

| Artifact | Path |
|----------|------|
| Registry | `Vidoe-Test/docs/plans/active-plan-phase-guards.json` |
| Commit hook | `Vidoe-Test/.ai-skill/project/rules/check-plan-phase-before-commit.py` |
| Cursor wiring | `Vidoe-Test/.cursor/hooks.json` → `beforeShellExecution` on `git commit` |
| BDD | `Vidoe-Test/tests/bdd/plan-phase-guard.test.mjs` |
| Opt-out | `PLAN_PHASE_GUARD_SKIP=1` |

**Guard rule**: staging `server_doc/.../ImmersivePlayerFrame.tsx` or player `page.module.css` **denies commit** until mapping + spec §Landscape exist.

Simulated bad commit (staged markdown with `Vidoe-Test`) blocked by sibling `server_doc` sanitization guard in same hook script.

**Ai-skill posture**: project overlay only; aligns with observation-period **no ai-skill commit-msg validator** decision.

---

## Phase 0 artifacts (structure phase)

| Artifact | Path |
|----------|------|
| Main plan | `docs/plans/2026-06-25-landscape-horizontal-player-mode.md` |
| Screen mapping | `docs/frontend-contracts/screen-mapping/player-landscape-mode.md` |
| Visual spec §4 | `docs/frontend-contracts/ui-style-reference/specs/player-spec.md` |
| Integration scaffold | `tests/integration/player-landscape-mode.integration.mjs` |
| Fixture | `tests/integration/player-fixtures.mjs` (`PLAYER_HORIZONTAL_EPISODE_ID`) |

---

## External pointer（blocking for Promoted only)

| Field | Value |
|-------|-------|
| Repo | `Vidoe-Test` (external; not vendored in Ai-skill workspace) |
| Commit SHA | *TBD — attach when Phase 0 artifacts land on remote* |
| Verifier | *TBD — name + date + command for guard simulation* |
| Reproducibility | **not independently reproducible from Ai-skill repo alone** (disclosed) |

**Not** blocking Phase 4 collected. **Not** blocking **Verified** (behavior can be proven without full pointer). **Blocks Promoted** until independently auditable.

---

## Failure-mode review lens

| failure | Project observation |
|---------|---------------------|
| intent oscillation | Plan orders Phase 0 → 1 → 3; guard blocks implementation-first commits |
| structure inflation | Phase 0 scoped to mapping + spec regions, not broad player rewrite |
| fake equivalence | **Gate A open** — artifact gate ≠ portrait regression proof |
| abstraction orphan | N/A — no unused parser/seam (contrast `01` prep-02) |
| compatibility collapse | Native plan uses Phase 0/1; no forced `execution_mode` on legacy plans |
| replacement parity misuse | Plan explicitly frontend-only; no replacement inventory |
| illegal transition | Not wired — would need native `steps[]` YAML + planvalidate advisory |

---

## What this evidence supports in `_plan.md`

| Plan section | Supported? | Notes |
|--------------|------------|-------|
| Phase 4 dual path collected | **yes** | Happy **partial** complements `01` force_exit **verified** |
| Phase 3 routing cross-links | **partial** | Concept via project plan wording, not ai-skill `loading_surfaces` |
| Phase 5 glossary | **no** | Project uses Phase 0/1 labels |
| Validator hook / enforcement | **explicitly no** | Project overlay ≠ ai-skill promotion |

---

## Recorded fields (Phase 4 checklist)

| Field | Value |
|-------|-------|
| change_kind | feature |
| execution_mode | preparatory_refactoring (project analog) |
| evidence_class | partial-verified |
| intent sequence | structure (Phase 0) → feature (Phase 1+, in progress) |
| transitions | structure → feature after artifact gate; equivalence **not** verified |
| exit_when | `target_test_becomes_expressible` (primary); `target_change_becomes_local` (partial) |
| force_exit trigger | none |
| checkpoint | `checkpoint_exists` yes; `checkpoint_valid` **pending** |
| blocking | project Cursor hook; ai-skill planvalidate `Blocking=false` |

---

## Appendix — related pilots (not equally mapped)

Same pattern **family** only; not submitted as Phase 4 evidence with equal rigor.

| Pilot | Maps to | Notes |
|-------|---------|-------|
| Navigation Phase B | structure before feature closure | completed; no retro-map in this file |
| Design Contract B0 | structure stub before enforcement | completed; pointer in ai-skill evidence-candidates |
| Plan-first ordering README | contract-before-implementation culture | `docs/plans/README.md` §Plan Phase 0 guard |

**Out of scope**: `server_doc` sanitization (documentation boundary — different failure family).
