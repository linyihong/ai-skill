> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/app-development-guidance/execution-flow.md`](../../../../workflow/app-development-guidance/execution-flow.md)

### 2026-05-07 - Analysis-To-SDK Contract Drift Gate

Status: candidate

#### One-line Summary

When APK analysis closes or reclassifies live-readiness boundaries, immediately audit downstream SDK contracts, BDD notes, live-test gates, README wording, and fixture provenance for drift.

#### Human Explanation

Reverse-engineering or APK-analysis work often updates facts such as opaque parameter origin, base endpoint ownership, private adapter requirements, session recovery behavior, or synthetic fixture availability. If the SDK repository already has parser code, transport skeletons, or BDD docs, stale statements can make maintainers overclaim live readiness or miss the private adapter gate.

The fix is not always code. Often the correct closure is a downstream contract audit that turns fresh analysis into stable SDK-facing wording: what SDK core owns, what the private adapter owns, which fixtures are authoritative, and which live negative tests remain unknown.

#### Trigger

- APK analysis changes a readiness status from `open` to provider-confirmed, boundary-confirmed, private-adapter-required, or static-path-confirmed.
- New sanitized fixtures, adapter smoke checklists, host/base endpoint notes, signing/decrypt boundaries, or session-recovery paths are added to analysis docs.
- An SDK repository already has contracts, BDD features, integration-test docs, README examples, or live adapter skeletons for the analyzed feature.

#### Evidence

- Tool: repository diff / contract audit.
- Sanitized excerpt: downstream SDK docs still described an opaque parameter as source-unknown after analysis had reclassified it as a private provider boundary; live adapter docs also omitted the base endpoint provider and fixture authority.
- Evidence path: `<PROJECT_ROOT>/docs/development/*-contract.md`, `<PROJECT_ROOT>/docs/plans/integration/*-live-tests.md`, `<PROJECT_ROOT>/<sdk>/README.md`.

#### Generalized Lesson

Treat APK analysis closure as an upstream contract event for SDK projects. After updating analysis docs, run a downstream drift gate before declaring development readiness:

- Read the SDK plan, API/interface contract, error-handling contract, test plan, live integration spec, README, and BDD status table.
- Replace stale "unknown" language with the latest classification, but keep public SDK wording conservative when raw materials remain private.
- Add or update private adapter responsibilities for base endpoint, route binding, opaque providers, signing, decrypt, session recovery, and negative live matrices.
- Link fixture provenance from SDK fixtures back to sanitized schema-derived analysis fixtures when applicable.
- Avoid adding public code that embeds secrets or overclaims production readiness.

#### Agent Action

When analysis findings are meant to support SDK/client development, inspect the downstream SDK docs in the same work session. If any file still implies outdated ownership or readiness, update it before finalizing the analysis task. If code changes are unnecessary, explicitly mark the closure as docs-only and record what was validated.

#### Goal / Action / Validation

- Goal: keep SDK contracts aligned with latest analysis without leaking private materials or overstating live readiness.
- Action: audit downstream contracts and live-test gates after analysis closure; update stale wording and fixture provenance.
- Validation or reference source: `git diff --check`, lints for edited docs when available, and a reviewed diff limited to contracts/readiness/README files.

#### Applies When

- A reverse-engineering or APK-analysis project feeds an SDK, API client, mock adapter, parser, BDD suite, or live integration test.
- Public SDK core and private live adapter responsibilities are intentionally separated.
- Analysis creates or updates sanitized fixtures that should become parser/BDD authority.

#### Does Not Apply When

- There is no downstream SDK/client repository or development-facing contract.
- The analysis change is purely local scratch evidence and has not changed any readiness classification, schema, adapter boundary, or fixture authority.
- The SDK docs already explicitly defer to the updated analysis gate and no stale claim remains.

#### Validation

Check that updated SDK docs answer these questions consistently:

- Can SDK core run parser/mock tests without private live materials?
- Which provider or adapter owns host/base endpoint, route binding, opaque values, signing, decrypt, and session recovery?
- Which sanitized fixtures are authoritative for parser/BDD drift?
- Which live negative matrices remain unknown or private?

#### Promotion Target

- `WORKFLOW.md`
- `process/`
- `checklists/`

#### Required Linked Updates

- Update `feedback_history/README.md` and `feedback_history/common/README.md` indexes when adding this lesson.
- If promoted, add a checklist item to the app-development workflow for "analysis-to-SDK contract drift" after reverse-engineering readiness changes.
- Keep concrete project evidence in project docs; this lesson only records the generalized gate.
