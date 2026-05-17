# Distributed

## Files

- [`async-job-lifecycle.yaml`](async-job-lifecycle.yaml) — Async job lifecycle state machine
- [`distributed-locks.yaml`](distributed-locks.yaml) — Distributed lock definitions
- [`multi-agent-coordination.yaml`](multi-agent-coordination.yaml) — Multi-agent coordination protocol

## Inbound References

- [`gate.checkpoint.distributed_sync_completed`](../gates/blocking-gates.yaml)
- [`obligation.checkpoint.sync_distributed_state`](../obligations/obligation-ledger.yaml)

## Source-of-Truth

- [`governance/lifecycle/distributed-runtime-philosophy.md`](../../governance/lifecycle/distributed-runtime-philosophy.md) — Design philosophy and core principles
