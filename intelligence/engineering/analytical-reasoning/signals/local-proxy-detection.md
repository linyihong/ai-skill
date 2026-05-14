# Local Proxy Detection Signals（本機代理偵測信號）

## 問題

如何判斷一個 APK 是否使用 local proxy/loopback 路由？需要哪些技術信號來確認分析路線？

## 判斷信號

### 主要信號（高可信度）

| 信號 | 檢查方式 | 可信度 |
|------|---------|-------|
| pcap 顯示 `127.0.0.1:<port>` 或 loopback HTTP | `tcpdump -r capture.pcap -X 'dst host 127.0.0.1'` | 高 |
| Logcat 出現 `ProxyServer` 或 Netty handler 日誌 | `adb logcat -d \| grep -iE "proxyserver\|netty"` | 高 |
| Logcat 出現 `local bridge`、`TUN`、`sing-box` | `adb logcat -d \| grep -iE "local bridge\|tun\|sing-box"` | 高 |

### 次要信號（中等可信度）

| 信號 | 檢查方式 | 可信度 |
|------|---------|-------|
| Wi-Fi MITM 看不到 business CONNECT，但 app 正常運作 | MITM proxy 無流量 vs app 有網路活動 | 中 |
| OkHttp/WebView hooks 不暴露 upstream | Frida Java hook 無 business 輸出 | 中 |
| 反編譯後發現 Netty、ProxyServer 相關類別 | `grep -r "ProxyServer\|Netty" jadx_output/` | 中 |
| pcap 顯示外部 TLS 但 MITM 無法攔截 | pcap 有 TLS handshake vs MITM 無對應流量 | 中 |

### 排除信號

| 信號 | 意義 |
|------|------|
| MITM proxy 出現 SSL handshake error | 可能是 TLS pinning，不是 local proxy |
| OkHttp hooks 可攔截所有流量 | 非 local proxy，流量直接走系統代理 |
| pcap 無 loopback 連線 | 不太可能是 local proxy |

## 判斷流程

```
Wi-Fi MITM 看不到 business CONNECT？
    ├── 是 → app 仍正常運作？
    │       ├── 是 → 檢查 pcap loopback 證據
    │       │       ├── 127.0.0.1:<port> 存在 → Local Proxy 路線
    │       │       └── 無 loopback → 檢查 TLS pinning
    │       └── 否 → MITM 連線失敗 → TLS Pinning 路線
    └── 否 → OkHttp/WebView hooks 有輸出 → 標準 HTTP 分析路線
```

## 相關 atoms

- `intelligence/engineering/analytical-reasoning/heuristics/local-proxy-routing-diagnosis.md`
- `analysis/apk/workflows/local-proxy-hook-flow.md`
- `analysis/apk/traffic-triage.md`

## Token 影響

低。此 atom 在分析初期 lazy-load，約 100-150 tokens。
