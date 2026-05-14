# Claude 使用說明

本檔說明 Claude Code 特有的配置與操作注意事項。通用配置原則見 [`ai-tools/README.md`](../README.md)；知識庫入口見 [`README.md`](../../README.md)；啟動流程見 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)。

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

工具中立規則見 [`enforcement/conversation-goal-ledger.md`](../../enforcement/conversation-goal-ledger.md)。Claude Code 是 CLI 工具，沒有 hooks 機制，但可以透過 `CLAUDE.md` 中的 Custom Instructions 和操作注意來實作對話目標閉環。

### 在 `CLAUDE.md` 中加入 goal ledger 提醒

本知識庫的 [`CLAUDE.md`](../../CLAUDE.md) 已包含以下內容：

```text
開始工作前，若 `<PROJECT_ROOT>/.agent-goals/` 存在，先讀取確認是否有未完成的 active goal。
若 goal 標示 `single-owner` 或 `non-parallelizable`，不要和其他 agent 分工同一流程；先取得使用者確認。
完成目標後，若仍有長期 roadmap 或治理狀態，先回寫到 durable planning 文件，再刪除 active goal。
```

### 建議操作

在 Claude Code 開始可中斷、可拆解或多目標工作時，或已看到 active project 有 modified / staged / untracked files、已建立 TodoWrite、使用者說「繼續」前一個多步驟任務時：

1. 讀取 `<PROJECT_ROOT>/.agent-goals/`，確認是否已有 active / blocked / needs-validation goal，以及 priority、owner、lock、parallelization mode、plan/todo links、missing/decision/strengthen。若使用者說 agent 中斷、突然關閉、要從哪裡重做、剩下什麼或下一步是什麼，必須先讀 `.agent-goals/README.md` 與對應 active goal，再用 transcripts、terminal output、git status 交叉確認；不要把 transcript/git 當成第一真相來源。
2. 若沒有 ledger 且任務不是單一回覆即可完成，使用本庫 helper 初始化；不要因為已有 TodoWrite 就跳過 goal ledger：

   ```bash
   <AI_SKILL_REPO>/scripts/agent-goals.sh --project <PROJECT_ROOT> init
   ```

3. 建立或更新本輪主要目標：

   ```bash
   <AI_SKILL_REPO>/scripts/agent-goals.sh --project <PROJECT_ROOT> start \
     --id P1-short-goal \
     --title "Short goal title" \
     --source "User request summary" \
     --parallelization single-owner \
     --next "Next concrete action" \
     --criteria "Observable completion condition"
   ```

4. 若使用者轉移目標，先 `pause` 或 `update --status superseded` 舊 goal，再建立新的 `P1`。
5. 若有 planning 文件或 TodoWrite todo，使用 `--plan` / `--todo` 連到 goal，並讓 `.agent-goals/README.md` 的主目標表可快速跳回該 goal。
6. 若發現需要拆小目標，使用 `split` 或在 goal 檔的 `Subgoals` 區塊記錄；若發現不能分工或需單一 owner，使用 `--parallelization single-owner|non-parallelizable` 更新。
7. 在回覆完成前，只有完成條件與驗證都成立時才 `complete --validated`；條件已成立時必須同輪刪除 goal 並刷新 `.agent-goals/README.md`，不要把 `completed` row 留在 active 表。若該 goal 完成後仍有長期 roadmap、phase、migration、promotion、deprecation 或治理狀態，先回寫到 durable planning 文件，再刪除 active goal。否則保留 goal，讓下一個 agent 可接手。

## Claude Code 與知識更新流程 Checkpoint

工具中立規則見 [`governance/lifecycle/knowledge-update-flow.md`](../../governance/lifecycle/knowledge-update-flow.md)。Claude Code 是 CLI 工具，沒有 hooks 機制，但可以透過 `CLAUDE.md` 中的 Custom Instructions 加入 checkpoint 提醒。

### 在 `CLAUDE.md` 中加入 checkpoint 提醒

本知識庫的 [`CLAUDE.md`](../../CLAUDE.md) 應加入以下內容（在 goal ledger 提醒之後）：

```text
## 知識更新流程 Checkpoint

每輪工作結束前、切回長時間專案工作前、或使用者說「繼續」展開下一輪前，必須執行知識更新檢查：

1. 讀取 [`<AI_SKILL_REPO>/governance/lifecycle/knowledge-update-flow.md`] 了解完整流程。
2. 自問：本輪是否新增可重用技巧、validation rule、replay knob、hook/runner guard、錯誤模式、或閉環缺口？
3. 若是，依 knowledge-update-flow.md 的 11 個步驟執行：
   - Step 1-2：觸發檢查 + 分類知識類型
   - Step 3：決定 Promotion Target（intelligence / workflow / analysis / shared-rules / runtime / memory）
   - Step 4：寫入 feedback/history/<domain>/<category>/ lesson（寫入前依 sanitization.md 去敏）
   - Step 5：更新目標層
   - Step 6-7：選擇性執行 Intelligence Extraction 或 Failure Learning
   - Step 8：執行 Linked Updates
   - Step 9：更新 Runtime Surfaces
   - Step 10：驗證（diff review、去敏檢查、link check）
   - Step 11：Commit / Push / Readback（關閉 writeback transaction）
4. 若否，簡短說明本輪只有 project-specific evidence 或尚未達可泛化標準。
```

### 與 Cursor 的差異

| 特性 | Cursor | Claude Code |
|------|--------|-------------|
| 自動提醒機制 | `.cursor/hooks.json`（sessionStart / preCompact / stop） | 無 hooks，需在 `CLAUDE.md` 中手動提醒 |
| Goal ledger 操作 | 可透過 hooks 自動檢查 | 需在每個 session 開始時手動讀取 |
| 知識更新流程 Checkpoint | 可在 `.cursor/rules/*.mdc` 中加入，可選 hooks 輔助 | 需在 `CLAUDE.md` 中加入 |
| 設定位置 | `.cursor/rules/*.mdc` + hooks.json | `CLAUDE.md`（Custom Instructions） |

## 與 Tool Adapter 的關係

若某個 skill 針對 Claude 有特殊執行策略（上下文載入順序、prompt chunking、工具輸出限制等），放在：

```
skills/<skill-name>/tool-adapters/claude.md
```

該 adapter 只寫 skill-specific 差異，並連回核心 `WORKFLOW.md` / `TOOLS.md`。

## 語言偏好設定（重要）

Claude Code 的語言偏好設定方式與 Roo Code（VS Code Extension）不同，因為 Claude Code 是 CLI 工具，**沒有** SQLite 全域資料庫。為了完整解決語言漂移問題，需要在語言偏好設定中加入**語言一致性強制規則**。

### 設定方式

Claude Code 的語言行為由 `CLAUDE.md` 中的 Custom Instructions 控制：

1. **在 `CLAUDE.md` 中設定**：本知識庫的 [`CLAUDE.md`](../CLAUDE.md) 已包含語言偏好設定。
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

- 讀了哪些 shared rules 與 skill 依賴。
- 哪些依賴不存在（標示 `not applicable`）。
- 目標是否完成，還有哪些 `.agent-goals` 未完成。
- 驗證方法：diff review、link check、commit/push/readback/clean status。

← [回到 AI 工具索引](../README.md)
