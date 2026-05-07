# Initial Development Docs Template

Use this template when turning a product brief into an implementation plan. Keep it lightweight at first; split into separate files when any section grows.

## Product Brief

| Field | Notes | Validation status | Evidence / reference source |
| --- | --- | --- | --- |
| Goal | | validated / assumption / open question / scoped out / invalidated | |
| Users / actors | | validated / assumption / open question / scoped out / invalidated | |
| Scope | | validated / assumption / open question / scoped out / invalidated | |
| Non-goals | | validated / assumption / open question / scoped out / invalidated | |
| Assumptions | Include owner, validation plan, and impact if false | validated / assumption / open question / scoped out / invalidated | |
| Success criteria | Test / metric / demo / release checklist / manual evidence | validated / assumption / open question / scoped out / invalidated | |
| Constraints | Security / privacy / platform / hardware / budget / schedule / compatibility | validated / assumption / open question / scoped out / invalidated | |
| Dependencies | External services / vendors / teams / generated clients / data / hardware | validated / assumption / open question / scoped out / invalidated | |
| Risks | Abuse / failure / safety / privacy / replay / data loss / operation | validated / assumption / open question / scoped out / invalidated | |
| Performance expectations | Latency / throughput / resource / startup / batch / concurrency / operational capacity | validated / assumption / open question / scoped out / invalidated | |
| Open questions | | open question | |

Validation status:

- `validated`: supported by evidence or explicit decision.
- `assumption`: plausible but not proven; needs owner and validation plan.
- `open question`: blocks impacted implementation until answered or scoped out.
- `scoped out`: explicitly not part of current work.
- `invalidated`: contradicted by evidence; revise brief before using it.

## Change Intake

Use this section before code changes.

| Field | Notes |
| --- | --- |
| Planning artifact reviewed | 企劃書 / product brief / issue / ticket / PRD / design note / BDD / API contract / other |
| Product brief validation result | validated / has assumptions / has blockers / not applicable |
| Change type | new requirement / bug fix / refactor / hardening / documentation-only |
| New requirement summary | |
| Bug expected behavior | |
| Bug actual behavior | |
| Reproduction / evidence | |
| Affected BDD scenarios | |
| Affected contracts | Domain Model / Architecture / API / Interface / Error Handling / Test Plan |
| Required planning updates before code | |
| Regression or validation test | |
| Performance impact | none expected / load / stress / spike / soak / smoke check needed |
| Blocker questions | |

## Contract Governance

Use this section when the project has multiple docs or was implemented before specs were backfilled.

| Field | Notes |
| --- | --- |
| Document precedence | Governance/framework contract -> product plan -> BDD -> contracts -> implementation -> tests, or project-specific order |
| Stable ID scheme | feature ID / rule ID / operation ID / route / command / diagnostic / scenario tag |
| Minimum doc-sync matrix | What docs/tests change for API, permission, DB, UI, generated client, vendor, CLI/tooling, release changes |
| Canceled / deferred / out-of-scope items | |
| Not-tool-enforceable / manual-only rules | |

## Test Strategy

| Field | Notes |
| --- | --- |
| Existing behavior to guard | |
| New / changed behavior to prove | |
| BDD scenarios written before code | yes/no |
| Failing tests or executable specs before code | unit / contract / property / integration / none with reason |
| Changed/new-code coverage target | |
| Total coverage context | |
| Mutation / negative checks | |
| Property or invariant tests | |
| Database / persistence fixtures | |
| Contract tests | |
| Generated client / SDK compile check | |
| Fixture-backed validation | |
| Performance budget | P95/P99 latency / throughput / error rate / resource ceiling |
| Performance test type | load / stress / spike / soak / CI smoke / not applicable with reason |
| Performance evidence location | CI run / report / dashboard / release checklist / manual note |
| Host fixture tests | |
| Hardware-in-loop / target evidence | |
| AI-generated code review focus | intent / edge cases / security / ownership / error handling |

## Existing Project Documentation Gap Audit

Use this section when the project is already implemented and documents are missing.

| Document | Status (`exists` / `partial` / `missing` / `unknown`) | Evidence source | Backfill action | Blocker questions |
| --- | --- | --- | --- | --- |
| Product Brief | | | | |
| Bounded Context Map | | | | |
| BDD Behavior | | | | |
| Domain Model Contract | | | | |
| Architecture Contract | | | | |
| API / Interface Contract | | | | |
| Hardware / Firmware Contract | | | | |
| Error Handling Contract | | | | |
| Test Plan | | | | |
| Contract Governance / Traceability | | | | |

Product Brief gaps may remain `unknown` if original intent is unavailable. BDD Behavior must be completed from observable implementation evidence.

## Blocker Questions

List every missing item that affects behavior, contracts, ownership, security, storage, error handling, or tests. Do not continue implementation until each blocker is answered, backed by evidence, or explicitly scoped out.

| Question | Blocks | Asked to | Answer / decision | Status |
| --- | --- | --- | --- | --- |

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

Traceability:

| Scenario / rule ID | Product source | Code refs | Test refs | Validation status |
| --- | --- | --- | --- | --- |

Validation status: `automated` / `fixture-backed` / `manual-evidence` / `pending-runner` / `not-automatable`.

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
- Embedded target / board / sensor constraints:
- Task / ISR / queue / timing constraints:

## Hardware / Firmware Contract

Use this section for embedded, firmware, sensor, board, protocol, or hardware-backed products.

| Field | Notes |
| --- | --- |
| Datasheet / vendor spec source | |
| Protocol / electrical interface | UART / I2C / SPI / BLE / CAN / GPIO / other |
| Board / module / sensor version | |
| Hardware context | Pins, bus, baud/rate, buffers, power, timing, injected config |
| Driver layer owns | |
| Service / parser layer owns | |
| Domain model owns | |
| Application layer owns | |
| Must not cross boundary | |
| Host fixture source | Vendor example / captured log / synthetic fixture / other |
| Hardware-in-loop evidence needed | |
| Bring-up notes | Board, wiring, firmware, flash command, logs, measurements, deviations |

## API / Interface Contract

| Operation | Consumer | Provider | Request / Input | Response / Output | Errors |
| --- | --- | --- | --- | --- | --- |

Contract format:

- OpenAPI / GraphQL / RPC / event schema / TypeScript type / CLI contract:
- Versioning:
- Compatibility rule:
- Mock / fixture source:
- Generated client / SDK command:
- Generated output checked in: yes/no:

## Vendor / Third-Party Integration Contract

Use this section only when an external provider is involved.

| Field | Notes |
| --- | --- |
| Vendor source doc location | Raw vendor docs stay in the project repository, not reusable skills |
| Sanitized integration excerpt | Operations actually used, without secrets or account-specific details |
| Auth / signing / credential boundary | |
| Idempotency / replay / retry / timeout | |
| Webhook / callback verification | Signature, timestamp, body binding, source, dedupe |
| Sandbox vs live differences | |
| Fixture source | |
| Live test gate | Environment/config required to run |
| Logging / redaction | |

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
