# Recovery Retry Governance

## Source Intelligence

source_intelligence:

- [`intelligence/engineering/agent-architecture/failure-recovery.md`](../../intelligence/engineering/agent-architecture/failure-recovery.md)
- [`intelligence/engineering/agent-architecture/cognitive-boundaries.md`](../../intelligence/engineering/agent-architecture/cognitive-boundaries.md)
- [`intelligence/engineering/agent-architecture/context-collapse.md`](../../intelligence/engineering/agent-architecture/context-collapse.md)

本文件把 failure recovery 與 cognitive boundaries 的 agent architecture intelligence 轉譯成 AI runtime recovery governance。原始思想回答「為什麼無策略重試會降低品質、為什麼 agent 不能可靠感知自己的認知邊界」；本文件定義 retry、escalation、source reload 與 recovery validation 的治理 gate。

## 觸發時機

在下列情況套用本治理：

- 同一操作、route、automation 或 checkpoint 連續失敗。
- 使用者指出 source miss、wrong workflow、guessing、或與 owner docs 衝突。
- Evidence 推翻原 execution graph，或 agent 開始用相同策略重試。
- 修改 `enforcement/escalation-policy.md`、`metadata/recovery/`、runtime recovery state 或 recovery validation。

## Runtime Gate

| Gate | 通過條件 |
| --- | --- |
| Retry limit | 同一操作最多一次 simple retry；第二次失敗前必須說明策略差異。 |
| Strategy change | 每次 retry 都要改變 evidence、tool、source reading 或 hypothesis；不得重複同一操作。 |
| Escalation threshold | 連續失敗、user contradiction、evidence conflict 或 source miss 時，進入 L3+ escalation。 |
| Source reload | L3+ 後必須讀 required source-of-truth，或標記 `not_applicable` / `source_missing`。 |
| Execution graph rebuild | 恢復執行前重建 `goal -> route -> dependencies -> checkpoint -> validation`。 |
| Recovery validation | Recovery 成功後必須用外部 evidence 驗證，不接受 agent 自評為完成。 |

## 分層判斷

| 內容類型 | 目標層 |
| --- | --- |
| 為什麼重複 retry 會污染 context、為什麼需要外部 gate | `intelligence/engineering/agent-architecture/` |
| Retry limit、strategy change、source reload、recovery validation gate | `governance/ai-runtime-governance/` |
| 即時 mismatch escalation 條文 | `enforcement/escalation-policy.md` |
| Domain-specific reload set 與 forbidden behavior metadata | `metadata/recovery/` |
| Machine-readable recovery state / repair strategy | `runtime/` |

## Runtime Mapping

- [`enforcement/escalation-policy.md`](../../enforcement/escalation-policy.md) — real-time mismatch escalation policy。
- [`metadata/recovery/README.md`](../../metadata/recovery/README.md) — domain-specific recovery metadata。
- [`runtime/README.md`](../../runtime/README.md) — runtime state machine 與 recovery compiled surfaces。
- [`enforcement/failure-learning-system.md`](../../enforcement/failure-learning-system.md) — recovery 後的 post-mortem prevention loop。

## Validation Candidate

後續若要 promotion 到 `validation/`，可建立 scenario 檢查：

- Agent 對相同錯誤重試 2+ 次但沒有策略變更。
- L3+ escalation 後未重讀 source-of-truth。
- User contradiction 後 agent 繼續沿用舊 execution graph。
- Recovery 完成宣告缺少外部 validation evidence。
