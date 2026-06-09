# UI Governance Evidence: <surface / change>

## Scope

- **Surface / flow**: <screen, component, route, CLI, SDK, job, or consumer surface>
- **Upstream UI contract**: <link to Screen / Consumer / ViewModel / Accessibility contract>
- **Claim**: <what UI compliance is being claimed>
- **Non-goals**: <visual taste, future polish, or scoped-out checks>

## Classification

| Field | Value | Evidence / Link |
| --- | --- | --- |
| Governance domain | <Contract / Design System / Accessibility / Responsive / Behavior / Closure / not_applicable> | <link> |
| Render context | <desktop / tablet / mobile / narrow_mobile / not_applicable> | <link to contract or validation matrix> |
| Collection method | <contract_readback / static_analysis / runtime_trace / browser_review / human_observation / not_applicable> | <link to acquired artifact or source> |
| Validation mechanism | <deterministic / screenshot_diff / ai_review / manual_review / not_applicable> | <link> |
| Evidence class | <contract / runtime / accessibility_scan / visual_diff / screenshot / responsive_capture / ai_review / human_review / not_applicable> | <link> |
| Severity | <block_candidate / warn / research / not_applicable> | <reason> |

## Project-Local Design System Policy

- **Policy source**: <project token / primitive / component policy>
- **Token / primitive expectation**: <what must be tokenized or use approved primitives>
- **Out of scope**: <global token scale or taste preference not governed here>

## Evidence Notes

- **Acquisition**: <how evidence was obtained; browser review / static scan / contract readback / runtime trace / human observation>
- **Responsive / render contexts**: <desktop + mobile evidence, validation matrix, viewport metrics, safe-area / keyboard / orientation evidence, or not_applicable>
- **Accessibility**: <keyboard, focus, semantics, labels, contrast, motion, or not_applicable>
- **Behavior pattern**: <loading, empty, error, destructive confirmation, retry, offline, permission denied, or not_applicable>
- **Visual baseline**: <baseline name, capture determinism, diff status, or not_applicable>
- **AI visual review**: <rubric, prompt scope, warning/research status, or not_applicable>

## Closure

- **Decision**: <pass / warn / block_candidate / deferred>
- **Linked validation**: <test, scan, screenshot, review report, or manual evidence>
- **Deferred scope**: <owner, follow-up, reason, or none>
- **Runtime projection**: advisory-only unless a separate promotion decision names executor and thresholds.
