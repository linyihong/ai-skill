# Model-Aware Execution Routing

> **狀態**: completed

> **建立日期**: 2026-05-20
> **目的**: 將 `models/` 從 model profile / compression documentation 升級為 model-aware execution strategy layer，讓 agent 能依 task complexity、cognitive state、autonomy mode、context budget 與 tool capability 選擇合適的 execution strategy；若工具支援 explicit model selection，再進一步選用具體模型或 subagent。

---

## 1. Problem Statement

目前 `models/` 已有：

- `models/profiles/README.md`：定義 `small`、`large`、`specialized` profile 與 context loading 深度。
- `models/compression/README.md`：定義 `index-only`、`summary-first`、`checklist-first`、`source-backed`、`graph-assisted` 等壓縮策略。
- `knowledge/runtime/model-context-report.md` 與 `model-checklists.md`：generated model-aware context loading view。
- `knowledge/runtime/routing-registry.yaml` 的 `route.models.model-aware-routing`。

缺口不是 profile 本身錯，而是 profile 尚未變成 execution routing contract。Agent 目前可以知道「讀多深」，但不知道：

1. 哪種 cognitive complexity 應該使用哪種 execution strategy。
2. 何時從 checklist execution 升級到 source-backed analysis。
3. 何時使用 subagent / explicit model selection，何時只能做 behavior-only adaptation。
4. 何時因 confidence decay、contamination、uncertainty、alignment required 而改變 model strategy。
5. Workflow 如何依 model capability 改寫成 small-model / large-model / coding-model / architecture-model / validation-model 的不同執行形狀。

本 plan 目標是讓 `models/` 成為「model-aware execution strategy layer」，而不是特定 provider 的 model picker。

---

## 2. Architecture Compatibility Preflight

| Field | Content |
| --- | --- |
| Trigger | 使用者指出 `models/` 目前像 documentation layer，缺少 execution routing、workflow adaptation 與 runtime routing contract，要求建立 plan、檢查是否與現行框架衝突、寫入建議 active plan 執行順序，並 commit / push。 |
| Checked sources | `models/README.md`、`models/profiles/README.md`、`models/compression/README.md`、`knowledge/runtime/routing-registry.yaml` 的 `route.models.model-aware-routing`、`runtime/README.md`、`plans/README.md`、`plans/archived/2026-05-20-1501-cognitive-state-evidence-governance.md`、`plans/archived/2026-05-20-1745-memory-retrieval-activation-governance.md`。 |
| Conflicts | 不應宣稱 `models/` 能強制 Cursor Auto 或任何工具的主對話模型切換。Actual model selection 屬於 tool capability / `ai-tools/` 邊界；`models/` 應定義 tool-neutral execution strategy。若工具支援 explicit model selection 或 subagent model selection，才可由 routing contract 建議具體模型。 |
| Decision | Proceed as separate active plan。Scope 應聚焦 model-aware execution strategy、routing contract、workflow adaptation、context budget orchestration 與 governance。不要直接新增 provider-specific model names 到 reusable core；必要時只在 tool adapter 或 runtime capability matrix 記錄可用性。 |
| Validation | Plan readback、diff review、ReadLints、Markdown link check。後續若修改 `knowledge/runtime/routing-registry.yaml`、generated model reports 或 runtime source，必須執行 `ai-skill runtime refresh` 與相關 validator。 |

### 2.1 與現有 Active Plan 的切分

| Concern | Cognitive State & Evidence Governance | Memory Retrieval & Activation Governance | 本 Model Plan |
| --- | --- | --- | --- |
| Cognitive state | 定義 `STABLE`、`UNCERTAIN`、`DEGRADED`、`CONTAMINATED`、`ALIGNMENT_REQUIRED` 等 state 與 autonomy modes。 | 使用 cognitive-state 的 confidence / contamination 語意來決定 memory replay 是否可用。 | 將 cognitive state 對應到 model-aware execution strategy，例如 validation-heavy、source-backed、rediscovery-only、human-facing summary。 |
| Runtime reduction | 定義 minimal runtime principle、signal normalization、避免過度 runtime 化。 | 避免 memory activation 過早變成 runtime guard。 | 避免 model routing 變成 provider-specific runtime lock；先定義 strategy primitives，再決定是否進 routing registry。 |
| Context contamination | 定義 stale frame / old checklist / prior route 的污染治理。 | 定義 memory replay contamination boundary。 | 定義 contaminated state 下的 model strategy：rediscovery-only、source-backed、no speculative patching。 |
| Execution contract | 檢查 current action 是否服務原始 goal / validation target。 | 防止 `.agent-goals/` 長期化成 memory。 | 讓 workflow 依 model capability 調整 execution shape，但不取代 `.agent-goals/` 或 workflow success criteria。 |
| Tool model selection | 不處理。 | 不處理。 | 明確區分 execution strategy routing 與 actual model selection；主對話 Auto 不保證可控，subagent / tool-supported model selection 才可用 explicit model。 |

