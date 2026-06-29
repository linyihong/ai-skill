# Change Brief: <title>

## Metadata
- **Change Type**: <feature|bugfix|refactor|replacement|migration|perf|docs|chore> <!-- legacy label; map to change_kind below -->
- **change_kind**: <feature|bugfix|replacement|migration|internal_refactor> <!-- canonical intake axis 1 -->
- **blocked_by_structure**: <true|false|n/a> <!-- feature only: untestable / god method / missing seam blocks landing feature locally -->
- **execution_mode** *(routed)*: <direct_change|preparatory_refactoring> <!-- axis 2; default direct_change when omitted -->
- **Priority**: <p0|p1|p2|p3>
- **Evidence Source**: <issue|incident|product-brief|customer-feedback>
- **Date**: <YYYY-MM-DD>

### change_kind ↔ Change Type（對齊表）

| Change Type | change_kind | blocked_by_structure | Typical execution_mode |
|-------------|-------------|----------------------|------------------------|
| `feature` | `feature` | `false` | `direct_change` |
| `feature` | `feature` | `true` | `preparatory_refactoring` |
| `bugfix` | `bugfix` | usually `false` | `direct_change` |
| `refactor` | `internal_refactor` | n/a | `direct_change` or structure-only batch |
| `replacement` | `replacement` | n/a | parity gate; not preparatory mode |
| `migration` | `migration` | n/a | parity gate; not preparatory mode |
| `perf` / `docs` / `chore` | *(out of execution_mode routing)* | n/a | n/a |

> `preparatory_refactoring` 是 execution mode，不是 change_kind。Routing 見 [`intake.md`](../intake.md) §change_kind × execution_mode 雙軸。

## Evidence Summary
<Link to source evidence, max 3 bullet points>

## Product Impact Alignment
- **Impact / journey artifact**: <link to product-impact-alignment-template.md if applicable>
- **Decision**: proceed | revise | reject | ask_user | not_applicable
- **Blocking mismatch**: <goal / actor / timing / pain / investment mismatch, if any>

## Scope
### In Scope
- <item>

### Out of Scope
- <item>

## Refactor / Replacement Parity
- **Applies**: <yes|no> <!-- yes for replacement/migration; no for feature+preparatory_refactoring -->
- **Inventory Artifact**: <link to parity inventory or not_applicable>
- **Legacy Surfaces Covered**: <old APIs / commands / scripts / UI flows / jobs / hooks / runtime surfaces>
- **Deferred / Not Planned**: <items and why they do not block this phase>
- **Validation Evidence**: <BDD / contract test / fixture / golden output / dry-run / manual review>

## Blocker Assessment
- [ ] No blocker — proceed to Requirements / Contract phase
- [ ] Blocker identified: <description>

## Traceability
- **Spec Source**: <link to spec-template.md if from Greenfield workflow>
- **Downstream**: → Product Impact Alignment → Contract → BDD Scenarios → Implementation Plan → Review Report
- **Linked Artifacts**:
  - **Product Impact Alignment**: <link after creation>
  - **Contract**: <link after creation>
  - **BDD Scenarios**: <link after creation>
  - **Implementation Plan**: <link after creation>
