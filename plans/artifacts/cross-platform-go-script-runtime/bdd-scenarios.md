# BDD-Lite Scenarios：Cross-Platform Go Script Runtime

> **上游計畫**：[`2026-05-21-0834-cross-platform-go-script-runtime.md`](../../active/2026-05-21-0834-cross-platform-go-script-runtime.md)

## Requirement Link

- **Source**：Cross-platform Go script runtime plan
- **Actor / system role**：Ai-skill maintainer、agent、CI runner、desktop contributor、mobile control-plane user
- **Behavior boundary**：CLI command behavior、dependency detection、runtime validation、safe close-loop
- **Ambiguity disposition**：draft; implementation cannot start until scenarios are reviewed

## Scenario: Missing Git Blocks Close-Loop

**Given** `ai-skill close-loop --commit` is executed on a desktop platform  
**When** `git` is not available in PATH  
**Then** the command exits with `missing_dependency`  
**And** no files are staged, committed, pushed, or modified  
**And** the output tells the user to install Git.

## Scenario: Doctor Reports Git Requirement

**Given** `ai-skill doctor --require-git --json` is executed  
**When** Git is missing  
**Then** JSON output contains `error.code = "missing_git"`  
**And** the exit code is stable across Windows, macOS, and Linux.

## Scenario: Runtime Compile Asserts Generated Surface

**Given** a runtime source file was modified  
**When** `ai-skill runtime compile --assert-source <path> --assert-keyword <keyword>` completes  
**Then** `runtime/runtime.db` contains the source path and keyword in the expected generated surface  
**And** validation fails if the keyword is absent.

## Scenario: Dry-Run Prevents Side Effects

**Given** a command can write files or run git operations  
**When** it is invoked with `--dry-run`  
**Then** it reports planned actions  
**And** no tracked file, untracked file, git index, commit, hook, mirror, or runtime DB is modified.

## Scenario: Unsafe Repo State Blocks Commit

**Given** the repository is in merge, rebase, or cherry-pick state  
**When** `ai-skill close-loop --commit` is executed  
**Then** the command exits with `unsafe_repo_state`  
**And** no commit is created.

## Scenario: iOS Native Binary Is Unsupported

**Given** the user asks to run `ai-skill` as a downloaded native binary on iOS  
**When** platform support is evaluated  
**Then** the result is unsupported  
**And** the recommended options are App-contained runtime, Browser/WASM, or SSH remote runner.

## Acceptance Criteria

- Missing Git cannot produce a partial close-loop.
- Dry-run commands do not mutate file system, git index, commits, hooks, mirrors, or runtime DB.
- Runtime compile can prove source-to-DB propagation with content assertions.
- Mobile support decisions do not promise iOS native arbitrary binary.

## Validation Target

- **Proof type**：fixture-backed automated tests
- **Test / fixture / checklist**：[`test-fixture-plan.md`](test-fixture-plan.md)
- **Limitations**：這些 scenarios 尚未證明 performance、release signing、mobile app feasibility 或 full compiler parity。

## Regression Scope

- [ ] Existing shell script behavior must remain available until replacement is validated.
- [ ] New tests required for missing Git, unsafe repo, dry-run, runtime assertion, and iOS unsupported decision.
- [ ] Test data / fixtures needed：temporary repo、PATH isolation、fake home、runtime source fixture。
