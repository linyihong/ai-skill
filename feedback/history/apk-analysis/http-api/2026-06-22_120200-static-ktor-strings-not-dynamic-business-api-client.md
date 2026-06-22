> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Static Ktor/OkHttp strings do not prove Ktor is the business API client

Status: candidate

#### One-line Summary

DEX 同時出現 `io.ktor`、`okhttp3`、`retrofit2` 時，不可假設業務 API 走 Ktor pipeline；動態以 `retrofit2.Invocation` request tag 或實際 OkHttp chain 命中為準，Ktor hook 未觸發時應優先 Retrofit/OkHttp 主線。

#### Human Explanation

靜態 triage 常把「有 Ktor」解讀成「主 API = Ktor Client」。實務上 Ktor 可能只服務廣告、H5、下載器或次要模組，而 **業務 REST 仍走 Retrofit2 → OkHttp**。只 hook `HttpSend` / Ktor engine 會得到 INIT 成功但 0 條業務 URL，浪費一輪。應在首輪動態檢查 request `tags` 是否含 `retrofit2.Invocation`。

#### Trigger

- 靜態：dex strings 含 `io/ktor/client` 與 `retrofit2`
- 動態：`HttpSend` / Ktor hook 0 事件；OkHttp chain hook 有大量 `POST` 業務 path
- Request `toString()` 含 `tags={class retrofit2.Invocation=...}`

#### Evidence

- Tool: Frida OkHttp `RealInterceptorChain` hook + spawn cold start
- Sanitized excerpt: 業務 path 全部帶 Retrofit invocation tag；Ktor DefaultSender hook 未觸發業務 host
- Evidence path: `<PROJECT_ROOT>/api/dynamic-*.md`、`<PROJECT_ROOT>/docs/static-analysis.md`

#### Generalized Lesson

**HTTP client 動態確認（靜態多棧並存）：**

```text
dex 有 ktor + okhttp + retrofit？
  → 動態首選：OkHttp chain hook + 看 retrofit2.Invocation tag
  → Ktor hook 作次要；命中為 0 不阻塞主線
  → 報告寫明「verified client: Retrofit2 | Ktor | both」
```

靜態 triage 決策樹應列 **候選** 而非單一主線。

#### Agent Action

1. W1 capture 腳本預設 OkHttp chain + trust bypass，Ktor hook 可選。
2. `domain-baseline.md` 分「靜態候選」與「動態 verified」兩欄。
3. Ai-skill 不寫 target host；client 結論寫泛化規則。

#### Goal / Action / Validation

- Goal: 避免 Ktor-first hook 策略在 Retrofit 主線 app 上空轉。
- Action: 修補 short-drama stack triage lesson 的 Java 主線描述。
- Validation: 動態報告標註 Retrofit tag 證據。

#### Applies When

- 多 DEX Java/Kotlin app，dex 同時有多種 HTTP client 字串

#### Does Not Apply When

- 已確認 Dart/Flutter HTTP（另走 AOT 主線）
- 動態已見 Ktor `HttpRequestData` 命中業務 host

#### Validation

- 至少一條業務 API 的 Frida 日誌含 `retrofit2.Invocation` 或等價 tag

#### Promotion Target

- `feedback/history/apk-analysis/common/2026-06-22_110600-short-drama-stack-not-always-flutter-ktor-okhttp-triage.md`（交叉引用）
- `analysis/apk/traffic-triage.md`

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` http-api 計數 +1
