> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-05 - In-app route map external transitions

Status: promoted

#### One-line Summary

UI route maps should only expand in-app pages; external transitions should be documented as boundaries, not treated as app screens.

#### Human Explanation

Operation maps are useful because they make app pages repeatable for API capture. If a tap opens system settings, a browser, a payment sheet, share sheet, third-party app, or external intent, the flow has crossed the app boundary. Continuing to map or automate that external surface as if it were part of the target app can create unsafe scripts, noisy evidence, and confusing UI-to-API attribution.

The safer pattern is to stop the app screen map at the boundary. Record the trigger, destination type, foreground package or visible evidence when available, whether manual handling is required, and the API capture window around the transition.

#### Trigger

- A documented tap opens a screen outside the app package.
- A route invokes system settings, browser/custom tab, payment/share sheet, third-party login, or another installed app.
- Automation scripts are being generated from a UI architecture map.

#### Evidence

- Tool: screenshot, `uiautomator dump`, foreground package/focus checks, pcap/MITM/Frida capture windows.
- Sanitized excerpt: `route_id=open-external-help destination_scope=external target=browser/custom-tab`.
- Evidence path: project-private UI evidence and capture logs; reusable skill stores only sanitized workflow guidance.

#### Generalized Lesson

Keep screen inventory scoped to the analyzed app. For route recipes that cross an app boundary, mark `Destination scope: external`, describe the external target type, and stop extending the automatic app page map. Continue API capture only if it is authorized and useful, and clearly separate app-originated requests from external app or system behavior.

#### Agent Action

- Add `Destination scope` and `External transition` fields to route recipes.
- Verify whether the destination is still in-app before adding it to screen inventory.
- Stop automation path expansion when the foreground package or visible UI leaves the app.
- Record trigger point, external target type, manual handling requirement, and API capture window.
- Do not classify system, browser, payment, share, or third-party screens as target-app screens.

#### Applies When

- Device/emulator control is authorized.
- UI maps or automation scripts are being built.
- The app can launch external intents or third-party flows.

#### Does Not Apply When

- The external surface is explicitly in scope and has separate authorization.
- The route never leaves the app package or embedded in-app WebView.
- The task is static-only and no UI route map is being created.

#### Validation

- In-app screens and external transitions are separated in documentation.
- Each external transition includes trigger, destination type, and boundary note.
- Automation scripts stop or require manual confirmation after leaving the app.
- API attribution notes distinguish target-app traffic from external/system traffic.

#### Promotion Target

- `SKILL.md`
- `WORKFLOW.md`
- `DOCUMENTATION.md`
- `TOOLS.md`
- `techniques/http-api/README.md`

#### Required Linked Updates

- Updated quick start and UI map workflow.
- Updated route recipe template with destination scope and external transition fields.
- Updated operation script safety guidance.
- Updated HTTP API UI automation flow.
- Updated feedback indexes and related lessons.
