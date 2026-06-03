# folder-shape-violations (fixture)

This directory holds plan files that violate the **warning-only** folder
convention enforced by `validatePlanTreeFolderConvention`:

| File / path | Violation |
|---|---|
| `bad-name.md` | filename does not match `_plan.md` or `^\d{2}-` prefix |
| `level3-bad/deep/extra-deep.md` | nested depth ≥ 3 under `plans/active/` |

These are **fixtures only** (path segment `fixtures/` excludes them from
cross-plan scans). Real plan folders should follow the `_plan.md` /
`NN-<slug>.md` convention with depth < 3.
