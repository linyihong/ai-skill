# Runtime Pipeline Philosophy

## Why Pipeline

Phase 1 建立了獨立的元件，但它們各自獨立運作：

```
Token Budget  ──→  Context Health  ──→  Circuit Breaker
                                                ↓
Tool Metadata  ──→  Tool Routing  ──→  Compression
                                                ↓
Memory (working/summary/decision)
```

Pipeline 將這些元件串接成**單一可執行流程**，確保：

1. **執行順序確定**：每個階段有明確的輸入/輸出。
2. **Guard chain 順序**：circuit breaker guards 按正確順序執行。
3. **Context 漸進擴展**：從 summary → module → detailed → raw，不一次載入全部。
4. **Session lifecycle 管理**：bootstrap → routing → execution → close-loop，每個階段有明確的進入/離開條件。
5. **元件間通訊**：Token Budget 的 hard stop 會觸發 Context Pollution 的 auto-archive，Context Health 的 critical 會觸發 Compression 的 minimal level。

## Pipeline Architecture

```
runtime/pipeline/
  session-lifecycle.yaml     ← Session lifecycle stages
  context-flow.yaml          ← Progressive context expansion
  guard-chain.yaml           ← Guard execution order
  relevance-engine.yaml      ← Skill Relevance Engine
```

## Component Communication

| 觸發事件 | 來源元件 | 目標元件 | 行為 |
|---------|---------|---------|------|
| Token usage > 70% | Token Budget | Context Health | 觸發 health score re-evaluation |
| Token usage > 90% | Token Budget | Context Pollution | 強制 auto-archive session |
| Context Health < 0.50 | Context Health | Tool Compression | 切換至 structured/minimal level |
| Context Health < 0.50 | Context Health | Circuit Breaker | 啟動 context growth guard |
| Recursive depth > 4 | Circuit Breaker | Session Lifecycle | 強制進入 close-loop stage |
| Tool calls > 20/task | Circuit Breaker | Tool Routing | 暫停工具呼叫，建議分解 |
| Pollution score critical | Context Pollution | Memory | auto-archive 到 memory/working/ |
| Pollution score critical | Context Pollution | Session Lifecycle | 建議新 session |
| Skill relevance < 0.5 | Relevance Engine | Skill Index | 跳過該 skill 的載入 |
| Compression active | Tool Compression | Tool Output | 輸出壓縮後餵回 context |

## 與既有文件的關係

- [`runtime/pipeline/README.md`](../../runtime/pipeline/README.md) — Runtime navigation entry point
- [`runtime/pipeline/session-lifecycle.yaml`](../../runtime/pipeline/session-lifecycle.yaml) — Session lifecycle stages
- [`runtime/pipeline/context-flow.yaml`](../../runtime/pipeline/context-flow.yaml) — Progressive context expansion
- [`runtime/pipeline/guard-chain.yaml`](../../runtime/pipeline/guard-chain.yaml) — Guard execution order
- [`runtime/pipeline/relevance-engine.yaml`](../../runtime/pipeline/relevance-engine.yaml) — Skill Relevance Engine
- [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md) — Bootstrap entry point
- [`runtime/guards/circuit-breaker.yaml`](../../runtime/guards/circuit-breaker.yaml) — Circuit breaker
- [`runtime/guards/context-pollution.yaml`](../../runtime/guards/context-pollution.yaml) — Context pollution detection
