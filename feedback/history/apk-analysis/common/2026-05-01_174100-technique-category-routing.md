> 遵守 [共用規則索引](../../../../enforcement/README.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

# Extracted — See [`enforcement/cross-skill-references.md`](../../../../enforcement/cross-skill-references.md)

### 2026-05-01 - Technique category routing

Status: promoted

#### One-line Summary

APK analysis techniques should be split by runtime/API family so an A-type analysis does not need to load B-type guidance.

#### Human Explanation

The `apk-analysis` skill started with heavy Flutter/Dart AOT usage, but future APKs may be Java/Kotlin, WebView, local proxy, media/HLS, native, Cronet, or other categories. If every lesson is promoted into the same top-level workflow, the skill becomes noisy and agents waste context reading irrelevant techniques.

The better pattern is a common routing workflow plus category folders. The common workflow answers "what kind of app/traffic is this?" Once evidence identifies a category, the agent reads only that category. Shared rules such as authorization, sanitization, API documentation, and evidence chains remain top-level.

#### Trigger

- A skill folder starts mixing Flutter-specific, media-specific, local proxy-specific, and generic HTTP guidance.
- A new app/runtime family appears and does not fit the existing dominant category.
- The user notes that future A-type analysis should not need to read B-type files.

#### Evidence

- Tool: repository structure review.
- Sanitized excerpt: `techniques/flutter-dart-aot/`, `techniques/http-api/`, `techniques/local-proxy/`, `techniques/media-hls/`.
- Evidence path: reusable skill docs only; no target-specific facts.

#### Generalized Lesson

Keep the default skill entrypoint small and route-based. Put category-specific guidance under `techniques/<category>/`, and link related feedback lessons from that category. Do not load every category by default; use evidence from the common workflow to choose the relevant category.

#### Agent Action

When a reusable APK analysis technique is discovered:

- Decide whether it is cross-cutting or category-specific.
- If cross-cutting, promote it to top-level `WORKFLOW.md`, `TOOLS.md`, or `DOCUMENTATION.md`.
- If category-specific, promote it to `techniques/<category>/README.md`.
- If no category exists, create one with a short README and link it from `techniques/README.md`.
- Keep the original dated lesson in `feedback_history/` as history.

#### Applies When

- Maintaining `apk-analysis` or a similar broad analysis skill.
- Adding guidance for runtime families, transport families, API documentation families, or media families.
- Trying to reduce irrelevant context for future agents.

#### Does Not Apply When

- A rule is truly universal, such as authorization, sanitization, feedback naming, or evidence requirements.
- The category is not yet evidenced and the technique is still a one-off observation; keep it as a candidate lesson first.

#### Validation

- A new APK analysis can start from `SKILL.md` and `WORKFLOW.md`, then route to one relevant `techniques/<category>/` folder.
- Top-level docs stay readable and do not become a long list of every runtime-specific trick.
- Each category README links its related lessons instead of duplicating all historical feedback content.

#### Promotion Target

- `SKILL.md`
- `README.md`
- `WORKFLOW.md`
- `RUNBOOK.md`
- `techniques/`

#### Required Linked Updates

- Added `techniques/README.md` category routing index.
- Added initial category READMEs for Flutter/Dart AOT, HTTP API, local proxy, and media/HLS.
- Updated `SKILL.md`, `README.md`, `WORKFLOW.md`, and `RUNBOOK.md` to route by category.
- Updated `feedback_history/README.md` index.
