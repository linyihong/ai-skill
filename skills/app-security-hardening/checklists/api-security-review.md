# API Security Review Checklist

- Authorization is enforced server-side for every sensitive object and action.
- Replay-sensitive operations use freshness, idempotency, or risk controls.
- Token scope, expiry, refresh, revocation, and audit behavior are documented.
- Error responses avoid sensitive implementation details.
- Rate limits and abuse monitoring are defined for high-risk endpoints.
- Client risk signals are treated as inputs, not trusted facts.
- Tests prove modified clients cannot bypass critical checks.
