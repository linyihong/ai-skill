# Ai-skill Framework Glossary

> Canonical entries for Ai-skill framework vocabulary, runtime semantics, cognitive vocabulary and architecture contracts.
>
> Schema spec：[`README.md`](README.md)。Validator：`ai-skill glossary validate`。
>
> 上游 plan：[`plans/active/2026-05-25-1000-context-language-glossary-system.md`](../../plans/active/2026-05-25-1000-context-language-glossary-system.md)（Phase 3）。
>
> 編寫慣例：entries 按 term snake_case 字母順序排列；`status: candidate` 標明為 economics plan 預留，待 [`plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md`](../../plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md) Phase 1 確認 owner path 後 promote。

---

## activation_fitness

```yaml
term: activation_fitness
status: candidate
owner-layer: ecosystem-adaptation
meaning: >
  Future fitness evidence for whether a specific activation combination
  (signals / workflow / intelligence / governance / memory) improved task
  outcome relative to its cognitive cost. In the current Gen3.5 boundary this
  is a placeholder vocabulary term, not a live scoring engine.
affects:
  - architecture/ai-native-cognitive-ecosystem-system.md
  - plans/active/2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md
anti-meaning: >
  Not an automatic activation policy, not a reinforcement signal, and not a
  reason to bypass runtime validation.
introduced-by: plans/active/2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md
```

## cognitive_cost

```yaml
term: cognitive_cost
status: canonical
owner-layer: runtime-cognition
meaning: >
  6-dim cognitive vector 的 summary cost class（LOW / MEDIUM / HIGH）。
  由 thinking / context / execution / knowledge 四個 cost 子項聚合而成，
  作為 commit-msg cognitive contract 的對外摘要欄位。
affects:
  - runtime/cognitive-modes.yaml
  - runtime/cognitive-modes-cost-class.yaml
  - models/cognitive-modes/README.md
anti-meaning: >
  不是 LLM token cost 的數值估算（那是 token_budget_pressure）；
  不是 reasoning latency。
related-terms:
  - { type: aggregates, target: thinking_cost }
  - { type: aggregates, target: context_cost }
  - { type: aggregates, target: execution_cost }
  - { type: aggregates, target: knowledge_cost }
introduced-by: plans/archived/2026-05-22-1629-runtime-cognitive-modes-system.md
```

## cognitive_slice

```yaml
term: cognitive_slice
status: canonical
owner-layer: governance
meaning: >
  可被獨立載入、驗證、路由的最小認知單元（minimum routable cognition unit）。
  每個 slice 有單一 cognitive phase 責任、明確 load_when / do_not_load_when、
  owner_layer 歸屬（workflow / analysis / intelligence 三層之一），並通過該層
  falsifiable membership predicate。
operational-wording: >
  對外文件使用 `execution surface` / `evidence surface` / `loading surface`
  作為較 runtime-oriented 的同義表達。`slice` 為 governance / 內部設計用語。
schema-spec: governance/cognitive-slice-taxonomy.md §1
placement-predicate: governance/cognitive-slice-taxonomy.md §4
affects:
  - governance/cognitive-slice-taxonomy.md
  - knowledge/runtime/routing-registry.yaml (loading_surfaces under hierarchical routes)
  - validation/scenarios/software-delivery/slice-load-scenario-*.yaml
anti-meaning: >
  不是 arbitrary chunk / static partition；不是 file split for cosmetic reasons；
  不是 ecosystem activation graph（Gen 4 vision，非本 term 涵蓋）。
introduced-by: plans/archived/2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md
promoted-at: 2026-05-30 (after Phase 4 validation A/B/C/D/E all PASS proving taxonomy stability)
```

## compression

```yaml
term: compression
status: canonical
owner-layer: runtime-cognition
meaning: >
  Cognitive context compression — 為了在 token 預算或 context window 限制下
  保留語義價值，將 context 壓縮為 summary、checklist 或 index reference。
  與 context_mode 的 SUMMARY_FIRST / CHECKLIST_FIRST 策略對應。
affects:
  - runtime/cognitive-modes.yaml
  - runtime/cognitive-modes-token-budget.yaml
anti-meaning: >
  不是 gzip / zstd 等資料壓縮算法；不是 LLM context window 的硬上限。
introduced-by: plans/archived/2026-05-22-1629-runtime-cognitive-modes-system.md
```

## context_cost

