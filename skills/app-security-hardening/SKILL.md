---
name: app-security-hardening
description: Convert authorized mobile app analysis findings into secure app development guidance. Use for mobile/API hardening, token and session handling, request signing, TLS and certificate pinning decisions, Flutter/Android release checks, local storage risks, sensitive logging, anti-tamper signals, and security review checklists.
---

# App Security Hardening

Use this skill when APK analysis, mobile API review, or app development work reveals a reusable security lesson that can improve future apps. The goal is to translate observed attack paths and failure modes into practical development requirements, checklists, and validation tests.

**Shared policy:** read [`shared-rules` index](../../shared-rules/README.md) and [`feedback-lessons.md`](../../shared-rules/feedback-lessons.md). Lessons in `feedback_history/` should reference those files, not duplicate shared rules.

## When To Use

- Turning APK analysis findings into secure development guidance.
- Reviewing mobile app/API design for replay resistance, token safety, transport security, local storage, logging, and release hardening.
- Creating PR/release checklists for Flutter/Android apps.
- Deciding what should be validated by tests, fixtures, runtime checks, or server-side controls.

## Out Of Scope

- Breaking into apps or services without authorization.
- Storing target-specific hosts, secrets, endpoints, tokens, real user data, or private business conclusions in this reusable skill.
- Treating client-side hardening as a replacement for server-side authorization, rate limits, fraud controls, or audit logs.

## Quick Start

1. Identify the observed weakness or reverse-engineering lesson.
2. Convert it into a developer-facing risk statement.
3. Choose the control layer:
   - API/server contract.
   - App runtime behavior.
   - Build/release configuration.
   - Monitoring or fraud signal.
4. Add a concrete validation method: unit test, integration test, release checklist item, fixture, or manual review step.
5. Classify the guidance:
   - Core security control: `controls/`.
   - Platform/app type detail: `platforms/`.
   - Language/runtime-specific trap: `languages/`.
   - Repeatable review step: `checklists/`.
6. If the lesson is reusable but not yet mature, add it to `feedback_history/` first; promote it into the structured folders when validated.

## Default Workflow

Read [WORKFLOW.md](WORKFLOW.md) to translate analysis evidence into hardening requirements.

Use [CHECKLIST.md](CHECKLIST.md) for the checklist index and [`checklists/`](checklists/) for focused design, PR, release, and API reviews.

Use [DOCUMENTATION.md](DOCUMENTATION.md) when writing reusable security guidance or project-specific hardening notes.

Use [`controls/`](controls/) as the primary home for cross-platform guidance, then link to [`platforms/`](platforms/) or [`languages/`](languages/) for implementation details.

## Output Style

When producing hardening guidance, include:

- Observed risk or failure mode.
- Why it matters for app development.
- Recommended control and the layer that owns it.
- Validation method.
- What not to overclaim.

## Feedback Loop

If a reusable app-security lesson emerges:

1. Create `feedback_history/YYYY-MM-DD_HHMMSS-<slug>.md` using [shared-rules/feedback-lessons.md](../../shared-rules/feedback-lessons.md).
2. Generalize the lesson so it is not tied to one APK or company.
3. Include evidence and validation criteria, but redact secrets and target-specific details.
4. Promote validated guidance into `controls/`, `platforms/`, `languages/`, `checklists/`, `WORKFLOW.md`, `CHECKLIST.md`, or `DOCUMENTATION.md` as appropriate.

**Cross-skill link:** if the lesson came from APK analysis, keep analysis mechanics in [`apk-analysis`](../apk-analysis/) and put development hardening guidance here.
