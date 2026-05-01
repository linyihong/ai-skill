# Dart Notes

Use this for Dart and Flutter-specific implementation details.

- Do not treat Dart AOT or obfuscation as a place to hide permanent secrets.
- Keep platform channel methods narrow and validate privileged operations server-side or with OS-backed controls.
- Review packages that touch storage, networking, analytics, crash reporting, or device identifiers.
- Use release build checks for debug logging, staging endpoints, and symbol handling.

See [`../platforms/mobile/flutter.md`](../platforms/mobile/flutter.md) for Flutter platform guidance.
