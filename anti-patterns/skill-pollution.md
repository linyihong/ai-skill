# Skill Pollution

## 症狀

- 不相關的 skill 被載入 context
- Agent 花 token 讀取與目前 task 無關的知識
- Context 被多個 skill 的 entrypoint 與 workflow 佔據

## 根本原因

1. **無 skill routing**：所有 skill 都被視為「可能相關」。
2. **無 relevance scoring**：沒有機制判斷 skill 與 task 的相關性。
3. **過度載入**：為了「以防萬一」載入過多 skill。

## 影響

- Token 浪費
- Context 被不相關知識污染
- Agent 可能被不相關的規則干擾

## 預防

1. 使用 `knowledge/runtime/routing-registry.yaml` 的 task_intent 做精準 routing。
2. 只載入 task intent 匹配的 route。
3. 使用 summary-first routing：先讀 summary（300-500 tokens），需要才展開。
4. 定期 prune 不相關的 context。

## 檢測

- Context 中 skill 數量 > 3 且與 task intent 無關
- Skill entrypoint 已讀但從未被使用
- Token 使用率與工作進度不成比例

## 恢復

1. Prune 不相關的 skill context。
2. 重新確認 task intent 與 routing-registry.yaml 的匹配。
3. 必要時建立新 session。

## 相關 Guards

- `knowledge/runtime/routing-registry.yaml`
- `runtime/context/ttl-policy.yaml`
- `runtime/health/context-health-score.yaml`
