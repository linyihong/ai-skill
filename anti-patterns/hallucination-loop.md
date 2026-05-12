# Hallucination Loop

## 症狀

- Agent 在沒有 canonical source 的情況下持續推論
- 使用 summary 代替 source 做精確判斷
- 引用過期或不存在的文件
- 對不確定的資訊做出確定性陳述

## 根本原因

1. **Summary-source 不同步**：Summary 與 source 不一致。
2. **Source 未讀**：Agent 沒讀 canonical source 就做結論。
3. **過度推論**：在資訊不足時「填空」。
4. **Tool output 截斷**：Output 被截斷但 agent 仍做判斷。

## 影響

- 錯誤的架構決策
- 難以回溯的錯誤
- 使用者對系統失去信任

## 預防

1. 重要判斷前先讀 canonical source（`dependency-reading.md`）。
2. Summary 只做 routing，不做精確判斷。
3. Tool output 截斷時標記不完整，不做確定性結論。
4. 使用 `runtime/guards/circuit-breaker.yaml` 的 hallucination risk detection。

## 檢測

- `runtime/guards/circuit-breaker.yaml` 的 `hallucination_risk` 維度
- Agent 使用「可能」、「應該」等不確定詞彙的頻率
- Source 引用但未實際讀取

## 恢復

1. 標記可疑結論。
2. 讀取相關 canonical source。
3. 修正或撤銷錯誤結論。

## 相關 Guards

- `runtime/guards/circuit-breaker.yaml`
- `shared-rules/dependency-reading.md`
