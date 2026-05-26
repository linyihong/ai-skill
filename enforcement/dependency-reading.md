# 依賴文件讀取鐵則

本規則適用於所有 agent 使用、修改或檢查 `enforcement/`、`workflow/`、`analysis/`、`intelligence/`、工具專用規則、模板、feedback lessons、同步腳本與根索引時。目的不是增加形式流程，而是避免 agent 只讀單一文件，卻忽略已更新的依賴規則。

## 核心規則

只要發現某個 workflow、enforcement rule、tool-specific rule、模板或 feedback lesson 已更新、將被更新、或可能影響目前任務，agent 必須讀取它的相關依賴文件後才能下結論或繼續修改。

最低讀取範圍：

| 發現或修改的項目 | 必須讀取或明確檢查 |
| --- | --- |
| 任一 `workflow/<domain>/execution-flow.md` | 先讀對應 `workflow/<domain>/README.md`、必要的 `analysis/<domain>/README.md`、`intelligence/<domain>/README.md` 或 `intelligence/engineering/<domain>/README.md`（若存在），以及 `enforcement/README.md`。不存在的檔案可標記為不適用。 |
| 任一 workflow / analysis / intelligence 子文件 | 最近的目錄 `README.md`、相關 workflow/checklist/template、`enforcement/linked-updates.md`。 |
| 任一 `enforcement/*.md` | `enforcement/README.md`、`enforcement/content-layering.md`、`enforcement/linked-updates.md`、`enforcement/rule-weight.md`（若涉及規則衝突、優先序或 default bootstrap）、`enforcement/reusable-guidance-boundary.md`（若涉及 reusable guidance / incident / feedback）、受影響 workflow 的 `execution-flow.md` 或模板。 |
| `enforcement/tool-neutral-documentation.md` | `enforcement/README.md`、`enforcement/content-layering.md`、`enforcement/linked-updates.md`、根 `README.md`、相關 `workflow/`、`analysis/`、`intelligence/` 入口、`ai-tools/README.md` 與受影響工具文件。 |
| `enforcement/rule-weight.md` | `enforcement/README.md`、`enforcement/content-layering.md`、`enforcement/linked-updates.md`、`enforcement/dependency-reading.md`、`enforcement/decision-efficiency.md`、`enforcement/goal-action-validation.md`、工具專用 always-apply agent rule、`ai-tools/` bootstrap 清單。 |
| `enforcement/decision-efficiency.md` | `enforcement/README.md`、`enforcement/content-layering.md`、`enforcement/linked-updates.md`、`enforcement/dependency-reading.md`、`governance/document-sizing.md`、相關 workflow / analysis / intelligence README（若該 route 有決策路由或 context-loading 指引）。 |
| `enforcement/escalation-policy.md` | `enforcement/README.md`、`enforcement/linked-updates.md`、`enforcement/failure-learning-system.md`、`enforcement/dependency-reading.md`、`runtime/README.md`、`runtime/runtime.db`（若要接 runtime guard）、相關 workflow primary source（若 domain-specific mismatch hook 受影響）。 |
| `enforcement/failure-learning-system.md` 或 `enforcement/failure-patterns/` | `enforcement/README.md`、`enforcement/content-layering.md`、`enforcement/linked-updates.md`、`enforcement/reusable-guidance-boundary.md`、`../../feedback/feedback-lessons.md`、`enforcement/goal-action-validation.md`、`enforcement/dependency-reading.md`、相關 tool 文件與被補強的 enforcement rule / workflow。 |
| 任一 tool adapter | 對應 workflow / tool 文件、`enforcement/tool-neutral-documentation.md`、`ai-tools/<tool>.md`（若存在）、`enforcement/linked-updates.md`。 |
| `enforcement/document-todo-list.md` | `enforcement/README.md`、`enforcement/content-layering.md`、`enforcement/linked-updates.md`、`enforcement/conversation-goal-ledger.md`、相關模板與 documentation/checklist 文件。 |
| `enforcement/conversation-goal-ledger.md` | `enforcement/README.md`、`enforcement/content-layering.md`、`enforcement/linked-updates.md`、`scripts/README.md`、相關 goal helper script、`ai-tools/` 中受影響工具文件；若同時改 tool-specific hook/rule，讀對應 hook/rule 文件。 |
| 任一工具專用規則檔 | 對應的 enforcement rule 正文、`enforcement/README.md`、受影響工具文件，以及受影響的 workflow 入口。 |
| 任一 template | 模板目錄 `README.md`、引用該模板的 workflow/documentation/checklist、`enforcement/linked-updates.md`。 |
| 任一 feedback lesson | 該分類 `README.md`、`feedback/history/<domain>/README.md`、`../../feedback/feedback-lessons.md`，以及 promotion target。 |
| 開始執行任一 `plans/active/*.md` plan | `plans/README.md` 的 Architecture Compatibility Preflight、該 plan 的 Phase 0 / candidate files、相關 layer README、相關 runtime/compiler/metadata/workflow source、`enforcement/linked-updates.md`。 |

