# Travel Planning Workflow

Use this workflow to turn a user's travel idea into a realistic, source-backed itinerary.

## 1. Intake

Capture the minimum frame before planning:

- Destination or candidate area.
- Dates, arrival/departure times, and flexibility.
- Party size, age constraints, accessibility needs, dietary needs, luggage, pets, or driving limits.
- Travel mode: public transit, rental car, private car, walking, cycling, ferry, flight, or mixed.
- Style: slow travel, food, nature, hiking, hot springs, city, photography, shopping, family, budget, luxury, 車中泊, camping, or event-focused.
- Must-do and must-avoid items.
- Budget and reservation tolerance.
- Cost preference: fastest, cheapest, fewest transfers, scenic route, luggage-friendly, child-friendly, or low walking load.
- Weather tolerance: rain, wind, heat, cold, snow, low visibility, sea conditions, or mountain conditions.
- Country/region-specific needs: navigation format, toll rules, driving side, parking constraints, permits, or local payment methods.

If key details are missing, make a clearly labeled draft assumption and ask only blocker questions that affect feasibility.

## 2. Source Triage

Classify every important recommendation:

| Claim Type | Required Source |
| --- | --- |
| Opening hours, closing days, last entry | Official facility page, official SNS, booking page, or tourism board page. |
| Exact place identity | Google Maps place link or coordinate pin, official facility page, official address, map service, or facility access page. |
| Transit schedule, ferry, bus, train | Operator timetable or official route planner. |
| Transport fare and pass value | Operator fare table, official reservation page, pass page, fare calculator, or booking platform. |
| Required transport booking | Operator reservation page, seat availability page, ferry/flight/bus booking page, rail pass seat rule, or timed-ticket page. |
| Road conditions, winter closure, tolls | Road authority, highway operator, local government, or official map. |
| Driving cost | Route distance, fuel/energy price source, toll calculator, parking operator, ferry/bridge operator, rental-car contract, or charging network. |
| Lodging, minshuku, guesthouse, campground, RV Park | Official lodging page, booking platform, tourism board, map listing, campground/RV Park listing, or direct facility page. |
| Route shape and backtracking | Map/route planner, transport timetable, drive route, day-by-day stop order, and next-day base logic. |
| Event dates and crowd risk | Official event page, venue page, tourism board, or local government. |
| Weather-sensitive activity | Weather agency, mountain/weather service, facility notice, or operator notice. |
| 車中泊 permission | Facility official page, RV Park listing, 道の駅 page, local notice, or recent rule notice. |
| Bathing, shower, laundry, fuel, charging | Facility official page, map listing, operator page, review recency, or local service page. |
| Country-specific navigation and parking | Official tourism/facility page, parking operator, map service, Mapcode lookup, local road authority, rental-car guidance, or facility access page. |
| Discovery idea | Maps, community map, blog, video, review site, or user-provided source; verify before treating as confirmed. |

Prefer current official sources. If only community information exists, label it `needs confirmation` and provide a safer backup.

## 3. Exact Location Verification

Before route optimization, make sure each recommended place points to the intended location.

1. Prefer Google Maps links that open a single exact place, or a coordinate/map pin for the exact entrance, parking lot, bath, laundromat, station, pier, trailhead, or overnight base.
2. Avoid generic Google search or map search URLs when they return multiple branches, similarly named facilities, broad areas, or uncertain pins.
3. Cross-check the Google Maps pin with official facility name, official address, official access page, and map listing details.
4. For Japan self-drive, cross-check Mapcode against Google Maps and official address/access details. If the Mapcode points to parking while Google Maps points to the attraction entrance, state that difference explicitly.
5. If the location is a large area, choose the practical target pin: visitor parking, ticket office, trailhead, ferry terminal, station exit, campground reception, RV Park entrance, bath entrance, laundromat, or viewpoint parking.
6. Mark ambiguity as `location needs confirmation` when multiple pins remain plausible, the official page lacks address detail, Mapcode and Google Maps disagree, or reviews suggest the pin is wrong.

Do not hide location ambiguity inside the itinerary. If the user could navigate to the wrong place, put the concern next to the stop and add a safer fallback pin or verification action.

## 4. Weather and Backup Pass

Use weather as a planning input, not an afterthought:

