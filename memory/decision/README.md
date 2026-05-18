# Decision Memory

`memory/decision/` 保存**跨 session 的重要決策記錄**。這是 Architecture Decision Records（ADR）的輕量版，讓 agent 記住「為什麼這樣做」。

## 用途

- 記錄架構決策與取捨
- 避免重複討論同一決策
- 提供決策的上下文與理由
- 支援 future agent 理解歷史脈絡

## 格式

```markdown
# ADR-{number}: {title}

## Status
{proposed | accepted | deprecated | superseded}

## Context
{為什麼需要這個決策}

## Decision
{我們決定了什麼}

## Consequences
{正面與負面影響}

## Alternatives Considered
- {alternative 1}：{為什麼不選}
- {alternative 2}：{為什麼不選}

## Related
- {related decision 1}
- {related file 1}
```

## 規則

1. **Immutable**：Decision 一旦 accepted 就不修改；需要變更時建立新的 ADR 並標記舊的為 superseded。
2. **Numbered**：ADR 使用流水號（ADR-001, ADR-002, ...）。
3. **Linked**：每個 decision 連結到相關的 source files 或 decisions。
4. **Searchable**：使用標準 front matter 格式。

## 索引（session 級）

| 日期 | 檔案 | 摘要 |
| --- | --- | --- |
| 2026-05-14 | [2026-05-14_architecture-plans-boundary.md](./2026-05-14_architecture-plans-boundary.md) | architecture/ vs plans/ 邊界 |
| 2026-05-18 | [2026-05-18_registry-first-and-decision-runtime.md](./2026-05-18_registry-first-and-decision-runtime.md) | ADR-006 + runtime decision-recording |

## 與既有層的關係

- `decisions/`：正式的 ADR 目錄（長期、跨專案的決策）
- `runtime/decisions/`：machine-readable 決策寫入路由（close-loop）
- `memory/summary/`：session summary 可連結到 decisions
- `intelligence/`：engineering intelligence 可引用 decisions
