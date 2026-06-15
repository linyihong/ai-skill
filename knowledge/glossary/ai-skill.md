# Ai-skill Framework Glossary

> Canonical entries for Ai-skill framework vocabulary, runtime semantics, cognitive vocabulary and architecture contracts.
>
> Schema spec：[`README.md`](README.md)。Validator：`ai-skill glossary validate`。
>
> 上游 plan：[`plans/active/2026-05-25-1000-context-language-glossary-system.md`](../../plans/active/2026-05-25-1000-context-language-glossary-system.md)（Phase 3）。
>
> 編寫慣例：entries 按 term snake_case 字母順序排列；`status: candidate` 標明為 economics plan 預留，待 [`plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md`](../../plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md) Phase 1 確認 owner path 後 promote。

---

## artifact_shape

```yaml
term: artifact_shape
status: candidate
owner-layer: validation-governance
meaning: >
  The proof carrier format for a validation claim, such as screenshot,
  navigation_log, poll_log, or dom_assertion_log. Artifact shape describes
  what the evidence looks like; it is not an evidence_type and must not
  appear in gate requires tokens.
affects:
  - validation/evidence-types/README.md
  - workflow/software-delivery/validation/evidence-gate-vocabulary.md
anti-meaning: >
  Not a gate token, not a collection_method, and not a subtype of
  evidence_type. Reject inheritance trees between type and shape.
introduced-by: plans/archived/2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime.md
related-terms:
  - { type: related_to, target: evidence_type }
  - { type: related_to, target: observable_evidence }
```

## authority_coupled_side_effects

```yaml
term: authority_coupled_side_effects
status: candidate
owner-layer: validation-governance
meaning: >
  Software-delivery validation evidence shape for side-effect claims where a
  low-authority proxy signal, such as a click, API 200, adapter success, log,
  or local counter, must be distinguished from the authority event that owns
  the business truth.
affects:
  - workflow/software-delivery/validation.md
  - workflow/software-delivery/artifact-gates.md
  - validation/scenarios/software-delivery/governance-authority-side-effect-low-authority.yaml
anti-meaning: >
  Not a rule that hardcodes a specific business event or counter. It requires
  naming the business truth, authority event, observable proxy, and evidence
  that proves the claim scope.
introduced-by: plans/active/2026-06-10-1718-software-delivery-governance-invariants.md
```

## authority_decision

```yaml
term: authority_decision
status: candidate
owner-layer: validation-governance
meaning: >
  The classification step that decides who may change what after a validation
  or projection failure: framework invariant, domain pattern, implementation
  defect, or env/deploy incident. Authority decision prevents every incident
  from becoming a framework change.
affects:
  - workflow/software-delivery/validation/authority-decision-table.md
  - workflow/software-delivery/validation/failure-evolution-catalog.md
anti-meaning: >
  Not blame assignment and not a substitute for root-cause analysis. It
  selects the eligible writeback targets for evolution.
introduced-by: plans/archived/2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime.md
related-terms:
  - { type: related_to, target: projection_break }
```

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
owner-layer: validation-governance
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
formalized-in: constitution/ADR-009-cognitive-slice-taxonomy.md (Accepted 2026-05-31)
deferred-items-tracker: constitution/ADR-009-cognitive-slice-taxonomy.md §Future Considerations (each item has explicit Trigger to revisit + Owner at revisit)
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

## configuration_readback_validation

