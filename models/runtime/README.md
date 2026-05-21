# Model Runtime Boundary

`models/runtime/` 定義 model-aware routing 可被 knowledge/runtime lookup 使用的 minimal primitives。這不是 provider model state，也不是 actual model selection log。

## 入口

- [`routing-primitives.md`](routing-primitives.md)：minimal strategy primitives。
- [`context-budgeting.md`](context-budgeting.md)：context loading 與 compression budget。
- [`execution-cost-strategy.md`](execution-cost-strategy.md)：strategy cost 與 escalation。
- [`adaptive-loading.md`](adaptive-loading.md)：summary-first 到 source-backed 的 loading rules。

## Runtime Rule

只有 lookup metadata、routing candidate 與 validation scenario 證明有用的 primitives 才能進 runtime。不得保存 raw model output、provider availability、main chat model state 或 unverified capability claims。
