---
name: app-development-guidance
description: Convert authorized app, API, embedded, firmware, and hardware-product findings into practical development guidance. Use when reverse-engineering lessons, product briefs, or hardware/firmware specs should become buildable app/API/SDK/firmware patterns, contracts, BDD/TDD plans, validation tests, release checks, and review checklists across mobile, web, backend, embedded, and future product types.
---

# App Development Guidance

Use this skill when APK analysis, app/API review, embedded/firmware review, or product development work reveals a reusable lesson that can improve future apps or hardware-backed products. The goal is to translate observed behavior, attack paths, failure modes, implementation tricks, hardware constraints, and protocol contracts into practical development requirements, buildable patterns, checklists, and validation tests.

**Shared policy:** read [`shared-rules` index](../../shared-rules/README.md), [`dependency-reading.md`](../../shared-rules/dependency-reading.md), [`feedback-lessons.md`](../../shared-rules/feedback-lessons.md), apply [`neutral-language.md`](../../shared-rules/neutral-language.md) when naming or summarizing docs, and apply [`goal-action-validation.md`](../../shared-rules/goal-action-validation.md) so important conclusions include a goal, action, validation, or reference source. If this skill or a related rule/template/lesson has changed, read the dependent docs before concluding, and close the Ai-skill writeback transaction before returning to project work. Lessons in `feedback_history/` should reference those files, not duplicate shared rules.

**Cross-skill references:** follow [`cross-skill-references.md`](../../shared-rules/cross-skill-references.md). When another skill hands off sanitized analysis artifacts, this skill consumes the development-relevant contract, risk, validation, and open-question details without copying the source skill's analysis workflow.

## When To Use

- Turning APK analysis findings into practical development guidance.
- Turning observed app/API techniques into buildable patterns for your own app.
- Turning APK analysis documents into app-related tools, SDKs, clients, mocks, fixture-driven implementations, contract tests, or rebuilt features.
- Turning product briefs into BDD, Domain Model Contracts, API Contracts, implementation slices, and tests.
- Validating product briefs before using them as implementation input: goals, users, scope, non-goals, assumptions, success criteria, constraints, dependencies, risks, and blocker questions.
- Turning hardware, sensor, protocol, firmware, or embedded product specs into datasheet/protocol contracts, hardware context contracts, firmware BDD, driver/service/application boundaries, host/target tests, and bring-up validation.
- Backfilling implemented-first projects with document precedence, product/rule traceability, BDD-to-test closure, minimum doc-sync matrices, generated-client checks, and explicit canceled/deferred/out-of-scope decisions.
- Turning OpenAPI/schema contracts into typed clients, SDKs, mocks, fixtures, and provider/consumer tests without hand-copied endpoints.
- Turning vendor or third-party API integrations into sanitized integration excerpts, fixture tests, live-test gates, webhook/idempotency checks, and secret-safe docs.
- Turning IDE extensions, CLIs, linters, static analyzers, code generators, and internal tools into rule catalogs, pure kernels, adapters, diagnostics/commands, fixtures, and integration tests.
- Classifying app changes as new requirements, bug fixes, refactors, or hardening work before code; new requirements must update planning docs before implementation.
- Separating legacy regression protection from new-code validation, including BDD/TDD, changed-code coverage, mutation tests, property-based tests, contract tests, database regression tests, and human review.
- Backfilling missing development documents for existing, already implemented projects; Product Brief gaps can be marked unknown, but BDD behavior must be completed from observed implementation evidence.
- Reviewing app/API design for replay resistance, token safety, transport security, local storage, logging, and release hardening.
- Creating PR/release checklists for mobile, web, backend/API, and future app types.
- Deciding what should be validated by tests, fixtures, runtime checks, or server-side controls.

## Out Of Scope

- Breaking into apps or services without authorization.
- Storing target-specific hosts, secrets, endpoints, tokens, real user data, or private business conclusions in this reusable skill.
- Treating client-side hardening as a replacement for server-side authorization, rate limits, fraud controls, or audit logs.

## Quick Start

1. Identify the source: product brief, observed weakness, reverse-engineering lesson, or app/API design question.
2. Before code changes, inspect and validate the project's 企劃書, product brief, planning docs, issue, ticket, PRD, design note, BDD, API contract, or equivalent artifact. Major Product Brief claims must be marked `validated`, `assumption`, `open question`, `scoped out`, or `invalidated` before they drive implementation.
3. Classify the request as new requirement, bug fix, refactor, hardening, or documentation-only.
4. If it is a new requirement or behavior change, update or create planning docs first: change brief, BDD scenarios, impacted Domain Model Contract, Architecture Contract, API / Interface Contract, Error Handling Contract, implementation slices, and tests. For embedded/hardware work, also update datasheet/protocol references, hardware context, driver/service/application ownership, fixture or hardware-in-loop validation, and bring-up notes. Do not start code until blocker questions are resolved.
5. If it is a bug fix, confirm expected vs actual behavior, reproduction/evidence, affected or missing BDD scenario, impacted contract/error handling, and regression test plan before code.
6. Define the test strategy before production code: distinguish existing-regression coverage from changed/new-code validation; prefer BDD first, then failing unit/contract/property/integration tests for new behavior before implementation.
7. If starting from a product brief, use [`process/`](process/) to draft or discuss the initial development docs: Product Brief validation, Bounded Contexts, BDD behavior, Domain Model Contract, Architecture Contract, API / Interface Contract, Error Handling Contract, implementation slices, and tests.
8. If opening this skill on an existing implemented project, audit missing documents and backfill them. Missing Product Brief fields may be marked `unknown` / `open question`, but BDD behavior must be completed from UI, API, code, tests, logs, fixtures, or observed behavior. Recover document precedence, traceability, BDD validation status, generated-client flow, vendor excerpts, and canceled/out-of-scope decisions.
9. If anything required for behavior, domain invariants, API/interface shape, error handling, security, storage, tests, ownership, document precedence, generated clients, vendor integration, or tool diagnostics is missing, ask the user or request evidence before continuing. Do not proceed with development while blocker questions remain unresolved.
10. Convert analysis findings into a developer-facing risk statement or implementation opportunity.
11. Choose the control layer:
   - API/server contract.
   - App runtime behavior.
   - Full-stack schema/codegen or provider/consumer contract.
   - Tooling, IDE extension, CLI, linter, or static-analysis kernel/adapter behavior.
   - Vendor or third-party API integration behavior.
   - Embedded firmware, hardware context, sensor/protocol driver, or board bring-up behavior.
   - Build/release configuration.
   - Monitoring or fraud signal.
