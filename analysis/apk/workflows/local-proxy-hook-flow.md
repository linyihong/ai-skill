# Local Proxy / Loopback Hook Flow（本機代理 Hook 操作流程）

`analysis/apk/workflows/local-proxy-hook-flow.md` 是從 `skills/apk-analysis/techniques/local-proxy/`（已刪除）拆解出的 **HOW TO DO** 操作流程。決策智慧（何時該懷疑 local proxy、如何區分導流 vs TLS 問題）請見 `intelligence/engineering/analysis/heuristics/local-proxy-routing-diagnosis.md`。

> **Intelligence Extracted**
> See:
> - `intelligence/engineering/analysis/heuristics/local-proxy-routing-diagnosis.md`
> - `intelligence/engineering/analysis/signals/local-proxy-detection.md`

## 前置準備

### 必要條件

- 目標 app 已安裝並可執行
- 具備 pcap/logcat/MITM 觀測能力
- Frida 可用於 hook Java 層

### 工具

```bash
# 流量觀測
tcpdump -i any port 443 -w capture.pcap
adb logcat | grep -E "ProxyServer|Netty|127\.0\.0\.1"

# MITM proxy
mitmproxy -p 8080
```

## 步驟 1：確認 Loopback 證據

在懷疑 local proxy 存在時，先用以下方式確認：

```bash
# pcap 檢查 loopback 流量
tcpdump -r capture.pcap -X 'dst host 127.0.0.1'

# logcat 檢查 ProxyServer/Netty handler 日誌
adb logcat -d | grep -iE "proxyserver|netty|local bridge|tun|sing-box"

# MITM timing 檢查：Wi-Fi MITM 看不到 business CONNECT
# 但 app 仍正常運作，pcap 顯示外部 TLS
```

## 步驟 2：識別 Handler 方法與 Request Object Type

不要直接加 broad socket hooks。先找出 handler 類別和方法：

```bash
# 使用 jadx 或類似工具反編譯 APK
# 搜尋 ProxyServer、Netty handler、local bridge 相關類別
# 找出 handler 方法的參數類型（如 FullHttpRequest、HttpRequest）
```

## 步驟 3：Hook Handler Argument

Hook handler 方法中代表 local request 的參數：

```javascript
// Frida script: hook_local_proxy.js
var ProxyHandler = Java.use("com.example.ProxyHandler");

ProxyHandler.handleRequest.implementation = function(req) {
    console.log("Local request:", {
        method: req.method(),
        uri: req.uri(),
        headers: req.headers().toString()
    });
    return this.handleRequest(req);
};
```

## 步驟 4：Cast Netty Request（如適用）

如果 `FullHttpRequest` 無法直接在 Frida 中讀取 method/path：

```javascript
// Cast to io.netty.handler.codec.http.HttpRequest or FullHttpRequest
var HttpRequest = Java.use("io.netty.handler.codec.http.HttpRequest");
var fullReq = Java.cast(req, HttpRequest);

console.log("Method:", fullReq.getMethod());
console.log("URI:", fullReq.getUri());
```

## 步驟 5：去敏與摘要

- 預設 redact query values
- 保留 path shape、method、header names
- 如果 `headers()` 或 `content()` 失敗，只在 private capture 中使用 `toString()`
- Public docs 應摘要 structure 和 `content-length`，不保留 raw sensitive values

## 步驟 6：歸因到 Upstream API

使用 local proxy metadata 將 local UI/client calls 對應到 upstream APIs：

```markdown
| Local Request | Upstream API | Evidence |
|---------------|-------------|----------|
| GET /local/feed | GET https://api.example.com/v1/feed | handler 轉發日誌 |
| POST /local/login | POST https://api.example.com/v1/auth/login | handler 轉發日誌 + response match |
```

如果 upstream request 是由 Flutter/Dart 在 local handler 之後建構的，請先確認證據指向 Flutter/Dart AOT 再切換路線。

## 成功產出格式

```text
local proxy evidence:
  pcap: 127.0.0.1:8080 loopback detected
  logcat: ProxyServer started on port 8080

handler hook:
  class: com.example.ProxyHandler
  method: handleRequest(FullHttpRequest)
  local request: GET /local/feed
  upstream: GET https://api.example.com/v1/feed
```

## 注意事項

- 不要將缺少 Wi-Fi MITM traffic 直接視為 pinning，先確認 routing
- 區分 local client-to-loopback requests 與 loopback-to-upstream requests
- Hook high-semantic handler objects 優先於 raw sockets
