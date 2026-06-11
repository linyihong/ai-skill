# Stale Derived State

Anti-pattern: a derived observation remains authoritative after the source reality has changed.

## Pattern

```text
source reality
↓
derive state
↓
cache / persist / replay derived value
↓
source reality changes
↓
derived value still drives behavior
↓
system acts on stale reality
```

## Common Variants

- `stale_layout_measurement`: viewport, container, or fixed-surface measurements drive layout after resize, reload, restore, or browser chrome changes.
- `stale_route_state`: navigation state is restored or reused after the route context changed.
- `stale_permission_state`: authorization or entitlement projection remains active after upstream identity or permission changed.
- `stale_feature_flag`: a cached flag value controls behavior after rollout state changed.
- `stale_session_projection`: session-derived UI or data state remains visible after session refresh, logout, or identity switch.

## Rule

- Derived state needs an invalidation contract.
- The invalidation contract should name source changes, transition events, and validation evidence.
- If no invalidation contract exists, prefer reading current source reality or letting the platform compute the derived value.
- Do not treat a successful initial render as proof that derived state survives transitions.

## Viewport Measurement Drift

Viewport measurement drift is a layout-specific stale derived state failure:

```text
current viewport / layout engine reality
↓
JavaScript measures viewport or element
↓
component caches derived width / position / inset
↓
resize, reload, browser chrome, route restore, or safe-area context changes
↓
cached measurement continues to position UI
```

The fix is not necessarily "add more resize listeners." First decide whether layout authority should belong to CSS, the browser layout engine, a container contract, or JavaScript measurement. See [`../ui/layout-authority-governance/`](../ui/layout-authority-governance/README.md).

## Validation

Validation must include at least one transition that can invalidate the derived state. For layout, that often means `fresh_load -> resize -> reload`, orientation change, route restore, or app resume.
