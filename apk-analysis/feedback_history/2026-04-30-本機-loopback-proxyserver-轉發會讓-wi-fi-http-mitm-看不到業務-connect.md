### 2026-04-30 - 本機 loopback「ProxyServer」轉發會讓 Wi‑Fi HTTP MITM 看不到業務 CONNECT

Status: candidate

#### One-line Summary

若流量先到 **`127.0.0.1:<local-port>`** 的 HTTP 代理，再由該代理連向真實 API host，電腦上的 Wi‑Fi HTTP MITM 通常只看到 loopback，業務 **CONNECT** 可能完全不進電腦代理。

#### Human Explanation

有些 App（含 Flutter + OkHttp 組合）會在裝置上開 **本機 proxy**，Dart／Java client 先對 loopback 發請求，再由服務端元件轉發 HTTPS。此時「系統 Wi‑Fi 手動代理」指向電腦時，**不一定**能複製到這段路徑：對 OS 而言出站連線可能仍是 **直連** API IP／走分流／經其他路由。

#### Trigger

- MITM／Proxyman 幾乎沒有業務流量。
- `adb logcat` 出現 **`ProxyServer`**／**`ProxyServerHandler`**、`Listening on port ... for target https://<api-host>`。

#### Evidence

- Tool：`adb logcat`（tag／line 依裝置而有差異）。
- Sanitized excerpt：存在「本地埠 → `https://<api-host>`」之類描述；MITM 側無對應 CONNECT。

#### Generalized Lesson

輔助抓網域時可同步搜 log：**proxy／handler／forward／localhost**。Wi‑Fi MITM 失敗時不要只猜 pinning；先確認是否存在 **loopback 中介**。

#### Agent Action

當使用者 MITM 為空但業務明顯有 HTTPS：建議 **logcat grep `ProxyServer`**／相關 tag；並提示原始 log 可能含敏感標頭，必須去敏後才寫入 repo。

#### Applies When

- App 內建本地 proxy／middleware／sing-box 類鏈路。
- OkHttp／Dart HttpClient 連 `127.0.0.1`。

#### Does Not Apply When

- 已確認業務 HTTPS **CONNECT** 目標為電腦 `<proxy-host>`。

#### Validation

- logcat 宣稱的 `<api-host>` 與 root **pcap SNI**／Frida hook URL 一致。

#### Promotion Target

- `TOOLS.md`（常見失敗判讀）
