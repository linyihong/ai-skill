# Claude Code Auto-Bootstrap

啟動時自動載入以下流程：

1. 讀 [CORE_BOOTSTRAP.md](CORE_BOOTSTRAP.md) — 3 條必讀核心規則（~800 tokens）
2. 讀 [README.md](README.md) — 超短入口，了解 OS layout
3. 依 [runtime/router/activation-rules.yaml](runtime/router/activation-rules.yaml) 決定 lazy-load rules
4. 載入 Runtime Phase 初始化（SQLite 為 agent 快速路徑）：
   - [runtime/runtime.db](runtime/runtime.db) → `phase_machine` / `phases` — 目前 phase 狀態
   - [runtime/runtime.db](runtime/runtime.db) → `obligation_ledger` / `obligations` — 本 phase 義務
   - [runtime/runtime.db](runtime/runtime.db) → `blocking_gates` / `gates` — 本 phase blocking gates
   - Runtime phase / obligation / gate / recovery 的 source 在 `runtime/phases/`、`runtime/obligations/`、`runtime/gates/`、`runtime/recovery/` YAML；修改後用 `ai-skill runtime compile` 重建 `runtime/runtime.db`。
6. 載入 Output Governance 初始化（YAML 為 source-of-truth，同時已編譯至 SQLite）：
   - [runtime/output-governance/language-policy.yaml](runtime/output-governance/language-policy.yaml) → `language_policy` 表 — 語言強制規則
   - [runtime/output-governance/output-rules.yaml](runtime/output-governance/output-rules.yaml) → `output_rules` 表 — 文件輸出規則
   - [runtime/output-governance/governance-gates.yaml](runtime/output-governance/governance-gates.yaml) → `governance_gates` 表 — 輸出品質 blocking gates
7. 先讀 `knowledge/summaries/` 對應 summary（300-500 tokens），需要時才展開全文

> **Runtime Config 已編譯至 SQLite**：所有 `runtime/**/*.yaml` 設定檔已由 compiler 編譯至 `runtime/runtime.db` 的專屬表格。Agent 可直接查 SQLite 取得結構化資料，YAML 檔案仍為 source-of-truth 供人類編輯。

詳細使用說明見 [ai-tools/agent/claude.md](ai-tools/agent/claude.md)。
