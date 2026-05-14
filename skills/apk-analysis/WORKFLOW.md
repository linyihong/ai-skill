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

### 0.1 Reset baseline / 起始狀態

當目標是「從 App 開始到某個具名功能的完整 API 流程」時，開始 capture 前先決定並記錄起始狀態，不要讓 warm cache 或已登入畫面變成隱藏前提：

| Reset level | 用途 | 注意事項 |
| --- | --- | --- |
| `force-stop only` | 保留帳號/session，只重新啟動 App 與網路 client。 | 適合避免 login/rate limit，同時觀察冷啟動與導航流程。 |
| `clear cache` | 減少資源/cache 對列表或詳情的干擾。 | 不一定清掉 DB/session；需記錄是否仍有本機資料。 |
| `clear app data` | 還原 first-run / session recovery / onboarding 狀態。 | 可能移除測試 session、觸發登入或限流；需使用授權測試帳號並記錄邊界。 |
| `reinstall` | 驗證安裝後首輪 bootstrap / permission / migration。 | 成本最高；不要在不需要 first-run 行為時使用。 |

每個 reset-to-feature capture 應拆成可反查的 window：

1. Reset / state preparation：`force-stop`、可選 data/cache clear、權限、proxy、Frida/MITM/pcap 狀態。
2. Cold start / bootstrap：launch、公告/ onboarding / login / session recovery、startup/background APIs。
3. Navigation：從 launcher/home 到目標 page/tab/module 的 tap/swipe/input。
4. Feature operations：列表、分類/filter、搜尋、分頁/scroll、詳情、評論/媒體/action 等。
5. Documentation closure：page map、operation map、API list、schema/correlation、feature handoff、unknowns。

Tab / category / filter strip coverage rule：當頁面有 top tabs、category chips、search result tabs、carousel-like tabs 或任何看起來可水平滑動的分類列時，要先記錄 first viewport 的可見項、hierarchy 暴露的總數（例如「第 X 个标签，共 N 个」）、左右滑動後新增項，以及每個 reachable tab 的 `captured` / `needs capture` / `no-network-observed` 狀態。只測可見 tab 不能宣稱完整 tab 面 API 覆蓋。

Post-selection lazy-load rule：選中 tab、category、filter、grid label 或 chip 後若短窗口沒有新增 feature API，不要直接判定該 UI action 沒有 API。先在已驗證 target package / feature context 下補一個低風險後續 gesture（例如列表 scroll、refresh 或 bounded wait），並把結果分成 `selection-only`、`post-selection-triggered`、`no-network-after-follow-up`。這類 replay knob 或 trigger pattern 一旦可重用，必須立即走 feedback lesson 檢查，不等專案收尾。

UI evidence package validation rule：每個 screenshot / hierarchy 用於 UI-to-API 對齊前，必須驗證 foreground package / activity 屬於目標 App。若 XML package 變成 launcher、browser、Google/search、settings、permission page 或其他外部 App，該 window 要標 `external` / `invalid for target UI`，automation 應中止或記錄明確轉場；Frida 仍命中目標 PID 只能證明目標進程內事件，不能自動證明是該 UI step 觸發。對重要 feature checkpoint，package 正確後還要驗證目標 feature context（例如穩定 tab label、page title、section heading、selected tab 或 route anchor）；同 package 但跑到充值、活動、WebView 或其它 module 的 window 要標 `wrong in-app screen` / `invalid for target feature`，不可當作該 feature evidence。

Checkpoint replay runner rule：同一 feature/page 需要反覆測 Frida、media、tab sweep 或 reset baseline 時，將已確認路徑固化成 replay script，並為 `launch`、目標 tab、列表、詳情、媒體區等節點提供 `--target` / checkpoint 停點。每個 checkpoint 都應截圖、dump XML、驗證 target package；如果跑歪，先修 selector、fallback coordinate、wait 或 scroll，再把後續 capture 當證據。

Post-reset window split rule：`clear app data` / reinstall 後同時需要 session recovery 與 feature API attribution 時，優先拆成「reset + startup/session recovery」與「已驗證導航後 attach feature hooks」兩個 capture window。若 Frida-from-launch 的長窗口導致外部 App、錯頁、公告、更新、WebView 或 timing drift，不要把 feature 操作硬接在同一窗口；先用 package / feature-context guard 證明 session recovery 成功，再從目標 feature checkpoint attach 低負載 hook 取得 feature API 證據。

