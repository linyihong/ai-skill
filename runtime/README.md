# Runtime

Executable runtime layer. Machine-oriented, query-oriented, deterministic.

## 放什麼

- Schema、registry、state definitions、transitions
- Activation graph、routing rules、guard definitions
- Compiled SQLite databases（`runtime.db`, `runtime-state.db`）
- Runtime metadata（budget、health、scheduler）
- Navigation-only README（files、refs、source-of-truth）

## 不放什麼

- 長篇解釋、philosophy、reasoning
- 教學、architecture discussion、design rationale
- Knowledge base 內容、領域知識
- 人類流程說明、操作步驟

這些內容屬於：
- `governance/` — Design philosophy、lifecycle、validation
- `workflow/` — Human-readable execution flows
- `intelligence/` — Heuristics、analytical reasoning
- `enforcement/` — Executable policy rules
- `analysis/` — Analysis 思路

## Domains

| Domain | Path | Description |
|--------|------|-------------|
| Activation | [`router/activation-rules.yaml`](router/activation-rules.yaml) | Lazy-load rules with activation conditions |
| Routing | [`router/`](router/) | Task intent → knowledge index → metadata → source-of-truth gate |
| Discovery | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Phase-aware capability discovery checkpoints (embedded source; compiled to `runtime.db`) |
| Phases | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Execution phase state machine (embedded source; compiled to `runtime.db`) |
| Obligations | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Per-phase atomic duties (embedded source; compiled to `runtime.db`) |
| Gates | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Phase transition prerequisites (embedded source; compiled to `runtime.db`) |
| Compiler | [`compiler/compiler-engine.rb`](compiler/compiler-engine.rb) | Prose → SQLite compilation |
| Runtime DB | [`runtime.db`](runtime.db) | Compiled immutable runtime registry — **must be committed** when changed |
| State DB | [`runtime-state.db`](runtime-state.db) | Mutable execution state (`.gitignore`) |
| Generated | [`generated/`](generated/) | Compiled runtime surfaces (legacy, migrated to SQLite) |
| Transactions | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Writeback transaction state machine (embedded source; compiled to `runtime.db`) |
| Pipeline | [`pipeline/`](pipeline/) | Context flow, guard chain, relevance engine |
| Recovery | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Recovery state machine, escalation levels, output schema, phase reconciliation, state repair, obligation rebuild (embedded source; compiled to `runtime.db`) |
| Scheduler | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Execution queue, priority scheduler (embedded source; compiled to `runtime.db`) |
| Guards | [`guards/`](guards/) | Circuit breaker, context pollution, mismatch escalation |
| Onboarding | [`onboarding/`](onboarding/) | New project/task setup guidance |
| Output Governance | [`output-governance/`](output-governance/) | Language policy, output rules, governance gates |
| Prompt Artifacts | [`prompt-artifacts/`](prompt-artifacts/) | Artifact templates, composition rules |
| Context | [`context/`](context/) | TTL policy and provider prompt cache layout |
| Budget | [`budget/`](budget/) | Token budget |
| Distributed | [`distributed/`](distributed/) | Multi-agent coordination, distributed locks |
| Intelligence | [`intelligence/`](intelligence/) | Intelligence routing |
| Decision recording | [`decisions/`](decisions/) | Close-loop tier routing（ADR / session / project） |

## Recovery Source Map

Runtime recovery 現在不是 standalone `runtime/recovery/*.yaml`。Agent 要處理 blocking gate、phase drift、stale generated surface 或 recovery retry 時，依下列分層讀取：

| Need | Read |
| --- | --- |
| 即時 escalation / recovery output | [`../enforcement/escalation-policy.md`](../enforcement/escalation-policy.md) |
| Retry limit、strategy change、source reload、validation gate | [`../governance/ai-runtime-governance/recovery-retry-governance.md`](../governance/ai-runtime-governance/recovery-retry-governance.md) |
| Domain-specific reload set / forbidden behaviors | [`../metadata/recovery/`](../metadata/recovery/) |
| Machine-readable recovery strategy / phase reconciliation / state repair | [`runtime.db`](runtime.db) tables: `recovery_strategies`, `phase_reconciliation`, `state_repair`, `obligation_rebuild` |
| 修改 runtime recovery 定義 | [`compiler/embedded_data.rb`](compiler/embedded_data.rb)，修改後重新編譯 `runtime.db` |

