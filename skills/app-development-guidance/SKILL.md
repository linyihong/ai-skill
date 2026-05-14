---
name: app-development-guidance
description: Convert authorized app, API, embedded, firmware, and hardware-product findings into practical development guidance. Use when reverse-engineering lessons, product briefs, or hardware/firmware specs should become buildable app/API/SDK/firmware patterns, contracts, BDD/TDD plans, validation tests, release checks, and review checklists across mobile, web, backend, embedded, and future product types.
---

# App Development Guidance

Use this skill when APK analysis, app/API review, embedded/firmware review, or product development work reveals a reusable lesson that can improve future apps or hardware-backed products. The goal is to translate observed behavior, attack paths, failure modes, implementation tricks, hardware constraints, and protocol contracts into practical development requirements, buildable patterns, checklists, and validation tests.

**Shared policy:** read [`shared-rules` index](../../enforcement/README.md), [`dependency-reading.md`](../../enforcement/dependency-reading.md), [`feedback-lessons.md`](../../enforcement/feedback-lessons.md), [`reusable-guidance-boundary.md`](../../enforcement/reusable-guidance-boundary.md), apply [`neutral-language.md`](../../enforcement/neutral-language.md) when naming or summarizing docs, and apply [`goal-action-validation.md`](../../enforcement/goal-action-validation.md) so important conclusions include a goal, action, validation, or reference source. If this skill or a related rule/template/lesson has changed, or the user asks to reload it, create a dependency read ledger (required files, files read, missing files marked `not applicable`, blocked items, validation), read the dependent docs before concluding, and close the Ai-skill writeback transaction before returning to project work. Lessons in `feedback_history/` should reference those files, not duplicate shared rules.

**Cross-skill references:** follow [`cross-skill-references.md`](../../enforcement/cross-skill-references.md). When another skill hands off sanitized analysis artifacts, this skill consumes the development-relevant contract, risk, validation, and open-question details without copying the source skill's analysis workflow.

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
- Adding performance testing to development and release gates when changes can affect latency, throughput, resource usage, startup, background work, database access, batching, caching, concurrency, or external-call volume.
- Backfilling missing development documents for existing, already implemented projects; Product Brief gaps can be marked unknown, but BDD behavior must be completed from observed implementation evidence.
- Reviewing app/API design for replay resistance, token safety, transport security, local storage, logging, and release hardening.
- Creating PR/release checklists for mobile, web, backend/API, and future app types.
- Deciding what should be validated by tests, fixtures, runtime checks, or server-side controls.

## Out Of Scope

- Breaking into apps or services without authorization.
- Storing target-specific hosts, secrets, endpoints, tokens, real user data, or private business conclusions in this reusable skill.
- Treating client-side hardening as a replacement for server-side authorization, rate limits, fraud controls, or audit logs.

## Quick Start（Routing）

Routing summary:
1. Identify source → 2. Validate planning docs → 3. Classify request → 4. New req: planning docs first → 5. Bug fix: before/after code rules → 6. Test strategy before code → 7. Product brief → `process/` → 8. Existing project: backfill → 9. Blocker questions → 10. Risk statement → 11. Choose control layer → 12. Validation method → 13. Classify guidance → 14. Linked updates → 15. Immature lesson → `feedback_history/`.

## Default Workflow

### 新分層路徑（優先讀取）

| 用途 | 路徑 |
|------|------|
| 開發流程（Product Brief → BDD → 實作 → 驗證） | [`workflow/software-delivery/development-process.md`](../../workflow/software-delivery/development-process.md) |
| 執行流程（Change Intake → Docs-First BDD → SDK Defect → Performance Gate） | [`workflow/software-delivery/execution-flow.md`](../../workflow/software-delivery/execution-flow.md) |
| 分析方法（Risk Translation、Controls Catalog、Platforms/Languages Catalog） | [`analysis/development-guidance/README.md`](../../analysis/development-guidance/README.md) |
| 工程智慧（Heuristics、Anti-patterns） | [`intelligence/engineering/development/README.md`](../../intelligence/engineering/development/README.md) |
| 產出格式與品質門檻 | [`workflow/software-delivery/artifact-gates.md`](../../workflow/software-delivery/artifact-gates.md) |
| Review Checklist | [`workflow/software-delivery/review-checklist.md`](../../workflow/software-delivery/review-checklist.md) |

### 舊路徑（保留向後相容）

| 用途 | 路徑 |
|------|------|
| WORKFLOW.md（舊執行流程） | [`WORKFLOW.md`](WORKFLOW.md) |
| process/（舊開發流程） | [`process/`](process/) |
| CHECKLIST.md + checklists/ | [`CHECKLIST.md`](CHECKLIST.md)、[`checklists/`](checklists/) |
| DOCUMENTATION.md（舊產出格式） | [`DOCUMENTATION.md`](DOCUMENTATION.md) |
| controls/ + platforms/ + languages/ + implementation/ | [`controls/`](controls/)、[`platforms/`](platforms/)、[`languages/`](languages/)、[`implementation/`](implementation/) |

## Output Style & Artifact Gates

See [`workflow/software-delivery/artifact-gates.md`](../../workflow/software-delivery/artifact-gates.md) for output format and quality gates.

## Feedback Loop

See [`enforcement/feedback-lessons.md`](../../enforcement/feedback-lessons.md) for the feedback lesson template and workflow. See [`feedback/`](../../feedback/) for the feedback promotion pipeline.

**Cross-skill link:** if the lesson came from APK analysis, keep analysis mechanics in [`apk-analysis`](../apk-analysis/) and put development guidance, implementation patterns, and validation checklists here.

When consuming a Feature Reconstruction Handoff from [`apk-analysis`](../apk-analysis/), this skill owns the development conversion: BDD, Domain Model Contract, Architecture Contract, API / Interface Contract, Error Handling Contract, implementation slices, fixtures, tests, and review checklists. If the requested output is an app-related tool, SDK, client, mock, fixture-driven implementation, contract test, or rebuilt feature, apply this skill before implementation planning and surface blocker questions. Keep raw APK evidence, target-specific hosts, tokens, accounts, and private product conclusions in the project analysis docs.
