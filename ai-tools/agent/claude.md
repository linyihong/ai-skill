# Claude 使用說明

本檔說明 Claude Code 特有的配置與操作注意事項。通用配置原則見 [`ai-tools/README.md`](../README.md)；知識庫入口見 [`README.md`](../../README.md)；啟動流程見 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)。

## Claude Code 配置實作

### 自動載入入口：`CLAUDE.md`

Claude Code 啟動時會自動讀取根目錄的 `CLAUDE.md`。本庫的 `CLAUDE.md` 已實作一行指向 `CORE_BOOTSTRAP.md`，Claude 啟動後會自動依啟動流程載入核心規則與 OS layout。

**設定一次 repo 即可**：只要 clone 本 repo，Claude Code 啟動時自動讀 `CLAUDE.md` → `CORE_BOOTSTRAP.md`，不需要每次手動指定。

### 工具配置：`.claude/settings.json`

`.claude/settings.json` 記錄 Claude Code 的工具特定設定（permissions、bootstrap 路徑等）。詳細內容見該檔案本身，此處不重複。

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

## 語言偏好設定

Claude Code 的語言偏好設定方式與 Roo Code（VS Code Extension）不同，因為 Claude Code 是 CLI 工具，**沒有** SQLite 全域資料庫。

### 設定方式

Claude Code 的語言行為由 `CLAUDE.md` 中的 Custom Instructions 控制：

1. **在 `CLAUDE.md` 中設定**：本知識庫的 [`CLAUDE.md`](../CLAUDE.md) 已包含語言偏好設定。
2. **語言偏好內容**：

```text
Language Preference: Default to English, but always match the user's language in conversation.
If the user writes in Chinese, respond in Chinese.
If the user writes in Japanese, respond in Japanese.
If the user switches languages, follow their switch.
```

### 與 Roo Code 的差異

| 特性 | Claude Code | Roo Code |
|------|------------|----------|
| 執行環境 | CLI terminal | VS Code extension |
| 設定位置 | `CLAUDE.md`（檔案） | `.roomodes` + SQLite 全域資料庫 |
| 全域語言欄位 | 無 | 有（`language` 欄位在 `state.vscdb`） |
| 設定方式 | 直接編輯 `CLAUDE.md` | 編輯 `.roomodes` + 修改 SQLite |

### 注意事項

- Claude Code 沒有「全域語言強制」的問題，只要 `CLAUDE.md` 中的語言偏好設定正確，Claude 就會跟隨使用者語言。
- 如果 Claude 仍然強制使用英文，請檢查 `CLAUDE.md` 中是否有固定的 `You should always speak and think in the "English" (en) language` 設定，改為上述軟性偏好即可。

## 驗證

使用 Claude 完成任務時，最後要求它回報：

- 讀了哪些 shared rules 與 skill 依賴。
- 哪些依賴不存在（標示 `not applicable`）。
- 目標是否完成，還有哪些 `.agent-goals` 未完成。
- 驗證方法：diff review、link check、commit/push/readback/clean status。

← [回到 AI 工具索引](../README.md)
