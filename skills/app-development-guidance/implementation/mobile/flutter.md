# Flutter Implementation Notes

Use this for Flutter implementation patterns that complement `platforms/mobile/flutter.md`.

## Common Implementation Areas

- Platform channel method allowlists and input validation.
- Release build checks for debug menus, staging endpoints, verbose logs, and symbol handling.
- Secure storage package choices and native fallback behavior.
- Plugin review for permissions, telemetry, and native SDK behavior.

## Required Linked Updates

When changing Flutter implementation guidance, also update or verify:

- [`../../platforms/mobile/flutter.md`](../../platforms/mobile/flutter.md)
- [`../../languages/dart.md`](../../languages/dart.md)
- Relevant files in [`../../controls/`](../../controls/)
- Relevant files in [`../../checklists/`](../../checklists/)