### 2.2 Recommended Active Plan Execution Order

建議順序：

1. 先執行 `plans/archived/2026-05-20-1501-cognitive-state-evidence-governance.md`（已完成 / archived）。
   - 原因：model-aware execution routing 需要引用 cognitive state、autonomy modes、confidence integrity、contamination、runtime reduction 與 minimal runtime principle。
   - 若上游 state model 尚未定穩，model strategy 容易變成另一套平行 governance。
2. 再執行本 plan 的 Phase 0-2。
   - 先建立 model-aware execution contract、capability dimensions、tool capability boundary。
   - 此階段可以與 memory plan Batch A 平行，但不得把 model routing runtime 化。
3. 接著執行 `plans/archived/2026-05-20-1745-memory-retrieval-activation-governance.md` 的 Batch A / Phase 0-2。
   - 原因：memory working buffer 與 model context budget 會互相影響，但兩者都依賴 cognitive-state 的上游治理。
4. 最後執行本 plan 的 Phase 3-6 與 memory plan 的 Batch B / C。
   - 先完成 model routing primitives 與 workflow adaptation，再接 knowledge/runtime routing、generated reports、validation scenarios。

Blocking note：

- Cognitive-state plan 已完成；後續執行本 plan 時應讀 archived plan 與 `governance/ai-runtime-governance/cognitive-state-governance.md` 作為上游語意。
- Autonomy-mode-to-model-strategy、runtime model routing primitives、fallback routing、multi-model handoff 的 runtime promotion，仍必須遵守 cognitive-state 的 runtime reduction / signal normalization，不得建立平行治理。

---

## 3. Target Layer Model

`models/` 的成熟定位：

```text
task intent / cognitive state / autonomy mode
  ↓
models/routing/
  execution strategy selection
  ↓
models/workflow-adaptation/
  workflow shape adjustment
  ↓
models/compression/
  context budget and loading depth
  ↓
tool capability check
  behavior-only adaptation OR explicit model / subagent selection
```

### 3.1 Execution Strategy vs Actual Model Selection

| Layer | Can Govern | Cannot Guarantee |
| --- | --- | --- |
| `models/` | Execution strategy、context depth、workflow adaptation、handoff contract、fallback rule。 | Provider Auto picker、main chat model switching、tool UI state。 |
| `ai-tools/` | Tool-specific model selector、available model names、subagent selection constraints、UI / SDK behavior。 | Tool-neutral model strategy semantics。 |
| `runtime/` | Machine-readable routing primitive、context budget、generated model checklist。 | Conceptual model philosophy 或 provider-specific claims without source。 |
| `workflow/` | How a task should execute under a given strategy。 | Model capability truth。 |

Rule：

> `models/` 決定 agent 應採用哪種 cognitive execution strategy；只有在 tool capability 明確支援時，才進一步指定 actual model 或 subagent model。

### 3.2 Capability Dimensions

現有 `small` / `large` / `specialized` profile 應保留，但需要補 execution-oriented capability dimensions：

| Dimension | Meaning |
| --- | --- |
| `reasoning_depth` | 是否適合 architecture、tradeoff、contradiction propagation、long-form planning。 |
| `context_stability` | 是否能在長 context 下保持指令、目標與證據鏈穩定。 |
| `instruction_following` | 是否適合嚴格遵守 checklist、format、patch scope。 |
| `long_chain_reliability` | 是否能跨多 step execution 保持 state。 |
| `diff_precision` | 是否適合小範圍程式碼 / 文件 patch。 |
| `spec_alignment` | 是否能把需求、驗證、契約與產出對齊。 |
| `hallucination_tolerance` | 低容忍任務需要 source-backed / validation-heavy strategy。 |
| `compression_resilience` | 是否能從 summary / checklist 中恢復足夠 execution context。 |
| `tool_call_reliability` | 是否適合多工具回圈、驗證、commit / push。 |
| `exploration_capability` | 是否適合未知 codebase / architecture discovery。 |

