# Software Delivery Requirements Workflow

本目錄保存 software-delivery 的 requirements stage。它先確認 product impact 與 customer journey 對齊，再把 BDD-lite 作為 default delivery governance：behavior contract、acceptance boundary、validation target、traceability 與 ambiguity handling。

## 流程

1. [`product-impact-discovery/`](product-impact-discovery/README.md)：用 Impact Map × Customer Journey Map 驗證 Why / Who / How / What 是否對準 journey pain。
2. [`behavior-driven-discovery/`](behavior-driven-discovery/README.md)：理解 actor intent、behavior boundary、shared language。
3. [`acceptance-definition/`](acceptance-definition/README.md)：建立 acceptance criteria、validation target、regression scope。
4. [`ambiguity-resolution/`](ambiguity-resolution/README.md)：標記 assumption / open question / scoped out / invalidated。

## 輸出

每個會改變 observable behavior 的需求都要能追到：impact / journey evidence → requirement → behavior contract → acceptance criteria → validation target → execution artifact。
