---
id: 2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation
plan_kind: main
status: draft
owner: linyihong
created: 2026-05-28
parent: null
---

# Gen4 Fitness & Optimization Memory Interface Reservation

**Status**: `draft`
**世代**：Gen 4 interface reservation（不是 Gen 3 current capability）
**建立日期**：2026-05-28
**最後更新**：2026-05-28

> 本 plan 只預留 Optimization Memory / Fitness System 的 architecture interface，不實作 autonomous optimizer、reinforcement loop、self-modifying governance、automatic workflow mutation 或完整 telemetry DB。其目的在於防止 Gen4 詞彙提前污染 Gen3 runtime，同時讓未來 economics / telemetry / activation contracts 有穩定落點。

---

## Decision Rationale

### Problem & Why Now

Ai-skill 目前已具備強大的 failure-derived evolution：

- failure-derived scenarios
- feedback history
- enforcement failure patterns
- recovery escalation
- forbidden routes
- stale frame / contradiction / routing drift 防護

這些形成 **Cognitive Immunity System**：避免再犯錯。但 SkillOpt-style textual optimization 指向另一條 evolution：成功 pattern 與 rejected candidate 都被保留，讓自然語言 skill / workflow / activation heuristic 能沿著 measured outcome 收斂。

目前缺口：

- 成功 cognition pattern 不會進入 durable memory lifecycle
- rejected optimization 不會以「optimization candidate」身份保留
- activation 組合是否有效沒有 fitness placeholder
- governance chain 是否過重沒有 positive/negative outcome memory
- economics / telemetry plan 會建立 primitive，但缺 downstream fitness taxonomy

現在不應實作 engine，但應預留 contract，避免未來每個 plan 各自發明 `fitness` / `optimization` / `positive evidence` 詞彙。

### Decision

建立 **Gen4-compatible interface reservation**：

1. 在 Gen4 vision 中明確加入 Optimization Memory criterion。
2. 在 Gen3 current architecture 中標註 boundary：fitness/optimizer 不是 current runtime capability。
3. 定義未來 taxonomy：
   - `failure-derived memory`
   - `optimization-derived memory`
   - `activation-derived memory`
   - `suppression-derived memory`
   - `fitness evidence`
4. 定義 minimum schema slots，但不啟用 runtime projection：
   - `fitness.score: unknown`
   - `fitness.evidence: []`
   - `outcome.delta: unknown`
   - `cost.delta: unknown`
   - `candidate.status: accepted | rejected | deferred | unknown`
5. 定義 rejected optimization memory：
   - candidate update
   - rejection reason
   - regression signal
   - cost/friction signal
   - do-not-repeat boundary
6. 在 economics plan 中標為 downstream dependency，不合併 scope。

### Alternatives Considered

- A. 直接做完整 Fitness Engine：reject。runtime trigger audit 與 economics foundation 尚未完成，會製造 telemetry / governance / token debt。
- B. 把 optimization memory 併入 current `feedback/history/`：reject。會把 safety evolution 與 optimization evolution 混在同一 lifecycle，降低未來可測性。
- C. 只在 Gen4 vision 寫概念，不開 plan：reject。會缺少可追蹤 sequencing，economics plan 也不知道 downstream contract。
- D. 先做 interface reservation plan：accept。保留詞彙與 schema slot，但不宣稱 runtime integration。

### Why Not an ADR Yet

此 plan 仍是 vocabulary / schema / sequencing reservation。尚未證明：

- positive optimization evidence 的最小可用 schema
- rejected optimization memory 的 retention / decay 規則
- fitness scoring 是否應屬於 `feedback/`、`memory/`、`ecosystem/optimization/` 或 `runtime/economics/`
- runtime 是否會真的消費這些 signals

待 future implementation plan 產生 executable contract、scenario evidence、runtime consumer 後，再評估 ADR promotion。

### ADR Promotion Criteria（completed 時驗證）

- [ ] Positive optimization memory 和 rejected optimization memory 的 owner path 已決定。
- [ ] 至少一個 successful execution case 能以 bounded pattern 形式被記錄。
- [ ] 至少一個 rejected optimization case 能保留 regression / cost / friction reason。
- [ ] Fitness schema 不與 feedback promotion score、cognitive_cost、runtime economics schema 重疊。
- [ ] 有 runtime consumer 或 scenario 證明 future contract 不是 dead surface。
- [ ] Open Questions 全解。

