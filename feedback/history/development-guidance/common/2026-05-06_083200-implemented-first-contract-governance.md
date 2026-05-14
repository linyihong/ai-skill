# Implemented-first projects need contract governance and BDD closure
# Extracted — See [`workflow/software-delivery/execution-flow.md`](../../../../workflow/software-delivery/execution-flow.md)

Status: promoted

## Lesson

Projects that are implemented first and documented later can still become reliable contract-first projects, but the backfill must recover more than BDD. It also needs document precedence, stable traceability IDs, minimum doc-sync rules, generated-client flow, third-party integration boundaries, and an explicit path from narrative BDD to executable or evidence-backed validation.

## Rule

When analyzing an implemented-first project:

1. Define which document wins when docs disagree: governance/framework contract, product plan, BDD, contracts, implementation, tests.
2. Link product/rule/operation/command/diagnostic IDs to BDD, code refs, fixtures, and tests.
3. Mark every critical BDD scenario as `automated`, `fixture-backed`, `manual-evidence`, `pending-runner`, or `not-automatable`.
4. Add a minimum doc-sync matrix for API, permission, database, UI, generated client, vendor integration, CLI/tooling, diagnostic, and release changes when those surfaces exist.
5. Require OpenAPI/schema/source-contract changes to regenerate typed clients, SDKs, mocks, or fixtures.
6. Keep vendor source docs separate from sanitized integration excerpts, fixtures, live-test gates, and secret-safe product docs.
7. For tools, IDE extensions, linters, and CLIs, separate pure kernel logic from adapters and keep rule catalogs aligned with diagnostics, fixtures, and tests.

## Required Linked Updates

- `process/README.md`: added contract governance, traceability, BDD execution closure, and implemented-first pipeline backfill rules.
- `WORKFLOW.md`: added evidence translation, owner layers, validation, and filing rules for generated clients, vendor integrations, and tooling.
- `SKILL.md` and `README.md`: added triggers and classification language.
- `templates/initial-development-docs.md`: added contract governance, traceability, generated client, and vendor integration fields.
- `CHECKLIST.md` and `checklists/contract-governance-review.md`: added repeatable review gates.
- `implementation/backend/contract-codegen.md`, `implementation/backend/vendor-integration.md`, and `implementation/tooling/README.md`: added buildable implementation patterns.

## Validation

Use implemented-first projects only as source evidence for reusable patterns. Do not copy project-specific hosts, credentials, business rules, vendor payloads, customer data, or internal policy text into reusable guidance.