12. Add a concrete validation method: brief evidence check, unit test, BDD scenario, API contract test, integration test, release checklist item, fixture, mutation/property check, or manual review step.
13. Classify the guidance:
   - Core security control: `controls/`.
   - Platform/app type detail: `platforms/`.
   - Language/runtime-specific trap: `languages/`.
   - Concrete implementation pattern: `implementation/`.
   - Product-to-contract development flow: `process/`.
   - Repeatable review step: `checklists/`.
14. Apply required linked updates from [`shared-rules/linked-updates.md`](../../shared-rules/linked-updates.md): when a process, control, platform, language, checklist, implementation pattern, or template is affected, update or explicitly verify the related files in the same change.
15. If the lesson is reusable but not yet mature, add it to the matching `feedback_history/<category>/` folder first; use `feedback_history/common/` for cross-cutting lessons. Promote it into the structured folders when validated.

## Default Workflow

Read [WORKFLOW.md](WORKFLOW.md) to translate analysis evidence into development requirements and implementation guidance.

Use [`process/`](process/) when starting from a product brief, planning a feature from BDD through Domain Model, Architecture, API / Interface, Error Handling, implementation, and tests, or backfilling missing documents for an existing implemented project.

Use [`templates/`](templates/) to choose a copyable template. Use [`templates/initial-development-docs.md`](templates/initial-development-docs.md) when the user wants the first draft of the development documents or when requirements are still being clarified through discussion.

Use [CHECKLIST.md](CHECKLIST.md) for the checklist index and [`checklists/`](checklists/) for focused design, PR, release, and API reviews.

Use [DOCUMENTATION.md](DOCUMENTATION.md) when writing reusable development guidance or project-specific hardening notes.

Use [`controls/`](controls/) as the primary home for cross-platform guidance, then link to [`platforms/`](platforms/) or [`languages/`](languages/) for implementation details.

Use [`implementation/`](implementation/) when the user asks how to build or implement a hardening control.

## Output Style

When producing development guidance, include:

- Goal, action, and validation or reference source for each important work unit or conclusion.
- Product Brief validation status when planning docs are used as implementation input.
- Observed risk or failure mode.
- Why it matters for app development.
- Change classification: new requirement, bug fix, refactor, hardening, or docs-only; include planning artifact reviewed.
- Test strategy: existing behavior guarded by regression tests, new/changed code validated by BDD/TDD and changed-code tests; mention mutation/property/contract/database tests when relevant.
- Document precedence and traceability for implemented-first projects, including BDD-to-test status.
- Generated-client, SDK, fixture, or schema sync status when API/interface contracts are involved.
- Vendor/third-party integration boundaries, live-test gates, and secret-safe documentation when external providers are involved.
- Missing questions that must be answered before implementation, if starting from a product brief.
- Existing-project documentation gaps and how they were backfilled; BDD must be complete when observable behavior exists.
- Blocker questions for any missing behavior, contract, error, security, storage, ownership, or test requirement; do not continue development until these are answered or explicitly scoped out.
- Draft documents or document sections that are ready to copy into the project repository.
- Recommended control and the layer that owns it.
- Implementation path or linked implementation doc, when applicable.
- Validation method.
- Required linked updates, if the change affects multiple folders.
- What not to overclaim.

## Feedback Loop

If a reusable app development lesson emerges:

1. Create `feedback_history/<category>/YYYY-MM-DD_HHMMSS-<slug>.md` using [shared-rules/feedback-lessons.md](../../shared-rules/feedback-lessons.md); use `common/` when no single category owns it.
2. Generalize the lesson so it is not tied to one APK or company.
3. Include evidence and validation criteria, but redact secrets and target-specific details.
4. Promote validated guidance into `controls/`, `platforms/`, `languages/`, `implementation/`, `checklists/`, `WORKFLOW.md`, `CHECKLIST.md`, or `DOCUMENTATION.md` as appropriate.
5. If the promotion creates linked updates, those updates are mandatory; do not leave related docs stale.

**Cross-skill link:** if the lesson came from APK analysis, keep analysis mechanics in [`apk-analysis`](../apk-analysis/) and put development guidance, implementation patterns, and validation checklists here.

When consuming a Feature Reconstruction Handoff from [`apk-analysis`](../apk-analysis/), this skill owns the development conversion: BDD, Domain Model Contract, Architecture Contract, API / Interface Contract, Error Handling Contract, implementation slices, fixtures, tests, and review checklists. If the requested output is an app-related tool, SDK, client, mock, fixture-driven implementation, contract test, or rebuilt feature, apply this skill before implementation planning and surface blocker questions. Keep raw APK evidence, target-specific hosts, tokens, accounts, and private product conclusions in the project analysis docs.
