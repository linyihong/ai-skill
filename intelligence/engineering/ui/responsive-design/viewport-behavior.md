# Viewport Behavior

Principle: viewport width and height are runtime inputs, not constants.

## Why

Mobile browsers change usable height as address bars, keyboards, browser chrome, and safe areas appear or disappear. Developer tools device emulation can also change dimensions between initial load, resize, and reload.

## Rule

- Prefer dynamic viewport units such as `dvh` when full-height UI must track the currently visible viewport.
- Provide fallback paths for environments where `dvh` is unavailable.
- Avoid depending on a single `vh` or `svh` rule when the UI must survive browser chrome changes.
- Test both fresh load and resize-then-reload when a defect report mentions device switching.

## Smells

- Full-page UI uses only `height: 100vh` or only `height: 100svh`.
- JavaScript computes viewport dimensions once at mount and never reconciles after resize.
- Content is clipped by a full-height scroll root with disabled vertical scrolling.
