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
| Embedded / firmware / hardware product | Datasheet or protocol contract, hardware context contract, driver/service/application boundary, BDD, host fixtures, hardware-in-loop checks. |
| Static analysis / IDE extension / developer tool | Rule catalog, diagnostic or command contract, pure kernel/adapter boundary, fixture pairs, editor/CLI integration tests. |

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
| Hardware / Firmware Contract | Datasheet/protocol truth, electrical interface, pin/context injection, driver/service/application boundary, target constraints | What hardware facts are fixed, what is injected per board, and how are host/target tests run? |
| Test Plan | Unit, BDD, contract, integration tests | What proves the behavior, invariants, and integration contract? |

These documents can start as lightweight Markdown drafts. If the project is small, keep them in one planning file; if they grow, split them into a folder with `README.md` and focused child files.

Use [`../templates/initial-development-docs.md`](../templates/initial-development-docs.md) for the first draft. Use [`../templates/README.md`](../templates/README.md) to choose between initial planning, reusable guidance notes, and quick threat-model reviews.

## Change Intake Gate

Before any code change driven by this skill, inspect the project's 企劃書, product brief, planning docs, issue, ticket, PRD, design note, BDD, API contract, or equivalent project artifact. Classify the request before implementation:

| Change type | Required before code |
| --- | --- |
| New requirement / feature / behavior change | Update or create planning docs first: Product Brief or change brief, BDD scenarios, impacted Domain Model Contract, Architecture Contract, API / Interface Contract, Error Handling Contract, implementation slices, and tests. Do not start code until blocker questions are resolved. |
| Bug fix | Confirm expected behavior vs actual behavior, reproduction or evidence, affected BDD scenario or missing scenario, impacted contract/error handling, and regression test plan. If the fix changes intended behavior or public contract, treat it as a new requirement too. |
| Refactor / internal cleanup | Confirm no behavior or public contract change. If behavior, data ownership, API, error handling, security, storage, or tests change, reclassify as new requirement or bug. |
| Security / hardening change | Confirm threat or failure mode, owner layer, required control, validation method, and whether behavior/API/contracts/checklists must change. |

If there is no planning artifact, create a lightweight change brief before implementation. If the request is a new requirement, missing planning docs are blockers; ask the user and fill BDD/contracts before writing code.

## Contract Governance Gate

Every project with multiple docs must define which artifact wins when documents disagree. Use this default precedence unless the project has a stronger local rule:

1. Governance / framework contract: repository-wide invariants, required update rules, dependency direction, naming, build/run constraints.
2. Product plan / accepted brief: product intent, scope, non-goals, canceled requirements, and business language.
3. BDD behavior: observable user/system behavior and acceptance criteria.
4. Domain, architecture, API/interface, error handling, hardware, or command contracts.
5. Implementation and generated clients.
6. Tests, fixtures, and examples.

If lower layers reveal the higher layer is wrong, do not silently "fix" code only. Classify the conflict as one of:

| Conflict type | Required action |
| --- | --- |
| Product intent changed | Update product brief or plan, then BDD/contracts/tests. |
| BDD missing or stale | Backfill or revise BDD from evidence and link impacted tests. |
| Contract stale | Update the contract and all consumers, mocks, generated clients, fixtures, and tests in the same change. |
| Implementation bug | Keep docs stable, add or update regression tests, then fix code. |
| Test or fixture stale | Update tests/fixtures to the current contract and cite the source. |

Document canceled, deferred, out-of-scope, and not-tool-enforceable items explicitly. Do not leave them as invisible absences that future agents may reintroduce.

## Traceability Gate

When a project was implemented first and docs are being backfilled, require traceability in both directions:

| Link | Purpose |
| --- | --- |
| Product or rule ID -> BDD | Shows which behavior proves the requirement. |
| BDD -> code refs | Shows where the behavior is implemented. |
| BDD -> test refs | Shows how the behavior is verified or what gap remains. |
| Contract operation / command / diagnostic -> fixture | Shows provider/consumer compatibility and edge cases. |
| Generated client or SDK method -> API/OpenAPI/source contract | Prevents hand-copied endpoints and drift. |

Stable IDs can be feature IDs, rule IDs, operation IDs, route names, command names, diagnostic codes, event names, or scenario tags. If a behavior is intentionally documented but not implemented, mark it as `TBD`, `noop`, `not enforceable by tool`, `manual-only`, or `out of scope`, with the reason and owner.

## BDD Execution Closure

Narrative BDD is acceptable during backfill, but it must not be treated as finished test coverage. For each critical scenario, record one of these statuses:

