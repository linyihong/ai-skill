# APK Analysis Tools & Failure Interpretation

本文件列出 APK 分析常用工具、適用時機與常見失敗判讀。工具名稱是通用建議，不要求每次全部使用。

## 基礎工具

| 工具 | 用途 | 何時使用 |
| --- | --- | --- |
| `adb` | 安裝、啟動、logcat、pull/push、查 PID、操作裝置。 | 每次分析都需要。 |
| `adb screencap` / `screenrecord` | 保存去敏 screenshot 或短影片，建立 UI/tab/操作證據。 | 需要把 API 對應回具體畫面與操作時。 |
| `uiautomator dump` | 匯出目前畫面的 view hierarchy、text/resource-id/content-desc。 | screenshot label 不清楚、需要確認 tab/screen 元素時。 |
| `apktool` | 解 resources、manifest、network security config、smali。 | 需要看 manifest、res、smali 或重打包時。 |
| `jadx` | Java/Kotlin 反編譯、搜尋 class/method/string。 | Java stack、OkHttp、WebView、加密 helper。 |
| `aapt` / `apkanalyzer` | APK metadata、package、version、permissions。 | 開始前盤點。 |
| `sqlite3` | 檢查本機 cache DB。 | App 有 SQLite cache、媒體 cache、離線資料時。 |
| `strings` / `rg` | 快速找 host、path、header、secret-like string。 | 靜態初篩。 |

## 動態觀測

| 工具 | 用途 | 注意 |
| --- | --- | --- |
| `tcpdump` | 裝置端 pcap，觀察 localhost、DNS、SNI、IP、port、timing。 | 需要 root 或可用抓包環境；HTTPS 內容仍是密文。 |
| Wireshark | 讀 pcap、看 SNI、TCP stream、時序。 | 用於證據，不是直接解密工具。 |
| Proxyman / Charles / Burp / mitmproxy | MITM、CONNECT、HTTP(S) 代理、HLS/media 可見流量。 | 先確認流量是否進代理，再處理 CA/pinning。 |
| Frida | Java/native/Dart hook，抓 request object、decode function、connect trace。 | 優先 hook 高語意點；注意 crash、ANR、輸出量。 |
| objection | 快速探索 Android runtime、常見 SSL pinning 檢查。 | 適合輔助，不取代定制 hook。 |

## Flutter / Dart AOT

| 工具 | 用途 | 注意 |
| --- | --- | --- |
| blutter | 分析 Flutter AOT `libapp.so`，產生 pseudo source、object pool、offset、Frida decoder。 | 對 Flutter release APK 很有用；若 SIGSEGV，保留失敗證據後改用替代 parser。 |
| unflutter | 不嵌入 Dart VM 的 Flutter/Dart AOT static parser，輸出 function map、call edges、string refs、metadata。 | 適合 blutter crash 或需要快速取得 function PC。 |
| reFlutter 類工具 | 改 Flutter engine / dump dart traffic 的路線。 | 侵入性較高，需評估是否符合授權範圍。 |
| IDA / Ghidra / radare2 | native disassembly、function offset、xref。 | 用於補足自動工具看不到的邏輯。 |

## 解密與資料處理

| 工具 / 套件 | 用途 |
| --- | --- |
| Python `cryptography` | AES-CBC/GCM、PKCS7、HMAC、hash KDF。 |
| Node.js `crypto` | 快速重放 signing / decrypt prototype。 |
| `jq` | JSON 摘要、schema 初步整理。 |
| `xxd` / `hexdump` | 檢查 magic bytes、prefix、binary wrapper。 |
| `ffprobe` / `ffmpeg` | 驗證媒體容器、HLS remux、frame count。 |

## 常見失敗判讀