```yaml
term: configuration_readback_validation
status: candidate
owner-layer: validation-governance
meaning: >
  Software-delivery validation evidence shape that separates desired
  configuration input, applied state, actual runtime/deployed readback state,
  and validation evidence before accepting a configuration-applied claim.
affects:
  - workflow/software-delivery/validation.md
  - workflow/software-delivery/artifact-gates.md
  - validation/scenarios/software-delivery/governance-configuration-readback-stale-runtime.yaml
anti-meaning: >
  Not a deploy-env checklist and not proof from editing a config file. It is an
  application of State Visibility Gap / Evidence Chain reasoning to deployed
  configuration claims.
introduced-by: plans/active/2026-06-10-1718-software-delivery-governance-invariants.md
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

## discovery_bridge

```yaml
term: discovery_bridge
status: candidate
owner-layer: runtime-cognition
meaning: >
  The mechanical fallback on the workflow-activation detector's MISS path.
  When deterministic activation finds no route, the Discovery Bridge runs
  Light (pre-Read cheap signals) then optionally Deep (piggyback content)
  discovery to rank candidate routes and inject a non-blocking advisory.
  Converts "detector miss = silent total failure" into "detector miss =
  advisory fallback", without ever becoming an activation path.
affects:
  - scripts/ai-skill-cli/internal/app/discovery.go
  - governance/workflow-activation-engine.md
  - enforcement/failure-patterns/detector-miss-no-fallback.md
anti-meaning: >
  Not an activation mechanism. A Discovery proposal never satisfies
  activation_triggers and never blocks a tool; activation stays limited to
  deterministic detector match or user manual-lock.
related-terms:
  - { type: related_to, target: discovery_signal }
  - { type: part_of, target: light_discovery }
  - { type: part_of, target: deep_discovery }
introduced-by: plans/active/2026-06-06-1700-workflow-activation-discovery-bridge.md
```

## light_discovery

```yaml
term: light_discovery
status: candidate
owner-layer: runtime-cognition
meaning: >
  Phase A of the Discovery Bridge. On detector miss, scores registry routes
  using only pre-Read cheap signals (user message tokens, artifact
  basenames/paths/extensions, frontmatter head bytes ≤200B, cwd, project
  overlay metadata) and emits top-3 candidates with a confidence score. Above
  threshold → inject advisory; below → mark awaiting_phase_b.
affects:
  - scripts/ai-skill-cli/internal/app/discovery.go
anti-meaning: >
  Not a content scan — Light Discovery issues zero new Reads. Reading the
  artifact body is Deep Discovery's job.
related-terms:
  - { type: part_of, target: discovery_bridge }
  - { type: related_to, target: deep_discovery }
introduced-by: plans/active/2026-06-06-1700-workflow-activation-discovery-bridge.md
```

## deep_discovery

```yaml
term: deep_discovery
status: candidate
owner-layer: runtime-cognition
meaning: >
  Phase B of the Discovery Bridge (deferred). Piggybacks the agent's natural
  next artifact Read (hijacks the content stream, issues no new Read) to run a
  content scan that refines the candidate set. Append-only evidence
  accumulation with re-score across the full evidence set, not max() over
  prior confidence, so early Light false positives can be overturned.
affects:
  - scripts/ai-skill-cli/internal/app/discovery.go
anti-meaning: >
  Does not initiate its own Read and does not overwrite prior proposals; it
  accumulates evidence and re-ranks.
related-terms:
  - { type: part_of, target: discovery_bridge }
  - { type: related_to, target: piggyback_read }
  - { type: related_to, target: light_discovery }
introduced-by: plans/active/2026-06-06-1700-workflow-activation-discovery-bridge.md
```

## discovery_proposal

```yaml
term: discovery_proposal
status: candidate
owner-layer: runtime-cognition
meaning: >
  A per-task Discovery Bridge record in runtime.db (`discovery_proposals`):
  top-N candidate routes with scores + evidence_set, a signal_snapshot for
  cross-version re-scoring, scoring_version, status
  (awaiting_phase_b | advised | dismissed | rejected | expired), miss_reason
  enum, and a 24h TTL. Confidence is never stored alone — always paired with
  scoring_version + signal_snapshot.
affects:
  - scripts/ai-skill-cli/internal/app/discovery.go
  - runtime/runtime.db
anti-meaning: >
  Not a routing decision and not a registry route. It is ephemeral runtime
  state, never auto-promoted into routing-registry.yaml.
