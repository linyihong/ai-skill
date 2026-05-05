# Cross-Skill References

Use cross-skill references when one skill needs another skill's output contract, documentation format, validation checklist, or implementation guidance.

## Principle

Skills may reference other skills, but they must not copy another skill's full rules. The referring skill should state when to read the other skill, what output it expects, and which boundary remains owned by the current skill.

## When To Reference Another Skill

Add an explicit cross-skill reference when:

- A workflow produces an artifact that another skill is responsible for consuming.
- A skill needs another skill's template, checklist, or contract format.
- A finding should be converted into another domain's guidance.
- A user asks for a handoff from one skill to another.

Example: `apk-analysis` owns authorized APK evidence, traffic attribution, schema recovery, and sanitized feature reconstruction handoff. `app-development-guidance` owns turning that handoff into BDD, Domain Model Contract, API / Interface Contract, Error Handling Contract, implementation slices, checks, and tests.

## Required Reference Shape

When adding a cross-skill reference, include:

| Field | Requirement |
| --- | --- |
| Target skill | Link to `../<skill-name>/` or the specific file/folder to read. |
| Trigger | State exactly when the agent should read the target skill. |
| Expected input/output | Name the artifact passed between skills, such as a handoff, checklist, contract, fixture, or implementation guidance. |
| Ownership boundary | State what remains in the current skill and what belongs in the target skill. |
| Sanitization boundary | State whether target-specific data, secrets, raw evidence, or product conclusions must stay in project docs. |
| Linked updates | Update both skills' entry points or explain why the target skill already covers the handoff. |

## Do Not

- Do not paste another skill's full workflow into the current skill.
- Do not make every skill read every other skill by default.
- Do not create circular "always read" chains.
- Do not move target-specific evidence into a reusable skill just because two skills reference each other.
- Do not describe a required cross-skill handoff as optional when it is needed for the requested output.

## Good Pattern

```markdown
Use `../target-skill/` when <specific trigger>. The current skill owns <current boundary>; the target skill owns <target boundary>. Pass <artifact name> with <required fields>. Keep <sensitive or target-specific data> in project docs.
```

## Linked Updates

When a cross-skill reference is added or changed:

- Update the referring skill's `SKILL.md` and any relevant `README.md`, `WORKFLOW.md`, `DOCUMENTATION.md`, or technique file.
- Update the target skill's entry point if it needs to recognize the incoming handoff.
- Update [`linked-updates.md`](linked-updates.md) if the relationship becomes a recurring repo-wide rule.
- Run the Cursor bundle sync when the changed files live under `shared-rules/` or `skills/`.

← [回到共用規則索引](README.md)
