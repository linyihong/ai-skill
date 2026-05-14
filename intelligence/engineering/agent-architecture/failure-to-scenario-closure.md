# Failure-to-Scenario Closure（錯誤經驗轉化為驗證場景）

**Status**: `candidate-intelligence`
**Source**: 本系統實際運作經驗（[`validation/README.md`](../../../validation/README.md) 的 Failure → Scenario 閉環章節）

## 原則

**Every AI-system failure that can be reproduced statelessly should be converted into a validation scenario.**

每個可以被無狀態重現的 AI 系統錯誤，都應該轉化為驗證場景。

## 為什麼

1. **錯誤經驗是單次的** — 如果只記錄在 failure pattern 中，下次 session 的 agent 不會知道這個錯誤曾經發生過。
2. **Validation scenario 是 stateless 的** — 它不依賴 agent 的記憶或上下文，任何 session 都可以執行。
3. **錯誤會重複發生** — 如果沒有驗證機制，同樣的 routing 錯誤、同樣的規則違反會在不同 session 中反覆出現。
4. **Scenario 是 failure pattern 的 executable 形式** — Failure pattern 記錄「為什麼錯」，scenario 驗證「還會不會錯」。

## 症狀

| 症狀 | 說明 | 可信度 |
|------|------|--------|
| **同樣的錯誤在不同 session 重複出現** | Agent 在不同時間犯了完全相同的 routing 錯誤 | 高 |
| **Failure pattern 累積但從未被驗證** | 有 5+ 個 failure pattern，但沒有對應的 validation scenario | 高 |
| **錯誤修復後又復發** | 修復了一個 routing 問題，但幾個 session 後又出現 | 高 |
| **Agent 說「我知道這個錯誤」但還是犯了** | Agent 讀了 failure pattern，但在實際決策時仍然重複錯誤 | 中 |

## 判斷是否要建立 scenario

| 條件 | 應該建立 scenario | 不應該建立 scenario |
|------|-------------------|-------------------|
| 錯誤可 stateless 重現 | ✅ 是 | ❌ 否 |
| 錯誤是 routing 或規則遵守問題 | ✅ 是 | ❌ 否 |
| 錯誤是工具或環境特定問題 | ❌ 否 | ✅ 是 |
| 錯誤發生 >= 2 次 | ✅ 是 | ❌ 可選 |
| 錯誤的根本原因已理解 | ✅ 是 | ❌ 先分析原因 |

## 閉環流程

```
Failure Occurs
    ↓
記錄 failure pattern（enforcement/failure-patterns/）
    ↓
判斷是否可 stateless 重現
    ↓
可重現 → 建立 validation scenario（validation/scenarios/failure-derived/）
    ↓
不可重現 → 記錄在 failure pattern 中，標註「無法建立 scenario」
    ↓
下次架構變更或模型升級後，執行 validation suite
    ↓
Scenario pass → 確認錯誤已修復
Scenario fail → 分析是否為相同 root cause
```

## 不建議的做法

| 不建議 | 原因 |
|--------|------|
| 為每個小錯誤都建立 scenario | Scenario 數量爆炸，維護成本高 |
| 只記錄 failure pattern 不建立 scenario | 錯誤仍會重複發生 |
| Scenario 與 failure pattern 不同步 | 修改 failure pattern 後忘記更新 scenario |

## 相關 atoms

- [`cognitive-boundaries.md`](cognitive-boundaries.md) — 認知邊界（agent 無法可靠檢測自身錯誤）
- [`failure-recovery.md`](failure-recovery.md) — 失效恢復
- [`task-routing.md`](task-routing.md) — 任務路由（routing 錯誤是最常見的可 scenario 化錯誤）

## Token Impact

建立 scenario 需要 500-1500 token，但每次錯誤重複發生可能消耗 5000-20000 token。Scenario 是低成本高回報的投資。

---

← [回到 agent-architecture/](README.md)
