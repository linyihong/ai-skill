> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-07 - UI Fast Path And Bounded Scroll

Status: candidate

#### One-line Summary

When validating list-to-detail subflows, first test app-owned shortcut controls and bounded scroll convergence before spending long windows on repeated generic swipes.

#### Human Explanation

Feature screens often expose high-semantic shortcuts near a list item, such as a comment, media, detail, or action icon. Tapping the item body may enter the detail screen at the top, while tapping the shortcut may jump directly to the relevant sub-section or trigger the target API faster.

For long scroll probes, accessibility hierarchy may only expose `scrollable=true` and not the runtime scroll offset or maximum extent. In that case, an operation script should avoid open-ended repeated swipes. Prefer bounded scroll counts, large-step swipes, and convergence checks such as unchanged sanitized hierarchy/screenshot hash, no new target API request for a configured post-wait, or a stable feature anchor.

#### Trigger

- A feature API is only triggered after repeatedly scrolling a detail/list page.
- A list card exposes visible shortcut icons or action areas near the target subflow.
- UI hierarchy confirms a scrollable container but does not expose max scroll extent or current offset.

#### Evidence

- Tool: UIAutomator hierarchy, screenshots, package/feature-context guards, and low-overhead request/decode hooks.
- Sanitized excerpt: The reusable pattern is to compare a generic card tap plus repeated scroll against a shortcut-area tap and bounded scroll/no-new-request observation. Project-specific coordinates, labels, and run results stay in project docs.
- Evidence path: `<PROJECT_ROOT>/docs/ui-operation-map.md` or equivalent operation map, plus sanitized capture logs under `<PROJECT_ROOT>/capture/`.

#### Generalized Lesson

Before concluding a pagination or comments/media subflow needs many scrolls, identify whether there is a faster app-owned UI path:

- Try the visible subflow shortcut area on the card or row.
- Parameterize tap coordinates so multiple visible items can be sampled without editing the script.
- Use bounded scroll convergence instead of unbounded loops.
- Treat no-new-request after bounded scrolling as UI-level evidence only, unless a server response explicitly returns an empty page or next/end flag.

#### Agent Action

Next time a feature subflow appears slow to reach:

1. Inspect the screenshot/hierarchy for visible shortcut controls or action areas.
2. Add replay-script parameters for shortcut/card coordinates instead of hard-coding one item.
3. Capture with target package and feature-context guards.
4. Stop scrolling by a configured bound or convergence signal, not by guesswork.
5. Document whether the shortcut path actually triggers the same API, and keep server terminal claims separate from UI no-new-request observations.

#### Goal / Action / Validation

- Goal: Reduce capture time and improve UI/API attribution for subflows hidden below a list or detail screen.
- Action: Compare shortcut tap versus generic item tap, then use bounded scroll/no-change/no-new-request checks.
- Validation or reference source: Target package/feature-context guard remains valid; request hook sees the target API or explicitly does not; sanitized evidence records tap point, scroll bound, and post-wait result.

#### Applies When

- The APK analysis is authorized and the operation is read-only or otherwise in scope.
- The target subflow has a visible shortcut/action area or a scrollable content area.
- UIAutomator does not expose reliable scroll extent metadata.

#### Does Not Apply When

- The shortcut triggers a write action, payment, messaging, posting, liking, deleting, or another high-risk operation outside scope.
- The UI leaves the target app or lands on the wrong in-app screen.
- Server-side pagination truth is required; UI no-new-request alone is not enough.

#### Validation

- Run one control path through the existing generic item/detail route.
- Run one shortcut path using parameterized tap coordinates.
- Compare target API sequence, response shape, foreground package, feature context, and elapsed operation count.
- If using bounded scroll, record the scroll count, post-wait, and whether sanitized UI/API evidence stopped changing.

#### Promotion Target

- `WORKFLOW.md`
- `TOOLS.md`

#### Required Linked Updates

- Updated `feedback_history/README.md` and `feedback_history/common/README.md` so the lesson is discoverable.
- Promotion to `WORKFLOW.md` / `TOOLS.md` is deferred until the shortcut path and bounded-scroll stop rule are validated across at least one replayed flow.
- Checked reusable-guidance boundary: this lesson contains only generalized method, validation, and applicability; target-specific evidence remains in project docs.