### 3.3 Cognitive State to Model Strategy

初版 mapping：

| Cognitive State | Model Strategy | Required Behavior |
| --- | --- | --- |
| `STABLE` | execution-heavy | 可使用 bounded edits、validation loop、checklist execution。 |
| `UNCERTAIN` | validation-heavy | 先讀 source、收集 evidence、避免 broad patch。 |
| `DEGRADED` | source-backed | 降低 autonomy，使用 primary source + validation gate。 |
| `CONTAMINATED` | rediscovery-only | 不沿用舊 route / memory / checklist；重新 discovery。 |
| `MISALIGNED` | goal-realignment | 回到 user goal / `.agent-goals` / workflow success criteria。 |
| `RECOVERY` | recovery-specialized | 使用 recovery workflow，不進行 unrelated optimization。 |
| `VALIDATION_REQUIRED` | validation-only | 只能收集 evidence、跑 validator、比對 source-of-truth。 |
| `ALIGNMENT_REQUIRED` | human-facing summary | 摘要 options / blockers，等待使用者決策。 |
| `READ_ONLY` | inspection-only | 只讀取、搜尋、分析，不寫檔、不 commit。 |

---

## 4. Proposed Directory Structure

Candidate structure：

```text
models/
├── README.md
├── profiles/
├── capabilities/
│   ├── README.md
│   ├── reasoning-depth.md
│   ├── context-stability.md
│   ├── long-chain-reliability.md
│   ├── tool-reliability.md
│   ├── hallucination-risk.md
│   ├── compression-resilience.md
│   └── planning-capability.md
├── compression/
├── routing/
│   ├── README.md
│   ├── task-routing.md
│   ├── escalation-routing.md
│   ├── fallback-routing.md
│   ├── multi-model-handoff.md
│   └── autonomy-routing.md
├── workflow-adaptation/
│   ├── README.md
│   ├── small-model-workflows.md
│   ├── large-model-workflows.md
│   ├── coding-workflows.md
│   ├── architecture-workflows.md
│   ├── validation-workflows.md
│   └── exploratory-workflows.md
├── governance/
│   ├── README.md
│   ├── model-selection-governance.md
│   ├── hallucination-boundaries.md
│   ├── context-budget-governance.md
│   └── model-confidence-governance.md
└── runtime/
    ├── README.md
    ├── routing-primitives.md
    ├── context-budgeting.md
    ├── execution-cost-strategy.md
    └── adaptive-loading.md
```

### 4.1 Minimal First Slice

不要一次建立所有文件。第一批應先建立最小可執行 contract：

- `models/routing/README.md`
- `models/routing/task-routing.md`
- `models/routing/autonomy-routing.md`
- `models/workflow-adaptation/README.md`
- `models/governance/model-selection-governance.md`
- `models/runtime/routing-primitives.md`

其餘 capability atoms 可在 validation 或 recurring task 需要時再補。

---

## 5. Routing Contract

### 5.1 Task Complexity to Execution Strategy

| Task Class | Strategy | Context Loading |
| --- | --- | --- |
| trivial answer / lookup | checklist-first | index-only / summary-first。 |
| bounded doc edit | source-backed execution | primary source + touched dependencies。 |
| code patch | coding workflow | source-backed + tests / lints。 |
| architecture planning | exploratory planning | source-backed + related architecture / governance / routing docs。 |
| migration / promotion / deprecation | graph-assisted | primary source + dependencies + lifecycle + routing / metadata。 |
| recovery / contradiction | validation-heavy / recovery-specialized | source-of-truth reload + evidence comparison。 |
| long context handoff | summary-first then source-backed | summary / goal ledger / plan + selected source reread。 |

### 5.2 Tool Capability Gate

Before claiming actual model selection, agent must determine:

