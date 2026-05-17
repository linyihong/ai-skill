# Router

## Files

- [`activation-engine.rb`](activation-engine.rb) — Activation rule evaluation engine
- [`activation-rules.yaml`](activation-rules.yaml) — Activation rule definitions
- [`activation-table.md`](activation-table.md) — Situation-to-activated-rules mapping table

## Inbound References

- [`gate.checkpoint.routing_completed`](../gates/blocking-gates.yaml)
- [`obligation.checkpoint.route_task`](../obligations/obligation-ledger.yaml)

## Source-of-Truth

- [`governance/lifecycle/router-philosophy.md`](../../governance/lifecycle/router-philosophy.md) — Design philosophy and core principles
