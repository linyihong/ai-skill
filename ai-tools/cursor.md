# Cursor 使用說明

本檔說明 Cursor 如何讀取與同步 Ai-skill 知識庫。分析流程與新技巧仍以中央庫為準，不要散落在各專案。

## 新專案啟用 apk-analysis

目標：新開業務專案時，讓 Cursor 容易辨識並套用 `apk-analysis`。

## 預設載入 shared rules

在 Cursor 中建議用 always-apply rule 或 sessionStart hook 先提醒 agent 載入 `shared-rules/README.md` 的 Default Bootstrap：

- `shared-rules/README.md`
- `shared-rules/dependency-reading.md`
- `shared-rules/linked-updates.md`
- `shared-rules/conversation-goal-ledger.md`
- `shared-rules/tool-neutral-documentation.md`
- `shared-rules/decision-efficiency.md`
- `shared-rules/failure-learning-system.md`
- `shared-rules/document-todo-list.md`
- `shared-rules/document-sizing.md`
- `shared-rules/goal-action-validation.md`
- `shared-rules/neutral-language.md`

Bootstrap 後仍要依任務讀 skill-specific README / WORKFLOW / TOOLS / DOCUMENTATION / CHECKLIST，以及 feedback、sanitization、authorization、cross-skill 等任務相關規則。
若使用者指出 agent 反覆失誤、更新到錯誤路徑、漏做 close-loop 或驗證不足，依 `failure-learning-system.md` 分類失效模式並沉澱成 failure pattern 或對應 skill lesson。

### 1. 中央庫先行

在新專案開始前或開始時，在本機中央庫目錄執行 `git pull`，確保 `skills/`、`shared-rules/` 與遠端一致。

### 2. 讓 Cursor 看得到 skill

預設策略改成 **reference-first**：專案 `.cursor` 先放短規則或提示，要求 Agent 直接讀 `<AI_SKILL_REPO>` 裡的 shared rules 與 skill 入口。這不需要把 skill 複製進每個專案，只要該中央庫路徑對目前 Cursor 工作區可讀即可。

這符合 repo-level 的 [`AI Native Operating System`](../architecture/ai-native-operating-system.md) 方向：Cursor bundle / copy snapshot 是相容層，不是預設模型。

若你還需要 Cursor 的原生 skill 掃描或本機全域載入，再把中央庫對應的 `skills/<name>/`（內含 `SKILL.md` 等）用 symlink 或同步腳本放到下列其一：

| 位置 | 用途 |
| --- | --- |
| `<PROJECT_ROOT>/.cursor/skills/<name>/` | 專案內：只有這個 repo 開工作區時也會載入；若是 symlink，可避免每次複製。 |
| `~/.cursor/skills/<name>/` | 本機共用：所有專案共用一份；建議指向中央庫或 bundle，不維護第二份正文。 |

本機若想把共用規則與 skill 都放在 `bundles/`，使用 `~/.cursor/bundles/shared-rules`（連到本庫 `shared-rules/`）與 `~/.cursor/bundles/ai-skill/`（各 skill），再讓 `~/.cursor/shared-rules`、`~/.cursor/skills/*` 指向上述路徑。本庫提供可選的 [`scripts/sync-cursor-bundle.sh`](../scripts/sync-cursor-bundle.sh) 建立這些 symlink；reference-only 時不需要執行。

資料來源優先順序：先用 `.cursor` 規則參照 `<AI_SKILL_REPO>`；需要原生掃描時用 symbolic link 指到 `<AI_SKILL_REPO>/skills/<name>` 與 `<AI_SKILL_REPO>/shared-rules`；只有在無法讀取中央庫或需要離線快照時才複製。另請依 [`shared-rules/cursor-sync.md`](../shared-rules/cursor-sync.md) 選擇 reference、symlink 或 copy snapshot，避免把 mirror 當成 source repo。

若只複製部分 skill 檔案，仍須帶上 `SKILL.md` 並另外同步 `shared-rules/`（至少含索引與 [`feedback-lessons.md`](../shared-rules/feedback-lessons.md)），否則缺共用底線。

### 3. 最穩用法

僅把檔案放在 `skills/` 目錄不等於 Agent 永遠會依你想要的順序執行。最穩是讓專案 `.cursor` 或開場提示明講要使用哪個 skill，並指定讀中央庫的流程檔：

```text
使用 apk-analysis skill。請先閱讀共用規則索引、feedback 格式與 skill 入口：
<AI_SKILL_REPO>/shared-rules/README.md
<AI_SKILL_REPO>/shared-rules/feedback-lessons.md
<AI_SKILL_REPO>/skills/apk-analysis/SKILL.md
<AI_SKILL_REPO>/skills/apk-analysis/RUNBOOK.md
之後依 WORKFLOW.md / TOOLS.md / DOCUMENTATION.md 進行分析（路徑皆在 skills/apk-analysis/）。
新技巧請依 shared-rules/feedback-lessons.md 寫入 skills/apk-analysis/feedback_history/（勿寫真實本機路徑或機密）。
完成後在 <AI_SKILL_REPO> commit 並 push。
若你同時改了業務專案裡的 .cursor，在該專案另行 commit／push。
```

若工作區已用多資料夾同時打開業務專案與中央庫，Agent 直接開檔最穩。

## Skill 與 `.cursor` 的真相來源

核心問題：若把 skill 複製進業務專案的 `.cursor`，複製品不會自動知道外層 repository 已更新。Cursor 也不會替你比對兩份檔案；除非你改成只維護一份實體，或固定重跑同步步驟。

