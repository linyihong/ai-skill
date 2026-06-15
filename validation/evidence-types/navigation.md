# evidence_type: `navigation`

## Proves

Navigation outcome is correct: resolved href, history back stack, route pathname, or basePath-safe client navigation.

## Non-goals

- Link exists in source
- Server route file exists without browser resolution

## Supported collection_methods

- `browser_observation`
- `runtime_trace`

## Supported artifact_shapes

- `navigation_log`
- `href_resolution_log`
- `history_stack_log`

## Proxy traps

- Rendered href double-applies `NEXT_PUBLIC_BASE_PATH` → navigation defect
- `router.push` string assert without resolved URL → `source_contract` only

## Example claim

`subscribe_href_resolves`
