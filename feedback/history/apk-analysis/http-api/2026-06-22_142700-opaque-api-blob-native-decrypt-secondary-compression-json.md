> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - API opaque blob field may be native decrypt + secondary compression before JSON

Status: candidate

#### One-line Summary

HTTP JSON 回應某欄位（如 `play_info`）為 opaque Base64 字串，而 Gson/客戶端內已見結構化 JSON 時，應假設 **native decrypt util（如 `decryptStr(blob, mode)`）+ 二級解碼（Base64 → zlib/gzip → UTF-8 JSON）**；host 端用 Frida RPC 呼叫同一 util 做 mode sweep，而非直接期待 `PlayURL` 在 raw HTTP body。

#### Human Explanation

短劇/視頻 API 常在 wire 上加密播放清單，僅在 app 內解密後才進 Gson。特徵：raw response 有 `play_info` 長字串；Frida `Gson.fromJson` hook 卻見 `[{PlayURL:…}]` 陣列。RE 路徑：對 blob 用與 app 相同的 native decrypt（常需試 mode 整數）；解密結果若仍非 JSON，再試 Base64 decode、zlib decompress。離線 SDK 在 decrypt RE 完成前，用 Frida CLI relay decrypt 與 sign/requestTime 並列。

#### Trigger

- API `code=0` 但播放 URL 欄位為長 Base64/opaque string
- Gson/converter hook 顯示解密後結構化 JSON
- `decryptStr` / `decryptByte` 出現在 DEX native util

#### Evidence

- Tool: hybrid HTTP client + Frida `decryptstr` RPC + Python zlib
- Sanitized excerpt: mode sweep → inner Base64 → zlib → JSON variant array with `PlayURL`
- Evidence path: `<PROJECT_ROOT>/scripts/sign/play_info.py`

#### Generalized Lesson

```text
Opaque response blob triage:
  1. Compare raw HTTP JSON vs in-app Gson plaintext
  2. Frida RPC: native decryptStr(blob, mode) for mode 0..N
  3. On printable inner: try base64decode → zlib/gzip → json.loads
  4. Document mode + chain in project; Ai-skill 不寫 key/mode 真值
  5. Downloader: relay decrypt until RE; CDN m3u8 often plain after step 3
```

#### Agent Action

1. Project `play_info.py` 封裝 decrypt chain；downloader 依賴 Frida interim。
2. Ai-skill 寫判斷樹，不寫具體 mode 數字。

#### Goal / Action / Validation

- Goal: downloader 不被 raw JSON 形狀誤導。
- Action: mode sweep + secondary decompression probe。
- Validation: resolved m3u8 URL + ffmpeg 成功。

#### Applies When

- Media play URL hidden behind encrypted API field
- Same APK util class handles sign and decrypt

#### Does Not Apply When

- PlayURL already plaintext in HTTP JSON
- DRM only at CDN `#EXT-X-KEY` layer (separate issue)

#### Validation

- m3u8 extracted and media download succeeds

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §response blob decrypt

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 已依 sanitization / reusable-guidance-boundary 自查
