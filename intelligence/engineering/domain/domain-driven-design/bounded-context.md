# Bounded Context

**Status**: `candidate-intelligence`

## 判斷原則

Bounded context 是業務語言與模型一致性的邊界，不是資料表、資料夾或 microservice 的同義詞。

當同一個詞在不同流程中代表不同規則時，應優先切 bounded context；當只是不同功能共用同一套規則時，不應為了看起來像 DDD 而拆分。

## 適用訊號

- 同一名詞在銷售、履約、帳務或客服中有不同生命週期。
- 不同團隊對同一實體有不同必填欄位、狀態或合法操作。
- 某個模型變更總是讓不相關功能一起修改。
- 外部系統的資料模型迫使內部模型使用錯誤語言。

## 不適用訊號

- 只是 CRUD entity 多。
- 只是資料表很多，但業務語言一致。
- 只是想拆 microservice 或 module。
- 只是因為資料夾太大。

## 決策規則

```text
需要新 bounded context？
  ├─ 語言是否衝突？否 → 不拆
  ├─ 不變量是否不同？否 → 不拆
  ├─ lifecycle 是否獨立？否 → 優先 module/package
  └─ 是否需要防外部模型污染？是 → 考慮 bounded context + ACL
```

## 邊界

Bounded context 可類比 cognitive governance 的 claim scope，但不可替代 evidence qualification 或 runtime claim scope gate。
