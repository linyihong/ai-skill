# Document TODO Lists

Reusable documents may include a short TODO list near the top so humans and agents can immediately see what is unfinished before reading the whole file.

This rule is for document-local unfinished work. It complements [`conversation-goal-ledger.md`](conversation-goal-ledger.md), which tracks conversation-level goals across files and sessions.

## When To Add

Add a TODO section when a document has any of these:

- Unfinished sections.
- Missing examples, tables, templates, links, or references.
- Claims that still need validation.
- Known weak spots that need strengthening.
- Follow-up work that affects the document's usefulness.
- Open questions that should be resolved before the document is considered complete.

If the document is complete and has no open document-local work, omit the TODO section or state `No open document TODOs` in a compact form.

## Placement

Place the TODO section where it is visible before the detailed body:

1. YAML/frontmatter, if any.
2. Title.
3. Short purpose/overview.
4. `## Document TODO` or `## TODO`.
5. Main content.

Do not bury document TODOs at the end. The purpose is immediate orientation after a long conversation or handoff.

## Template

Use a compact table:

```markdown
## Document TODO

| Priority | Status | TODO | Link | Owner / Goal |
| --- | --- | --- | --- | --- |
| P1 | pending | Add validation checklist for package sync | [Validation](#validation) | `.agent-goals/goals/P1-example.md` |
```

Recommended fields:

| Field | Meaning |
| --- | --- |
| `Priority` | `P0`, `P1`, `P2`, or `P3`, aligned with the goal ledger priority vocabulary. |
| `Status` | `pending`, `in_progress`, `blocked`, `needs-validation`, `done`, or `cancelled`. |
| `TODO` | Concrete unfinished work, not a vague reminder. |
| `Link` | Anchor to the related section in the same document, or a related file/section. |
| `Owner / Goal` | Optional owner, goal file, issue, or todo ID that explains who should decide or finish it. |

## Linking Rules

Every TODO should point somewhere actionable:

- Link to a heading in the same document when the work belongs to a section.
- Link to another document when the missing work lives elsewhere.
- Link to `.agent-goals/goals/<goal-id>.md` when the TODO is part of a conversation-level active goal.
- Link to an issue, planning document, or checklist item when that is the source of truth.

If no useful link exists yet, write `needs anchor` and create the relevant section before closing the TODO.

## Relationship To Goal Ledger

Document TODOs and goal ledger entries should reinforce each other:

- A document TODO is local to one document.
- A goal ledger entry tracks a user-facing objective across one or more documents, tools, or sessions.
- If a document TODO becomes a user-facing objective or spans multiple files, create or link a goal ledger entry.
- If a goal ledger entry depends on a document section, link back to the document TODO or heading.

When a goal is completed, update or remove related document TODOs. When a document TODO remains open, do not delete the linked goal unless the goal's completion criteria explicitly exclude that TODO.

## Maintenance

- Keep the TODO table short. Move large task breakdowns to `.agent-goals/`, an issue tracker, or a planning document.
- Remove or mark TODOs `done` only after the linked work is actually complete and validated.
- Prefer `blocked` or `needs-validation` over deleting unclear work.
- During reviews, check the TODO section before claiming a document is complete.

← [Back to shared rules index](README.md)
