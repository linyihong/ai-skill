# Agent Architecture Intelligence

放 **AI Agent 自身運作的智慧**。這裡收集從實際使用中累積的、關於 AI 如何思考、如何決策、如何失敗、如何恢復的經驗。

## 核心

AI 不是黑盒子 — 它的行為模式可以被觀察、分類、預測、改善。

## 目前 atoms

| Atom | 原則 | 狀態 | 來源 |
|------|------|------|------|
| [`context-collapse.md`](context-collapse.md) | When context exceeds effective window, agent loses earlier decisions and repeats or contradicts them. | `candidate-intelligence` | 本系統實際運作觀察 |
| [`rule-overload.md`](rule-overload.md) | When too many rules compete for attention, agent follows the most recently loaded or most concrete rule, not the most important one. | `candidate-intelligence` | 本系統實際運作觀察 |
| [`task-routing.md`](task-routing.md) | Agent routing decisions are determined by signal strength, not signal correctness. | `candidate-intelligence` | 本系統實際運作觀察 |
| [`attention-budgeting.md`](attention-budgeting.md) | Agent has finite attention per session; every unnecessary file read or tool call consumes budget that could be spent on reasoning. | `candidate-intelligence` | 本系統實際運作觀察 |
| [`failure-recovery.md`](failure-recovery.md) | Agent's first recovery attempt after failure is usually the most reliable; repeated retries without strategy change degrade output quality. | `candidate-intelligence` | 本系統實際運作觀察 |
| [`cognitive-boundaries.md`](cognitive-boundaries.md) | Agent cannot reliably detect its own cognitive boundaries; external gates (validation, checklists, failure patterns) are required. | `candidate-intelligence` | 本系統實際運作觀察 |
| [`pilot-first-validation.md`](pilot-first-validation.md) | Before abstracting a pattern into a generalized pipeline, validate it with a concrete pilot first. | `candidate-intelligence` | Phase 28-31 實際運作經驗 |
| [`failure-to-scenario-closure.md`](failure-to-scenario-closure.md) | Every AI-system failure that can be reproduced statelessly should be converted into a validation scenario. | `candidate-intelligence` | 本系統實際運作經驗 |
| [`linked-updates-completeness.md`](linked-updates-completeness.md) | When modifying one file, agent must identify and update all files that reference or depend on it. | `candidate-intelligence` | Phase 33-35 實際運作經驗 |
| [`decomposition-strategy-selection.md`](decomposition-strategy-selection.md) | Not all content benefits from the same extraction strategy; choose based on content structure. | `candidate-intelligence` | Phase 26-33 實際運作經驗 |
| [`stateless-validation-necessity.md`](stateless-validation-necessity.md) | AI decision path validation must be stateless — it cannot depend on conversation memory. | `candidate-intelligence` | 本系統實際運作經驗 |
| [`index-first-documentation.md`](index-first-documentation.md) | Documentation that agents repeatedly read should expose routeable indexes before detailed prose. | `candidate-intelligence` | documentation workflow extraction |
| [`multi-layer-enforcement.md`](multi-layer-enforcement.md) | 可靠強制 agent obligation 需三層協同：context 注入（資料精確性）、禁止 prompt（意識觸發）、機械關卡（兜底）— 每層防止不同失效模式，單層均不足。 | `candidate-intelligence` | Bootstrap 三層架構 empirical 測試 2026-05-27 |
| [`wish-to-task-list-translation.md`](wish-to-task-list-translation.md) | 當 user prompt 是 wish（模糊願望）時，agent 必須先翻譯為 measure / process / verify 三段任務單再執行，否則進入 shotgun mode。 | `candidate-intelligence` | 外部 infographic「Prompt 與其寫成願望，不如寫成任務單」沉澱 2026-05-27 |

## 與其他層的關係

- `enforcement/failure-patterns/` — 記錄具體的 agent 失效模式，本層提供背後為什麼會發生的認知原理
- `enforcement/decision-efficiency.md` — 提供決策效率的執行規則，本層提供為什麼需要這些規則的認知基礎
- [`governance/document-sizing.md`](../../governance/document-sizing.md) — 提供文件拆分的治理規則，本層提供為什麼拆分能改善 agent 注意力的認知基礎
- `validation/` — 提供 stateless 驗證場景，本層提供為什麼需要驗證的認知邊界理論
- `governance/lifecycle/` — 提供知識生命週期管理，本層提供 agent 如何處理資訊的認知模型

← [回到 engineering/](../README.md)
