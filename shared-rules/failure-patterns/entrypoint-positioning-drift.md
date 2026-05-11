# Entrypoint Positioning Drift

Status: validated
Class: validation-gap

## Trigger

When a repository, skill, shared rule, architecture document, or tool adapter is renamed, re-scoped, or promoted to a new top-level concept.

## Failure Mode

The agent updates secondary references, mid-document sections, links, or filenames but misses the primary entrypoint positioning: the root title, opening paragraph, index summary, or first screen a future reader will see.

## Risk

Users and future agents see stale framing first, even though deeper links point to the new concept. This creates confusion about the source of truth, weakens architecture adoption, and can make a completed update look incomplete.

## Required Agent Action

For naming, branding, architecture, or top-level scope changes:

1. Identify the primary entrypoint files before editing.
2. Update the title and opening positioning, not only links or mid-document sections.
3. Search for old names, old slugs, and old framing after edits.
4. Re-read the entrypoint as a user would see it from the top.
5. If a user points out the drift, run the failure learning loop immediately instead of only patching the missed line.

## Prevention Gate

Before claiming completion, answer:

- Does the first heading use the new canonical name?
- Does the first paragraph describe the new role/scope?
- Do root indexes and architecture/tool indexes point to the new canonical file?
- Does a search for old names and old slugs return only intentional historical references?

## Validation

Read back the root entrypoint and the changed architecture/tool/shared indexes from line 1. Run an exact search for the old name and old slug. Confirm `git status --short --branch` is clean after commit/push/readback.

## Linked Rules

- [`../failure-learning-system.md`](../failure-learning-system.md)
- [`../goal-action-validation.md`](../goal-action-validation.md)
- [`../linked-updates.md`](../linked-updates.md)
- [`../content-layering.md`](../content-layering.md)
