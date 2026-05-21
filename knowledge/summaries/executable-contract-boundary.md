## governance.executable-contract-boundary

| 甈? | ??|
| --- | --- |
| Atom ID | `governance.executable-contract-boundary` |
| Source path | `governance/lifecycle/executable-contract-boundary.md`, `governance/lifecycle/executable-contract-boundary.yaml` |
| Lifecycle | `candidate` |
| Summary | Defines the boundary for executable YAML contracts: source stays in the owner layer, Markdown explains, YAML carries executable triggers/steps/gates/evidence, and execution-affecting contracts opt into `runtime.db` projection with `runtime_projection.enabled: true`. |
| When to read | When adding or changing governance/enforcement/workflow documents that include steps, activation, dependencies, exit gates, blocking gates, required evidence, or failure actions. |
| Do not use for | Do not move all YAML into `runtime/`; do not project ordinary metadata, graph, validation, or philosophy YAML unless it explicitly opts in. |
| Validation signal | `runtime/runtime.db.generated_surfaces` contains opt-in contracts; `knowledge-update-flow.yaml` and this boundary contract both project to runtime. |
| Last checked | 2026-05-21 |
