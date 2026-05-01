# APK 分析工作流

本工作流用於授權 APK 分析。它的目的不是一次就猜中答案，而是用證據逐層排除錯誤方向。

## 0. 開始前確認

記錄：

- APK 名稱與版本。
- package name。
- CPU 架構：`arm64-v8a`、`armeabi-v7a` 等。
- 裝置或 emulator。
- 是否 root。
- 是否允許 Frida / MITM / pcap / decompile。
- 是否有測試帳號與授權範圍。

禁止：

- 保存完整 token、cookie、真實 device id、私密 key。
- 把第三方或非目標流量當成分析資料。
- 把單一 App 的私有 host / secret 寫進 reusable skill。

## 1. 建立 UI 架構地圖

當可以控制裝置或 emulator 時，先把 App 的可見結構文件化，讓後續 API 分析能回到具體操作，而不是只留 endpoint 清單。

記錄：

- 入口與登入狀態：冷啟動、已登入、未登入、權限彈窗、地區/語言。
- 可見 navigation：bottom tabs、top tabs、drawer、profile/menu、search、detail page、player/media page。
- 每個主要 screen 的 screenshot evidence：去敏後保存路徑、時間戳、screen label。
- 操作序列：`Home > Tab: Discover > item tap > Detail` 這類可重放路徑。
- 捕獲時間窗：操作開始/結束時間、對應 pcap/MITM/Frida log window。
- API 關聯：每個操作觸發的 method/path、response schema、cache/local-only 判斷。

判斷原則：

- 截圖只能證明可見 UI，不直接證明 API；必須用同窗 request/response、pcap timing 或 hook sequence 對齊。
- 同一個 endpoint 可能被多個 screen 共用；文件要保留多個 UI path，不要硬塞單一來源。
- 啟動畫面、快取、預載與背景同步要標清楚，避免把 startup/cache JSON 誤寫成使用者操作 API。
- 若可用 `uiautomator` hierarchy、accessibility labels 或 route/logcat 線索，應與 screenshot 互相校正 tab 名稱與 screen label。

## 2. 先判斷流量在哪一層

不要一開始就假設是 certificate pinning。先回答這些問題：

```text
App 是否有 localhost bridge？
  -> 抓 lo / loopback

裝置是否真的對外連線？
  -> 抓 whole-device pcap，看 DNS/SNI/IP/port/timing

流量是否進系統代理或 MITM？
  -> 看代理工具是否收到 CONNECT / request
  -> 若沒有 CONNECT，先不要談 CA / pinning

Java HTTP stack 是否命中？
  -> hook OkHttp / HttpURLConnection / WebView 類

native backtrace 落在哪裡？
  -> libflutter / libapp / Cronet / custom native / Java
```

判斷原則：

- pcap 看得到 host / SNI / timing，但看不到 HTTPS 內文。
- MITM 看不到核心流量，不代表沒有流量；可能是 client 不走代理。
- Java hook 沒命中，不代表沒有 HTTP；Flutter/Dart、Cronet、native client 可能繞過 Java HTTP stack。
- 看到 proxy CONNECT 後 TLS 失敗，才進入 CA / pinning / custom trust 排查。

## 3. 選擇主線

| 證據 | 優先主線 |
| --- | --- |
| WebView / OkHttp / HttpURLConnection 命中 | Java hook + MITM + request/response logging |
| 已對 **`OkHttpClient.newCall`**／**`Request$Builder.url`**／**`RealCall.enqueue`** 廣覆蓋（含延遲重試），**使用者操作下仍無業務 host／path** | **勿**再假設「只有時間不夠」；升級為 **Flutter／Dart `dart:io`**、**native `connect`／pcap SNI**，或 **MITM（僅在流量進代理時有效）** |
| Flutter / Dart AOT native path | blutter / reFlutter 類工具 + Frida Dart object hook |
| Native C/C++ custom client | native symbol/string/disassembly + connect/send/recv 輔助 hook |
| Cronet / QUIC | Cronet hooks、flags、HTTP/2/QUIC telemetry、必要時停用 QUIC |
| MITM CONNECT 成功但 handshake failed | CA trust、network security config、pinning 排查 |
| 只有 pcap host/timing | 反編譯搜尋 host/path/header，再找高語意 hook 點 |

## 4. 找高語意 hook 點

優先找已組裝完成的業務物件：

