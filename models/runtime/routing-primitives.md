# Routing Primitives

Model runtime primitives 是 lookup-friendly 的 execution strategy labels，不是 provider model identifiers。

## Primitive Set

| Primitive | Use |
| --- | --- |
| `checklist_first` | Low-risk bounded tasks, small context. |
| `source_backed` | Canonical source edits, rules, docs, code patches. |
| `validation_heavy` | Uncertain evidence, claim scope risk, degraded confidence. |
| `rediscovery_only` | Contaminated or stale execution frame. |
| `behavior_only_fallback` | Actual model selection unavailable. |
| `human_alignment` | Missing authority or user decision required. |

## Promotion Boundary

新增 primitive 前必須有 validation scenario 或 recurring routing need。不要為單次偏好新增 runtime primitive。
