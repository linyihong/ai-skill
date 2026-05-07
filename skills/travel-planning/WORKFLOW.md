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
| Travel agency tour / package tour / model course | Travel agency itinerary page, official tourism model course, operator tour page, brochure page, booking itinerary, price page, cancellation terms, or included/excluded items page. |
| Exact place identity | Google Maps place link or coordinate pin, official facility page, official address, map service, or facility access page. |
| What to do, local food, restaurants | Official tourism page, facility page, local tourism board, restaurant page, market page, recent map listing, or local guide. |
| Transit schedule, ferry, bus, train | Operator timetable or official route planner. |
| Long-distance transport comparison | Airline page, airport access page, Shinkansen/rail operator, limited express route, highway bus operator, ferry operator, rental-car company, toll/fuel/parking calculator, or route planner. |
| Transport fare and pass value | Operator fare table, official reservation page, pass page, fare calculator, or booking platform. |
| Required transport booking | Operator reservation page, seat availability page, ferry/flight/bus booking page, rail pass seat rule, or timed-ticket page. |
| Road conditions, winter closure, tolls | Road authority, highway operator, local government, or official map. |
| Driving cost | Route distance, fuel/energy price source, toll calculator, parking operator, ferry/bridge operator, rental-car contract, or charging network. |
| Lodging, minshuku, guesthouse, campground, RV Park | Official lodging page, booking platform, tourism board, map listing, campground/RV Park listing, or direct facility page. |
| Route shape and backtracking | Map/route planner, transport timetable, drive route, day-by-day stop order, and next-day base logic. |
| Schedule feasibility | Opening hours, last entry, visit duration, transfer/drive time, check-in deadline, meal hours, sunset/sunrise, and fatigue risk. |
| 車中泊 quietness | Recent reviews, map context, road/truck traffic, idling risk, nearby facilities, lighting, late-night activity, and official quiet-hour rules. |
| Event dates and crowd risk | Official event page, venue page, tourism board, or local government. |
| Weather-sensitive activity | Weather agency, mountain/weather service, facility notice, or operator notice. |
| 車中泊 permission | Facility official page, RV Park listing, 道の駅 page, local notice, or recent rule notice. |
| Bathing, shower, laundry, fuel, charging | Facility official page, map listing, operator page, review recency, route distance, opening hours, or local service page. |
| Country-specific navigation and parking | Official tourism/facility page, parking operator, map service, Mapcode lookup, local road authority, rental-car guidance, or facility access page. |
| Discovery idea | Maps, community map, blog, video, review site, or user-provided source; verify before treating as confirmed. |

Prefer current official sources. If only community information exists, label it `needs confirmation` and provide a safer backup.

## 3. Travel Agency and Model-Course Benchmark / Direct Option

Use travel agency tours and official model courses as planning references, or as direct recommended options when they match the region, season, transport mode, budget, or trip length.

1. Search agency tours, package tours, bus tours, self-drive model courses, and tourism-board model routes for the target area.
2. If recommending a package directly, show the package price, price basis (per person / group / vehicle), what is included, what is excluded, cancellation/change conditions, meeting point, departure/return time, and booking deadline.
3. Warn the user that a package tour may assume charter buses, fixed group movement, group meals, shopping stops, reduced free time, and agency-controlled timing.
4. Compare the package against self-planning when possible: likely cheaper/easier/safer, what flexibility is lost, and which costs remain outside the package.
5. If using the agency route as a benchmark only, extract useful structure: route order, common base areas, typical stop duration, lunch/meal placement, seasonal highlights, transport mode, booking/ticket requirements, and how long operators allocate between stops.
6. Do not treat agency content as verified fact. Re-check all opening hours, access, map pins, parking, weather suitability, ticket rules, and current availability through official/current sources.
7. Do not copy a tour wholesale unless the user chooses the package as the plan. Adapt benchmarked routes to the user's pace, transport mode, budget, lodging style, car-stay needs, weather, and backtracking constraints.
8. If agency routes skip a place the user wants, check whether the skip is due to time, access, parking, season, or route inefficiency.
9. Record whether the agency item is a `direct package option` or `benchmark only`, and what was changed.

