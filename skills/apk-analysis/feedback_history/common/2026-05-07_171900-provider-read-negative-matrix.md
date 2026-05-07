> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-07 - Provider Read Negative Matrix

Status: candidate

#### One-line Summary

When classifying an opaque provider's bad/empty behavior, prefer a short no-writeback provider-read override over mutating outgoing request maps directly.

#### Human Explanation

Some opaque request fields are derived from app settings or cached provider state before signing. If the goal is to learn whether empty, stale, or bad provider values are tolerated, rejected, or routed into session recovery, changing the final query map can bypass the same initialization and signing path you are trying to validate. A provider-read override keeps the app-owned normalization, signing, transport, and decode path intact while avoiding persistent local storage changes.

#### Trigger

Use this when an opaque request field has a known upstream provider, equality/parity with outgoing requests is already established, and the remaining blocker is a negative matrix for empty/stale/bad provider values.

#### Evidence

- Tool: Frida Java/Android storage hook or equivalent provider-read hook in an authorized short window.
- Sanitized excerpt: Store only mode (`empty`, `bad-fixed`, `stale-candidate`), value-present flags, length/hash fingerprints, request key set, wrapper code class, and UI/session reaction class.
- Evidence path: Keep target-specific run logs and operation ids in `<PROJECT_ROOT>` analysis docs; do not copy app names, endpoints, raw values, or live result details into this lesson.

#### Generalized Lesson

Do not start a provider negative matrix by editing the final query map unless no upstream provider-read boundary exists. First override the provider read result without writing storage, then let the app build, sign, send, and decrypt the request. Keep the override disabled by default, bounded by a small read count, and scoped to harmless read-only routes.

#### Agent Action

When planning this diagnostic:

- Add a disabled-by-default flag and explicit mode.
- Override only reads, not writes, unless the test explicitly requires persisted state.
- Log only redacted fingerprints and status classes.
- Run each mode in a separate short capture.
- Restart with the flag disabled and verify the app/session returns to normal before documenting conclusions.

#### Goal / Action / Validation

- Goal: Classify provider-value failure behavior without polluting storage or bypassing app-owned signing/decode paths.
- Action: Override the upstream provider read result for one authorized read-only route and record redacted request/response/UI classes.
- Validation or reference source: A valid run proves the modified provider read is reflected in the signed request key set or equality diagnostic, and a follow-up clean run proves recovery with the override disabled.

#### Applies When

- The provider source is known or strongly bounded.
- The app signs or normalizes requests after reading the provider.
- The route is read-only and authorized for controlled negative testing.
- The evidence can be recorded without raw provider values or raw responses.

#### Does Not Apply When

- The provider source is unknown and the diagnostic would be guessing.
- The only available hook would change persistent account/device/session state without a restore plan.
- The route performs writes, purchases, messaging, or other state-changing actions.
- The user has not authorized negative replay or tamper-style testing.

#### Validation

Confirm all of the following before promoting a result:

- The override is disabled by default in committed tooling.
- The override does not write back to storage.
- The captured request used the app-owned signing/decode path.
- The app recovered after a clean restart or documented restore step.
- Project docs store only redacted mode/status/fingerprint classes.

#### Promotion Target

- `TOOLS.md`
- `WORKFLOW.md`
- `DOCUMENTATION.md`

#### Required Linked Updates

- Updated the common feedback indexes so this candidate lesson is discoverable.
- Promotion into main workflow is deferred until at least one additional authorized project validates the no-writeback provider-read pattern.
- Project-specific evidence must remain in project docs per [`reusable-guidance-boundary`](../../../../shared-rules/reusable-guidance-boundary.md).
