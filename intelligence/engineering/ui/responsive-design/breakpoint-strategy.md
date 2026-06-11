# Breakpoint Strategy

Principle: breakpoints describe layout constraints, not device SKUs.

## Why

Device names age quickly and do not cover browser zoom, split-screen windows, embedded webviews, safe-area changes, or dynamic browser chrome. Breakpoints should name the constraint the UI must survive.

## Render Context Vocabulary

- `narrow_mobile`: very constrained width and height.
- `mobile`: common single-column touch layout.
- `large_mobile`: wider phone or compact tablet-like width.
- `tablet`: larger viewport with room for denser layout.
- `desktop`: wide pointer-capable context.
- `landscape`: reduced height or altered aspect ratio.
- `safe_area`: viewport affected by notches, home indicators, or system overlays.

## Rule

- Prefer mobile-first `min-width` breakpoints for progressive enhancement.
- Use `max-width` only for correcting narrow-context density.
- Prefer container queries when a component is embedded in variable-width parents.
- Validate at representative viewport sizes, but do not make the viewport size the semantic contract.
