> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-12 - Non-Standard TLS Implementation Bypasses BoringSSL/OpenSSL Frida Hooks

Status: candidate

#### One-line Summary

當 Android 應用的代理或 HTTP 客戶端使用非標準 TLS 實作（非 BoringSSL/OpenSSL）時，標準的 `SSL_write`/`SSL_read` Frida hook 永遠不會觸發，需要改用 Netty `SslHandler` 或 Java `SSLSocketFactory` 層級的 hook。

#### Human Explanation

在分析 TATA 應用的代理架構時，發現代理伺服器（Java Netty-based）到外部 API 伺服器的 TLS 連線完全無法透過標準的 `SSL_write`/`SSL_read` hook 捕獲。經過 7 分鐘的完整 capture（v10），SSL hooks 註冊了 0 次觸發。

這不是 hook 寫錯或時機問題——而是該代理使用了非標準的 TLS 實作（非 BoringSSL/OpenSSL）。Android 通常使用 BoringSSL（Google 的 OpenSSL fork），但某些自訂 Java Netty 實作可能使用 Bouncy Castle TLS、Conscrypt、或完全自訂的 TLS 實作，這些都不會呼叫 `SSL_write`/`SSL_read`。

這意味著：
1. **標準 SSL hook 不適用於非標準 TLS**——`SSL_write`/`SSL_read` 只對 OpenSSL/BoringSSL 有效。
2. **需要更高層級的 hook**——Netty `SslHandler`（`write()`/`decode()`）或 Java `SSLSocketFactory.createSocket()`。
3. **代理對外流量不可見**——如果代理使用非標準 TLS，libc-level 只能看到加密後的 TLS record，無法看到解密後的 HTTP body。

#### Trigger

- Frida hook 註冊了 `SSL_write`/`SSL_read` 但長時間（>7 分鐘）完全沒有觸發。
- 代理伺服器使用 Java Netty，但 Netty `SslHandler` 的 hook 也無法捕獲解密資料。
- 代理的對外連線是 TLS 加密的，但無法透過任何標準 SSL/TLS hook 看到解密後的內容。
- 代理版本標頭為 `Server: PWS/8.3.1.0.8`，這不是標準的 HTTP 伺服器。

#### Evidence

- **工具**：Frida 17.9.1，hook `libssl.so`/`libboringssl.so` 的 `SSL_write`/`SSL_read` export。
- **去敏摘要**：
  - v9 capture（7 分鐘）：SSL hooks 註冊 0 次觸發。
  - v10 capture（7 分鐘）：SSL hooks 註冊 0 次觸發。
  - 代理使用 Java Netty 實作，類別為 `cn.stool.skyshield.proxy.ProxyServerHandler`。
  - 代理版本：`PWS/8.3.1.0.8`。
  - 代理對外連線使用非標準 TLS（非 BoringSSL/OpenSSL）。
- **證據路徑**：
  - `capture/self_gen_phase1_v9_analysis.md` — SSL hooks 0 觸發的分析。
  - `capture/self_gen_phase1_v10_analysis.md` — SSL hooks 0 觸發的確認。
  - `TATA/scripts/frida/hook_self_generation_phase1_v10.js` — SSL hook 實作（`tryRegisterSslHooks`）。

#### Generalized Lesson

當分析 Android 應用的 TLS 流量時，如果標準的 `SSL_write`/`SSL_read` hook 無法捕獲任何資料：

1. **先確認是否有 TLS 流量**：檢查 libc `connect()` hook 是否顯示連線到外部伺服器的 443 port。
2. **檢查 SSL hook 是否正確註冊**：確認 `Module.findExportByName` 成功找到 `SSL_write`/`SSL_read` 的位址。
3. **如果 SSL hook 註冊成功但從未觸發**，可能原因：
   - 應用使用非標準 TLS 實作（Bouncy Castle TLS、Conscrypt、自訂 TLS）。
   - 應用使用 Netty 的 `SslHandler`，但底層不是 BoringSSL/OpenSSL。
   - 代理伺服器使用自訂 TLS 實作（如 `PWS/8.3.1.0.8`）。
