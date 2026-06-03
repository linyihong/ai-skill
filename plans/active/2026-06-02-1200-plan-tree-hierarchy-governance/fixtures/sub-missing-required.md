---
id: 2026-01-01-0097-fixture-sub-missing-required
plan_kind: sub
status: draft
owner: example-owner
created: 2026-01-01
parent: 2026-01-01-0000-example-main-plan
sub_plan_reason: deliberately omits required_for_completion field
---

# Fixture: sub plan missing `required_for_completion`

> **Fixture purpose**: should be rejected by `validatePlanTreeFrontmatter`
> with `required_for_completion` missing. Authors must explicitly choose
> whether the sub-plan blocks parent archive — no default.
