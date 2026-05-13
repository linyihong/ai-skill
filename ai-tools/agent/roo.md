# Roo Code 使用說明

本檔說明 Roo Code（VS Code extension）與其他工具的差異。通用配置原則見 [`ai-tools/README.md`](../README.md)；知識庫入口見 [`README.md`](../../README.md)；啟動流程見 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)。

## Roo Code 與其他工具的差異

| 特性 | Claude Code | Cursor | Roo Code |
|------|------------|--------|----------|
| 執行環境 | CLI terminal | 獨立 IDE | VS Code extension |
| 自動載入入口 | `CLAUDE.md` | `.cursor/rules/*.mdc` | Custom instructions（VS Code settings 或 `.roomodes`） |
| 工具配置 | `.claude/settings.json` | `.cursor/hooks.json` | VS Code settings（`roo-code.*`） |
| Modes | 單一模式 | 單一模式 | 多 modes（code/architect/ask/debug 等） |
| File restrictions | 無 | 無 | 每個 mode 可限制可讀寫的 file patterns |

## Roo Code 配置實作

### 自動載入入口：Custom Instructions

Roo Code 沒有像 `CLAUDE.md` 那樣的自動載入機制。你需要手動設定 **Custom Instructions**：

1. 在 VS Code 中開啟本 repo
2. 點擊 Roo Code extension 的設定圖示（⚙️）
3. 在 **Custom Instructions** 中貼入以下內容：

```text
你是一個運行在 Roo Code（VS Code AI extension）的 AI agent。

開始工作前，請依 CORE_BOOTSTRAP.md 的啟動流程載入核心規則與 OS layout。
```

Roo Code 會自動讀取根目錄的 `CORE_BOOTSTRAP.md`、`README.md`、`skills-index.yaml` 等檔案，不需要在 custom instructions 中列出完整 bootstrap 流程。

### Modes 設定

Roo Code 支援多種 modes，每個 mode 可以有不同的 custom instructions 與 file restrictions。
你可以在 VS Code 的 `settings.json` 中設定：

```json
{
  "roo-code.modes": [
    {
      "slug": "code",
      "name": "💻 Code",
      "roleDefinition": "You are Roo, a highly skilled software engineer...",
      "customInstructions": "依 CORE_BOOTSTRAP.md 啟動流程載入核心規則與 OS layout。",
      "groups": ["read", "edit", "command", "mcp"]
    },
    {
      "slug": "architect",
      "name": "🏗️ Architect",
      "roleDefinition": "You are Roo, an expert software architect...",
      "customInstructions": "依 CORE_BOOTSTRAP.md 啟動流程載入核心規則與 OS layout。",
      "groups": ["read"],
      "fileRestrictions": ["**/*.md"]
    },
    {
      "slug": "ask",
      "name": "❓ Ask",
      "roleDefinition": "You are Roo, a knowledgeable technical advisor...",
      "customInstructions": "依 CORE_BOOTSTRAP.md 啟動流程載入核心規則與 OS layout。",
      "groups": ["read"]
    },
    {
      "slug": "debug",
      "name": "🪲 Debug",
      "roleDefinition": "You are Roo, a systematic debugger...",
      "customInstructions": "依 CORE_BOOTSTRAP.md 啟動流程載入核心規則與 OS layout。",
      "groups": ["read", "edit", "command", "mcp"]
    }
  ]
}
```

或者，你也可以在專案根目錄建立 `.roomodes` 檔案（Roo Code 會自動讀取）：

```json
{
  "customModes": [
    {
      "slug": "ai-skill-maintainer",
      "name": "Ai-skill Maintainer",
      "roleDefinition": "你負責維護 Ai-skill 知識庫的架構、規則與技能。",
      "customInstructions": "依 CORE_BOOTSTRAP.md 啟動流程載入核心規則與 OS layout。\n\nLanguage Preference: Default to English, but always match the user's language in conversation. If the user writes in Chinese, respond in Chinese. If the user writes in Japanese, respond in Japanese. If the user switches languages, follow their switch.",
      "groups": ["read", "edit", "command", "mcp"]
    }
  ]
}
```

> **注意**：`.roomodes` 中的 `customInstructions` 會**覆蓋** Roo Code Extension 設定面板中的全域 Custom Instructions。如果你已經在全域設定中寫了固定的語言偏好，`.roomodes` 的設定會優先。

### 建議的 `.roomodes` 設定

由於本知識庫有多個操作面向，建議建立以下自訂 modes：

| Mode | Slug | 用途 | Groups |
|------|------|------|--------|
| 🏗️ Architect | `architect` | 規劃架構、設計 pipeline | read |
| 💻 Code | `code` | 實作元件、寫規則、建立檔案 | read, edit, command, mcp |
| ❓ Ask | `ask` | 查詢知識、解釋架構 | read |
| 🪲 Debug | `debug` | 診斷問題、追蹤失效模式 | read, edit, command, mcp |
| 🪃 Orchestrator | `orchestrator` | 協調多步驟任務 | read, edit, command, mcp |
| 📝 Skill Writer | `skill-writer` | 撰寫 skill、feedback lesson | read, edit |
| 🧹 Governance | `governance` | 維護 lifecycle、validation | read, edit |

