# Cognitive State Governance

## Purpose

This document translates cognitive state and evidence governance into AI runtime governance. It closes the gap before recovery: the agent may not have failed yet, but its assumptions, evidence quality, claim scope, or intent can already be unstable.

This is governance source, not a persistent runtime state model. Runtime promotion must stay minimal and must follow the compression rule in this document.

## Boundary

| Concern | Owner |
| --- | --- |
| Why belief drift and evidence mismatch matter | `intelligence/engineering/agent-architecture/` |
| Assumption, evidence, confidence, claim scope, autonomy, and contamination governance | This document |
| Executable MUST / forbidden behavior | `enforcement/evidence-hierarchy.md` and `enforcement/escalation-policy.md` |
| Domain-specific reload and evidence policy | `metadata/recovery/` and `metadata/evidence/` |
| Machine-enforced runtime guards | `runtime/README.md` after signal compression |
| Failure reproduction | `validation/scenarios/failure-derived/` |

## Cognitive Governance Model

### Assumption Lifecycle

| State | Meaning | Allowed use |
| --- | --- | --- |
| `assumption` | Unvalidated framing or inference | May guide discovery only |
| `tentative_belief` | Supported by limited evidence | May guide low-risk local action |
| `validated_belief` | Supported by qualified evidence with matching scope | May drive execution |
| `contradicted_belief` | Overturned by stronger evidence | Must not drive execution |
| `suspended_belief` | Held while recovery or rediscovery runs | May be referenced as prior frame only |
| `deprecated_belief` | Superseded after recovery / replan | Must be pruned from active execution |

Rules:

- An assumption must not be written as a conclusion.
- Low-confidence assumptions must not drive high-risk actions.
- Contradicted beliefs must be suspended before execution continues.
- Deprecated beliefs must be removed from the active execution graph.

### Assumption Ledger

The assumption ledger is execution-scoped and lightweight. It does not replace `.agent-goals/`, the dependency-read ledger, durable plans, or validation evidence.

Minimum fields:

| Field | Meaning |
| --- | --- |
| `assumption` | The belief being used. |
| `source` | `user`, `doc`, `tool_output`, `live_observation`, `inference`, or `memory`. |
| `confidence` | `low`, `medium`, or `high`. |
| `validated_by` | Evidence that supports the belief. |
| `contradicted_by` | Evidence that overturns the belief. |
| `expires_when` | Condition requiring revalidation. |
| `execution_dependency` | Actions, claims, or checkpoints that depend on the belief. |

Use the ledger only for assumptions that affect execution decisions.

## Evidence Qualification

Evidence quality is multi-axis. Authority alone is insufficient.

| Axis | Meaning |
| --- | --- |
| `authority` | Who owns the truth: live observation, owner contract, runtime DB, test output, tool log, inference. |
| `freshness` | Whether the evidence still reflects the current execution moment. |
| `validity` | Whether the evidence still applies to this version, scope, or task. |
| `scope` | Local hook, single screen, single API, repo, workflow, feature, or domain. |
| `observability` | Direct observation, derived signal, summary, memory, or inference. |

Core rule:

> Low-quality or low-scope evidence cannot justify a higher-scope claim or override a higher-quality contradiction.

## Belief Ownership

| Belief type | Authority owner |
| --- | --- |
| UI route / screen state | Live runtime observation |
| API contract | Owner contract or canonical API spec |
| Repo architecture | Current repo files and owner README |
| Execution success | Validation evidence and artifact gates |
| Workflow intent | User goal, goal ledger, and workflow primary source |
| Runtime state | `runtime/runtime.db` plus source and compiler timestamp |

Override rules:

- Lower-owner evidence may trigger review but cannot silently overwrite higher-owner belief.
- User intent can override workflow intent, but live runtime evidence still owns UI / API reality.
- Automation logs can support runtime event claims, but cannot own UI route or feature success.
- Ownership conflict requires evidence-conflict handling.

## Confidence Integrity

| Integrity | Meaning | Required action |
| --- | --- | --- |
| `aligned` | Confidence matches evidence quality and claim scope | Continue normally |
| `inflated` | Confidence exceeds evidence quality | Narrow claim or validate |
| `degraded` | Conflict, staleness, or repeated patch lowers confidence | Downgrade autonomy |
| `unsupported` | Claim has no external evidence | Validate before claiming success |

Temporal decay applies only when one of these appears:

- Long session without recap.
- Multi-recovery chain.
- Context compaction preserved conclusions without evidence.
- Evidence validity window expired.
- Stale repo, UI, workflow, or memory is reused after a task boundary.

Short trivial tasks do not require time-based scoring unless evidence instability appears.

## Claim Scope

Evidence coverage must be at least as broad as the claim.

Forbidden compressions:

- Single hook success -> full feature success.
- Single API pass -> workflow correctness.
- Single grep result -> implementation complete.
- Local test pass -> global behavior verified.

Required completion claims:

- State the claim scope.
- Name evidence used.
- Name uncovered areas when evidence is partial.
- Downgrade confidence if the evidence is local.

## Intent Stability

Intent chain:

```text
goal -> execution target -> active subtask -> current action -> validation target
```

Intent drift signals:

- Subtask expands beyond the accepted plan.
- Recovery becomes the main task.
- Current action no longer serves validation target.
- Repeated patches do not reconnect to success criteria.

Required action:

