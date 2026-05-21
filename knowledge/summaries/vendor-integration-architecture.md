## intelligence.vendor-integration-architecture

| 欄位 | 值 |
| --- | --- |
| Atom ID | `intelligence.vendor-integration-architecture` |
| Source path | `intelligence/engineering/architecture/vendor-integration-architecture.md` |
| Lifecycle | `candidate` |
| Summary | 整合超過 3 個外部廠商（支付聚合、社群登入、IM、博弈聚合、廣告聯播等）時的整合策略選型。五種策略：A. Adapter/Strategy（單模組多實作）/ B. Compile-time submodule per vendor / C. Plugin SPI（runtime 載入）/ D. Out-of-process service / E. Hybrid 分層。N ≥ 10 必須跳出 compile-time module per vendor 模式，否則編譯時間、IDE、升級成本爆炸。 |
| When to read | 評估如何整合新廠商；專案內 vendor 相關 Maven/Gradle module 數量 ≥ 10；升級共用框架需要動多個 vendor 模組；廠商需要熱換、灰度、動態啟停；廠商 SDK 相依衝突。 |
| Do not use for | 不可取代特定廠商 SDK 的整合文件；不可作為 OSGi/ServiceLoader 等具體 plugin 框架的實作教學。 |
| Context cost | ~320 tokens |
| Estimated full cost | ~2500 tokens |
| Validation signal | 廠商數對應策略符合軸向；廠商相關 module 數 ≤ 廠商數 × 4；全量 build < 5 分鐘；新增廠商成本 < 3 工作日；升級共用相依不需動所有 vendor 模組。 |
| Last checked | 2026-05-21 |
