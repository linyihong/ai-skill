---
id: 2026-01-01-0094-fixture-archive-required-pending
plan_kind: sub
status: in-progress
owner: example-owner
created: 2026-01-01
parent: 2026-01-01-0000-example-main-plan
required_for_completion: true
sub_plan_reason: simulates an in-progress required sub that should block parent archive
---

# Fixture: required sub-plan still in-progress

> **Fixture purpose**: if a main plan with id `2026-01-01-0000-example-main-plan`
> were moved into `plans/archived/` while this sub-plan still has
> `status: in-progress` and `required_for_completion: true`, then
> `validatePlanTreeArchiveOrder` must block the commit.
