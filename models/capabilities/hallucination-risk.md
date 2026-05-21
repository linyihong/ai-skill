# Hallucination Risk

Hallucination risk 描述信任 model output、path、claim 或 routing decision 前需要多少 validation。

## Risk Classes

| Class | 意義 | 必要行為 |
| --- | --- | --- |
| `low` | Claim 由 current source 或 tool output 直接支持。 | 正常 cite 或 validate。 |
| `medium` | Claim 從 partial source 推論而來。 | 行動前縮小 scope 並 validate。 |
| `high` | Claim 影響 durable rules、runtime、security、data 或 user-facing completion。 | 使用 source-backed validation 與 evidence hierarchy。 |
| `unknown` | Evidence 缺失或 stale。 | 視為 unsupported，不當成 fact 行動。 |

## High-risk Patterns

- Invented file path、command、model capability 或 API behavior。
- 用 local evidence 當作 global success。
- Tool log 覆蓋 live observation 或 owner contract。
- Summary 被當作 source-of-truth replacement。
- 只從 provider model name 假設 capability。

## Required Response

Risk 為 medium 或 higher 時：

1. 讀取 relevant source。
2. 說明 claim scope。
3. 執行或命名 validation。
4. 標記 uncovered areas。
5. Evidence local 時 downgrade confidence。
