> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/http-api-documentation-flow.md`](../../../../analysis/apk/workflows/http-api-documentation-flow.md)

### 2026-05-05 - UI automation operation scripts for API capture

Status: promoted

#### One-line Summary

輕量 UI 架構地圖可以為關鍵 flow 產生小型可重放操作腳本，搭配 pcap/MITM/Frida 時間窗穩定抓 API。

#### Human Explanation

手動點 App 很容易造成 API attribution 不穩：操作時間不一致、背景預載混入、或同一段 capture 內有太多 UI 動作。對少量高價值流程建立 `operation_id` 與最小 adb/uiautomator 腳本，可以讓 API capture 更容易重跑與比較。

自動化不應變成全 App crawler。批量遍歷容易讓裝置卡住，也可能觸發登入限流、付費、刪除、發文、下單、私訊等高風險行為。腳本要保持可中止、限量、只做授權範圍內的一個 flow，並輸出開始/結束時間戳。

#### Trigger

- 需要重複抓同一 UI flow 的 API。
- 已有 API 清單但缺少穩定 UI attribution。
- 手動操作造成 pcap/MITM/Frida log 時間窗不穩。
- 使用者希望「搭配自動化寫腳本抓 API」。

#### Evidence

- Tool: `adb shell input`, `adb shell screencap`, pcap/MITM/Frida log window.
- Sanitized excerpt: `operation=open-detail phase=start/end ts=<utc>; POST /<path>; source=hook; triggerConfidence=high`.
- Evidence path: project-private `evidence/ui/` and API capture logs; reusable skill stores only sanitized template guidance.

#### Generalized Lesson

Use small operation scripts to stabilize UI-to-API capture. Each script should represent one operation, print UTC start/end timestamps, optionally save a sanitized screenshot/UI hierarchy, and be correlated with network/hook logs. Automation is optional and should not replace manual authorization checks or API field documentation.

#### Agent Action

When building a UI map for API analysis:

- Add `operation_id` to each key flow.
- If repeated capture is needed, create a minimal operation script for that flow.
- Run pcap/MITM/Frida capture in the same start/end window.
- Record script path, timestamp window, screenshot path, and API evidence in the operation-to-API matrix.
- Avoid login loops, destructive actions, payments, posting, messaging, account changes, or anything outside authorization.

#### Applies When

- The device/emulator can be controlled.
- A test account or authorized state exists.
- The goal is API attribution, replay, fixture building, or regression comparison.

#### Does Not Apply When

- Authorization does not allow automated operation.
- The flow includes high-risk side effects that cannot be safely sandboxed.
- The app is unstable under automation; use manual single-step operation and timestamps instead.

#### Validation

- Re-running the operation script produces the same main method/path or a documented cache/preload difference.
- Capture logs can be bounded by the script's start/end timestamps.
- API documentation records `operation_id`, capture window, source, response shape, and trigger confidence.

#### Promotion Target

- `WORKFLOW.md`
- `TOOLS.md`
- `DOCUMENTATION.md`
- `techniques/http-api/README.md`

#### Required Linked Updates

- Updated `WORKFLOW.md` lightweight UI map guidance with optional automation scripts.
- Updated `TOOLS.md` with a safe adb operation script template.
- Updated `DOCUMENTATION.md` UI architecture map template with automation fields.
- Updated `techniques/http-api/README.md` with UI automation capture flow.
- Updated `feedback_history/http-api/README.md` and root feedback index.
