> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/app-development-guidance/development-process.md`](../../../../workflow/app-development-guidance/development-process.md)

### 2026-05-05 - Missing requirements block development

Status: promoted

#### One-line Summary

If required behavior, contracts, errors, security, storage, ownership, or tests are missing, the agent must ask for answers or evidence before development continues.

#### Human Explanation

Contract-first development only works if unknowns are made explicit. Some missing information can be backfilled from evidence, and some product intent can remain marked as `unknown`. But when a missing item changes what should be built or tested, continuing implementation would force the agent to invent behavior.

The safe default is to stop at blocker questions: list what is missing, explain what it blocks, ask the user or request evidence, then update the documents before implementation proceeds.

#### Trigger

- Product brief, BDD, domain, API/interface, error handling, storage, security, ownership, or test details are missing.
- An existing project has implementation evidence but not enough to recover critical BDD behavior.
- The user asks to continue implementation while required documents or decisions are incomplete.

#### Evidence

- Tool: project docs, source code, UI behavior, API schemas, tests, logs, fixtures, manual verification.
- Sanitized excerpt: `missing item -> blocker question -> answer/evidence -> document update -> implementation`.
- Evidence path: project repository planning docs; reusable skill stores only the generalized rule.

#### Generalized Lesson

Classify every missing item:

- Recoverable from evidence: backfill and cite evidence.
- Product intent unavailable: mark `unknown` / `open question`.
- Affects behavior, domain invariants, API/interface shape, error handling, security, storage, ownership, or tests: blocker question.
- Does not affect behavior or contracts: non-blocking open question with reason.

Do not proceed with development while blocker questions remain unresolved.

#### Agent Action

- Ask all blocker questions before implementation continues.
- Do not invent missing behavior or contract decisions.
- Update the relevant document after the user answers or evidence is found.
- If the user intentionally scopes out a blocker, record that decision and its impact.
- Keep non-blocking unknowns separate from blockers.

#### Applies When

- Starting a new feature from incomplete requirements.
- Backfilling docs for an existing implemented project.
- Converting APK analysis handoff into development documents.
- Planning implementation slices or tests.

#### Does Not Apply When

- The missing item is unrelated context that does not affect behavior, contracts, or tests.
- The user explicitly asks only for a narrow, non-implementation review and the unknown is outside scope.

#### Validation

- Planning docs include a blocker question table.
- No implementation slice depends on an unanswered blocker.
- BDD and contracts are updated after answers or evidence arrive.
- Non-blocking unknowns explain why they do not affect behavior or contracts.

#### Promotion Target

- `process/README.md`
- `templates/initial-development-docs.md`
- `CHECKLIST.md`
- `WORKFLOW.md`
- `SKILL.md`
- `feedback_history/README.md`
- `feedback_history/common/README.md`

#### Required Linked Updates

- Updated `process/README.md` with Missing Information Gate.
- Updated `templates/initial-development-docs.md` with blocker question tracking.
- Updated `CHECKLIST.md` with blocker checks.
- Updated `WORKFLOW.md` with missing-information blocking behavior.
- Updated `SKILL.md` quick start and output style.
- Updated `feedback_history/README.md` and `feedback_history/common/README.md`.
