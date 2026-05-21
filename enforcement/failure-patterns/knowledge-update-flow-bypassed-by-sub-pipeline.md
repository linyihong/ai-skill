# Knowledge-Update-Flow Bypassed by Sub-Pipeline（以子流程文件替代主流程）

Status: validated
Class: `process-gap` / `validation-gap`

## Trigger

Agent 在 Ai-skill 加入新 intelligence atom / failure pattern / validation scenario / workflow step / memory entry 等可重用知識時，**只讀並只跑 sub-pipeline 文件（如 `intelligence-extraction-pipeline.md`、`failure-learning-system.md`、`linked-updates.md`），把它當成完整流程**，未先讀取 `governance/lifecycle/knowledge-update-flow.md` 作為 11-step master flow 總索引。

具體觸發信號：

- Agent 直接打開 sub-pipeline 文件，並在跑完該文件內部 step 後就準備收尾 commit
- Commit message 引用 sub-pipeline（例：「依 intelligence-extraction-pipeline.md Step 6–7 執行」），但未對照 master flow 11 step
- Master flow 中的 Step 4（寫 feedback lesson）、Step 7（failure learning 判斷）、Step 9（`ai-skill runtime compile/validate` + `runtime.db` 同 commit）、Step 11（push 後 readback）完全沒被觸發
- Sub-pipeline 內部的判斷表格被誤判（例：Step 6a「新 intelligence atom 應建立 scenario」被誤判為「純內容更新」）

## Failure Mode

Sub-pipeline 文件是「某個 master step 的內部展開」，不是完整流程。Agent 把 sub-pipeline 當成完整流程時，會發生以下漏洞：

1. **Feedback lesson 從未寫入** — Step 4 不在任何 sub-pipeline 中，只在 master flow
2. **Failure learning 未 capture** — Step 7 的「強制判斷是/否」未執行，agent 自己的 close-loop gap 從未被記錄
3. **`runtime.db` 與 source 不同步** — Step 9 的 `ai-skill runtime compile/validate` 未跑，knowledge runtime 仍是舊狀態
4. **Readback 未執行** — Step 11 push 後沒重讀更新過的入口，下一輪可能仍用舊記憶
5. **Sub-pipeline 內部判斷表格被誤讀** — 因為缺少 master flow 的「強制判斷」框架（Step 6/7 必須明確回答是/否並記錄理由），agent 用直覺替代逐欄對照

## Risk

- Knowledge runtime 與 source 永久不一致（`generated_surfaces` 過期，agent 讀到 stale routing）
- 同類型錯誤無法被未來 agent 偵測（沒寫進 `enforcement/failure-patterns/`）
- 新加入的 atom 失去 validation scenario 保護（行為退化無法被自動 catch）
- 信任損失：使用者必須親自抓 master flow 缺失，agent 沒有主動對照

## Required Agent Action

當本輪自問「Step 1: 本輪是否有新知識？」回答為「是」時：

1. **第一個 Read 必須是 `governance/lifecycle/knowledge-update-flow.md`**，不論該知識最終放到哪一層
2. 在腦中（或 task list 中）建立 11 個 step 的對照清單
3. 走到 Step 6 / 7 時才打開 sub-pipeline 文件（`intelligence-extraction-pipeline.md` / `failure-learning-system.md`），把它當成「Step 6/7 內部展開」而非「整體流程」
4. Step 6 / 7 / 9 的「強制判斷」必須明文寫出「是/否 + 理由」，不可省略
5. Commit message 必須含 11 個 step 對照註記（做了 / 不適用 / 原因），讓閉環缺口無法沉默跳過

## Prevention

- **首讀順序強制**：任何「新知識寫入」任務的 first Read 必須命中 `knowledge-update-flow.md`
- **Commit message 模板**：包含「本次 master flow 11 step 對照」段落（即使某些 step 不適用也要明列）
- **Step 6a / Step 7 判斷表格逐欄對照**：不得直覺判定，必須對照表格的「應/不應」兩欄
- **Step 9 阻斷處理**：sqlite3 / ai-skill CLI 不可用時，依 `enforcement/failure-patterns/mandatory-step-blocker-bypass.md` 立即停下並通知使用者

## Detection Signals

- Commit message 提到 sub-pipeline 名稱但未提到 master flow 11 step
- 同一 PR 內無對應 `feedback/history/` lesson 新增
- 同一 PR 內無 `runtime/runtime.db` 變更但有 `routing-registry.yaml` / `knowledge/graphs/` 變更
- 使用者在 PR 後追加質疑「為什麼沒走 X 流程」

## Related

- [`governance/lifecycle/knowledge-update-flow.md`](../../governance/lifecycle/knowledge-update-flow.md) — Master flow 11 step
- [`governance/lifecycle/intelligence-extraction-pipeline.md`](../../governance/lifecycle/intelligence-extraction-pipeline.md) — Step 6 內部展開
- [`enforcement/failure-learning-system.md`](../failure-learning-system.md) — Step 7 內部展開
- [`enforcement/linked-updates.md`](../linked-updates.md) — Step 8 內部展開
- [`enforcement/failure-patterns/mandatory-step-blocker-bypass.md`](mandatory-step-blocker-bypass.md) — 同類問題的環境阻斷版本

## Source

- 2026-05-21 session：agent 加入 3 個 intelligence atom 後，直接套 `intelligence-extraction-pipeline.md` 收尾，跳過 master flow Step 4 / 7 / 9 / 11。使用者明確指出未依 `knowledge-update-flow.md` 後才回頭補救。
- Corresponding feedback lesson: `feedback/history/development-guidance/common/2026-05-21_220226-knowledge-update-flow-master-doc-required.md`
