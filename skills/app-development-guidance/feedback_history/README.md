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
| `common/2026-05-06_083200-implemented-first-contract-governance.md` | promoted | Implemented-first contract governance | Backfill document precedence, traceability, BDD closure, generated-client flow, vendor boundaries, and tooling kernel/adapter rules for projects implemented before specs. |
| `common/2026-05-06_103200-product-brief-validation-gate.md` | promoted | Product Brief validation gate | Validate Product Brief claims or mark assumptions, blockers, scoped-out items, and invalidated claims before using them for BDD, contracts, implementation, or tests. |
| `common/2026-05-07_122800-performance-test-release-gate.md` | promoted | Performance test release gate | Performance-sensitive changes need budgets and load, stress, spike, soak, or smoke-size evidence before release. |
| `common/2026-05-07_152100-private-live-adapter-smoke-gate.md` | candidate | Private live adapter smoke gate | Live-facing SDK/client work needs a private adapter smoke checklist when public analysis docs redact host/base endpoint, service, session, signing, or decrypt material. |
| `common/2026-05-07_153600-schema-derived-synthetic-fixtures.md` | candidate | Schema-derived synthetic fixtures | Use schema-compatible placeholder fixtures to unblock parser/BDD work when raw live payloads are private or non-committable. |
| `common/2026-05-07_154400-analysis-sdk-contract-drift-gate.md` | candidate | Analysis-to-SDK contract drift gate | Audit downstream SDK contracts, live-test gates, README wording, and fixture provenance whenever APK analysis reclassifies readiness boundaries. |
| `common/2026-05-07_160200-media-metadata-private-decrypt-boundary.md` | candidate | Media metadata private decrypt boundary | Preserve safe media metadata in SDK core while signed URLs, key unwrap, and decrypt/download stay in private adapters or media modules. |
| `controls/2026-05-01_142100-client-encrypted-header-not-boundary.md` | promoted | Client encrypted header is not a security boundary | Client-side encrypted or signed headers are recoverable from shipped apps; backend authorization, replay protection, token hygiene, and monitoring must provide the real boundary. |
