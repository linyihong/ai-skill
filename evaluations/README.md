# Evaluations

`evaluations/` 存放 AI decision path 的評估結果記錄。每個 evaluation 對應一次 validation scenario 的實際執行結果，記錄 route correctness、heuristic obedience、forbidden routes 等指標。

## 用途

- 記錄 validation scenario 的執行結果（passed / failed）
- 提供 regression 追蹤：同一 scenario 在不同時間或不同模型的執行結果比較
- 不取代 validation scenario 定義（scenario 定義在 [`validation/`](../validation/README.md)）

## 目錄結構

```
evaluations/
  failure-derived/          # Failure-derived scenarios 的執行結果
    <scenario-name>-<timestamp>.yaml
```

## 目前記錄

| 檔案 | Scenario | 結果 |
|------|----------|------|
| [`failure-derived/feedback-history-consolidation-v1-2026-05-13.yaml`](failure-derived/feedback-history-consolidation-v1-2026-05-13.yaml) | feedback-history-consolidation-v1 | ✅ passed (6/6) |

## Evaluation 格式

每個 evaluation YAML 包含：

- `scenario`：對應的 scenario ID
- `timestamp`：執行時間
- `model`：執行模型
- `result`：包含 route_correctness、heuristic_obedience、forbidden_routes_used、intelligence_usage 等指標

詳細格式請見 [`validation/evaluations/template.yaml`](../validation/evaluations/template.yaml)。

## 與既有層的關係

- [`validation/`](../validation/README.md)：scenario 定義與驗證框架
- [`traces/`](../traces/README.md)：AI 的實際 decision trace（與 evaluation 成對）
- [`shared-rules/failure-learning-system.md`](../shared-rules/failure-learning-system.md)：Failure → Scenario 閉環流程
