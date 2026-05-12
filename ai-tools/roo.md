# Roo Code 使用說明

本檔說明 Roo Code（VS Code extension）如何讀取與同步 Ai-skill 知識庫。Roo Code 是 VS Code 內的 AI agent extension，支援自訂 modes、custom instructions 與 file restrictions。

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

## Core Bootstrap（每個 session 必讀）

開始工作前，請依序執行以下 bootstrap 流程：

1. 讀取 CORE_BOOTSTRAP.md → 3 條核心規則（rule-weight、dependency-reading、goal-ledger）
2. 讀取 README.md → 了解 OS layout 與 quickstart
3. 依 user intent 查詢 skills-index.yaml → 找到對應 skill
4. 讀取 runtime/router/activation-rules.yaml → 決定哪些 lazy-load rules 需要 activate
5. 先讀 knowledge/summaries/ 對應的 summary（300-500 tokens），需要時才展開全文

## Session Lifecycle

遵循 runtime/pipeline/session-lifecycle.yaml 的四階段流程：
- Bootstrap（~2000 tokens）：載入核心規則 + OS layout
- Routing（~2500 tokens）：解析 intent → 查 skill index → relevance scoring → 載入 summary
- Execution（~100000 tokens）：執行工具呼叫 + health check + compression
- Close-loop（~1000 tokens）：寫入 working memory + session summary + ADR

## Context Expansion

遵循 runtime/pipeline/context-flow.yaml 的漸進式 context 擴展：
- Summary level（~500 tokens）：只載入 summary
- Module summary level（~1500 tokens）：需要時載入 README
- Detailed source level（~4500 tokens）：需要時載入完整 source
- Raw source level（~10000 tokens）：僅 debug/deep analysis 時使用

## Skill Relevance

使用 runtime/pipeline/relevance-engine.yaml 的 scoring 邏輯：
- 計算每個 skill 的 relevance score（trigger_match 0.5 + domain_match 0.3 + weight 0.2）
- Score ≥ 0.5 才載入 summary，< 0.5 跳過
- 衝突的 skill 會套用 penalty（×0.5）

## Guards

遵循 runtime/pipeline/guard-chain.yaml 的 guard 執行順序：
- Critical guards（token_budget、recursive_depth、tool_calls）：每個 tool call 前檢查
- High guards（context_growth、hallucination_risk、conversation_length）：每個 task 完成後檢查
- Medium guards（conflict_rules、repetitive_edits、module_count 等）：定期檢查

## Close-loop

每個 session 結束時：
1. 寫入 memory/working/ 的 working memory
2. 寫入 memory/summary/ 的 session summary（≤500 tokens）
3. 記錄架構決策到 decisions/（如有）
4. 如果 pollution critical，auto-archive 到 memory/working/session-archive-{timestamp}.md
5. 執行 git add → git commit → git push → 讀回確認 → git status clean
```

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
      "customInstructions": "（貼入上面的 Core Bootstrap 內容）",
      "groups": [
        "read",
        "edit",
        "command",
        "mcp"
      ]
    },
    {
      "slug": "architect",
      "name": "🏗️ Architect",
      "roleDefinition": "You are Roo, an expert software architect...",
      "customInstructions": "（同上，但強調規劃與設計）",
      "groups": [
        "read"
      ],
      "fileRestrictions": [
        "**/*.md"
      ]
    },
    {
      "slug": "ask",
      "name": "❓ Ask",
      "roleDefinition": "You are Roo, a knowledgeable technical advisor...",
      "customInstructions": "（同上，但強調解釋與分析）",
      "groups": [
        "read"
      ]
    },
    {
      "slug": "debug",
      "name": "🪲 Debug",
      "roleDefinition": "You are Roo, a systematic debugger...",
      "customInstructions": "（同上，但強調問題診斷）",
      "groups": [
        "read",
        "edit",
        "command",
        "mcp"
      ]
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
      "customInstructions": "（貼入上面的 Core Bootstrap 內容）",
      "groups": [
        "read",
        "edit",
        "command",
        "mcp"
      ]
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

### 2. 與 Pipeline 的整合

Roo Code 的 session lifecycle 與本知識庫的 pipeline 完全對應：

| Pipeline Stage | Roo Code 對應行為 |
|---------------|------------------|
| Bootstrap | 讀取 CORE_BOOTSTRAP.md + README.md |
| Routing | 解析 user intent → 查 skills-index.yaml → relevance scoring |
| Execution | 執行工具呼叫（read_file、write_to_file、execute_command 等） |
| Close-loop | 寫入 memory → commit → push → readback → clean status |

### 3. Token 管理

Roo Code 在 VS Code 中運行，token 限制取決於使用的模型。
建議遵循 `runtime/budget/token-budget.yaml` 的預算分配：

| Layer | Budget | 說明 |
|-------|--------|------|
| Bootstrap | 2000 tokens | CORE_BOOTSTRAP.md + README.md |
| Skill index | 500 tokens | skills-index.yaml |
| Activation rules | 300 tokens | runtime/router/activation-rules.yaml |
| Summaries | 500 tokens per skill | knowledge/summaries/ |
| Full source | 4500 tokens per skill | 需要時才載入 |
| Tool output | 動態壓縮 | 依 context health 調整 compression level |

### 4. 與既有層的關係

| 本庫路徑 | Roo Code 用途 |
|---------|-------------|
| `CORE_BOOTSTRAP.md` | 每個 session 的必讀入口（3 條核心規則） |
| `README.md` | OS layout 與 quickstart |
| `skills-index.yaml` | Skill routing index |
| `runtime/pipeline/` | Session lifecycle、context expansion、guard chain、relevance engine |
| `runtime/budget/token-budget.yaml` | Token 預算管理 |
| `runtime/health/context-health-score.yaml` | Context 健康度評估 |
| `runtime/guards/` | Circuit breaker 與 pollution detection |
| `tools/compression/README.md` | Tool output compression |
| `memory/` | Working / summary / decision 記憶層 |
| `decisions/` | ADR 架構決策記錄 |
| `anti-patterns/` | 已知失效模式 |
| `feedback/pipeline/` | Feedback promotion pipeline |
| `shared-rules/` | 共用規則 |
| `knowledge/summaries/` | 知識摘要 |
| `knowledge/indexes/` | 知識索引與路由 |
| `knowledge/runtime/` | Runtime routing registry、SQLite index |

### 5. 同步與更新

由於 Roo Code 直接操作 VS Code 工作區，同步流程如下：

1. 開始工作前：`git pull` 確保本庫最新
2. 修改檔案：直接編輯本庫內的檔案
3. 完成後：`git add` → `git commit` → `git push`
4. 確認：`git log --oneline -3` + `git status`

### 6. 與其他工具的協作

如果你同時使用 Claude Code（CLI）和 Roo Code（VS Code）：

- **Claude Code** 適合 CLI 操作、批次腳本、git 操作
- **Roo Code** 適合 VS Code 內的開發、檔案編輯、即時預覽
- 兩者共用同一份 `CORE_BOOTSTRAP.md`、`skills-index.yaml` 與 `runtime/pipeline/`
- 修改規則時，只需修改本庫一份，兩個工具都會讀到最新版本

---

← [回到 AI 工具索引](README.md)
