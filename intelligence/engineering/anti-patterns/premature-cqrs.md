# Premature CQRS

**Status**: `candidate-intelligence`

## 反模式

在沒有 read/write divergence、scale pressure、audit、eventual consistency 或複雜查詢需求時，預設導入 CQRS。

## 訊號

- Command / query 分離只讓檔案數增加。
- Read model 與 write model 幾乎相同。
- 測試和 transaction 更複雜，但沒有對應收益。
- 團隊無法說明 eventual consistency boundary。

## 修正

先使用單一 model 或 explicit query method。當讀寫模型真的分化、讀取效能或查詢形狀成為 recurring bottleneck 時再升級。
