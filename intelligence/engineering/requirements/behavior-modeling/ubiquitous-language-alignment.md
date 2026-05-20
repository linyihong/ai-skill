# Ubiquitous Language Alignment

**Status**: `candidate-intelligence`

## 判斷原則

BDD 與 DDD 的交集是 shared language。需求語言若沒有先對齊，後續 bounded context、scenario、contract 與 test 都會漂移。

## 適用訊號

- 使用者、產品文件、BDD、API contract 與實作使用不同名詞描述同一行為。
- 同一名詞在不同 actor 或 workflow 中含義不同。
- Scenario 看似通過，但實際驗證的是不同語意。

## 行動

1. 擷取會影響行為的名詞。
2. 標記 actor、狀態、操作、限制與反例。
3. 若名詞在不同 context 中含義不同，交給 `architecture/domain-modeling/bounded-context.md`。
4. 若名詞只是 acceptance wording，留在 requirements layer。

## 邊界

本文件不建立 DDD bounded context；它只先穩定 requirement / behavior language。
