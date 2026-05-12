> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/app-development-guidance/execution-flow.md`](../../../../workflow/app-development-guidance/execution-flow.md)

### 2026-05-12 - Dart `dart:io` HttpClient Bypasses Java-Level Frida Hooks

Status: candidate

#### One-line Summary

當 Flutter/Dart 應用使用 `dart:io` HttpClient 發送 HTTP 請求時，這些請求會完全繞過 Java 層（OkHttp、Netty、`java.net.Socket`、`URL.openConnection`），因為 Dart 的 `dart:io` HttpClient 直接使用原生 libc socket 函數（`connect`、`send`、`write`），無法透過 Java-level Frida hook 捕獲。

#### Human Explanation

在分析 TATA 應用的 guest login HTTP 請求時，發現該請求完全無法透過任何 Java-level Frida hook 捕獲——包括 OkHttp `RealCall.getResponseWithInterceptorChain`、OkHttpClient `newCall`、Netty `ProxyServerHandler.channelRead`、`java.net.Socket.connect`、`javax.net.ssl.SSLSocketFactory.createSocket`、`java.net.URL.openConnection`。

經過 v8 版本使用原生 libc `connect()` hook 後，才發現 Dart 的 `dart:io` HttpClient 直接透過原生 libc socket 函數發送請求到本地代理（`127.0.0.1:36279`），完全繞過 Java 層。這意味著：

1. **Java-level hook 無法捕獲 Dart `dart:io` HttpClient 流量**——因為 Dart 的 socket 實作直接呼叫 libc 的 `connect()`、`send()`、`write()`，不經過任何 Java class。
2. **需要原生 libc hook**——hook `libc.so` 的 `connect()`、`send()`、`write()` 才能捕獲這類流量。
3. **代理架構**——Dart `dart:io` HttpClient 發送請求到本地代理（plain HTTP），代理再轉發到外部 API 伺服器（HTTPS/TLS）。
4. **TLS 加密阻擋 body 可見性**——代理的對外連線是 TLS 加密的，所以 libc-level 只能看到加密後的 TLS record，看不到 HTTP body。

#### Trigger

- Flutter/Dart 應用的 HTTP 請求無法透過任何 Java-level Frida hook 捕獲（OkHttp、Netty、java.net.Socket、URL.openConnection）。
- 請求路徑顯示請求發送到 `127.0.0.1:<dynamic-port>`（本地代理），而非直接發送到外部 API 伺服器。
- 請求 header 包含 `user-agent: Dart/3.11 (dart:io)`，這是 Dart `dart:io` HttpClient 的典型特徵。
- 所有 Java-level hook 都正常註冊但從未觸發該請求。

#### Evidence

- **工具**：Frida 17.9.1，hook `libc.so` 的 `connect()`、`send()`、`write()` export。
- **去敏摘要**：
  - v4-v7：使用 OkHttp、Netty、java.net.Socket、URL.openConnection hook → guest login 請求完全未捕獲。
  - v8：加入 libc `connect()` hook → 捕獲到 `127.0.0.1:36279` 的連線，以及 `user-agent: Dart/3.11 (dart:io)` 的 HTTP 請求。
  - 代理架構：Dart `dart:io` HttpClient → `127.0.0.1:<port>` (proxy) → Netty `ProxyServerHandler.channelRead` → `Handler.a(req, URI)` → `https://api-n.52zyp.com` → TLS encrypted outbound。
- **證據路徑**：
  - `TATA/scripts/frida/hook_self_generation_phase1.js` — v8 加入 `hookLibcConnect()` 和 `hookLibcSend()`。
  - `capture/self_gen_phase1_v8_20260512_130953.log` — 588 events，包含 libc-connect 到 `127.0.0.1:36279` 和 libc-write 帶 `user-agent: Dart/3.11 (dart:io)`。
  - `capture/self_gen_phase1_v8_analysis.md` — 分析文件。

#### Generalized Lesson

當分析 Flutter/Dart 應用的 HTTP 流量時，如果 Java-level Frida hook（OkHttp、Netty、java.net.Socket、URL.openConnection）無法捕獲請求，應考慮以下可能性：

