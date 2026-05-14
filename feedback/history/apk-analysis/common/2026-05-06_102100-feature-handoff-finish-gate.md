# Extracted — See [`workflow/apk-analysis/artifact-gates.md`](../../../../workflow/apk-analysis/artifact-gates.md) (Feature Reconstruction Handoff) and [`analysis/apk/workflows/http-api-documentation-flow.md`](../../../../analysis/apk/workflows/http-api-documentation-flow.md)

> 遵守 [共用規則索引](../../../../enforcement/README.md)、[neutral-language](../../../../enforcement/neutral-language.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-06 - Feature handoff finish gate

Status: promoted

#### One-line Summary

When a named APK feature is understood enough to explain its UI/API behavior, create or update a project-level feature handoff document before calling the analysis complete.

#### Human Explanation

API tables, schema catalogs, hook logs, and chat summaries are not enough for feature-level APK analysis. Once the agent understands the main user path and major API flows, that understanding must be turned into a durable feature artifact that connects UI behavior, domain concepts, API contracts, state/error behavior, data lifecycle, evidence, and unknowns. Otherwise future agents can see endpoints but still miss how the feature works.

#### Trigger

Use this when:

- a user asks to analyze a named page, tab, module, or feature;
- core flows move from `Candidate` to `Confirmed`;
- the agent can explain entry path, UI blocks, request keys, response schema, and state/error behavior;
- the user asks whether API docs or feature architecture have been produced.

#### Evidence

- Tool: project documentation review after dynamic UI/API attribution.
- Sanitized excerpt:
  - API/schema/page-map files were updated first;
  - the feature was understood, but the feature-level handoff was created only after the user pointed out the gap;
  - adding the handoff made the analysis more reusable because it connected behavior, domain, API, state, evidence, and unknowns in one place.
- Evidence path: project docs only; target-specific service hashes and schema shapes remain in project files, not reusable skill docs.

#### Generalized Lesson

Treat feature handoff as a finish gate for named-feature analysis. A feature is not fully documented just because endpoint tables are updated. If the analysis has enough evidence to describe the feature, the agent must create or update a project-level handoff/architecture document in the same session.

#### Agent Action

Before reporting a named feature as complete:

1. Check whether a project-level feature handoff / architecture file exists.
2. If missing, create one with feature/capability, behavior scenarios, domain candidates, API contracts, state/error handling, data lifecycle, validation evidence, and unknowns.
3. If evidence is incomplete, write a skeleton and mark gaps as `needs capture`, `candidate`, or `low confidence`.
4. Link the handoff from the page map, operation map, and API summary docs.
5. Only then summarize the analysis result to the user.

#### Applies When

- The task targets a named feature, tab, module, or user-visible page.
- UI/API attribution has enough evidence to explain at least one core flow.
- The project has a docs area where feature artifacts can be stored.

#### Does Not Apply When

- The turn is only a narrow tool/environment fix with no feature understanding.
- The task is explicitly limited to a one-off command or single evidence lookup.
- The feature cannot be operated or identified yet; then document `needs capture` in the closest existing map instead.

#### Validation

The lesson is validated when:

- the feature handoff exists or is updated;
- it links to page map, operation map, API/schema/correlation docs;
- it separates confirmed behavior from assumptions and unknowns;
- no raw secrets, tokens, raw responses, or personal content are copied into reusable docs.

#### Promotion Target

- `SKILL.md`
- `DOCUMENTATION.md`

#### Required Linked Updates

- Promoted into `SKILL.md` durable asset list as a finish gate.
- Promoted into `DOCUMENTATION.md` feature handoff section.
- Updated `feedback_history/common/README.md`.