文件中要把 API 標為 `startup/preload`、`session-recovery`、`navigation`、`feature-triggered`、`cache-hydration` 或 `background/ambiguous`，避免把啟動期或預載 request 誤判成當前點擊觸發。

Read-only argument override rule：若要驗證分頁、排序、語言、filter 這類 read-only 行為，但完整 UI 很難自然觸發下一頁或邊界值，優先選擇高語意函式參數覆寫的短窗口，並保留 App 自己的 session、signing、gateway 與 decrypt path。覆寫腳本必須預設短窗、輸出 schema/hash/key set，不輸出 raw token、raw signature、raw service 或 raw response value；文件要標 `app-owned signing/decrypt preserved`，不能把結果當成 standalone replay parity。

Redacted sample-targeting classifier rule：若 UI 盲抽樣本低收益，但 decrypted response 裡有可判斷樣本可用性的欄位（例如 count、availability、status、type），可新增 disabled-by-default classifier，只輸出 value class 與 item index（例如 `zero/nonzero/missing/other`），再用 UI replay 點擊候選項。不要輸出 raw id、title、body、comment、user、URL、token、完整 count 或其他內容值；文件要標明此 classifier 只是 sample targeting aid，不是資料擷取或 standalone replay parity。

Articles-first live adapter smoke rule：當 APK 分析輸出要接 SDK/private adapter 並驗證真實 read-only 資料時，先選一條核心 read route 做最小 smoke（通常是 list/page 1），只把 base endpoint、該 route binding、opaque/session provider、identity readiness、signing、decrypt/plaintext boundary 設為必填。分類、詳情、留言、媒體、next-page 等 secondary routes 應是 optional follow-up，除非當前目標明確是 full route parity；不要讓 secondary binding 缺失阻塞第一條 live proof。

## 1. 先判斷流量在哪一層

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

## 2. 選擇主線

通用決策原則見 [`../../enforcement/decision-efficiency.md`](../../enforcement/decision-efficiency.md)：先界定當前未知，依 time-to-evidence、語意距離、安全/可回退性、驗證信號、context cost 與 user value 選下一步；不要一次讀完所有分類或把低收益路線硬跑到底。

先做「最高收益路線」檢查：預設工作流是 routing aid，不是固定順序。每個 checkpoint 先用一句話寫出當前未知，再比較可用路線的 **time-to-evidence**、證據語意距離、安全性、可回退性與驗證信號。若已有 App-owned boundary、read-only API replay、高語意 hook 或靜態 xref 能更直接回答問題，主力應轉向該路線；UI / MITM / pcap / broad hooks 等較慢路線可保留作 attribution、對照或 fallback，不要因為一開始選了某種分法就硬跑到底。

| 證據 | 優先主線 |
| --- | --- |
| WebView / OkHttp / HttpURLConnection 命中 | Java hook + MITM + request/response logging |
| 已對 **`OkHttpClient.newCall`**／**`Request$Builder.url`**／**`RealCall.enqueue`** 廣覆蓋（含延遲重試），**使用者操作下仍無業務 host／path** | **勿**再假設「只有時間不夠」；升級為 **Flutter／Dart `dart:io`**、**native `connect`／pcap SNI**，或 **MITM（僅在流量進代理時有效）** |
| Flutter / Dart AOT native path | blutter / reFlutter 類工具 + Frida Dart object hook |
| Native C/C++ custom client | native symbol/string/disassembly + connect/send/recv 輔助 hook |
| Cronet / QUIC | Cronet hooks、flags、HTTP/2/QUIC telemetry、必要時停用 QUIC |
| MITM CONNECT 成功但 handshake failed | CA trust、network security config、pinning 排查 |
| 只有 pcap host/timing | 反編譯搜尋 host/path/header，再找高語意 hook 點 |

## 3. 找高語意 hook 點

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

## 4. Flutter / Dart AOT 常見流程

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

## 5. MITM / Proxy 判斷流程

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

## 6. Response 解碼與離線化

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

## 7. Session / Token 重新取得

