# Compiler

## Files

- [`compiler-engine.rb`](compiler-engine.rb) — Prose-to-YAML compiler engine
- [`compiler-rules.yaml`](compiler-rules.yaml) — Compiler rules and mapping

## Inbound References

- [`obligation.checkpoint.run_compiler`](../obligations/obligation-ledger.yaml)
- [`gate.checkpoint.compiler_executed`](../gates/blocking-gates.yaml)
- [`phase.checkpoint`](../phases/phase-machine.yaml) — allowed_actions includes `run_compiler`

## Source-of-Truth

- [`governance/lifecycle/compiler-philosophy.md`](../../governance/lifecycle/compiler-philosophy.md) — Design philosophy and core principles
