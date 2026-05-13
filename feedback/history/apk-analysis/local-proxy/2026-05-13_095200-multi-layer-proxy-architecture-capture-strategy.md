> 遵守 [共用規則索引](../../../shared-rules/README.md)、[dependency-reading](../../../shared-rules/dependency-reading.md)、[neutral-language](../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-13 - Multi-Layer Proxy Architecture Requires Layered Capture Strategy

Status: candidate

#### One-line Summary

當 App 使用多層代理架構（OkHttp health check → proxy1 → proxy2 → external TLS），單一 hook 層級無法捕獲完整流量，需要針對每一層使用不同的 hook 策略。

#### Human Explanation

TATA App 的網路架構包含三層代理：

```
Dart dart:io HttpClient → 127.0.0.1:55382 (proxy1, plain HTTP via OkHttp)
    ↓
proxy1 → 127.0.0.1:64571 (proxy2, internal, plaintext HTTP via Netty)
    ↓
proxy2 → [NON-Standard TLS] → api-n.52zyp.com (external, encrypted)
```

每一層使用不同的技術棧，因此需要不同的 hook 策略：

1. **Layer 1 (Dart → proxy1)**：Dart `dart:io` HttpClient 發送 plain HTTP 到 proxy1。需要 libc `connect()`/`write()` hook 才能捕獲，因為 Dart 不經過 Java 層。
2. **Layer 2 (proxy1 → proxy2)**：內部轉發，plaintext HTTP via Netty。可以透過 Netty `channelRead` on `HttpObjectAggregator`/`HttpResponseDecoder` 捕獲完整的 HTTP headers。
3. **Layer 3 (proxy2 → external)**：非標準 TLS 加密連線到外部 API。SSL_write/SSL_read hook 不觸發，需要 libc `write`/`sendto` hook 或 Netty channelRead 層級。

**關鍵教訓**：不要假設單一 hook 策略可以捕獲所有流量。必須先理解代理架構的每一層，再針對每一層選擇正確的 hook 方法。

#### Trigger

- App 使用多層代理架構（health check → proxy1 → proxy2 → external）
- 不同層使用不同的技術棧（Dart dart:io / OkHttp / Netty / 非標準 TLS）
- 單一 hook 策略只能捕獲部分流量
- 需要同時捕獲內部請求和外部回應才能理解完整流程

#### Evidence

- Tool: Frida (hook_self_generation_phase1_v10.js), tcpdump
- Sanitized excerpt: v10 capture 顯示完整的代理鏈：OkHttp health check → proxy1 (evt=38), OkHttp health check → proxy2 (evt=47), Dart dart:io guest login POST → proxy2 (evt=69), proxy2 response (evt=79, 107)
- Evidence path: `<PROJECT_ROOT>/capture/self_gen_phase1_v10_analysis.md`（Confirmed: Proxy Architecture）

#### Generalized Lesson

當分析使用多層代理架構的 App 時：

1. **先繪製代理架構圖**：確定每一層使用的技術棧（Dart/Java/Native）、通訊協定（HTTP/HTTPS/TLS）、以及埠號。
2. **為每一層選擇正確的 hook**：
   - Dart `dart:io` HttpClient → libc `connect()`/`write()` hook
   - OkHttp → Java-level OkHttp hook（`RealCall.getResponseWithInterceptorChain`）
   - Netty proxy → Netty `channelRead` on `HttpObjectAggregator`/`HttpResponseDecoder`
   - 非標準 TLS → libc `write`/`sendto` 或 Netty channelRead
3. **從最內層開始 hook**：先 hook 最接近外部 API 的層（通常是 Netty channelRead），因為這一層提供最完整的已解密 HTTP 訊息。
4. **交叉比對各層事件**：使用 fnv1a32 hash 或 request ID 比對同一請求在不同層的表現。
5. **記錄代理埠和路由**：每次 capture 都記錄代理埠號，因為動態埠會在 App 重啟後改變。

#### Agent Action

1. 在開始 hook 之前，先分析 App 的代理架構（使用 tcpdump 或 libc connect hook 快速繪製連線圖）
2. 為每一層準備對應的 hook 策略，不要只依賴單一 hook
3. 優先使用 Netty channelRead hook（如果 App 使用 Netty），因為它提供最完整的 HTTP 訊息
4. 如果 Dart `dart:io` HttpClient 參與其中，必須加入 libc hook
5. 記錄完整的代理鏈到分析文件，包括每一層的埠號、技術棧、hook 策略

#### Goal / Action / Validation

- Goal: 讓 agent 在多層代理架構中能系統性地選擇正確的 hook 策略組合
- Action: 先繪製代理架構圖 → 為每一層選擇 hook → 從最內層開始 → 交叉比對
- Validation or reference source: v10 capture 成功捕獲所有四層的流量（OkHttp health check、Dart dart:io POST、Netty channelRead response、libc write data）

#### Applies When

- App 使用多層代理架構（2+ 層代理）
- 不同代理層使用不同的技術棧
- 需要捕獲完整的請求/回應流程（不僅是單一層的流量）

#### Does Not Apply When

- App 使用單一網路層（直接 OkHttp 或直接 URL.openConnection）
- 所有流量都經過同一技術棧（例如全部是 OkHttp）
- 只需要捕獲外部 API 流量（不需要內部代理層）

#### Validation

- 使用 libc `connect()` hook 列出所有連線目標（埠號和 IP）
- 根據連線目標繪製代理架構圖
- 為每一層選擇對應的 hook 策略
- 確認每一層的 hook 都有事件產生
- 使用 fnv1a32 hash 比對同一請求在不同層的表現

#### Promotion Target

- `analysis/apk/workflows/local-proxy-hook-flow.md`（新分層）

#### Required Linked Updates

- 無需連動更新；此為新 lesson，尚未 promote 到正式文件