## 語言偏好設定（重要）

Roo Code 的語言偏好設定涉及**兩個層級**，需要分別處理才能完整解決語言漂移問題。

### 問題

如果 Custom Instructions 中寫了固定的語言偏好（例如 `You should always speak and think in the "English" (en) language`），agent 會**無視使用者實際使用的語言**，強制用該語言回應。這是因為 Custom Instructions 在 system prompt 中的優先級高於對話上下文。

此外，Roo Code 的全域設定中有一個 `language` 欄位（預設 `"en"`），此欄位會被傳入 system prompt 建構過程，進一步強化英文偏好。

### 解決方案：兩個層級

#### 層級 1：Custom Instructions（`.roomodes`）

在 `.roomodes` 的每個 mode 的 `customInstructions` 中，將語言偏好設為**軟性預設值**，加上「跟隨使用者語言」的例外：

```text
Language Preference: Default to English, but always match the user's language in conversation.
If the user writes in Chinese, respond in Chinese.
If the user writes in Japanese, respond in Japanese.
If the user switches languages, follow their switch.
```

> **原理**：`.roomodes` 中的 `customInstructions` 會**覆蓋** Roo Code Extension 設定面板中的全域 Custom Instructions，因此這是設定語言偏好的主要方式。

#### 層級 2：全域 `language` 欄位（SQLite）

Roo Code 的全域設定儲存在 VS Code 的 **globalStorage SQLite 資料庫**中。其中的 `language` 欄位（預設 `"en"`）會影響 system prompt 的建構。

修改方法見 [`ai-tools/ide/vscode-extension-global-state.md`](../ide/vscode-extension-global-state.md)（通用 VS Code Extension 全域設定修改指南），Roo Code 專屬的 key 為 `RooVeterinaryInc.roo-cline`，需修改的欄位為 `language`（設為 `"zh-CN"`）。

### 如何預設語言偏好以跟隨文件規則

本知識庫的文件規則要求：**預設英文，但跟隨使用者語言**。要讓 Roo Code 遵守此規則，需完成以下步驟：

1. **建立 `.roomodes`**（已建立）：在專案根目錄建立 `.roomodes`，每個 mode 的 `customInstructions` 中包含上述軟性語言偏好設定。
2. **修改全域 `language` 欄位**（已修改）：將 SQLite 中的 `language` 從 `"en"` 改為 `"zh-CN"`，避免 system prompt 強制使用英文。
3. **驗證**：啟動新 session 後，用中文提問確認 agent 以中文回應。

> **注意**：如果只做步驟 1 不做步驟 2，Roo Code 的 system prompt 仍可能因 `language: "en"` 而傾向英文。兩個層級都設定才能完整解決。

### 相關資源

- VS Code Extension 全域設定修改通用指南: [`ai-tools/ide/vscode-extension-global-state.md`](../ide/vscode-extension-global-state.md)
- Failure pattern: [`shared-rules/failure-patterns/language-preference-drift.md`](../shared-rules/failure-patterns/language-preference-drift.md)
- Feedback lesson: [`feedback/history/roo-code/2026-05-13_124800-language-preference-drift.md`](../feedback/history/roo-code/2026-05-13_124800-language-preference-drift.md)
- 設定檔: [`.roomodes`](../.roomodes)

## Roo Code 操作注意

### 1. 工作區設定

Roo Code 直接在 VS Code 中操作，所以工作區就是目前開啟的資料夾。
建議用多資料夾工作區同時開啟業務專案與本 repo：

```text
<PROJECT_ROOT>/          ← 業務專案
<AI_SKILL_REPO>/         ← 本知識庫
```

### 2. 同步與更新

由於 Roo Code 直接操作 VS Code 工作區，同步流程如下：

1. 開始工作前：`git pull` 確保本庫最新
2. 修改檔案：直接編輯本庫內的檔案
3. 完成後：`git add` → `git commit` → `git push`
4. 確認：`git log --oneline -3` + `git status`

### 3. 與其他工具的協作

如果你同時使用 Claude Code（CLI）和 Roo Code（VS Code）：

- **Claude Code** 適合 CLI 操作、批次腳本、git 操作
- **Roo Code** 適合 VS Code 內的開發、檔案編輯、即時預覽
- 兩者共用同一份 `CORE_BOOTSTRAP.md`、`skills-index.yaml` 與 `runtime/pipeline/`
- 修改規則時，只需修改本庫一份，兩個工具都會讀到最新版本

---

← [回到 AI 工具索引](../README.md)
