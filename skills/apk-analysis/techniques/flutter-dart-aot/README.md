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

## Related Lessons

- `../../feedback_history/2026-05-01_133900-dart-aot-interceptor-strings-after-java-helper-miss.md`
- `../../feedback_history/2026-05-01_140900-unflutter-aot-offset-hook-after-blutter-crash.md`
- `../../feedback_history/2026-05-01_145948-dart-aot-callsite-hooks-are-not-function-hooks.md`
- `../../feedback_history/2026-05-01_151551-schema-only-jsondecode-hook.md`
- `../../feedback_history/2026-05-01_153437-sequence-jsondecode-before-api-response.md`
- `../../feedback_history/2026-05-01_155200-dart-compressed-response-fields.md`
- `../../feedback_history/2026-05-01_164741-dart-inline-onebyte-string-smi-length.md`
