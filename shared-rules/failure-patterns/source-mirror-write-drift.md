# Source / Mirror Write Drift

Status: validated
Class: `source-mirror-drift`

## Trigger

Use this pattern when a user asks an agent to update rules, skills, tool setup, feedback lessons, templates, or Ai-skill guidance, especially when the visible file is under a project `.cursor`, `~/.cursor`, generated bundle, runtime copy, or other tool mirror.

## Failure Mode

The agent edits the copy that the current tool can see instead of editing the canonical `<AI_SKILL_REPO>` source repository. The result looks fixed in the current session, but the reusable knowledge does not reliably persist, sync, commit, or propagate to other projects and tools.

## Risk

- The next agent reads stale canonical rules and repeats the same mistake.
- A project-local `.cursor` copy diverges from the shared knowledge base.
- Tool sync, commit, push, and readback gates are skipped because the agent thinks the mirror is the source.
- Private project paths or tool-local details can leak into reusable docs.

## Required Agent Action

1. Stop broad editing as soon as source/mirror confusion is suspected.
2. Locate `<AI_SKILL_REPO>` and confirm it is the git root.
3. Check `git status --short --branch` in `<AI_SKILL_REPO>`.
4. Apply the reusable change in the canonical source file.
5. Treat project `.cursor`, `~/.cursor/skills*`, `~/.cursor/shared-rules`, `~/.cursor/bundles/*`, and generated bundles as deployment or runtime surfaces.
6. Run the configured tool sync only after source changes are correct.
7. Commit, push, reread updated entries, and confirm clean status before claiming completion.

## Prevention Gate

Before the first write to rules, skills, templates, feedback lessons, or tool deployment paths, the agent must be able to answer:

| Check | Required answer |
| --- | --- |
| Canonical repo | Which path is `<AI_SKILL_REPO>` and is it the git root? |
| Current file role | Is this file source, project config, tool mirror, runtime copy, or generated output? |
| Source edit | Which canonical source file will be edited first? |
| Sync strategy | Is this reference-only, symlink/bundle, or copy snapshot? |
| Close loop | What diff review, sync, commit, push, readback, and clean-status checks will run? |

If any answer is unknown, do not edit the mirror. Read [`dependency-reading.md`](../dependency-reading.md), [`cursor-sync.md`](../cursor-sync.md), and the relevant `ai-tools/` document first.

## Validation

This pattern is validated when:

- The canonical source file contains the reusable change.
- Mirror or runtime paths are either unchanged, symlinked to source, or updated only by the configured sync.
- `git diff` and `git status --short --branch` were checked in `<AI_SKILL_REPO>`.
- The change was committed, pushed, and read back when it affects Ai-skill.
- The final response names any project-local mirror changes separately from the Ai-skill source update.

## Linked Rules

- [`failure-learning-system.md`](../failure-learning-system.md)
- [`dependency-reading.md`](../dependency-reading.md)
- [`cursor-sync.md`](../cursor-sync.md)
- [`tool-neutral-documentation.md`](../tool-neutral-documentation.md)
- [`linked-updates.md`](../linked-updates.md)

← [Back to failure patterns](README.md)
