# Stale Layout Measurement

Principle: measured layout state becomes stale unless every relevant transition invalidates it.

## Pattern

```text
measure viewport or element
↓
cache derived layout value
↓
viewport / shell / route / browser state changes
↓
cached value remains authoritative
↓
layout drifts from current reality
```

## Rule

- Do not promote a measurement to layout authority without an invalidation contract.
- Prefer declarative CSS for values that should automatically follow current layout reality.
- Validate transition sequences, not only fresh steady state.

## Common Variants

- stale viewport measurement;
- stale fixed surface position;
- stale scroll-root inset;
- stale container width;
- stale restored route shell dimensions.

See also [`../../anti-patterns/stale-derived-state.md`](../../anti-patterns/stale-derived-state.md) for the broader engineering anti-pattern.
