# Vendor And Third-Party API Integration

Use this pattern when a project connects to payment providers, identity providers, shipping services, analytics vendors, app stores, external webhooks, or other third-party APIs.

## Documentation Split

Keep these documents separate:

| Document | Owns | Must not contain |
| --- | --- | --- |
| Raw vendor source | Vendor PDF, portal docs, public docs, or private contract kept in the project repository. | Reusable skill content, copied secrets, or unrelated implementation decisions. |
| Integration excerpt | Sanitized operations actually used by the product: request/response shape, auth, signing, idempotency, callbacks, retries, sandbox/live differences. | Full vendor docs, credentials, account-specific hosts, private business terms. |
| Product contract | Domain behavior, BDD, error mapping, audit, settlement/reconciliation, user-visible states. | Vendor-only vocabulary when a product term exists. |
| Test evidence | Fixtures, sandbox logs, gated live tests, replay/idempotency checks, webhook verification evidence. | Long-lived secrets or raw personal data. |

## Required Flow

1. Identify the exact vendor operations used and mark unused/canceled operations out of scope.
2. Convert vendor terms into product domain terms before application code consumes them.
3. Define auth, signing, timestamp/nonce, idempotency, retry, timeout, callback/webhook, and audit behavior.
4. Store credentials and account-specific config outside docs and fixtures.
5. Add fixture-backed tests for request construction, response parsing, error mapping, and signature/webhook verification.
6. Gate live connectivity tests behind explicit environment variables or config so normal CI does not hit real services.
7. Record sandbox/live differences and operational runbook notes in the project repository.

## Security And Reliability Checks

- Requests are bound to account/session/order context where relevant.
- Replay-sensitive operations use idempotency keys, nonces, timestamps, or provider-supported duplicate detection.
- Webhooks validate signature, timestamp, source, body binding, and idempotent processing.
- Errors distinguish retryable, user-fixable, vendor outage, fraud/risk, and permanent failures.
- Logs redact tokens, keys, signatures, account identifiers, personal data, and raw vendor payloads unless explicitly approved.
- Settlement, reconciliation, or final-state polling is defined for asynchronous providers.

## Test Strategy

| Test | Proves |
| --- | --- |
| Fixture request test | Product code signs/builds the expected request without secrets. |
| Fixture response test | Product code maps vendor payloads to domain states and errors. |
| Replay/idempotency test | Duplicate callback or request does not double-apply business state. |
| Webhook verification test | Mutated body/signature/timestamp is rejected. |
| Gated live test | Sandbox or live connectivity works only when explicitly enabled. |

## Required Linked Updates

Follow [`../../../../shared-rules/linked-updates.md`](../../../../shared-rules/linked-updates.md). Vendor integration changes must update or verify [`../../controls/api-transport.md`](../../controls/api-transport.md), [`../../controls/logging-telemetry.md`](../../controls/logging-telemetry.md), [`../../platforms/backend/api.md`](../../platforms/backend/api.md), [`../../checklists/api-security-review.md`](../../checklists/api-security-review.md), and project-specific sanitized docs.
