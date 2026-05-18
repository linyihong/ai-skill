# Contract: <feature/component>

## Domain Model
<Key entities and relationships>

## Architecture Decision
**Context**: <problem statement>
**Decision**: <chosen approach>
**Consequence**: <trade-offs accepted>

## API Contract
| Endpoint | Method | Request | Response | Error Codes |
|----------|--------|---------|----------|-------------|
| <path> | <GET/POST/PUT/DELETE> | <body/params> | <shape> | <list> |

## Error Handling
| Scenario | Error | Handling Strategy |
|----------|-------|-------------------|
| <condition> | <code> | <retry/fallback/fail-fast> |

## Data Model Changes
- **Schema**: <description>
- **Migration**: <forward/backward plan>
- **Backward Compatibility**: <yes/no + reason>

## Traceability
- **Upstream**: <link to change-brief-template.md>
- **Downstream**: → BDD Scenarios → Implementation Plan → Review Report
- **Linked Artifacts**:
  - **Change Brief**: <link>
  - **BDD Scenarios**: <link>
  - **Implementation Plan**: <link after creation>
