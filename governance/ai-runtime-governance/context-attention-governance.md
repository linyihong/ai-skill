# Context Attention Governance

## Source Intelligence

source_intelligence:

- [`intelligence/engineering/agent-architecture/context-collapse.md`](../../intelligence/engineering/agent-architecture/context-collapse.md)
- [`intelligence/engineering/agent-architecture/attention-budgeting.md`](../../intelligence/engineering/agent-architecture/attention-budgeting.md)
- [`intelligence/engineering/agent-architecture/rule-overload.md`](../../intelligence/engineering/agent-architecture/rule-overload.md)

本文件把 context collapse 與 attention budgeting 的 agent architecture intelligence 轉譯成 AI runtime context governance。原始思想回答「為什麼長 session 會遺失早期決策、為什麼無效讀取會稀釋注意力」；本文件定義 context loading、recap、prune 與 task boundary 需要通過的治理 gate。

## 觸發時機

在下列情況套用本治理：

- 任務需要讀取多個 workflow、governance、runtime、analysis 或 intelligence 文件。
- Session 已跨多個 task、commit、tool call 或 user redirect。
- Agent 開始重複讀檔、重複工具呼叫、忘記早期決策或改變已確認 scope。
- 修改 context loading、TTL、summary-first、prompt-cache 或 runtime routing surface。

## Runtime Gate

| Gate | 通過條件 |
| --- | --- |
| Summary-first loading | 先讀 index / summary / route，再展開 full source；不得一開始 full-context 掃描。 |
| Attention budget cap | 任務開始前界定必要讀取集合；新增讀取必須能解釋要降低的 uncertainty。 |
| Decision externalization | 影響後續行動的 scope、route、commit boundary 或 user constraint 必須寫入 plan、goal、todo 或 durable doc。 |
| Recap checkpoint | 跨 task、長工具鏈、commit 前或 user redirect 後，重新確認 goal、dirty state、route 與 validation。 |
| Task-boundary prune | 任務完成後，不把已完成 task 的 context 當成下一個任務的 silent source。 |
| Escalation on collapse | 出現重複、矛盾、忘記 user constraint 或工具試錯膨脹時，停止擴張 context，重建 minimal route。 |

## 分層判斷

| 內容類型 | 目標層 |
| --- | --- |
| 為什麼 agent 會 context collapse、注意力被無效讀取稀釋 | `intelligence/engineering/agent-architecture/` |
| Context loading、recap、prune、attention budget 的治理 gate | `governance/ai-runtime-governance/` |
| TTL policy、prompt cache layout、generated runtime surface | `runtime/context/` |
| 具體路由、summary、compression 與 model context view | `knowledge/runtime/`、`models/` |
| 可測 long-session drift、stale-context 或 repeated-tool-loop | `validation/` 或 `anti-patterns/` |

## Runtime Mapping

- [`governance/lifecycle/context-ttl-philosophy.md`](../lifecycle/context-ttl-philosophy.md) — context TTL 與 prune 的 design layer。
- [`runtime/context/ttl-policy.yaml`](../../runtime/context/ttl-policy.yaml) — machine-readable TTL policy。
- [`runtime/context/prompt-cache-playbook.md`](../../runtime/context/prompt-cache-playbook.md) — provider prompt cache layout。
- [`knowledge/runtime/routing-registry.yaml`](../../knowledge/runtime/routing-registry.yaml) — route context cost、TTL 與 compression metadata。

## Validation Candidate

後續若要 promotion 到 `validation/`，可建立 scenario 檢查：

- Agent 在長 session 後忘記早期 user constraint。
- Agent 未先讀 summary/index 就展開大量 full source。
- Agent 在 task boundary 後沿用 stale context，導致錯誤 route 或錯誤 commit boundary。
- Agent 出現重複工具呼叫但未 recap / prune / rebuild minimal route。
