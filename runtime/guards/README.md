# Runtime Guards

`runtime/guards/` 負責「AI Runtime 安全護欄」。本層保存 circuit breaker、context pollution detection 等 runtime 保護機制，防止 agent 陷入遞迴循環、tool explosion、context 失控或 hallucination 高風險情境。

## 目前文件

- [`circuit-breaker.yaml`](circuit-breaker.yaml)：AI Runtime Circuit Breaker — 防止遞迴深度過深、tool call 爆炸、context 膨脹與 hallucination 風險。
- [`context-pollution.yaml`](context-pollution.yaml)：Context Pollution Detection — 偵測 context 是否已被污染（過長、過多修改、過多 modules）。

## 與既有層的關係

- `anti-patterns/` 記錄 runtime anti-patterns 的症狀與預防方式；本層提供具體的 guard 實作。
- `shared-rules/failure-patterns/` 記錄 agent 常犯的錯誤模式；本層提供 prevention gate。
- `runtime/` 定義 runtime 整體設計；本層是 runtime 的安全護欄子系統。