## Agent 行為

1. 先讀 `enforcement/README.md` 的 **Default Bootstrap**，並載入 bootstrap 表列出的最小必讀集合；再依任務讀相關 enforcement rule 全文。
2. 若任務碰到 workflow，讀該 workflow 入口與依賴文件；不要只依賴 `description`、單檔或單一段落。
3. 若看到文件有 cross-link、promotion target、required linked updates、template reference、feedback index，或 reusable guidance / project incident 邊界，就循連結讀到任務所需的規則載入完成。
4. 若依賴文件不存在，記錄為 `not applicable`；若存在但未讀，不可宣稱已完成檢查。
5. 回覆或提交前，說明依賴讀取與連動更新的驗證方式。
6. 完成 `git commit`、`git push` 與必要的 tool sync 後，必須重新讀取本次更新過的 skill/enforcement-rule 入口文件與主要依賴文件，確認目前 agent context 已載入最新版；不可只依賴提交前讀過的內容。Reference-only 策略不需要 tool sync。
7. 最終回覆前必須執行 `git status --short --branch`。若 `Ai-skill` repo 仍有 modified/untracked/staged changes，或 branch 仍 ahead/behind remote，不得回覆「已完成」；必須先完成驗證、sync、commit、push、讀回，或明確說明被什麼阻塞。
8. 同一輪對話可以依獨立邏輯單元建立多個 commit；不要求每個 commit 後立刻 push。但在使用者表示任務完成、切換新任務、agent 輸出最終 summary 或任何會讓本輪 writeback 交易關閉的回覆前，必須完成 `git push`、push 後讀回與 clean status，並確認 `git log origin/<branch>..HEAD` 為空。
9. 若使用者未明確要求 push / merge，而更新後發現 `Ai-skill` 有尚未推送、尚未合併、ahead/behind、或其他 pending commit 狀態，最終回覆必須主動提醒使用者目前狀態與下一步（例如需要 push、pull/rebase、或處理既有 dirty changes），不可讓使用者以為規則已完全進入遠端主線。

### Plan Execution Preflight

