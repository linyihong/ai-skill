# App Security Hardening Checklist

Use this checklist during design review, PR review, and release review. It is intentionally high level: project-specific requirements should live in the product repository.

## API And Transport

- Sensitive flows use HTTPS only; cleartext traffic is disabled in release builds.
- Certificate pinning is considered for high-risk apps, with a rotation and incident plan.
- Backend authorization does not trust client-only flags, roles, prices, balances, or feature gates.
- Replay-sensitive requests have server-side timestamp, nonce, idempotency, or risk checks.
- Request signing, if used, signs the right fields and does not rely on a static client secret.
- Error responses avoid leaking stack traces, internal hostnames, or sensitive business rules.

## Auth, Tokens, And Sessions

- Access tokens are scoped, time-limited, and revocable.
- Refresh flows are rate-limited and bound to account/session context.
- Logout, password change, and risk events invalidate relevant sessions.
- Tokens are not written to debug logs, crash reports, analytics, or screenshots.
- Session identifiers are not used as long-term device identifiers.

## Local Storage

- Secrets are stored only when needed, and use platform-backed secure storage where possible.
- Cache files, SQLite, shared preferences, and downloaded media are reviewed for sensitive content.
- Backups and screenshots are configured according to product risk.
- Offline data has an expiry, encryption plan, or clear business justification.

## Flutter And Android Build

- Release builds disable debug flags, test endpoints, dev menus, and verbose network logging.
- Obfuscation/minification is enabled where compatible with crash symbolication and support needs.
- Native symbols are stripped unless needed for diagnostics.
- Flutter platform channels do not expose privileged operations without server-side or OS-level checks.
- Third-party SDKs are reviewed for permissions, telemetry, and data collection.

## Logging And Telemetry

- Logs redact tokens, cookies, authorization headers, device IDs, and personal data.
- Crash reports and analytics events avoid raw request/response payloads.
- Debug logging is gated away from release builds.
- Security-relevant failures are observable without exposing secrets.

## Anti-Tamper And Risk Signals

- Root/jailbreak/emulator/hook detection is treated as a risk signal, not a sole access-control decision.
- Server-side risk scoring can tolerate false positives and false negatives.
- Critical operations still require backend authorization and abuse controls.
- The app avoids storing static secrets that become permanent bypass targets.

## Release Gate

- A reviewer can point to tests, build checks, or documented evidence for every required control.
- Known residual risks are documented in the project repository.
- Reusable lessons are generalized into this skill only after sanitization.
