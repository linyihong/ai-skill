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

## 錯誤查詢索引

當遇到以下情境時，可快速查詢對應的 ADR 或 failure pattern：

| 情境 | 查什麼 | 預期找到 |
|------|--------|----------|
| 舊 skill 內容要搬到新分層 | [`ADR-001`](ADR-001-reference-first-migration-strategy.md) | Reference-first migration 策略 |
| Intelligence 與 Knowledge 分不清楚 | [`ADR-002`](ADR-002-intelligence-vs-knowledge-separation.md) | 分離原則與邊界 |
| Analysis / Workflow / Intelligence 三層如何分工 | [`ADR-003`](ADR-003-three-layer-architecture.md) | 三層架構定義 |
| Feedback lesson 如何提升為 reusable rule | [`ADR-004`](ADR-004-feedback-promotion-pipeline.md) | Promotion pipeline |
| 記憶模型如何分層 | [`ADR-005`](ADR-005-memory-architecture.md) | 6 子層記憶模型 |
| Agent 犯了重複錯誤 | [`shared-rules/failure-patterns/README.md`](../shared-rules/failure-patterns/README.md) | 對應的 failure pattern 與 prevention gate |
| 某個架構決策需要修改 | 建立新 ADR（`ADR-006-{title}.md`）並標記舊 ADR 為 superseded | 新 ADR 記錄變更理由 |
| Session-level 的輕量決策 | [`memory/decision/`](../memory/decision/) | 跨 session 但非架構級的決策記錄 |

## 與 Failure Patterns 的關係

本層 ADR 記錄「正確的架構決策」，而 [`shared-rules/failure-patterns/`](../shared-rules/failure-patterns/README.md) 記錄「agent 常犯的錯誤模式」。兩者互補：

- ADR 告訴 agent **應該怎麼做**
- Failure pattern 告訴 agent **不要怎麼做**
- 新增 ADR 時，應檢查是否有對應的 failure pattern 需要更新
- 新增 failure pattern 時，應檢查是否有相關的 ADR 可以引用

## 規則

1. **Immutable**：Decision 一旦 accepted 就不修改；需要變更時建立新的 ADR 並標記舊的為 superseded。
2. **Numbered**：ADR 使用流水號（ADR-001, ADR-002, ...）。
3. **Linked**：每個 decision 連結到相關的 source files 或 decisions。
4. **Minimal**：每個 ADR 不超過 1000 tokens。需要詳細技術分析時引用外部文件。

## 誰會參考這裡（Inbound References）

- [`route.decisions.adr`](../knowledge/runtime/routing-registry.yaml:687) — primary_source 為 `decisions/README.md`
- [`route.architecture.permanent-docs`](../knowledge/runtime/routing-registry.yaml:723) — candidate_sources 引用 `decisions/ADR-001`、`decisions/ADR-003`

## 與既有層的關係

- `memory/decision/`：輕量版決策記錄（session-level）
- `intelligence/`：engineering intelligence 可引用 ADR
- `architecture/`：架構規劃文件可引用 ADR
