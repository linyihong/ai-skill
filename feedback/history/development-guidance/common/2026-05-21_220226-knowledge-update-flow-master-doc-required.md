> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-21 — 知識更新必須先讀 `knowledge-update-flow.md` 11-step master flow，不得以子流程文件替代

Status: validated

#### One-line Summary

當「有新知識要寫入 Ai-skill」時，必須先讀 `governance/lifecycle/knowledge-update-flow.md` 作為流程總索引；**不得只讀子流程文件（如 `intelligence-extraction-pipeline.md`）並將其視為完整流程**，否則會跳過 Step 4（feedback lesson）、Step 7（failure learning）、Step 9（runtime.db 重編）、Step 11（readback）等 master flow 才有的強制步驟。

#### Human Explanation

Ai-skill 的知識更新流程設計成兩層：

1. **Master flow** = `governance/lifecycle/knowledge-update-flow.md`，定義 11 個強制 step（觸發 → 分類 → promotion → lesson → 更新目標層 → intelligence extraction 判斷 → failure learning 判斷 → linked updates → runtime surfaces → 驗證 → commit/push/readback）。
2. **Sub-pipelines** = 各個 Step 內部展開的細節文件（如 `intelligence-extraction-pipeline.md` 是 Step 6 的子流程；`failure-learning-system.md` 是 Step 7 的子流程；`linked-updates.md` 是 Step 8 的子流程）。

Agent 容易犯的錯：直覺讀到「我這次的任務是新增 intelligence atom」→ 直接打開 `intelligence-extraction-pipeline.md` 並把它的 7 個內部 step 當成完整流程跑完。**結果跳過 master flow 中的 Step 4 / 7 / 9 / 11**，因為這些 step 不在 sub-pipeline 範圍內。

#### Trigger

本輪在 Ai-skill 加入 3 個 intelligence atom 後，agent 直接套用 `intelligence-extraction-pipeline.md` 跑完 Step 6 內部流程後就準備收尾（commit + push）。使用者明確指出「沒有依照 `knowledge-update-flow.md`」，agent 才回頭對照 11 個 step，發現完全沒寫 feedback lesson、沒做 failure learning capture、沒跑 `ai-skill runtime compile/validate`、沒做 readback。

#### Evidence

- Tool: 讀取 `governance/lifecycle/knowledge-update-flow.md` 與 `governance/lifecycle/intelligence-extraction-pipeline.md` 對比
- Sanitized excerpt: master flow 的 Step 4 / 7 / 9 / 11 未在 sub-pipeline 中出現
- Evidence path: `<AI_SKILL_REPO>/governance/lifecycle/knowledge-update-flow.md`

#### Generalized Lesson

**當「有新知識要寫入可重用層」時，第一個讀取的文件必須是 `governance/lifecycle/knowledge-update-flow.md`**，不論該知識最終會放到 intelligence / workflow / analysis / enforcement / memory 哪一層。sub-pipeline 文件只負責某個 step 內部展開，不取代 master flow。

判斷句：

| 觸發條件 | 必須讀取 |
|---------|---------|
| 有新知識（atom / pattern / heuristic / failure / scenario） | `knowledge-update-flow.md`（master）|
| Master flow 走到 Step 6 且決定執行 extraction | + `intelligence-extraction-pipeline.md`（sub）|
| Master flow 走到 Step 7 且決定執行 failure learning | + `failure-learning-system.md`（sub）|
| Master flow 走到 Step 8 | + `linked-updates.md`（sub）|

順序反過來（先讀 sub、後想起 master）幾乎必定漏 step。

#### Agent Action

1. 當本輪自問「Step 1: 本輪是否有新知識？」回答為「是」時，**立刻打開 `governance/lifecycle/knowledge-update-flow.md` 並逐條對照 11 個 step**，而非直接跳到「該怎麼 extraction」。
2. 在 commit message 中明列「本次 master flow 11 個 step 對照表」，每個 step 註明「做了 / 不適用 / 原因」，使閉環缺口無法沉默跳過。
3. Step 6 與 Step 7 的「強制判斷」必須明文寫出「是/否 + 理由」，不可省略。

#### Goal / Action / Validation

- Goal: 確保 Ai-skill 知識更新永遠走完 master flow 11 step，不被 sub-pipeline 替代
- Action: 任何「新知識寫入」任務的第一個 Read 必須是 `knowledge-update-flow.md`
- Validation: commit message 含 11 個 step 對照註記；下次同類任務重讀本 lesson 應能 prevent 重複錯誤

#### Applies When

- Agent 即將在 Ai-skill 加入新 intelligence atom、failure pattern、validation scenario、workflow step、memory entry
- Agent 即將修改 `routing-registry.yaml` 或 `knowledge/graphs/`
- Agent 即將在 `feedback/history/` 寫 lesson

#### Does Not Apply When

- 純 project-specific evidence 收集（不回饋 Ai-skill）
- 修正 typo / 格式而無實質知識變更

#### Validation

下次同類任務（新增 intelligence atom）執行時，第一個 Read 應指向 `governance/lifecycle/knowledge-update-flow.md`；如未指向，視為 lesson 未生效，需重新檢視觸發條件。

#### Promotion Target

- `enforcement/failure-patterns/knowledge-update-flow-bypassed-by-sub-pipeline.md`（新增 failure pattern，跨 skill 重演風險高）
- `governance/lifecycle/knowledge-update-flow.md`（在文件最頂端加 callout，強調「即使任務看起來只涉及單一 sub-pipeline，仍必須以本文件為總索引」）

#### Required Linked Updates

- 已依 [`linked-updates.md`](../../../enforcement/linked-updates.md) 檢查：
  - `enforcement/failure-patterns/knowledge-update-flow-bypassed-by-sub-pipeline.md`：本輪新建（見 Step 7 產出）
  - `governance/lifecycle/knowledge-update-flow.md`：建議加 callout（本輪暫不修改，避免擴大 commit 範圍；列為 follow-up）
  - `feedback/history/development-guidance/README.md`：加入本 lesson 索引列

- Step 6（Intelligence Extraction）：**否** — 本 lesson 是 process / governance 知識，不屬於工程判斷 atom；已在 `intelligence-extraction-pipeline.md` 已存在 governance 範圍內。
- Step 7（Failure Learning）：**是** — 已新建 `enforcement/failure-patterns/knowledge-update-flow-bypassed-by-sub-pipeline.md`。
