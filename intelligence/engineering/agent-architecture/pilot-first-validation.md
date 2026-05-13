# Pilot-First Validation（先驗證再抽象化）

**Status**: `candidate-intelligence`
**Source**: 本系統 Phase 28-31 實際運作經驗（[`plans/archived/technique-intelligence-pilot.md`](../../../plans/archived/technique-intelligence-pilot.md)）

## 原則

**Before abstracting a pattern into a generalized pipeline or framework, validate it with a concrete pilot first.**

在將某個模式抽象化為通用 pipeline 或框架之前，先用具體的 pilot 驗證它。

## 為什麼

1. **Premature abstraction 是 agent 常見的錯誤** — Agent 傾向於在看到 1-2 個案例後就建立通用框架，但這些框架往往遺漏了邊界案例。
2. **Pilot 提供真實的邊界案例** — 只有實際執行後，才能發現哪些步驟是通用的、哪些是 domain-specific 的。
3. **抽象化的成本很高** — 建立 pipeline 需要設計介面、處理邊界條件、撰寫文件，如果模式本身不成熟，這些成本就浪費了。
4. **Pilot 結果可以修正抽象方向** — 本系統的 Intelligence Extraction Pipeline 在 pilot 前原本設計為 5 步，pilot 後才發現需要 Step 7a（Shared-Rules 同步檢查）。

## 症狀

| 症狀 | 說明 | 可信度 |
|------|------|--------|
| **過早建立 pipeline** | 在只有 1 個案例時就開始設計通用框架 | 高 |
| **pipeline 需要頻繁修改** | 每遇到一個新案例就要調整 pipeline 結構 | 高 |
| **抽象層與實作脫節** | Pipeline 文件寫得很漂亮，但實際執行時從未完全遵循 | 中 |
| **忽略不做範圍** | 計畫中寫了「不做 X」，但抽象化時又把 X 包含進來 | 中 |

## 預防方式

1. **Pilot-first 原則** — 在抽象化之前，先用一個具體案例完整執行一次
2. **定義不做範圍** — 在 pilot 計畫中明確列出「不做什麼」，避免 scope creep
3. **Pilot 後回顧** — 完成 pilot 後，記錄哪些部分需要調整、哪些邊界案例被發現
4. **逐步泛化** — 不要一次從 1 跳到 N，先做 1 → 2 → 3，確認模式穩定後再抽象化
5. **保留 pilot 記錄** — Pilot 的執行記錄（如 [`notes/intelligence-extraction-observations.md`](../../../notes/intelligence-extraction-observations.md)）是抽象化的重要輸入

## 不建議的做法

| 不建議 | 原因 |
|--------|------|
| 看到一個案例就開始設計通用框架 | 遺漏邊界案例，框架需要反覆修改 |
| 跳過 pilot 直接建立 pipeline | Pipeline 的設計缺乏實戰驗證 |
| Pilot 成功後不記錄學到的教訓 | 抽象化時會遺漏 pilot 中發現的關鍵細節 |

## 相關 atoms

- [`premature-optimization.md`](../heuristics/premature-optimization.md) — 過早最佳化經驗法則
- [`cognitive-boundaries.md`](cognitive-boundaries.md) — 認知邊界（agent 無法可靠判斷「何時該抽象化」）
- [`failure-recovery.md`](failure-recovery.md) — 失效恢復（pilot 失敗後的策略調整）

## Token Impact

Pilot-first 策略在短期內看起來較慢（需要先做一個完整案例），但長期能節省大量 token：
- 避免建立錯誤的 pipeline（節省 5000-15000 token）
- 避免後續修改 pipeline（每次修改 2000-5000 token）
- Pilot 本身消耗 3000-8000 token，但這是投資而非成本

---

← [回到 agent-architecture/](README.md)
