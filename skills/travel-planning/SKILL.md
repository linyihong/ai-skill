---
name: travel-planning
description: Plan evidence-based travel itineraries with current operating hours, exact Google Maps place links, seasonal closures, transport timing, booking needs, fare estimates, lodging or car-stay options, lodging/guesthouse recommendations, anti-backtracking route checks, weather-driven route choices, backup plans, bathing and laundry stops, country-specific driving checks such as Japan Mapcode and visitor parking, driving cost estimates, local risk checks, and source-backed recommendations. Use when the user asks where to travel, what to do on specific dates, whether places are open, how to arrange a trip, Google Maps, map pin, public transport, train, bus, ferry, ticket booking, fare, hotel, guesthouse, minshuku, lodging, route optimization, road trip, Japan travel, Mapcode, parking, car camping, overnight parking, RV Park, michi-no-eki, bathhouse, laundromat, or 車中泊 planning.
---

# Travel Planning

Use this skill when a user wants a practical trip plan that depends on dates, place availability, transport feasibility, local conditions, or current web information. The goal is to produce an itinerary that is enjoyable, realistic, and traceable to sources.

**Shared policy:** read [`shared-rules` index](../../shared-rules/README.md), apply [`dependency-reading.md`](../../shared-rules/dependency-reading.md) when this skill or related rules change, apply [`neutral-language.md`](../../shared-rules/neutral-language.md) for titles and summaries, and apply [`goal-action-validation.md`](../../shared-rules/goal-action-validation.md) so major recommendations include sources, validation, and uncertainty. Lessons in `feedback_history/` should reference those files, not duplicate shared rules.

## When To Use

- Planning trips from destination, date range, budget, travel style, or constraints.
- Checking whether attractions, restaurants, parking areas, campgrounds, RV Parks, hot springs, events, ferries, or transit are open on planned dates.
- Marking exact locations with Google Maps place links, coordinates, or precise pins rather than ambiguous search-result links.
- Designing road trip, public transit, walking, cycling, family, food, nature, photography, or 車中泊 itineraries.
- Optimizing non-driving trips with trains, buses, ferries, flights, taxis, walking, passes, reservations, departure times, transfer buffers, and fare estimates.
- Recommending overnight bases, hotels, guesthouses, minshuku, campgrounds, RV Parks, or car-stay bases that fit the route and next-day plan.
- Checking route order for avoidable backtracking, especially A to B and then returning to a point between A and B.
- Comparing alternative areas, route orders, overnight bases, backup plans, rainy-day options, and weather-dependent timing.
- Planning road-trip support stops such as hot springs, public baths, showers, laundromats, fuel, charging, toilets, and groceries.
- Applying country- or region-specific travel requirements, such as Mapcode and visitor parking checks for Japan self-drive trips.
- Estimating rough driving costs such as fuel, tolls, parking, ferry/bridge fees, rental-car add-ons, and charging fees.
- Turning a user-provided travel site, map, blog, or official page into a sourced itinerary.

## Out Of Scope

- Guaranteeing availability without checking official or reservation sources.
- Legal, medical, immigration, insurance, or safety-critical advice beyond pointing to official sources.
- Storing user identity, exact home address, passport details, reservation codes, or payment details in reusable skill docs.
- Relying on a single blog, map pin, AI summary, or outdated page when official or recent sources are available.

## Quick Start

