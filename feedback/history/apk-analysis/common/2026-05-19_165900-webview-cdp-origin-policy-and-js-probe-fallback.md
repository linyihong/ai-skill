> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../feedback/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-19 - WebView CDP Origin Policy And JS Probe Fallback

Status: candidate

#### One-line Summary

WebView response capture 不要只押 `evaluateJavascript` 或程式化 CDP；JS probe 可能讓 Frida agent crash，CDP WebSocket 也可能被 Chrome/WebView 的 `remote-allow-origins` policy 擋住。

#### Human Explanation

嵌入式 H5 API 分析在 request path 已觀測後，下一步常想抓 response body。兩條看似方便的路徑都有風險：第一，JS `fetch` / `XMLHttpRequest` probe 需要注入 WebView runtime，特定裝置或 WebView 版本可能因 Frida agent / JNI callback lifecycle 問題崩潰；第二，`WebView.setWebContentsDebuggingEnabled(true)` 只能讓 `webview_devtools_remote_<pid>` 出現，不保證自製 CDP client 能連上，Chrome/WebView 可能拒絕非允許 origin 的 WebSocket handshake。即使已把代理 CA 放進 system trust store，特定 H5 API host 仍可能在 WebView 端觸發 SSL error，需要以短窗、可回復的方式處理。

比較穩定的決策是：先用 `shouldInterceptRequest` 完成 request metadata；若 response/schema 是必要 gate，優先準備 MITM 或人工 Chrome DevTools frontend。程式化 CDP 可嘗試，但若遇到 `remote-allow-origins` 403，不要在同一輪反覆調整不確定的 WebView command-line 旗標，避免改壞測試設備狀態。若 MITM 只剩 WebView SSL error，可用 opt-in `onReceivedSslError(...).proceed()` 作為授權短窗 schema probe，但不得常態啟用。

#### Trigger

- WebView H5 endpoint 已用 `shouldInterceptRequest` 看到 request path / headers，但缺 response body。
- Agent 打算啟用 JS response probe、`evaluateJavascript`、或自製 CDP WebSocket client 來抓 schema。
- CDP `/json` 能看到 target，但 WebSocket handshake 回 403，訊息要求 `--remote-allow-origins=*` 或指定 origin。
- MITM 能看到 target host，但 client TLS handshake 顯示不信任 proxy certificate，即使代理 CA 已進 system CA。

#### Evidence

- Tool: Frida Java WebView hook, Android WebView DevTools socket forwarding, CDP `/json`, Python WebSocket handshake, mitmproxy response wrapper extraction.
- Sanitized excerpt: `webview_devtools_remote_<pid>` 可列出 H5 target；自製 WebSocket 連線可能因 origin policy 被拒絕。`JS_PROBE=1` 在特定 Pixel / WebView 版本上造成 Frida agent stack overflow crash。對單一 H5 API host，`onReceivedSslError` 短窗 proceed 可讓 MITM 取得 HTTP 200 JSON wrapper，但仍只記 schema，不保存 body。
- Evidence path: 專案 raw log、target URL、token、user id、完整 request URL 留在 gitignored evidence；本 lesson 只保存 generalized failure mode。

#### Generalized Lesson

1. **把 JS probe 當實驗，不當主路線**：預設關閉；即使只記 schema，也可能在 Frida agent 或 WebView callback 邊界崩潰。
2. **分清 CDP target discovery 與 CDP control**：`/json` 能看到 target 只代表 DevTools socket 存在，不代表 custom WebSocket client 可以控制 Network domain。
3. **遇到 origin 403 先降級**：若 WebSocket handshake 被 `remote-allow-origins` policy 擋住，改走 MITM 或人工 DevTools frontend，不要無限嘗試設備 command-line flag。
4. **SSL proceed 必須 opt-in**：只有在授權 MITM response/schema 實驗中才 hook `WebViewClient.onReceivedSslError` 並呼叫 `handler.proceed()`；預設路徑不能改變 WebView TLS 行為。
5. **清理臨時設備旗標**：若測過 `/data/local/tmp/webview-command-line`、Android global proxy 或類似設定，結束前移除或恢復，避免污染後續 capture。
6. **文件狀態維持 pending/partial**：只要沒有 runtime 內成功 payload shape，就不要把 API 升級成 `schema-validated` 或 SDK-ready；只拿到 wrapper 可標成 partial。

#### Agent Action

1. 先用 `WebViewClient.shouldInterceptRequest` 抓 request metadata，避免直接注入 JS。
2. 若必須抓 response，嘗試 MITM 或人工 DevTools；CDP WebSocket 只作為快速探針。
3. 若 JS probe 或 CDP 失敗，記錄失敗類型：agent crash、JNI object error、origin policy 403、target missing。
4. 若 MITM 只剩 WebView SSL error，在短窗實驗中 opt-in hook `onReceivedSslError`，確認 log 有 proceed 標記，再抽取 response schema。
5. 將失敗原因和成功 wrapper 回填專案 capture checklist，並把下一步改成 MITM / manual DevTools / safe response hook。
6. 移除測試設備上的臨時 WebView command-line flag，並清掉 Android global proxy。

#### Goal / Action / Validation

- Goal: 避免在 H5 response capture gate 上因高風險 probe 卡住或誤判 SDK readiness。
- Action: 將 JS probe 與程式化 CDP 視為可失敗探針；失敗後明確降級到 MITM / manual DevTools；MITM SSL error 只用 opt-in WebView SSL proceed 短窗處理。
- Validation or reference source: 只有成功 response wrapper keys / payload shape 被捕獲後，才能升級 API schema 狀態；只看到 request metadata 或 CDP target 不算成功 schema。

#### Applies When

- Android WebView / Flutter InAppWebView 內有 H5 API。
- 需要從授權 runtime 裡取得 response body / schema。
- 裝置 WebView 版本較新，或 target WebSocket handshake 出現 origin policy 403。
- system CA 已安裝但個別 WebView API host 仍拒絕 MITM certificate。

#### Does Not Apply When

- 已有穩定 MITM HAR 或代理能直接取得 response body。
- App 使用可控測試 build，能安全加入 debug bridge 或內建 logging。
- 只需要 H5 entry URL，不需要 H5 internal API schema。
- 測試範圍未授權修改 TLS trust / SSL error 行為。

#### Validation

- Default capture script 不啟用 JS response probe。
- `JS_PROBE=1` 或 CDP 失敗不會升級 readiness gate。
- SSL proceed 預設關閉；只有明確 opt-in 的短窗 capture 會啟用。
- 專案文件註明失敗模式與下一個可靠抓包選項。
- 臨時設備 command-line flag 與 Android global proxy 已移除或恢復。

#### Promotion Target

- `analysis/apk/workflows/frida-hook-flow.md` for WebView response-capture fallback guidance.
- `workflow/apk-analysis/execution-flow.md` for adding CDP-origin-policy and JS-probe-crash as known downgrade reasons.

#### Required Linked Updates

- Project H5 capture docs should mention whether JS probe is safe on the tested device / WebView version.
- If CDP target discovery works but WebSocket control fails, record this separately from "DevTools unavailable".
- If MITM only works with WebView SSL proceed, project docs should mark it as an experiment mode and keep payload/body values out of tracked files.
- `feedback/history/apk-analysis/README.md` category count should be updated with this lesson.
