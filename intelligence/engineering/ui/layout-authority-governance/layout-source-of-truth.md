# Layout Source Of Truth

Principle: each layout decision should have one authoritative source for dimensions, position, and insets.

## Authority Contract

```yaml
layout_authority:
  layout_source_of_truth: css
  allowed_sources:
    - browser_layout_engine
    - css_custom_properties
    - container_query
  forbidden_sources:
    - stale_javascript_measurement
    - persisted_viewport_measurement
    - duplicated_component_constant
```

The exact values are project-local, but the authority model must be explicit when multiple systems can affect the same layout.

## Rule

- CSS should own stable layout relationships when the browser layout engine can express them.
- JavaScript may observe layout for behavior, but should not become a second layout source unless the contract names why.
- Shared dimensions should come from one contract surface, such as design tokens, CSS custom properties, or container contracts.
- Duplicated width, inset, or height constants are authority drift.

## Smells

- CSS centers the shell while JavaScript separately computes fixed navigation position.
- A component stores viewport dimensions and does not invalidate them on relevant transitions.
- Multiple files define the same shell max width, tab inset, or safe-area rule.
