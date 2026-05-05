# Android Hardening Notes

Use this for Android-specific implementation details.

## Focus Areas

- Network Security Config and cleartext traffic policy.
- Manifest permissions, exported components, deeplinks, and intent handling.
- Keystore-backed storage and biometric gate decisions.
- Release build flags, minification, resource shrinking, symbol stripping, and debug menu removal.
- WebView configuration and JavaScript bridge exposure.

## Review Prompts

- Are exported activities/services/receivers intentional?
- Can deeplinks trigger privileged flows without server-side validation?
- Does release config differ safely from debug config?
- Are logs and crash reports clean in release builds?

See also:

- [`../../controls/api-transport.md`](../../controls/api-transport.md)
- [`../../controls/local-storage.md`](../../controls/local-storage.md)
- [`../../controls/release-build.md`](../../controls/release-build.md)
- [`../../implementation/mobile/android.md`](../../implementation/mobile/android.md)
