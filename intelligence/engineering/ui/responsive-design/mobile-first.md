# Mobile-First Responsive Design

Principle: design the narrowest viable viewport first, then progressively enhance for larger or less constrained contexts.

## Why

Narrow viewports have the fewest degrees of freedom: less width, less height, touch input, browser chrome, safe areas, and higher risk of clipped controls. If the smallest context is stable, larger contexts usually add layout options. If a desktop layout is built first, mobile often becomes a late-stage subtraction exercise.

## Rule

- Start with the smallest supported content width and height.
- Add enhancements with `min-width` or container-based conditions.
- Treat desktop-only success as incomplete when the contract includes mobile or narrow mobile.
- Do not encode device names as the rule; encode the render constraint.

## Smells

- `width: 1200px` as the base layout.
- Mobile rules that mostly undo desktop assumptions.
- Critical controls only tested in wide viewports.
