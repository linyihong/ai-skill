# Travel Analysis Methods（旅遊規劃分析方法）

本文件定義旅遊規劃的核心目標、範圍與使用模式。承接 [`skills/travel-planning/README.md`](../../skills/travel-planning/README.md) 的內容，提取為 tool-neutral 的分析方法。

> **相容性規則**：`skills/travel-planning/README.md` 仍為 active skill entrypoint。本文件為 reference target，兩者應保持同步。

## 核心目標

- Turn a destination and date range into a practical route.
- Use travel agency tours, package tours, and official model courses as benchmarks, or recommend them directly when they fit the user's needs.
- Show package price, included/excluded items, booking/cancellation conditions, and self-planned cost comparison when presenting agency options.
- Compare long-distance transport options by door-to-door time and total cost, including flights, Shinkansen, limited express, highway bus, ferry, rental car, self-drive, and mixed modes.
- Explain what to do at each recommended stop, how long to stay, and what local food or restaurants to consider.
- Filter restaurant candidates with the destination country's commonly used review/rating tools, such as Google Maps plus 食べログ for Japan when practical.
- Keep daily schedules realistic with stop-duration, travel-buffer, meal, check-in, last-entry, sunset, and fatigue checks.
- Mark places with exact Google Maps place links, coordinate pins, or precise map URLs instead of ambiguous search-result links.
- For self-drive routes, use the nearest confirmed visitor-usable parking lot or official designated parking as the Google Maps navigation point when it differs from the attraction, restaurant, trailhead, or facility entrance.
- Provide calendar-ready and app-ready fields when useful: event titles, start/end times, time zones, map links, reminders, reservation notes, map-list grouping, and offline-map needs.
- Verify time-sensitive details before recommending a stop.
- Separate confirmed facts from assumptions and open questions.
- Use weather forecasts to choose better route order and realistic backup plans.
- Optimize non-driving routes with transport schedules, transfer buffers, booking needs, and fare estimates.
- Recommend overnight bases and lodging candidates that fit the route, next-day plan, budget, and travel style.
- Avoid inefficient route shapes, especially going from A to B and then returning to an intermediate point unless clearly justified.
- Evaluate car-stay quietness and sleep quality, including traffic, truck idling, late-night crowds, lighting, and nearby facilities.
- Apply country- and region-specific requirements, such as Mapcode and visitor parking checks for Japan self-drive routes.
- Estimate rough self-drive costs from distance, fuel/charging, tolls, parking, ferries, bridges, and rental-car add-ons.
- Plan fuel or EV charging stops when the route crosses rural, mountain, island, night, winter, or long-distance areas with sparse supply.
- Use community sources for discovery while grounding decisions in official or current sources.
- Make car-stay and road-trip plans realistic: legal overnight status, toilets, bathing, laundry, trash rules, noise rules, weather, road conditions, and backup lodging.

## 範圍

### 屬於本層

- Reusable workflows for planning trips.
- Source hierarchy, travel agency/model-course benchmark or direct-package checks, package price comparison, long-distance transport comparison, stop-level activity/food checks, country-specific restaurant rating/review checks, exact-location and driving-parking-pin checks, calendar/app-ready output checks, schedule-feasibility checks, lodging/base checks, car-stay quietness checks, anti-backtracking route checks, transport booking/cost checks, fuel/charging checks, country-specific checks, and verification rules.
- Output templates for itineraries, agency/model-course comparison, package price comparison, long-distance transport comparison, stop recommendations, schedule feasibility, lodging candidates, route-shape warnings, quietness notes, transport plans, cost estimates, weather/backup logic, support-stop tables, source tables, and day-before checklists.
- Reusable lessons about travel planning quality, not private trip details.

### 不屬於本層

- Passport, payment, reservation code, home address, or traveler identity details.
- One-off live availability results that only apply to a specific user's trip.
- Claims that a facility is open, bookable, or legal for overnight stay without source and timestamp context.
- Legal, medical, immigration, or insurance advice beyond linking official sources.

## 使用模式

1. Start from the user's destination, dates, style, and constraints.
2. Gather official and current sources for every time-sensitive recommendation.
3. Check travel agency tours, package tours, and official model courses when useful. Either benchmark them or list them as direct options with price, inclusions/exclusions, booking conditions, and user-facing caveats.
4. Verify exact place identity with Google Maps and official name/address; for driving, choose the nearest confirmed visitor parking/official parking pin as the practical navigation target and cross-check Mapcode when relevant.
5. Add what to do, expected stay time, and food/local-specialty ideas for key stops. Screen restaurant candidates with local review/rating platforms and route/timing constraints.
6. For cross-city, inter-prefecture, island, airport, or 2+ hour transfers, compare long-distance transport options by time, money, baggage burden, and booking risk.
7. Check weather and local disruption risks before locking route order.
8. If not driving, optimize transport routes and list booking/ticket requirements with fare estimates.
9. If driving, estimate rough transport costs, apply country-specific navigation/access/parking checks, and plan fuel/charging stops.
10. If overnighting, recommend lodging bases or candidates that reduce next-day friction and fit the route.
11. For car-stay or overnight parking, evaluate quietness and sleep-quality risk.
12. Plan route order with buffers, backups, support stops, anti-backtracking checks, and time feasibility checks.
13. Add calendar/app-ready fields when the user may import the trip into a calendar, map list, reminder app, notes app, travel planning app, or offline map.
14. Mark confidence and unresolved checks clearly.
15. Convert new reusable planning lessons into `feedback_history/` when they generalize.

## 與其他層的關係

- `workflow/travel-planning/execution-flow.md` 提供執行流程，本文件提供流程中的分析目標與範圍。
- `analysis/travel/sources-and-tools.md` 提供來源選擇策略與工具知識。
- `workflow/travel-planning/artifact-gates.md` 提供產出格式規範。
- `intelligence/travel/` 提供旅遊規劃的啟發式規則。
- `skills/travel-planning/README.md` 是原始來源，仍為 active entrypoint。
