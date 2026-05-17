# Compiler

## Files

- [`compiler-engine.rb`](compiler-engine.rb) — Prose-to-SQLite compiler engine
- [`embedded_data.rb`](embedded_data.rb) — All runtime YAML data embedded as Ruby hash constants
- [`compiler-rules.yaml`](compiler-rules.yaml) — Compiler rules and mapping (embedded in `embedded_data.rb`)

## Inbound References

- [`obligation.checkpoint.run_compiler`](../obligations/obligation-ledger.yaml)
- [`gate.checkpoint.compiler_executed`](../gates/blocking-gates.yaml)
- [`phase.checkpoint`](../phases/phase-machine.yaml) — allowed_actions includes `run_compiler`

## Source-of-Truth

- [`governance/lifecycle/compiler-philosophy.md`](../../governance/lifecycle/compiler-philosophy.md) — Design philosophy and core principles
