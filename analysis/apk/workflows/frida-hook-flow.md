# Frida Hook Flow（Frida Hook 操作流程）

`analysis/apk/workflows/frida-hook-flow.md` 是從 `skills/apk-analysis/techniques/flutter-dart-aot/`（已刪除）拆解出的 **HOW TO DO** 操作流程。決策智慧（何時該用哪個步驟）請見 `intelligence/engineering/analytical-reasoning/`。

> **Intelligence Extracted**
> See:
> - `intelligence/engineering/analytical-reasoning/heuristics/hook-selection.md`
> - `intelligence/engineering/analytical-reasoning/anti-patterns/early-hook-instability.md`
> - `intelligence/engineering/analytical-reasoning/failure/frida-spawn-race.md`
> - `intelligence/engineering/analytical-reasoning/signals/flutter-dart-aot-detection.md`

## 前置準備

### 工具安裝

```bash
# blutter（Dart AOT 偽源碼生成）
pip install blutter  # 或從源碼編譯

# unflutter（靜態解析器，blutter 失敗時的備用）
pip install unflutter

# Frida
pip install frida-tools
```

### APK 解包與確認

```bash
# 解包 APK
apktool d target.apk -o target_unpacked

# 確認 Flutter/Dart AOT 存在
ls -la target_unpacked/lib/<arch>/libapp.so
ls -la target_unpacked/lib/<arch>/libflutter.so
```

## 步驟 1：生成偽源碼與函數偏移

```bash
# 使用 blutter（首選）
blutter libapp.so output_dir

# 如果 blutter 識別到 snapshot 但 crash，保留失敗證據
# 切換到 unflutter
unflutter libapp.so --output output_dir
```

產出：
- 偽源碼（pseudo source）
- object pool
- function offsets
- string refs
- call edges

## 步驟 2：搜尋關鍵字

在生成的偽源碼中搜尋：

```bash
# 搜尋 host/base URL
grep -r "https\?://" output_dir/

# 搜尋 Dio/HttpClient
grep -r "Dio\|HttpClient\|RequestOptions\|Interceptor" output_dir/

# 搜尋加密相關
grep -r "encrypt\|decrypt\|AES\|base64\|hash" output_dir/

# 搜尋 response interceptor
grep -r "intercept\|response\|decoder\|parser" output_dir/
```

## 步驟 3：Hook Request Options

```javascript
// Frida script: hook_request_options.js
// 使用 libapp.so base + PC 附加到少量 app-owned function

var libapp = Module.findBaseAddress("libapp.so");
var requestFunc = libapp.add(<PC_OFFSET>);

Interceptor.attach(requestFunc, {
    onEnter: function(args) {
        console.log("Request:", {
            method: args[0].readCString(),
            baseUrl: args[1].readCString(),
            path: args[2].readCString(),
            headers: args[3].readCString(),
            query: args[4].readCString()
        });
    },
    onLeave: function(retval) {
        console.log("Response:", retval.readCString());
    }
});
```

## 步驟 4：Hook Response Decode/Decrypt

在 response decode/decrypt 函數的 return value 上 hook，不要嘗試重建 TLS/socket bytes。

```javascript
// Frida script: hook_response_decode.js
var decodeFunc = libapp.add(<DECODE_PC_OFFSET>);

Interceptor.attach(decodeFunc, {
    onLeave: function(retval) {
        var decoded = retval.readCString();
        console.log("Decoded response:", decoded);
    }
});
```

## 步驟 5：Dart String Decoding（如有需要）

如果 Dart string decoding 失敗：

1. 收集有限的 private hexdump 來推測 layout
2. 修正 decoder
3. 關閉 noisy dumps

```bash
# 收集 hexdump（僅限必要範圍）
dd if=libapp.so bs=1 skip=<OFFSET> count=<LENGTH> | xxd
```

## 步驟 6：對齊與去敏

將 raw wrapper 與 decrypted payload 對齊成 sanitized fixture。

### Image / Resource Decrypt Fixture Sampling

當目標是圖片、封面、avatar 或其他 resource decrypt，而單一指定樣本已被 memory/disk/image cache 吃掉時，不要把指定樣本當作唯一證據來源。先把它當定位點：確認 UI route、loader branch、path family、hook offset 或 cache key 形狀；真正建立 decrypt fixture 時，放寬 hook 條件到同一資源族群，例如同一 path family 或同一 loader branch。

建議 capture 形狀：

1. 在 loader / cache boundary 記錄 path hash、event type、byte length 與去敏 path family。
2. 針對 decrypt function 同時抓 encrypted input、key / params、decrypt output；若函式參數是 wrapper object，先用 AOT 確認 field layout，不要直接把參數當 bytes。
3. 停用 noisy memory-cache dump，除非需要 final payload parity；優先短窗口 attach，讓人工或已驗證 checkpoint 操作列表。
4. 離線 probe 先用任一乾淨同族樣本驗證算法，再回到原指定樣本做 parity。
5. 成功條件是 offline output 可由 magic bytes / container length 驗證，而不是只看到 UI 顯示圖片。

## 成功產出格式

```text
request hook:
  method / baseUrl / path / headers / query

response decode hook:
  decrypted JSON/string

resource decrypt hook:
  path hash / event / encrypted input / key or params / decrypted bytes / final magic
```

## 嵌入式 H5（`flutter_inappwebview`）

當業務主體在 WebView、原生只有殼層 API 時，**不要**用本流程步驟 3–4 的全量 Dio/解密 hook 來找「列表 API」。改走：

1. blutter 搜 `*h5_page*`、`jumpToNext`、domain 鍵 `*H5`；靜態還原 query（常見：`tt` token、`uu`/`un` 使用者、`au` API base、`aId` 渠道）。
2. Frida **minimal**：`Java.perform` hook `android.webkit.WebView.loadUrl`（filter `http` / `tt=`）。
3. 可選：單點 hook blutter 定位的「H5 page URL assign」offset（`onEnter` 讀 `x0` 字串一次），**勿** hook 全域 `Uri::parse`。
4. attach **已啟動** PID；錄製 60–120s 內由使用者手動進入 H5。
5. H5 內 XHR：第二階段對 H5 host 做 MITM 或 `WebViewClient.shouldInterceptRequest`。

詳見 [`traffic-triage.md`](../traffic-triage.md) §混合功能、[`feedback/.../hybrid-native-shell-plus-embedded-h5-frida-routing.md`](../../feedback/history/apk-analysis/common/2026-05-19_101500-hybrid-native-shell-plus-embedded-h5-frida-routing.md)。

## 注意事項

- 避免 broad hooks on global Dart runtime helpers（如 `LinkedHashMap._set`）
- 除非使用 short filtered observation window，否則不要 broad hook
- Dart AOT callsite `BL` addresses 是 navigation hints，不是 function hook entry points
- 嵌入式 H5：**禁止**啟動期 `Uri::parse` object walk；優先 Java WebView（見上節）
