# Rule Weight And Conflict Precedence

This rule defines how agents should weigh rules when several instructions, documents, tool adapters, or compatibility paths appear to conflict.

## Core Principle

Not every rule has the same weight. When rules compete, choose the highest-weight rule that directly governs the current risk, then use lower-weight rules only if they do not weaken the higher-weight requirement.

Rule weight is not about which file was read most recently. It is about what user-facing or repository risk the rule controls.

## Weight Order

| Weight | Rule type | Examples | Agent behavior |
| --- | --- | --- | --- |
| P0 | Safety, authorization, secrets, data loss, destructive actions, source-of-truth integrity | Authorization scope, sanitization, no secrets, no destructive git without approval, source vs mirror boundaries | Must not be bypassed for convenience, speed, or tool limitations. Stop and ask if blocked. |
| P1 | User explicit current request and active goal closure | Latest user instruction, accepted plan, active `.agent-goals/` goal, completion criteria | Drives the current task unless it conflicts with P0. Update or pause older goals when redirected. |
| P1 | Canonical repository writeback and validation gates | Dependency reading, linked updates, diff review, commit/push/readback, clean status | Required before claiming repository changes are complete. Tool reload or local sync cannot replace it. |
| P2 | Cross-repo operating policy | Tool-neutral documentation, failure learning, goal/action/validation, document sizing, neutral language | Apply consistently, but adapt output size to the task. |
| P2 | Skill-specific workflow and checklists | `skills/<name>/SKILL.md`, `WORKFLOW.md`, `DOCUMENTATION.md`, `CHECKLIST.md` | Follow after shared operating rules; do not override shared rules unless the shared rule explicitly delegates to the skill. |
| P3 | Tool adapter and compatibility guidance | `ai-tools/`, `.cursor/rules/`, sync scripts, symlink/bundle/copy snapshot details | Applies only for the active tool or compatibility path. Must not redefine canonical source or make optional sync universal. |
| P3 | Efficiency and style preferences | Decision efficiency, output shape, formatting preferences, optional cleanup | Optimize only after higher-weight requirements are satisfied. |

## Conflict Rules

1. **Higher weight wins.** A lower-weight instruction cannot weaken a higher-weight safety, source, validation, or user-goal requirement.
2. **Specific beats generic within the same weight.** A task-specific accepted plan can refine a general workflow, but it cannot bypass required validation.
3. **Current user request beats stale context.** If the latest user message redirects the task, update goal state and follow the new request unless doing so violates a P0 rule.
4. **Tool adapters do not define source truth.** Tool-specific files may explain how a tool reads or syncs content, but canonical edits still happen in `<AI_SKILL_REPO>`.
5. **Efficiency never skips required dependencies.** Context-saving and decision-efficiency rules can choose the order of reading, not remove required reads.
6. **Compatibility layers stay conditional.** Symlink, bundle, and copy snapshot flows apply only when intentionally in use; reference-first remains the default.

## When Unsure

If the agent cannot determine which rule has higher weight:

1. State the conflict in one sentence.
2. Identify the risk each rule controls.
3. Choose the path that preserves safety, canonical source, validation, and the latest user goal.
4. Ask the user only when the conflict affects scope, permissions, destructive action, or incompatible outcomes.

## Validation

Before closing work that involved rule conflicts, verify:

- The final action did not bypass P0 safety/source/secret rules.
- The latest user request is reflected in the active goal or final answer.
- Required dependency reading and linked updates were completed or explicitly marked not applicable.
- Tool-specific sync or compatibility behavior was applied only when intentionally in use.
- The final response reports validation rather than relying on implied compliance.

← [Back to shared rules index](README.md)
