# Legacy Script Disposition

## Closure Decision

The Go CLI is now the primary desktop runtime entrypoint. Runtime Ruby report/index/query/validation entrypoints, the Ruby runtime compiler, the Roo Python adapter, and the legacy Cursor bundle shell sync have been deleted. Remaining shell surfaces are legacy-retained only because they still require write-mode parity; they must not receive new functionality.

Developer migration details live in [`legacy-to-go-migration-map.md`](legacy-to-go-migration-map.md). Before deleting another legacy surface, update that map first so future maintainers can see the new Go owner, editable source, and validation evidence.

## Disposition Table

| Legacy Entry | Disposition | Reason / Removal Condition |
| --- | --- | --- |
| `scripts/init-new-project.sh` | Retained, pending parity | Delete after `ai-skill init-project` write mode and template parity are complete |
| `scripts/agent-goals.sh` | Retained, pending parity | Delete after full goal lifecycle write parity is complete |
| `scripts/install-hooks.sh` / `.githooks/` | Deleted | Replaced by `ai-skill hooks install` dry-run planner using `scripts/git-hooks/`; write mode remains blocked until parity fixtures |
| `scripts/sync-cursor-bundle.sh` | Deleted | Removed as legacy shell surface; future Cursor mirror writes must be implemented in `ai-skill sync-cursor-bundle` Go write mode |
| `scripts/ai-skill-close-loop.sh` | Legacy retained, pending Go write-mode parity | Delete or thin-wrap after `ai-skill close-loop` supports commit/push/private-path scan/plan closure/readback parity |
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

Future additions are stricter: new repository automation must be implemented in Go CLI first. A new shell entrypoint is allowed only as a Git hook adapter or temporary binary bootstrap wrapper with a documented deletion condition.

## Current Closure Scope

This plan closes the cross-platform desktop runtime migration and deletes runtime Ruby/Python scripts that now have native CLI coverage. `scripts/sync-cursor-bundle.sh` is also deleted by user decision; its Go command remains dry-run only until write-mode parity is implemented.
