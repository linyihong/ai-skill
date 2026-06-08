# Sanitization Mechanical Enforcement

This file is the companion for the mechanical executor behind [`sanitization.md`](sanitization.md). The prose rule stays canonical for what must not enter reusable shared-layer documents; runtime surfaces define how the rule is checked.

## Runtime Sources

| Surface | Role |
|---|---|
| [`../runtime/repository-topology.yaml`](../runtime/repository-topology.yaml) | Defines shared-layer vs project-local paths. |
| [`../runtime/sanitization-patterns.yaml`](../runtime/sanitization-patterns.yaml) | Defines deterministic generic regex patterns and placeholder allowlist. |
| Project `.ai-skill-project.yaml` | Declares private project tokens/entities that compile into `runtime.db.derived_forbidden_tokens`. |

## Enforcement Boundary

- Scanner runs at pre-commit on staged content.
- Shared-layer paths are classified by `runtime.db.repository_topology`; the scanner does not hardcode folder globs.
- Project-specific token blocking only uses explicit project metadata. Absence from reusable docs is never treated as proof that a token is private.
- Generic regex patterns cover email, phone, OS absolute path, and credential-like tokens.
- Phase 2.5 incident-score heuristics are intentionally deferred and warning-only.

## Placeholders

Allowed placeholder forms are defined in `runtime/sanitization-patterns.yaml`; current examples include `<USER>`, `<PROJECT_ROOT>`, `<AI_SKILL_REPO>`, and `<WORKSPACE>`.
