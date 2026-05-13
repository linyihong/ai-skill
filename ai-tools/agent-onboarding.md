# AI Agent 加入設定集

當一個新的 AI agent 工具要與本知識庫整合時，依以下 checklist 完成設定。每個項目都標註了參考來源，確保設定一致且可追溯。

## 必要設定

| # | 設定項目 | 強度 | 參考來源 | 說明 |
|---|---------|------|----------|------|
| 1 | 自動載入入口 | 必要 | [`CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md) | 設定工具的自動載入機制，指向 Core Bootstrap 流程。不同工具的入口機制不同（`CLAUDE.md` / `.cursor/rules/*.mdc` / Custom Instructions）。 |
| 2 | 語言偏好 | 必要 | [`shared-rules/failure-patterns/language-preference-drift.md`](../shared-rules/failure-patterns/language-preference-drift.md) | 設定軟性語言偏好（預設英文但跟隨使用者語言），避免 agent 強制使用固定語言。注意有些工具（如 Roo Code）有雙層設定（Custom Instructions + 全域 SQLite 欄位）。 |
| 3 | 對話目標閉環 | 必要 | [`shared-rules/conversation-goal-ledger.md`](../shared-rules/conversation-goal-ledger.md) | 實作 goal ledger 讀取與更新流程，讓 agent 能在中斷、context compaction、multi-agent handoff 後恢復未完成工作。每個工具的實作方式不同（hooks / 操作注意 / custom instructions）。 |
| 4 | 工具使用說明文件 | 必要 | [`ai-tools/README.md`](README.md) | 在 `ai-tools/agent/` 下建立工具使用說明文件（`<tool>.md`），只記錄該工具特有的差異，不重複中央庫內容。 |
| 5 | 驗證流程 | 建議 | [`scripts/validate-knowledge-runtime.rb`](../scripts/validate-knowledge-runtime.rb) | 設定 commit 前驗證（pre-commit hook），確保修改不破壞知識庫結構。 |
| 6 | 共用規則索引 | 建議 | [`shared-rules/README.md`](../shared-rules/README.md) | 若工具需要特殊規則，在 `shared-rules/` 中建立對應規則，並在索引中註冊。 |

## 各工具實作對照

| 設定項目 | Claude Code | Cursor | Roo Code |
|---------|------------|--------|----------|
| 自動載入入口 | `CLAUDE.md` → `CORE_BOOTSTRAP.md` | `.cursor/rules/dependency-reading.mdc`（alwaysApply） | Custom Instructions 或 `.roomodes` |
| 語言偏好 | `CLAUDE.md` 中設定軟性偏好 | `.cursor/rules/*.mdc` 中設定軟性偏好 | `.roomodes` + SQLite `language` 欄位 |
| 對話目標閉環 | 操作注意（claude.md 第 19-23 行） | 完整章節含 hooks 範本（cursor.md） | 操作注意（roo.md） |
| 工具文件 | `ai-tools/agent/claude.md` | `ai-tools/agent/cursor.md` | `ai-tools/agent/roo.md` |
| 驗證流程 | pre-commit hook（共用） | pre-commit hook（共用） | pre-commit hook（共用） |

## 新增工具的步驟

1. **確認工具類型**：是 AI Agent（CLI / IDE 內建）還是 IDE 設定？Agent 放 `ai-tools/agent/`，IDE 設定放 `ai-tools/ide/`。
2. **建立工具使用說明**：在對應目錄下建立 `<tool>.md`，只記錄該工具特有的差異。
3. **設定自動載入入口**：依工具的機制設定入口，指向 `CORE_BOOTSTRAP.md`。
4. **設定語言偏好**：依工具的設定方式，加入軟性語言偏好。
5. **實作對話目標閉環**：依工具的能力（hooks / custom instructions / 操作注意），實作 goal ledger 整合。
6. **更新 `ai-tools/README.md`**：在對應類別的表格中加入新工具的連結與用途說明。
7. **驗證**：執行 `scripts/validate-knowledge-runtime.rb` 確認無誤。

## 不建議設定的項目

以下項目不應放在工具設定中，因為它們由中央庫統一管理：

| 項目 | 應在何處 |
|------|---------|
| Core Bootstrap 流程細節 | [`CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md) |
| 共用規則本體 | [`shared-rules/`](../shared-rules/) |
| 知識庫 OS layout | [`README.md`](../README.md) |
| Goal ledger 操作流程 | [`shared-rules/conversation-goal-ledger.md`](../shared-rules/conversation-goal-ledger.md) |
| 驗證邏輯 | [`scripts/validate-knowledge-runtime.rb`](../scripts/validate-knowledge-runtime.rb) |

← [回到 AI 工具索引](README.md)
