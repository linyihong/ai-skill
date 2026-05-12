# Risk Translation & Control Selection

本文件定義如何將分析觀察轉換為開發者視角的風險陳述，並選擇最適合的控制層。承接 [`skills/app-development-guidance/WORKFLOW.md`](../../skills/app-development-guidance/WORKFLOW.md) §2-4 的內容。

## 1. Translate To Risk

Write the risk in developer terms:

| Observation | Developer takeaway |
| --- | --- |
| Request can be replayed with the same body and token | Backend contract may lack nonce, timestamp, idempotency, or risk scoring. |
| Token is long-lived and broadly scoped | Account takeover blast radius is too large if token leaks. |
| Sensitive values appear in logs | Debug or telemetry pipeline may leak secrets. |
| App relies on hidden client logic for authorization | Server may be trusting client state that attackers can modify. |
| Local storage contains recoverable secrets | Device compromise or backups can expose credentials. |
| Release build exposes symbols, debug flags, or test endpoints | Reverse engineering cost is unnecessarily low. |
| Generated client diverges from OpenAPI/schema | Frontend, SDK, or tool consumers may call stale routes or deserialize wrong shapes. |
| Gherkin exists without executable linkage | The team may believe behavior is tested when it is only documented. |
| Vendor docs are copied directly into product flow | Secrets, irrelevant fields, or unstable third-party assumptions can leak into implementation and reusable docs. |

## 2. Choose The Owner Layer

Prefer controls owned by the strongest layer:

| Layer | Good use |
| --- | --- |
| Backend/API | Authorization, replay defense, rate limits, fraud signals, token rotation, audit logs. |
| Client app | Safe storage, secure defaults, pinning where justified, UX friction for risky flows, telemetry hygiene. |
| Full-stack contract | OpenAPI/schema generation, typed clients, provider/consumer fixtures, contract tests, compatibility rules. |
| Tooling / extension | Pure rule engine or command kernel, adapter boundary, diagnostics/commands, fixture-backed rules, editor/CLI integration tests. |
| Third-party integration | Sanitized vendor excerpt, credentials boundary, live-test gate, replay/idempotency/webhook/settlement behavior, audit logs. |
| Embedded/firmware | Sensor/protocol parsing, hardware context injection, driver/service/application boundary, RTOS/task ownership, host fixtures, hardware-in-loop checks. |
| Build/release | Obfuscation, symbol stripping, debug flag enforcement, dependency review, secret scanning. |
| Monitoring | Anomaly detection, device risk signals, abuse pattern alerts. |

Client-side hardening raises cost and improves signal quality, but it must not be the only control for authorization or financial/business integrity.

## 3. Define Controls

For each risk or useful implementation lesson, define:

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

## 4. File The Guidance

Classify the outcome before writing docs:

| If the lesson is about... | Put it in |
| --- | --- |
| Security property that applies across stacks | `controls/` |
| Mobile, web, backend, embedded, firmware, hardware, or OS-specific implementation | `platforms/` |
| Dart, Kotlin/Java, Swift, TypeScript, or runtime-specific pitfalls | `languages/` |
| Concrete buildable pattern or how-to | `implementation/` |
| A repeated design, PR, release, or API review step | `checklists/` |
| A reusable but still emerging lesson | `feedback_history/<category>/` or `feedback_history/common/` |
| A copyable documentation shape | `templates/` |
| Missing development documents in an implemented project | `process/` and `templates/initial-development-docs.md` |
| Blocker questions for missing requirements or contracts | `process/` and the current planning document |
| Change intake before code | `process/` and the current planning document |
| Test strategy for new or AI-generated code | `process/`, `CHECKLIST.md`, and the current planning document |
| Embedded/hardware product flow | `platforms/embedded/`, `implementation/embedded/`, `process/`, and hardware-aware checklists |
| OpenAPI/schema/codegen or full-stack provider/consumer contract | `implementation/backend/`, `process/`, `CHECKLIST.md`, and relevant API checklists |
| Tool, CLI, IDE extension, linter, or static-analysis architecture | `implementation/tooling/`, `process/`, and relevant review checklists |
| Vendor or third-party API integration | `implementation/backend/`, `controls/`, `checklists/`, and project-specific sanitized docs |

Prefer linking between folders over duplicating the same guidance.

## 5. Apply Required Linked Updates

Before finishing a change, follow the repo-wide rule in [`shared-rules/linked-updates.md`](../../shared-rules/linked-updates.md). If the update affects related folders, those linked updates are **required** in the same change:

| Changed area | Must update or verify |
| --- | --- |
| `controls/` | Relevant `implementation/`, `platforms/`, `languages/`, and `checklists/` docs. |
| `implementation/` | Relevant `controls/`, `platforms/`, `languages/`, and `checklists/` docs. |
| `platforms/` | Relevant `controls/`, `implementation/`, language notes, and checklists. |
| `languages/` | Relevant platform and implementation docs. |
| `checklists/` | Relevant controls and implementation docs. |
| `templates/` | `templates/README.md`, `DOCUMENTATION.md`, and any docs that instruct users to copy the template. |
| `process/` governance or backfill rules | `templates/`, `CHECKLIST.md`, `WORKFLOW.md`, and relevant `implementation/` docs. |

If no linked update is needed, state why. Do not call linked updates optional when they are required for consistency.
