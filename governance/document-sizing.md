# 文件大小與拆分原則

本規則定義文件何時該拆分、如何拆分。適用於本知識庫所有文件，也可跨專案通用。

## 核心原則

- **單一檔案只承載一個清楚目的**：入口、索引、流程、工具、模板、單一技巧或單一 lesson。一個檔案不該同時是教學、規範、範例和 checklist。
- **檔案變大 = token 成本變高**：每次 session agent 都需要讀取檔案內容。檔案越大，無關內容越多，token 浪費越嚴重。
- **拆分不是為了整齊，是為了降低 context 成本**：讓 agent 只讀到任務相關的內容，而不是整份文件。
- **不要為了拆分而拆分**：短小且高關聯的內容可以留在同檔。拆分也有成本（增加檔案數量、維護索引）。

## 權重與關係

本規則在 [`../enforcement/rule-weight.md`](../enforcement/rule-weight.md) 中屬於 **P2（Cross-repo operating policy）**。詳細的權重衝突處理請見 rule-weight.md 的「常見範例」表格。

Token 成本模型與 context loading 分層策略請見 [`../enforcement/decision-efficiency.md`](../enforcement/decision-efficiency.md) 的 Token 成本模型章節。

內容分層與跨專案適用說明請見 [`../enforcement/content-layering.md`](../enforcement/content-layering.md)。

## 何時要拆

符合任一情況時，優先考慮拆成資料夾與多檔：

### 內容層面

- 文件已明顯變長，新增內容只跟其中一小段有關。
- 一句或一段規則開始展開成多個步驟、例外、模板或範例。
- 同檔混合多種分類（例如 Flutter、HTTP、local proxy、media），且任務通常只需要其中一類。
- 寫 skill 的規範、工具教學、工作流、文件模板、feedback lesson 開始互相混在一起。
- agent 每次都需要讀大量無關內容才能找到當前任務的規則。

### Token 層面

- 檔案超過 **150 行**且內容主題不單一。
- 檔案超過 **300 行**（警戒線）。
- 該檔案在 routing registry 中被標記為高頻讀取（每次 session 都會載入），但實際只有部分內容是通用的。

### 維護層面

- 多人協作時，同一檔案頻繁發生 merge conflict。
- 檔案的修改歷史顯示不同時期加入的主題彼此無關。
- 檔案的「目錄」部分已經比實際內容還長。

## 決策流程

```
檔案是否超過 150 行？
├── 否 → 保持單檔
└── 是 → 檔案是否只有一個主題？
    ├── 是 → 保持單檔，但考慮是否可濃縮
    └── 否 → 需要拆分
        ├── 每個主題獨立成子檔
        ├── 建立 README.md 做目錄與路由
        ├── 更新所有引用此檔案的連結
        └── 執行連動更新（linked-updates.md）
```

## 拆分後必須做

- **父層 `README.md`** 要說明每個子檔何時讀，不要只列檔名。
- **舊連結要同步更新**；必要時保留短入口連到新位置。
- **若是全庫共用規則**，正文只放在 `governance/`，其他檔案引用它。
- **若改動影響索引、模板、skill 入口或分類文件**，依 [`../enforcement/linked-updates.md`](../enforcement/linked-updates.md) 同步檢查。
- **若改動影響 routing registry**，執行 `ruby scripts/refresh-knowledge-runtime.rb` 更新 runtime index。
- **若檔案在 routing registry 中被引用**，更新 registry 中的路徑（否則 validator 會報錯）。

在業務專案或其他 repository **從零撰寫 agent 友善文件**時的步驟、分類維度與驗證自查，見 [`../workflow/documentation/execution-flow.md`](../workflow/documentation/execution-flow.md)（本檔定閾值與拆分形狀，該檔定操作順序與與 `enforcement/` 的對齊方式）。

## 建議結構

任何主題或類別，建議拆分為：

```text
<domain>/
  README.md              # 目錄、路由、何時讀哪個子檔
  execution-flow.md      # 流程或決策樹
  artifact-gates.md      # 產出格式與完成定義
  quickstart.md          # 快速入門（選用）
```

若該主題需要更細的分類，可再往下分：

```text
<domain>/
  README.md
  <sub-category>/
    README.md
    <specific-flow>.md
```

本知識庫實際使用的三層分離範例（操作流程層 / 分析方法層 / 決策智慧層）請見 [`../enforcement/content-layering.md`](../enforcement/content-layering.md) 的「文件變大時」章節。

## 常見錯誤

| 錯誤 | 說明 | 正確做法 |
|------|------|----------|
| 拆分後不更新索引 | 檔案移走了，但 `README.md` 和 routing registry 還指向舊路徑 | 拆分後立即執行 linked-updates 檢查 |
| 拆分後父 README 只列檔名 | agent 不知道何時該讀哪個子檔 | 每個子檔要說明用途和觸發條件 |
| 為了拆分而拆分 | 50 行的檔案也拆成 3 個檔案，增加維護成本 | 只有超過 150 行或混合多主題時才考慮拆分 |
| 拆分後不更新 validator | routing registry 中的路徑失效，validation 失敗 | 執行 `ruby scripts/refresh-knowledge-runtime.rb` |
| 忽略 token 成本 | 不考慮 agent 每次讀取的成本 | 見 `decision-efficiency.md` 的 Token 成本模型 |

## 檢查點

本規則應在以下時間點被檢查：

| 檢查點 | 時機 | 檢查內容 |
|--------|------|----------|
| 新增文件後 | 建立任何 `.md` 檔案後 | 檔案是否超過 150 行？是否混合多主題？是否需要拆分？ |
| 大幅修改文件後 | 對既有檔案新增大量內容後 | 同上，並檢查既有拆分結構是否仍合理 |
| 文件合併或重組後 | 搬移、合併或重組目錄結構後 | 新結構是否符合拆分原則？父層 README 是否需更新？ |
| 任何 commit 前 | `git add` 前 | 快速掃描本次新增/修改的檔案行數，確認無遺漏拆分需求 |
| 知識提取 pipeline 完成後 | intelligence-extraction-pipeline.md 執行完 Step 5 後 | 提取產出的新檔案是否需要拆分？ |

← [回到 governance 索引](README.md)
