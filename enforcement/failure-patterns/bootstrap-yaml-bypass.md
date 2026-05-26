# Bootstrap YAML Bypass（agent 跳過 generated_surfaces 直讀 prose）

Status: validated
Class: `governance-drift` / `source-of-truth-duplication`

## Trigger

Agent 在 session 啟動或執行 cognitive obligation lookup 時，**直接讀 `CORE_BOOTSTRAP.md` prose** 或 individual contract YAML 檔，而**沒**透過 `generated_surfaces[runtime.core_bootstrap.contract]` query 或 `ai-skill runtime obligations` CLI 取得 machine-readable obligations list，使用此 pattern。

具體訊號：

- Agent first-turn 沒查 `runtime/runtime.db` 取得 per_turn / per_commit obligations 名稱列表
- Bootstrap Receipt 只列 count（`obligations=3 gates=2`），沒列 obligation names（無 enhanced two-line form）
- Per-obligation dispatcher refactor 完成後，仍有 validator 從 hardcode list 跑而非從 generated_surfaces JSON 載入
- 新加 obligation 到 YAML 後 hook 行為沒變（因為 hook 沒讀 generated_surfaces）

## Failure Mode

`runtime/core-bootstrap.yaml` 是 canonical executable contract，projected to `generated_surfaces[runtime.core_bootstrap.contract]`。Markdown companion 是 human rationale，**不是 source-of-truth**。

繞過 generated_surfaces 直接讀 prose 會：

1. **Schema 同步落後** — 加 obligation 到 YAML 後 agent / hook 不知道（讀的 .md 沒更新）
2. **List vs count drift** — Bootstrap Receipt 應列 obligation names（per Phase 5 enhanced form），但 agent 從 prose 只看到 format example，無法列出實際 IDs
3. **Per-obligation dispatcher 失效** — Phase 6 dispatcher refactor 從 generated_surfaces 讀 per_commit_obligations 動態 dispatch；若 agent / hook 沒 query SQLite，dispatcher 無作用

## Required Agent Action

Session 啟動 first turn：

1. Query `runtime/runtime.db generated_surfaces[runtime.core_bootstrap.contract]` 或執行 `ai-skill runtime obligations`
2. Bootstrap Receipt 包含 `Active per-turn obligations: <comma-separated ids>` 行（per Phase 5）
3. 若任務涉及 commit，hook 已 dispatch validators by per_commit_obligations 列表（per Phase 6 dispatcher refactor）

修改 obligation 規則時：

1. 改 `runtime/core-bootstrap.yaml`（canonical）
2. `runtime compile + refresh`
3. 同步 markdown companion（CORE_BOOTSTRAP.md）

## Prevention Gate

- **Mechanical**: Phase 6 dispatcher refactor — `runCommitMsgHook` 從 `generated_surfaces` 讀 per_commit_obligations 動態 dispatch validators by id，**不再 hardcode** list。新 obligation 加到 YAML + Go registry 即可，hook 行為自動跟上。
- **Audit CLI**: `ai-skill runtime obligations` 列當前 obligations，agent 可隨時驗證實際載入 vs YAML 宣告。
- **Bootstrap Receipt 形式**: enhanced two-line form（含 `Active per-turn obligations:` 行）要 query 才寫得出來 — forcing function。

## Validation

符合下列條件即此 pattern 已被防止：

- Session first turn 含 enhanced two-line Bootstrap Receipt（含 obligation IDs）
- `ai-skill runtime obligations` 輸出與 hook 實際 dispatch 一致
- 加 obligation 到 YAML 後，hook 行為自動跟上（不需要改 hooks.go）

## Source

- 2026-05-26 session：Phase 5 enhanced Receipt 落地 + Phase 6 per_commit_obligations 11-entry enumeration 完成。Dispatcher refactor（Go 改 runCommitMsgHook 讀 generated_surfaces dispatch）為配對 prevention gate。

## Related

- [`runtime/core-bootstrap.yaml`](../../runtime/core-bootstrap.yaml) — canonical contract
- [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md) — companion only
- [`bootstrap-bypass-on-resume.md`](bootstrap-bypass-on-resume.md) — 上一層 bootstrap bypass（resume 沒走 flow）
- [`runtime-yaml-unprojected.md`](runtime-yaml-unprojected.md) — 同類「YAML 規則沒生效」的不同表現

## Linked Validation Scenarios

- `bootstrap-obligations-enumerable-v1`（per_commit_obligations 必須 JSON1 query-able）
- `receipt-includes-active-obligations-v1`（enhanced Receipt format）

← [Back to failure patterns](README.md)
