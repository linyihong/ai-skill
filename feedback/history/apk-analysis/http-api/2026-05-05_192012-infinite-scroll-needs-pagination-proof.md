> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/http-api-documentation-flow.md`](../../../../analysis/apk/workflows/http-api-documentation-flow.md)

### 2026-05-05 - Infinite Scroll Needs Pagination Proof

Status: validated

#### One-line Summary

分析可下拉 feed 時，要證明下拉如何分頁：請求參數、下一頁判斷、去重與停止條件都要記錄。

#### Human Explanation

很多 App feed 首屏只是一頁資料。使用者下拉後，App 可能持續用 `pageNumber+1`、cursor、timestamp、offset 或 search-after token 拉下一頁；也可能插入廣告、置頂內容或混合推薦。若只記首屏 API，文件會漏掉長時間採集最重要的分頁規則。

#### Trigger

遇到瀑布流、社群 feed、商品列表、影片列表、評論區、tag 列表或任何「一直往下拉還會出新資料」的畫面。

#### Evidence

- Tool: authorized API replay paired with scroll-depth screenshots.
- Sanitized excerpt: page 1/2/3 returned different item ids while `hasNext=true`; visible scrolled cards matched the paginated API response.
- Evidence path: project-specific API docs, not this reusable skill.

#### Generalized Lesson

UI-to-API binding for scrollable feeds is incomplete until pagination is tested. Record the request fields that change between pages, the response field that says whether more data exists, and whether inserted ads or sticky cards are API items or UI-only inserts.

#### Agent Action

For every scrollable feed:

1. Capture top/mid/deep UI evidence.
2. Replay page 1, 2, and 3 or equivalent cursor transitions.
3. Compare first item ids/titles between pages to prove they are distinct.
4. Record `hasNext`, `nextCursor`, `total`, empty-list behavior, or any other stop signal.
5. Document whether tools support manual page navigation, auto infinite scroll, or both.

#### Applies When

- The UI keeps loading new cards while scrolling.
- A response contains fields such as `hasNext`, `pageNumber`, `cursor`, `next`, `offset`, or `total`.
- The goal includes collection, SDK pagination, or external tooling.

#### Does Not Apply When

- The screen is static or only has local scroll over already-loaded content.
- Pagination cannot be tested due missing authorization or unstable network; then document the gap instead.

#### Validation

At minimum, provide:

- two or more distinct pages/cursors;
- a stop/continue field such as `hasNext`;
- proof that visible scrolled content maps to a later page or already-loaded appended data.

#### Promotion Target

- `WORKFLOW.md`
- `DOCUMENTATION.md`

#### Required Linked Updates

- Project API docs should describe pagination request/response fields and stop conditions.
- Tool docs should state whether the tool auto-loads on scroll or requires manual page navigation.
