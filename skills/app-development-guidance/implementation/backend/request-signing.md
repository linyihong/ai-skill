# Request Signing Implementation

Use this when request signing is justified by the threat model.

## Pattern

Request signing should bind the request to server-verifiable facts:

- Method.
- Path and normalized query.
- Timestamp.
- Nonce or idempotency key.
- Body hash.
- Account/session/client context where applicable.

The signing secret must not be a permanent static secret embedded only in the client app. Prefer server-issued, scoped, revocable material or platform-backed keys when the product risk justifies it.

## Canonicalization Checklist

- Define exact header names and casing behavior.
- Define query parameter sorting and encoding.
- Define body hashing for empty and streaming bodies.
- Define timestamp format and clock-skew tolerance.
- Define what happens on duplicate headers or repeated query keys.

## Validation

- Signed request succeeds.
- Mutating method, path, body, timestamp, or nonce fails.
- Expired timestamp fails.
- Replayed nonce fails or returns an idempotent result.

## Required Linked Updates

When changing this pattern, also update or verify:

- [`../../controls/api-transport.md`](../../controls/api-transport.md)
- [`api-replay-defense.md`](api-replay-defense.md)
- [`../../checklists/api-security-review.md`](../../checklists/api-security-review.md)
