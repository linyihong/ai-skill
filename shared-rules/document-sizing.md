# 文件大小與拆分（通用規則）

本規則適用於本庫所有文件，也**可跨專案通用**。目標是讓 agent 只讀到任務相關內容，避免單一 Markdown 檔累積成難維護的大雜燴。無論專案使用何種目錄結構，拆分原則與建議結構皆可直接套用。

## 原則

- 單一檔案只承載一個清楚目的：入口、索引、流程、工具、模板、單一技巧或單一 lesson。
- 當檔案開始同時承載多個主題，或需要讀完整檔才找得到某個小規則時，應拆分。
- 拆分後用資料夾包裝：`README.md` 做目錄與路由，子檔放具體內容。
- 不要為了拆分而拆分；短小且高關聯的內容可以留在同檔。

## 何時要拆

符合任一情況時，優先考慮拆成資料夾與多檔：

- 文件已明顯變長，新增內容只跟其中一小段有關。
- 一句或一段規則開始展開成多個步驟、例外、模板或範例。
- 同檔混合多種分類，例如 Flutter、HTTP、local proxy、media，且任務通常只需要其中一類。
- 寫 skill 的規範、工具教學、工作流、文件模板、feedback lesson 開始互相混在一起。
- agent 每次都需要讀大量無關內容才能找到當前任務的規則。

## 建議結構（通用）

任何主題或類別，建議拆分為：

```text
<domain>/
  README.md        # 目錄、路由、何時讀哪個子檔
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

### 本庫實際分層範例（可作為其他專案參考）

本庫使用三層分離架構，各層獨立拆分：

```text
# 操作流程層（HOW TO DO）
workflow/<domain>/
  README.md
  execution-flow.md
  artifact-gates.md

# 分析方法層（HOW TO DO 的細節）
analysis/<domain>/
  README.md
  workflows/
    README.md
    <specific-flow>.md

# 決策智慧層（HOW TO THINK）
intelligence/<domain>/
  README.md
  heuristics/
    README.md
    <heuristic>.md
  anti-patterns/
    README.md
    <anti-pattern>.md
  signals/
    README.md
    <signal>.md
  failure/
    README.md
    <failure>.md
```

### 舊結構（向後相容參考，僅限尚未遷移的 skill）

若專案仍使用傳統 skill 結構，常見模式（注意：已遷移 skill 的 `WORKFLOW.md`、`TOOLS.md`、`DOCUMENTATION.md`、`techniques/` 已被刪除，不再保留）：

```text
skills/<skill-name>/
  README.md
  SKILL.md
  WORKFLOW.md
  TOOLS.md
  DOCUMENTATION.md
  techniques/
    README.md
    <category>/
      README.md
  feedback_history/
    README.md
    <category>/
      README.md
      YYYY-MM-DD_HHMMSS-<slug>.md
```

## 拆分後必須做

- 父層 `README.md` 要說明每個子檔何時讀，不要只列檔名。
- 舊連結要同步更新；必要時保留短入口連到新位置。
- 若是全庫共用規則，正文只放在 `shared-rules/`，其他檔案引用它。
- 若改動影響索引、模板、skill 入口或分類文件，依 [`linked-updates.md`](linked-updates.md) 同步檢查。

← [回到共用規則索引](README.md)