- request options：method、base URL、path、query、headers、body。
- response wrapper：status、outer JSON、raw bytes。
- response decoder / interceptor：解密後 inner JSON。
- token/session provider：刷新流程、device identity、header provider。
- local proxy handler：若 App 內有 `ProxyServer` / Netty / loopback server，優先 hook handler 的 `FullHttpRequest` + resolved `URI`，確認本機請求如何映射到上游 API。
  - 若 `FullHttpRequest` 直接讀不到 method/path，先 cast 到 Netty `HttpRequest` / `FullHttpRequest` interface，再讀 `method/getMethod`、`uri/getUri`，並預設去敏 query。
  - 若 method/path 已可見但 `headers()` / `content()` 讀不到，可用實際 Java request 物件的 `toString()` 做 fallback；raw output 可能含完整 header，僅放私有 capture，文件只寫去敏結構。

避免一開始就 hook：

- `send` / `recv` 全量 bytes。
- TLS read/write 全量事件。
- 太多 native function offset。

低層 hook 事件多、容易造成卡頓，也需要自己重組協議。只有在高語意點找不到時才降層。

## 5. Flutter / Dart AOT 常見流程

如果 evidence 指向 Flutter：

1. 解 APK，確認 `lib/<arch>/libapp.so` 與 `libflutter.so`。
2. 用 Dart AOT 分析工具產生 pseudo source、object pool、function offsets；若 `blutter` 識別 snapshot 後 SIGSEGV，改用 `unflutter` 等 static parser 先拿 function map／call edges／string refs。
3. 搜尋：
   - host / base URL。
   - `Dio`、`HttpClient`、`RequestOptions`、`Interceptor`。
   - `encrypt`、`decrypt`、`AES`、`base64`、`hash`、`ResponseInterceptor`。
4. 用 Frida hook request options；若已取得 AOT function PC，可用 `libapp.so` base + PC 對少量 `RequestInterceptor`／sign／encrypt/decrypt 函式做 native offset hook。`call_edges` 裡 caller 內部的 `BL` 位址只當導航線索，不要預設可直接 `Interceptor.attach()`。
5. 用 Frida hook response decode/decrypt return value。
6. 若 Dart String decoder 失敗，先在私有 capture 限量 hexdump 物件推導 layout（例如 OneByteString raw length/data offset），修好後關閉 hexdump。
7. 把 raw wrapper + decrypted payload 對齊成 fixture。

避免把全域 Dart runtime/collection helper（例如 `LinkedHashMap._set`、常見 string helper）當第一個 hook 點；這些 helper 高頻且噪音大，可能讓 App 卡頓或提前結束。若一定要觀察內部 helper，先用短窗、嚴格 filter，或改用 Stalker／更高語意邊界驗證。

若 local proxy/Netty hook 已看到自訂加密／簽名 header，但同窗 Java plugin/helper hook（例如 AES/RC2/getNMKey/query map）沒有命中，應把 Java plugin 視為橋接或設定層，轉向 `libapp.so` Dart AOT interceptor 字串、object pool xref、blutter/offset hook；不要在 Java helper 層無限加 hook。

成功特徵：

```text
request hook:
  method / baseUrl / path / headers / query

response decode hook:
  decrypted JSON/string
```

## 6. MITM / Proxy 判斷流程

先分兩層：

```text
流量有沒有進 proxy？
  否 -> client 沒走代理、注入太晚、proxy host/port 錯，或流程未重建 client
  是 -> 再看 TLS trust / CA / pinning
```

如果需要讓 Flutter/Dart client 進 proxy，常見策略是冷啟動前注入 proxy env 或修改 debug network behavior。驗收不是「代理工具有沒有明文」，而是先看 connect target 是否變成 `<proxy-host>:<proxy-port>`。

不要把「PC 端代理正在監聽」當成裝置已導流。先記錄裝置 proxy 狀態（例如 global proxy、Wi-Fi proxy、reverse/port forward），再看是否有任何流量進代理；最後才驗證**業務 host** 是否也進同一條 proxy。若代理已收到校時／統計／第三方流量，但 native `getaddrinfo`／pcap 同窗顯示業務 host 直連外部 IP，應判讀為**核心業務路由繞過 PC MITM**，不要直接歸因 pinning。

如果代理工具顯示 `SSL Handshake Failed`：

