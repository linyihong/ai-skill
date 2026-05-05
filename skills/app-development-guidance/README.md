# App Development Guidance Skill

This skill turns authorized app analysis findings into reusable development guidance for better, safer apps.

It complements [`apk-analysis`](../apk-analysis/):

- `apk-analysis` asks: "How does this APK behave, communicate, protect, or fail?"
- `app-development-guidance` asks: "What should we build, test, or avoid in our own apps because of that lesson?"

## Goals

Capture practical guidance for:

- App implementation patterns learned from analysis.
- App/API and transport security.
- Product brief to contract-first development flow.
- Existing project documentation gap backfill, with complete BDD recovery from implemented behavior.
- Token, session, and replay handling.
- Local storage and secret exposure risks.
- Mobile, web, backend/API, and release hardening.
- Sensitive logging and telemetry hygiene.
- Anti-tamper signals without false confidence.
- Developer checklists, buildable patterns, and validation tests.

## What Belongs Here

- Reusable development patterns learned from APK analysis, app/API review, or product development.
- Documentation backfill rules for already implemented projects, especially when original planning documents are missing.
- High-level controls that can be implemented by app, API, backend, or release engineering teams.
- Checklists that help prevent repeat mistakes.
- Guidance that clearly names validation steps and limitations.

## What Does Not Belong Here

- Target-specific API hosts, endpoints, tokens, device identifiers, or private response schemas.
- Raw request/response data.
- One-off product conclusions that do not generalize.
- Advice that relies only on client-side secrecy.

## Files

| File | Purpose |
| --- | --- |
| `SKILL.md` | Cursor/agent entry point. |
| `WORKFLOW.md` | Translate observations into development requirements and validation steps. |
| `process/` | Product brief to BDD, Domain Model, Architecture, API / Interface, Error Handling, implementation, and test flow. |
| `CHECKLIST.md` | Development, PR, and release checklist. |
| `DOCUMENTATION.md` | How to document reusable development guidance. |
| `controls/` | Cross-platform security controls; use this as the primary home for core guidance. |
| `platforms/` | Platform or app-type implementation guidance such as mobile, web, and backend API. |
| `languages/` | Language/runtime-specific pitfalls only. |
| `checklists/` | Focused design, PR, release, and API review checklists. |
| `implementation/` | Concrete implementation patterns for app, backend, and platform teams. |
| `templates/` | Copyable templates for initial development docs, guidance notes, and lightweight threat models; start with [`templates/README.md`](templates/README.md). |
| `FEEDBACK.md` | Short entry pointing to shared feedback rules. |
| `feedback_history/` | Reusable lessons, categorized by the same primary folders when applicable. |

## Classification Rules

When adding new guidance:

1. Put the core security property in [`controls/`](controls/) first.
2. Add platform details in [`platforms/`](platforms/) only when the implementation differs by app type or OS.
3. Add language notes in [`languages/`](languages/) only for language/runtime-specific traps.
4. Put concrete implementation patterns in [`implementation/`](implementation/) when engineers need buildable steps.
5. Put development process and contract-first flow guidance in [`process/`](process/).
6. Put repeatable review steps in [`checklists/`](checklists/).
7. Put draft reusable lessons in the matching [`feedback_history/<category>/`](feedback_history/) folder before promoting them into the structured folders; use `feedback_history/common/` for cross-cutting lessons.

For existing implemented projects, use [`process/`](process/) to audit and backfill missing documents. Product Brief gaps can remain `unknown`; BDD behavior must be completed from observable behavior and implementation evidence.

This keeps the skill readable as it grows across mobile, web, backend, and future app types.

## Required Linked Updates

Repo-wide rule: [`../../shared-rules/linked-updates.md`](../../shared-rules/linked-updates.md). Some changes must move together. If a change touches a control, platform, language, checklist, implementation pattern, or template, the related files **must** be updated or explicitly checked in the same change.

Examples:

- New replay-defense implementation -> update or verify `controls/api-transport.md`, `platforms/backend/api.md`, and `checklists/api-security-review.md`.
- New Flutter storage implementation -> update or verify `platforms/mobile/flutter.md`, `languages/dart.md`, `controls/local-storage.md`, and the relevant checklist.
- New review checklist item -> update or verify the matching `controls/` and `implementation/` docs.
- New contract-first process guidance -> update or verify `process/`, `checklists/`, `templates/`, and relevant implementation docs.
- New initial planning template -> update or verify `templates/README.md`, `process/README.md`, `DOCUMENTATION.md`, and `CHECKLIST.md`.
- New existing-project backfill rule -> update or verify `process/README.md`, `templates/initial-development-docs.md`, `CHECKLIST.md`, and `WORKFLOW.md`.

Do not describe these updates as optional. If they are relevant, they are required.

## Minimum Useful Output

A good development guidance note should include:

- The observed risk.
- The development consequence.
- The recommended control.
- The owner layer: client, API, backend, build, or monitoring.
- The validation method.
- Limits and non-goals.
