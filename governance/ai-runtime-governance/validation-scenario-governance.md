# Validation Scenario Governance

## Source Intelligence

source_intelligence:

- [`intelligence/engineering/agent-architecture/stateless-validation-necessity.md`](../../intelligence/engineering/agent-architecture/stateless-validation-necessity.md)
- [`intelligence/engineering/agent-architecture/failure-to-scenario-closure.md`](../../intelligence/engineering/agent-architecture/failure-to-scenario-closure.md)
- [`intelligence/engineering/agent-architecture/cognitive-boundaries.md`](../../intelligence/engineering/agent-architecture/cognitive-boundaries.md)

本文件把 stateless validation 與 failure-to-scenario closure 的 agent architecture intelligence 轉譯成 AI runtime validation governance。原始思想回答「為什麼 AI decision path 驗證不能依賴 conversation memory」；本文件定義何時需要 scenario、scenario 必須滿足哪些 gate，以及 failure 如何 promotion 成可重放驗證。

## 觸發時機

在下列情況套用本治理：

- Routing、heuristic obedience、rule obedience 或 recovery behavior 出現可重現 failure。
- 新增或修改 `validation/scenarios/`、`validation/rules/`、`validation/traces/` 或 `validation/evaluations/`。
- 模型升級、routing registry 變更、workflow 分層變更後，需要確認 AI decision path 未退化。
- Failure pattern 或 feedback lesson 已成熟，且具有 stateless reproduction value。

## Runtime Gate

| Gate | 通過條件 |
| --- | --- |
| Stateless reproduction | Scenario 的 `given` 包含完整決策條件，不依賴前文、session memory 或先前對話。 |
| No answer leakage | 給 agent 的 scenario input 不包含 expected route、forbidden route 或 evaluator-only answer hints。 |
| Failure class clarity | Scenario 對應的 failure class、source pattern 或 source lesson 已清楚標註。 |
| Expected and forbidden behavior | 可描述 expected route / heuristic，以及至少一個 forbidden route 或 forbidden behavior。 |
| Traceability | Scenario 可產出 trace，包含 signals、loaded intelligence、rejected routes 與 final route。 |
| Maintenance boundary | 一次性、環境依賴、工具故障或不可重放事件不升格為 scenario。 |

## 分層判斷

| 內容類型 | 目標層 |
| --- | --- |
| 為什麼 AI validation 必須 stateless、為什麼 failure 要變 scenario | `intelligence/engineering/agent-architecture/` |
| Scenario promotion 條件、answer leakage gate、traceability gate | `governance/ai-runtime-governance/` |
| Scenario / trace / evaluation 檔案格式與執行方式 | `validation/` |
| Failure capture、classification、promotion decision | `enforcement/failure-learning-system.md` |
| 實際 failure patterns | `enforcement/failure-patterns/` |

## Validation Mapping

- [`validation/README.md`](../../validation/README.md) — scenario、trace 與 evaluation 格式。
- [`enforcement/failure-learning-system.md`](../../enforcement/failure-learning-system.md) — failure capture 與 promotion decision。
- [`validation/scenarios/failure-derived/`](../../validation/scenarios/failure-derived/) — failure-derived scenarios。

## Validation Candidate

本治理本身可用 scenario 檢查：

- Scenario 依賴「還記得之前」之類 conversation state。
- Scenario input 洩漏 expected route，讓 agent 直接猜答案。
- Failure pattern 已標為 repeated，但沒有 scenario 或明確 not-applicable 理由。
- Scenario 沒有 traceable rejected routes 或 final route。
