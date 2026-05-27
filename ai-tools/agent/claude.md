# Claude 使用說明

本檔說明 Claude Code 特有的配置與操作注意事項。通用配置原則見 [`ai-tools/README.md`](../README.md)；知識庫入口見 [`README.md`](../../README.md)；啟動流程見 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)。

## 使用者快速上手（onboarding）

第一次在這個 repository 用 Claude Code，建議流程：

1. **開啟 Claude Code session**（`claude` CLI 或 IDE 整合）。
2. 第一個 assistant response 應直接回報 Bootstrap Receipt + Cognitive Mode 報告；SessionStart hook 會預先注入 receipt 所需資料。
   - Receipt 範例：`Bootstrap: rules=✓ phase=phase.bootstrap obligations=23 gates=25`
   - 看到這兩個區塊代表規則已載入完成。
3. **再給實際任務**（例如「review 這次變更」「init 一個新 workflow」）。

若 SessionStart hook 沒有觸發，或第一個 response 沒有 receipt，可手動執行 `/bootstrap` 補救。

### 為什麼仍保留 `/bootstrap`

Claude Code 的主要路徑現在是 hooks 自動化：SessionStart 注入 bootstrap context、PreToolUse 阻擋 receipt 前的非讀取工具、PostToolUse 補提醒、Stop 檢查 Cognitive Mode。`/bootstrap` 保留為 fallback，避免 Claude Code 版本差異、hook 未註冊或使用者在非標準環境執行時卡住。

`.claude/hooks/` 是本 repo 的 Claude-specific enforcement 例外；這些細節只放在 Claude adapter 與 `.claude/` 設定，不推廣到其他工具入口。

### 可用的 slash commands

| Command | 用途 |
|---|---|
| `/bootstrap` | 完整執行 CORE_BOOTSTRAP.md + 3 條必讀規則 + runtime.db 查詢，並輸出 Bootstrap Receipt |
| 其他 | 在 Claude Code 介面輸入 `/` 可看到所有可用 commands（含 `.claude/commands/` 自訂） |

## Bootstrap 自動觸發架構（empirically verified 2026-05-27）

> **設計原理**：見 [`intelligence/engineering/agent-architecture/multi-layer-enforcement.md`](../../intelligence/engineering/agent-architecture/multi-layer-enforcement.md)。本節只記錄 Claude Code 的具體實作對應。

Claude Code 的 bootstrap 自動化由**三層協同**達成（對應 multi-layer-enforcement 的 L1/L2/L3 模式）：

| 層 | 對應 Claude Code 檔案 | 角色 |
|---|---|---|
| **L1 資料注入** | `.claude/hooks/auto-bootstrap.sh`（SessionStart hook） | session 開啟時把預算好的 Receipt 資料注入 context，Claude 才有**精確數值**可以 echo，不需要自己查 |
| **L2 禁止 prompt** | `CLAUDE.md` 最前面的 `# CRITICAL RUNTIME RULES` | 禁止語氣 + 顯式排序 + execution prohibition，把 Receipt echo 推到第一行 |
| **L3 機械關卡** | `.claude/hooks/check-bootstrap-receipt.sh`（PreToolUse hook） | Claude 偷跑非-Read 工具時以 `exit 2` 攔截，逼它回去走完整序列 |

### 關鍵格式要求（Claude Code 專屬）

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

### 已知限制（Claude Code 版本相關）

- SessionStart hook 在 Claude Code 不同版本上的 matcher 支援度不一；目前驗證有效的 matcher 字串：`"startup|resume|clear"`
- Per-turn obligations 定義在 `runtime/core-bootstrap.yaml`，**未** projected 到 `runtime/runtime.db` 的 `obligations` table；auto-bootstrap.sh 直接 hardcode 與 YAML source 一致
- ~~雙重 bootstrap~~（已修）：agent 在同一 turn 同時輸出 Receipt 文字 + 呼叫工具時，工具比文字先 fire，PreToolUse 掃 transcript 找不到 Receipt 而 block。修法：SessionStart 成功後寫 `/tmp/ai-skill-sessionstart-{project_hash}.flag`（TTL 120s），PreToolUse 優先檢查 flag，避免重讀 3 個檔案

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

工具中立規則見 [`enforcement/conversation-goal-ledger.md`](../../enforcement/conversation-goal-ledger.md)。Claude Code hooks 目前只處理 bootstrap receipt、per-turn reminder、PreToolUse gate 與 Cognitive Mode stop check；goal ledger 的 truth 仍在 `<PROJECT_ROOT>/.agent-goals/` 與 runtime / CLI helper，不放進 Claude adapter 重述。

**Goal ledger 操作流程已由 runtime 管理**，請參考：
- [`runtime/runtime.db`](../../runtime/runtime.db) — `phase_machine` / `obligation_ledger` / `blocking_gates` 快速查詢
- [`runtime/runtime.db`](../../runtime/runtime.db) — phase / obligation / gate / transaction / recovery 的 embedded source
- `ai-skill goals`（source: [`scripts/ai-skill-cli/internal/app/goals.go`](../../scripts/ai-skill-cli/internal/app/goals.go)）— goal ledger CLI helper

Claude Code 專屬注意事項：
- hook 可以提醒與阻擋部分 bootstrap / output contract 問題，但不會自動完成或刪除 `.agent-goals/`。
- goal ledger 的完成與刪除仍需依共用規則驗證。

## Claude Code 與知識更新流程 Checkpoint

工具中立規則見 [`governance/lifecycle/knowledge-update-flow.md`](../../governance/lifecycle/knowledge-update-flow.md)。**快速路徑**請優先查詢 runtime.db：[`runtime/runtime.db`](../../runtime/runtime.db)（`generated_surfaces` 表）。Claude Code hooks 只提供 bootstrap / output enforcement，不取代 knowledge-update-flow 的 source checks、linked updates、runtime refresh 或 close-loop。

**Knowledge update flow 已由 runtime 管理**，請參考：
- `runtime/runtime.db → generated_surfaces (type='knowledge_update_phases')` — 11 個步驟的結構化記錄（快速路徑）
- `runtime/runtime.db → recovery_strategies / phase_reconciliation / state_repair` — runtime recovery strategy（快速路徑）
- [`runtime/runtime.db`](../../runtime/runtime.db) — recovery / transaction state machine 的 source

Claude Code 專屬注意事項：
- 如果修改 `.claude/settings.json` 或 `.claude/hooks/*`，必須實測 hook 是否 fire，不能只看設定檔。
- Knowledge update checkpoint 的內容仍由 runtime / governance source 決定，不在 hook script 中複製。

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
