# Media / HLS Analysis Flow（媒體 / HLS 分析操作流程）

`analysis/apk/workflows/media-hls-analysis-flow.md` 是從 `skills/apk-analysis/techniques/media-hls/`（已刪除）拆解出的 **HOW TO DO** 操作流程。決策智慧（媒體類型判斷信號）請見 `intelligence/engineering/analytical-reasoning/signals/media-type-detection.md`。

> **Intelligence Extracted**
> See:
> - `intelligence/engineering/analytical-reasoning/signals/media-type-detection.md`

## 前置準備

### 必要條件

- API responses 中包含 media source paths、HLS playlist URLs、key URLs、segment lists、signed media URLs 或 image/audio/video payloads
- 目標是 reconstruct playable media 或 verify media container type

### 工具

```bash
# 媒體驗證
ffprobe -v quiet -print_format json -show_format -show_streams media_file
ffmpeg -i media_file -f null -

# 容器分析
file media_file           # magic bytes 判斷
hexdump -C media_file | head -20  # 檢查 header

# 下載工具
curl -O "https://example.com/media/file.ts"
wget "https://example.com/media/playlist.m3u8"
```

## 步驟 1：分離控制面與資料面

先區分 control-plane APIs 與 data-plane media downloads：

| 類型 | 內容 | 分析方式 |
|------|------|---------|
| Control-plane API | title, cover, source path, auth | API documentation flow |
| Data-plane media | playlist, key, segments, final media | Media chain analysis |

## 步驟 2：分析 Detail/Control API

記錄 media 相關 API 的：
- API path、required auth、source field meaning
- 哪些欄位指向 media source（如 `video_url`、`cover_url`、`manifest_url`）

## 步驟 3：分析 Playlist（HLS `.m3u8`）

```bash
# 下載 playlist
curl -O "https://example.com/media/playlist.m3u8"

# 檢查 playlist 內容
cat playlist.m3u8
```

記錄：
- Key URI（如有加密）
- Segment count 與 duration
- Base URL 與 expiration behavior
- 備用/alternative playlist

## 步驟 4：分析 Key（如適用）

```bash
# 下載 key file
curl -O "https://example.com/media/key.bin"

# 檢查 key length
hexdump -C key.bin
```

記錄：
- Key length（通常 16 bytes for AES-128）
- Retrieval conditions（auth required?）
- IV handling（來自 playlist 或固定值）

## 步驟 5：下載與分析 Segments

```bash
# 下載 segment
curl -O "https://example.com/media/segment001.ts"

# 檢查 segment 類型
file segment001.ts
ffprobe segment001.ts
```

記錄：
- URL lifetime 與 query meaning
- Download order 與 sequence gaps
- Segment 格式（`.ts`、chunk、signed URL）

## 步驟 6：合併與解碼 Final Media

```bash
# 合併 segments（如有加密需先解密）
ffmpeg -i "playlist.m3u8" -c copy output.mp4

# 或手動合併 TS segments
cat segment*.ts > combined.ts
ffmpeg -i combined.ts -c copy output.mp4

# 驗證最終媒體
ffprobe output.mp4
```

## 步驟 7：容器驗證

不要假設 extension 等於格式。使用 magic bytes 和 container metadata 驗證：

```bash
# 方法 1：file command（magic bytes）
file output.mp4
# 輸出範例：output.mp4: ISO Media, MP4 v2, Base Media

# 方法 2：ffprobe
ffprobe -v quiet -print_format json -show_format output.mp4

# 方法 3：hexdump header
hexdump -C output.mp4 | head -5
# MP4: 00 00 00 18 66 74 79 70
# WebP: 52 49 46 46 xx xx xx xx 57 45 42 50
# GIF: 47 49 46 38
# PNG: 89 50 4E 47
```

## 成功產出格式

```markdown
## Media Chain

### Control API
- Endpoint: GET /api/v1/media/detail
- Response fields: { title, cover_url, video_url, manifest_url }

### Playlist
- URL: https://cdn.example.com/hls/playlist.m3u8
- Segments: 12
- Duration: 120s
- Encryption: AES-128 (key URI: https://cdn.example.com/hls/key.bin)

### Key
- Length: 16 bytes
- Auth: none
- IV: from playlist (IV=0x...)

### Segments
- Format: .ts
- Count: 12
- URL pattern: https://cdn.example.com/hls/segment{seq}.ts?token=...

### Final Media
- Format: MP4 (verified by ffprobe)
- Resolution: 1920x1080
- Codec: H.264 + AAC
- Playable: yes
```

## 注意事項

- 分離 control-plane APIs 與 data-plane media downloads
- 不要假設格式來自 extension；驗證 magic bytes 和 container metadata
- HLS 未完成直到 playlist、key/IV（如需要）、segments、decryption、merge/remux、playback/container validation 都處理完
- Redact URLs 可直接 replay private 或 paid content 的 URL
