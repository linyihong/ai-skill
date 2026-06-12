# Reference Reasoning

**Scope**: Reference-Driven Development — generic reasoning for turning implicit reference material into explicit specifications before implementation.

This directory is **not** the Visual Reference Workflow. UI-specific templates and per-screen specs live in project docs (e.g. Vidoe-Test `docs/frontend-contracts/`).

## Phase 0.5 placement decision

| Distills | Location |
|----------|----------|
| Reference-Driven Development (cross-domain) | Ai-skill `reference-reasoning/` |
| Visual Reference Workflow (UI domain) | Project `visual-reference-spec.template.md` + `specs/*-spec.md` |

Future umbrella candidate: `execution/model-before-action/` — not created yet; see [`reference-decomposition.md`](reference-decomposition.md) § Related Reasoning Families.

## Entries

| File | Purpose |
| --- | --- |
| [`reference-decomposition.md`](reference-decomposition.md) | Core pattern: Reference → Decompose → Explicit Spec → Implement |
| [`reference-types.md`](reference-types.md) | Illustrative type hints only (not a canonical taxonomy) |
