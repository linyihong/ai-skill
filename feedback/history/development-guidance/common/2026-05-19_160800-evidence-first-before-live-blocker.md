> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-19 - Evidence-first before live blocker

Status: candidate

#### One-line Summary

Before declaring a live/integration prerequisite missing, first search project docs, sanitized API entries, captures, and sibling modules for evidence that can derive a safe default.

#### Human Explanation

Live SDK gates often mix true private-adapter blockers with values that were already classified in earlier analysis. Treating every unset env var as equally missing can make a runnable smoke look more blocked than it is and can contradict existing domain evidence.

The agent should distinguish "unknown and private" from "known, derivable, or shared with a sibling flow" before reporting blockers.

#### Trigger

- A live test, smoke runner, SDK adapter, or integration gate reports missing host, locale, route, signing, decrypt, identity, or session material.
- The value name looks opaque but may be documented as a language, environment, route, or shared runtime constant.
- A sibling module in the same product already has working live evidence for a similar parameter.

#### Evidence

- Tool: project docs, sanitized API docs, capture summaries, existing live tests, and sibling-module env resolvers.
- Sanitized excerpt: `l` looked like a required opaque env value, but docs and captures classified it as server-language code and existing SDK code used `zh-cn`.
- Evidence path: keep project-specific route names, raw values, captures, and secrets in project docs only.

#### Generalized Lesson

Classify each live prerequisite before calling it missing:

- Derivable shared constant: default it and document the source.
- Documented but overrideable runtime value: provide a conservative default plus env override.
- Sanitized route/service/signing/decrypt/identity material: keep it as private-adapter-required.
- Negative matrix or alternate locale behavior: do not infer; require authorized dynamic signing/decrypt evidence.

#### Agent Action

1. Search the owning docs and sibling live implementations for the parameter name, provider, and captured examples.
2. If evidence supports a default, update code/docs/tests so the gate no longer asks the user to hand-fill it.
3. If evidence only proves presence or hash class, keep the blocker and state exactly what evidence is still missing.
4. Re-run the focused gate and verify the skipped/missing list changed as expected.

#### Applies When

- Maintaining live/integration gates for SDKs, private adapters, CLI smoke tests, or replay tools.
- Working from APK/API analysis where public docs intentionally redact secrets but preserve key names and provider classes.

#### Does Not Apply When

- No project evidence exists for the value.
- The value is security-sensitive, account/session/device-bound, or route-specific raw material.
- The task is an intentionally offline parser/mock test.

#### Validation

- Skip/setup-failure output no longer lists derivable defaults.
- Docs cite the evidence class and override mechanism.
- Remaining blockers are limited to truly private or unproven material.
