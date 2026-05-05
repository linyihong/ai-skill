# Local Storage Controls

Use this for secrets, caches, offline data, backups, and on-device exposure.

## Core Guidance

- Store secrets only when needed, and prefer platform-backed secure storage.
- Review shared preferences, SQLite, cache directories, downloaded files, and media storage for sensitive data.
- Define expiry and cleanup for offline data.
- Configure backup, screenshot, and clipboard behavior according to product risk.
- Avoid storing static API secrets or signing keys in the app bundle.

## Validation Ideas

- Inspect app data directories on a test device after using sensitive flows.
- Check backup behavior and exported files.
- Search logs, caches, and databases for token-like strings.
- Add regression tests for redaction and cache cleanup where the codebase supports it.

## Common Overclaims

- Local encryption does not protect secrets if the key is recoverable from the same app.
- Obfuscation does not make embedded secrets durable.
- Secure storage still needs token expiry, revocation, and server-side controls.
