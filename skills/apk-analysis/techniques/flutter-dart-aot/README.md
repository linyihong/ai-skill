# Flutter / Dart AOT Techniques

Use this category only when evidence points to Flutter or Dart AOT. Do not read it for a Java/Kotlin-only, WebView-only, plain native, or media-only task unless routing evidence leads here.

## When To Use

- APK contains `lib/<arch>/libflutter.so` and `lib/<arch>/libapp.so`.
- Java/OkHttp hooks miss business HTTP while pcap/native timing shows network activity.
- Static strings or logs mention `Dio`, `Interceptor`, Dart package paths, encrypted response handlers, or Dart AOT function names.
- Tools such as blutter/unflutter are needed to recover function maps, string refs, or native offsets.

## Read With

- Common routing: `../../WORKFLOW.md`.
- Tool setup: `../../TOOLS.md`.
- API documentation output: `../../DOCUMENTATION.md`.

## Core Guidance

- Start from high-semantic request/response boundaries: Dio options, interceptors, response decoder, token/header provider.
- If Java helper/plugin hooks miss but local proxy or route metadata shows encrypted headers, inspect Dart AOT strings and function refs before adding more Java hooks.
- Use blutter when it works; if it identifies the snapshot but crashes, keep failure evidence and switch to static parsers such as unflutter.
- Treat Dart AOT callsite `BL` addresses as navigation hints, not function hook entry points.
- Avoid broad hooks on global Dart runtime helpers unless using a short, filtered observation window.
- When decoding Dart strings or compressed object fields, validate layouts with limited private hexdumps, then turn noisy dumps off.

## Common Flow

1. Decode or unpack the APK and confirm `lib/<arch>/libapp.so` plus `libflutter.so`.
2. Generate pseudo source, object pool, function offsets, string refs, or call edges with Dart AOT tooling. If `blutter` identifies the snapshot but crashes, keep the failure evidence and switch to a static parser such as `unflutter`.
3. Search for host/base URL clues, `Dio`, `HttpClient`, `RequestOptions`, `Interceptor`, `encrypt`, `decrypt`, `AES`, `base64`, `hash`, and response interceptor names.
4. Hook request options first. If function PCs are known, attach only to a small set of app-owned request/sign/encrypt/decrypt functions using `libapp.so` base + PC.
5. Hook response decode/decrypt return values before trying to reconstruct TLS/socket bytes.
6. If Dart string decoding fails, collect a limited private hexdump to infer layout, fix the decoder, then disable noisy dumps.
7. Align raw wrapper and decrypted payload into a sanitized fixture.

Success shape:

```text
request hook:
  method / baseUrl / path / headers / query

response decode hook:
  decrypted JSON/string
```

## Pitfalls

- Do not treat call edges pointing to internal `BL` instructions as function hook entry points.
- Do not start with global Dart runtime or collection helpers such as `LinkedHashMap._set`; they are hot, noisy, and can destabilize the app.
- If local proxy or Netty evidence shows custom encrypted/signed headers but Java helper hooks miss, treat Java plugin/helper code as bridge/setup until proven otherwise and inspect Dart AOT interceptors.

## Related Lessons

- `../../feedback_history/flutter-dart-aot/2026-05-01_133900-dart-aot-interceptor-strings-after-java-helper-miss.md`
- `../../feedback_history/flutter-dart-aot/2026-05-01_140900-unflutter-aot-offset-hook-after-blutter-crash.md`
- `../../feedback_history/flutter-dart-aot/2026-05-01_145948-dart-aot-callsite-hooks-are-not-function-hooks.md`
- `../../feedback_history/flutter-dart-aot/2026-05-01_151551-schema-only-jsondecode-hook.md`
- `../../feedback_history/flutter-dart-aot/2026-05-01_153437-sequence-jsondecode-before-api-response.md`
- `../../feedback_history/flutter-dart-aot/2026-05-01_155200-dart-compressed-response-fields.md`
- `../../feedback_history/flutter-dart-aot/2026-05-01_164741-dart-inline-onebyte-string-smi-length.md`