4. **替代方案**：
   - hook Netty `SslHandler.write()`（outbound，加密前）和 `SslHandler.decode()`（inbound，解密後）。
   - hook `javax.net.ssl.SSLSocketFactory.createSocket()`。
   - hook `java.net.Socket` 的 `connect()` 和 `getOutputStream()`/`getInputStream()`。
   - 如果以上都無法捕獲解密資料，考慮在代理層級（如 Netty `channelRead`）捕獲 HTTP 請求和回應。

**偵測模式**：
- 檢查代理伺服器版本標頭（如 `Server: PWS/8.3.1.0.8`）。
- 檢查 Java 類別名稱是否包含非標準 TLS 實作（如 `BouncyCastle`、`Conscrypt`、`SkyShield`）。
- 如果 `SSL_write`/`SSL_read` hook 註冊成功但從未觸發，且 libc `connect()` 顯示有外部 TLS 連線，則很可能使用了非標準 TLS。

#### Agent Action

1. 先嘗試標準 SSL hook（`SSL_write`/`SSL_read`）捕獲 TLS 流量。
2. 如果長時間（>5 分鐘）無觸發，但 libc `connect()` 顯示有外部 TLS 連線：
   - 檢查應用是否使用非標準 TLS 實作。
   - 嘗試 Netty `SslHandler` hook。
   - 嘗試 `SSLSocketFactory.createSocket()` hook。
3. 如果高層級 hook 也無法捕獲解密資料：
   - 考慮在代理層級（Netty `channelRead`）捕獲 HTTP 請求和回應。
   - 記錄代理使用的 TLS 實作類型，以便後續分析。
4. 記錄無法捕獲解密資料的限制，並在分析文件中明確標註。

#### Goal / Action / Validation

- **目標**：當標準 SSL hook 無法捕獲 TLS 流量時，能快速診斷是否為非標準 TLS 實作，並改用正確的 hook 層級。
- **行動**：
  1. 檢查 libc `connect()` 是否有外部 TLS 連線（443 port）。
  2. 檢查 SSL hook 是否正確註冊。
  3. 如果 SSL hook 註冊成功但無觸發，嘗試 Netty `SslHandler` 或 `SSLSocketFactory` hook。
  4. 如果仍無法捕獲，在代理層級捕獲 HTTP 請求和回應。
- **驗證或參考來源**：v9/v10 capture 確認 SSL hooks 註冊成功但 0 觸發，同時 libc `connect()` 顯示有外部 TLS 連線。

#### Applies When

- 分析 Android 應用的 TLS 流量。
- 標準 `SSL_write`/`SSL_read` hook 註冊成功但長時間無觸發。
- libc `connect()` hook 顯示有外部伺服器的 443 port 連線。
- 應用使用 Java Netty 或自訂代理實作。

#### Does Not Apply When

- 應用使用標準 BoringSSL/OpenSSL 進行 TLS 連線（SSL hook 正常觸發）。
- 應用不使用 TLS（純 HTTP 或非加密連線）。
- 分析的是 iOS 應用（使用不同的 SSL/TLS 實作）。

#### Validation

在認為本 lesson 已適用之前，確認以下所有項目：

- [ ] 已確認 libc `connect()` hook 顯示有外部伺服器的 443 port 連線。
- [ ] 已確認 `SSL_write`/`SSL_read` hook 正確註冊（`Module.findExportByName` 成功）。
- [ ] 已確認在足夠長時間（>5 分鐘）的 capture 中 SSL hook 完全無觸發。
- [ ] 已嘗試 Netty `SslHandler` 或 `SSLSocketFactory` hook 作為替代方案。
- [ ] 已記錄無法捕獲解密資料的限制。
- [ ] Feedback history 索引已更新。

#### Promotion Target

- `WORKFLOW.md`
- `checklists/frida-hook-checklist.md`

#### Required Linked Updates

- 更新 skill feedback 根索引與 `common/README.md`。
- 如果 lesson 來自特定專案，將專案特定的類別名稱、endpoint、port 保留在專案文件中。
