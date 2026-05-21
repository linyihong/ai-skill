# Reasoning Depth

Reasoning depth 描述 execution profile 是否能支撐 multi-step tradeoffs、contradiction handling 與 architectural decisions。

## Levels

| Level | 意義 | Execution impact |
| --- | --- | --- |
| `low` | 可 follow checklist 或做 bounded formatting。 | 使用 `checklist-first`；避免 architecture claims。 |
| `medium` | 可比對少量 sources 並做 scoped decisions。 | 使用 `source-backed`；保持 patch scope narrow。 |
| `high` | 可跨 layers、tradeoffs、recovery 與 migration 推理。 | Complex work 使用 `validation-heavy` 或 `graph-assisted`。 |
| `unknown` | 任務沒有 reliable evidence。 | 預設 source-backed 並增加 validation。 |

## Escalate When

- 任務跨越 `workflow/`、`runtime/`、`metadata/`、`knowledge/` 或 `enforcement/`。
- 工作會改 durable rule、plan、migration 或 generated runtime surface。
- 出現 evidence conflict、stale source 或 user correction。

## Validation

只有 output 能保留 source paths、tradeoffs、assumptions 與 validation target，且不擴大 claim scope 時，reasoning depth 才算足夠。
