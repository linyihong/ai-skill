# Flutter / Dart AOT Feedback Lessons

Lessons for Flutter, Dart AOT, `libapp.so`, Dart object decoding, and Dart/native offset hook workflows.

| File | Status | Topic | Notes |
| --- | --- | --- | --- |
| `2026-05-01_133900-dart-aot-interceptor-strings-after-java-helper-miss.md` | validated | Dart AOT interceptor strings after Java helper miss | Route from missed Java helper hooks to Dart AOT. |
| `2026-05-01_140900-unflutter-aot-offset-hook-after-blutter-crash.md` | promoted | Unflutter AOT offset hook after blutter crash | Use static parser/output after blutter instability. |
| `2026-05-01_142000-exhaustive-java-okhttp-hooks-may-still-miss-flutter-business-http.md` | validated | 廣覆蓋 Java OkHttp 仍無業務 host 時應轉 Dart／native／pcap | Evidence-based switch away from Java OkHttp. |
| `2026-05-01_145948-dart-aot-callsite-hooks-are-not-function-hooks.md` | validated | Dart AOT Callsite Hooks Are Not Function Hooks | Treat callsites as navigation hints. |
| `2026-05-01_151551-schema-only-jsondecode-hook.md` | validated | Schema-only jsonDecode Hook | Schema-only Dart JSON decode logging. |
| `2026-05-01_153437-sequence-jsondecode-before-api-response.md` | validated | Sequence jsonDecode Before Calling It API Response | Sequence request/decrypt/json before attribution. |
| `2026-05-01_155200-dart-compressed-response-fields.md` | candidate | Dart AOT compressed response fields | Restore compressed pointers before field reads. |
| `2026-05-01_164741-dart-inline-onebyte-string-smi-length.md` | validated | Dart inline one-byte string Smi length | Decode Dart inline one-byte strings correctly. |
