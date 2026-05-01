# Backend Replay Defense Implementation

Use this when implementing replay resistance for sensitive API operations.

## Pattern

For replay-sensitive requests, the server should validate a combination of:

- Timestamp within an accepted window.
- Nonce or unique request ID scoped to account/session/client.
- Idempotency key for operations where retry should be safe.
- Body hash or canonical request digest when request signing is used.
- Server-side authorization and risk checks independent of client claims.

## Minimal Flow

```text
client builds request
client adds timestamp + nonce or idempotency key
client optionally signs canonical request
server checks auth/session
server checks timestamp window
server checks nonce/idempotency replay store
server verifies signature/digest if used
server processes once or returns idempotent result
```

## Validation

- Replay the exact same request and expect rejection or idempotent behavior.
- Modify a signed field and expect signature validation failure.
- Use an expired timestamp and expect rejection.
- Reuse a nonce under the same account/session and expect rejection.

## Required Linked Updates

When changing this pattern, also update or verify:

- [`../../controls/api-transport.md`](../../controls/api-transport.md)
- [`../../platforms/backend/api.md`](../../platforms/backend/api.md)
- [`../../checklists/api-security-review.md`](../../checklists/api-security-review.md)
