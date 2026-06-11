# Responsive Failure Taxonomy

This taxonomy names what failed. It does not decide who has authority or whether release is blocked; see [`authority-mapping.md`](authority-mapping.md) and [`severity-policy.md`](severity-policy.md).

## Failure Classes

| Failure | Description |
| --- | --- |
| `horizontal_overflow` | Document or shell becomes wider than the allowed render context, excluding declared horizontal scrollers. |
| `vertical_clipping` | Content or controls become unreachable because height, scroll-root, or fixed-surface rules clip them. |
| `fixed_surface_drift` | Fixed or sticky UI does not align to the same container contract as content. |
| `safe_area_overlap` | System UI, notch, or home indicator overlaps interactive content. |
| `dynamic_resize_staleness` | Layout remains based on an old viewport after resize, reload, browser chrome change, or orientation change. |
| `density_collapse` | Text, cards, controls, or spacing no longer fit the declared context while preserving readable hierarchy. |
| `unscoped_framework_mix` | A framework pattern or class set is introduced without matching project ownership, reset, spacing, or precedence rules. |

## Classification Rule

Classify the failure by the user-observable broken invariant. Avoid classifying by the tool that found it, the device model, or the framework used to fix it.
