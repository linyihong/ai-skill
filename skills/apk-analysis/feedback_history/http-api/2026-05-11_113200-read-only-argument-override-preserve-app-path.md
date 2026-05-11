> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-11 - Read-Only Argument Override Preserve App Path

Status: candidate

#### One-line Summary

驗證 read-only 分頁或邊界值時，可短窗覆寫高語意函式參數，但必須保留 App 自己的 session、signing、gateway 與 decrypt path。

#### Human Explanation

完整 UI 不一定容易自然觸發下一頁、空頁、極端 filter 或錯誤矩陣列。直接重放 HTTP 又可能缺簽章、session 或 opaque material。介於兩者之間的低風險方法，是在已授權動態分析中只覆寫高語意函式的 read-only 參數，讓後續 request normalization、signing、gateway、response decrypt 仍由 App runtime 處理。

#### Trigger

- 需要驗證 pagination、sort、filter、locale 或 read-only query boundary。
- UI 很難穩定自然觸發該值。
- Direct replay 尚未具備 session/signing/decrypt parity。

#### Evidence

- Tool: short-window Frida/native hook on an app-owned function entry.
- Sanitized excerpt: original argument class and overridden target class were logged; signed request key set reflected the target page; decrypted response kept schema-only wrapper summary.
- Evidence path: concrete project logs stay in `<PROJECT_ROOT>/capture/`; reusable lesson records the method and boundaries only.

#### Generalized Lesson

Read-only argument override can validate behavior that sits behind private signing/decrypt boundaries without claiming standalone replay. The analysis must explicitly state that app-owned signing/decrypt was preserved and that the result proves only that runtime path's behavior.

#### Agent Action

- Prefer function-entry parameters over final query-map mutation.
- Keep override values limited, deterministic, and read-only.
- Log only key sets, page/class labels, length/hash, and schema shape.
- Stop if the operation becomes a write path, auth mutation, payment, destructive action, or leaves target feature context.

#### Goal / Action / Validation

- Goal: classify read-only boundary behavior when UI or direct replay is insufficient.
- Action: override a high-semantic argument in a short dynamic window and let app-owned request/decrypt continue.
- Validation or reference source: request key set contains the target class, response decrypt shape is captured, and raw private material is not logged.

#### Applies When

- The target argument controls read-only pagination, sort, filter, locale, or similar query behavior.
- The app path remains authorized and observed in the target feature.

#### Does Not Apply When

- The argument affects account state, purchases, moderation, content submission, identity material, or destructive operations.
- The goal is to prove standalone SDK replay; this only proves app-runtime path behavior.

#### Validation

Require a hook log that shows original argument class, target class, request key set or page marker, and schema-only response. Mark the evidence as `app-owned signing/decrypt preserved`.

#### Promotion Target

- `WORKFLOW.md`

#### Required Linked Updates

- Updated `feedback_history/README.md` and `feedback_history/http-api/README.md`.
- Promoted the rule into `WORKFLOW.md`.
- Checked reusable guidance boundary: no project names, hosts, endpoints, raw payloads, or local paths are included.
