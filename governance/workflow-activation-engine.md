# Workflow Activation Engine

> Operational spec for the registry-driven workflow detector. Canonical schema
> lives in [`knowledge/runtime/routing-registry.yaml`](../knowledge/runtime/routing-registry.yaml)
> (§`route_type_spec`, §`activation_mode_spec`, §`activation_triggers_spec`).
> Decision record: [`constitution/ADR-012`](../constitution/ADR-012-route-type-activation-behavior-family.md)
> (route_type semantics) + [`constitution/ADR-006`](../constitution/ADR-006-registry-first-workflow-activation.md)
> (registry-first activation). Source plan:
> [`plans/archived/2026-05-31-1900-workflow-activation-engine.md`](../plans/archived/2026-05-31-1900-workflow-activation-engine.md).

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

**accepted format ≠ canonical format.** The legacy flat form is permanently
*accepted* (always parseable, always runnable) but is **not** a second
first-class schema. The two-phase form is canonical; tooling (compile / export /
rewrite) normalizes toward it — the same pattern as JSON5/TypeScript: read many
shapes, emit one. This prevents the flat form from degrading into an
ever-larger special case as the model grows (`activation_all_of` /
`activation_none_of` / …). `deprecation_policy.remove_after: never`, but new
authoring uses the canonical form.

> Executor status: the read-path normalizer (`normalizeRouteTriggers`) is
> implemented and tested (legacy≡canonical, canonical idempotent, merge+dedupe)
> — backward-compat is runtime-verified, not just spec-claimed. Emitting
> canonical on the write path (rewrite-on-export) is a documented future
> enhancement; no consumer needs it today.

Routes that predate the two-phase schema declare triggers directly under
`activation_triggers`; the normalizer folds them:

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

### Activated ≠ Activatable — the three-layer selection (invariant)

The detector reports a route as `Activated` when an **activation signal
matched** — a detector fact. That is NOT the same as **Activatable** — whether
the route's `activation_mode` permits taking activation ownership
(`can_activate`). Conflating the two lets an `advisory` route lock
`active_route`, which would falsely trigger the Phase 5 primary_source gate and
block tool calls (`can_activate: false` silently becoming `true` + `can_block:
true`). `RuntimeContext` therefore selects in three layers:

| Layer | Field | Meaning |
| --- | --- | --- |
| Matched | `DetectedRoutes` | every route whose signals matched (incl. advisory **suggestions**) |
| Activatable | `CandidateRoutes` | the subset whose `EffectiveMode.CanActivate()` (auto-detect only) |
| Selected | `ActiveRoute` | the single candidate (empty on 0 / conflict) |

**Invariant:** `ActiveRoute != "" (detector path) ⇒ EffectiveMode.CanActivate()`.
Advisory routes appear only in `DetectedRoutes` as suggestions; they never
become candidates, never lock, never fire the gate. Conflict is counted over
candidates, not over matched routes. The invariant is enforced structurally
(`ActiveRoute` is only assigned from `CandidateRoutes`) plus a fail-safe guard
and tests (`TestAdvisoryNeverLocksActiveRoute`,
`TestAutoDetectPlusAdvisory_AdvisoryIsSuggestionOnly`).

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

## RuntimeContext Lifecycle (Phase 4.0)

`RuntimeContext` is the in-memory workflow-activation state for one task,
derived from the transcript by `BuildRuntimeContext` (Go: `runtime_context.go`).
Inspect it with `ai-skill runtime workflow-context --transcript <jsonl>`.

**Persistence model — rebuilt, not stored.** The PreToolUse hook runs as a
fresh process per tool call, so there is no live in-memory object shared across
calls. Because the detector is deterministic, the context is *rebuilt from the
transcript on every invocation* and yields the same answer — no store is
required for in-task correctness. This is precisely why SQLite persistence is
deferred (Phase 4.1, conditional): nothing in-task needs it. A store becomes
necessary only for cross-session replay, analytics, or multi-agent handoff.

**States:** `no-match` → nothing activated; `detected` → ≥1 route activated
(`ActiveRoute` set iff exactly one; `Conflict=true` when >1, and `ActiveRoute`
stays empty — the engine never auto-picks); `locked` → user manual-lock.

**Lifecycle rules:**

1. **Substantive gate (vocabulary, not length).** A turn carries task intent if
   it contains a domain noun (aggregated live from every participating route's
   `user_signals` — the registry *is* the vocabulary) or an action verb. An
   8-char Chinese message can be a full task; a longer greeting is not. Length
   thresholds are explicitly rejected (they mis-fire).
2. **Explicit pivot** (`換任務` / `現在我要` / `new task` / `switch to` …) →
   detection re-runs over post-pivot turns only; pre-pivot routes do not linger.
3. **Manual lock** (`鎖定` / `之後都用` / `lock to` …) → if the lock turn names
   exactly one participating route's signals, `ActiveRoute` locks with
   `EffectiveMode=manual-lock` (sticky; the detector yields). Ambiguous locks
   (zero or >1 matching routes) do **not** lock — no guessing.
