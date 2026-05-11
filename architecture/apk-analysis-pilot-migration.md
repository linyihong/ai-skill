# APK Analysis Pilot Migration Map

本文件定義 `apk-analysis` 作為下一階段 Workflow / Analysis / Intelligence 分離的 pilot。目標是先建立 reference-first 目的地與 mapping，不搬移大量內容、不破壞既有 `skills/apk-analysis/` 入口。

## Migration Status

| Field | Value |
| --- | --- |
| Pilot skill | `skills/apk-analysis/` |
| Status | `candidate-map` |
| Old entrypoint | `skills/apk-analysis/SKILL.md` remains the active skill entrypoint |
| New reference paths | `analysis/apk/`, `workflow/apk-analysis/`, `intelligence/engineering/apk-analysis/` |
| Bulk migration | Not started |

## Compatibility Rules

- `skills/apk-analysis/SKILL.md` remains the canonical trigger for tools that load skills.
- New top-level layer paths are reference targets, not replacements, until content has been copied, validated, indexed, and linked back.
- Do not move or delete existing `skills/apk-analysis/` files during the pilot.
- Any future move must preserve old links with redirect notes or index rows.
- Target-specific APK hosts, endpoints, raw responses, tokens, device identifiers, and private run evidence stay in project docs.

## Source To Target Map

| Existing source | Current role | Candidate target | Migration action | Status |
| --- | --- | --- | --- | --- |
| `skills/apk-analysis/SKILL.md` | Tool skill trigger, authorization boundary, output style, cross-skill handoff | `workflow/apk-analysis/` | Keep as active entrypoint; later extract tool-neutral execution flow into workflow while retaining skill trigger | candidate |
| `skills/apk-analysis/WORKFLOW.md` | Evidence-first traffic/runtime decision tree | `analysis/apk/` + `workflow/apk-analysis/` | Split observation / triage methods into analysis; keep execution sequencing in workflow | candidate |
| `skills/apk-analysis/TOOLS.md` | Tool choice, command templates, failure interpretation | `analysis/apk/` + `workflow/apk-analysis/` | Move tool-selection reasoning to analysis; keep setup steps as workflow references or tool docs if tool-specific | candidate |
| `skills/apk-analysis/DOCUMENTATION.md` | Project artifact templates and documentation gates | `workflow/apk-analysis/` + `intelligence/engineering/apk-analysis/` | Keep artifact production flow in workflow; extract stable engineering lessons into intelligence | candidate |
| `skills/apk-analysis/techniques/` | Route-specific analysis methods | `analysis/apk/` | Keep categories as-is during pilot; future analysis layer can reference or gradually absorb category summaries | candidate |
| `skills/apk-analysis/feedback_history/` | Lesson history and validated/candidate learning | `intelligence/engineering/apk-analysis/` + `feedback/` + `memory/` | Keep lesson files in skill history; promote stable conclusions by reference, not bulk copy | candidate |
| Cross-skill handoff to `app-development-guidance` | Boundary between evidence recovery and development guidance | `workflow/apk-analysis/` | Preserve handoff artifact and ownership boundary | candidate |

## First Reference-First Paths

| Task intent | New path | Still read |
| --- | --- | --- |
| Decide how to observe an APK traffic/runtime path | `analysis/apk/README.md` | `skills/apk-analysis/WORKFLOW.md`, relevant `skills/apk-analysis/techniques/` category |
| Execute an APK analysis session or handoff | `workflow/apk-analysis/README.md` | `skills/apk-analysis/SKILL.md`, `WORKFLOW.md`, `DOCUMENTATION.md` |
| Reuse engineering lessons from APK analysis | `intelligence/engineering/apk-analysis/README.md` | `skills/apk-analysis/feedback_history/README.md`, validated lesson files |

## Completion Criteria For Future Migration

- Each moved atom has metadata using `metadata/schema.md`.
- `knowledge/indexes/README.md` routes task intents to the new path and the old skill entrypoint.
- Old `skills/apk-analysis/` links still resolve.
- Shared rules and skill dependencies remain readable from the old entrypoint.
- Validation includes Markdown link check, diff review, close-loop dry run, commit, push, readback, and clean status.

## Open Follow-Ups

- Decide whether technique category summaries should live under `analysis/apk/<category>/` or remain skill-local with metadata atoms.
- Create Knowledge Atom candidates for the highest-value `apk-analysis` workflow and technique entries.
- Define when a feedback lesson graduates from `skills/apk-analysis/feedback_history/` into `intelligence/engineering/apk-analysis/`.
