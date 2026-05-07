# Travel Planning Skill

This skill supports source-backed travel planning: itinerary design, stop-level activity and food recommendations, exact location verification, current opening-hour checks, weather-aware route ordering, seasonal feasibility, public-transport optimization, fare and driving-cost estimates, lodging or guesthouse recommendations, car-stay routing, anti-backtracking route checks, fuel/charging gap planning, country-specific driving checks, road-trip support stops, and backup planning.

## Goals

- Turn a destination and date range into a practical route.
- Explain what to do at each recommended stop, how long to stay, and what local food or restaurants to consider.
- Mark places with exact Google Maps place links, coordinate pins, or precise map URLs instead of ambiguous search-result links.
- Verify time-sensitive details before recommending a stop.
- Separate confirmed facts from assumptions and open questions.
- Use weather forecasts to choose better route order and realistic backup plans.
- Optimize non-driving routes with transport schedules, transfer buffers, booking needs, and fare estimates.
- Recommend overnight bases and lodging candidates that fit the route, next-day plan, budget, and travel style.
- Avoid inefficient route shapes, especially going from A to B and then returning to an intermediate point unless clearly justified.
- Apply country- and region-specific requirements, such as Mapcode and visitor parking checks for Japan self-drive routes.
- Estimate rough self-drive costs from distance, fuel/charging, tolls, parking, ferries, bridges, and rental-car add-ons.
- Plan fuel or EV charging stops when the route crosses rural, mountain, island, night, winter, or long-distance areas with sparse supply.
- Use community sources for discovery while grounding decisions in official or current sources.
- Make car-stay and road-trip plans realistic: legal overnight status, toilets, bathing, laundry, trash rules, noise rules, weather, road conditions, and backup lodging.

## What Belongs Here

- Reusable workflows for planning trips.
- Source hierarchy, stop-level activity/food checks, exact-location checks, lodging/base checks, anti-backtracking route checks, transport booking/cost checks, fuel/charging checks, country-specific checks, and verification rules.
- Output templates for itineraries, stop recommendations, lodging candidates, route-shape warnings, transport plans, cost estimates, weather/backup logic, support-stop tables, source tables, and day-before checklists.
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
3. Verify exact place identity with Google Maps and official name/address; cross-check Mapcode when relevant.
4. Add what to do, expected stay time, and food/local-specialty ideas for key stops.
5. Check weather and local disruption risks before locking route order.
6. If not driving, optimize transport routes and list booking/ticket requirements with fare estimates.
7. If driving, estimate rough transport costs, apply country-specific navigation/access/parking checks, and plan fuel/charging stops.
8. If overnighting, recommend lodging bases or candidates that reduce next-day friction and fit the route.
9. Plan route order with buffers, backups, support stops, and anti-backtracking checks.
10. Mark confidence and unresolved checks clearly.
11. Convert new reusable planning lessons into `feedback_history/` when they generalize.
