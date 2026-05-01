# Auth, Tokens, And Session Controls

Use this for identity, authentication, token lifecycle, and session integrity.

## Core Guidance

- Keep access tokens scoped, short-lived, revocable, and bound to the relevant account/session context.
- Refresh token flows should be rate-limited and observable.
- Logout, password change, device loss, and risk events should invalidate relevant sessions.
- Do not log tokens, cookies, authorization headers, refresh tokens, or one-time codes.
- Avoid long-term device identifiers as authentication substitutes.
- Design token rotation so clients can recover from expired credentials without broad retries that look like abuse.

## Validation Ideas

- Unit test token parsing and expiry handling.
- Integration test session invalidation after logout or password change.
- Abuse test refresh endpoint rate limits.
- Log scan for authorization headers and session values.

## Common Overclaims

- A hidden token format is not an access-control boundary.
- Device binding reduces risk but does not replace account authorization.
- Client-side logout without server invalidation is not enough for high-risk products.
