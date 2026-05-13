> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/frida-hook-flow.md`](../../../../analysis/apk/workflows/frida-hook-flow.md)

### 2026-05-01 - Unflutter AOT offset hook after blutter crash

Status: promoted

#### One-line Summary

`blutter` 能識別 Dart snapshot 但在載入／分析時 SIGSEGV 時，可改用 `unflutter` 取得 Dart AOT function PC，再用 Frida native offset hook 驗證 request/sign/decrypt 流程。

#### Human Explanation

Flutter release APK 的業務請求常在 Dart AOT interceptor 內組裝；若 Java HTTP／Java plugin hooks 已排除，下一步需要 function offsets。`blutter` 很有用，但在新 Dart 版本或特定 snapshot 上可能 crash。此時不要卡在同一工具；可用不嵌入 Dart VM 的 static parser 產生 `functions.jsonl`／`call_edges.jsonl`／`string_refs.jsonl`，再挑高語意函式做 Frida native hook。

#### Trigger

- Local proxy／Netty hook 已看到自訂加密 header。
- Java helper／plugin hook 未命中 header 生成點。
- `libapp.so` 字串已出現 `RequestInterceptor`、`_generate...`、`_encrypt...`、header 名稱或 Dart package path。
- `blutter` 可偵測 Dart 版本，但 full 或 `--no-analysis` 仍 SIGSEGV。

#### Evidence

- Tool: `unflutter`, Frida native `Interceptor.attach`, Dart OneByteString best-effort decoder.
- Sanitized excerpt: AOT function map identified an interceptor chain like `_generateEhHeader -> _encodeHeader -> _encryptAES`; dynamic hook decoded plaintext shaped like `<prefix>|<timestamp>|<nonce>|<path>|<query/form material...>` and an encrypted/base64-like return.
- Evidence path: target-specific raw logs remain under `<PROJECT_ROOT>/capture/`; reusable lesson stores only generalized structure.

#### Generalized Lesson

When `blutter` crashes after confirming the Dart snapshot version, use an independent AOT parser such as `unflutter` to recover function PCs and string refs. Hook only the narrow candidate functions, not every Dart stub. For Dart 3.x compressed-pointer ARM64 samples, check OneByteString layouts empirically with a tiny hexdump; some strings may store raw byte length at object offset `+0x08` and data at `+0x10`, rather than a Smi length at later offsets.

#### Agent Action

1. Record the `blutter` failure as evidence, then pivot to another AOT parser instead of retrying random flags indefinitely.
2. Search recovered `functions.jsonl`, `call_edges.jsonl`, and `string_refs.jsonl` for `Interceptor`, header names, `encrypt`, `base64`, `sign`, `decrypt`, and Dart package paths.
3. Build a small Frida script from `libapp.so` base + recovered PC, and log sanitized argument/return summaries.
4. If string decoding fails, dump only a few object bytes in private capture to infer layout, then disable hexdump for normal runs.
5. Keep raw values in project-private capture only; documentation should use lengths, hashes, shapes, and redacted examples.

#### Applies When

- Authorized Flutter/Dart AOT APK analysis.
- You have `libapp.so` and ideally `libflutter.so`.
- Static strings or prior hooks point to a small set of candidate Dart functions.

#### Does Not Apply When

- The app is not Flutter AOT or the issue is clearly in Java/Kotlin.
- You need source-level decompilation before any dynamic validation; offset hooks can often validate faster.
- The recovered function PC is not stable for the installed APK build being tested.

#### Validation

- `unflutter doctor` or equivalent confirms snapshot parsing support.
- Candidate function PCs are present in the tool output and match surrounding call edges/string refs.
- Frida hook hits during the same UI/network action that produces the target header or decoded response.
- Hooked plaintext/return shape correlates with handler-level headers, request paths, or response decrypt points.

#### Promotion Target

- `WORKFLOW.md`
- `TOOLS.md`
