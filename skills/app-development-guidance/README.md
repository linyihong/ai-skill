# App Development Guidance Skill

This skill turns authorized app, API, embedded, firmware, and hardware-product findings into reusable development guidance for better, safer products.

It complements [`apk-analysis`](../apk-analysis/):

- `apk-analysis` asks: "How does this APK behave, communicate, protect, or fail?"
- `app-development-guidance` asks: "What should we build, test, or avoid in our own apps, SDKs, firmware, APIs, or hardware-backed products because of that lesson?"

## Goals

Capture practical guidance for:

- App implementation patterns learned from analysis.
- App/API and transport security.
- Embedded, firmware, sensor, protocol, and hardware context development flow.
- Product brief to contract-first development flow.
- Product Brief validation before implementation: goals, users, scope, non-goals, assumptions, success criteria, constraints, dependencies, risks, and blocker questions.
- Implemented-first project governance: document precedence, traceability, BDD closure, minimum doc-sync matrices, and **same-session closure** when code changes observable behavior (see `WORKFLOW.md`, `CHECKLIST.md`).
- Change intake before code: classify new requirements, bug fixes, refactors, hardening, and docs-only work from planning artifacts.
- Test strategy that separates legacy regression protection from new-code validation with BDD/TDD, changed-code coverage, mutation/property/contract/database tests, and human review.
- Performance test strategy that treats latency, throughput, error rate, resource usage, load, stress, spike, and soak evidence as release inputs when the change can affect capacity or responsiveness.
- OpenAPI/schema/codegen flows for typed clients, SDKs, mocks, fixtures, and provider/consumer tests.
- Vendor and third-party API integration docs, fixtures, live-test gates, webhook/idempotency checks, and secret-safe boundaries.
- Developer tooling patterns for IDE extensions, CLIs, linters, static analyzers, rule catalogs, kernels, adapters, diagnostics, and fixtures.
- Existing project documentation gap backfill, with complete BDD recovery from implemented behavior.
- Token, session, and replay handling.
- Local storage and secret exposure risks.
- Mobile, web, backend/API, and release hardening.
- Hardware bring-up, host/target tests, protocol fixtures, and firmware release gates.
- Sensitive logging and telemetry hygiene.
- Anti-tamper signals without false confidence.
- Developer checklists, buildable patterns, and validation tests.

## What Belongs Here

- Reusable development patterns learned from APK analysis, app/API review, or product development.
- Reusable development patterns for embedded or hardware-backed products, including datasheet/protocol contracts, hardware context, driver/service/application layering, and validation.
- Rules for validating planning artifacts before code and deciding whether work is a new requirement or bug fix.
- Rules for validating new or AI-generated code beyond total project coverage.
- Documentation backfill rules for already implemented projects, especially when original planning documents are missing.
- Rules for document precedence, traceability, BDD execution closure, and explicit canceled/deferred/out-of-scope decisions.
- Reusable implementation patterns for generated clients, vendor integrations, and developer tools.
- High-level controls that can be implemented by app, API, backend, or release engineering teams.
- Checklists that help prevent repeat mistakes.
- Guidance that clearly names validation steps and limitations.
- Generalized incident lessons that explain the reusable cause, decision rule, and validation method without preserving the triggering project's concrete details; follow [`reusable-guidance-boundary.md`](../../enforcement/reusable-guidance-boundary.md).

## What Does Not Belong Here

- Target-specific API hosts, endpoints, tokens, device identifiers, or private response schemas.
- Raw request/response data.
- One-off product conclusions that do not generalize.
- Project names, module names, local paths, sample IDs, class names, live environment quirks, execution results, or BDD/test file names from a single incident.
- Advice that relies only on client-side secrecy.

## Files

| File | Purpose |
| --- | --- |
| `SKILL.md` | Agent entry point. |
| `WORKFLOW.md` | Translate observations into development requirements and validation steps. |
| `process/` | Product brief to BDD, Domain Model, Architecture, API / Interface, Error Handling, implementation, and test flow. |
| `CHECKLIST.md` | Development, PR, and release checklist. |
| `DOCUMENTATION.md` | How to document reusable development guidance. |
| `controls/` | Cross-platform security controls; use this as the primary home for core guidance. |
| `platforms/` | Platform or product-type implementation guidance such as mobile, web, backend API, and embedded/firmware. |
| `languages/` | Language/runtime-specific pitfalls only. |
| `checklists/` | Focused design, PR, release, API, contract governance, embedded, and tooling review checklists. |
| `implementation/` | Concrete implementation patterns for app, backend, and platform teams. |
| `templates/` | Copyable templates for initial development docs, guidance notes, and lightweight threat models; start with [`templates/README.md`](templates/README.md). |
| `FEEDBACK.md` | Short entry pointing to shared feedback rules. |
| `feedback_history/` | Reusable lessons, categorized by the same primary folders when applicable. |

## Classification Rules

When adding new guidance:

1. Put the core security property in [`controls/`](controls/) first.
2. Add platform details in [`platforms/`](platforms/) only when the implementation differs by app type, OS, firmware runtime, board, sensor, or hardware interface.
3. Add language notes in [`languages/`](languages/) only for language/runtime-specific traps.
4. Put concrete implementation patterns in [`implementation/`](implementation/) when engineers need buildable steps.
5. Put development process and contract-first flow guidance in [`process/`](process/).
6. Put repeatable review steps in [`checklists/`](checklists/).
7. Put draft reusable lessons in the matching [`feedback_history/<category>/`](feedback_history/) folder before promoting them into the structured folders; use `feedback_history/common/` for cross-cutting lessons.
8. When the source is a project incident, split the generalized method into this skill and keep concrete reproduction evidence in the project repository, per [`reusable-guidance-boundary.md`](../../enforcement/reusable-guidance-boundary.md).

