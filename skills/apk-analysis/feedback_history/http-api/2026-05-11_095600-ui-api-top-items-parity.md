> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/http-api-documentation-flow.md`](../../../../analysis/apk/workflows/http-api-documentation-flow.md)

### 2026-05-11 - UI/API top items parity before API blame

Status: candidate

#### One-line Summary

Before blaming an API for a list/feed mismatch, compare the first few visible UI items with the first few decoded API items from the same operation window.

#### Human Explanation

List/feed bugs can come from API data, UI sorting, cache hydration, tab/filter mismatch, stale local state, or SDK field mapping. If analysis only inspects the API response, it may miss that the user-visible UI is showing a different tab, a cached page, or a transformed order. A small top-items parity check anchors the API to the actual screen.

#### Trigger

A user reports wrong list/feed behavior, missing items, incorrect ordering, stale content, media mismatch, or SDK output that seems inconsistent with the app UI.

#### Evidence

- Tool: UI capture plus decoded API response comparison.
- Sanitized excerpt: a feed investigation relied on API/SDK evidence without first proving that the top visible UI items matched the decoded API rows for the same operation.
- Evidence path: keep concrete screenshots, hierarchy dumps, API samples, and item ids in the project docs or private evidence folder.

#### Generalized Lesson

For list/feed analysis, capture a small UI/API parity sample before assigning fault: usually the first three visible UI cards and the first three decoded API rows are enough to validate operation binding, ordering, selected tab/filter, and field mapping. If the top items do not match, investigate UI route, cache, sort, pagination, or mapper before calling the API wrong.

#### Agent Action

When analyzing feed/list errors:

1. Record the UI operation id, selected tab/filter, and capture window.
2. Capture the first visible UI items using sanitized stable fields such as title hash, item id if visible, thumbnail/path family, or row position.
3. Decode the API response from the same operation and compare the first three items by id/title/path/order.
4. Only after parity is established, attribute remaining mismatch to API data or SDK mapping; if parity fails, document the UI/API binding gap.

#### Goal / Action / Validation

- Goal: avoid misattributing UI/list bugs to the API when operation binding is not proven.
- Action: add a top-items UI/API parity check to feed/list investigations.
- Validation or reference source: project evidence shows UI top rows and API top rows match or clearly explains the mismatch class.

#### Applies When

- The target feature is a list/feed/search result/tabbed page.
- The conclusion depends on whether API order/content matches what the app user sees.
- SDK/client behavior is being validated against app UI behavior.

#### Does Not Apply When

- The task is a pure schema inventory with no claim about UI behavior.
- The API call is a non-list action where UI top-item ordering is irrelevant.

#### Validation

Use sanitized screenshots/UI hierarchy plus decoded API samples. Do not copy raw private content, tokens, accounts, or full media URLs into reusable skill docs.

#### Promotion Target

- `WORKFLOW.md`
- `DOCUMENTATION.md`
- `SKILL.md`

#### Required Linked Updates

- Indexed in `feedback_history/http-api/README.md` and `feedback_history/README.md`.
- Concrete app/page evidence stays in project docs, not this lesson.
