# ADR-004: Feedback Promotion Pipeline

## Status

**Accepted**

## Context

Feedback lessons 是系統演化的核心原料。每次 agent 在執行任務時學到的經驗教訓，都應該有機會被系統吸收、泛化、並推廣到其他情境。

最初 feedback lessons 只存在於 `skills/*/feedback_history/`，問題是：

1. **Skill-bound**：lesson 綁定在特定 skill 下，無法跨 skill 推廣。
2. **無結構化**：lesson 格式自由，難以自動判斷是否成熟到可 promotion。
3. **無 lifecycle**：lesson 一旦寫入就永遠存在，沒有 archive / deprecate 機制。
4. **無自動化**：promotion 完全依賴人工判斷，無法規模化。

需要一個 pipeline 讓 feedback lesson 可以從原始觀察，經過結構化評估，自動或半自動地 promotion 到目標層。

## Decision

建立 **Feedback Promotion Pipeline**，包含三個前置階段與一個 promotion 引擎：

```text
feedback_history (raw)
       ↓
feedback/replay/    ← 經驗回放（分析失敗模式與成功模式）
       ↓
feedback/extraction/ ← 智慧提取（從 replay 中提取可泛化的 lesson）
       ↓
feedback/refinement/ ← 流程精煉（根據 lesson 改進 workflow）
       ↓
feedback/promotion/  ← Promotion Engine（scoring → decision → write → validate）
       ↓
Target Layers: workflow/ | intelligence/ | enforcement/ | memory/ | skill-doc/
```

### 核心設計決策

1. **Replay 不是 re-execution**：Replay 是分析過往 session 的經驗，不是重新執行。它回答「當時發生了什麼、為什麼、學到了什麼」。
2. **Failure first**：Replay 優先處理失敗模式（failure repeat ≥ 2、session blocker），成功模式可選。
3. **必須泛化**：Extraction 產出的 intelligence atom 必須是泛化的（不包含專案特定細節），否則退回。
4. **有成本**：Replay 和 extraction 都有 token 成本，不應對每個 session 都執行。使用 trigger conditions 控制頻率。
5. **Scoring-driven promotion**：Promotion 使用 5 維度 scoring（impact 0.30 + maturity 0.25 + frequency 0.20 + freshness 0.15 + urgency 0.10），threshold 0.7 立即 promotion，0.5 進 backlog。

### Promotion 目標層

| Score | Target | 條件 |
| --- | --- | --- |
| ≥ 0.7 | `enforcement/` | 跨 skill 驗證過的執行規則 |
| ≥ 0.7 | `intelligence/` | 跨 skill 的工程智慧 |
| ≥ 0.7 | `workflow/` | 可改進執行流程的 lesson |
| ≥ 0.5 | `skill-doc/` | Skill-specific 的改進建議 |
| < 0.5 | archive | 尚未成熟的 lesson |

## Consequences

### 正面

- **Feedback 有完整 lifecycle**：從 raw observation → replay → extraction → refinement → promotion → archive。
- **跨 skill 推廣**：成熟的 lesson 可以 promotion 到 intelligence/ 或 enforcement/，被多個 skill 使用。
- **自動化門檻**：Scoring 機制讓 promotion 可部分自動化，降低人工負擔。
- **Token 效率**：Replay 和 extraction 只在 trigger conditions 滿足時執行，不浪費 token。

### 負面

- **Pipeline 複雜度**：5 階段 pipeline 需要 agent 理解完整流程才能正確操作。
- **Scoring 主觀性**：5 維度 scoring 的權重設定可能不適合所有情境，需要持續調整。
- **初期資料不足**：在 feedback_history 累積足夠 lesson 之前，pipeline 的效益不明顯。

## Alternatives Considered

- **直接 promotion**：從 feedback_history 直接寫入目標層，跳過 replay/extraction/refinement。但 lesson 未經泛化，可能包含專案特定細節。不採用。
- **純人工 promotion**：完全依賴開發者判斷。無法規模化。不採用。
- **只保留 feedback_history**：維持現狀，不建立 pipeline。但 lesson 無法跨 skill 推廣。不採用。

## Related

- [`feedback/replay/README.md`](../feedback/replay/README.md) — 經驗回放階段
- [`feedback/extraction/README.md`](../feedback/extraction/README.md) — 智慧提取階段
- [`feedback/refinement/README.md`](../feedback/refinement/README.md) — 流程精煉階段
- [`feedback/promotion/README.md`](../feedback/promotion/README.md) — Promotion Engine
- [`feedback/pipeline/promotion-engine.yaml`](../feedback/pipeline/promotion-engine.yaml) — Scoring 規則
- [`feedback/pipeline/promotion-workflow.yaml`](../feedback/pipeline/promotion-workflow.yaml) — Promotion workflow
- [`feedback/pipeline/lifecycle-automation.yaml`](../feedback/pipeline/lifecycle-automation.yaml) — Lifecycle automation
- [`enforcement/feedback-lessons.md`](../enforcement/feedback-lessons.md) — Feedback lesson 格式
