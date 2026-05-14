> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/apk-analysis/execution-flow.md`](../../../../workflow/apk-analysis/execution-flow.md)

### 2026-05-06 - Post-reset Window Split

Status: promoted

#### One-line Summary

`clear app data` 後若同時要抓 session recovery 與 feature API，先拆成 startup/session window 與 feature-checkpoint attach window，避免長 Frida-from-launch 窗口污染 UI attribution。

#### Human Explanation

重置資料後，App 常會先經過 splash、公告、更新檢查、guest/session recovery、permission 或外部 provider 轉場。若把這些狀態和目標 feature 操作塞進同一個 Frida-from-launch 長窗口，容易把啟動期 request 誤歸因到 feature 點擊，或讓 hook/timing 造成外部 App、WebView、錯頁、spinner。更穩定的做法是先獨立確認 reset 後的 session recovery，再用 replay runner 把 UI 帶到目標 feature checkpoint，最後 attach 低負載 feature hooks 捕獲同窗 API。

#### Trigger

- 需要驗證 `pm clear` / reinstall 後能否從 App 起始狀態進入某個 feature。
- Frida-from-launch 長窗口在進入 feature 前後發生外部 App、錯頁、公告/更新停留、spinner 或 timing drift。
- 需要把 startup / session recovery API 與 feature-triggered API 分開歸因。

#### Evidence

- Tool: UIAutomator checkpoint replay + Frida Dart AOT hook.
- Sanitized excerpt: post-reset startup window confirmed guest/session recovery; a later feature-window attach confirmed feature list/detail/comment APIs. A direct Frida-from-launch feature attempt hit an external foreground package and was excluded by package guard.
- Evidence path: project-local capture logs/screenshots only; reusable lesson does not include target host, raw endpoint, token, media URL, account, or content values.

#### Generalized Lesson

Reset-to-feature analysis should be split into at least two capture windows when session recovery is involved:

1. Reset + startup/session recovery: prove app data state, launch path, session/login request shape, response schema, and UI/package state.
2. Feature checkpoint + attach: use no-Frida or low-impact replay to reach the feature, validate target package and feature context, then attach feature-specific low-overhead hooks for API attribution.

Only merge these windows into one conclusion after evidence shows the session recovery completed and the feature window stayed in the target app/context.

#### Agent Action

When a user asks for "from reset/startup to feature" API coverage:

- Do record the reset level and classify startup/session APIs separately.
- Do stop at stable checkpoints (`launch`, after notice/login, target tab/page) and capture screenshot/XML.
- Do attach feature hooks after a validated feature checkpoint if launch-time hooks destabilize UI.
- Do exclude windows where foreground package or feature context is wrong.
- Do not claim a feature API was caused by a tap just because the target process still had Frida events.

#### Goal / Action / Validation

- Goal: prevent post-reset startup/session events and feature-triggered APIs from being conflated.
- Action: split capture windows and use checkpoint/package/context validation before attaching feature hooks.
- Validation or reference source: accepted evidence requires reset-state record, same-window startup/session request/decrypt schema, later target feature checkpoint validation, and feature-specific request/decrypt hooks.

#### Applies When

- App data/cache is cleared, app is reinstalled, or first-run/session recovery is under test.
- The target feature requires an existing session/guest profile before it can load.
- Startup Frida hooks are heavier than feature-specific hooks or affect UI timing.

#### Does Not Apply When

- The task only needs preserved-session feature behavior and reset/session recovery is out of scope.
- A single Frida-from-launch capture remains stable and already has package/context-validated feature evidence.
- The app requires an interactive credential login that is not authorized or not available.

#### Validation

Confirm:

- Startup/session window has request keys or equivalent high-level request shape, response schema, and UI/package evidence.
- Feature window starts from a validated feature checkpoint and stays in target package/context.
- Feature API evidence uses feature-specific hooks and decrypted schema.
- Invalid external or wrong-screen windows are marked excluded, not silently reused.

#### Promotion Target

- `WORKFLOW.md`

#### Required Linked Updates

- Updated `WORKFLOW.md` reset baseline section with the post-reset window split rule.
- Updated `feedback_history/README.md` and `feedback_history/common/README.md` indexes.
