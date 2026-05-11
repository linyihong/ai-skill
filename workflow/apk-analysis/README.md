# APK Analysis Workflow

`workflow/apk-analysis/` is the candidate home for tool-neutral APK analysis execution flow. During the pilot, `skills/apk-analysis/SKILL.md` remains the active tool skill entrypoint.

## Scope

This layer owns:

- Agent execution sequence for authorized APK analysis.
- Task decomposition, capture windows, documentation closure, and handoff flow.
- References from workflow steps to `analysis/apk/` methods and `intelligence/engineering/apk-analysis/` lessons.
- Cross-skill handoff rules when APK findings become app/API/SDK development guidance.

This layer does not own:

- Deep runtime or traffic analysis technique content; use `analysis/apk/`.
- Stable engineering wisdom and anti-patterns; use `intelligence/engineering/apk-analysis/`.
- Tool-specific UI, local mirror, hook installation, or sync details; use `ai-tools/` or project docs.
- Project-specific findings, raw evidence, or private service details.

## Current Source References

| Workflow concern | Current source | Pilot target status |
| --- | --- | --- |
| Skill trigger and authorization boundary | `../../skills/apk-analysis/SKILL.md` | remains active |
| Default decision tree | `../../skills/apk-analysis/WORKFLOW.md` | candidate for workflow extraction |
| Environment and tool preparation | `../../skills/apk-analysis/TOOLS.md` | reference from workflow; do not duplicate |
| Documentation and artifact gates | `../../skills/apk-analysis/DOCUMENTATION.md` | candidate for workflow extraction |
| First-day project runbook | `../../skills/apk-analysis/RUNBOOK.md` | candidate for onboarding workflow |
| Development guidance handoff | `../../skills/app-development-guidance/` | referenced only when trigger applies |

## Reference-First Workflow Shape

1. Confirm authorization, scope, APK version, device/emulator, and allowed actions.
2. Route through `analysis/apk/` to identify the traffic/runtime path.
3. Load only the matching technique category from `skills/apk-analysis/techniques/`.
4. Produce sanitized project artifacts: UI map, operation-to-API matrix, API catalog, runtime baseline, fixtures, and open questions as applicable.
5. If findings must become app/API/SDK implementation guidance, hand off to `skills/app-development-guidance/` with a sanitized Feature Reconstruction Handoff.
6. If a reusable lesson emerges, keep the lesson in `skills/apk-analysis/feedback_history/` until promotion rules move it into intelligence or feedback layers.

## Compatibility Notes

- Existing tools should still load `skills/apk-analysis/SKILL.md`.
- This directory is safe to read as a high-level workflow map, but it is not yet a replacement for the skill files.
- Future extraction should preserve old links, update `knowledge/indexes/README.md`, and attach `metadata/schema.md` metadata to each atom.
