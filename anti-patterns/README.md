# Anti-patterns

`anti-patterns/` 保存 AI Runtime 中常見的 anti-patterns。這些是從實際操作中沉澱的失效模式，幫助 agent 避免重複犯錯。

## 為什麼需要 Anti-patterns

AI 在長時間運作中容易陷入特定失效模式：

| Anti-pattern | 症狀 | 後果 |
| --- | --- | --- |
| Context Explosion | Context 持續成長不 pruning | Token 爆量、品質暴跌 |
| Recursive Tool Loop | 同一 tool 反覆呼叫無新結果 | Token 浪費、無進展 |
| Hallucination Loop | Agent 在無 source 情況下持續推論 | 錯誤結論、難以回溯 |
| Stale Summary | Summary 與 source 不一致 | Agent 基於過期資訊做決策 |
| Skill Pollution | 不相關的 skill 被載入 context | Token 浪費、干擾判斷 |
| Source Mirror Drift | Mirror 與 source 不同步 | 不一致的行為 |

## 目錄

| Anti-pattern | 嚴重度 | 檔案 |
| --- | --- | --- |
| Context Explosion | critical | [`context-explosion.md`](context-explosion.md) |
| Recursive Tool Loop | critical | [`recursive-tool-loop.md`](recursive-tool-loop.md) |
| Hallucination Loop | critical | [`hallucination-loop.md`](hallucination-loop.md) |
| Stale Summary | high | [`stale-summary.md`](stale-summary.md) |
| Skill Pollution | high | [`skill-pollution.md`](skill-pollution.md) |

## 格式

每個 anti-pattern 使用以下格式：

```markdown
# {Anti-pattern Name}

## 症狀
{如何辨識}

## 根本原因
{為什麼發生}

## 影響
{對系統的負面影響}

## 預防
{如何避免}

## 檢測
{如何自動檢測}

## 恢復
{發生後如何處理}

## 相關 Guards
- runtime/guards/{相關 guard}
```

## 誰會參考這裡（Inbound References）

變更本層內容時，需要一併檢查以下依賴方：

| 來源 | 關係 |
|------|------|
| [`route.anti-patterns.runtime-patterns`](../knowledge/runtime/routing-registry.yaml) | Routing registry record，agent 依此找到 anti-patterns/ |
| [`enforcement/failure-patterns/`](../enforcement/failure-patterns/README.md) | 操作層級 failure patterns 與本層互補 |
| [`runtime/guards/`](../runtime/guards/README.md) | Runtime guards 實作本層 anti-pattern 的預防方式 |

## 與既有層的關係

- `enforcement/failure-patterns/`：現有的 failure patterns（偏重操作層級）
- `anti-patterns/`：更高層級的 runtime anti-patterns（偏重系統層級）
- `runtime/guards/`：對應的 runtime guards（自動防護）
