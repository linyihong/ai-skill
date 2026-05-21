# Legacy Script Disposition

## Closure Decision

The Go CLI is now the primary desktop runtime entrypoint. Legacy scripts are not all deleted in this closure because several non-runtime surfaces still require write-mode parity. Runtime Ruby report/index/query/validation entrypoints have been deleted after native CLI coverage became the default. The Ruby compiler engine remains only because the Go compiler is still snapshot mode, not a true source-to-DB compiler.

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
| `runtime/compiler/compiler-engine.rb` | Retained as compiler source/parity reference | Delete only after true source-to-DB Go compiler parity, not just snapshot mode |
| `scripts/migrate-runtime-config-to-sqlite.rb` | Deleted | Obsolete migration helper; supported path is `ai-skill runtime compile` / compiler integration |
| `scripts/init-runtime-state-db.rb` | Deleted | Mutable runtime-state scope is not active; future support must be Go-native |
| `scripts/sync-runtime-yaml-from-embedded.rb` | Deleted | Prevents accidental rollback from embedded data into stale YAML; restore only via a dedicated source restoration plan |
| `scripts/set-roo-global-custom-instructions.py` | Tool-specific adapter retained | User editor globalStorage mutation is not a general CLI default |
| `scripts/git-hooks/*` | Hook adapter retained | Git hook surface remains an adapter installed by CLI |

## Removal Policy

Future deletion is allowed only when all are true:

1. Replacement command has write-mode parity.
2. Tests cover dry-run, success, failure, and missing dependency behavior.
3. Docs point users to the new CLI command.
4. `script-parity-inventory.md` row is updated to `deleted`, `hook adapter retained`, `tool-specific adapter retained`, or `explicitly out of scope`.

## Current Closure Scope

This plan closes the cross-platform desktop runtime migration and deletes the runtime Ruby scripts that now have native CLI coverage. Remaining script deletion is limited to non-runtime write-mode parity, tool adapters, and the Ruby compiler engine follow-up.