related-terms:
  - { type: part_of, target: discovery_bridge }
  - { type: related_to, target: advisory_injection }
introduced-by: plans/active/2026-06-06-1700-workflow-activation-discovery-bridge.md
```

## advisory_injection

```yaml
term: advisory_injection
status: candidate
owner-layer: runtime-cognition
meaning: >
  The Discovery Bridge's output mechanism: when a proposal clears threshold,
  a ≤200-token, non-blocking advisory is written into the PreToolUse hook's
  additionalContext, listing top candidate routes and their primary_source so
  the agent may CHOOSE to Read the workflow. Explicitly marked optional /
  non-blocking.
affects:
  - scripts/ai-skill-cli/internal/app/discovery.go
  - scripts/ai-skill-cli/internal/app/hooks.go
anti-meaning: >
  Not a gate, not a deny, not a forced Read. The agent remains free to ignore
  it; activation is never triggered by an advisory.
related-terms:
  - { type: related_to, target: discovery_proposal }
  - { type: part_of, target: discovery_bridge }
introduced-by: plans/active/2026-06-06-1700-workflow-activation-discovery-bridge.md
```

## piggyback_read

```yaml
term: piggyback_read
status: candidate
owner-layer: runtime-cognition
meaning: >
  Deep Discovery's zero-marginal-cost signal source: instead of issuing its
  own Read, it subscribes to the agent's already-occurring artifact Read
  (PostToolUse:Read) and scans that content stream. Each new Read appends
  fresh evidence and triggers a re-score, making Discovery a stream-y,
  append-only accumulation rather than a one-shot decision.
affects:
  - scripts/ai-skill-cli/internal/app/hooks.go
  - scripts/ai-skill-cli/internal/app/discovery.go
anti-meaning: >
  Not a new file read initiated by the runtime; it only observes Reads the
  agent was already going to perform.
related-terms:
  - { type: part_of, target: deep_discovery }
introduced-by: plans/active/2026-06-06-1700-workflow-activation-discovery-bridge.md
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

## experience_runtime

```yaml
term: experience_runtime
status: candidate
owner-layer: workflow-orchestration
meaning: >
  Cross-cutting model for immersive client surfaces that span runtime state,
  journey specification, validation execution, and UI contracts. Documented as
  templates under workflow/cross-cutting/experience-runtime/; not a
  software-delivery slice until multiple surface pilots converge.
affects:
  - workflow/cross-cutting/experience-runtime/README.md
  - workflow/cross-cutting/experience-runtime/player.yaml
anti-meaning: >
  Not a replacement for journey BDD, not an integration runner, and not a
  gate token namespace. Do not register sd-experience-runtime prematurely.
introduced-by: plans/archived/2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime.md
related-terms:
  - { type: related_to, target: journey_specification }
  - { type: related_to, target: validation_scope }
```

## expected_outcomes

```yaml
term: expected_outcomes
status: candidate
owner-layer: validation-governance
meaning: >
  Journey Validation 中要被證明成立的真實狀態或產品結果，例如 entitlement
  active、playback allowed、account state changed。它描述 claim 的 outcome，
  不描述用來證明 outcome 的 artifact。
affects:
  - workflow/software-delivery/test-strategy.md
  - workflow/software-delivery/validation.md
  - workflow/software-delivery/artifact-gates.md
anti-meaning: >
  Not screenshot evidence, not UI copy, not API status code, and not a generic
  readback bucket. Observable proof belongs in observable_evidence.
introduced-by: plans/active/2026-06-10-0908-user-journey-validation-integration.md
```

## evidence_type

