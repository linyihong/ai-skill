# Viewport Measurement Authority

Principle: viewport measurement is evidence, not authority by default.

## Rule

- Treat viewport measurements as observations that can become stale.
- Do not persist viewport measurements across sessions unless the product explicitly defines a restore contract.
- If JavaScript measurement drives layout, list the state transitions that invalidate it.
- Prefer current browser layout primitives before adding measurement-driven positioning.

## Required Invalidation Triggers

If measurement-driven layout is unavoidable, account for:

- resize;
- orientation change;
- browser chrome expansion or collapse;
- visual viewport scroll or resize;
- page reload after viewport change;
- app restore or route restore;
- keyboard or system overlay when relevant.

## Failure Class

`stale_viewport_measurement` means the UI is positioned or sized from an old observation rather than the current layout reality.