1. Check forecast by area and time of day, not only the destination city.
2. Identify weather-sensitive items: viewpoints, hiking, cycling, ferries, ropeways, beaches, mountain roads, night driving, outdoor markets, and photo spots.
3. Put outdoor or scenic activities in the best available weather window.
4. Move indoor, food, shopping, museum, hot spring, laundry, and driving-transfer blocks into rain, heat, or low-visibility windows.
5. Prepare backups close to the planned route so a bad-weather switch does not create a new long detour.
6. Mark weather confidence: `stable`, `changeable`, `needs day-before check`, or `unsafe / should avoid`.

When weather could affect safety, transport, or road access, do not merely add a note. Reorder the itinerary, downgrade the activity, or add a concrete alternative.

## 5. Transport Mode Decision

Before route details, decide whether the plan is non-driving, self-drive, or mixed.

For non-driving plans:

1. Optimize by the user's priority: fastest, cheapest, fewest transfers, luggage-friendly, scenery, late-start, or low walking.
2. Build each travel leg with departure time, arrival time, transfer station/stop, buffer, platform/terminal risk when known, and last-return option.
3. Identify required bookings: limited express/reserved seats, highway bus, ferry, flight, airport transfer, event shuttle, ropeway/cable car slot, timed attraction ticket, taxi, or luggage service.
4. Record when to book: now, after lodging is fixed, same day, day-before, or no reservation needed.
5. Estimate fare per person and total group fare. Compare passes only when the break-even route is clear; do not recommend a pass just because it is popular.
6. Flag fragile legs: last bus, sparse rural route, short transfer, weather-cancelable ferry, sold-out seat risk, cash-only bus, or luggage-heavy transfer.

For self-drive plans:

1. Estimate rough driving cost before claiming driving is better than public transport.
2. Include distance, fuel economy or EV efficiency, fuel/charging unit price, tolls, paid parking, ferry/bridge fees, snow-chain or winter tire needs, and rental-car add-ons when relevant.
3. Provide a range when prices depend on route, vehicle, discount pass, ETC, day, season, or parking duration.
4. Compare time/cost tradeoffs against non-driving options when both are plausible.

For mixed plans, separate each leg and cost source so the user can see which parts require a car and which can be booked as transit.

## 6. Overnight Base and Lodging Planning

When the trip spans overnight stays, choose the overnight base as part of the route logic, not as an afterthought.

1. Decide the best overnight area from the route shape: near the last stop, near the next morning's first stop, near a transport hub, near reliable parking, or near bathing/laundry/food support.
2. Provide lodging options when useful: hotel, guesthouse, minshuku, ryokan, campground, RV Park, cabin, hostel, or business hotel. Match them to budget, parking/transit access, check-in time, luggage, family needs, and next-day route.
3. Prefer lodging that reduces next-day detours and avoids returning to an already-passed area.
4. Include basic booking checks: availability window, check-in deadline, parking or station access, cancellation risk, bath/laundry availability, breakfast timing, and late-arrival rules.
5. If the best lodging area is not the cheapest or most famous area, explain the route reason.
6. If no lodging candidate is verified yet, recommend a lodging area/base and mark specific lodging as `needs availability check`.

## 7. Route Shape and Backtracking Check

Before finalizing each day, inspect the route shape:

1. Prefer one-way progression, clean loops, or hub-and-spoke days with short local returns.
2. Watch for avoidable backtracking: going from A to B, then returning to a point between A and B, or passing a stop and visiting it later without a reason.
3. If backtracking appears, first try to reorder stops, move the overnight base, or split the day differently.
4. If the backtrack remains, label it `backtracking warning` and explain why: strong recommendation, fixed opening hours, weather window, transport timetable, sold-out lodging, sunset timing, or better parking/road access.
5. For strongly recommended points that create detours, state the tradeoff in time/cost and give a shorter alternative.
6. Do not hide route inefficiency inside a polished schedule; the user should be able to follow the flow without surprise returns.

## 8. Country and Region Specific Checks

Apply local driving and navigation rules before finalizing route order.

For Japan self-drive plans:

1. For each drive-to destination, search for a Mapcode when available. If unavailable, provide a phone number, address, official facility name, or map link as fallback navigation input.
2. Cross-check Mapcode, Google Maps exact place/pin, official address, and parking/access page. If they point to different entrances or lots, describe the difference and choose the best navigation target for the route.
3. Prefer attractions, restaurants, baths, laundromats, viewpoints, and overnight bases with ordinary visitor parking, public parking, facility parking, 道の駅 parking, service-area parking, or clearly listed paid coin parking.
4. Do not count 月極 parking, resident-only lots, staff-only lots, apartment parking, permit-only lots, or unclear private lots as usable visitor parking.
5. If parking is unclear, either downgrade the stop, add a nearby confirmed visitor parking option, or mark it `parking needs confirmation`.
6. Include parking caveats for narrow streets, height limits, winter closure, event restrictions, shuttle-only access, and popular lots that fill early.
7. When a stop is chosen mainly because it has reliable parking, say so in the itinerary.

