> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - MITM shows CDN/media but not primary API — pinning tier split

Status: candidate

#### One-line Summary

系統代理 + 系統 CA 已就緒時，若 MITM 只見 CDN/圖片域名、完全沒有主 API host，多半不是「導流失敗」，而是 **主 API 層 TLS pinning / custom trust**；CDN 與 API 應分線排查，下一步用 Java trust bypass + OkHttp chain hook，而非先換 Ktor/SSL_write hook。

#### Human Explanation

延續「導流 vs TLS 兩層」：當 proxy 已收到流量、圖片/CDN GET 可解密，但業務 REST host 0 條，常見誤判是「Frida hook 點錯」。實務上 API 與 CDN 常採 **不同 trust 策略**——媒體走標準 CA，API 走 pinning 或內建 trust store。此時 MITM 單獨無法列舉 API path，需並行 Java 層 hook（必要時 Conscrypt `TrustManagerImpl` bypass）在 app 內觀察明文 URL。

#### Trigger

- `mitmdump` / Proxyman 有 CONNECT，且有 CDN/image host 的 HTTP/2 200
- 同一 capture 窗口內 **零** 主 API host（靜態 dex 已列出該 host）
- Frida 僅 hook 標準 `okhttp3.Request` 時 0 條；換混淆 overload 或加 trust bypass 後出現 API

#### Evidence

- Tool: mitmdump + adb global http_proxy + Frida spawn
- Sanitized excerpt: MITM flow 檔僅含 media CDN 類 host；Frida chain hook 同窗口出現多條 `POST` 業務 path（path 留 project docs）
- Evidence path: `<PROJECT_ROOT>/capture/mitm_*.mitm`、`<PROJECT_ROOT>/api/dynamic-*.md`

#### Generalized Lesson

```text
MITM 有流量？
  否 → 導流 / proxy 設定 / 冷啟動時機
  是 → 有 CDN、無 API？
         是 → API pinning tier；MITM 作 CDN 佐證，API 用 Java hook + trust bypass
         否 → 正常 MITM 解碼或查 request body 加密
```

不要因「MITM 沒 API」就結論「沒有 HTTP API」或「全走 native」。

#### Agent Action

1. 並跑 MITM 與 Frida：MITM 驗證代理與 CDN；Frida 驗證 API path 與 header schema。
2. 若需解密 MITM 中的 API，再加 trust bypass；勿只調 SSL_write。
3. 報告分欄：**CDN（MITM）** vs **API（hook）** 證據來源。

#### Goal / Action / Validation

- Goal: 避免 pinning tier 誤判導致錯 hook 主線。
- Action: 與 `2026-04-30-proxy-failure-導流與-tls-兩層` 串聯使用。
- Validation: MITM 仍無 API 但 Frida 有 API → 記錄 tier split 成立。

#### Applies When

- 短劇/內容 App：靜態有 API host + 多 CDN host
- 已確認 proxy 有流量

#### Does Not Apply When

- MITM 完全無流量（先修導流）
- API 與 CDN 同 host（無 tier 可分）

#### Validation

- 同一時間窗 MITM host 列表 vs Frida URL 列表對照表寫入 project baseline

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §MITM triage
- `analysis/apk/traffic-triage.md` §pinning

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` local-proxy 計數 +1
- 已依 sanitization 自查
