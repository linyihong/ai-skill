# Dogfood Evidence — Vidoe-Test project-layer (Change Intent Lock pilot)

**Plan**: [`_plan.md`](_plan.md) Phase 4 (second evidence path)  
**Date**: 2026-06-25 (pilot design) · **verified** 2026-06-29  
**Repo**: `Vidoe-Test` (project-layer; **not** Ai-skill `plans/active/`)

Complements [`01-dogfood-evidence.md`](01-dogfood-evidence.md) (Ai-skill planvalidate `force_exit` path). This file records **cross-repo** evidence that `preparatory_refactoring` + Change Intent Lock can constrain real H5 work without commit-msg runtime projection.

---

## Intake routing (mapped to dual-axis)

```yaml
change_kind: feature
blocked_by_structure: true
# Rationale: ImmersivePlayerFrame is preview-gated, swipe-mutex, tab-shell z-index,
# and HLS metadata coupling — feature cannot land safely without contracts first.
execution_mode: preparatory_refactoring   # project analog; plan uses Phase 0/1 wording
```

**Not** `replacement` / `migration` — no parity inventory; orthogonal per `intake.md`.

---

## Implementation plan (as designed + partially executed)

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
        evidence: existing portrait player + vertical-snap + preview-gate integrations unchanged

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

**Intent transition**: `structure → feature` only after Phase 0 artifacts exist (mechanical + human gate).

---

## Path classification

| Segment | Classification | Status (2026-06-29 verify) |
|---------|----------------|----------------------------|
| Phase 0 contracts | structure + observable equivalence intent | **done** — see artifacts below |
| Phase 1 implementation | feature | **in progress** — WIP in Vidoe-Test root working tree (not yet server_doc commit) |
| force_exit | N/A | Not triggered — contracts unblocked local feature work |

**Valid evidence path**: **Happy path (partial)** — structure checkpoint met → feature work started. Differs from `01-dogfood` force_exit teaching case.

---

## Mechanical enforcement (project overlay)

| Artifact | Path |
|----------|------|
| Registry | `Vidoe-Test/docs/plans/active-plan-phase-guards.json` |
| Commit hook | `Vidoe-Test/.ai-skill/project/rules/check-plan-phase-before-commit.py` |
| Cursor wiring | `Vidoe-Test/.cursor/hooks.json` → `beforeShellExecution` on `git commit` |
| BDD | `Vidoe-Test/tests/bdd/plan-phase-guard.test.mjs` |
| Opt-out | `PLAN_PHASE_GUARD_SKIP=1` |

**Guard rule**: staging `server_doc/.../ImmersivePlayerFrame.tsx` or player `page.module.css` **denies commit** until:

- `docs/frontend-contracts/screen-mapping/player-landscape-mode.md` exists
- `docs/frontend-contracts/ui-style-reference/specs/player-spec.md` contains `landscape`

**Verified 2026-06-29**:

```text
screen-mapping/player-landscape-mode.md -> ok
player-spec §Landscape -> ok
```

Simulated bad commit (staged markdown with `Vidoe-Test`) still blocked by sibling `server_doc` sanitization guard in same hook script.

---

## Phase 0 artifacts (structure checkpoint)

| Artifact | Path |
|----------|------|
| Main plan | `docs/plans/2026-06-25-landscape-horizontal-player-mode.md` |
| Screen mapping | `docs/frontend-contracts/screen-mapping/player-landscape-mode.md` |
| Visual spec §4 | `docs/frontend-contracts/ui-style-reference/specs/player-spec.md` |
| Integration scaffold | `tests/integration/player-landscape-mode.integration.mjs` |
| Fixture | `tests/integration/player-fixtures.mjs` (`PLAYER_HORIZONTAL_EPISODE_ID`) |

---

## Failure-mode review lens

| failure | Project observation |
|---------|---------------------|
| intent oscillation | Plan orders Phase 0 → 1 → 3; guard blocks implementation-first commits |
| structure inflation | Phase 0 scoped to mapping + spec regions, not broad player rewrite |
| fake equivalence | Checkpoint cites **non-regression** of portrait / snap / preview paths — integration-owned |
| abstraction orphan | N/A — no unused parser/seam left behind (contrast `01-dogfood` prep-02) |
| replacement parity misuse | Plan explicitly frontend-only; no replacement inventory |
| illegal transition | Not yet wired — would need `steps[]` YAML in project plan for planvalidate |

---

## Secondary pilots (same pattern family)

| Pilot | Maps to | Evidence |
|-------|---------|----------|
| Navigation Phase B (`2026-06-18-1030-ui-continuation-runtime-pilot.md`) | bounded structure refactor before feature closure | **completed** — contracts split before code churn |
| Design Contract B0 (`2026-06-16-design-contract-validation-pilot.md`) | structure stub before scanner/feature | **completed** — authority stub before static enforcement |
| Plan-first ordering README | contract-before-implementation culture | `docs/plans/README.md` §Plan Phase 0 guard |

**Out of scope for this plan**: `server_doc` sanitization (documentation boundary governance — different failure family).

---

## What this evidence supports in `_plan.md`

| Plan section | Supported? | Notes |
|--------------|------------|-------|
| Phase 4 dogfood | **yes** — second path | Project-layer happy/partial path vs ai-skill force_exit |
| Phase 3 routing cross-links | **partial** | Vidoe-Test consumes `execution-modes` **concept** via project plan wording, not ai-skill `loading_surfaces` |
| Phase 5 glossary | **no** | Project uses `Phase 0/1`; glossary terms not registered in Vidoe-Test |
| Validator hook / enforcement | **explicitly no** | Aligns with observation-period decision — mechanical guard is **project overlay**, not ai-skill commit-msg |

---

## Recorded fields (Phase 4 checklist)

| Field | Value |
|-------|-------|
| change_kind | feature |
| execution_mode | preparatory_refactoring (project analog) |
| intent sequence | structure (Phase 0) → feature (Phase 1+, in progress) |
| transitions | Phase 0 complete → Phase 1 allowed (human + file_exists gate) |
| force_exit trigger | none |
| checkpoint | portrait player regression suite + spec/mapping existence |
| blocking | project Cursor hook only; `Blocking=false` for ai-skill planvalidate |
