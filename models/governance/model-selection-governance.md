# Model Selection Governance

Model selection governance 防止 agent 把 behavior adaptation 說成實際 model switch。

## Rules

1. Main chat fixed 或 Auto 時，只能聲明 behavior-only adaptation。
2. Explicit model selector 只有在 tool docs 或 tool API 證明可控制時才可使用。
3. 使用者要求 unavailable model 時，必須回報不可用與 available options。
4. 不得 silent substitute。
5. Provider-specific names 不放在 reusable `models/`，除非只是引用 tool adapter source path。

## Required Statement

當工具能力不明或無法控制時，說明：

```text
Actual model selection: unavailable / not controlled here
Behavior adaptation: <strategy>
Validation added: <checks>
```

## Escalation

若 task 需要更強 reasoning，但不能切 model，改用 source-backed、smaller patch scope、more validation 或 user alignment。
