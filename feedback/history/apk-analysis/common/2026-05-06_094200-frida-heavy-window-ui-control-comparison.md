> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[neutral-language](../../../../shared-rules/neutral-language.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

# Extracted — See [`workflow/apk-analysis/execution-flow.md`](../../../../workflow/apk-analysis/execution-flow.md) (Section 1: Capture Window 詳細規則)

### 2026-05-06 - Frida-heavy window UI control comparison

Status: validated

#### One-line Summary

When a Flutter app shows retry/spinner states only during heavy Frida hook windows, run a no-Frida UI control and a reduced attach-after-load probe before blaming backend timeout or missing data.

#### Human Explanation

Dynamic instrumentation can change the app's UI timing enough to alter controller state, especially when many high-semantic Dart hooks log request, decrypt, and model parsing events. If the UI displays a retry panel or indefinite spinner while decrypted API responses are still arriving, the failure may be caused by instrumentation overhead or a state-machine interaction rather than by absent backend data. A same-route no-Frida control separates app behavior from observation overhead.

#### Trigger

Use this when an authorized APK analysis session shows:

- a retry button or spinner during a Frida-heavy same-window capture;
- decrypted response JSON arrives in hooks, but the expected UI list/detail does not render;
- target high-level API hooks are missing even though lower-level request/decrypt hooks fire;
- the user suspects a short app-side timeout or slow capture.

#### Evidence

- Tool: adb UIAutomator timing captures, screenshots, Frida Dart AOT hooks.
- Sanitized excerpt:
  - retry tap cleared the error panel;
  - first request began about 2 seconds after tap;
  - decrypted category JSON arrived about 4 seconds after tap;
  - the Frida-heavy window stayed on a spinner and never hit the expected article-list hooks;
  - the same UI path without Frida loaded category, list, detail, and comments normally.
  - a reduced feature-specific hook attached after list render captured the target detail/comment API hooks and schemas without triggering the retry/spinner state.
- Evidence path: `<PROJECT_ROOT>/capture/` UI hierarchy and Frida logs only; no raw titles, user content, tokens, hosts, or private values are written into reusable skill docs.

#### Generalized Lesson

Do not treat a retry/spinner during instrumentation as immediate proof of backend timeout. First compare:

1. Frida-heavy same-window behavior;
2. no-Frida UI-only behavior on the same route;
3. reduced attach-after-load behavior after the UI has already reached the target screen.

If no-Frida succeeds and Frida-heavy gets stuck while decrypt hooks still receive data, classify the problem as instrumentation-sensitive timing or controller-state interaction until a lower-overhead hook set reproduces the target API call.

#### Agent Action

Next time this symptom appears:

1. Capture the retry button selector/bounds and tap timestamp.
2. Record UI snapshots at several delays after retry.
3. Compare request-start and decrypt/handoff timestamps against the tap.
4. Run the same UI path with Frida disabled.
5. If UI-only succeeds, reduce hook volume or attach after the route is loaded before promoting any timeout conclusion.
6. Document the retry behavior as diagnostic evidence unless it also hits the target high-level API function.

#### Applies When

- The app/device is authorized for dynamic analysis.
- The target route can be replayed with and without Frida.
- Hooks are high-volume enough to plausibly affect UI timing.
- Decrypted responses arrive but UI rendering or controller progression stalls.

#### Does Not Apply When

- The no-Frida control also fails in the same way.
- Network capture shows the backend never responds.
- The app explicitly returns a business error JSON that matches the UI error.
- The target route requires instrumentation to become reachable.

#### Validation

The lesson is validated when:

- the no-Frida route renders data that the Frida-heavy run failed to render;
- Frida logs show response/decrypt events during the failed UI window;
- reducing hook volume or attaching later changes the UI behavior or prevents the retry/spinner state.
- the reduced hook can capture the originally missing target high-level API calls without broad model/parser hooks.

#### Promotion Target

- `WORKFLOW.md`
- `TOOLS.md`

#### Required Linked Updates

- Updated `feedback_history/common/README.md`.
- Project-specific retry timing and confidence changes belong in project docs, not this reusable lesson.
