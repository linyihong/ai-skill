> 遵守 [共用規則索引](../../../../shared-rules/README.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/app-development-guidance/development-process.md`](../../../../workflow/app-development-guidance/development-process.md)

### 2026-05-05 - Contract-first development flow

Status: promoted

#### One-line Summary

Start development from product intent, split bounded contexts, write BDD, define Domain, Architecture, API/Interface, and Error Handling Contracts, then let implementation and tests proceed against shared contracts.

#### Human Explanation

Jumping directly from a 企劃書 to frontend/backend tasks often loses domain language and creates mismatched assumptions. A contract-first flow gives AI agents and engineers a stable sequence: understand product intent, split bounded contexts, describe behavior, define the domain model, define integration contracts, then implement and test each side.

The flow should not assume every project has both frontend and backend. The same pattern applies to CLI tools, libraries, workers, mobile-only apps, backend-only APIs, and frontend-only apps by replacing frontend/backend with provider/consumer roles.

#### Trigger

- A product brief or feature idea needs to become implementation work.
- Multiple modules, teams, agents, or app surfaces need to work in parallel.
- API or schema mismatch is likely if coding starts too early.
- The app shape may not have a conventional frontend/backend split.

#### Evidence

- Tool: product brief, BDD scenarios, bounded context map, domain model notes, OpenAPI/GraphQL/schema/event contracts, test results.
- Sanitized excerpt: `brief -> bounded-contexts -> bdd -> domain-contract -> api-contract -> provider/consumer implementation -> contract tests -> integration tests`.
- Evidence path: project repository planning docs and tests; reusable skill stores only the generalized workflow.

#### Generalized Lesson

Use a contract-first development flow:

1. Product brief.
2. AI analysis and Bounded Context split.
3. BDD behavior.
4. Domain Model Contract.
5. Architecture Contract.
6. API, event, command, or public interface Contract.
7. Error Handling Contract.
8. Parallel implementation planning.
9. Provider implementation.
10. Consumer implementation with mocks, fixtures, or schema-first clients.
11. Unit, BDD, and contract tests.
12. Integration tests.

When no frontend/backend split exists, define provider/consumer roles that match the architecture.

#### Agent Action

- Do not start implementation from a product brief alone.
- Ask for or draft BDD scenarios before detailed implementation.
- Define the Domain Model Contract before the API Contract.
- Define Architecture Contract and Error Handling Contract before implementation details harden.
- Treat API Contract as a key integration artifact, not generated after coding by accident.
- If the architecture lacks frontend/backend, map the flow to provider/consumer roles.
- Add unit, BDD, contract, and integration tests to the plan.

#### Applies When

- Starting a new feature, app, module, API, SDK, CLI, worker, or app surface.
- AI agents need to split work into implementation tasks.
- Multiple implementation surfaces depend on shared behavior or schemas.

#### Does Not Apply When

- The change is a small typo, content edit, or isolated refactor with no behavior/API impact.
- The product has no new behavior and no contract boundary changes.
- Emergency hotfixes where the contract already exists and must be patched narrowly.

#### Validation

- Bounded contexts and ownership are documented.
- BDD scenarios cover critical behavior.
- Domain Model Contract names core invariants.
- Architecture Contract names dependency direction, ownership, and runtime boundaries.
- API/public interface contract exists before parallel implementation.
- Error Handling Contract names failure taxonomy, retry behavior, user messaging, logging, and redaction.
- Provider and consumer tests pass against the same contract.
- Integration test covers at least one critical happy path and one failure path.

#### Promotion Target

- `process/README.md`
- `implementation/README.md`
- `templates/README.md`
- `templates/initial-development-docs.md`
- `SKILL.md`
- `README.md`
- `CHECKLIST.md`
- `shared-rules/linked-updates.md`

#### Required Linked Updates

- Added `process/README.md`.
- Updated `implementation/README.md` with contract-to-implementation slice guidance.
- Added `templates/README.md` as the template routing index.
- Added `templates/initial-development-docs.md`.
- Updated `SKILL.md` quick start and default workflow links.
- Updated `README.md` file index and classification rules.
- Updated `CHECKLIST.md` with product-to-contract readiness checks.
- Updated `DOCUMENTATION.md` with the initial development docs template link.
- Updated `shared-rules/linked-updates.md` with `process/` and `templates/` linked updates.
- Updated `feedback_history/README.md` and `feedback_history/common/README.md`.
