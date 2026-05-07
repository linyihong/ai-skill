# Travel Planning Skill

This skill supports source-backed travel planning: itinerary design, current opening-hour checks, weather-aware route ordering, seasonal feasibility, transport timing, lodging or car-stay routing, road-trip support stops, and backup planning.

## Goals

- Turn a destination and date range into a practical route.
- Verify time-sensitive details before recommending a stop.
- Separate confirmed facts from assumptions and open questions.
- Use weather forecasts to choose better route order and realistic backup plans.
- Use community sources for discovery while grounding decisions in official or current sources.
- Make car-stay and road-trip plans realistic: legal overnight status, toilets, bathing, laundry, trash rules, noise rules, weather, road conditions, and backup lodging.

## What Belongs Here

- Reusable workflows for planning trips.
- Source hierarchy and verification rules.
- Output templates for itineraries, weather/backup logic, support-stop tables, source tables, and day-before checklists.
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
3. Check weather and local disruption risks before locking route order.
4. Plan route order with buffers, backups, and support stops.
5. Mark confidence and unresolved checks clearly.
6. Convert new reusable planning lessons into `feedback_history/` when they generalize.
