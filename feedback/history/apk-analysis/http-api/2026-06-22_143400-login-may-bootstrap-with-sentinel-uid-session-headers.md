> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Login route may bootstrap with sentinel uid/session headers

Status: candidate

#### One-line Summary

登入 API **不一定**需要已有 session；訪客路徑常可用 **`uid=0` / `session=0`（或空）+ 新 device id** 換回正式 session。Capture 裡「登入請求帶舊 session」可能是**已持久化帳號 sync**，不是 first-time bootstrap。

#### Human Explanation

Chicken-and-egg 常來自混淆兩條路：(1) 已安裝 app 從本地讀 session 再 call login sync；(2) 全新 device id 用 sentinel header 註冊訪客並拿 response 裡的 session/token。離線 SDK 應測 (2)。endpoint 名、form 欄位、login_mode 值、腳本路徑**只寫專案 evidence**。

#### Trigger

- Session assumed to require Frida / local store only
- Login capture always shows non-empty session header
- Off-device sign works but blocked on session relay

#### Evidence

- Tool: signed POST login + downstream API without app
- Sanitized excerpt: `code=0` + session field in response body
- Evidence path: `<PROJECT_ROOT>/api/dynamic-w4-session.md`

#### Generalized Lesson

```text
When login seems to require existing session:
  1. Try sentinel bootstrap headers (uid=0, session=0) + fresh device id
  2. Separate persisted-account sync vs first-time guest register in notes
  3. Parse response for session/token before assuming client-only generation
  4. Keep endpoint names and form schemas in project docs only
```

#### Agent Action

1. 專案 evidence 記錄 login form + response JSON paths。
2. 交叉引用 `142500`（partial SDK）、`141700`（sign gate）。
3. Ai-skill 不寫真實 session 值。

#### Goal / Action / Validation

- Goal: 解除「必須 attach app 取 session」假設。
- Validation: off-device downstream API success (project doc).

#### Applies When

- Guest/visitor login_mode exists
- Signed headers include session + separate login POST body

#### Does Not Apply When

- Multi-step OAuth only, no sentinel bootstrap
- Server rejects sentinel without prior challenge

#### Validation

- Project E2E documented; lesson has no target-app names

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §session bootstrap

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引
- 已依 sanitization / reusable-guidance-boundary 自查
