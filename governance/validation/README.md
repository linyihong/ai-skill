# Knowledge Validation Gates

`governance/validation/` defines validation gates for new AI-native layer changes. It complements executable rules in `shared-rules/`; it does not replace them.

## Required Gates

| Gate | Required check | Applies when |
| --- | --- | --- |
| Source boundary | Confirm canonical repository paths are edited, not tool mirrors or runtime copies. | Every Ai-skill writeback. |
| Old entrypoint | Confirm old `skills/`, `shared-rules/`, `ai-tools/`, or `scripts/` entrypoints remain reachable. | Candidate maps, atom promotions, migrations. |
| Linked updates | Check affected README, roadmap, index, metadata, and source entry files. | Any new layer path or routing change. |
| Metadata | Confirm `metadata/schema.md` fields exist when a Knowledge Atom is introduced or promoted. | Candidate atom, validated atom, promoted atom. |
| Navigation | Confirm `knowledge/indexes/README.md` can route to the new path when it should be discoverable. | Routing surfaces and promoted reference paths. |
| Link check | Resolve Markdown links in touched docs. | Documentation changes. |
| Lints | Run available lints for touched files. | Documentation or code changes. |
| Diff review | Review for secrets, private hostnames, raw evidence, local absolute paths, and unrelated changes. | Before commit. |
| Close-loop dry run | Run `./scripts/ai-skill-close-loop.sh` to verify dirty path grouping. | Before commit. |
| Commit / push / readback | Commit, push, reread changed entries, and confirm clean `git status --short --branch`. | Ai-skill repository updates. |

## Migration Validation Checklist

Use this checklist for any move from old `skills/` content into a new layer:

```text
Goal:
- What user-visible or runtime outcome does this change support?

Source:
- Old source path:
- New candidate/promoted path:
- Source-of-truth state:

Linked updates:
- Layer README:
- Knowledge index:
- Metadata:
- Roadmap:
- Old entrypoint:

Validation:
- Lints:
- Markdown links:
- Diff review:
- Close-loop dry run:
- Commit/push/readback:
- Clean status:
```

## Pass / Block Rules

- If old entrypoints break, the change is blocked.
- If metadata is missing for a promoted atom, the change remains `candidate-map` or `candidate-atom`.
- If links fail, fix links before commit.
- If validation cannot be run, record the blocker and do not mark the lifecycle state as promoted.
- If a change is reference-only and no tool mirror was used, tool sync is not applicable.

## Relationship To Shared Rules

- Dependency reading, canonical writeback, commit / push / readback, and clean status remain governed by `shared-rules/dependency-reading.md`.
- Linked update requirements remain governed by `shared-rules/linked-updates.md`.
- Rule priority remains governed by `shared-rules/rule-weight.md`.
- This file provides architecture-layer validation shape for the new knowledge system.
