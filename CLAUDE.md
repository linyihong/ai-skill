# Claude Code Auto-Bootstrap

啟動時自動載入以下流程：

1. 讀 [CORE_BOOTSTRAP.md](CORE_BOOTSTRAP.md) — 3 條必讀核心規則（~800 tokens）
2. 讀 [README.md](README.md) — 超短入口，了解 OS layout
3. 依 [runtime/runtime.db](runtime/runtime.db) 決定 lazy-load rules
4. 載入 Runtime Phase 初始化（SQLite 為 agent 快速路徑）：
   - [runtime/runtime.db](runtime/runtime.db) → `phase_machine` / `phases` — 目前 phase 狀態
   - [runtime/runtime.db](runtime/runtime.db) → `obligation_ledger` / `obligations` — 本 phase 義務
   - [runtime/runtime.db](runtime/runtime.db) → `blocking_gates` / `gates` — 本 phase blocking gates
   - Runtime phase / obligation / gate / recovery 的 committed canonical copy 在 `runtime/runtime.db` 的 `runtime_config_documents` 與 projection tables；不要保留 `runtime/**/*.yaml` mirror。
6. 載入 Output Governance 初始化（SQLite 為 source-of-truth）：
   - [runtime/runtime.db](runtime/runtime.db) → `language_policy` 表 — 語言強制規則
   - [runtime/runtime.db](runtime/runtime.db) → `output_rules` 表 — 文件輸出規則
   - [runtime/runtime.db](runtime/runtime.db) → `governance_gates` 表 — 輸出品質 blocking gates
7. 先讀 `knowledge/summaries/` 對應 summary（300-500 tokens），需要時才展開全文

> **Runtime Config 以 SQLite 為 canonical**：committed runtime config 只保留在 `runtime/runtime.db`，由 `runtime_config_documents` 保存完整 canonical documents。Agent 直接查 SQLite；不要提交 `runtime/**/*.yaml` mirror。

詳細使用說明見 [ai-tools/agent/claude.md](ai-tools/agent/claude.md)。
