# Embedded hardware product flow belongs in app development guidance
# Extracted — See [`workflow/software-delivery/execution-flow.md`](../../../../workflow/software-delivery/execution-flow.md)

Status: promoted

## Lesson

Some products are not mobile/web/backend apps, but they still follow the same contract-first development discipline. Embedded firmware and hardware-backed products need product behavior, BDD, domain models, public interfaces, error handling, implementation slices, and tests, plus hardware-specific contracts for datasheets, protocols, board context, driver boundaries, target validation, and bring-up evidence.

## Rule

When a project involves firmware, sensors, boards, UART/I2C/SPI/BLE/CAN/GPIO, RTOS tasks, or hardware-in-loop validation:

1. Keep it in `app-development-guidance` as an app/product development flow unless the user needs a separate skill for hardware lab operations, flashing automation, schematic/PCB work, or toolchain-specific runbooks.
2. Add embedded guidance under `platforms/embedded/` and implementation details under `implementation/embedded/`.
3. Require datasheet/protocol contracts, hardware context contracts, driver/service/application boundaries, host fixtures, and hardware-in-loop evidence.
4. Do not rename the skill only because one hardware-backed project appears; prefer broadening language from app-only to app/product while preserving existing cross-skill links.

## Required Linked Updates

- `SKILL.md`: updated to include embedded, firmware, and hardware-product triggers.
- `README.md`: updated to describe app/product scope and embedded linked docs.
- `WORKFLOW.md`: updated with embedded owner layers, blockers, classification, and validation.
- `process/README.md`: updated with Embedded / Hardware Product Flow.
- `platforms/embedded/README.md`: added platform guidance.
- `implementation/embedded/README.md`: added contract-to-code implementation guidance.
- `CHECKLIST.md` and `checklists/embedded-firmware-review.md`: updated with review gates.
- `templates/initial-development-docs.md`: updated with Hardware / Firmware Contract fields.

## Validation

Use an existing embedded project as evidence only for reusable patterns. Do not copy project-specific pin choices, logs, board conclusions, or private hardware details into the reusable skill.
