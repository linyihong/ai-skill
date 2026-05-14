# Traces

`traces/` 存放 AI 的實際 decision trace 記錄。每個 trace 對應一次 validation scenario 的執行過程，記錄 signals detected、heuristics used、rejected routes、final route 與 decision rationale。

## 用途

- 記錄 AI 在執行 validation scenario 時的完整決策過程
- 提供 decision path 的可追溯性與審計能力
- 與 evaluation 成對：trace 記錄「怎麼決策」，evaluation 記錄「決策是否正確」

## 目錄結構

```
traces/
  failure-derived/          # Failure-derived scenarios 的 decision trace
    <scenario-name>-<timestamp>.yaml
```

## 目前記錄

| 檔案 | Scenario |
|------|----------|
| [`failure-derived/feedback-history-consolidation-v1-2026-05-13.yaml`](failure-derived/feedback-history-consolidation-v1-2026-05-13.yaml) | feedback-history-consolidation-v1 |

## Trace 格式

每個 trace YAML 包含：

- `scenario`：對應的 scenario ID
- `timestamp`：執行時間
- `model`：執行模型
- `trace`：包含 signals_detected、heuristics_used、rejected_routes、final_route、intelligence_loaded、decision_rationale

詳細格式請見 [`validation/traces/template.yaml`](../validation/traces/template.yaml)。

## 誰會參考這裡（Inbound References）

變更本層內容時，需要一併檢查以下依賴方：

| 來源 | 關係 |
|------|------|
| [`route.traces.decision-traces`](../knowledge/runtime/routing-registry.yaml) | Routing registry record，agent 依此找到 traces/ |
| [`validation/`](../validation/README.md) | Scenario 執行後寫入 trace 至此 |
| [`enforcement/failure-learning-system.md`](../enforcement/failure-learning-system.md) | Failure → Scenario 閉環完成後產出 trace |

## 與既有層的關係

- [`validation/`](../validation/README.md)：scenario 定義與驗證框架
- [`evaluations/`](../evaluations/README.md)：scenario 執行結果（與 trace 成對）
- [`enforcement/failure-learning-system.md`](../enforcement/failure-learning-system.md)：Failure → Scenario 閉環流程
