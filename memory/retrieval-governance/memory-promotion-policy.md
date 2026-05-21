# Memory Promotion Policy

Memory promotion 將 temporary session context 逐步壓縮、抽象化，最後只把可重用內容提升到更穩定 layer。Promotion 不得把 raw transcript、private evidence 或 active execution state 長期化。

## Pipeline

```text
memory/working/
→ compress
memory/summary/
→ select episode / project / decision / failure
memory/episodic/ | memory/project/ | memory/decision/ | memory/failure/
→ abstract / generalize
intelligence/ | knowledge/ | enforcement/ | workflow/
```

## Promotion Criteria

| Condition | Required |
| --- | --- |
| Reusable | yes |
| Generalized | yes |
| Non-project-secret | yes |
| Low contamination risk | yes |
| Repeated utility | preferred |
| Source compatibility known | yes |

## Forbidden Promotion

- Raw transcript。
- Temporary blocker。
- Active runtime assumption。
- Unstable execution graph。
- Unresolved contradiction。
- Project-secret / private evidence。
- Old `.agent-goals/` owner、next action、lock 或 active blocker。

## Destination Rules

| Destination | Use when |
| --- | --- |
| `memory/summary/` | Session outcome 需要 compressed recovery。 |
| `memory/episodic/` | 特定情境未來可能 replay，但尚未抽象化。 |
| `memory/project/` | Same project / repo 的 durable context 仍有用。 |
| `memory/decision/` | 決策已 accepted，且需要 status / supersession tracking。 |
| `memory/failure/` | 失效模式需要長期追蹤與 prevention evolution。 |
| `knowledge/` | 內容已抽象成 reusable navigation / summary / graph。 |
| `intelligence/` | 內容已抽象成可重用 reasoning / heuristic。 |
| `enforcement/` | 已有 recurring validated failure，且需要可執行 policy。 |

## Pruning

Memory 應在 superseded、expired、scope mismatch 或 contamination risk 過高時降級、封存或刪除。Pruning 後若仍需保留經驗，只保留 generalized lesson 與 source compatibility note。
