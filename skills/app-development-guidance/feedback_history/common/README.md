# Common Feedback Lessons

Cross-cutting lessons for app development guidance that affect multiple categories such as process, implementation, checklists, templates, or controls.

| File | Status | Topic | Summary |
| --- | --- | --- | --- |
| `2026-05-05_194400-contract-first-development-flow.md` | promoted | Contract-first development flow | Start from product brief, split bounded contexts, write BDD, define Domain, Architecture, API/Interface, and Error Handling Contracts, then build and test provider/consumer sides. |
| `2026-05-05_200500-existing-project-doc-backfill-bdd-required.md` | promoted | Existing project doc backfill requires complete BDD | Existing projects must audit and backfill missing development docs; Product Brief gaps can be unknown, but BDD must be completed from observable implementation evidence. |
| `2026-05-05_201000-missing-requirements-block-development.md` | promoted | Missing requirements block development | Missing behavior, contract, error, security, storage, ownership, or test requirements must be asked and resolved before development continues. |
| `2026-05-06_081600-change-intake-before-code.md` | promoted | Change intake before code | Review planning artifacts and classify work as new requirement, bug, refactor, hardening, or docs-only before code; new requirements must update planning docs first. |
| `2026-05-06_082000-separate-regression-from-new-code-validation.md` | promoted | Separate regression from new code validation | Distinguish legacy regression coverage from new-code validation; use BDD/TDD, changed-code coverage, mutation/property/contract/database tests, and human review. |
| `2026-05-06_083000-embedded-hardware-product-flow.md` | promoted | Embedded hardware product flow | Treat embedded/firmware/hardware-backed products as app/product development flow with additional datasheet, protocol, hardware context, host fixture, and hardware-in-loop contracts. |
