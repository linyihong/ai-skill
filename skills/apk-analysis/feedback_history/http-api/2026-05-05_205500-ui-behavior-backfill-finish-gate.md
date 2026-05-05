> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-05 - UI Behavior backfill as finish gate

Status: promoted

#### One-line Summary

APK analysis is not complete until observed UI behavior is backfilled into the project UI behavior/page-map docs.

#### Human Explanation

API findings, screenshots, and test tools can drift apart if UI behavior is only mentioned in chat or hidden inside endpoint docs. A user looking for "UI Behavior" needs a stable project entry point that explains where page maps live, what actions were observed, and which UI labels map to which API parameters. If no UI evidence was captured, the documentation should say so explicitly instead of silently omitting the UI behavior.

#### Trigger

- An analysis includes App operation, screenshots, UI hierarchy, visible sort labels, tap/swipe/input behavior, or page-to-API mapping.
- A tool or fixture is changed to mimic App operation behavior.
- A user asks where UI behavior documentation is or says the tool should behave like the App.

#### Evidence

- Tool: project UI behavior docs and live API replay.
- Sanitized excerpt: a home page sort row was observed in App screenshots, then mapped to `moduleSort` values and a top-level `chosenVideoInfo[]` list without storing secrets.
- Evidence path: project-level UI behavior/page-map docs and sanitized screenshot references.

#### Generalized Lesson

Keep a dedicated project UI Behavior entry/index and page-level UI maps. At the end of an analysis, update them with entry path, visible UI blocks, App labels, operation steps, data source mapping, evidence, and unknowns. If UI capture was skipped or blocked, write `needs capture`, `needs replay`, or `Trigger confidence: low` in the project docs.

#### Agent Action

Before reporting an APK-analysis task complete:

1. Check whether UI behavior was observed or inferred from App operation.
2. Update the project UI Behavior entry/index and relevant page-level map.
3. Cross-link API docs when a UI label, sort option, card, detail action, comment flow, or playback action maps to an endpoint or response field.
4. If no UI capture exists, document the gap explicitly rather than leaving the UI Behavior section absent.

#### Applies When

- Authorized APK analysis where the App can be operated or has existing UI evidence.
- UI-to-API binding, page/tab/module analysis, playback/detail/comment workflows, or local tools that emulate App behavior.

#### Does Not Apply When

- Pure static setup work before any UI behavior or endpoint is known.
- Non-UI protocol-only analysis where the project explicitly scopes UI capture out; even then, record the UI gap if the protocol is tied to a user-facing capability.

#### Validation

- The project has a stable UI Behavior entry or equivalent index.
- The relevant page map includes operation steps, visible App labels, data source mapping, evidence, and unknowns.
- API/tool docs link back to UI behavior where the behavior explains how the endpoint is triggered.

#### Promotion Target

- `SKILL.md`
- `WORKFLOW.md`
- `DOCUMENTATION.md`

#### Required Linked Updates

- Promoted into `SKILL.md`, `WORKFLOW.md`, and `DOCUMENTATION.md`.
- Updated `feedback_history/README.md` and `feedback_history/http-api/README.md`.
