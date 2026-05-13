> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/apk-analysis/artifact-gates.md`](../../../../workflow/apk-analysis/artifact-gates.md)

### 2026-05-06 - API Catalog finish gate

Status: promoted

#### One-line Summary

API list work is not complete when endpoints are only listed; confirmed APIs need a catalog with grouped indexes, per-API detail, coverage gaps, UI/API mapping, SDK/client field usage, evidence, and validation.

#### Human Explanation

APK analysis often discovers many endpoints through hooks, replay, decoded fixtures, pcap timing, or static strings. A flat endpoint table helps at first, but it is not enough for SDKs, mocks, contract tests, or feature reconstruction. The API reference must answer where each API belongs, how it is triggered, what request and response fields mean, how confident the mapping is, and what remains untested.

The reusable pattern is to create an API Catalog: one total API entry, grouped indexes, per-operation detail files, coverage/gap status, UI/API mapping, SDK/client field usage, and sanitized evidence. If a flow is not fully confirmed, create a skeleton and label gaps instead of waiting for perfect coverage.

#### Trigger

- The user asks for an API list, API reference, endpoint inventory, SDK input, mock API, contract test, or rebuildable feature.
- Dynamic or static analysis finds multiple APIs that need long-term documentation.
- A project has API rows but no per-operation details, field meanings, or coverage/gap status.

#### Evidence

- Tool: hook logs, MITM export, pcap, replay output, decrypted fixtures, UI map, operation script, SDK/client tests.
- Sanitized excerpt: `observed API -> group index -> per-API detail -> coverage/gap -> validation`.
- Evidence path: project repository API docs; reusable skill stores only the generalized catalog rule.

#### Generalized Lesson

API Catalog artifacts should include:

| Artifact | Required content |
| --- | --- |
| API entry | Hosts/base URLs, traffic families, wrapper/decode rules, shared headers, coverage, UI map, SDK/client links. |
| Group index | APIs grouped by path prefix, domain, feature, or protocol family. |
| Per-API detail | Request, response, field meaning, behavior, evidence, validation, open questions. |
| Coverage / gap matrix | Static candidates, observed, replayed, decoded, UI-bound, tested, missing, scoped out. |
| SDK/client mapping | Consumed fields, compatibility expectations, raw JSON strategy, fixture/test status. |

#### Agent Action

- Do not stop at a flat endpoint list when the task is API reference, SDK, mock, contract test, or feature reconstruction.
- Create or update the project API Catalog and per-API detail skeletons for high-value APIs.
- Mark incomplete items as `candidate`, `needs capture`, `needs replay`, `meaning unknown`, `low confidence`, `out of scope`, or `not observed`.
- Link API docs back to UI operation ids, capture windows, evidence paths, and SDK/client tests when available.

#### Goal / Action / Validation

- Goal: Make API list output usable for future analysis, SDK/client work, mocks, contract tests, and feature reconstruction.
- Action: Promote API Catalog structure into `DOCUMENTATION.md`, `techniques/http-api/README.md`, `WORKFLOW.md`, `README.md`, and `SKILL.md`.
- Validation or reference source: Compare against project-level API catalog patterns and require lints, Markdown link check, diff review, and necessary tool sync; reference-only tool usage does not require bundle sync.

#### Applies When

- HTTP/HTTPS API metadata is visible, decoded, replayed, or statically enumerated.
- The output will be used by another agent, SDK/client, mock server, fixture, contract test, or feature handoff.

#### Does Not Apply When

- A single low-value endpoint is noted only as a temporary lead and not used as project reference.
- The user asks only for a quick read-only status; still mention that API Catalog completion has not been performed.

#### Validation

- API docs have a total entry, grouped index, per-API detail, coverage/gap matrix, and evidence links.
- High-value APIs include request fields, response fields, field meanings, UI trigger confidence, and validation/open questions.
- `apk-analysis` docs mention the API Catalog finish gate and per-API skeleton rule.

#### Promotion Target

- `DOCUMENTATION.md`
- `techniques/http-api/README.md`
- `WORKFLOW.md`
- `README.md`
- `SKILL.md`
- `feedback_history/http-api/README.md`

#### Required Linked Updates

- Updated `DOCUMENTATION.md` with API Catalog / API list requirements.
- Updated `techniques/http-api/README.md` with catalog shape and finish gate.
- Updated `WORKFLOW.md`, `README.md`, and `SKILL.md` with API Catalog completion requirements.
- Updated `feedback_history/http-api/README.md`.
