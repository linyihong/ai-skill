# Overflow Management

Principle: overflow is diagnostic evidence. Do not hide it before finding the source layout failure.

## Rule

- Investigate horizontal overflow by locating the oversized element.
- Fix the source: container width, min-width, grid/flex shrink behavior, long text wrapping, fixed element alignment, or media sizing.
- Use `overflow: hidden` only when clipping is the intended component behavior, not as a page-level bandage.
- Treat horizontal document overflow as a release blocker candidate for constrained render contexts.

## Common Causes

- Missing `min-width: 0` in grid or flex children.
- Long unwrapped tokens, URLs, order IDs, or addresses.
- Fixed-width cards inside a narrow container.
- Fixed or absolute positioned elements using viewport width while content uses a centered shell.
- Media without `max-width: 100%`.

## Validation Signal

Useful responsive validation includes both DOM metrics and visual evidence:

```text
document.documentElement.scrollWidth <= window.innerWidth
```

Pair this with checking key shell, fixed footer, and primary content bounding boxes.
