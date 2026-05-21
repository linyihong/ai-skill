# AI 工具使用說明

本目錄放各種 AI coding / agent 工具如何讀取與同步本知識庫。工具依類型分為以下類別：

| 類別 | 說明 | 包含工具 |
|------|------|----------|
| [`agent/`](agent/) | AI Agent 工具（CLI 或 IDE 內建的 AI 助手） | Roo Code、Cursor、Claude Code、Codex |

> **注意**：VS Code Extension 全域設定的通用知識（SQLite 資料庫位置、結構、查詢/修改方法）已昇華到 [`intelligence/ide/vscode-extension-global-state.md`](../intelligence/ide/vscode-extension-global-state.md) 作為可重複使用的工程智慧。各 agent 工具（如 Roo Code）的專屬設定細節留在對應的工具文件中。

新增其他工具時，依 [`agent-onboarding.md`](agent-onboarding.md) 的 checklist 完成設定。

開新專案時，用 [`new-project-onboarding.md`](new-project-onboarding.md) 了解如何一次設定所有工具，或直接執行 `ai-skill init-project`。

## 配置思想與邊界定義

每個 AI 工具的配置應遵循**薄配置層**設計：

| 層級 | 位置 | 內容 |
|------|------|------|
| **自動載入入口** | 工具專屬入口檔（如 `CLAUDE.md`、`.cursorrules`） | 一行，指向 `README.md` |
| **工具配置** | 工具配置檔（如 `.claude/settings.json`） | 僅放 permissions、hooks 等工具特定設定 |
| **工具使用說明** | `ai-tools/<tool>.md` | 此工具的配置實作與特殊操作注意 |
| **共用規則** | `enforcement/` | 所有規則本體，集中管理 |
| **知識庫入口** | `README.md` | OS layout 與導航，所有工具的共同起點 |

### 工具文件不得重複中央庫內容

`ai-tools/<tool>.md` **只能記錄該工具特有的差異**，不得包含以下已在中央庫可發現的內容：

| ❌ 不應放入工具文件 | ✅ 已在何處 |
|---|---|
| Core Bootstrap 流程（讀 CORE_BOOTSTRAP.md → README.md → ...） | [`CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md) |
| Pipeline 階段（Bootstrap → Routing → Execution → Close-loop） | [`runtime/README.md`](../runtime/README.md) |
| Context expansion 層級（Summary → Module → Detailed → Raw） | [`runtime/runtime.db`](../runtime/runtime.db) |
| Relevance scoring 邏輯 | [`runtime/runtime.db`](../runtime/runtime.db) |
| Guard chain 執行順序 | [`runtime/runtime.db`](../runtime/runtime.db) |
| Token budget 分配表 | [`runtime/runtime.db`](../runtime/runtime.db) |
| 知識庫路徑對照表（CORE_BOOTSTRAP.md → README.md → ...） | [`README.md`](../README.md) OS layout |
| 共用規則清單 | [`enforcement/README.md`](../enforcement/README.md) |
| Goal ledger 操作流程 | [`enforcement/conversation-goal-ledger.md`](../enforcement/conversation-goal-ledger.md) |
| Close-loop 流程（commit/push/readback） | [`enforcement/dependency-reading.md`](../enforcement/dependency-reading.md) |

**原則**：每個工具文件應假設 reader 已讀過 `README.md` 的 OS layout 與 `CORE_BOOTSTRAP.md` 的啟動流程。工具文件只回答：「這個工具跟其他工具有什麼不同？它的入口檔、配置檔、特殊操作要注意什麼？」

**不應放在工具配置或工具說明中的內容：**
- Bootstrap 規則清單（由 `enforcement/README.md` 管理）
- 情境路由表或架構層級表（由 `README.md` / `enforcement/README.md` 管理）
- 任何已在 enforcement 或 README.md 中的重複內容

---

共用原則：**中央庫是真相來源**，工具端只做參照、symlink 或同步快照。

通用 enforcement rules、workflow / analysis / intelligence source、templates 與根 README 應保持工具中立；工具專屬路徑、hook、UI、reload、設定與同步命令放在本目錄對應工具文件中。

Repo-level 載入與同步方向見 [`architecture/ai-native-knowledge-operating-system.md`](../architecture/ai-native-knowledge-operating-system.md)：reference-first 是預設，工具 mirror / bundle / copy snapshot 是相容層。

若工具文件、workflow、使用者目標或 enforcement rules 看似衝突，依 [`enforcement/rule-weight.md`](../enforcement/rule-weight.md) 判斷權重；工具 adapter 不得覆蓋 safety、source-of-truth、validation 或最新使用者目標。

若某個 workflow 對某工具有必要的特殊執行策略，應在對應 workflow 或 `ai-tools/` 連結說明中記錄差異；本目錄仍只放該工具的全域設定、同步與操作方式。

## Agent 工具

| 工具 | 文件 | 用途（僅記錄該工具特有差異） |
| --- | --- | --- |
| Claude Code | [`agent/claude.md`](agent/claude.md) | `CLAUDE.md` 自動載入入口、`.claude/settings.json` 工具配置、tool adapter 機制、對話目標閉環。 |
| Cursor | [`agent/cursor.md`](agent/cursor.md) | `.cursor/rules/*.mdc` 自動載入、`.cursor/hooks.json` 設定、對話目標閉環（含 hooks 範本）。 |
| Roo Code | [`agent/roo.md`](agent/roo.md) | Custom instructions 手動設定、多 modes 與 file restrictions、`.roomodes` 自訂 mode 定義、對話目標閉環、語言偏好雙層設定。 |
| Codex | [`agent/codex.md`](agent/codex.md) | `AGENTS.md` 自動載入入口、runtime SQLite source-of-truth、YAML contract projection、commit / push 更新流程。 |
| **新增工具指引** | [`agent-onboarding.md`](agent-onboarding.md) | 新 AI agent 工具加入時的設定 checklist，含必要項目與參考來源對照。 |

← [回到根目錄](../README.md)
