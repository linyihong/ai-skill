# Documenting App Security Hardening Notes

Use this guide when turning an APK analysis lesson or design review into reusable development guidance.

## Reusable Note Structure

```markdown
### Short title

Status: candidate | validated | promoted | deprecated | experimental

#### Observed Risk

What pattern was observed, without target-specific secrets.

#### Development Consequence

Why this matters when building our own app or API.

#### Recommended Control

What to implement, and which layer owns it.

#### Validation

How to prove the control exists or fails safely.

#### Limits

What this does not solve.
```

## Keep Separate

| Content | Put it in |
| --- | --- |
| Reusable development principle | This skill. |
| APK analysis method or hook technique | [`apk-analysis`](../apk-analysis/). |
| Product-specific API host, endpoint, schema, or token detail | Project repository docs. |
| Shared sanitization or feedback rules | [`shared-rules`](../../shared-rules/README.md). |

## Good Guidance

Good hardening guidance is:

- Actionable by an engineer.
- Explicit about the owner layer.
- Testable.
- Honest about residual risk.
- Sanitized and free of target-specific details.

## Avoid

- "Use obfuscation" without naming what it protects and what it does not.
- "Add pinning" without a rotation plan or threat model.
- "Detect root" as a hard authorization decision.
- "Hide the secret in the app" as a durable security control.
- Copying raw findings from a third-party APK into reusable docs.