- Recap original goal.
- Compare current action with the execution graph.
- Prune unrelated subtasks.
- Replan if the drift is material.

## Lineage And Contradiction Propagation

Minimum lineage fields:

| Field | Meaning |
| --- | --- |
| `belief` | Claim or conclusion used for execution. |
| `derived_from` | Evidence or upstream belief. |
| `transformed_by` | Inference, summary, classification, or routing decision. |
| `validated_by` | Evidence that supports the belief. |
| `superseded_by` | Stronger evidence or replacement belief. |
| `dependents` | Claims, checkpoints, or execution nodes that rely on it. |

Minimum propagation scope:

```text
contradicted evidence
-> dependent execution beliefs downgraded
-> dependent checkpoints invalidated
-> execution graph marked unstable
-> autonomy mode reevaluated
```

Do not propagate to every note or thought. Propagate only to execution claims and checkpoints.

## Autonomy Modes

| Mode | Allowed actions | Forbidden actions |
| --- | --- | --- |
| `FULL_AUTONOMY` | Scoped patch, validation, commit-flow work | None beyond normal gates |
| `LIMITED_AUTONOMY` | Small edits, source reload, targeted validation | Broad refactor, unbounded automation |
| `VALIDATION_REQUIRED` | Gather evidence, run checks, compare source-of-truth | Claim success or continue patching before validation |
| `HUMAN_ALIGNMENT_REQUIRED` | Summarize options and ask a blocker question | Autonomous execution beyond agreed next step |
| `READ_ONLY_MODE` | Discovery, read docs, inspect evidence | Write files, commit, or run production actions |

Mode transitions:

- `inflated` or `unsupported` confidence -> `VALIDATION_REQUIRED`.
- Ownership ambiguity or missing canonical authority -> `HUMAN_ALIGNMENT_REQUIRED`.
- Session-global contamination -> `READ_ONLY_MODE` until rediscovery.
- Fresh validation may return autonomy to `LIMITED_AUTONOMY` or `FULL_AUTONOMY`.

## Recovery Exit And Re-entry

Recovery can exit only when:

- Old assumption is invalidated or downgraded.
- New source-of-truth is loaded or marked `source_missing` / `not_applicable`.
- Execution graph is rebuilt.
- Dependent beliefs / checkpoints are downgraded or revalidated.
- Validation evidence is reacquired.
- Autonomy mode is reevaluated.

Re-entry safety:

- The same contradiction class cannot immediately trigger the identical recovery path without new evidence.
- A repeated recovery path must show a strategy change or new source-of-truth.
- Two repeated re-entries for the same class require human alignment or blocker output.

## Cognitive Contamination And GC

Contamination scopes:

| Scope | Meaning | Default action |
| --- | --- | --- |
| `workflow-local` | One workflow frame is stale | Reload workflow and downgrade assumptions |
| `domain-local` | Domain route or checklist may be stale | Reload domain policy and candidate sources |
| `session-global` | Prior frame may pollute the whole session | Read-only rediscovery or human alignment |

Belief GC rules:

- Recovery / replan must remove or downgrade superseded assumptions.
- Deprecated beliefs must not remain execution dependencies.
- Long recovery frames need recap and prune.

## Governance Minimality

Governance depth must stay proportional to risk, ambiguity, autonomy level, evidence instability, and blast radius.

Cost classes:

| Mechanism | Cost |
| --- | --- |
| Confidence signal | Low |
| Evidence qualification | Medium |
| Contradiction propagation | Medium |
| Human alignment | Medium |
| Rediscovery / recovery | High |
| Session-global contamination reset | High |

Rules:

- Use the smallest governance depth that restores safe execution.
- Do not launch full recovery for small tasks unless high-risk conflict appears.
- Cost governance cannot override safety, source-of-truth, security, or data-loss gates.
- If governance work itself drifts from the user goal, recap and prune it.

## Tier Boundary And Meta-stop

| Tier | Purpose | Blocking rule |
| --- | --- | --- |
| Tier 0 | Safety / source-of-truth | Always blocking when violated |
| Tier 1 | Evidence correctness | Blocking when action or claim depends on invalid evidence |
| Tier 2 | Recovery stabilization | Blocking while recovery exit criteria are unmet |
| Tier 3 | Cognitive optimization | Does not block Tier 0-2 execution unless promoted by concrete failure |
| Tier 4 | Meta-governance | Documentation by default; runtime enforcement requires recurring validated failure |

Meta-stop rule:

> If a governance mechanism mainly governs other governance mechanisms, keep it in governance unless a concrete runtime failure requires enforcement.

## Runtime Reduction

Before any runtime guard addition, normalize cognitive signals into the smallest viable primitive set.

Initial runtime primitive set:

| Primitive | Covers | Runtime action target |
| --- | --- | --- |
| `evidence_quality_mismatch` | Low-quality evidence, ownership conflict, source freshness mismatch | `VALIDATION_REQUIRED` or recovery |
| `execution_reliability_degradation` | Repeated patch, stale belief, contradiction accumulation, intent instability | Autonomy downgrade |
| `claim_scope_overreach` | Local evidence used for global claim | Narrow claim or block success declaration |
| `recovery_loop_risk` | Budget exhaustion, unsafe re-entry, no strategy change | Human alignment |
| `cognitive_contamination_risk` | Stale frame crossing workflow/domain/session boundary | Rediscovery / read-only mode |

Runtime should enforce only these minimal signals until validation scenarios justify more.
