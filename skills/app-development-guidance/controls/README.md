# Controls

This directory is the primary home for cross-platform security controls.

Use `controls/` before `platforms/`, `languages/`, or `implementation/` when the lesson is mostly about a security property rather than framework-specific or buildable steps.

| File | Topic |
| --- | --- |
| `api-transport.md` | HTTPS, certificate pinning decisions, request signing, replay resistance, API error hygiene. |
| `auth-session.md` | Token scope, refresh, revocation, logout, session invalidation, account binding. |
| `local-storage.md` | Secure storage, cache, backups, screenshots, offline data. |
| `logging-telemetry.md` | Log redaction, crash reports, analytics, security observability. |
| `anti-tamper-risk.md` | Root/hook/emulator signals, anti-tamper limits, risk scoring. |
| `release-build.md` | Obfuscation, debug flag removal, symbol stripping, dependency and secret checks. |

If a control has concrete implementation details, keep the core principle here and link to `implementation/`, `platforms/`, or `languages/` as required by [`../../../enforcement/linked-updates.md`](../../../enforcement/linked-updates.md).

Common implementation links:

- Replay defense and request signing: [`../implementation/backend/api-replay-defense.md`](../implementation/backend/api-replay-defense.md), [`../implementation/backend/request-signing.md`](../implementation/backend/request-signing.md)
- Token/session lifecycle: [`../implementation/backend/token-session.md`](../implementation/backend/token-session.md)
- Mobile platform implementation: [`../implementation/mobile/`](../implementation/mobile/)
