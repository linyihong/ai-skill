# Embedded Firmware Review Checklist

Use this checklist for firmware, sensor, board, protocol, and hardware-in-loop changes.

## Planning

- Change was classified as new requirement, bug fix, refactor, hardening, or docs-only.
- Datasheet/protocol source, board revision, hardware context, and BDD behavior were reviewed before code.
- Missing behavior, pin/bus mapping, timing, safety, protocol, error, or test information is resolved or explicitly scoped out.

## Contracts

- Datasheet/protocol facts are separated from product behavior and measured deviations.
- BDD scenarios describe observable domain/device behavior.
- Domain types define units, ranges, timestamps, and invariants without HAL/RTOS leakage.
- Hardware context is injectable per board and records pins, bus settings, buffers, power, and timing assumptions.
- Architecture contract defines driver/service/application boundaries, task/ISR ownership, queues, locks, and lifecycle.
- Public API contract defines context create/destroy, commands, callbacks/subscriptions, errors, and multi-device behavior.

## Implementation

- Raw bytes, registers, and bus calls stay in driver/protocol layers.
- Application code consumes domain objects, events, or public service APIs.
- Parser/session/subscriber state is per context when multiple devices or board variants are possible.
- Blocking APIs document timeout and calling-context rules.
- Debug/lab modes cannot silently ship as production defaults.

## Tests And Evidence

- Host tests cover parser fixtures, malformed frames, range/boundary cases, domain invariants, and API/error mapping.
- Rule-heavy behavior has property, invariant, mutation, or negative tests when ordinary examples are insufficient.
- Hardware-in-loop or manual bring-up is limited to target-only proof and records board, wiring, firmware, command, logs, measurements, and deviations.
- Safety behavior has fail-safe, timeout, debounce/cooldown, retry, and recovery validation.
- BDD, contracts, tests, and implementation paths are traceable.

## Release

- Target build and flash/run command are documented.
- Production config defaults match the intended board or deployment.
- Logs and raw payload dumps are gated or redacted for release.
- Known hardware limitations, calibration needs, errata, and rollback/recovery path are documented.
