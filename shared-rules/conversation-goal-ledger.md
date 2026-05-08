# Conversation Goal Ledger

This rule defines a tool-neutral temporary ledger for active conversation goals. It helps a later agent recover unfinished work after an interrupted session, context compaction, model switch, or multi-agent handoff.

The ledger is for project-local state, not reusable knowledge. It should live under the current project, be excluded from git, and be deleted when the goal is fully complete.

The goal ledger gives humans and agents an immediate overview after a long conversation: what is still unfinished, what needs a decision, what should be prioritized next, and what needs more work before it can be considered complete.

## Purpose

Use a goal ledger when a conversation has work that spans more than one action, includes multiple goals, can be interrupted, or needs clear completion criteria.

The ledger must answer:

| Field | Required content |
| --- | --- |
| Goal | What user-visible outcome is being pursued. |
| Priority | `P0`, `P1`, `P2`, or `P3`. |
| Status | `active`, `paused`, `blocked`, `needs-validation`, `superseded`, or `complete-pending-delete`. |
| Owner | Current agent/tool owner and timestamp. |
| Source | User request or instruction that created the goal. |
| Scope | In scope, out of scope, and affected project/repo. |
| Subgoals | Child goals or checklist items when the goal is decomposed. |
| Planning / todo links | Planning document path, plan section, TodoWrite IDs, checklist items, or external issue links related to this goal. |
| Open work / decisions | What is not done yet, what decision is needed, or what needs strengthening. |
| Dependencies | Required user answer, external command, file, agent, or upstream goal. |
| Next action | The next concrete step a new agent should take. |
| Completion criteria | What must be true before the file can be deleted. |
| Validation | How completion was or will be verified. |

## Location

Store ledgers in the project being worked on:

```text
<PROJECT_ROOT>/.agent-goals/
  README.md              # main goal table / quick locator
  goals/
    P1-<slug>.md
  locks/
    <goal-id>.lock/
```

Do not store the canonical ledger in a tool-specific configuration directory, because this workflow must work across different agent tools. A tool may read or remind about this directory, but it is not the source of truth.

The `.agent-goals/` directory is temporary project state and should not be committed. Prefer excluding it through `.git/info/exclude` so business repositories do not receive policy churn. A project may add `.agent-goals/` to `.gitignore` only when the team wants that convention tracked.

## When To Create Or Update

Before substantive work, first run or perform the equivalent of:

```text
<AI_SKILL_REPO>/scripts/agent-goals.sh --project <PROJECT_ROOT> status
```

If the ledger does not exist and any trigger below applies, initialize it before continuing:

```text
<AI_SKILL_REPO>/scripts/agent-goals.sh --project <PROJECT_ROOT> init
```

Create or update a goal file when:

- A user asks for implementation, analysis, planning, review, debugging, or repository updates that can span more than one tool call.
- A task has multiple goals or priorities.
- The agent observes modified, staged, untracked, or otherwise dirty project files and intends to continue work in that project.
- The agent creates a tool-level todo list or resumes a previous todo list whose items are not all complete.
- The user says to continue a prior multi-step task, especially after context compaction, interruption, or a different side quest.
- A task is paused, blocked, superseded, or waiting for user input.
- A goal is decomposed into smaller goals.
- An agent is about to stop, compact context, switch mode, launch subagents, or hand off.
- A user changes priority, adds a new target, or redirects the conversation.

For very small one-message answers, the ledger is optional. If any work remains after the response, if files were changed, or if the working tree is dirty for the active task, it is no longer optional. Do not treat the tool todo list as a substitute for this project-local ledger; todos track execution steps, while `.agent-goals/` tracks user-visible goals and handoff state.

## Goal File Template

Use Markdown so any tool can read it:

```markdown
---
id: P1-short-slug
priority: P1
status: active
owner: <agent/tool/session>
created: <ISO-8601 timestamp>
updated: <ISO-8601 timestamp>
project: <PROJECT_ROOT or project label>
---

# <Goal title>

## Source Request
<User request or concise quote.>

## Scope
- In:
- Out:
- Affected paths/repos:

## Subgoals
- [ ] <subgoal>

## Planning / Todo Links
| Type | Reference | Status / Note |
| --- | --- | --- |
| plan | <path#section or none> | <why it matters> |
| todo | <todo id / checklist item / issue> | <pending / in_progress / completed / blocked> |

## Open Work / Decisions
- Missing work:
- Decision needed:
- Needs strengthening:

## Dependencies
- <none / user answer / external state / parent goal>

## Progress
- <timestamp>: <what changed>

## Next Action
<The next concrete action for a future agent.>

## Completion Criteria
- [ ] <observable completion condition>

## Validation
- <diff review / test / lint / source checked / user confirmation / not yet validated>

## Handoff Notes
<Risks, blockers, assumptions, and recovery hints.>
```

