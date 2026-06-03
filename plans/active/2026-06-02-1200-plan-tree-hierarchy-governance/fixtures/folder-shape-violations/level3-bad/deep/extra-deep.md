---
id: 2026-01-01-0092-fixture-depth-three
plan_kind: sub
status: draft
owner: example-owner
created: 2026-01-01
parent: 2026-01-01-0000-example-main-plan
required_for_completion: false
sub_plan_reason: nested 3 levels deep — warning-only depth advisory
---

# Fixture: plan nested at depth ≥ 3

> **Fixture purpose**: this file lives at depth 4 under `plans/active/`
> (`level3-bad/deep/extra-deep.md`). `validatePlanTreeFolderConvention`
> emits a non-blocking advisory suggesting the tree be split into an
> independent main plan instead of nesting deeper.
