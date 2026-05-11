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
| `2026-05-06_081400-auto-app-guidance-for-sdk-tools.md` | promoted | Auto app-development-guidance for SDK/tool outputs | Automatically apply `app-development-guidance` when APK analysis docs are used to build app tools, SDKs, clients, mocks, contract tests, or rebuilt features. |
| `2026-05-06_104300-api-catalog-finish-gate.md` | promoted | API Catalog finish gate | API list work needs grouped indexes, per-API detail, coverage gaps, UI/API mapping, SDK/client field usage, evidence, and validation. |
| `2026-05-07_125600-api-first-pagination-replay.md` | candidate | API-first pagination replay | When replay prerequisites are known, validate pagination through direct read-only API replay before long UI scroll probes, then return to UI for attribution. |
| `2026-05-11_095600-ui-api-top-items-parity.md` | candidate | UI/API top items parity before API blame | Compare the first visible UI items with the first decoded API rows before blaming API data or SDK mapping. |
| `2026-05-11_113100-post-selection-gesture-for-lazy-api.md` | candidate | Post-selection gesture for lazy API | After selecting a tab/category/filter, add a bounded follow-up gesture before concluding no feature API fired. |
| `2026-05-11_113200-read-only-argument-override-preserve-app-path.md` | candidate | Read-only argument override preserve app path | Override read-only function arguments only in short windows while preserving app-owned signing/decrypt, and document the boundary. |
| `2026-05-11_135000-redacted-sample-targeting-classifier.md` | candidate | Redacted sample targeting classifier | Use disabled-by-default value-class/index classifiers to target rare samples without logging raw response values. |
| `2026-05-11_135700-articles-first-live-adapter-smoke.md` | candidate | Articles-first live adapter smoke | Validate one core read-only list route before making secondary routes mandatory for SDK/private adapter live proof. |
