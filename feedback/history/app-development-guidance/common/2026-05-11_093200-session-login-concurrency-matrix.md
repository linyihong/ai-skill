> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/app-development-guidance/execution-flow.md`](../../../../workflow/app-development-guidance/execution-flow.md)

### 2026-05-11 - Session login concurrency matrix

Status: candidate

#### One-line Summary

SDKs that self-login must test multi-process startup, per-device session reuse, refresh single-flight, and rate-limit behavior before release.

#### Human Explanation

Single-process lazy initialization can pass while real users open several app instances or tools at once. If sessions are memory-only and login HTTP is not coordinated across processes, every startup may hit the login endpoint and trigger server-side frequency controls.

#### Trigger

Users report "login too frequently" or rate-limit errors when multiple clients start at nearly the same time.

#### Evidence

- Tool: SDK code review and concurrency test planning.
- Sanitized excerpt: device identity is persisted, but session is process memory and login occurs outside the device-pool lock.
- Evidence path: SDK auth/session/device classes in the project repo.

#### Generalized Lesson

Auth-capable SDKs need an explicit session concurrency matrix: cache hit avoids login, cache miss single-flights login, expired session single-flights refresh, waiting requests do not continue with known-expired material, and rate-limit responses produce cooldown/backoff.

#### Agent Action

When adding self-login to an SDK, design token/session persistence and refresh coordination together. Add tests for concurrent first calls, process-like multi-client startup, refresh in flight, failed refresh, and server rate-limit classification.

#### Goal / Action / Validation

- Goal: prevent login storms and stale-token API fanout.
- Action: add session cache, login/refresh single-flight, and concurrency tests.
- Validation or reference source: test matrix proves concurrent clients trigger one login/refresh or wait for the shared result.

#### Applies When

- SDK/client generates or refreshes its own API session.
- Multiple app/process instances can share a device pool or identity material.

#### Does Not Apply When

- The host app owns session lifecycle and passes fresh tokens into the SDK for every request.

#### Validation

Use offline mocked HTTP counters for deterministic concurrency, then optional authorized live probes for server rate-limit behavior.

#### Promotion Target

- `WORKFLOW.md`
- `CHECKLIST.md`
- `process/README.md`

#### Required Linked Updates

- Indexed in `feedback_history/common/README.md`.
- Concrete endpoint names, tokens, and product incidents stay in project docs.
