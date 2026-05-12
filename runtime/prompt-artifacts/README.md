# Prompt Artifact Generator

`runtime/prompt-artifacts/` 負責 **task-specific prompt artifact 的定義與組合**。本層不是 prompt 模板本身，而是**定義 AI 如何根據任務類型自動組合最佳 prompt 結構**。

## 為什麼需要

現有系統已有：

| 層 | 提供 | 但缺少 |
|---|------|--------|
| `skills-index.yaml` | Skill routing（用哪個 skill） | 每個 skill 的 prompt 結構 |
| `workflow/` | 執行流程步驟 | 如何將步驟組合成 prompt |
| `intelligence/` | 工程智慧 atoms | 何時在 prompt 中引用哪些 atoms |
| `analysis/` | 分析方法 | 如何在 prompt 中引用方法 |
| `runtime/pipeline/` | Session lifecycle | 每個 stage 的 prompt 產出格式 |

**Prompt Artifact Generator** 填補這個缺口：讓 AI 在 routing stage 之後、execution stage 之前，自動組裝出一個**針對當前任務類型優化的 prompt artifact**。

## 核心概念

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

### 與現有層的關係

| 層 | 在 Artifact 中的角色 |
|---|---------------------|
| `workflow/` | 提供執行步驟（what to do） |
| `intelligence/` | 提供判斷依據（how to think） |
| `analysis/` | 提供分析方法（how to do） |
| `skills/` | 提供 entrypoint 與 tool-specific 內容 |
| `metadata/schema.md` | 提供 context cost 估算 |
| `runtime/pipeline/` | 提供 session lifecycle 整合點 |

## 檔案結構

```text
runtime/prompt-artifacts/
  README.md                    ← 本檔：概覽與設計
  artifact-templates.yaml      ← 各 task type 的 artifact 模板定義
  composition-rules.yaml       ← 如何從各層組合 artifact
```

## 使用方式

Agent 在 routing stage 完成後、execution stage 開始前：

1. 讀取 `artifact-templates.yaml`，根據 task type 找到對應模板。
2. 根據模板的 `sections` 定義，依序載入各 section 的內容。
3. 根據 `composition-rules.yaml` 的規則，決定哪些 intelligence atoms / analysis methods 需要嵌入。
4. 組裝成最終 prompt artifact，作為 execution stage 的 prompt 起點。

## 與既有層的差異

| 層 | 偏 | 範例 |
|---|----|------|
| `skills/` | 工具綁定的完整 skill 文件 | SKILL.md 包含 trigger、workflow、tools |
| `workflow/` | 工具中立的執行流程 | apk-analysis execution flow |
| `intelligence/` | 工程智慧 atoms | modular-monolith-vs-microservices.md |
| `runtime/prompt-artifacts/` | **動態組合的 prompt 結構** | 根據 task type 自動組裝上述各層 |

## 第一批 Task Types

| Task Type | 對應 Skill | 典型 Prompt 結構 |
|-----------|-----------|-----------------|
| `apk-analysis` | apk-analysis | Context → Workflow → Intelligence → Analysis Methods → Artifact Gates → Output |
| `app-development-guidance` | app-development-guidance | Context → Review Type → Checklists → Controls → Intelligence → Output |
| `repo-analysis` | repo-analysis | Context → Analysis Type → Methods → Intelligence → Output |
| `travel-planning` | travel-planning | Context → Planning Scope → Workflow → Intelligence → Output |
| `repo-governance` | repo-governance | Context → Governance Scope → Lifecycle → Validation → Output |
| `knowledge-navigation` | knowledge-navigation | Context → Query Type → Index → Summary → Output |
| `feedback-promotion` | feedback-promotion | Context → Lesson → Promotion Target → Validation → Output |

---

← [回到 Runtime](../README.md)
