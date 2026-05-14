> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/apk-analysis/execution-flow.md`](../../../../workflow/apk-analysis/execution-flow.md)

### 2026-05-06 - State reset baseline before feature capture

Status: promoted

#### One-line Summary

Feature-level APK API analysis should record whether the run starts from a clean app state, because cache/session leftovers can hide startup and first-use APIs.

#### Human Explanation

Named-feature analysis often starts from an already-open app or a warm authenticated session. That is useful for low-friction UI exploration, but it can miss first-use APIs, bootstrap calls, permission/notice gates, cache hydration, session recovery, and the exact route from launcher to the target feature. If the user asks for the complete flow from app start, the analysis needs an explicit reset baseline: what app state was cleared, what state was intentionally kept, how the app was launched, and which API windows belong to reset/startup/navigation/feature operations.

Resetting state is not always safe or necessary. It can remove a test session, trigger login/rate limits, or require user consent. The lesson is to document and choose the reset level deliberately, not to clear data blindly.

#### Trigger

Use this when:

- the user asks for a complete flow from app start to a named feature;
- a feature may be served from cache after warm navigation;
- startup, notice, onboarding, login/session recovery, or feature bootstrap APIs must be attributed;
- previous captures started from an already-running app or attach-after-load state.

#### Evidence

- Tool: project documentation review after feature-level UI/API attribution.
- Sanitized excerpt:
  - feature read flows were confirmed from UI action windows and Dart decode hooks;
  - the route proof used `force-stop` and clean UI operation captures;
  - a dedicated reset-to-feature baseline was still missing as a reusable skill rule, so future agents could skip cache/session state documentation.
- Evidence path: project-level UI/API docs only; target-specific package names, selectors, service identifiers, tokens, and response values stay in project files, not reusable skill docs.

#### Generalized Lesson

Treat reset baseline as part of end-to-end feature capture. Before calling a feature flow complete, record the starting app state and separate these windows:

1. Reset / state preparation: `force-stop`, optional app data/cache clear, permission state, network/proxy state, and whether session/account state is preserved.
2. Cold start / bootstrap: launch, first screen, notice/onboarding/login/session recovery, startup/background APIs.
3. Navigation to target feature: tab/menu/search/deeplink steps and UI evidence.
4. Feature operations: list, filter/category, pagination/scroll, search, detail, comments/media/actions, and read/write boundaries.
5. Documentation closure: page map, operation map, API list, schema/correlation, feature handoff, gaps.

#### Agent Action

When the user asks for app-start-to-feature analysis:

1. Ask or infer the allowed reset level: `force-stop only`, `clear cache`, `clear app data`, `reinstall`, or `preserve logged-in session`.
2. If clearing data may remove credentials or trigger login/rate limits, warn and prefer a controlled test account/session plan.
3. Create an operation id for the reset baseline, such as `reset-state-to-feature`.
4. Capture or document the full route from launcher/startup to the target page before attributing feature APIs.
5. Mark API events as `startup/preload`, `navigation`, `feature-triggered`, `cache-hydration`, or `background/ambiguous`.
6. Update project docs with the reset baseline and confirmed API chain; do not leave the flow only in chat.

#### Goal / Action / Validation

- Goal: Make app-start-to-feature API analysis reproducible instead of depending on a hidden warm app state.
- Action: Add a reset baseline requirement to the common workflow and document the reset/navigation/feature API windows separately.
- Validation or reference source: Confirm docs show the starting state, operation id, UI path, capture evidence, API confidence labels, and explicit gaps for APIs not seen after reset.

#### Applies When

- The task is authorized APK analysis.
- The user asks for a complete feature flow from app launch or clean state.
- Caches, persisted session, startup bootstrap, or first-use UI gates may change observed APIs.
- The feature has a UI route that can be operated or documented.

#### Does Not Apply When

- The task is only a narrow static code lookup.
- Clearing data is out of scope, unsafe for the test account, or would destroy required evidence; in that case document the warm-state baseline instead.
- The feature is intentionally tested as a warm-session behavior; then record that baseline and avoid claiming clean-start coverage.

#### Validation

The lesson is validated when:

- project docs state whether app data/cache was cleared, force-stopped, reinstalled, or preserved;
- the operation map includes reset/startup/navigation/feature windows;
- startup/preload APIs are not mislabeled as current tap-triggered APIs;
- confirmed feature APIs still have per-API request/response docs and evidence links;
- raw tokens, account values, private IDs, and raw responses are not copied into reusable skill docs.

#### Promotion Target

- `WORKFLOW.md`
- `DOCUMENTATION.md`

#### Required Linked Updates

- Promoted into `WORKFLOW.md` as a reset baseline and end-to-end feature capture rule.
- `DOCUMENTATION.md` already requires UI operation maps, API docs, and feature handoff; no template change needed beyond using the existing operation/state fields for reset baseline.
- Updated `feedback_history/common/README.md`.
