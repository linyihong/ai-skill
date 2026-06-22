> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Per-request custom signatures block standalone SDK until RE or in-app relay

Status: candidate

#### One-line Summary

Primary business API 若每請求帶 **custom signature header + encrypted timestamp/session/device/trace headers**，在演算法與 session bootstrap 未逆向前，**無法**從零實作 standalone HTTP SDK；過渡策略是 **in-app OkHttp relay**（app 代簽）或 interceptor/crypto hook RE；CDN/media URL 可見 **不能**替代 API 授權。

#### Human Explanation

動態 capture 常先拿到 path、form body、response JSON（尤其 converter hook），但離線 SDK 仍卡在簽名層。Signature 通常隨 path/body/time 變化；encrypted timestamp blob 與 plain epoch 可能並存。簽名邏輯常在 R8 混淆的 OkHttp Interceptor。MITM 對 pinning tier 的 API 不可見時，更不能用 CDN 流量推 API 契約。Paywall/entitlement 應預設 **server-side**：content API 對 locked item 可能拒絕或空 payload；client flag 篡改不等於長期 bypass。

#### Trigger

- 需求：非 app 進程呼叫 business API / 批量下載
- Header 含 sign（固定長度 hex）+ 非 plain requestTime + session
- jadx 找不到明文 SignUtil；MITM 無 primary host
- Path/body schema 已有，但 replay 失敗

#### Evidence

- Tool: Frida OkHttp chain hook + redacted header schema 寫 project docs
- Sanitized excerpt: sign 長度與 per-request 變化；requestTime 為密文 blob；CDN tier 與 API tier 分離
- Evidence path: `<PROJECT_ROOT>/api/dynamic-*.md`、`<PROJECT_ROOT>/docs/domain-baseline.md`

#### Generalized Lesson

```text
Standalone SDK readiness gate:
  HAVE: paths, form fields, response JSON, CDN URL shape (if API succeeded)
  NEED: sign canonical string + secret/key, encrypted timestamp algo, session bootstrap
  INTERIM: Frida proceed relay OR addHeader("sign") stack + MessageDigest/Hmac/Cipher trace
  RE order: addHeader hook → crypto hooks → diff two requests same session different body
  ENTITLEMENT: assume server-side; document locked-item content API matrix in project only
```

#### Agent Action

1. Project 維護 SDK gap matrix（headers / session / sign / entitlement）。
2. Ai-skill 不寫 crack 步驟、target endpoint、live header 值。
3. 與 pinning-tier split lesson 並用：API 證據來自 hook，CDN 來自 MITM。

#### Goal / Action / Validation

- Goal: 防止低估 SDK 工作量或誤用 CDN URL 當授權。
- Action: 與 `wire-path-vs-signing-canonical-path`、`sha256-hash-verify-python-not-shell` 交叉引用。
- Validation: gap matrix 標 sign=blocking；relay 可重放至少 1 業務 path 後再談離線 SDK。

#### Applies When

- Mobile REST 有 custom sign + session headers
- Pinning tier 使 MITM 看不到 primary API

#### Does Not Apply When

- API 無 custom sign（僅 Bearer/static token）
- 已有可重放 sign 實作或官方 SDK

#### Validation

- Relay 或 RE 後：同 header 族可離線重放 1+ path；locked item 行為在 project docs 有對照

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §SDK readiness gate
- `development-guidance/`（若跨到 external SDK 設計）

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 已依 sanitization / reusable-guidance-boundary 自查
