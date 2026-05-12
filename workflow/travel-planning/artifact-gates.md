# Travel Planning Artifact Gates

本文件定義旅行規劃產出的格式規範與品質門檻。承接 [`skills/travel-planning/DOCUMENTATION.md`](../../skills/travel-planning/DOCUMENTATION.md) 的內容，提取為 tool-neutral 的 artifact gates。

> **相容性規則**：`skills/travel-planning/DOCUMENTATION.md` 仍為 active skill entrypoint。本文件為 reference target，兩者應保持同步。

## 1. Itinerary Summary

```markdown
## Trip Frame

- Destination:
- Dates:
- Party:
- Transport:
- Transport priority:
- Long-distance transport comparison:
- Pace:
- Must-do:
- Constraints:
- Weather assumptions:
- Country / region rules:
- Travel agency / model-course references:
- Location verification:
- Lodging / overnight base:
- Route shape:
- Activity / food focus:
- Restaurant screening:
- Calendar / app output:
- Fuel / charging assumptions:
- Schedule pressure:
- Planning assumptions:

## Recommendation

<Short explanation of the route logic, weather-aware ordering, support stops, and why this plan fits the user.>
```

## 2. Day Plan

```markdown
## Day <n> - <theme / area>

| Time | Plan | Travel / Access | Validation | Backup |
| --- | --- | --- | --- | --- |
| 09:00 | <stop + Google Maps exact place/pin> | <route / drive / transit> | <hours/source/location confidence> | <nearby fallback> |

### Notes

- Reservations:
- Transport booking:
- Fare / cost:
- Tickets / passes:
- Food / rest:
- What to do:
- Local food / restaurant:
- Rating / review source:
- Weather or seasonal risk:
- Exact location / Google Maps:
- Driving parking pin:
- Calendar / app fields:
- Navigation / parking:
- Lodging / overnight base:
- Route shape:
- Schedule pressure:
- Bathing / laundry / fuel / charging:
- Day-before check:
```

## 3. Weather-Aware Options

```markdown
## Weather Strategy

| Condition | Better Plan | Backup / Swap | Source / Confidence |
| --- | --- | --- | --- |
| Sunny / clear morning | <outdoor / scenic / ferry / mountain plan> | <later indoor option> | <forecast source / confidence> |
| Rain / strong wind / low visibility | <indoor / food / bath / laundry / city plan> | <route swap or rest option> | <forecast source / confidence> |

### Route Order Reason

- Why this order fits the forecast:
- What changes if the forecast worsens:
- Day-before check:
```

## 4. Source Table

```markdown
| Item | Source | Checked | Claim | Confidence | Follow-up |
| --- | --- | --- | --- | --- | --- |
| <place> | <official page / operator / map / community source> | <date TZ> | <hours / rule / schedule> | confirmed / likely / needs day-before check / unknown | <reserve / call / recheck> |
```

## 5. Calendar / App-Ready Table

Use this when the user wants to add the itinerary to a calendar, reminders, map list, notes app, travel planning app, or offline map.

```markdown
| Day / Time | Event Title | Start / End | Time Zone | Location / Pin | Notes | Reminder | App / Map Group | Import Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| <day + time> | <short title> | <start-end> | <TZ> | <exact place / driving parking pin / station> | <reservation / source / caveat> | <when to alert> | <Day 1 sightseeing / food / backup / support> | ready / needs recheck / do not import yet |
```

Calendar/app notes:

- Use stable, short event titles that make sense on a phone lock screen.
- For self-drive, put the practical parking pin in the location field when it differs from the attraction entrance.
- Mark weather-dependent, reservation-pending, unverified, or backup-only items as `needs recheck` or `do not import yet`.
- Group map pins by day and category: sightseeing, food, lodging, parking, support stops, backups.
- Include reminder offsets for departure, booking, last entry, last order, check-in, fuel/charging, and day-before weather checks.

## 6. Offline / Save-Before-Departure Checklist

```markdown
| Item | Save / Prepare | Why | Confidence |
| --- | --- | --- | --- |
| Offline map | <area / route> | <low signal / rural / overseas> | confirmed / recommended |
| Reservation | <ticket / hotel / ferry / restaurant> | <check-in / boarding / proof> | confirmed / needs booking |
| Route backup | <screenshot / note / phone number> | <bad signal / late arrival / closure> | recommended |
```

## 7. Travel Agency / Model-Course Benchmark

Use this when agency tours, package tours, bus tours, or official model courses inform the route.

```markdown
| Source | Use Type | Price | Included / Excluded | Borrowed Idea | Verification Needed | Change Made / Caveat |
| --- | --- | --- | --- | --- | --- | --- |
| <agency/model course> | direct package option / benchmark only | <per person / group / unknown> | <transport / meal / ticket / lodging included; exclusions> | <route order / lunch / stop duration / base area> | <hours / pin / parking / ticket / weather / cancellation> | <adapted for user pace / self-drive / weather / no backtracking / user warning> |
```

Benchmark note:

- Direct package options must show price, included/excluded items, booking/cancellation notes, and user-facing caveats.
- Do not treat an agency itinerary as proof that a stop is open, reachable, or suitable today unless independently verified.
- Explain when an agency route assumes charter buses, group meals, shopping stops, or faster group movement.

## 8. Stop Experience Table

Use this for key sightseeing, food, and support stops.

```markdown
| Stop | Why Stop | What To Do | Suggested Time | Food / Local Specialty | Restaurant Screening | Backup / Nearby Alternative | Confidence |
| --- | --- | --- | --- | --- | --- | --- | --- |
| <place> | <scenery / food / hot spring / route support> | <1-3 concrete actions> | <duration> | <restaurant / market / specialty / fallback> | <Google Maps / local rating tool / reservation / route fit> | <nearby option> | confirmed / needs confirmation |
```

