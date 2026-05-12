# AI 工具使用說明

本目錄放各種 AI coding / agent 工具如何讀取與同步本知識庫。新增其他工具時，在這裡加新的子檔，例如 CLI agent、IDE agent、雲端 agent 或其他支援 skill/rules 的工具。

## 配置思想與邊界定義

每個 AI 工具的配置應遵循**薄配置層**設計：

| 層級 | 位置 | 內容 |
|------|------|------|
| **自動載入入口** | 工具專屬入口檔（如 `CLAUDE.md`、`.cursorrules`） | 一行，指向 `README.md` |
| **工具配置** | 工具配置檔（如 `.claude/settings.json`） | 僅放 permissions、hooks 等工具特定設定 |
| **工具使用說明** | `ai-tools/<tool>.md` | 此工具的配置實作與 Claude 特殊操作注意 |
| **共用規則** | `shared-rules/` | 所有規則本體，集中管理 |
| **知識庫入口** | `README.md` | OS layout 與導航，所有工具的共同起點 |

**不應放在工具配置或工具說明中的內容：**
- Bootstrap 規則清單（由 `shared-rules/README.md` 管理）
- 情境路由表或架構層級表（由 `README.md` / `shared-rules/README.md` 管理）
- 任何已在 shared-rules 或 README.md 中的重複內容

---

共用原則：**中央庫是真相來源**，工具端只做參照、symlink 或同步快照。

通用 shared rules、skills、templates 與根 README 應保持工具中立；工具專屬路徑、hook、UI、reload、設定與同步命令放在本目錄對應工具文件中。

Repo-level 載入與同步方向見 [`architecture/ai-native-knowledge-operating-system.md`](../architecture/ai-native-knowledge-operating-system.md)：reference-first 是預設，工具 mirror / bundle / copy snapshot 是相容層。

若工具文件、skill workflow、使用者目標或 shared rules 看似衝突，依 [`shared-rules/rule-weight.md`](../shared-rules/rule-weight.md) 判斷權重；工具 adapter 不得覆蓋 safety、source-of-truth、validation 或最新使用者目標。

若某個 skill 對某工具有必要的特殊執行策略，skill 內可用 `skills/<skill>/tool-adapters/<tool>.md` 記錄差異；本目錄仍只放該工具的全域設定、同步與操作方式。

| 工具 | 文件 | 用途 |
| --- | --- | --- |
| Claude | [claude.md](claude.md) | Claude 類工具如何明確讀取 shared rules、skill 入口、依賴文件、goal ledger 與 Ai-skill writeback 流程。 |
| Cursor | [cursor.md](cursor.md) | Cursor 如何啟用 apk-analysis、參照或同步 `.cursor`、維持中央庫一致性，並以工具中立 `.agent-goals/` 做對話目標閉環提醒。 |
| Roo Code | [roo.md](roo.md) | Roo Code（VS Code AI extension）如何設定 custom instructions、modes、file restrictions，以及與 runtime pipeline 的整合。 |

← [回到根目錄](../README.md)
