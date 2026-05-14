# Embedded / Firmware / Hardware Products

Use this guidance when a project includes firmware, sensors, boards, UART/I2C/SPI/BLE/CAN/GPIO, RTOS tasks, hardware bring-up, or hardware-in-loop validation.

## Core Contracts

| Contract | Purpose |
| --- | --- |
| Datasheet / Protocol Contract | Records vendor spec, electrical interface, wire format, command/ACK/report shapes, timing, valid ranges, defaults, and errata. |
| Protocol Parsing Contract | Defines state machines, byte layout, checksums or tails, resynchronization, fixtures, and invalid-frame behavior. |
| BDD Behavior | Describes observable device/product behavior in domain terms, not raw registers or bus calls. |
| Domain Model Contract | Defines pure DTOs, units, ranges, timestamps, validity windows, and invariants without HAL/RTOS types. |
| Hardware Context Contract | Defines board-specific pins, buses, baud/rates, buffers, power modes, and injected configuration. |
| Embedded Architecture Contract | Defines driver/service/application boundaries, task/ISR ownership, queues, lifecycle, concurrency, and error escalation. |
| Public API / Interface Contract | Defines context lifecycle, callbacks/subscriptions, commands, errors, ownership, and multi-device behavior. |
| Test Plan | Splits host-repeatable tests from simulator, target, bench, or hardware-in-loop evidence. |

## Development Rules

- Treat the datasheet or protocol spec as the source of truth for bytes, timing, ranges, and electrical facts; record measured deviations separately.
- Keep raw bus data in driver/protocol layers. Application behavior should consume domain objects, events, or public APIs.
- Inject board-specific hardware context through config or context objects. Do not make one board's pins, UART, or bus settings the only production source of truth.
- Keep domain types free of HAL, RTOS, GPIO, UART, or task-handle types unless the contract explicitly defines an adapter boundary.
- Define task, ISR, queue, callback, and lock ownership before code. Blocking and timeout behavior must be documented on public APIs.
- For multi-device or multi-bus support, make state per context/instance. Avoid global parser, subscription, or command-session state unless the contract proves only one device can exist.
- Safety-related behavior needs explicit fail-safe state, timeout, debounce/cooldown, retry, recovery, and operator/user-visible behavior.

## Validation

Prefer this order:

1. Host-side unit tests for parsing, domain invariants, command/API contracts, and error mapping.
2. Fixture tests from vendor examples, captured logs, or synthetic positive and negative frames.
3. Property/invariant or fuzz-style tests for length, resynchronization, range, and malformed payload handling.
4. Simulator or mock-driver tests for application behavior.
5. Hardware-in-loop or bench tests only for evidence that requires the target board, sensor, timing, power, or physical environment.
6. Bring-up records with board revision, wiring, pins, bus settings, firmware version, flash command, logs, measurements, and observed deviations.

## Release Gate

- Target build and flashing path are documented and reproducible.
- Production defaults are separated from lab/debug modes.
- Debug logs, raw payload dumps, secrets, and calibration-only commands are gated for release.
- Hardware context defaults match the intended board revision.
- Protocol and domain fixtures are traceable to BDD or contract entries.
- Known hardware limitations, calibration requirements, errata, and rollback/recovery steps are documented.

## Required Linked Updates

Follow [`../../../../enforcement/linked-updates.md`](../../../../enforcement/linked-updates.md). Embedded/platform changes must update or verify [`../../implementation/embedded/`](../../implementation/embedded/), [`../../process/`](../../process/), [`../../CHECKLIST.md`](../../CHECKLIST.md), and relevant templates or review checklists.
