# 知識更新流程（Knowledge Update Flow）

本文件整合 Ai-skill 系統中從「學到新知識」到「commit/push 完成」的完整端到端流程。它不取代既有文件的細節，而是作為**執行順序的總索引**，讓 agent 在每輪 checkpoint 知道下一步該做什麼。

> 遵守 [`enforcement/README.md`](../../enforcement/README.md)、[`enforcement/dependency-reading.md`](../../enforcement/dependency-reading.md)、[`enforcement/linked-updates.md`](../../enforcement/linked-updates.md)、[`enforcement/feedback-lessons.md`](../../enforcement/feedback-lessons.md)、[`enforcement/failure-learning-system.md`](../../enforcement/failure-learning-system.md)；本檔只定義流程順序，不重複貼上各規則全文。

---

## 流程總覽

```
每輪 checkpoint
    │
    ▼
Step 1: 觸發檢查 ─────────── 本輪是否有新知識？
    │                           (feedback-lessons.md §原則)
    ▼
Step 2: 分類知識類型 ──────── 這是什麼類型的知識？
    │                           (feedback-lessons.md §判斷流程)
    ▼
Step 3: 決定 Promotion Target ─ 應該放到哪一層？
    │                           (feedback/promotion/README.md §Promotion Targets)
    ▼
Step 4: 寫入 Feedback Lesson ── 寫入 feedback/history/<domain>/
    │                           (feedback-lessons.md §模板)
    ▼
Step 5: 更新目標層 ────────── 同步更新 intelligence / workflow / analysis / shared-rules / runtime
    │                           (promotion/README.md §Promotion Checklist)
    ▼
Step 6: 執行 Intelligence Extraction ─ 若需要提取 intelligence atoms
    │                           (governance/lifecycle/intelligence-extraction-pipeline.md)
    ▼
Step 7: 檢查 Failure Learning ── 是否需要建立 failure pattern？
    │                           (failure-learning-system.md)
    ▼
Step 8: 執行 Linked Updates ── 同步更新所有相關文件
    │                           (linked-updates.md)
    ▼
Step 9: 更新 Runtime Surfaces ── 更新 registry / summaries / graphs / reports
    │                           (knowledge/runtime/README.md)
    ▼
Step 10: 驗證 ──────────────── 執行 validators、link check、diff review
    │                           (dependency-reading.md §回寫完成門檻)
    ▼
Step 11: Commit / Push / Readback ─ 關閉 writeback transaction
                                (dependency-reading.md §Writeback Transaction Guard)
```

---

## Step 1：觸發檢查（Per-Round Checkpoint）

**來源**：[`enforcement/feedback-lessons.md`](../../enforcement/feedback-lessons.md) 第 11 行

**時機**：每個有實質進展的工作回合結束前、切回長時間專案工作前、提交 project-only evidence 前、或使用者說「繼續」展開下一輪前。

**Agent 必須自問**：

> 1. 本輪是否新增可重用技巧、validation rule、replay knob、hook/runner guard、錯誤模式、或閉環缺口？
> 2. 本輪是否涉及新增或改名目錄？→ 若是，先執行 [`directory-structure-governance.md`](directory-structure-governance.md) 的 5 步驟 Checkpoint

**判斷結果**：

| 結果 | 下一步 |
|------|--------|
| 有新知識 | 進入 **Step 2** |
| 只有 project-specific evidence | 留在 project docs，不回饋到 Ai-skill |
| 不確定 | 先標 `candidate`，進入 Step 2 |

---

## Step 2：分類知識類型

**來源**：[`enforcement/feedback-lessons.md`](../../enforcement/feedback-lessons.md) §判斷流程

**判斷流程**：

1. **確認 domain 歸屬**：這個技巧屬於哪個 skill 的 scope？
   - APK 分析技術 → `apk-analysis`
   - 開發指引 → `app-development-guidance`
   - 旅遊規劃 → `travel-planning`
   - 若不確定，讀取該 skill 的 `SKILL.md` 確認 scope 描述

2. **確認 domain 下的分類**：對應 `feedback/history/<domain>/` 下是否有對應分類目錄（如 `common/`、`flutter-dart-aot/`、`controls/`）

3. **若尚無對應分類目錄，應主動建立**

---

## Step 3：決定 Promotion Target

**來源**：[`feedback/promotion/README.md`](../../feedback/promotion/README.md) §Promotion Targets

根據知識類型決定最終目標層：

