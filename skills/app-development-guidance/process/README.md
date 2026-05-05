# Development Process

Use this process when a product idea,企劃書, or feature brief needs to become implementation work that AI agents and engineers can execute without losing domain intent.

The sequence is contract-first: clarify behavior and domain language before teams start coding, then let each implementation surface work against the same contracts.

## Default Flow

| Step | Output | Notes |
| --- | --- | --- |
| 1. 企劃書 / product brief | Goals, users, scope, non-goals, constraints | Keep assumptions explicit; mark unknowns instead of inventing certainty. |
| 2. AI analysis + module split | Bounded Context map, module ownership, integration points | Split by domain responsibility, not by UI pages or database tables only. |
| 3. BDD behavior | Feature files or scenario tables | Describe user/system behavior in domain language. |
| 4. Domain Model Contract | Entities, value objects, commands, events, invariants | This is the core contract; define what must always be true. |
| 5. Architecture Contract | Boundaries, dependencies, data ownership, runtime/deployment shape | Define which layers may depend on each other and where decisions live. |
| 6. API Contract | OpenAPI/GraphQL/schema/events/RPC/message contracts | Critical when multiple surfaces or services integrate; write request/response/error shapes before implementation. |
| 7. Error Handling Contract | Error taxonomy, retry rules, validation errors, user-visible messages, logging | Design failure behavior before implementation; errors are part of the contract. |
| 8. Parallel implementation planning | Work slices per context and contract | Frontend/backend can start together only after shared contracts are stable enough. |
| 9. Backend / service implementation | Behavior + domain + API contract implementation | If no backend exists, replace with local service, library, worker, or platform implementation. |
| 10. Frontend / client implementation | Mock API, schema-first client, UI behavior | If no frontend exists, replace with CLI, SDK, mobile screen, job runner, or consumer integration. |
| 11. Automated tests | Unit, BDD, API contract, schema tests | Tests should prove domain invariants and contract compatibility. |
| 12. Integration test | End-to-end or component integration evidence | Verify real adapters, auth/session, error paths, and cross-context flows. |

## Required Contracts

Do not assume every project has a frontend/backend split. Pick contracts that match the architecture:

| App shape | Contract to prioritize |
| --- | --- |
| Frontend + backend | API Contract, Domain Model Contract, BDD scenarios. |
| Backend/API only | Domain Model Contract, API Contract, contract tests, integration tests with consumers. |
| Frontend-only app | UI behavior contract, local state/domain contract, mocked API/schema contract. |
| Mobile app | Screen/flow behavior, local storage/session contract, API Contract when remote services exist. |
| CLI / desktop / tool | Command contract, input/output schema, domain model, fixture-based tests. |
| Library / SDK | Public API contract, type/schema contract, examples, compatibility tests. |
| Event-driven / worker | Event schema, command/event contract, idempotency and retry behavior. |

## Initial Documentation Pack

When this skill is opened for a new feature or project, the agent should help produce a first draft or ask the missing questions for these documents:

| Document | Purpose | If missing, ask about |
| --- | --- | --- |
| Product Brief | Goal, users, scope, non-goals, assumptions | Who is this for, what problem is solved, what is explicitly out of scope? |
| Bounded Context Map | Module/domain split and ownership | What concepts change together, and what boundaries should not leak? |
| BDD Behavior | User/system behavior scenarios | What are the critical happy paths and failure paths? |
| Domain Model Contract | Core business objects and invariants | What must always be true, and what state transitions are allowed? |
| Architecture Contract | Layers, dependencies, ownership, runtime constraints | Which layer owns data, side effects, security, persistence, and external calls? |
| API / Interface Contract | Requests, responses, events, commands, public methods | Who consumes this contract, how is compatibility tested, and how are versions handled? |
| Error Handling Contract | Error types, recovery, user messaging, logging | Which errors are retryable, user-fixable, fatal, or security-sensitive? |
| Test Plan | Unit, BDD, contract, integration tests | What proves the behavior, invariants, and integration contract? |

These documents can start as lightweight Markdown drafts. If the project is small, keep them in one planning file; if they grow, split them into a folder with `README.md` and focused child files.

Use [`../templates/initial-development-docs.md`](../templates/initial-development-docs.md) for the first draft. Use [`../templates/README.md`](../templates/README.md) to choose between initial planning, reusable guidance notes, and quick threat-model reviews.

## Existing Project Documentation Backfill

