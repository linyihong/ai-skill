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
| Default decision tree | `../../skills/apk-analysis/WORKFLOW.md` | ✅ extracted to [`execution-flow.md`](execution-flow.md) |
| Capture window detailed rules | `../../skills/apk-analysis/WORKFLOW.md` | ✅ extracted to [`execution-flow.md`](execution-flow.md) |
| Environment and tool preparation | `../../skills/apk-analysis/TOOLS.md` | reference from `analysis/apk/`; do not duplicate |
| Documentation and artifact gates | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ extracted to [`artifact-gates.md`](artifact-gates.md) |
| SDK live self-generation audit | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ extracted to [`artifact-gates.md`](artifact-gates.md) |
| Identity material self-generation audit | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ extracted to [`artifact-gates.md`](artifact-gates.md) |
| UI architecture map template | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ extracted to [`artifact-gates.md`](artifact-gates.md) |
| API catalog detail requirements | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ extracted to [`artifact-gates.md`](artifact-gates.md) |
| Sanitization rules | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ extracted to [`artifact-gates.md`](artifact-gates.md) |
| Developer guidance notes | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ extracted to [`artifact-gates.md`](artifact-gates.md) |
| Feedback lesson writing tips | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ extracted to [`artifact-gates.md`](artifact-gates.md) |
| Backfill rules | `../../skills/apk-analysis/DOCUMENTATION.md` | ✅ extracted to [`artifact-gates.md`](artifact-gates.md) |
| First-day project runbook | `../../skills/apk-analysis/RUNBOOK.md` | candidate for onboarding workflow |
| Development guidance handoff | `../../skills/app-development-guidance/` | referenced only when trigger applies |

## Reference-First Workflow Shape

1. Confirm authorization, scope, APK version, device/emulator, and allowed actions.
2. Route through `analysis/apk/` to identify the traffic/runtime path.
3. Load the matching workflow from `analysis/apk/workflows/` and intelligence atoms from `intelligence/engineering/apk-analysis/`.
4. Produce sanitized project artifacts: UI map, operation-to-API matrix, API catalog, runtime baseline, fixtures, and open questions as applicable.
5. If findings must become app/API/SDK implementation guidance, hand off to `skills/app-development-guidance/` with a sanitized Feature Reconstruction Handoff.
6. If a reusable lesson emerges, keep the lesson in `skills/apk-analysis/feedback_history/` until promotion rules move it into intelligence or feedback layers.

## Compatibility Notes

- Existing tools should still load `skills/apk-analysis/SKILL.md`.
- This directory is now a substantial reference-first workflow layer. Most content from `WORKFLOW.md` and `DOCUMENTATION.md` has been extracted here.
- The original skill files still contain the authoritative content; this directory provides a reference-first view.
- Future extraction should preserve old links, update `knowledge/indexes/README.md`, and attach `metadata/schema.md` metadata to each atom.
