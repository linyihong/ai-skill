# Routing Signal Governance

## Source Intelligence

source_intelligence:

- [`intelligence/engineering/agent-architecture/task-routing.md`](../../intelligence/engineering/agent-architecture/task-routing.md)
- [`intelligence/engineering/agent-architecture/attention-budgeting.md`](../../intelligence/engineering/agent-architecture/attention-budgeting.md)
- [`intelligence/engineering/agent-architecture/cognitive-boundaries.md`](../../intelligence/engineering/agent-architecture/cognitive-boundaries.md)

本文件把 task-routing 的 agent architecture intelligence 轉譯成 AI runtime 路由治理。原始思想回答「為什麼 agent 會被強信號帶偏」；本文件定義 routing decision 需要通過的治理 gate。

## 觸發時機

在下列情況套用本治理：

- 新增或修改 `knowledge/runtime/routing-registry.yaml` route。
- 修改 `workflow/workflow-routing.md`、`knowledge/indexes/README.md` 或 routing discovery 流程。
- 多個 route 同時命中，或使用者指出 agent 走錯 workflow / source。
- 最近讀取檔、相似目錄名、工具偏好或第一個搜尋結果可能覆蓋 primary source。

## Runtime Gate

| Gate | 通過條件 |
| --- | --- |
| Task intent before path signal | 先用使用者目標定義 task intent；不得讓最近打開的檔案或相似路徑名稱直接決定 route。 |
| Primary source first | 讀 registry 時先檢查 `primary_source` 與 `required_dependencies`，再看 `candidate_sources`。 |
| Negative signal check | 若 route 有明確 out-of-scope、deprecated、forbidden route 或 stale source 信號，必須列入排除理由。 |
| Multi-route disambiguation | 多個 route 命中時，列出候選 route、拒絕理由與阻擋問題；不可默默採第一個匹配。 |
| Route validation signal | 選定 route 前，確認該 route 有 validation signal，且能解釋本輪 deliverable。 |
| Recovery re-entry | 發現 route stale 或 source-of-truth mismatch 時，停止 patch，重新跑 routing discovery。 |

## 分層判斷

| 內容類型 | 目標層 |
| --- | --- |
| Agent 為什麼會受強信號影響、路徑漂移、第一匹配偏誤 | `intelligence/engineering/agent-architecture/` |
| Routing decision 必須通過的 gate、排除信號、re-entry 條件 | `governance/ai-runtime-governance/` |
| Workflow route 選擇表與常見歧義 | `workflow/workflow-routing.md` |
| Machine-readable route、trigger、primary source、validation signal | `knowledge/runtime/routing-registry.yaml` |
| 可測的 wrong-route / forbidden-route failure mode | `validation/` |

## Workflow Mapping

- [`workflow/workflow-routing.md`](../../workflow/workflow-routing.md) — workflow route 選擇與多 route 歧義裁決。
- [`governance/lifecycle/routing-philosophy.md`](../lifecycle/routing-philosophy.md) — routing pipeline 的 design layer。
- [`knowledge/runtime/routing-registry.yaml`](../../knowledge/runtime/routing-registry.yaml) — machine-readable route records。

## Validation Candidate

後續若要 promotion 到 `validation/`，可建立 scenario 檢查：

- Agent 因最近讀取的檔案而忽略 registry `primary_source`。
- 多個 `route.workflow.*` 命中時，agent 未列出候選與拒絕理由。
- Agent 使用 `candidate_sources` 取代 `required_dependencies`。
- Route 選擇無法對應到該 route 的 validation signal。
