# Controls

This directory is the primary home for cross-platform security controls.

Use `controls/` before `platforms/` or `languages/` when the lesson is mostly about a security property rather than a framework-specific implementation.

| File | Topic |
| --- | --- |
| `api-transport.md` | HTTPS, certificate pinning decisions, request signing, replay resistance, API error hygiene. |
| `auth-session.md` | Token scope, refresh, revocation, logout, session invalidation, account binding. |
| `local-storage.md` | Secure storage, cache, backups, screenshots, offline data. |
| `logging-telemetry.md` | Log redaction, crash reports, analytics, security observability. |
| `anti-tamper-risk.md` | Root/hook/emulator signals, anti-tamper limits, risk scoring. |
| `release-build.md` | Obfuscation, debug flag removal, symbol stripping, dependency and secret checks. |

If a control has platform-specific implementation details, keep the core principle here and link to `platforms/` or `languages/`.
