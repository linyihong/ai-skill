# Failure Patterns

This directory stores reusable cross-skill agent failure patterns. Each pattern records a generalized failure mode, trigger, required action, prevention gate, and validation method.

Read this index when [`failure-learning-system.md`](../failure-learning-system.md) tells you to promote or look up a reusable failure pattern.

| Pattern | Class | Status | Summary |
| --- | --- | --- | --- |
| [Entrypoint positioning drift](entrypoint-positioning-drift.md) | `validation-gap` | validated | Prevent agents from updating secondary links or sections while leaving the root title, opening paragraph, or primary entrypoint framing stale after naming or architecture changes. |
| [Source / mirror write drift](source-mirror-write-drift.md) | `source-mirror-drift` | validated | Prevent agents from updating project-local tool mirrors or runtime copies instead of the canonical Ai-skill source repo. |

## Maintenance

- Keep project-specific evidence out of this directory.
- Add a new pattern when a failure mode can recur across projects, tools, skills, or agents.
- If a pattern becomes skill-specific, move the lesson to that skill's `feedback_history/` and link back here only if the cross-skill trigger remains useful.
- If a pattern becomes long, split examples into separate pattern files instead of expanding the index.

← [Back to shared rules index](../README.md)
