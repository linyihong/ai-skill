# Contract Codegen And Typed Clients

Use this pattern when a provider API, schema, SDK, generated client, or frontend consumer depends on a shared contract such as OpenAPI, GraphQL schema, protobuf, event schema, CLI schema, or TypeScript type package.

## Rule

The source contract owns operation names, routes, DTOs, errors, auth/session shape, versioning, and compatibility. Consumers should use generated or schema-derived code instead of hand-copying endpoints, field names, or response wrappers.

## Required Flow

1. Update the provider contract first: OpenAPI/schema/event/command/public method contract.
2. Update BDD and error-handling docs for changed behavior.
3. Regenerate typed clients, SDKs, mocks, fixtures, or schema packages from the source contract.
4. Compile or test every consumer that imports the generated output.
5. Add provider/consumer contract tests for critical operations and at least one important failure path.
6. Commit generated artifacts only if the project intentionally checks them in; otherwise document the generation command.

## Drift Checks

- No hand-copied URLs, route strings, DTOs, enum values, or response envelopes in consumers when generated equivalents exist.
- Generated client version matches the provider contract version or commit.
- Mocks and fixtures are derived from the current contract or explicitly pinned.
- BDD scenario tags, operation IDs, generated method names, and tests are traceable.
- Error codes and response wrappers are tested, not only happy-path payloads.
- Breaking changes require migration notes, compatibility tests, or explicit version bump.

## Test Strategy

| Test | Proves |
| --- | --- |
| Provider contract test | The server/provider emits shapes matching the contract. |
| Consumer compile/type check | Generated types are consumable by frontend, SDK, CLI, or jobs. |
| Fixture round-trip | Mocked responses and requests match the current schema. |
| Integration test | Real provider and consumer agree on auth, errors, and representative data. |
| Regression test | Existing consumers still compile and critical scenarios still pass. |

## Anti-Patterns

- Updating frontend or SDK request code before updating the provider contract.
- Hand-writing "temporary" DTOs that duplicate generated types.
- Treating generated client success as proof that auth, permission, idempotency, or error behavior is correct.
- Committing stale generated files without the source contract change.

## Required Linked Updates

Follow [`../../../../shared-rules/linked-updates.md`](../../../../shared-rules/linked-updates.md). Contract-codegen changes must update or verify [`../../process/`](../../process/), [`../../CHECKLIST.md`](../../CHECKLIST.md), [`../../checklists/api-security-review.md`](../../checklists/api-security-review.md), backend/API platform notes, and any template fields that describe API/interface contracts.
