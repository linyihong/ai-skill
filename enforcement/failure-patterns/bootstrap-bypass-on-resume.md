# Bootstrap Bypass on Resume（resume session 跳過 bootstrap）

Status: validated
Class: `process-gap` / `governance-drift`

## Trigger

當 agent 從 conversation summary 被喚起（resume / continuation session），直接跳進任務執行，沒有讀 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md) + [`README.md`](../../README.md)，沒有查 [`runtime/runtime.db`](../../runtime/runtime.db) 的 phase / obligations / gates，也沒有輸出 Bootstrap Receipt，使用此 pattern。

具體觸發訊號：

- Summary 包含 "Resume directly, do not acknowledge" 類字眼
- Resume 後 first user-facing message 不含 `Bootstrap: rules=✓ phase=<id> obligations=<n> gates=<n>`
- Resume 後第一個工具呼叫直接是 Edit / Write / Bash / git，非 Read CORE_BOOTSTRAP.md
- 任務執行時 agent 沒有目前 phase 的 obligations / gates count
- Knowledge update flow 11-step 沒走過 Step 1（read bootstrap）

## Failure Mode

把 summary 的 "Resume directly" **對話 framing 指示** 當成可以跳過 **runtime / governance bootstrap**，導致：

1. **Phase machine 失盲**：不知道目前 phase 的 allowed/forbidden_actions
2. **Obligations 漏跑**：本 phase 的 obligation 直接被略過
3. **Gates 失效**：blocking gates 沒被檢查就 commit
4. **Cognitive mode 未解析**：與 [`cognitive-mode-resolution-bypass`](cognitive-mode-resolution-bypass.md) 連鎖觸發
5. **Knowledge-update-flow 子管道跳過 master**：直接跳到 "linked updates" sub-pipeline，與 [`knowledge-update-flow-bypassed-by-sub-pipeline`](knowledge-update-flow-bypassed-by-sub-pipeline.md) 同類

## Risk

- **Governance silent drift**：每次 resume 都偏離一點，累積後 agent 行為與 governance 設計脫鉤
- **Failure pattern 重演**：剛寫進 repo 的失敗模式立刻被當前 session 觸發（self-inconsistency）
- **Audit 困難**：無 Bootstrap Receipt 留下，事後無法追溯 agent 當時是否查過 runtime state

## Required Agent Action

每次 session 啟動（含 resume）的 first turn：

1. **讀 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)** + [`README.md`](../../README.md)
2. **查 [`runtime/runtime.db`](../../runtime/runtime.db)**：
   - `SELECT id FROM phase_machine WHERE active=1` → phase
   - `SELECT COUNT(*) FROM obligations WHERE phase=<current>` → obligations
   - `SELECT COUNT(*) FROM gates WHERE phase=<current>` → gates
3. **輸出 Bootstrap Receipt**（在 first user-facing message）：
   ```
   Bootstrap: rules=✓ phase=<id> obligations=<n> gates=<n>
   ```
4. **才能**執行非-Read 工具（Edit / Write / Bash / git / ...）

Summary 的 "Resume directly, do not acknowledge" 是 **對話 framing**（不要寒暄、不要 recap），不是 **runtime bootstrap 豁免**。

## Prevention Gate

每次 session 啟動，agent 必須能回答：

| Check | Required answer |
|-------|-----------------|
| 是否為 resume session | 即使是，仍要做 bootstrap |
| CORE_BOOTSTRAP.md 已讀 | 是 |
| README.md 已讀 | 是 |
| runtime.db 查詢完成 | phase / obligations / gates 已取得 |
| Bootstrap Receipt 已輸出 | 在 first user-facing message |
| 才執行非-Read 工具 | 是 |

對應 runtime enforcement：

- `obligation.bootstrap.receipt_acknowledged` — first-turn 必須輸出 Bootstrap Receipt
- `gate.bootstrap.receipt_present` — 未輸出即 block_execution

## Validation

符合下列條件時，此 pattern 已被防止：

- Resume session 的 first message 含 `Bootstrap: rules=✓ phase=...`
- `obligations` / `gates` 表 phase=phase.bootstrap 有對應 row
- [`CLAUDE.md`](../../CLAUDE.md) 與 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md) 含 Resume clause 與 Bootstrap Receipt 格式

## Source

- 2026-05-25 session：agent 從 summary 喚起後直接補做 linked updates，沒讀 CORE_BOOTSTRAP.md / README.md，沒查 runtime.db。使用者指出「你剛剛開始的時候 沒有一開始就讀 bootstrap.md 以及走更新流程」，觸發本 pattern。同 session 補上 Phase A+B+C enforcement（強化語氣 + Bootstrap Receipt + runtime obligation/gate）。
- Related scenario: `validation/scenarios/bootstrap/bootstrap-receipt-enforcement-v1.yaml`

## Related

- [`CLAUDE.md`](../../CLAUDE.md) — IMPORTANT 啟動條款 + Resume clause + Bootstrap Receipt
- [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md) — Bootstrap Receipt 格式與 §驗證 checklist
- [`runtime/runtime.db`](../../runtime/runtime.db) — `obligation.bootstrap.receipt_acknowledged` + `gate.bootstrap.receipt_present`
- [`cognitive-mode-resolution-bypass.md`](cognitive-mode-resolution-bypass.md) — 同 session 解析相關，常連鎖觸發
- [`knowledge-update-flow-bypassed-by-sub-pipeline.md`](knowledge-update-flow-bypassed-by-sub-pipeline.md) — 同類 process-gap

## Linked Validation Scenarios

- `bootstrap-receipt-enforcement-v1` — 驗證 CLAUDE.md / CORE_BOOTSTRAP.md 含 Bootstrap Receipt 條款，obligation/gate 存在於 runtime.db

← [Back to failure patterns](README.md)