| Lesson 類型 | 目標層 | 門檻 |
|------------|--------|------|
| 單一 skill 技巧 | `skills/<skill>/WORKFLOW.md`、`TOOLS.md`、`DOCUMENTATION.md` 或 `techniques/` | Lesson 已 generalized、去敏，且 skill index 已更新 |
| 工程判斷 | `intelligence/` | 影響 trade-off、anti-pattern、route selection 或 cross-project decision |
| 執行流程 | `workflow/` | 影響 agent 如何執行 planning、review、handoff 或 validation |
| 跨 skill 或全庫規則 | `enforcement/` 或 `enforcement/failure-patterns/` | Failure class 或 prevention gate 可跨 skill 重演 |
| Runtime 導航 | `knowledge/`、`metadata/`、`runtime/` | 需要被 registry、summary、graph 或 model context report route 到 |
| 長期記憶 | `memory/` | 需要保留 replay / episodic / project abstraction boundary |

---

## Step 4：寫入 Feedback Lesson

**來源**：[`enforcement/feedback-lessons.md`](../../enforcement/feedback-lessons.md) §模板、§檔名規則

**位置**：`feedback/history/<domain>/<category>/YYYY-MM-DD_HHMMSS-<slug>.md`

**必須包含**：
- One-line Summary
- Human Explanation
- Trigger
- Evidence（已去敏）
- Generalized Lesson
- Agent Action
- Goal / Action / Validation
- Applies When / Does Not Apply When
- Validation
- **Promotion Target**（指向 Step 3 決定的目標層）
- **Required Linked Updates**（列出 Step 8 需要同步的文件）

**同步索引**：更新 `feedback/history/<domain>/README.md` 和 `feedback/history/<domain>/<category>/README.md`

---

## Step 5：更新目標層

**來源**：[`feedback/promotion/README.md`](../../feedback/promotion/README.md) §Promotion Checklist

根據 Step 3 的 Promotion Target，更新對應層的文件：

| 目標層 | 更新內容 |
|--------|---------|
| `workflow/<domain>/execution-flow.md` | 加入新的執行步驟或判斷規則 |
| `workflow/<domain>/artifact-gates.md` | 加入新的產出規範或驗證 gate |
| `intelligence/<domain>/` | 新增 intelligence atom（heuristics / anti-patterns / failure / signals） |
| `analysis/<domain>/` | 更新分析方法或技術流程 |
| `enforcement/` | 新增或更新全庫規則 |
| `enforcement/failure-patterns/` | 新增 failure pattern（見 Step 7） |
| `knowledge/` / `metadata/` / `runtime/` | 更新 registry / summary / graph（見 Step 9） |
| `memory/` | 更新長期記憶條目 |

**Promotion Checklist**：
1. ✅ 保留原 `feedback/history/<domain>/` lesson，不刪除歷史
2. ✅ 檢查 lesson 只含 generalized rule，不含 project incident raw evidence
3. ✅ 決定最小 durable target
4. ✅ 若 promotion 變成 runtime route，更新 runtime surfaces
5. ✅ 若 source lesson 仍 active，在新 layer 標 `candidate` / `dual-reference`
6. ✅ 執行 validation（見 Step 10）

---

## Step 6：Intelligence Extraction（選擇性）

**來源**：[`governance/lifecycle/intelligence-extraction-pipeline.md`](../../governance/lifecycle/intelligence-extraction-pipeline.md)

**適用時機**：當新知識來自 technique 文件、SKILL.md 或 feedback history，且需要提取 intelligence atoms 時。

**Pipeline 步驟**：

```
Step 1: 內容審計（Content Audit）
    │   識別所有可拆解元素（操作步驟、判斷決策、工具命令、失敗模式等）
    ▼
Step 2: 類型判斷（Type Classification）
    │   HOW TO DO → analysis/workflows/
    │   HOW TO THINK → intelligence/{heuristics,anti-patterns,failure,signals}/
    │   Execution Flow → workflow/execution-flow.md
    │   Artifact Gate → workflow/artifact-gates.md
    ▼
Step 3: 拆解執行（Decomposition）
    ▼
Step 4: 格式轉換（Format Transformation）
    ▼
Step 5: 標註來源（Source Annotation）
    ▼
Step 6: 驗證（Validation）
    ▼
Step 6a: 建立 Validation Scenario（架構變更／新 extraction 專用）
    ▼
Step 7: 更新索引（Index Update）
    ▼
Step 7a: Shared-Rules 同步檢查（架構變更專用）
```

---

## Step 7：Failure Learning（選擇性）

**來源**：[`enforcement/failure-learning-system.md`](../../enforcement/failure-learning-system.md)

**適用時機**：當新知識來自 agent 錯誤、close-loop gap、或重複失效模式時。

**Failure Learning Loop**：