- 導流可能已成功。
- 問題可能在 CA trust、Android user CA、network security config、custom trust 或 pinning。
- 若只要 App 可用與時序證據，可先對核心 host pass-through。

### 冷啟動代理導流流程

當目標是「讓 Proxyman / Burp / mitmproxy 連得到 App 核心流量」時，優先用冷啟動流程，而不是 attach 已啟動 App：

```text
1. 確認 proxy tool 正在監聽 <proxy-host>:<proxy-port>
2. 確認裝置目前 proxy 狀態，避免殘留錯誤 port
3. force-stop App
4. 用 Frida spawn 啟動 App
5. 在網路 client 初始化前注入 proxy env 或 hook proxy selector
6. hook connect 或觀察 proxy CONNECT
7. 驗證 connect target 是否變成 <proxy-host>:<proxy-port>
8. 若已進 proxy 但 TLS failed，再查 CA trust / pinning
```

判讀：

- `connect <proxy-host>:<proxy-port>`：導流成功。
- `connect <api-host>:443`：仍是直連，優先查注入時機、proxy host/port、client 是否已初始化。
- proxy 有 CONNECT 但無 HTTPS 明文：導流與 TLS 要分開查，不要直接回到「Proxyman 沒用」的結論。

## 7. Response 解碼與離線化

遇到 outer response 包 encrypted `data` 時：

1. 找 response interceptor / decoder / decrypt function。
2. 記錄輸入格式：base64、prefix/salt、ciphertext、version field。
3. 記錄演算法：KDF、AES mode、padding、MAC、compression。
4. 用 hook 取得 decrypted output。
5. 若 decrypt return 是 wrapper / Future / Map 而非可直接讀的 String，改 hook `jsonDecode`、`JsonDecoder.convert` 或 app 的 parse/decode helper，先輸出 schema-only 摘要（length/hash/top-level keys/types，不印 values）。
   - 同時記錄 request/decrypt/json 的 sequence/timestamp；`jsonDecode` 若出現在第一個業務 request 前，先標成 local/cache/startup schema，不要寫成 API response。
6. 寫離線 decoder。
7. 建立 raw encrypted -> decrypted fixture。
8. 用 fixture 驗證 SDK/client mapping。

離線化完成後，後續不應每次依賴 Frida 才能跑測試。

## 8. Session / Token 重新取得

遇到 token 過期、no token、invalid token，不要先假設有標準 refresh-token。應還原 App 的真實流程：

1. response interceptor 如何處理錯誤碼。
2. token 存在哪裡，何時清空。
3. device identity 來源是什麼。
4. login / device login body 如何組。
5. request signing 的 canonical path 是否正確。
6. 成功後 token 如何寫回。

如果 live 測試需要登入：

- 同一輪測試優先重用 session/context。
- 記錄每次 login attempt 的時間、device identity、User-Agent/device 是否同步。
- 遇到 login too frequently，先停止 tight-loop，再分析 server-side bucket 可能維度。
- 不要在沒有證據時假設旋轉單一欄位可以解限流。

## 9. 媒體 / HLS 分析

影片與音訊資源要分控制面與資料面：

| 層 | 例子 | 文件要記錄 |
| --- | --- | --- |
| 詳情 API | 回 title、cover、source path | API path、必要 auth、source field |
| playlist | HLS `.m3u8` | key URI、segment count、duration、base URL |
| key | AES key endpoint 或 key file | key 長度、取得條件、是否需要 auth |
| segments | `.ts` / chunk / signed URL | segment URL 是否短效、query 意義、下載順序 |
| final media | mp4/mp3/image/webp/gif | 解密、解碼、remux、`ffprobe`/header 驗證 |

不要只看副檔名判斷格式。應用 magic bytes、container probe 或 frame count 驗證。例如 WebP 動圖、靜態 GIF、animated GIF 都要分清楚。

## 10. 分析結束定義

一次分析可以收斂時，應具備：

- 清楚知道核心流量走哪個 stack。
- 有 UI architecture map，能說明主要 tabs/screens 與已測操作。
- 有 request metadata 或已證明拿不到的原因。
- 有 response outer shape。
- 若有加密，有解碼點或下一步定位計畫。
- 有去敏樣本或 fixture。
- 有文件回填位置。
- 有新的 reusable lesson 時，已在 **`feedback_history/`** 新增對應檔案（規則見 `shared-rules/feedback-lessons.md`）。
