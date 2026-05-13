# Stateless Validation Necessity（無狀態驗證的必要性）

**Status**: `candidate-intelligence`
**Source**: 本系統實際運作經驗（[`validation/README.md`](../../../validation/README.md) 的 Stateless Deterministic Validation 章節）

## 原則

**AI decision path validation must be stateless — it cannot depend on conversation memory or accumulated context.**

AI 決策路徑驗證必須是無狀態的 — 它不能依賴對話記憶或累積的上下文。

## 為什麼

1. **Agent 的決策受上下文影響** — 同一個 agent 在同一個 session 中，前 5 輪和後 50 輪的決策品質不同。如果驗證依賴特定上下文狀態，結果無法重現。
2. **Stateless 驗證是可重現的** — 任何 session、任何時間執行同一個 scenario，應該得到相同的結果（或至少相同的決策路徑）。
3. **Context 殘留會污染驗證結果** — 如果驗證 scenario 之前有前文提示，agent 可能會「猜到」預期答案，而不是真正根據 scenario 條件做決策。
4. **Stateless 是 contract testing 的核心** — 就像 unit test 不依賴資料庫狀態一樣，AI decision contract test 不依賴 conversation state。

## 症狀

| 症狀 | 說明 | 可信度 |
|------|------|--------|
| **同一個 scenario 每次結果不同** | 沒有修改任何檔案，但驗證結果時 pass 時 fail | 高 |
| **Scenario 依賴前文提示** | Scenario 中包含了「還記得我們之前討論的...」這類語句 | 高 |
| **驗證結果受 session 長度影響** | Session 開始時驗證 pass，50 輪後同一個 scenario fail | 高 |
| **Agent 在驗證時「作弊」** | Agent 直接從 scenario 的 `expected_route` 推斷答案 | 中 |

## 如何確保 Stateless

1. **Scenario 必須包含所有必要條件** — 不能假設 agent 已經知道任何背景資訊
2. **Scenario 不能包含預期答案** — `expected_route`、`expected_heuristics` 等欄位是給 evaluator 用的，不是給 agent 看的
3. **使用獨立 session 執行** — 每個 scenario 在全新的 session 中執行，不帶任何前文
4. **Scenario 的 `given` 區塊必須完整** — 包含 app type、artifacts、signals、constraints 等所有決策所需的輸入
5. **避免在 scenario 中使用「還記得...」這類語句** — 這暗示 agent 應該有前文記憶

## Stateless vs Stateful 驗證

| 面向 | Stateless Validation | Stateful Validation |
|------|---------------------|-------------------|
| 依賴 | 只有 scenario 的 `given` 條件 | Conversation memory + context |
| 可重現性 | ✅ 高（任何 session 結果相同） | ❌ 低（受上下文影響） |
| 測試對象 | AI 的 routing 決策 | AI 的長期行為一致性 |
| 適用場景 | Decision contract testing | Regression testing、長期專案 |
| 執行成本 | 低（獨立 session，快速） | 高（需要長時間 session） |

## 不建議的做法

| 不建議 | 原因 |
|--------|------|
| 在同一個 session 中執行多個 scenario | 前一個 scenario 的 context 會污染下一個 |
| Scenario 中包含「預期答案」的提示 | Agent 會「作弊」，不是真正根據條件做決策 |
| 依賴 agent 記住 scenario 的背景 | Agent 的記憶不可靠，尤其是跨多輪 tool call |
| 使用 conversation 中的歷史錯誤作為 scenario 條件 | 歷史錯誤的 context 無法在 stateless 環境中重現 |

## 相關 atoms

- [`cognitive-boundaries.md`](cognitive-boundaries.md) — 認知邊界（stateless validation 是外部閘門的一種形式）
- [`failure-to-scenario-closure.md`](failure-to-scenario-closure.md) — 錯誤經驗轉化為驗證場景（stateless 是轉化的前提）
- [`task-routing.md`](task-routing.md) — 任務路由（routing 決策最適合 stateless 驗證）

## Token Impact

Stateless validation 每次執行消耗 500-2000 token（讀取 scenario + 執行決策 + 產生 trace）。Stateful validation 需要 5000-20000 token（需要先建立 context）。Stateless 是更經濟的驗證方式。

---

← [回到 agent-architecture/](README.md)
