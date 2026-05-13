# Travel Planning Quick Start（旅行規劃快速入門）

本文件提取 [`workflow/travel-planning/execution-flow.md`](../../workflow/travel-planning/execution-flow.md) 與 [`workflow/travel-planning/artifact-gates.md`](../../workflow/travel-planning/artifact-gates.md) 中 Quick Start 的操作步驟，作為 `runtime/onboarding/` 層的執行指引。

> **遷移狀態**：此文件為新分層的 reference target，舊 `skills/travel-planning/` 已不再作為 active entrypoint。新內容請直接寫入此文件。

## 快速入門步驟

### 步驟 1：釐清旅行框架

確認以下基本資訊：
- Destination（目的地）
- Dates（日期範圍）
- Party size（人數）
- Transport（交通方式）
- Pace（步調：緊湊/悠閒/混合）
- Budget（預算）
- Must-do interests（必去興趣點）
- Dietary/accessibility needs（飲食/無障礙需求）
- Lodging style（住宿風格：飯店/民宿/車中泊/混合）
- Tolerance for long drives or early starts（長途駕駛或早起的容忍度）

### 步驟 2：識別時間敏感檢查

確認以下時間敏感項目：
- Operating days（營業日）
- Reservation windows（預約窗口）
- Seasonal closures（季節性關閉）
- Event crowding（活動人潮）
- Weather forecast（天氣預報）
- Road/weather constraints（道路/天氣限制）
- Parking rules（停車規則）
- Public transport schedules（大眾運輸時刻表）
- Last-entry times（最後入場時間）

### 步驟 3：使用當前網路來源

- 優先使用官方設施、旅遊局、交通營運商、天氣、道路管理局、預約頁面
- 使用社群地圖或部落格進行探索性發現，然後在其他來源驗證細節

### 步驟 4：搜尋旅遊社套裝行程

搜尋旅遊社套裝行程、package tours 與官方 model course（同一區域/日期/季節）：

**直接推薦套裝**：
- 顯示價格、包含/排除項目、預約/取消注意事項、集合/出發點、交通假設
- 說明為什麼比自助規劃更簡單或更安全

**作為 Benchmark**：
- 提取路線順序、停留時間、季節亮點、用餐/交通模式、隱藏限制
- 對每個宣稱與官方/當前來源進行驗證

### 步驟 5：驗證地點精確度

- 優先使用 **Google Maps place link** 或 **coordinate pin**，指向一個精確地點
- 避免回傳多個可能點的通用搜尋 URL
- 駕駛路線使用 **最近的可確認訪客停車場**、**官方停車場** 或 **實用到達停車場** 作為導航目標（若與景點/餐廳入口不同）
- 交叉檢查地圖 pin 與官方名稱/地址
- 日本自駕行程：如可用，包含 **Mapcode**

### 步驟 6：每個停留點的建議

對每個推薦的停留點，提供：
- 做什麼、為什麼值得
- 預期停留時間
- 附近替代方案
- 當地美食或餐廳建議（如相關）

**餐廳建議**：
- 使用適合該國家的當地評論/評分工具 + Google Maps
- 日本：交叉檢查 Google Maps 與 **食べログ**
- 考慮：評分、評論數、近期度、營業時間、最後點餐、預約需求、價格範圍、排隊風險、停車/交通、路線適合度、附近備案

### 步驟 7：長距離交通比較

對跨城市、跨縣、島嶼、機場或 2+ 小時的移動：
- 在選擇主要路線前比較多種交通選項
- 包含 **door-to-door 時間** 與 **總費用**，不僅是票價

### 步驟 8：建立可行路線

- 加入移動緩衝
- 天氣感知排序：將戶外、景觀、渡輪、山區、步行較多的行程安排在最佳天氣窗口

### 步驟 9：非自駕交通最佳化

若非自駕，最佳化交通：
- 時刻表可靠性
- 總旅行時間
- 轉乘風險
- 營業時間
- 預約需求
- 票價
- 末班選項
- 識別哪些火車、巴士、渡輪、航班、pass、座位或 timed tickets 需要預約及截止時間

### 步驟 10：自駕成本估算

若自駕，估算交通成本：
- 距離-based 燃料或充電費用
- 過路費
- 停車費
- 渡輪/橋樑費
- 租車費用（如適用）
- 不確定性範圍
- 稀疏燃料區域警告
- 建議的加油/充電站

### 步驟 11：住宿推薦

若行程需要過夜：
- 推薦減少隔天移動的住宿基地或候選
- 避免路線折返
- 符合使用者的預算/風格
- 說明每個基地為什麼適合路線

### 步驟 12：檢查路線形狀

檢查不必要的折返：
- 若行程從 A 到 B 然後回到 A 與 B 之間的地點
- 重新排序當天行程、移動過夜基地，或明確標示折返並解釋為什麼仍然值得

### 步驟 13：檢查行程可行性

檢查每日行程是否太滿：
- 若太滿，移動、縮短或降級停留點
- 解釋取捨

### 步驟 14：加入備案

為以下情況加入 fallback plans：
- Rain（雨）
- Wind（強風）
- Heat（高溫）
- Snow（積雪）
- Closures（關閉）
- Full parking lots（停車場滿）
- Sold-out meals（售罄）
- Transport disruption（交通中斷）

### 步驟 15：國家/地區特定檢查

日本自駕行程：
- 包含 Mapcode（如可用）
- 優先選擇有一般訪客停車的目的地或停留點
- 使用最近的可確認訪客停車或官方指定停車作為 Google Maps 駕駛點
- 不可將月極停車場、居民專用、員工停車場或不明確的私人停車位視為可用停車

### 步驟 16：車中泊 / Road Trip 檢查

驗證以下項目：
- 過夜許可
- 安靜度
- 廁所
- 營業時間
- 噪音規定
- 洗澡選項
- 洗衣選項
- 垃圾規定
- 冬季道路狀況
- 附近備用住宿

### 步驟 17：日曆/App 就緒輸出

當有用時，加入：
- Stable event titles
- Start/end times
- Time zone
- Practical location or parking pin
- Notes
- Reminders
- Reservation references
- Map-list grouping
- Offline-map needs
- 哪些項目在驗證前不應加入

### 步驟 18：提供完整行程

最終產出應包含：
- 來源與信心標籤
- 假設
- 替代方案
- 需要預約的項目
- 地點信心
- 行程風險備註
- 路線形狀警告
- 住宿理由
- 推薦活動/美食
- 長距離交通比較（如相關）
- 燃料/充電計劃
- 旅遊社/model-course benchmark 備註
- 日曆/App 就緒欄位
- 成本估算與假設

## 與其他層的關係

- `workflow/travel-planning/execution-flow.md` 提供執行流程，本文件提供快速入門的操作步驟。
- `workflow/travel-planning/artifact-gates.md` 提供產出規範與品質門檻。
- `skills/travel-planning/SKILL.md` 是原始來源，已不再作為 active entrypoint（舊 `skills/` 結構已於 2026-05-13 標記為 deprecated）。
