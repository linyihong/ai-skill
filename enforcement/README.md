# 強制執行規則（分類索引）

> **Layer**: Runtime Enforcement Layer
> **Role**: 定義 AI agent 必須遵守的執行政策
> **Position**: `governance/` (Policy Architecture) → `enforcement/` (Runtime Enforcement) → `runtime/` (Runtime Engine)

## 範圍邊界（Scope Boundary）

### ✅ 什麼該放這裡

- **AI agent 必須遵守的執行政策**：跨 workflow 共用的強制規則，例如授權範圍、去敏規則、依賴讀取鐵則、連動更新義務
- **Core Bootstrap 規則**：每個 session 啟動時必須載入的規則（rule-weight、dependency-reading、conversation-goal-ledger）
- **Lazy-load 規則**：依條件 activate 的規則（linked-updates、failure-learning、sanitization、feedback-lessons 等）
- **Failure Patterns**：跨 workflow / 全庫可重用的 agent 失效模式（`failure-patterns/`）
- **Runtime Activation Model**：定義 Core Bootstrap + Lazy-load 的載入策略

### ❌ 什麼不該放這裡

- **知識治理的架構設計** → 放 [`governance/`](../governance/README.md)（生命週期、驗證關卡、清理策略、萃取管線）
- **AI 系統的動態載入與路由實作** → 放 [`runtime/`](../runtime/README.md)（context routing、activation rules、onboarding）
- **可重用的工程智慧（heuristics、signals、tradeoffs）** → 放 [`intelligence/`](../intelligence/README.md)
- **專案特定的 incident 證據** → 放 [`feedback/history/`](../feedback/history/) 或專案文件
- **工具特定的設定與同步腳本** → 放 [`ai-tools/`](../ai-tools/README.md) 或工具設定檔
- **業務專案的分析產出** → 放 [`analysis/`](../analysis/README.md)

### 層間互動規則

1. `governance/` 定義 `enforcement/` 的治理架構（如何建立/更新/廢棄規則），但不定義規則的具體內容。
2. `enforcement/` 定義 `runtime/` 的載入政策（哪些必讀、何時 activate），但不定義如何載入。
3. `runtime/` 實作 `enforcement/` 的載入要求，但不定義載入什麼。
4. 上層可引用下層（governance → enforcement → runtime），下層不應反向引用上層的具體內容。

---

本目錄放**所有 skill 共用**的政策與約定，依主題分檔維護。**不要**在每一則 `feedback_history` lesson 裡重複貼上全文；條目頂部用相對路徑**連回此處**即可（模板與檔名規則見 [feedback-lessons.md](feedback-lessons.md)）。

## Agents（必讀）

