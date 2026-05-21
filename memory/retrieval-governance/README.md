# Memory Retrieval Governance

`memory/retrieval-governance/` 定義 memory 如何被選擇性 retrieval、qualification、activation、replay、discard 或 promotion。Memory 預設 dormant；只有符合 trigger、qualification 與 replay budget 時，才可進入 `memory/working/` 的 session cognition buffer。

## 邊界

| Layer | 責任 | 不得成為 |
| --- | --- | --- |
| `memory/` | Historical replay archive。 | Canonical source、active execution state、永久 context dump。 |
| `memory/retrieval-governance/` | Retrieval trigger、activation threshold、replay budget、freshness、contamination、promotion policy。 | Runtime state machine、evidence hierarchy replacement。 |
| `memory/working/` | Session cognition buffer，保存已 qualification 的 activated context。 | `.agent-goals/`、owner / next action / lock、runtime-state。 |
| `.agent-goals/` | Active execution contract。 | Long-term memory 或 session summary。 |
| `knowledge/` | Structured reusable navigation / lookup。 | Historical replay 或 transcript recap。 |
| `runtime/` | Deterministic executable lookup / state。 | Raw historical memory。 |

## Activation Pipeline

```text
trigger
→ retrieval
→ qualification
→ replay budget check
→ activation
→ memory/working/ buffer
→ execution usage
→ revalidation / discard / promotion
```

## 入口文件

- [`activation-thresholds.md`](activation-thresholds.md)：何時允許 replay memory。
- [`retrieval-routing.md`](retrieval-routing.md)：trigger signal 到 memory type 的 routing。
- [`replay-cost-governance.md`](replay-cost-governance.md)：replay 成本與最小回放原則。
- [`replay-budget.md`](replay-budget.md)：replay depth、token budget 與 recursion ban。
- [`freshness-and-decay.md`](freshness-and-decay.md)：freshness、scope、confidence defaults。
- [`contamination-boundary.md`](contamination-boundary.md)：contamination classes 與 forbidden replay。
- [`memory-promotion-policy.md`](memory-promotion-policy.md)：working → summary → long-term memory → reusable knowledge 的 promotion policy。

## 基本規則

1. 未通過 qualification 的 memory 不得進入 active execution frame。
2. Episodic memory 預設只能作 weak guidance。
3. Project memory 只在同 repo / 同 architecture boundary / 同 workflow family 下使用。
4. Decision memory 必須檢查 status、supersession 與 compatibility scope。
5. Summary memory 只能恢復 session context，不證明 current truth。
6. Full session replay 預設禁止，除非使用者明確要求或沒有其他足夠 source。
7. Memory-derived conclusion 不得取代 canonical source。

## Decision Record

每次非 trivial memory activation 應記錄：

```text
Trigger:
Candidate memory:
Qualification:
Replay budget:
Activation target:
Current source revalidation:
Discard / promotion decision:
```
