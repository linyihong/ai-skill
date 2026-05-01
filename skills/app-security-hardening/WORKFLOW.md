# App Security Hardening Workflow

Use this workflow to convert analysis observations into development guidance without overclaiming what client-side defenses can guarantee.

## 1. Start From Evidence

Record the reusable observation:

- What behavior was observed?
- Which layer exposed it: client code, transport, API contract, storage, logs, build config, or runtime behavior?
- Is the issue confirmed, suspected, or only a risk pattern?

Do not copy target-specific endpoints, tokens, secrets, device IDs, or raw user data into this skill.

## 2. Translate To Risk

Write the risk in developer terms:

| Observation | Developer risk |
| --- | --- |
| Request can be replayed with the same body and token | Backend contract may lack nonce, timestamp, idempotency, or risk scoring. |
| Token is long-lived and broadly scoped | Account takeover blast radius is too large if token leaks. |
| Sensitive values appear in logs | Debug or telemetry pipeline may leak secrets. |
| App relies on hidden client logic for authorization | Server may be trusting client state that attackers can modify. |
| Local storage contains recoverable secrets | Device compromise or backups can expose credentials. |
| Release build exposes symbols, debug flags, or test endpoints | Reverse engineering cost is unnecessarily low. |

## 3. Choose The Owner Layer

Prefer controls owned by the strongest layer:

| Layer | Good use |
| --- | --- |
| Backend/API | Authorization, replay defense, rate limits, fraud signals, token rotation, audit logs. |
| Client app | Safe storage, secure defaults, pinning where justified, UX friction for risky flows, telemetry hygiene. |
| Build/release | Obfuscation, symbol stripping, debug flag enforcement, dependency review, secret scanning. |
| Monitoring | Anomaly detection, device risk signals, abuse pattern alerts. |

Client-side hardening raises cost and improves signal quality, but it must not be the only control for authorization or financial/business integrity.

## 4. Define Controls

For each risk, define:

- Required control.
- Owner.
- Implementation note.
- Validation method.
- Residual risk.

Example:

```text
Risk: API request can be replayed.
Control: Server verifies timestamp, nonce, account/session binding, and idempotency key.
Owner: Backend/API.
Validation: Integration test replays the same signed request and expects rejection or idempotent handling.
Residual risk: Device compromise can still steal valid sessions; monitoring remains required.
```

## 5. File The Guidance

Classify the outcome before writing docs:

| If the lesson is about... | Put it in |
| --- | --- |
| Security property that applies across stacks | `controls/` |
| Mobile, web, backend, or OS-specific implementation | `platforms/` |
| Dart, Kotlin/Java, Swift, TypeScript, or runtime-specific pitfalls | `languages/` |
| A repeated design, PR, release, or API review step | `checklists/` |
| A reusable but still emerging lesson | `feedback_history/` |
| A copyable documentation shape | `templates/` |

Prefer linking between folders over duplicating the same guidance.

## 6. Validate

Use at least one validation method:

- Unit or integration test.
- Release checklist item.
- Static scan or build assertion.
- Manual review with evidence.
- Runtime or backend telemetry query.

## 7. Feed Back Reusable Lessons

If a lesson generalizes beyond one product:

1. Add a file under `feedback_history/`.
2. Link shared rules instead of duplicating them.
3. Promote validated guidance into the structured folders, [CHECKLIST.md](CHECKLIST.md), or this workflow.

If the lesson came from APK analysis, keep the analysis method in [`apk-analysis`](../apk-analysis/) and the development action here.
