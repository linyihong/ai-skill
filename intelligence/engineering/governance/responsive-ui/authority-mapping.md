# Responsive Failure Authority Mapping

Authority mapping answers who can judge a responsive finding as failed, blocked, warning, or not applicable.

## Authority Types

| Authority | Can Decide | Cannot Decide Alone |
| --- | --- | --- |
| `contract` | Whether the render context and invariant are in scope. | Runtime failure without evidence. |
| `validation` | Objective failure proven by metrics, DOM, screenshot baseline, runtime trace, or accessibility evidence. | Subjective design preference without criteria. |
| `review` | UX judgment, density trade-off, visual hierarchy, or design preference. | Mechanical pass/fail unless backed by objective rubric. |
| `project_policy` | Framework adoption, design-system token policy, supported contexts, release thresholds. | One-off evidence that contradicts runtime observations. |

## Mapping

| Finding | Authority | Default Severity | Release Posture |
| --- | --- | --- | --- |
| `horizontal_overflow` | validation | high | block candidate |
| `primary_cta_clipped` | validation | critical | blocked |
| `fixed_surface_drift` | validation | high | block candidate |
| `safe_area_overlap` | validation | high | block candidate |
| `dynamic_resize_staleness` | validation | high | block candidate |
| `minor_typography_shift` | review | low | warning |
| `density_collapse` | review + contract | medium | fix or defer |
| `unsupported_context_claim` | contract + project_policy | not_applicable or high | no impact if explicit; block candidate if implicit |
| `unscoped_framework_mix` | project_policy + review | medium | fix or defer |

## Rule

If authority is unclear, do not silently downgrade the issue. Name the missing authority: contract, validation evidence, review rubric, or project policy.
