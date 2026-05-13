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
<AI_SKILL_REPO>/shared-rules/README.md
之後依 workflow/<domain>/execution-flow.md 進行。
完成後在 <AI_SKILL_REPO> commit 並 push。
```

若工作區已用多資料夾同時打開業務專案與中央庫，Agent 直接開檔最穩。

## 語言偏好設定

Cursor 的語言行為由 `.cursor/rules/*.mdc` 中的規則控制：

1. **在 `.cursor/rules/` 中設定**：本知識庫的 `.cursor/rules/` 目錄已包含語言偏好設定。
2. **語言偏好內容**（在 `alwaysApply` 的規則中）：

```text
Language Preference: Default to English, but always match the user's language in conversation.
If the user writes in Chinese, respond in Chinese.
If the user writes in Japanese, respond in Japanese.
If the user switches languages, follow their switch.
```

### 與 Roo Code 的差異

| 特性 | Cursor | Roo Code |
|------|--------|----------|
| 執行環境 | 獨立 IDE（基於 VS Code） | VS Code extension |
| 設定位置 | `.cursor/rules/*.mdc`（檔案） | `.roomodes` + SQLite 全域資料庫 |
| 全域語言欄位 | 無公開 API 直接修改 | 有（`language` 欄位在 `state.vscdb`） |
| 設定方式 | 編輯 `.cursor/rules/` 中的規則 | 編輯 `.roomodes` + 修改 SQLite |

### 注意事項

- Cursor 沒有公開的 SQLite 全域資料庫修改方式，語言偏好完全由 `.cursor/rules/` 中的規則控制。
- 如果 Cursor 仍然強制使用英文，請檢查 `.cursor/rules/` 中是否有固定的語言偏好設定，改為上述軟性偏好即可。
- Cursor 的 `alwaysApply` 規則會在每個 session 啟動時自動載入，因此語言偏好設定放在 `alwaysApply: true` 的規則中最有效。

## 公用更新流程

以本 repository 為準：

1. 在 `<AI_SKILL_REPO>` 執行 `git pull`（若與他人共用或換機）。
2. 只在本庫編輯文件。
3. 依 [`shared-rules/linked-updates.md`](../../shared-rules/linked-updates.md) 檢查連動更新。
4. 在 `<AI_SKILL_REPO>` 執行 `git add`、`git commit`、`git push`。

## Cursor 連動更新

目標是：一邊在業務專案裡分析，一邊在這份知識庫裡寫入或回饋，且不跟遠端脫節。

- 多資料夾工作區：在 Cursor 同時打開業務專案與本 repository，兩邊檔案、終端、Git 狀態都能在同一次工作階段檢閱。
- 遠端同步節奏：開始改內容前在本 repo 執行 `git pull`；改完依 Git 規則 commit / push。
- Cursor 裡完成 Git：可用 Source Control 或整合終端執行同等 git 指令。
- 本機路徑與提示詞：規則裡的 `<AI_SKILL_REPO>` 請指到實際 clone 路徑；路徑固定、工作區內含本 repo 時最不容易錯位。
- 多裝置：Cursor/VSCode Settings Sync 不會同步這份 git 知識庫；內容仍靠 `git pull` / `git push`。

## Cursor 與對話目標閉環

工具中立規則見 [`shared-rules/conversation-goal-ledger.md`](../../shared-rules/conversation-goal-ledger.md)。Cursor 只是其中一種操作環境；goal ledger 的真相來源仍是業務專案本地的 `<PROJECT_ROOT>/.agent-goals/`，不要放在 `.cursor/`，也不要把 goal 檔 commit。

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

### Cursor hooks 範本方向

Cursor hook 可以輔助提醒，但不應成為唯一真相。可選的 project hook：

| Event | 用途 | 行為 |
| --- | --- | --- |
| `sessionStart` | 開局提醒 | 檢查 `.agent-goals/README.md` 與 `.agent-goals/goals/*.md`，若有 active goal，提醒目前未完成項、待決策、owner/lock、parallelization mode、優先順序與需要補強的地方。 |
| `preCompact` | 壓縮前檢查 | 若有 active goal，提醒 agent 先更新 `Next Action`、`Progress`、`Validation`。 |
| `stop` | 停止前檢查 | 若 goal 未完成，提醒保留或更新；不要自動刪除。若發現重疊 goal 被其他 owner/lock 處理，停止並請使用者決定。 |

Hook 建議使用 command hook，fail-open，避免 hook 故障阻塞正常工作。若要建立專案 hook，放在：

```text
<PROJECT_ROOT>/.cursor/hooks.json
<PROJECT_ROOT>/.cursor/hooks/goal-ledger-reminder.sh
```

範例 `hooks.json` 方向：

```json
{
  "version": 1,
  "hooks": {
    "sessionStart": [{ "command": ".cursor/hooks/goal-ledger-reminder.sh", "timeout": 5, "failClosed": false }],
    "preCompact": [{ "command": ".cursor/hooks/goal-ledger-reminder.sh", "timeout": 5, "failClosed": false }],
    "stop": [{ "command": ".cursor/hooks/goal-ledger-reminder.sh", "timeout": 5, "failClosed": false }]
  }
}
```

Hook script 只應讀 `.agent-goals/` 並提醒，除非團隊明確要求自動修改 goal 檔。若 hook 需要修改 goal，仍要遵守 `conversation-goal-ledger.md` 的 lock、TTL、完成驗證與刪除條件。

← [回到 AI 工具索引](../README.md)
