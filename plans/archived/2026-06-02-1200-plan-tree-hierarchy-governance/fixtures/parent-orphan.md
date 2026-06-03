---
id: 2026-01-01-0096-fixture-parent-orphan
plan_kind: sub
status: draft
owner: example-owner
created: 2026-01-01
parent: 0000-00-00-0000-nonexistent-main-plan
required_for_completion: false
sub_plan_reason: parent id deliberately does not resolve to any real plan
---

# Fixture: sub plan with dangling `parent` pointer

> **Fixture purpose**: should be rejected by `validatePlanTreeParentReference`
> because `parent` references an id that has no corresponding plan file
> in `plans/active/` or `plans/archived/`.
