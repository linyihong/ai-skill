> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/software-delivery/execution-flow.md`](../../../../workflow/software-delivery/execution-flow.md)

### 2026-05-07 - Schema-Derived Synthetic Fixtures

Status: candidate

#### One-line Summary

When live payloads are private but response schemas are known, create synthetic schema-compatible fixtures so parser and BDD work can start without copying raw captured data.

#### Human Explanation

APK/API analysis often identifies route ids, envelopes, fields, types, empty-list variants, and pagination shapes before a private live adapter is ready. Waiting for replayable live payloads blocks useful SDK parser work. Copying captured payloads is unsafe. The middle path is to create synthetic fixtures with placeholder values that preserve shape and type but contain no secrets, user content, media URLs, tokens, hosts, or replayable material.

#### Trigger

- Analysis docs contain confirmed or candidate response schemas, but no checked-in parser fixtures.
- Public docs intentionally redact raw payload values or live replay material.
- SDK, mock transport, contract tests, or BDD parser scenarios are about to start.
- Empty, error, or pagination variants are documented as shape-only evidence.

#### Evidence

- Tool: schema catalog, API contract docs, fixture JSON validation.
- Sanitized excerpt: fixtures may include route id, schema id, service hash, field names, placeholder strings, booleans, numbers, arrays, and empty variants.
- Evidence path: project fixture directories such as `<PROJECT_ROOT>/api/fixtures/` or `<PROJECT_ROOT>/test/fixtures/`; raw captures remain outside committed fixtures.

#### Generalized Lesson

Treat synthetic fixtures as a bridge between analysis and implementation. They are not proof of live behavior and not substitutes for adapter smoke tests, but they let parser shape, null/empty handling, and BDD scenarios become executable earlier. Each fixture should declare that it is synthetic, name the source schema/route, and preserve only the minimum field/type shape needed for tests.

#### Agent Action

1. Identify schema shapes that are ready for parser work.
2. Create synthetic fixtures for happy path, empty path, and notable type variants.
3. Add a fixture README with redaction rules and source schema references.
4. Validate fixture JSON syntax or schema shape before committing.
5. Update SDK/readiness docs from "fixtures missing" to "synthetic fixtures present" without claiming live replay readiness.

#### Goal / Action / Validation

- Goal: unblock safe parser/BDD work without leaking captured payloads.
- Action: create schema-derived synthetic fixtures and link them from API/schema/readiness docs.
- Validation or reference source: JSON parses successfully, fixture fields match documented schema, and docs state that fixtures are synthetic and non-replayable.

#### Applies When

- Building SDK parsers, mock transports, contract tests, or BDD from reverse-engineered or observed API schemas.
- Live payloads contain user content, tokens, hosts, signatures, media URLs, or other private values.
- Schema confidence is high enough for parser shape tests, even if live adapter readiness is not complete.

#### Does Not Apply When

- Tests need exact cryptographic bytes, media segments, signatures, or server replay.
- Field semantics are unknown enough that even placeholder type choices would mislead implementation.
- Official public fixtures are already available and allowed to be committed.

#### Validation

Confirm all of the following:

- Fixtures parse as JSON or the target fixture format.
- Fixtures include synthetic/redaction metadata or a fixture README.
- Every placeholder value is non-sensitive and non-replayable.
- Empty and variant shapes are represented where documented.
- Owning API/schema/readiness docs link to the fixture location and preserve live-readiness caveats.

#### Promotion Target

- `WORKFLOW.md`
- `checklists/`
- `implementation/`

#### Required Linked Updates

- Update the skill feedback root index and `common/README.md`.
- Keep project-specific route names, hashes, schema ids, and fixture paths in project docs; this lesson only records the generalized rule.
- When promoted, cross-link from SDK/client and contract-test guidance so fixture creation happens before parser implementation when live payloads cannot be committed.
