# TypeScript Notes

Use this for TypeScript app, web frontend, or SDK concerns.

- Avoid embedding permanent secrets in frontend or mobile JS bundles.
- Keep generated API clients aligned with server-side authorization and error handling.
- Redact sensitive values before logs, analytics, and error reporting.
- Treat type safety as input-shape help, not an authorization boundary.

See [`../platforms/web/frontend.md`](../platforms/web/frontend.md), [`../platforms/backend/api.md`](../platforms/backend/api.md), and backend implementation notes in [`../implementation/backend/`](../implementation/backend/) when TypeScript clients depend on API security contracts.
