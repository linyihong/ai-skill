> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-20 - WebView DevTools JSON No-JS Fallback

Status: candidate

#### One-line Summary

當 WebView JS probe 有崩潰風險、Frida request hook 又沒有輸出時，可用 WebView DevTools `/json` 作為無 JS 注入的 target metadata fallback。

#### Human Explanation

嵌入式 H5 分析常需要確認「目前是否真的在目標 WebView 頁」以及 H5 URL / route。最直接的 JS `fetch` / `XMLHttpRequest` probe 可能改變 WebView runtime，甚至造成 Frida agent 或 JNI callback 崩潰；Java WebView hook 也可能因 attach 時機、WebView wrapper、renderer 狀態或 framework 封裝而只輸出初始化，不列出現有 request。

若 Android WebView debugging 已啟用，且 device 上存在 `webview_devtools_remote_<pid>` socket，DevTools `/json` 可以讀取 target metadata（title、url、visible state、webSocketDebuggerUrl）。這不需要注入 page JS，也不需要抓 response body，適合作為「H5 target 是否存在」的低侵入確認手段。

但 `/json` 的 `title` 或 `url` 可能包含 token、uid、user name、session query、route-specific material。可重用文件只能寫去敏結論與欄位類型，不得保存原始 URL。

#### Trigger

- JS probe 已知不穩或被禁用。
- Frida WebView hook 已 attach，但只有 `[INIT]`，沒有 `loadUrl` / request metadata。
- 畫面看起來是 H5 / WebView，但 Java heap snapshot 沒列出 WebView instance。
- `/proc/net/unix` 或等價手段可看到 WebView DevTools remote socket。

#### Evidence

- Tool: Android WebView DevTools remote socket, local TCP forward, HTTP `/json`.
- Sanitized excerpt: DevTools `/json` 可顯示 target `type=page`、`visible=true`、sanitized H5 host / route；raw query 需丟棄。
- Evidence path: raw `/json` output 留在受控 project evidence 或不保存；tracked docs 只保留 sanitized target metadata。

#### Generalized Lesson

WebView target confirmation 可以分層：

1. 先用 Frida minimal hook 確認 App 不崩潰。
2. 若 hook 沒有 request output，檢查 WebView DevTools socket。
3. 只讀 `/json` target metadata，確認 page route / visibility。
4. 立即去敏，不把 raw URL、title、token、uid、user name 或 websocket URL 寫入 reusable docs。
5. `/json` 只能證明 target 存在，不等於 response schema 已捕獲。

#### Agent Action

下次遇到 WebView hook silent 且 JS probe 不安全時：

- 不要先開 JS probe。
- 不要把 silent hook 直接判成「沒有 H5」。
- 先用 DevTools `/json` 確認 target metadata。
- 若需要 response body，再另走 MITM、人工 DevTools、或已授權的短窗 response capture。

#### Goal / Action / Validation

- Goal: 在不注入 JS 的情況下確認嵌入式 H5 target 是否存在。
- Action: 讀 WebView DevTools `/json` 並只保留去敏 metadata。
- Validation or reference source: `/json` 顯示 target type / visibility / sanitized route；Frida no-JS attach 沒有造成 target App crash。

#### Applies When

- Android WebView / Flutter InAppWebView / hybrid App H5。
- 已啟用 WebView debugging 或可由 hook 啟用。
- 只需要 target confirmation 或 route metadata，不需要 response body。

#### Does Not Apply When

- WebView debugging 不可用或不允許。
- 目標是 response payload schema，`/json` 只能作輔助，不能取代 MITM / hook / fixture。
- `/json` raw output 無法安全去敏或會暴露使用者資料。

#### Validation

- tracked docs 不含 raw URL / token / uid / user name。
- project docs 明確區分 target metadata observed 與 response schema observed。
- JS probe 仍維持 disabled-by-default。

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` 的 WebView / H5 fallback path。
- `analysis/apk/workflows/` 中 WebView response-capture 或 target-confirmation workflow。

#### Required Linked Updates

- Project H5 capture checklist 可補充 `/json` 為 no-JS fallback。
- 若 promotion 到 workflow，需同步 sanitization guidance，避免 raw target URL 流入 reusable docs。
