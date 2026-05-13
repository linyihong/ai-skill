> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/frida-hook-flow.md`](../../../../analysis/apk/workflows/frida-hook-flow.md)

### 2026-05-01 - Dart AOT Compressed Response Fields

Status: candidate

#### One-line Summary

Hook Dart/Dio response handoff objects by reading compressed fields from the tagged object address, not from a generic `x28` heap-base assumption.

#### Human Explanation

Dart AOT response interceptors often decrypt a payload, write it back into a response object's data field, then pass the object to the next handler. If a Frida hook only logs function arguments, the decrypted return value may appear disconnected from the downstream response. Reading object fields can bridge that gap, but compressed pointer reconstruction is easy to get wrong.

In observed ARM64 AOT code, field loads used instructions like `LDUR Wn, [taggedObj,#off]`, meaning the compressed value should be read from the tagged object address at the same offset. Reconstructing the full pointer from `x28` produced bogus `0x800...` addresses in this case. Using the containing object's upper 32 bits plus the compressed field value produced plausible Dart tagged pointers and class ids.

#### Trigger

You can hook a response decrypt function and a downstream handler such as `ResponseInterceptorHandler.next`, but the handler argument only shows an undecoded response object and the decrypted string/wrapper does not appear in obvious function arguments.

#### Evidence

- Tool: Frida native offset hooks on Dart AOT functions.
- Sanitized excerpt: decrypt return and downstream response object fields were logged with event sequence/timestamps; wrong reconstruction showed `0x800...` pseudo-pointers, corrected reconstruction produced normal Dart heap pointers and class ids.
- Evidence path: project-local `capture/dart_response_next_fields_*.log` and `capture/dart_response_payload_stats_*.log` (raw project evidence, not reusable skill content).

#### Generalized Lesson

When inspecting Dart AOT object fields from Frida:

1. Use the assembly addressing mode as ground truth. If the code says `LDUR Wn, [taggedObj,#off]`, read from `taggedObj.add(off)`, not from an untagged object plus guessed offsets.
2. Reconstruct compressed heap refs using the containing object's high address bits when `x28` does not match observed heap pointers.
3. Preserve only sanitized summaries: class id, string length/hash/features, segment shape, and event sequence.
4. Validate field identity by observing a before/after transition across the decrypt call and downstream handler call.

#### Agent Action

For Flutter/Dart AOT response analysis, after finding a decrypt function, hook the next response handoff function and log sanitized response object field summaries. If field pointers look like `0x800...` or all decode as Smi-like values, revisit pointer reconstruction before concluding the field is not useful.

#### Applies When

- Flutter/Dart AOT ARM64 code uses compressed pointers.
- Static disassembly shows response object field stores/loads around decrypt and handler calls.
- You need to prove decrypted output is handed to upper layers without dumping raw payloads.

#### Does Not Apply When

- The target uses uncompressed pointers or a different architecture/layout.
- You already have a high-level safe hook that exposes the parsed response schema directly.
- The analysis scope does not allow dynamic instrumentation or object memory inspection.

#### Validation

The corrected field reader should produce heap pointers in the same address range as other Dart objects and stable class ids across repeated handler calls. A candidate data field should change after decrypt and remain stable across repeated downstream `next` calls for the same response object.

#### Promotion Target

- `WORKFLOW.md`
- `TOOLS.md`

#### Required Linked Updates

- Project docs should record only sanitized response handoff evidence.
- No immediate promotion until this compressed-field method is validated on at least one more Dart AOT sample or a second response-object layout.
