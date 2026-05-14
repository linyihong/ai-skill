> Follow [shared rules](../../../../enforcement/README.md), [dependency-reading](../../../../enforcement/dependency-reading.md), [neutral-language](../../../../enforcement/neutral-language.md), [goal-action-validation](../../../../enforcement/goal-action-validation.md), and [feedback-lessons](../../../../enforcement/feedback-lessons.md). This lesson records only generalized guidance.
# Extracted — See [`analysis/apk/workflows/frida-hook-flow.md`](../../../../analysis/apk/workflows/frida-hook-flow.md)

### 2026-05-07 - Entry Boundary Before Unstable AOT Hooks

Status: candidate

#### One-line Summary

When Flutter/Dart AOT callsite or broad collection hooks are unstable, move evidence gathering to stable function-entry boundaries plus static string/asm review before trying more invasive hooks.

#### Human Explanation

Dart AOT analysis often tempts agents to hook the exact `BL` instruction or a broad `Map._set` method after static xrefs show where a key or value is written. Those probes can be useful in short windows, but they may be rejected by Frida or crash hot Flutter paths. Stable function-entry hooks usually give enough sequence evidence to narrow the boundary, while static strings and asm explain the likely normalization stage without destabilizing the app.

#### Trigger

- A broad Dart collection hook causes app crashes.
- Frida rejects a non-function-entry callsite hook.
- Static xrefs show a likely request/key injection site, but the upstream value provider is still unclear.

#### Evidence

- Tool: Frida function-entry hooks, Dart AOT asm/string references, short-window experimental hooks.
- Sanitized excerpt: broad collection hooks observed request-key writes but destabilized the process; non-entry callsite hooks were rejected; a function-entry hook on the surrounding request-normalization function remained stable and established the order before signing/header generation.
- Evidence path: project-specific logs and offsets stay in `<PROJECT_ROOT>/capture/` and `<PROJECT_ROOT>/docs/`.

#### Generalized Lesson

Treat Dart AOT callsites and broad collection methods as experimental probes, not default instrumentation. Prefer high-semantic function-entry hooks at request normalization, interceptor, storage, decrypt, or response boundaries, and use static asm/string refs to decide the next narrower target.

#### Agent Action

If an exact callsite or broad collection hook is attractive, first add a feature flag and keep it disabled by default. Use stable entry hooks to prove call order and only run the invasive hook in a short, explicitly experimental window. If the hook is rejected or crashes, document that as a limitation and continue with static boundary narrowing rather than retrying increasingly broad hooks.

#### Goal / Action / Validation

- Goal: Preserve app stability while still moving closer to the opaque value or request mutation source.
- Action: Hook stable function entries, correlate them with static strings/asm and downstream request-shape evidence, and gate unstable probes behind disabled flags.
- Validation or reference source: A stable run should reach target package/UI validation and show the expected call sequence; project docs should mark rejected/crashing hook paths as non-default evidence.

#### Applies When

- Analyzing Flutter/Dart AOT request construction, signing, storage, or response decode paths.
- Static AOT output identifies an internal callsite but Frida cannot safely attach there.
- Broad collection or helper hooks are too hot for normal UI replay.

#### Does Not Apply When

- The target function entry itself is the exact stable provider and can be hooked safely.
- A short crash-prone probe is explicitly needed and isolated from target UI evidence.

#### Validation

Confirm the final default hook script runs without crash, logs only sanitized shapes or keys, and reaches the target capture stop. Keep experimental logs separate from valid UI/API evidence when they crash or stop before package validation.

#### Promotion Target

- `techniques/flutter-dart-aot/`
- `WORKFLOW.md`

#### Required Linked Updates

- Checked `feedback_history/flutter-dart-aot/README.md` and the root `feedback_history/README.md`; both should index this lesson.
- No immediate promotion: keep as `candidate` until repeated across another Flutter AOT target.
