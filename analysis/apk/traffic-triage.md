# APK 流量與執行路徑分流

本文件定義如何判斷 APK 的網路流量走哪一層。這是所有 APK 分析的第一步，在選擇具體工具或 hook 策略之前。

## 核心原則

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

## 主線選擇

通用決策原則：先界定當前未知，依 time-to-evidence、語意距離、安全/可回退性、驗證信號、context cost 與 user value 選下一步。

| 證據 | 優先主線 |
| --- | --- |
| WebView / OkHttp / HttpURLConnection 命中 | Java hook + MITM + request/response logging |
| 已對 OkHttpClient.newCall／Request$Builder.url／RealCall.enqueue 廣覆蓋，使用者操作下仍無業務 host/path | 升級為 Flutter/Dart dart:io、native connect／pcap SNI，或 MITM（僅在流量進代理時有效） |
| Flutter / Dart AOT native path | blutter / reFlutter 類工具 + Frida Dart object hook |
| Native C/C++ custom client | native symbol/string/disassembly + connect/send/recv 輔助 hook |
| Cronet / QUIC | Cronet hooks、flags、HTTP/2/QUIC telemetry、必要時停用 QUIC |
| MITM CONNECT 成功但 handshake failed | CA trust、network security config、pinning 排查 |
| 只有 pcap host/timing | 反編譯搜尋 host/path/header，再找高語意 hook 點 |

## 高語意 Hook 點選擇

優先找已組裝完成的業務物件：

- request options：method、base URL、path、query、headers、body。
- response wrapper：status、outer JSON、raw bytes。
- response decoder / interceptor：解密後 inner JSON。
- token/session provider：刷新流程、device identity、header provider。
- local proxy handler：若 App 內有 ProxyServer / Netty / loopback server，優先 hook handler 的 FullHttpRequest + resolved URI。

避免一開始就 hook：

- `send` / `recv` 全量 bytes。
- TLS read/write 全量事件。
- 太多 native function offset。

低層 hook 事件多、容易造成卡頓，也需要自己重組協議。只有在高語意點找不到時才降層。

## Flutter / Dart AOT 判斷流程

如果 evidence 指向 Flutter：

1. 解 APK，確認 `lib/<arch>/libapp.so` 與 `libflutter.so`。
2. 用 Dart AOT 分析工具產生 pseudo source、object pool、function offsets。
3. 搜尋：host / base URL、Dio、HttpClient、RequestOptions、Interceptor、encrypt、decrypt、AES、base64、hash。
4. 用 Frida hook request options。
5. 用 Frida hook response decode/decrypt return value。
6. 把 raw wrapper + decrypted payload 對齊成 fixture。

避免把全域 Dart runtime/collection helper 當第一個 hook 點；這些 helper 高頻且噪音大。

## MITM / Proxy 判斷流程

先分兩層：

```text
流量有沒有進 proxy？
  否 -> client 沒走代理、注入太晚、proxy host/port 錯，或流程未重建 client
  是 -> 再看 TLS trust / CA / pinning
```

不要把「PC 端代理正在監聽」當成裝置已導流。先記錄裝置 proxy 狀態，再看是否有任何流量進代理；最後才驗證業務 host 是否也進同一條 proxy。

### 冷啟動代理導流流程

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
- proxy 有 CONNECT 但無 HTTPS 明文：導流與 TLS 要分開查。

## Response 解碼流程

遇到 outer response 包 encrypted `data` 時：

1. 找 response interceptor / decoder / decrypt function。
2. 記錄輸入格式：base64、prefix/salt、ciphertext、version field。
3. 記錄演算法：KDF、AES mode、padding、MAC、compression。
4. 用 hook 取得 decrypted output。
5. 寫離線 decoder。
6. 建立 raw encrypted -> decrypted fixture。
7. 用 fixture 驗證 SDK/client mapping。

離線化完成後，後續不應每次依賴 Frida 才能跑測試。

## Session / Token 重新取得

遇到 token 過期、no token、invalid token，不要先假設有標準 refresh-token。應還原 App 的真實流程：

1. response interceptor 如何處理錯誤碼。
2. token 存在哪裡，何時清空。
3. device identity 來源是什麼。
4. login / device login body 如何組。
5. request signing 的 canonical path 是否正確。
6. 成功後 token 如何寫回。

## 媒體 / HLS 分析

影片與音訊資源要分控制面與資料面：

| 層 | 例子 | 文件要記錄 |
| --- | --- | --- |
| 詳情 API | 回 title、cover、source path | API path、必要 auth、source field |
| playlist | HLS `.m3u8` | key URI、segment count、duration、base URL |
| key | AES key endpoint 或 key file | key 長度、取得條件、是否需要 auth |
| segments | `.ts` / chunk / signed URL | segment URL 是否短效、query 意義、下載順序 |
| final media | mp4/mp3/image/webp/gif | 解密、解碼、remux、ffprobe/header 驗證 |

不要只看副檔名判斷格式。應用 magic bytes、container probe 或 frame count 驗證。

---

← [回到 analysis/apk/](README.md)