```yaml
term: context_cost
status: candidate
owner-layer: ecosystem-adaptation
meaning: >
  Source-backed reads、graph traversal、memory loading、routing lookup 所
  產生的 context expansion cost。是 cognitive_cost split 的子項之一。
affects:
  - runtime/cognitive-modes.yaml
  - plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
introduced-by: plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
```

## context_mode

```yaml
term: context_mode
status: canonical
owner-layer: runtime-cognition
meaning: >
  Runtime control plane 對 context expansion strategy 的 deterministic enum
  （INDEX_ONLY / SUMMARY_FIRST / CHECKLIST_FIRST / SOURCE_BACKED /
  GRAPH_ASSISTED）。決定 agent 在 task entry 時以何種深度載入 source /
  summary / graph / checklist。是 cognitive mode 6 維 vector 的其中一維。
affects:
  - runtime/cognitive-modes.yaml
  - runtime/cognitive-modes-discovery.yaml
  - knowledge/runtime/routing-registry.yaml
aliases:
  - ctx_mode
anti-meaning: >
  不是 LLM context window 的 token budget（那是 token_budget_pressure 或
  context_cost）；不是 IDE / editor 的 "context menu"。
excludes:
  - discovery_mode
related-terms:
  - { type: related_to, target: execution_mode }
  - { type: related_to, target: discovery_signal }
introduced-by: plans/archived/2026-05-22-1629-runtime-cognitive-modes-system.md
```

## discovery_mode

```yaml
term: discovery_mode
status: candidate
owner-layer: ecosystem-adaptation
meaning: >
  Knowledge discovery strategy enum（STATIC_ROUTE / HEURISTIC_DISCOVERY /
  ARCHAEOLOGY / DOMAIN_MAPPING / TOOL_CAPABILITY_DISCOVERY /
  KNOWLEDGE_GAP_DETECTION）。決定 agent 如何尋找新知識來源；與
  context_mode 互補，但職責分明：context_mode 決定載入深度，discovery_mode
  決定載入哪些來源。
affects:
  - plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
related-terms:
  - { type: related_to, target: knowledge_mode }
introduced-by: plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
```

## discovery_signal

```yaml
term: discovery_signal
status: canonical
owner-layer: runtime-cognition
meaning: >
  Named runtime-cognition trigger declared in `runtime/cognitive-modes-discovery.yaml`.
  A discovery signal maps user keywords, staged file scopes, git/session state,
  active goals or token budget pressure to cognitive mode dimension overrides.
  Cognitive Contract v2 `activation_reason` / compact `Sig:` must cite known
  discovery signal names, not free-form labels.
affects:
  - runtime/cognitive-modes-discovery.yaml
  - runtime/cognitive-modes-phase-integration.yaml
  - runtime/cognitive-modes-governance-integration.yaml
  - scripts/ai-skill-cli/internal/app/hooks.go
anti-meaning: >
  Not a routing registry route, not validation evidence, and not an arbitrary
  prose explanation for why a task felt complex.
related-terms:
  - { type: related_to, target: execution_mode }
  - { type: related_to, target: context_mode }
  - { type: related_to, target: governance_mode }
introduced-by: plans/archived/2026-05-22-1629-runtime-cognitive-modes-system.md
```

## ecosystem

```yaml
term: ecosystem
status: candidate
owner-layer: ecosystem-adaptation
meaning: >
  Cross-layer interaction layer，承載 models / tools / memory / workflow
  之間的 resource interaction、economic pressure、adaptation、feedback。
  不是 source-of-truth layer；source 仍歸屬各原始 layer。
affects:
  - plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
anti-meaning: >
  不是 npm / pip 等 software ecosystem；不是業務領域中的 "vendor ecosystem"。
introduced-by: plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
```

## execution_cost

```yaml
term: execution_cost
status: candidate
owner-layer: ecosystem-adaptation
meaning: >
  Tool calls、mutations、validation、runtime refresh、test runs 所產生的
  execution cost。是 cognitive_cost split 的子項之一。
affects:
  - runtime/cognitive-modes.yaml
  - plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
introduced-by: plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
```

## execution_mode

```yaml
term: execution_mode
status: canonical
owner-layer: runtime-cognition
meaning: >
  Runtime control plane 對 reasoning depth 的 deterministic enum
  （FAST / NORMAL / DEEP / FORENSIC / RECOVERY）。決定 agent 在 task entry
  時投入的 reasoning steps 規模。是 cognitive mode 6 維 vector 的其中一維。
affects:
  - runtime/cognitive-modes.yaml
  - runtime/cognitive-modes-discovery.yaml
aliases:
  - reasoning_mode
anti-meaning: >
  不是 process / thread 的 OS 執行模式；不是 IDE 的 "run / debug" mode。
related-terms:
  - { type: related_to, target: context_mode }
  - { type: related_to, target: discovery_signal }
introduced-by: plans/archived/2026-05-22-1629-runtime-cognitive-modes-system.md
```