| 現象 | 可能原因 | 下一步 |
| --- | --- | --- |
| UI map / 截圖流程讓 App 或分析環境變卡 | 截圖、錄影、UI dump、自動遍歷、hook logging 同時進行，I/O 與主執行緒壓力太高。 | 降成 lightweight overview：只截主要 tabs/關鍵 screen；先停錄影與批量 dump，保留 API hook/pcap 主線。 |
| 自動化操作抓 API 時結果不穩 | 操作未等待畫面穩定、背景預載/cache 混入、同一腳本含多個 action。 | 每個 operation script 只做一個 flow；輸出開始/結束 timestamp；必要時先 force-stop/冷啟動並加入短等待。 |
| 抓到 API 但不知道是哪個操作觸發 | 沒有建立 UI 操作時間窗，或 startup/preload/background sync 混在一起。 | 先補 screenshot/UI hierarchy 與 operation id；每次只操作一個 screen/action。 |
| 截圖看起來是某個 tab，但 API timing 對不上 | tab 預載、快取、背景同步、或同 endpoint 被多個 screen 共用。 | 標成 trigger confidence low/medium；用冷啟動、清 cache、單步操作或 hook sequence 重新驗證。 |
| Proxyman 沒有核心 API | client 不走系統代理、attach 太晚、流程沒觸發。 | pcap 確認 host；用 cold-start injection 或高語意 hook。 |
| PC 代理正在監聽但完全沒流量 | 裝置未設 proxy、adb reverse/port forward 未建立、Wi-Fi/global proxy 狀態與預期不符。 | 先查裝置 proxy 狀態與 reverse，再冷啟動做短窗驗證。 |
| MITM 有校時／三方流量但沒有業務 host | 只有部分 stack 尊重系統代理；業務可能走 Dart/native/local proxy/TUN 類路由。 | 同窗跑 native getaddrinfo/connect 或 pcap/SNI。 |
| 有 CONNECT 但 SSL handshake failed | CA 不被信任、Android user CA 不生效、custom trust、pinning。 | 先 pass-through 保 App 可用；再處理 CA/system trust/pinning。 |
| Java hook 沒命中 | 流量不在 Java HTTP stack。 | native connect trace；查 Flutter/Cronet/native client。 |
| 只看到 127.0.0.1:\<port\> loopback，沒有上游 API path | App 內建 local ProxyServer/Netty handler 先接本機請求，再由 handler 選上游。 | 反射/Frida 探測 ProxyServerHandler 方法；優先 hook FullHttpRequest + URI 類參數。 |
| Frida 只有 banner 沒輸出 | hook 未命中、script 沒載入、sandbox/權限、attach 時機錯。 | 最小 hook 測試；spawn；降低 hook 數量。 |
| App 卡住或 ANR | hook 太低層、輸出太多、代理 TLS 卡住。 | 限制輸出、pass-through、改高語意 hook。 |
| 解密結果亂碼 | key/IV/KDF/padding/壓縮順序錯。 | hook decrypt return value 建對照 fixture。 |
| Token 重新簽了仍失敗 | token 本身失效，不是簽章問題。 | 還原 App 的 login/device-login/session refresh 流程。 |
| Login too frequently | 短時間 tight-loop、device/session/IP/app fingerprint 風控。 | 停止重試、重用 session、記錄 login attempt metadata。 |
| HLS 只保存 m3u8 仍不能播放 | 缺 key、segments、base URL、AES 解密或 remux。 | 分開抓 playlist/key/segments，最後用 ffprobe/ffmpeg 驗證。 |
| Wi‑Fi 代理 MITM 幾乎沒有業務流量，但全機 pcap 仍有 TLS／App 功能正常 | 內建 TUN／sing-box／embedded VPN 等可能繞過系統 HTTP 代理。 | 字串搜 singbox/MethodChannel；改看 pcap SNI／Frida hook 高語意 client。 |
| Wi‑Fi MITM 空，但 logcat 有 ProxyServer／對 127.0.0.1:\<port\> 轉發至 https://\<api-host\> | 本機 loopback 中介先於對外連線。 | adb logcat 搜尋 ProxyServer／handler。 |
| blutter 能偵測 Dart version/snapshot，但 full 或 --no-analysis SIGSEGV | Dart VM introspection 路線對該 snapshot/tool 版本不穩。 | 改用 unflutter 等 static parser 產生 functions.jsonl/call_edges.jsonl/string_refs.jsonl。 |
| Dart AOT offset hook 命中，但 Dart String 解碼全是空或亂碼 | String layout 假設錯，尤其 Dart 3.x compressed pointer OneByteString。 | 私有 capture 中限量 hexdump 物件，驗證 length/data offset。 |
| Dart AOT call_edges 指向的 caller 內部 BL 位址無法 Frida attach | callsite 不是函式入口；全域 runtime/helper 太熱且噪音大。 | 把 callsite 當靜態導航線索；優先 hook app-owned function PC 或更高語意 Java/native boundary。 |

## 工具選擇原則

- 能用 request object hook，就不要先重組 socket bytes。
- 能用 response decoder hook，就不要只靠 MITM 解 outer response。
- 能離線驗證，就不要每次都依賴動態 hook。
- 任何工具失敗都要留下「失敗證據」，因為排除路徑也是重要結論。

## Frida 健康檢查

完整 hook 沒有輸出時，先用最小步驟確認工具鏈：

1. `frida-ps -D <device-serial>` — 確認裝置與 Frida server。
2. `adb -s <device-serial> shell pidof <package-name>` — 確認目標 PID。
3. `frida -D <device-serial> -p <pid> -e 'console.log("HOOK_LOADED")'` — 最小 attach。
4. `frida -D <device-serial> -f <package-name> -e 'console.log("SPAWN_LOADED")'` — 最小 spawn。

只有最小 hook 可用後，才逐步加入完整 Java/native/Dart hook。

## 命令模板

查 APK metadata：
```bash
aapt dump badging app.apk
```

查預設 launcher component：
```bash
adb -s <device-serial> shell cmd package resolve-activity --brief <package-name>
```

查裝置與 App PID：
```bash
adb devices
adb -s <device-serial> shell pidof <package-name>
```

建立 UI architecture map 的 screenshot 與 hierarchy evidence：
```bash
adb -s <device-serial> shell screencap -p /sdcard/screen.png
adb -s <device-serial> pull /sdcard/screen.png ./evidence/ui/<operation-id>.png
adb -s <device-serial> shell uiautomator dump /sdcard/window.xml
adb -s <device-serial> pull /sdcard/window.xml ./evidence/ui/<operation-id>.xml
```

抓全機 pcap：
```bash
adb -s <device-serial> shell su -c 'tcpdump -i any -s 0 -w /sdcard/app.pcap'
adb -s <device-serial> pull /sdcard/app.pcap .
```

查 proxy 狀態：
```bash
adb -s <device-serial> shell settings get global http_proxy
adb -s <device-serial> shell dumpsys connectivity
```

Frida attach：
```bash
PID="$(adb -s <device-serial> shell pidof <package-name> | tr -d '\r')"
frida -D <device-serial> -p "$PID" -l hook.js
```

Frida spawn：
```bash
adb -s <device-serial> shell am force-stop <package-name>
frida -D <device-serial> -f <package-name> -l hook.js
```

---

← [回到 analysis/apk/](README.md)