1. **本檔是指索引，不是全文。**請先讀此 README，再依你將要做的行為，**打開並讀完下方表格中所有相關連結的全文**（不得只讀本頁摘要就當已遵守規則）。
2. **啟動流程：**每次 session 啟動時，先讀 [`CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md)（3 條核心規則，~800 tokens），再讀此 README 的索引與 **[Runtime Activation Model](#runtime-activation-model)**。不要一次把所有 enforcement rules 當作已讀。若規則看似衝突，先依 [rule-weight.md](rule-weight.md) 判斷權重。
3. **依賴讀取鐵則：**只要發現某個 workflow、enforcement rule、tool-specific rule、模板或 feedback lesson 已更新、將被更新、或會影響目前任務，必須依 **[dependency-reading.md](dependency-reading.md)** 讀取相關依賴文件；包含 `enforcement/` 本身，不可只讀單一 workflow 檔。若使用者要求「重新讀 workflow」或指出 enforcement rules 可能漏讀，必須先建立 dependency read ledger，列出已讀、缺檔不適用與仍阻塞項目。
4. **撰寫或修改可重用文件、撰寫或新增 `feedback_history` lesson、或任何「回饋／沉澱技巧」行為前**，必須已讀過 **[feedback-lessons.md](feedback-lessons.md)**、**[reusable-guidance-boundary.md](reusable-guidance-boundary.md)**、**[sanitization.md](sanitization.md)**、**[neutral-language.md](neutral-language.md)** 與 **[goal-action-validation.md](goal-action-validation.md)**；其中 `neutral-language.md` 的目標文件語言確認是寫入前 blocking gate。若涉及授權邊界，另讀 **[authorization-scope.md](authorization-scope.md)**。
5. 索引與各分檔中的 **Markdown 連結**視為規則的一部分；請**循連結讀取**，直到該任務需要的條文都已載入為止。
6. **Git：**若你改動的是 **Ai-skill 這個 git repository**（`enforcement/`、`workflow/`、`analysis/`、`intelligence/`、根目錄 README 等），**除非使用者明講不要提交**，否則**必須**在該 repo 根目錄完成 **`git add` → `git commit` → `git push`**；請自行申請所需工具權限。重新載入或重啟任一工具**不取代**這一步。第一次寫入 Ai-skill 或其工具同步路徑時，必須依 [dependency-reading.md](dependency-reading.md) 先定位 canonical `<AI_SKILL_REPO>` git repo，再開啟 writeback transaction；工具部署 / mirror 路徑（如 `~/.cursor/`、`~/.claude/`、專案本機設定檔）不可當成 source repo 完成回寫。具體工具部署路徑見 [`ai-tools/agent/`](../ai-tools/agent/) 中各工具文件。在切回專案分析或回覆完成前關閉 transaction。同一輪對話可依獨立邏輯單元建立多個 commit，不必每個 commit 立刻 push；但任務完成、切換新任務或最終回覆前必須 push、讀回，並確認 clean status 與 `git log origin/<branch>..HEAD` 為空。Commit/push 與必要 tool sync 完成後，還必須重新讀取本次更新過的 workflow/shared-rule 入口與主要依賴文件，並用 `git status --short --branch` 確認工作樹乾淨且沒有 ahead/behind。若使用者未授權 push / merge 而 repo 仍有 pending commit、ahead/behind 或未合併狀態，最終回覆必須主動提醒使用者。**優先使用 `ai-skill close-loop --commit --push` 執行分組提交與推送**，而非直接下 `git add` / `git commit` / `git push` 指令；close-loop command 會檢查 active lock、依 owner 分組提交、避免混入不相干變更。詳見 [dependency-reading.md](dependency-reading.md) 的 writeback transaction 章節與 [scripts/README.md](../scripts/README.md) 的 close-loop automation 說明。
7. **工具同步：**預設 reference-first 時不需要同步工具 mirror。只有本機明確使用工具特定的 symlink / bundle / copy mirror，且希望該 mirror 立即跟上時，才執行對應同步流程；具體工具路徑與命令放在 [`ai-tools/`](../ai-tools/README.md)。**Agent** 不得因為改了 `enforcement/` 或 `workflow/` 就自動跑工具同步，除非使用者要求、環境變數啟用，或目前任務就是維護該 mirror。
8. **連動更新：**任何會影響其他文件、索引、workflow 入口、同步流程或分類文件的改動，都**必須**依 [linked-updates.md](linked-updates.md) 同步更新或明確檢查；不得把必要連動說成「可選」。
9. **修復後必須加入測試：**每次修復錯誤後，必須執行以下步驟：(1) 泛化錯誤模式，(2) 檢查 `ai-skill runtime validate` 中是否有測試可以檢測這類錯誤，(3) 若無則加入新的測試方法，(4) 驗證測試有效性。詳見 [failure-patterns/failure-to-validator-closure.md](failure-patterns/failure-to-validator-closure.md)。

## Runtime Activation Model

本索引採用 **runtime activation model**：規則分為 **Core Bootstrap**（每個 session 必讀）與 **Lazy-load**（只在特定條件 activate）。

## Executable YAML Contracts

所有 active enforcement rule 都有 companion YAML contract，供 agent 先讀結構化 activation、required sources、steps、gates、failure modes 與 final status requirements；Markdown 保留完整說明與維護脈絡。Runtime 不再維護 `activation_rules` / `activation_rules_mirror` tables，避免與 owner-layer contract 雙寫。

| Rule | Contract | 用途 |
| --- | --- | --- |
| [`authorization-scope.md`](authorization-scope.md) | [`authorization-scope.yaml`](authorization-scope.yaml) | 授權範圍、第三方資料與可重用證據安全 gate。 |
| [`content-layering.md`](content-layering.md) | [`content-layering.yaml`](content-layering.yaml) | 內容 owner layer、project evidence boundary 與索引更新 gate。 |
| [`conversation-goal-ledger.md`](conversation-goal-ledger.md) | [`conversation-goal-ledger.yaml`](conversation-goal-ledger.yaml) | Active goal lifecycle、owner/lock、handoff、completion/delete gate。 |
| [`cross-skill-references.md`](cross-skill-references.md) | [`cross-skill-references.yaml`](cross-skill-references.yaml) | Cross-workflow reference、handoff artifact、ownership 與 sanitization boundary gate。 |
| [`decision-efficiency.md`](decision-efficiency.md) | [`decision-efficiency.yaml`](decision-efficiency.yaml) | 決策點、路線比較、context cost 與 validation signal gate。 |
| [`dependency-reading.md`](dependency-reading.md) | [`dependency-reading.yaml`](dependency-reading.yaml) | 依賴讀取 ledger、source-of-truth miss recovery、writeback final status gate。 |
| [`document-todo-list.md`](document-todo-list.md) | [`document-todo-list.yaml`](document-todo-list.yaml) | Document-local TODO placement、linking 與 completion validation gate。 |
| [`escalation-policy.md`](escalation-policy.md) | [`escalation-policy.yaml`](escalation-policy.yaml) | Escalation trigger、source reload、execution graph rebuild 與 recovery output gate。 |
| [`evidence-hierarchy.md`](evidence-hierarchy.md) | [`evidence-hierarchy.yaml`](evidence-hierarchy.yaml) | Evidence quality、claim scope、confidence integrity 與 autonomy downgrade gate。 |
| [`failure-learning-system.md`](failure-learning-system.md) | [`failure-learning-system.yaml`](failure-learning-system.yaml) | Failure capture、classification、durable promotion、guardrail strengthening 與 validation gate。 |
| [`feedback-lessons.md`](feedback-lessons.md) | [`feedback-lessons.yaml`](feedback-lessons.yaml) | 舊 feedback lesson redirect 與 canonical feedback source gate。 |
| [`goal-action-validation.md`](goal-action-validation.md) | [`goal-action-validation.yaml`](goal-action-validation.yaml) | 重要工作單元的目標、執行、驗證與 completion claim gate。 |
| [`linked-updates.md`](linked-updates.md) | [`linked-updates.yaml`](linked-updates.yaml) | 連動更新矩陣、runtime sync、writeback closure 與 final report gate。 |
| [`neutral-language.md`](neutral-language.md) | [`neutral-language.yaml`](neutral-language.yaml) | 可重用文件語言一致性、中性用語與敏感詞處理 gate。 |
| [`prompt-cache-efficiency.md`](prompt-cache-efficiency.md) | [`prompt-cache-efficiency.yaml`](prompt-cache-efficiency.yaml) | Context layout、stable prefix、provider cache metadata 與 required-read preservation gate。 |
| [`reusable-guidance-boundary.md`](reusable-guidance-boundary.md) | [`reusable-guidance-boundary.yaml`](reusable-guidance-boundary.yaml) | 可重用 guidance 與 project evidence boundary、sanitization 與 durable location gate。 |
| [`sanitization.md`](sanitization.md) | [`sanitization.yaml`](sanitization.yaml) | secrets、本機路徑、project incident evidence 與 prompt injection 去敏 gate。 |
| [`tool-neutral-documentation.md`](tool-neutral-documentation.md) | [`tool-neutral-documentation.yaml`](tool-neutral-documentation.yaml) | Tool-neutral reusable docs、adapter isolation 與 tool-specific detail boundary gate。 |

## Mechanical Enforcement Registry（Layer 2.5）

[`enforcement-registry.yaml`](enforcement-registry.yaml) + companion [`enforcement-registry.md`](enforcement-registry.md) 是 Layer 2.5（Meta Governance / Framework Self-Audit）的 canonical binding 表。每條 rule_class 必須宣告 `coverage` ∈ `mechanical | behavioral_only | not_mechanizable | pending_implementation | research_required | deprecated`，並依 coverage 填對應必填 metadata（mechanical 要 executor symbol、behavioral_only 要雙必填 sunset_decision、pending 要 child_plan、research 要 research_questions、deprecated 要 replaced_by 或 removal_date）。

修改 `enforcement/*.yaml`、`runtime/*.yaml`、`governance/**/*.yaml`、`knowledge/runtime/routing-registry.yaml` 或 `scripts/ai-skill-cli/internal/app/hooks.go` 時，**同一 commit** 必須同步更新 registry binding。Phase 3 compile-time lint 會 hard-fail orphan rule / orphan executor / missing executor symbol / incomplete sunset_decision。

對應 meta-pattern：[`failure-patterns/rule-without-executor.md`](failure-patterns/rule-without-executor.md)。對應 plan：[`plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md`](../plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md)。

### Core Bootstrap（每個 session 必讀）

每次 session 啟動時，先讀 [`CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md)，包含 3 條核心規則：

