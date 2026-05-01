> 遵守 [共用規則索引](../../../shared-rules/README.md) 與 [feedback-lessons](../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-01 - Sequence jsonDecode Before Calling It API Response

Status: validated

#### One-line Summary

`jsonDecode` hook 看到的 schema 不一定是 API response；必須用 request/decrypt sequence 證明它發生在業務 response 路徑上。

#### Human Explanation

Flutter apps often parse local cache, preferences, bundled config, or restored session state during cold start. A global `jsonDecode` or `JsonDecoder.convert` hook can capture rich user/profile-looking JSON before any business request is sent. Without sequencing, it is easy to misclassify these startup/local JSON objects as decrypted API responses.

Schema-only logging is still useful, but it needs `reqSeq`, timestamps, and nearest request/decrypt context. If `jsonDecode` appears with `nearestReqSeq=0` or before the first relevant network request, keep it as a startup/local-cache candidate. Promote it to API response evidence only when it follows a matching request/decrypt cluster.

#### Trigger

- `jsonDecode` hook captures plausible response schema immediately after app launch.
- Request hooks / Netty headers / decrypt hooks have not yet fired.
- The schema looks like user/session/profile data, tempting the analyst to write endpoint conclusions.

#### Evidence

- Tool: Frida native offset hooks on Dart AOT request/decrypt/json functions with event sequence and timestamps.
- Sanitized excerpt: `jsonDecode:enter ... nearestReqSeq=0` before `RequestInterceptor._generateEhHeader` and Netty `eh` events; later API decrypt output appears as a separate non-JSON payload shape.
- Evidence path: project-private Frida logs under `<PROJECT_ROOT>/capture/`.

#### Generalized Lesson

Treat global JSON parse hooks as broad visibility, not endpoint attribution. For API documentation:

- keep startup/local JSON schemas separate from API response schemas;
- record `reqSeq`, `respSeq`, event timestamp, nearest request sequence, and decrypt sequence;
- require a request/decrypt/json time window before labeling a schema as endpoint response;
- document uncertain schemas as local/cache candidates or generic decoded JSON, not endpoint contracts.

#### Agent Action

When adding a `jsonDecode` hook, also add sequence metadata from request generation, handler headers, and response decrypt hooks. If the JSON event occurs before the first request cluster or with no nearby decrypt event, do not write it into endpoint specs. Update project docs to show this distinction.

#### Applies When

- Flutter/Dart AOT APK analysis.
- Broad `jsonDecode`, `JsonDecoder.convert`, or parser hooks are used.
- The goal is API endpoint/schema correlation.

#### Does Not Apply When

- The hook is attached to a narrow app-owned response parser already proven to receive only network responses.
- The endpoint path and response body are visible in the same high-level request/response object.
- The app has no local cache/session restore path and this has been validated.

#### Validation

- A schema is API-confirmed only when it appears after a matching request/decrypt cluster in repeated captures.
- Startup/local schemas are labeled separately when `nearestReqSeq=0` or no request/decrypt context exists.
- No raw JSON values are written to reusable docs.

#### Promotion Target

- `WORKFLOW.md`
- `DOCUMENTATION.md`

#### Required Linked Updates

- Update response decode workflow to require sequence context before promoting schema-only JSON to endpoint response evidence.
