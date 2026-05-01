# Swift Notes

Use this for Swift and iOS implementation details.

- Use Keychain intentionally, with access control appropriate to the data and user experience.
- Keep URL scheme and universal link handlers narrow and server-validated for privileged actions.
- Review logging, crash reporting, and analytics at the call site before data leaves the process.
- Treat jailbreak checks as risk signals, not standalone authorization.

See [`../platforms/mobile/ios.md`](../platforms/mobile/ios.md) for iOS platform guidance and [`../implementation/mobile/ios.md`](../implementation/mobile/ios.md) for buildable implementation notes.
