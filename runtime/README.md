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
| Activation | `runtime.db` tables `activation_rules`, `activation_rules_mirror`, `core_bootstrap_rules` | Lazy-load rules with activation conditions |
| Routing | `runtime.db` activation/runtime config tables | Runtime routing support; knowledge routing remains in `knowledge/runtime/` |
| Discovery | `runtime.db` tables `discovery_checkpoints`, `capability_checkpoints` | Phase-aware capability discovery checkpoints |
| Phases | `runtime.db` tables `phases`, `phase_machine`, `phase_transitions` | Execution phase state machine |
| Obligations | `runtime.db` tables `obligations`, `obligation_ledger` | Per-phase atomic duties |
| Gates | `runtime.db` tables `gates`, `blocking_gates` | Phase transition prerequisites |
| Compiler | `runtime.db` tables `compiler_rules`, `runtime_config_documents`, `runtime_config_projections` + `ai-skill runtime compile` | SQLite canonical config + deterministic prose → SQLite projections |
| Runtime DB | [`runtime.db`](runtime.db) | Compiled immutable runtime registry — **must be committed** when changed |
| State DB | [`runtime-state.db`](runtime-state.db) | Mutable execution state (`.gitignore`) |
| Generated | [`generated/`](generated/) | Compiled runtime surfaces (legacy, migrated to SQLite) |
| Transactions | `runtime.db` tables `transaction_*` | Writeback transaction state machine and templates |
| Pipeline | `runtime.db` tables `pipeline_context_flow`, `guard_chain`, `relevance_engine`, `session_lifecycle` | Context flow, guard chain, relevance engine |
| Recovery | `runtime.db` tables `recovery_strategies`, `state_repair`, `obligation_rebuild`, `phase_reconciliation` | Recovery strategy, phase reconciliation, state repair, obligation rebuild |
| Scheduler | `runtime.db` tables `execution_queue`, `priority_scheduler` | Execution queue, priority scheduler |
| Guards | `runtime.db` tables `circuit_breaker`, `context_pollution` | Circuit breaker, context pollution, mismatch escalation |
| Onboarding | [`onboarding/`](onboarding/) | New project/task setup guidance |
| Output Governance | `runtime.db` tables `language_policy`, `output_rules`, `governance_gates` | Language policy, output rules, governance gates |
| Prompt Artifacts | `runtime.db` tables `prompt_artifact_templates`, `prompt_composition_rules` | Artifact templates, composition rules |
| Context | `runtime.db` table `context_ttl_policy` | TTL policy and provider prompt cache layout |
| Budget | `runtime.db` table `runtime_budget` | Token budget |
| Distributed | `runtime.db` tables `distributed_locks`, `multi_agent_coordination`, `async_job_lifecycle` | Multi-agent coordination, distributed locks |
| Intelligence | `runtime.db` table `intelligence_routing` | Intelligence routing |
| Decision recording | `runtime.db` table `decision_recording` | Close-loop tier routing（ADR / session / project） |

## Recovery Source Map

Runtime recovery 的 machine-readable source 已收斂到 `runtime/runtime.db` canonical documents。Agent 要處理 blocking gate、phase drift、stale generated surface 或 recovery retry 時，依下列分層讀取：

| Need | Read |
| --- | --- |
| 即時 escalation / recovery output | [`../enforcement/escalation-policy.md`](../enforcement/escalation-policy.md) |
| Retry limit、strategy change、source reload、validation gate | [`../governance/ai-runtime-governance/recovery-retry-governance.md`](../governance/ai-runtime-governance/recovery-retry-governance.md) |
| Domain-specific reload set / forbidden behaviors | [`../metadata/recovery/`](../metadata/recovery/) |
| Machine-readable recovery strategy / phase reconciliation / state repair | `runtime.db` recovery tables |
| 修改 runtime recovery 定義 | 更新 `runtime.db` canonical config document，然後執行 `ai-skill runtime compile` refresh projections |

`runtime/runtime.db` 是 committed runtime config 的 canonical copy。不要再保留 committed runtime YAML mirror；若需要人類可讀 diff，先匯出臨時 JSON/YAML 檢視，不把 mirror commit 回 repo。

## Inbound References

- [`route.runtime.activation-rules`](../knowledge/runtime/routing-registry.yaml:77)
- [`route.runtime.context-ttl`](../knowledge/runtime/routing-registry.yaml:102)
- [`route.runtime.context-loading`](../knowledge/runtime/routing-registry.yaml:161)
- [`route.metadata.knowledge-atom-schema`](../knowledge/runtime/routing-registry.yaml:191)
- [`route.models.model-aware-routing`](../knowledge/runtime/routing-registry.yaml:319)
- [`route.runtime.router-flow`](../knowledge/runtime/routing-registry.yaml:348)
- [`route.runtime.context-ttl-doc`](../knowledge/runtime/routing-registry.yaml:407)
- `gate.checkpoint.capability_discovery_completed` in [`runtime.db`](runtime.db)
- `obligation.checkpoint.run_capability_discovery` in [`runtime.db`](runtime.db)

## Databases

Runtime uses two SQLite databases with different lifecycles:

### `runtime.db` (Immutable — Compiled Registry)

Generated by Go-native `ai-skill runtime compile` from canonical config documents stored inside `runtime.db` plus deterministic prose mappings also stored inside `runtime.db`. Rebuilt when canonical runtime documents or mapped prose sources change.

All committed runtime config is stored in SQLite. Do not keep a committed YAML mirror for data whose canonical copy is `runtime.db`.

> **⚠️ Commit 規則：`runtime.db` 必須包含在 commit 中。**
> 當 canonical runtime config 變更時，`runtime.db` 本身就是 source-of-truth 與 projection output 的 commit artifact。
> 若手動 commit（跳過 hook），**必須**確認 `runtime validate` 通過且 `runtime.db` 已 `git add`。

| Table | Canonical Source | Purpose |
|-------|--------|---------|
| `runtime_config_documents` | `runtime.db` | Full canonical JSON documents replacing former runtime YAML documents |
| `runtime_config_projections` | `runtime.db` | Projection metadata from canonical documents to runtime tables |
| `runtime_source_files` | `runtime.db` | Compatibility manifest for former source paths, marked `source_kind='db'` |
| `phases`, `phase_transitions` | `runtime_config_documents` | Execution phase definitions and transition rules |
| `obligations`, `obligation_ledger` | `runtime_config_documents` | Per-phase atomic duties with verification criteria |
| `gates`, `blocking_gates` | `runtime_config_documents` | Phase transition prerequisites with severity and failure actions |
| `transaction_states`, `transaction_transitions`, `transaction_rules`, `transaction_templates`, `transaction_templates_ext` | `runtime_config_documents` | Transaction state definitions, rules, and templates |
| `activation_rules`, `activation_rules_mirror`, `core_bootstrap_rules` | `runtime_config_documents` | Lazy-load rule definitions with activation conditions |
| `discovery_checkpoints`, `discovery_search_strategy`, `capability_checkpoints` | `runtime_config_documents` | Phase-aware capability discovery checkpoints |
| `generated_surfaces` | Compiled from prose sources | Extracted structured data from workflow, enforcement, governance documents |
| `compiler_metadata` | Auto-generated | Compiler version, compilation timestamp, schema version |
| Remaining runtime config tables | `runtime_config_documents` | Normalized projections for budget, TTL, guards, health, intelligence routing, output governance, pipeline, prompt artifacts, recovery, scheduler, distributed runtime, and capability checkpoints |

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
