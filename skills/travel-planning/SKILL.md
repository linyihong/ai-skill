---
name: travel-planning
description: Plan evidence-based travel itineraries with current operating hours, seasonal closures, transport timing, lodging or car-stay options, weather-driven route choices, backup plans, bathing and laundry stops, country-specific driving checks such as Japan Mapcode and visitor parking, local risk checks, and source-backed recommendations. Use when the user asks where to travel, what to do on specific dates, whether places are open, how to arrange a trip, road trip, Japan travel, Mapcode, parking, car camping, overnight parking, RV Park, michi-no-eki, bathhouse, laundromat, or 車中泊 planning.
---

# Travel Planning

Use this skill when a user wants a practical trip plan that depends on dates, place availability, transport feasibility, local conditions, or current web information. The goal is to produce an itinerary that is enjoyable, realistic, and traceable to sources.

**Shared policy:** read [`shared-rules` index](../../shared-rules/README.md), apply [`dependency-reading.md`](../../shared-rules/dependency-reading.md) when this skill or related rules change, apply [`neutral-language.md`](../../shared-rules/neutral-language.md) for titles and summaries, and apply [`goal-action-validation.md`](../../shared-rules/goal-action-validation.md) so major recommendations include sources, validation, and uncertainty. Lessons in `feedback_history/` should reference those files, not duplicate shared rules.

## When To Use

- Planning trips from destination, date range, budget, travel style, or constraints.
- Checking whether attractions, restaurants, parking areas, campgrounds, RV Parks, hot springs, events, ferries, or transit are open on planned dates.
- Designing road trip, public transit, walking, cycling, family, food, nature, photography, or 車中泊 itineraries.
- Comparing alternative areas, route orders, overnight bases, backup plans, rainy-day options, and weather-dependent timing.
- Planning road-trip support stops such as hot springs, public baths, showers, laundromats, fuel, charging, toilets, and groceries.
- Applying country- or region-specific travel requirements, such as Mapcode and visitor parking checks for Japan self-drive trips.
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
4. Build a feasible route with travel buffers and weather-aware ordering. Move outdoor, scenic, ferry, mountain, and walking-heavy plans into the best weather windows when possible.
5. Add fallback plans for rain, wind, heat, snow, closures, full parking lots, sold-out meals, and transport disruption.
6. Apply country/region-specific checks. For Japan self-drive plans, include Mapcode where available and prefer destinations or stops with ordinary visitor parking; do not treat 月極 parking, resident-only parking, staff parking, or unclear private lots as usable parking.
7. For 車中泊 or road trips, verify overnight permission, toilets, opening hours, noise rules, bathing options, laundry options, trash rules, winter road conditions, and nearby backup lodging.
8. Provide an itinerary with sources, confidence labels, assumptions, alternatives, and what still needs reservation or day-before confirmation.

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
- Weather, season, crowd, road, transit, and overnight-stay risks when relevant, including why the recommended order fits the forecast.
- Country/region-specific navigation and driving notes when relevant, including Japan Mapcode, visitor parking status, and parking caveats.
- Road-trip support points when relevant: bathing, shower, laundry, fuel, charging, toilet, grocery, and backup lodging options.
- Practical next actions: reservations, ticket purchase, route save, day-before checks, and fallback choices.

## Feedback Loop

If a reusable planning method emerges, create a lesson under [`feedback_history/`](feedback_history/) using [`shared-rules/feedback-lessons.md`](../../shared-rules/feedback-lessons.md). Keep specific traveler details, reservation codes, private addresses, and one-off live results out of reusable docs.
