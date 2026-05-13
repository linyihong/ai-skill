> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/apk-analysis/execution-flow.md`](../../../../workflow/apk-analysis/execution-flow.md)

### 2026-05-05 - Screen reachability operation recipes

Status: promoted

#### One-line Summary

UI architecture maps should document how to reach important screens, not only what each screen looks like.

#### Human Explanation

Screenshots and screen inventories explain the visible UI, but later API capture needs repeatable operations. If the documentation only says `item.detail` exists, future agents still need to rediscover how to get there. A route recipe records the entry state and step-by-step launch, tap, and swipe actions needed to reach a screen.

Keep screen identity and route identity separate. The screen id describes the destination. The route id describes how to reach it. Operation scripts and operation-to-API matrices should reference the route id so manual replay, automation, and API attribution use the same path.

#### Trigger

- A page must be revisited later for API capture or validation.
- A key API needs UI attribution.
- Automation scripts are being created from a UI architecture map.
- Multiple screens or tabs can lead to the same endpoint.

#### Evidence

- Tool: screenshots, `uiautomator dump`, `adb shell input tap`, `adb shell input swipe`, pcap/MITM/Frida capture windows.
- Sanitized excerpt: `route_id=open-detail start=home.feed step=1 type=tap target=item-card expected=item.detail`.
- Evidence path: project-private UI evidence and capture logs; reusable skill stores only sanitized workflow guidance.

#### Generalized Lesson

Add a `Screen Reachability / Operation Recipe` section to UI architecture maps. For each important target screen, document start state, ordered steps, action type, target or gesture, selector/coordinate source, expected destination, and notes. API operation matrices should reference the route id rather than duplicating fragile natural-language paths.

#### Agent Action

- Create stable `screen_id`, `route_id`, and `operation_id` values.
- Use `Screen Inventory` for destination metadata.
- Use `Screen Reachability / Operation Recipe` for step-by-step route documentation.
- Use `Operation To API Matrix` to bind route id and operation id to captured APIs.
- Make automation scripts print route id, step number, action type, target, and timestamps when useful.

#### Applies When

- Device/emulator operation is authorized.
- Later replay, API capture, or validation depends on reaching the same screen again.
- UI navigation has tabs, drawers, lists, detail pages, search, filters, or player pages.

#### Does Not Apply When

- The task is static-only and no UI operation is allowed.
- UI attribution is explicitly out of scope.
- The route includes unsafe actions that cannot be sandboxed or manually guarded.

#### Validation

- Important screens have a screen id and at least one route id.
- A human can follow the route recipe without rediscovering the UI.
- A script can be written from the same route recipe without adding hidden steps.
- API capture notes reference route id / operation id / time window consistently.

#### Promotion Target

- `SKILL.md`
- `WORKFLOW.md`
- `DOCUMENTATION.md`
- `TOOLS.md`
- `techniques/http-api/README.md`

#### Required Linked Updates

- Updated `SKILL.md` quick start.
- Updated `WORKFLOW.md` UI map records and principles.
- Updated `DOCUMENTATION.md` with `Screen Reachability / Operation Recipe`.
- Updated `TOOLS.md` operation script logging.
- Updated `techniques/http-api/README.md` automation flow.
- Updated feedback indexes.
