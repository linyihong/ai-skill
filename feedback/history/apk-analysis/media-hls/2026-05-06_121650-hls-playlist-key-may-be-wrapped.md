> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/media-hls-analysis-flow.md`](../../../../analysis/apk/workflows/media-hls-analysis-flow.md)

### 2026-05-06 - HLS Playlist Key May Be Wrapped

Status: promoted

#### One-line Summary

HLS playlist 的 `#EXT-X-KEY` URI 不一定直接提供可解 segment 的 AES key；先檢查控制 API 是否另有 wrapped key，並 hook App 端 unwrap 函式驗證。

#### Human Explanation

HLS 分析常會照標準流程下載 `.m3u8`、讀 `#EXT-X-KEY`、下載 key URI，再用 AES-128-CBC + IV 解 `.ts` segment。若解密後沒有 MPEG-TS sync byte `0x47`，不要立刻判定 segment 不是 HLS 或 IV 算錯；有些 App 會在控制 API 另外下發 `encrypted_key` / `decrypt_key` / `customKey`，播放器或 helper 會先把這個值 unwrap 成真正 segment key，playlist key URI 可能只是 placeholder、fallback 或被 App 改寫。

#### Trigger

授權 APK 媒體分析時，詳情 API 的 content block 內有 `.m3u8` URL 與 `encrypted_key`。playlist 宣告 `#EXT-X-KEY:METHOD=AES-128`，key URI 可下載，但用該 key + playlist IV 解第一段 segment 沒有 MPEG-TS sync bytes。Hook App 端 key unwrap helper 後，取得另一個 16-byte key；用它配合 playlist IV 可解出有效 MPEG-TS。

#### Evidence

- Tool: schema-only response hook, private media URL capture, Frida hook on App key unwrap helper, offline AES-CBC segment probe, `ffprobe`.
- Sanitized excerpt: control API block contains `url=.m3u8`, `encrypted_key=<44-char string>`, `video_duration=<number>`；playlist key response length looked valid but failed TS sync; App unwrap helper returned a base64 string whose decoded 16 bytes decrypted segment 0 to TS.
- Evidence path: project-private `capture/` logs and sample segment only; reusable lesson contains no target host, auth query, raw key, raw URL, segment bytes, or user media.

#### Generalized Lesson

對 HLS/AES-128 影片，不要把 playlist key URI 視為唯一 key source。控制 API、model fields、player wrappers、header helpers、`getDecryptionKey`、`decryptBase64...`、`FFAES` 類函式都可能參與真正 key derivation。離線驗證必須以 segment 解密結果為準：

1. 下載 playlist/key/segment。
2. 用 playlist key URI 的 bytes 嘗試解密。
3. 若沒有 `0x47` TS sync 或 container probe 失敗，回到控制 API 找 `encrypted_key` / `decrypt_key` / `customKey`。
4. Hook App 端 unwrap/decrypt helper，取得去敏長度/hash與本地私有 raw key。
5. 用 unwrap 後 key + playlist IV 解 segment，再以 sync byte、packet interval、`ffprobe` 驗證。

#### Agent Action

下次分析 HLS：

1. 先分控制面欄位、playlist、key URI、segments、final container。
2. 若 playlist key 解不出 TS，不要只反覆調整 IV；先搜尋並 hook `encrypted_key` / `decrypt_key` / `customKey` 相關 unwrap 函式。
3. raw URL、auth query、key bytes 與 media bytes 只放 project private capture；公開/skill 文件只寫 shape、length、hash、驗證結果。
4. 成功標準是 decrypted segment 能通過 magic/container probe，而不是 hook 到某個 key function。

#### Goal / Action / Validation

- Goal: 避免把 wrapped-key HLS 誤判為無法下載、IV 錯誤或非標準 playlist。
- Action: 在 playlist key 失敗時，回查控制 API key material，hook App unwrap helper，再做離線 AES-CBC 驗證。
- Validation or reference source: decrypted segment starts with MPEG-TS sync bytes at 188-byte intervals and `ffprobe` identifies expected audio/video streams.

#### Applies When

- HLS playlist 宣告 `METHOD=AES-128`，但 key URI 解密失敗。
- 控制 API 或 media model 有 `encrypted_key`、`decrypt_key`、`customKey`、`videoId`、`videoUrl` 等欄位。
- App 內有 key helper、media player wrapper、Dart AOT `FFAES` / AES / base64 decrypt 函式。

#### Does Not Apply When

- Playlist key URI 解密後已直接得到有效 TS/fMP4 container。
- Playlist 使用 SAMPLE-AES、DRM、FairPlay/Widevine，或 segment 不是 AES-128-CBC HLS。
- 分析範圍不允許下載 segments 或 hook key unwrap function。

#### Validation

- 至少一個 segment 用 unwrap key + playlist IV 解密後，`0x47` sync byte 出現在 188-byte packet interval。
- `ffprobe` 或等效 container probe 能辨識 video/audio stream。
- Public docs 不包含 raw key、raw media URL、token、auth query 或 media bytes。

#### Promotion Target

- `WORKFLOW.md` / 媒體與 HLS 分析章節。
- Project media docs and private capture notes.

#### Required Linked Updates

- 已同步更新 `WORKFLOW.md` 的 HLS key unwrap 注意。
- 已更新 `feedback_history/README.md` 與 `feedback_history/media-hls/README.md` 索引。