## 4. Exact Location Verification

Before route optimization, make sure each recommended place points to the intended location.

1. Prefer Google Maps links that open a single exact place, or a coordinate/map pin for the exact entrance, parking lot, bath, laundromat, station, pier, trailhead, or overnight base.
2. Avoid generic Google search or map search URLs when they return multiple branches, similarly named facilities, broad areas, or uncertain pins.
3. Cross-check the Google Maps pin with official facility name, official address, official access page, and map listing details.
4. For Japan self-drive, cross-check Mapcode against Google Maps and official address/access details. If the Mapcode points to parking while Google Maps points to the attraction entrance, state that difference explicitly.
5. If the location is a large area, choose the practical target pin: visitor parking, ticket office, trailhead, ferry terminal, station exit, campground reception, RV Park entrance, bath entrance, laundromat, or viewpoint parking.
6. Mark ambiguity as `location needs confirmation` when multiple pins remain plausible, the official page lacks address detail, Mapcode and Google Maps disagree, or reviews suggest the pin is wrong.

Do not hide location ambiguity inside the itinerary. If the user could navigate to the wrong place, put the concern next to the stop and add a safer fallback pin or verification action.

## 5. Stop Experience and Food Planning

For each recommended stop, give the user enough context to know what to do there.

1. State the reason to stop: scenery, food, hot spring, museum, walk, market, viewpoint, local culture, seasonal event, rest, or route convenience.
2. List 1-3 concrete things to do, not only the place name.
3. Estimate time needed: quick photo stop, 30-60 minutes, half day, meal stop, or overnight base.
4. Add food ideas when relevant: local specialty, market, restaurant/cafe candidate, roadside station food, convenience fallback, opening/last-order risk, and reservation need.
5. Align food stops with route timing; avoid recommending lunch after a restaurant's last order or dinner after rural closing time.
6. If a recommended stop has no clear activity value, mark it as a support stop rather than a sightseeing stop.

## 6. Weather and Backup Pass

Use weather as a planning input, not an afterthought:

1. Check forecast by area and time of day, not only the destination city.
2. Identify weather-sensitive items: viewpoints, hiking, cycling, ferries, ropeways, beaches, mountain roads, night driving, outdoor markets, and photo spots.
3. Put outdoor or scenic activities in the best available weather window.
4. Move indoor, food, shopping, museum, hot spring, laundry, and driving-transfer blocks into rain, heat, or low-visibility windows.
5. Prepare backups close to the planned route so a bad-weather switch does not create a new long detour.
6. Mark weather confidence: `stable`, `changeable`, `needs day-before check`, or `unsafe / should avoid`.

When weather could affect safety, transport, or road access, do not merely add a note. Reorder the itinerary, downgrade the activity, or add a concrete alternative.

## 7. Long-Distance Transport Comparison Gate

Apply this gate for cross-city, inter-prefecture, island, airport, or any transfer expected to take 2+ hours.

1. Compare all plausible modes: flight, Shinkansen, limited express/train, highway bus, ferry, rental car, full self-drive, and mixed options such as `rail + local rental car`.
2. Use door-to-door time, not just in-vehicle time. Include home/hotel to station/airport, check-in/security, transfer wait, luggage pickup, local transit, parking, and rental-car pickup/return.
3. Use total cost, not just ticket fare. Include seat/reservation fees, airport access, baggage fees, local transit, taxi, rental car, tolls, fuel/charging, parking, ferry vehicle fare, and cancellation/change fees when relevant.
4. Compare practical burden: luggage, walking, transfers, missed-connection risk, weather/delay risk, late-night arrival, car fatigue, and whether the destination still needs a car.
5. Mark each option as `recommended`, `viable`, `backup`, or `not recommended`, with the reason.
6. If the recommended option is not cheapest or fastest, explain the tradeoff.

## 8. Transport Mode Decision

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
3. Check fuel or EV charging gaps. If the route crosses rural, mountain, island, winter, night, or long-distance areas, identify the last reliable fuel/charging point and the next reliable point.
4. Recommend where to fuel/charge before entering sparse areas, including opening hours, payment caveats, and backup options.
5. Provide a range when prices depend on route, vehicle, discount pass, ETC, day, season, or parking duration.
6. Compare time/cost tradeoffs against non-driving options when both are plausible.

