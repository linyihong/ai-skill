# Engineering Governance Intelligence

Engineering governance intelligence defines reusable decision systems: contracts, validation shapes, failure authority, severity, and review checklists. It is not an execution workflow by itself.

## Modules

| Module | Use When |
| --- | --- |
| [`responsive-ui/`](responsive-ui/README.md) | UI behavior must remain valid across render contexts, viewport variance, fixed/sticky layout, overflow risk, safe areas, or responsive framework decisions. |

## Boundary

- Governance knowledge answers what must be classified and who can decide failure.
- `workflow/` answers when to run the checks and how to close the task.
- `validation/` or workflow validation slices answer how evidence is acquired and evaluated.
- Project-specific evidence stays in the project repository.