When this skill is opened for a project that is already fully or mostly implemented, first audit the existing documents and backfill any missing development documents. Do not skip the process because implementation already exists.

| Missing document | Backfill rule |
| --- | --- |
| Product Brief | Reconstruct only what evidence supports: visible goals, users/actors, scope, non-goals, assumptions, and constraints. If original intent is unavailable, mark fields as `unknown` or `open question`; do not invent business rationale. |
| Bounded Context Map | Infer modules from code ownership, runtime boundaries, database tables, API groups, UI areas, queues, SDK/public APIs, and deployment units. |
| BDD Behavior | **Must be completed.** Reconstruct critical happy paths, failure paths, permissions, empty states, edge cases, and cross-context flows from the implemented product, tests, UI, API behavior, and logs. |
| Domain Model Contract | Infer entities, value objects, commands, events, invariants, and state transitions from code, schemas, storage, UI states, and tests; mark uncertain vocabulary as candidate. |
| Architecture Contract | Document actual dependency direction, data ownership, side-effect boundaries, integrations, runtime/deployment shape, and known violations. |
| API / Interface Contract | Extract actual request/response schemas, public methods, events, commands, auth/session behavior, versioning, compatibility, fixtures, and consumers. |
| Error Handling Contract | Backfill observed error taxonomy, retry rules, user messages, logging/redaction behavior, security-sensitive failures, and gaps. |
| Test Plan | Map existing tests to behavior/contracts and list required tests for uncovered BDD scenarios, invariants, contracts, and integration paths. |

Backfill order for existing projects:

1. Inventory existing docs, source folders, tests, schemas, API specs, fixtures, release notes, and observed behavior.
2. Create a documentation gap table with status: `exists`, `partial`, `missing`, or `unknown`.
3. Backfill BDD Behavior first when product brief is missing, because implemented behavior is the strongest available source of truth.
4. Backfill Domain Model, Architecture, API / Interface, and Error Handling Contracts from the completed behavior and implementation evidence.
5. Mark unknown product intent separately from observed behavior. Unknown intent does not block BDD completion.
6. Add tests or test TODOs for any critical BDD scenario that lacks coverage.

## Contract-First Rules

- BDD describes behavior; it should not lock in framework or database choices.
- Domain Model Contract owns invariants, business vocabulary, and state transitions.
- Architecture Contract owns dependency direction, runtime boundaries, data ownership, and allowed integration paths.
- API Contract owns integration shape: request, response, error, auth/session, versioning, and compatibility.
- Error Handling Contract owns failure taxonomy, retry policy, user messaging, logging, and security redaction.
- Implementation can run in parallel only when the shared contracts are versioned enough for mock, stub, or schema-first work.
- If a contract changes, update BDD, implementation, mocks, and tests in the same change or explicitly record why not.
- For already implemented projects, BDD becomes the required behavioral recovery document. Product Brief may contain unknowns, but BDD must be filled from observable product behavior and implementation evidence.

## When Frontend And Backend Do Not Both Exist

Replace "frontend" and "backend" with producer/consumer roles:

| Original role | Generic role |
| --- | --- |
| Frontend | Consumer: UI, CLI, SDK, job, mobile screen, test harness. |
| Backend | Provider: API, domain service, library function, local adapter, worker. |

The flow still applies:

1. Define behavior.
2. Define domain invariants.
3. Define the provider/consumer contract.
4. Build each side against mock, fixture, or schema.
5. Prove compatibility with contract and integration tests.

## Minimum Definition Of Ready

Before implementation starts, the feature should have:

- Product brief with scope and non-goals.
- Bounded Context or module split.
- BDD scenarios for critical behaviors.
- Domain Model Contract for core invariants.
- Architecture Contract for dependencies, ownership, and runtime boundaries.
- API, event, command, or public interface contract for integrations.
- Error Handling Contract for expected failures and recovery behavior.
- Test plan covering unit, behavior, contract, and integration levels.

For an already implemented project, "ready" means the missing-document audit is complete and BDD covers the implemented critical behavior, even if original product intent remains partly unknown.

## Minimum Definition Of Done

Before shipping or merging:

- Domain invariants are tested.
- Contract tests pass for provider and consumer.
- Mocks/fixtures match the latest contract.
- Integration test covers at least the critical happy path and one important failure path.
- Residual unknowns or deferred behavior are documented in the project repository.

← [Back to App Development Guidance](../README.md)
