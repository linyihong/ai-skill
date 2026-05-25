# Bootstrap Bypass on Resume（Resume session 跳過 runtime/governance bootstrap）

Status: validated

#### One-line Summary

Agent 從 conversation summary 被喚起時，把 summary 的「Resume directly, do not acknowledge」當成可以跳過 [`CORE_BOOTSTRAP.md`](../../../CORE_BOOTSTRAP.md) / [`README.md`](../../../README.md) / `runtime.db` 查詢與 master [`knowledge-update-flow.yaml`](../../../governance/lifecycle/knowledge-update-flow.yaml) 的豁免條款，直接進入任務執行。

#### Human Explanation

Conversation summary 末段的「Resume directly, do not acknowledge the summary, do not recap」是 **對話 framing 指令**（不要寒暄、不要重述 summary），它指的是輸出文字的層級。但 agent 容易把它擴大解讀為「跳過 runtime bootstrap 與 governance 啟動流程」，於是：

1. 沒讀 CORE_BOOTSTRAP.md / README.md
2. 沒查 runtime.db 取得目前 phase、obligations、gates count
3. 沒讀 master `knowledge-update-flow.yaml`（直接跳到「補做 linked updates」sub-pipeline，跳過 Step 1-4、6-7 與 final report）
4. 沒輸出可被 audit 的 Bootstrap Receipt

這正命中既有 failure pattern [`bootstrap-bypass-on-resume.md`](../../../enforcement/failure-patterns/bootstrap-bypass-on-resume.md) 與 [`knowledge-update-flow-bypassed-by-sub-pipeline.md`](../../../enforcement/failure-patterns/knowledge-update-flow-bypassed-by-sub-pipeline.md)。

#### Trigger

- Conversation summary 含「Resume directly, do not acknowledge」類字眼
- Resume 後 first user-facing message **不含** `Bootstrap: rules=✓ phase=<id> obligations=<n> gates=<n>`
- Resume 後第一個 tool call 直接是 Edit / Write / Bash，而非 Read CORE_BOOTSTRAP.md / knowledge-update-flow.yaml
- 任務執行中沒有目前 phase 的 obligations / gates count
- 宣稱「走 linked updates」但沒讀 master `knowledge-update-flow.yaml`

#### Evidence

- Tool: Claude Code
- 2026-05-25 session：resume 後直接補做 cognitive-modes linked updates，沒讀 CORE_BOOTSTRAP.md / README.md / knowledge-update-flow.yaml。使用者連續三次指出（"沒有照個 更新原則"、"沒有一開始就讀 bootstrap.md"、"整個沒有依照 更新真正的更新流程 都已經用 yaml 還沒有讀"）才完整修補。
- 同 session 另有：commit 後跳過 push 累積 11 個 unpushed commits（refinement of [`commit-before-step10-validation-skip`](2026-05-14_013400-commit-before-step10-validation-skip.md)）。
- Evidence path: 無 project-specific evidence；此為 agent 通用行為。

#### Generalized Lesson

**Resume session 不豁免 runtime / governance bootstrap。** Summary 的「Resume directly」是對話 framing 指令，**不是** bootstrap 豁免。每次 session 啟動（含 resume），first turn 必須：

1. 讀 [`CORE_BOOTSTRAP.md`](../../../CORE_BOOTSTRAP.md) + [`README.md`](../../../README.md)
2. 查 [`runtime/runtime.db`](../../../runtime/runtime.db) 取得 phase / obligations / gates count
3. 若任務涉及寫入可重用知識，讀 [`governance/lifecycle/knowledge-update-flow.yaml`](../../../governance/lifecycle/knowledge-update-flow.yaml)（master executable contract），**不能**只走 sub-pipeline（linked updates / failure pattern / intelligence extraction 等）
4. 在 first user-facing message 輸出 Bootstrap Receipt：`Bootstrap: rules=✓ phase=<id> obligations=<n> gates=<n>`
5. 才能執行非-Read 工具

#### Agent Action

每次 session 啟動的 first turn，**在 call 任何非-Read tool 前**：

