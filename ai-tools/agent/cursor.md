# Cursor 使用說明

本檔說明 Cursor 特有的配置與操作注意事項。通用配置原則見 [`ai-tools/README.md`](../README.md)；知識庫入口見 [`README.md`](../../README.md)；啟動流程見 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)。

## 預設載入（Core Bootstrap）

在 Cursor 中，`.cursor/rules/dependency-reading.mdc`（`alwaysApply: true`）已實作 Core Bootstrap 自動載入流程，指向 `CORE_BOOTSTRAP.md`。

**設定一次 repo 即可**：只要 clone 本 repo，Cursor 啟動時自動套用 `dependency-reading.mdc`（alwaysApply），不需要每次手動指定載入哪些規則。

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
  ├── 指向 Ai-skill 的 CORE_BOOTSTRAP.md（絕對路徑）
  ├── 語言偏好設定
  └── 語言一致性強制規則

專案 .cursor/rules/（層級 B，可選）
  ├── 只在需要專案特定規則時建立
  └── 加上該專案特有的規則
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

### 建議操作

在 Cursor 開始可中斷、可拆解或多目標工作時，或已看到 active project 有 modified / staged / untracked files、已建立 TodoWrite、使用者說「繼續」前一個多步驟任務時：

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

## Cursor 與知識更新流程 Checkpoint

工具中立規則見 [`governance/lifecycle/knowledge-update-flow.md`](../../governance/lifecycle/knowledge-update-flow.md)。Cursor 可以透過 `.cursor/rules/*.mdc`（alwaysApply）加入 checkpoint 提醒，並可選用 hooks 輔助。

### 在 `.cursor/rules/*.mdc` 中加入 checkpoint 提醒

在 `dependency-reading.mdc` 或其他 alwaysApply 規則檔中，加入以下內容：

```markdown
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

### Cursor hooks 輔助提醒（可選）

Cursor 的 hooks 機制可以用來輔助提醒 checkpoint，但不應取代規則檔中的強制提醒：

| Event | 用途 | 行為 |
| --- | --- | --- |
| `sessionStart` | 開局提醒 | 檢查是否有未完成的知識更新流程，提醒 agent 先完成 writeback transaction 再開始新工作。 |
| `preCompact` | 壓縮前檢查 | 若有未關閉的 writeback transaction，提醒 agent 先 commit/push 再壓縮。 |
| `stop` | 停止前檢查 | 若有 dirty Ai-skill repo 或未關閉的 transaction，提醒 agent 處理或記錄狀態。 |

Hook 建議使用 command hook，fail-open，避免 hook 故障阻塞正常工作。若要建立專案 hook，放在：

```text
<PROJECT_ROOT>/.cursor/hooks.json
<PROJECT_ROOT>/.cursor/hooks/knowledge-update-reminder.sh
```

範例 `hooks.json` 方向：

```json
{
  "version": 1,
  "hooks": {
    "sessionStart": [{ "command": ".cursor/hooks/knowledge-update-reminder.sh", "timeout": 5, "failClosed": false }],
    "preCompact": [{ "command": ".cursor/hooks/knowledge-update-reminder.sh", "timeout": 5, "failClosed": false }],
    "stop": [{ "command": ".cursor/hooks/knowledge-update-reminder.sh", "timeout": 5, "failClosed": false }]
  }
}
```

Hook script 只應檢查 `<AI_SKILL_REPO>` 的 git status 和 `.agent-goals/` 狀態並提醒，除非團隊明確要求自動執行。

← [回到 AI 工具索引](../README.md)
