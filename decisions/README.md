# Architecture Decision Records

`decisions/` 保存**正式的 Architecture Decision Records（ADR）**。這些是跨 session、跨專案的重要架構決策，讓 agent 永遠記得「為什麼這樣做」。

## 為什麼需要 ADR

AI 最大的問題之一是：

> 以前做過的決策忘記了。

ADR 解決這個問題：

- 每個重要決策都有記錄
- 新 session 的 agent 可以快速回顧歷史決策
- 避免重複討論同一決策
- 提供決策的上下文、取捨與替代方案

## ADR 生命週期

```text
proposed → accepted → deprecated → superseded
                ↓
          (不再修改)
```

## 現有 ADR

| ADR | Title | Status | Date |
| --- | --- | --- | --- |
| [ADR-001](ADR-001-reference-first-migration-strategy.md) | Reference-First Migration Strategy | accepted | 2026-05-12 |
| [ADR-002](ADR-002-intelligence-vs-knowledge-separation.md) | Intelligence vs Knowledge Separation | accepted | 2026-05-12 |
| [ADR-003](ADR-003-three-layer-architecture.md) | Three-Layer Architecture（Knowledge / Skills / Intelligence） | accepted | 2026-05-12 |
| [ADR-004](ADR-004-feedback-promotion-pipeline.md) | Feedback Promotion Pipeline | accepted | 2026-05-12 |
| [ADR-005](ADR-005-memory-architecture.md) | Memory Architecture（6 子層記憶模型） | accepted | 2026-05-12 |

## 格式

每個 ADR 使用 `ADR-{number}-{short-title}.md` 命名：

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
4. **Minimal**：每個 ADR 不超過 1000 tokens。需要詳細技術分析時引用外部文件。

## 與既有層的關係

- `memory/decision/`：輕量版決策記錄（session-level）
- `intelligence/`：engineering intelligence 可引用 ADR
- `architecture/`：架構規劃文件可引用 ADR
