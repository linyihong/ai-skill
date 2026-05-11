# AI-native Knowledge Operating System

This document defines the repository-level architecture direction for Ai-skill. It is a roadmap and boundary document, not an executable shared rule. Operational policies still live in `shared-rules/`; tool-specific setup still lives in `ai-tools/`.

## Purpose

The AI-native Knowledge Operating System is the layer that lets agents reliably load knowledge, track goals, improve reusable guidance, validate changes, and close git writeback loops without treating any one AI tool or local mirror as the source of truth.

Its default model is **reference-first**:

1. Agents read the canonical `<AI_SKILL_REPO>` directly.
2. `shared-rules/` provides the operating rules.
3. `skills/` provides capability modules.
4. `ai-tools/` and tool config provide adapters for specific tools.
5. `.agent-goals/` records temporary project-local execution state.
6. Scripts support validation, goal tracking, commit/push/readback, and optional tool sync.

Symlink, bundle, and copy snapshot flows remain compatibility layers for tools that need native scanning or offline snapshots. They are not the default architecture and should not replace direct reads from the canonical repository.

## Layers

| Layer | Location | Responsibility |
| --- | --- | --- |
| Source of truth | `<AI_SKILL_REPO>` git repository | Canonical rules, skills, templates, scripts, and architecture docs. |
| Operating rules | `shared-rules/` | Cross-skill policy, dependency reading, linked updates, goal ledger rules, failure learning, validation, and documentation boundaries. |
| Capability modules | `skills/` | Domain-specific workflows, checklists, documentation templates, techniques, and feedback lessons. |
| Tool adapters | `ai-tools/`, tool config, optional skill adapters | Tool-specific paths, hooks, UI behavior, sync strategy, and troubleshooting. |
| Temporary execution state | `<PROJECT_ROOT>/.agent-goals/` | Active user goals, owner/lock decisions, open work, next actions, and completion validation. |
| Close-loop automation | `scripts/` | Conservative helpers for goal ledger operations, grouped commits, push, readback support, and optional tool sync. |

## Reference-First Loading

Reference-first means an agent is expected to load the central repo by path instead of relying on a copied skill package inside each project.

Minimum startup shape:

```text
<AI_SKILL_REPO>/shared-rules/README.md
<AI_SKILL_REPO>/skills/<skill-name>/SKILL.md
```

After the shared-rule bootstrap, the agent follows dependency routing to read only the skill workflows, tool adapters, feedback lessons, templates, or architecture docs needed for the current task.

This keeps updates simple:

- Update the canonical repo.
- Validate linked updates.
- Commit and push.
- Re-read updated entries after push.
- Let future sessions read the canonical repo instead of refreshing copied packages.

## Compatibility Layers

Use compatibility layers only when the active tool cannot reliably reference the central repo or needs native scan behavior.

| Strategy | When to use | Close-loop rule |
| --- | --- | --- |
| Reference-first | Default for normal agent work. | No tool mirror sync required. Confirm the canonical repo is current and readable. |
| Symlink / bundle bridge | A tool needs native skill discovery, but local paths can point back to `<AI_SKILL_REPO>`. | Sync only when this bridge is intentionally in use or explicitly requested. |
| Copy snapshot | The tool cannot read the central repo, cannot use symlinks, or needs an offline snapshot. | Record source commit/date and refresh intentionally; do not treat the copy as source. |

Tool-specific details belong in `ai-tools/`. Generic rules should say "configured tool sync" or "optional tool sync" instead of naming a single tool as the default.

## Migration Roadmap

### Phase 1: Reference-First Default

- Keep `shared-rules/README.md` as the bootstrap index.
- Keep tool docs clear that reference-first does not require copying or bundle sync.
- Keep close-loop automation conservative: commit/push/readback is mandatory for Ai-skill changes; tool sync is conditional.
- Keep compatibility scripts available for users who still need symlink, bundle, or copy snapshot workflows.

### Phase 2: Compatibility Inventory

- Identify any active workflows that still depend on native tool scanning or copied skill directories.
- Document each remaining compatibility case in the relevant `ai-tools/` file.
- Prefer symlink or reference strategies over copy snapshots when the tool allows it.
- Update stale docs when they make copy or bundle sync sound like the default path.

### Phase 3: Deprecation Readiness

Copy and bundle flows can be deprecated when all are true:

- Active tool docs point to reference-first as the normal path.
- No active project workflow requires copied skill packages for day-to-day use.
- Remaining native-scan needs are covered by symlink/reference adapters or explicitly documented exceptions.
- Close-loop validation no longer depends on mirror refresh checks except for compatibility tests requested by the user.
- A deprecation note and migration path exist before scripts or docs are removed.

## Removal Criteria For Copy And Bundle Sync

Do not remove compatibility scripts solely because reference-first exists. Remove or archive them only after:

1. A search confirms no required workflow instructs agents to run them by default.
2. Tool-specific docs list any replacement strategy.
3. Users who depend on native scanning have a symlink/reference path or an acknowledged manual snapshot path.
4. `shared-rules/dependency-reading.md` and `shared-rules/linked-updates.md` still explain how to handle necessary tool sync without making it universal.
5. The close-loop process has been validated without accidental source/mirror drift.

## Relationship To Existing Documents

- `shared-rules/` remains the executable policy layer.
- `skills/` remains the capability layer.
- `ai-tools/` remains the tool adapter layer.
- `scripts/` remains helper automation, not the architecture itself.
- `.agent-goals/` remains temporary project state and is deleted when goals complete.

This architecture document should be updated when the repo changes how agents load rules, how skills are discovered, how goal state is tracked, or how source/mirror boundaries are enforced.
