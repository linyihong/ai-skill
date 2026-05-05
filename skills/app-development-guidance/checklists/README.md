# Checklists

Use this directory for review checklists. Keep them short enough to run during real development.

| File | When to use |
| --- | --- |
| `mobile-design-review.md` | Before implementing a mobile feature or security-sensitive flow. |
| `mobile-pr-review.md` | During code review. |
| `mobile-release-review.md` | Before shipping a mobile release. |
| `api-security-review.md` | When mobile/web clients depend on API security properties. |
| `contract-governance-review.md` | When multiple planning, BDD, contract, generated, and test docs must stay traceable. |
| `embedded-firmware-review.md` | When firmware, sensors, boards, protocols, or hardware-in-loop validation are involved. |

The top-level [`../CHECKLIST.md`](../CHECKLIST.md) is the quick index and release-gate summary.

Checklist items must stay linked to the implementation docs they ask reviewers to verify. When adding a check, update or verify [`../implementation/`](../implementation/) and the relevant control docs in the same change.
