# Working Memory

`memory/working/` 保存**目前 session 的進行中狀態**。這是短期記憶，session 結束後可清除或存檔到 `memory/summary/`。

## 用途

- 目前 task 的進行中 context
- 尚未完成的子任務狀態
- 暫存的分析中間結果
- Session archive（context pollution 或 hard stop 時自動存檔）

## 規則

1. **Session-local**：Working memory 只存活於目前 session。
2. **可丟棄**：Working memory 的內容可安全丟棄，不影響長期知識。
3. **Session 結束時**：重要內容應 promotion 到 `memory/summary/` 或 `knowledge/`。
4. **Context pollution / hard stop 時**：自動存檔到 `memory/working/session-archive-{timestamp}.md`。

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

## 與既有層的關係

- `memory/summary/`：session 結束後 promotion 到 summary memory
- `memory/decision/`：重要決策 promotion 到 decision memory
- `runtime/guards/context-pollution.yaml`：自動存檔觸發條件
