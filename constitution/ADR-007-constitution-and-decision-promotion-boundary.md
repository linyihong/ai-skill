# ADR-007: Constitution and Decision Promotion Boundary

## Status

accepted

## Context

The previous outer `decisions/` directory stored formal ADRs, while `memory/decision/` stored session-level decision memory. The name `decisions/` was too broad and could be confused with runtime decisions, project decisions, or memory decisions.

At the same time, decision promotion needed a clearer target rule. A repeated or durable decision should not automatically become an ADR.

## Decision

Rename the formal ADR layer from `decisions/` to `constitution/`.

Use content-based decision promotion:

| Decision content | Target |
| --- | --- |
| Executable rule or cross-agent policy | `enforcement/` |
| Reasoning heuristic, tradeoff, signal, anti-pattern, or failure judgment | `intelligence/` |
| Operational process or repeatable workflow | `workflow/` |
| Runtime gate, activation, phase, obligation, or executable contract projection | `runtime/runtime.db` |
| Architecture-level irreversible or foundational decision | `constitution/ADR-*` |
| Session-scoped replay decision | `memory/decision/` |
| Project-specific decision | `<PROJECT_ROOT>/docs/decisions/` |

Runtime decision-recording canonical config is renamed to `runtime/constitution/decision-recording.yaml` to align the architecture-tier naming.

## Consequences

- `constitution/` is the formal ADR / architecture constitution layer.
- `memory/decision/` remains a session-level decision replay layer.
- `<PROJECT_ROOT>/docs/decisions/` remains the project-local decision tier.
- ADR creation is no longer the default endpoint for every promoted decision.
- Runtime-affecting decisions must update `runtime.db` or an executable YAML contract projected into `runtime.db`.

## Alternatives Considered

- Keep `decisions/`: rejected because it kept confusing formal ADRs with memory and runtime decisions.
- Promote every durable decision to ADR: rejected because executable policy, reasoning heuristics, workflows, and runtime gates have better owner layers.
- Move project decision folders to `constitution/`: rejected because project-local decisions are not part of the Ai-skill constitution.

## Related

- [`README.md`](README.md)
- [`../governance/lifecycle/decision-promotion-pipeline.md`](../governance/lifecycle/decision-promotion-pipeline.md)
- [`../governance/lifecycle/decision-promotion-pipeline.yaml`](../governance/lifecycle/decision-promotion-pipeline.yaml)
- [`../memory/decision/README.md`](../memory/decision/README.md)
- [`../runtime/runtime.db`](../runtime/runtime.db)
