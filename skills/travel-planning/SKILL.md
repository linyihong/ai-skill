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

## Quick Start

1. Clarify the trip frame: destination, dates, party size, transport, pace, budget, must-do interests, dietary/accessibility needs, lodging style, and tolerance for long drives or early starts.
2. Identify time-sensitive checks: operating days, reservation windows, seasonal closures, event crowding, weather forecast, road/weather constraints, parking rules, public transport schedules, and last-entry times.
3. Use current web sources. Prefer official facility, tourism board, transit operator, weather, road authority, and reservation pages; use community maps or blogs for discovery, then verify details elsewhere.
4. Search travel agency tours, package tours, and official model courses for the same area/date/season when helpful. They may be used either as benchmarks or as direct recommended options, but the user must be told when a plan is based on an agency/package itinerary.
   - For direct package recommendations, show price, what is included/excluded, booking/cancellation notes, meeting/departure point, transport assumptions, and why it may be easier or safer than self-planning.
   - For benchmark use, extract route order, stop duration, seasonal highlights, meal/transport patterns, and hidden constraints, then verify each claim against official/current sources.
5. Verify exact location identity. Prefer a Google Maps place link or coordinate pin that opens one exact place; avoid generic search URLs that return many possible points. For driving routes, use the nearest confirmed visitor-usable parking lot, official parking lot, or practical arrival lot as the navigation target when it differs from the attraction/restaurant entrance. Cross-check the map pin against official name/address and, for Japan self-drive, Mapcode when available.
6. For each recommended stop, add what to do there, why it is worth stopping, expected stay time, nearby alternatives, and food or local-specialty options when relevant.
   - For restaurant recommendations, use local review/rating tools appropriate to the country plus Google Maps when available. For Japan, cross-check Google Maps with 食べログ when practical; consider rating, review count, recency, opening hours, last order, reservation needs, price range, queue risk, parking/access, route fit, and nearby backups.
7. For cross-city, inter-prefecture, island, airport, or 2+ hour transfers, compare long-distance transport options before choosing the main route. Include door-to-door time and total cost, not only ticket price.
8. Build a feasible route with travel buffers and weather-aware ordering. Move outdoor, scenic, ferry, mountain, and walking-heavy plans into the best weather windows when possible.
9. If the user is not driving, optimize transport by schedule reliability, total travel time, transfer risk, operating hours, reservation needs, fare, and last-return options. Identify which trains, buses, ferries, flights, passes, seats, or timed tickets need booking and by when.
10. If the user is driving, estimate rough transport costs and plan refuel/charging: distance-based fuel or charging, tolls, parking, ferry/bridge fees, rental-car fees when relevant, uncertainty ranges, sparse-fuel areas, and recommended fuel/charging stops.
11. If the trip requires overnight stays, recommend lodging bases or accommodation candidates that reduce next-day travel, avoid route backtracking, and match the user's budget/style. Include why each base fits the route.
12. Check route shape for unnecessary backtracking. If a plan goes from A to B and then returns to a point between A and B, either reorder the day, move the overnight base, or explicitly flag the backtrack and explain why it is still worth it.
13. Check schedule feasibility. If the day is too packed, move, shorten, or downgrade stops and explain the tradeoff.
14. Add fallback plans for rain, wind, heat, snow, closures, full parking lots, sold-out meals, and transport disruption.
15. Apply country/region-specific checks. For Japan self-drive plans, include Mapcode where available and prefer destinations or stops with ordinary visitor parking. Use the nearest confirmed visitor parking or official designated parking as the Google Maps driving point; do not treat 月極 parking, resident-only parking, staff parking, or unclear private lots as usable parking.
16. For 車中泊 or road trips, verify overnight permission, quietness, toilets, opening hours, noise rules, bathing options, laundry options, trash rules, winter road conditions, and nearby backup lodging.
17. Add calendar/app-ready output when useful: stable event titles, start/end times, time zone, practical location or parking pin, notes, reminders, reservation references, map-list grouping, offline-map needs, and what should not be added until verified.
18. Provide an itinerary with sources, confidence labels, assumptions, alternatives, what needs reservation, location confidence, schedule-risk notes, route-shape warnings, lodging rationale, recommended activities/food, long-distance transport comparison when relevant, fuel/charging plan, agency/model-course benchmark notes, calendar/app-ready fields, and cost estimates with assumptions.

