# Travel Planning Sources and Tools

Use current web sources for travel planning. Do not rely on model memory for operating hours, closures, prices, route disruptions, or overnight-stay rules.

## Source Hierarchy

| Priority | Source Type | Use For |
| --- | --- | --- |
| 1 | Official facility, transport operator, event, park, campground, tourism board, local government, weather, marine, mountain, road authority, official address, and facility access pages | Confirming facts that affect feasibility and identity. |
| 2 | Official reservation platforms, ticketing pages, RV Park listings, transit operators, rail/bus/ferry/airline companies, fare calculators, pass pages, toll calculators, bath/shower facilities, laundromats, parking operators, charging networks, and fuel station operators | Availability, booking rules, fees, cancellation rules, schedules, access, cost. |
| 3 | Google Maps exact place links, coordinate pins, current map services, Mapcode lookup services, and review platforms | Discovery, exact location identity, navigation identifiers, parking/access notes, recent user signals; verify key claims elsewhere. |
| 4 | Blogs, community maps, videos, forum posts, social posts | Discovery and qualitative context; label as unconfirmed unless cross-checked. |

## Japan Travel Sources

Use official or operator pages first:

- Facility official websites and official social accounts.
- Google Maps exact place links or coordinate pins; avoid generic search-result links when they can open many similarly named places.
- Local tourism association pages.
- Train, bus, ferry, ropeway, highway bus, airport operators, official route planners, fare tables, reservation pages, seat rules, and pass pages.
- Japan Meteorological Agency or reputable weather services; for islands, coasts, ferries, mountains, or snow routes, also check marine, wind, wave, snowfall, or road-condition sources.
- Highway and road authority pages for closures, snow, tolls, ETC/pass options, and rest areas.
- National park, prefecture, city, and 道の駅 official pages.
- RV Park, campground, and booking sites for overnight rules and reservations.
- Mapcode lookup services and facility access pages for self-drive navigation. If Mapcode is unavailable, capture fallback navigation input: official name, address, phone number, Google Maps exact place link, or coordinate pin.
- Onsen, sento, super sento, public bath, coin shower, campground shower, laundromat, fuel, EV charging, and grocery pages or map listings for support-stop planning.
- Parking operator pages, official facility access pages, and map listings for visitor parking. Prefer ordinary visitor parking, public parking, facility parking, coin parking, 道の駅 parking, or service-area parking. Exclude 月極 parking, resident-only parking, staff-only parking, apartment parking, permit-only parking, and unclear private lots.

For 車中泊 discovery, a user may provide community maps such as `https://syachuhaku.fxtec.info/`. Treat these as discovery sources, then verify overnight rules, toilets, gates, fees, and recent notices through official pages, RV Park listings, local pages, or direct contact when needed.

## Checks By Travel Type

| Travel Type | Minimum Checks |
| --- | --- |
| City trip | Opening days, last entry, restaurant reservations, transit pass value, crowd/event calendar. |
| Rural road trip | Drive time, fuel, parking, road closure, weather, limited food hours, backup lodging. |
| 車中泊 | Overnight permission, toilets, bath/shower, laundry, trash/cooking/idling rules, safety, backup lodging. |
| Hiking/nature | Trail status, weather, daylight, transport return, facility closure, gear and emergency limits. |
| Ferry/island | Sailing schedule, weather cancellation, vehicle reservation, last return, port access. |
| Seasonal events | Official event dates, ticketing, crowd control, special transit, road restrictions. |

## Transport Booking and Cost Sources

For non-driving plans, use official or operator-backed sources for:

- Timetable: train, bus, ferry, flight, shuttle, ropeway/cable car, airport transfer, or taxi reservation window.
- Required booking: reserved seat, highway bus, ferry vehicle/passenger booking, flight, airport transfer, event shuttle, timed attraction entry, luggage delivery, or taxi.
- Fare: base fare, seat fee, limited express fee, bus/ferry fare, pass price, local day pass, child fare, luggage/bicycle fee, and cancellation/change rules when important.
- Final-leg risk: last bus/train, sparse rural service, weather-cancelable ferry, cash-only bus, sold-out seat, and transfer buffer.

For driving plans, estimate with:

- Distance and time from a route planner.
- Fuel economy or EV efficiency assumptions.
- Fuel price, charging rate, toll calculator, ferry/bridge fee, parking fee, rental-car fee, insurance/add-ons, winter tire or chain fee when relevant.
- Range or caveat when ETC discounts, vehicle class, seasonal rates, parking duration, or detours can change the number.

Cost estimates should be labeled as rough unless all operator prices and route choices are fixed.

## Country-Specific Driving Sources

When the trip uses a car, check whether the destination country has local navigation or access conventions.