不要引用已移除的 `runtime/recovery/recovery-strategies.yaml`、`runtime/recovery/phase-reconciliation.yaml`、`runtime/phases/phase-machine.yaml`、`runtime/obligations/obligation-ledger.yaml` 或 `runtime/gates/blocking-gates.yaml` 作為 standalone source。

## Inbound References

- [`route.runtime.activation-rules`](../knowledge/runtime/routing-registry.yaml:77)
- [`route.runtime.context-ttl`](../knowledge/runtime/routing-registry.yaml:102)
- [`route.runtime.context-loading`](../knowledge/runtime/routing-registry.yaml:161)
- [`route.metadata.knowledge-atom-schema`](../knowledge/runtime/routing-registry.yaml:191)
- [`route.models.model-aware-routing`](../knowledge/runtime/routing-registry.yaml:319)
- [`route.runtime.router-flow`](../knowledge/runtime/routing-registry.yaml:348)
- [`route.runtime.context-ttl-doc`](../knowledge/runtime/routing-registry.yaml:407)
- `gate.checkpoint.capability_discovery_completed` in [`runtime.db`](runtime.db) / [`compiler/embedded_data.rb`](compiler/embedded_data.rb)
- `obligation.checkpoint.run_capability_discovery` in [`runtime.db`](runtime.db) / [`compiler/embedded_data.rb`](compiler/embedded_data.rb)

## Databases

Runtime uses two SQLite databases with different lifecycles:

### `runtime.db` (Immutable — Compiled Registry)

Generated by [`compiler-engine.rb`](compiler/compiler-engine.rb) from runtime YAML sources and embedded runtime data in [`compiler/embedded_data.rb`](compiler/embedded_data.rb). When a runtime YAML source exists, the compiler reads that YAML first; embedded data is the fallback for legacy domains whose standalone YAML source has not been restored. Rebuilt on every commit that touches runtime sources. **Do not edit manually.**

Some legacy runtime YAML sources have been embedded into `compiler/embedded_data.rb` and no longer exist as standalone files. When the `Source` column below points to `compiler/embedded_data.rb`, edit the embedded source and recompile `runtime.db`; do not create a new YAML file with the old path unless a dedicated source restoration migration is planned. When the `Source` column points to an existing YAML file, edit the YAML and recompile; do not update only `embedded_data.rb`.

> **⚠️ Commit 規則：`runtime.db` 必須包含在 commit 中。**
> 當 runtime YAML 來源或 compiler 規則變更時，pre-commit hook 會自動重新編譯並 `git add runtime.db`。
> 若手動 commit（跳過 hook），**必須**確認 `runtime.db` 已 `git add`，否則 runtime 與來源不一致。
> 驗證方式：`git diff --cached --name-only | grep runtime.db` 應回傳非空。

