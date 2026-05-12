# Travel Planning Workflow

`workflow/travel-planning/` 負責「旅行規劃的執行流程」。本目錄保存 agent 在規劃旅行時可照著執行的 intake、source triage、route optimization、feasibility check 與 output formatting 流程。

## Scope

本 workflow 涵蓋以下規劃類型：

- **Itinerary Planning**：從使用者需求到完整行程的端到端流程。
- **Transportation Research**：長距離交通比較、非自駕/自駕/混合模式決策。
- **Budget Planning**：費用估算（交通、住宿、餐飲、停車、燃料）。
- **Accommodation Planning**：住宿區域選擇與 lodging candidate 推薦。
- **Route Optimization**：路線形狀檢查、避免 backtracking、天氣驅動的順序調整。
- **車中泊 / Road Trip Planning**：過夜許可、安靜度評估、支援站點規劃。

## 核心原則

1. **Source-backed planning**：每個重要推薦必須有官方或當前來源支撐。
2. **Exact location first**：每個推薦地點必須有精確的 Google Maps place link 或 coordinate pin。
3. **Feasibility over ambition**：行程必須在實際時間限制內可執行。
4. **Weather-aware ordering**：戶外活動應放在最佳天氣窗口。
5. **Transparent uncertainty**：所有不確定的聲明必須標註 confidence label。

## 與既有層的關係

- `skills/travel-planning/` 目前仍是 active skill entrypoint；本層只承接逐步抽出的通用執行流程。
- `skills/travel-planning/WORKFLOW.md` 是目前的 workflow source of truth。
- `skills/travel-planning/DOCUMENTATION.md` 是目前的 artifact template source of truth。
- `intelligence/travel/` 可被本 workflow 引用來輔助規劃判斷。

## 第一批候選遷移來源

- `skills/travel-planning/WORKFLOW.md` — ✅ 已提取（execution-flow.md）
- `skills/travel-planning/DOCUMENTATION.md` — ✅ 已提取（artifact-gates.md）

## 已提取內容

| 檔案 | 來源 | 說明 |
|------|------|------|
| [`execution-flow.md`](execution-flow.md) | `WORKFLOW.md` §1-17 | 完整 17 步驟執行流程：Intake、Source Triage、Agency Benchmark、Location Verification、Stop Planning、Weather、Transport、Lodging、Route Shape、Country Checks、Feasibility、Schedule、Calendar Output、車中泊、Recommendation Pass、Final Verification |
| [`artifact-gates.md`](artifact-gates.md) | `DOCUMENTATION.md` | 14 個產出模板與 final verification checklist：Itinerary Summary、Day Plan、Weather Strategy、Source Table、Calendar/App Table、Offline Checklist、Agency Benchmark Table、Stop Experience Table、Restaurant Table、Location Table、Transport Plan、Cost Estimate、車中泊 Quietness Table、Verification Checklist |

## 建議 Workflow 流程

### Itinerary Planning Flow

```
1. Intake → capture trip frame (destination, dates, party, transport, pace).
2. Source triage → classify every important claim by required source type.
3. Agency/model-course benchmark → search and compare package tours.
4. Exact location verification → Google Maps place links, parking pins, Mapcode.
5. Stop experience and food planning → what to do, how long, what to eat.
6. Weather and backup pass → weather-aware ordering, concrete alternatives.
7. Long-distance transport comparison gate → for 2+ hour transfers.
8. Transport mode decision → non-driving / self-drive / mixed.
9. Overnight base and lodging planning → route-logic-driven base selection.
10. Route shape and backtracking check → avoid A→B→middle-point returns.
11. Country/region specific checks → Mapcode, visitor parking, local rules.
12. Feasibility build → anchor immovable items, add buffers, place support stops.
13. Schedule feasibility check → label each day comfortable/tight/too packed.
14. Calendar/app-ready output pass → structured fields for import.
15. 車中泊 / road trip checks → permission, quietness, support stops.
16. Recommendation pass → 30+ point checklist before finalizing.
17. Final verification → goal, action, validation for every conclusion.
```

### Transportation Research Flow

```
1. Identify all plausible modes for each 2+ hour transfer.
2. Compare door-to-door time and total cost.
3. Evaluate practical burden: luggage, transfers, delay risk.
4. Mark each option: recommended / viable / backup / not recommended.
5. For non-driving: build legs with departure/arrival/transfer/booking.
6. For self-drive: estimate cost, check fuel/charging gaps.
```

### Budget Planning Flow

```
1. Long-distance transport: compare fares across modes.
2. Daily transport: local transit, tolls, parking.
3. Lodging: per-night estimates for recommended areas.
4. Food: per-meal budget based on area and style.
5. Fuel/charging: distance × efficiency × unit price.
6. Activities: entry fees, tickets, reservations.
7. Provide range when prices depend on season, booking time, or route.
```

## 產出格式

每次規劃應產出：

- **Trip Frame**（≤100 tokens）：目的地、日期、人數、交通、步調。
- **Day-by-day Itinerary**（每 day ≤500 tokens）：時間區塊、地點、交通、驗證、備案。
- **Source Table**（≤300 tokens）：每個重要聲明的來源、檢查時間、信心標籤。
- **Calendar/App-ready Fields**（可選）：事件標題、時間、時區、位置、提醒。
- **Final Verification**（≤100 tokens）：確認所有檢查點通過。
