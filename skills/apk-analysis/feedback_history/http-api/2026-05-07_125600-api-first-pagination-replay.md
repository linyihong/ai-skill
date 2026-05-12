> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/http-api-documentation-flow.md`](../../../../analysis/apk/workflows/http-api-documentation-flow.md)

### 2026-05-07 - API-First Pagination Replay

Status: candidate

#### One-line Summary

When API envelope, session, and signing prerequisites are known enough, validate pagination with direct read-only API replay before spending long windows on UI scrolling.

#### Human Explanation

UI scrolling is useful for attribution, but it is a slow way to prove pagination ground truth. Once the analysis has recovered enough runtime factors to replay a read-only list endpoint, direct API calls can test page sequences, empty-page behavior, explicit next flags, and count fields much faster than repeated device gestures.

This does not replace UI evidence. UI still confirms which screen and operation trigger the API, verifies app context, and catches hidden client-side behavior such as preloading, filtering, cache hydration, or a different trigger path. The safer sequence is often: use UI + hooks to recover the request contract, replay the API to stress pagination, then return to UI for final attribution and edge-case confirmation.

#### Trigger

- A scrollable list or comments/media subflow requires many slow UI swipes.
- Request fields, route/service mapping, session/bootstrap dependency, opaque parameters, and signing/decode boundaries are already documented enough for a controlled read-only replay.
- The remaining question is pagination truth: page numbers, cursors, empty page, `hasNext`-like flags, count fields, or server end behavior.

#### Evidence

- Tool: sanitized request metadata from hooks, project API docs, fixture or replay script, and a UI operation map.
- Sanitized excerpt: Instead of continuing unbounded UI scrolls, run controlled page requests such as page 1, later pages, and beyond-known-end candidates, then compare response shape and item identity against UI-observed samples.
- Evidence path: project-specific request values, hosts, opaque parameters, and replay outputs stay under `<PROJECT_ROOT>` docs/captures; reusable skill docs keep only the method.

#### Generalized Lesson

Use direct API replay as a pagination accelerator only after the live-readiness factors are known or explicitly supplied:

- Endpoint/path family and route/service mapping are known.
- Session/bootstrap dependency is satisfied.
- Required opaque query/header values are available from a safe provider or documented fixture.
- Signing/header generation and response decode/unwrap are replayable or delegated to an app-owned/private adapter.
- Replay is read-only, rate-limited, and uses sanitized logging.

Then use API replay to test pagination quickly, and use UI only to confirm trigger mapping and app behavior.

#### Agent Action

Next time pagination depends on slow UI scrolling:

1. Check the project domain/runtime baseline for replay readiness.
2. If blockers remain, document them instead of fabricating a client.
3. If ready, build a small read-only replay harness or private adapter call that logs only page key, item count, wrapper keys, response hash, and stable item ids.
4. Probe page 1, page 2, later pages, and one beyond-known-end candidate.
5. Cross-check at least one replayed page against a UI-triggered sample before upgrading the conclusion.

#### Goal / Action / Validation

- Goal: Prove pagination behavior faster than repeated UI gestures while preserving UI-to-API attribution.
- Action: Use direct read-only API replay for page sequence and terminal behavior; use UI replay for trigger confirmation.
- Validation or reference source: Request keys match UI-captured API docs; response shapes match decrypted UI samples; replay logs are sanitized; UI operation map links the API back to a real screen.

#### Applies When

- The analysis is authorized and the endpoint is read-only.
- Required session, opaque parameters, signing, and decode dependencies are known or safely delegated.
- The task is to validate pagination, count semantics, or terminal-page behavior.

#### Does Not Apply When

- The endpoint mutates server state, sends messages, writes reactions, purchases, deletes, follows, posts, or performs another high-risk action.
- Required session/signing/opaque values are unknown and would require guessing secrets.
- UI behavior itself is the unknown, such as whether a control exists, whether a route is reachable, or whether a scroll triggers preloading/cache behavior.

#### Validation

- Compare direct replay page 1 with a UI-triggered page 1 sample by schema, item count, stable ids, and response hash when safe.
- Confirm later-page behavior by page/cursor sequence rather than only by UI no-new-request observations.
- Mark UI-only no-new-request as a weaker observation unless direct replay confirms an empty page or explicit terminal flag.

#### Promotion Target

- `WORKFLOW.md`
- `TOOLS.md`
- `DOCUMENTATION.md`

#### Required Linked Updates

- Updated `feedback_history/README.md` and `feedback_history/http-api/README.md` so the lesson is discoverable.
- Promotion to workflow/tooling is deferred until at least one project validates the API-first replay path against a UI-triggered sample.
- Checked reusable-guidance boundary: no target APK name, host, endpoint, request value, opaque value, or live result is included here.