| Table | Source | Purpose |
|-------|--------|---------|
| `phases` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Execution phase definitions with entry conditions, allowed/forbidden actions |
| `phase_transitions` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Phase transition rules (blocked transitions, recovery rules) |
| `obligations` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Per-phase atomic duties with verification criteria |
| `gates` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Phase transition prerequisites with severity and failure actions |
| `transaction_states` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Transaction state definitions |
| `transaction_transitions` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Allowed state transitions |
| `transaction_rules` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Transaction rules (lock check, canonical first, etc.) |
| `transaction_templates` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Transaction templates (skill_update, new_skill, feedback_lesson) |
| `activation_rules` | [`router/activation-rules.yaml`](router/activation-rules.yaml) | Lazy-load rule definitions with activation conditions |
| `core_bootstrap_rules` | [`router/activation-rules.yaml`](router/activation-rules.yaml) | Core bootstrap rules (always loaded) |
| `discovery_checkpoints` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Phase-aware capability discovery checkpoints |
| `discovery_search_strategy` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Search strategy configuration |
| `generated_surfaces` | Compiled from prose sources | Extracted structured data from workflow, enforcement, governance documents |
| `compiler_metadata` | Auto-generated | Compiler version, compilation timestamp, schema version |
| `runtime_budget` | [`budget/token-budget.yaml`](budget/token-budget.yaml) | Per-model token budget configuration |
| `context_ttl_policy` | [`context/ttl-policy.yaml`](context/ttl-policy.yaml) | Context TTL policy per context type |
| `circuit_breaker` | [`guards/circuit-breaker.yaml`](guards/circuit-breaker.yaml) | Circuit breaker and mismatch escalation guard definitions |
| `context_pollution` | [`guards/context-pollution.yaml`](guards/context-pollution.yaml) | Context pollution detection signals |
| `context_health_score` | [`health/context-health-score.yaml`](health/context-health-score.yaml) | Context health scoring dimensions |
| `intelligence_routing` | [`intelligence/intelligence-routing.yaml`](intelligence/intelligence-routing.yaml) | Intelligence routing rules |
| `obligation_ledger` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Obligation ledger entries |
| `language_policy` | [`output-governance/language-policy.yaml`](output-governance/language-policy.yaml) | Language enforcement rules |
| `output_rules` | [`output-governance/output-rules.yaml`](output-governance/output-rules.yaml) | Document output formatting rules |
| `governance_gates` | [`output-governance/governance-gates.yaml`](output-governance/governance-gates.yaml) | Output governance blocking gates |
| `blocking_gates` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Blocking gates (runtime config mirror) |
| `phase_machine` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Phase machine (runtime config mirror) |
| `pipeline_context_flow` | [`pipeline/context-flow.yaml`](pipeline/context-flow.yaml) | Progressive context expansion levels |
| `guard_chain` | [`pipeline/guard-chain.yaml`](pipeline/guard-chain.yaml) | Guard execution order per stage |
| `relevance_engine` | [`pipeline/relevance-engine.yaml`](pipeline/relevance-engine.yaml) | Skill relevance scoring configuration |
| `session_lifecycle` | [`pipeline/session-lifecycle.yaml`](pipeline/session-lifecycle.yaml) | Session lifecycle stage definitions |
| `prompt_artifact_templates` | [`prompt-artifacts/artifact-templates.yaml`](prompt-artifacts/artifact-templates.yaml) | Task type prompt artifact templates |
| `prompt_composition_rules` | [`prompt-artifacts/composition-rules.yaml`](prompt-artifacts/composition-rules.yaml) | Prompt composition rules |
| `recovery_strategies` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Recovery state machine, escalation levels, output schema, and strategy definitions |
| `state_repair` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | State repair procedures |
| `obligation_rebuild` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Obligation rebuild procedures |
| `phase_reconciliation` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Phase reconciliation procedures |
| `execution_queue` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Task execution queue configuration |
| `priority_scheduler` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Priority scheduling levels |
| `activation_rules_mirror` | [`router/activation-rules.yaml`](router/activation-rules.yaml) | Activation rules (runtime config mirror) |
| `transaction_templates_ext` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Transaction templates (runtime config mirror) |
| `distributed_locks` | [`distributed/distributed-locks.yaml`](distributed/distributed-locks.yaml) | Distributed lock/lease definitions |
| `multi_agent_coordination` | [`distributed/multi-agent-coordination.yaml`](distributed/multi-agent-coordination.yaml) | Multi-agent coordination rules |
| `async_job_lifecycle` | [`distributed/async-job-lifecycle.yaml`](distributed/async-job-lifecycle.yaml) | Async job lifecycle state definitions |
| `capability_checkpoints` | [`compiler/embedded_data.rb`](compiler/embedded_data.rb) | Capability discovery checkpoints (runtime config mirror) |

### `runtime-state.db` (Mutable — Execution State)

Future mutable execution-state surface. The old Ruby initializer has been removed to avoid accidental use; do not create or mutate this database until a Go-native `ai-skill runtime state ...` command and validation contract exist.

| Table | Purpose |
|-------|---------|
| `execution_state` | Current phase, status, sub-phase, metadata |
| `obligation_status` | Per-obligation completion tracking (pending/in_progress/completed/blocked/skipped) |
| `transaction_state` | Active transaction lifecycle tracking |
| `execution_log` | Append-only event log |

### Query Examples

```sql
-- What phase am I in?
SELECT phase, status FROM execution_state ORDER BY id DESC LIMIT 1;

-- Which obligations are still pending?
SELECT obligation_id, phase FROM obligation_status WHERE status = 'pending';

-- What gates block the current phase?
SELECT g.name, g.severity FROM gates g
JOIN phases p ON g.phase = p.id
WHERE p.id = (SELECT phase FROM execution_state ORDER BY id DESC LIMIT 1);

-- What transitions are allowed from the current state?
SELECT to_state, condition FROM transaction_transitions
WHERE from_state = (SELECT state FROM transaction_state ORDER BY id DESC LIMIT 1);
```

## Source-of-Truth

Runtime does not hold conceptual explanations. Source-of-truth for runtime design:

- `governance/` — Design philosophy, lifecycle, validation
- `workflow/` — Human-readable execution flows
- `intelligence/` — Heuristics, analytical reasoning
- `enforcement/` — Executable policy rules