Do not write secrets, tokens, raw private data, reservation codes, personal addresses, or private host details into the ledger. Use redacted labels or project-local references.

## Main Goal Table

Keep `<PROJECT_ROOT>/.agent-goals/README.md` as the primary locator for active goals. It should contain a compact table that links to each goal file:

```markdown
| Priority | Status | Goal | Planning / Todo Links | Next Action | Updated |
| --- | --- | --- | --- | --- | --- |
| P1 | active | [Short title](goals/P1-short-slug.md) | plan: docs/plan.md#section; todo: implement-api | Run validation | 2026-05-08T00:00:00Z |
```

The main table is for quick recovery. It should not replace the detail in each goal file.

Update the table when a goal is created, paused, split, linked to a todo, or completed. When a goal file is deleted after validation, remove it from the table.

## Planning And Todo Links

When a planning document, checklist, or tool-level todo list exists, connect it to the goal ledger:

1. Put the goal ID next to the relevant plan section, checklist item, or todo when practical.
2. Record the plan path, section anchor, TodoWrite ID, checklist item, or issue ID under `Planning / Todo Links` in the goal file.
3. If a todo becomes a separate resumable work item, either add it as a subgoal or split it into a child goal.
4. When a todo is completed, update the goal progress and validation notes before deleting the goal.
5. If a todo is cancelled because the user changed direction, mark the linked goal `paused` or `superseded` and record the reason.

The goal ledger tracks user-facing intent; todo tools track execution steps. Keep both connected so a future agent can jump from a high-level goal to the exact plan/todo item and back.

Document-level TODO lists are local to a file and should appear near the top of that file. See [`document-todo-list.md`](document-todo-list.md). When a document TODO is part of a larger user goal, link it from `Planning / Todo Links` or `Open Work / Decisions`.

## Priority Rules

Use these priorities:

| Priority | Meaning |
| --- | --- |
| `P0` | User-blocking, safety/secret risk, data-loss risk, or explicit urgent request. |
| `P1` | Current primary user goal. |
| `P2` | Important follow-up or validation needed after the primary path. |
| `P3` | Nice-to-have cleanup, optional refactor, or low-risk follow-up. |

Only one `P1` should normally be active per conversation. If a new `P1` arrives, pause or supersede the previous `P1` with a reason and next action.

## Decomposition

When a goal becomes too broad, split it:

1. Keep a parent goal with the user-facing outcome.
2. Add child goals or checklist items for independently resumable work.
3. Record dependencies between child goals.
4. Promote a child goal to `P1` only when it is the current work focus.

Do not hide a discovered subgoal only in chat. If it affects completion, record it in the ledger.

## Goal Transfer

When the user redirects the task:

1. Update the old goal to `paused` or `superseded`.
2. Record why it was paused and what would resume it.
3. Create or promote the new goal with the correct priority.
4. Make the final response clear about which goal is now active.

If the new goal conflicts with a high-risk unfinished goal, flag the conflict before switching.

## Multi-Agent Safety

Agents must coordinate through lock directories:

```text
<PROJECT_ROOT>/.agent-goals/locks/<goal-id>.lock/
  owner
  pid
  startedAt
```

Use atomic directory creation for locks. If another active lock exists, do not modify that goal. Report the owner, age, and intended next step. A stale lock may be removed only after checking the recorded PID/session is no longer active or after user approval.

Recommended default TTL: 30 minutes. A tool may override TTL when a legitimate long-running task is active.

## Completion And Deletion

Delete a goal file only when all are true:

1. The completion criteria are met.
2. Validation has run or the user explicitly accepts the result.
3. No child goal remains active, blocked, or needs validation.
4. The final answer or handoff states the outcome.

If the work is done but validation is missing, set status to `needs-validation` rather than deleting it.

If the goal was superseded, keep it until the user accepts the new direction or the reason is clear enough for a future agent. Then it may be deleted or archived according to project preference.

## Relationship To Ai-skill Writeback Transactions

This ledger is separate from the Ai-skill writeback transaction in [`dependency-reading.md`](dependency-reading.md).

- Conversation goal ledger: project-local, temporary, not committed, tracks user goals and handoff state.
- Ai-skill writeback transaction: repository-specific, committed/pushed, tracks changes to this knowledge base and configured tool sync/mirrors.

When updating this repository, both may apply: the project goal ledger tracks the user-facing task, while the Ai-skill transaction must still close through diff review, linked updates, sync, commit, push, reread, and clean status.

## Tool Integration

Tools may automate ledger checks, but automation is advisory unless it can deterministically validate the goal. A hook or script should remind, create, or inspect goals; it should not silently mark goals complete without the completion criteria and validation evidence.

For tool-specific handling, see the matching documentation under `ai-tools/`.

← [Back to shared rules index](README.md)
