> 遵守 [共用規則索引](../../../../enforcement/README.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/frida-hook-flow.md`](../../../../analysis/apk/workflows/frida-hook-flow.md)

### 2026-05-01 - 對 `OkHttpClient.newCall`／`Request$Builder.url`／`enqueue` 廣覆蓋仍無業務 host 時

Status: validated

#### One-line Summary

Flutter 類 App 在 **`libapp.so`** 主導 HTTP 時，僅 hook **Java `okhttp3`** 可能長時間 **得不到 `https://api…` 業務 URL**；需在授權下切 **Dart／native TLS／pcap SNI** 或 **MITM（若進代理）**。

#### Human Explanation

許多分析流程預設「業務走 OkHttp」。當已同時 hook **`newCall`**、**`Builder.url`**、**`RealCall.enqueue/execute`**（並延遲重試）且在 **使用者操作列表／滑動** 後仍 **零筆**業務域—代表 **假設錯誤**或 **業務在另一 client**，而不是「hook 時間不夠」單一因素。

#### Trigger

- **`dump_okhttp_headers`／`RealCall` chain** 長期只有校時、DoH、loopback、設定鏈。
- 已升级到 **`dump_okhttp_full_url`** 仍無 **`service=`**／**`/api/public/`**／**業務根域名**。
- 靜態已有 **Flutter AOT**、`dart:io`、`MethodChannel` 證據。

#### Evidence

- Tool: Frida，`okhttp3.OkHttpClient.newCall`、`okhttp3.Request$Builder.url`、`enqueue`／`execute`。
- Sanitized excerpt: 側錄見 **DoH**、**`127.0.0.1:<port>/check`**、**資源 GET**，但 **無** **`<api-business-host>`** 明文 URL。
- Evidence path: `<PROJECT_ROOT>/capture/frida_okhttp_full_url*.log`（勿提交）。

#### Generalized Lesson

**廣覆蓋 Java OkHttp** 可作為 **否定假設** 實驗：若結果為空，應 **升高優先級**到 **Dart HTTP**、**Cronet**（若存在）、**native `connect`**／**pcap SNI**，而非無限加長 OkHttp 掛載時間。

#### Agent Action

1. 先跑 **`dump_okhttp_full_url.js`**（含手動 UX），保存 log。
2. 若仍無業務 host → 記錄為 **「Java OkHttp 非業務主路徑」**，改開 **`min_native_dns_connect.js`** 或 **pcap**，並在專案筆記更新 **§管道對照**。
3. 勿向使用者保證「再加長秒數就一定能側錄到 REST」。

#### Applies When

- APK 含 **`libapp.so`**／Flutter embedding，且文档已懷疑 **Dart 主線**。

#### Does Not Apply When

- 業務明確走 **Java Retrofit/OkHttp**（靜態／動態已有 **業務 URL** 命中）。

#### Validation

同一份操作路徑下 **parallel**：OkHttp log 空 + **`getaddrinfo`／SNI** 出現業務域 → 方向一致。

#### Promotion Target

- `WORKFLOW.md`（決策樹：OkHttp 否證後分支）
- `DOCUMENTATION.md`
