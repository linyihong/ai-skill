# Workflow Activation Engine

> Operational spec for the registry-driven workflow detector. Canonical schema
> lives in [`knowledge/runtime/routing-registry.yaml`](../knowledge/runtime/routing-registry.yaml)
> (§`route_type_spec`, §`activation_mode_spec`, §`activation_triggers_spec`).
> Decision record: [`constitution/ADR-012`](../constitution/ADR-012-route-type-activation-behavior-family.md)
> (route_type semantics) + [`constitution/ADR-006`](../constitution/ADR-006-registry-first-workflow-activation.md)
> (registry-first activation). Source plan:
> [`plans/active/2026-05-31-1900-workflow-activation-engine.md`](../plans/active/2026-05-31-1900-workflow-activation-engine.md).

## Philosophy

A route should become active because the **registry** says its triggers fired —
not because a human remembered to wire one more `activation #` into a hook. The
engine reads each route's self-declared `route_type` + `activation_triggers`,
derives the route's `activation_mode`, and lets a deterministic detector decide
which route(s) the current task matches. Adding a workflow = adding a registry
record, not editing detector code.

This preserves the ADR-006 invariant ("trigger conditions live in the registry,
not in per-workflow hook branches") and extends it with an explicit activation
behavior contract per route.

## Two-Phase Activation (circular-dependency break)

The detector must answer "which workflow applies?" **before** the agent reads
files, because one purpose of activation is to force the agent to read the
workflow `primary_source` first. But the richest signal — what a file actually
contains (`artifact_signals`) — is only observable **after** a Read. Depending
on it creates a deadlock:

```
need to Read a file → to know which workflow → but the workflow says
"read me before you Read anything" → which requires already knowing the workflow
```

The engine breaks this by classifying every signal by **when it is observable**:

| Phase | Field group | Observable | Role |
| --- | --- | --- | --- |
| 1 — activation | `activation_any_of` (`user_signals`, `context_signals`) | pre-Read — inbox text, filenames, paths, cwd | **can gate activation** — any hit locks `active_route` |
| 2 — reinforcement | `reinforcement_any_of` (`artifact_signals`) | post-Read — content of files already read | **confirmation only** — raises confidence, never activates alone |

- **Phase 1 signals** are available without reading file content: conversation
  text (`user_signals`) and filename/path globs (`context_signals`).
- **Phase 2 signals** accumulate naturally as the agent reads. A Phase-1 hit
  followed by a Phase-2 hit means "direction confirmed". A Phase-1 *miss* with a
  Phase-2 hit is a `late-detected` event — logged for trigger-coverage analysis,
  but it does **not** retroactively rewrite activation history.

### Deterministic rule (no scoring)

> ANY single entry in ANY `activation_any_of` sub-array hits → activate.

No weighting, no score thresholds, no ranked tie-breaks. When more than one
route activates, the engine does **not** auto-pick — it routes to
[`workflow/workflow-routing.md`](../workflow/workflow-routing.md) Stage 2
ambiguity adjudication and lets the agent choose. A zero-match task is not
blocked; it is recorded as `no-match` for future trigger-coverage review.

## Backward Compatibility

The legacy flat form stays valid **indefinitely**. Routes that predate the
two-phase schema declare triggers directly under `activation_triggers`; the
detector normalizes them:

| Legacy flat field | Normalized to |
| --- | --- |
| top-level `user_signals` | `activation_any_of.user_signals` |
| top-level `context_signals` / `file_change_globs` | `activation_any_of.context_signals` |
| top-level `artifact_signals` | `reinforcement_any_of.artifact_signals` |
| top-level `task_intents` | unchanged (shared field, orthogonal to phases) |

The routes already carrying `activation_triggers` are **not rewritten** when the
schema lands — the normalizer handles them. New triggers may use the explicit
two-phase form.

## Activation Mode Capability Matrix

`route_type` derives a default `activation_mode`; the mode determines what the
engine may do with the route. (Full matrix: registry §`activation_mode_spec`.)

| Mode | Activates | Preload | Reinforce | Can conflict | Auto-expire | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| `always-on` | n/a (always loaded) | yes | — | no | no | bootstrap / runtime_core only |
| `auto-detect` | detector, on trigger hit | no | yes | yes → Stage 2 | task end | requires `activation_triggers` |
| `on-demand` | only on explicit user invocation | no | no | no | per-turn | reference docs |
| `advisory` | never standalone-locks `active_route` | no | yes | no | task end | weak signals OK; may suggest promotion |
| `manual-lock` | user only (runtime-assigned) | — | — | no (user adjudicated) | sticky | overrides detector; never author-declared |

## Session-Entry-Point Classification Heuristic

To assign a mode to a mixed-layer (`must-declare`) route, apply the mechanical
test (canonical in registry §`route_type_spec.classification_heuristic`,
durable in ADR-012):

- The route can stand alone as a user's **primary task entry** → primary
  candidate → default **`auto-detect`**.
- The route only has meaning **attached to another route** → secondary
  candidate → **`advisory`** or **`on-demand`**.

This supersedes subjective "primary/secondary" prose as route count scales to
the hundreds.

## Scope Boundaries

**IS for**: the registry schema, the deterministic activation contract, the
two-phase signal classification, and how a route's mode governs detector
behavior.

**NOT for**: the detector's Go implementation (Phase 3 — `detector.go`),
in-memory `RuntimeContext` state (Phase 4), the per-turn obligation wiring
(Phase 5), or the Discovery→Detector feedback loop (Phase 6). Those are tracked
in the source plan and will link back here as they land.

## Related

- [`constitution/ADR-006`](../constitution/ADR-006-registry-first-workflow-activation.md) — registry-first activation invariant
- [`constitution/ADR-012`](../constitution/ADR-012-route-type-activation-behavior-family.md) — `route_type` = activation behavior family
- [`knowledge/runtime/routing-registry.yaml`](../knowledge/runtime/routing-registry.yaml) — canonical schema (`route_type_spec` / `activation_mode_spec` / `activation_triggers_spec`)
- [`workflow/workflow-routing.md`](../workflow/workflow-routing.md) — Stage 2 ambiguity adjudication for multi-route hits
