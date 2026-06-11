# State Transition Validation

Principle: many correctness failures occur during transitions, not in a steady state.

## Why

A fresh-load pass proves only that one initial state can render or execute. It does not prove the system survives changes in viewport, route, session, permission, feature flag, external state, or app lifecycle.

## Transition Families

| Transition | Use When |
| --- | --- |
| `fresh_load` | Establish baseline state. |
| `resize` | Viewport, container, split-screen, or browser chrome can change layout. |
| `reload` | Reload after state or viewport change can expose stale initialization assumptions. |
| `rotate` | Orientation or aspect ratio changes can alter layout and reachability. |
| `route_restore` | Browser history, deep-link fallback, or app navigation restores prior route state. |
| `session_restore` | Identity or entitlement projection is restored from session state. |
| `background_resume` | Native wrapper, mobile browser, or app lifecycle resumes from paused state. |
| `source_change` | Permission, feature flag, remote config, or external source of truth changes. |

## Rule

- If the defect report names a transition, validation must include that transition.
- If implementation introduces cached or derived state, validation must include the transition that can stale it.
- If a claim depends on layout authority, include layout transition evidence such as bounding boxes, overflow metrics, and final content reachability.
- A steady-state screenshot or API success is not enough for transition-only claims.

## Evidence Shape

```yaml
state_transition_validation:
  sequence:
    - fresh_load
    - resize
    - reload
  authority_expected_to_update:
    - css_layout_engine
    - derived_state_cache
  evidence:
    before:
      - source_reality_snapshot
      - user_visible_state
    after:
      - source_reality_snapshot
      - user_visible_state
      - authority_readback
```

## Relationship

- Use [`stale-derived-state.md`](../../anti-patterns/stale-derived-state.md) when a derived observation remains authoritative after source reality changes.
- Use [`evidence-chain-validation.md`](evidence-chain-validation.md) when the transition must propagate through multiple system layers.
- Use [`state-visibility-gap.md`](state-visibility-gap.md) when observed state and true system state can diverge.
