---
id: 2026-01-01-0099-fixture-sub-missing-parent
plan_kind: sub
status: draft
owner: example-owner
created: 2026-01-01
required_for_completion: true
sub_plan_reason: deliberately omits the parent field for validatePlanTreeFrontmatter testing
---

# Fixture: sub plan missing `parent`

> **Fixture purpose**: should be rejected by `validatePlanTreeFrontmatter`
> with `missing: parent`. Files inside `fixtures/` are excluded from
> cross-plan scans so the bad shape does not pollute the live registry.
