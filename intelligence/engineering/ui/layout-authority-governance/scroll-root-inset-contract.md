# Scroll Root Inset Contract

Principle: scroll roots must know which fixed surfaces and safe areas reduce usable content space.

## Contract Fields

```yaml
scroll_root_contract:
  owns_vertical_scroll: true
  bottom_inset_sources:
    - fixed_bottom_navigation
    - safe_area
  final_content_reachable: required
  intended_scroll_lock:
    allowed: false
    reason: null
```

## Rule

- If a page has fixed bottom UI, the scroll root must reserve bottom space.
- Disabling vertical scroll is a layout authority decision, not a cosmetic flag.
- Safe-area padding must protect interactive content, not only background fill.
- Final content reachability is part of validation evidence.

## Smells

- `overflow: hidden` on a full-page root hides a clipping defect.
- A scroll view disables vertical scroll while content height can exceed the visible viewport.
- Bottom sentinel height is unrelated to fixed surface height and safe-area inset.
