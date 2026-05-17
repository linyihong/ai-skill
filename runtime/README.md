# Runtime

Executable runtime layer. Machine-oriented, query-oriented, deterministic.

## Domains

| Domain | Path | Description |
|--------|------|-------------|
| Activation | [`router/activation-rules.yaml`](router/activation-rules.yaml) | Lazy-load rules with activation conditions |
| Routing | [`routing/README.md`](routing/README.md) | Task intent → knowledge index → metadata → source-of-truth gate |
| Discovery | [`discovery/README.md`](discovery/README.md) | Phase-aware capability discovery checkpoints |
| Phases | [`phases/phase-machine.yaml`](phases/phase-machine.yaml) | Execution phase state machine |
| Obligations | [`obligations/obligation-ledger.yaml`](obligations/obligation-ledger.yaml) | Per-phase atomic duties |
| Gates | [`gates/blocking-gates.yaml`](gates/blocking-gates.yaml) | Phase transition prerequisites |
| Compiler | [`compiler/compiler-engine.rb`](compiler/compiler-engine.rb) | Prose → YAML compilation |
| Generated | [`generated/`](generated/) | Compiled runtime surfaces |
| Transactions | [`transactions/transaction-machine.yaml`](transactions/transaction-machine.yaml) | Writeback transaction state machine |
| Pipeline | [`pipeline/`](pipeline/) | Context flow, guard chain, relevance engine |
| Recovery | [`recovery/`](recovery/) | Phase reconciliation, state repair, obligation rebuild |
| Scheduler | [`scheduler/`](scheduler/) | Execution queue, priority scheduler |
| Guards | [`guards/`](guards/) | Circuit breaker, context pollution |
| Onboarding | [`onboarding/`](onboarding/) | New project/task setup guidance |
| Output Governance | [`output-governance/`](output-governance/) | Language policy, output rules, governance gates |
| Prompt Artifacts | [`prompt-artifacts/`](prompt-artifacts/) | Artifact templates, composition rules |
| Context | [`context/`](context/) | TTL policy |
| Budget | [`budget/`](budget/) | Token budget |
| Distributed | [`distributed/`](distributed/) | Multi-agent coordination, distributed locks |
| Intelligence | [`intelligence/`](intelligence/) | Intelligence routing |

## Inbound References

- [`route.runtime.activation-rules`](../knowledge/runtime/routing-registry.yaml:77)
- [`route.runtime.context-ttl`](../knowledge/runtime/routing-registry.yaml:102)
- [`route.runtime.context-loading`](../knowledge/runtime/routing-registry.yaml:161)
- [`route.metadata.knowledge-atom-schema`](../knowledge/runtime/routing-registry.yaml:191)
- [`route.models.model-aware-routing`](../knowledge/runtime/routing-registry.yaml:319)
- [`route.runtime.router-flow`](../knowledge/runtime/routing-registry.yaml:348)
- [`route.runtime.context-ttl-doc`](../knowledge/runtime/routing-registry.yaml:407)
- [`gate.checkpoint.capability_discovery_completed`](gates/blocking-gates.yaml)
- [`obligation.checkpoint.run_capability_discovery`](obligations/obligation-ledger.yaml)

## Source-of-Truth

Runtime does not hold conceptual explanations. Source-of-truth for runtime design:

- `governance/` — Design philosophy, lifecycle, validation
- `workflow/` — Human-readable execution flows
- `intelligence/` — Heuristics, analytical reasoning
- `enforcement/` — Executable policy rules
