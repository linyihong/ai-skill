# Local Proxy Routing Diagnosis Heuristic（本機代理路由診斷啟發式）

## 問題

APK 分析中，何時該懷疑 app 使用 local proxy/loopback 路由？如何區分「Wi-Fi MITM 看不到流量是因為 local proxy」vs「是因為 TLS pinning」？

## 原則

- Local proxy 會讓 Wi-Fi MITM 看不到 business CONNECT，因為流量先經過本機 loopback
- TLS pinning 會讓 MITM proxy 連線失敗（SSL handshake error），但 pcap 仍看得到 CONNECT
- 兩者的處理方式完全不同：local proxy 需要 hook handler 層，TLS pinning 需要 bypass 或 Frida hook
- 不要將缺少 Wi-Fi MITM traffic 直接視為 pinning

## 決策表

### 判斷 Local Proxy vs TLS Pinning

| 情境 | 診斷 | 下一步 |
|------|------|--------|
| Wi-Fi MITM 看不到 business CONNECT，但 app 正常運作，pcap 顯示外部 TLS | 可能是 local proxy | 檢查 loopback 證據（pcap/logcat） |
| pcap 顯示 `127.0.0.1:<port>` 或 loopback HTTP | 確認 local proxy | 開始 handler hook flow |
| MITM proxy 出現 SSL handshake error，pcap 有 CONNECT | 可能是 TLS pinning | 檢查 pinning bypass 或 Frida hook |
| Logcat 提到 `ProxyServer`、Netty handler、local bridge、TUN、sing-box | 確認 local proxy | 找出 handler 類別 |
| Hooking OkHttp/WebView 不暴露 business upstream | 可能是 local proxy | 檢查 local handler objects |
| 同時有 loopback 證據 + MITM SSL error | 兩者皆有 | 先處理 local proxy handler hook，再處理 pinning |

### 何時該開始 Handler Hook

| 情境 | 建議做法 | 判斷信號 |
|------|---------|---------|
| pcap 顯示 `127.0.0.1:<port>` | 開始 handler hook flow | loopback 連線 |
| logcat 有 ProxyServer 日誌 | 找出 handler 類別與方法 | ProxyServer/Netty 日誌 |
| OkHttp/WebView hooks 無輸出但 app 有網路活動 | 先確認 routing 再 hook | Java hooks 無輸出 vs pcap 有流量 |
| 已知 local proxy 但 handler 不明 | 反編譯搜尋 handler 類別 | 靜態分析結果 |

## 不建議的做法

- 不要因為 Wi-Fi MITM 看不到流量就直接宣告 pinning
- 不要對已知 local proxy 的 app 花時間在 OkHttp/WebView hook
- 不要 broad hook raw sockets；先找出 handler 層

## 相關 atoms

- `intelligence/engineering/analysis/signals/local-proxy-detection.md`
- `analysis/apk/workflows/local-proxy-hook-flow.md`
- `analysis/apk/traffic-triage.md`

## Token 影響

低。此 atom 在遇到疑似 local proxy 情境時 lazy-load，約 150-200 tokens。
