# Tool-Neutral Documentation

Reusable rules, skills, templates, indexes, and feedback lessons should be tool-neutral by default. Tool-specific setup, paths, hooks, commands, UI labels, and synchronization details belong under `ai-tools/` or the tool's own configuration files.

## Core Rule

When writing reusable documentation:

1. Describe the portable behavior, decision rule, workflow, or artifact first.
2. Use generic terms such as `agent`, `AI tool`, `tool-specific adapter`, `local tool mirror`, `project tool config`, and `<PROJECT_ROOT>`.
3. Do not make a specific tool sound required unless the rule truly only applies to that tool.
4. Put tool-specific instructions in the matching `ai-tools/<tool>.md` file.
5. Link from reusable docs to `ai-tools/` when users need concrete setup steps.

## Allowed Tool-Specific Locations

Tool names, paths, and UI instructions may appear in:

| Location | Allowed content |
| --- | --- |
| `ai-tools/<tool>.md` | Tool setup, sync paths, UI steps, hooks, settings, and troubleshooting. |
| Tool config files such as `.cursor/rules/*.mdc` | Rules loaded by that tool. |
| Tool-specific scripts or script docs | Commands that exist only for that tool, with a link from generic docs when needed. |
| Project-local tool files | Project-specific adapter configuration, if it is safe to commit in that project. |

## Avoid In Reusable Docs

Avoid tool-specific wording in these places unless the section is explicitly about tool integration:

- Root `README.md`.
- `shared-rules/README.md` index summaries.
- Skill `README.md`, `SKILL.md`, `WORKFLOW.md`, `TOOLS.md`, `DOCUMENTATION.md`.
- Skill templates and `skills/ADDING_SKILLS.md`.
- Feedback lessons and reusable checklists.

Examples:

| Less portable | Prefer |
| --- | --- |
| `Cursor agent entry point` | `Agent entry point` |
| `Reload Cursor` | `Reload or refresh the active tool if it caches skills/rules` |
| `copy to .cursor/skills` | `deploy to the active tool's skill/rule location; see ai-tools/` |
| `run sync-cursor-bundle.sh` in a generic rule | `run the configured tool sync; Cursor details live in ai-tools/cursor.md` |

## Relationship To Existing Tool Files

This repository currently includes Cursor-specific documentation and helper scripts. They are still valid, but references to them should be centralized in `ai-tools/cursor.md`, `shared-rules/cursor-sync.md`, `.cursor/rules/`, and script docs when the script itself is tool-specific.

Reusable docs may mention a tool only as an example, not as the default requirement, unless the file is already inside that tool's documentation area.

## Review Checklist

Before finishing a reusable documentation change, check:

- Does the root or skill-level wording depend on a single IDE or agent product?
- Are tool-specific paths such as `.cursor/` or `~/.cursor/` only in tool docs, tool config, or explicitly tool-specific scripts?
- Does a generic rule say "configured tool sync" first and link to `ai-tools/` for concrete commands?
- Did any new skill or shared rule accidentally copy a tool-specific setup section instead of linking to it?

← [Back to shared rules index](README.md)
