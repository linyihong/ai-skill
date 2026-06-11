# Responsive UI Governance Checklist

Use this checklist before claiming a responsive UI issue is resolved or a responsive design foundation is ready.

## Contract

- Supported render contexts are named with shared vocabulary, not device-only labels.
- Intended unsupported contexts are explicit.
- Fixed, sticky, modal, sheet, and scroll-root contracts are named when present.
- Intended horizontal scrollers are scoped.

## Implementation Review

- Containers are fluid first and capped by a shared max-width or container contract.
- Breakpoints encode layout constraints rather than device SKUs.
- Viewport height behavior accounts for dynamic browser chrome when relevant.
- Fixed or sticky surfaces align to the same shell as content.
- Safe-area insets are applied where interactive surfaces need them.

## Validation

- Evidence includes at least one wide and one constrained render context.
- `narrow_mobile`, `safe_area`, `landscape`, or `dynamic_resize` are included when the defect or contract depends on them.
- Document and app-shell overflow metrics distinguish intended scrollers from unintended overflow.
- Primary actions and navigation remain reachable.

## Closure

- Failure class, authority, severity, and release posture are recorded.
- Any deferred context has an owner and reason.
- Project-specific evidence remains in the project repository, not reusable governance docs.
