# Software Delivery Architecture Workflow

本目錄保存 software-delivery 的 architecture stage。它接收 requirements stage 產生的 behavior boundary 與 acceptance baseline，再判斷 domain boundary、consistency boundary 與 architecture fit。

## 流程

1. [`domain-driven-design-fit/`](domain-driven-design-fit/README.md)：判斷是否需要 CRUD、DDD Lite 或 Full DDD。
2. [`bounded-context-discovery/`](bounded-context-discovery/README.md)：從 language / invariant / lifecycle 判斷 context boundary。
3. [`consistency-boundary-design/`](consistency-boundary-design/README.md)：根據 invariant 決定 transaction / eventual consistency。
4. [`architecture-escalation/`](architecture-escalation/README.md)：處理 architecture underfit / overfit。

## 兼容文件

本目錄保留既有 flat workflow files 作為具體 gate / framework：

- [`architecture-fit-analysis.md`](architecture-fit-analysis.md)
- [`architecture-decision-framework.md`](architecture-decision-framework.md)
- [`architecture-minimality-principle.md`](architecture-minimality-principle.md)
- [`architecture-overengineering-detection.md`](architecture-overengineering-detection.md)

## 輸出

Architecture recommendation 至少包含：chosen strategy、rejected lighter option、rejected heavier option、fit evidence、validation plan、upgrade/downgrade trigger。
