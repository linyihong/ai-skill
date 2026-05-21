# Runtime Guards Philosophy

## Purpose

`runtime/guards/` 負責「AI Runtime 安全護欄」。本層保存 circuit breaker、context pollution detection 等 runtime 保護機制，防止 agent 陷入遞迴循環、tool explosion、context 失控或 hallucination 高風險情境。

## Current Guards

- **Circuit Breaker**：防止遞迴深度過深、tool call 爆炸、context 膨脹與 hallucination 風險。
- **Context Pollution Detection**：偵測 context 是否已被污染（過長、過多修改、過多 modules）。

## 與既有文件的關係

- [`runtime/guards/`](../../runtime/guards/) — Runtime navigation entry point (data files: `circuit-breaker.yaml`, `context-pollution.yaml`)
- [`runtime/runtime.db`](../../runtime/runtime.db) — Circuit breaker definitions
- [`runtime/runtime.db`](../../runtime/runtime.db) — Context pollution detection
- [`anti-patterns/`](../../anti-patterns/) — Runtime anti-patterns 的症狀與預防方式
- [`enforcement/failure-patterns/`](../../enforcement/failure-patterns/) — Agent 常犯的錯誤模式
