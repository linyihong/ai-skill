> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/software-delivery/execution-flow.md`](../../../../workflow/software-delivery/execution-flow.md)

### 2026-05-11 - Public mirror drift gate

Status: candidate

#### One-line Summary

Public SDK delivery repositories can accumulate release-facing runtime or README fixes; check both directions before overwriting them from private source.

#### Human Explanation

Treating a public package repo as a disposable mirror is unsafe once maintainers can patch it directly. A later private-to-public sync can silently erase public hotfixes unless agents first compare allowlisted runtime/package files in both directions.

#### Trigger

A package sync task finds public runtime/docs that differ from private source after the last private commit.

#### Evidence

- Tool: repository diff review.
- Sanitized excerpt: a public SDK package carried runtime model/pagination and package README changes not yet present in private source.
- Evidence path: project sync docs and package diffs under `<PROJECT_ROOT>`.

#### Generalized Lesson

For private/public SDK pairs, run a reverse drift gate before syncing: public-only runtime or package docs must be merged back into private source or explicitly documented as public-only before private source overwrites public.

#### Agent Action

Before `rsync` or copying allowlisted package files to a public SDK repo, compare public vs private paths, classify differences, merge public-only source changes back, then run the normal public build/closed-loop verification.

#### Goal / Action / Validation

- Goal: prevent public package hotfixes from being lost during source-to-public sync.
- Action: add a reverse-diff checkpoint to SDK sync workflows.
- Validation or reference source: diff review shows no unexpected public-only runtime/docs before sync, or documented exceptions exist.

#### Applies When

- A source repository syncs runtime SDK code/docs into a public delivery repository.
- Public commits may happen outside the private source repo.

#### Does Not Apply When

- The public repo is generated from an immutable artifact and cannot receive direct changes.

#### Validation

Check both repos' git status, diff allowlisted files in both directions, and verify public build/temp dependency resolution after sync.

#### Promotion Target

- `WORKFLOW.md`
- `CHECKLIST.md`
- `process/README.md`

#### Required Linked Updates

- Indexed in `feedback_history/common/README.md`.
- Project-specific paths and commit ids remain in project docs, not this lesson.
