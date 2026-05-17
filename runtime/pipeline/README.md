# Pipeline

## Files

- [`context-flow.yaml`](context-flow.yaml) — Context flow pipeline definitions
- [`guard-chain.yaml`](guard-chain.yaml) — Guard chain pipeline definitions
- [`relevance-engine.yaml`](relevance-engine.yaml) — Relevance engine pipeline definitions
- [`session-lifecycle.yaml`](session-lifecycle.yaml) — Session lifecycle pipeline definitions

## Inbound References

- [`gate.checkpoint.pipeline_executed`](../gates/blocking-gates.yaml)
- [`obligation.checkpoint.execute_pipeline`](../obligations/obligation-ledger.yaml)

## Source-of-Truth

- [`governance/lifecycle/pipeline-philosophy.md`](../../governance/lifecycle/pipeline-philosophy.md) — Design philosophy and core principles