## Default Workflow

Read [WORKFLOW.md](WORKFLOW.md) for the planning decision tree.

Use [TOOLS.md](TOOLS.md) when choosing current information sources for official hours, travel conditions, transport, weather, events, and 車中泊 discovery.

Use [DOCUMENTATION.md](DOCUMENTATION.md) when writing itinerary outputs, source tables, open-question lists, and day-before checklists.

Use [README.md](README.md) for human guidance and skill boundaries.

## Output Style

When producing a plan, include:

- Trip assumptions: dates, area, transport, party, pace, and constraints.
- A day-by-day itinerary with time blocks, travel time, opening hours, last entry, reservation status, and backup options.
- Travel agency / model-course notes when used: source, matching season/area, whether it is a direct package option or benchmark only, price, included/excluded items, what was borrowed, and what was changed after verification.
- Schedule feasibility notes: visit duration, movement buffers, meal timing, sunset/last-entry/check-in constraints, fatigue risk, and what was shortened or moved.
- Stop-level recommendations: what to do, why it is worth it, expected visit duration, local food/restaurant ideas, country-appropriate rating/review signals, and nearby alternatives.
- Source-backed validation for time-sensitive claims.
- Confidence labels: `confirmed`, `likely`, `needs day-before check`, or `unknown`.
- Exact location notes when relevant: Google Maps place link or coordinate pin, official name/address match, driving parking pin, Mapcode cross-check, and any ambiguity.
- Calendar/app-ready notes when relevant: event title, start/end time, timezone, location, notes, reminder timing, reservation reference, map-list grouping, and whether the item is safe to import now or needs recheck.
- Weather, season, crowd, road, transit, and overnight-stay risks when relevant, including why the recommended order fits the forecast.
- Transport plan when relevant: route, departure/arrival windows, transfers, booking deadlines, required reservations, pass/ticket options, last-return risk, and fare estimate.
- Long-distance transport comparison when relevant: flight, Shinkansen, limited express, highway bus, ferry, rental car, self-drive, and mixed-mode options with door-to-door time, total cost, luggage burden, booking/cancellation, and delay/weather risk.
- Lodging recommendations when overnighting: area/base logic, hotel/guesthouse/minshuku/RV Park candidates, access, parking or transit fit, check-in timing, and why the base avoids unnecessary backtracking.
- Route-shape notes: whether the day is mostly one-way/loop/backtracking, any A→B→middle-point return, and whether the detour is avoidable or strongly recommended.
- 車中泊 quietness notes when relevant: expected noise level, traffic/idling/crowd/light risk, sleep quality confidence, and quieter alternatives.
- Country/region-specific navigation and driving notes when relevant, including Japan Mapcode, visitor parking status, and parking caveats.
- Driving cost estimate when relevant: assumed distance, fuel/energy unit cost, fuel economy or efficiency, tolls, parking, ferry/bridge fees, rental-car add-ons, and confidence range.
- Road-trip support points when relevant: bathing, shower, laundry, fuel, charging, toilet, grocery, and backup lodging options, including warnings for sparse fuel/charging areas.
- Practical next actions: reservations, ticket purchase, route save, day-before checks, and fallback choices.

## Feedback Loop

If a reusable planning method emerges, create a lesson under [`feedback_history/`](feedback_history/) using [`shared-rules/feedback-lessons.md`](../../shared-rules/feedback-lessons.md). Keep specific traveler details, reservation codes, private addresses, and one-off live results out of reusable docs.
