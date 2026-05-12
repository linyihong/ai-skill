# Roo Code 使用說明

本檔說明 Roo Code（VS Code extension）與其他工具的差異。通用配置原則見 [`ai-tools/README.md`](README.md)；知識庫入口見 [`README.md`](../README.md)；啟動流程見 [`CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md)。

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
      "customInstructions": "依 CORE_BOOTSTRAP.md 啟動流程載入核心規則與 OS layout。",
      "groups": ["read", "edit", "command", "mcp"]
    }
  ]
}
```

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

← [回到 AI 工具索引](README.md)
