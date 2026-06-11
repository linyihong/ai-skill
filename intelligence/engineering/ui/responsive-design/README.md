# Responsive Design Intelligence

Responsive design intelligence captures framework-neutral layout principles. Bootstrap is useful evidence because it has encoded these economics for years, but the reusable primitive is not Bootstrap itself.

## Atoms

| Atom | Principle |
| --- | --- |
| [`mobile-first.md`](mobile-first.md) | Start from the narrowest viable viewport, then progressively enhance. |
| [`fluid-container.md`](fluid-container.md) | Containers should be fluid first and capped later with `max-width`. |
| [`breakpoint-strategy.md`](breakpoint-strategy.md) | Breakpoints describe layout constraints, not device SKUs. |
| [`viewport-behavior.md`](viewport-behavior.md) | Viewport height and browser chrome are dynamic runtime inputs. |
| [`fixed-sticky-layout.md`](fixed-sticky-layout.md) | Fixed and sticky UI must stay aligned to the same container contract as content. |
| [`overflow-management.md`](overflow-management.md) | Overflow is a signal to fix the source layout, not something to hide by default. |
| [`safe-area-layout.md`](safe-area-layout.md) | Mobile layouts must explicitly account for safe areas and constrained system UI. |

## Use

Workflow documents may reference these atoms when defining UI contracts, responsive validation, or review checklists. Do not copy the full atom text into workflow files.

For governance classification, authority mapping, severity, and closure checklist, use [`../../governance/responsive-ui/`](../../governance/responsive-ui/README.md). For shared context names, use [`../../render-contexts/`](../../render-contexts/README.md).
