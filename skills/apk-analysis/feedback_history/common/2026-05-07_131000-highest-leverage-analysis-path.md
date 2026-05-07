> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-07 - Highest Leverage Analysis Path

Status: candidate

#### One-line Summary

At each analysis checkpoint, choose the path with the best time-to-evidence ratio instead of continuing a fixed technique sequence.

#### Human Explanation

APK analysis often starts with a default workflow, but the default workflow is a routing aid, not a rigid order. Once evidence shows what the app already exposes, the agent should compare available routes by expected time saved, evidence strength, safety, and reversibility.

For example, long UI scrolling may be useful for attribution but poor for proving pagination truth. If an app-owned API boundary can safely control a page parameter while preserving session, signing, and decrypt behavior, that path may answer the core question faster. Likewise, if a high-semantic request/decode hook is available, it should usually beat broad low-level socket logging.

#### Trigger

- The current route is slow, repetitive, noisy, or only indirectly answers the question.
- Multiple supported techniques are available, such as UI automation, app-owned adapter calls, direct API replay, AOT hooks, pcap, MITM, or static xrefs.
- The user asks whether there is a faster or more advantageous way to analyze the target.
- Evidence already narrows the unknown to a specific contract, parameter, response shape, or boundary.

#### Evidence

- Tool: workflow notes, project API docs, operation map, hook logs, replay harness output, and blocker list.
- Sanitized excerpt: Instead of repeatedly applying one technique, rank candidate paths such as "UI-only scroll", "app-owned parameter override", "direct read-only replay", "static xref first", or "low-level packet capture" by expected time, evidence strength, safety, and setup cost.
- Evidence path: project-specific values and results stay in `<PROJECT_ROOT>` docs/captures; reusable skill docs keep only the decision method.

#### Generalized Lesson

Before investing in another long capture window, run a short strategy check:

- What exact question remains?
- Which supported route can answer it most directly?
- Which route preserves the app-owned boundary for session, signing, decrypt, and UI attribution?
- Which route is read-only, reversible, and easiest to validate?
- What evidence would stop this branch and move to the next route?

Prefer the route with the highest useful evidence per unit time. Do not keep using UI gestures, MITM, broad hooks, static reading, or any other single method just because it was the initial plan.

#### Agent Action

When an analysis task starts or gets stuck:

1. List the active unknowns in one sentence each.
2. List 2-4 available routes supported by current evidence.
3. Choose the route with the best time-to-evidence ratio and explain the validation signal.
4. Keep a fallback route ready if the chosen path fails or produces ambiguous evidence.
5. Document why a slower route is deferred, especially when it is still needed later for UI attribution or edge-case confirmation.

#### Goal / Action / Validation

- Goal: Reduce analysis time while improving evidence quality and avoiding rigid technique sequences.
- Action: Select analysis routes by time-to-evidence, semantic proximity, safety, and validation clarity.
- Validation or reference source: The chosen route answers the current unknown with sanitized evidence, preserves necessary app-owned boundaries or explicitly documents blockers, and cross-checks with UI/API attribution when required.

#### Applies When

- The analysis is authorized and multiple safe techniques are available.
- The current question is specific enough to compare routes.
- A faster route can be validated without collecting secrets or mutating production state.

#### Does Not Apply When

- Authorization or scope is unclear.
- The faster route would require guessing secrets, bypassing safety boundaries, or performing write actions outside scope.
- UI behavior itself is the question; in that case UI capture remains the primary evidence, though it can still be optimized.

#### Validation

- Record the chosen route, fallback route, and validation signal in project notes or the final report.
- Compare at least one accelerated result with UI-triggered or app-owned evidence when attribution matters.
- If the route fails, document the blocker and switch rather than repeating the same slow method.

#### Promotion Target

- `WORKFLOW.md`
- `SKILL.md`
- `TOOLS.md`

#### Required Linked Updates

- Update `WORKFLOW.md` so route selection explicitly includes time-to-evidence and evidence-strength comparison.
- Update `feedback_history/README.md` and `feedback_history/common/README.md` so the lesson is discoverable.
- Checked reusable-guidance boundary: this lesson contains only generalized decision criteria; target-specific probe details remain in project docs.