```
1. Capture  ── 記錄發生什麼、在哪裡被發現、造成什麼風險
2. Classify ── 用 taxonomy 分類（source-mirror-drift / dependency-miss / goal-ledger-miss / validation-gap / scope-drift / handoff-gap / tool-strategy-gap / parallelization-risk）
3. Contain  ── 繼續廣泛工作前先控制當前風險
4. Promote  ── 把可重用 lesson 放到正確位置
5. Strengthen ── 補強原本可防止它的 rule / workflow / checklist / tool adapter / validation gate
6. Validate ── 確認未來 agent 找得到並能套用這個 prevention
```

**Promotion Decision**：

| Failure scope | Promotion target |
|--------------|-----------------|
| 只影響單一 conversation | `.agent-goals/` |
| 單一文件有局部 gap | 該文件的 Document TODO |
| Cross-document 或 cross-agent | `enforcement/failure-patterns/` |
| Skill-specific 重複錯誤 | `feedback/history/<domain>/` |
| Tool-specific 執行錯誤 | `ai-tools/<tool>.md` |
| 架構重構後 shared-rules 未同步 | `failure-patterns/shared-rules-architecture-drift.md` + Step 7a |
| AI 系統面執行錯誤 | `validation/scenarios/failure-derived/` |

---

## Step 8：Linked Updates

**來源**：[`enforcement/linked-updates.md`](../../enforcement/linked-updates.md)

**核心原則**：修改一個文件時，必須同步更新所有相關文件，或明確寫出「已檢查，無需更新」的理由。

**常見連動關係（摘要）**：

| 改動位置 | 必須同步更新或檢查 |
|----------|------------------|
| `enforcement/` | 根 `README.md`、相關 skill 入口、`feedback/history/` 模板引用 |
| `feedback/history/` lesson | `feedback/history/<domain>/README.md`、promotion target |
| `workflow/<domain>/` | 對應 `analysis/`、`intelligence/`、`runtime/onboarding/` |
| `intelligence/<domain>/` | `knowledge/indexes/`、`knowledge/summaries/`、`knowledge/graphs/` |
| `knowledge/` / `runtime/` | 執行 `refresh-knowledge-runtime.rb` |
| 架構重構 | 建立 validation scenario + shared-rules 同步檢查（Step 6a + Step 7a） |

> 完整表格請見 [`linked-updates.md`](../../enforcement/linked-updates.md) §常見連動關係

---

## Step 9：更新 Runtime Surfaces

**來源**：[`knowledge/runtime/README.md`](../../knowledge/runtime/README.md)

**適用時機**：當新知識影響到 routing、summary、graph 或 model context 時。

**必須更新的 surfaces**：

| Surface | 更新方式 |
|---------|---------|
| `knowledge/runtime/routing-registry.yaml` | 新增或更新 routing record |
| `knowledge/runtime/refresh-policy.yaml` | 更新 refresh / revalidate / downgrade 規則 |
| `knowledge/summaries/` | 更新對應 domain 的 summary |
| `knowledge/graphs/` | 更新對應 domain 的 graph edges |
| `knowledge/runtime/runtime-report.md` | `ruby scripts/generate-knowledge-runtime-report.rb --write` |
| `knowledge/runtime/model-context-report.md` | `ruby scripts/generate-model-context-report.rb --write` |
| `knowledge/runtime/model-checklists.md` | `ruby scripts/generate-model-checklists.rb --write` |

**驗證**：
```bash
ruby scripts/refresh-knowledge-runtime.rb
ruby scripts/validate-knowledge-runtime.rb
```

---

## Step 10：驗證（Validation）

**來源**：[`enforcement/dependency-reading.md`](../../enforcement/dependency-reading.md) §Ai-skill 回寫完成門檻、[`enforcement/sanitization.md`](../../enforcement/sanitization.md)

**驗證檢查清單**：

