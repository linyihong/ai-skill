> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-08 - Browser layout overflow integration

Status: validated

#### One-line Summary

RWD 或視覺縮放問題不能只靠 CSS marker、snapshot 或人工感覺驗證；要用真實瀏覽器 viewport 量 `document`、`body`、app shell、固定導覽與 scroll content 的 computed layout。

#### Human Explanation

小螢幕問題常被誤判成「整個頁面沒有縮小」，但根因可能不是外層容器寬度，而是某個子列、負 margin、固定寬度、grid 欄位或橫向 rail 把內層 scroll content 撐寬。只檢查 CSS 是否存在 `@media`、`width: 100%` 或 viewport meta，無法證明瀏覽器實際 layout 沒有 overflow。

應該用 headless browser 或等效瀏覽器 automation 設定目標 viewport，直接讀 computed DOM metrics，再判斷 overflow 是預期的內部橫向 scroll，還是非預期的文件或內容盒撐寬。

#### Trigger

- 使用者回報 RWD、小螢幕、mobile viewport、device emulation 或「整體沒有縮小」。
- CSS/BDD marker 已通過，但實機或瀏覽器 devtools 仍顯示布局不對。
- 頁面包含 horizontal tabs、carousel、rail、fixed tab bar、negative margin、grid poster list 或 nested scroll view。
- 先前修正只調整 media query，沒有量瀏覽器實際 `scrollWidth` / `clientWidth`。

#### Evidence

- Tool: agent-assisted frontend debugging
- Sanitized excerpt: A mobile layout appeared not to shrink even after media-query updates. Browser measurement showed the top-level document and shell matched the viewport, while an internal horizontally scrollable row widened the scroll content because of edge-to-edge negative margins.
- Evidence path: Project-specific routes, class names, live outputs, and screenshots stay in `<PROJECT_ROOT>` project tests or incident notes; this lesson only records the generalized browser-layout validation rule.

#### Generalized Lesson

When diagnosing responsive layout, measure rendered layout in a real browser context before making another CSS-only guess.

| Metric | Why it matters |
| --- | --- |
| `window.innerWidth` / `visualViewport.width` | Confirms the emulated or real viewport width. |
| `document.documentElement.scrollWidth` / `document.body.scrollWidth` | Detects page-level horizontal overflow. |
| App shell and fixed navigation bounding boxes | Confirms the outer app chrome fits the viewport. |
| Main scroll content `scrollWidth` vs `clientWidth` | Detects child rows that widen a nested scroll container. |
| Overflowing element list | Identifies the exact row, grid, image, rail, or button causing the layout risk. |
| Horizontal scroller ancestry | Separates intended internal scroll from unexpected document/content overflow. |

#### Agent Action

When RWD or visual shrinking is disputed:

1. Launch a real browser engine through the project’s existing runner, Playwright/Puppeteer, WebDriver, or browser DevTools Protocol.
2. Set representative viewports, including the reported device width and a smaller lower-bound width such as 320px.
3. Open the real route, wait for hydration/layout, and collect viewport, document, body, shell, fixed navigation, and main content metrics.
4. Build an overflow report that filters out descendants of intentionally horizontal scroll containers, then inspect the remaining offenders.
5. Fix the actual overflowing element, not only the nearest media query.
6. Add or update an integration test that asserts page-level and main-content widths stay within the viewport while allowing intentional internal scrollers.

#### Goal / Action / Validation

- Goal: Prevent responsive-layout fixes from being accepted based only on CSS markers when the browser-rendered layout still overflows.
- Action: Require browser-based computed layout measurement for disputed RWD bugs and encode the metric checks in an integration test.
- Validation or reference source: The test fails when `document` or `body` width exceeds the viewport, when app shell or fixed navigation exceeds viewport width, or when main scroll content is widened by a child row outside an intentional horizontal scroller.

#### Applies When

- The issue is visual layout, viewport fitting, RWD, nested scroll, mobile emulation, fixed bottom navigation, carousel, tab rail, or edge-to-edge content.
- A previous static test or CSS marker passed but the user still observes a layout problem.
- The page has client-side hydration, CSS modules, component libraries, or nested scroll containers that static CSS inspection may not model.

#### Does Not Apply When

- The change is purely content copy, API shape, server data, or a layout-independent behavior.
- The project already has an equivalent visual regression or layout integration suite covering the affected viewport and route.
- The bug is a known browser-specific rendering engine issue that requires manual device evidence beyond generic headless layout measurement.

#### Validation

The prevention worked when:

- A browser integration test visits the affected route at representative mobile widths.
- The test asserts document/body/app shell/fixed navigation widths do not exceed the viewport.
- The test asserts main content is not widened by child elements, except for explicitly allowed horizontal scrollers.
- Removing the real overflowing CSS fix causes the integration test to fail.

#### Promotion Target

- `workflow/software-delivery/test-strategy.md`
- `workflow/software-delivery/validation.md`
- `intelligence/engineering/execution/validation-reasoning/`

#### Required Linked Updates

- Updated `workflow/software-delivery/test-strategy.md` with a browser layout / RWD row in the test strategy gate and project test strategy checklist.
- Updated `feedback/history/development-guidance/README.md` count for committed development-guidance common lessons.
- Checked reusable guidance boundary: this lesson contains generalized browser layout measurement guidance only; project-specific routes, class names, screenshots, and live run details remain outside Ai-skill reusable docs.
