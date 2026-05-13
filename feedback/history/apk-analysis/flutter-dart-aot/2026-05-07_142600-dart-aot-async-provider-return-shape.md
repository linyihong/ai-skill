> Follow [shared rules](../../../../shared-rules/README.md), [dependency-reading](../../../../shared-rules/dependency-reading.md), [neutral-language](../../../../shared-rules/neutral-language.md), [goal-action-validation](../../../../shared-rules/goal-action-validation.md), and [feedback-lessons](../../../../shared-rules/feedback-lessons.md). This lesson records only generalized guidance.
# Extracted — See [`analysis/apk/workflows/frida-hook-flow.md`](../../../../analysis/apk/workflows/frida-hook-flow.md)

### 2026-05-07 - Dart AOT Async Provider Return Shape

Status: candidate

#### One-line Summary

When a Dart AOT provider is async, a function-exit hook may see only a Future object; validate the awaited value through static continuation flow or the downstream consumer instead.

#### Human Explanation

For Flutter/Dart AOT, high-semantic provider names can look like normal getters but compile as async state machines. Hooking function entry/exit is still useful for sequencing, but `onLeave` can report a pointer/Future shape rather than the final string, token, config, or language value. Treat that as an async boundary, not a failed decoder.

#### Trigger

- A high-semantic Dart provider hook attaches and logs before the target request.
- The return hook only shows a pointer or Future-like object.
- Static call edges show SharedPreferences, mapping, await helpers, or continuation stubs below the provider.

#### Evidence

- Tool: Frida function return hooks, Dart AOT asm, call edges, static string references.
- Sanitized excerpt: a provider ran before request signing, but its return hook only produced a pointer shape; static asm showed the provider awaited persisted settings, mapped the result, and a downstream cached field was later consumed by request normalization.
- Evidence path: project-specific offsets and logs stay in `<PROJECT_ROOT>/capture/` and `<PROJECT_ROOT>/docs/`.

#### Generalized Lesson

Do not conclude a provider's returned value is unreadable or absent just because `onLeave` sees a pointer. For async Dart AOT providers, the function return may be a Future. Use static async continuation flow, call_edges, awaited helper calls, and the downstream consumer field or request boundary to validate semantics.

#### Agent Action

If a provider return hook logs only `ptr`, inspect whether the provider is async. Read its asm for storage reads, mapping helpers, fallback constants, and continuation calls. Prefer validating the final value shape at the consumer boundary or cached field after await. Keep raw values out of logs.

#### Goal / Action / Validation

- Goal: Avoid misclassifying async provider evidence when the direct return hook observes a Future.
- Action: Combine static async-chain review with a stable consumer-boundary hook and document the return-hook limitation.
- Validation or reference source: A valid run should show the provider executes before the target request, the downstream consumer receives the expected key/field shape, and default hooks remain stable after disabling diagnostic return hooks.

#### Applies When

- Flutter/Dart AOT providers use SharedPreferences, network/bootstrap state, locale/config mappings, or other awaited helpers.
- Function exit hooks show only pointer/Future shape.
- The requested proof is provider semantics or timing rather than raw value recovery.

#### Does Not Apply When

- The provider is synchronous and returns a directly decodable value.
- Raw value recovery is required and static semantics are insufficient.
- The downstream consumer cannot be identified.

#### Validation

Confirm the diagnostic hook is disabled by default after the probe, and keep project docs explicit that the return hook observed an async/Future boundary while static/consumer evidence provides the semantic proof.

#### Promotion Target

- `techniques/flutter-dart-aot/`
- `WORKFLOW.md`

#### Required Linked Updates

- Checked `feedback_history/flutter-dart-aot/README.md` and root `feedback_history/README.md`; both should index this lesson.
- No immediate promotion: keep as `candidate` until repeated across another Flutter/Dart AOT target.