| Capability | Behavior |
| --- | --- |
| Main chat model fixed or Auto | Use behavior-only adaptation; do not claim model switched. |
| Subagent supports explicit model | Route deep analysis / coding / review to available model if user requested or strategy requires. |
| Tool exposes model selector but agent cannot control it | Document recommended model class; ask user to switch if necessary. |
| Tool model unavailable | Do not substitute silently; report available models or fallback to behavior-only adaptation. |

### 5.3 Subagent / Multi-Model Handoff

Use explicit model subagent only when:

- User requests a model or model class.
- Task complexity exceeds current execution profile.
- Independent deep analysis can run without mutating shared files.
- Handoff summary can preserve source paths, assumptions, validation targets, and open questions.

Do not use subagent when:

- The task is small and direct.
- It would duplicate context cost without improving confidence.
- The subagent cannot access required source or tool capability.
- Multiple agents editing same files would violate `.agent-goals` / lock / single-owner boundaries.

---

## 6. Suggested Implementation Phases

### Phase 0 — Boundary Confirmation

Status: completed.

Tasks:

- [x] Confirm `models/` controls execution strategy, not guaranteed provider model selection.
- [x] Confirm tool-specific model names and selector behavior belong in `ai-tools/`.
- [x] Confirm model routing does not replace `workflow/`, `.agent-goals/`, `runtime/`, or cognitive-state governance.
- [x] Confirm cognitive-state plan provides autonomy / state semantics before runtime model routing promotion.
- [x] Confirm memory plan uses model routing only for context budget / replay strategy after memory boundaries are stable.

Exit criteria:

- [x] Execution strategy vs actual model selection boundary written.
- [x] Active plan execution order documented.
- [x] Candidate directories scoped to minimal first slice.

### Phase 1 — Models Routing Contract

Candidate files:

- `models/routing/README.md`
- `models/routing/task-routing.md`
- `models/routing/autonomy-routing.md`
- `models/routing/fallback-routing.md`
- `models/routing/multi-model-handoff.md`

Tasks:

- [x] Define task class → strategy mapping.
- [x] Define cognitive state / autonomy mode → strategy mapping.
- [x] Define fallback routing when explicit model selection is unavailable.
- [x] Define subagent / multi-model handoff contract.
- [x] Update `models/README.md` to point to routing.

### Phase 2 — Capability Dimensions

Candidate files:

- `models/capabilities/README.md`
- `models/capabilities/reasoning-depth.md`
- `models/capabilities/context-stability.md`
- `models/capabilities/tool-reliability.md`
- `models/capabilities/hallucination-risk.md`
- `models/capabilities/compression-resilience.md`

Tasks:

- [x] Replace overly broad `small` / `large` assumptions with capability dimensions.
- [x] Mark unverified capability claims as confidence-scoped.
- [x] Define how capability dimensions influence context loading and workflow adaptation.

### Phase 3 — Workflow Adaptation

Candidate files:

- `models/workflow-adaptation/README.md`
- `models/workflow-adaptation/small-model-workflows.md`
- `models/workflow-adaptation/large-model-workflows.md`
- `models/workflow-adaptation/coding-workflows.md`
- `models/workflow-adaptation/architecture-workflows.md`
- `models/workflow-adaptation/validation-workflows.md`

Tasks:

- [x] Define checklist-first workflow shape for small / constrained models.
- [x] Define source-backed / exploratory workflow shape for high reasoning tasks.
- [x] Define coding workflow boundaries: diff precision, tests, lints, patch scope.
- [x] Define architecture workflow boundaries: planning, tradeoff, contradiction analysis.
- [x] Define validation workflow boundaries: evidence-first, claim scope, no premature success.

### Phase 4 — Governance and Tool Boundary

Candidate files:

- `models/governance/README.md`
- `models/governance/model-selection-governance.md`
- `models/governance/hallucination-boundaries.md`
- `models/governance/context-budget-governance.md`
- `models/governance/model-confidence-governance.md`
- `ai-tools/agent/cursor.md` or other adapters only if tool-specific behavior must be documented.

Tasks:

- [x] Define model selection governance without provider-specific assumptions.
- [x] Define when agent must say model selection is unavailable or behavior-only.
- [x] Define hallucination boundary for low-confidence model outputs.
- [x] Define context budget governance for long tasks.
- [x] Add tool adapter notes only where actual tool behavior matters.

### Phase 5 — Runtime / Knowledge Integration

Candidate files:

