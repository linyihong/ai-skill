---
id: 2026-01-01-0098-fixture-sub-empty-reason
plan_kind: sub
status: draft
owner: example-owner
created: 2026-01-01
parent: 2026-01-01-0000-example-main-plan
required_for_completion: true
sub_plan_reason: ""
---

# Fixture: sub plan with empty `sub_plan_reason`

> **Fixture purpose**: should be rejected by `validatePlanTreeFrontmatter`
> with `sub_plan_reason (non-empty)` missing. Empty string is treated
> identical to a missing field — `sub_plan_reason` must carry an actual
> rationale for why the sub-plan exists.
