# Contract Governance Review Checklist

Use this checklist when a project has multiple planning, BDD, contract, implementation, generated, and test documents, especially when the product was implemented first and documented later.

## Document Precedence

- The project defines which document wins when docs disagree.
- Governance/framework contracts own repository-wide invariants and required linked updates.
- Product plans own intent, scope, non-goals, canceled requirements, and accepted business language.
- BDD owns observable behavior.
- Domain/API/interface/error/hardware/command contracts own integration shapes and invariants.
- Implementation and tests are updated when higher-level docs change, or the exception is recorded.

## Traceability

- Product, feature, rule, operation, command, or diagnostic IDs are stable.
- IDs link to BDD scenarios or behavior specs.
- BDD scenarios link to implementation refs and tests.
- API/command/diagnostic contracts link to fixtures or generated clients.
- Unimplemented, noop, process-only, manual-only, and out-of-scope items are explicitly labeled.

## BDD Closure

- Critical BDD scenarios are marked `automated`, `fixture-backed`, `manual-evidence`, `pending-runner`, or `not-automatable`.
- Narrative Gherkin without a runner is not counted as automated coverage.
- Missing runner or step-definition work has an owner or alternative executable test type.
- Manual-only scenarios have review or release-gate evidence.

## Same-session sync (avoid implementation–documentation drift)

- A coding task is not closed while tests pass but contracts or BDD still describe superseded behavior—unless docs were verified unchanged or deferral is explicitly recorded.
- Bug fixes and refactors are explicitly classified; refactors with accidental behavior change are reclassified and trigger contract/BDD updates in the **same change batch**.
- Project-specific **Linked updates** matrices (often in `docs/development/README.md` or architecture contracts) are consulted when touching listed modules.

## Contract Sync

- API/schema changes regenerate typed clients, SDKs, mocks, fixtures, or schema packages.
- Vendor integration changes update sanitized excerpts, fixtures, gated live-test notes, and secret/redaction rules.
- Permission/auth/tenant/session changes update BDD, API contracts, policy docs, tests, and release checks.
- Database/persistence changes update domain, migration, fixture, repository, and rollback notes.
- Tooling or extension rule changes update rule catalog, diagnostic/command codes, fixtures, and adapter tests.
