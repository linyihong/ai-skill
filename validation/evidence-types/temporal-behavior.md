# evidence_type: `temporal_behavior`

## Proves

An observable state transition occurs within a time boundary: preview cutoff, debounce window, poll interval firing, timed overlay appearance, buffer-stall recovery.

Answers **what transitioned**, not **why the class of bug exists**.

## Non-goals

- Failure class label as gate token (rejected: `timing_gate`)
- Infinite wait without bounded timeout in test

## Supported collection_methods

- `browser_observation`
- `runtime_trace`

## Supported artifact_shapes

- `poll_log`
- `timestamped_state_log`
- `bounded_wait_assertion`

## Proxy traps

- Overlay appears at 60s when contract says 15s → temporal evidence fails even if `user_visible` eventually passes
- Poll log without initial pathname / video time columns → incomplete artifact

## Example claim

`preview_limit_enforced`

## Note on naming

`timing_gate` was rejected as an evidence type because it described a failure mode, not a provable outcome. Use `temporal_behavior` with an explicit claim id instead.
