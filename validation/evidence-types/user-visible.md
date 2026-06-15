# evidence_type: `user_visible`

## Proves

A user-visible UI state exists: overlay, mask, CTA, modal, disabled control, or layout region with non-zero visible box.

## Non-goals

- Source code contains class name
- API returned 200 without DOM observation
- Screenshot without linked **claim**

## Supported collection_methods

- `browser_observation`
- `human_review` (advisory only)

## Supported artifact_shapes

- `screenshot`
- `dom_assertion_log`
- `accessibility_tree_excerpt`

## Proxy traps

- Headless DOM height > 0 but opacity 0 / off-screen → failed `user_visible`
- Screenshot with no `claim` id → artifact without trace chain

## Example claim

`preview_overlay_shown`
