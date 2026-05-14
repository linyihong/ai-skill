# Implementation

This directory explains how to turn hardening guidance into buildable patterns.

Use this folder when the question is "how do we implement it?" rather than "what security property do we need?"

When work starts from the contract-first process in [`../process/`](../process/), use implementation docs to turn contracts into build slices:

1. Map each Domain Model invariant to provider-side code and unit tests.
2. Map each API, event, command, or public interface contract to provider/consumer fixtures, mocks, or schema checks.
3. Map each Error Handling Contract entry to implementation behavior, logging redaction, and tests.
4. Keep implementation slices linked to the latest contract before teams or agents build in parallel.

| Directory | Scope |
| --- | --- |
| `backend/` | Server/API implementation patterns that mobile and web clients depend on. |
| `mobile/` | Android, iOS, Flutter, React Native implementation patterns. |
| `embedded/` | Firmware, sensor/protocol, hardware context, driver/service/application, and bring-up implementation patterns. |
| `tooling/` | IDE extensions, CLIs, linters, static analyzers, code generators, and internal automation. |
| `examples/` | Cross-cutting implementation patterns and snippets in pseudocode. |

Start with the sub-indexes:

- [`backend/README.md`](backend/README.md)
- [`embedded/README.md`](embedded/README.md)
- [`mobile/README.md`](mobile/README.md)
- [`tooling/README.md`](tooling/README.md)

## Required Linked Updates

Implementation docs are not standalone. Follow [`../../../enforcement/linked-updates.md`](../../../enforcement/linked-updates.md). When adding or changing an implementation pattern, you **must** update or verify the linked guidance:

1. `controls/` for the security principle.
2. `platforms/` for platform-specific behavior, if relevant.
3. `languages/` for language/runtime-specific traps, if relevant.
4. `checklists/` if reviewers need a repeatable check.
5. `templates/` if the documentation shape changes.

If no linked update is needed, say why in the change note or commit message.
