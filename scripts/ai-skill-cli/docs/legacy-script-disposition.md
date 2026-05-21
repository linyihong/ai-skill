# Legacy Script Disposition

## Closure Decision

The Go CLI is now the primary desktop runtime entrypoint. Runtime Ruby report/index/query/validation entrypoints, the Ruby runtime compiler, and the Roo Python adapter have been deleted after native CLI coverage became the default. Some non-runtime shell surfaces remain because they still require write-mode parity.

## Disposition Table

| Legacy Entry | Disposition | Reason / Removal Condition |
| --- | --- | --- |
| `scripts/init-new-project.sh` | Retained, pending parity | Delete after `ai-skill init-project` write mode and template parity are complete |
| `scripts/agent-goals.sh` | Retained, pending parity | Delete after full goal lifecycle write parity is complete |
| `scripts/install-hooks.sh` / `.githooks/` | Deleted | Replaced by `ai-skill hooks install` dry-run planner using `scripts/git-hooks/`; write mode remains blocked until parity fixtures |
| `scripts/sync-cursor-bundle.sh` | Tool-specific adapter retained | Cursor mirror sync remains opt-in and tool-specific |
| `scripts/ai-skill-close-loop.sh` | Retained, pending parity | Delete or thin-wrap after commit/push/private-path parity |
| Runtime report/index/query Ruby scripts | Deleted | Replaced by `ai-skill runtime refresh` and `ai-skill runtime query`; `--legacy-wrapper` is removed for refresh |
| Runtime validators | Deleted | Replaced by `ai-skill runtime validate`; `--legacy-wrapper` is removed for validate |
| `runtime/compiler/compiler-engine.rb` | Deleted | Replaced by Go-native `ai-skill runtime compile` source-to-DB compiler using runtime YAML and `runtime/compiler/compiler-rules.yaml` |
| `scripts/migrate-runtime-config-to-sqlite.rb` | Deleted | Obsolete migration helper; supported path is `ai-skill runtime compile` / compiler integration |
| `scripts/init-runtime-state-db.rb` | Deleted | Mutable runtime-state scope is not active; future support must be Go-native |
| `scripts/sync-runtime-yaml-from-embedded.rb` | Deleted | Prevents accidental rollback from embedded data into stale YAML; restore only via a dedicated source restoration plan |
| `scripts/set-roo-global-custom-instructions.py` | Deleted | Replaced by guarded `ai-skill roo set-global-custom-instructions` with fake VS Code SQLite DB tests |
| `scripts/git-hooks/*` | Hook adapter retained | Git hook surface remains an adapter installed by CLI |

## Removal Policy

Future deletion is allowed only when all are true:

1. Replacement command has write-mode parity.
2. Tests cover dry-run, success, failure, and missing dependency behavior.
3. Docs point users to the new CLI command.
4. `script-parity-inventory.md` row is updated to `deleted`, `hook adapter retained`, `tool-specific adapter retained`, or `explicitly out of scope`.

## Current Closure Scope

This plan closes the cross-platform desktop runtime migration and deletes runtime Ruby/Python scripts that now have native CLI coverage. Remaining script deletion is limited to non-runtime shell write-mode parity and tool adapters that are intentionally not yet native.
