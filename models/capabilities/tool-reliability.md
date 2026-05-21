# Tool Reliability

Tool reliability 描述 execution profile 是否能安全協調 tool calls、file edits、validation、git operations 與 handoffs。

## Reliability Classes

| Class | 意義 | Execution impact |
| --- | --- | --- |
| `single-step` | 適合一次 read、一次 edit 或一次 validation。 | 保持 tasks bounded，避免 multi-tool loops。 |
| `multi-step` | 可協調多次 reads、edits 與 checks。 | 使用 source-backed strategy 與 explicit validation。 |
| `close-loop-capable` | 可安全處理 validation、commit、push、readback 與 clean status。 | 只有 repository state safe 且 gates 已知時使用。 |
| `unknown` | 沒有 reliable tool orchestration evidence。 | 避免 automation；改詢問或用 dry-run。 |

## Required Gates

對 commit / push / generated runtime work，tool reliability 必須包含：

- Git availability 與 safe repo state。
- Diff review。
- Relevant lints / tests / runtime validation。
- Commit / push 後 readback。
- Clean `git status --short --branch`。

## Handoff Boundary

Subagents 可以 inspect 或 analyze，但除非使用者明確委派，parent agent 負責 shared file edits、commits、pushes 與 final claims。
