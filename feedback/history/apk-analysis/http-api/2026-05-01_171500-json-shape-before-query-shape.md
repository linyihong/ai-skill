> 遵守 [共用規則索引](../../../../enforcement/README.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/http-api-documentation-flow.md`](../../../../analysis/apk/workflows/http-api-documentation-flow.md)

### 2026-05-01 - JSON Shape Before Query Shape

Status: validated

#### One-line Summary

Schema-only response logging should classify JSON before running query-string extraction, otherwise embedded URLs inside JSON can be mistaken for request query evidence.

#### Human Explanation

APK traffic analysis often adds generic `queryKeys` helpers to summarize path/query/form material without printing values. If the same helper runs on every decoded string, a decrypted JSON response that contains signed image URLs, CDN URLs, or other URL fields may produce plausible-looking `queryKeys` and `pathHash` output. That output is response content shape, not request evidence, and can corrupt endpoint correlation.

#### Trigger

A decoded response string is already classified as JSON (`jsonLike`, top-level keys/types present), but the same log line also shows query keys such as URL auth parameters.

#### Evidence

- Tool: Frida Dart AOT offset hook with schema-only logging.
- Sanitized excerpt: a decrypted `ret` / `data` / `msg` JSON line produced both `jsonNested=...info:array[...]` and `queryKeys=...` because one JSON field contained an embedded URL query.
- Evidence path: `<PROJECT_ROOT>/capture/dart_info_item_shape_*.log` and the corresponding hook script diff.

#### Generalized Lesson

Treat JSON schema extraction and query/form extraction as mutually exclusive classifiers unless there is a specific reason to inspect embedded URLs. Run `jsonShape` first; if it succeeds, do not attach request-oriented `queryKeys`, `serviceHash`, or `pathHash` to that same decoded string.

#### Agent Action

When building sanitized payload loggers, order classifiers from highest semantic certainty to lower-level heuristics. For decoded strings:

1. Detect and summarize JSON as keys/types/nested item shape.
2. Only run query/form summarizers when the string is not JSON-like.
3. Keep embedded URL analysis in a separate, clearly named field if it is intentionally needed.

#### Applies When

- Hooking decrypted response strings, `jsonDecode` inputs, Dio response data, or similar high-level application payloads.
- Adding endpoint-correlation metadata such as `queryKeys`, `serviceHash`, or `pathHash`.
- Redacting raw values while preserving enough structure to correlate requests and responses.

#### Does Not Apply When

- The target string is known request canonical material, form data, or header plaintext rather than response JSON.
- The analysis goal is explicitly to inventory embedded URLs inside response bodies; in that case use a separate `embeddedUrlKeys` style label.

#### Validation

After applying the rule, request plaintext still reports expected request keys and service hashes, while decrypted JSON response lines keep only JSON keys/types/nested shape and no longer emit request-oriented query metadata from embedded URL fields.

#### Revision - 2026-05-01

Binary-like or pipe-separated decrypted payloads can also contain incidental `&` / `=` characters and produce fake query keys. The safer generalized rule is: run request-oriented query extraction only on strings that first look like request material, such as containing `service=` or an API path marker. Do not treat arbitrary non-JSON strings with `&` / `=` as request evidence.

#### Promotion Target

- `WORKFLOW.md`
- `DOCUMENTATION.md`

#### Required Linked Updates

- Project docs were updated to state that `queryKeys` / `serviceHash` are limited to non-JSON decoded strings.
- No `TOOLS.md` update is needed; this is a logging/classification rule, not a tool installation change.
