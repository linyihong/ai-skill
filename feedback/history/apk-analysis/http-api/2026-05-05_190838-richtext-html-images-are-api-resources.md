> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/apk-analysis/artifact-gates.md`](../../../../workflow/apk-analysis/artifact-gates.md)

### 2026-05-05 - RichText HTML Images Are API Resources

Status: validated

#### One-line Summary

詳情 API 的可下載圖片不一定只在 JSON array，還可能藏在 rich text / HTML 欄位的 `<img src>` 裡。

#### Human Explanation

很多 feed / blog / post 詳情 API 會同時回結構化欄位與富文本欄位。結構化欄位可能只有封面、縮圖或預覽圖；真正正文中的圖片、GIF、長圖或步驟圖，可能在 `richText`、`html`、`contentHtml`、`body` 等欄位裡以 `<img src="...">` 形式出現。若 client 或文件只讀 `images[]`、`covers[]`、`seriesCover[]`，就會漏掉使用者在詳情頁實際看到的圖片資源。

#### Trigger

App 詳情頁顯示圖文混排、blog-like post、商品詳情、文章詳情、社群貼文或富文本內容，但 API response 的圖片陣列數量和 UI 看到的圖片數不一致。

#### Evidence

- Tool: authorized API replay and UI screenshot comparison.
- Sanitized excerpt: a detail response contained both structured image arrays and a rich text field; parsing `<img src>` from rich text found additional image/GIF resources that matched the detail page content flow.
- Evidence path: project-specific API docs, not this reusable skill.

#### Generalized Lesson

When documenting or implementing resource download for detail pages, scan both structured JSON image fields and rich text / HTML fields. Treat discovered `<img src>` values as first-class resource candidates after URL normalization, deduplication, and sanitization.

#### Agent Action

When a detail response has rich text:

1. Search likely HTML fields such as `richText`, `html`, `contentHtml`, `body`, `descriptionHtml`.
2. Extract `<img src>` values with an HTML parser or carefully scoped extraction.
3. Normalize relative paths using the same image host/prefix rules as structured image fields.
4. Merge and dedupe with `cover`, `thumbnail`, `images[]`, `seriesCover[]`, or similar arrays.
5. Document which fields feed gallery/download behavior.

#### Applies When

- The UI shows rich text or image-heavy detail content.
- The API response contains HTML-like strings.
- The goal includes downloading or rendering all visible resources.

#### Does Not Apply When

- Rich text is unavailable or sanitized server-side with no embedded resources.
- The client intentionally supports only cover/thumbnail resources.
- The HTML is untrusted and must not be rendered directly; in that case extract URLs safely but do not inject raw HTML.

#### Validation

Confirm by comparing:

- count and type of structured image fields;
- count and type of rich text `<img src>` resources;
- UI screenshots across scroll depth;
- successful normalized image fetch/decode for at least one structured image and one rich text image.

#### Promotion Target

- `DOCUMENTATION.md`
- `techniques/http-api/README.md`

#### Required Linked Updates

- Project API docs should name the rich text fields and resource normalization rule.
- Client/tool docs should state whether gallery/download includes rich text images.