1. ✅ `git status --short --branch` 檢查變更
2. ✅ **去敏檢查（Sanitization）** — 依 [`enforcement/sanitization.md`](../../enforcement/sanitization.md) 檢查所有新增/修改的可重用文件：
   - 不得包含本機真實絕對路徑（改用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>` 占位符）
   - 不得包含使用者帳號名稱、私用工作目錄、git clone 實體路徑
   - 不得包含 secrets、raw tokens、私人 host、個資
   - 不得包含 project-specific evidence（依 [`reusable-guidance-boundary.md`](../../enforcement/reusable-guidance-boundary.md)）
3. ✅ `git diff` 檢查將提交的內容，確認上述去敏項目已處理
4. ✅ 執行適用的 lints / Markdown link check / required linked updates 檢查
5. ✅ **目錄結構命名檢查** — 若本輪涉及新增或改名目錄，執行 `scripts/validate-knowledge-runtime.rb` 的 `validate_directory_naming`：
   - 檢查 `intelligence/engineering/` 下是否有與根目錄同名的目錄（跨層名稱衝突）
   - 檢查目錄名稱是否為舊技能名稱的縮寫（慣性命名）
   - 檢查目錄深度是否超過 4 層
6. ✅ 若本輪使用或更新 tool mirror，執行對應 tool sync；reference-only 只需確認 `<AI_SKILL_REPO>` 可讀
7. ✅ 若有多個 owner group，使用 `scripts/ai-skill-close-loop.sh --commit` 分組提交

---

## Step 11：Commit / Push / Readback

**來源**：[`enforcement/dependency-reading.md`](../../enforcement/dependency-reading.md) §Writeback Transaction Guard、§Commit/Push 後讀回 Gate

**交易關閉條件**：

1. ✅ `git status --short --branch` 與 `git diff` 已檢查
2. ✅ 必要的 linked updates 已同步或明確寫出不適用理由
3. ✅ 若本輪使用或更新 tool mirror，必要的 tool sync 已執行
4. ✅ 相關檔案已 `git add`、`git commit`、`git push`
5. ✅ Push 後已重新讀取更新過的入口、主要依賴、索引與 promotion target
6. ✅ 最後一次 `git status --short --branch` 顯示 clean，且 branch 沒有 ahead/behind

**Commit/Push 後讀回 Gate**：

| 更新類型 | Commit/push 後必須重新讀取 |
|---------|--------------------------|
| `enforcement/` | 更新過的 shared rule、`enforcement/README.md`、`enforcement/linked-updates.md` |
| `skills/<name>/` | 該 skill 的 `SKILL.md`，以及本次更新過的 workflow / documentation / checklist |
| 工具專用規則 | 更新過的工具規則檔，以及對應的 shared rule 正文 |
| template 或 feedback lesson | 更新過的 template/lesson、索引 README、promotion target |

---

## 快速參考：每輪 Checkpoint 執行摘要

```
□ Step 1:  本輪是否有新知識？（feedback-lessons.md §原則）
□ Step 2:  分類知識類型（feedback-lessons.md §判斷流程）
□ Step 3:  決定 Promotion Target（feedback/promotion/README.md）
□ Step 4:  寫入 feedback/history/<domain>/<category>/  lesson
□ Step 5:  更新目標層（workflow / intelligence / analysis / shared-rules / runtime）
□ Step 6:  若需要，執行 Intelligence Extraction Pipeline
□ Step 7:  若需要，執行 Failure Learning Loop
□ Step 8:  執行 Linked Updates（linked-updates.md）
□ Step 9:  更新 Runtime Surfaces + 執行 scripts
□ Step 10: 驗證（diff review、link check、lint）
□ Step 11: Commit / Push / Readback（dependency-reading.md）
```

---

## 與既有文件的關係

| 文件 | 在本流程中的角色 |
|------|----------------|
| [`enforcement/feedback-lessons.md`](../../enforcement/feedback-lessons.md) | Step 1-2, 4：觸發檢查、分類、寫入 lesson |
| [`feedback/promotion/README.md`](../../feedback/promotion/README.md) | Step 3, 5：決定 promotion target、更新目標層 |
| [`feedback/pipeline/README.md`](../../feedback/pipeline/README.md) | Pipeline 架構設計，agent 在 close-loop 階段遵循的執行模型 |
| [`enforcement/failure-learning-system.md`](../../enforcement/failure-learning-system.md) | Step 7：failure capture → classify → promote |
| [`enforcement/linked-updates.md`](../../enforcement/linked-updates.md) | Step 8：多文件同步更新規則 |
| [`enforcement/dependency-reading.md`](../../enforcement/dependency-reading.md) | Step 10-11：writeback transaction、驗證、commit/push/readback |
| [`governance/lifecycle/intelligence-extraction-pipeline.md`](../../governance/lifecycle/intelligence-extraction-pipeline.md) | Step 6：從 technique/feedback/SKILL.md 提取 intelligence atoms |
| [`knowledge/runtime/README.md`](../../knowledge/runtime/README.md) | Step 9：更新 runtime surfaces |
| [`enforcement/failure-patterns/`](../../enforcement/failure-patterns/) | Step 7 的 promotion target：跨 skill failure pattern |
| [`feedback/history/apk-analysis/common/2026-05-11_125615-per-round-feedback-checkpoint.md`](../../feedback/history/apk-analysis/common/2026-05-11_125615-per-round-feedback-checkpoint.md) | 既有 lesson：per-round checkpoint 的原始記錄 |

---

← [回到 Knowledge Lifecycle](README.md)
