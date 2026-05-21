# Activation Thresholds

Memory replay 不是 default context loading。只有目前 task 真的需要 historical replay，且比重新讀 canonical source 更有價值時，才 activation。

## Activation Levels

| Level | 意義 | 允許行為 |
| --- | --- | --- |
| `none` | 沒有足夠 trigger。 | 不讀 memory；改讀 current source / routing registry。 |
| `weak_hint` | Memory 可能提示風險或方向。 | 只能影響 discovery，不得支撐 conclusion。 |
| `scoped_context` | Memory 與目前 repo / workflow / failure class 匹配。 | 可進入 `memory/working/`，但需 current source revalidation。 |
| `decision_reference` | Accepted decision 或 active project memory 仍相容。 | 可作 decision context，但仍需檢查 supersession。 |
| `recovery_frame` | Context compaction、contamination 或 repeated failure 需要 replay。 | 只 replay 最小摘要，並重建 execution graph。 |

## Trigger Examples

- Repeated failure class 再次出現。
- 同 repo / project 的 architecture boundary 需要恢復。
- 使用者要求「接續上一段」但 current source 不足。
- Context compaction 後，summary 可補足目標與 validation target。
- Active decision 需要查 accepted / superseded 狀態。
- Stale assumption suspicion 需要找過去失敗模式當 weak hint。

## 不得 Activation

- 只是一般 lookup，current source 可以直接回答。
- Memory 來源比 canonical source 更舊且成本更高。
- Memory 會把 old blocker 當成 active blocker。
- Memory 會覆蓋使用者最新目標。
- Memory 會讓 agent 跳過 source-of-truth gate。

## Gate

Activation 前至少確認：

```text
Why memory now:
Why current source alone is insufficient:
Replay level:
Qualification source:
Current-source revalidation target:
```
