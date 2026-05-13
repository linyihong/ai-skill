> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/app-development-guidance/development-process.md`](../../../../workflow/app-development-guidance/development-process.md)

### 2026-05-06 - Separate regression from new code validation

Status: promoted

#### One-line Summary

Separate tests that guard existing behavior from tests that prove new code; high total coverage does not prove newly written or AI-generated code is correct.

#### Human Explanation

Teams often see 80% or 90% coverage and assume a pull request is safe. That number usually measures how much old code is guarded, not whether the new behavior was specified, tested, and reviewed. For new requirements and AI-generated code, the safer workflow is BDD first, then tests that fail without the new behavior, then implementation.

Validation should be layered. Use existing regression tests to protect known behavior. Use changed/new-code coverage, unit tests, contract tests, property/invariant checks, mutation/negative checks, database-backed tests, and human review to prove new behavior.

#### Trigger

- A new requirement, bug fix, SDK/tool implementation, or AI-generated code is about to be written.
- The team cites total coverage as proof that new code is safe.
- Rule-heavy, persistence-heavy, or safety/security-sensitive logic is changing.

#### Evidence

- Tool: coverage reports, changed-code coverage, BDD scenarios, unit/contract/integration tests, mutation testing, property-based tests, database fixtures, human review.
- Sanitized excerpt: `old behavior regression coverage != new code validation`.
- Evidence path: project test reports and planning docs; reusable skill stores only the generalized strategy.

#### Generalized Lesson

Before implementation, define:

1. Existing behavior to guard with regression tests.
2. New or changed behavior to prove with BDD and failing tests.
3. Changed/new-code coverage expectation.
4. Mutation, property-based, invariant, or negative tests for rule-heavy logic.
5. Contract tests for provider/consumer compatibility.
6. Database/repository/migration tests for persistence behavior.
7. Human review focus for AI-generated code.

#### Agent Action

- Do not treat total project coverage as proof of new behavior.
- Write or update BDD before production code.
- Add failing tests or executable specs before code when feasible.
- Use mutation/property/contract/database tests when examples alone do not prove the rule.
- Call out when tests only guard old behavior and do not validate the new change.

#### Applies When

- Implementing new requirements or behavior changes.
- Fixing bugs with expected-vs-actual behavior.
- Reviewing AI-generated code.
- Changing domain rules, API contracts, validation, persistence, or security-sensitive behavior.

#### Does Not Apply When

- The change is a pure typo or formatting-only edit.
- The request is a read-only explanation.

#### Validation

- Test plan separates regression coverage from new-code validation.
- BDD exists for the new or changed behavior.
- Changed/new-code tests fail without the intended behavior.
- Critical rules have mutation, property, invariant, or negative checks where practical.
- Human review compares implementation against planning docs, BDD, contracts, and tests.

#### Promotion Target

- `process/README.md`
- `templates/initial-development-docs.md`
- `CHECKLIST.md`
- `WORKFLOW.md`
- `SKILL.md`
- `README.md`
- `feedback_history/README.md`
- `feedback_history/common/README.md`

#### Required Linked Updates

- Updated `process/README.md` with Test Strategy Gate.
- Updated `templates/initial-development-docs.md` with Test Strategy fields.
- Updated `CHECKLIST.md` with test strategy checks.
- Updated `WORKFLOW.md` with test strategy routing and validation notes.
- Updated `SKILL.md` quick start and output style.
- Updated `README.md` goals, scope, and linked update examples.
- Updated `feedback_history/README.md` and `feedback_history/common/README.md`.
