# Travel Planning Output Templates

Use these templates when producing trip plans or documenting reusable planning work.

## Itinerary Summary

```markdown
## Trip Frame

- Destination:
- Dates:
- Party:
- Transport:
- Pace:
- Must-do:
- Constraints:
- Weather assumptions:
- Planning assumptions:

## Recommendation

<Short explanation of the route logic, weather-aware ordering, support stops, and why this plan fits the user.>
```

## Day Plan

```markdown
## Day <n> - <theme / area>

| Time | Plan | Travel / Access | Validation | Backup |
| --- | --- | --- | --- | --- |
| 09:00 | <stop> | <route / drive / transit> | <hours/source/confidence> | <nearby fallback> |

### Notes

- Reservations:
- Tickets / passes:
- Food / rest:
- Weather or seasonal risk:
- Bathing / laundry / fuel / charging:
- Day-before check:
```

## Weather-Aware Options

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

## Source Table

```markdown
| Item | Source | Checked | Claim | Confidence | Follow-up |
| --- | --- | --- | --- | --- | --- |
| <place> | <official page / operator / map / community source> | <date TZ> | <hours / rule / schedule> | confirmed / likely / needs day-before check / unknown | <reserve / call / recheck> |
```

## 車中泊 Candidate Table

```markdown
| Candidate | Overnight Status | Toilet | Bath / Shower | Laundry | Rules / Fees | Risk | Backup | Confidence |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| <place> | allowed / listed / unclear / not allowed | <hours> | <nearby option> | <coin laundry / none / unknown> | <quiet hours / trash / cooking / gate> | <weather / safety / noise> | <legal backup> | confirmed / needs confirmation |
```

## Road-Trip Support Table

```markdown
| Need | Candidate | Timing | Checks | Confidence | Backup |
| --- | --- | --- | --- | --- | --- |
| Bath / shower | <onsen / sento / coin shower> | <before overnight / morning> | <hours / last entry / parking / fee> | confirmed / needs recheck | <alternative> |
| Laundry | <coin laundry> | <while eating / bathing / morning> | <hours / parking / drying time> | confirmed / likely / unknown | <alternative> |
| Fuel / charging | <station / charger> | <before rural route> | <hours / payment / plug type> | confirmed / likely | <alternative> |
```

## Risk and Backup Section

```markdown
## Risks and Backups

- Closure risk:
- Weather risk:
- Weather-based route swap:
- Crowd / event risk:
- Transport risk:
- Overnight-stay risk:
- Bathing / laundry support risk:
- Backup route:
- What to check the day before:
```

## Final Answer Shape

When replying to the user, keep the plan readable:

1. State the assumptions.
2. Give the recommended itinerary and explain weather-aware ordering.
3. List key source-backed checks and confidence.
4. Call out reservations, support stops, and day-before checks.
5. Offer 1-2 meaningful alternatives when the current plan has weather, closure, or availability uncertainty.

Do not bury blockers. If a core attraction, transport leg, or overnight candidate is unverified, say that before the detailed schedule.
