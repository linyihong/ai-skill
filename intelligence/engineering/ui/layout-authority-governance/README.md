# Layout Authority Governance

UI layout is a governed system. Viewport, layout engine, scroll root, fixed surfaces, safe area, and validation evidence must share a consistent authority model.

This module captures layout authority failures that are deeper than responsive breakpoints. Responsive UI is one application of this governance model; fixed navigation, embedded webviews, full-height shells, drawers, modals, and browser state restoration are others.

For layout, the reusable chain is:

```text
layout authority
→ layout contract
→ layout state
→ validation
→ evidence
```

This module owns the layout authority, contract, and state portions. Validation and evidence remain cross-domain concerns referenced from validation reasoning and software-delivery workflow surfaces.

## Surfaces

| File | Use When |
| --- | --- |
| [`layout-source-of-truth.md`](layout-source-of-truth.md) | A UI has more than one possible owner for dimensions, position, or insets. |
| [`css-layout-ownership.md`](css-layout-ownership.md) | Deciding whether CSS or JavaScript should own layout calculation. |
| [`viewport-measurement-authority.md`](viewport-measurement-authority.md) | JavaScript reads viewport or element measurements and may cache or replay them. |
| [`fixed-surface-contract.md`](fixed-surface-contract.md) | Fixed, sticky, floating, or docked surfaces must align with the content shell. |
| [`scroll-root-inset-contract.md`](scroll-root-inset-contract.md) | Scroll roots must reserve space for fixed surfaces and safe areas. |
| [`stale-layout-measurement.md`](stale-layout-measurement.md) | Layout state can go stale after resize, reload, restore, or browser chrome changes. |

## Boundary

- Traditional responsive patterns such as mobile-first, breakpoints, and fluid containers remain in [`../responsive-design/`](../responsive-design/README.md).
- Shared context names live in [`../../render-contexts/`](../../render-contexts/README.md).
- Responsive severity and failure authority live in [`../../governance/responsive-ui/`](../../governance/responsive-ui/README.md).
- Transition validation is cross-domain and lives in [`../../execution/validation-reasoning/state-transition-validation.md`](../../execution/validation-reasoning/state-transition-validation.md); this module only references it when layout state can stale across resize, reload, rotate, restore, or resume.
- Evidence acquisition and validation execution live in workflow validation surfaces.
