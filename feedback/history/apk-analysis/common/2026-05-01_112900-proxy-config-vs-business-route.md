> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

# Extracted — See [`workflow/apk-analysis/execution-flow.md`](../../../../workflow/apk-analysis/execution-flow.md) (Section 1: 開始前確認, proxy triage)

### 2026-05-01 - Proxy Config Is Not Business Route Proof

Status: validated

#### One-line Summary

PC 端代理正在監聽不代表裝置已導流；即使部分流量進 proxy，也不能推出核心業務主線一定進 proxy。

#### Human Explanation

APK 側錄時常見誤判是：看到 mitmproxy、Burp 或 Proxyman 已啟動，就假設裝置流量會進來；或看到校時、統計、第三方 SDK 已進代理，就假設所有業務 HTTPS 也會進代理。實務上要拆成三層驗證：

1. 裝置是否真的設了 proxy／reverse。
2. 是否有任何 proxy-aware 流量進 PC MITM。
3. 核心業務 host 是否也在同一時間窗進 PC MITM。

如果第 2 層成立但第 3 層不成立，同時 native `getaddrinfo`／`connect` 可見業務候選 host 解析與直連外部 IP，優先判讀為核心路由繞過 PC MITM，而不是「沒有業務流量」或直接跳到 certificate pinning。

#### Trigger

- PC 端代理工具正在監聽，但 MITM 日誌沒有業務 host。
- 裝置上 App 功能正常，Frida 或 pcap 仍可見業務候選 DNS／TCP。
- MITM 只看到校時、統計、第三方 SDK、媒體或系統流量。

#### Evidence

- Tool: `adb settings get global http_proxy`, `adb reverse --list`, Frida native `getaddrinfo`/`connect`, PC MITM log.
- Sanitized excerpt:
  - Before setup: `global http_proxy=:0`, `adb reverse --list` empty.
  - After setup: PC MITM receives proxy-aware third-party traffic.
  - Same window: native hook shows business-candidate hostnames followed by `connect` to external IPs, not to `<proxy-host>:<proxy-port>`.
- Evidence path: project-local `capture/` logs and sanitized project notes; do not copy target hosts or private paths into reusable skill docs.

#### Generalized Lesson

Treat proxy capture as a route-validation problem, not a binary "MITM works/doesn't work" problem. Validate proxy configuration first, then prove whether the specific business route enters that proxy. A mixed result is common: some stacks obey Android/global proxy while Flutter/Dart/native/local-proxy/TUN-like paths bypass it.

#### Agent Action

When MITM lacks business traffic:

1. Check device proxy state with `settings get global http_proxy` and `adb reverse --list`.
2. If needed, run a controlled short window with known proxy settings and cold-start the app.
3. In the same window, compare PC MITM with native `getaddrinfo`/`connect` or pcap/SNI.
4. Report separately:
   - proxy not configured,
   - proxy configured and third-party traffic enters,
   - business route enters proxy,
   - business route bypasses proxy.

Do not conclude pinning until the business host actually reaches the proxy and TLS/handshake or certificate validation fails.

#### Applies When

- Android APK dynamic traffic analysis with PC-side MITM.
- App may use Flutter/Dart, native sockets, local ProxyServer, embedded VPN/TUN, DoH, or custom DNS.
- Some non-business traffic appears in MITM but target business APIs are absent.

#### Does Not Apply When

- The target business host is already visible as a CONNECT/server-connect in the proxy and only TLS decryption fails.
- The device is intentionally using a full-device VPN capture that is known to see all routes.
- The analysis goal is only to confirm third-party SDK traffic, not business API capture.

#### Validation

Validated when one controlled time window shows:

- proxy state was explicitly recorded before and after,
- proxy-aware traffic reached PC MITM,
- native DNS/connect or pcap showed business-candidate traffic in the same window,
- those business destinations did not appear as proxy CONNECT/server-connect targets.

#### Promotion Target

- `WORKFLOW.md`
- `TOOLS.md`
