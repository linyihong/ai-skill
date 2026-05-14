> 遵守 [共用規則索引](../../../../enforcement/README.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/apk-analysis/artifact-gates.md`](../../../../workflow/apk-analysis/artifact-gates.md)

### 2026-05-05 - Feature reconstruction handoff

Status: promoted

#### One-line Summary

APK analysis documents should preserve enough feature, behavior, domain, API, state, error, fixture, and unknown details for `app-development-guidance` to rebuild the functionality.

#### Human Explanation

Endpoint lists are useful for traffic analysis, but they are not enough to recreate an app feature. To rebuild a function, a later agent needs to know what user capability the endpoint supports, how the user reaches it, what domain concepts appear, how state changes, which errors exist, which data is stored or refreshed, and what evidence proves each claim.

The analysis document should therefore include a handoff section that bridges observed APK behavior to development contracts. It should not overclaim. Unclear domain names, hidden business rules, or untested states must be marked as candidates or open questions.

#### Trigger

- The user wants APK findings to become a buildable app feature, not just an API reference.
- A project analysis document has endpoint names but no feature/capability mapping.
- UI actions, response schemas, state transitions, or errors are observed but not connected.
- A later `app-development-guidance` pass is expected to draft BDD, Domain Model Contract, API / Interface Contract, Error Handling Contract, implementation slices, or tests.

#### Evidence

- Tool: UI map, operation-to-API matrix, hook/MITM/pcap logs, schema-only response summaries, fixtures, replay or contract tests.
- Sanitized excerpt: `operation -> capability -> domain concept candidates -> API/interface contract -> state/error behavior -> fixture -> open questions`.
- Evidence path: project analysis docs; reusable skill stores only the generalized handoff structure.

#### Generalized Lesson

For every high-value feature or API cluster, add a Feature Reconstruction Handoff with:

- Capability and user goal.
- Screen, route, and operation ids.
- Behavior scenario candidates.
- Domain concept candidates with confidence.
- API or public interface contract summary.
- State, empty, loading, and error behavior.
- Data lifecycle: source, local storage, refresh, expiry, sensitivity.
- Fixtures, replay, contract tests, screenshots, UI hierarchy, or sanitized log evidence.
- Open questions and assumptions.

#### Agent Action

- Do not stop at method/path and schema when the user wants functionality rebuilt.
- Add capability and operation mapping to each important API.
- Convert observations into BDD/domain/API/error-contract inputs, but mark low-confidence inferences clearly.
- Keep target-specific hosts, tokens, raw responses, personal data, and private business conclusions out of reusable skill docs.
- Hand sanitized development implications to `app-development-guidance`; keep APK-analysis mechanics here.

#### Applies When

- APK analysis is meant to support feature recreation, SDK/client implementation, API mocks, fixtures, or contract tests.
- UI-to-API binding, response fields, state behavior, or error codes have been observed.
- A later development-guidance pass will convert findings into implementation documents.

#### Does Not Apply When

- The task is only to prove a traffic path, proxy route, or decryption point.
- The feature is out of authorization scope.
- The only available evidence is too thin to infer behavior; record unknowns instead.

#### Validation

- A reader can draft BDD scenarios from the handoff without guessing the main user behavior.
- A reader can draft Domain Model Contract candidates and identify low-confidence terms.
- A reader can draft API / Interface Contract and Error Handling Contract sections from the documented evidence.
- Fixtures or sanitized evidence exist for the high-value API responses.
- Open questions are explicit instead of hidden in prose.

#### Promotion Target

- `DOCUMENTATION.md`
- `WORKFLOW.md`
- `SKILL.md`
- `README.md`
- `techniques/http-api/README.md`
- `feedback_history/README.md`
- `feedback_history/http-api/README.md`

#### Required Linked Updates

- Updated `DOCUMENTATION.md` with the Feature Reconstruction Handoff section and template fields.
- Updated `WORKFLOW.md` analysis completion criteria.
- Updated `SKILL.md` durable-assets and output guidance.
- Updated `README.md` goals, usage, principles, and minimum output.
- Updated `techniques/http-api/README.md` with functional contract mapping requirements.
- Updated `feedback_history/README.md` and `feedback_history/http-api/README.md` indexes.
