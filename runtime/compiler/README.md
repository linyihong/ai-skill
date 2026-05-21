# Compiler

## Files

- [`compiler-engine.rb`](compiler-engine.rb) — Prose-to-SQLite compiler engine
- [`embedded_data.rb`](embedded_data.rb) — All runtime YAML data embedded as Ruby hash constants
- [`compiler-rules.yaml`](compiler-rules.yaml) — Compiler rules and mapping (embedded in `embedded_data.rb`)

## Inbound References

- `obligation.checkpoint.run_compiler` in [`runtime.db`](../runtime.db) / [`embedded_data.rb`](embedded_data.rb)
- `gate.checkpoint.compiler_executed` in [`runtime.db`](../runtime.db) / [`embedded_data.rb`](embedded_data.rb)
- `phase.checkpoint` in [`runtime.db`](../runtime.db) / [`embedded_data.rb`](embedded_data.rb) — allowed_actions includes `run_compiler`

## Source-of-Truth

- [`governance/lifecycle/compiler-philosophy.md`](../../governance/lifecycle/compiler-philosophy.md) — Design philosophy and core principles