For mixed plans, separate each leg and cost source so the user can see which parts require a car and which can be booked as transit.

## 9. Overnight Base and Lodging Planning

When the trip spans overnight stays, choose the overnight base as part of the route logic, not as an afterthought.

1. Decide the best overnight area from the route shape: near the last stop, near the next morning's first stop, near a transport hub, near reliable parking, or near bathing/laundry/food support.
2. Provide lodging options when useful: hotel, guesthouse, minshuku, ryokan, campground, RV Park, cabin, hostel, or business hotel. Match them to budget, parking/transit access, check-in time, luggage, family needs, and next-day route.
3. Prefer lodging that reduces next-day detours and avoids returning to an already-passed area.
4. Include basic booking checks: availability window, check-in deadline, parking or station access, cancellation risk, bath/laundry availability, breakfast timing, and late-arrival rules.
5. If the best lodging area is not the cheapest or most famous area, explain the route reason.
6. If no lodging candidate is verified yet, recommend a lodging area/base and mark specific lodging as `needs availability check`.

## 10. Route Shape and Backtracking Check

Before finalizing each day, inspect the route shape:

1. Prefer one-way progression, clean loops, or hub-and-spoke days with short local returns.
2. Watch for avoidable backtracking: going from A to B, then returning to a point between A and B, or passing a stop and visiting it later without a reason.
3. If backtracking appears, first try to reorder stops, move the overnight base, or split the day differently.
4. If the backtrack remains, label it `backtracking warning` and explain why: strong recommendation, fixed opening hours, weather window, transport timetable, sold-out lodging, sunset timing, or better parking/road access.
5. For strongly recommended points that create detours, state the tradeoff in time/cost and give a shorter alternative.
6. Do not hide route inefficiency inside a polished schedule; the user should be able to follow the flow without surprise returns.

## 11. Country and Region Specific Checks

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

## 12. Feasibility Build

Plan from constraints outward:

1. Anchor immovable items: flights, booked lodging, event times, reservation slots, sunset/sunrise, ferry or bus times.
2. Add must-do stops with opening windows and last-entry times.
3. Calculate travel time using the user's transport mode and add buffers.
4. Place overnight base candidates where they reduce same-day and next-day friction.
5. Place meals, fuel, toilet, bathing, shower, laundry, charging, groceries, and rest stops where they naturally fit; do not leave long rural legs without fuel/charging planning.
6. Add backups near the same route instead of far detours.
7. Remove or downgrade items that depend on perfect timing, unverified access, or inefficient backtracking.

Use conservative buffers:

- Urban transit or walking transfers: at least 10-20 minutes for unfamiliar places.
- Rural driving: add time for parking, narrow roads, fuel, snow, or mountain routes.
- Popular restaurants and attractions: account for queues, reservation checks, or sold-out risk.
- 車中泊: arrive before dark when rules, toilets, or parking layout are uncertain.

## 13. Schedule Feasibility Check

Before finalizing time blocks, make sure the day can be followed in practice.

1. For each stop, estimate minimum and comfortable visit duration.
2. Add movement buffers: parking search, walking from parking/station, ticketing, bathroom, luggage, fuel/charging, and rural road delays.
3. Check hard time limits: last entry, last order, last bus/train/ferry, check-in deadline, sunset, winter road closure, and bath/laundry closing time.
4. Flag fatigue risk when driving hours, walking load, early starts, late arrivals, or consecutive long days stack up.
5. If the day is too packed, downgrade lower-priority stops, move food/support stops, split the day, or suggest an overnight base change.
6. Label schedule pressure as `comfortable`, `tight`, `too packed`, or `needs day-before check`.

## 14. 車中泊 / Road Trip Checks

For each overnight candidate, verify:

