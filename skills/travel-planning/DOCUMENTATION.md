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
- Planning assumptions:

## Recommendation

<Short explanation of the route logic and why this plan fits the user.>
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
| Candidate | Overnight Status | Toilet | Bath / Shower | Rules / Fees | Risk | Backup | Confidence |
| --- | --- | --- | --- | --- | --- | --- | --- |
| <place> | allowed / listed / unclear / not allowed | <hours> | <nearby option> | <quiet hours / trash / cooking / gate> | <weather / safety / noise> | <legal backup> | confirmed / needs confirmation |
```

## Risk and Backup Section

```markdown
## Risks and Backups

- Closure risk:
- Weather risk:
- Crowd / event risk:
- Transport risk:
- Overnight-stay risk:
- Backup route:
- What to check the day before:
```

## Final Answer Shape

When replying to the user, keep the plan readable:

1. State the assumptions.
2. Give the recommended itinerary.
3. List key source-backed checks and confidence.
4. Call out reservations and day-before checks.
5. Offer 1-2 meaningful alternatives when the current plan has uncertainty.

Do not bury blockers. If a core attraction, transport leg, or overnight candidate is unverified, say that before the detailed schedule.
