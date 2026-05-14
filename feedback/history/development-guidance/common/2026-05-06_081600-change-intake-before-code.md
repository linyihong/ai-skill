> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/software-delivery/development-process.md`](../../../../workflow/software-delivery/development-process.md)

### 2026-05-06 - Change intake before code

Status: promoted

#### One-line Summary

Before code changes, review the planning artifact and classify the work as new requirement, bug fix, refactor, hardening, or docs-only; new requirements must update planning docs before implementation.

#### Human Explanation

Development work often starts from a vague request. If the agent writes code before checking the 企劃書, product brief, issue, ticket, PRD, design note, BDD, or API contract, it may implement the wrong behavior or blur bug fixes with new requirements.

The safe flow is to classify the change first. New requirements and behavior changes require updated planning docs, BDD, contracts, implementation slices, and tests before code. Bug fixes require expected vs actual behavior, reproduction/evidence, affected BDD or missing scenario, impacted contracts/errors, and a regression test plan. Refactors must be proven behavior-neutral or reclassified.

#### Trigger

- The user asks to change app behavior, add a feature, fix a bug, refactor, harden, or implement from planning docs.
- The request is ambiguous about whether it is a new requirement or a bug.
- The project has planning docs that should be checked before code.

#### Evidence

- Tool: product brief, issue/ticket, PRD, design note, BDD, API contract, source code, tests, bug report, logs.
- Sanitized excerpt: `planning artifact -> change classification -> required doc updates -> implementation`.
- Evidence path: project repository planning docs; reusable skill stores only the generalized workflow.

#### Generalized Lesson

Before code:

1. Read the relevant planning artifact.
2. Classify the change.
3. If new requirement: update planning docs, BDD, affected contracts, implementation slices, and tests before code.
4. If bug: record expected vs actual behavior, reproduction/evidence, affected scenario/contract, and regression test.
5. If refactor: prove no behavior or public contract change, or reclassify.
6. Ask blocker questions before implementation if required information is missing.

#### Agent Action

- Do not start code from an ambiguous request.
- Ask whether the change is a new requirement or bug when unclear.
- Do not treat new requirements as bug fixes to skip planning docs.
- Do not treat behavior-changing refactors as internal cleanup.
- Update docs before implementation when the change affects behavior or contracts.

#### Applies When

- Any code implementation is requested under `app-development-guidance`.
- A planning document, issue, or product brief exists.
- The change could affect behavior, API/interface, domain model, errors, storage, security, ownership, or tests.

#### Does Not Apply When

- The user asks only for a typo or formatting edit with no behavior impact.
- The user asks for a read-only explanation and no implementation.

#### Validation

- The output names the planning artifact reviewed.
- The change type is explicit.
- New requirements include updated planning/BDD/contract/test docs before code.
- Bug fixes include expected vs actual behavior and regression tests.
- Blocker questions are resolved or explicitly scoped out.

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

- Updated `process/README.md` with Change Intake Gate.
- Updated `templates/initial-development-docs.md` with Change Intake fields.
- Updated `CHECKLIST.md` with change intake checks.
- Updated `WORKFLOW.md` with before-code intake.
- Updated `SKILL.md` quick start and output style.
- Updated `README.md` goals, classification notes, and linked update examples.
- Updated `feedback_history/README.md` and `feedback_history/common/README.md`.
