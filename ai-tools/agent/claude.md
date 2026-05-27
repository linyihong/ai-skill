# Claude 使用說明

本檔說明 Claude Code 特有的配置與操作注意事項。通用配置原則見 [`ai-tools/README.md`](../README.md)；知識庫入口見 [`README.md`](../../README.md)；啟動流程見 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)。

## 使用者快速上手（onboarding）

第一次在這個 repository 用 Claude Code，建議流程：

1. **開啟 Claude Code session**（`claude` CLI 或 IDE 整合）。
2. **第一句先打** `/bootstrap`，等 Claude 回報 Bootstrap Receipt + Cognitive Mode 報告。
   - Receipt 範例：`Bootstrap: rules=✓ phase=phase.bootstrap obligations=23 gates=25`
   - 看到這兩個區塊代表規則已載入完成。
3. **再給實際任務**（例如「review 這次變更」「init 一個新 workflow」）。

若忘記跑 `/bootstrap`，Claude 預期會在第一個回覆**主動提醒**（見 [`CLAUDE.md`](../../CLAUDE.md) §「第一輪使用者互動」）。若沒提醒，代表 Claude 跳過了 onboarding 指示——可手動打 `/bootstrap` 補救。

### 為什麼需要手動觸發

Cursor、Roo Code 等工具有 always-apply rule 可在每個 turn 機械注入規則；Claude Code 的 `CLAUDE.md` 與 `.claude/rules/*.md` 是注入 system prompt 的 prose，模型可能在「任務看起來簡單」時跳過 bootstrap 序列。`/bootstrap` slash command 是使用者明確觸發，是目前最可靠的方式。

`.claude/hooks/` 雖然提供 PreToolUse / Stop 等機械 gate，但實測在部分環境下不穩定觸發（依 Claude Code 版本而異），不能取代手動 `/bootstrap`。

### 可用的 slash commands

| Command | 用途 |
|---|---|
| `/bootstrap` | 完整執行 CORE_BOOTSTRAP.md + 3 條必讀規則 + runtime.db 查詢，並輸出 Bootstrap Receipt |
| 其他 | 在 Claude Code 介面輸入 `/` 可看到所有可用 commands（含 `.claude/commands/` 自訂） |

## Bootstrap 自動觸發架構（empirically verified 2026-05-27）

Claude Code 的 bootstrap 自動化由**三層協同**達成，empirical 測試確認**單獨任何一層都不夠**：

| 層 | 檔案 | 機制 | 角色 |
|---|---|---|---|
| **L1 SessionStart hook** | `.claude/hooks/auto-bootstrap.sh` | 機械（hook 注入 context） | session 一開啟就把 bootstrap data + 預算好的 Receipt 塞進 context，Claude 才有**精確數值**可以 echo（phase / obligation count / gate count），不需要自己查 |
| **L2 CLAUDE.md 禁止規則** | `CLAUDE.md` 最前面的 `# CRITICAL RUNTIME RULES` | prompt-based | 用「禁止語氣 + 顯式排序 + execution prohibition」把 Receipt echo 推到**第一行**，避免 Claude 直接答題 |
| **L3 PreToolUse hook** | `.claude/hooks/check-bootstrap-receipt.sh` | 機械（hook 攔截 tool call） | 兜底安全網。Claude 偷跑非-Read 工具時 mechanically block（`exit 2`），逼它回去走完整序列 |

### 為什麼三層都需要

- **沒 L1**：Claude 會編造 Receipt 數值（看起來對但其實是猜的）
- **沒 L2**：Claude 可能跳過 Receipt 直接答題（特別是任務看起來簡單時）
- **沒 L3**：Claude 仍然會偶爾偷跑（純 prompt-based 強度不夠）

### 關鍵格式要求（踩過的雷）

**`.claude/settings.json` 必須用三層 nested 格式 + `matcher` 欄位**，flat 格式會被 Claude Code silently 忽略，整個 hook 不註冊：

```json
// ❌ 錯誤：flat 格式（會被忽略）
{
  "hooks": {
    "PreToolUse": [{ "type": "command", "command": "..." }]
  }
}

// ✅ 正確：三層 nested + matcher
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "",
        "hooks": [{ "type": "command", "command": "...", "timeout": 10 }]
      }
    ]
  }
}
```

判斷 hook 有沒有 fire 的方法：在每個 hook script 第一行加 `echo "fired" >> /tmp/<name>.log`，然後檢查 log 是否更新。如果完全沒 entry，就是 settings 格式錯誤導致 hook 沒註冊。

