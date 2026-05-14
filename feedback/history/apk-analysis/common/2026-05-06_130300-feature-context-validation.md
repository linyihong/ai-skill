> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/apk-analysis/execution-flow.md`](../../../../workflow/apk-analysis/execution-flow.md)

### 2026-05-06 - Feature Context Validation

Status: promoted

#### One-line Summary

UI replay 不只要驗證 foreground package；同一 App 內也可能跑到錯頁，重要 checkpoint 應再驗證目標 feature context。

#### Human Explanation

Package guard 可以排除 launcher、browser、settings 等外部 App，但不能排除同一 App 內的錯誤頁面。例如分類 tab replay 若誤點浮動廣告、充值入口、內部 WebView 或其它 module，XML package 仍然是目標 App；若只看 package，就會把錯頁截圖誤算成目標 tab 成功。

#### Trigger

授權 APK UI/API 對齊時，某個 replay step 仍停在目標 package，但畫面已離開目標 feature，例如從 tab/category/list 操作跳到充值、活動、WebView、彈窗或其它 App 內 module。

#### Evidence

- Tool: UIAutomator screenshot/XML, checkpoint replay runner, package guard, feature-context text guard.
- Sanitized excerpt: category-grid one-by-one replay 中，lower-row tap 仍在目標 package，但畫面轉到非目標 module；加入 feature context guard 後，後續 replay 要求 XML 同時包含目標 feature 的穩定上下文文字。
- Evidence path: project-private capture and operation map only; reusable lesson contains no target host, token, raw API payload, article text, user data, or media URL.

#### Generalized Lesson

對 UI-to-API evidence，至少分兩層驗證：

1. `package/activity` guard：確認不是外部 App 或系統頁。
2. `feature context` guard：確認仍在目標 module/page，例如穩定 tab label、page title、section heading、route marker、grid title、selected tab 或其它不含私密值的 UI anchor。

若 package 正確但 context 不符，該 window 應標為 `wrong in-app screen` / `invalid for target feature`，不可用來支撐該 feature 的 API attribution。

#### Agent Action

下次寫 replay script：

1. Package guard 仍必須保留。
2. 對 tab sweep、grid picker、detail/media、search result 等 checkpoint，加 feature-specific context guard。
3. Context guard 只使用穩定、非敏感 UI anchor；避免把文章標題、留言、使用者名稱等 raw content 寫進通用規則。
4. 若跑到同 package 的其它 module，先修 selector/bounds/overlay handling，再重抓，不要把該 window 算成功。

#### Goal / Action / Validation

- Goal: 避免把同一 App 內錯頁誤當目標 feature 成功 evidence。
- Action: 在 package validation 後追加 feature context validation。
- Validation or reference source: checkpoint XML 同時包含目標 package 與目標 feature 的穩定 anchor；錯頁應中止或被標 invalid。

#### Applies When

- App 內有浮動廣告、充值入口、活動頁、WebView、內部跳轉或多 module 入口。
- UI replay 需要支撐 API/UI attribution。
- 操作含座標 fallback、scroll、grid picker、tab strip 或容易誤點的 lower-row controls。

#### Does Not Apply When

- 只是粗略確認 App 是否啟動，不做 feature attribution。
- 目標 feature 沒有穩定非敏感 anchor，只能先人工檢視或改用其它 evidence。
- 分析本身就是要記錄跨 module / external transition，且已明確標示轉場。

#### Validation

- 成功 checkpoint 同時通過 package guard 與 feature context guard。
- 同 package wrong-screen capture 會被 runner abort 或文件標 invalid。
- 文件只引用通過兩層 guard 的 target-feature windows。

#### Promotion Target

- `WORKFLOW.md` / UI evidence validation and replay runner guidance.
- Project operation maps and page-level UI maps.

#### Required Linked Updates

- 已同步更新 `WORKFLOW.md` 的 UI evidence validation rule。
- 已更新 `feedback_history/README.md` 與 `feedback_history/common/README.md` 索引。