## fitness_system

```yaml
term: fitness_system
status: candidate
owner-layer: ecosystem-adaptation
meaning: >
  Future Gen4 contract family for evaluating whether cognition patterns,
  activation sets, memory replay, workflow sequences, or textual skill updates
  improve outcomes relative to their cost. Current usage is interface
  reservation only; no autonomous scoring engine is active.
affects:
  - architecture/ai-native-cognitive-ecosystem-system.md
  - plans/active/2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md
anti-meaning: >
  Not feedback promotion score, not cognitive_cost, and not a live telemetry
  database.
introduced-by: plans/active/2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md
```

## generated_surface

```yaml
term: generated_surface
status: canonical
owner-layer: runtime-projection
meaning: >
  Runtime compiler 從 owner-layer source 投影出的 derived data，例如
  `runtime/runtime.db generated_surfaces[*]` 中的 JSON document。
  consumers 只能讀；不能直接編輯 generated surface 作為 canonical source。
affects:
  - runtime/runtime.db
  - runtime/core-bootstrap.yaml
  - scripts/ai-skill-cli/internal/app/runtime_compiler.go
related-terms:
  - { type: related_to, target: projection }
introduced-by: plans/archived/2026-05-22-0855-executable-yaml-contract-migration.md
```

## governance_mode

```yaml
term: governance_mode
status: canonical
owner-layer: runtime-cognition
meaning: >
  Runtime control plane 對 governance / validation strictness 的 deterministic
  enum（LIGHT / STANDARD / STRICT / LOCKDOWN）。決定 agent 在 task entry 時
  必須通過的 validation gate set。是 cognitive mode 6 維 vector 的其中一維。
affects:
  - runtime/cognitive-modes.yaml
  - runtime/cli-modification-policy.yaml
related-terms:
  - { type: related_to, target: discovery_signal }
introduced-by: plans/archived/2026-05-22-1629-runtime-cognitive-modes-system.md
```

## intelligence_mode

```yaml
term: intelligence_mode
status: candidate
owner-layer: ecosystem-adaptation
meaning: >
  Intelligence activation strategy enum（ATOM_ONLY / WORKFLOW_GUIDED /
  HEURISTIC_ENFORCED / CROSS_INTELLIGENCE / FAILURE_AUGMENTED /
  DOMAIN_REASONING）。決定 agent 在解任務時啟用哪些 intelligence layer。
affects:
  - plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
related-terms:
  - { type: related_to, target: knowledge_mode }
introduced-by: plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
```

## knowledge_cost

```yaml
term: knowledge_cost
status: candidate
owner-layer: ecosystem-adaptation
meaning: >
  Discovery、source refresh、cross-domain synthesis、intelligence activation、
  memory promotion 所產生的 knowledge acquisition cost。是 cognitive_cost
  split 的子項之一。
affects:
  - runtime/cognitive-modes.yaml
  - plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
  - governance/lifecycle/knowledge-update-flow.md
introduced-by: plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
```

## knowledge_mode

```yaml
term: knowledge_mode
status: candidate
owner-layer: ecosystem-adaptation
meaning: >
  Knowledge acquisition strategy enum（REUSE_ONLY / SOURCE_REFRESH /
  DISCOVERY_REQUIRED / CROSS_DOMAIN_SYNTHESIS / FAILURE_LEARNING /
  MEMORY_PROMOTION）。決定 agent 在解任務時對知識的需求型態，並影響
  knowledge-update-flow 的 promotion 動作。
affects:
  - plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
  - governance/lifecycle/knowledge-update-flow.md
related-terms:
  - { type: related_to, target: discovery_mode }
  - { type: related_to, target: intelligence_mode }
introduced-by: plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
```

## memory_mode

