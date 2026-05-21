# Validation Workflows

Validation workflow 在 evidence 不足、confidence degraded、runtime / generated surface 受影響、或 user-facing completion claim 前使用。

## Workflow Shape

1. 明確 validation target。
2. 選擇最接近 claim 的 evidence。
3. 執行 checks 或標記 not run / not applicable。
4. 比對 claim scope 與 evidence scope。
5. 若 validation local，降低 final claim 範圍。

## No Premature Success

不得用以下 evidence 做 global completion：

- Single grep result。
- Single hook success。
- Summary match。
- Memory replay。
- Local test pass for unrelated path。

## Model-aware Adjustment

- `UNCERTAIN`：validation-heavy。
- `CONTAMINATED`：rediscovery-only，再 validation。
- Low compression resilience：不要只看 summary。
- High hallucination risk：優先 source-backed evidence。
