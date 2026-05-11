# Metadata Rules

`metadata/rules/` is the index for operational metadata rules. These rules explain how fields from `metadata/schema.md` should be used by routing, governance, summaries, and future graph/runtime surfaces.

## Rule Set

| Rule | Purpose |
| --- | --- |
| [`ranking/`](../ranking/README.md) | Decide which atom or source to read first when several are relevant. |
| [`confidence/`](../confidence/README.md) | Describe evidence strength and when an atom can move from candidate to validated or stable. |
| [`compatibility/`](../compatibility/README.md) | Preserve old skill entrypoints and tool compatibility while new layers evolve. |

## Boundary

- Metadata rules do not override executable shared rules.
- If metadata and shared rules conflict, follow `shared-rules/rule-weight.md`.
- Metadata can reduce context loading cost, but it cannot skip required dependency reading.
- Old `skills/` source files remain source of truth until lifecycle promotion gates pass.
