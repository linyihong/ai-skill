---
name: travel-planning
description: Plan evidence-based travel itineraries with current operating hours, exact Google Maps place links, driving navigation pins that target the nearest usable visitor parking instead of monthly/private lots, calendar-ready schedule blocks, reminder and app-ready itinerary fields, what to do at each stop, local food and restaurant recommendations filtered through country-appropriate review tools such as Google Maps and Japan Tabelog, travel agency itinerary benchmarking or direct package-tour adoption, package prices and included/excluded costs, long-distance transport comparison across flights, Shinkansen, limited express trains, highway buses, ferries, rental cars, self-drive and mixed modes, seasonal closures, transport timing, schedule feasibility checks, booking needs, fare estimates, lodging or car-stay options, lodging/guesthouse recommendations, anti-backtracking route checks, weather-driven route choices, backup plans, bathing and laundry stops, fuel/charging gap planning, car-stay quietness checks, country-specific driving checks such as Japan Mapcode and visitor parking, driving cost estimates, local risk checks, and source-backed recommendations. Use when the user asks where to travel, what to do on specific dates, whether places are open, how to arrange a trip, add to calendar, calendar-ready itinerary, itinerary app, travel app, map list, offline map, reminders, travel agency tour, package tour, tour price, model course, Google Maps, map pin, parking pin, food, restaurant, local specialty, Tabelog, 食べログ, local restaurant rating, fuel, gas station, charging, flight, airplane, Shinkansen, limited express, highway bus, ferry, public transport, train, bus, ticket booking, fare, hotel, guesthouse, minshuku, lodging, route optimization, road trip, Japan travel, Mapcode, parking, quiet overnight, car camping, overnight parking, RV Park, michi-no-eki, bathhouse, laundromat, or 車中泊 planning.
---

# Travel Planning

Use this skill when a user wants a practical trip plan that depends on dates, place availability, transport feasibility, local conditions, or current web information. The goal is to produce an itinerary that is enjoyable, realistic, and traceable to sources.

**Shared policy:** read [`shared-rules` index](../../shared-rules/README.md), apply [`dependency-reading.md`](../../shared-rules/dependency-reading.md) when this skill or related rules change, apply [`neutral-language.md`](../../shared-rules/neutral-language.md) for titles and summaries, and apply [`goal-action-validation.md`](../../shared-rules/goal-action-validation.md) so major recommendations include sources, validation, and uncertainty. If this skill is reloaded after an update, create a dependency read ledger covering required files, files read, missing files marked `not applicable`, blocked items, and validation. Lessons in `feedback_history/` should reference those files, not duplicate shared rules.

## When To Use

- Planning trips from destination, date range, budget, travel style, or constraints.
- Benchmarking or directly recommending travel agency tours, package tours, official model courses, and local tourism itineraries when they fit the user's needs.
- Checking whether attractions, restaurants, parking areas, campgrounds, RV Parks, hot springs, events, ferries, or transit are open on planned dates.
- Recommending what to do at each stop, expected visit duration, nearby highlights, local food, restaurants, cafes, markets, or regional specialties, filtered with country-appropriate restaurant rating and review tools.
- Marking exact locations with Google Maps place links, coordinates, or precise pins rather than ambiguous search-result links; for driving, the practical pin should be the nearest confirmed visitor-usable parking lot or official designated parking, not an unusable monthly/private lot.
- Designing road trip, public transit, walking, cycling, family, food, nature, photography, or 車中泊 itineraries.
- Comparing long-distance transport options such as flights, Shinkansen, limited express trains, highway buses, ferries, rental cars, self-drive, and mixed modes by door-to-door time and total cost.
- Preparing calendar-ready and app-ready itinerary outputs: event titles, start/end times, time zones, locations, map links, notes, reminders, reservation references, and route/map list grouping.
- Optimizing non-driving trips with trains, buses, ferries, flights, taxis, walking, passes, reservations, departure times, transfer buffers, and fare estimates.
- Recommending overnight bases, hotels, guesthouses, minshuku, campgrounds, RV Parks, or car-stay bases that fit the route and next-day plan.
- Checking route order for avoidable backtracking, especially A to B and then returning to a point between A and B.
- Comparing alternative areas, route orders, overnight bases, backup plans, rainy-day options, and weather-dependent timing.
- Planning road-trip support stops such as hot springs, public baths, showers, laundromats, fuel, charging, toilets, and groceries.
- Planning fuel or EV charging stops when rural, mountain, island, night, winter, or long-distance routes have sparse supply.
- Evaluating 車中泊 quietness and sleep quality risks such as truck traffic, idling, road noise, late-night crowds, lighting, nearby facilities, and early-morning activity.
- Checking schedule feasibility: realistic stop duration, travel buffers, meal timing, check-in, last entry, sunset, fatigue, and whether the day is too packed.
- Applying country- or region-specific travel requirements, such as Mapcode and visitor parking checks for Japan self-drive trips.
- Estimating rough driving costs such as fuel, tolls, parking, ferry/bridge fees, rental-car add-ons, and charging fees.
- Turning a user-provided travel site, map, blog, or official page into a sourced itinerary.

## Out Of Scope

- Guaranteeing availability without checking official or reservation sources.
- Legal, medical, immigration, insurance, or safety-critical advice beyond pointing to official sources.
- Storing user identity, exact home address, passport details, reservation codes, or payment details in reusable skill docs.
- Relying on a single blog, map pin, AI summary, or outdated page when official or recent sources are available.

## Quick Start（Routing）

See [`workflow/travel-planning/execution-flow.md`](../../workflow/travel-planning/execution-flow.md) for the full execution flow.

Routing summary:
1. Trip frame → 2. Time-sensitive checks → 3. Current web sources → 4. Travel agency tours → 5. Exact location → 6. Stop recommendations → 7. Long-distance transport comparison → 8. Feasible route → 9. Non-driving optimization → 10. Driving cost estimate → 11. Lodging → 12. Route shape check → 13. Schedule feasibility → 14. Fallback plans → 15. Country-specific checks → 16. 車中泊 checks → 17. Calendar/app-ready output → 18. Full itinerary.

## Default Workflow

### 新分層路徑（優先讀取）

| 用途 | 路徑 |
|------|------|
| 執行流程（Intake → Source Triage → Feasibility → Output） | [`workflow/travel-planning/execution-flow.md`](../../workflow/travel-planning/execution-flow.md) |
| 分析方法（Sources & Tools、分析方法說明） | [`analysis/travel/README.md`](../../analysis/travel/README.md) |
| 工程智慧（Heuristics） | [`intelligence/travel/README.md`](../../intelligence/travel/README.md) |
| 產出格式與品質門檻 | [`workflow/travel-planning/artifact-gates.md`](../../workflow/travel-planning/artifact-gates.md) |

### 舊路徑（保留向後相容）

| 用途 | 路徑 |
|------|------|
| WORKFLOW.md（舊執行流程） | [`WORKFLOW.md`](WORKFLOW.md) |
| TOOLS.md（舊工具參考） | [`TOOLS.md`](TOOLS.md) |
| DOCUMENTATION.md（舊產出格式） | [`DOCUMENTATION.md`](DOCUMENTATION.md) |
| README.md（舊說明） | [`README.md`](README.md) |

## Output Style & Artifact Gates

See [`workflow/travel-planning/artifact-gates.md`](../../workflow/travel-planning/artifact-gates.md) for output format and quality gates.

## Feedback Loop

See [`shared-rules/feedback-lessons.md`](../../shared-rules/feedback-lessons.md) for the feedback lesson template and workflow. See [`feedback/`](../../feedback/) for the feedback promotion pipeline.
