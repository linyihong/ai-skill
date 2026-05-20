## skill.travel-planning

| 欄位 | 值 |
| --- | --- |
| Atom ID | `workflow.travel-planning` |
| Source path | `workflow/travel-planning/execution-flow.md` |
| Lifecycle | `candidate` |
| Summary | 依目的地、日期、交通與玩法規劃行程，包含營業時間查證、交通比較、住宿與備案。支援 itinerary 結構化輸出與可行性檢查。日本自駕含 Mapcode 粒度規則（沿線景點 2km+ 需各停車點獨立一行）與查詢工具鏈。 |
| When to read | 使用者要求旅遊路線、交通、餐飲、住宿或行程規劃時；或日本自駕行程需要 Mapcode 表時。 |
| Do not use for | 不可取代即時票價/房價查詢 API。不可用於未經查證的營業時間或交通時刻表。 |
| Context cost | ~280 tokens |
| Estimated full cost | ~2800 tokens |
| Validation signal | Workflow entrypoint links 可解析，execution flow 與 artifact gates 結構完整。 |
| Last checked | 2026-05-21 |
