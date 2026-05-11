# AI-native Knowledge Operating System 下一階段升級規劃

本文件是下一階段架構升級規劃書。它承接 [`ai-native-knowledge-operating-system.md`](ai-native-knowledge-operating-system.md) 的 reference-first、goal ledger、failure learning、rule weight 與 close-loop 基礎，規劃如何從現有 skill-centered repository 演進成 AI Knowledge Runtime System。

## 目前走到哪裡

已完成的基礎層：

- Root `README.md` 已是 AI-native Knowledge Operating System dashboard。
- `shared-rules/` 已建立 dependency reading、linked updates、conversation goal ledger、failure learning、rule weight、language consistency 等 operating rules。
- `architecture/ai-native-knowledge-operating-system.md` 已定義 reference-first、compatibility inventory、Phase 3 deprecation checklist。
- `.agent-goals/` 已作為 project-local active goal ledger 使用，完成後刪除，不進 git。
- Cursor / Claude tool docs 已指向 central repository 與 shared-rule bootstrap。

尚未完成的下一階段：

- 尚未建立 `analysis/`、`intelligence/`、`workflow/`、`runtime/`、`memory/`、`feedback/`、`models/`、`governance/`、`knowledge/`、`metadata/` 等正式目錄。
- 既有 `skills/` 仍同時承載 workflow、analysis 方法、工程智慧、templates 與 feedback lessons。
- 尚未建立 Knowledge Atom metadata schema。
- 尚未建立可供 runtime 使用的 indexes、summaries、graphs 與 routing metadata。
- 尚未定義 multi-model routing / compression strategy。

## 核心問題

下一階段要回答的不只是「有哪些 prompts 或 skills」，而是：

- AI 如何工作。
- AI 如何學習。
- AI 如何沉澱知識。
- AI 如何找到正確知識。
- AI 如何演化知識。
- AI 如何多模型協作。
- AI 如何長期維護知識。

因此整體方向要從 **Skill Collection** 升級為 **AI Knowledge Runtime System**。

## 目標架構分層

下一階段建議正式拆分：

```text
analysis/
intelligence/
workflow/
runtime/
memory/
feedback/
models/
governance/
knowledge/
metadata/
```

這些目錄不是一次搬完所有內容，而是先建立責任邊界、metadata schema 與 navigation layer，再逐批遷移。

## 各層責任

### `analysis/`

負責「如何觀察與拆解」。

建議結構：

```text
analysis/
  apk/
  repo/
  production/
  issue/
```

核心責任：

- reverse engineering。
- 流程拆解。
- 技術觀察。
- pattern extraction。
- 分析方法。

不應承載過多：

- trade-off。
- architecture lesson。
- anti-pattern conclusion。

這些應抽取到 `intelligence/`。

### `intelligence/`

負責「沉澱工程智慧與領域知識」。

建議結構：

```text
intelligence/
  engineering/
    domain/
    architecture/
    failure/
    realtime/
    erp/
  travel/
  business/
```

核心責任：

- engineering decision。
- trade-off。
- anti-pattern。
- architecture lesson。
- reusable domain knowledge。

`intelligence/` 是 Senior Engineer Brain。

### `workflow/`

負責「AI 如何執行工作」。

建議結構：

```text
workflow/
  app-development-guidance/
  apk-analysis/
  repo-analysis/
  travel-planning/
```

核心責任：

- planning flow。
- task decomposition。
- review flow。
- orchestration flow。
- execution flow。

`workflow/` 應 reference `intelligence/`，而不是內嵌大量知識。

### `runtime/`

負責「AI 系統如何運作」。

建議結構：

```text
runtime/
  scheduler/
  routing/
  orchestration/
  context/
```

核心責任：

- dynamic loading。
- context injection。
- orchestration。
- task routing。
- context pruning。
- agent coordination。

### `memory/`

負責「長期記憶」。

建議結構：

```text
memory/
  short-term/
  episodic/
  project/
  failure/
```

核心責任：

- experience replay。
- long-term memory。
- historical context。

### `feedback/`

負責「系統如何持續演化」。

建議結構：

```text
feedback/
  replay/
  extraction/
  refinement/
  promotion/
```

核心責任：

- workflow refinement。
- intelligence extraction。
- lesson replay。
- knowledge evolution。

### `models/`

負責「不同模型如何協作」。

建議結構：

```text
models/
  claude/
  gpt/
  gemini/
  qwen/
  small-model/
```

核心責任：

- capability profile。
- reasoning strength。
- context limit。
- routing strategy。
- compression strategy。
- prompt adaptation。

### `governance/`

負責「知識治理與系統維護」。

建議結構：

```text
governance/
  cleanup/
  splitting/
  lifecycle/
  validation/
```

核心責任：

- duplicate cleanup。
- lifecycle management。
- validation。
- splitting rules。
- dependency maintenance。

### `knowledge/`

負責「知識導航與知識圖譜」。

建議結構：

```text
knowledge/
  atoms/
  indexes/
  summaries/
  graphs/
  runtime/
```

核心思想是 Atomic Knowledge。真正目標不是單純拆小文件，而是支援 Dynamic Context Composition。

不要讓系統變成 Knowledge Fragment Hell；每個 atom 都必須能被 index、summary、graph 與 runtime metadata 找到。

### `metadata/`

負責「知識控制系統」。

建議結構：

```text
metadata/
  rules/
  ranking/
  confidence/
  compatibility/
```

`metadata/` 是 Rule Metadata System 的核心。Metadata 不是只描述文件，而是控制 runtime 行為。

每個 Knowledge Atom 應包含：

