# Android Implementation Notes

Use this for Android implementation patterns that complement `platforms/mobile/android.md`.

## Common Implementation Areas

- Network Security Config for cleartext policy and scoped certificate behavior.
- Manifest review for exported components and deeplink entry points.
- Keystore-backed storage for high-risk tokens or keys.
- Release build configuration for logging, debug flags, symbols, and endpoints.

## Required Linked Updates

When changing Android implementation guidance, also update or verify:

- [`../../platforms/mobile/android.md`](../../platforms/mobile/android.md)
- [`../../languages/kotlin-java.md`](../../languages/kotlin-java.md)
- Relevant files in [`../../controls/`](../../controls/)
- Relevant files in [`../../checklists/`](../../checklists/)
