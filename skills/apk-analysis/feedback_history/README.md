# feedback_history 索引

下列為條目檔名與摘要；完整內容請開啟對應 `.md`。檔名格式：**`YYYY-MM-DD_HHMMSS-<slug>.md`**（`HHMMSS` 為 24 小時制本機時間）。

| 檔案 | Status | 標題 | 一句話摘要 |
| --- | --- | --- | --- |
| `2026-04-30_120000-proxy-failure-要先拆成導流與-tls-兩層.md` | promoted | Proxy failure 要先拆成導流與 TLS 兩層 | 代理看不到明文時，先確認「有沒有進代理」，再談憑證或 pinning。 |
| `2026-04-30_120001-冷啟動比-attach-更適合抓初始化與代理設定.md` | promoted | 冷啟動比 attach 更適合抓初始化與代理設定 | 要影響 App 的網路 client 初始化，通常要用冷啟動 `spawn`，不要等 App 跑起來後才 attach。 |
| `2026-04-30_120002-高語意-hook-優先於-socket-bytes.md` | promoted | 高語意 hook 優先於 socket bytes | 能 hook request/response 物件，就不要先從 socket bytes 開始拼。 |
| `2026-04-30_120003-動態-hook-結果要離線化.md` | promoted | 動態 hook 結果要離線化 | Frida 是拿樣本的手段，不應成為長期測試唯一依賴。 |
| `2026-04-30_120004-frida-只有-banner-時先做最小-hook-健康檢查.md` | validated | Frida 只有 banner 時先做最小 hook 健康檢查 | Frida 沒輸出不一定是 hook 點錯，可能是 client、權限、sandbox、attach 時機或 App 沒觸發流程。 |
| `2026-04-30_120005-session-refresh-要還原-app-的真實登入-裝置流程.md` | validated | Session refresh 要還原 App 的真實登入/裝置流程 | Token 過期不一定有 refresh-token；要看 App 實際怎麼重新取得 session。 |
| `2026-04-30_120006-登入限流要避免-tight-loop-優先-session-reuse.md` | validated | 登入限流要避免 tight-loop，優先 session reuse | 遇到 login too frequently，不要盲目旋轉單一參數；先重用 session 並記錄風控維度。 |
| `2026-04-30_120007-媒體播放鏈要分控制面-金鑰與資料面.md` | validated | 媒體播放鏈要分控制面、金鑰與資料面 | 影片不是只有一個 URL；HLS 需要分開記錄 playlist、key、segment、解密與合併驗證。 |
| `2026-04-30_120008-aapt-sdk-build-tools-resolve-activity.md` | validated | APK metadata：`aapt` 不在 PATH 時走 SDK build-tools；launcher 用 `resolve-activity` | `apkanalyzer` 或環境找不到 build-tools 時，改用 `$ANDROID_HOME/build-tools/<version>/aapt`；badging 若沒有 launcher 行，用 `cmd package resolve-activity` 取得 … |
| `2026-04-30_120009-內建-sing-box-tun-類通道可能繞過-wi-fi-系統代理.md` | candidate | 內建 sing-box／TUN 類通道可能繞過 Wi‑Fi 系統代理 | 當 `libapp.so`／`classes.dex` 出現 **sing-box**、**MethodChannel `.../singbox`** 或 embedded VPN 類命名時，MITM 若「完全收不到業務 host」，不一定是 App 沒連線，而可能是流量**不經… |
| `2026-04-30_120010-本機-loopback-proxyserver-轉發會讓-wi-fi-http-mitm-看不到業務-connect.md` | candidate | 本機 loopback「ProxyServer」轉發會讓 Wi‑Fi HTTP MITM 看不到業務 CONNECT | 若流量先到 **`127.0.0.1:<local-port>`** 的 HTTP 代理，再由該代理連向真實 API host，電腦上的 Wi‑Fi HTTP MITM 通常只看到 loopback，業務 **CONNECT** 可能完全不進電腦代理。 |
| `2026-05-01_101500-doh-dns-query-param-side-channel-from-okhttp-log.md` | validated | DoH 的 `dns=` 參數可作為 MITM 業務 host 空白時的側信道 | 業務 HTTPS 未進 MITM 時，若 OkHttp 仍發 DoH GET，可從 **`dns=<base64url>`** 離線解出問句 **QNAME**（A/AAAA），確認 App 仍解析哪些 API 根。 |
| `2026-05-01_142000-exhaustive-java-okhttp-hooks-may-still-miss-flutter-business-http.md` | validated | 廣覆蓋 Java OkHttp 仍無業務 host 時應轉 Dart／native／pcap | 同時 hook **`newCall`／`Builder.url`／`enqueue`** 並延遲重試後，若使用者操作下仍無業務 URL，視為 **Java OkHttp 非主路徑**；改 **Dart／native TLS／pcap SNI** 或 **MITM**。 |
