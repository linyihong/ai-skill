# Model Profiles

`models/profiles/` 定義工具中立的 model profile。Profile 用來描述 context loading 深度、適合任務與壓縮策略，不指定特定品牌、版本或工具 UI。

## Profile Types

| Profile | 適合任務 | Context loading 深度 | 必要 guardrail |
| --- | --- | --- | --- |
| `small` | 快速 routing、檢查清單、格式套用、低風險摘要。 | 先讀 index、summary、registry；只在 validation 需要時讀 primary source。 | 不可跳過 required bootstrap、source-of-truth gate 或 validation gate。 |
| `large` | 跨層規劃、規則更新、migration、複雜 debugging、需要多來源整合的任務。 | 讀 primary source、required dependencies、related sources；必要時讀 graph / summaries 對照。 | 必須回報 deferred sources 與 validation signal。 |
| `specialized` | 需要特定工具、語言、domain 或資料格式能力的任務。 | 先讀 routing registry 與 primary source，再讀該 domain 的 workflow / technique / adapter。 | 不得讓工具能力覆蓋 shared rules 或 source-of-truth。 |

## Routing Rules

1. 任務若涉及 safety、source-of-truth、commit/push/readback 或 shared rule 更新，最低需套用 `large` profile 的讀取深度。
2. 任務若只是定位入口、查詢狀態或使用已驗證 checklist，可使用 `small` profile。
3. 任務若需要 APK、app guidance、travel planning 或 tool adapter 的專門流程，使用 `specialized` profile，但仍要先遵守 shared-rule bootstrap。
4. Profile 只決定 context loading 深度，不決定規則權重。規則衝突仍依 `enforcement/rule-weight.md`。

## Metadata Mapping

| Metadata field | Profile use |
| --- | --- |
| `models.small` | 小模型可讀的 summary / checklist / compressed guidance。 |
| `models.large` | 大模型可讀的 full source、graph 與 related references。 |
| `models.specialized` | 需要特定 domain / tool / data format 的補充說明。 |
| `context_cost` | 決定是否先讀 summary 或直接讀 full source。 |
| `complexity` | 決定是否需要 large 或 specialized profile。 |

## Validation

Profile 選擇有效時，agent 應能說明：

- 使用哪個 profile。
- 先讀哪些 sources。
- 哪些 sources 被延後。
- 哪個 validation signal 證明讀取深度足夠。
