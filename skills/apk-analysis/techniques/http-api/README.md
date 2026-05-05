# HTTP API Documentation Techniques

Use this category when the current task is to document, replay, validate, or map HTTP/HTTPS APIs. It applies whether the API was observed through MITM, Frida hooks, local proxy handlers, Flutter/Dart interceptors, or offline fixtures.

## When To Use

- Method/path/header/request/response metadata is visible or decoded.
- The goal is API reference, SDK mapping, replay, fixture creation, or contract tests.
- A UI flow needs to be bound to HTTP requests after core API behavior is understood.

## Core Output

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

If UI binding is not done yet, write `UI path: unknown` and `Trigger confidence: low`; later use screenshots and operation windows to raise confidence.

When the goal is to rebuild a feature, API docs should be ready for [`app-development-guidance`](../../../app-development-guidance/) to turn into BDD, Domain Model Contract, API / Interface Contract, Error Handling Contract, implementation slices, and tests. Mark uncertain field meaning or domain vocabulary as `candidate` instead of inventing final product language.

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

- `../../feedback_history/http-api/2026-05-01_171500-json-shape-before-query-shape.md`
- `../../feedback_history/http-api/2026-05-01_171650-ui-architecture-map-from-screenshots.md`
- `../../feedback_history/http-api/2026-05-01_173800-api-field-documentation-after-analysis.md`
- `../../feedback_history/http-api/2026-05-05_183400-ui-automation-operation-scripts-for-api-capture.md`
- `../../feedback_history/http-api/2026-05-05_183700-scrollable-clickable-screen-mapping.md`
- `../../feedback_history/http-api/2026-05-05_184700-screen-reachability-operation-recipes.md`
- `../../feedback_history/http-api/2026-05-05_184900-in-app-route-map-external-transitions.md`