```yaml
term: evidence_type
status: candidate
owner-layer: validation-governance
meaning: >
  The validation dimension describing what must be proven for a completion
  claim, such as user_visible, navigation, temporal_behavior, or
  state_persistence. Gate requires tokens use evidence_type only
  (evidence:<name>); collection_method and artifact_shape are separate layers.
affects:
  - validation/evidence-types/README.md
  - workflow/software-delivery/validation/evidence-gate-vocabulary.md
anti-meaning: >
  Not collection_method, not artifact_shape, and not a browser_review
  activity label. Reject token inheritance between type and method/shape.
introduced-by: plans/archived/2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime.md
related-terms:
  - { type: related_to, target: artifact_shape }
  - { type: related_to, target: observable_evidence }
  - { type: related_to, target: validation_capability }
```

## feedback_decision

```yaml
term: feedback_decision
status: canonical
owner-layer: validation-governance
meaning: >
  Feedback / Learning Report 的 close-out decision enum（NONE / NEEDED /
  UNKNOWN），回答本輪是否有值得沉澱的 reusable learning，以及是否有足夠
  evidence 判斷。它不分類 lesson 品質，也不追蹤 promotion lifecycle。
affects:
  - runtime/core-bootstrap.yaml
  - scripts/ai-skill-cli/internal/app/hooks.go
  - plans/archived/2026-06-08-1047-feedback-learning-report-obligation.md
aliases:
  - learning_decision
anti-meaning: >
  不是 knowledge_mode、不是 memory promotion decision、不是 semantic quality score。
related-terms:
  - { type: related_to, target: feedback_learning_report }
  - { type: related_to, target: repo_context }
  - { type: related_to, target: writeback_status }
introduced-by: plans/archived/2026-06-08-1047-feedback-learning-report-obligation.md
```

## feedback_learning_report

```yaml
term: feedback_learning_report
status: canonical
owner-layer: validation-governance
meaning: >
  Final close-out obligation that reports learning disposition for the current
  turn: feedback_decision, repo_context, writeback_status, optional target, and
  a short reason when needed. It makes learning disposition visible without
  forcing a durable lesson on every turn.
affects:
  - runtime/core-bootstrap.yaml
  - CORE_BOOTSTRAP.md
  - scripts/ai-skill-cli/internal/app/hooks.go
  - validation/scenarios/runtime/feedback-report-required-v1.yaml
  - plans/archived/2026-06-08-1047-feedback-learning-report-obligation.md
anti-meaning: >
  Not a knowledge acquisition system, not an economics/telemetry surface, and
  not a replacement for feedback/history lessons or failure learning records.
related-terms:
  - { type: related_to, target: feedback_decision }
  - { type: related_to, target: repo_context }
  - { type: related_to, target: writeback_status }
introduced-by: plans/archived/2026-06-08-1047-feedback-learning-report-obligation.md
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

## journey_specification

```yaml
term: journey_specification
status: candidate
owner-layer: validation-governance
meaning: >
  BDD-owned definition of a project-defined multi-step user outcome, including
  journey name, criticality, triggering action, expected side-effect chain,
  expected outcomes, and observable evidence expectations.
affects:
  - workflow/software-delivery/test-strategy.md
  - workflow/software-delivery/templates/bdd-scenario-template.md
  - workflow/software-delivery/validation.md
anti-meaning: >
  Not a framework canonical journey list, not a UI governance domain, and not
  proof that the journey executed. Execution evidence belongs to Journey
  Validation.
introduced-by: plans/active/2026-06-10-0908-user-journey-validation-integration.md
```

## journey_validation

```yaml
term: journey_validation
status: candidate
owner-layer: validation-governance
meaning: >
  Validation-owned execution and evidence evaluation for a Journey Specification.
  It verifies that the user action produces the expected state transition chain,
  expected outcomes, and observable evidence at the required evidence depth.
affects:
  - workflow/software-delivery/validation.md
  - workflow/software-delivery/execution-flow.yaml
  - workflow/software-delivery/artifact-gates.yaml
  - validation/scenarios/software-delivery/journey-validation-*.yaml
anti-meaning: >
  Not validation_domain, not screen-level UI validation, and not API success by
  itself. It is scoped outcome validation over a project-defined journey.
introduced-by: plans/active/2026-06-10-0908-user-journey-validation-integration.md
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

