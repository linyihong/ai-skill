# Travel Planning Skill

This skill supports source-backed travel planning: itinerary design, travel agency itinerary benchmarking, stop-level activity and food recommendations, exact location verification, current opening-hour checks, schedule feasibility checks, weather-aware route ordering, seasonal feasibility, public-transport optimization, fare and driving-cost estimates, lodging or guesthouse recommendations, car-stay routing, car-stay quietness checks, anti-backtracking route checks, fuel/charging gap planning, country-specific driving checks, road-trip support stops, and backup planning.

## Goals

- Turn a destination and date range into a practical route.
- Use travel agency tours, package tours, and official model courses as benchmarks for route order, stop duration, seasonal highlights, meal timing, and booking constraints.
- Explain what to do at each recommended stop, how long to stay, and what local food or restaurants to consider.
- Keep daily schedules realistic with stop-duration, travel-buffer, meal, check-in, last-entry, sunset, and fatigue checks.
- Mark places with exact Google Maps place links, coordinate pins, or precise map URLs instead of ambiguous search-result links.
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

## What Belongs Here

- Reusable workflows for planning trips.
- Source hierarchy, travel agency/model-course benchmark checks, stop-level activity/food checks, exact-location checks, schedule-feasibility checks, lodging/base checks, car-stay quietness checks, anti-backtracking route checks, transport booking/cost checks, fuel/charging checks, country-specific checks, and verification rules.
- Output templates for itineraries, agency/model-course comparison, stop recommendations, schedule feasibility, lodging candidates, route-shape warnings, quietness notes, transport plans, cost estimates, weather/backup logic, support-stop tables, source tables, and day-before checklists.
- Reusable lessons about travel planning quality, not private trip details.

## What Does Not Belong Here

- Passport, payment, reservation code, home address, or traveler identity details.
- One-off live availability results that only apply to a specific user's trip.
- Claims that a facility is open, bookable, or legal for overnight stay without source and timestamp context.
- Legal, medical, immigration, or insurance advice beyond linking official sources.

## Files

| File | Purpose |
| --- | --- |
| `SKILL.md` | Cursor/agent entry point and trigger rules. |
| `WORKFLOW.md` | Planning and verification workflow. |
| `TOOLS.md` | Source categories and preferred lookup strategy. |
| `DOCUMENTATION.md` | Itinerary, source table, risk, and checklist templates. |
| `FEEDBACK.md` | Short entry pointing to shared feedback rules. |
| `feedback_history/` | Reusable travel-planning lessons. |

## Use Pattern

1. Start from the user's destination, dates, style, and constraints.
2. Gather official and current sources for every time-sensitive recommendation.
3. Check travel agency tours, package tours, and official model courses when useful, then verify and adapt rather than copying.
4. Verify exact place identity with Google Maps and official name/address; cross-check Mapcode when relevant.
5. Add what to do, expected stay time, and food/local-specialty ideas for key stops.
6. Check weather and local disruption risks before locking route order.
7. If not driving, optimize transport routes and list booking/ticket requirements with fare estimates.
8. If driving, estimate rough transport costs, apply country-specific navigation/access/parking checks, and plan fuel/charging stops.
9. If overnighting, recommend lodging bases or candidates that reduce next-day friction and fit the route.
10. For car-stay or overnight parking, evaluate quietness and sleep-quality risk.
11. Plan route order with buffers, backups, support stops, anti-backtracking checks, and time feasibility checks.
12. Mark confidence and unresolved checks clearly.
13. Convert new reusable planning lessons into `feedback_history/` when they generalize.
