> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-13 - Non-Standard TLS Implementation Bypasses Standard SSL Hooks

Status: candidate

#### One-line Summary

當 App 使用非標準的 TLS 實作（非 OpenSSL、非 Java SSLSocketFactory、非 Netty SslHandler），標準的 SSL_write/SSL_read Frida hook 不會觸發，需要改用 libc write/sendto hook 或 Netty channelRead 層級來捕獲解密後的資料。

#### Human Explanation

在逆向分析 App 網路通訊時，標準做法是 hook SSL 函式庫的 `SSL_write` / `SSL_read` 來捕獲解密後的明文。但這個做法依賴一個假設：App 使用標準的 TLS 實作（OpenSSL / BoringSSL / Java SSLSocket / Netty SslHandler）。

**實務案例**：TATA App 內部的 Netty proxy 使用非標準的 TLS 實作。Frida 的 `SSL_write` / `SSL_read` hook 在 7 分鐘的 capture 中產生 **0 個事件**，即使 App 確實有 TLS 加密連線到 production API。最終捕獲明文的方式是：

1. **hook libc `write` / `sendto`** — 在系統呼叫層級捕獲 Netty 寫入 socket 的資料（但需要自行處理 HTTP 分幀）
2. **hook Netty `channelRead`** — 在 Netty pipeline 層級捕獲 `HttpObjectAggregator` / `HttpResponseDecoder` 處理後的 HTTP 訊息（已解密、已組裝）

**關鍵教訓**：不要假設 App 使用標準 TLS。先確認 TLS 實作方式，再選擇 hook 策略。

#### Trigger

- Frida 的 `SSL_write` / `SSL_read` hook 註冊成功但從不觸發
- App 確實有網路連線（tcpdump 看到 TLS 流量），但 SSL hook 沒有事件
- App 使用自製的 Netty proxy 或自訂的 TLS wrapper

#### Evidence

- Tool: Frida (hook_self_generation_phase1_v10.js), tcpdump
- Sanitized excerpt: `SSL_write` / `SSL_read` hook 在 7 分鐘 capture 中產生 0 事件；改用 libc `write` hook 後成功捕獲 HTTP 明文
- Evidence path: `<PROJECT_ROOT>/capture/self_gen_phase1_v10_analysis.md`（Key Finding: Native SSL hooks registered but NEVER Fire）

#### Generalized Lesson

當需要捕獲 App 的 HTTPS 明文時：

1. **不要假設 TLS 實作是標準的** — 先測試 SSL hook 是否觸發
2. **分層 hook 策略**：
   - Layer 1: SSL_write/SSL_read（OpenSSL/BoringSSL）— 最快，但可能不觸發
   - Layer 2: libc write/sendto（系統呼叫層）— 通用，但需要自行處理 HTTP 分幀
   - Layer 3: Netty channelRead（應用層）— 最可靠，捕獲已解密的 HTTP 訊息
3. **Netty channelRead 是最穩定的選擇**：如果 App 使用 Netty，`channelRead` 上的 `HttpObjectAggregator` / `HttpResponseDecoder` 會提供完整的已解密 HTTP 請求和回應
4. **非標準 TLS 的 fingerprint 也是 anti-bot 繞過的關鍵**：因為 TLS 實作是自訂的，其 fingerprint 與任何標準 library 都不同，這正是 PerimeterX 等 anti-bot 服務無法阻擋的原因

#### Agent Action

下次 agent 需要捕獲 App HTTPS 流量時：

1. 先註冊 SSL_write/SSL_read hook 並觀察是否觸發
2. 如果 SSL hook 在 30 秒內沒有事件，立即切換到 libc write/sendto hook
3. 如果 App 使用 Netty，優先使用 Netty channelRead hook（最穩定）
4. 記錄 TLS 實作方式到分析文件，避免後續分析者重複嘗試 SSL hook
5. 如果發現非標準 TLS，記錄這個事實作為「anti-bot 繞過」的潛在原因

#### Goal / Action / Validation

- Goal: 讓 agent 在 SSL hook 不工作時能快速切換到有效的 hook 策略
- Action: 分層 hook（SSL → libc → Netty channelRead），每層設定 timeout
- Validation or reference source: 實測驗證 — SSL hook 0 events → libc write hook 成功捕獲 HTTP 明文

#### Applies When

- 逆向分析 App 網路通訊，需要捕獲 HTTPS 明文
- App 使用 Netty 作為網路框架
- SSL_write/SSL_read hook 註冊成功但不觸發

#### Does Not Apply When

- App 使用標準 OpenSSL/BoringSSL（SSL hook 正常觸發）
- App 使用標準 Java SSLSocket（Java Frida hook 可捕獲）
- 只需要 TLS 層級的資訊（不需要解密後的 HTTP 內容）

#### Validation

- 註冊 SSL_write/SSL_read hook → 等待 30 秒 → 檢查事件計數
- 如果為 0，切換到 libc write hook → 確認捕獲到 HTTP 明文
- 如果 libc write 也失敗，切換到 Netty channelRead → 確認捕獲到 HttpObjectAggregator 訊息

#### Promotion Target

- `skills/apk-analysis/techniques/local-proxy/non-standard-tls-hook.md`（舊結構）

#### Required Linked Updates

- 無需連動更新；此為新 lesson，尚未 promote 到正式文件
