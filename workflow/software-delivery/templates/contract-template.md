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

## Screen Mapping
| Scenario | Screens / Surfaces | APIs / Events | Tables / Stores | Validation Target |
|----------|--------------------|---------------|-----------------|-------------------|
| <BDD scenario> | <screen/route/CLI/SDK/job> | <operation/event/provider> | <table/store/cache/queue> | <BDD/contract/UI/integration proof> |

## Consumer Contract
| Consumer | Intent | Needs | Freshness | Loading | Empty / Error Behavior | Permissions |
|----------|--------|-------|-----------|---------|------------------------|-------------|
| <screen/CLI/SDK/job> | <actor intent> | <data/command/event> | <refresh/cache policy> | <skeleton/progress> | <empty/retry/fallback> | <role/capability> |

## UI Behavior / Screen Contract
| Screen / Flow | State | Action | Validation | Feedback | Navigation | Events |
|---------------|-------|--------|------------|----------|------------|--------|
| <screen id> | <loading/empty/error/success> | <create/edit/delete/submit/retry> | <rules> | <toast/inline/banner/focus> | <route/modal/back> | <emitted/consumed events> |

## Frontend ViewModel Contract
| ViewModel | Source | Field | Derivation / Formatting | Null / Error Behavior | Fixture |
|-----------|--------|-------|-------------------------|-----------------------|---------|
| <VM name> | <API/domain/event/local state> | <display field> | <rule> | <fallback/mapping> | <source -> expected VM> |

## Accessibility Contract
| Surface | Keyboard / Focus | Semantics | Assistive Feedback | Validation |
|---------|------------------|-----------|--------------------|------------|
| <screen/component> | <tab/focus/escape behavior> | <role/label/live region> | <success/error/progress announcement> | <lint/manual/test evidence> |

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
- **Downstream**: → BDD Scenarios → Screen Mapping → Implementation Plan → Review Report
- **Screen Traceability**: <BDD -> Screen -> Consumer -> API/Event -> Table/Store -> Validation Target>
- **Linked Artifacts**:
  - **Change Brief**: <link>
  - **BDD Scenarios**: <link>
  - **UI Contracts**: <link when consumer surface exists>
  - **Implementation Plan**: <link after creation>
