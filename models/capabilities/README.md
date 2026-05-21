# Model Capability Dimensions

`models/capabilities/` 把粗略的 `small`、`large`、`specialized` profiles 細分成 capability dimensions，用來指引 execution strategy，同時避免 unsupported provider-specific claims。

## Core Dimensions

| Dimension | 意義 | 使用場景 |
| --- | --- | --- |
| [`reasoning-depth.md`](reasoning-depth.md) | 處理 tradeoffs、contradiction propagation、architecture 與 long-form planning 的能力。 | Planning、migration、recovery、model routing。 |
| [`context-stability.md`](context-stability.md) | 長 context 中維持 goals、evidence 與 instructions 的能力。 | Handoff、long tasks、contamination avoidance。 |
| [`tool-reliability.md`](tool-reliability.md) | 跨 tool calls、validation loops 與 commit / push workflows 的可靠性。 | Coding、close-loop、runtime validation。 |
| [`hallucination-risk.md`](hallucination-risk.md) | Unsupported claims、invented source paths 或 overbroad conclusions 的風險。 | Evidence hierarchy、claim scope、validation strategy。 |
| [`compression-resilience.md`](compression-resilience.md) | 從 summaries、checklists、generated reports 工作時不遺失 required gates 的能力。 | Context loading、summary-first routing、small profile work。 |

## Confidence Rules

- Capability claim 必須 scoped 到 observed behavior、tool documentation 或 validation evidence。
- Evidence missing 時，將 capability 標成 `unknown`，並使用更安全的 behavior。
- Coarse profiles 是 defaults，不是 guarantees。
- Capability dimensions 只能 tighten execution strategy；不得 override source-of-truth、validation 或 user intent。

## Capability Record

記錄 capability decision 時使用此結構：

```text
Profile:
Capability dimension:
Observed strength:
Observed limitation:
Confidence:
Execution impact:
Validation signal:
```

## Relation To Profiles

`small`、`large`、`specialized` 仍用於 context-loading defaults。Capability dimensions 用來決定任務是否需要 stricter source loading、smaller patch scope、more validation 或 handoff。
