> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[feedback-lessons](../../../../shared-rules/feedback-lessons.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/software-delivery/development-process.md`](../../../../workflow/software-delivery/development-process.md)

### 2026-05-06 - Product Brief validation gate

Status: promoted

#### One-line Summary

Product Brief / 企劃書 claims must be validated, labeled as assumptions, asked as open questions, scoped out, or revised before they drive BDD, contracts, implementation, or tests.

#### Goal / Action / Validation

| Field | Content |
| --- | --- |
| Goal | Prevent agents and engineers from treating an unverified Product Brief as implementation truth. |
| Action | Add a Product Brief validation gate before BDD, contracts, implementation slices, and tests. |
| Validation | Each major brief claim has evidence, explicit decision, validation plan, blocker status, or scoped-out status. |

#### Human Explanation

A Product Brief is a planning artifact, not proof by itself. It can contain unclear goals, broad scope, hidden assumptions, missing non-goals, untested success criteria, or dependencies that block implementation.

Before using the brief as the source for development, validate goals, users, scope, non-goals, assumptions, success criteria, constraints, dependencies, and risks. If a claim affects behavior, contracts, security, storage, ownership, tests, schedule, or release gates and cannot be validated, it becomes a blocker question.

#### Trigger

- The user asks to start from a 企劃書, Product Brief, PRD, issue, ticket, design note, or planning document.
- A new requirement or behavior change depends on product intent.
- Existing implementation docs are being backfilled and the original product intent is incomplete.

#### Evidence

- Tool: Product Brief, PRD, user request, stakeholder decision, BDD, UI/API behavior, source code, tests, logs, fixtures, architecture contract, policy, vendor or platform docs.
- Sanitized excerpt: `brief claim -> validation status -> evidence or blocker -> downstream BDD/contract/test`.
- Evidence path: project repository planning docs; reusable skill stores only the generalized workflow.

#### Generalized Lesson

For every major brief claim, record one status:

| Status | Meaning |
| --- | --- |
| `validated` | Supported by evidence or explicit decision. |
| `assumption` | Plausible but not proven; needs owner and validation plan. |
| `open question` | Needed before impacted implementation can proceed. |
| `scoped out` | Explicitly not part of current work. |
| `invalidated` | Contradicted by evidence and must be revised. |

#### Agent Action

- Validate the brief before converting it into BDD, contracts, implementation slices, or tests.
- Do not invent missing product intent.
- Ask blocker questions for unvalidated claims that affect behavior, contracts, risks, ownership, tests, schedule, or release gates.
- For implemented-first projects, validate backfilled brief claims against observable evidence and keep unrecoverable intent as `unknown` or `open question`.

#### Applies When

- Starting a new feature or project from planning docs.
- Reviewing whether a feature is ready for implementation.
- Backfilling documentation for an already implemented project.
- Explaining why development is blocked by unclear product intent.

#### Does Not Apply When

- The task is a typo, formatting edit, or docs-only index update with no behavior or contract impact.
- The user explicitly asks for brainstorming only; in that case, label conclusions as assumptions or reference-source judgments.

#### Validation

- `process/README.md` contains the Product Brief Validation Gate.
- `WORKFLOW.md` requires brief validation during change intake.
- `CHECKLIST.md` has review checks for validated brief claims and blockers.
- `templates/initial-development-docs.md` captures validation status and evidence for brief fields.
- `SKILL.md` and `README.md` mention brief validation as part of app-development-guidance.

#### Promotion Target

- `process/README.md`
- `WORKFLOW.md`
- `CHECKLIST.md`
- `templates/initial-development-docs.md`
- `templates/README.md`
- `DOCUMENTATION.md`
- `SKILL.md`
- `README.md`
- `shared-rules/linked-updates.md`
- `feedback_history/README.md`
- `feedback_history/common/README.md`

#### Required Linked Updates

- Updated `process/README.md` with Product Brief Validation Gate.
- Updated `WORKFLOW.md` with validated planning artifact intake.
- Updated `CHECKLIST.md` with Product Brief validation review checks.
- Updated `templates/initial-development-docs.md` with validation status and evidence columns.
- Updated `templates/README.md`, `DOCUMENTATION.md`, `SKILL.md`, and `README.md`.
- Updated `shared-rules/linked-updates.md`.
- Updated feedback indexes.
