# iOS Hardening Notes

Use this for iOS-specific implementation details.

## Focus Areas

- App Transport Security decisions and exceptions.
- Keychain storage and access control classes.
- Universal links and URL scheme handling.
- Jailbreak and hook signals as risk inputs.
- Release build configuration, symbol handling, and logging.

## Review Prompts

- Are ATS exceptions narrowly scoped and justified?
- Are sensitive tokens stored in Keychain with appropriate accessibility?
- Can custom URL schemes trigger privileged actions?
- Are crash reports and analytics free of secrets?

See also:

- [`../../controls/api-transport.md`](../../controls/api-transport.md)
- [`../../controls/auth-session.md`](../../controls/auth-session.md)
- [`../../controls/logging-telemetry.md`](../../controls/logging-telemetry.md)
- [`../../implementation/mobile/ios.md`](../../implementation/mobile/ios.md)
