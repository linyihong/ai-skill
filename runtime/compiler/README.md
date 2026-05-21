# Compiler

## Files

- [`compiler-rules.yaml`](compiler-rules.yaml) — Compiler rules and source-to-target mapping for Go-native `ai-skill runtime compile`
- Go implementation: [`../../scripts/ai-skill-cli/internal/app/runtime_compiler.go`](../../scripts/ai-skill-cli/internal/app/runtime_compiler.go)

## Inbound References

- `obligation.checkpoint.run_compiler` in [`runtime.db`](../runtime.db) / [`../obligations/obligation-ledger.yaml`](../obligations/obligation-ledger.yaml)
- `gate.checkpoint.compiler_executed` in [`runtime.db`](../runtime.db) / [`../gates/blocking-gates.yaml`](../gates/blocking-gates.yaml)
- `phase.checkpoint` in [`runtime.db`](../runtime.db) / [`../phases/phase-machine.yaml`](../phases/phase-machine.yaml) — allowed_actions includes `run_compiler`

## Source-of-Truth

- [`governance/lifecycle/compiler-philosophy.md`](../../governance/lifecycle/compiler-philosophy.md) — Design philosophy and core principles