## observable_evidence

```yaml
term: observable_evidence
status: candidate
owner-layer: validation-governance
meaning: >
  Evidence artifact, readback path, or observable product surface used to prove
  declared expected_outcomes in Journey Validation, such as DB readback,
  profile badge, protected resource access, event record, or external
  confirmation.
affects:
  - workflow/software-delivery/test-strategy.md
  - workflow/software-delivery/validation.md
  - workflow/software-delivery/artifact-gates.md
anti-meaning: >
  Not the expected outcome itself and not a catch-all readback field. It must
  map to a specific outcome claim and evidence scope.
introduced-by: plans/active/2026-06-10-0908-user-journey-validation-integration.md
```

## outcome_validation

```yaml
term: outcome_validation
status: candidate
owner-layer: validation-governance
meaning: >
  Validation concern that asks whether the product or business outcome actually
  happened, beyond local API, adapter, mock, or screen success signals. Journey
  Validation is the current workflow-local expression of this pressure.
affects:
  - workflow/software-delivery/validation.md
  - plans/archived/2026-06-09-1040-experience-validation-pipeline-evolution.md
  - plans/archived/2026-06-10-0908-user-journey-validation-integration.md
anti-meaning: >
  Not a replacement for contract, behavior, accessibility, or responsive
  validation. It complements those targets by requiring evidence that the
  intended outcome materialized.
introduced-by: plans/active/2026-06-10-0908-user-journey-validation-integration.md
```

## ownership_awareness

```yaml
term: ownership_awareness
status: candidate
owner-layer: validation-governance
meaning: >
  Software-delivery intake and closure evidence shape that asks whether a
  change crosses repo, module, security, platform, generated artifact, shared
  library, or team ownership boundaries that require explicit review or
  escalation.
affects:
  - workflow/software-delivery/intake.md
  - workflow/software-delivery/closure.md
  - workflow/software-delivery/review-checklist.md
  - validation/scenarios/software-delivery/governance-ownership-awareness-cross-boundary.yaml
anti-meaning: >
  Not the same as task_scope_validation. A change can be inside the task scope
  while crossing ownership boundaries, or outside task scope without crossing
  owners.
introduced-by: plans/active/2026-06-10-1718-software-delivery-governance-invariants.md
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

## parent

```yaml
term: parent
status: canonical
owner-layer: validation-governance
meaning: >
  Sub-plan / spike frontmatter 指向其主計畫 id 的 pointer，是 plan tree
  hierarchy 的單一 source of truth。Runtime 由 parent pointer 動態建樹，
  不維護反向 children 欄位；main plan 的 parent 為 null。
affects:
  - governance/lifecycle/plan-tree-hierarchy.md
  - plans/README.md
anti-meaning: >
  不是 folder 路徑（folder 只是 UI convention）；不是 sub-plan 之間的
  dependency edge（那是未來 depends_on，不在 plan tree 治理範圍）。
related-terms:
  - { type: related_to, target: plan_tree }
introduced-by: plans/active/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md
```

## plan_kind

```yaml
term: plan_kind
status: canonical
owner-layer: validation-governance
meaning: >
  Plan 的類型 enum：main（主計畫）/ sub（子計畫）/ spike（短期 PoC）。
  決定 frontmatter 必填欄位集合與 validator 行為（sub / spike 須帶 parent /
  required_for_completion / sub_plan_reason）。
affects:
  - governance/lifecycle/plan-tree-hierarchy.md
  - plans/README.md
anti-meaning: >
  不是 plan status（draft / in-progress / completed 是生命週期）；不是 priority。
related-terms:
  - { type: related_to, target: plan_tree }
introduced-by: plans/active/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md
```

## plan_tree

```yaml
term: plan_tree
status: canonical
owner-layer: validation-governance
meaning: >
  主計畫與 sub-plan 由 frontmatter parent pointer 動態建立的階層樹結構。
  `ai-skill plans tree` 純讀 active + archived 兩種狀態的 frontmatter 建樹，
  即使 folder 放錯仍能渲染正確 hierarchy；referential integrity 由 commit-msg
  validator 保證。