| Status | Meaning | Required next step |
| --- | --- | --- |
| `automated` | Scenario is covered by unit, contract, API, integration, E2E, fixture, or runner test. | Link the test path/name. |
| `fixture-backed` | Scenario is proven by checked-in input/output fixtures but not a full runner. | Link fixture and assertion owner. |
| `manual-evidence` | Scenario requires manual, UI, bench, or external service evidence. | Record run steps, evidence, and limits. |
| `pending-runner` | Gherkin exists but no runner/step definition is wired. | Add runner choice or map to an executable test type. |
| `not-automatable` | Tooling cannot enforce it directly. | State the manual review or release checklist item. |

BDD closure does not require every scenario to use a Cucumber-style runner. It does require every critical scenario to have an explicit validation path and no ambiguous "documented but untested" state.

## Test Strategy Gate

Separate "guarding old behavior" from "validating new code" before implementation. High total coverage can prove old behavior is protected, but it does not prove newly generated or newly written code is correct.

| Target | Purpose | Required validation |
| --- | --- | --- |
| Existing / legacy behavior | Prevent regressions and protect known contracts. | Run existing unit, BDD, contract, integration, and regression tests that cover impacted behavior. |
| New requirement or new code | Prove the new behavior is correct, safe, and aligned with docs. | Write or update BDD first, then add failing tests or executable specs before production code when feasible. Track changed/new-code coverage separately from whole-project coverage. |
| AI-generated code | Compensate for plausible but wrong code. | Require BDD scenario, unit/contract tests, and human review focused on intent, edge cases, and security/ownership boundaries. |
| Business rules / algorithms | Catch examples that pass but rules that are wrong. | Add property-based tests, invariant tests, or table-driven edge cases. |
| Critical conditionals / validation logic | Prove tests fail when logic is wrong. | Add mutation testing where practical, or manually test negative cases that would fail if guards were removed. |
| Database / persistence behavior | Protect real state transitions and migrations. | Add fixture-backed repository tests, migration tests, or integration tests against representative data. |

Recommended order for new requirements:

1. BDD scenarios.
2. Failing unit, contract, property, or integration tests for the new behavior.
3. Production code.
4. Mutation/negative checks for critical rules.
5. Human review with the planning docs, BDD, and tests side by side.

## Embedded / Hardware Product Flow

Use this flow when the project involves firmware, sensors, boards, UART/I2C/SPI/BLE/CAN/GPIO, RTOS tasks, hardware bring-up, or host/target validation:

| Layer | Contract | Notes |
| --- | --- | --- |
| Datasheet / vendor spec | Electrical interface, protocol bytes, timing, default parameters, valid ranges, errata. | Treat vendor docs as the protocol truth; record observed deviations separately. |
| Protocol Parsing Contract | Frame format, state machine, length/checksum rules, command/ACK/report shapes, fixtures, invalid frames. | Keep byte-level parsing separate from product meaning. |
| BDD Behavior | User/system behavior, device states, setup/config flows, fault handling, target events. | BDD uses domain terms, not raw registers or UART calls. |
| Domain Model Contract | Pure DTOs, units, ranges, invariants, timestamps, validity windows. | Keep HAL/RTOS types out of domain objects. |
| Hardware Context Contract | Board-specific pins, UART/I2C/SPI bus, baud/rates, buffers, interrupts, power modes, injected configuration. | Board changes should change context/config, not protocol/domain logic. |
| Embedded Architecture Contract | Driver/service/application layering, task/ISR boundaries, queues, ownership, concurrency, lifecycle, error escalation. | Drivers handle bytes; services parse; applications decide product behavior. |
| Public API / Interface Contract | Context lifecycle, callbacks/subscriptions, commands, errors, consumer ownership, multi-device rules. | Avoid parallel second context APIs unless contracts are revised first. |
| Test Plan | Host unit tests, protocol fixtures, negative cases, property/invariant tests, simulator/mocks, hardware-in-loop, bring-up log evidence. | Separate host-repeatable proof from bench-only evidence. |

Before firmware code:

1. Read the datasheet/protocol spec and project contracts.
2. Confirm hardware context is injectable per board and not hard-coded as the only source of truth.
3. Write or update BDD and protocol/domain/API contracts.
4. Add host-side fixtures for protocol parsing and negative cases.
5. Define target or hardware-in-loop validation only for evidence that cannot be proven on host.
6. Record bring-up evidence: board revision, wiring, pins, bus settings, firmware version, logs, and known deviations.

## Missing Information Gate

Before development planning or implementation continues, missing information must be handled explicitly:

