# Responsive UI Governance

Responsive UI Governance defines how to classify and judge UI correctness across render contexts. It is governance knowledge, not a standalone workflow.

## Read When

- A UI claim depends on narrow mobile, mobile, tablet, desktop, landscape, safe-area, embedded, or dynamic resize behavior.
- A defect report mentions clipping, horizontal overflow, fixed/sticky drift, browser chrome, safe area, or reload after viewport changes.
- A project considers Bootstrap or another mature responsive framework as a design source.
- A workflow needs objective severity and failure-authority mapping for responsive UI findings.

## Surfaces

| File | Purpose |
| --- | --- |
| [`contract.md`](contract.md) | Defines the responsive UI contract fields a project or workflow should name. |
| [`validation-matrix.md`](validation-matrix.md) | Maps render contexts to evidence and validation expectations. |
| [`failure-taxonomy.md`](failure-taxonomy.md) | Names common responsive failure classes. |
| [`severity-policy.md`](severity-policy.md) | Maps responsive findings to release risk. |
| [`authority-mapping.md`](authority-mapping.md) | Defines who has authority to judge each failure type. |
| [`checklist.md`](checklist.md) | Compact review and closure checklist. |

## Source Intelligence

- [`responsive-design/`](../../ui/responsive-design/README.md)
- [`render-contexts/`](../../render-contexts/README.md)
- [`mobile-first-economics.md`](../../economics/mobile-first-economics.md)
- [`responsive-cost-curve.md`](../../economics/responsive-cost-curve.md)
- [`viewport-variance-management.md`](../../heuristics/viewport-variance-management.md)
- [`framework-pattern-extraction.md`](../../heuristics/framework-pattern-extraction.md)
