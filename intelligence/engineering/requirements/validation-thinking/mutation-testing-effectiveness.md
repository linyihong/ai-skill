# Mutation Testing Effectiveness

**Status**: `candidate-intelligence`

## 判斷原則

高 coverage 不等於測試有效。Mutation testing 的價值是刻意製造小型錯誤，反問現有測試是否真的能抓到 behavior、contract 或 invariant 的破壞。

## 使用時機

不要把 mutation testing 當成所有小改動的必跑 gate。當以下訊號出現時才升級：

- AI 生成或大幅改寫核心邏輯。
- 權限、安全、金流、庫存、資料一致性或狀態轉換規則。
- Domain invariant 密集，example-based tests 可能只覆蓋 happy path。
- Refactor 宣稱無行為變更，但缺少強 regression proof。
- Coverage 很高，但測試只驗證路徑執行或 snapshot shape。
- BDD scenario 有 validation target，但缺少 negative path 或 edge case。

## 流程

1. 描述風險：先說明想防止哪類錯誤，而不是盲目提高 mutant count。
2. 產生 mutants：針對 branch、boundary、comparison、boolean、nullability、error handling、invariant guard 產生小型錯誤版本。
3. 過濾 equivalent mutants：移除語意相同或無法由觀察行為區分的 mutants。
4. 補有效測試：若 mutant survived，補 BDD、unit、property、contract、fixture 或 invariant test，直到至少一個有價值 mutant 被 killed。

## 判讀

| 結果 | 意義 | 行動 |
| --- | --- | --- |
| `killed` | 測試能抓到該錯誤 | 保留測試，記錄它保護的 risk |
| `survived` | 測試沒有抓到該錯誤 | 補 validation target 或縮小 correctness claim |
| `equivalent` | mutant 與原行為不可區分 | 不把它當測試缺口，記錄過濾原因 |

## 邊界

- Mutation score 不是新的 coverage KPI；它是 test effectiveness signal。
- 不要求所有專案導入 mutation runner；可用手動 negative check 或 targeted mutant review。
- 不讓 runtime 理解 mutant syntax；只可能產生 `missing_validation_target`、`behavior_scope_overclaim` 或 `invariant_violation` 類壓縮信號。
