> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/software-delivery/execution-flow.md`](../../../../workflow/software-delivery/execution-flow.md)

### 2026-05-07 - Media Metadata Private Decrypt Boundary

Status: candidate

#### One-line Summary

When live media URLs, signed query values, encrypted files, or wrapped media keys are private, SDK core should preserve safe metadata and content blocks while media fetch/decrypt stays in a private adapter or media module.

#### Human Explanation

SDK teams often want parser and domain work to continue after analysis confirms media fields. The risky mistake is to fold encrypted media download, short-lived signed URLs, key unwrap, or app-specific decrypt into the public SDK core too early.

Keep the split explicit: the core SDK may expose media metadata, cover fields, raw JSON content blocks, descriptors, or resource references. The private adapter or media module owns host, auth query values, key unwrap, decrypt, cache invalidation, and live playback/download validation.

#### Trigger

- Analysis confirms media fields such as image URL, avatar URL, cover URL, HLS playlist URL, duration, resolution, signed query, wrapped key, or encrypted file shape.
- Raw media URLs, auth keys, decrypt keys, or playlist/key bytes cannot be committed.
- The downstream SDK already parses article/detail/comment models and could accidentally overclaim media download support.

#### Evidence

- Tool: analysis-to-SDK contract audit.
- Sanitized excerpt: core SDK models preserved cover strings and raw mixed content blocks, while private analysis showed encrypted images and HLS key unwrap required app-specific materials.
- Evidence path: `<PROJECT_ROOT>/docs/development/*-contract.md`, `<PROJECT_ROOT>/<sdk>/README.md`, `<ANALYSIS_ROOT>/api/media-download.md`.

#### Generalized Lesson

Use a three-layer media boundary:

- Metadata layer: public SDK core may preserve safe fields such as cover id, cover URL placeholder, avatar field, duration, resolution, content block JSON, or neutral descriptors.
- Adapter/media module layer: private code resolves signed URLs, auth query values, wrapped keys, session headers, decryptors, and expiring materials.
- Output layer: only after bytes are decrypted should reusable standard output transforms or cross-SDK storage formats run.

Do not treat "URL field exists" as "SDK supports download". Do not expose raw signed URLs or keys in public README examples, fixtures, tests, logs, or cache keys.

#### Agent Action

When adding SDK docs or models from analysis that includes media fields, state whether each field is metadata-only, fetchable through a private adapter, or backed by a committed fixture. If decrypt/download is not implemented, mark the media module optional/experimental and avoid adding public facade methods that imply live media readiness.

#### Goal / Action / Validation

- Goal: allow parser/domain SDK progress without leaking private media materials or overclaiming media support.
- Action: document metadata preservation separately from private media decrypt/download ownership.
- Validation or reference source: review API/interface contract, README, BDD, and live-test gate for explicit "metadata-only vs private adapter" wording.

#### Applies When

- Media fields are present in API JSON but raw URLs, keys, tokens, or decrypt steps are sensitive.
- A public SDK core is being built from private app/API analysis.
- Decrypted media output may later be converted into a shared offline format.

#### Does Not Apply When

- The media endpoint is public, stable, unauthenticated, and already covered by committed non-sensitive contract tests.
- The task is only to document private analysis notes, with no downstream SDK or app contract.
- The SDK is intentionally private and allowed to contain the full authorized decrypt/download implementation.

#### Validation

Check the downstream SDK docs answer:

- Which media fields are preserved as metadata?
- Which component owns signed URL freshness, key unwrap, decrypt, and download?
- Are raw URLs, auth query values, tokens, and keys excluded from logs, fixtures, and public examples?
- Is media download support described as implemented, experimental, or out of scope?

#### Promotion Target

- `WORKFLOW.md`
- `implementation/`
- `checklists/`

#### Required Linked Updates

- Update `feedback_history/README.md` and `feedback_history/common/README.md` indexes when adding this lesson.
- If promoted, add a media boundary checklist to the SDK/client development flow.
- Keep project-specific media URLs, keys, and captures in project docs or private evidence stores only.
