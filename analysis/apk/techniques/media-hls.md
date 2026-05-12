# Media / HLS Techniques

> **Deprecated** — This technique has been fully decomposed. See the new locations below.
> See:
> - `analysis/apk/workflows/media-hls-analysis-flow.md`（HOW TO DO 操作流程）
> - `intelligence/engineering/apk-analysis/signals/media-type-detection.md`（媒體類型偵測信號）

Use this category when the target flow involves video, audio, images, playlists, media keys, segments, or container validation.

> **相容性規則**：`skills/apk-analysis/techniques/media-hls/` 仍為 active skill entrypoint。本文件為 reference target，兩者應保持同步。

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

## Media Chain

| Layer | Example | Document |
| --- | --- | --- |
| Detail/control API | title, cover, source path | API path, required auth, source field meaning. |
| Playlist | HLS `.m3u8` | key URI, segment count, duration, base URL, expiration behavior. |
| Key | AES key endpoint or key file | key length, retrieval conditions, auth requirement, IV handling. |
| Segments | `.ts` / chunk / signed URL | URL lifetime, query meaning, download order, sequence gaps. |
| Final media | mp4/mp3/image/webp/gif | decryption, decode, remux, `ffprobe`/header validation. |

Use magic bytes, container probes, or frame counts to distinguish static image, animated WebP/GIF, audio, and video containers.

## Related Lessons

- `skills/apk-analysis/feedback_history/media-hls/2026-04-30_120007-媒體播放鏈要分控制面-金鑰與資料面.md`
