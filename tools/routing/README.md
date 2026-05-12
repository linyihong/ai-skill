# Tool Routing & Lazy Activation

`tools/routing/` 負責決定哪些 tool 需要 activate、哪些可以 deferred。目標是避免不必要的 tool loading 與 tool explosion。

## 核心原則

1. **Lazy activation by default** — 只 activate 目前 task 需要的 tool。
2. **Tool explosion detection** — 偵測 recursive tool loop（search → search → search）。
3. **Cost-aware routing** — 根據 tool 的 token 成本決定是否使用。

## Tool Activation 流程

```
Task Start
  │
  ├─ 1. 查 tools/metadata/ 取得 tool 清單與成本
  │
  ├─ 2. 依 task intent 決定需要哪些 tool
  │
  ├─ 3. 只 activate 必要的 tool
  │     ├─ preload: read-file, write-to-file（基本 IO）
  │     ├─ lazy: repo-search, execute-command（依 task 決定）
  │     └─ on_demand: ask-followup（使用者觸發）
  │
  ├─ 4. 監控 tool call 次數（circuit-breaker）
  │
  └─ 5. 偵測 tool explosion：
        - 同一 tool 連續呼叫 > 3 次
        - 5 分鐘內 tool call > 15 次
        - 遞迴深度 > 4 層
```

## Tool Explosion Detection

```yaml
explosion_signals:
  - signal: recursive_search
    pattern: search → search → search（無新結果）
    action: halt tool, suggest consolidate

  - signal: repetitive_read
    pattern: 重複讀取同一檔案 > 3 次
    action: suggest cache or summary

  - signal: tool_chain_too_long
    pattern: tool call chain > 10
    action: suggest decompose task

  - signal: output_too_large
    pattern: 單一 tool output > 10000 tokens
    action: suggest compression
```

## 與既有層的關係

- `ai-tools/`：工具配置與同步細節
- `tools/metadata/`：tool metadata schema 與成本資訊
- `tools/compression/`：tool output compression 策略
- `runtime/guards/circuit-breaker.yaml`：tool call 次數限制