### Claude 對 prompt-based 規則的反應特性

實測發現 Claude 對以下 prompt 模式較難跳過（依強度遞增）：

1. 一般敘述（「請先 bootstrap」）— 容易跳過
2. IMPORTANT block（「IMPORTANT: must run bootstrap first」）— 仍可能跳過
3. **禁止語氣 + 顯式排序 + execution prohibition**（CLAUDE.md `# CRITICAL RUNTIME RULES` 採用此形式）— 最有效
4. + 機械 hook 兜底（exit 2 block tool）— 最可靠

放在文件**最前面**比放在中段或末尾顯著有效（primacy effect）。

### 已知限制

- SessionStart hook 在 Claude Code 不同版本上的 matcher 支援度不一；目前驗證有效的 matcher 字串：`"startup|resume|clear"`
- Per-turn obligations 定義在 `runtime/core-bootstrap.yaml`，**未** projected 到 `runtime/runtime.db` 的 `obligations` table；auto-bootstrap.sh 直接 hardcode 與 YAML source 一致

## Claude Code 配置實作

### 自動載入入口：`CLAUDE.md`

Claude Code 啟動時會自動讀取根目錄的 `CLAUDE.md`。本庫的 `CLAUDE.md` 已實作一行指向 `CORE_BOOTSTRAP.md`，Claude 啟動後會自動依啟動流程載入核心規則與 OS layout。

**設定一次 repo 即可**：只要 clone 本 repo，Claude Code 啟動時自動讀 `CLAUDE.md` → `CORE_BOOTSTRAP.md`，不需要每次手動指定。

### 工具配置：`.claude/settings.json`

`.claude/settings.json` 記錄 Claude Code 的工具特定設定（permissions、bootstrap 路徑等）。詳細內容見該檔案本身，此處不重複。

## 全域設定 vs 專案設定

Claude Code 的設定方式與 Roo Code 不同，沒有「全域 Custom Instructions」的概念。但可以透過以下方式達到「設定一次，所有專案生效」的效果：

### 設定層級說明

| 層級 | 範圍 | 設定位置 | 說明 |
|------|------|----------|------|
| 層級 A：全域（所有專案） | 所有專案 | `~/.claude/` 中的 `claude.md` 或 `projects.json` | 放在使用者的 home 目錄中，所有專案都會載入 |
| 層級 B：專案（單一專案） | 單一專案 | `<PROJECT_ROOT>/CLAUDE.md` | 只對該專案生效 |

### 建議策略

```
全域 ~/.claude/claude.md（層級 A）
  ├── 指向 Ai-skill 的 CORE_BOOTSTRAP.md（絕對路徑）
  ├── 語言偏好設定
  └── 語言一致性強制規則

專案 CLAUDE.md（層級 B，可選）
  ├── 只在需要專案特定規則時建立
  └── 加上該專案特有的規則
```

### 注意事項

- Claude Code 支援 `~/.claude/claude.md` 作為全域設定檔，所有專案啟動時都會載入。
- 如果 Ai-skill 路徑變更，需要更新 `~/.claude/claude.md` 中的路徑。
- 專案 `CLAUDE.md` 中的內容會與全域 `~/.claude/claude.md` 合併（不會覆蓋），因此不需要像 Roo Code 那樣擔心覆蓋問題。

## Claude Code 與對話目標閉環

工具中立規則見 [`enforcement/conversation-goal-ledger.md`](../../enforcement/conversation-goal-ledger.md)。Claude Code 是 CLI 工具，沒有 hooks 機制，但可以透過 `CLAUDE.md` 中的 Custom Instructions 加入 goal ledger 提醒。

**Goal ledger 操作流程已由 runtime 管理**，請參考：
- [`runtime/runtime.db`](../../runtime/runtime.db) — `phase_machine` / `obligation_ledger` / `blocking_gates` 快速查詢
- [`runtime/runtime.db`](../../runtime/runtime.db) — phase / obligation / gate / transaction / recovery 的 embedded source
- `ai-skill goals`（source: [`scripts/ai-skill-cli/internal/app/goals.go`](../../scripts/ai-skill-cli/internal/app/goals.go)）— goal ledger CLI helper

Claude Code 專屬注意事項：
- 無 hooks 機制，需在 `CLAUDE.md` 中手動加入 goal ledger 提醒
- `CLAUDE.md` 已包含基本 goal ledger 提醒

