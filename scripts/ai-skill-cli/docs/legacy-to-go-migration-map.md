# Legacy To Go Migration Map

This document is the developer-facing handoff map for legacy runtime surfaces that were moved into the Go CLI. It answers four questions before any legacy file is deleted:

1. Which old surface existed?
2. Which Go command or package owns the replacement?
3. Which source-of-truth file should future developers edit?
4. Which tests or validation prove the move is complete?

Use this map together with [`script-parity-inventory.md`](script-parity-inventory.md), [`command-contract.md`](command-contract.md), and [`legacy-script-disposition.md`](legacy-script-disposition.md). The parity inventory keeps the full old-script table; this file keeps the active developer route.

## Runtime Compiler

| Old surface | New owner | Source to edit | Validation |
| --- | --- | --- | --- |
| `runtime/compiler/compiler-engine.rb` | `ai-skill runtime compile` | `runtime/compiler/compiler-rules.yaml` and runtime YAML sources | `go test ./...`, `ai-skill runtime compile`, `ai-skill runtime validate` |
| `runtime/compiler/embedded_data.rb` | `scripts/ai-skill-cli/internal/app/runtime_compiler.go` | `runtime/phases/`, `runtime/obligations/`, `runtime/gates/`, `runtime/recovery/`, `runtime/scheduler/`, `runtime/transactions/`, `runtime/discovery/` | source-to-DB compiler tests, runtime DB row-count and JSON validation |
| Ruby compiler schema creation | `createGoRuntimeSchema` in `runtime_compiler.go` | Go compiler schema code | native runtime DB validation and stable compiler snapshot tests |
| Ruby generated surface extraction | Go prose extraction in `runtime_compiler.go` | `runtime/compiler/compiler-rules.yaml` mappings and prose sources | generated surface assertions and golden runtime fixture |
| Ruby compiler metadata writes | Go compiler metadata writes | `goRuntimeCompilerVersion` and runtime compiler code | `compiler_metadata` validation |

### Developer Rule

Do not recreate Ruby fallback files or copy an existing `runtime.db` as a compiler shortcut. A valid compiler change must compile from source files into a fresh DB and pass native runtime DB validation.

## Restored Runtime YAML Sources

| Runtime concern | Source path | Runtime DB tables |
| --- | --- | --- |
| Phase machine | `runtime/phases/phase-machine.yaml` | `phases`, `phase_transitions`, `phase_machine` |
| Obligations | `runtime/obligations/obligation-ledger.yaml` | `obligations`, `obligation_ledger` |
| Blocking gates | `runtime/gates/blocking-gates.yaml` | `gates`, `blocking_gates` |
| Transactions | `runtime/transactions/transaction-machine.yaml`, `runtime/transactions/transaction-templates.yaml` | `transaction_states`, `transaction_transitions`, `transaction_rules`, `transaction_templates`, `transaction_templates_ext` |
| Recovery | `runtime/recovery/*.yaml` | `recovery_strategies`, `state_repair`, `obligation_rebuild`, `phase_reconciliation` |
| Scheduler | `runtime/scheduler/*.yaml` | `execution_queue`, `priority_scheduler` |
| Discovery checkpoints | `runtime/discovery/capability-checkpoints.yaml` | `discovery_checkpoints`, `discovery_search_strategy`, `capability_checkpoints` |

When any source above changes, run:

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 runtime compile
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 runtime validate
```

Use the binary matching the host platform.

## Roo Global Instructions

| Old surface | New owner | Source to edit | Validation |
| --- | --- | --- | --- |
| `scripts/set-roo-global-custom-instructions.py` | `ai-skill roo set-global-custom-instructions` | `scripts/ai-skill-cli/internal/app/roo.go` | fake VS Code `state.vscdb` tests in `roo_test.go` |

The Roo command is a guarded tool adapter, not a general onboarding default. It refuses by default when VS Code appears to be running, checks that the SQLite DB exists, updates `RooVeterinaryInc.roo-cline`, checkpoints WAL, and verifies the write.

## Hooks And Release

| Surface | New behavior |
| --- | --- |
| `scripts/git-hooks/pre-commit` | Calls repo-local `ai-skill runtime compile` and `ai-skill runtime validate`; no Ruby compiler call remains. |
| `scripts/ai-skill-cli/bin/*` | Rebuilt only after CLI source changes; `go test ./...` verifies `BUILDINFO`, `SHA256SUMS`, and host binary smoke. |

## Deleted Shell Surfaces

| Legacy surface | Go owner | Current Go coverage | Follow-up |
| --- | --- | --- | --- |
| `scripts/sync-cursor-bundle.sh` | `ai-skill sync-cursor-bundle` / `internal/app/sync_cursor_bundle.go` | Dry-run planning, explicit target requirement, target-outside-repo check, copy-fallback strategy, skill mirror planning | Future Cursor mirror writes must be implemented as Go write mode with managed mirror writes, fake home fixture, unmanaged target safety, copy fallback, and symlink policy |

## Retained Shell Surfaces

These shell files are legacy-retained surfaces, not places for new functionality. New behavior must be added to the Go CLI first.

| Legacy surface | Go owner | Current Go coverage | Deletion condition |
| --- | --- | --- | --- |
| `scripts/ai-skill-close-loop.sh` | `ai-skill close-loop` / `internal/app/close_loop.go` | Dry-run inspection, missing Git block, unsafe Git state block, active lock block, dirty owner grouping | Delete or reduce to a short binary bootstrap wrapper after Go supports commit, push, private-path scan, plan closure, optional Cursor sync, and readback parity |

## Required Close Loop

Before deleting or replacing another legacy surface:

1. Add or update the row in this file.
2. Update `script-parity-inventory.md` with old behavior, new command, status, and required validation.
3. Update `legacy-script-disposition.md` with the deletion or retention reason.
4. Add fixture-backed tests when the surface writes files, touches Git, mutates tool settings, or updates `runtime.db`.
5. Run the relevant Go tests and runtime compile / refresh / validate commands.
6. Rebuild repo-local binaries if Go CLI source changed.

For new automation, create or extend a Go CLI command instead of adding a shell script. If a shell wrapper is unavoidable, it must only locate and invoke the repo-local binary and must include a deletion condition in this map.
