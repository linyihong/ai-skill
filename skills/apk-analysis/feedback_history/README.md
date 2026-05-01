# feedback_history 索引

下列為條目檔名與摘要；完整內容請開啟對應 `.md`。檔名格式：**`YYYY-MM-DD_HHMMSS-<slug>.md`**（`HHMMSS` 為 24 小時制本機時間）。

## 分類索引

此 skill 已拆分 `techniques/`，因此 feedback 也按分類存放。每個分類 README 負責該類 lesson 的索引；新 lesson 若屬於特定分類，應優先放入 `feedback_history/<category>/`。

| 分類 | 用途 |
| --- | --- |
| [`common/`](common/) | 跨分類、共通路由、工具健康檢查、proxy/session 等通用分析規則。 |
| [`flutter-dart-aot/`](flutter-dart-aot/) | Flutter、Dart AOT、Dart object/string/offset hook 技巧。 |
| [`http-api/`](http-api/) | HTTP API 文件化、request/response schema、UI-to-API binding。 |
| [`local-proxy/`](local-proxy/) | loopback、Netty/local proxy、embedded route/TUN 類技巧。 |
| [`media-hls/`](media-hls/) | HLS、playlist/key/segments、媒體容器與下載驗證。 |

## 歷史平面索引

| 檔案 | Status | 標題 | 一句話摘要 |
| --- | --- | --- | --- |
| `common/2026-04-30_120000-proxy-failure-要先拆成導流與-tls-兩層.md` | promoted | Proxy failure 要先拆成導流與 TLS 兩層 | 代理看不到明文時，先確認「有沒有進代理」，再談憑證或 pinning。 |
| `common/2026-04-30_120001-冷啟動比-attach-更適合抓初始化與代理設定.md` | promoted | 冷啟動比 attach 更適合抓初始化與代理設定 | 要影響 App 的網路 client 初始化，通常要用冷啟動 `spawn`，不要等 App 跑起來後才 attach。 |
| `common/2026-04-30_120002-高語意-hook-優先於-socket-bytes.md` | promoted | 高語意 hook 優先於 socket bytes | 能 hook request/response 物件，就不要先從 socket bytes 開始拼。 |
| `common/2026-04-30_120003-動態-hook-結果要離線化.md` | promoted | 動態 hook 結果要離線化 | Frida 是拿樣本的手段，不應成為長期測試唯一依賴。 |
| `common/2026-04-30_120004-frida-只有-banner-時先做最小-hook-健康檢查.md` | validated | Frida 只有 banner 時先做最小 hook 健康檢查 | Frida 沒輸出不一定是 hook 點錯，可能是 client、權限、sandbox、attach 時機或 App 沒觸發流程。 |
| `common/2026-04-30_120005-session-refresh-要還原-app-的真實登入-裝置流程.md` | validated | Session refresh 要還原 App 的真實登入/裝置流程 | Token 過期不一定有 refresh-token；要看 App 實際怎麼重新取得 session。 |
| `common/2026-04-30_120006-登入限流要避免-tight-loop-優先-session-reuse.md` | validated | 登入限流要避免 tight-loop，優先 session reuse | 遇到 login too frequently，不要盲目旋轉單一參數；先重用 session 並記錄風控維度。 |
| `media-hls/2026-04-30_120007-媒體播放鏈要分控制面-金鑰與資料面.md` | validated | 媒體播放鏈要分控制面、金鑰與資料面 | 影片不是只有一個 URL；HLS 需要分開記錄 playlist、key、segment、解密與合併驗證。 |
| `common/2026-04-30_120008-aapt-sdk-build-tools-resolve-activity.md` | validated | APK metadata：`aapt` 不在 PATH 時走 SDK build-tools；launcher 用 `resolve-activity` | `apkanalyzer` 或環境找不到 build-tools 時，改用 `$ANDROID_HOME/build-tools/<version>/aapt`；badging 若沒有 launcher 行，用 `cmd package resolve-activity` 取得 … |
| `local-proxy/2026-04-30_120009-內建-sing-box-tun-類通道可能繞過-wi-fi-系統代理.md` | candidate | 內建 sing-box／TUN 類通道可能繞過 Wi‑Fi 系統代理 | 當 `libapp.so`／`classes.dex` 出現 **sing-box**、**MethodChannel `.../singbox`** 或 embedded VPN 類命名時，MITM 若「完全收不到業務 host」，不一定是 App 沒連線，而可能是流量**不經… |
| `local-proxy/2026-04-30_120010-本機-loopback-proxyserver-轉發會讓-wi-fi-http-mitm-看不到業務-connect.md` | candidate | 本機 loopback「ProxyServer」轉發會讓 Wi‑Fi HTTP MITM 看不到業務 CONNECT | 若流量先到 **`127.0.0.1:<local-port>`** 的 HTTP 代理，再由該代理連向真實 API host，電腦上的 Wi‑Fi HTTP MITM 通常只看到 loopback，業務 **CONNECT** 可能完全不進電腦代理。 |
| `common/2026-05-01_101500-doh-dns-query-param-side-channel-from-okhttp-log.md` | validated | DoH 的 `dns=` 參數可作為 MITM 業務 host 空白時的側信道 | 業務 HTTPS 未進 MITM 時，若 OkHttp 仍發 DoH GET，可從 **`dns=<base64url>`** 離線解出問句 **QNAME**（A/AAAA），確認 App 仍解析哪些 API 根。 |
| `common/2026-05-01_112900-proxy-config-vs-business-route.md` | validated | Proxy config is not business route proof | PC 代理正在監聽不代表裝置已導流；即使部分流量進 proxy，也不能推出核心業務主線一定進 proxy。 |
| `local-proxy/2026-05-01_114300-local-proxy-handler-uri-hook.md` | validated | Hook local proxy handler URI, not just OkHttp | App 內建 loopback ProxyServer 時，hook handler 的 `FullHttpRequest` + `URI` 參數可直接確認上游業務 host。 |
| `local-proxy/2026-05-01_131000-cast-netty-request-for-handler-route.md` | validated | Cast Netty request interfaces for handler routes | Frida hook 到 Netty request 參數時，先 cast 到 `HttpRequest` / `FullHttpRequest` 再讀 method/URI，可取得去敏 method/path。 |
| `local-proxy/2026-05-01_132400-netty-aggregated-request-tostring-headers.md` | validated | Netty aggregated request toString can expose headers | Netty request accessor 讀不到 headers/body 時，Java `toString()` 仍可暴露 request line、headers 與 `content-length` 結構。 |
| `flutter-dart-aot/2026-05-01_133900-dart-aot-interceptor-strings-after-java-helper-miss.md` | validated | Dart AOT interceptor strings after Java helper miss | Netty/local proxy 已看到加密 header，但 Java plugin/helper hook 未命中生成點時，轉向 `libapp.so` 的 Dart AOT interceptor 字串與函式名線索。 |
| `flutter-dart-aot/2026-05-01_140900-unflutter-aot-offset-hook-after-blutter-crash.md` | promoted | Unflutter AOT offset hook after blutter crash | `blutter` 識別 Dart snapshot 但 SIGSEGV 時，改用 `unflutter` 取 function PC/string refs，再以 Frida native offset hook 驗證 request/sign/decrypt 流程。 |
| `flutter-dart-aot/2026-05-01_142000-exhaustive-java-okhttp-hooks-may-still-miss-flutter-business-http.md` | validated | 廣覆蓋 Java OkHttp 仍無業務 host 時應轉 Dart／native／pcap | 同時 hook **`newCall`／`Builder.url`／`enqueue`** 並延遲重試後，若使用者操作下仍無業務 URL，視為 **Java OkHttp 非主路徑**；改 **Dart／native TLS／pcap SNI** 或 **MITM**。 |
| `flutter-dart-aot/2026-05-01_145948-dart-aot-callsite-hooks-are-not-function-hooks.md` | validated | Dart AOT Callsite Hooks Are Not Function Hooks | Dart AOT call edge 裡的 caller 內部 `BL` 位址只適合導航；Frida 動態驗證應優先 hook function entry 或高語意邊界。 |
| `flutter-dart-aot/2026-05-01_151551-schema-only-jsondecode-hook.md` | validated | Schema-only jsonDecode Hook | 回包明文只在 Dart decode 階段出現時，hook `jsonDecode` input 並只輸出 length/hash/top-level keys/types，避免落 raw JSON value。 |
| `flutter-dart-aot/2026-05-01_153437-sequence-jsondecode-before-api-response.md` | validated | Sequence jsonDecode Before Calling It API Response | `jsonDecode` hook 看到的 schema 不一定是 API response；必須用 request/decrypt sequence 證明它發生在業務 response 路徑上。 |
| `flutter-dart-aot/2026-05-01_155200-dart-compressed-response-fields.md` | candidate | Dart AOT compressed response fields | 讀 Dart/Dio response handoff object field 時，要依 tagged object field load 還原 compressed pointer，避免誤判為 `0x800...` 假指標。 |
| `flutter-dart-aot/2026-05-01_164741-dart-inline-onebyte-string-smi-length.md` | validated | Dart inline one-byte string Smi length | Dart AOT inline one-byte string 可能用 Smi length；raw32 讀法會把有效 JSON/path/base64 字串誤判為 undecoded。 |
| `http-api/2026-05-01_171500-json-shape-before-query-shape.md` | validated | JSON Shape Before Query Shape | Schema-only response logging should classify JSON before query extraction and gate query parsing to request-like strings, otherwise embedded URLs or binary pipe payloads can be mistaken for request query evidence. |
| `http-api/2026-05-01_171650-ui-architecture-map-from-screenshots.md` | promoted | UI architecture map from screenshots | 在可控制裝置時，用輕量 screenshot/UI hierarchy 建立 App 架構地圖；若截圖造成卡頓，可先解 API 再對關鍵流程補 UI binding。 |
| `http-api/2026-05-01_173800-api-field-documentation-after-analysis.md` | promoted | API field documentation after analysis | 分析完 API 後要文件化 HTTP headers、request、response，並逐字段說明 type、meaning、required/source；截圖只作為 UI 觸發輔助證據。 |
| `common/2026-05-01_174100-technique-category-routing.md` | promoted | Technique category routing | APK analysis techniques should be split by runtime/API family so an A-type analysis does not need to load B-type guidance. |
| `common/2026-05-01_175500-workflow-as-routing-not-technique-dump.md` | promoted | Workflow as routing, not technique dump | Top-level APK workflow should stay as common routing; runtime/API-family details belong in `techniques/<category>/`. |