| 順序 | 規則 | 用途 | 預估 tokens |
| --- | --- | --- | --- |
| 1 | [rule-weight.md](rule-weight.md) | 規則權重與衝突優先序 | ~300 |
| 2 | [dependency-reading.md](dependency-reading.md) | 依賴讀取、dependency read ledger、Ai-skill writeback transaction | ~400 |
| 3 | [conversation-goal-ledger.md](conversation-goal-ledger.md) | 對話目標閉環、owner/lock、parallelization mode、完成刪除 | ~100 |

**總計：~800 tokens**

### Lazy-load Rules（依條件 activate）

以下規則**不預設載入**，只在符合各自 companion YAML contract 的 `activation` 條件時才 activate；contract projection 可在 [`runtime/runtime.db`](../runtime/runtime.db) 的 `generated_surfaces` 查詢：

| 規則 | 觸發條件範例 | 預估 tokens | 優先權 |
| --- | --- | --- | --- |
| [linked-updates.md](linked-updates.md) | multi-file change、architecture update | ~800 | P1 |
| [escalation-policy.md](escalation-policy.md) | repeated failure、user contradiction、evidence conflict、source-of-truth miss | ~900 | P1 |
| [failure-learning-system.md](failure-learning-system.md) | user 指出失誤、close-loop gap | ~1200 | P1 |
| [decision-efficiency.md](decision-efficiency.md) | 多條可行路線需選擇 | ~600 | P2 |
| [prompt-cache-efficiency.md](prompt-cache-efficiency.md) | prompt cache、context loading 或 token cost layout | ~700 | P2 |
| [tool-neutral-documentation.md](tool-neutral-documentation.md) | 建立或修改可重用文件 | ~500 | P2 |
| [governance/document-sizing.md](../governance/document-sizing.md) | 文件接近拆分門檻 | ~400 | P2 |
| [document-todo-list.md](document-todo-list.md) | 文件有未完成項目 | ~300 | P2 |
| [goal-action-validation.md](goal-action-validation.md) | 重要工作單元需要驗證 | ~500 | P2 |
| [evidence-hierarchy.md](evidence-hierarchy.md) | evidence conflict、claim scope overreach、confidence decay、assumption 被當成 fact | ~900 | P1 |
| [neutral-language.md](neutral-language.md) | 撰寫或審查文件用語 | ~300 | P2 |
| [sanitization.md](sanitization.md) / [sanitization-mechanical.md](sanitization-mechanical.md) | 撰寫 feedback lesson；pre-commit shared-layer sanitization scanner | ~400 | P1 |
| [authorization-scope.md](authorization-scope.md) | 涉及授權邊界 | ~300 | P0 |
| [cross-skill-references.md](cross-skill-references.md) | 引用其他 workflow | ~400 | P2 |
| [feedback-lessons.md](feedback-lessons.md) | 撰寫或 promotion lesson | ~600 | P2 |
| [content-layering.md](content-layering.md) | 遷移或重構內容 | ~300 | P2 |
| [reusable-guidance-boundary.md](reusable-guidance-boundary.md) | 泛化 lesson 為 reusable rule | ~400 | P2 |
| [failure-patterns/failure-to-validator-closure.md](failure-patterns/failure-to-validator-closure.md) | 修復錯誤後需加入 validator 測試 | ~400 | P1 |

