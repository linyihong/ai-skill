# Error Learning System

This rule turns repeated agent mistakes into durable prevention. It connects error capture, classification, reusable pattern records, linked updates, and validation so the same failure mode is less likely to recur.

Use this when a user points out an agent mistake, a close-loop gap, wrong source/mirror update, missed dependency, incomplete validation, forgotten goal, unsafe parallel work, or any repeated behavior that should become part of the Ai-skill operating system.

## Core Rule

When an error is discovered, do not only fix the immediate file. Run the error learning loop:

1. **Capture** the error in the current work context: what went wrong, where it was detected, and what user-facing risk it created.
2. **Classify** the failure mode using the taxonomy below.
3. **Contain** the current risk before broad work continues.
4. **Promote** the reusable lesson to the right durable location.
5. **Strengthen** the rule, workflow, checklist, tool adapter, or validation gate that would have prevented it.
6. **Validate** that the prevention can be found and applied by a future agent.

The goal is not to archive blame. The goal is to convert an observed failure into a reusable guardrail with a clear trigger and validation method.

## Error Taxonomy

| Class | Meaning | Common prevention |
| --- | --- | --- |
| `source-mirror-drift` | Agent updated a local tool mirror, project `.cursor`, runtime copy, or generated bundle instead of the canonical source repository. | Require canonical repo check, source-first edit, then tool sync. |
| `dependency-miss` | Agent changed or used a rule/skill without reading required linked dependencies. | Add or strengthen dependency read ledger and linked updates. |
| `goal-ledger-miss` | Multi-step or resumable user goal was not recorded, updated, split, paused, or completed correctly. | Update `.agent-goals/` before continuing and link todos/plans. |
| `validation-gap` | Agent claimed completion without diff review, lints, tests, link check, source check, sync, push, readback, or clean status. | Add a concrete validation gate and report what ran. |
| `scope-drift` | Agent mixed unrelated changes, project incident details, or local absolute paths into reusable docs. | Apply reusable guidance boundary and sanitization. |
| `handoff-gap` | Agent left unclear next actions, blockers, owner/lock state, or remaining decisions. | Update goal ledger, Document TODO, or handoff notes. |
| `tool-strategy-gap` | A reusable rule assumed one tool's behavior instead of isolating tool-specific execution. | Move tool-specific details to `ai-tools/` or a skill adapter. |
| `parallelization-risk` | Multiple agents could edit shared state, git history, migrations, release steps, or rules independently. | Mark non-parallelizable or single-owner and stop on conflicting locks. |

If a failure does not fit a class, add a new class only after checking whether an existing class can describe it clearly.

## Storage Rules

| Content | Durable location |
| --- | --- |
| Current unfinished fix, owner, lock, next action | `<PROJECT_ROOT>/.agent-goals/` |
| Reusable cross-skill error pattern | `shared-rules/error-patterns/` |
| Skill-specific technique or failure lesson | `skills/<skill>/feedback_history/` |
| Tool-specific reminder, hook, prompt, or UI detail | `ai-tools/<tool>.md` or tool config |
| Project incident evidence, raw logs, exact private paths, hosts, tokens | Project docs, issue tracker, or private evidence, not reusable docs |

Do not store secrets, real tokens, raw private data, or local absolute paths in error patterns. Use placeholders such as `<AI_SKILL_REPO>`, `<PROJECT_ROOT>`, `<tool-mirror>`, and `<runtime-copy>`.

## Error Pattern Record

Create or update a reusable pattern when the same error could recur across projects, agents, tools, or skills.

Recommended file:

```text
shared-rules/error-patterns/<short-slug>.md
```

Recommended shape:

```markdown
# <Pattern Title>

Status: candidate | validated | promoted | deprecated
Class: <taxonomy class>

## Trigger
When should an agent suspect this error?

## Failure Mode
What goes wrong in generalized terms?

## Risk
What user-facing, repo, validation, or handoff risk does it create?

## Required Agent Action
What must the agent do next time?

## Prevention Gate
What check would have stopped the mistake?

## Validation
How to confirm the prevention worked?

## Linked Rules
- <shared rule / skill / tool doc links>
```

Keep pattern records short. If a pattern becomes long, split examples into smaller pattern files and keep `error-patterns/README.md` as the index.

## Promotion Decision

After classifying an error, choose the smallest durable promotion target that prevents recurrence:

| Error scope | Promotion target |
| --- | --- |
| One active conversation only | `.agent-goals/` progress or handoff note |
| One reusable document has an open local gap | Document TODO near the top of that document |
| Cross-document or cross-agent workflow failure | `shared-rules/error-patterns/` plus the relevant shared rule |
| Skill-specific repeated mistake | The skill's `feedback_history/` and, when mature, its workflow/checklist |
| Tool-specific execution failure | `ai-tools/<tool>.md`, tool config, or skill tool adapter |

Do not promote a project incident directly into reusable docs. First generalize the cause, trigger, required action, and validation.

## Source And Mirror Errors

`source-mirror-drift` is a high-priority error class for this repository.

When a user asks to update rules, skills, feedback lessons, templates, or Ai-skill guidance, the agent must:

1. Locate the canonical `<AI_SKILL_REPO>` git root.
2. Confirm `git status --short --branch` in that repo.
3. Edit the canonical source files first.
4. Treat `.cursor`, `~/.cursor/skills*`, `~/.cursor/shared-rules`, `~/.cursor/bundles/*`, generated bundles, and project-local mirrors as deployment surfaces.
5. Sync mirrors only after the source repo change is complete.
6. Commit, push, read back, and confirm clean status before claiming the update is complete.

Reference-first tool setup helps because it reduces duplicate copies, but it does not replace the source check. An agent can still write the wrong place unless this gate is explicit.

## Validation

Before closing an error-learning update, verify:

- The immediate issue is contained or explicitly recorded as still open.
- The error class is named.
- The durable location is correct.
- The prevention gate is written where future agents will read it.
- Linked updates have been checked.
- If Ai-skill changed, tool sync, commit, push, readback, and clean status are complete.

← [Back to shared rules index](README.md)
