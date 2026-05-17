# Prompt Artifact Generator Philosophy

## Why Needed

現有系統已有：

| 層 | 提供 | 但缺少 |
|---|------|--------|
| `knowledge/runtime/routing-registry.yaml` | Route routing（用哪個 route） | 每個 route 的 prompt 結構 |
| `workflow/` | 執行流程步驟 | 如何將步驟組合成 prompt |
| `intelligence/` | 工程智慧 atoms | 何時在 prompt 中引用哪些 atoms |
| `analysis/` | 分析方法 | 如何在 prompt 中引用方法 |
| `runtime/pipeline/` | Session lifecycle | 每個 stage 的 prompt 產出格式 |

**Prompt Artifact Generator** 填補這個缺口：讓 AI 在 routing stage 之後、execution stage 之前，自動組裝出一個**針對當前任務類型優化的 prompt artifact**。

## Core Concept

### Artifact = 組合式 prompt 結構

```
Task Intent
    │
    ▼
┌─────────────────────────────────────────────┐
│         Prompt Artifact Generator            │
│                                              │
│  1. Identify task type (from intent)         │
│  2. Load artifact template (YAML)            │
│  3. Compose sections:                        │
│     ├─ Task Context (from user input)        │
│     ├─ Workflow Steps (from workflow/)       │
│     ├─ Intelligence Atoms (from intelligence/)│
│     ├─ Analysis Methods (from analysis/)     │
│     ├─ Artifact Gates (from workflow/)       │
│     └─ Output Format (from template)         │
│  4. Return composed prompt artifact          │
└─────────────────────────────────────────────┘
    │
    ▼
Execution Stage (uses composed artifact as prompt)
```

## 與既有文件的關係

- [`runtime/prompt-artifacts/`](../../runtime/prompt-artifacts/) — Runtime navigation entry point (data files: `artifact-templates.yaml`, `composition-rules.yaml`)
- [`runtime/prompt-artifacts/artifact-templates.yaml`](../../runtime/prompt-artifacts/artifact-templates.yaml) — Task type artifact template definitions
- [`runtime/prompt-artifacts/composition-rules.yaml`](../../runtime/prompt-artifacts/composition-rules.yaml) — Composition rules
- [`workflow/`](../../workflow/) — 提供執行步驟（what to do）
- [`intelligence/`](../../intelligence/) — 提供判斷依據（how to think）
- [`analysis/`](../../analysis/) — 提供分析方法（how to do）
- [`metadata/schema.md`](../../metadata/schema.md) — 提供 context cost 估算
- [`runtime/pipeline/`](../../runtime/pipeline/) — 提供 session lifecycle 整合點
