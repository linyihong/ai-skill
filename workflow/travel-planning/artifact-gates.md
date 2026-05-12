# Travel Planning Artifact Gates（旅行規劃產出規範）

本文件定義旅行規劃的產出格式與品質門檻。承接 [`skills/travel-planning/SKILL.md`](../../skills/travel-planning/SKILL.md) 的 Output Style 內容，提取為 workflow 層的 artifact gates。

> **相容性規則**：`skills/travel-planning/SKILL.md` 仍為 active skill entrypoint。本文件為 reference target，兩者應保持同步。

## 1. 產出必備項目

當產出旅行計劃時，必須包含以下項目：

| # | 項目 | 說明 |
|---|------|------|
| 1 | **Trip assumptions** | 日期、區域、交通、人數、步調、限制條件 |
| 2 | **Day-by-day itinerary** | 時間區塊、移動時間、營業時間、最後入場、預約狀態、備案 |
| 3 | **Travel agency / model-course notes** | 來源、匹配季節/區域、是否為直接套裝或僅 benchmark、價格、包含/排除項目、借用了什麼、驗證後改了什麼 |
| 4 | **Schedule feasibility notes** | 停留時間、移動緩衝、用餐時間、日落/最後入場/入住限制、疲勞風險、縮短或移動了什麼 |
| 5 | **Stop-level recommendations** | 做什麼、為什麼值得、預期停留時間、當地美食/餐廳建議、當地評分/評論信號、附近替代方案 |
| 6 | **Source-backed validation** | 時間敏感資訊需有來源驗證 |
| 7 | **Confidence labels** | `confirmed`、`likely`、`needs day-before check`、`unknown` |
| 8 | **Exact location notes** | Google Maps place link 或 coordinate pin、官方名稱/地址比對、駕駛停車 pin、Mapcode 交叉檢查、任何模糊點 |
| 9 | **Calendar/app-ready notes** | 事件標題、開始/結束時間、時區、地點、備註、提醒時間、預約參考、地圖列表分組、是否可安全匯入或需再確認 |
| 10 | **Weather/season/crowd/road/transit/overnight risks** | 包含為什麼推薦順序符合天氣預報 |
| 11 | **Transport plan** | 路線、出發/到達窗口、轉乘、預訂截止、必要預約、pass/票券選項、末班風險、票價估算 |
| 12 | **Long-distance transport comparison** | 飛機、新幹線、特急、高速巴士、渡輪、租車、自駕、混合模式，含 door-to-door 時間、總費用、行李負擔、預約/取消、延誤/天氣風險 |
| 13 | **Lodging recommendations** | 區域/基地邏輯、飯店/民宿/RV Park 候選、交通、停車或大眾運輸配合、入住時間、為什麼基地避免不必要折返 |
| 14 | **Route-shape notes** | 當天是否多為單向/環狀/折返、是否有 A→B→中間點折返、繞路是否可避免或強烈建議 |
| 15 | **車中泊 quietness notes** | 預期噪音、交通/怠速/人群/燈光風險、睡眠品質信心、較安靜替代方案 |
| 16 | **Country/region-specific navigation notes** | 日本 Mapcode、訪客停車狀態、停車注意事項 |
| 17 | **Driving cost estimate** | 假設距離、燃料/能源單價、油耗或效率、過路費、停車費、渡輪/橋樑費、租車附加費、信心範圍 |
| 18 | **Road-trip support points** | 洗澡、淋浴、洗衣、加油、充電、廁所、超市、備用住宿，含稀疏燃料/充電區警告 |
| 19 | **Practical next actions** | 預約、購票、路線儲存、出發前檢查、備案選擇 |

## 2. 品質門檻

### 2.1 地點精確度

- 每個推薦地點必須有 **Google Maps place link** 或 **coordinate pin**，而非模糊的搜尋結果連結。
- 駕駛路線的導航目標必須是 **最近的可確認訪客停車場** 或 **官方指定停車場**，而非月極停車場、居民專用、員工停車場或不明確的私人停車位。
- 日本自駕行程必須包含 **Mapcode**（如可用）。

### 2.2 來源驗證

- 時間敏感資訊（營業時間、季節性關閉、預約窗口、天氣、道路狀況）必須有 **官方或近期來源**。
- 旅遊社套裝行程或官方 model course 可作為 benchmark 或直接推薦，但必須標示來源、匹配程度、價格、包含/排除項目。
- 單一部落格、地圖 pin、AI summary 或過期頁面不可作為唯一來源。

### 2.3 可行性檢查

- 每日行程必須有 **schedule feasibility check**：停留時間、移動緩衝、用餐時間、日落/最後入場/入住限制、疲勞風險。
- 路線形狀必須檢查 **不必要的折返**：若 A→B 後回到 A 與 B 之間的地點，必須重新排序或明確標示繞路原因。
- 長距離移動（跨城市、跨縣、島嶼、機場、2+ 小時）必須比較 **多種交通方式**（door-to-door 時間 + 總費用）。

### 2.4 備案要求

- 每個行程必須有 **fallback plans**：雨天、強風、高溫、積雪、關閉、停車場滿、售罄、交通中斷。
- 車中泊行程必須驗證：過夜許可、安靜度、廁所、營業時間、噪音規定、洗澡選項、洗衣選項、垃圾規定、冬季道路狀況、附近備用住宿。

## 3. 產出格式範例

### Trip Frame
```markdown
## Trip Frame
- Destination: [城市/區域]
- Dates: [YYYY-MM-DD] ~ [YYYY-MM-DD]
- Party: [人數]
- Transport: [自駕/大眾運輸/混合]
- Pace: [緊湊/悠閒/混合]
- Budget: [預算範圍]
- Must-do: [必去景點/活動]
- Dietary/accessibility: [特殊需求]
- Lodging style: [飯店/民宿/車中泊/混合]
```

### Day Block
```markdown
### Day <n> - <主題/區域>
- Date: YYYY-MM-DD (星期)
- Weather forecast: [天氣概況]

| Time | Activity | Location | Duration | Notes |
|------|----------|----------|----------|-------|
| 09:00 | [活動] | [place link] | 2h | [備註] |
| 11:00 | [移動] | → [下個地點] | 30min | [交通方式] |
| ... | ... | ... | ... | ... |

**Feasibility**: ✅/⚠️/❌ [說明]
**Backup**: [備案]
```

### Confidence Label
```markdown
| Item | Confidence | Source |
|------|-----------|--------|
| [景點] 營業時間 | confirmed | [官方網站連結] |
| [餐廳] 評價 | likely | Google Maps 4.2★ / 食べログ 3.5★ |
| [活動] 預約狀況 | needs day-before check | 當天電話確認 |
| [道路] 狀況 | unknown | 出發前再確認 |
```

## 4. 與其他層的關係

- `workflow/travel-planning/execution-flow.md` 提供執行流程，本文件定義產出規範。
- `skills/travel-planning/SKILL.md` 是原始來源，仍為 active entrypoint。
- `skills/travel-planning/DOCUMENTATION.md` 提供詳細的輸出模板。