1. Clarify the trip frame: destination, dates, party size, transport, pace, budget, must-do interests, dietary/accessibility needs, lodging style, and tolerance for long drives or early starts.
2. Identify time-sensitive checks: operating days, reservation windows, seasonal closures, event crowding, weather forecast, road/weather constraints, parking rules, public transport schedules, and last-entry times.
3. Use current web sources. Prefer official facility, tourism board, transit operator, weather, road authority, and reservation pages; use community maps or blogs for discovery, then verify details elsewhere.
4. Verify exact location identity. Prefer a Google Maps place link or coordinate pin that opens one exact place; avoid generic search URLs that return many possible points. Cross-check the map pin against official name/address and, for Japan self-drive, Mapcode when available.
5. Build a feasible route with travel buffers and weather-aware ordering. Move outdoor, scenic, ferry, mountain, and walking-heavy plans into the best weather windows when possible.
6. If the user is not driving, optimize transport by schedule reliability, total travel time, transfer risk, operating hours, reservation needs, fare, and last-return options. Identify which trains, buses, ferries, flights, passes, seats, or timed tickets need booking and by when.
7. If the user is driving, estimate rough transport costs: distance-based fuel or charging, tolls, parking, ferry/bridge fees, rental-car fees when relevant, and uncertainty ranges.
8. If the trip requires overnight stays, recommend lodging bases or accommodation candidates that reduce next-day travel, avoid route backtracking, and match the user's budget/style. Include why each base fits the route.
9. Check route shape for unnecessary backtracking. If a plan goes from A to B and then returns to a point between A and B, either reorder the day, move the overnight base, or explicitly flag the backtrack and explain why it is still worth it.
10. Add fallback plans for rain, wind, heat, snow, closures, full parking lots, sold-out meals, and transport disruption.
11. Apply country/region-specific checks. For Japan self-drive plans, include Mapcode where available and prefer destinations or stops with ordinary visitor parking; do not treat 月極 parking, resident-only parking, staff parking, or unclear private lots as usable parking.
12. For 車中泊 or road trips, verify overnight permission, toilets, opening hours, noise rules, bathing options, laundry options, trash rules, winter road conditions, and nearby backup lodging.
13. Provide an itinerary with sources, confidence labels, assumptions, alternatives, what needs reservation, location confidence, route-shape warnings, lodging rationale, and cost estimates with assumptions.

## Default Workflow

Read [WORKFLOW.md](WORKFLOW.md) for the planning decision tree.

Use [TOOLS.md](TOOLS.md) when choosing current information sources for official hours, travel conditions, transport, weather, events, and 車中泊 discovery.

Use [DOCUMENTATION.md](DOCUMENTATION.md) when writing itinerary outputs, source tables, open-question lists, and day-before checklists.

Use [README.md](README.md) for human guidance and skill boundaries.

## Output Style

When producing a plan, include:

- Trip assumptions: dates, area, transport, party, pace, and constraints.
- A day-by-day itinerary with time blocks, travel time, opening hours, last entry, reservation status, and backup options.
- Source-backed validation for time-sensitive claims.
- Confidence labels: `confirmed`, `likely`, `needs day-before check`, or `unknown`.
- Exact location notes when relevant: Google Maps place link or coordinate pin, official name/address match, Mapcode cross-check, and any ambiguity.
- Weather, season, crowd, road, transit, and overnight-stay risks when relevant, including why the recommended order fits the forecast.
- Transport plan when relevant: route, departure/arrival windows, transfers, booking deadlines, required reservations, pass/ticket options, last-return risk, and fare estimate.
- Lodging recommendations when overnighting: area/base logic, hotel/guesthouse/minshuku/RV Park candidates, access, parking or transit fit, check-in timing, and why the base avoids unnecessary backtracking.
- Route-shape notes: whether the day is mostly one-way/loop/backtracking, any A→B→middle-point return, and whether the detour is avoidable or strongly recommended.
- Country/region-specific navigation and driving notes when relevant, including Japan Mapcode, visitor parking status, and parking caveats.
- Driving cost estimate when relevant: assumed distance, fuel/energy unit cost, fuel economy or efficiency, tolls, parking, ferry/bridge fees, rental-car add-ons, and confidence range.
- Road-trip support points when relevant: bathing, shower, laundry, fuel, charging, toilet, grocery, and backup lodging options.
- Practical next actions: reservations, ticket purchase, route save, day-before checks, and fallback choices.

## Feedback Loop

If a reusable planning method emerges, create a lesson under [`feedback_history/`](feedback_history/) using [`shared-rules/feedback-lessons.md`](../../shared-rules/feedback-lessons.md). Keep specific traveler details, reservation codes, private addresses, and one-off live results out of reusable docs.
