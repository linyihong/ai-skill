# Model Confidence Governance

Model confidence governance 管理 capability confidence、tool capability confidence 與 execution confidence。它不評分品牌模型，只治理本次 execution 是否有足夠 evidence。

## Confidence Inputs

- Current source read。
- Tool documentation / tool output。
- Tests、lints、runtime validation。
- Prior validated behavior。
- User correction 或 contradiction。

## Downgrade Rules

| Signal | Action |
| --- | --- |
| Capability unknown | 使用 safer behavior。 |
| Tool selector unknown | behavior-only adaptation。 |
| Evidence stale | reread source。 |
| Repeated patch fail | downgrade autonomy。 |
| Local evidence only | narrow claim scope。 |

## Completion Claim

Final claim 必須說明 validation evidence。若只完成 docs 或 local checks，不得宣稱 full runtime behavior verified。
