> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/apk-analysis/execution-flow.md`](../../../../workflow/apk-analysis/execution-flow.md)

### 2026-05-05 - Page analysis requires a UI map artifact

Status: promoted

#### One-line Summary

當使用者指定某個 App 頁面做 APK/API 分析時，分析結果必須產出或更新一份專案內的 page-level 輕量 UI 架構地圖，不能只把結論散落在 API 或工具文件。

#### Human Explanation

頁面分析通常同時涉及入口、tab、可點擊卡片、scroll pagination、詳情頁、評論、媒體資源與 SDK mapping。若只更新 endpoint 文件或工具說明，下一位 agent 雖然能看到 API，但不知道 UI 怎麼進、哪些區塊尚未抓包、哪幾種內容格式已驗證，也不容易複用成其他頁面的補文流程。

可重用規則是：只要任務是「針對某個頁面」而且已經做了實機操作、截圖、滑動、API 對照或 SDK 驗證，就必須在專案 docs 內留下一份輕量 UI 地圖，或明確更新既有地圖。這份地圖是 target-specific artifact，不放進 reusable skill；skill 只保存模板與提醒。

#### Trigger

使用者要求以下任一類工作時適用：

- 分析某個具名頁面、tab、module 或功能入口。
- 比對手機畫面和 API / docs/tools / SDK 是否一致。
- 要求滑動、截圖、點詳情、下載資源或確認頁面功能是否缺失。
- 從 APK 分析結果延伸到 SDK、final validation 或交付包。

#### Evidence

- Page artifact example: `docs/UI架構地圖/<page-name>.md`
- Evidence types: sanitized screenshots, UI path, operation notes, `/modules/list`, list API, detail API, comment API, final validation output.
- Reusable skill should only reference generalized template; project-specific endpoints, module IDs, titles, screenshots, and output paths stay in the project docs.

#### Generalized Lesson

Page-level APK analysis should end with a durable UI map artifact in the project repository. The artifact should connect UI structure to API operations and implementation/testing surfaces. API docs explain paths; the UI map explains how a user reaches and exercises the feature.

#### Agent Action

Before ending a page-specific APK analysis task, check:

- Is there a `docs/UI架構地圖/` or equivalent project docs area?
- Did this task identify page entry path, tabs, cards, feed, detail, comments, media, or pagination?
- Did UI observations lead to docs/tools/SDK/test changes?
- If yes, create or update a page map and link it from a nearby index/tool/API doc.

Minimum page map sections:

1. Entry path.
2. UI block map.
3. Scroll / pagination behavior.
4. Detail content formats.
5. API chain.
6. SDK or client field mapping.
7. Validation output or test evidence.
8. Follow-up rules for analyzing similar pages.

#### Applies When

- The app can be operated or screenshots/UI hierarchy are available.
- The analysis target is a screen, page, module, tab, or named feature.
- Results are expected to guide future agents, SDK work, docs/tools, or tests.

#### Does Not Apply When

- The task is pure static analysis with no specific screen or feature.
- The user only asks for a one-off command result or narrow endpoint lookup.
- Authorization forbids saving UI descriptions; in that case, record an abstract route map without sensitive labels.

#### Validation

- A future reader can follow the map to reach the page and know which API each major UI block uses.
- The map identifies unknown UI blocks as `needs capture` rather than guessing.
- The map links to or names validation artifacts such as final output, tests, screenshots, or sanitized API docs.

#### Promotion Target

- `SKILL.md`
- `DOCUMENTATION.md`
- `feedback_history/README.md`

#### Required Linked Updates

- Add a Quick Start checklist item requiring page-level UI map artifacts for page-specific analyses.
- Add a documentation rule explaining where the artifact belongs and what sections it must contain.
- Add this lesson to the feedback index.
