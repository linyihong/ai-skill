> Follow [enforcement rules](../../../../enforcement/README.md) and [feedback-lessons](../../../../enforcement/feedback-lessons.md). This lesson is generalized and contains no target-specific hosts, endpoints, tokens, or user data.
# Extracted — See [`analysis/development-guidance/controls-catalog.md`](../../../../analysis/development-guidance/controls-catalog.md)

### 2026-05-01 - Client Encrypted Header Is Not A Security Boundary

Status: promoted

#### One-line Summary

If a mobile client generates encrypted or signed request headers, assume an authorized analyst or attacker with the app binary can recover the plaintext construction and keys or inputs; backend controls must provide the real security boundary.

#### Human Explanation

Mobile request signing, encrypted parameter headers, obfuscation, and local proxy routing can raise reverse-engineering cost and provide useful tamper signals. They should not be treated as authorization, replay defense, or data confidentiality against someone who can inspect the shipped client. The app binary contains enough logic to build the same header, and dynamic hooks can often observe plaintext immediately before encryption or immediately after decryption.

#### Trigger

- APK analysis finds custom headers that carry encrypted request parameters.
- Static AOT/native analysis recovers functions such as request interceptors, signing, encoding, encryption, or decryption.
- Dynamic hooks can observe plaintext before encryption or decoded response data after decryption.

#### Evidence

- Tool: Flutter/Dart AOT function mapping plus Frida native offset hooks.
- Sanitized excerpt: a client-side request interceptor produced plaintext shaped like `<prefix>|<timestamp>|<nonce>|<path>|<query/form material...>` immediately before encryption, and returned an encrypted/base64-like header value.
- Evidence path: raw evidence remains in the project-private capture folder; this reusable lesson only records the generalized development conclusion.

#### Generalized Lesson

Do not design backend trust around the assumption that client-side signing or encrypted headers are secret. The backend should validate account/session authorization, freshness, nonce or idempotency, token scope, rate limits, and anomaly signals independently. Client signatures can be one input into risk scoring, but a correctly formatted client-generated header should not prove that the request is legitimate.

#### Agent Action

When converting APK analysis into hardening guidance:

1. Put reverse-engineering mechanics in `apk-analysis`.
2. Put the development action here: server-side authorization, replay protection, token hygiene, and logging controls.
3. Avoid writing target-specific header names, hosts, endpoint paths, token values, or plaintext examples.
4. Add a validation method that proves replay, token leakage, or modified-client bypasses are rejected server-side.

#### Applies When

- Mobile, desktop, or web clients generate request signatures or encrypted parameter wrappers.
- The protected action has account, payment, entitlement, wallet, creator income, private data, or abuse implications.
- The same client binary is distributed to untrusted devices.

#### Does Not Apply When

- The client-side wrapper is only for compression or compatibility and not presented as a security control.
- Secrets are generated and held only on a trusted server or hardware-backed attestation flow, and the client never receives signing authority.

#### Validation

- Replay a previously valid signed/encrypted request and verify rejection or idempotent handling.
- Change account/session-bound fields and verify server-side authorization rejects the request.
- Scan logs and telemetry for tokens, encrypted header plaintext, query credentials, and decrypted payloads.
- Confirm monitoring treats client signatures as risk signals, not as the sole source of truth.

#### Promotion Target

- `controls/api-transport.md`
- `controls/auth-session.md`
- `checklists/api-security-review.md`