affects:
  - governance/lifecycle/plan-tree-hierarchy.md
  - plans/README.md
anti-meaning: >
  不是 DAG（sub-plan 之間若需 depends_on 依賴將引入 graph，不在本治理範圍）；
  不是 folder 目錄結構本身（folder 是 UI convention，parent pointer 才是 truth）。
related-terms:
  - { type: related_to, target: parent }
  - { type: related_to, target: plan_kind }
  - { type: related_to, target: required_for_completion }
  - { type: related_to, target: sub_plan_reason }
introduced-by: plans/active/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md
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

## projection_break

```yaml
term: projection_break
status: candidate
owner-layer: validation-governance
meaning: >
  A structural gap where L2 project concretization (BDD, workflow YAML,
  overlay rules) passes but L3 validation capability and evidence artifacts
  do not exist or were not executed, so UX or runtime behavior can regress
  undetected.
affects:
  - validation/evidence-types/README.md
  - workflow/software-delivery/validation/failure-evolution-catalog.md
  - validation/scenarios/failure-derived/projection-break-missing-browser-evidence-v1.yaml
anti-meaning: >
  Not a single bug class, not "missing unit test" alone, and not equivalent
  to contract drift without a missing validation layer.
introduced-by: plans/archived/2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime.md
related-terms:
  - { type: related_to, target: validation_capability }
  - { type: related_to, target: evidence_type }
  - { type: related_to, target: authority_decision }
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

## repo_context

```yaml
term: repo_context
status: canonical
owner-layer: validation-governance
meaning: >
  Feedback / Learning Report dimension that states whether the agent can
  identify the current repository/project context as LOCAL, NON_LOCAL, or
  UNKNOWN. It is orthogonal to feedback_decision and does not imply git clean,
  pushed, or readback state.
affects:
  - runtime/core-bootstrap.yaml
  - scripts/ai-skill-cli/internal/app/hooks.go
  - plans/archived/2026-06-08-1047-feedback-learning-report-obligation.md
anti-meaning: >
  Not git status, not close-loop completion, and not permission to write
  reusable docs.
related-terms:
  - { type: related_to, target: feedback_learning_report }
  - { type: related_to, target: feedback_decision }
  - { type: related_to, target: writeback_status }
introduced-by: plans/archived/2026-06-08-1047-feedback-learning-report-obligation.md
```

## required_for_completion

```yaml
term: required_for_completion
status: canonical
owner-layer: validation-governance
meaning: >
  Sub-plan frontmatter boolean，描述該 sub-plan 是否屬於 parent 的 acceptance
  criteria。validatePlanTreeArchiveOrder 由此推導 archive blocking：主計畫
  archive 時，所有 required_for_completion: true 的 sub-plan 必須 status:
  completed（只看 status，不看 location）。
affects:
  - governance/lifecycle/plan-tree-hierarchy.md
  - plans/README.md
anti-meaning: >
  不是描述機制的 completion_blocks_parent（已棄用）；描述業務語意，由 validator
  推導 archive blocker，不直接命名機制。
related-terms:
  - { type: related_to, target: plan_tree }
introduced-by: plans/active/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md
```

## side_effect_chain

```yaml
term: side_effect_chain
status: candidate
owner-layer: validation-governance
meaning: >
  Ordered sequence of state transitions expected after a user action in Journey
  Validation, such as order created, payment event recorded, membership updated,
  and entitlement granted.
affects:
  - workflow/software-delivery/test-strategy.md
  - workflow/software-delivery/validation.md
  - intelligence/engineering/execution/validation-reasoning/evidence-chain-validation.md
anti-meaning: >
  Not a list of pages visited and not proof by itself. Each important transition
  still needs evidence appropriate to its risk and scope.
