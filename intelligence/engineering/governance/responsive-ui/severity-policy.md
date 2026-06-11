# Responsive Severity Policy

Responsive severity is based on user impact, evidence objectivity, and release risk.

## Severity Levels

| Severity | Use When | Release Posture |
| --- | --- | --- |
| `critical` | Primary action, navigation, payment, authentication, consent, destructive action, or required content is clipped or unreachable. | Blocked |
| `high` | Unintended horizontal overflow, fixed navigation drift, or dynamic resize staleness affects a supported render context. | Block candidate |
| `medium` | Layout is usable but degraded in a supported context, with clear evidence and no critical action blocked. | Fix or explicitly defer |
| `low` | Minor visual shift, subjective density concern, or review-only issue without objective contract break. | Warning |
| `not_applicable` | Context is explicitly unsupported or the apparent overflow is an intended scoped scroller. | No release impact |

## Rules

- Objective overflow or clipped controls should not be downgraded to visual taste.
- Subjective visual preference should not block release unless translated into an objective project policy.
- Severity can rise when the affected surface is part of a critical journey.
- Unsupported contexts must be explicit; silence is not an unsupported-context declaration.