開始執行 `plans/active/*.md` 的任何 implementation phase 前，必須先依 [`plans/README.md`](../plans/README.md#plan-執行前架構相容性檢查architecture-compatibility-preflight) 完成 Architecture Compatibility Preflight。此檢查用來確認 candidate files、source-of-truth、compiler / generated surfaces、layer responsibility 與現行架構一致。

若 preflight 發現衝突，agent 必須先暫停執行，更新 plan 或向使用者確認；不得直接套用 stale plan。若涉及 `runtime.db`、generated reports、SQLite index 或 compiler outputs，必須查詢或驗證 source 變更已進入 generated surface。

### Dependency Read Ledger

當使用者要求「重新讀 workflow」、指出「enforcement rules / shared workflow 是否漏讀」、或 agent 自己發現某個 workflow/rule 已更新時，必須在繼續業務專案前建立一個簡短的 dependency read ledger（可在回覆、todo 或工作筆記中呈現），至少列：

| 欄位 | 必填內容 |
| --- | --- |
| Trigger | 例如 `workflow/<domain>/execution-flow.md changed`、user asked to reload workflow、enforcement rule changed。 |
| Required set | 依本檔「最低讀取範圍」列出應讀文件。 |
| Read | 實際已讀文件。 |
| Not applicable | 不存在的檔案，例如該 workflow 沒有 `CHECKLIST.md`；必須明寫，不可假裝已讀。 |
| Deferred / blocked | 因權限、缺檔、衝突或使用者決策而未讀的項目。 |
| Validation | 連動更新檢查、diff review、sync、commit/push/readback 或純判斷的參考來源。 |

若 ledger 顯示最低讀取範圍仍有缺口，agent 不得宣稱「已按更新後 skill 執行」或長時間切回專案分析；必須先補讀、標 `not applicable`，或明確向使用者說明阻塞。

### Source-of-truth Miss Escalation

若 agent 已開始執行，但發現重要操作前未讀 canonical workflow、owner docs、UI map、API contract、runtime source 或 routing registry，這不只是一般 dependency miss。依 [`escalation-policy.md`](escalation-policy.md) 將它視為 `source-of-truth-miss`：

1. 停止目前 patch / automation / guess。
2. 在 dependency read ledger 中標明缺少的 source。
3. 讀取或標記 `not applicable` / `source missing`。
4. 重建 execution graph 後才恢復執行。

若使用者指出「你沒看文件」「不是這個流程」「你又在猜」，同樣依 escalation policy 進入 recovery frame，不只補一個檔案讀取。

## Default Bootstrap Boundary

Default bootstrap 是每次 agent/session 開始時的最低共用上下文，不等於所有規則全文都已讀。

啟動時先讀 [`README.md`](README.md) 的 Default Bootstrap，至少載入：

### primary_entrypoint 優先規則

查詢 [`knowledge/runtime/routing-registry.yaml`](../knowledge/runtime/routing-registry.yaml) 找到對應 route 後，**必須先檢查是否有 `primary_source` 欄位**：

- **有 `primary_entrypoint`** → 優先讀該路徑（指向新分層：`workflow/<domain>/`、`analysis/<domain>/`、`intelligence/<domain>/`）
- **無 `primary_entrypoint`** → 讀 `entrypoint` 指向的舊路徑（向後相容）
- `entrypoint` 保留給 tool adapter 使用，AI 不應優先讀取

### 必讀規則

- `enforcement/README.md`
- `enforcement/dependency-reading.md`
- `enforcement/linked-updates.md`
- `enforcement/conversation-goal-ledger.md`
- `enforcement/tool-neutral-documentation.md`
- `enforcement/rule-weight.md`
- `enforcement/decision-efficiency.md`
- `enforcement/failure-learning-system.md`
- `enforcement/document-todo-list.md`
- `governance/document-sizing.md`
- `enforcement/goal-action-validation.md`
- `enforcement/neutral-language.md`

之後按任務補讀：

- **Workflow 編排**（activation **#27** 或 registry `route.workflow.*` 命中）：**Blocking** — 先 [`routing-philosophy.md`](../governance/lifecycle/routing-philosophy.md) → 比對 `activation_triggers` → 必要時 [`workflow-routing.md`](../workflow/workflow-routing.md) 歧義裁決 → 選定 route 的 `primary_source` 與 `required_dependencies`；未完成不得寫可觀察產品行為變更。
- 寫 feedback / lesson 時補讀 `feedback-lessons.md`、`reusable-guidance-boundary.md`、`sanitization.md`，必要時 `authorization-scope.md`。
- 引用其他 skill 時補讀 `cross-skill-references.md`。
- 修改 skill 時補讀該 skill 的 README / WORKFLOW / TOOLS / DOCUMENTATION / CHECKLIST / FEEDBACK / feedback index。

工具可用 hook、always-apply rule 或固定提示詞自動提醒 bootstrap；但工具提醒不取代實際讀取與 dependency read ledger。

## Ai-skill Writeback Transaction Guard → Runtime Transaction Machine

> **已升級為 SQLite-canonical runtime state machine**：快速路徑與 committed canonical config 都在 [`runtime/runtime.db`](../runtime/runtime.db)。完整 documents 存於 `runtime_config_documents`，transaction / phase / gate tables 是 projection；不要提交 `runtime/**/*.yaml` mirror。

本節原有的 procedural writeback transaction 邏輯（7 個狀態、entry conditions、blocking gates、transition rules）已從 prose 規則升級為 SQLite-canonical runtime state machine。Agent 應查詢 `runtime/runtime.db` 作為 authoritative runtime reference；需要修改定義時更新 `runtime_config_documents`，然後執行 `ai-skill runtime compile` refresh projections。

### 本節保留的邊界規則

下列規則不屬於 transaction state machine 的範疇，仍保留在此：

1. **Canonical source first**：修改 Ai-skill 規則時，必須先在 canonical `<AI_SKILL_REPO>` 修改、diff、commit/push；只有 source 變更成立後，才依工具同步策略更新 runtime copy。若 agent 已誤改 runtime / mirror copy，必須立刻停止擴大修改，定位 canonical repo，將同等規則補回 source。
2. **Lock 檢查**：若其他 agent / user 正在操作（active close-loop lock），不得自動 commit、push 或清理其變更，只能回報 lock owner、狀態與下一步。
3. **去敏檢查**：依 [`sanitization.md`](sanitization.md) 檢查所有新增/修改的可重用文件，不得包含本機絕對路徑（改用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>` 占位符）、使用者名稱、私有工作目錄、clone 位置、secrets、raw tokens、私人 host、個資或 project-specific evidence。
4. **Tool sync 邊界**：若本輪明確使用或更新 tool mirror / symlink / copy snapshot，必要的 tool sync 已執行；reference-only 則記為不適用。
5. **Commit / push cadence**：同一輪可以多次 commit，不必每 commit 後立即 push；但任務完成、切換任務或最終回覆前必須 push、讀回，並確認本地 branch 沒有未推送 commit。
6. **最終狀態確認**：最後一次 `git status --short --branch` 顯示 clean，且 branch 沒有 ahead/behind。

### 交易關閉條件（由 runtime tables 定義）

Transaction 的完整關閉條件已由 `runtime/runtime.db` 的 transaction tables 定義（canonical document：`runtime_config_documents`）。Agent 應直接查詢 SQLite，而非依賴本節的 prose 摘要。

若 transaction 未關閉，agent 不得把注意力長時間切回業務專案，也不得把「已更新 skill」當作完成。可以繼續工作的唯一例外是：使用者明確要求暫停 Ai-skill close-loop；此時必須說明目前 dirty/ahead/behind/unmerged 狀態與下一步。

### State-based Enforcement（狀態化強制規則）

下列邊界規則已對應到 runtime state machine 的 phase/gate 定義，agent 應優先查詢 runtime SQLite tables：

```yaml
# State-based enforcement mapping for dependency-reading.md
# 這些規則已由 runtime state machine 管理，agent 不應再以 prose 方式逐條檢查。
state_based_enforcement:
  version: v1
  status: active
  owner_layer: enforcement/dependency-reading
  description: >
    將 Ai-skill writeback transaction 的邊界規則對應到 runtime state machine。
    Agent 應優先查詢 runtime/runtime.db；需要修改時更新 runtime_config_documents。

  # Canonical source first → 由 runtime/runtime.db 的 phase.execution 管理
  # allowed_actions 包含 write_file/apply_diff/execute_command
  # forbidden_actions 包含 commit/push/finalize
  - rule: canonical_source_first
    phase: phase.execution
    enforcement: |
      Agent 必須在 canonical <AI_SKILL_REPO> 中修改檔案。
      若誤改 runtime/mirror copy，必須停止並定位 canonical repo。
    runtime_ref: runtime/runtime.db
    runtime_section: "phase.execution.allowed_actions"

  # Lock 檢查 → 由 runtime/runtime.db 的 tx.rule.lock_check 管理
  - rule: lock_check
    phase: phase.execution
    enforcement: |
      開始 transaction 前必須檢查是否有 active lock。
      若有其他 agent/user 正在操作，不得自動 commit/push。
    runtime_ref: runtime/runtime.db
    runtime_section: "tx.rule.lock_check"

  # 去敏檢查 → 由 runtime/runtime.db 的 gate.validation.artifacts_complete 管理
  - rule: sanitization_check
    phase: phase.validation
    enforcement: |
      Commit 前必須執行去敏檢查。
      依 sanitization.md 檢查所有新增/修改的可重用文件。
    runtime_ref: runtime/runtime.db
    runtime_section: "gate.validation.artifacts_complete"

  # Tool sync 邊界 → 由 runtime/runtime.db 的 tx.rule.linked_updates_check 管理
  - rule: tool_sync_boundary
    phase: phase.validation
    enforcement: |
      若本輪使用或更新 tool mirror/symlink/copy snapshot，
      必要的 tool sync 已執行；reference-only 記為不適用。
    runtime_ref: runtime/runtime.db
    runtime_section: "tx.rule.linked_updates_check"

  # 最終狀態確認 → 由 runtime/runtime.db 的 verified → closed 管理
  - rule: final_status_check
    phase: phase.readback
    enforcement: |
      最後一次 git status --short --branch 顯示 clean，
      branch 沒有 ahead/behind。
    runtime_ref: runtime/runtime.db
    runtime_section: "state.verified → state.closed"
```

## Conversation Goal Ledger Boundary

[`conversation-goal-ledger.md`](conversation-goal-ledger.md) 管的是使用者對話目標是否完成；本檔的 Ai-skill writeback transaction 管的是本知識庫改動是否完成 sync / commit / push / reread / clean status。兩者不可互相取代：

- `.agent-goals/` 是專案本地暫存狀態，不應 commit。
- 當目前任務會跨多個 tool call、已建立 TodoWrite、使用者要求「繼續」前一個多步驟任務，或 agent 已看到 active project 有 modified / staged / untracked files 時，必須先檢查或初始化 `.agent-goals/`，再長時間繼續專案工作。
- Ai-skill writeback transaction 是本 repository 的 git 閉環，必須 commit / push。
- 當使用者目標是「修改 Ai-skill 規則或 skill」時，agent 可能需要同時維護 `.agent-goals/` 中的 user goal，並完成本檔要求的 Ai-skill writeback transaction。
- 不可因為 `.agent-goals/` 目標刪除了，就跳過本庫的 diff review、linked updates、必要 tool sync、commit、push、讀回與 clean status；reference-only 時 tool sync 應明確標為不適用。

## Failure Learning Boundary

[`failure-learning-system.md`](failure-learning-system.md) 管的是「已發現失效模式如何變成可重用防呆規則」。當使用者指出 agent 反覆失誤、寫錯 source/mirror、漏讀依賴、漏做驗證、忘記 goal、混入專案細節或閉環不完整時：

- 先依本檔與相關 enforcement rule 補救當前 close-loop。
- 再用 `failure-learning-system.md` 分類失效模式、選擇 promotion target，必要時新增或更新 `enforcement/failure-patterns/`。
- 若失效模式是 skill-specific，將 lesson 放到該 skill 的 `feedback_history/`；若是 cross-skill 或全庫行為，放到 `enforcement/failure-patterns/` 或對應 enforcement rule。
- 不可把 project incident 的 raw evidence、私人路徑、host、token 或一次性細節寫進 failure pattern；先依 `reusable-guidance-boundary.md` 抽象化。

## Tool-Neutral Documentation Boundary

[`tool-neutral-documentation.md`](tool-neutral-documentation.md) 要求可重用文件預設保持工具中立。新增或修改 enforcement rule、workflow、template、README、lesson 時：

- 先寫通用 agent/tool 行為，不把單一 IDE、CLI 或 agent 產品寫成預設需求。
- 工具專屬路徑、hook、同步命令、UI 操作、reload 步驟，放在 `ai-tools/<tool>.md`、工具設定檔或工具專用腳本文件。
- 若 generic rule 需要提同步，用「configured tool sync」等中立詞，再連到 `ai-tools/` 取得具體工具做法。
- 若某 workflow 需要工具差異，使用 Strategy-style adapter：核心 workflow 保存通用契約，工具差異放在 `ai-tools/` 或 workflow 明確連結的 adapter 說明。
- Commit/push 後讀回時，也要確認沒有把工具專屬段落誤放進 root README、workflow README、enforcement rule index 或 reusable lesson。

## Document TODO List Boundary

[`document-todo-list.md`](document-todo-list.md) 管的是文件本身的未完成、待決策、待補強與待驗證項目；[`conversation-goal-ledger.md`](conversation-goal-ledger.md) 管的是跨文件、跨工具或跨對話的使用者目標。兩者需要互相連結：

- 修改文件時若留下未完成內容，應在文件前段加入或更新 `Document TODO`。
- 如果 TODO 變成跨文件或 user-facing 目標，應連到 `.agent-goals/` goal。
- 不能因為 goal ledger 已記錄，就把文件內明顯未完成的 TODO 藏在對話裡。
- Commit/push 前若宣稱文件完成，應檢查該文件的 TODO 表沒有未處理的 `pending`、`blocked` 或 `needs-validation` 項目。

## Ai-skill 回寫完成門檻

只要 agent 在本庫回寫任何 `enforcement/`、`workflow/`、`analysis/`、`intelligence/`、工具專用規則、模板、feedback lessons、README 或同步腳本，最終回覆前必須完成整個更新閉環：

1. `git status --short --branch` 檢查變更。
2. `git diff` 檢查將提交的內容，不得包含 secrets、raw tokens、私人 host、個資或本機絕對路徑。
3. 執行適用的 lints / Markdown link check / required linked updates 檢查。
4. 若本輪明確使用或更新本機工具 mirror / symlink / copy snapshot，執行對應 tool sync；reference-only 只需確認 `<AI_SKILL_REPO>` 可讀，不跑同步。
5. 若有多個 owner group，優先使用 `ai-skill close-loop --commit` 分組提交；若手動提交，仍需按 enforcement、scripts、workflow / analysis / intelligence owner 分開提交，避免把不相干內容混成一包。
6. `git add` 相關檔案。
7. **若 runtime SQLite canonical config或 compiler 規則有變更，確認 `runtime.db` 已 `git add`**（pre-commit hook 會自動處理，但手動 commit 時需自行確認）。
8. `git commit`；同一輪可為後續獨立邏輯單元重複第 1-8 步建立多個 commit。
9. 在本輪 writeback 交易關閉前執行 `git push`。不要求每個 commit 立刻 push，但不可在最終回覆、任務切換或完成宣告前留下 `origin/<branch>..HEAD` 非空。
10. Push 後重新讀取更新過的入口與主要依賴文件。
11. 再跑 `git status --short --branch` 與 `git log origin/<branch>..HEAD`，必須看到沒有未提交變更、branch 不再 ahead/behind remote，且沒有未推送 commit。

若第 10 步不乾淨，agent 必須回到第 1 步處理剩餘變更。不可在 dirty tree 或未 push 狀態下回覆「完成」。若使用者沒有授權 push 或 merge，必須明確提醒「本地已提交但尚未推送 / 合併」以及需要使用者決定下一步。

## Commit / Push 後讀回 Gate

當本庫變更已完成 `git commit`、`git push`，且改動涉及 `enforcement/`、`workflow/`、`analysis/`、`intelligence/`、工具專用規則、模板或 feedback lessons 時，agent 必須在最終回覆前做一次讀回：

| 更新類型 | Commit / push 後必須重新讀取 |
| --- | --- |
| `enforcement/` | 更新過的 enforcement rule、`enforcement/README.md`、`enforcement/linked-updates.md`；若有工具專用規則，也讀對應工具規則檔。 |
| `workflow/` / `analysis/` / `intelligence/` | 更新過的入口 README、execution-flow / analysis source / intelligence source，以及本次更新過的 documentation / checklist / template / feedback index。 |
| 工具專用規則 | 更新過的工具規則檔，以及對應的 enforcement rule 正文。 |
| template 或 feedback lesson | 更新過的 template/lesson、索引 README、promotion target 或引用它的 workflow/documentation。 |

讀回目的：

- 確認提交後工作樹、必要 tool sync 與 agent context 沒有停在舊版本；reference-only 時確認已讀回 canonical source。
- 讓同一輪最終回覆能基於最新規則，而不是 commit 前暫存理解。
- 若下一個 agent 接手，最終回覆要明確說明已讀回哪些入口文件；若未能讀回，必須列為未完成驗證。

## 與連動更新的關係

本規則是「先讀依賴」；[`linked-updates.md`](linked-updates.md) 是「讀完後該同步更新或明確檢查哪些文件」。兩者都必須遵守：

- 沒有讀依賴，就不能可靠判斷是否需要連動更新。
- 已讀依賴但發現需要同步，就必須依 `linked-updates.md` 更新或說明無需更新的理由。
- 若改動會影響正在使用的本機工具 mirror / symlink / copy snapshot，必須執行對應 tool sync；reference-only 時不要自動同步。

## 驗證

每次套用本規則時，至少要能回答：

| 欄位 | 必填內容 |
| --- | --- |
| 目標 | 這次要確認哪個 skill/rule/template 的依賴沒有漏讀。 |
| 執行 | 實際讀取或明確檢查了哪些依賴文件。 |
| 驗證 | `git diff`、Markdown link check、lints、required linked updates 檢查、必要 tool sync，或說明哪些文件不存在 / reference-only 所以不適用。 |

← [回到共用規則索引](README.md)
