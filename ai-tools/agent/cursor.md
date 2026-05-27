# Cursor 使用說明

本檔說明 Cursor 特有的配置與操作注意事項。通用配置原則見 [`ai-tools/README.md`](../README.md)；知識庫入口見 [`README.md`](../../README.md)；啟動流程見 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)。

## 預設載入（Core Bootstrap）

在 Cursor 中，`.cursor/rules/dependency-reading.mdc`（`alwaysApply: true`）只作為 thin entry pointer，指向 `CORE_BOOTSTRAP.md` 與 `runtime/core-bootstrap.yaml`。

**設定一次 repo 即可**：只要 clone 本 repo，Cursor 啟動時自動套用 `dependency-reading.mdc`（alwaysApply）。Bootstrap obligations、format、enum、examples 與後續 runtime 讀取都不寫在 Cursor rule；以 `runtime/core-bootstrap.yaml` 為 canonical。

### 中央庫先行

在新專案開始前或開始時，在本機中央庫目錄執行 `git pull`，確保知識庫與遠端一致。

通用開場提示先用根 [`README.md`](../../README.md) 的說明；本檔只補 Cursor 的工具差異。

### 最穩用法

讓專案 `.cursor` 或開場提示明講要使用哪個 domain，並指定讀中央庫的流程檔：

```text
請先閱讀共用規則索引：
<AI_SKILL_REPO>/enforcement/README.md
之後依 workflow/<domain>/execution-flow.md 進行。
完成後在 <AI_SKILL_REPO> commit 並 push。
```

若工作區已用多資料夾同時打開業務專案與中央庫，Agent 直接開檔最穩。

## 全域設定 vs 專案設定

Cursor 的設定方式與 Roo Code 不同，沒有「全域 Custom Instructions」的概念。但可以透過以下方式達到「設定一次，所有專案生效」的效果：

### 設定層級說明

| 層級 | 範圍 | 設定位置 | 說明 |
|------|------|----------|------|
| 層級 A：全域（所有專案） | 所有專案 | `~/.cursor/rules/` 中的 `alwaysApply: true` 規則 | 放在使用者的 home 目錄 `.cursor/rules/` 中，所有專案都會載入 |
| 層級 B：專案（單一專案） | 單一專案 | `<PROJECT_ROOT>/.cursor/rules/*.mdc` | 只對該專案生效 |

### 建議策略

```
全域 .cursor/rules/（層級 A，放在 ~/.cursor/rules/）
  └── 指向 Ai-skill 的 CORE_BOOTSTRAP.md、runtime/core-bootstrap.yaml（絕對路徑）

專案 .cursor/rules/（層級 B，可選）
  ├── thin bootstrap pointer
  └── 專案特定 gates / overlays
```

### 注意事項

- Cursor 的 `alwaysApply` 規則會在所有專案中載入，因此放在 `~/.cursor/rules/` 中的規則可以達到「設定一次，所有專案生效」的效果。
- 如果 Ai-skill 路徑變更，需要更新 `~/.cursor/rules/` 中的規則路徑。
- 專案 `.cursor/rules/` 中的規則會與全域規則合併（不會覆蓋），因此不需要像 Roo Code 那樣擔心覆蓋問題。

## 語言偏好設定（重要）

Cursor 的語言行為由 `.cursor/rules/*.mdc` 中的規則控制。為了完整解決語言漂移問題，需要在語言偏好設定中加入**語言一致性強制規則**。

### 設定方式

1. **在 `~/.cursor/rules/` 中設定**（全域生效）：將語言偏好設定放在 `~/.cursor/rules/` 中的 `alwaysApply: true` 規則中，所有專案都會載入。
2. **語言偏好內容**（在 `alwaysApply` 的規則中）：

```text
Language Preference: Default to English, but always match the user's language in conversation.
If the user writes in Chinese, respond in Chinese.
If the user writes in Japanese, respond in Japanese.
If the user switches languages, follow their switch.

語言一致性強制規則：所有輸出（包含 attempt_completion 結果、技術分析、表格欄位、章節標題、commit message）都必須與使用者當前語言一致。如果使用者使用中文，所有內容（包括技術關鍵詞、程式碼註解、分析報告）都必須使用中文。在 attempt_completion 前必須先確認語言一致性。
```

### 與 Roo Code 的差異

| 特性 | Cursor | Roo Code |
|------|--------|----------|
| 執行環境 | 獨立 IDE（基於 VS Code） | VS Code extension |
| 設定位置 | `.cursor/rules/*.mdc`（檔案） | `.roomodes` + SQLite 全域資料庫 |
| 全域語言欄位 | 無公開 API 直接修改 | 有（`language` 欄位在 `state.vscdb`） |
| 設定方式 | 編輯 `.cursor/rules/` 中的規則 | 編輯 `.roomodes` + 修改 SQLite |
| 語言一致性強制規則 | 需手動加入 `alwaysApply` 規則 | 需手動加入 `.roomodes` 或全域 Custom Instructions |

