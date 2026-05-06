> Follow [shared rules](../../../../shared-rules/README.md), [feedback lessons](../../../../shared-rules/feedback-lessons.md), and [goal/action/validation](../../../../shared-rules/goal-action-validation.md). This lesson is generalized and avoids project-specific details.

### 2026-05-07 - Keep Project Incidents Out Of Skills

Status: promoted

#### One-line Summary

Reusable skills should capture generalized causes, decisions, and validation loops, while concrete project incidents stay in project documentation.

#### Human Explanation

It is useful to convert field bugs into better development guidance, but copying the triggering project's names, paths, payloads, endpoints, test classes, or live-data quirks into a reusable skill makes the skill narrow and brittle. A skill should teach what to do next time, not preserve the exact incident.

#### Trigger

A reusable guidance update was written with details from a specific SDK bug investigation instead of only the generalized workflow.

#### Evidence

- Tool: User review of the skill update.
- Sanitized excerpt: A skill update included project-specific examples where a generic SDK defect closure rule was enough.
- Evidence path: The reusable guidance was corrected by removing project names and replacing them with neutral categories.

#### Generalized Lesson

When promoting an incident into a skill, split the content:

- Put the **why** and **repeatable method** in the skill.
- Put the **who/where/sample/result** in the project repository, runbook, issue, or integration notes.

#### Agent Action

Before editing a reusable skill, scan the proposed text for project names, private paths, endpoint strings, payload fragments, class names, sample IDs, or live environment quirks. If any appear, move them to the project docs and keep the skill phrased as a neutral rule with validation criteria.

#### Goal / Action / Validation

- Goal: Keep reusable guidance portable across projects.
- Action: Generalize incident-derived lessons before writing them into skill docs.
- Validation or reference source: Search the skill for project-specific strings after the edit; project-specific evidence should remain in project docs.

#### Applies When

- Updating reusable skills from app/API/SDK bugs, live integration failures, reverse-engineering observations, or product implementation lessons.
- Promoting feedback lessons into workflow, checklist, documentation, or template files.

#### Does Not Apply When

- Writing project-specific BDD, integration notes, runbooks, issue reports, or architecture contracts inside the project repository.

#### Promotion Target

- `shared-rules/reusable-guidance-boundary.md`.
- `shared-rules/README.md`, `content-layering.md`, `feedback-lessons.md`, `linked-updates.md`, `sanitization.md`, `goal-action-validation.md`, `dependency-reading.md`.
- `DOCUMENTATION.md` § Reusable Guidance Boundary.
- `CHECKLIST.md` § Reusable Guidance Boundary.
- `WORKFLOW.md` Docs-first BDD closure loop / SDK defect closure loop.
- `README.md` What Belongs Here / What Does Not Belong Here / Classification Rules.
- `SKILL.md` Feedback Loop.

#### Required Linked Updates

- `DOCUMENTATION.md`: updated with the split between generalized skill guidance and project-specific incident details.
- `CHECKLIST.md`: updated so reviews can catch project-specific details before skill changes are considered complete.
- `WORKFLOW.md`: generalized BDD and SDK defect closure loops so they describe reusable flow and delegate concrete paths to project governance docs.
- `README.md`: updated so the skill index states what belongs and what does not belong in reusable guidance.
- `SKILL.md`: updated so the feedback loop requires a post-edit search for project-specific strings.
- `feedback_history/common/README.md`: indexed this lesson.
- `shared-rules/reusable-guidance-boundary.md`: promoted the boundary to a global rule because it applies to all skills and shared docs, not only app development guidance.
- `shared-rules/README.md`, `content-layering.md`, `feedback-lessons.md`, `linked-updates.md`, `sanitization.md`, `goal-action-validation.md`, `dependency-reading.md`: updated or linked so agents must read and apply the global boundary, analyze incomplete closure causes, and perform required linked updates.
