# CSS Layout Ownership

Principle: when layout can be expressed declaratively, CSS should own it.

## Rule

- Use CSS for width, max-width, centering, safe-area padding, fixed-surface alignment, responsive units, and container-based adaptation.
- Use JavaScript for user interaction, data state, observation, and behavior that CSS cannot express.
- If JavaScript writes layout values, name the invalidation triggers and validation evidence.
- Prefer CSS custom properties for shared layout contracts over duplicated component constants.

## Why

The browser layout engine already recalculates on viewport, zoom, font, safe-area, and container changes. JavaScript layout measurement creates a second authority that must manually track all the same transitions.

## Smells

- `visualViewport` measurement is used to position a fixed surface that CSS can center.
- JavaScript calculates width from constants also present in CSS.
- A component has mount-time layout state but no resize, orientation, safe-area, or restore invalidation plan.
