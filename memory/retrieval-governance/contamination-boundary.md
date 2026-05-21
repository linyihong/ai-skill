# Contamination Boundary

Memory replay 可能把 stale frame 帶回 active execution。Contamination boundary 定義哪些 memory 只能作 weak hint，哪些必須重新驗證，哪些不得 replay。

## Boundary Classes

| Boundary | 意義 | Replay rule |
| --- | --- | --- |
| `workflow-local` | 僅適用同一 workflow family。 | 可作 checklist hint，需 current source check。 |
| `domain-local` | 僅適用同一 domain / architecture family。 | 不可跨 domain 自動 replay。 |
| `project-local` | 僅適用同一 repo / project。 | Repo refactor / migration 後需 revalidation。 |
| `session-global` | 會影響整個 session frame。 | 需 recap / prune / human alignment。 |

## Forbidden Replay

- Replay stale blockers as active blockers。
- Replay old `.agent-goals/` state as current execution contract。
- Replay deprecated architecture frame without compatibility check。
- Replay old execution graph without revalidation。
- Replay memory-derived conclusion as canonical source。
- Replay project-private evidence into reusable docs。
- Replay workaround as permanent policy without abstraction。

## Contamination Response

若 replay 後發現 contamination：

1. Stop current patch / claim。
2. Mark replay as weak or invalidated。
3. Reread current source-of-truth。
4. Rebuild execution graph。
5. Downgrade autonomy if needed。
6. Discard or quarantine the memory candidate。

## Safe Usage

Memory 最安全的用途是提示「可能要檢查什麼」，而不是直接回答「現在是什麼」。任何會影響 commit、runtime、rules、user-facing completion 的 replay，都必須重新驗證 current source。
