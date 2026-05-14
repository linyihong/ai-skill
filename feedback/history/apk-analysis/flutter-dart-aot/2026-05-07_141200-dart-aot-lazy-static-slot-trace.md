> Follow [shared rules](../../../../enforcement/README.md), [dependency-reading](../../../../enforcement/dependency-reading.md), [neutral-language](../../../../enforcement/neutral-language.md), [goal-action-validation](../../../../enforcement/goal-action-validation.md), and [feedback-lessons](../../../../enforcement/feedback-lessons.md). This lesson records only generalized guidance.
# Extracted — See [`analysis/apk/workflows/frida-hook-flow.md`](../../../../analysis/apk/workflows/frida-hook-flow.md)

### 2026-05-07 - Dart AOT Lazy Static Slot Trace

Status: candidate

#### One-line Summary

When a Dart AOT lazy static getter does not fire in the target window, trace the cached ISO-group static slot from asm and hook the consumer boundary instead of assuming the getter is the live provider.

#### Human Explanation

Dart AOT code often lowers static fields into ISO-group slots with a lazy getter only on first initialization. By the time a request interceptor or feature handler runs, the hot path may read the cached slot directly and bypass the lazy getter. A Frida hook on the getter can therefore look silent even when the value is actively used.

#### Trigger

- A static xref shows a common lazy static getter near a value provider.
- The getter hook attaches successfully but logs nothing in a later request window.
- The consuming function reads `THR` / `ISO_GROUP` offsets before loading object fields.

#### Evidence

- Tool: static Dart AOT asm review, function-entry Frida hooks, shape-only field diagnostics.
- Sanitized excerpt: a request normalizer loaded a cached object from an ISO-group slot and read a field into an outgoing parameter; the matching lazy static getter only ran during first initialization and was bypassed once the slot was populated.
- Evidence path: project-specific offsets, field numbers, and logs stay in `<PROJECT_ROOT>/capture/` and `<PROJECT_ROOT>/docs/`.

#### Generalized Lesson

Treat lazy static getter hooks as initialization probes, not proof that a static value is or is not used later. If the asm shows `THR+... -> ISO_GROUP+slot` loads, record the slot and field path, then validate at the consumer boundary with sanitized shape logging.

#### Agent Action

When a lazy static getter is silent, inspect the consuming asm for the cached slot read. Prefer documenting the slot/field path and hooking the high-semantic consumer entry. If a diagnostic field hook is needed, gate it behind a disabled-by-default feature flag and log only shapes, booleans, lengths, or key names.

#### Goal / Action / Validation

- Goal: Avoid misclassifying a cached static field as unused just because the lazy getter no longer runs.
- Action: Trace ISO-group slot reads statically, then validate the consumer sequence with stable entry hooks or short diagnostic field probes.
- Validation or reference source: A valid run should reach target package/UI validation with default hooks quiet, while diagnostic runs should show the consumer reading the expected cached slot/field without raw values.

#### Applies When

- Analyzing Flutter/Dart AOT request construction, session state, signing, or singleton-backed caches.
- Static asm includes `THR` and `ISO_GROUP` loads near a value used by a request or decoder.
- The value is initialized before the target UI/API window.

#### Does Not Apply When

- The provider is a normal function call that executes every target request.
- The static slot or field path cannot be identified from asm.
- Raw value recovery is required and shape-only evidence is insufficient.

#### Validation

Confirm the final default hook script disables diagnostic slot/field probes and still reaches target validation. Keep any raw pointer or field-shape diagnostics out of public docs except as sanitized summaries.

#### Promotion Target

- `techniques/flutter-dart-aot/`
- `WORKFLOW.md`

#### Required Linked Updates

- Checked `feedback_history/flutter-dart-aot/README.md` and root `feedback_history/README.md`; both should index this lesson.
- No immediate promotion: keep as `candidate` until repeated across another Flutter/Dart AOT target.