- `models/runtime/README.md`
- `models/runtime/routing-primitives.md`
- `models/runtime/context-budgeting.md`
- `models/runtime/execution-cost-strategy.md`
- `models/runtime/adaptive-loading.md`
- `knowledge/runtime/routing-registry.yaml`
- `knowledge/summaries/model-routing.md`
- `knowledge/runtime/model-context-report.md` and `model-checklists.md` if regenerated.

Tasks:

- [x] Define minimal model routing primitives.
- [x] Update `route.models.model-aware-routing` candidate sources.
- [x] Decide whether generated model reports need new fields for execution strategy.
- [x] Run knowledge runtime refresh if routing registry or generated views change.
- [x] Avoid provider-specific runtime state unless tool capability source exists.

### Phase 6 — Validation Scenarios

Candidate files:

- `validation/scenarios/models/`

Tasks:

- [x] Test Auto/fixed main model fallback to behavior-only adaptation.
- [x] Test user-requested explicit model routes to subagent when available.
- [x] Test unavailable model does not get silently substituted.
- [x] Test uncertain state selects validation-heavy strategy.
- [x] Test contaminated state selects rediscovery-only strategy.
- [x] Test small-model workflow uses checklist-first and bounded context.
- [x] Test architecture task uses source-backed / exploratory strategy.
- [x] Test model routing does not override source-of-truth or enforcement rules.

### Phase 7 — Plan Completion Closure

Tasks:

- [x] Confirm all phases complete or marked blocked.
- [x] Run ReadLints, Markdown link check, and `ai-skill runtime refresh` if routing / generated surfaces changed.
- [x] Check linked updates.
- [x] Update `plans/README.md` status.
- [x] Move plan to `plans/archived/` if completed.
- [x] Commit / push / readback / clean status.

---

## 7. Concrete Work Breakdown

### Batch A — Contract and Routing Skeleton

Scope:

- Create `models/routing/`.
- Update `models/README.md`.
- Define execution strategy vs actual model selection.
- Define task class / autonomy mode routing.

Why first:

- Gives Auto mode an executable behavior contract even when actual model switching is unavailable.

### Batch B — Workflow Adaptation and Governance

Scope:

- Create `models/workflow-adaptation/`.
- Create `models/governance/`.
- Define small / large / coding / architecture / validation workflow shapes.
- Define model selection unavailable / behavior-only rules.

Why second:

- Lets workflow execution change based on model capability without claiming provider control.

### Batch C — Runtime / Knowledge / Validation

Scope:

- Update routing registry / summaries / generated reports if needed.
- Add validation scenarios.
- Add runtime primitives only after cognitive-state runtime reduction is stable.

Why last:

- Prevents model routing from becoming premature runtime state or provider-specific lock-in.

---

## 8. Open Questions

1. Should actual provider model names appear anywhere in reusable `models/`?
   - Current recommendation: no. Put provider-specific names and availability in `ai-tools/` or tool adapter docs.
2. Should `models/routing/` be runtime-enforced immediately?
   - Current recommendation: no. Start as design-layer routing contract, then promote minimal primitives after validation.
3. Should subagent model selection be automatic?
   - Current recommendation: only when user requests it, task complexity justifies it, or routing contract marks it necessary and tool capability is known.
4. Should `small` / `large` profiles remain?
   - Current recommendation: yes, but treat them as coarse context loading profiles; add capability dimensions for execution decisions.
5. Should model-aware routing depend on cognitive-state governance?
   - Current recommendation: yes. Autonomy and state semantics should come from cognitive-state plan, not be redefined in `models/`.

---

## 9. Completion Definition

This plan is complete when:

- `models/` clearly distinguishes execution strategy from actual model selection.
- Auto / fixed model mode has behavior-only adaptation rules.
- Tool-supported explicit model selection has a safe handoff contract.
- Task complexity maps to execution strategy.
- Cognitive state / autonomy mode maps to model strategy.
- Context budget maps to compression mode.
- Workflow adaptation exists for small / large / coding / architecture / validation tasks.
- Model selection governance prevents unsupported model claims and silent substitution.
- Routing registry / summaries / generated reports reflect the new model-aware routing surfaces if promoted.
- Validation scenarios cover unavailable models, explicit model requests, contaminated state, uncertain state, small-model workflow, architecture workflow, and source-of-truth override prevention.
