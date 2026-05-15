# Issue Triage

`analysis/issue/` 負責「Issue 分類與優先級判斷方法」。本目錄保存如何系統化地分類、優先級排序與初步診斷 issue 的分析框架，讓 agent 能快速判斷 issue 的性質、嚴重程度與下一步行動。

## 核心責任

- Issue 分類（bug / feature / improvement / question / incident）。
- 優先級評估（severity × urgency × impact）。
- 初步診斷流程（從 issue 描述到可重現步驟）。
- 標籤與路由規則（哪個 team、哪個 skill、哪個 workflow）。
- 重複 issue 偵測（避免重複工作）。

## 與其他層的關係

- `workflow/` 可引用本層的分類與路由規則，但不複製分析方法細節。
- `intelligence/engineering/failure/` 承接從 issue 中萃取的抽象化失敗模式。
- `intelligence/engineering/domain/` 承接從 issue 中萃取的領域知識。
- `skills/` 目前仍是相容入口；本層只承接逐步抽出的分析方法。

## 第一批候選遷移來源

- `enforcement/failure-learning-system.md` 中偏 issue triage 的方法。
- `plans/archived/2026-05-11-1112-next-stage-upgrade-plan.md` 中 `analysis/` 的分層說明。

## 建議分析方法

### 1. Issue 分類

```
1. 讀取 issue title 與 description。
2. 判斷類型：
   ├─ Bug：既有功能行為不符合預期。
   ├─ Feature：需要新功能或增強。
   ├─ Improvement：既有功能的優化（效能、UX、安全性）。
   ├─ Question：需要釐清或解釋。
   └─ Incident：正在影響 production 的緊急問題。
3. 判斷子類型（如 bug 的 crash / logic error / UI issue / performance regression）。
4. 標記相關領域（security / performance / UX / API / database 等）。
```

### 2. 優先級評估

```
Severity × Urgency × Impact = Priority

Severity:
  - critical：資料遺失、安全漏洞、全部服務不可用。
  - major：主要功能不可用、大量用戶受影響。
  - minor：次要功能問題、少量用戶受影響。
  - cosmetic：UI 問題、文件錯誤、不影響功能。

Urgency:
  - immediate：需要立即修復（production 正在受影響）。
  - high：需要在本次 sprint 修復。
  - medium：可以在下個 sprint 修復。
  - low：可以排入 backlog。

Impact:
  - wide：影響所有用戶或所有功能。
  - moderate：影響部分用戶或部分功能。
  - narrow：影響特定情境或少數用戶。
  - minimal：幾乎無實際影響。

Priority 計算：
  - P0：critical + (immediate or high) + (wide or moderate)
  - P1：critical + medium + wide / major + immediate + wide
  - P2：major + high + moderate / critical + low + narrow
  - P3：minor + medium + narrow / cosmetic + high + minimal
  - P4：cosmetic + low + minimal
```

### 3. 初步診斷流程

```
1. 確認 issue 是否可重現：
   ├─ 有提供重現步驟 → 嘗試重現。
   ├─ 無重現步驟 → 詢問 reporter 提供。
   └─ 無法重現 → 標記為 intermittent / unreproducible。

2. 收集環境資訊：
   ├─ 版本號（app version、OS version、browser version）。
   ├─ 裝置資訊（device model、screen size、network type）。
   └─ 用戶角色與權限。

3. 判斷是否與近期變更有關：
   ├─ 檢查最近 deploy / config change / dependency update。
   ├─ 檢查是否有相關的 regression test 失敗。
   └─ 檢查是否有相關的 monitoring alert。

4. 決定下一步：
   ├─ 可重現 + 有明確根因 → 指派修復。
   ├─ 可重現 + 無明確根因 → 需要深入分析（參考 analysis/production/）。
   ├─ 無法重現 + 高影響 → 加 monitoring 等待再次發生。
   └─ 無法重現 + 低影響 → 標記為 intermittent，排入 backlog。
```

### 4. 重複 Issue 偵測

```
1. 搜尋既有 issue 中相同或相似的 title / description。
2. 比對關鍵字（error message、stack trace、功能名稱）。
3. 比對 affected version / component。
4. 如果找到重複：
   ├─ 關閉新 issue，comment 指向既有 issue。
   └─ 如果既有 issue 已修復但新 issue 仍發生，可能是 regression。
```

### 5. 標籤與路由規則

```
1. 根據分類與領域決定標籤：
   ├─ type/bug、type/feature、type/improvement、type/question。
   ├─ area/security、area/performance、area/ux、area/api。
   ├─ severity/critical、severity/major、severity/minor、severity/cosmetic。
   └─ priority/p0、priority/p1、priority/p2、priority/p3、priority/p4。

2. 根據標籤決定路由：
   ├─ security issue → 安全團隊（或 security skill）。
   ├─ API bug → backend team（或 API workflow）。
   ├─ UI bug → frontend team（或 UI workflow）。
   └─ performance regression → performance team（或 analysis/production/）。
```

## 產出格式

每次 issue triage 應產出：

- **Issue 摘要**（≤100 tokens）：類型、severity、urgency、impact、priority。
- **初步診斷結果**（≤200 tokens）：是否可重現、環境資訊、與近期變更的關聯。
- **下一步行動**（≤100 tokens）：指派給誰、需要哪些深入分析、預計修復時間。
- **相關 issue 連結**（≤100 tokens）：重複 issue、相關 issue、已知根因。
