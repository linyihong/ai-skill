# Magic Bytes Reference Heuristic（Magic Bytes 快速參考）

**Status**: `candidate-intelligence`
**Source**: 從 `intelligence/engineering/analytical-reasoning/signals/media-type-detection.md` 提取的跨領域通用部分

## 問題

在分析未知檔案時，如何快速判斷其實際容器類型？檔案 extension 不可信時，如何可靠地辨識媒體格式？

## 原則

- 檔案 extension 不可信；extension 與實際 bytes 可能不一致
- Magic bytes（file header）是判斷容器類型的最可靠方式
- `ffprobe` 可提供 container metadata、codec、resolution、duration 等詳細資訊
- 靜態圖片、動畫、音訊、影片需要不同的處理流程

## Magic Bytes 快速參考

| 容器類型 | Magic Bytes（hex） | 檔案 extension（常見） |
|---------|-------------------|---------------------|
| MP4 | `00 00 00 18 66 74 79 70` | `.mp4` |
| WebP | `52 49 46 46 xx xx xx xx 57 45 42 50` | `.webp` |
| GIF | `47 49 46 38` | `.gif` |
| PNG | `89 50 4E 47` | `.png` |
| JPEG | `FF D8 FF` | `.jpg`、`.jpeg` |
| MP3 (ID3) | `49 44 33` | `.mp3` |
| FLAC | `66 4C 61 43` | `.flac` |
| TS (MPEG-TS) | `47` (sync byte) | `.ts` |
| WebM | `1A 45 DF A3` | `.webm` |

## 靜態 vs 動畫判斷

| 類型 | 檢查方式 | 判斷依據 |
|------|---------|---------|
| 靜態圖片 | `ffprobe -show_streams` | 無 duration 或 duration=0，單一 frame |
| 動畫 WebP | `ffprobe -show_streams` | VP8/VP8/VP8X codec，有 duration |
| 動畫 GIF | `ffprobe -show_streams` | GIF codec，有 duration，多 frame |
| 音訊 | `ffprobe -show_streams` | 只有 audio stream，無 video stream |
| 影片 | `ffprobe -show_streams` | 有 video stream（可能同時有 audio stream） |

## Container Probe 指令

```bash
# 快速判斷容器類型
file media_file

# 詳細 container metadata
ffprobe -v quiet -print_format json -show_format media_file

# 詳細 stream 資訊
ffprobe -v quiet -print_format json -show_streams media_file

# 檢查是否可播放
ffmpeg -i media_file -f null -
```

## 判斷流程

```
取得媒體檔案
    ├── 執行 `file` command
    │   ├── 結果明確 → 記錄容器類型
    │   └── 結果不明確或 "data" → 檢查 magic bytes
    ├── 執行 `ffprobe -show_streams`
    │   ├── 有 video stream → 影片
    │   ├── 只有 audio stream → 音訊
    │   ├── 有 duration + 單一 video stream → 動畫
    │   └── 無 duration / 單一 frame → 靜態圖片
    └── 記錄最終判斷與驗證方式
```

## 相關 atoms

- `intelligence/engineering/analytical-reasoning/signals/media-type-detection.md`（領域特定版本，包含 HLS/TS 等 APK 分析專用信號）

## Token 影響

低。此 atom 在任何需要判斷檔案類型的任務中 lazy-load，約 150-200 tokens。
