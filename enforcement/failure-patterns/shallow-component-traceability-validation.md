# Shallow Component Traceability Validation（共用元件追蹤驗證過淺）

Status: candidate
Class: `validation-gap`

## Trigger

當任務新增或提升共用 UI component，且驗收只檢查 component inventory 名稱、fixture 名稱、檔案是否存在或 broad directory ref，而沒有檢查具體 source path 與最小實作語意 marker 時，使用此 pattern。

常見訊號：

- 使用者詢問「feature / test 怎麼沒有指到這個 component」。
- BDD passed，但 feature refs 只指到 component 目錄，沒有指到實作檔。
- 測試只 assert `"ComponentName"`，沒有讀取 component source 或檢查 props / state / interaction primitive。
- 新增 common component 後只更新 component index，未更新 feature traceability。

## Failure Mode

Agent 把 component traceability 當成名稱清單維護，導致 BDD / contract test 即使通過，也不能證明：

- feature 或 contract 能追到具體 component 檔案；
- component 的核心責任仍存在；
- route page 確實引用了共用 component；
- component 行為的關鍵 props、狀態、accessibility、權限提示、空態或互動控制沒有被移除。

## Risk

- UI component 看似已納入企劃與測試，但實際實作可以漂移而不被 BDD 發現。
- 後續 agent 只看到 inventory 名稱，無法定位 owner file 或最小驗收責任。
- Review 需要人工記憶補洞，traceability 失去防回歸價值。
- Shared component promotion 可能形成「放到共用資料夾但未定義責任」的長期維護負擔。

## Required Agent Action

1. 將目前驗收重新分類為 traceability depth 問題，不只當成少補一行 ref。
2. 在 feature / contract 中補上 component 名稱與具體 source path。
3. 在 component index 或 owner README 中補上 component purpose 與 scope。
4. 在 BDD / contract test 中讀取具體 component source，檢查最小語意 marker。
5. 若 component 是由 route page 使用，測試 route page 是否 import / render 該共用 component。
6. 若此缺口來自 agent 未主動回饋，新增或更新 feedback lesson / failure pattern。

## Prevention Gate

當新增或提升 shared UI component 時，agent 在宣稱完成前必須回答：

| Check | Required answer |
| --- | --- |
| Feature trace | 哪個 feature / contract 以具體 path 指到 component source？ |
| Ownership | 哪個 component index / owner README 說明它為何在共用層？ |
| Semantic marker | BDD / contract test 檢查哪些 props、state、interaction primitive、accessibility label、empty state 或 permission marker？ |
| Usage marker | 需要時，哪個 route page / parent component 被檢查有引用它？ |
| Negative drift | 移除核心 marker 時，哪個 test 會失敗？ |

## 驗證

此 pattern 已套用時，應可反查：

- feature / contract refs 包含具體 component source path；
- BDD / contract test 讀取該 component source 或 focused fixture；
- 測試包含至少一個與 component 核心責任相關的 marker；
- component index / owner README 已同步；
- 若是 route-level component composition，route page 的 import / render 也有測試 marker。

## Linked Rules

- [`../failure-learning-system.md`](../failure-learning-system.md)
- [`../linked-updates.md`](../linked-updates.md)
- [`../reusable-guidance-boundary.md`](../reusable-guidance-boundary.md)
- [`../../feedback/feedback-lessons.md`](../../feedback/feedback-lessons.md)

## Linked Validation Scenarios

- candidate: shared-component-traceability-depth-v1 — 檢查 shared component 新增時，feature refs、component index、source semantic markers 與 route usage marker 是否同時存在。

← [Back to failure patterns](README.md)