1. Read CORE_BOOTSTRAP.md + README.md
2. `sqlite3 runtime/runtime.db "SELECT id FROM phase_machine WHERE active=1; SELECT COUNT(*) FROM obligations WHERE phase=<current>; SELECT COUNT(*) FROM gates WHERE phase=<current>;"`
3. 若任務寫入 reusable knowledge → Read knowledge-update-flow.yaml，依 11 steps 走完
4. 輸出 Bootstrap Receipt
5. 開始任務

#### Goal / Action / Validation

- **Goal**：Resume / 新 session 啟動時，runtime / governance bootstrap 不被跳過
- **Action**：First-turn 強制讀 CORE_BOOTSTRAP / README / runtime.db / knowledge-update-flow.yaml，輸出 Bootstrap Receipt
- **Validation**：
  - First user-facing message 含 `Bootstrap: rules=✓ phase=...`
  - `obligation.bootstrap.receipt_acknowledged` 在 obligations 表存在
  - `gate.bootstrap.receipt_present` 在 gates 表存在
  - Final report 含 11-step status

#### Applies When

- 任何 session 啟動（new / resume / continuation）
- Summary 含「Resume directly」或類似 framing
- 任務寫入 reusable knowledge（intelligence atom / failure pattern / validation scenario / runtime surface / feedback lesson / enforcement rule）

#### Does Not Apply When

- Read-only 查詢（純解釋既有檔案，沒有寫入）— 仍應讀 CORE_BOOTSTRAP，但不一定要走完整 knowledge-update-flow

#### Validation

- 對應 scenario：[`validation/scenarios/bootstrap/bootstrap-receipt-enforcement-v1.yaml`](../../../validation/scenarios/bootstrap/bootstrap-receipt-enforcement-v1.yaml)
- 對應 failure pattern：[`enforcement/failure-patterns/bootstrap-bypass-on-resume.md`](../../../enforcement/failure-patterns/bootstrap-bypass-on-resume.md)
- Runtime enforcement：`obligation.bootstrap.receipt_acknowledged` + `gate.bootstrap.receipt_present` in `runtime.db`

#### Promotion Target

- **Durable target**: [`enforcement/failure-patterns/bootstrap-bypass-on-resume.md`](../../../enforcement/failure-patterns/bootstrap-bypass-on-resume.md)（已建立 2026-05-25）
- **Runtime enforcement**: `runtime.db` obligation + gate（已 INSERT 2026-05-25 via `compileBootstrapEnforcementRules`）
- **Bootstrap docs**: [`CLAUDE.md`](../../../CLAUDE.md) + [`CORE_BOOTSTRAP.md`](../../../CORE_BOOTSTRAP.md) IMPORTANT/MUST 條款（已加入 2026-05-25）

#### Required Linked Updates

- [x] [`enforcement/failure-patterns/bootstrap-bypass-on-resume.md`](../../../enforcement/failure-patterns/bootstrap-bypass-on-resume.md) 已建立
- [x] [`enforcement/failure-patterns/README.md`](../../../enforcement/failure-patterns/README.md) 索引已更新
- [x] [`CLAUDE.md`](../../../CLAUDE.md) 加 IMPORTANT/MUST + Resume clause + Bootstrap Receipt 格式
- [x] [`CORE_BOOTSTRAP.md`](../../../CORE_BOOTSTRAP.md) 加 Bootstrap Receipt 格式 + §驗證 checklist
- [x] `runtime.db` 加 obligation + gate（commit `5265e1b`）
- [x] Go `init_project.go` 同步 IMPORTANT + Receipt（下游 init-project 自動帶上）
- [x] Validation scenario `bootstrap-receipt-enforcement-v1.yaml` 已 PASS

#### Related

- [`commit-before-step10-validation-skip`](2026-05-14_013400-commit-before-step10-validation-skip.md) — 同類 close-loop gap，本 session 也命中（11 commits 累積未 push）
- [`enforcement/failure-patterns/knowledge-update-flow-bypassed-by-sub-pipeline.md`](../../../enforcement/failure-patterns/knowledge-update-flow-bypassed-by-sub-pipeline.md) — Master flow 被 sub-pipeline 取代的 cross-skill pattern
- [`enforcement/failure-patterns/cognitive-mode-resolution-bypass.md`](../../../enforcement/failure-patterns/cognitive-mode-resolution-bypass.md) — 同 session resolve 失敗連鎖