For existing implemented projects, use [`process/`](process/) to audit and backfill missing documents. Product Brief gaps can remain `unknown`; BDD behavior must be completed from observable behavior and implementation evidence.

For implemented-first projects, also recover document precedence, stable IDs, traceability links, BDD validation status, minimum doc-sync matrix, generated-client flow, vendor integration excerpts, and explicit canceled/deferred/out-of-scope items.

Before code work, use [`process/`](process/) to validate the planning artifact and classify the change. Product Brief claims that affect behavior, contracts, risk, ownership, tests, schedule, or release gates must be supported by evidence, explicit decision, validation plan, or blocker status before implementation starts. New requirements must update planning docs, BDD, contracts, implementation slices, and tests before implementation starts. Bug fixes must document expected vs actual behavior and regression tests first.

When planning tests, separate old behavior from new code. Existing tests guard regressions; new or AI-generated code needs BDD/TDD and changed-code validation. Use mutation, property-based, contract, database-backed, or performance tests when ordinary example tests do not prove the rule. For performance-sensitive changes, define a budget and choose load, stress, spike, soak, or smoke-size CI validation.

For embedded or hardware-backed products, use [`platforms/embedded/`](platforms/embedded/) and [`implementation/embedded/`](implementation/embedded/) to capture datasheet/protocol truth, electrical/hardware context, firmware architecture, driver/service/application ownership, host-side fixtures, target/hardware-in-loop checks, and bring-up evidence.

For full-stack API work, use [`implementation/backend/contract-codegen.md`](implementation/backend/contract-codegen.md) to keep OpenAPI/schema/source contracts, generated clients, SDKs, mocks, fixtures, and provider/consumer tests aligned. Use [`implementation/backend/vendor-integration.md`](implementation/backend/vendor-integration.md) for third-party integrations. Use [`implementation/tooling/`](implementation/tooling/) for IDE extensions, CLIs, linters, static analyzers, and code generators.

This keeps the skill readable as it grows across mobile, web, backend, and future app types.

## Required Linked Updates

Repo-wide rule: [`../../enforcement/linked-updates.md`](../../enforcement/linked-updates.md). Some changes must move together. If a change touches a control, platform, language, checklist, implementation pattern, or template, the related files **must** be updated or explicitly checked in the same change.

Examples:

- New replay-defense implementation -> update or verify `controls/api-transport.md`, `platforms/backend/api.md`, and `checklists/api-security-review.md`.
- New Flutter storage implementation -> update or verify `platforms/mobile/flutter.md`, `languages/dart.md`, `controls/local-storage.md`, and the relevant checklist.
- New review checklist item -> update or verify the matching `controls/` and `implementation/` docs.
- New contract-first process guidance -> update or verify `process/`, `checklists/`, `templates/`, and relevant implementation docs.
- New Product Brief validation rule -> update or verify `process/README.md`, `templates/initial-development-docs.md`, `CHECKLIST.md`, `WORKFLOW.md`, and `SKILL.md`.
- New initial planning template -> update or verify `templates/README.md`, `process/README.md`, `DOCUMENTATION.md`, and `CHECKLIST.md`.
- New existing-project backfill rule -> update or verify `process/README.md`, `templates/initial-development-docs.md`, `CHECKLIST.md`, and `WORKFLOW.md`.
- New implemented-first governance rule -> update or verify `process/README.md`, `templates/initial-development-docs.md`, `CHECKLIST.md`, `WORKFLOW.md`, `implementation/`, and `checklists/contract-governance-review.md`.
- New generated-client or vendor-integration rule -> update or verify `implementation/backend/`, `platforms/backend/api.md`, `checklists/api-security-review.md`, `process/README.md`, and `CHECKLIST.md`.
- New tooling/extension rule -> update or verify `implementation/tooling/`, `process/README.md`, `CHECKLIST.md`, and relevant templates.
- New change-intake rule -> update or verify `process/README.md`, `templates/initial-development-docs.md`, `CHECKLIST.md`, `WORKFLOW.md`, and `SKILL.md`.
- New test-strategy rule -> update or verify `process/README.md`, `templates/initial-development-docs.md`, `CHECKLIST.md`, `WORKFLOW.md`, and `SKILL.md`.
- New performance-test strategy rule -> update or verify `process/README.md`, `templates/initial-development-docs.md`, `templates/README.md`, `CHECKLIST.md`, `WORKFLOW.md`, `DOCUMENTATION.md`, and `SKILL.md`.
- New embedded/hardware flow -> update or verify `platforms/embedded/`, `implementation/embedded/`, `process/README.md`, `templates/initial-development-docs.md`, `CHECKLIST.md`, `WORKFLOW.md`, and `SKILL.md`.

Do not describe these updates as optional. If they are relevant, they are required.

## Minimum Useful Output

A good development guidance note should include:

- The observed risk.
- The development consequence.
- The recommended control.
- The owner layer: client, API, backend, build, or monitoring.
- The validation method.
- Limits and non-goals.
