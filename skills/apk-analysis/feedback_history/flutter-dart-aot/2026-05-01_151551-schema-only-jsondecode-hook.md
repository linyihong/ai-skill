> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/frida-hook-flow.md`](../../../../analysis/apk/workflows/frida-hook-flow.md)

### 2026-05-01 - Schema-only jsonDecode Hook

Status: validated

#### One-line Summary

回包明文只在 Dart decode 階段出現時，hook `jsonDecode` input 並只輸出 length/hash/top-level keys/types，可先取得安全 schema，而不落 raw JSON value。

#### Human Explanation

Flutter/Dart app 常把 response decrypt、wrapper handling、JSON parsing 都放在 Dart AOT 內。PC MITM 或 Java OkHttp hook 可能看不到明文 response body；即使 hook 到 decrypt function，return 也可能是 Future、wrapper、Map 或其他 Dart object，未必能直接當 String 讀。

此時 `jsonDecode` 是很好的高語意觀測點：它接收的通常已是解密後 JSON string。為了避免保存個資、token 或可重放資料，不要預設印 raw JSON；先只印字串長度、SHA-256、`jsonLike` 特徵、top-level key 與 key 的型別。這足以確認「明文在哪一層」、建立 API schema 初稿、並用 hash 與前後事件做同窗關聯。

#### Trigger

- Response decrypt function 命中，但 return 是 undecoded wrapper 或 async object。
- MITM/Java hook 看不到業務 response JSON。
- 靜態 AOT xref 顯示 app 會呼叫 `jsonDecode` / `JsonDecoder.convert` / `HttpManager._parseAndDecode`。

#### Evidence

- Tool: `unflutter` function map/call edges plus Frida native offset hooks.
- Sanitized excerpt: `jsonDecode:enter` receives Dart String inputs with `jsonLike`, length/hash, and schema-only `keys=... types=...`; no `value=` output.
- Evidence path: project-private Frida logs under `<PROJECT_ROOT>/capture/`.

#### Generalized Lesson

When response JSON is not visible in MITM/Java hooks, add a schema-only hook at Dart JSON parse boundaries before trying lower-level socket/TLS capture. Keep the output to:

- byte/char length;
- hash of the whole JSON string;
- top-level object keys or array first-item keys;
- type categories (`string`, `number`, `boolean`, `object`, `array[n]`, `null`);
- no field values.

If JSON parsing in the hook fails due to wrapper padding or NUL/trailing bytes, parse a sanitized copy for shape only, while hashing the original decoded string so correlation remains stable.

#### Agent Action

Next time a Flutter APK response decrypt hook returns wrapper-like objects, search for `jsonDecode`, `JsonDecoder.convert`, response manager parse/decode helpers, or equivalent function PCs. Hook those function entries and emit schema-only summaries first. Only collect raw JSON in private capture when explicitly needed and authorized.

#### Applies When

- Flutter/Dart AOT app.
- Response content is encrypted/wrapped before JSON parsing.
- The goal is API/schema discovery or path-to-response correlation.

#### Does Not Apply When

- The JSON is already visible safely in MITM or a high-level response object.
- The parse boundary receives binary/protobuf/msgpack rather than JSON string.
- Raw values are required for a narrow authorized fixture; in that case keep raw capture private and document only sanitized summaries.

#### Validation

- Hook output shows `jsonLike` strings entering `jsonDecode` or `JsonDecoder.convert`.
- Top-level keys/types are stable across repeated short captures.
- Raw values are absent from logs intended for documentation (`value=` should not appear).

#### Promotion Target

- `WORKFLOW.md`
- `DOCUMENTATION.md`

#### Required Linked Updates

- Add response-decode workflow guidance and documentation sanitization guidance for schema-only JSON summaries. Keep endpoint-specific keys and any raw values in project docs/capture, not reusable skill files.
