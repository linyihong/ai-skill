> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - OkHttp Interceptor sorted-map canonical string for request signing

Status: candidate

#### One-line Summary

R8 混淆的 OkHttp **Signing Interceptor** 常把 **header 與 POST form 欄位** 合併進同一 `HashMap`，**字母序排序**後拼成 `key1=value1&key2=value2&…`（去尾 `&`），再交給 native/crypto util 產生 `sign` header；RE 應 hook **interceptor 入口 + crypto util 入參**，而非先猜 HMAC 公式。

#### Human Explanation

動態 capture 已看到 `sign` header，但離線重放仍失敗，常因未還原 **canonical string 組裝規則**。此模式在 Retrofit + OkHttp 短劇/REST App 中常見：interceptor 在 `chain.proceed` 前注入 `uid`、`session`、`ts`、`requestTime`、`clientTraceId` 等，並把 `application/x-www-form-urlencoded` body 的 key 一併放入 map。空 body 的 POST 則只有 header 欄位。排序與 `&` 拼接是機械步驟，可先用 Frida 打印 **crypto util 的第一個 String 參數** 與最終 `sign` 對照，再在 host 端實作 canonical builder。

#### Trigger

- DEX：`Interceptor.intercept` 實作類 + `sign` / `requestTime` header literal
- Frida：每次 v-api POST 前出現固定 header 族 + 64 hex `sign`
- 同 session、不同 form body → `sign` 變化；同 body、不同 `requestTime` → `sign` 變化
- Python `hashlib.sha256(canonical)` 與 capture 的 `sign` 不一致

#### Evidence

- Tool: Frida hook obfuscated Interceptor + util `sha256Encrypt(String, int)`（或等效 native crypto）
- Sanitized excerpt: `[SIGN_IN] canonical=apiVersion=…&book_id=…&requestTime=…&session=…` → `[SIGN_OUT] <64 hex>`；POST body keys 出現在 canonical 中
- Evidence path: `<PROJECT_ROOT>/api/signing-re.md`、`<PROJECT_ROOT>/capture/`（gitignore）

#### Generalized Lesson

```text
Signing canonical RE order:
  1. Locate OkHttp Interceptor (implements Interceptor, before chain.proceed)
  2. Hook crypto util first String arg + returned hex sign
  3. Confirm map sources: standard headers + POST form fields (if any)
  4. Confirm sort: alphabetical key order, skip empty values, join with &
  5. Host-side: build_canonical(sorted_map) — still need native/crypto for final sign
  6. Do NOT assume plain SHA256(canonical) until verified
```

#### Agent Action

1. Project：維護 canonical field list 與範例 shape（placeholder 值，無 live session）。
2. Ai-skill：只寫組裝判斷樹；class 名、channelId、apiVersion 留 project。
3. 與 `141700`（SDK gate）、`wire-path-vs-signing-canonical-path` 並用。

#### Goal / Action / Validation

- Goal: 縮短 sign RE 從「猜算法」到「還原 canonical + 定位 crypto」。
- Action: Frida probe 腳本只 hook interceptor + crypto util，避免 class enumeration timeout。
- Validation: ≥2 不同 endpoint 的 canonical 可從 log 重組且欄位規則一致；crypto 仍 native 則標 relay 為 interim。

#### Applies When

- Retrofit/OkHttp business API 有 custom `sign` header
- POST `application/x-www-form-urlencoded` 與 header 同時存在

#### Does Not Apply When

- Sign 只覆蓋 URL path/query，不含 body（見 wire-path lesson）
- GraphQL / JSON body 且無 form merge（需另驗證 JSON canonicalization）

#### Validation

- Frida log 中 canonical 字串可離線重組且與 crypto 入參一致
- 至少一個含 body 與一個不含 body 的 POST 對照通過

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §sign RE
- `analysis/apk/traffic-triage.md` §custom headers

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 已依 sanitization / reusable-guidance-boundary 自查