**硬規則：** 使用者要求「同步 skill repo」、「回饋到 skill repo」或修正可重用 skill / shared rule 時，先定位並修改 `<AI_SKILL_REPO>` 這個 git repository。`~/.cursor/skills*`、`~/.cursor/shared-rules`、`~/.cursor/bundles/*` 與專案 `.cursor/` 都是 Cursor 可讀的部署 / mirror 路徑；除非已確認它們是 symlink 指回 `<AI_SKILL_REPO>`，否則不能把這些路徑當成 source repo，也不能只改這些路徑就回覆「已同步」。

建議先決定單一真相來源：

| 策略 | 做法 | 優點 | 注意 |
| --- | --- | --- | --- |
| 參照（預設） | `.cursor` 裡只放短規則：要求 Agent 一律先讀 `<AI_SKILL_REPO>/shared-rules/README.md`、`skills/apk-analysis/SKILL.md`（及 RUNBOOK 等）。工作區用多資料夾同時打開業務專案與本 repo。 | 永遠讀到同一份檔案；`git pull` 本庫即更新技巧與共用規則；不需每個專案複製。 | 必須能開到本庫路徑；若換機，先 clone / pull 中央庫並更新 `<AI_SKILL_REPO>` 指向。 |
| 符號連結 | 將 `.cursor/skills/apk-analysis` 連結到本庫的 `skills/apk-analysis`；另將 `.cursor/shared-rules` 連結到本庫 `shared-rules` 或 bundle。 | Cursor 原生掃描可看到 skill，同時仍只有一份正文。 | `shared-rules` 與 `skills` 通常要分開佈署；避免在 repo 內建立反向 symlink。 |
| 複製快照 | `shared-rules/` 整包到 `.cursor/shared-rules/`；`skills/apk-analysis/` 整包到 `.cursor/skills/apk-analysis/`。 | 離線、不能 symlink、或工具限制時可行。 | 每次 `pull` 後需重跑同步；最好記錄來源 commit hash，否則 `.cursor` 容易過期。 |

## 公用更新流程

以本 repository 為準，所有策略共用前半段：

1. 在 `<AI_SKILL_REPO>` 執行 `git pull`（若與他人共用或換機）。
2. 只在本庫編輯 `shared-rules/`、`skills/apk-analysis/SKILL.md`、`RUNBOOK.md`、`DOCUMENTATION.md` 等。
3. 依 [`shared-rules/linked-updates.md`](../shared-rules/linked-updates.md) 檢查連動更新。
4. 依你選的策略處理 `.cursor`：reference-only 不必複製，只需確認 `<AI_SKILL_REPO>` 可讀；symlink / copy snapshot 才需要同步。
5. 在 `<AI_SKILL_REPO>` 執行 `git add`、`git commit`、`git push`。
6. 若業務專案的 `.cursor` 有變更，在該專案 git 另行 commit / push。

若要讓複製策略也能追蹤是否落後，可在同步後於 `.cursor` 內保留同步日期或本庫 commit hash；真正一致性仍靠流程與單向同步。

## Cursor 連動更新

目標是：一邊在業務專案裡分析，一邊在這份 skill 知識庫裡寫入或回饋，且不跟遠端脫節。

- 多資料夾工作區：在 Cursor 同時打開業務專案與本 repository，兩邊檔案、終端、Git 狀態都能在同一次工作階段檢閱。
- 遠端同步節奏：開始改 skill 前在本 repo 執行 `git pull`；改完依 Git 規則 commit / push。
- Cursor 裡完成 Git：可用 Source Control 或整合終端執行同等 git 指令。
- 本機路徑與提示詞：規則裡的 `<AI_SKILL_REPO>` 請指到實際 clone 路徑；路徑固定、工作區內含本 repo 時最不容易錯位。
- 多裝置：Cursor/VSCode Settings Sync 不會同步這份 git 知識庫；內容仍靠 `git pull` / `git push`。
- `.cursor` 與本庫一致：可重用技巧以 `skills/`、共用政策以 `shared-rules/` 為真相來源；`.cursor` 應參照或同步該來源。

## Cursor 與對話目標閉環

工具中立規則見 [`shared-rules/conversation-goal-ledger.md`](../shared-rules/conversation-goal-ledger.md)。Cursor 只是其中一種操作環境；goal ledger 的真相來源仍是業務專案本地的 `<PROJECT_ROOT>/.agent-goals/`，不要放在 `.cursor/`，也不要把 goal 檔 commit。

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
7. 在回覆完成前，只有完成條件與驗證都成立時才 `complete --validated`；條件已成立時必須同輪刪除 goal 並刷新 `.agent-goals/README.md`，不要把 `completed` row 留在 active 表。否則保留 goal，讓下一個 agent 可接手。

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
    "sessionStart": [
      {
        "command": ".cursor/hooks/goal-ledger-reminder.sh",
        "timeout": 5,
        "failClosed": false
      }
    ],
    "preCompact": [
      {
        "command": ".cursor/hooks/goal-ledger-reminder.sh",
        "timeout": 5,
        "failClosed": false
      }
    ],
    "stop": [
      {
        "command": ".cursor/hooks/goal-ledger-reminder.sh",
        "timeout": 5,
        "failClosed": false
      }
    ]
  }
}
```

Hook script 只應讀 `.agent-goals/` 並提醒，除非團隊明確要求自動修改 goal 檔。若 hook 需要修改 goal，仍要遵守 `conversation-goal-ledger.md` 的 lock、TTL、完成驗證與刪除條件。

← [回到 AI 工具索引](README.md)
