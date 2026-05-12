# Claude 使用說明

本檔說明「配置思想與邊界定義」在 Claude Code 的具體實作，以及 Claude 特有的操作注意事項。通用配置原則見 [`ai-tools/README.md`](README.md)。

## Claude Code 配置實作

### 自動載入入口：`CLAUDE.md`

Claude Code 啟動時會自動讀取根目錄的 `CLAUDE.md`。本庫的 `CLAUDE.md` 現在實作 **Core Bootstrap 自動載入流程**：

```text
1. 讀 CORE_BOOTSTRAP.md → 3 條必讀核心規則（~800 tokens）
2. 讀 README.md → OS layout
3. 依 task intent 查 skills-index.yaml → 找到對應 skill
4. 依 activation-rules.yaml → 決定 lazy-load rules
5. 先讀 knowledge/summaries/ 對應 summary → 需要時才展開全文
```

這取代了舊的「CLAUDE.md → README.md → shared-rules/README.md (12 rules)」鏈，token 消耗從 ~5000 降至 ~800。

### 工具配置：`.claude/settings.json`

`.claude/settings.json` 現在包含完整的 bootstrap 與 reference 路徑：

```json
{
  "bootstrap": {
    "core": "CORE_BOOTSTRAP.md",
    "layout": "README.md",
    "skill_index": "skills-index.yaml",
    "activation_rules": "runtime/router/activation-rules.yaml"
  },
  "reference": {
    "token_budget": "runtime/budget/token-budget.yaml",
    "context_health": "runtime/health/context-health-score.yaml",
    "circuit_breaker": "runtime/guards/circuit-breaker.yaml",
    "tool_metadata": "tools/metadata/README.md",
    "tool_compression": "tools/compression/README.md",
    "memory_working": "memory/working/README.md",
    "memory_summary": "memory/summary/README.md",
    "memory_decision": "memory/decision/README.md",
    "decisions": "decisions/README.md",
    "anti_patterns": "anti-patterns/README.md"
  }
}
```

**設定一次 repo 即可**：只要 clone 本 repo，Claude Code 啟動時自動讀 `CLAUDE.md` → Core Bootstrap → skills-index.yaml，不需要每次手動指定。

## Claude 操作注意

- Claude 若只能看到單一專案，請同時提供 `<AI_SKILL_REPO>` 的可讀路徑，或把必要 skill/shared-rules 以工具支援的方式同步成可讀上下文。
- 如果 Claude 已經長時間對話，請先要求它讀 `<PROJECT_ROOT>/.agent-goals/`，確認未完成項、優先順序與 owner/lock 狀態。
- 如果 goal 標示 `single-owner` 或 `non-parallelizable`，不要讓 Claude 和其他 agent 分工同一流程；先取得使用者確認。
- 若 Claude 要改本庫，提醒它不要只更新文件；還要跑驗證、commit、push、讀回和 clean status。
- 若 Claude 完成 goal 後仍留下長期 roadmap 或治理狀態，要求它先回寫到 durable planning 文件，再刪除 active goal。

## 與 Tool Adapter 的關係

若某個 skill 針對 Claude 有特殊執行策略（上下文載入順序、prompt chunking、工具輸出限制等），放在：

```
skills/<skill-name>/tool-adapters/claude.md
```

該 adapter 只寫 skill-specific 差異，並連回核心 `WORKFLOW.md` / `TOOLS.md`。

## 驗證

使用 Claude 完成任務時，最後要求它回報：

- 讀了哪些 shared rules 與 skill 依賴。
- 哪些依賴不存在（標示 `not applicable`）。
- 目標是否完成，還有哪些 `.agent-goals` 未完成。
- 驗證方法：diff review、link check、commit/push/readback/clean status。

← [回到 AI 工具索引](README.md)