## Claude Code 與知識更新流程 Checkpoint

工具中立規則見 [`governance/lifecycle/knowledge-update-flow.md`](../../governance/lifecycle/knowledge-update-flow.md)。**快速路徑**請優先查詢 runtime.db：[`runtime/runtime.db`](../../runtime/runtime.db)（`generated_surfaces` 表）。Claude Code 是 CLI 工具，沒有 hooks 機制，但可以透過 `CLAUDE.md` 中的 Custom Instructions 加入 checkpoint 提醒。

**Knowledge update flow 已由 runtime 管理**，請參考：
- `runtime/runtime.db → generated_surfaces (type='knowledge_update_phases')` — 11 個步驟的結構化記錄（快速路徑）
- `runtime/runtime.db → recovery_strategies / phase_reconciliation / state_repair` — runtime recovery strategy（快速路徑）
- [`runtime/runtime.db`](../../runtime/runtime.db) — recovery / transaction state machine 的 source

Claude Code 專屬注意事項：
- 無 hooks 機制，需在 `CLAUDE.md` 中手動加入 checkpoint 提醒

## 與 Tool Adapter 的關係

若某個 workflow 針對 Claude 有特殊執行策略（上下文載入順序、prompt chunking、工具輸出限制等），放在對應 workflow 文件或 `ai-tools/` 的 tool-specific 說明中，例如：

```
workflow/<domain>/tool-adapters/claude.md
```

該 adapter 只寫 workflow-specific 差異，並連回核心 `README.md` / `execution-flow.md`。

## 語言偏好設定（重要）

Claude Code 的語言偏好設定方式與 Roo Code（VS Code Extension）不同，因為 Claude Code 是 CLI 工具，**沒有** SQLite 全域資料庫。為了完整解決語言漂移問題，需要在語言偏好設定中加入**語言一致性強制規則**。

### 設定方式

Claude Code 的語言行為由 `CLAUDE.md` 中的 Custom Instructions 控制：

1. **在 `CLAUDE.md` 中設定**：本知識庫的 [`CLAUDE.md`](../../CLAUDE.md) 已包含語言偏好設定。
2. **語言偏好內容**：

```text
Language Preference: Default to English, but always match the user's language in conversation.
If the user writes in Chinese, respond in Chinese.
If the user writes in Japanese, respond in Japanese.
If the user switches languages, follow their switch.

語言一致性強制規則：所有輸出（包含 attempt_completion 結果、技術分析、表格欄位、章節標題、commit message）都必須與使用者當前語言一致。如果使用者使用中文，所有內容（包括技術關鍵詞、程式碼註解、分析報告）都必須使用中文。在 attempt_completion 前必須先確認語言一致性。
```

### 與 Roo Code 的差異

| 特性 | Claude Code | Roo Code |
|------|------------|----------|
| 執行環境 | CLI terminal | VS Code extension |
| 設定位置 | `CLAUDE.md`（檔案） | `.roomodes` + SQLite 全域資料庫 |
| 全域語言欄位 | 無 | 有（`language` 欄位在 `state.vscdb`） |
| 設定方式 | 直接編輯 `CLAUDE.md` | 編輯 `.roomodes` + 修改 SQLite |
| 語言一致性強制規則 | 需手動加入 `CLAUDE.md` | 需手動加入 `.roomodes` 或全域 Custom Instructions |

### 注意事項

- Claude Code 沒有「全域語言強制」的問題，只要 `CLAUDE.md` 中的語言偏好設定正確，Claude 就會跟隨使用者語言。
- 如果 Claude 仍然強制使用英文，請檢查 `CLAUDE.md` 中是否有固定的 `You should always speak and think in the "English" (en) language` 設定，改為上述軟性偏好 + 語言一致性強制規則即可。
- **語言一致性強制規則**是為了解決「作者習慣漂移」問題（Type B），即 agent 在描述技術細節時會不自覺使用英文。加入此規則後，所有輸出（包括 attempt_completion、表格、commit message）都會強制跟隨使用者語言。

## 驗證

使用 Claude 完成任務時，最後要求它回報：

- 讀了哪些 enforcement rules 與 skill 依賴。
- 哪些依賴不存在（標示 `not applicable`）。
- 目標是否完成，還有哪些 `.agent-goals` 未完成。
- 驗證方法：diff review、link check、commit/push/readback/clean status。

← [回到 AI 工具索引](../README.md)