- Overnight stay is allowed or at least not explicitly prohibited by current official rules.
- Parking access hours, gates, height limits, fees, quiet hours, and vehicle size limits.
- Quietness and sleep quality: distance from major roads, truck parking/idling, late-night traffic, nearby convenience store/restaurant activity, lighting, train/port noise, event/crowd risk, and early-morning activity.
- Toilet availability overnight.
- Bathing or shower option nearby and its closing time.
- Laundromat or washing option when the trip spans multiple nights, includes hiking/beach/snow activities, or the user asks for light packing.
- Trash, cooking, generator, idling, tent, table, and chair rules.
- Safety, lighting, winter closure, flood, wind, or isolation risk.
- Nearby legal backup: RV Park, campground, hotel, capsule hotel, ferry terminal, or 24-hour rest facility.
- Fuel or charging before the overnight area if late-night or early-morning supply is uncertain.

Do not describe a place as safe or permitted for 車中泊 unless the source supports it. Use `candidate` or `needs confirmation` when relying on community maps.

For quietness, use practical labels:

- `quiet`: low traffic/activity, no obvious truck/idling concentration, suitable for sleep.
- `moderate`: some traffic, lighting, or facilities nearby; usable with caveats.
- `noisy`: road/truck/late-night activity likely affects sleep; prefer only if necessary.
- `unknown`: insufficient recent evidence; add a quieter backup.

For comfort planning, group support stops into the route:

- Bathing: onsen, sento, super sento, day-use hotel bath, campground shower, or coin shower.
- Laundry: coin laundry near the overnight base, bath stop, grocery stop, or morning departure route.
- Recovery: late food, convenience store, supermarket, fuel/charging, toilet, trash rule, and dry indoor rest option.
- Timing: choose support stops with enough buffer before closing; avoid plans that require late-night bathing or laundry unless hours are confirmed.

## 15. Recommendation Pass

Before finalizing, check:

- Does each day have a clear theme and realistic pace?
- If agency/model-course references were used, did the plan state whether each source is a direct package option or benchmark only, price/inclusions when direct, and what was borrowed or changed after verification?
- For 2+ hour or cross-region transfers, did the plan compare plausible long-distance transport modes by door-to-door time, total cost, luggage/transfer burden, and delay/weather risk?
- Is each day schedule `comfortable`, `tight`, or `too packed`, and were overloaded days simplified?
- For each key stop, does the plan say what to do, how long to stay, and whether there is a food/local-specialty recommendation or support-stop role?
- Does each recommended place have an exact Google Maps place link or precise pin, and are any ambiguous locations clearly marked?
- If the trip includes overnight stays, is each lodging area/candidate placed for route logic, check-in timing, parking/transit access, and next-day flow?
- Does the route avoid A→B→middle-point backtracking? If not, is there a `backtracking warning`, reason, and shorter alternative?
- Is the order optimized for forecast, daylight, and weather-sensitive activities?
- Are opening hours and travel times compatible with the user's dates?
- Are closures, seasonal access, holidays, and maintenance notices checked?
- If not driving, are departure times, transfer buffers, last-return options, required reservations, and fare estimates listed?
- If driving, is rough transport cost estimated with distance, fuel/charging, tolls, parking, ferry/bridge, and rental assumptions?
- If driving through sparse areas, are fuel/charging gaps identified with recommended refuel/charge points and backups?
- For Japan self-drive, does each drive-to stop include Mapcode or fallback navigation input, has it been cross-checked against Google Maps and official address/access details, and is parking ordinary visitor parking rather than 月極 or private-only parking?
- Are food and rest breaks realistic?
- Are rainy-day, heat, snow, wind, transport, and closure backups close enough to use?
- For road trips or 車中泊, are bathing, laundry, fuel/charging, groceries, and backup lodging covered?
- For 車中泊, is quietness/sleep quality labeled and is there a quieter backup if the main candidate is noisy or unknown?
- Are all uncertain claims labeled?
- Are next actions specific: reserve, call, check official notice, buy pass, download offline map, or prepare cash?

## 16. Final Verification

For every important conclusion, provide:

- Goal: what the plan is optimizing for.
- Action: sources checked and route logic used.
- Validation or reference source: official/current source, or an explicit uncertainty label.

If current web access is unavailable, say so and provide a planning skeleton plus a checklist of sources the user should verify before travel.
