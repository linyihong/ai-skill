# HTTP API / UI Binding Feedback Lessons

Lessons for HTTP/API documentation, request/response schema classification, UI-to-API attribution, and field-level API references.

| File | Status | Topic | Notes |
| --- | --- | --- | --- |
| `2026-05-05_205500-ui-behavior-backfill-finish-gate.md` | promoted | UI Behavior backfill as finish gate | Analysis is not complete until observed UI behavior is backfilled into the project UI behavior/page-map docs or marked as a documented gap. |
| `2026-05-01_171500-json-shape-before-query-shape.md` | validated | JSON Shape Before Query Shape | Avoid treating embedded URLs or binary-like payloads as request query evidence. |
| `2026-05-01_171650-ui-architecture-map-from-screenshots.md` | promoted | UI architecture map from screenshots | Use lightweight UI maps and bind key APIs to operations. |
| `2026-05-01_173800-api-field-documentation-after-analysis.md` | promoted | API field documentation after analysis | Document headers, request fields, response fields, and field meanings. |
| `2026-05-05_183400-ui-automation-operation-scripts-for-api-capture.md` | promoted | UI automation operation scripts for API capture | Use small replayable operation scripts with timestamps to stabilize UI-to-API capture. |
| `2026-05-05_183700-scrollable-clickable-screen-mapping.md` | promoted | Scrollable and clickable screen mapping | Classify scrollable regions and clickable entry points before writing UI automation scripts. |
| `2026-05-05_184218-playable-list-item-still-needs-detail-api.md` | validated | Playable list item still needs detail API | Do not skip item-tap detail capture just because list items already contain playable URLs. |
| `2026-05-05_185921-scroll-depth-before-api-conclusion.md` | validated | Scroll depth before API conclusion | Capture top/mid/deep positions before deciding a scrollable page's API source. |
| `2026-05-05_190838-richtext-html-images-are-api-resources.md` | validated | RichText HTML images are API resources | Parse rich text `<img src>` resources, not only structured image arrays. |
| `2026-05-05_192012-infinite-scroll-needs-pagination-proof.md` | validated | Infinite scroll needs pagination proof | Prove page/cursor changes, `hasNext`, and tool pagination behavior for scrollable feeds. |
| `2026-05-05_195200-feature-reconstruction-handoff.md` | promoted | Feature reconstruction handoff | Preserve feature, behavior, domain, API, state/error, fixture, and unknown details so `app-development-guidance` can rebuild the functionality. |
