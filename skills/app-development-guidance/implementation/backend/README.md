# Backend Implementation

Backend implementation docs explain how server/API teams build the controls that mobile, web, and SDK clients rely on.

| File | Scope | Linked controls |
| --- | --- | --- |
| `api-replay-defense.md` | Timestamp, nonce, idempotency, and replay handling. | `../../controls/api-transport.md` |
| `contract-codegen.md` | OpenAPI/schema/source-contract to typed clients, SDKs, mocks, fixtures, and provider/consumer tests. | `../../process/README.md`, `../../checklists/api-security-review.md` |
| `request-signing.md` | Canonical request signing and validation boundaries. | `../../controls/api-transport.md` |
| `token-session.md` | Token lifecycle, refresh, revocation, and session invalidation. | `../../controls/auth-session.md`, `../../controls/logging-telemetry.md` |
| `vendor-integration.md` | Third-party API documentation split, credentials boundary, fixtures, gated live tests, webhooks, idempotency, and audit. | `../../controls/api-transport.md`, `../../controls/logging-telemetry.md` |

Required linked updates: changes here must update or verify `../../platforms/backend/api.md`, the linked controls, and `../../checklists/api-security-review.md`.
