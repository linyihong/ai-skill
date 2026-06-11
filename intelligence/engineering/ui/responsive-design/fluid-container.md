# Fluid Container Pattern

Principle: layout containers should be fluid first and capped later.

## Transferable Shape

```css
.container {
  width: 100%;
  max-width: var(--container-max-width);
  margin-inline: auto;
  padding-inline: var(--container-padding);
}
```

This shape is more resilient than starting from a fixed width and trying to patch overflow later.

## Rule

- Base container width should respond to the available inline size.
- `max-width` is a cap, not the starting point.
- Padding should shrink or clamp for narrow contexts.
- Nested containers should not independently redefine the page width contract unless they are a deliberate sub-layout.

## Smells

- Multiple components each define their own shell width.
- A fixed footer, modal, or scroll root uses a different max width from the page shell.
- Horizontal overflow is fixed by hiding the document overflow instead of correcting the oversized container.
