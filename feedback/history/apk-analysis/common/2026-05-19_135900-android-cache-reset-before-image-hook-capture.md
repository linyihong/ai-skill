> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-19 - Android Cache Reset Before Image Hook Capture

Status: candidate

#### One-line Summary

當影像解密 hook 看不到事件但 UI 已顯示圖片時，先清記憶體與 App cache 再冷啟動抓取；不要直接用會清登入狀態的 `pm clear`。

#### Human Explanation

Android App 的圖片載入常同時有記憶體快取、磁碟 cache、Flutter image cache 或自定義解密 cache。若目標圖片已在 UI 中可見，後續滑動或切 tab 可能只命中 cache，不會再進入 `ImageLoader`、decryptor 或 HTTP 路徑。此時應把「狀態重置」視為動態分析的一部分：先 `am force-stop` 清記憶體，再用 Frida 早期注入，在 App 進程內刪除 `getCacheDir()`、`getExternalCacheDir()`、`getCodeCacheDir()` 等 cache 目錄，重新導向使用者進入目標畫面。

#### Trigger

- UI 已能顯示目標圖片，但 `_addToMemoryCache`、`getCachedImageData` 或 decryptor hook 沒有新事件。
- 重新滑動列表、切換 tab、進出 detail 仍只看到既有 fixture，目標 hash 不出現。
- HTTP 端能下載加密圖片，但 App 端已快取解密結果，導致 hook 無法觀察解密鏈。
- 需要保留登入、裝置綁定或實驗狀態，不適合執行 `pm clear <package>`。

#### Evidence

- Tool: `adb shell am force-stop`、Frida spawn/early attach、Java `ActivityThread.currentApplication()`、`Context` cache dirs。
- Sanitized excerpt: Frida 啟動期輸出 `[cache-reset] cleared app cache dirs` 後，下一輪 UI 進入目標列表才重新出現 image-cache hook 事件。
- Evidence path: 專案內 `capture/`、`scripts/frida/`、runner log；lesson 不含 package name、裝置 serial、host、auth query 或 sample ID。

#### Generalized Lesson

1. **分清 reset 層級**：`am force-stop` 清記憶體；進程內刪 cache dirs 清磁碟 cache；`pm clear` 會清 App data，除非使用者明確允許，否則不要用。
2. **早期注入優先**：cache 清除與 hook 安裝應在使用者進入目標畫面前完成，否則首次載入可能又先被快取吃掉。
3. **保留互動窗口**：runner 要留足操作時間，讓使用者手動進入目標 tab/detail；自動 tap 只能做輔助，不可當唯一證據。
4. **hook 要最小化**：優先 hook image cache 邊界；decryptor onLeave 或全量 Dart AOT hook 若曾造成崩潰，改為可選開關。
5. **驗證不要只看 UI**：完成後必須用檔名 hash、magic bytes、尺寸或 fixture index 驗證目標圖片真的落地。

#### Agent Action

1. 啟動 capture 前確認前景 App 與目標畫面路線。
2. 執行 `am force-stop <package>` 清記憶體，避免舊 Flutter image cache 影響。
3. 用 Frida spawn 或早期 attach 執行 Java cache reset：
   - 取得 current `Application`。
   - 取得 `getCacheDir()`、`getExternalCacheDir()`、`getCodeCacheDir()`。
   - 遞迴刪除目錄內容。
4. 安裝最小 image hook，讓使用者手動進入目標畫面。
5. Pull artifacts，檢查目標 hash 或 manifest；若缺失，回到 UI 路線或 hook 邊界，不直接加重 hook。

#### Goal / Action / Validation

- Goal: 消除記憶體與磁碟 cache 對 image/decrypt hook 的遮蔽，讓目標圖片重新走可觀察路徑。
- Action: `force-stop` + 進程內 cache dirs delete + early hook + bounded manual navigation。
- Validation: runner log 有 cache reset 訊息；capture 期間有新的 image hook event；目標 artifact 以 hash/magic/fixture index 驗證存在；App data 未被清除。

#### Applies When

- Android App 圖片、封面、avatar、富文本圖等資源有 App 端解密或自定義 cache。
- UI 已顯示資源但動態 hook 沒有事件。
- 需要保留登入與 App data，只能清 cache。

#### Does Not Apply When

- 目標是登入流程、偏好設定或 database state，本來就需要完整清資料重跑。
- App cache 位於外部共享目錄且需人工授權或 root 才能刪除。
- 圖片由系統瀏覽器或外部 WebView 進程載入，目標 App cache reset 無法涵蓋。

#### Validation

- 使用去敏 runner log 確認 cache reset 已執行。
- 至少一次新 capture 產生非舊 timestamp 的 artifact。
- 若 capture 仍缺目標 hash，記錄為 UI 路線或 hook 邊界問題，而不是宣稱圖片不存在。

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md`（state reset / dynamic capture 小節）
- `analysis/apk/workflows/frida-hook-flow.md`（Frida image/decrypt hook 前置 reset）

#### Required Linked Updates

- `feedback/extraction/apk-analysis-index.md` 加入本 lesson 的 candidate row。
- `analysis/apk/workflows/frida-hook-flow.md` 已補入 image/resource capture 前置 cache reset 流程。
