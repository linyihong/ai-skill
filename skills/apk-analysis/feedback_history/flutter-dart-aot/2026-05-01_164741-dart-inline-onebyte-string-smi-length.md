> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/frida-hook-flow.md`](../../../../analysis/apk/workflows/frida-hook-flow.md)

### 2026-05-01 - Dart inline one-byte string Smi length

Status: validated

#### One-line Summary

Dart AOT inline one-byte strings may store length as a Smi at object `+8` and bytes at `+16`; reading `+8` as raw32 can make valid strings look undecodable.

#### Human Explanation

When a Dart object class repeatedly appears as an undecoded string-like object, do not assume it is a wrapper just because the usual raw32 length decoder fails. Some AOT one-byte string layouts store a tagged Smi length. If that tagged value is interpreted as raw length, the decoder reads past the real string into adjacent object data, lowers printable percentage, and rejects the value.

#### Trigger

Use this when Frida summaries show:

- a stable Dart class id for string-like values;
- shallow field probes contain ASCII chunks at consecutive offsets;
- the first field looks like a tagged small integer length;
- existing raw32 string layout candidates fail or return `undecoded`.

#### Evidence

- Tool: Frida native offset hooks on Dart AOT request/decrypt/JSON functions.
- Sanitized excerpt: a string-like object initially showed `cid=<string-like-cid>` with shallow fields containing ASCII chunks; after adding a `smi32@8/16` candidate, the same object decoded to schema-only summaries such as `layout=utf8:smi32@8/16`, `features=jsonLike`, and top-level keys/types.
- Evidence path: project-private capture logs under `<PROJECT_ROOT>/capture/`; do not copy raw values into reusable skill docs.

#### Generalized Lesson

For Dart AOT string decoding, try both raw and tagged-Smi length variants for common inline layouts. A useful candidate order is:

1. `lenOff=8`, `dataOff=16`, UTF-8/one-byte, Smi length.
2. `lenOff=8`, `dataOff=16`, UTF-8/one-byte, raw32 length.
3. Other observed UTF-16 or shifted layouts.

Use shallow field probes to validate the hypothesis before broadening raw decoding.

#### Agent Action

When an object looks string-like but remains undecoded, add a shallow, sanitized field summary first. If fields reveal printable ASCII chunks and a plausible Smi length, add a `smi32@8/16` decode candidate and re-run a short capture. Keep output to length/hash/features/schema shape unless the user explicitly requests raw private capture.

#### Applies When

- Flutter/Dart AOT reverse engineering with native Frida hooks.
- Object fields or inline bytes reveal ASCII/JSON/path/base64-like material.
- Sanitized schema or payload shape is enough for API documentation.

#### Does Not Apply When

- The value is clearly a normal heap object, list, map, or typed-data buffer rather than inline string bytes.
- The target runtime uses a different object layout verified by disassembly or heap dumps.
- Raw values are not authorized to be collected.

#### Validation

The decode is credible when the decoded length matches the tagged Smi, printable percentage is high, hash is stable across the same object handoff, and downstream hooks (`jsonDecode`, decrypt return, or response handoff) show the same string hash/shape.

#### Promotion Target

- `WORKFLOW.md`
- `TOOLS.md`

#### Required Linked Updates

- Project docs should record only sanitized length/hash/schema/segment shape.
- Reusable skill docs should not include target-specific paths, services, hosts, tokens, or raw payload content.
