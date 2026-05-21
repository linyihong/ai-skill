# Summary Memory

`memory/summary/` 保存**已完成的 session 摘要**。這是 compressed 歷史記憶，讓 agent 在後續 session 中快速回顧過去做了什麼；summary 只恢復 session context，不證明 current truth。

## 用途

- 記錄每個 session 的目標、成果與關鍵決策
- 跨 session 的上下文銜接
- 減少重複工作
- Context compaction / handoff 後的最小 replay source

## 格式

```markdown
# Session Summary: {date}

## Goals
- {goal 1}
- {goal 2}

## Achievements
- {achievement 1}
- {achievement 2}

## Key Decisions
- {decision 1}（→ decisions/{decision-file}）
- {decision 2}

## Files Modified
- {file 1}
- {file 2}

## Token Usage
- Total: {tokens}
- By layer: {bootstrap: 800, skill: 3500, ...}

## Next Steps
- {next step 1}
- {next step 2}
```

## 規則

1. **Session boundary**：每個 session 結束時建立一個 summary。
2. **Compressed**：每個 summary 不超過 500 tokens。
3. **Link to decisions**：重要決策連結到 `memory/decision/` 或 `decisions/`。
4. **Searchable**：Summary 使用標準格式，支援全文搜尋。
5. **Not current truth**：Summary 不得取代 current source reading。
6. **Replay budget**：Summary replay 預設只讀 latest relevant summary；不足時回到 current source 或詢問使用者，不追溯 full transcript chain。
7. **Promotion gate**：Summary 中的 reusable lesson 必須先抽象化，才可 promotion 到 `knowledge/`、`intelligence/`、`workflow/` 或 `enforcement/`。

## 與既有層的關係

- `memory/working/`：session archive 的 promotion 目標
- `memory/decision/`：重要決策的持久化位置
- `memory/retrieval-governance/`：定義 summary replay 的 activation threshold 與 budget
- `knowledge/summaries/`：knowledge-level summaries（不同於 session summaries）