### Consequences

#### 正面

- 讓 Gen4 optimization vocabulary 有受控入口，不污染 Gen3 current runtime。
- economics / telemetry plan 可把 future fitness 作為 downstream consumer，而不是把 scope 無限擴張。
- 系統開始承認 positive evidence 是 evolution input，不只有 failure evidence。
- Rejected optimization 成為 first-class memory，避免反覆重試「全量 activation」「過度 telemetry」「過度 governance」等 optimization hallucination。

#### 負面

- 新增一個 active plan，短期增加 roadmap 複雜度。
- 需要維護 `current / candidate / vision / forbidden` 邊界，避免概念漂移。

#### 風險

| 風險 | 緩解 |
|---|---|
| Interface reservation 被誤解為已實作 fitness engine | Runtime Execution Path 明確標 doc-only trial；完成條件不包含 runtime integration |
| `feedback/`、`memory/`、`ecosystem/` owner path 重疊 | Phase 1 先做 owner path decision，不先建檔 |
| 過早設計過細 schema | Phase 2 僅保留 placeholder fields，不加 scoring algorithm |
| Optimization memory 變成另一種 failure memory | taxonomy 強制區分 failure-derived / optimization-derived / activation-derived / suppression-derived |

---

## Runtime Execution Path

### Doc-only Trial Statement

目前狀態：**doc-only interface reservation**。本 plan 不新增 runtime generated surface、不新增 discovery signal、不新增 validator、不改 cognitive mode behavior。

### Future Runtime Graduation Entry Conditions

本 plan 只有在以下前置完成後，才允許進入 runtime implementation plan：

1. [`2026-05-28-1200-gen3-runtime-trigger-audit-and-completion.md`](2026-05-28-1200-gen3-runtime-trigger-audit-and-completion.md) 至少完成 audit tooling 與 `validateRuntimeTriggerWiring`。
2. [`2026-05-27-1557-tool-runtime-signal-economics-integration.md`](2026-05-27-1557-tool-runtime-signal-economics-integration.md) 至少完成 economics / telemetry owner boundary 與 initial contract primitive。
3. Future implementation plan 能宣告 named consumer：scenario / runtime query / CLI validator / cognitive state report / telemetry report 之一。

### Future Trigger Flow（not active）

```text
successful or rejected execution
  -> outcome evidence captured
  -> bounded optimization candidate classified
  -> positive / rejected optimization memory retained
  -> fitness placeholder updated
  -> future economics / activation layer may consume as signal
```

### Generated Surfaces

None in this plan. Future candidates:

- `ecosystem.optimization_memory.contract`
- `ecosystem.activation_fitness.contract`
- `ecosystem.rejected_optimization_memory.contract`

These keys are **reserved names only** until a future implementation plan wires named consumers.

### Validation Scenarios

Future candidate scenarios:

- `positive-optimization-memory-retained-v1`
- `rejected-optimization-memory-retained-v1`
- `fitness-placeholder-does-not-claim-score-v1`
- `optimization-memory-not-autonomous-engine-v1`

---

## Phase 0: Pre-Build Interrogation & Architecture Compatibility

- [ ] Confirm this plan remains interface reservation only.
- [ ] Confirm no runtime generated surface is added in this plan.
- [ ] Confirm Gen3 architecture remains current and does not claim fitness engine.
- [ ] Confirm Gen4 vision owns optimization / fitness vocabulary.
- [ ] Confirm economics plan references this plan as downstream, not merged scope.
- [ ] Confirm no new owner layer is created before owner path decision.

完成條件：

- [ ] Architecture docs and economics plan sequencing agree on boundary.

## Phase 1: Define Owner Path Options

- [ ] Evaluate `feedback/optimization/` for evidence lifecycle ownership.
- [ ] Evaluate `memory/optimizations/` for replay / retention ownership.
- [ ] Evaluate future `ecosystem/optimization/` for cross-layer fitness ownership.
- [ ] Evaluate whether `runtime/economics/` should only consume fitness signals, not own them.
- [ ] Record accepted owner path and rejected alternatives.

完成條件：

- [ ] Owner path decision can be used by a future implementation plan without duplicating source-of-truth.

## Phase 2: Reserve Minimal Schema Vocabulary

