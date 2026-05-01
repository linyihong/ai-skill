# Media / HLS Techniques

Use this category when the target flow involves video, audio, images, playlists, media keys, segments, or container validation.

## When To Use

- API responses contain media source paths, HLS playlist URLs, key URLs, segment lists, signed media URLs, or image/audio/video payloads.
- The goal is to reconstruct playable media or verify media container type.
- File extension and actual bytes may disagree.

## Core Guidance

- Separate control-plane APIs from data-plane media downloads.
- Document playlist, key, segment, and final media as separate layers.
- Do not assume format from extension; verify magic bytes and container metadata.
- HLS is not complete until playlist, key/IV if needed, segments, decryption, merge/remux, and playback/container validation are handled.
- Redact URLs that can directly replay private or paid content.

## Related Lessons

- `../../feedback_history/2026-04-30_120007-媒體播放鏈要分控制面-金鑰與資料面.md`
