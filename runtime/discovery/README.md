# Discovery

Phase-aware capability discovery checkpoints.

## Files

- [`capability-checkpoints.yaml`](capability-checkpoints.yaml) — Discovery checkpoint definitions per phase

## Inbound References

- [`gate.checkpoint.capability_discovery_completed`](../gates/blocking-gates.yaml)
- [`obligation.checkpoint.run_capability_discovery`](../obligations/obligation-ledger.yaml)
- [`phase.checkpoint`](../phases/phase-machine.yaml) — allowed_actions includes `run_capability_discovery`

## Source-of-Truth

- [`governance/lifecycle/capability-discovery-philosophy.md`](../../governance/lifecycle/capability-discovery-philosophy.md) — Conceptual explanation of the Capability Discovery Problem
