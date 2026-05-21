# Replay Budget

Replay budget 定義 memory replay 的最大深度與 token cost 邊界。Budget 不是硬性 token 計算器，而是防止 historical context 無限制展開的治理規則。

## Budget Defaults

| Situation | Default budget | Max depth |
| --- | --- | --- |
| Lookup current decision status | Decision index + one decision record | 1 hop |
| Same project context recovery | Project memory summary + current source | 1-2 hops |
| Context compaction recovery | Latest session summary only | 1 hop |
| Repeated failure suspicion | Failure memory summary + one linked episode if needed | 1-2 hops |
| Recovery after contamination | Minimal summary + source-of-truth reload | 1 hop before validation |
| User-requested full recap | Summary first, ask before transcript-level replay | gated |

## Depth Rules

- 第一層讀 memory type README 或 index。
- 第二層只讀一個最匹配 candidate。
- 第三層以上需要明確 validation target 或 user request。
- Transcript-level replay 預設 blocked。
- 每次 replay 後都要決定 discard、keep in working buffer、或 promotion candidate。

## Budget Record

```text
Replay purpose:
Max memory files:
Max hops:
Expected current-source validation:
Stop condition:
```

## Recursion Ban

禁止為了理解 memory 而不斷打開更舊 memory。若 summary 不足，優先回到 current source 或詢問使用者，而不是追溯 full transcript chain。
