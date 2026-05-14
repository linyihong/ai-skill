> 遵守 [共用規則索引](../../../../enforcement/README.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/frida-hook-flow.md`](../../../../analysis/apk/workflows/frida-hook-flow.md)

### 2026-05-01 - Dart AOT Interceptor Strings After Java Helper Miss

Status: validated

#### One-line Summary

當 Netty/local proxy 已看到加密 header，但 Java plugin/helper hook 未命中生成點時，應轉向 `libapp.so` 的 Dart AOT interceptor 字串與函式名線索。

#### Human Explanation

Flutter APK 可能用 Java plugin 提供 host、時間或橋接能力，但真正的 request signing / header assembly 仍在 Dart AOT 層完成。若 local proxy handler 已穩定看到自訂 header，而同窗 hook Java helper（例如 AES/RC2/getNMKey/query-map helper）沒有命中，不要繼續把 Java 當主線。

此時先抽出 `libapp.so` 做字串盤點，搜尋 interceptor log、Dio/HttpManager 檔名、sign metadata/result、header 名稱與函式名。Dart AOT 常保留 `package:<app>/.../*.dart` 路徑與 `_generate...@<hash>` 形式的函式名，足以決定下一步要用 blutter/AOT xref/offset hook 哪些函式。

#### Trigger

- Netty/local proxy request 已可見 method/path/header shape。
- 自訂 header 長度與內容每次變動，顯示它是 request-specific signing/encryption output。
- Java plugin/helper hooks 成功安裝，但同窗業務 header 出現時未命中生成函式。
- 靜態 `libapp.so` 有 request interceptor、Dio、sign 或 header 相關字串。

#### Evidence

- Tool: Frida on local proxy handler plus Java plugin/helper hooks.
- Tool: `libapp.so` string extraction.
- Sanitized finding:
  - local proxy sees custom encrypted/signature headers with empty body.
  - Java helper hooks remain quiet during those request events.
  - `libapp.so` contains Dart package path for request interceptor plus sign-generation function-name strings.
- Evidence path: project-local `capture/` and sanitized project docs only. Do not copy raw header values, keys, hostnames, tokens, or private paths into reusable skill docs.

#### Generalized Lesson

In Flutter APKs, Java plugin classes may be bridge/config providers rather than the business request signer. If Java hooks do not correlate with observed request headers, treat Dart AOT as the likely owner. Use static strings to find the Dart source path and candidate function names before attempting lower-level byte hooks.

#### Agent Action

When this pattern appears:

1. Keep the local proxy/Netty hook active and log only header names, lengths, hashes, and body length in shared notes.
2. Hook obvious Java helpers once to test correlation, not as an indefinite search loop.
3. Extract or locate `libapp.so`.
4. Search strings for `Interceptor`, `Dio`, `HttpManager`, header names, `sign metadata`, `sign result`, `encrypt`, `decrypt`, and package paths.
5. Record candidate Dart function names and source paths.
6. Move next to blutter/AOT xref/Frida offset hooks for those Dart functions.

#### Applies When

- APK is Flutter/Dart AOT.
- Business requests use Dio/dart:io and pass through a local Java/Netty proxy.
- Request body is empty but custom headers carry encrypted/signature-like data.
- Java bridge methods exist but do not fire during header generation.

#### Does Not Apply When

- Java OkHttp/Retrofit hooks already capture the full request builder and signer.
- The header is added by a Java interceptor or native C/C++ library with clear hook evidence.
- `libapp.so` has no useful Dart package/function strings and AOT tooling is unavailable.

#### Validation

Validated when the same capture window shows custom headers at the local proxy while Java signing/helper hooks do not fire, and `libapp.so` contains matching Dart interceptor/signing strings.

#### Promotion Target

- `WORKFLOW.md`
- `TOOLS.md`
