# Rule Dependency Graphs

`knowledge/graphs/rules/` 存放 enforcement rule 之間的依賴關係圖，將 [`enforcement/linked-updates.md`](../../enforcement/linked-updates.md) 中的文字依賴表轉換為 machine-readable graph records。

## Source of Truth

所有依賴關係的 **source of truth** 是 [`enforcement/linked-updates.md`](../../enforcement/linked-updates.md) 的「常見連動關係」表格（第 19-52 行）。Graph records 是該表格的 machine-readable 映射，不可與表格內容矛盾。

## Graph Records

| Graph record | 用途 | 狀態 |
| --- | --- | --- |
| [`core-bootstrap.yaml`](core-bootstrap.yaml) | rule-weight、dependency-reading、conversation-goal-ledger 之間的關係 | `candidate` |
| [`linked-updates-deps.yaml`](linked-updates-deps.yaml) | linked-updates 依賴的所有 rules | `candidate` |
| [`failure-learning-deps.yaml`](failure-learning-deps.yaml) | failure-learning-system 依賴的所有 rules | `candidate` |
| [`content-layering-deps.yaml`](content-layering-deps.yaml) | content-layering 被哪些 rules 依賴 | `candidate` |
| [`full-rule-graph.yaml`](full-rule-graph.yaml) | 完整的 rule dependency graph，含所有 16 條 enforcement rules | `candidate` |

## Edge Types Used

| Edge | 意義 |
| --- | --- |
| `depends_on` | 修改此 rule 時，必須同步更新或檢查 target rule（對應 linked-updates.md 的「必須同步更新或檢查」） |
| `related_to` | Target rule 可能有幫助，但不是強制連動 |
| `routes_to` | 此 rule 的索引或 routing 指向 target |

## 查詢方式

```bash
ai-skill runtime query --graph --source knowledge/graphs/rules/ --limit 10
```

## 驗證

- 每個 edge 必須能對應到 `enforcement/linked-updates.md` 中的一行。
- 新增或刪除 enforcement rule 時，必須同步更新對應的 graph records。
- Graph records 不可與 `metadata/rules/*.yaml` 中的 `depends` 欄位矛盾。
