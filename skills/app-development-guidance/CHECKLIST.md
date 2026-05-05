# App Development Guidance Checklist

Use this checklist during design review, PR review, and release review. It is intentionally high level: project-specific requirements should live in the product repository.

For focused checklists, use:

- [`checklists/mobile-design-review.md`](checklists/mobile-design-review.md)
- [`checklists/mobile-pr-review.md`](checklists/mobile-pr-review.md)
- [`checklists/mobile-release-review.md`](checklists/mobile-release-review.md)
- [`checklists/api-security-review.md`](checklists/api-security-review.md)

When a checklist item changes because of a new control or implementation pattern, the linked `controls/` and `implementation/` docs must be updated or explicitly verified in the same change.

## Change Intake

- The project's 企劃書, product brief, planning doc, issue, ticket, PRD, design note, BDD, API contract, or equivalent artifact was reviewed before code.
- The change is classified as new requirement, bug fix, refactor, hardening, or documentation-only.
- New requirements update planning docs, BDD, affected contracts, implementation slices, and test plan before code starts.
- Bug fixes document expected vs actual behavior, reproduction/evidence, affected BDD scenario or missing scenario, impacted contracts/errors, and regression test plan.
- Refactors are confirmed to have no behavior or public contract change; otherwise they are reclassified.
- Blocker questions are answered, backed by evidence, or explicitly scoped out before implementation.

## Test Strategy

- Existing/legacy behavior has regression coverage for the affected paths.
- New or changed behavior has BDD scenarios before production code.
- New code has failing unit, contract, property, integration, or executable spec coverage before implementation when feasible.
- Changed/new-code coverage is checked separately from total project coverage.
- Mutation testing, property-based testing, invariant tests, or negative cases cover rule-heavy or safety-sensitive logic.
- Database, repository, migration, or persistence behavior is verified with fixtures or integration tests when state matters.
- AI-generated code receives human review against planning docs, BDD, contracts, edge cases, and security/ownership boundaries.

## Product To Contract Flow

- Product brief names goals, users, scope, non-goals, assumptions, and constraints.
- Bounded Contexts or modules are split by domain responsibility and integration boundary.
- Critical behavior is written as BDD scenarios before implementation.
- Domain Model Contract defines entities, value objects, commands, events, and invariants.
- Architecture Contract defines dependency direction, data ownership, runtime boundaries, and allowed integrations.
- API, event, command, or public interface contract is defined before parallel implementation.
- Error Handling Contract defines error taxonomy, retry rules, user messaging, logging, and redaction.
- If there is no frontend/backend split, producer and consumer roles are still named.
- Mock APIs, fixtures, schemas, or stubs are generated from the latest contract.
- Unit, BDD, contract, and integration test responsibilities are assigned before build work starts.
- Missing behavior, domain, API/interface, error handling, security, storage, ownership, or test requirements are asked as blocker questions before development continues.

## Existing Project Documentation Backfill

- Existing project docs are inventoried and marked `exists`, `partial`, `missing`, or `unknown`.
- Missing Product Brief fields are reconstructed only from evidence; unavailable original intent is marked `unknown` or `open question`.
- BDD behavior is complete for implemented critical happy paths, failure paths, permissions, empty states, edge cases, and cross-context flows.
- BDD scenarios cite evidence from UI behavior, API behavior, code paths, tests, logs, fixtures, or manual verification.
- Domain Model, Architecture, API / Interface, and Error Handling Contracts are backfilled from observed behavior and implementation evidence.
- Every critical BDD scenario maps to existing test coverage or a required test gap.
- Any gap that cannot be backfilled from evidence and affects behavior or contracts is asked before implementation proceeds.

## API And Transport

- Sensitive flows use HTTPS only; cleartext traffic is disabled in release builds.
- Certificate pinning is considered for high-risk apps, with a rotation and incident plan.
- Backend authorization does not trust client-only flags, roles, prices, balances, or feature gates.
- Replay-sensitive requests have server-side timestamp, nonce, idempotency, or risk checks.
- Request signing, if used, signs the right fields and does not rely on a static client secret.
- Error responses avoid leaking stack traces, internal hostnames, or sensitive business rules.

## Auth, Tokens, And Sessions

- Access tokens are scoped, time-limited, and revocable.
- Refresh flows are rate-limited and bound to account/session context.
- Logout, password change, and risk events invalidate relevant sessions.
- Tokens are not written to debug logs, crash reports, analytics, or screenshots.
- Session identifiers are not used as long-term device identifiers.

## Local Storage

- Secrets are stored only when needed, and use platform-backed secure storage where possible.
- Cache files, SQLite, shared preferences, and downloaded media are reviewed for sensitive content.
- Backups and screenshots are configured according to product risk.
- Offline data has an expiry, encryption plan, or clear business justification.

## Flutter And Android Build

- Release builds disable debug flags, test endpoints, dev menus, and verbose network logging.
- Obfuscation/minification is enabled where compatible with crash symbolication and support needs.
- Native symbols are stripped unless needed for diagnostics.
- Flutter platform channels do not expose privileged operations without server-side or OS-level checks.
- Third-party SDKs are reviewed for permissions, telemetry, and data collection.

## Logging And Telemetry

- Logs redact tokens, cookies, authorization headers, device IDs, and personal data.
- Crash reports and analytics events avoid raw request/response payloads.
- Debug logging is gated away from release builds.
- Security-relevant failures are observable without exposing secrets.

## Anti-Tamper And Risk Signals

- Root/jailbreak/emulator/hook detection is treated as a risk signal, not a sole access-control decision.
- Server-side risk scoring can tolerate false positives and false negatives.
- Critical operations still require backend authorization and abuse controls.
- The app avoids storing static secrets that become permanent bypass targets.

## Release Gate

- A reviewer can point to tests, build checks, or documented evidence for every required control.
- No unresolved blocker questions remain for behavior, contracts, errors, security, storage, ownership, or tests.
- Known residual risks are documented in the project repository.
- Reusable lessons are generalized into this skill only after sanitization.
