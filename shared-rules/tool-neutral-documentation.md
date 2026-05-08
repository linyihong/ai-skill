# Tool-Neutral Documentation

Reusable rules, skills, templates, indexes, and feedback lessons should be tool-neutral by default. Tool-specific setup, paths, hooks, commands, UI labels, and synchronization details belong under `ai-tools/` or the tool's own configuration files.

## Core Rule

When writing reusable documentation:

1. Describe the portable behavior, decision rule, workflow, or artifact first.
2. Use generic terms such as `agent`, `AI tool`, `tool-specific adapter`, `local tool mirror`, `project tool config`, and `<PROJECT_ROOT>`.
3. Do not make a specific tool sound required unless the rule truly only applies to that tool.
4. Use a Strategy-style adapter when behavior differs by tool: keep the common contract in the skill/shared rule, and isolate only the tool-specific execution details.
5. Put tool-wide setup instructions in the matching `ai-tools/<tool>.md` file.
6. Put skill-specific tool execution differences in a small skill-local adapter document only when the difference belongs to that skill.
7. Link from reusable docs to `ai-tools/` or the skill-local adapter when users need concrete setup steps.

## Allowed Tool-Specific Locations

Tool names, paths, and UI instructions may appear in:

| Location | Allowed content |
| --- | --- |
| `ai-tools/<tool>.md` | Tool setup, sync paths, UI steps, hooks, settings, and troubleshooting. |
| `skills/<skill>/tool-adapters/<tool>.md` | Skill-specific execution differences for one tool, when the core skill remains tool-neutral. |
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

## Strategy-Style Tool Adapters

Use this pattern when a skill has real execution differences across AI tools:

```text
skills/<skill>/
  README.md                 # tool-neutral overview
  WORKFLOW.md               # tool-neutral workflow
  tool-adapters/
    README.md               # index of supported adapters
    <tool>.md               # only the execution differences for that tool
```

The core skill document acts like the strategy interface:

- Trigger conditions.
- Inputs and outputs.
- Required evidence and validation.
- Safety, sanitization, and handoff rules.
- Tool-neutral workflow and terminology.

Each tool adapter acts like one strategy implementation:

- Which tool events, commands, hooks, prompts, or settings are used.
- What the tool can automate versus what remains manual.
- Tool-specific failure modes and validation.
- Links back to the exact core workflow steps it implements.

Do not copy the full core workflow into each adapter. If a tool-specific adapter needs to restate common behavior, move that behavior back to the core skill/shared rule and link to it.

Choose placement by scope:

| Scope | Placement |
| --- | --- |
| Tool-wide setup, sync, global hooks, UI, settings | `ai-tools/<tool>.md` |
| Skill-specific execution details for one tool | `skills/<skill>/tool-adapters/<tool>.md` |
| Project-specific tool config | Project docs or project tool config |
| Reusable cross-tool policy | `shared-rules/` |

## Review Checklist

Before finishing a reusable documentation change, check:

- Does the root or skill-level wording depend on a single IDE or agent product?
- Are tool-specific paths such as `.cursor/` or `~/.cursor/` only in tool docs, tool config, or explicitly tool-specific scripts?
- Does a generic rule say "configured tool sync" first and link to `ai-tools/` for concrete commands?
- If a skill needs tool-specific behavior, is it isolated in `tool-adapters/<tool>.md` and linked back to the core workflow?
- Did any new skill or shared rule accidentally copy a tool-specific setup section instead of linking to it?

← [Back to shared rules index](README.md)
