# Working Memory

`memory/working/` 保存**目前 session 的 cognition buffer**。這是 session-scoped、semi-stable、discardable 的 activated context frame，session 結束後可清除，或經 compression / qualification 後 promotion 到 `memory/summary/`。

## 用途

- 目前 task 的 activated assumptions、recent evidence 與 temporary workflow context
- 經 `memory/retrieval-governance/` qualification 後的 activated memory refs
- 暫存的 risk assessment 與 current architecture frame
- Session archive（context pollution 或 hard stop 時自動存檔）

## 規則

1. **Session-local**：Working memory 只存活於目前 session。
2. **可丟棄**：Working memory 的內容可安全丟棄，不影響長期知識。
3. **不是 goal ledger**：不得保存 owner、lock、next action 或 active blocker；這些屬於 `.agent-goals/`。
4. **不是 runtime-state**：不得保存 persistent execution state 或 machine guard state。
5. **Session 結束時**：重要內容先 compression 到 `memory/summary/`，再依 policy promotion。
6. **Context pollution / hard stop 時**：自動存檔到 `memory/working/session-archive-{timestamp}.md`，且後續 replay 必須重新 qualification。

## 格式

```markdown
# Session Archive: {timestamp}

## Task
{current task description}

## Progress
- [x] Step 1
- [ ] Step 2

## Key Decisions
- {decision 1}
- {decision 2}

## Pending
- {pending item 1}
```

## Buffer 欄位

```yaml
active_assumptions:
recent_evidence:
current_architecture_frame:
active_repo_topology:
temporary_workflow_context:
current_risk_assessment:
activated_memory_refs:
discard_after:
```

## 不保存

- Current blocker as durable truth。
- Next action、owner、lock。
- Long-term project state。
- Canonical decision。
- Runtime execution state。
- 未經 qualification 的 memory-derived conclusion。

## 與既有層的關係

- `memory/summary/`：session 結束後 promotion 到 summary memory
- `memory/decision/`：重要決策 promotion 到 decision memory
- `memory/retrieval-governance/`：決定哪些 memory refs 可以 activation 到 working buffer
- `.agent-goals/`：active execution contract，不由 working memory 取代
- `runtime/guards/context-pollution.yaml`：自動存檔觸發條件
