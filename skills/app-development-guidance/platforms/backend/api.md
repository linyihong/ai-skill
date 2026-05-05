# Backend API Hardening Notes

Use this for backend controls that mobile and web apps depend on.

## Focus Areas

- Server-side authorization and object-level access checks.
- Replay resistance, idempotency, nonce/timestamp checks, and request signing validation.
- Token issuance, refresh, revocation, scope, and auditability.
- OpenAPI/schema as the source contract for generated clients, SDKs, mocks, and fixtures.
- Third-party/vendor API integration with sanitized excerpts, credential boundaries, webhook verification, and gated live tests.
- Rate limits, abuse controls, fraud signals, and anomaly detection.
- Safe API errors and security telemetry.

## Review Prompts

- Can the server reject a request even if the client is modified?
- Are replay-sensitive operations fresh, bound, or idempotent?
- Can sessions be revoked after risk events?
- Are generated clients rebuilt from the current OpenAPI/schema instead of hand-copied routes or DTOs?
- Are vendor integrations documented with sanitized operation excerpts, fixture tests, secret redaction, and opt-in live checks?
- Are security events observable without storing secrets?

See also:

- [`../../controls/api-transport.md`](../../controls/api-transport.md)
- [`../../controls/auth-session.md`](../../controls/auth-session.md)
- [`../../controls/anti-tamper-risk.md`](../../controls/anti-tamper-risk.md)
- [`../../implementation/backend/api-replay-defense.md`](../../implementation/backend/api-replay-defense.md)
- [`../../implementation/backend/contract-codegen.md`](../../implementation/backend/contract-codegen.md)
- [`../../implementation/backend/request-signing.md`](../../implementation/backend/request-signing.md)
- [`../../implementation/backend/token-session.md`](../../implementation/backend/token-session.md)
- [`../../implementation/backend/vendor-integration.md`](../../implementation/backend/vendor-integration.md)
