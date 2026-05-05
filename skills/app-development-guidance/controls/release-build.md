# Release Build Controls

Use this for release configuration, dependency checks, build artifacts, and shipping gates.

## Core Guidance

- Release builds should disable debug flags, test endpoints, dev menus, and verbose network logging.
- Obfuscation/minification should be enabled where compatible with crash symbolication and support needs.
- Native symbols should be stripped unless there is a documented diagnostics need.
- Build pipelines should scan for secrets and high-risk config drift.
- Third-party SDKs should be reviewed for permissions, data collection, network behavior, and supply-chain risk.

## Validation Ideas

- CI check for debug flags and endpoint configuration.
- Artifact inspection for symbols, test URLs, and known secret patterns.
- Dependency review before release.
- Smoke test release builds, not only debug builds.

## Common Overclaims

- Obfuscation raises cost but does not protect business logic by itself.
- A clean debug build does not prove release configuration is safe.
- Dependency risk includes runtime behavior, not just package version.
