> 遵守 [共用規則索引](../../../../enforcement/README.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/apk-analysis/execution-flow.md`](../../../../workflow/apk-analysis/execution-flow.md)

### 2026-05-06 - Auto app-development-guidance for SDK/tool outputs

Status: promoted

#### One-line Summary

When APK analysis documents will be used to build app tools, SDKs, clients, mocks, contract tests, or rebuilt features, automatically apply `app-development-guidance`.

#### Human Explanation

`apk-analysis` can recover evidence, UI/API attribution, schemas, fixtures, and confidence labels. But once the requested output becomes a tool, SDK, client, mock, fixture-driven implementation, contract test, or rebuilt app feature, the work has crossed into development planning.

At that point the agent must load `app-development-guidance` so BDD, Domain Model Contract, API / Interface Contract, Error Handling Contract, implementation slices, tests, and blocker questions are handled by the development skill instead of being improvised inside APK analysis notes.

#### Trigger

- The user wants to build an app-related tool or SDK from APK analysis documents.
- The user asks for a client, mock API, fixture-driven implementation, contract test, or rebuilt feature based on analysis findings.
- A Feature Reconstruction Handoff is being used as input to implementation work.

#### Evidence

- Tool: APK analysis docs, API/schema docs, UI operation map, fixtures, replay/contract tests.
- Sanitized excerpt: `analysis docs -> Feature Reconstruction Handoff -> app-development-guidance -> BDD/contracts/tests`.
- Evidence path: project docs and sanitized fixtures; reusable skill stores only the generalized handoff rule.

#### Generalized Lesson

Use `apk-analysis` for evidence recovery and documentation. Automatically switch to `app-development-guidance` when the analysis output is meant to become:

- App-related tool.
- SDK or client.
- Mock API.
- Fixture-driven implementation.
- Contract test.
- Rebuilt feature or app behavior.

#### Agent Action

- Do not draft implementation plans only from `apk-analysis`.
- Read and apply `app-development-guidance/SKILL.md` before implementation planning.
- Pass a sanitized Feature Reconstruction Handoff as the input artifact.
- Let `app-development-guidance` surface missing BDD, contract, error, storage, security, ownership, and test blocker questions.
- Keep raw APK evidence and target-specific secrets in project analysis docs.

#### Applies When

- Analysis documents are used as implementation input.
- The requested output is development-facing, not just analysis-facing.
- The user mentions tools, SDKs, clients, mocks, fixtures, tests, or rebuilding.

#### Does Not Apply When

- The task is only to locate traffic, decode a response, or document an API.
- The output remains an analysis report with no implementation request.

#### Validation

- `app-development-guidance` was read/applied in the same task.
- The handoff includes behavior, domain/API contract candidates, state/error behavior, fixtures, and open questions.
- Development output includes blocker questions instead of invented missing requirements.

#### Promotion Target

- `SKILL.md`
- `DOCUMENTATION.md`
- `WORKFLOW.md`
- `README.md`
- `../app-development-guidance/SKILL.md`
- `feedback_history/README.md`
- `feedback_history/http-api/README.md`

#### Required Linked Updates

- Updated `SKILL.md` automatic cross-skill trigger.
- Updated `DOCUMENTATION.md` handoff rules and backfill rules.
- Updated `WORKFLOW.md` analysis completion criteria.
- Updated `README.md` usage and minimum output notes.
- Updated `app-development-guidance/SKILL.md` receiver-side trigger.
- Updated `feedback_history/README.md` and `feedback_history/http-api/README.md`.
