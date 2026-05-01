# APK 分析工具與必要條件

本文件列出常用工具、適用時機與失敗判讀。工具名稱是通用建議，不要求每次全部使用。

## 基礎工具

| 工具 | 用途 | 何時使用 |
| --- | --- | --- |
| `adb` | 安裝、啟動、logcat、pull/push、查 PID、操作裝置。 | 每次分析都需要。 |
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
| blutter | 分析 Flutter AOT `libapp.so`，產生 pseudo source、object pool、offset、Frida decoder。 | 對 Flutter release APK 很有用；若可識別 Dart 版本但 SIGSEGV，保留失敗證據後改用替代 parser。 |
| unflutter | 不嵌入 Dart VM 的 Flutter/Dart AOT static parser，輸出 function map、call edges、string refs、metadata。 | 適合 `blutter` crash 或需要快速取得 function PC 以做 Frida native offset hook。 |
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
| Proxyman 沒有核心 API | client 不走系統代理、attach 太晚、流程沒觸發。 | pcap 確認 host；用 cold-start injection 或高語意 hook。 |
| PC 代理正在監聽但完全沒流量 | 裝置未設 proxy、`adb reverse`/port forward 未建立、Wi-Fi/global proxy 狀態與預期不符。 | 先查裝置 proxy 狀態與 reverse，再冷啟動做短窗驗證。 |
| MITM 有校時／三方流量但沒有業務 host | 只有部分 stack 尊重系統代理；業務可能走 Dart/native/local proxy/TUN 類路由。 | 同窗跑 native `getaddrinfo`/`connect` 或 pcap/SNI；分開記錄「proxy 可用」與「業務是否進 proxy」。 |
| 有 CONNECT 但 SSL handshake failed | CA 不被信任、Android user CA 不生效、custom trust、pinning。 | 先 pass-through 保 App 可用；再處理 CA/system trust/pinning。 |
| Java hook 沒命中 | 流量不在 Java HTTP stack。 | native connect trace；查 Flutter/Cronet/native client。 |
| 只看到 `127.0.0.1:<port>` loopback，沒有上游 API path | App 內建 local ProxyServer/Netty handler 先接本機請求，再由 handler 選上游。 | 反射/Frida 探測 `ProxyServerHandler` 方法；優先 hook `FullHttpRequest` + `URI` 類參數，只記去敏 route metadata。 |
| Netty `FullHttpRequest` hook 到了但 method/path 是空或讀取失敗 | Frida 參數未直接暴露 interface 方法，或 Netty 版本有 `method/getMethod`、`uri/getUri` 差異。 | `Java.cast` 到 `io.netty.handler.codec.http.HttpRequest` / `FullHttpRequest` 後再讀 method/URI；query 值預設去敏。 |
| Netty request method/path 可見，但 headers/body accessor 讀不到 | Frida wrapper/interface dispatch 失敗，或 headers/content 不易直接取出。 | 對實際 Java request 物件呼叫 `toString()`；raw 僅放私有 capture，用去敏摘要判斷 header 名稱與 `content-length`。 |
| Netty/local proxy 已看到加密 header，但 Java plugin/helper hook 未命中生成點 | Flutter Java plugin 可能只是橋接／設定層，實際 signing 在 Dart AOT。 | 抽 `libapp.so` 搜 `Interceptor`、`Dio`、header 名、`sign metadata/result`、`package:<app>/...dart` 與 `_generate...@<hash>`；再用 blutter/AOT xref/offset hook。 |
| Frida 只有 banner 沒輸出 | hook 未命中、script 沒載入、sandbox/權限、attach 時機錯。 | 最小 hook 測試；spawn；降低 hook 數量。 |
| App 卡住或 ANR | hook 太低層、輸出太多、代理 TLS 卡住。 | 限制輸出、pass-through、改高語意 hook。 |
| 解密結果亂碼 | key/IV/KDF/padding/壓縮順序錯。 | hook decrypt return value 建對照 fixture。 |
| Token 重新簽了仍失敗 | token 本身失效，不是簽章問題。 | 還原 App 的 login/device-login/session refresh 流程。 |
| Login too frequently | 短時間 tight-loop、device/session/IP/app fingerprint 風控。 | 停止重試、重用 session、記錄 login attempt metadata。 |
| HLS 只保存 m3u8 仍不能播放 | 缺 key、segments、base URL、AES 解密或 remux。 | 分開抓 playlist/key/segments，最後用 `ffprobe` / `ffmpeg` 驗證。 |
| `apkanalyzer` 報無法定位 build tools；或 PATH 沒有 `aapt` | 只有 platform-tools 在 PATH，或 SDK 配置不完整。 | 改用 `$ANDROID_HOME/build-tools/<version>/aapt dump badging`（見下）。 |
| `aapt dump badging` 沒有 `launchable-activity` | manifest 較複雜或工具輸出差異。 | 裝置已安裝該 App 時用 `cmd package resolve-activity --brief`；僅有 APK 時用 `aapt dump xmltree` / apktool。 |
| Wi‑Fi 代理 MITM 幾乎沒有業務流量，但全機 pcap 仍有 TLS／App 功能正常 | **內建 TUN／sing-box／embedded VPN** 等可能繞過系統 HTTP 代理。 | 字串搜 `singbox`/`MethodChannel`；改看 **pcap SNI**／**Frida hook 高語意 client**；不要先歸因 pinning。 |
| Wi‑Fi MITM 空，但 **logcat** 有 **`ProxyServer`**／對 **`127.0.0.1:<port>`** 轉發至 `https://<api-host>` | **本機 loopback 中介**先於對外連線。 | `adb logcat` 搜尋 `ProxyServer`／handler；勿將含標頭的原始 log 提交公開 repo。 |
| `blutter` 能偵測 Dart version/snapshot，但 full 或 `--no-analysis` SIGSEGV | Dart VM introspection 路線對該 snapshot/tool 版本不穩。 | 改用 `unflutter` 等 static parser 產生 `functions.jsonl`/`call_edges.jsonl`/`string_refs.jsonl`；再 hook 少量高語意 Dart function PC。 |
| Dart AOT offset hook 命中，但 Dart String 解碼全是空或亂碼 | String layout 假設錯，尤其 Dart 3.x compressed pointer OneByteString 可能使用 raw byte length + inline bytes。 | 私有 capture 中限量 hexdump 物件，驗證 length/data offset；常見候選包含 untagged `+0x08` raw length、`+0x10` data；修好 decoder 後關閉 hexdump。 |

