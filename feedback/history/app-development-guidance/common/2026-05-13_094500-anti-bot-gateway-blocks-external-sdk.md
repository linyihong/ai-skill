> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-13 - Anti-Bot Gateway Blocks External SDK Calls via TLS Fingerprint

Status: candidate

#### One-line Summary

當目標 API 使用 anti-bot gateway（如 PerimeterX）時，外部 JVM SDK 無法直接呼叫 API，因為 gateway 會驗證 TLS fingerprint；唯一繞過方式是透過目標 App 本身的 proxy 轉發請求。

#### Human Explanation

在逆向分析一個有 anti-bot 保護的 App 時，常見的假設是「只要把 API 的 request/response 格式搞清楚，就可以用外部 SDK 直接呼叫」。但這個假設在遇到 anti-bot gateway 時會失敗：

1. **PerimeterX / DataDome / Cloudflare 等 anti-bot 服務** 會在 TLS handshake 階段採集客戶端的 TLS fingerprint（JA3/JA3S）
2. 外部 JVM（Java SDK、OkHttp、curl 等）的 TLS fingerprint 與目標 App 的 fingerprint 不同
3. Gateway 直接拒絕（HTTP 403 / 429 / challenge page），不會進入業務邏輯

**實務案例**：TATA App 使用 PerimeterX 保護 guest login API。外部 Java SDK 發送的請求即使 request body、headers、cookies 完全正確，仍然被 PerimeterX 阻擋（HTTP 403）。繞過方式是利用 `adb forward tcp:18881 tcp:18881` 將裝置上的 Netty proxy 埠轉發到 host，然後透過這個 proxy 發送請求 — proxy 使用 App 本身的 TLS stack，fingerprint 與 App 一致。

#### Trigger

- 外部 SDK 發送的請求被 API gateway 回 HTTP 403 / 429 / challenge page
- Request body、headers、cookies 都與 App 實際發送的請求一致，但仍然被拒絕
- 同一個請求透過 App 內部的 proxy 轉發就可以成功

#### Evidence

- Tool: Frida hook (SSL_write/SSL_read), tcpdump capture, Java SDK live test
- Sanitized excerpt: 外部 JVM 發送 guest login POST 請求 → PerimeterX 回 HTTP 403；透過 `adb forward tcp:18881 tcp:18881` + device proxy 發送相同請求 → HTTP 200 + Set-Cookie
- Evidence path: `<PROJECT_ROOT>/TATA/docs/domain-baseline.md`（PerimeterX 章節）、`<PROJECT_ROOT>/TATA/api/API列表/public/guest_login.md`（2026-05-13 發現記錄）

#### Generalized Lesson

當開發 SDK 或外部 client 呼叫有 anti-bot 保護的 API 時：

1. **先確認 API 是否有 anti-bot gateway**：檢查 response 是否包含 `_px` 相關 cookie、`x-px` header、或 403 challenge page
2. **TLS fingerprint 是繞不過的**：外部 JVM 的 TLS stack 與 App 不同，gateway 可以輕易識別
3. **唯一繞過方式**：透過 App 本身的 proxy/agent 轉發請求，讓 proxy 使用 App 的 TLS stack
4. **Proxy 轉發的侷限**：需要裝置在線、adb 連線、且 proxy 埠可轉發；不適合 production 使用，僅適用於開發/測試階段

#### Agent Action

下次 agent 遇到「外部 SDK 無法呼叫 API」的情況時：

1. 先檢查 API response 是否有 anti-bot gateway 的跡象（`_px`、`x-px`、403 challenge）
2. 如果有 anti-bot gateway，不要花時間嘗試「修正 request 格式」，因為問題不在 request 內容而在 TLS fingerprint
3. 直接建議使用 App 內部的 proxy 轉發，或告知使用者外部 SDK 無法直接繞過 anti-bot 保護
4. 記錄這個限制到專案文件，避免後續開發者重複嘗試

#### Goal / Action / Validation

- Goal: 讓 agent 在遇到 anti-bot gateway 時能快速識別並給出正確的解決方案
- Action: 檢查 API response 是否有 anti-bot 特徵 → 若有，告知使用者外部 SDK 無法直接繞過 → 提供 proxy 轉發方案
- Validation or reference source: 實測驗證 — 相同 request 透過 device proxy 成功、直接發送失敗

#### Applies When

- 開發外部 SDK 或 client 呼叫有 anti-bot 保護的 API
- 逆向分析中發現 API 使用 PerimeterX / DataDome / Cloudflare 等 anti-bot 服務
- 需要從外部環境（非 App 本身）發送請求到受保護的 API

#### Does Not Apply When

- API 沒有 anti-bot 保護（直接 HTTP 200）
- 可以在 App 內部執行 SDK 程式碼（如同一個 process）
- 目標是測試/開發環境，anti-bot 已關閉

#### Validation

- 發送相同 request 到 API：直接發送 vs 透過 device proxy 轉發
- 如果直接發送成功（HTTP 200），表示沒有 anti-bot 保護，此 lesson 不適用
- 如果直接發送失敗但透過 proxy 成功，確認此 lesson

#### Promotion Target

- `intelligence/engineering/development/anti-bot-gateway.md`（新分層）
- `workflow/app-development-guidance/execution-flow.md`（SDK defect closure loop 中增加 anti-bot 檢查步驟）

#### Required Linked Updates

- 無需連動更新；此為新 lesson，尚未 promote 到正式文件
