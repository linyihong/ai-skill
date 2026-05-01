> 遵守 [共用規則索引](../../../shared-rules/README.md) 與 [feedback-lessons](../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-01 - UI architecture map from screenshots

Status: promoted

#### One-line Summary

在可控制裝置時，先用去敏 screenshot、UI hierarchy 與操作時間窗建立 App 架構地圖，再把 API 對應回具體 tab、screen 與 action。

#### Human Explanation

只列出 API endpoint 會讓分析結果難以使用，因為後續重放、SDK 實作或產品理解通常需要知道「哪個畫面、哪個 tab、哪個操作」觸發了請求。Screenshot 能快速建立 bottom tab、drawer、列表、詳情頁與播放器等可見架構；但 screenshot 不是 API 證據，必須與 pcap/MITM/Frida 的時間窗或 sequence 一起使用。

常見誤判是看到某個畫面後，就把同一時間出現的所有 request 都歸因給這個 screen。啟動預載、背景同步、快取刷新與多 screen 共用 endpoint 都可能混在同一段流量裡，所以文件要保留 trigger confidence 與 unknowns。

Revision 2026-05-01:

UI map 不一定要先做完整截圖或完整遍歷。若 screenshot、screenrecord、UI dump 或自動操作導致 App、裝置、proxy、hook 變慢，應改成輕量架構盤點，或先解核心 API/response/token，再回頭只對高價值 endpoint 補 UI binding。

#### Trigger

遇到以下情境時應使用：

- 已能操作 App 或 emulator，但 API 文件缺少 UI path。
- 抓到 request/response，卻不知道是哪個 tab、screen 或 action 觸發。
- 要把 APK 分析結果轉成可重放的 API/reference docs、SDK mapping 或測試案例。
- 使用者要求「靠截圖看出 app 架構」、「知道有哪些 tab」或「判斷 API 是在哪個操作上」。

#### Evidence

- Tool: `adb screencap`, `uiautomator dump`, pcap/MITM/Frida log time window.
- Sanitized excerpt: `Operation open-detail: Home > item tap, window <start-end>, POST /<path>, source hook, response top-level keys only`.
- Evidence path: 專案私有 evidence 目錄；reusable skill 只保存模板與去敏方法。

#### Generalized Lesson

APK traffic analysis should include a UI architecture map whenever the app can be operated. The map should list visible navigation, screen inventory, screenshots, operation IDs, capture windows, and an operation-to-API matrix. API attribution requires runtime evidence such as hook sequence, pcap/MITM timing, or replay validation; screenshots alone only establish UI context.

#### Agent Action

下次 agent 進行授權 APK 分析時，如果能控制裝置或使用者提供 app 截圖，應先建立 `App Architecture Map`：

- 盤點 bottom tabs、top tabs、drawer/menu、主要列表、詳情頁、播放器/媒體頁與設定頁。
- 先決定 capture strategy：`lightweight overview`、`API-first then bind` 或 `full operation map`。
- 若截圖或 UI 遍歷造成卡頓，暫停批量截圖，保留 API hook/pcap 主線。
- 為每個操作建立 stable operation id 與 UI path。
- 每次只操作一個 screen/action，記錄操作前後時間戳。
- 將 request/response 依 operation id 回填到 `Operation To API Matrix`。
- 對預載、背景同步、cache/local-only、未確認 trigger 標出 confidence，不要硬寫成確定操作。

#### Applies When

- 有授權控制實機、emulator 或遠端裝置。
- 有 screenshot、screen recording、UI hierarchy、accessibility labels 或可手動觀察的 UI。
- 分析目標包含 API mapping、SDK 行為、產品架構、重放測試或文件化。

#### Does Not Apply When

- 只能做純靜態分析且沒有畫面/操作證據。
- 授權範圍不允許操作帳號、截圖或保存畫面。
- Screenshot 含敏感個資且無法安全去敏；此時只記 abstract screen label 與操作，不保存圖片。

#### Validation

- 每個 operation id 至少有 screenshot/UI label、操作時間窗、以及一筆可對齊的 pcap/MITM/hook/replay 證據。
- 同一操作重跑時，主要 API path 或 schema 能穩定重現；若不穩定，文件標出 background/preload/cache 可能性。
- API 文件中的 `UI path`、`Operation ID`、`Source` 與 `Trigger confidence` 欄位完整。

#### Promotion Target

- `WORKFLOW.md`
- `TOOLS.md`
- `DOCUMENTATION.md`
- `SKILL.md`

#### Required Linked Updates

- 已更新 `SKILL.md` Quick Start，加入 UI architecture map 與 operation-to-API matrix。
- 已更新 `WORKFLOW.md`，新增 UI 架構地圖步驟。
- 已更新 `TOOLS.md`，加入 screenshot/UI hierarchy 工具與命令。
- 已更新 `DOCUMENTATION.md`，加入 App Architecture Map 與 Operation To API Matrix 模板。
- 已更新 `feedback_history/README.md` 索引。
