# app-development-guidance feedback history

Each file in this directory is one reusable app development guidance lesson.

Follow [`shared-rules/feedback-lessons.md`](../../../shared-rules/feedback-lessons.md). Do not duplicate shared rules in every lesson.

## Categories

This skill already separates promoted guidance into `controls/`, `platforms/`, `languages/`, `implementation/`, and `checklists/`. New feedback lessons should use the matching `feedback_history/<category>/` folder when the primary category is clear; cross-cutting lessons can use `feedback_history/common/`.

| Category | Purpose |
| --- | --- |
| [`common/`](common/) | Cross-cutting lessons that affect process, implementation, checklists, templates, or multiple guidance areas. |
| [`controls/`](controls/) | Lessons whose primary promotion target is a reusable security control. |

## Historical Flat Index

| File | Status | Topic | Summary |
| --- | --- | --- | --- |
| `common/2026-05-05_194400-contract-first-development-flow.md` | promoted | Contract-first development flow | Start from product brief, split bounded contexts, write BDD, define Domain, Architecture, API/Interface, and Error Handling Contracts, then build and test provider/consumer sides. |
| `common/2026-05-05_200500-existing-project-doc-backfill-bdd-required.md` | promoted | Existing project doc backfill requires complete BDD | Existing projects must audit and backfill missing development docs; Product Brief gaps can be unknown, but BDD must be completed from observable implementation evidence. |
| `common/2026-05-05_201000-missing-requirements-block-development.md` | promoted | Missing requirements block development | Missing behavior, contract, error, security, storage, ownership, or test requirements must be asked and resolved before development continues. |
| `common/2026-05-06_081600-change-intake-before-code.md` | promoted | Change intake before code | Review planning artifacts and classify work as new requirement, bug, refactor, hardening, or docs-only before code; new requirements must update planning docs first. |
| `common/2026-05-06_082000-separate-regression-from-new-code-validation.md` | promoted | Separate regression from new code validation | Distinguish legacy regression coverage from new-code validation; use BDD/TDD, changed-code coverage, mutation/property/contract/database tests, and human review. |
| `common/2026-05-06_083000-embedded-hardware-product-flow.md` | promoted | Embedded hardware product flow | Treat embedded/firmware/hardware-backed products as app/product development flow with additional datasheet, protocol, hardware context, host fixture, and hardware-in-loop contracts. |
| `controls/2026-05-01_142100-client-encrypted-header-not-boundary.md` | promoted | Client encrypted header is not a security boundary | Client-side encrypted or signed headers are recoverable from shipped apps; backend authorization, replay protection, token hygiene, and monitoring must provide the real boundary. |