遇到 token 過期、no token、invalid token，不要先假設有標準 refresh-token。應還原 App 的真實流程：

1. response interceptor 如何處理錯誤碼。
2. token 存在哪裡，何時清空。
3. device identity 來源是什麼。
4. login / device login body 如何組。
5. request signing 的 canonical path 是否正確。
6. 成功後 token 如何寫回。

若這條流程要支撐 SDK、client、automation 或 live integration，還必須補一張 identity material self-generation audit：逐一判斷 device/install/account/session seed/vendor attestation/server-issued session 相關 key group 是否 `sdk-generatable`，若能，寫 sanitized generation recipe 與驗證矩陣；若不能，標明 `caller-provided`、`server-issued`、`trusted-bridge`、`private-adapter-required` 或 `unknown`，並把缺口列為 live development blocker 或 private adapter scope。

如果 live 測試需要登入：

- 同一輪測試優先重用 session/context。
- 記錄每次 login attempt 的時間、device identity、User-Agent/device 是否同步。
- 遇到 login too frequently，先停止 tight-loop，再分析 server-side bucket 可能維度。
- 不要在沒有證據時假設旋轉單一欄位可以解限流。

## 8. 媒體 / HLS 分析

影片與音訊資源要分控制面與資料面：

| 層 | 例子 | 文件要記錄 |
| --- | --- | --- |
| 詳情 API | 回 title、cover、source path | API path、必要 auth、source field |
| playlist | HLS `.m3u8` | key URI、segment count、duration、base URL |
| key | AES key endpoint 或 key file | key 長度、取得條件、是否需要 auth |
| segments | `.ts` / chunk / signed URL | segment URL 是否短效、query 意義、下載順序 |
| final media | mp4/mp3/image/webp/gif | 解密、解碼、remux、`ffprobe`/header 驗證 |

不要只看副檔名判斷格式。應用 magic bytes、container probe 或 frame count 驗證。例如 WebP 動圖、靜態 GIF、animated GIF 都要分清楚。

HLS `#EXT-X-KEY` URI 不一定就是最終 segment 解密 key。若用 playlist key bytes + IV 解 `.ts` 後沒有 MPEG-TS sync byte 或 container probe 失敗，先回查控制 API / model 是否有 `encrypted_key`、`decrypt_key`、`customKey`、`videoId` 等 key material，並 hook App 端 unwrap helper（例如 base64/AES helper、`getDecryptionKey`、player wrapper）取得私有 key 樣本；最終仍以 segment 解密後 `0x47` sync / `ffprobe` 驗證為準。

## 9. 分析結束定義

一次分析可以收斂時，應具備：

- 清楚知道核心流量走哪個 stack。
- 若使用者要求完整 app-start-to-feature 流程，已記錄 reset/cache/session baseline，並把 reset、startup/session-recovery、navigation、feature API windows 分開回填。
- 有 request metadata 或已證明拿不到的原因。
- 有 response outer shape。
- 若有加密，有解碼點或下一步定位計畫。
- 有去敏樣本或 fixture。
- 有文件回填位置。
- **若下游要程式化取數／接 SDK／寫 integration：** 專案內已有或可指向的 **Domain／執行環境基線**（`DOCUMENTATION.md` § Domain／執行環境基線）：環境維度、session 對列表參數的依存、opaque 參數由來、分頁地面真相、簽章／gateway 前置（均去敏），並連回 API Catalog 條目。缺則視為收口不完整，只可標 `baseline: skeleton` + 明示 open 項。
- **若要開始開發 live-facing SDK/client/app tool：** 上述 baseline 必須通過 development readiness gate，至少回答 endpoint/path family、route/service mapping 或 adapter strategy、session/bootstrap、opaque 參數來源與時效、signing/gateway、decrypt/unwrap、pagination、error/session recovery、replay checklist。缺口未解或未 scoped out 時，不得宣稱分析已足以開發；只能做離線 parser、fixture、mock 或文件補齊。
- 有新的 reusable lesson，或使用者 / reviewer 提出可泛化改進時，已在 **`feedback_history/`** 新增對應檔案；若尚未驗證，標 `candidate` / `experimental` 並寫 validation criteria（規則見 `enforcement/feedback-lessons.md`）。
