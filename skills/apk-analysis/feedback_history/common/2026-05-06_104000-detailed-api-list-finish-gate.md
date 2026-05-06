> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-06 - Detailed API list finish gate

Status: promoted

#### One-line Summary

Confirmed APK API flows need per-API request/response documents, not only schema summaries or endpoint correlation tables.

#### Human Explanation

Schema catalogs and correlation tables are good analysis tools, but they are not enough as API reference. Once an API flow is confirmed, future agents need a durable document that answers: how is this API requested, what fields are required, what wrapper is returned, what inner fields mean, what evidence proves it, and what gaps remain. Without that per-API list, feature docs become difficult to use for SDKs, fixtures, replay, or implementation handoff.

#### Trigger

Use this when:

- a flow is promoted to `Confirmed`;
- the user asks for detailed request / response docs;
- the project has only schema/correlation/feature summaries;
- the analysis output may be used by SDKs, clients, mocks, replay tools, or future agents.

#### Evidence

- Tool: project documentation review after feature-level UI/API attribution.
- Sanitized excerpt:
  - feature summaries had request keys, service hashes, and schema shape;
  - they did not yet provide a per-API reference similar to an API list;
  - adding an API list made each confirmed flow directly searchable by request/response contract and evidence.
- Evidence path: project API docs only; target-specific raw values remain in project capture, not reusable skill docs.

#### Generalized Lesson

Treat detailed API list docs as a finish gate for confirmed API flows. Correlation tables prove the mapping; per-API docs make the mapping usable.

#### Agent Action

Before reporting confirmed API analysis as complete:

1. Check whether the project has an API list/index location.
2. Create or update a group index and one per-API document for each confirmed flow.
3. Include method/path family, endpoint/service identifier, request fields, response wrapper, inner schema, field meanings, evidence, confidence, and gaps.
4. Use placeholders/hashes for sensitive values; do not copy tokens, raw signatures, raw service names, private query values, user content, or media URLs.
5. Link the API list from feature summary, schema catalog, endpoint correlation, UI/operation maps, and feature handoff.

#### Goal / Action / Validation

- Goal: Make confirmed API flows usable as durable request/response references.
- Action: Require per-API documents or skeletons once flow confidence reaches `Confirmed`.
- Validation or reference source: API docs link back to evidence/correlation and mark gaps instead of leaving only schema summaries.

#### Applies When

- The APK uses REST, gateway, GraphQL-like, RPC, local bridge, or service-selector APIs.
- The agent has enough evidence to mark a flow `Confirmed`.
- The project is building durable API reference docs.

#### Does Not Apply When

- The task is only environment setup or a failed capture with no confirmed API.
- The API remains a weak candidate; then create a skeleton only if useful and label it `Candidate` / `needs capture`.

#### Validation

The lesson is validated when:

- each confirmed API has a request/response reference entry;
- field meanings and unknowns are explicit;
- docs link back to evidence and correlation;
- sensitive/raw values are redacted or replaced by placeholders.

#### Promotion Target

- `SKILL.md`
- `DOCUMENTATION.md`

#### Required Linked Updates

- Promoted into `SKILL.md` durable asset list as a finish gate.
- Promoted into `DOCUMENTATION.md` as detailed API request/response document rules.
- Updated `feedback_history/common/README.md`.
