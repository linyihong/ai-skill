# Context Budgeting

Context budgeting 將 model profile、capability dimensions 與 compression level 對齊。

## Budget Tiers

| Tier | Inputs | Escalate when |
| --- | --- | --- |
| `index` | Registry / index only. | Need source-backed action. |
| `summary` | Summary / checklist. | Conflict or edit needed. |
| `source` | Primary source. | Dependencies affect result. |
| `source_plus_deps` | Primary + required dependencies. | Cross-layer change. |
| `graph_assisted` | Source + graph / related sources. | Promotion, migration, deprecation. |

## Memory Boundary

Memory replay budget must not exceed context budget. If memory is more expensive than current source, read current source first.
