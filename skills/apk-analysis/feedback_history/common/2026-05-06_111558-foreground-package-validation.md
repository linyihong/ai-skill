> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/apk-analysis/execution-flow.md`](../../../../workflow/apk-analysis/execution-flow.md)

### 2026-05-06 - Foreground Package Validation

Status: promoted

#### One-line Summary

UI 截圖/XML 只有在前景 package 屬於目標 App 時，才能作為該 App 的操作證據。

#### Human Explanation

APK 分析常用 `adb input` / UIAutomator / fallback bounds replay。若 App 被外部 intent、搜尋、瀏覽器、launcher、系統彈窗或 crash 帶離前景，後續 tap 仍會繼續執行，但截圖與 XML 已經不再屬於目標 App。這種情況會讓動態 hook 仍命中目標進程的背景/預載事件，卻被錯誤對齊到另一個 App 的 UI 操作。

#### Trigger

使用者提醒要確認截圖是否真的跑在目標 App。回查證據後發現一批 replay XML package 變成 launcher、Google app、Chrome，不能作為目標 App UI 證據。

#### Evidence

- Tool: UIAutomator hierarchy package check.
- Sanitized excerpt: XML root package can be `com.android.chrome` / launcher / Google app while Frida still attaches to the target process.
- Evidence path: project capture should keep invalid XML as local diagnostics only and explicitly exclude it from UI/API alignment.

#### Generalized Lesson

每個 replay step 的 screenshot / hierarchy 都要驗證 foreground package 或 activity。若 package 不屬於目標 App，也不屬於明確允許的系統 permission/dialog transition，該 capture window 必須標為 invalid/external，後續 automation 應中止，不可繼續把 tap 結果記成目標 App 行為。

#### Agent Action

下次做 UI-to-API capture 時：

1. 每次 `uiautomator dump` 後讀 XML package 集合。
2. 若 package 不含目標 package，立即 abort 或明確標 `external transition`。
3. 只有 system permission / Android dialog 等已知允許場景可暫時放行，且下一步必須回到目標 package。
4. Frida/PID 命中與 UI package 驗證要分開記錄；Frida 命中只能證明目標進程內事件，不能自動證明是當前 UI step 觸發。
5. Project docs 中不要用 invalid/external screenshots 支撐功能結論。

#### Goal / Action / Validation

- Goal: 避免把外部 App 畫面誤當目標 App 操作證據。
- Action: 在 runner / 手動流程中加入 foreground package guard。
- Validation or reference source: 每個 operation evidence 的 XML package 包含目標 package；若不包含，operation 標為 invalid/external 並重測。

#### Applies When

- 使用 adb input、UIAutomator、Appium、OCR、自動 tap/swipe 或 fallback bounds replay。
- 流程可能打開瀏覽器、搜尋、launcher、系統設定、登入頁或外部 intent。
- 同時使用 Frida hook 與 UI capture 對齊 API 觸發。

#### Does Not Apply When

- 分析目標本來就是外部 intent / browser handoff，且文件明確標記轉場。
- 只做純靜態分析，沒有 UI evidence 對齊。

#### Validation

- Operation map 中每個 replay window 都能引用 package-validated XML/screenshot。
- 若發現 package mismatch，重新 capture 或降低該證據 confidence。

#### Promotion Target

- `WORKFLOW.md`
- Project runner scripts
- Project UI operation maps

#### Required Linked Updates

- 已同步更新 `WORKFLOW.md` 的 UI evidence package validation 規則。
- 目標專案 runner 應加入 package guard；已有錯誤 capture 應在專案文件中標為 excluded/invalid diagnostics。
