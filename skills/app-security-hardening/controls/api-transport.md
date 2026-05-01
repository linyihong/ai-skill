# API And Transport Controls

Use this for controls that protect mobile app communication and backend API contracts.

## Core Guidance

- Enforce HTTPS and disable cleartext traffic in release builds.
- Treat certificate pinning as a risk-based control, not a default checkbox; include rotation, outage, and incident plans.
- Do not rely on a static secret embedded in the app for request signing.
- Replay-sensitive operations need server-side freshness controls: timestamp, nonce, idempotency key, account/session binding, or risk scoring.
- Backend authorization must not trust client-only flags, roles, balances, prices, or feature gates.
- API errors should avoid stack traces, internal hostnames, sensitive identifiers, and business-rule leakage.

## Validation Ideas

- Integration test: replay a captured signed request and verify rejection or idempotent handling.
- Release test: verify cleartext traffic is disabled.
- API test: mutate client-controlled fields and verify server-side authorization rejects them.
- Operations drill: validate pin rotation or rollback path before enabling pinning on critical flows.

## Common Overclaims

- Pinning does not prevent all reverse engineering.
- Request signing does not protect a static client secret once the app is analyzed.
- Client checks do not replace server-side authorization.