```yaml
id:
type:
domain:
tags:
priority:
confidence:
stability:
complexity:
context_cost:
depends:
related:
conflicts:
models:
summary:
checklist:
```

Runtime 依賴 metadata 進行：

1. Context Routing：現在該載入哪些知識。
2. Priority Selection：哪些規則優先。
3. Conflict Resolution：規則衝突時如何仲裁。
4. Dynamic Loading：根據 task 載入知識。
5. Model-aware Compression：小模型只讀 checklist 或 compressed knowledge。
6. Knowledge Promotion：`candidate` → `validated` → `stable`。
7. Knowledge Cleanup：找出過期知識。
8. Dependency Graph Construction：自動建立 knowledge graph。

## Knowledge Navigation System

Atomic Knowledge 必須搭配 navigation + index system。

建議建立：

```text
knowledge/indexes/
knowledge/summaries/
knowledge/graphs/
knowledge/runtime/
```

真正重要的不是知識量，而是 AI 能否找到正確知識。

## Intelligence Feedback Loop

系統應形成閉環：

```text
Analysis -> Extraction -> Intelligence -> Workflow -> Feedback
```

例：

```text
apk-analysis
  -> intelligence extraction
  -> realtime intelligence
  -> workflow reference
  -> future refinement
```

## Multi-model Runtime Architecture

未來模型一定是混用，因此 workflow 應 model-aware。

範例：

```yaml
small-model:
  use:
    - checklist
    - compressed knowledge

large-model:
  use:
    - full intelligence graph
```

## Knowledge Lifecycle System

知識一定會熵增，因此每個知識單元需要 lifecycle：

```text
temporary/
candidate/
validated/
stable/
deprecated/
```

## 遷移原則

1. 不一次搬完所有檔案。
2. 先建立 top-level directory README，定義責任邊界。
3. 先定義 metadata schema，再遷移 content。
4. 先選一個 skill 做示範遷移，再擴展到其他 skill。
5. 保留 `skills/` 與 `shared-rules/` 相容層，直到 workflow / intelligence / metadata / runtime 的 reference path 穩定。
6. 每次搬移都必須保留舊連結或提供 redirect / index。
7. 每次遷移都要經過 `.agent-goals`、linked updates、diff review、commit/push/readback、clean status。

## 建議遷移階段

### Phase 0：目前已完成的基礎

- OS dashboard。
- `reference-first`。
- `rule-weight`。
- goal ledger。
- failure learning。
- language consistency。
- compatibility inventory。
- Phase 3 deprecation checklist。

### Phase 1：建立新架構目錄

建立下列目錄與 README：

```text
analysis/
intelligence/
workflow/
runtime/
memory/
feedback/
models/
governance/
knowledge/
metadata/
```

每個 README 只定義：

- 該層責任。
- 放什麼。
- 不放什麼。
- 與現有 `skills/`、`shared-rules/`、`ai-tools/` 的關係。
- 第一批候選遷移來源。

### Phase 2：Metadata System

新增：

```text
metadata/schema.md
metadata/rules/
metadata/ranking/
metadata/confidence/
metadata/compatibility/
```

定義 Knowledge Atom schema 與 required/optional 欄位。

### Phase 3：Knowledge Navigation

新增：

```text
knowledge/indexes/
knowledge/summaries/
knowledge/graphs/
knowledge/runtime/
```

先做 index 與 summary，不急著做完整 graph runtime。

### Phase 4：Workflow / Intelligence 分離

第一個示範對象建議使用 `apk-analysis`：

- `analysis/apk/`：保留觀察、拆解、traffic/runtime 分析方法。
- `workflow/apk-analysis/`：保留 agent 執行流程、task decomposition、review flow。
- `intelligence/engineering/failure/`：抽取反覆失效模式與 anti-pattern。
- `intelligence/engineering/architecture/`：抽取架構與 trade-off lessons。

### Phase 5：Runtime / Models

定義：

- context routing。
- dynamic loading。
- context pruning。
- model capability profiles。
- small-model / large-model 使用策略。

### Phase 6：Lifecycle / Governance

定義：

- knowledge lifecycle。
- duplicate cleanup。
- dependency graph maintenance。
- validation gates。
- deprecation / archive process。

## `.agent-goals` 拆解

下一階段應建立並追蹤以下 goals：

| Priority | Goal | Next action | Completion criteria |
| --- | --- | --- | --- |
| P1 | 建立 next-stage upgrade plan | 本文件完成並連到入口文件 | 規劃書 commit/push/readback，root/architecture 入口可找到 |
| P1 | 建立 top-level architecture directories | 建立 10 個目錄與 README skeleton | 每個目錄責任邊界清楚，不搬移大量內容 |
| P2 | 設計 metadata schema | 建立 `metadata/schema.md` | Schema 可套用到第一批 Knowledge Atom |
| P2 | 建立 knowledge navigation index | 建立 `knowledge/indexes/README.md` 與初版索引格式 | Agent 能從 index 找到 task-relevant knowledge |
| P2 | 遷移第一個 skill 作為示範 | 選 `apk-analysis` 拆出 analysis/workflow/intelligence | 舊入口仍可用，新路徑可被 reference-first 找到 |

## 最終目標

AI-native Knowledge Operating System 的最終目標不只是讓 AI 產生內容，而是建立：

- AI-native Engineering System。
- Knowledge Graph Runtime。
- Multi-model Orchestration。
- Engineering Intelligence Platform。
- Long-term AI Learning System。

未來真正瓶頸不會只是模型強度，而是知識是否能被正確管理、導航、組合與演化。

這是本 repository 下一階段的核心方向。
