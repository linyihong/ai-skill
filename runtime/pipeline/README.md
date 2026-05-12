# Runtime Pipeline

本目錄定義 **Runtime Pipeline** — 將所有 Runtime Quality & Safety 元件（Token Budget、Context Health、Circuit Breaker、Tool Routing、Compression、Memory）串接成可執行的 orchestration flow。

## 為什麼需要 Pipeline

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

## Pipeline 架構

```text
runtime/pipeline/
  README.md                  ← 本檔：pipeline 概覽
  session-lifecycle.yaml     ← Session lifecycle stages
  context-flow.yaml          ← Progressive context expansion
  guard-chain.yaml           ← Guard execution order
  relevance-engine.yaml      ← Skill Relevance Engine
```

## Pipeline 流程圖

```
┌─────────────────────────────────────────────────────────────┐
│                    Session Lifecycle                         │
│                                                             │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐ │
│  │Bootstrap │ → │ Routing  │ → │Execution │ → │Close-loop│ │
│  │  Stage   │   │  Stage   │   │  Stage   │   │  Stage   │ │
│  └────┬─────┘   └────┬─────┘   └────┬─────┘   └────┬─────┘ │
│       │              │              │              │        │
│       ▼              ▼              ▼              ▼        │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐ │
│  │Core      │   │Skill     │   │Tool      │   │Memory    │ │
│  │Bootstrap │   │Relevance │   │Execution │   │Summary   │ │
│  │~800 tok  │   │Engine    │   │+ Guards  │   │Writeback │ │
│  └──────────┘   └──────────┘   └──────────┘   └──────────┘ │
└─────────────────────────────────────────────────────────────┘
         │              │              │              │
         ▼              ▼              ▼              ▼
   ┌─────────────────────────────────────────────────────────┐
   │                 Cross-stage Components                   │
   │                                                         │
   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐ │
   │  │Token Budget │  │Context      │  │Context Pollution│ │
   │  │(continuous) │  │Health Score │  │Detection        │ │
   │  │             │  │(per stage)  │  │(on threshold)   │ │
   │  └─────────────┘  └─────────────┘  └─────────────────┘ │
   │                                                         │
   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐ │
   │  │Circuit      │  │Tool Output  │  │Memory           │ │
   │  │Breaker      │  │Compression  │  │(working/summary │ │
   │  │(per action) │  │(per output) │  │ /decision)      │ │
   │  └─────────────┘  └─────────────┘  └─────────────────┘ │
   └─────────────────────────────────────────────────────────┘
```

## 元件間通訊

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

## 與既有層的關係

| Pipeline 元件 | 依賴的既有層 | 關係 |
|-------------|------------|------|
| session-lifecycle | `CORE_BOOTSTRAP.md`, `skills-index.yaml` | 定義 bootstrap → routing → execution → close-loop |
| context-flow | `knowledge/summaries/`, `runtime/context/ttl-policy.yaml` | 定義 progressive expansion 順序 |
| guard-chain | `runtime/guards/circuit-breaker.yaml`, `runtime/guards/context-pollution.yaml` | 定義 guard 執行順序 |
| relevance-engine | `skills-index.yaml` (weight/domains/conflicts) | 使用 skills v2 metadata 做 scoring |

## 使用方式

Pipeline 不是一個需要「啟動」的服務，而是 Agent 在每個 session 中遵循的**執行模型**：

1. **Session 開始**：Agent 讀取 `session-lifecycle.yaml` 了解當前階段。
2. **每個階段**：Agent 檢查 `guard-chain.yaml` 確認哪些 guards 需要執行。
3. **Context 載入**：Agent 遵循 `context-flow.yaml` 的 progressive expansion。
4. **Skill 選擇**：Agent 使用 `relevance-engine.yaml` 的 scoring 邏輯決定載入哪些 skill。
5. **Session 結束**：Agent 執行 close-loop，寫入 memory summary 與 decision record。

---

← [回到 Runtime](../README.md)
