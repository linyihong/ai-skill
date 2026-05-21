# Replay Cost Governance

Replay cost governance 用來避免 agent 因為「想回想」而載入過多歷史 context。Memory replay 的成本必須低於或等於它帶來的風險降低與目標推進價值。

## Cost Classes

| Memory type | Cost | Default |
| --- | --- | --- |
| Failure pattern memory | low | 可頻繁使用，但必須 scoped。 |
| Decision memory | low | Status check 後可用。 |
| Summary memory | medium | Conditional，用於 handoff / context recovery。 |
| Episodic memory | medium | Weak guidance only。 |
| Project memory | medium-high | 限 same project / repo / architecture boundary。 |
| Old execution recap | high | On-demand only。 |
| Full transcript / full session replay | very high | Avoid。 |

## 最小回放原則

- Prefer smallest sufficient replay scope。
- Replay depth 不得超過 current task 需要的最小 evidence。
- 若 replay 成本高於重新讀 canonical source，優先讀 canonical source。
- Replay 不得形成 recap recursion：summary → old summary → old transcript → older summary。
- Replay 不得替代 tests、lints、runtime validation 或 current source read。

## Cost Escalation

Replay cost 升高時，需要更強 qualification：

| Cost | Qualification |
| --- | --- |
| low | Trigger + scope match。 |
| medium | Trigger + scope match + freshness check。 |
| high | User goal relevance + current source insufficiency + explicit replay budget。 |
| very high | 使用者明確要求，或 recovery 沒有其他 sufficient source。 |

## Stop Rule

當 memory replay 已回答 routing decision 或 validation target，停止讀更多 memory；接著回到 current source、execution 或 user alignment。
