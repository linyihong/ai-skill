# Complexity Evaluation

**Status**: `candidate-intelligence`

## 判斷原則

複雜度不是檔案數或資料表數，而是 business rule、state transition、invariant、coordination 與 ownership 的組合。

## 分級

| 等級 | 訊號 | 合適策略 |
| --- | --- | --- |
| Low | CRUD、單一流程、低 invariant | CRUD、feature modules、simple service layer |
| Medium | 有 domain language、部分 invariant、流程開始分支 | DDD Lite、selective aggregate、modular monolith |
| High | 多 bounded context、高 invariant、長期演化、高整合壓力 | Full DDD、ACL、domain events、context map |

## 防誤判

- 不因資料表多就升級 DDD。
- 不因 framework 支援 CQRS 就使用 CQRS。
- 不因多人開發就拆 microservice；先看 ownership boundary。
