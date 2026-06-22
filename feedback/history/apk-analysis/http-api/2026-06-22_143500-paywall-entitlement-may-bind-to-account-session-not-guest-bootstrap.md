> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Paywall entitlement may bind to account session, not guest bootstrap

Status: candidate

#### One-line Summary

付費／解鎖 entitlement **可能綁在帳號 session（uid/session 對）**，而非「能簽名就能拿內容」。**訪客 bootstrap session** 與 **裝置已登入 session** 對同一 unlock API 可能都回 `code=0`，但只有後者讓 content API 回傳媒體 payload；client unlock flag  alone 不足。

#### Human Explanation

`143400` 解決「如何拿到第一個 session」；本條解決「session **是哪一種帳號**」。常見誤判：unlock 端點回成功 + chapter id 列表 ⇒ 已解鎖。實務上 server 可能分兩層：(1) **acknowledge** unlock intent（任何合法簽名）；(2) **grant** `play_info` / media URL（僅已建立 entitlement 的 account session）。離線 RE 應 **A/B 同請求、不同 session 來源**（guest bootstrap vs device relay），不要只測一種。endpoint 名、欄位名、`is_lock` 數值語意、腳本名留在專案 evidence。

#### Trigger

- Unlock POST returns `code=0` but content GET still empty payload
- Guest off-device SDK signs correctly yet locked items never populate media field
- Same unlock path works on device relay session but not fresh guest login
- Client sets unlock-intent flag; server still withholds blob

#### Evidence

- Tool: paired API probe — guest credentials vs device-relay credentials (Frida/Java session reader)
- Sanitized excerpt: guest → empty media field; device session → non-empty blob + decodable stream URL
- Evidence path: `<PROJECT_ROOT>/api/dynamic-w4-paywall.md`

#### Generalized Lesson

```text
When paywall unlock API succeeds but media payload stays empty:
  1. Do not treat unlock endpoint success as entitlement grant
  2. A/B test same unlock + content chain with:
       (a) guest/bootstrap session
       (b) persisted or device-relay account session
  3. Compare content response: media blob length, lock-state field, follow-up flag semantics
  4. Document capability matrix in project: sign ✅ / session type / unlock ack ✅ / media grant ❌|✅
  5. Keep endpoint names, error codes, and lock-state enums in project docs only
```

#### Agent Action

1. 專案 evidence 記錄 guest vs account session 對照表與 E2E 結果。
2. 交叉引用 `143400`（guest bootstrap）、`141700`（server-side paywall）、`142500`（session relay）。
3. Ai-skill 不寫真實 uid/session、error code 真值、API path。

#### Goal / Action / Validation

- Goal: 避免把「unlock API 200」誤當「可下載全劇」。
- Action: 實作 session-source matrix 再測 content API。
- Validation: project doc shows divergent guest vs account outcome on same chapter id.

#### Applies When

- Separate unlock/intent API + content API with opaque media field
- Guest login exists alongside persisted account on device
- Server returns success envelope with empty media for locked tier

#### Does Not Apply When

- Content API returns media solely from signature (no session binding observed)
- Unlock is purely client-side UI with no server ack endpoint
- Entitlement is CDN token only with no account session in API tier

#### Validation

- Project paired probe documented; lesson has no app/endpoint literals

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §paywall / entitlement matrix

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引
- 交叉引用 `143400` Agent Action
- 已依 sanitization / reusable-guidance-boundary 自查
