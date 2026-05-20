# Metadata Ranking

`metadata/ranking/` 定義 agent 在 context loading 時，如何排序 candidate Knowledge Atoms 與 source files。

## 排序輸入

使用 `metadata/schema.md` 的下列欄位：

| 欄位 | 排序效果 |
| --- | --- |
| `priority` | 高優先序先載入：`P0`、`P1`、`P2`、`P3`。 |
| `status` | 優先 `stable` 與 `validated`，其次 `candidate`；除非為了相容性，避免載入 `deprecated`。 |
| `confidence` | 任務允許選擇時，優先 `high`，再來 `medium`，最後 `low`。 |
| `context_cost` | 兩個來源能回答同一問題時，優先低成本來源。 |
| `depends` | Atom 使用前必須先載入的依賴。 |
| `conflicts` | 若排序會載入不相容 atom，先暫停並解決衝突。 |
| `when_to_read` | 只有觸發條件符合任務時，才把 atom 納入排序。 |

## 預設排序

1. 必要的 safety、source-of-truth、dependency reading 與 validation rules。
2. 最新 user goal 與 active `.agent-goals/` state。
3. 目前 source-of-truth entrypoints，尤其是 `enforcement/` 與 `workflow/<domain>/execution-flow.md`。
4. 直接符合 task intent 的 `validated` 或 `stable` Knowledge Atoms。
5. 能協助導航、但不取代 source 行為的 candidate maps 與 summaries。
6. Examples、background references 與 optional optimization notes。

## 同分排序

多個來源都相關時：

- 優先選擇與目前決策語意距離最低的來源。
- 優先選擇能提供具體 validation signal 的來源。
- 優先選擇 canonical repository paths，而不是 tool mirrors。
- 只有短 index 或 summary 指向 canonical source 時，才優先先讀短內容。
- 不可為了節省 context 跳過 required dependencies。

## 停止條件

符合下列情況時，停止載入更多 context：

- 目前來源已回答決策點。
- 額外來源信心更低，或只是重複同一 guidance。
- 衝突需要 rule-weight 或 user decision。
- Candidate map 顯示舊 skill 仍是 source of truth，且本次不在 promotion scope。

## 驗證

排序路線有效時，final answer 或 commit 應能說明：

- 先載入哪個 source。
- 哪些 dependencies 是必讀。
- 哪些 sources 被延後，以及原因。
- 哪個 validation signal 證明所選 source 已足夠。
