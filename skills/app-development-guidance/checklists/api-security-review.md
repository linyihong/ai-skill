# API Security Review Checklist

- Authorization is enforced server-side for every sensitive object and action.
- Replay-sensitive operations use freshness, idempotency, or risk controls.
- Token scope, expiry, refresh, revocation, and audit behavior are documented.
- Bearer tokens, session IDs, and long-lived credentials are not accepted in URL query parameters.
- Client-generated encrypted headers or signatures are treated as tamper/risk signals, not as the only proof of authorization.
- Error responses avoid sensitive implementation details.
- Rate limits and abuse monitoring are defined for high-risk endpoints.
- Client risk signals are treated as inputs, not trusted facts.
- OpenAPI/schema changes regenerate typed clients, SDKs, mocks, or fixtures from the source contract.
- Consumers do not hand-copy endpoints, DTOs, enum values, or response envelopes when generated equivalents exist.
- Vendor or third-party API docs are split into sanitized integration excerpts, fixture examples, and gated live tests.
- Webhooks and callbacks verify signature, timestamp, body binding, idempotency, and replay behavior.
- Live integration tests are opt-in through environment/config and do not run against real services by default CI.
- Tests prove modified clients cannot bypass critical checks.
- Tests replay captured signed/encrypted requests and mutate account-bound parameters to verify server-side rejection or idempotent handling.