| Missing item type | Required action |
| --- | --- |
| Can be recovered from evidence | Backfill it and cite the evidence source. |
| Product intent not recoverable | Mark `unknown` / `open question`, ask the user, and do not invent intent. |
| Affects BDD behavior, domain invariants, API/interface shape, error handling, security, storage, or tests | Treat as a blocker: ask the user or request evidence before continuing implementation. |
| Nice-to-have context that does not change behavior or contracts | Record as non-blocking open question and state why it does not block. |

Do not continue development with unresolved blockers. The agent must list missing items as questions, wait for answers or evidence, then update the documents before proceeding.

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

For implemented-first projects, also recover the delivery pipeline:

| Pipeline artifact | Backfill rule |
| --- | --- |
| Plan index / product radar | Map source product docs, PDFs, tickets, screenshots, or legacy notes to modules, controllers, screens, commands, or packages. Mark canceled or superseded requirements. |
| Contract taxonomy | List which documents govern build/run, HTTP/API shape, auth/tenant/session, persistence, domain layering, frontend/backend integration, third-party integration, testing, and documentation sync. |
| Minimum doc sync matrix | For each change type, state the minimum docs/tests to update: API, permission, database, UI flow, generated client, vendor integration, CLI command, diagnostic rule, release setting. |
| OpenAPI / schema / generated client | Verify the generated consumer code comes from the source contract and does not hand-copy endpoints or DTOs. |
| Vendor / third-party integration | Separate raw vendor docs from sanitized integration excerpts, request/response contracts, fixture examples, live-test gates, and secret handling. |
| Tooling / extension rule catalog | Map catalog order, rule IDs, diagnostics/commands, fixtures, and tests; mark process-only or non-enforceable rules explicitly. |

Backfill order for existing projects:

1. Inventory existing docs, source folders, tests, schemas, API specs, fixtures, release notes, and observed behavior.
2. Create a documentation gap table with status: `exists`, `partial`, `missing`, or `unknown`.
3. Backfill BDD Behavior first when product brief is missing, because implemented behavior is the strongest available source of truth.
4. Backfill Domain Model, Architecture, API / Interface, and Error Handling Contracts from the completed behavior and implementation evidence.
5. Mark unknown product intent separately from observed behavior. Unknown intent does not block BDD completion.
6. If BDD cannot be completed from available evidence, stop and ask for the missing behavior, screen/API examples, logs, test cases, or user decisions before continuing development.
7. Add tests or test TODOs for any critical BDD scenario that lacks coverage.

## Contract-First Rules

- BDD describes behavior; it should not lock in framework or database choices.
- Domain Model Contract owns invariants, business vocabulary, and state transitions.
- Architecture Contract owns dependency direction, runtime boundaries, data ownership, and allowed integration paths.
- API Contract owns integration shape: request, response, error, auth/session, versioning, and compatibility.
- Error Handling Contract owns failure taxonomy, retry policy, user messaging, logging, and security redaction.
- Contract Governance owns document precedence, conflict handling, canceled/deferred scope, and minimum linked updates.
- New requirements must update planning docs, BDD, contracts, implementation slices, and tests before code starts.
- Bug fixes must identify expected vs actual behavior and the regression test before code starts.
- New or AI-generated code must be validated with tests that target the changed behavior, not only total project coverage.
- Use mutation, property-based, contract, or database-backed tests when ordinary examples do not prove the rule.
- Embedded changes must distinguish datasheet/protocol truth, hardware context, driver/service/application ownership, host-testable logic, and target-only evidence.
- Implementation can run in parallel only when the shared contracts are versioned enough for mock, stub, or schema-first work.
- If a contract changes, update BDD, implementation, mocks, and tests in the same change or explicitly record why not.
- If an API/schema contract changes, regenerate typed clients or SDKs from the source contract; do not hand-copy routes, DTOs, or operation names.
- If a third-party integration changes, update sanitized integration docs, fixtures, live-test gates, and secret/redaction notes without copying private vendor or account details into reusable guidance.
- For already implemented projects, BDD becomes the required behavioral recovery document. Product Brief may contain unknowns, but BDD must be filled from observable product behavior and implementation evidence.
- Any missing information that changes behavior, contracts, ownership, error handling, storage, security, or tests blocks development until it is answered or explicitly scoped out.

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
- Test strategy that distinguishes existing-regression coverage from changed/new-code validation.
- No unresolved blocker questions that affect implementation behavior or contracts.

For an already implemented project, "ready" means the missing-document audit is complete and BDD covers the implemented critical behavior, even if original product intent remains partly unknown.

## Minimum Definition Of Done

Before shipping or merging:

- Domain invariants are tested.
- Contract tests pass for provider and consumer.
- Mocks/fixtures match the latest contract.
- Integration test covers at least the critical happy path and one important failure path.
- Residual unknowns or deferred behavior are documented in the project repository.

← [Back to App Development Guidance](../README.md)
