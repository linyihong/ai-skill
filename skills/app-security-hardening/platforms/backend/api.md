# Backend API Hardening Notes

Use this for backend controls that mobile and web apps depend on.

## Focus Areas

- Server-side authorization and object-level access checks.
- Replay resistance, idempotency, nonce/timestamp checks, and request signing validation.
- Token issuance, refresh, revocation, scope, and auditability.
- Rate limits, abuse controls, fraud signals, and anomaly detection.
- Safe API errors and security telemetry.

## Review Prompts

- Can the server reject a request even if the client is modified?
- Are replay-sensitive operations fresh, bound, or idempotent?
- Can sessions be revoked after risk events?
- Are security events observable without storing secrets?

See also:

- [`../../controls/api-transport.md`](../../controls/api-transport.md)
- [`../../controls/auth-session.md`](../../controls/auth-session.md)
- [`../../controls/anti-tamper-risk.md`](../../controls/anti-tamper-risk.md)
