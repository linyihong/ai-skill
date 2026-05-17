> Follow [enforcement rules](../../../../enforcement/README.md), [dependency-reading](../../../../enforcement/dependency-reading.md), [neutral-language](../../../../enforcement/neutral-language.md), [goal-action-validation](../../../../enforcement/goal-action-validation.md), and [feedback-lessons](../../../../enforcement/feedback-lessons.md). This lesson records only generalized guidance.
# Extracted — See [`workflow/apk-analysis/artifact-gates.md`](../../../../workflow/apk-analysis/artifact-gates.md)

### 2026-05-07 - UI Route Backfill After Runtime Probe

Status: candidate

#### One-line Summary

When a runtime/API probe discovers a new in-app route, update the UI architecture map in the same closure loop instead of only documenting runtime evidence.

#### Human Explanation

Runtime investigations often require navigating through settings, profile, dialogs, tabs, or global menus to trigger a provider or API. Those routes are UI architecture evidence too. If the agent only updates operation logs, API docs, or runtime baseline files, later work loses the reusable route recipe and screen inventory.

#### Trigger

- A non-UI main task uses manual or scripted app navigation.
- The route reaches a previously undocumented app screen, settings page, tab, dialog, or global menu.
- The new route affects API/session/signing/storage behavior or future replay.

#### Evidence

- Tool: UIAutomator XML/screenshots, Frida timing logs, operation maps, runtime baseline docs.
- Sanitized excerpt: a provider validation navigated through a global settings path to a language screen; runtime docs captured the provider update, but the UI architecture map initially lacked the new settings/language route.
- Evidence path: project-specific screenshots, XML, and logs stay in `<PROJECT_ROOT>/capture/`; UI route summaries belong in `<PROJECT_ROOT>/docs/UI架構地圖/`.

#### Generalized Lesson

Treat newly discovered in-app routes as first-class architecture findings even when the primary goal is API or runtime analysis. Before closing the task, check whether each new route/screen has a page-level or global UI architecture map entry and whether related operation/API/runtime docs link to it.

#### Agent Action

After any dynamic probe, compare the operated UI path against existing `docs/UI架構地圖/` files. If missing, create or update the relevant route/page map with entry path, UI blocks, operations, evidence, runtime/API attribution, and gaps. Then add cross-links from operation map, runtime baseline, API docs, or feature docs touched in the same turn.

#### Goal / Action / Validation

- Goal: Prevent route knowledge from being trapped in logs or chat summaries.
- Action: Backfill UI architecture docs whenever a new in-app route is discovered during runtime/API work.
- Validation or reference source: A valid closure includes a UI map file or explicit route entry, evidence paths, and links from any runtime/API docs that rely on that route.

#### Applies When

- Settings, profile, global menus, dialogs, bottom tabs, feature tabs, or state-changing pages are opened during an APK probe.
- The screen is inside the target app and can be reached again.
- The route matters for future validation, replay, testing, or feature reconstruction.

#### Does Not Apply When

- The navigation leaves the target app; external destinations should be documented as external transitions only.
- The route is already documented and only evidence paths need appending.
- The probe fails before reaching a stable app screen.

#### Validation

Confirm the UI architecture map is updated in the same work closure, and include a short reason if the route is intentionally not mapped.

#### Promotion Target

- `DOCUMENTATION.md`
- `WORKFLOW.md`

#### Required Linked Updates

- Updated `DOCUMENTATION.md` UI architecture map requirements to include this backfill gate.
- Checked `feedback_history/common/README.md` and root `feedback_history/README.md`; both should index this lesson.
- No immediate promotion: keep as `candidate` until repeated across another runtime/API probe.
