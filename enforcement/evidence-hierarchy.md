# Evidence Hierarchy

## Purpose

This rule defines how agents must compare evidence quality, claim scope, confidence, and ownership before taking action or declaring success.

Use it when:

- Evidence conflicts.
- A claim is broader than the evidence.
- An assumption is being used as a fact.
- Repeated patching or retrying suggests confidence decay.
- User correction or live evidence contradicts the current execution frame.

## Evidence Qualification

Evidence must be evaluated on five axes:

| Axis | Required question |
| --- | --- |
| Authority | Who owns this truth? |
| Freshness | Is it current enough for this execution? |
| Validity | Does it still apply to this version / task / environment? |
| Scope | What claim boundary does it actually cover? |
| Observability | Is it direct observation, derived signal, memory, or inference? |

Authority ranking is not enough. A high-authority source can be stale or out of scope; a direct live observation can be narrow.

## MUST Rules

1. Unvalidated assumptions MUST remain marked as assumptions.
2. Low-quality or low-scope evidence MUST NOT override a higher-quality contradiction.
3. Local evidence MUST NOT justify a global completion claim.
4. Contradicted beliefs MUST be suspended before dependent execution continues.
5. Confidence MUST be downgraded when evidence is stale, contradictory, missing, or narrower than the claim.
6. Evidence ownership conflicts MUST enter escalation / recovery instead of being silently merged.
7. Recovery MUST NOT exit until contradicted assumptions, dependent checkpoints, validation evidence, and autonomy mode are reconciled.

## Forbidden Behaviors

- Treating a hook/log event as proof of UI route, API contract, or feature success.
- Treating one passing test as proof that untested actors, paths, or error modes are correct.
- Treating grep/search presence as proof of implementation completeness.
- Continuing to patch after user contradiction without reloading source-of-truth.
- Reusing a stale checklist or prior execution frame across a workflow or domain boundary without validation.
- Declaring success while confidence integrity is `inflated` or `unsupported`.

## Claim Scope Gate

Before declaring completion, answer:

| Field | Required content |
| --- | --- |
| Claim | What is being asserted? |
| Claim scope | Local file, command, test, API, screen, workflow, feature, or system? |
| Supporting evidence | What evidence supports that exact scope? |
| Uncovered areas | What does the evidence not cover? |
| Confidence integrity | `aligned`, `inflated`, `degraded`, or `unsupported` |

If evidence scope is smaller than claim scope, narrow the claim or gather more evidence.

## Confidence Integrity

| Integrity | Trigger | Required action |
| --- | --- | --- |
| `aligned` | Evidence quality matches claim scope | Continue |
| `inflated` | Claim confidence exceeds evidence coverage | Narrow claim or validate |
| `degraded` | Contradiction, stale source, repeated patch, or user correction | Downgrade autonomy |
| `unsupported` | No external validation | Validate before acting or claiming |

## Ownership Rule

| Belief | Owner |
| --- | --- |
| UI route / screen state | Live observation |
| API behavior | Actual response or owner contract |
| Repo architecture | Current files and owner docs |
| Execution success | Validation evidence and artifact gates |
| Workflow intent | User goal plus workflow source |
| Runtime state | `runtime/runtime.db` plus source timestamp |

Lower-owner evidence may trigger review; it cannot silently overwrite higher-owner belief.

## Autonomy Downgrade

| Signal | Minimum mode |
| --- | --- |
| Unsupported or inflated claim | `VALIDATION_REQUIRED` |
| Ownership ambiguity | `HUMAN_ALIGNMENT_REQUIRED` if source cannot resolve it |
| Repeated patch without new evidence | `LIMITED_AUTONOMY` or recovery |
| Same contradiction class re-enters without new evidence | `HUMAN_ALIGNMENT_REQUIRED` |
| Session-global cognitive contamination | `READ_ONLY_MODE` until rediscovery |

## Relation To Escalation

`evidence-hierarchy.md` prevents weak evidence from becoming a false fact. `escalation-policy.md` controls the real-time response after contradiction or source-of-truth miss is detected.

When evidence conflict or source miss is active, follow `escalation-policy.md` recovery output requirements.

## Validation

Before closing work that uses this rule, confirm:

- Claim scope matches supporting evidence.
- Contradicted assumptions are suspended or removed.
- Required source-of-truth was read or marked missing / not applicable.
- Autonomy mode matches evidence quality.
- Any reusable failure was linked to a validation scenario or documented deferral.

← [Back to enforcement index](README.md)