- [ ] Define `optimization_pattern` placeholder fields.
- [ ] Define `rejected_optimization` placeholder fields.
- [ ] Define `activation_fitness` placeholder fields.
- [ ] Define `fitness.score: unknown` and `fitness.evidence: []` semantics.
- [ ] Define status enum: `accepted`, `rejected`, `deferred`, `unknown`.
- [ ] Define forbidden fields that imply live scoring before implementation.

完成條件：

- [ ] Schema vocabulary is stable enough for future executable contract design, but not projected.

## Phase 3: Define Positive Evidence Flow

- [ ] Define successful execution evidence shape.
- [ ] Define bounded promotion path from successful execution to reusable activation heuristic.
- [ ] Define how positive evidence differs from feedback lesson and intelligence atom.
- [ ] Define minimum evidence required before promoting a winning pattern.

完成條件：

- [ ] Positive evidence is accepted as a first-class evolution input without bypassing governance.

## Phase 4: Define Rejected Optimization Memory

- [ ] Define rejection reasons: regression, token explosion, governance friction, telemetry overhead, activation overreach, workflow inflation.
- [ ] Define retention behavior for rejected optimization candidates.
- [ ] Define do-not-repeat boundary.
- [ ] Define relationship to failure memory.

完成條件：

- [ ] Rejected optimization memory is distinct from failure-derived memory.

## Phase 5: Link to Economics / Telemetry / Suppression

- [ ] Define how economics plan outputs may become future fitness inputs.
- [ ] Define how suppression events may become optimization evidence.
- [ ] Define how telemetry must remain budgeted before fitness scoring can be trusted.
- [ ] Define guardrail: telemetry cost cannot exceed optimization value.

完成條件：

- [ ] Fitness interface can consume future economics signals without requiring a full telemetry DB now.

## Phase 6: Validation and Closure

- [ ] Re-read Gen3 and Gen4 architecture docs for boundary consistency.
- [ ] Re-read economics plan sequencing.
- [ ] Run link/search check for new plan path.
- [ ] Run `ai-skill runtime validate --repo . --json`.
- [ ] If all phases are completed in a future execution, run Plan Completion Closure.

完成條件：

- [ ] No current doc claims autonomous optimization exists.
- [ ] No runtime surface is projected without named consumer.

---

## Open Questions

- Should positive optimization evidence live first under `feedback/`, `memory/`, or future `ecosystem/optimization/`?
- What is the minimal evidence unit for a winning cognition pattern?
- Should rejected optimization memory decay, or remain permanent like failure anti-patterns?
- What distinguishes `fitness_score` from feedback promotion score and cognitive cost?
- Can fitness remain qualitative until telemetry exists, or must it be strictly `unknown`?
- What future CLI or runtime query would consume optimization memory?

## Stakeholder 同意項目

- [ ] This plan reserves interfaces only.
- [ ] Gen3 remains current and does not include autonomous optimization.
- [ ] Positive evidence matters, but does not bypass validation.
- [ ] Rejected optimization memory is first-class.
- [ ] Economics / telemetry primitives must exist before scoring.
- [ ] No full telemetry DB in this plan.
- [ ] No self-modifying governance in this plan.
- [ ] No automatic workflow mutation in this plan.

## 與其他 plans 的關係

- Runs after [`2026-05-28-1200-gen3-runtime-trigger-audit-and-completion.md`](2026-05-28-1200-gen3-runtime-trigger-audit-and-completion.md), because runtime trigger audit prevents new Gen4 surfaces from becoming orphan.
- Runs alongside or after [`2026-05-27-1557-tool-runtime-signal-economics-integration.md`](2026-05-27-1557-tool-runtime-signal-economics-integration.md), because economics / telemetry primitives are future fitness inputs.
- Updates [`architecture/ai-native-cognitive-ecosystem-system.md`](../../architecture/ai-native-cognitive-ecosystem-system.md), because Optimization Memory is Gen4 vision.
- Adds only a boundary note to [`architecture/ai-native-cognitive-execution-system.md`](../../architecture/ai-native-cognitive-execution-system.md), because Gen3 current must not claim fitness engine capability.

## 完成條件

- [ ] Gen4 architecture names Optimization Memory / Rejected Optimization Memory as future criteria.
- [ ] Gen3 architecture explicitly keeps autonomous optimization out of current runtime.
- [ ] Economics active plan references this plan as downstream dependency.
- [ ] Plans README lists this active plan and sequencing.
- [ ] Runtime validation passes.
