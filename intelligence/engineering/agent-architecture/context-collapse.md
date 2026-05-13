# Context Collapse（上下文崩塌）

**Status**: `candidate-intelligence`
**Source**: 本系統實際運作觀察

## 原則

**When context exceeds effective window, agent loses earlier decisions and repeats or contradicts them.**

當上下文超過有效窗口時，agent 會遺失早期決策，並重複或矛盾它們。

## 為什麼

1. **Agent 的注意力是滑動窗口** — 新內容推入時，舊內容被擠出。不是所有 LLM 都有無限上下文，即使有，attention 分佈也會偏向近期內容。
2. **早期決策最脆弱** — 在 session 開始時做的決策（如「這次只改 shared-rules，不碰 workflow」）在 50 次 tool call 後容易被遺忘。
3. **崩塌是漸進的** — 不是突然無法運作，而是先遺失次要細節，再遺失主要決策，最後產生矛盾行為。

## 症狀

| 症狀 | 說明 | 可信度 |
|------|------|--------|
| **重複已做過的事** | Agent 再次執行已經完成的步驟，如再次提取已提取過的檔案 | 高 |
| **矛盾決策** | Agent 先說「不刪除舊檔案」，幾輪後又說「刪除舊檔案」 | 高 |
| **遺失 task plan** | Agent 忘記自己規劃的步驟順序，開始自由發揮 | 中 |
| **忽略 user 早期指示** | User 在 session 開始時說的「不要動 A」，後續被 agent 修改 A | 高 |
| **工具呼叫模式改變** | 從精準的單一工具呼叫變成大量試錯 | 中 |

## 預防方式

1. **外部化決策** — 重要決策不要只存在 agent 記憶中，要寫入檔案（`.agent-goals/`、`plans/`、`Document TODO`）
2. **定期 recap** — 每完成一個階段，用 `git log` 或 summary 重新確認當前狀態
3. **縮小 session scope** — 一個 session 只做一件事，做完就 commit 並結束
4. **使用 validation gate** — 在關鍵決策點設置外部檢查（如 `git status`、grep 驗證）
5. **避免 long session** — 超過 30-40 次 tool call 的 session 應考慮拆分

## 不建議的做法

| 不建議 | 原因 |
|--------|------|
| 依賴 agent 記住所有上下文 | Agent 的記憶不可靠，尤其是跨多輪 tool call |
| 在同一個 session 做太多不相關的事 | 增加上下文混雜的風險 |
| 不 recap 就繼續下一階段 | 遺失早期決策的細節 |

## 相關 atoms

- [`attention-budgeting.md`](attention-budgeting.md) — 注意力預算管理
- [`cognitive-boundaries.md`](cognitive-boundaries.md) — 認知邊界
- [`failure-recovery.md`](failure-recovery.md) — 失效恢復

## Token Impact

Context collapse 是 token 浪費的最大來源之一。一個崩塌的 agent 可能花費 50-100 次 tool call 做重複或矛盾的工作，消耗數萬 token 而無實際進展。

---

← [回到 agent-architecture/](README.md)
