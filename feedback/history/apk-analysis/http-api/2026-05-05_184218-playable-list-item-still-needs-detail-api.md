> 遵守 [共用規則索引](../../../../enforcement/README.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/http-api-documentation-flow.md`](../../../../analysis/apk/workflows/http-api-documentation-flow.md)

### 2026-05-05 - Playable List Item Still Needs Detail API

Status: validated

#### One-line Summary

列表 item 已經帶可播放 URL 時，仍要用 UI 點擊行為比對是否還有詳情 API，不要直接判斷「沒有詳情頁」。

#### Human Explanation

很多影音 App 會在列表 item 內預先放入 `sourceURL`、`playUrl`、`videoUrl` 或圖片資源，讓列表可以快速起播或預覽。但實際 UI 點擊 item 進入詳情頁時，常會再打詳情 API 取得描述、互動統計、作者、關聯人物、推薦、評論入口或權限狀態。若工具或 SDK 因為列表已可播放就直接播放，會漏掉詳情頁 API 與欄位。

#### Trigger

遇到列表 item 已有播放來源，測試工具點擊後只播放、不顯示詳情資料；使用者或 UI 行為顯示應該存在詳情頁。

#### Evidence

- Tool: authorized local API replay through a project runner.
- Sanitized excerpt: list item exposed an m3u8-like `sourceURL`; a separate `GET /<video-detail-path>?id=<item-id>` returned richer fields such as description, publisher, interaction counts, status, gallery resources, and relation data from another item-specific API.
- Evidence path: project-specific API/tool docs, not this reusable skill.

#### Generalized Lesson

Treat playable list URLs as optimization hints, not as proof that detail APIs are unnecessary. UI-to-API mapping should compare at least these operation windows: list load, item tap before playback, explicit play tap, comments/related panels, and media segment loading.

#### Agent Action

When a list item has a playable source, still test item tap with capture/replay and probe known detail-path patterns from existing API docs. In tools, prefer "hydrate detail first, then play" for normal video/detail entries unless the item type is known to be self-contained.

#### Applies When

- The item card has a playable URL in the list response.
- The UI has a distinct detail screen or right-side detail pane.
- Existing docs or captures mention item-specific detail or relation endpoints.

#### Does Not Apply When

- The product explicitly uses a self-contained feed item with no detail navigation.
- The item type is documented to use direct resources only, and item tap/play capture shows no additional non-media API.
- The goal is only a minimal playback smoke test, not API completeness.

#### Validation

Confirm by replaying or capturing:

- list load API;
- item tap API before playback;
- explicit play/HLS/media APIs;
- optional relation/comments APIs.

The detail API should return fields not present in the list item, or the capture should record that no extra non-media request appeared.

#### Promotion Target

- `WORKFLOW.md`
- `DOCUMENTATION.md`

#### Required Linked Updates

- Project-specific API docs should record the concrete list -> detail -> relation/comment flow.
- Tool docs should state whether clicking an item hydrates detail first or plays directly.
