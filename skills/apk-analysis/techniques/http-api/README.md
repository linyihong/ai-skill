# HTTP API Documentation Techniques

Use this category when the current task is to document, replay, validate, or map HTTP/HTTPS APIs. It applies whether the API was observed through MITM, Frida hooks, local proxy handlers, Flutter/Dart interceptors, or offline fixtures.

## When To Use

- Method/path/header/request/response metadata is visible or decoded.
- The goal is API reference, SDK mapping, replay, fixture creation, or contract tests.
- A UI flow needs to be bound to HTTP requests after core API behavior is understood.

## Core Output

For each HTTP API, write project documentation that includes:

- Method and path shape.
- Auth/session requirements.
- Request headers with meaning, source, sensitivity, and whether they are required.
- Request query/path/body fields with type/shape, meaning, required/optional, source, and signing/encryption participation.
- Response headers when visible, especially content type, cache, rate-limit, session, and status behavior.
- Response wrapper fields and decrypted/inner payload fields with type, meaning, optional/nullability, and list item shape.
- Evidence and validation: hook/MITM/pcap/replay/fixture, plus UI path or screenshot only when it helps attribution.

Screenshots can support UI trigger attribution, but they do not replace HTTP header/request/response field analysis.

## API Documentation Flow

When an API has been observed or decoded, do not stop at the endpoint name. Document at least:

| Area | Required Notes |
| --- | --- |
| Identity | Method, host/path shape, auth conditions, evidence source, UI path if confirmed. |
| Request headers | Header name, purpose, required/optional, source, sensitivity, token/sign/device/session involvement. |
| Request query/body | Field type, meaning, required/optional, example shape, sensitivity, signing/encryption participation. |
| Response headers | Status behavior, content type, cache/rate-limit/session headers; if invisible, state why. |
| Response wrapper | `status`, `code`, `message`, `data`, `error`, and other outer fields with type and meaning. |
| Inner payload | Field type, meaning, nullability, list item shape, media/source fields, derived values. |
| Validation | Replay, fixture, contract test, or hook/pcap/MITM sequence proving request/response alignment. |

If UI binding is not done yet, write `UI path: unknown` and `Trigger confidence: low`; later use screenshots and operation windows to raise confidence.

## Related Lessons

- `../../feedback_history/http-api/2026-05-01_171500-json-shape-before-query-shape.md`
- `../../feedback_history/http-api/2026-05-01_171650-ui-architecture-map-from-screenshots.md`
- `../../feedback_history/http-api/2026-05-01_173800-api-field-documentation-after-analysis.md`