introduced-by: plans/active/2026-06-10-0908-user-journey-validation-integration.md
```

## sub_plan_reason

```yaml
term: sub_plan_reason
status: canonical
owner-layer: validation-governance
meaning: >
  Sub-plan / spike frontmatter 的 free-text 欄位，說明「為什麼拆出這個
  sub-plan」。validatePlanTreeFrontmatter 只強制非空字串，不審內容；取代硬
  enum sub_plan_trigger，避免 framework 隨情境膨脹。
affects:
  - governance/lifecycle/plan-tree-hierarchy.md
  - plans/README.md
anti-meaning: >
  不是 enum（沒有白名單；recommended triggers 只是參考不強制）；不是 plan 的
  Decision Rationale（那是 parent 的職責）。
related-terms:
  - { type: related_to, target: plan_tree }
introduced-by: plans/active/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md
```

## task_scope_validation

```yaml
term: task_scope_validation
status: candidate
owner-layer: validation-governance
meaning: >
  Software-delivery intake evidence shape that checks whether changed surfaces
  still belong to the user-approved task scope, or whether the diff requires a
  scope expansion decision, split, or narrowing before commit.
affects:
  - workflow/software-delivery/intake.md
  - workflow/software-delivery/closure.md
  - validation/scenarios/software-delivery/governance-task-scope-overreach.yaml
anti-meaning: >
  Not ownership_awareness. It answers whether the change belongs to this task,
  not whether the touched surface belongs to another owner or governance
  boundary.
introduced-by: plans/active/2026-06-10-1718-software-delivery-governance-invariants.md
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

## validation_capability

```yaml
term: validation_capability
status: candidate
owner-layer: validation-governance
meaning: >
  The executable layer that produces evidence for a behavior claim: browser
  observation, timing poll, navigation readback, DB readback, or integration
  envelope. Sits between Behavior (L2) and Evidence artifacts (L3 output).
affects:
  - validation/evidence-types/README.md
  - workflow/software-delivery/validation/README.md
anti-meaning: >
  Not BDD source asserts alone, not journey specification text, and not
  synonymous with evidence artifacts.
introduced-by: plans/archived/2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime.md
related-terms:
  - { type: related_to, target: projection_break }
  - { type: related_to, target: evidence_type }
  - { type: related_to, target: journey_validation }
```

## validation_scope

```yaml
term: validation_scope
status: candidate
owner-layer: validation-governance
meaning: >
  The bounded path, audience, journey, state slice, or context slice over which
  validation evidence is claimed. For first Journey Validation landing, journey
  is modeled as validation_scope consuming state, context, and evidence coverage.
affects:
  - workflow/software-delivery/validation.md
  - workflow/software-delivery/execution-flow.yaml
  - plans/active/2026-06-10-0908-user-journey-validation-integration.md
anti-meaning: >
  Not validation_domain and not a quality attribute. It selects what path or
  slice is covered, while validation targets describe what quality or outcome is
  being validated.
introduced-by: plans/active/2026-06-10-0908-user-journey-validation-integration.md
```

## writeback_status

```yaml
term: writeback_status
status: canonical
owner-layer: validation-governance
meaning: >
  Feedback / Learning Report dimension that reports current-turn writeback
  capability/result for a needed durable target: COMPLETED, DEFERRED,
  UNAVAILABLE, or N/A. It describes this turn's close-out state only.
affects:
  - runtime/core-bootstrap.yaml
  - scripts/ai-skill-cli/internal/app/hooks.go
  - plans/archived/2026-06-08-1047-feedback-learning-report-obligation.md
anti-meaning: >
  Not promotion lifecycle tracking, not linked update completion, not
  economics/telemetry, and not evidence that the learning quality was scored.
related-terms:
  - { type: related_to, target: feedback_learning_report }
  - { type: related_to, target: feedback_decision }
  - { type: related_to, target: repo_context }
introduced-by: plans/archived/2026-06-08-1047-feedback-learning-report-obligation.md
```
