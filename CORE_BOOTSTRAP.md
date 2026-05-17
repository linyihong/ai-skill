# Core Bootstrap

本文件是 AI-native Knowledge Operating System 的**最小必讀啟動集合**。每個 session 啟動時，Agent 只需載入這 3 條核心規則，其餘規則依任務 lazy-load。

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
3. [Runtime Phase 初始化] 讀取 runtime/phases/phase-machine.yaml（或查 SQLite phase_machine 表）：
   - 載入目前 phase、allowed_actions、forbidden_actions、blocking_gates
   - 載入 obligation-ledger.yaml（或查 SQLite obligation_ledger 表）確認本 phase 的未完成義務
   - 載入 blocking-gates.yaml（或查 SQLite blocking_gates 表）確認本 phase 的阻斷閘門
   - 查詢 runtime.db（generated_surfaces 表）作為知識更新快速路徑
   - 若 phase 為 bootstrap → 繼續往下；若為其他 phase → 先檢查 blocking gates
4. [Output Governance 初始化] 載入輸出規則（YAML 或 SQLite 皆可）：
   - 讀取 runtime/output-governance/language-policy.yaml（或查 SQLite language_policy 表）— 語言強制規則
   - 讀取 runtime/output-governance/output-rules.yaml（或查 SQLite output_rules 表）— 文件輸出規則
   - 讀取 runtime/output-governance/governance-gates.yaml（或查 SQLite governance_gates 表）— 輸出品質 blocking gates
   - 確認目前 phase 的 governance gates 狀態
5. [新專案檢查] 檢查目前專案是否已設定 Ai-skill：
   - 檢查 .roomodes 是否存在且包含 CORE_BOOTSTRAP.md 參考
   - 檢查 .cursor/rules/ 下是否有 ai-skill bootstrap 規則
   - 檢查 CLAUDE.md 是否存在且包含 CORE_BOOTSTRAP.md 參考
   - 若以上皆無 → 主動告知使用者：「此專案尚未設定 Ai-skill 知識庫。
     是否要執行初始化？(./scripts/init-new-project.sh <PROJECT_ROOT>)」
   - 若使用者同意，執行初始化腳本
6. 依任務 intent 查詢 skills-index.yaml 找到對應 skill
7. 檢查該 skill 是否有 primary_entrypoint 欄位：
   - 有 → 優先讀 primary_entrypoint 指向的新分層路徑
   - 無 → 讀 entrypoint 指向的舊路徑（向後相容）
8. 依 activation rules 決定哪些 lazy-load rules 需要載入
9. 先讀 knowledge/summaries/ 對應 summary（300-500 tokens）
10. 需要時才展開完整 source
```

> **Runtime Config 已編譯至 SQLite**：所有 `runtime/**/*.yaml` 設定檔已由 compiler 編譯至 `runtime/runtime.db` 的專屬表格（如 `phase_machine`、`obligation_ledger`、`blocking_gates`、`language_policy`、`output_rules`、`governance_gates` 等）。Agent 可直接查 SQLite 取得結構化資料，YAML 檔案仍為 source-of-truth 供人類編輯。

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
- 任務 intent 已對應到 skills-index.yaml
- Lazy-load rules 的 activation conditions 已檢查
