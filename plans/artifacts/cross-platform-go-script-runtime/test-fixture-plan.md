# Test Fixture Plan：Cross-Platform Go Script Runtime

> **上游計畫**：[`2026-05-21-0834-cross-platform-go-script-runtime.md`](../../active/2026-05-21-0834-cross-platform-go-script-runtime.md)

## Fixture Strategy

Fixtures must avoid the user's real home directory, real git config, real Cursor bundle, and real working tree. All destructive or write-capable behavior starts with dry-run fixtures.

## Required Fixtures

| Fixture | Purpose | Required Assertions |
| --- | --- | --- |
| `fixture/temp-repo-clean` | clean repo close-loop behavior | no-op dry-run, status clean |
| `fixture/temp-repo-dirty-owned` | owner group classification | only expected files are grouped |
| `fixture/temp-repo-merge-state` | unsafe repo state | blocks commit / push |
| `fixture/missing-git-path` | Git absent from PATH | `doctor` / `close-loop` block with install guidance |
| `fixture/fake-home` | home / Cursor bundle isolation | no real user home writes |
| `fixture/windows-paths` | path separator and drive handling | normalized paths match contract |
| `fixture/runtime-source-change` | runtime.db assertion | source keyword appears in generated surface |

## Missing Git Fixture

Purpose：prove Git is external dependency and missing Git cannot create partial linked-update / close-loop.

Setup:

- Execute command with PATH isolation.
- Ensure `git` is not discoverable.
- Use a temporary repo-like directory with files but no actual Git execution.

Assertions:

- `ai-skill doctor --require-git` returns `missing_dependency`.
- `ai-skill close-loop --commit` returns `missing_dependency`.
- Output includes install guidance.
- No file, index, commit, push, hook, mirror, or runtime DB is modified.

## Runtime DB Assertion Fixture

Purpose：prove source changes enter `runtime/runtime.db`.

Setup:

- Use fixture source content with a unique keyword.
- Run compile in dry-run and real fixture mode.
- Query generated surface by source path and keyword.

Assertions:

- Compile succeeds.
- `generated_surfaces` contains expected source path.
- `data` contains expected keyword.
- Missing keyword returns `validation_failed`.

## Windows Path Fixture

Purpose：prove path normalization does not assume POSIX shell.

Setup:

- Use Windows-style paths, drive letters, spaces, and mixed separators.
- Avoid shell-specific quoting assumptions.

Assertions:

- Path normalization is deterministic.
- JSON output uses stable path representation.
- Commands do not require Git Bash / WSL.

## Fake Home Fixture

Purpose：prevent accidental writes to real user config.

Setup:

- Override home/config roots with fixture paths.
- Run `init-project`, `goals`, `sync-cursor-bundle`, and `hooks install` in dry-run.

Assertions:

- No writes occur outside fixture root.
- Planned writes are listed in JSON output.
- Permission denial returns `permission_denied`.

## Artifact Gate

Before Phase 1 Go implementation:

- [ ] Fixtures have stable names and expected assertions.
- [ ] Missing Git fixture exists in the planned test suite.
- [ ] Runtime DB assertion fixture is linked to command contract.
- [ ] Windows path fixture does not require POSIX shell.
- [ ] Fake home fixture prevents real local configuration writes.
