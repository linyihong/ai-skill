# Responsive UI Contract

A responsive UI contract names the render contexts a UI must support and the layout invariants that must hold in each context.

## Required Fields

```yaml
responsive_ui_contract:
  supported_render_contexts:
    - narrow_mobile
    - mobile
    - tablet
    - desktop
  layout_invariants:
    - no_unintended_horizontal_overflow
    - primary_actions_visible_and_reachable
    - fixed_surfaces_align_to_content_shell
    - scroll_root_preserves_access_to_final_content
  authority_model:
    layout_source_of_truth: css
    fixed_surface_position_source: css_layout_engine
    viewport_measurement_source: observation_only
    forbidden_sources:
      - stale_viewport_measurement
      - persisted_layout_measurement
  exception_policy:
    intended_horizontal_scrollers:
      - component_name_or_surface
    unsupported_contexts:
      - context_with_reason
```

## Contract Rules

- Render contexts describe constraints, not device names.
- Every fixed, sticky, modal, sheet, or full-height scroll root must name its alignment and inset contract.
- Every responsive contract must name layout authority when JavaScript, CSS, browser APIs, or persisted state can affect the same dimensions.
- Intended overflow must be explicit and scoped to a component.
- A desktop-only contract is valid only when the project states it explicitly.

## Non-Goals

- Do not define a global breakpoint scale here.
- Do not require a specific framework such as Bootstrap or Tailwind.
- Do not store project screenshots, routes, hostnames, or one-off live evidence in reusable governance.
