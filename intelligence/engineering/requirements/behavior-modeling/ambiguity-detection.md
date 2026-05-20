# Ambiguity Detection

**Status**: `candidate-intelligence`

## 判斷原則

Ambiguity 是 autonomy downgrade signal，不是 agent 自行補完需求的機會。

## 高風險 ambiguity

- Acceptance criteria 不可驗證。
- Actor 或權限不明。
- 行為與既有 contract 衝突。
- Success metric 缺失但會影響 implementation。
- Business language 與 domain language 不一致。

## 行動

低風險 ambiguity 可標記 assumption；高風險 ambiguity 必須 human alignment。若 agent 自行補需求，屬於 requirement hallucination。
