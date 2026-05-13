> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/app-development-guidance/development-process.md`](../../../../workflow/app-development-guidance/development-process.md)

### 2026-05-05 - Existing project doc backfill requires complete BDD

Status: promoted

#### One-line Summary

When opening app-development-guidance on an already implemented project, audit and backfill missing development documents; missing Product Brief details may remain unknown, but BDD behavior must be completed from observable evidence.

#### Human Explanation

Some projects are already fully developed before this skill is introduced. In that case, the agent should not treat missing planning documents as a reason to skip documentation. The project already contains behavior in UI flows, API contracts, code paths, tests, fixtures, logs, and release behavior.

Original product intent can be hard to recover. It is acceptable to mark Product Brief fields as `unknown` or `open question` when evidence does not support them. BDD is different: for implemented behavior, the agent must recover complete behavior scenarios from the product itself, including happy paths, failures, permissions, empty states, edge cases, and cross-context flows.

#### Trigger

- The user opens `app-development-guidance` on an existing project.
- A project has working implementation but missing Product Brief, BDD, Domain Model, Architecture, API, Error Handling, or Test Plan docs.
- The user asks to "補齊文件", "檢查缺失文件", or make an old project align with the development process.

#### Evidence

- Tool: source code, tests, API schemas, UI behavior, logs, fixtures, screenshots, release notes, manual verification.
- Sanitized excerpt: `existing implementation -> document gap audit -> BDD recovery -> contract backfill -> test gap mapping`.
- Evidence path: project repository docs and implementation; reusable skill stores only the generalized process.

#### Generalized Lesson

For existing projects:

1. Inventory current documents and mark each as `exists`, `partial`, `missing`, or `unknown`.
2. Backfill Product Brief only from evidence; mark unavailable intent as `unknown` / `open question`.
3. Complete BDD from observable implementation behavior.
4. Backfill Domain Model, Architecture, API / Interface, Error Handling, and Test Plan from BDD plus implementation evidence.
5. Map every critical BDD scenario to existing tests or required test gaps.

#### Agent Action

- Do not stop because the Product Brief is missing.
- Do not invent product intent to fill a blank Product Brief.
- Do complete BDD for implemented critical behavior using available evidence.
- Do record the source of each recovered scenario or contract.
- Do list remaining unknowns separately from observed behavior.

#### Applies When

- The project already has implementation, tests, APIs, UI flows, fixtures, or release behavior.
- The request is documentation recovery, process alignment, or development document completion.

#### Does Not Apply When

- The project is only an idea with no implementation or observed behavior.
- The user asks only for a narrow security checklist unrelated to development docs.
- Evidence is unavailable; then document the gap and ask for access or examples.

#### Validation

- A document gap table exists.
- Product Brief unknowns are explicitly marked.
- BDD covers implemented critical happy paths, failure paths, permissions, empty states, edge cases, and cross-context flows.
- Contracts and test plan reference recovered BDD scenarios.
- Test gaps are listed for critical uncovered behavior.

#### Promotion Target

- `process/README.md`
- `templates/initial-development-docs.md`
- `CHECKLIST.md`
- `WORKFLOW.md`
- `SKILL.md`
- `README.md`

#### Required Linked Updates

- Updated `process/README.md` with Existing Project Documentation Backfill.
- Updated `templates/initial-development-docs.md` with a document gap audit table.
- Updated `CHECKLIST.md` with backfill checklist items.
- Updated `WORKFLOW.md` with document gap audit behavior.
- Updated `SKILL.md` quick start and output style.
- Updated `README.md` goals, classification notes, and linked update examples.
- Updated `feedback_history/README.md` and `feedback_history/common/README.md`.
