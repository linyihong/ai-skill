> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-05 - Scroll Depth Before API Conclusion

Status: validated

#### One-line Summary

可滑動頁面不能只用首屏截圖判斷 API 來源；至少要做 top / mid / deep 滑動抽樣再下結論。

#### Human Explanation

許多 App 的 feed 頁會在首屏放置入口宮格、廣告、公告、熱門卡片或 sticky tabs，真正的列表資料在下方才開始，並且滑動後可能出現更多貼文、插入廣告、分頁載入或不同資料型態。只看首屏容易把入口 / 廣告誤判成主要 API response，也可能漏掉某個 tab 實際共用另一個 sort/filter 的證據。

#### Trigger

遇到有垂直列表、sticky tab、瀑布流、社群 feed、商品列表、媒體列表或多區塊首頁時，首屏資料不足以說明整頁資料來源。

#### Evidence

- Tool: ADB screenshot and swipe sampling, paired with authorized API replay.
- Sanitized excerpt: first screen showed entry tiles and a small number of feed cards; after multiple swipes, additional feed cards and inserted banners appeared. Text visible only after scrolling matched a replayed paginated/sorted API response, changing the initial interpretation of one UI tab.
- Evidence path: project-specific UI screenshots and API docs, not this reusable skill.

#### Generalized Lesson

UI-to-API attribution needs scroll depth coverage. Capture at least top, mid, and deeper positions for scrollable screens, then compare visible labels, titles, authors, counters, and card shapes against candidate API responses. Treat sticky headers separately from list content.

#### Agent Action

When analyzing a scrollable page:

1. Capture the initial screen.
2. Swipe through representative depths and capture screenshots after each stable state.
3. Note sticky elements that remain fixed.
4. Match visible text from scrolled cards against API page/sort candidates.
5. Only then document whether a tab has its own API, shares a sort/filter, or remains unknown.

#### Applies When

- The page scrolls vertically or horizontally.
- The first screen contains mixed UI blocks such as ads, shortcut tiles, banners, and feed items.
- The task is API completeness, not only a visual overview.

#### Does Not Apply When

- The screen is static and fully visible.
- The goal is only to identify top-level navigation labels.
- Device control is unavailable and the limitation is explicitly documented.

#### Validation

Evidence is stronger when:

- top / mid / deep screenshots show different feed cards;
- visible text from scrolled cards can be found in candidate API responses;
- page/sort parameters explain the order of cards;
- unmatched regions are explicitly marked for request capture instead of guessed.

#### Promotion Target

- `WORKFLOW.md`
- `DOCUMENTATION.md`

#### Required Linked Updates

- Project UI/API docs should record scroll depth evidence and any remaining unmatched regions.
- If tools render only the first page by default, tool docs should say whether scrolling/pagination was tested.
