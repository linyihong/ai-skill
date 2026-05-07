# Travel Planning Skill

This skill supports source-backed travel planning: itinerary design, travel agency itinerary benchmarking or direct package-tour adoption, package price comparison, long-distance transport comparison, stop-level activity and food recommendations, country-appropriate restaurant rating/review screening, exact location verification, current opening-hour checks, schedule feasibility checks, weather-aware route ordering, seasonal feasibility, public-transport optimization, fare and driving-cost estimates, lodging or guesthouse recommendations, car-stay routing, car-stay quietness checks, anti-backtracking route checks, fuel/charging gap planning, country-specific driving checks, road-trip support stops, and backup planning.

## Goals

- Turn a destination and date range into a practical route.
- Use travel agency tours, package tours, and official model courses as benchmarks, or recommend them directly when they fit the user's needs.
- Show package price, included/excluded items, booking/cancellation conditions, and self-planned cost comparison when presenting agency options.
- Compare long-distance transport options by door-to-door time and total cost, including flights, Shinkansen, limited express, highway bus, ferry, rental car, self-drive, and mixed modes.
- Explain what to do at each recommended stop, how long to stay, and what local food or restaurants to consider.
- Filter restaurant candidates with the destination country's commonly used review/rating tools, such as Google Maps plus 食べログ for Japan when practical.
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
- Source hierarchy, travel agency/model-course benchmark or direct-package checks, package price comparison, long-distance transport comparison, stop-level activity/food checks, country-specific restaurant rating/review checks, exact-location checks, schedule-feasibility checks, lodging/base checks, car-stay quietness checks, anti-backtracking route checks, transport booking/cost checks, fuel/charging checks, country-specific checks, and verification rules.
- Output templates for itineraries, agency/model-course comparison, package price comparison, long-distance transport comparison, stop recommendations, schedule feasibility, lodging candidates, route-shape warnings, quietness notes, transport plans, cost estimates, weather/backup logic, support-stop tables, source tables, and day-before checklists.
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
3. Check travel agency tours, package tours, and official model courses when useful. Either benchmark them or list them as direct options with price, inclusions/exclusions, booking conditions, and user-facing caveats.
4. Verify exact place identity with Google Maps and official name/address; cross-check Mapcode when relevant.
5. Add what to do, expected stay time, and food/local-specialty ideas for key stops. Screen restaurant candidates with local review/rating platforms and route/timing constraints.
6. For cross-city, inter-prefecture, island, airport, or 2+ hour transfers, compare long-distance transport options by time, money, baggage burden, and booking risk.
7. Check weather and local disruption risks before locking route order.
8. If not driving, optimize transport routes and list booking/ticket requirements with fare estimates.
9. If driving, estimate rough transport costs, apply country-specific navigation/access/parking checks, and plan fuel/charging stops.
10. If overnighting, recommend lodging bases or candidates that reduce next-day friction and fit the route.
11. For car-stay or overnight parking, evaluate quietness and sleep-quality risk.
12. Plan route order with buffers, backups, support stops, anti-backtracking checks, and time feasibility checks.
13. Mark confidence and unresolved checks clearly.
14. Convert new reusable planning lessons into `feedback_history/` when they generalize.
