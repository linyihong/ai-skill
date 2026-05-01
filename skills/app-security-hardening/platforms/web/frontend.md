# Web Frontend Hardening Notes

Use this for browser/frontend-specific guidance when the same product family includes web apps.

## Focus Areas

- Token storage and session handling in browser contexts.
- XSS impact reduction, CSP, and trusted rendering boundaries.
- Sensitive logging and analytics redaction.
- Build-time secret checks.
- API contract alignment with mobile clients.

## Review Prompts

- Are tokens stored with the least exposure practical for the app architecture?
- Does XSS exposure turn into account takeover or payment/action abuse?
- Are frontend-only authorization checks backed by server-side enforcement?
- Are build artifacts free of secrets and internal endpoints?

See also:

- [`../../controls/auth-session.md`](../../controls/auth-session.md)
- [`../../controls/api-transport.md`](../../controls/api-transport.md)
- [`../../controls/logging-telemetry.md`](../../controls/logging-telemetry.md)
