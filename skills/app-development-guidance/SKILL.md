---
name: app-development-guidance
description: Convert authorized app analysis findings into practical app development guidance. Use when reverse-engineering lessons should become buildable app/API patterns, security controls, token and session handling, request signing, TLS and certificate pinning decisions, release checks, local storage guidance, logging hygiene, anti-tamper signals, and review checklists across mobile, web, backend, and future app types.
---

# App Development Guidance

Use this skill when APK analysis, app/API review, or app development work reveals a reusable lesson that can improve future apps. The goal is to translate observed app behavior, attack paths, failure modes, and implementation tricks into practical development requirements, buildable patterns, checklists, and validation tests.

**Shared policy:** read [`shared-rules` index](../../shared-rules/README.md) and [`feedback-lessons.md`](../../shared-rules/feedback-lessons.md). Lessons in `feedback_history/` should reference those files, not duplicate shared rules.

**Cross-skill references:** follow [`cross-skill-references.md`](../../shared-rules/cross-skill-references.md). When another skill hands off sanitized analysis artifacts, this skill consumes the development-relevant contract, risk, validation, and open-question details without copying the source skill's analysis workflow.

## When To Use

- Turning APK analysis findings into practical development guidance.
- Turning observed app/API techniques into buildable patterns for your own app.
- Turning product briefs into BDD, Domain Model Contracts, API Contracts, implementation slices, and tests.
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
2. If starting from a product brief, use [`process/`](process/) to draft or discuss the initial development docs: Bounded Contexts, BDD behavior, Domain Model Contract, Architecture Contract, API / Interface Contract, Error Handling Contract, implementation slices, and tests.
3. If opening this skill on an existing implemented project, audit missing documents and backfill them. Missing Product Brief fields may be marked `unknown` / `open question`, but BDD behavior must be completed from UI, API, code, tests, logs, fixtures, or observed behavior.
4. Convert analysis findings into a developer-facing risk statement or implementation opportunity.
5. Choose the control layer:
   - API/server contract.
   - App runtime behavior.
   - Build/release configuration.
   - Monitoring or fraud signal.
6. Add a concrete validation method: unit test, BDD scenario, API contract test, integration test, release checklist item, fixture, or manual review step.
7. Classify the guidance:
   - Core security control: `controls/`.
   - Platform/app type detail: `platforms/`.
   - Language/runtime-specific trap: `languages/`.
   - Concrete implementation pattern: `implementation/`.
   - Product-to-contract development flow: `process/`.
   - Repeatable review step: `checklists/`.
8. Apply required linked updates from [`shared-rules/linked-updates.md`](../../shared-rules/linked-updates.md): when a process, control, platform, language, checklist, implementation pattern, or template is affected, update or explicitly verify the related files in the same change.
9. If the lesson is reusable but not yet mature, add it to the matching `feedback_history/<category>/` folder first; use `feedback_history/common/` for cross-cutting lessons. Promote it into the structured folders when validated.

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

- Observed risk or failure mode.
- Why it matters for app development.
- Missing questions that must be answered before implementation, if starting from a product brief.
- Existing-project documentation gaps and how they were backfilled; BDD must be complete when observable behavior exists.
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

When consuming a Feature Reconstruction Handoff from [`apk-analysis`](../apk-analysis/), this skill owns the development conversion: BDD, Domain Model Contract, Architecture Contract, API / Interface Contract, Error Handling Contract, implementation slices, fixtures, tests, and review checklists. Keep raw APK evidence, target-specific hosts, tokens, accounts, and private product conclusions in the project analysis docs.
