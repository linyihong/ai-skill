# React Native Hardening Notes

Use this for React Native app hardening.

## Focus Areas

- JavaScript bundle exposure and source map handling.
- Native module boundaries and bridge-exposed capabilities.
- OTA update controls and integrity checks.
- Release build logging and endpoint configuration.
- Secure storage choices across JS and native modules.

## Review Prompts

- Are native modules exposing only intended capabilities?
- Are source maps and debug bundles handled according to product risk?
- Are OTA updates signed or otherwise integrity controlled?
- Are secrets and tokens kept out of JS logs and analytics?

See also:

- [`../../controls/local-storage.md`](../../controls/local-storage.md)
- [`../../controls/release-build.md`](../../controls/release-build.md)
- [`../../controls/logging-telemetry.md`](../../controls/logging-telemetry.md)
- [`../../implementation/mobile/react-native.md`](../../implementation/mobile/react-native.md)