4. **Manual unlock** (`回到自動偵測` / `unlock` / `解鎖` …) → restores
   auto-detection.
5. **NO implicit keyword-drift invalidation.** Consecutive turns without the
   original trigger keyword do NOT invalidate the active route — a normal
   drill-down (many in-domain sub-questions) would otherwise be mis-read as a
   topic shift. Keyword absence only updates `LastReinforcedAt` (optional soft
   warning), never auto-invalidates.

## Discovery Bridge (miss-path advisory fallback)

The deterministic detector covers the **hit path**: when a task matches a
route's `activation_triggers`, the gate blocks until the route's
`primary_source` is Read. The **miss path** (no trigger matches) is `fail-open`
by design — activation must not depend on reading content (pre-Read
circular-dependency break). But "fail-open" left the miss path with *no
mechanical fallback at all*: the agent was expected to self-route by intuition,
and when that failed the task completed with governance gates silently unmet
(2026-06-05 incident; see
[`enforcement/failure-patterns/detector-miss-no-fallback.md`](../enforcement/failure-patterns/detector-miss-no-fallback.md)).

The **Discovery Bridge** is the miss-path fallback. On a detector miss it
mechanically runs **Light Discovery** over pre-Read cheap signals (user message
tokens, artifact basenames/paths/extensions, frontmatter head bytes, cwd,
project-overlay metadata), scores the registry's routes, and — if the best
candidate clears a confidence threshold — injects a PreToolUse **advisory**
listing the top candidate route(s) and their `primary_source`. A deferred
**Deep Discovery** phase piggybacks the agent's natural next Read to accumulate
content signals. Landed surface: `scripts/ai-skill-cli/internal/app/discovery.go`
(`RunDiscoveryBridge`); plan `2026-06-06-1700-workflow-activation-discovery-bridge`.

Two invariants keep Discovery from corrupting the activation contract:

- **Advisory, never a gate.** A Discovery proposal MUST NOT satisfy
  `activation_triggers` and MUST NOT block a tool. Activation has exactly two
  paths — (a) deterministic detector match, (b) user manual-lock — and Discovery
  is never a third. `if proposal.confidence > X then auto_activate()` is
  permanently forbidden in `detector.go` / `hooks.go` / `routing-registry.yaml`.
- **Scoring is legal here precisely because it is advisory.** The
  deterministic-rule / no-scoring requirement applies to the *activation gate*;
  Discovery is *ranking*, not gating, so weighted scoring is appropriate and
  must stay on the advisory side of the line.

**Why Discovery Bridge is NOT a mechanical `enforcement-registry.yaml`
rule_class.** Every `coverage: mechanical` rule_class binds an executor that
*blocks* (compile/commit/PreToolUse deny). Discovery never blocks — it only
advises — so registering it as a mechanical rule_class would misrepresent an
advisory ranking as an enforcement gate, and would invite the very
auto-activation the invariants above forbid. The Discovery Bridge is therefore
deliberately *unregistered*; its governance lives here and in its plan, and its
efficacy is measured empirically (Phase D KPI), not asserted by a registry
coverage cell. The *gap it mitigates* (`detector-miss-no-fallback`) is the
governed object; the advisory itself is not an enforcement obligation.

## Scope Boundaries

**IS for**: the registry schema, the deterministic activation contract, the
two-phase signal classification, how a route's mode governs detector behavior,
and the Discovery Bridge's advisory contract + non-registration rationale.

**NOT for**: the detector's Go implementation (Phase 3 — `detector.go`),
in-memory `RuntimeContext` state (Phase 4), the per-turn obligation wiring
(Phase 5), or the Discovery Bridge's runtime implementation details and
empirical tuning (Phase 6 — `discovery.go`; Phase A landed, Phase B deferred,
Phase D = 3-week empirical). Those are tracked in the source plans and link
back here as they land.

## Related

- [`constitution/ADR-006`](../constitution/ADR-006-registry-first-workflow-activation.md) — registry-first activation invariant
- [`constitution/ADR-012`](../constitution/ADR-012-route-type-activation-behavior-family.md) — `route_type` = activation behavior family
- [`knowledge/runtime/routing-registry.yaml`](../knowledge/runtime/routing-registry.yaml) — canonical schema (`route_type_spec` / `activation_mode_spec` / `activation_triggers_spec`)
- [`workflow/workflow-routing.md`](../workflow/workflow-routing.md) — Stage 2 ambiguity adjudication for multi-route hits
- [`enforcement/failure-patterns/detector-miss-no-fallback.md`](../enforcement/failure-patterns/detector-miss-no-fallback.md) — the miss-path failure the Discovery Bridge mitigates
- [`governance/lifecycle/capability-discovery-philosophy.md`](lifecycle/capability-discovery-philosophy.md) — Discovery→Detector feedback-loop philosophy
- [`plans/active/2026-06-06-1700-workflow-activation-discovery-bridge.md`](../plans/active/2026-06-06-1700-workflow-activation-discovery-bridge.md) — Discovery Bridge plan (Phase A landed)
