# Cursor 使用說明

本檔說明 Cursor 如何讀取與同步 Ai-skill 知識庫。分析流程與新技巧仍以中央庫為準，不要散落在各專案。

## 新專案啟用 apk-analysis

目標：新開業務專案時，讓 Cursor 容易辨識並套用 `apk-analysis`。

### 1. 中央庫先行

在新專案開始前或開始時，在本機中央庫目錄執行 `git pull`，確保 `skills/`、`shared-rules/` 與遠端一致。

### 2. 讓 Cursor 看得到 skill

Cursor 會掃描特定路徑下的 skill；把中央庫對應的 `skills/<name>/`（內含 `SKILL.md` 等）放到下列其一：

| 位置 | 用途 |
| --- | --- |
| `<PROJECT_ROOT>/.cursor/skills/<name>/` | 專案內：只有這個 repo 開工作區時也會載入；可進業務專案 git。 |
| `~/.cursor/skills/<name>/` | 本機共用：所有專案共用一份，不必每個專案複製。 |

本機若想把共用規則與 skill 都放在 `bundles/`，使用 `~/.cursor/bundles/shared-rules`（連到本庫 `shared-rules/`）與 `~/.cursor/bundles/ai-skill/`（各 skill），再讓 `~/.cursor/shared-rules`、`~/.cursor/skills/*` 指向上述路徑。本庫提供 [`scripts/sync-cursor-bundle.sh`](../scripts/sync-cursor-bundle.sh)。

資料來源：從中央庫的 `skills/<name>/` 整包複製過去，或對該目錄做 symbolic link 指到 `<AI_SKILL_REPO>/skills/<name>`。另請依 [`shared-rules/cursor-sync.md`](../shared-rules/cursor-sync.md) 把 `shared-rules/` 一併部署，Agent 才讀得到分類後的共用規則。

若只複製部分 skill 檔案，仍須帶上 `SKILL.md` 並另外同步 `shared-rules/`（至少含索引與 [`feedback-lessons.md`](../shared-rules/feedback-lessons.md)），否則缺共用底線。

### 3. 最穩用法

僅把檔案放在 `skills/` 目錄不等於 Agent 永遠會依你想要的順序執行。請明講要使用 `apk-analysis`，並指定讀中央庫的流程檔：

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

建議先決定單一真相來源：

| 策略 | 做法 | 優點 | 注意 |
| --- | --- | --- | --- |
| 參照（建議） | `.cursor` 裡只放短規則：要求 Agent 一律先讀 `<AI_SKILL_REPO>/shared-rules/README.md`、`skills/apk-analysis/SKILL.md`（及 RUNBOOK 等）。工作區用多資料夾同時打開業務專案與本 repo。 | 永遠讀到同一份檔案；`git pull` 本庫即更新技巧與共用規則。 | 必須能開到本庫路徑。 |
| 符號連結 | 將 `.cursor/skills/apk-analysis` 連結到本庫的 `skills/apk-analysis`；另將 `.cursor/shared-rules` 連結或複製自本庫 `shared-rules`。 | skill 與共用規則可依連結各別處理。 | `shared-rules` 與 `skills` 通常要分開佈署。 |
| 複製 | `shared-rules/` 整包到 `.cursor/shared-rules/`；`skills/apk-analysis/` 整包到 `.cursor/skills/apk-analysis/`。 | 離線快照可行。 | 每次 `pull` 後需重跑同步；否則 `.cursor` 過期。 |

## 公用更新流程

以本 repository 為準，所有策略共用前半段：

1. 在 `<AI_SKILL_REPO>` 執行 `git pull`（若與他人共用或換機）。
2. 只在本庫編輯 `shared-rules/`、`skills/apk-analysis/SKILL.md`、`RUNBOOK.md`、`DOCUMENTATION.md` 等。
3. 依 [`shared-rules/linked-updates.md`](../shared-rules/linked-updates.md) 檢查連動更新。
4. 依你選的策略同步或參照 `.cursor`。
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

← [回到 AI 工具索引](README.md)