### 注意事項

- Cursor 沒有公開的 SQLite 全域資料庫修改方式，語言偏好完全由 `.cursor/rules/` 中的規則控制。
- 如果 Cursor 仍然強制使用英文，請檢查 `.cursor/rules/` 中是否有固定的語言偏好設定，改為上述軟性偏好 + 語言一致性強制規則即可。
- Cursor 的 `alwaysApply` 規則會在每個 session 啟動時自動載入，因此語言偏好設定放在 `alwaysApply: true` 的規則中最有效。
- **語言一致性強制規則**是為了解決「作者習慣漂移」問題（Type B），即 agent 在描述技術細節時會不自覺使用英文。加入此規則後，所有輸出（包括 attempt_completion、表格、commit message）都會強制跟隨使用者語言。

## 公用更新流程

以本 repository 為準：

1. 在 `<AI_SKILL_REPO>` 執行 `git pull`（若與他人共用或換機）。
2. 只在本庫編輯文件。
3. 依 [`enforcement/linked-updates.md`](../../enforcement/linked-updates.md) 檢查連動更新。
4. 在 `<AI_SKILL_REPO>` 執行 `git add`、`git commit`、`git push`。

## Cursor 連動更新

目標是：一邊在業務專案裡分析，一邊在這份知識庫裡寫入或回饋，且不跟遠端脫節。

- 多資料夾工作區：在 Cursor 同時打開業務專案與本 repository，兩邊檔案、終端、Git 狀態都能在同一次工作階段檢閱。
- 遠端同步節奏：開始改內容前在本 repo 執行 `git pull`；改完依 Git 規則 commit / push。
- Cursor 裡完成 Git：可用 Source Control 或整合終端執行同等 git 指令。
- 本機路徑與提示詞：規則裡的 `<AI_SKILL_REPO>` 請指到實際 clone 路徑；路徑固定、工作區內含本 repo 時最不容易錯位。
- 多裝置：Cursor/VSCode Settings Sync 不會同步這份 git 知識庫；內容仍靠 `git pull` / `git push`。

## Cursor 與對話目標閉環

工具中立規則見 [`enforcement/conversation-goal-ledger.md`](../../enforcement/conversation-goal-ledger.md)。Cursor 只是其中一種操作環境；goal ledger 的真相來源仍是業務專案本地的 `<PROJECT_ROOT>/.agent-goals/`，不要放在 `.cursor/`，也不要把 goal 檔 commit。

**Goal ledger 操作流程已由 runtime 管理**，請參考：
- [`runtime/runtime.db`](../../runtime/runtime.db) — `phase_machine` / `obligation_ledger` / `blocking_gates` 快速查詢
- [`runtime/runtime.db`](../../runtime/runtime.db) — phase / obligation / gate / transaction / recovery 的 embedded source
- `ai-skill goals`（source: [`scripts/ai-skill-cli/internal/app/goals.go`](../../scripts/ai-skill-cli/internal/app/goals.go)）— goal ledger CLI helper

Cursor 專屬注意事項：
- 可透過 hooks 自動檢查 goal ledger（sessionStart / preCompact / stop）
- `.agent-goals/` 不應放在 `.cursor/` 目錄下

## Cursor 與知識更新流程 Checkpoint

工具中立規則見 [`governance/lifecycle/knowledge-update-flow.md`](../../governance/lifecycle/knowledge-update-flow.md)。**快速路徑**請優先查詢 runtime.db：[`runtime/runtime.db`](../../runtime/runtime.db)（`generated_surfaces` 表）。Cursor 可以透過 `.cursor/rules/*.mdc`（alwaysApply）加入 checkpoint 提醒，並可選用 hooks 輔助。

**Knowledge update flow 已由 runtime 管理**，請參考：
- `runtime/runtime.db → generated_surfaces (type='knowledge_update_phases')` — 11 個步驟的結構化記錄（快速路徑）
- `runtime/runtime.db → recovery_strategies / phase_reconciliation / state_repair` — runtime recovery strategy（快速路徑）
- [`runtime/runtime.db`](../../runtime/runtime.db) — recovery / transaction state machine 的 source

Cursor 專屬注意事項：
- 可在 `.cursor/rules/*.mdc` 中加入 checkpoint 提醒（alwaysApply）
- 可選用 hooks 輔助（sessionStart / preCompact / stop），fail-open 避免阻塞正常工作
- Hook script 只應檢查 `<AI_SKILL_REPO>` 的 git status 和 `.agent-goals/` 狀態並提醒，除非團隊明確要求自動執行

← [回到 AI 工具索引](../README.md)
