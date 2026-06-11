# UI Engineering Intelligence

Reusable UI engineering judgment lives here. Keep this layer framework-neutral: record principles, trade-offs, and failure modes that survive Bootstrap, Tailwind, MUI, native CSS, or future UI systems.

## Topics

| Topic | Use When |
| --- | --- |
| [`responsive-design/`](responsive-design/README.md) | A UI must remain correct as viewport width, height, safe area, orientation, or container size changes. |

## Boundaries

- Workflow steps belong in `workflow/`.
- Governance classification, authority, severity, and closure checklist belong in [`../governance/responsive-ui/`](../governance/responsive-ui/README.md).
- Shared render context vocabulary belongs in [`../render-contexts/`](../render-contexts/README.md).
- Project-specific screenshots, routes, hostnames, and one-off evidence stay in the project repository.
- Framework docs may be cited as evidence, but this layer should store the generalized engineering primitive, not framework-specific class names.
