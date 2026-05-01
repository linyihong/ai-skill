# Implementation

This directory explains how to turn hardening guidance into buildable patterns.

Use this folder when the question is "how do we implement it?" rather than "what security property do we need?"

| Directory | Scope |
| --- | --- |
| `backend/` | Server/API implementation patterns that mobile and web clients depend on. |
| `mobile/` | Android, iOS, Flutter, React Native implementation patterns. |
| `examples/` | Cross-cutting implementation patterns and snippets in pseudocode. |

## Required Linked Updates

Implementation docs are not standalone. When adding or changing an implementation pattern, you **must** update or verify the linked guidance:

1. `controls/` for the security principle.
2. `platforms/` for platform-specific behavior, if relevant.
3. `languages/` for language/runtime-specific traps, if relevant.
4. `checklists/` if reviewers need a repeatable check.
5. `templates/` if the documentation shape changes.

If no linked update is needed, say why in the change note or commit message.
