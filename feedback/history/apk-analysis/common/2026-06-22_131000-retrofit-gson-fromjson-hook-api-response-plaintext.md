> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Retrofit Gson.fromJson hook captures API JSON when obfuscated Response peek fails

Status: candidate

#### One-line Summary

OkHttp chain hook 已命中 business URL，但混淆 `Response`（`okhttp3.u0`）無法 peekBody 時，改 hook **`com.google.gson.Gson.fromJson(String, Class|Type)`** 可在 Retrofit `GsonResponseBodyConverter` 路徑拿到 **API JSON 明文**（如 HLS `PlayURL` 陣列）。

#### Human Explanation

`130100` lesson 指出 Response 層可能被 R8 掏空。實務下一層是 Retrofit Gson converter：在 Gson 解析前，第一參數就是 response body 字串。Hook 此處比 hook `GsonResponseBodyConverter.convert(y0)` 更穩（`y0.string()` 會 consume body）。注意同一 hook 也會看到 telemetry/analytics JSON，需用 path 關鍵字或 `/api/` 特徵過濾。

#### Trigger

- REQ 日誌豐富，RSP/`u0.h` = 0
- APK 含 `retrofit2.converter.gson.GsonResponseBodyConverter`
- 需要 `getChapterContent` 等 response schema

#### Evidence

- Tool: Frida hook `Gson.fromJson` overloads on running app
- Sanitized excerpt: 出現 `[{PlayURL, Encode, Dpi, Bitrate, MultiBit, …}]` 結構之 JSON 陣列（具體 URL/host 留 project docs）
- Evidence path: `<PROJECT_ROOT>/scripts/frida/hook_capture_full.js`、`<PROJECT_ROOT>/api/dynamic-w2-playback.md`

#### Generalized Lesson

```text
REQ 命中 + body 為 0？
  1. 確認 retrofit2.converter.gson 存在
  2. Hook Gson.fromJson(String, Class/Type) — log json 前 N 字元
  3. interest 過濾：chapter/book/play_url/m3u8 或業務 path 關鍵字
  4. 仍不足 → jadx 找自訂 Interceptor / 非 Gson converter
```

與 `RealInterceptorChain.proceed(o0)` **並用**，不互斥。

#### Agent Action

1. 增量腳本：在既有 chain hook 上加 Gson 層，不刪 REQ hook。
2. Capture log 可能含 analytics；gitignore + redact uid/device 後再摘 schema 寫 project docs。
3. Ai-skill 只寫 hook 點與過濾策略，不寫 target host 或 sample URL。

#### Goal / Action / Validation

- Goal: 關閉「response body unknown」gap。
- Action: 更新 `130100` 交叉引用；promotion 至 traffic-triage response 小節。
- Validation: Gson hook 後至少 1 條業務 JSON（非純 telemetry）含可辨識 schema 欄位。

#### Applies When

- Retrofit2 + Gson + OkHttp3 Java 主線
- Obfuscated Response 無 peekBody

#### Does Not Apply When

- Moshi/Kotlin serialization only（改 hook 對應 converter）
- Protobuf/encrypted body（JSON hook 無 plaintext）

#### Validation

- `[JSON]` 行含預期 schema 欄位（如 stream variant 陣列）
- 與 REQ 中同 path 的請求成對出現

#### Promotion Target

- `analysis/apk/traffic-triage.md` §response body
- 交叉引用 `common/2026-06-22_130100-r8-obfuscated-okhttp-response-needs-converter-hook.md`

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` common +1
