# Architecture Selection Governance Workflow

## 目的

在 software delivery planning 中避免 default architecture。Agent 必須先做 fit analysis，再推薦 CRUD、DDD Lite、Full DDD、event-driven 或 microservices。

## 執行步驟

1. 收集 evidence：product brief、BDD、domain contract、API contract、team / lifecycle / integration constraints。
2. 評估 complexity axes：domain、invariant、workflow、integration、lifecycle、team、delivery speed。
3. 選擇最小可行架構。
4. 明確列出拒絕的更重架構與更輕架構。
5. 定義 validation：review、BDD、contract test、performance evidence 或 architecture decision record。

## 完成條件

若建議使用 DDD、CQRS、event sourcing 或 microservices，必須能說明它解決的 business complexity；否則降級。
