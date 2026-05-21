# Context Stability

Context stability 描述 execution profile 是否能在 long tasks 中保存 goals、evidence、constraints 與 open questions。

## Signals

| Signal | 意義 | 必要行動 |
| --- | --- | --- |
| `stable` | Goal、source paths 與 validation target 仍 aligned。 | 以 selected strategy 繼續。 |
| `strained` | Context 很長或已 compressed，但 evidence 仍可恢復。 | Edits 前 recap assumptions 並 reread primary source。 |
| `unstable` | Prior frame 可能 contaminate current task。 | 使用 rediscovery-only behavior。 |
| `unknown` | 無法評估 stability。 | 降低 autonomy，使用 source-backed validation。 |

## Instability Triggers

- Context compaction 保留 conclusions，但沒有 evidence。
- Previous route、checklist 或 memory 跨 task boundaries 重用。
- Current action 不再服務 user goal。
- 多個 active plans 使用相似詞彙，但 source-of-truth 不同。

## Validation

繼續 long 或 compressed task 前，先確認 current goal、active source 與 validation target。任一項缺失時，使用 `goal-realignment` 或 `rediscovery-only`。
