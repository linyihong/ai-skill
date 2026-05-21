# Episodic Memory

`memory/episodic/` 保存**跨 session 的情境記憶**。不同於 `memory/working/`（session-local）與 `memory/summary/`（壓縮摘要），episodic memory 記錄的是「在特定情境下發生過什麼、如何應對、結果如何」，讓 agent 在類似情境再次出現時能快速 recall 過往經驗。

Episodic replay 預設只能作 weak guidance。它可以提示要檢查的風險或路線，但不得直接支撐 current truth、completion claim 或 source-of-truth override。

## 用途

- 記錄跨 session 的 recurring scenarios（例如：某類分析任務反覆出現的瓶頸）
- 保存特定情境下的成功應對模式（例如：某種錯誤模式的最佳處理路徑）
- 提供 `feedback/replay/` 的素材來源（replay 使用 episodic memory 作為分析原料）
- 支援 agent 在類似情境中快速 recall 過往經驗，避免從零開始

## 不放什麼

- Session-local 進行中狀態 → `memory/working/`
- 壓縮的 session 摘要 → `memory/summary/`
- 不可變的決策記錄 → `memory/decision/`
- 可執行的 feedback lesson → `feedback/`
- 抽象化的 intelligence atom → `intelligence/`

## 格式

```markdown
# Episodic: {情境名稱}

## Trigger Context
{什麼情境觸發了這個 episodic memory}

## What Happened
{發生了什麼事，按時間順序}

## Response
{agent 如何應對}

## Outcome
{結果如何：成功／部分成功／失敗}

## Key Signals
- {signal 1}：{這個 signal 的意義}
- {signal 2}：{這個 signal 的意義}

## Replay Candidate
{是否適合 promotion 到 feedback/replay/，以及原因}

## Linked Episodes
- {related episode 1}
- {related episode 2}
```

## 規則

1. **情境驅動**：Episodic memory 以「情境」為單位，不以「session」為單位。一個 session 可能產生多個 episodic records。
2. **可 recall**：每個 record 必須包含足夠的 trigger context，讓 agent 在未來能判斷「這個情境是否匹配」。
3. **不重複 summary**：Episodic memory 不是 session summary 的副本。它聚焦於「特定情境的經驗」，而非「整個 session 做了什麼」。
4. **Replay pipeline**：成熟的 episodic memory 應考慮 promotion 到 `feedback/replay/` 進行 generalized lesson 提取。
5. **Token-aware**：每個 episodic record 不超過 300 tokens，保留關鍵信號即可。
6. **Weak guidance default**：Replay 前必須確認情境、workflow family、domain boundary 與 current task 是否匹配。
7. **Revalidation required**：任何從 episodic memory 得到的 conclusion 都必須重新驗證 current source。
8. **No active blocker replay**：不得把過去 episode 的 blocker 當成現在仍 active。

## 與既有層的關係

- `memory/working/`：episodic memory 的原始素材來源（session archive 中的關鍵情境）
- `memory/summary/`：episodic memory 提供 summary 無法捕捉的「情境細節」
- `feedback/replay/`：episodic memory 是 replay 的分析原料
- `intelligence/engineering/failure/`：失敗類的 episodic memory 可抽象化為 failure intelligence
- `enforcement/failure-learning-system.md`：failure episodic 的 promotion 路徑
- `memory/retrieval-governance/`：定義 episodic replay 的 trigger、budget 與 contamination boundary
