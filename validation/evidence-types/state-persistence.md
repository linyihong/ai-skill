# evidence_type: `state_persistence`

## Proves

Client or server state survives navigation, remount, or session boundary as specified (sessionStorage, cookie, DB readback).

## Non-goals

- In-memory React state within same mount
- One-shot write without readback

## Supported collection_methods

- `browser_observation`
- `contract_readback`
- `runtime_trace`

## Supported artifact_shapes

- `storage_readback_log`
- `cookie_trace`
- `db_readback_row`

## Proxy traps

- SSR prop present once but lost after client navigation → missing persistence evidence
- DB row from test setup without user action path → weak journey evidence

## Example claim

`playback_prefs_persist_after_episode_swipe`
