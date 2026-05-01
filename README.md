# Ai-skill

這個 repository 是**專業共用**的 AI skill 知識庫，一般透過 **Git 與遠端同步**（與團隊或跨機器共用流程一致）。

未來不同專案需要 APK 分析、抓包、Frida、Proxyman、Dart AOT、解密或分析方法沉澱時，優先從這裡讀 skill，並把新的可重用技巧回饋回這裡。

**路徑約定：**在本機上，中央庫位置就是你心裡的 `<AI_SKILL_REPO>`（clone 所在目錄）。**請勿把真實本機絕對路徑寫進本庫任何會 commit 的文件**（含 `feedback_history`、`shared-rules`、規則範本）；對外一律用占位符 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`。

## 現有 Skills

| Skill | 用途 |
| --- | --- |
| `skills/apk-analysis/` | 授權 APK 流量分析、動態抓包、Flutter/Dart AOT、response 解碼、文件化與 `feedback_history/` 技巧條目。 |

**目錄約定：**  
- **`skills/`**：各 skill 技巧包（目前為 `apk-analysis/`）；之後新增 skill 放在 `skills/<name>/`，步驟見 [`skills/ADDING_SKILLS.md`](skills/ADDING_SKILLS.md)。  
- **`shared-rules/`**：**共用規則**（依主題分檔：授權、去敏、內容分層、**feedback 檔名／模板／agent 行為**、Cursor 同步等）；索引為 [`shared-rules/README.md`](shared-rules/README.md)，feedback 流程與模板集中在 [`feedback-lessons.md`](shared-rules/feedback-lessons.md)。各 skill 目錄下 **`FEEDBACK.md`** 若存在，僅為**一行入口**，不必重複維護正文。  
- **每一則 `feedback_history`**：**不要**重複貼上共用規則全文，頂部引用 `shared-rules/` 即可。**同步到 `.cursor`** 時：先複製 **`shared-rules/`** 整包，再同步 **`skills/apk-analysis/`**（見 [`shared-rules/cursor-sync.md`](shared-rules/cursor-sync.md)）。

## 新專案啟用 apk-analysis（Cursor）

目標：**新開業務專案時，讓 Cursor 容易辨識並套用 apk-analysis**；分析流程與新技巧仍以**中央庫**為準，不要散落在各專案。

### 1. 中央庫先行（網路同步）

在新專案開始前或開始時，在本機中央庫目錄執行 `git pull`，確保 `skills/`、`shared-rules/` 與遠端一致。

### 2. 讓 Cursor「看得到」skill 的兩種位置（擇一或並用）

Cursor 會掃描特定路徑下的 skill；把中央庫 **`skills/apk-analysis/`**（內含 `SKILL.md` 等）放到下列**其一**，提到 APK 分析、Frida、Proxyman、Dart AOT、抓包、解密等較容易自動對應到這份 skill：

| 位置 | 用途 |
| --- | --- |
| `<PROJECT_ROOT>/.cursor/skills/apk-analysis/` | **專案內**：只有這個 repo 開工作區時也會載入；可進業務專案 git（若以複製方式放入，要記得與中央庫同步策略）。 |
| `~/.cursor/skills/apk-analysis/` | **本機共用**：所有專案共用一份，不必每個專案複製；適合單人固定一台機器開發。 |

**本機若想「共用規則與 skill 都放在 `bundles/`」**（與 `~/.cursor` 其他內容分流）：用 **`~/.cursor/bundles/shared-rules`**（連到本庫 `shared-rules/`）與 **`~/.cursor/bundles/ai-skill/`**（各 skill），再讓 `~/.cursor/shared-rules`、`~/.cursor/skills/*` 指向上述路徑；本庫提供 **`scripts/sync-cursor-bundle.sh`**（見 [`scripts/README.md`](scripts/README.md) 與 [`shared-rules/cursor-sync.md`](shared-rules/cursor-sync.md)）。

**資料從哪來：**從中央庫的 **`skills/apk-analysis/`** **整包複製**過去，或對該目錄做 **symbolic link** 指到 `<AI_SKILL_REPO>/skills/apk-analysis`。另請依 [`shared-rules/cursor-sync.md`](shared-rules/cursor-sync.md) 把 **`shared-rules/`** 一併部署（專案內可放 `<PROJECT_ROOT>/.cursor/shared-rules/`；本機建議搭配 bundle 腳本），Agent 才讀得到分類後的共用規則。

若**只複製部分 skill 檔案**，仍須帶上 **`SKILL.md`** 並另外同步 **`shared-rules/`**（至少含索引與 [`feedback-lessons.md`](shared-rules/feedback-lessons.md)），否則缺共用底線。

### 3. 最穩用法（建議每次開案或開長對話時做一次）

僅把檔案放在 `skills/` 目錄**不等於** Agent 永遠會依你想要的順序執行。請**明講**要使用 apk-analysis，並指定讀中央庫的流程檔（路徑用占位符，實際開檔靠本機 `<AI_SKILL_REPO>` 或多根工作區）：

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

（若工作區已用「多資料夾」同時打開業務專案與中央庫，Agent 直接開檔最穩。）

### 4. 與「真相來源」的關係

- **流程與泛用技巧**：以中央庫 **`skills/apk-analysis/`** 為準；**共用規則**以 **`shared-rules/`** 為準；回饋與修訂在此 repo **commit / push**。  
- **`~/.cursor/skills` 或 `.cursor/skills`**：目的是讓 Cursor **更容易套用** skill；若採複製，須與中央庫更新節奏對齊，詳見下一節。

## 回饋規則

可以回饋：

- 新的分析流程。
- 新的工具使用技巧。
- 新的失敗判讀方式。
- 新的去敏與 fixture 沉澱方式。
- 通用的媒體、解密、session、proxy、Frida 經驗。

不要回饋：

- 特定產品的完整 host / endpoint。
- token、secret、device id、帳號或個資。
- 未去敏 raw response。
- 本機絕對路徑、使用者名稱、私有工作目錄、clone 位置；請改用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>`。
- 只對單一產品有效、沒有泛化價值的結論。

## Skill 與 `.cursor`：真相來源與同步策略

核心問題：**若把 skill「複製」進業務專案的 `.cursor`，複製品不會自動知道外層（本 repository）已更新。** Cursor 也不會替你比對兩份檔案；除非你改成「只維護一份實體」或「固定重跑同步步驟」。

建議先決定 **單一真相來源（single source of truth）**：

| 策略 | 做法 | 優點 | 注意 |
| --- | --- | --- | --- |
| **參照（建議）** | `.cursor` 裡只放**短規則**：要求 Agent **一律先讀** `<AI_SKILL_REPO>/shared-rules/README.md`、`skills/apk-analysis/SKILL.md`（及 RUNBOOK 等）。工作區用「多資料夾」同時打開業務專案 + 本 repo。 | 永遠讀到同一份檔案；`git pull` 本庫即更新技巧與共用規則。 | 必須能開到本庫路徑。 |
| **符號連結** | 將 `.cursor/skills/apk-analysis` **連結到**本庫的 `skills/apk-analysis`；**另**將 `.cursor/shared-rules` **連結或複製**自本庫 `shared-rules`（見 [`cursor-sync.md`](shared-rules/cursor-sync.md)）。 | skill 與共用規則可依連結各別處理。 | `shared-rules` 與 `skills` 通常要**分開**佈署（連結 skill 不會自動帶上上一層目錄）。 |
| **複製** | **`shared-rules/`** 整包 → `.cursor/shared-rules/`；**`skills/apk-analysis/`** 整包 → `.cursor/skills/apk-analysis/`（順序見 [`cursor-sync.md`](shared-rules/cursor-sync.md)）。若只拷片段，至少帶 `shared-rules/README.md`、`feedback-lessons.md` 與 `SKILL.md`。 | 離線快照可行。 | 每次 `pull` 後需重跑同步；否則 `.cursor` 過期。 |

**「公用更新流程」— 以本 repository 為準（所有策略共用前半段）**

1. 在 `<AI_SKILL_REPO>`：`git pull`（若與他人共用或換機）。
2. **只在本庫**編輯：`shared-rules/` 內各分類檔（含 `feedback-lessons.md`）、`skills/apk-analysis/SKILL.md`、`RUNBOOK.md`、`DOCUMENTATION.md` 等（遵守上方回饋規則）。
3. 依你選的策略：
   - **參照**：通常不需改 `.cursor`；若新增「必讀檔案路徑」或專案特例，才改業務專案 `.cursor/rules` 補一句。
   - **符號連結**：確認連結仍指向正確本庫路徑；本庫 commit 後無須再複製檔案。
   - **複製**：從本庫**單向**同步到 `.cursor`（永遠以本庫為準覆寫快照，避免在複製品上長期單獨編輯）。
4. 在 `<AI_SKILL_REPO>`：`git add` → `commit` → `push`。
5. 若業務專案的 `.cursor` 有變更（規則文字或複製內容）：在該專案 git **另行** `commit`／`push`。

若要讓「複製策略」也能追蹤是否落後，可在同步後於 `.cursor` 內保留一行註記（例如同步日期或本庫 `git rev-parse --short HEAD` 的輸出）；真正的一致性仍靠流程與單向同步，無法只靠 Cursor 自動完成。

以上策略與「新專案啟用 apk-analysis」「在 Cursor 中連動更新」一併使用。**較省事且不漂移的組合**：`skills/apk-analysis` **symlink** 到 `~/.cursor/skills/apk-analysis` 或專案內對應路徑；**另** symlink 或複製 **`shared-rules/`** → `.cursor/shared-rules/`；再加 **多根工作區**與對話時 **明講必讀** `shared-rules/README.md`、`SKILL.md`、`RUNBOOK.md`。

## 在 Cursor 中連動更新

目標是：**一邊在業務專案裡分析、一邊在這份 skill 知識庫裡寫入或回饋，且不跟遠端脫節。**

- **多資料夾工作區（建議）**  
  在 Cursor 用「檔案 → 將資料夾新增至工作區…」同時打開你的業務專案與本 repository（`Ai-skill`）。兩邊的檔案、終端、Git 狀態都會一併可見，AI 改寫 `skills/apk-analysis/` 或 `shared-rules/` 時，你能在同一次編輯工作階段裡檢閱、commit。  
  儲存工作區（「檔案 → 將工作區儲存為…」）之後，下次一鍵還原同一套「專案 + 知識庫」佈局。

- **和遠端同步的節奏**  
  若多人共用或你在多台機器上 clone，在開始改 skill 前在本 repo 執行 `git pull`；改完用下方 Git 流程 `commit` / `push`。這樣本機、Cursor 裡看到的內容與遠端不會各走各的。

- **在 Cursor 裡完成 Git**  
  用左側「原始檔控制」（Source Control）可檢視變更、寫訊息、commit；`push` / `pull` 可從畫面操作或在整合終端執行 `git` 指令，與下方命令列等價。

- **本機路徑與提示詞**  
  規則裡的 `<AI_SKILL_REPO>` 請指到你實際 clone 的路徑；路徑固定、工作區內含本 repo 時，提示與實際開檔最不容易錯位。

- **多裝置時（選用）**  
  Cursor/VSCode 的「Settings Sync」只同步編輯器與帳戶相關設定，**不會**自動幫你同步這份 git 知識庫；內容仍靠 `git pull` / `git push` 或你慣用的 Git 用戶端。

- **`.cursor` 與本庫要一致（重要）**  
  你在 Cursor 裡給 Agent 的長效指引通常放在**業務專案**的 `.cursor`。可重用技巧以 **`skills/apk-analysis/`**、共用政策以 **`shared-rules/`** 為**真相來源**；`.cursor` 應**參照或同步**該來源，詳見上一節「Skill 與 `.cursor`」。  
  **每回沉澱或修正技巧時**：若採**參照**策略，先把變更寫進本庫再 commit／push；`.cursor` 只在補專案特例或調整「必讀路徑」時才改。若採**複製**策略，改完本庫後務必執行該節的同步步驟，再在業務專案 commit `.cursor`。  
  - 本庫：更新 `shared-rules/`、`SKILL.md`、`RUNBOOK.md`、`DOCUMENTATION.md` 等，然後在 `<AI_SKILL_REPO>` **commit 並 push**（見下方「Git 規則」）。  
  - 業務專案：依策略更新 `.cursor`（規則檔或複製快照），再在該專案 git **另行 commit**（若該專案有 remote，再 **push**）。

## Git 規則

這個 repository 有自己的 git remote。更新 skill 後：

```bash
cd <AI_SKILL_REPO>
git status
git add .
git commit -m "Update apk analysis skill"
git push
```

不要把專案私有資料、抓包原始檔或未去敏樣本 commit 到這裡。

`.cursor` 的變更若屬於業務專案目錄，請勿誤把該專案才有的機密或絕對路徑抄進本庫；本庫維持泛化與 `<AI_SKILL_REPO>` / `<PROJECT_ROOT>` 等占位符。
