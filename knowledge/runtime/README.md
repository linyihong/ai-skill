# Knowledge Runtime Surfaces

`knowledge/runtime/` defines the future runtime-facing view of knowledge indexes, summaries, graphs, and metadata. During the current phase, this directory defines format and boundaries only; no automated runtime is implemented here.

## Runtime Inputs

| Input | Source |
| --- | --- |
| Task intent routing | `knowledge/indexes/README.md` |
| Atom metadata | `metadata/schema.md` |
| Ranking rules | `metadata/ranking/README.md` |
| Confidence rules | `metadata/confidence/README.md` |
| Compatibility rules | `metadata/compatibility/README.md` |
| Lifecycle and validation gates | `governance/lifecycle/README.md`, `governance/validation/README.md` |
| Runtime routing design | `runtime/routing/README.md` |

## Runtime View Format

Future runtime views should answer:

| Field | Purpose |
| --- | --- |
| `task_intent` | What the agent is trying to do. |
| `primary_source` | First canonical source to read. |
| `required_dependencies` | Required shared rules, skill entries, or metadata. |
| `candidate_sources` | Optional maps, summaries, or atoms. |
| `source_of_truth_gate` | Whether old entrypoint still wins. |
| `ranking_reason` | Why this source is first. |
| `validation_signal` | How to confirm the route is safe. |

## Runtime Rules

- Runtime views cannot skip required shared-rule bootstrap.
- Runtime views cannot replace old skill behavior unless lifecycle promotion gates pass.
- Runtime views should prefer low-cost summaries only when they link to canonical sources.
- Runtime views must treat tool mirrors as deployment surfaces, not source paths.
- Runtime views should record deferred sources when context is intentionally not loaded.

## Not Implemented Yet

- Automatic graph construction.
- Generated summaries.
- Machine-readable routing registry.
- Model-aware compression output.

These remain future work after the governance, metadata, and routing surfaces stabilize.
