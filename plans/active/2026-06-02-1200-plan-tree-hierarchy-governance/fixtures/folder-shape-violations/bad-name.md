---
id: 2026-01-01-0093-fixture-bad-filename
plan_kind: sub
status: draft
owner: example-owner
created: 2026-01-01
parent: 2026-01-01-0000-example-main-plan
required_for_completion: false
sub_plan_reason: filename violates NN- prefix convention (warning, not block)
---

# Fixture: bad filename inside plan folder

> **Fixture purpose**: filename `bad-name.md` (no `NN-` prefix and not
> `_plan.md`) triggers `validatePlanTreeFolderConvention` to emit an
> advisory warning — the commit is **not** blocked.
