# HTTP API Documentation Techniques

Use this category when the current task is to document, replay, validate, or map HTTP/HTTPS APIs. It applies whether the API was observed through MITM, Frida hooks, local proxy handlers, Flutter/Dart interceptors, or offline fixtures.

> **相容性規則**：`skills/apk-analysis/techniques/http-api/` 仍為 active skill entrypoint。本文件為 reference target，兩者應保持同步。

## When To Use

- Method/path/header/request/response metadata is visible or decoded.
- The goal is API reference, SDK mapping, replay, fixture creation, or contract tests.
- A UI flow needs to be bound to HTTP requests after core API behavior is understood.

## Core Output

For API-list or API-reference tasks, create or update a project-level API Catalog instead of leaving endpoints only in logs, chat, or one flat table. A good catalog has a total API entry, grouped indexes, per-operation detail files, coverage/gap status, UI/API mapping, SDK/client field usage, evidence, and validation.

For each HTTP API, write project documentation that includes:

- Method and path shape.
- Capability or feature this API supports, with operation id and trigger confidence.
- Auth/session requirements.
- Request headers with meaning, source, sensitivity, and whether they are required.
- Request query/path/body fields with type/shape, meaning, required/optional, source, and signing/encryption participation.
- Response headers when visible, especially content type, cache, rate-limit, session, and status behavior.
- Response wrapper fields and decrypted/inner payload fields with type, meaning, optional/nullability, and list item shape.
- Candidate domain concepts, state impact, error/empty behavior, pagination/cache rules, and fixture needs when the API is needed to rebuild a feature.
- Evidence and validation: hook/MITM/pcap/replay/fixture, plus UI path or screenshot only when it helps attribution.

Screenshots can support UI trigger attribution, but they do not replace HTTP header/request/response field analysis.

## API Catalog Shape

Use the project naming convention, but keep these artifacts explicit:

| Artifact | Required content |
| --- | --- |
| API entry | Hosts/base URLs, traffic families, wrapper/decode rules, shared headers, links to coverage, UI map, SDK/client notes, and grouped indexes. |
| Group index | API grouped by path prefix, domain, feature, or protocol family; each row links to per-API detail. |
| Per-API detail | Request, response, field meaning, behavior, evidence, validation, open questions. |
| Coverage / gap matrix | Static candidates, observed APIs, replayed APIs, decoded APIs, UI-bound APIs, tested APIs, missing parameters, untriggered flows, scope-out decisions. |
| SDK/client mapping | Fields actually consumed, compatibility expectations, raw JSON strategy, fixtures/tests. |

If a high-value endpoint has only a row in a table, create a per-API detail skeleton and mark missing sections as `needs capture`, `needs replay`, `meaning unknown`, or `low confidence`.

## API Documentation Flow

When an API has been observed or decoded, do not stop at the endpoint name. Document at least:

| Area | Required Notes |
| --- | --- |
| Identity | Method, host/path shape, auth conditions, evidence source, UI path if confirmed. |
| Capability mapping | Feature/capability, operation id, user-visible behavior, trigger confidence, and whether this API is startup/preload/background or direct user action. |
| Request headers | Header name, purpose, required/optional, source, sensitivity, token/sign/device/session involvement. |
| Request query/body | Field type, meaning, required/optional, example shape, sensitivity, signing/encryption participation. |
| Response headers | Status behavior, content type, cache/rate-limit/session headers; if invisible, state why. |
| Response wrapper | `status`, `code`, `message`, `data`, `error`, and other outer fields with type and meaning. |
| Inner payload | Field type, meaning, nullability, list item shape, media/source fields, derived values. |
| Functional contract | Candidate domain concepts, commands/events, state impact, empty/error behavior, pagination/cache semantics, and open questions. |
| Validation | Replay, fixture, contract test, or hook/pcap/MITM sequence proving request/response alignment. |
| Catalog status | grouped, per-API detail exists, coverage status, SDK/client mapping status when relevant. |

If UI binding is not done yet, write `UI path: unknown` and `Trigger confidence: low`; later use screenshots and operation windows to raise confidence.

When the goal is to rebuild a feature, API docs should be ready for `workflow/app-development-guidance/` to turn into BDD, Domain Model Contract, API / Interface Contract, Error Handling Contract, implementation slices, and tests. Mark uncertain field meaning or domain vocabulary as `candidate` instead of inventing final product language.

## API Catalog Finish Gate

Before reporting an API-list task complete, check:

- Every observed or decoded API is in a group index or a coverage/gap file.
- High-value APIs have per-operation detail, not only method/path rows.
- Each per-operation detail includes request fields, response fields, field meaning, evidence, validation, and open questions.
- Shared headers, wrapper/decode behavior, auth/session, and sensitivity rules are documented once and linked from API details.
- UI/API mapping records operation id, capture window, trigger confidence, and startup/preload/background status.
- SDK/client/tool usage records consumed fields and fixture/test status when the API is used for implementation.
- Unverified APIs are explicitly marked `candidate`, `needs capture`, `needs replay`, `meaning unknown`, `low confidence`, `out of scope`, or `not observed`.

## UI Automation For API Capture

For high-value flows, a small operation script can make API capture repeatable:

1. Give each flow a stable `operation_id`.
2. Give each route a stable `route_id` that explains how to reach the target screen.
3. Keep route maps scoped to in-app pages; if a step opens a system screen, browser, payment/share sheet, third-party app, or external intent, mark it as an external transition and stop extending the app screen map.
4. Classify each in-app screen as scrollable or not, and record clickable entry points before writing a script.
5. Keep each script to one UI path or action group, such as `open-home`, `scroll-feed`, `open-detail`, or `start-playback`.
6. Convert the documented route recipe into explicit launch/tap/swipe steps.
7. For scrollable screens, use bounded sampling such as top/mid/bottom rather than crawling to the end.
8. For clickable screens, use labels, resource IDs, content descriptions, hierarchy bounds, or verified coordinates.
9. Print UTC start/end timestamps before and after the operation, and include step-level route logs when useful.
10. Run pcap/MITM/Frida capture in the same window.
11. Save one sanitized screenshot or UI hierarchy at the end of the operation.
12. Fill the operation-to-API matrix with route id, method/path, source, response shape, and trigger confidence.

Use automation to stabilize capture, not to crawl the whole app. Avoid scripts that perform login loops, payment, destructive actions, posting, messaging, account changes, or any flow outside authorization.

## Related Lessons

- `skills/apk-analysis/feedback_history/http-api/2026-05-01_171500-json-shape-before-query-shape.md`
- `skills/apk-analysis/feedback_history/http-api/2026-05-01_171650-ui-architecture-map-from-screenshots.md`
- `skills/apk-analysis/feedback_history/http-api/2026-05-01_173800-api-field-documentation-after-analysis.md`
- `skills/apk-analysis/feedback_history/http-api/2026-05-05_183400-ui-automation-operation-scripts-for-api-capture.md`
