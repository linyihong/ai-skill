# Failure: Resume Bootstrap Bypass + Deferred Push (Same Session)

## Status

monitored

## Trigger Context

2026-05-25 session：從 conversation summary 喚起後，連續觸發三個 process gap：

1. **Resume bootstrap bypass**：把 summary 的「Resume directly, do not acknowledge」當成 runtime/governance bootstrap 豁免，沒讀 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md) / [`README.md`](../../README.md) / `runtime.db`
2. **Master flow not read**：直接補做 linked updates sub-pipeline，沒讀 [`knowledge-update-flow.yaml`](../../governance/lifecycle/knowledge-update-flow.yaml) master executable contract
3. **Push deferred too long**：累積 11 個 commits 才 push（其中包含多個獨立工作單元），不是「對話結束才 push」而是「忘了 push」

## Symptoms

- Resume 後 first user-facing message 不含 `Bootstrap: rules=✓ phase=... obligations=... gates=...`
- 直接執行 Edit/Write tool 而非先 Read bootstrap docs
- 宣稱「做完 linked updates」但未走 11-step master flow（缺 Step 4 feedback lesson、final report）
- `git log origin/main..HEAD` 累積遠超過單一邏輯工作單元的 commits
- 使用者連續三次糾正才完整修補

## Root Cause

三個 process gap 共享一個根因：**把 governance 文件當 overhead，不是 task 的一部分**。

子原因：

- Summary ≠ 讀過。Summary 是對話歷史壓縮，不是 runtime state，也不是 process contract。
- "Resume directly" scope creep。它指輸出文字 framing（不寒暄、不重述），不是流程豁免。
- Cost optimization 偏見壓過 governance contract — 估算「直接做事比較快」，但跳過反而要使用者糾正三次，總成本更高。

## Impact

- **Audit 困難**：無 Bootstrap Receipt 留下，無法追溯 agent 啟動時的 runtime state
- **Master flow drift**：sub-pipeline 取代 master，11-step 缺漏未被偵測
- **Push deferral risk**：未 push 的 commits 若本機掛掉就蒸發；CI 沒驗證；其他工作站讀不到

## Detection Signals

- First-turn 沒有 `Bootstrap: rules=✓ phase=` 字串：bootstrap bypass
- 寫 reusable knowledge 但沒有 final report 11-step status：master flow bypass
- 對話結束時 `git log origin/<branch>..HEAD` 非空：push deferred past conversation end

## Mitigation

本 session 已部署 Phase A+B+C：

- **A**：[`CLAUDE.md`](../../CLAUDE.md) + [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md) 加 IMPORTANT/MUST + Resume clause
- **B**：First-turn 強制輸出 Bootstrap Receipt
- **C**：`runtime.db` 加 `obligation.bootstrap.receipt_acknowledged` + `gate.bootstrap.receipt_present`

Feedback lesson：[`feedback/history/roo-code/2026-05-25_145500-bootstrap-bypass-on-resume.md`](../../feedback/history/roo-code/2026-05-25_145500-bootstrap-bypass-on-resume.md)

## Prevention

**Bootstrap 規則（已落實）**：
- 每 session 啟動（含 resume），first turn 讀 CORE_BOOTSTRAP + README + `runtime.db` + 必要時 `knowledge-update-flow.yaml`，輸出 Bootstrap Receipt 才執行非-Read tool。
- Summary 的 "Resume directly" 是對話 framing，**不豁免** bootstrap。

**Push 規則（per user 2026-05-25 clarification）**：
- 對話過程中 **可以多次 commit**，每個 commit 是獨立邏輯單元即可
- 但 **對話結束前必須 push + readback + clean status**（`git log origin/<branch>..HEAD` 為空、`git status` clean）
- 不必每個 commit 立刻 push（先前的 [`feedback_close_loop_must_push`](#) 草稿表述為「每 commit 後立即 push」是過嚴解讀，已修正為本規則）
- 對話結束 = 使用者表示任務完成 / 切換新任務 / agent 輸出最終 summary 前

## Generalized Lesson

Resume / continuation session 不豁免 governance bootstrap。**Master executable contract**（`knowledge-update-flow.yaml`）必須讀，不能用 sub-pipeline 取代。Push 的閉環粒度是**對話**，不是每個 commit；但對話結束前必須 push + readback。

## Occurrences

- 2026-05-25：本 session 連續三次觸發（resume bypass、master flow bypass、push deferred）；使用者三次糾正後完整修補。Phase A+B+C 上線。

## Linked Patterns

- [`enforcement/failure-patterns/bootstrap-bypass-on-resume.md`](../../enforcement/failure-patterns/bootstrap-bypass-on-resume.md)
- [`enforcement/failure-patterns/knowledge-update-flow-bypassed-by-sub-pipeline.md`](../../enforcement/failure-patterns/knowledge-update-flow-bypassed-by-sub-pipeline.md)
- [`enforcement/failure-patterns/commit-before-validation-skip.md`](../../enforcement/failure-patterns/commit-before-validation-skip.md)
- [`enforcement/failure-patterns/cognitive-mode-resolution-bypass.md`](../../enforcement/failure-patterns/cognitive-mode-resolution-bypass.md)
