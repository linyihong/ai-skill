# Core Bootstrap

本文件是本系統的**最小必讀啟動集合**。本系統已升級為以 runtime state machine、routing registry、governance translation 與 validation gates 驅動的認知執行系統。每個 session 啟動時，Agent 只需載入這 3 條核心規則，其餘規則依任務 lazy-load。

## 必讀規則

| 順序 | 規則 | 用途 | 預估 tokens |
| --- | --- | --- | --- |
| 1 | [`enforcement/rule-weight.md`](enforcement/rule-weight.md) | 規則衝突時如何判斷優先序（P0/P1/P2/P3） | ~300 |
| 2 | [`enforcement/dependency-reading.md`](enforcement/dependency-reading.md) | 依賴文件讀取鐵則與 Ai-skill writeback transaction | ~400 |
| 3 | [`enforcement/conversation-goal-ledger.md`](enforcement/conversation-goal-ledger.md) | 對話目標閉環與 `.agent-goals/` 使用方式 | ~100 |

**總計：~800 tokens**

## 啟動流程

```
1. 讀取 CORE_BOOTSTRAP.md（本檔）
2. 讀取 README.md（超短入口，了解 OS layout）
3. [Runtime Phase 初始化] 查詢 [`runtime/runtime.db`](runtime/runtime.db)（SQLite）：
   - `phase_machine` / `phases`：目前 phase、allowed_actions、forbidden_actions
   - `obligation_ledger` / `obligations`：本 phase 的未完成義務
   - `blocking_gates` / `gates`：本 phase 的阻斷閘門
   - `generated_surfaces`：知識更新與 governance surface 快速路徑
   - Runtime config 的 committed canonical copy 是 `runtime/runtime.db`；不要保留 `runtime/runtime.db canonical documents` mirror。若需修改 phase / obligation / gate / recovery 定義，更新 SQLite canonical config document，然後用 `ai-skill runtime compile` refresh projections。
   - 若 phase 為 bootstrap → 繼續往下；若為其他 phase → 先檢查 blocking gates
4. [Output Governance 初始化] 從 SQLite 載入輸出規則：
   - `runtime/runtime.db → language_policy` — 語言強制規則
   - `runtime/runtime.db → output_rules` — 文件輸出規則
   - `runtime/runtime.db → governance_gates` — 輸出品質 blocking gates
   - 確認目前 phase 的 governance gates 狀態
5. [新專案檢查] 檢查目前專案是否已設定 Ai-skill：
   - 檢查 .roomodes 是否存在且包含 CORE_BOOTSTRAP.md 參考
   - 檢查 .cursor/rules/ 下是否有 ai-skill bootstrap 規則
   - 檢查 CLAUDE.md 是否存在且包含 CORE_BOOTSTRAP.md 參考
   - 若以上皆無 → 主動告知使用者：「此專案尚未設定 Ai-skill 知識庫。
     是否要執行初始化？(`ai-skill init-project --project <PROJECT_ROOT>`)」
   - 若使用者同意，執行初始化腳本
6. 依 activation rules 決定哪些 lazy-load rules 需要載入
7. 先讀 knowledge/summaries/ 對應 summary（300-500 tokens）
8. 需要時才展開完整 source
```

> **Runtime Config 以 SQLite 為 canonical**：committed runtime config 只保留在 `runtime/runtime.db`，並透過 `runtime_config_documents` 保存完整 canonical documents、透過專屬表格（如 `phase_machine`、`obligation_ledger`、`blocking_gates`、`language_policy`、`output_rules`、`governance_gates` 等）提供查詢。不要提交 `runtime/runtime.db canonical documents` mirror。

> **Governance Translation 已升級**：可重用判斷智慧優先放在 `intelligence/`，AI runtime gate 放在 `governance/ai-runtime-governance/`，具體操作順序放在 `workflow/`，可機讀或可驗證狀態再進 `runtime/` / `validation/`。遇到任務分層或治理疑義時，先讀 [`governance/ai-runtime-governance/README.md`](governance/ai-runtime-governance/README.md)。

## 與舊 Default Bootstrap 的關係

舊 Default Bootstrap（12 條規則）已拆分為：

- **Core Bootstrap**（本檔）：3 條必讀規則，每個 session 載入
- **Lazy-load rules**：9 條規則，只在特定情境 activate（詳見 [`enforcement/README.md`](enforcement/README.md) 的 activation model）

## 不變的原則

- `reference-first` 仍是預設：agent 直接讀本 repository
- `rule-weight` 的 P0/P1/P2/P3 權重體系不變
- `dependency-reading` 的 dependency read ledger 與 writeback transaction 不變
- `conversation-goal-ledger` 的 `.agent-goals/` 使用方式不變

## 驗證

- Core Bootstrap 三條規則已讀
- 新專案檢查已完成（若為新專案，已詢問使用者是否初始化）
- Lazy-load rules 的 activation conditions 已檢查
