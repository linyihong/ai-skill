# UI Engineering Intelligence

Reusable UI engineering judgment lives here. Keep this layer framework-neutral: record principles, trade-offs, and failure modes that survive Bootstrap, Tailwind, MUI, native CSS, or future UI systems.

## Topics

| Topic | Use When |
| --- | --- |
| [`responsive-design/`](responsive-design/README.md) | A UI must remain correct as viewport width, height, safe area, orientation, or container size changes. |
| [`layout-authority-governance/`](layout-authority-governance/README.md) | A UI has competing owners for layout dimensions, position, insets, viewport state, scroll roots, or fixed surfaces. |

## Boundaries

- Workflow steps belong in `workflow/`.
- Responsive failure classification, authority, severity, and closure checklist belong in [`../governance/responsive-ui/`](../governance/responsive-ui/README.md).
- Layout source-of-truth and viewport authority rules belong in [`layout-authority-governance/`](layout-authority-governance/README.md).
- Shared render context vocabulary belongs in [`../render-contexts/`](../render-contexts/README.md).
- Project-specific screenshots, routes, hostnames, and one-off evidence stay in the project repository.
- Framework docs may be cited as evidence, but this layer should store the generalized engineering primitive, not framework-specific class names.
