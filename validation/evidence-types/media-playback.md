# evidence_type: `media_playback`

## Proves

Media element or stream behavior: play/pause, currentTime progression, HLS load, stall recovery, mute/rate application.

## Non-goals

- `videoUrl` present in SSR props only
- CDN reachable from curl without browser media pipeline

## Supported collection_methods

- `browser_observation`
- `runtime_trace`

## Supported artifact_shapes

- `media_state_log`
- `hls_event_log`
- `screenshot` (frame-visible proof only)

## Proxy traps

- `video.play()` resolved but `currentTime` frozen (buffer stall) → needs `temporal_behavior` companion claim
- Mock player unit test without real `<video>` → not browser validation

## Example claim

`preview_playback_reaches_boundary`