1. **Dart `dart:io` HttpClient**：Dart 的 `dart:io` HttpClient 直接使用原生 libc socket 函數，完全繞過 Java 層。特徵是 `user-agent: Dart/<version> (dart:io)`。
2. **解決方案**：hook `libc.so` 的 `connect()` 來捕獲所有 socket 連線，hook `send()`/`write()` 來捕獲 HTTP 請求資料。
3. **代理架構**：Dart `dart:io` HttpClient 通常發送請求到本地代理（plain HTTP），代理再轉發到外部 API 伺服器（HTTPS/TLS）。libc-level 只能看到代理對外的 TLS 加密流量。
4. **TLS 限制**：如果代理的對外連線是 TLS 加密的，libc-level 無法看到 HTTP body。需要 hook SSL/TLS 層（如 Netty `SslHandler`、`SSLSocketFactory`、或 OpenSSL/BoringSSL 的 `SSL_write`/`SSL_read`）來捕獲解密後的資料。

**偵測模式**：
- 檢查請求 header 是否包含 `user-agent: Dart/<version> (dart:io)`。
- 檢查請求是否發送到 `127.0.0.1:<dynamic-port>`（本地代理模式）。
- 如果 Java-level hook 正常註冊但從未觸發特定請求，很可能該請求來自 Dart 層。

#### Agent Action

1. 先嘗試 Java-level hook（OkHttp、Netty、java.net.Socket、URL.openConnection）捕獲 HTTP 流量。
2. 如果特定請求無法捕獲，檢查請求特徵：
   - 是否包含 `user-agent: Dart/<version> (dart:io)`？
   - 是否發送到 `127.0.0.1:<dynamic-port>`？
3. 如果是 Dart `dart:io` HttpClient，改用原生 libc hook：
   - `libc.so` 的 `connect()` — 捕獲所有 socket 連線。
   - `libc.so` 的 `send()`、`sendto()`、`write()` — 捕獲 HTTP 請求資料。
4. 如果 libc-level 只能看到 TLS 加密資料，需要 hook SSL/TLS 層：
   - Netty `SslHandler.write()`（outbound，加密前）和 `SslHandler.decode()`（inbound，解密後）。
   - 或 `javax.net.ssl.SSLSocketFactory.createSocket()`。
   - 或原生 OpenSSL/BoringSSL 的 `SSL_write()`/`SSL_read()`。
5. 記錄捕獲到的代理埠、外部 API 端點、以及請求/回應的 fnv1a32 hash 以供後續比對。

#### Goal / Action / Validation

- **目標**：當 Java-level Frida hook 無法捕獲 Flutter/Dart 應用的 HTTP 請求時，能快速診斷並改用正確的 hook 層級。
- **行動**：
  1. 檢查請求 header 是否包含 `Dart/<version> (dart:io)`。
  2. 如果是，hook `libc.so` 的 `connect()` 和 `send()`/`write()`。
  3. 如果 libc-level 只能看到 TLS 加密資料，hook SSL/TLS 層。
- **驗證或參考來源**：libc hook 成功捕獲 Dart `dart:io` HttpClient 的連線和 HTTP 請求 header（不含 body，因為 TLS 加密）。

#### Applies When

- 分析 Flutter/Dart 應用的 HTTP 流量。
- Java-level Frida hook（OkHttp、Netty、java.net.Socket、URL.openConnection）無法捕獲特定請求。
- 請求特徵符合 Dart `dart:io` HttpClient（`user-agent: Dart/<version> (dart:io)`、發送到本地代理）。

#### Does Not Apply When

- 應用使用 Java OkHttp 或 Netty 發送 HTTP 請求（非 Dart `dart:io` HttpClient）。
- 請求已經可以透過 Java-level hook 成功捕獲。
- 分析的是純 Java/Kotlin 應用（無 Flutter/Dart 層）。

#### Validation

在認為本 lesson 已適用之前，確認以下所有項目：

- [ ] 已確認請求 header 包含 `user-agent: Dart/<version> (dart:io)`。
- [ ] 已確認請求發送到 `127.0.0.1:<dynamic-port>`（本地代理）。
- [ ] 已嘗試 libc `connect()` hook 並成功捕獲 socket 連線。
- [ ] 已嘗試 libc `send()`/`write()` hook 並成功捕獲 HTTP 請求 header。
- [ ] 如果 libc-level 只能看到 TLS 加密資料，已記錄此限制並考慮 SSL/TLS 層 hook。
- [ ] Feedback history 索引已更新。

#### Promotion Target

- `WORKFLOW.md`
- `checklists/frida-hook-checklist.md`

#### Required Linked Updates

- 更新 skill feedback 根索引與 `common/README.md`。
- 如果 lesson 來自特定專案（如 TATA self-generation investigation），將專案特定的 class 名稱、endpoint、port 保留在專案文件中。
- 晉升時，從 Frida hook 策略指引建立 cross-link。
