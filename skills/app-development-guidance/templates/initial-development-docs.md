# Initial Development Docs Template

Use this template when turning a product brief into an implementation plan. Keep it lightweight at first; split into separate files when any section grows.

## Product Brief

| Field | Notes |
| --- | --- |
| Goal | |
| Users / actors | |
| Scope | |
| Non-goals | |
| Assumptions | |
| Open questions | |

## Bounded Context Map

| Context / Module | Responsibility | Owns | Depends on | Does not own |
| --- | --- | --- | --- | --- |

## BDD Behavior

```gherkin
Feature:

Scenario:
  Given
  When
  Then
```

## Domain Model Contract

| Concept | Type | Responsibility | Invariants | Notes |
| --- | --- | --- | --- | --- |

Commands:

| Command | Input | Preconditions | Result / Event | Failure cases |
| --- | --- | --- | --- | --- |

## Architecture Contract

| Layer / Component | Owns | May depend on | Must not depend on | Notes |
| --- | --- | --- | --- | --- |

Runtime / deployment constraints:

- Persistence:
- External services:
- Auth/session:
- Background jobs/events:
- Observability:

## API / Interface Contract

| Operation | Consumer | Provider | Request / Input | Response / Output | Errors |
| --- | --- | --- | --- | --- | --- |

Contract format:

- OpenAPI / GraphQL / RPC / event schema / TypeScript type / CLI contract:
- Versioning:
- Compatibility rule:
- Mock / fixture source:

## Error Handling Contract

| Error code / type | Trigger | User-visible message | Retryable | Logged fields | Redaction |
| --- | --- | --- | --- | --- | --- |

Rules:

- Validation errors:
- Auth/session errors:
- Permission errors:
- Rate limit / abuse errors:
- External dependency failures:
- Unexpected errors:

## Test Plan

| Test level | Proves | Owner | Required before merge |
| --- | --- | --- | --- |
| Unit | Domain invariants and pure logic | | yes/no |
| BDD | Critical behavior scenarios | | yes/no |
| Contract | Provider/consumer compatibility | | yes/no |
| Integration | Real adapters and cross-context flow | | yes/no |

## Implementation Slices

| Slice | Context | Provider work | Consumer work | Contract dependency | Test dependency |
| --- | --- | --- | --- | --- | --- |
