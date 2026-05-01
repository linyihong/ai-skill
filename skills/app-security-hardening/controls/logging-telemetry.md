# Logging And Telemetry Controls

Use this for logs, analytics, crash reports, traces, and security observability.

## Core Guidance

- Redact tokens, cookies, authorization headers, refresh tokens, device identifiers, and personal data.
- Avoid raw request/response payloads in crash reports and analytics.
- Gate verbose network logging away from release builds.
- Security-relevant failures should be observable without exposing secrets.
- Document which fields are allowed in telemetry and which are forbidden.

## Validation Ideas

- Automated log scan for token-like values in test runs.
- Manual crash report review for high-risk flows.
- Release build assertion that debug logging flags are disabled.
- Analytics schema review for sensitive event fields.

## Common Overclaims

- "Internal only" logs can still leak through support tools, crash reporters, or developer devices.
- Redaction must happen before export, not only in dashboards.
- Debug logs often survive longer than intended unless release checks enforce removal.
