# Architecture Decision Records

`constitution/` 保存**正式的 Architecture Decision Records（ADR）**。這些是跨 session、跨專案的重要架構決策，讓 agent 永遠記得「為什麼這樣做」。

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
accepted → deprecated → superseded
              ↓
        (不再修改)
```

**No proposed ADRs**：constitution/ 只放 **accepted** ADRs。架構決策的提案、討論、alternatives 評估在對應 `plans/active/` 的 `Decision Rationale` section 完成；plan completed 且通過 `ADR Promotion Criteria` 後才升級為 accepted ADR。此規則於 2026-05-22 確立，理由：

- ADR-007 §Decision 已明文「ADR is NOT the default endpoint」
- Proposed ADR 是「未驗證的憲法」— 與 constitution 性質矛盾
- 改用 plan-first / ADR-after-completion 避免「廢棄憲法」累積
- 完整規則見 [`governance/lifecycle/decision-promotion-pipeline.md`](../governance/lifecycle/decision-promotion-pipeline.md) §No-Proposed-ADR Rule

## 現有 ADR

| ADR | Title | Status | Date | Framework Generation |
| --- | --- | --- | --- | --- |
| [ADR-001](ADR-001-reference-first-migration-strategy.md) | Reference-First Migration Strategy | accepted | 2026-05-12 | cross-generation（Gen 2 起延伸至 Gen 3；skills/ 路徑已搬遷） |
| [ADR-002](ADR-002-intelligence-vs-knowledge-separation.md) | Intelligence vs Knowledge Separation | accepted | 2026-05-12 | Gen 3 確立 |
| [ADR-003](ADR-003-three-layer-architecture.md) | Three-Layer Architecture（Knowledge / Skills / Intelligence） | accepted | 2026-05-12 | **Gen 2 / Gen 3 transition**（核心精神保留，Skills 已演化為 workflow/+analysis/，見 ADR-003 §Vocabulary Evolution） |
| [ADR-004](ADR-004-feedback-promotion-pipeline.md) | Feedback Promotion Pipeline | accepted | 2026-05-12 | cross-generation（pipeline 結構保留，路徑詞彙演化，見 ADR-004 §Vocabulary Evolution） |
| [ADR-005](ADR-005-memory-architecture.md) | Memory Architecture（6 子層記憶模型） | accepted | 2026-05-12 | Gen 3 確立 |
| [ADR-006](ADR-006-registry-first-workflow-activation.md) | Registry-First Workflow Activation | accepted | 2026-05-18 | Gen 3 |
| [ADR-007](ADR-007-constitution-and-decision-promotion-boundary.md) | Constitution and Decision Promotion Boundary | accepted | 2026-05-21 | Gen 3 |

每個 ADR 內含 `Framework Generation` section，標註世代分類與當前世代文件連結。新增世代時依 [`governance/lifecycle/system-upgrade-governance.md`](../governance/lifecycle/system-upgrade-governance.md) §3 規則 7 評估每個 ADR 的延伸狀態。

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
| Agent 犯了重複錯誤 | [`enforcement/failure-patterns/README.md`](../enforcement/failure-patterns/README.md) | 對應的 failure pattern 與 prevention gate |
| 某個架構決策需要修改 | 建立新 ADR（下一個流水號）並標記舊 ADR 為 superseded | 新 ADR 記錄變更理由 |
| Session-level 的輕量決策 | [`memory/decision/`](../memory/decision/) | 跨 session 但非架構級的決策記錄 |
| Agent 跳過 workflow discovery 就寫碼 | [ADR-006](ADR-006-registry-first-workflow-activation.md)、[`workflow/workflow-routing.md`](../workflow/workflow-routing.md) | registry-first + #27 閘門 |
| 每加 workflow 就加 activation # | [ADR-006](ADR-006-registry-first-workflow-activation.md)、[`runtime/runtime.db`](../runtime/runtime.db) | 觸發條件在 registry |
| SDK Page／cache／catalog 為何這樣設計 | `<PROJECT_ROOT>/docs/decisions/`（如 apk-analysis-sdk） | 專案 tier；見 runtime decision-recording |
| `decisions/` 與 `memory/decision/` 容易混淆 | [ADR-007](ADR-007-constitution-and-decision-promotion-boundary.md) | formal ADR layer 改名為 `constitution/`，decision promotion 依內容選 target |

## 與 Failure Patterns 的關係

本層 ADR 記錄「正確的架構決策」，而 [`enforcement/failure-patterns/`](../enforcement/failure-patterns/README.md) 記錄「agent 常犯的錯誤模式」。兩者互補：

- ADR 告訴 agent **應該怎麼做**
- Failure pattern 告訴 agent **不要怎麼做**
- 新增 ADR 時，應檢查是否有對應的 failure pattern 需要更新
- 新增 failure pattern 時，應檢查是否有相關的 ADR 可以引用

## 規則

1. **Immutable**：Decision 一旦 accepted 就不修改；需要變更時建立新的 ADR 並標記舊的為 superseded。
2. **Numbered**：ADR 使用流水號（ADR-001, ADR-002, ...）。
3. **Linked**：每個 decision 連結到相關的 source files 或 constitution records。
4. **Minimal**：每個 ADR 不超過 1000 tokens。需要詳細技術分析時引用外部文件。

## Runtime 決策紀錄（close-loop）

閉環時 agent 應讀 [`runtime/runtime.db`](../runtime/runtime.db)：

- **architecture** → 本目錄 `ADR-{n}-*.md`
- **session** → `memory/decision/{date}_*.md`
- **project** → `<PROJECT_ROOT>/docs/decisions/`

## Decision Promotion Target

Decision promotion 不一定以 ADR 為終點；依內容選 target：

| Decision content | Target |
| --- | --- |
| 可執行規則 / cross-agent policy | `enforcement/` |
| reasoning heuristic / tradeoff / signal / anti-pattern | `intelligence/` |
| 操作流程 / repeatable work sequence | `workflow/` |
| runtime gate / activation / obligation / executable contract projection | `runtime/runtime.db` |
| 架構級不可逆、跨 session / project 的 foundational decision | `constitution/ADR-*` |
| session-scoped replay decision | `memory/decision/` |
| project-specific decision | `<PROJECT_ROOT>/docs/decisions/` |

詳細 promotion gate 見 [`governance/lifecycle/decision-promotion-pipeline.md`](../governance/lifecycle/decision-promotion-pipeline.md)。

## 誰會參考這裡（Inbound References）

- [`route.runtime.decision-recording`](../knowledge/runtime/routing-registry.yaml) — 何時寫入哪一 tier
- [`route.constitution.adr`](../knowledge/runtime/routing-registry.yaml) — primary_source 為 `constitution/README.md`
- [`route.architecture.permanent-docs`](../knowledge/runtime/routing-registry.yaml:723) — candidate_sources 引用 `constitution/ADR-001`、`constitution/ADR-003`

## 與既有層的關係

- `memory/decision/`：輕量版決策記錄（session-level）
- `intelligence/`：engineering intelligence 可引用 ADR
- `architecture/`：架構規劃文件可引用 ADR