Stop note:

- If a stop is mainly for fuel, toilet, laundry, parking, or rest, label it as a support stop.
- Do not list food that is closed at the planned time or far off-route without explaining the detour.
- Use restaurant rating/review tools appropriate to the country. For Japan, include Google Maps and 食べログ when practical, plus official or reservation pages for hours and booking.

## 9. Restaurant Recommendation Table

Use this when meal stops materially affect the trip, the user asks for restaurant filtering, or the area has popular restaurants with queues, reservations, limited hours, or parking constraints.

```markdown
| Meal / Area | Candidate | Cuisine / Specialty | Google Maps Signal | Local Rating Tool Signal | Hours / Last Order | Reservation / Queue | Access / Parking | Route Fit | Backup | Confidence |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| <lunch / dinner / area> | <restaurant> | <local food> | <rating + review count / recent signal> | <Tabelog or local equivalent rating + review count / not available> | <hours / last order / closed day> | <book now / likely queue / walk-in> | <station / visitor parking / not 月極> | on-route / detour / optional | <nearby fallback> | confirmed / needs day-before check |
```

Restaurant note:

- Do not select only by numeric score. Explain the tradeoff among cuisine fit, review quality, review volume, recency, opening window, price, reservation/queue risk, and route efficiency.
- For Japan, use 食べログ as a local screening source when available and Google Maps for exact place/access/recent practical signals; note meaningful disagreements.
- For self-drive meal stops, include visitor-usable parking or a practical nearby paid parking option.

## 10. Exact Location Table

Use this for every meaningful stop when location identity matters.

```markdown
| Stop | Google Maps Exact Place / Pin | Driving Parking Pin | Official Name / Address Check | Mapcode Check | Practical Navigation Target | Location Confidence | Concern |
| --- | --- | --- | --- | --- | --- | --- | --- |
| <place> | <place URL / coordinate pin, not broad search> | <nearest visitor/official parking URL or not driving> | match / mismatch / not found | match / different / unavailable / not needed | entrance / visitor parking / station exit / pier / trailhead / reception | confirmed / needs confirmation | <ambiguity or none> |
```

Location note:

- Prefer exact Google Maps place links or coordinate pins that open one location.
- Avoid generic search links when many results can appear.
- For driving, the practical Google Maps target should be the nearest confirmed visitor-usable parking or official designated parking when available, not 月極 or unclear private parking.
- If Google Maps, Mapcode, official address, and practical parking/entrance point differ, explain which one should be used for navigation and why.

## 11. Non-Driving Transport Plan

Use this when the trip does not use a car, or when a day has public-transport legs.

```markdown
| Leg | Recommended Transport | Depart / Arrive | Transfer Buffer | Booking Needed | Fare Estimate | Risk / Backup |
| --- | --- | --- | --- | --- | --- | --- |
| <A to B> | <train / bus / ferry / flight / taxi> | <time window> | <minutes / station note> | <reserve by when / no reservation> | <per person / group> | <last return / weather / sold out / backup> |
```

Transport summary:

- Optimization reason: fastest / cheapest / fewer transfers / luggage-friendly / scenic / low walking.
- Tickets or passes to buy:
- Booking timing:

## 12. Self-Drive Cost Estimate

Use this when the user is driving.

```markdown
| Cost Item | Assumption | Estimate | Source / Confidence |
| --- | --- | --- | --- |
| Distance | <km> | <total km> | <route planner / map> |
| Fuel / charging | <unit price> | <total> | <price site / operator / range> |
| Tolls | <route / discount> | <total> | <toll calculator / operator> |
| Parking | <stops / overnight> | <total> | <parking operator / map> |
| Ferry / bridge | <crossings> | <total> | <operator page> |
| Rental car | <days / type / insurance> | <total> | <rental site / contract> |
| **Total** | | **<sum>** | |
```

## 13. 車中泊 Quietness Table

Use this for each 車中泊 candidate.

```markdown
| Candidate | Quietness Label | Noise Source | Toilet | Bathing / Laundry | Backup | Confidence |
| --- | --- | --- | --- | --- | --- | --- |
| <place> | quiet / moderate / noisy / unknown | <traffic / truck / crowd / none> | available / none / unknown | <nearby / none> | <nearby backup> | confirmed / candidate / needs confirmation |
```

## 14. Final Verification Checklist

Before delivering the itinerary, verify:

- [ ] Trip frame is documented (destination, dates, party, transport, pace, constraints)
- [ ] Source table covers all time-sensitive claims
- [ ] Agency/model-course references are labeled as `direct package option` or `benchmark only`
- [ ] Each stop has an exact Google Maps place link or precise pin
- [ ] Driving stops use visitor-usable parking as the navigation target
- [ ] Japan self-drive stops have Mapcode or fallback navigation input
- [ ] Location ambiguities are explicitly marked
- [ ] Weather strategy is documented with concrete swaps
- [ ] Long-distance transport comparison is done for 2+ hour transfers
- [ ] Transport mode decision is explicit (non-driving / self-drive / mixed)
- [ ] Overnight base is chosen for route logic, not convenience
- [ ] Route shape is checked for avoidable backtracking
- [ ] Schedule feasibility is labeled (comfortable / tight / too packed)
- [ ] Calendar/app-ready fields are included when useful
- [ ] 車中泊 quietness is labeled with backup
- [ ] Fuel/charging gaps are identified for self-drive plans
- [ ] Driving cost is estimated with assumptions
- [ ] All uncertain claims are labeled with confidence
- [ ] Next actions are specific (reserve, call, check, buy, download)
