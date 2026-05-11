> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-11 - Media prefix matrix and stale CDN classification

Status: candidate

#### One-line Summary

For APK media fields, document path-prefix-to-final-URL rules and classify static CDN 403 separately from URL composition bugs.

#### Human Explanation

API samples often show several relative path families. If docs only describe the common prefix, SDKs can miss exceptional path shapes. Conversely, a stale CDN object may fail under every valid URL composition and should not be treated as a mapper regression.

#### Trigger

Image/avatar/media downloads fail for a small subset of relative resource paths while other resources work.

#### Evidence

- Tool: API sample review plus live byte/header comparison.
- Sanitized excerpt: reachable samples produced identical encoded bytes through two URL forms, while one stale object returned 403 through both.
- Evidence path: project API docs, resource transport contract, and SDK integration notes.

#### Generalized Lesson

APK analysis should include a media prefix matrix: path family, final URL construction, required headers/session, decode/encryption rule, and stale-resource behavior. 403 classes must be separated into signed/authorization failures versus static object missing or CDN-denied.

#### Agent Action

When media fields feed an SDK downloader, collect representative samples for every path prefix, write the URL construction matrix, add SDK contract tests, and make bulk validation record per-resource errors instead of failing the whole run on stale avatars.

#### Goal / Action / Validation

- Goal: prevent media resource path-shape gaps from reaching SDK users.
- Action: add prefix matrix and stale CDN classification to project docs and SDK tests.
- Validation or reference source: contract tests cover each prefix; live spot checks compare bytes/headers for ambiguous URL forms.

#### Applies When

- APK/API fields expose relative media paths.
- SDKs or tools will download/decode those resources.

#### Does Not Apply When

- The API always returns fully signed absolute URLs with server-defined expiry and no local composition.

#### Validation

Use sanitized samples and short range/metadata checks; do not store raw private media in reusable skill docs.

#### Promotion Target

- `WORKFLOW.md`
- `DOCUMENTATION.md`
- `TOOLS.md`

#### Required Linked Updates

- Indexed in `feedback_history/media-hls/README.md`.
- Project-specific path examples and CDN host details stay in project docs.
