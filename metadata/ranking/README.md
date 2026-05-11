# Metadata Ranking

`metadata/ranking/` defines how agents should prioritize candidate Knowledge Atoms and source files during context loading.

## Ranking Inputs

Use these fields from `metadata/schema.md`:

| Field | Ranking effect |
| --- | --- |
| `priority` | Higher priority loads first: `P0`, `P1`, `P2`, `P3`. |
| `status` | Prefer `stable` and `validated` over `candidate`; avoid `deprecated` unless needed for compatibility. |
| `confidence` | Prefer `high`, then `medium`, then `low` when the task allows choice. |
| `context_cost` | Prefer lower cost when two sources answer the same question. |
| `depends` | Load required dependencies before the atom. |
| `conflicts` | Pause and resolve conflict when ranking would load incompatible atoms. |
| `when_to_read` | Only rank an atom after its trigger condition matches the task. |

## Default Ranking Order

1. Required safety, source-of-truth, dependency reading, and validation rules.
2. Latest user goal and active `.agent-goals/` state.
3. Current source-of-truth entrypoints, especially `shared-rules/` and `skills/<name>/SKILL.md`.
4. Validated or stable Knowledge Atoms that directly match the task intent.
5. Candidate maps and summaries that help navigate without replacing source behavior.
6. Examples, background references, and optional optimization notes.

## Tie Breakers

When multiple sources appear relevant:

- Prefer the source with the lowest semantic distance to the current decision.
- Prefer the source that gives a concrete validation signal.
- Prefer canonical repository paths over tool mirrors.
- Prefer a short index or summary first only if it points to the canonical source.
- Do not skip required dependencies to save context.

## Stop Conditions

Stop loading more context when:

- The current source answers the decision point.
- Additional sources are lower confidence or duplicate the same guidance.
- A conflict requires rule-weight or user decision.
- A candidate map says the old skill remains source of truth and no promotion is in scope.

## Validation

A ranked route is valid when the final answer or commit can state:

- Which source was loaded first.
- Which dependencies were required.
- Which sources were deferred and why.
- What validation signal confirmed the selected source was sufficient.
