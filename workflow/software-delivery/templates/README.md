# Software Delivery Templates

Use these templates as focused artifact shapes. Load only the template that matches the current artifact.

| Template | Use when |
| --- | --- |
| [`change-brief-template.md`](change-brief-template.md) | Capturing intake, scope, assumptions, acceptance, and validation target before implementation. |
| [`contract-template.md`](contract-template.md) | Defining domain/API/UI/consumer contracts and traceability before parallel implementation. |
| [`bdd-scenario-template.md`](bdd-scenario-template.md) | Writing behavior scenarios and acceptance examples. |
| [`implementation-plan-template.md`](implementation-plan-template.md) | Planning implementation slices, validation, and same-session closure. |
| [`review-report-template.md`](review-report-template.md) | Reporting review findings, evidence, residual risk, and closure status. |
| [`product-impact-alignment-template.md`](product-impact-alignment-template.md) | Aligning product impact, journey evidence, assumptions, and acceptance. |
| [`ui-governance-evidence-template.md`](ui-governance-evidence-template.md) | Classifying UI compliance evidence by governance domain, collection method, validation mechanism, evidence class, severity, and project-local design-system policy. |

Do not merge UI governance evidence back into `contract-template.md` unless the artifact is defining expected UI behavior. Compliance evidence belongs in the focused UI governance template.
