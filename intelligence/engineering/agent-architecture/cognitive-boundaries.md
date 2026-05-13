# Cognitive Boundaries（認知邊界）

**Status**: `candidate-intelligence`
**Source**: 本系統實際運作觀察

## 原則

**Agent cannot reliably detect its own cognitive boundaries; external gates (validation, checklists, failure patterns) are required.**

Agent 無法可靠地檢測自己的認知邊界；需要外部關卡（驗證、檢查清單、失效模式）來補償。

## 為什麼

1. **無內省能力** — Agent 無法準確判斷「我是否知道這個問題的答案」或「我是否遺漏了關鍵資訊」。它傾向於過度自信。
2. **確認偏誤** — Agent 傾向於尋找支持自己當前假設的證據，而不是尋找反證。這在除錯和 routing 時尤其明顯。
3. **無法感知上下文耗盡** — Agent 不會主動意識到「我的上下文已經太滿了，需要 recap 或 commit」。它會繼續運作，直到品質明顯下降。
4. **無法感知規則衝突** — Agent 不會主動檢測當前載入的規則之間是否有矛盾。它只是選擇最近看到的那條。
5. **無法評估自己的輸出品質** — Agent 對自己輸出的品質評估不可靠。它傾向於認為自己的輸出是正確的，除非被外部驗證反駁。

## 症狀

| 症狀 | 說明 | 可信度 |
|------|------|--------|
| **過度自信** | Agent 在不確定的情況下給出確定的答案，沒有標記不確定性 | 高 |
| **忽略矛盾證據** | Agent 找到與當前假設矛盾的資訊時，傾向於忽略而不是重新評估 | 高 |
| **不檢查自己的輸出** | Agent 完成修改後不驗證結果是否正確 | 高 |
| **不請求幫助** | Agent 在不確定的情況下繼續前進，而不是向 user 請求澄清 | 中 |
| **不接受限制** | Agent 在被指出錯誤時，傾向於辯解而不是接受 | 中 |

## 外部關卡設計原則

1. **驗證必須是 stateless 的** — 不依賴 agent 的記憶，只依賴檔案內容和命令輸出
2. **檢查清單必須在 task plan 中** — 不要依賴 agent 記住要檢查什麼
3. **failure pattern 必須有明確的 trigger** — Agent 需要知道「什麼時候該懷疑」
4. **關卡必須在 agent 做出不可逆決策之前** — 刪除檔案、修改 source of truth 之前要有 gate
5. **關卡必須可被自動執行** — `git status`、`grep`、diff review 等命令式檢查

## 已知的認知邊界

| 邊界 | 說明 | 外部補償 |
|------|------|---------|
| **上下文容量** | Agent 無法感知上下文何時耗盡 | `git log` recap、`.agent-goals/`、定期 commit |
| **規則優先級** | Agent 無法可靠排序規則 | `rule-weight.md`、`dependency-reading.md` 的 Bootstrap Boundary |
| **路由正確性** | Agent 無法判斷是否選了正確路徑 | `primary_entrypoint`、`forbidden_routes`、validation scenario |
| **輸出品質** | Agent 無法評估自己的輸出 | validation gate、diff review、grep 驗證 |
| **完整性** | Agent 無法判斷是否完成了所有連動更新 | `linked-updates.md` 表格、Step 7a 檢查 |
| **錯誤檢測** | Agent 無法判斷自己的推理是否正確 | `failure-patterns/`、validation scenario |

## 不建議的做法

| 不建議 | 原因 |
|--------|------|
| 依賴 agent 自我檢測錯誤 | Agent 沒有可靠的內省能力 |
| 假設 agent 知道自己的限制 | Agent 傾向於過度自信 |
| 不使用外部驗證就接受 agent 的完成宣告 | Agent 的完成宣告不可靠 |

## 相關 atoms

- [`context-collapse.md`](context-collapse.md) — 上下文崩塌
- [`rule-overload.md`](rule-overload.md) — 規則超載
- [`task-routing.md`](task-routing.md) — 任務路由
- [`failure-recovery.md`](failure-recovery.md) — 失效恢復

## Token Impact

認知邊界是 token 浪費的根本原因之一。如果 agent 能可靠地檢測自己的邊界，可以避免大量因過度自信、錯誤路由、不完整驗證而產生的無效工作。保守估計，良好的外部關卡可以減少 20-40% 的無效 token 消耗。

---

← [回到 agent-architecture/](README.md)
