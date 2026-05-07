> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-07 - Private Live Adapter Smoke Gate

Status: candidate

#### One-line Summary

When analysis supports SDK core design but live replay needs private host, service, session, signing, or decrypt material, create a private adapter smoke gate before live-facing development starts.

#### Human Explanation

It is common for reverse-engineering or integration analysis to identify routes, schemas, pagination, and parser contracts while intentionally redacting live secrets. That is enough to build a typed SDK core, mocks, fixtures, and parser tests, but not enough to claim a standalone live client. The development handoff needs a small adapter-level smoke checklist that names the private providers, BDD acceptance, redacted evidence, and failure classes before code begins.

#### Trigger

- Project docs say the SDK can model routes and response schemas, but raw base endpoints, route `service`, tokens, device material, signing headers, decrypt keys, or equivalent private materials are absent.
- A live-facing client, SDK, CLI, backend job, or automation is about to be implemented against a production-like service.
- The team wants to keep public docs sanitized while still making private implementation verifiable.

#### Evidence

- Tool: analysis documentation review, BDD/checklist drafting, redacted adapter contract.
- Sanitized excerpt: public docs can contain route ids, key names, schema ids, status classes, lengths, hashes, and pass/fail; private adapter logs keep raw secrets out of reusable docs.
- Evidence path: project docs under `<PROJECT_ROOT>/docs/` should contain the public readiness gate; private matrices stay outside public docs or in access-controlled implementation notes.

#### Generalized Lesson

Separate "SDK core ready" from "live adapter ready." The core may own typed queries, parsers, mock transport, fixtures, pagination semantics, and redacted logging. A private adapter must own base endpoint or host rotation, raw route binding, session/device values, locale/opaque providers, signing, decrypt, and any negative error matrix. Do not start live-facing code until the missing private adapter capabilities are either provided, scoped out, or represented as setup-failure behavior.

#### Agent Action

1. Classify readiness as SDK core, private adapter, or public standalone live SDK.
2. Before live-facing implementation, create or update a private adapter smoke checklist with provider interfaces, BDD acceptance, redacted evidence rules, and failure classifications.
3. Require setup failures for missing private material instead of letting the SDK send malformed live requests.
4. Keep automatic retry, refresh, relogin, or tamper behavior out of SDK core until an authorized negative matrix proves the live behavior.
5. Update owning project docs and feedback indexes in the same session.

#### Goal / Action / Validation

- Goal: prevent live SDK implementation from inventing missing private runtime behavior.
- Action: add a private adapter smoke gate before live-facing code and keep public docs sanitized.
- Validation or reference source: project docs distinguish SDK core readiness, private adapter readiness, public standalone readiness, and contain BDD/pass-fail evidence rules.

#### Applies When

- Building an SDK, client, mock API, fixture-driven parser, backend integration, or automation from APK/API analysis.
- Public documentation intentionally omits secrets or private replay material.
- The live service requires private host/base endpoint selection, signing, opaque request parameters, session/device identity, encrypted headers, or response decrypt.

#### Does Not Apply When

- The implementation is strictly offline parser/mock work with no live backend access.
- The service already has an official documented public API and test credentials covering signing/session/decrypt.
- The private adapter already exists and has a current smoke checklist plus passing evidence.

#### Validation

Confirm all of the following before declaring live-facing readiness:

- Public docs identify route ids, request key sets, schema/envelope shape, and sanitized logging rules.
- Private adapter capabilities are named as interfaces or setup prerequisites.
- At least one read-only smoke scenario defines pass/fail without raw values.
- Missing private capabilities fail before network send.
- Error/session recovery remains passthrough unless a negative matrix is authorized and captured.

#### Promotion Target

- `WORKFLOW.md`
- `checklists/`
- `process/`

#### Required Linked Updates

- Update the skill feedback root index and `common/README.md`.
- If the lesson came from an APK/API analysis handoff, keep project-specific routes, hashes, captures, and service names in project docs only.
- When promoted, cross-link from the APK-analysis-to-development handoff guidance so live-facing SDK work sees this gate before code.
