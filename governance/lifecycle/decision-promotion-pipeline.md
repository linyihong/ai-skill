# Decision Promotion Pipeline

This document defines how a session/runtime decision is promoted. The endpoint is not always an ADR.

## Core Rule

Promote decisions by content type, not by a fixed ladder.

| Decision content | Promotion target |
| --- | --- |
| Executable policy or cross-agent rule | `enforcement/` |
| Reasoning heuristic, tradeoff, signal, anti-pattern, or failure judgment | `intelligence/` |
| Operational procedure or repeatable work sequence | `workflow/` |
| Runtime gate, activation, phase, obligation, policy surface, or executable contract projection | `runtime/runtime.db` |
| Architecture-level irreversible or foundational decision | `constitution/ADR-*` |
| Session-scoped decision with future replay value | `memory/decision/` |
| Project-specific decision | `<PROJECT_ROOT>/docs/decisions/` |

## Promotion Flow

```text
runtime/session decision
  -> classify content type
  -> validate repeated or durable value
  -> choose target by content
  -> update linked surfaces
  -> write feedback lesson when reusable
  -> refresh runtime projection when execution-affecting
```

## ADR Boundary

Use `constitution/ADR-*` only when the decision is:

- foundational to the Ai-skill architecture,
- cross-session and cross-project,
- expected to remain stable,
- expensive to reverse,
- needed to explain why the system is shaped this way.

Do not promote every repeated decision to ADR. Many decisions are better represented as enforcement rules, intelligence atoms, workflow contracts, or runtime gates.

## Runtime Rule

If the target affects agent execution, it must either:

- have an executable YAML contract with `runtime_projection.enabled: true`, or
- update a canonical runtime document inside `runtime/runtime.db`.

## Related

- [`decision-promotion-pipeline.yaml`](decision-promotion-pipeline.yaml)
- [`executable-contract-boundary.md`](executable-contract-boundary.md)
- [`../../constitution/README.md`](../../constitution/README.md)
- [`../../memory/decision/README.md`](../../memory/decision/README.md)
