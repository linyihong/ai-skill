> Follow [shared rules](../../../../enforcement/README.md), [dependency-reading](../../../../enforcement/dependency-reading.md), [neutral-language](../../../../enforcement/neutral-language.md), [goal-action-validation](../../../../enforcement/goal-action-validation.md), and [feedback-lessons](../../../../enforcement/feedback-lessons.md). This lesson records only generalized guidance.
# Extracted — See [`analysis/apk/workflows/frida-hook-flow.md`](../../../../analysis/apk/workflows/frida-hook-flow.md)

### 2026-05-07 - Dio RequestOptions URI Keyset Probe

Status: candidate

#### One-line Summary

For Flutter/Dio AOT flows, hook `RequestOptions.get:uri` helpers and log only query key sets to prove whether parameters existed before or after interceptor normalization.

#### Human Explanation

When an interceptor mutates query parameters before signing, downstream request/header hooks may only show the final key set. A targeted URI keyset probe can compare pre-normalization and post-normalization states without printing raw query values: filter `Transformer.urlEncodeQueryMap` and `Uri.parse` to the `RequestOptions.get:uri` callsite, then summarize only key names and presence flags.

#### Trigger

- A required query parameter is visible at signing/header generation but its upstream source is unclear.
- Static AOT analysis shows a Dio `RequestOptions.get:uri` path with `urlEncodeQueryMap` / `Uri.parse`.
- Broad map hooks are too noisy or unstable for normal UI replay.

#### Evidence

- Tool: Frida function-entry hooks on `RequestOptions.get:uri`, `Transformer.urlEncodeQueryMap`, and `Uri.parse`, filtered by caller offset.
- Sanitized excerpt: pre-normalization URI summaries contained only baseline keys; post-normalization URI encode/parse summaries contained the opaque key immediately before signing/header generation.
- Evidence path: project-specific logs and offsets stay in `<PROJECT_ROOT>/capture/` and `<PROJECT_ROOT>/docs/`.

#### Generalized Lesson

Use Dio URI encode/parse helpers as a high-semantic boundary for query mutation timing. Log only sorted query key names, selected page/cursor values when safe, and boolean presence markers for opaque parameters. Do not log raw values, tokens, device ids, service strings, or signatures.

#### Agent Action

Add the URI helper probe behind a disabled-by-default feature flag because it can be verbose. Filter by `RequestOptions.get:uri` caller offsets so unrelated URI parsing does not pollute the capture. Compare the logged key sets before interceptor entry, after normalization, and at the signing/header boundary.

#### Goal / Action / Validation

- Goal: Determine whether a query parameter originates in the original request options or is added by an interceptor/normalizer.
- Action: Hook URI encode/parse helper function entries, filter to the request-options callsite, and emit key-set summaries only.
- Validation or reference source: A valid run should reach package/UI validation, show no raw values, and produce an ordered sequence across request options, interceptor, URI parse, and signing/header hooks.

#### Applies When

- The app uses Flutter/Dart AOT with Dio-like `RequestOptions`.
- Static call edges identify `Transformer.urlEncodeQueryMap` and `Uri.parse` below `RequestOptions.get:uri`.
- The question is parameter presence/timing rather than raw parameter value.

#### Does Not Apply When

- The parameter value itself must be recovered and key-set timing is insufficient.
- The app does not use Dio-like request options or the URI helper cannot be scoped safely.
- The helper path is so hot that even filtered key-set logging disrupts UI replay.

#### Validation

Run once with the helper enabled to collect evidence, then disable it by default and rerun a short capture to prove the default hook script remains stable and quiet.

#### Promotion Target

- `techniques/flutter-dart-aot/`
- `WORKFLOW.md`

#### Required Linked Updates

- Checked `feedback_history/flutter-dart-aot/README.md` and root `feedback_history/README.md`; both should index this lesson.
- No immediate promotion: keep as `candidate` until repeated across another Flutter/Dio target.
