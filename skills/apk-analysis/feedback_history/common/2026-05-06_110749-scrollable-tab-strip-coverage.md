> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/apk-analysis/execution-flow.md`](../../../../workflow/apk-analysis/execution-flow.md)

### 2026-05-06 - Scrollable Tab Strip Coverage

Status: promoted

#### One-line Summary

分析 tab / category / result tabs 時，必須檢查是否可水平滑動；只點可見 tab 不能宣稱完整覆蓋。

#### Human Explanation

Mobile UI 常把功能藏在水平 tab strip、segmented control、category chip row 或結果分類列中。UIAutomator XML 可能只顯示目前 viewport 的可見 tab，或顯示「共 N 個」但最後幾個在螢幕外。若只根據第一屏可見 tab 建 API 清單，容易漏掉 off-screen tab 的 API、預載、分頁或媒體入口。

#### Trigger

使用者指出：tab 可能可以左右滑動，導致已分析功能不完整。

#### Evidence

- Tool: UIAutomator hierarchy / screenshot / read-only swipe replay.
- Sanitized excerpt: result tabs showed visible labels but hierarchy reported more total tabs than visible viewport; category strip showed only first visible categories while API returned a larger category count.
- Evidence path: project docs should record viewport-visible tabs, total tab count if exposed, and horizontal swipe evidence separately.

#### Generalized Lesson

任何 named feature 的 UI/API 覆蓋若包含 tab、top tabs、category chips、filter strip、search result tabs、carousel-like modules，必須把「可見項」和「可水平滑動後可達項」分開記錄。只有完成左右滑動探測並對每個 reachable tab 標記 captured / no-network / needs capture，才可宣稱 tab 面完整。

#### Agent Action

下次看到 tab strip 時：

1. 記錄 first viewport 的 visible labels / selected state / bounds。
2. 檢查 XML 是否有 `第 X 个标签，共 N 个`、scrollable node、off-screen label、或 API 回傳 count 大於可見 count。
3. 做最小 read-only horizontal swipe：向左滑 tab strip 直到沒有新 tab，再必要時向右回復。
4. 對每個 tab 建 operation id 或 gap row，並標 `Confirmed`、`Candidate+`、`no same-window API observed` 或 `needs capture`。
5. 不要把同一 API family 的一個 tab 推廣成所有 tabs 的結論，除非每個 tab 都有同窗證據或已證明只是 client-side filter。

#### Goal / Action / Validation

- Goal: 避免因 off-screen tabs 漏掉功能/API。
- Action: 在 UI map 和 operation map 中加入 horizontal tab sweep 或明確 gap。
- Validation or reference source: 比對可見 tab 數、hierarchy total count、swipe 後新增 labels、以及每個 tab 的 API capture window。

#### Applies When

- App 使用 top tabs、bottom sub-tabs、category chips、filter strip、search result tabs、carousel tabs。
- 使用者要求完整分析某一頁、某一功能、或所有 API。
- UI hierarchy 顯示 total count 大於可見 count，或設計上看起來可水平捲動。

#### Does Not Apply When

- tab 容器已證明不可滑動，且所有 tab 完整可見。
- 分析範圍明確只限某個指定 tab，文件已標明不代表其他 tab。

#### Validation

- 截圖/XML 證明滑動前後 viewport。
- 操作記錄包含 horizontal swipe 座標與方向。
- API 文件對每個 reachable tab 都有 evidence 或 gap。

#### Promotion Target

- `WORKFLOW.md`
- Project UI map / operation map docs

#### Required Linked Updates

- 已同步更新 `WORKFLOW.md` 的 UI coverage 規則。
- 專案文件應同步更新目標頁 UI map 和 operation map；不要把 raw tab labels、結果標題或個資寫進 reusable skill。
