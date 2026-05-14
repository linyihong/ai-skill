> 遵守 [共用規則索引](../../../../enforcement/README.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

# Extracted — See [`enforcement/cross-skill-references.md`](../../../../enforcement/cross-skill-references.md)

### 2026-05-01 - Workflow as routing, not technique dump

Status: promoted

#### One-line Summary

Top-level APK workflow should stay as common routing; runtime/API-family details belong in `techniques/<category>/`.

#### Human Explanation

As `apk-analysis` grows, it is tempting to keep adding full Flutter, local proxy, HTTP API, media, or future category flows into `WORKFLOW.md`. That makes the default workflow too long and causes agents to load irrelevant context. The top-level workflow should answer "what evidence do we have, and which route should we take?" Once a route is known, the detailed steps should live in the matching category folder.

#### Trigger

- `WORKFLOW.md` contains full runtime-specific or API-family-specific procedures.
- A category already exists in `techniques/`, but the same detailed flow is duplicated in top-level docs.
- A future A-type analysis would have to read B-type guidance before it can proceed.

#### Evidence

- Tool: repository documentation review.
- Sanitized excerpt: Flutter/Dart AOT flow, local proxy Netty details, HTTP API field documentation, and media/HLS chain were moved out of `WORKFLOW.md` into matching `techniques/` folders.
- Evidence path: reusable skill docs only.

#### Generalized Lesson

Keep `WORKFLOW.md` short, common, and route-oriented. Promote category-specific details into `techniques/<category>/README.md`; keep only minimal pointers in the workflow table. If a category-specific lesson is discovered, write it under `feedback_history/<category>/` and promote it to that category, not to the top-level workflow.

#### Agent Action

When editing `apk-analysis`:

- Put cross-cutting authorization, routing, evidence, sanitization, and completion criteria in top-level docs.
- Put Flutter/Dart, HTTP API, local proxy, media/HLS, or future family details in `techniques/<category>/`.
- If top-level `WORKFLOW.md` grows with category steps, refactor those steps into a category file and leave a routing pointer.
- Add feedback lessons under the matching `feedback_history/<category>/` or `feedback_history/common/`.

#### Applies When

- Maintaining broad analysis skills with multiple runtime/API families.
- Adding new categories or promoting feedback lessons into structured docs.
- Reducing context load for future agents.

#### Does Not Apply When

- The rule is truly universal and should be read before category routing.
- A category does not exist yet and the guidance is still experimental; keep it as a feedback lesson first.

#### Validation

- A new analysis can read `SKILL.md` and `WORKFLOW.md` without loading all category details.
- Category-specific procedures are discoverable from `techniques/README.md`.
- `WORKFLOW.md` points to category files instead of duplicating their detailed steps.

#### Promotion Target

- `WORKFLOW.md`
- `techniques/flutter-dart-aot/README.md`
- `techniques/http-api/README.md`
- `techniques/local-proxy/README.md`
- `techniques/media-hls/README.md`

#### Required Linked Updates

- Trimmed `WORKFLOW.md` to route-oriented guidance.
- Moved detailed Flutter/Dart AOT guidance into `techniques/flutter-dart-aot/README.md`.
- Moved detailed local proxy handler guidance into `techniques/local-proxy/README.md`.
- Moved detailed HTTP API documentation guidance into `techniques/http-api/README.md`.
- Moved detailed media/HLS chain guidance into `techniques/media-hls/README.md`.
- Updated feedback indexes.