## 命令模板

查 APK metadata：

```bash
aapt dump badging app.apk
```

若 shell 找不到 `aapt`，但已安裝 Android SDK（環境變數 `ANDROID_HOME` 或標準 SDK 根目錄存在 `build-tools/<version>/aapt`）：

```bash
# 將 <version> 換成已存在的 build-tools 目錄（通常選最新版）
"$ANDROID_HOME/build-tools/<version>/aapt" dump badging app.apk
```

在裝置上查預設 launcher component（便於 `am start -n`，比 `monkey` 易重現）：

```bash
adb -s <device-serial> shell cmd package resolve-activity --brief <package-name>
```

查裝置與 App PID：

```bash
adb devices
adb -s <device-serial> shell pidof <package-name>
```

抓全機 pcap：

```bash
adb -s <device-serial> shell su -c 'tcpdump -i any -s 0 -w /sdcard/app.pcap'
adb -s <device-serial> pull /sdcard/app.pcap .
```

查 proxy 是否殘留：

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

## Frida 健康檢查

完整 hook 沒有輸出時，先用最小步驟確認工具鏈，不要直接判斷 hook 點錯。

1. 確認裝置與 Frida server：

```bash
frida-ps -D <device-serial>
```

2. 確認目標 PID：

```bash
adb -s <device-serial> shell pidof <package-name>
```

3. 最小 attach：

```bash
frida -D <device-serial> -p <pid> -e 'console.log("HOOK_LOADED")'
```

4. 最小 spawn：

```bash
frida -D <device-serial> -f <package-name> -e 'console.log("SPAWN_LOADED")'
```

只有最小 hook 可用後，才逐步加入完整 Java/native/Dart hook。

## 媒體驗證工具

檢查圖片 magic bytes：

```bash
xxd -l 16 output.bin
```

檢查影片 / GIF / WebP 容器：

```bash
ffprobe -v error -show_entries format=format_name,duration -of default=noprint_wrappers=1 output.mp4
ffprobe -v error -select_streams v:0 -count_frames -show_entries stream=codec_name,nb_read_frames,nb_frames,duration output.gif
```

HLS 下載後若要提供給一般使用者播放，通常要把 segments 依 playlist/key/IV 解密合併，再 remux 成常見容器，例如 MP4。只保存 m3u8 不等於完成播放檔驗證。

## 工具選擇原則

- 能用 request object hook，就不要先重組 socket bytes。
- 能用 response decoder hook，就不要只靠 MITM 解 outer response。
- 能離線驗證，就不要每次都依賴動態 hook。
- 任何工具失敗都要留下「失敗證據」，因為排除路徑也是重要結論。
