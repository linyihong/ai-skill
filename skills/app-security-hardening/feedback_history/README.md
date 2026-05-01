# app-security-hardening feedback history

Each file in this directory is one reusable app security hardening lesson.

Follow [`shared-rules/feedback-lessons.md`](../../../shared-rules/feedback-lessons.md). Do not duplicate shared rules in every lesson.

## Categories

This skill already separates promoted guidance into `controls/`, `platforms/`, `languages/`, `implementation/`, and `checklists/`. New feedback lessons should use the matching `feedback_history/<category>/` folder when the primary category is clear; cross-cutting lessons can use `feedback_history/common/`.

| Category | Purpose |
| --- | --- |
| [`controls/`](controls/) | Lessons whose primary promotion target is a reusable security control. |

## Historical Flat Index

| File | Status | Topic | Summary |
| --- | --- | --- | --- |
| `controls/2026-05-01_142100-client-encrypted-header-not-boundary.md` | promoted | Client encrypted header is not a security boundary | Client-side encrypted or signed headers are recoverable from shipped apps; backend authorization, replay protection, token hygiene, and monitoring must provide the real boundary. |