### Japan self-drive

For every drive-to stop, try to record:

- Mapcode for car navigation when available.
- Google Maps exact place link or coordinate pin.
- Fallback navigation input: official facility name, address, phone number, or map link.
- Parking type: visitor/facility/public/coin/RV Park/道の駅/service area, or `unclear`.
- Parking caveats: fee, hours, height limit, distance to destination, event restrictions, shuttle-only access, winter closure, or fills-early risk.
- Cross-check status: Mapcode, Google Maps, official address, and access/parking page match or have noted differences.

Parking source handling:

- Prefer official facility access pages and parking operator pages.
- Map listings and reviews can help discover parking, but verify whether the lot is visitor-usable.
- Treat 月極, resident-only, staff-only, apartment, permit-only, and unclear private lots as unavailable unless an official source says visitors can use them.
- If no reliable visitor parking exists, choose a different stop, suggest public transport/taxi for that segment, or add a confirmed nearby paid parking lot.

## Exact Location Sources

Use precise location sources when marking stops:

- Google Maps place URLs that open one exact facility, not a search results page.
- Coordinate pins when the practical destination is an entrance, trailhead, parking lot, pier, station exit, or viewpoint instead of the facility's general address.
- Official address and access pages to confirm the map pin.
- Mapcode for Japan self-drive, cross-checked against the Google Maps pin and official address/access page.

Location red flags:

- Search URL or broad area link instead of a specific place.
- Multiple facilities with the same or similar name.
- Google Maps pin, Mapcode, and official address point to different roads, entrances, lots, or nearby towns.
- Attraction pin is correct but the practical navigation target should be the visitor parking lot, shuttle stop, pier, trailhead, or reception.
- Map listing is user-generated and the official page gives a different address or access route.

## Weather-Aware Planning Sources

Check weather at the smallest practical area and time window:

- General forecast: temperature, precipitation probability, wind, humidity, heat index, and hourly trend.
- Mountain or highland: visibility, wind, snow, trail alerts, ropeway operation, and road closures.
- Coast or island: wind, wave, ferry cancellation risk, tide, storm, and marine warnings.
- Winter road trip: snow, ice, chain/tire requirements, pass closures, and daylight.
- Urban trip: heavy rain, heat, crowd/event congestion, indoor alternatives, and transit disruption.

Use weather to decide route order. If the forecast is uncertain, provide Plan A / Plan B rather than a single fragile schedule.

## Road-Trip Support Stops

For car-stay or long driving days, search and verify:

| Support Need | Examples | Checks |
| --- | --- | --- |
| Bathing / shower | Onsen, sento, super sento, day-use bath, coin shower, campground shower | Hours, last entry, closed days, towel rental, parking, tattoo policy when relevant, crowd risk. |
| Laundry | Coin laundry, campground laundry, hotel day-use laundry when available | Hours, parking, machine count if known, detergent, drying time, nearby food/bath stop. |
| Food / groceries | Supermarket, convenience store, local market, late restaurant | Closing time, parking, rest-day risk, cash/card. |
| Fuel / charging | Fuel station, EV charger, roadside station, service area | Hours, charger type, payment app/card, rural gaps. |
| Rest / toilet | 道の駅, service area, park, public toilet, bath facility | Overnight access, cleanliness recency, gates, lighting. |

## Source Note Format

When citing a time-sensitive source, record:

```markdown
- Source: <site/page name>
- URL: <url>
- Checked: <date and timezone>
- Claim verified: <hours / closure / rule / schedule / price>
- Confidence: confirmed | likely | needs day-before check | unknown
```

## Red Flags

- Opening hours shown only on a map listing, with no official confirmation.
- Itinerary stop has only a generic Google Maps search link or broad area link.
- Mapcode and Google Maps pin disagree without a note explaining which navigation target to use.
- Community reports older than the current season.
- A place described as 車中泊-friendly without explicit overnight, parking, or toilet details.
- Japan self-drive stops without Mapcode or fallback navigation input.
- Parking candidates labeled 月極, resident-only, staff-only, permit-only, apartment parking, or unclear private lot.
- Non-driving itineraries that omit required reservation timing, last-return risk, or fare estimate.
- Driving itineraries that compare routes without fuel/charging, toll, parking, ferry/bridge, or rental-cost assumptions.
- A bath, shower, or laundromat stop placed after its last entry or likely closing time.
- Winter, typhoon, heavy rain, wildfire, landslide, or road restriction risk.
- Ferry, ropeway, mountain road, or viewpoint plans that ignore wind, visibility, snow, or wave forecasts.
- Last-entry time too close to arrival.
- Backup location more than 30-60 minutes away for late-night plans.
