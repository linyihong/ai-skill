# Behavior Scope Governance

**Status**: `candidate-intelligence`

## 原則

Local behavior validation 不等於 global feature correctness。

## 反模式

```text
single scenario pass
→ full workflow declared successful
```

## 規則

Agent 必須說明 scenario coverage 的 scope，並列出未覆蓋的 actor、state、error path 或 integration path。
