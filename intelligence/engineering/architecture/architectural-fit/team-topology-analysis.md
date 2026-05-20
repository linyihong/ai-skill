# Team Topology Analysis

**Status**: `candidate-intelligence`

## 判斷原則

Team boundary 可以支持架構拆分，但不能單獨證明 microservice 或 bounded context。拆分需要 ownership、release cadence、domain language 或 compliance boundary 的證據。

## 評估問題

- 是否有獨立 owner 可維護該 boundary？
- 是否需要不同 release cadence？
- 是否有不同 domain language 或 invariant？
- 是否有不同 scaling、security 或 compliance 要求？
- 拆分後 observability、testing、deployment 成本是否可承擔？

## 決策

若只是多人協作，優先 modular monolith 或 module boundary。只有當 ownership 和 runtime boundary 同時成立時，才考慮 service split。
