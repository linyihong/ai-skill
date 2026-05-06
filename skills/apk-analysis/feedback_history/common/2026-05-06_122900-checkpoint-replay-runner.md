> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-06 - Checkpoint Replay Runner

Status: promoted

#### One-line Summary

重複驗證 UI/API flow 時，把已確認路徑寫成可停在 checkpoint 的 replay runner；跑歪時先修 selector / bounds，再重抓證據。

#### Human Explanation

APK UI 動態分析常需要反覆回到同一頁、同一 tab 或同一媒體區塊。如果每次都讓 agent 重新推理點擊順序，容易浪費時間，也容易把偶發 drift 誤判成 API 行為差異。更穩定的做法是把路徑固化成腳本，讓腳本在每個 checkpoint 截圖、dump hierarchy、驗證 foreground package，並支援 `--target` 類參數提前停住。

#### Trigger

已經知道從 launch 到目標 feature 的 tap/swipe 順序，而且需要重複測試 Frida hook、media capture、tab coverage、search replay 或 reset baseline。使用者或分析者發現「每次重新用 AI 分析路徑」比固定腳本慢，且 UI 漂移時需要可定位是哪一步跑歪。

#### Evidence

- Tool: UIAutomator / adb replay script, screenshot + XML capture, foreground package validation, optional Frida hook.
- Sanitized excerpt: runner can stop at `gossip-tab`, `detail`, `detail-media`, or `comments`; smoke test reached a target tab and emitted `[target-ok]` after all XML package checks passed.
- Evidence path: project operation map and private capture outputs only; reusable lesson contains no target host, token, UI text values, raw API payloads, or media URLs.

#### Generalized Lesson

對已探索過的 UI flow，建立「checkpoint replay runner」：

1. 每個主要節點命名成 checkpoint，例如 `launch`、`feature-tab`、`category`、`detail`、`media-probe`。
2. 每到一個 checkpoint 就截圖、dump XML、驗證 foreground package / activity。
3. 提供 `--target` 或等效參數，讓測試只跑到本次需要的頁面。
4. 需要 Frida/pcap/MITM 時，把 hook/capture 接在同一 runner 的固定 timing 上。
5. 如果跑歪，先修 selector、fallback coordinate、等待時間或 scroll 次數，再把新 capture 當成有效證據。

#### Agent Action

下次重複分析同一 App 功能時：

1. 在路徑已穩定後，優先固化成 replay script，不要每輪重新推理全部 UI 操作。
2. 為重要節點加 `--target`/checkpoint，方便快速回到特定頁面。
3. 每個 checkpoint 必須有 package validation；外部 app / browser / launcher 畫面要中止或標 invalid。
4. 文件引用 replay 結果時，引用 operation id、checkpoint、capture tag 與驗證結果，而不是只寫「點到某頁」。

#### Goal / Action / Validation

- Goal: 提高重複 capture 的速度與可驗證性，減少 UI drift 對 API attribution 的污染。
- Action: 將已知路徑寫成 checkpoint replay runner，支援定點停止、截圖/XML、package guard 與可調整 scroll/tap 參數。
- Validation or reference source: runner 到達目標 checkpoint 後輸出成功標記，且該 checkpoint 的 XML foreground package/activity 屬於目標 App。

#### Applies When

- 同一 feature/page 需要反覆跑 Frida hook、API schema、media decrypt、tab sweep 或 reset baseline。
- UI path 已經被至少一次人工或半自動操作確認。
- 裝置解析度、語系或測試帳號狀態相對固定，或 runner 已把差異參數化。

#### Does Not Apply When

- 尚未知道入口路徑，需要先探索 UI 架構。
- 目標行為高度依賴即時推薦、外部 App、WebView 跳轉或一次性驗證流程，固定 replay 會造成誤導。
- 沒有能力驗證 foreground package/activity，無法判斷 capture 是否仍在目標 App。

#### Validation

- 至少一次 smoke test 可從 launch 跑到指定 checkpoint 並輸出成功標記。
- 每個保存為證據的 screenshot/XML 都通過 target package guard。
- 跑歪時可從最後一個成功 checkpoint 定位需要修正的 selector、coordinate、wait 或 scroll。

#### Promotion Target

- `WORKFLOW.md` / UI architecture map and operation replay guidance.
- Project operation maps that need stable API/UI attribution.

#### Required Linked Updates

- 已同步更新 `WORKFLOW.md` 的 replay runner guidance。
- 已更新 `feedback_history/README.md` 與 `feedback_history/common/README.md` 索引。
