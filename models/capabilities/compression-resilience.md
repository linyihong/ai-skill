# Compression Resilience

Compression resilience 描述 model profile 是否能安全使用 summaries、checklists、generated reports 或 indexes，而不遺失 required gates。

## Levels

| Level | 意義 | Allowed compression |
| --- | --- | --- |
| `low` | 容易遺失 required dependencies，或過度泛化 compressed context。 | 使用 `source-backed`。 |
| `medium` | 可用 summaries orient，但 edits 前需要 source。 | 使用 `summary-first`，再 escalate。 |
| `high` | 可對 low-risk bounded work 使用 checklist-first execution。 | 不需 source edits 時使用 `checklist-first`。 |
| `unknown` | 沒有 evidence。 | Material work 從 source-backed 開始。 |

## Escalation Triggers

- 任務修改 canonical docs、rules、runtime source、generated artifacts 或 plans。
- Generated report 可能 stale。
- Summary 與 source 衝突。
- 工作需要 commit / push / readback。

## Validation

Compressed route 只有在記錄已讀內容、deferred sources，以及哪些 trigger 會強制 source-backed loading 時才有效。
