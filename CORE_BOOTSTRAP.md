# Core Bootstrap（companion）

> **本檔為 companion markdown**。Canonical executable contract 在 [`runtime/core-bootstrap.yaml`](runtime/core-bootstrap.yaml)（projected to `runtime/runtime.db generated_surfaces[runtime.core_bootstrap.contract]`）。Agent **必須** query SQLite 或執行 `ai-skill runtime obligations` 取得 machine-readable `required_reads` / `per_session_obligations` / `per_turn_obligations` / `per_commit_obligations` / `phase_state_init`。本 .md 只保留人類讀的 rationale + 必讀 3 條規則表 + lazy-load 9 條規則敘述。Format examples / enum values / template details **全部在 YAML**。
>
> 修改 obligations 規則的順序：先改 `runtime/core-bootstrap.yaml` → `ai-skill runtime compile + refresh` → 本檔同步說明（companion only）。
>
> **Resume / continuation session 同樣須走完 bootstrap** — summary 的「Resume directly」是對話 framing，**不豁免** runtime / governance bootstrap。

## 必讀規則（lazy-load 入口）

| 順序 | 規則 | 用途 | 預估 tokens |
| --- | --- | --- | --- |
| 1 | [`enforcement/rule-weight.md`](enforcement/rule-weight.md) | 規則衝突時如何判斷優先序（P0/P1/P2/P3） | ~300 |
| 2 | [`enforcement/dependency-reading.md`](enforcement/dependency-reading.md) | 依賴文件讀取鐵則與 Ai-skill writeback transaction | ~400 |
| 3 | [`enforcement/conversation-goal-ledger.md`](enforcement/conversation-goal-ledger.md) | 對話目標閉環與 `.agent-goals/` 使用方式 | ~100 |

## 啟動序列

1. 讀本檔（companion） + 執行 `ai-skill runtime receipt` / `ai-skill runtime obligations` 取得 Bootstrap Receipt 與 obligations machine-readable list（不要臨時拼 SQLite schema）
2. 讀 [`README.md`](README.md) — OS layout
3. 讀 [`runtime/runtime.db`](runtime/runtime.db) `phase_machine` / `obligations` / `gates` / `language_policy` / `output_rules` / `governance_gates`
4. **新專案檢查**：若無 `CLAUDE.md` / `.cursor/rules/*` / `.roomodes` / `AGENTS.md`，主動詢問是否執行 `ai-skill init-project`
5. 依 `runtime/cognitive-modes-discovery.yaml` 14 signals 解析 cognitive mode；先讀 `knowledge/summaries/` 對應 summary（300-500 tokens）再展開全文

## Obligations（query YAML for format / examples / enum）

- **Bootstrap Receipt**（first-turn）— format / example / fields 見 `runtime/core-bootstrap.yaml` §`per_session_obligations[obligation.bootstrap.receipt]`；未輸出即執行非-Read 工具命中 [`gate.bootstrap.receipt_present`](enforcement/failure-patterns/bootstrap-bypass-on-resume.md)
- **Cognitive Mode 報告**（final close-out, v2）— 兩種形式：
  - **Compact**（全 6 維預設）：`Cognitive: NORMAL·SUMMARY_FIRST·STANDARD·NONE / V:CHECKLIST / Cost:LOW / Sig:<signal>`
  - **Full table**（任一維度非預設 或 高風險 mode）：`### Cognitive Mode 報告` + 6-dim markdown table + `activation_reason`
  - 完整 format / enum / template 見 `runtime/core-bootstrap.yaml` §`per_turn_obligations[obligation.cognitive.mode_report]` + [`models/cognitive-modes/README.md`](models/cognitive-modes/README.md)；commit 階段由 `commit-msg` hook 機械強制（[ADR-008](constitution/ADR-008-runtime-cognitive-modes.md)）。支援 session stop / final-response hook 的工具，應在 `init-project` 或工具 adapter 安裝 close-out check，避免 chat final response 漏報。
- **Feedback / Learning Report**（final close-out）— 每個 final user-facing response 必須回報 learning disposition：
  - **Compact**：`FeedbackDecision` + `RepoContext` + `Writeback`，`FeedbackDecision: NEEDED` 時加 `Target`
  - **Full table**：`### Feedback / Learning Report` + markdown table
  - 完整 format / enum / fixed order / required field combinations 見 `runtime/core-bootstrap.yaml` §`per_turn_obligations[obligation.feedback.learning_report]`。此 report 只回答本輪是否有值得沉澱的 reusable learning，以及 repo/writeback 可見狀態；不追蹤 knowledge classification、promotion lifecycle、economics、memory、telemetry 或 linked update completion。Commit message 不要求此 report。
- **Close-loop 終局檢查**（per-task）— 在回報任務完成或結束對話前，**必須**確認 `git status` 為 clean 且所有 commit 已推送（`git log origin/<branch>..HEAD` 為空）。若 push 未獲授權，必須明確說明 pending 狀態。見 `enforcement/linked-updates.yaml` §`gate.linked_updates.writeback_closed`。
- **Contextual activations**（依情境載入）— `runtime/core-bootstrap.yaml` §`contextual_activations` 定義輕量觸發，例如 Markdown 變大或混多主題時載入 [`governance/document-sizing.md`](governance/document-sizing.md)。
- **Per-commit validators**（11 個 enumerated in YAML）— `ai-skill runtime obligations` 列當前 active list；`ai-skill runtime receipt` 輸出 canonical Bootstrap Receipt；validator dispatch via registry in `scripts/ai-skill-cli/internal/app/hooks.go`

## Lazy-load rules（9 條依情境 activate）

舊 Default Bootstrap 12 條已拆為：core（本檔上方 3 條）+ contextual / lazy-load activations。機器可讀觸發先看 `runtime/core-bootstrap.yaml`；完整 lazy-load 說明見 [`enforcement/README.md`](enforcement/README.md) 的 activation model + [`runtime/runtime.db`](runtime/runtime.db) `language_policy` / `output_rules` / `governance_gates` tables。

## 不變的原則

- `reference-first` 仍是預設：agent 直接讀本 repository
- `rule-weight` 的 P0/P1/P2/P3 權重體系不變
- `dependency-reading` 的 dependency read ledger 與 writeback transaction 不變
- `conversation-goal-ledger` 的 `.agent-goals/` 使用方式不變

## 驗證

- 3 條必讀規則已讀（或 lazy-load activation 已 trigger）
- Bootstrap Receipt 已在 first user-facing message 輸出（query YAML 取格式）
- Cognitive Mode 報告已附在 final response / session close-out（query YAML 取格式）
- Feedback / Learning Report 已附在 final response / session close-out（query YAML 取格式）
- 新專案檢查已完成（若為新專案，已詢問使用者是否初始化）
