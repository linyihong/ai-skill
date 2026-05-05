> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-05 - Scrollable and clickable screen mapping

Status: promoted

#### One-line Summary

UI-to-API mapping should classify each screen's scrollable regions and clickable entry points before writing automation scripts.

#### Human Explanation

Screens are not just static screenshots. Some pages reveal API behavior only after swiping lists, paginated feeds, carousels, or related-content sections. Other pages trigger APIs through taps on cards, buttons, tabs, search fields, filters, players, or detail links. If the UI map only records a single screenshot, API attribution can miss pagination/preload APIs or confuse tap-triggered APIs with background traffic.

The reliable approach is to record a lightweight interaction inventory: which screen is scrollable, how far to sample, which visible entries are clickable, and what result is expected. Automation scripts can then replay bounded `swipe` and verified `tap` steps without crawling the whole app.

#### Trigger

- A UI screen contains feeds, lists, carousels, tabs, related items, filters, or detail cards.
- API capture needs to distinguish initial page load, pagination/preload, and tap-triggered detail/action requests.
- Automation scripts are being created for UI-to-API capture.

#### Evidence

- Tool: `adb shell input swipe`, `adb shell input tap`, `uiautomator dump`, screenshots, pcap/MITM/Frida windows.
- Sanitized excerpt: `operation=scroll-feed type=swipe depth=mid` followed by `GET /<path>`; `operation=open-detail type=tap target=item-card` followed by `POST /<path>`.
- Evidence path: project-private UI evidence and API capture logs; reusable skill stores only sanitized method guidance.

#### Generalized Lesson

When building a UI architecture map for API analysis, record scrollable regions and clickable entry points. Use bounded scroll depths such as top/mid/bottom and verified tap targets from screenshots or UI hierarchy. Link these interactions to API capture windows through operation IDs and timestamps.

#### Agent Action

- Add `Scrollable` and `Clickable entries` columns to screen inventory.
- Add an `Interaction Inventory` with `tap` / `swipe` action types.
- For scrollable screens, limit sampling and avoid infinite list crawling.
- For clickable entries, record the selector, label, content-desc, bounds, or coordinate source.
- Include tap/swipe steps in operation scripts and operation-to-API matrices.

#### Applies When

- Device/emulator control is authorized.
- UI-to-API attribution is needed.
- The app has scrollable or tap-heavy screens.

#### Does Not Apply When

- Static analysis only; no UI operation is allowed.
- The screen contains high-risk actions that cannot be safely sandboxed.
- The goal is only network stack triage and UI attribution is not needed yet.

#### Validation

- The UI map shows whether each key screen is scrollable.
- Key clickable entries have documented targets and expected results.
- Operation scripts identify `tap` / `swipe` steps and bounded capture windows.
- Captured APIs can be attributed to initial load, scroll/pagination, or tap-triggered action with confidence.

#### Promotion Target

- `WORKFLOW.md`
- `DOCUMENTATION.md`
- `TOOLS.md`
- `techniques/http-api/README.md`

#### Required Linked Updates

- Updated `WORKFLOW.md` UI map guidance.
- Updated `DOCUMENTATION.md` templates with scrollable/clickable fields and interaction inventory.
- Updated `TOOLS.md` operation script template with swipe support and safety limits.
- Updated `techniques/http-api/README.md` automation flow.
- Updated `feedback_history/http-api/README.md` and root feedback index.