For other countries or regions, identify equivalent local requirements before planning: navigation identifiers, low-emission zones, toll systems, vignette/permit needs, parking restrictions, road permits, ferry reservations, or seasonal access rules.

## 9. Feasibility Build

Plan from constraints outward:

1. Anchor immovable items: flights, booked lodging, event times, reservation slots, sunset/sunrise, ferry or bus times.
2. Add must-do stops with opening windows and last-entry times.
3. Calculate travel time using the user's transport mode and add buffers.
4. Place overnight base candidates where they reduce same-day and next-day friction.
5. Place meals, fuel, toilet, bathing, shower, laundry, charging, groceries, and rest stops where they naturally fit.
6. Add backups near the same route instead of far detours.
7. Remove or downgrade items that depend on perfect timing, unverified access, or inefficient backtracking.

Use conservative buffers:

- Urban transit or walking transfers: at least 10-20 minutes for unfamiliar places.
- Rural driving: add time for parking, narrow roads, fuel, snow, or mountain routes.
- Popular restaurants and attractions: account for queues, reservation checks, or sold-out risk.
- 車中泊: arrive before dark when rules, toilets, or parking layout are uncertain.

## 10. 車中泊 / Road Trip Checks

For each overnight candidate, verify:

- Overnight stay is allowed or at least not explicitly prohibited by current official rules.
- Parking access hours, gates, height limits, fees, quiet hours, and vehicle size limits.
- Toilet availability overnight.
- Bathing or shower option nearby and its closing time.
- Laundromat or washing option when the trip spans multiple nights, includes hiking/beach/snow activities, or the user asks for light packing.
- Trash, cooking, generator, idling, tent, table, and chair rules.
- Safety, lighting, winter closure, flood, wind, or isolation risk.
- Nearby legal backup: RV Park, campground, hotel, capsule hotel, ferry terminal, or 24-hour rest facility.

Do not describe a place as safe or permitted for 車中泊 unless the source supports it. Use `candidate` or `needs confirmation` when relying on community maps.

For comfort planning, group support stops into the route:

- Bathing: onsen, sento, super sento, day-use hotel bath, campground shower, or coin shower.
- Laundry: coin laundry near the overnight base, bath stop, grocery stop, or morning departure route.
- Recovery: late food, convenience store, supermarket, fuel/charging, toilet, trash rule, and dry indoor rest option.
- Timing: choose support stops with enough buffer before closing; avoid plans that require late-night bathing or laundry unless hours are confirmed.

## 11. Recommendation Pass

Before finalizing, check:

- Does each day have a clear theme and realistic pace?
- Does each recommended place have an exact Google Maps place link or precise pin, and are any ambiguous locations clearly marked?
- If the trip includes overnight stays, is each lodging area/candidate placed for route logic, check-in timing, parking/transit access, and next-day flow?
- Does the route avoid A→B→middle-point backtracking? If not, is there a `backtracking warning`, reason, and shorter alternative?
- Is the order optimized for forecast, daylight, and weather-sensitive activities?
- Are opening hours and travel times compatible with the user's dates?
- Are closures, seasonal access, holidays, and maintenance notices checked?
- If not driving, are departure times, transfer buffers, last-return options, required reservations, and fare estimates listed?
- If driving, is rough transport cost estimated with distance, fuel/charging, tolls, parking, ferry/bridge, and rental assumptions?
- For Japan self-drive, does each drive-to stop include Mapcode or fallback navigation input, has it been cross-checked against Google Maps and official address/access details, and is parking ordinary visitor parking rather than 月極 or private-only parking?
- Are food and rest breaks realistic?
- Are rainy-day, heat, snow, wind, transport, and closure backups close enough to use?
- For road trips or 車中泊, are bathing, laundry, fuel/charging, groceries, and backup lodging covered?
- Are all uncertain claims labeled?
- Are next actions specific: reserve, call, check official notice, buy pass, download offline map, or prepare cash?

## 12. Final Verification

For every important conclusion, provide:

- Goal: what the plan is optimizing for.
- Action: sources checked and route logic used.
- Validation or reference source: official/current source, or an explicit uncertainty label.

If current web access is unavailable, say so and provide a planning skeleton plus a checklist of sources the user should verify before travel.
