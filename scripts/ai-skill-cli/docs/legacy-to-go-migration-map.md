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
| `runtime/compiler/compiler-engine.rb` | `ai-skill runtime compile` | `runtime/runtime.db` and SQLite canonical runtime documents | `go test ./...`, `ai-skill runtime compile`, `ai-skill runtime validate` |
| `runtime/compiler/embedded_data.rb` | `scripts/ai-skill-cli/internal/app/runtime_compiler.go` | `runtime/runtime.db` canonical documents | source-to-DB compiler tests, runtime DB row-count and JSON validation |
| Ruby compiler schema creation | `createGoRuntimeSchema` in `runtime_compiler.go` | Go compiler schema code | native runtime DB validation and stable compiler snapshot tests |
| Ruby generated surface extraction | Go prose extraction in `runtime_compiler.go` | `runtime/runtime.db` mappings and prose sources | generated surface assertions and golden runtime fixture |
| Ruby compiler metadata writes | Go compiler metadata writes | `goRuntimeCompilerVersion` and runtime compiler code | `compiler_metadata` validation |

### Developer Rule

Do not recreate Ruby fallback files or committed runtime YAML mirrors. A valid compiler change must refresh projections from SQLite canonical runtime documents into a fresh DB and pass native runtime DB validation.

## SQLite Canonical Runtime Sources

| Runtime concern | Canonical source | Runtime DB tables |
| --- | --- | --- |
| Phase machine | `runtime/runtime.db` | `phases`, `phase_transitions`, `phase_machine` |
| Obligations | `runtime/runtime.db` | `obligations`, `obligation_ledger` |
| Blocking gates | `runtime/runtime.db` | `gates`, `blocking_gates` |
| Transactions | `runtime/runtime.db` | `transaction_states`, `transaction_transitions`, `transaction_rules`, `transaction_templates`, `transaction_templates_ext` |
| Recovery | `runtime/runtime.db` | `recovery_strategies`, `state_repair`, `obligation_rebuild`, `phase_reconciliation` |
| Scheduler | `runtime/runtime.db` | `execution_queue`, `priority_scheduler` |
| Discovery checkpoints | `runtime/runtime.db` | `discovery_checkpoints`, `discovery_search_strategy`, `capability_checkpoints` |

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

## Shell To Go Migration Targets

These shell files are current migration targets. They must be deleted after Go write-mode parity, fixtures, documentation updates, binary rebuild, and readback are complete.

| Legacy surface | Go owner | Current Go coverage | Deletion condition |
| --- | --- | --- | --- |
| deleted `scripts/init-new-project.sh` | `ai-skill init-project` / `internal/app/init_project.go` | Dry-run and write mode, tool selection, conflict / force handling, bootstrap templates | Deleted after Go write mode and selected-tools fixtures passed |
| deleted `scripts/agent-goals.sh` | `ai-skill goals` / `internal/app/goals.go` | `init`, `status`, `start`, `update`, `split`, `pause`, `complete --validated`, `cleanup`, locks, index, git exclude | Deleted after lifecycle and lock fixtures passed |
| deleted `scripts/ai-skill-close-loop.sh` | `ai-skill close-loop` / `internal/app/close_loop.go` | Dry-run, commit, push, Git safety, owner grouping, private-path scan, plan closure, runtime validation, readback | Deleted after grouped commit and plan closure fixtures passed |

## Git Hook Adapter Boundary

Git hook files may remain only as integration adapters when Git requires a hook file path. Business logic must live in Go commands:

| Hook surface | Go owner | Required behavior |
| --- | --- | --- |
| `scripts/git-hooks/pre-commit` | `ai-skill hooks run pre-commit` | Hook adapter locates repo-local binary and delegates staged runtime compile / knowledge validation to Go. |
| `scripts/git-hooks/post-commit` | `ai-skill hooks run post-commit` | Hook adapter delegates post-commit behavior to Go; reference-only remains no-op unless a Go write mode is explicitly enabled. |

If a native binary hook can be installed safely on Windows, macOS, and Linux, `ai-skill hooks install` should install the binary entrypoint. Otherwise the retained hook file must stay a minimal adapter with no reusable automation logic.

## Required Close Loop

Before deleting or replacing another legacy surface:

1. Add or update the row in this file.
2. Update `script-parity-inventory.md` with old behavior, new command, status, and required validation.
3. Update `legacy-script-disposition.md` with the deletion or retention reason.
4. Add fixture-backed tests when the surface writes files, touches Git, mutates tool settings, or updates `runtime.db`.
5. Run the relevant Go tests and runtime compile / refresh / validate commands.
6. Rebuild repo-local binaries if Go CLI source changed.

For new automation, create or extend a Go CLI command instead of adding a shell script. If a hook adapter or bootstrap wrapper is unavoidable, it must only locate and invoke the repo-local binary and must include a deletion condition in this map. `ai-skill runtime validate` enforces this boundary for `.sh` bootstrap/helper files through `gate.tool_bootstrap_shell_requires_cli_decision`.