### Activation 判斷流程

```
1. 讀取 CORE_BOOTSTRAP.md（3 條核心規則）
2. 檢查 enforcement companion YAML contracts：
   - 目前 task 是否符合任何 contract 的 `activation` 條件？
   - 符合 → 先載入該 YAML contract，再依 `source_markdown` / `required_sources` 載入 Markdown
   - 不符合 → deferred，不載入
3. 先讀 knowledge/summaries/ 對應 summary（300-500 tokens）
4. 需要時才展開完整 source
```

### 完整規則索引

| 分類 | 檔案 | 內容摘要 | Load 策略 |
| --- | --- | --- | --- |
| 授權與範圍 | [authorization-scope.md](authorization-scope.md) | 僅在授權範圍內分析、資料邊界。 | lazy (P0) |
| 去敏與占位符 | [sanitization.md](sanitization.md) | 什麼不可寫進可重用文件、占位符約定。 | lazy (P1) |
| 中性、低爭議與語言一致性 | [neutral-language.md](neutral-language.md) | 文件標題、檔名、摘要與正文避免高風險或爭議詞；可重用文件正文預設繁體中文，必要英文保留給路徑、指令、API 欄位與專有名詞。 | lazy (P2) |
| 工具中立文件 | [tool-neutral-documentation.md](tool-neutral-documentation.md) | 可重用文件預設保持工具中立；工具專屬路徑、hook、同步與 UI 步驟放到 `ai-tools/` 或工具設定檔。 | lazy (P2) |
| 規則權重與衝突優先序 | [rule-weight.md](rule-weight.md) | 當 enforcement rules、workflow、tool adapter、使用者目標或效率規則看似衝突時，依安全/source/validation/user-goal/tool adapter/效率的權重排序處理。 | **core** |
| 決策效率 | [decision-efficiency.md](decision-efficiency.md) | 先界定未知，再依 time-to-evidence、語意距離、風險、驗證信號與 context 成本選下一步。 | lazy (P2) |
| Prompt cache efficiency | [prompt-cache-efficiency.md](prompt-cache-efficiency.md) | 組裝 context 時維持 stable prefix / semi-stable middle / volatile suffix，避免高變動內容破壞 provider prompt cache 命中率。 | lazy (P2) |
| Escalation policy | [escalation-policy.md](escalation-policy.md) | repeated failure、使用者反證、evidence conflict 或 source-of-truth miss 發生時，停止局部 patch，重讀 source-of-truth 並重建 execution graph。 | lazy (P1) |
| Evidence hierarchy | [evidence-hierarchy.md](evidence-hierarchy.md) | 依 authority、freshness、validity、scope、observability 評估 evidence；防止 assumption 當 fact、local evidence 支持 global claim、弱證據覆蓋高品質 contradiction。 | lazy (P1) |
| 失效學習系統 | [failure-learning-system.md](failure-learning-system.md) / [failure-patterns](failure-patterns/README.md) | 使用 failure taxonomy、pattern records、promotion targets 與 validation gate，將重複 agent 失效模式沉澱成可重用防呆規則；例如 source/mirror 寫入漂移。 | lazy (P1) |
| 文件 TODO | [document-todo-list.md](document-todo-list.md) | 文件有未完成、待決策、待補強或待驗證內容時，在前段放可掃描 TODO 表並連到相關章節、goal 或 issue。 | lazy (P2) |
| 目標、執行、驗證 | [goal-action-validation.md](goal-action-validation.md) | 每個重要工作單元要能反查目標、執行內容與驗證方式；純判斷題用參考來源與推論邊界。 | lazy (P2) |
| 對話目標閉環 | [conversation-goal-ledger.md](conversation-goal-ledger.md) | 使用專案本地 `.agent-goals/` 暫存帳本追蹤 active goals、優先權、parallelization mode、owner/lock 決策、plan/todo links、missing/decision/strengthen、拆解、轉移、完成驗證與刪除條件；多步驟、已有 todo、使用者要求繼續或看到 dirty files 時要先 status/init；長期 roadmap / lifecycle 狀態必須落到 durable planning 文件，不保存在 completed goal row。 | **core** |
| 依賴文件讀取 | [dependency-reading.md](dependency-reading.md) | 發現 workflow/enforcement rule/tool-specific rule/template/lesson 更新時，必須讀相關依賴文件，包含 enforcement/，並用 dependency read ledger 防止漏讀。 | **core** |
| 內容分層 | [content-layering.md](content-layering.md) | 共用規則／技巧／業務專案各自放哪；intelligence 內部 entry/solution 分層。 | lazy (P2) |
| 可重用規則與專案證據邊界 | [reusable-guidance-boundary.md](reusable-guidance-boundary.md) | 技巧只沉澱通用原因、規則與驗證；專案 incident 證據留在專案文件；指出閉環不完整時必須分析原因並強化規則。 | lazy (P2) |
| 文件大小與拆分 | [governance/document-sizing.md](../governance/document-sizing.md) | 文件拆分原則、決策流程、拆分後必做事項；token 成本模型見 decision-efficiency.md，跨專案適用見 content-layering.md。 | lazy (P2) |
| Cross-workflow references | [cross-skill-references.md](cross-skill-references.md) | 一個 workflow 需要引用另一個 workflow 的規範、模板、交接產物或驗證流程時怎麼寫。 | lazy (P2) |
| 連動更新 | [linked-updates.md](linked-updates.md) | 全庫必須連動更新規則：改一處影響多處時，相關文件必須同步更新或明確檢查。 | lazy (P1) |
| Feedback 與技巧條目 | [feedback-lessons.md](feedback-lessons.md) | **檔名規則、模板、agent 行為、索引**（唯一正文）；所有 lesson 統一放在 `feedback/history/<domain>/`，舊 `skills/<name>/feedback_history/` 已於 2026-05-13 刪除。 | lazy (P2) |
| VS Code Extension 全域設定修改 | [vscode-extension-global-state.md](../intelligence/ide/vscode-extension-global-state.md) | VS Code Extension 的全域設定儲存在 SQLite 資料庫中，修改方法與注意事項。此為可重複使用的工程智慧，非工具設定。 | lazy (P2) |


**單一真相來源：**只在本庫 **`enforcement/`** 維護共用規則正文；部署到工具或專案時優先參照中央庫或 symlink，需要離線快照時才複製整個 `enforcement/` 資料夾。
