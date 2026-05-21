# Fallback Routing

Fallback routing 用於 preferred execution strategy 無法對應到 explicit model 或 tool capability 的情況。

## Tool Capability Outcomes

| Capability outcome | 必要行為 |
| --- | --- |
| Main chat 是 fixed 或 Auto | 使用 behavior-only adaptation。不得宣稱 model changed。 |
| Tool exposes selector 但 agent 無法控制 | 推薦 model class，必要時請使用者切換。 |
| Subagent supports explicit model | 只有 task complexity 與 single-owner boundaries 合理時才使用 subagent。 |
| Requested model unavailable | 回報 unavailable model 並列出 available options；不得 silent substitute。 |
| Tool capability unknown | 在 tool docs 證明前，當成 behavior-only adaptation。 |

## Behavior-only Adaptation

Actual model selection 不可用時，改調整 execution shape：

- 對 uncertain 或 high-risk tasks 增加 source-backed context。
- 對 code 或 rules 使用更小 patch scope 與更多 validation。
- Simple tasks 優先使用 checklist-first execution。
- Ambiguity 或 missing authority 時使用 human-facing summary。
- Contaminated state 使用 rediscovery-only behavior。

## Silent Substitution Ban

若使用者要求某個 model 或 model class，但它不可用，agent 必須說明。Fallback strategy 只能在說明限制並選擇最安全的 behavior-only route 後繼續。

## Validation

Fallback 只有在 response 或 handoff 記錄下列項目時才有效：

```text
Requested capability:
Available capability:
Fallback strategy:
新增風險:
新增 validation:
```
