# Fixed Surface Contract

Principle: fixed, sticky, floating, and docked surfaces must share the same authority model as the content they control.

## Contract Fields

```yaml
fixed_surface_contract:
  aligns_to: content_shell
  width_source: css_container_contract
  position_source: css_layout_engine
  inset_source: safe_area_contract
  overlap_policy: reserved_scroll_space
```

## Rule

- A fixed surface should not independently define shell width if the content has a shell contract.
- Bottom navigation and floating actions must reserve scroll space or prove they do not hide final content.
- Fixed surfaces near system UI must account for safe-area constraints.
- The validation target includes both surface bounding box and content reachability.

## Smells

- Fixed footer uses viewport width while content is centered and capped.
- A sticky header or bottom bar drifts after reload, zoom, or resize.
- Last content item, submit button, or primary navigation is hidden behind a fixed surface.
