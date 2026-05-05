# Flutter Hardening Notes

Use this for Flutter and Dart-specific app hardening.

## Focus Areas

- Platform channel boundaries and native capability exposure.
- Dart AOT reverse-engineering implications.
- Plugin permissions, telemetry, and native SDK behavior.
- Release build flags, obfuscation, split debug symbols, and symbol management.
- Storage choices across Dart packages and native platform storage.

## Review Prompts

- Do platform channels expose privileged operations without server-side checks?
- Are debug menus, staging endpoints, and verbose logs removed from release builds?
- Are package/plugin permissions and network behavior reviewed?
- Are API controls server-owned rather than hidden in Dart code?

See also:

- [`../../controls/release-build.md`](../../controls/release-build.md)
- [`../../controls/local-storage.md`](../../controls/local-storage.md)
- [`../../languages/dart.md`](../../languages/dart.md)
- [`../../implementation/mobile/flutter.md`](../../implementation/mobile/flutter.md)
