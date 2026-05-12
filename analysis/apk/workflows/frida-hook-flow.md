# Frida Hook Flow（Frida Hook 操作流程）

`analysis/apk/workflows/frida-hook-flow.md` 是從 `skills/apk-analysis/techniques/flutter-dart-aot/` 拆解出的 **HOW TO DO** 操作流程。決策智慧（何時該用哪個步驟）請見 `intelligence/engineering/apk-analysis/`。

> **Intelligence Extracted**
> See:
> - `intelligence/engineering/apk-analysis/heuristics/hook-selection.md`
> - `intelligence/engineering/apk-analysis/anti-patterns/early-hook-instability.md`
> - `intelligence/engineering/apk-analysis/failure/frida-spawn-race.md`
> - `intelligence/engineering/apk-analysis/signals/flutter-dart-aot-detection.md`

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

## 成功產出格式

```text
request hook:
  method / baseUrl / path / headers / query

response decode hook:
  decrypted JSON/string
```

## 注意事項

- 避免 broad hooks on global Dart runtime helpers（如 `LinkedHashMap._set`）
- 除非使用 short filtered observation window，否則不要 broad hook
- Dart AOT callsite `BL` addresses 是 navigation hints，不是 function hook entry points