```yaml
term: memory_mode
status: canonical
owner-layer: runtime-cognition
meaning: >
  Runtime control plane 對 memory layer 的 activation enum
  （NONE / EPISODIC / DECISION_REPLAY / FAILURE_REPLAY / PROJECT_CONTEXT）。
  決定 agent 在 task entry 時引入哪類 memory replay。是 cognitive mode
  6 維 vector 的其中一維。
affects:
  - runtime/cognitive-modes.yaml
  - memory/README.md
anti-meaning: >
  不是 LLM internal context memory；不是 RAM / disk 的硬體 memory。
introduced-by: plans/archived/2026-05-22-1629-runtime-cognitive-modes-system.md
```

## owner_layer

```yaml
term: owner_layer
status: canonical
owner-layer: architecture-contracts
meaning: >
  Semantic ownership designation — 每個 canonical rule、contract、glossary
  entry 必須宣告其 owner-layer。其他 layer 只能引用、alias 或標記 local
  usage，不得 inline redefine。Owner 與 storage location 解耦：owner 反映
  semantic responsibility，storage 反映 file topology。
affects:
  - knowledge/glossary/README.md
  - architecture/README.md
  - governance/lifecycle/executable-contract-boundary.md
aliases:
  - owner-layer
anti-meaning: >
  不只是檔案存放的資料夾名稱（那是 storage layer）；不是 ACL / RBAC 的
  ownership 概念。
introduced-by: plans/active/2026-05-25-1000-context-language-glossary-system.md
```

## optimization_memory

```yaml
term: optimization_memory
status: candidate
owner-layer: ecosystem-adaptation
meaning: >
  Future memory lifecycle for successful cognition or textual strategy patterns:
  successful execution produces a bounded winning pattern, which may later be
  promoted into a reusable activation heuristic or skill-like instruction after
  validation. Current status is schema/interface reservation only.
affects:
  - architecture/ai-native-cognitive-ecosystem-system.md
  - plans/active/2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md
anti-meaning: >
  Not failure memory, not generic feedback history, and not an automatic
  self-evolving prompt updater.
introduced-by: plans/active/2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md
```

## pressure_model

```yaml
term: pressure_model
status: candidate
owner-layer: ecosystem-adaptation
meaning: >
  Cross-layer pressure interaction model，例如 context_explosion、
  memory_amplification、governance_overhead、validation_fatigue。描述
  source-of-truth layers 互動時產生的 emergent cost，是 ecosystem-adaptation
  的核心建模單位。
affects:
  - plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
introduced-by: plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
```

## projection

```yaml
term: projection
status: canonical
owner-layer: runtime-projection
meaning: >
  從 canonical source（owner-layer YAML / Markdown）派生出的 derived data
  structure，例如 SQLite index、glossary semantic index、generated report。
  Projection 是 read-only consumer surface；修改必須回到 canonical source。
affects:
  - runtime/runtime.db
  - knowledge/runtime/sqlite/runtime-index.sqlite
  - knowledge/glossary/README.md
related-terms:
  - { type: related_to, target: generated_surface }
introduced-by: plans/archived/2026-05-22-0855-executable-yaml-contract-migration.md
```

## runtime_refresh

```yaml
term: runtime_refresh
status: canonical
owner-layer: runtime-projection
meaning: >
  `ai-skill runtime refresh` 命令所執行的流程：從 canonical source
  重建 knowledge runtime reports、SQLite runtime index 與 projection
  tables。預設 native mode（pure Go），不依賴外部 sqlite3 CLI。
affects:
  - scripts/ai-skill-cli/internal/app/runtime.go
  - knowledge/runtime/runtime-report.md
  - knowledge/runtime/sqlite/runtime-index.sqlite
introduced-by: plans/archived/2026-05-21-0834-cross-platform-go-script-runtime.md
```

## rejected_optimization_memory

```yaml
term: rejected_optimization_memory
status: candidate
owner-layer: ecosystem-adaptation
meaning: >
  Future memory lifecycle for optimization candidates that were explicitly
  rejected because they regressed quality, increased token cost, raised
  governance friction, caused telemetry overhead, or over-activated context.
  The purpose is to prevent repeated optimization hallucinations.
affects:
  - architecture/ai-native-cognitive-ecosystem-system.md
  - plans/active/2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md
anti-meaning: >
  Not a generic failure pattern; it records rejected improvement attempts and
  their regression or cost evidence.
introduced-by: plans/active/2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md
```

## thinking_cost

```yaml
term: thinking_cost
status: candidate
owner-layer: ecosystem-adaptation
meaning: >
  Reasoning depth、recursive analysis、validation chain 所產生的 thinking
  cost。是 cognitive_cost split 的子項之一。
affects:
  - runtime/cognitive-modes.yaml
  - plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
introduced-by: plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
```
