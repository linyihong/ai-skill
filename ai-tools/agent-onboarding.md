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

## 全域設定 vs 專案設定

每個 AI agent 工具都有兩種設定層級，需要根據使用情境決定用哪一種。

### 設定層級說明

| 層級 | 設定位置 | 生效範圍 | 優先順序 | 適用時機 |
|------|---------|---------|---------|---------|
| **全域** | 工具設定面板（如 Roo Code Custom Instructions） | 所有專案 | 低（可被專案設定覆蓋） | 希望所有專案都自動啟用 Ai-skill 系統 |
| **專案** | 專案根目錄的設定檔（如 `.roomodes`、`.cursor/rules/`） | 單一專案 | 高（覆蓋全域設定） | 專案需要自訂 mode 定義或覆蓋全域設定 |

### 建議策略

1. **全域設定一次**：在工具的全域設定中，用**絕對路徑**指向 Ai-skill 的 `CORE_BOOTSTRAP.md`，加上語言偏好與語言一致性規則
2. **專案設定只在需要時才建立**：如果專案需要自訂 mode 定義或 file restrictions，才在專案根目錄建立設定檔
3. **注意覆蓋行為**：部分工具（如 Roo Code）的專案設定檔會**完全覆蓋**全域設定（不會合併），所以專案設定檔中必須包含全域設定的所有內容

### 路徑注意事項

全域設定中使用的是**絕對路徑**（如 `/Users/larrylin/Documents/Ai-skill/CORE_BOOTSTRAP.md`）。如果：
- **Ai-skill 移動位置** → 更新全域設定中的路徑
- **在其他電腦使用** → 修改為對應的絕對路徑
- **使用相對路徑** → 只能在 Ai-skill repo 內生效，不適合全域設定

## 新增工具的步驟

1. **確認工具類型**：是 AI Agent（CLI / IDE 內建）？Agent 放 `ai-tools/agent/`。
2. **決定設定層級**：全域設定一次（所有專案生效）還是只做專案設定？
3. **建立工具使用說明**：在 `ai-tools/agent/` 下建立 `<tool>.md`，記錄：
   - 全域設定的內容與位置
   - 專案設定檔的格式與位置
   - 兩個層級的覆蓋關係
4. **設定自動載入入口**：依工具的機制設定入口，指向 `CORE_BOOTSTRAP.md`。
5. **設定語言偏好**：依工具的設定方式，加入軟性語言偏好。
6. **實作對話目標閉環**：依工具的能力（hooks / custom instructions / 操作注意），實作 goal ledger 整合。
7. **更新 `ai-tools/README.md`**：在 agent 類別的表格中加入新工具的連結與用途說明。
8. **驗證**：執行 `scripts/validate-knowledge-runtime.rb` 確認無誤。

> **注意**：IDE 生態系統的通用知識（如 VS Code Extension 全域設定的 SQLite 儲存機制）屬於可重複使用的工程智慧，應放在 `intelligence/ide/`，而非 `ai-tools/` 下。`ai-tools/` 只放工具特有的設定與操作方式。

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
