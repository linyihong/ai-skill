# app-development-guidance feedback history

Each file in this directory is one reusable app development guidance lesson.

Follow [`shared-rules/feedback-lessons.md`](../../../shared-rules/feedback-lessons.md). Do not duplicate shared rules in every lesson.

## Categories

This skill already separates promoted guidance into `controls/`, `platforms/`, `languages/`, `implementation/`, and `checklists/`. New feedback lessons should use the matching `feedback_history/<category>/` folder when the primary category is clear; cross-cutting lessons can use `feedback_history/common/`.

| Category | Purpose |
| --- | --- |
| [`common/`](common/) | Cross-cutting lessons that affect process, implementation, checklists, templates, or multiple guidance areas. |
| [`controls/`](controls/) | Lessons whose primary promotion target is a reusable security control. |

## Historical Flat Index

| File | Status | Topic | Summary |
| --- | --- | --- | --- |
| `common/2026-05-05_194400-contract-first-development-flow.md` | promoted | Contract-first development flow | Start from product brief, split bounded contexts, write BDD, define Domain, Architecture, API/Interface, and Error Handling Contracts, then build and test provider/consumer sides. |
| `common/2026-05-05_200500-existing-project-doc-backfill-bdd-required.md` | promoted | Existing project doc backfill requires complete BDD | Existing projects must audit and backfill missing development docs; Product Brief gaps can be unknown, but BDD must be completed from observable implementation evidence. |
| `controls/2026-05-01_142100-client-encrypted-header-not-boundary.md` | promoted | Client encrypted header is not a security boundary | Client-side encrypted or signed headers are recoverable from shipped apps; backend authorization, replay protection, token hygiene, and monitoring must provide the real boundary. |
