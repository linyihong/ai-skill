# Embedded Implementation Patterns

Use this page to turn embedded or hardware-product contracts into buildable firmware slices.

## Contract-To-Code Slices

| Contract | Implementation slice | Tests |
| --- | --- | --- |
| Datasheet / Protocol Contract | Constants, frame definitions, command tables, valid ranges, and documented deviations. | Compile-time checks, fixture examples, boundary tests. |
| Protocol Parsing Contract | Stream parser, state machine, serializer, ACK/report decoder, resynchronization, malformed-frame handling. | Host fixture tests, negative frames, length/range/property tests. |
| Domain Model Contract | Pure structs/enums/value objects with units and invariants. | Unit tests for conversion, invariant, and boundary rules. |
| Hardware Context Contract | Context/config struct for pins, buses, baud/rates, buffers, power, and board-specific options. | Mock or host tests for config validation; target smoke test for actual board wiring. |
| Embedded Architecture Contract | Driver/service/application files, task/ISR boundary, queues, locks, lifecycle, callbacks. | Concurrency tests where possible; bench logs for target-only timing. |
| Public API / Interface Contract | Context create/destroy, commands, subscriptions/callbacks, error mapping, multi-device behavior. | Provider/consumer contract tests, examples, mock-driver tests. |
| BDD Behavior | Application logic that maps domain events to product behavior. | BDD/executable specs, simulator tests, hardware-in-loop only for physical proof. |

## Implementation Order

1. Add or update contracts before code: datasheet/protocol, BDD, domain, hardware context, architecture, public API, error handling, and test plan.
2. Write host-side fixtures and failing tests for parser/domain/API behavior.
3. Implement pure parsing and domain conversion without target hardware dependencies.
4. Implement hardware context and driver boundaries with injectable board configuration.
5. Add service/application logic that consumes domain objects instead of raw bytes or registers.
6. Add target or hardware-in-loop checks for wiring, timing, sensor behavior, power, and real device interactions.
7. Record bring-up evidence and update the contract if measured behavior differs from the spec.

## Design Traps

- Do not let application code call bus read/write APIs directly when a service/domain boundary should own protocol semantics.
- Do not hide product behavior inside driver code.
- Do not duplicate context lifecycle APIs for the same device path unless the contracts are revised first.
- Do not use global parser/session/subscriber state when multiple devices, buses, or board variants are possible.
- Do not treat bench logs as a replacement for host-repeatable fixtures when parsing and domain behavior can be tested offline.
- Do not hard-code lab pins, baud rates, calibration values, debug modes, or fixture payloads into production paths.

## Bring-Up Note Shape

Record target-only evidence in the product repository:

```markdown
## Hardware Bring-Up

- Board / revision:
- Sensor / module / firmware:
- Wiring / pins / bus settings:
- Firmware build / commit:
- Flash or run command:
- Scenario:
- Expected behavior:
- Observed logs / measurements:
- Deviations from datasheet or contract:
- Follow-up blocker questions:
```

## Required Linked Updates

Follow [`../../../../shared-rules/linked-updates.md`](../../../../shared-rules/linked-updates.md). Embedded implementation patterns must update or verify [`../../platforms/embedded/`](../../platforms/embedded/), [`../../CHECKLIST.md`](../../CHECKLIST.md), [`../../checklists/`](../../checklists/), and templates when the project-document shape changes.
