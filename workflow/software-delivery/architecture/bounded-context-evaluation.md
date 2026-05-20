# Bounded Context Evaluation

## 目的

判斷是否需要 bounded context，避免把資料表、資料夾或 microservice 誤認為 domain boundary。

## 檢查問題

- 同一名詞是否在不同流程中有不同意義？
- 不同 context 是否有不同 invariant 或 lifecycle？
- 是否有外部模型污染內部模型？
- 是否有不同 team ownership 或 release cadence？
- 若不拆 context，錯誤會如何擴散？

## 決策

- 語言衝突 + invariant 差異成立：建立 bounded context。
- 只有資料結構不同：優先 module / package boundary。
- 只有 deployment 想拆：回到 team topology 與 lifecycle evaluation。

## 連結

詳細判斷見 `intelligence/engineering/domain/domain-driven-design/bounded-context.md`。
