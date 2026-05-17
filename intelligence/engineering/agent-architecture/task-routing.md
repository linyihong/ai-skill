# Task Routing（任務路由）

**Status**: `candidate-intelligence`
**Source**: 本系統實際運作觀察

## 原則

**Agent routing decisions are determined by signal strength, not signal correctness.**

Agent 的路由決策是由信號強度決定的，而不是信號正確性。

## 為什麼

1. **信號強度 > 信號正確性** — Agent 對「看起來相關」的信號反應強烈，即使該信號是誤導性的。例如看到 `app-development-guidance` 目錄名，即使任務是 APK 分析，agent 也可能被吸引過去。
2. **路徑名稱是強信號** — 一個目錄或檔案的名稱對 agent 的 routing 影響力遠大於該目錄的實際內容。
3. **最近訪問的檔案是強信號** — Agent 傾向於繼續操作最近讀取或修改過的檔案，而不是根據任務需求重新選擇。
4. **第一個找到的匹配是強信號** — Agent 通常不會搜尋所有可能選項後再做 routing 決策，而是接受第一個看似合理的匹配。

## 症狀

| 症狀 | 說明 | 可信度 |
|------|------|--------|
| **路徑漂移** | Agent 被相似名稱的目錄吸引，去了錯誤的位置 | 高 |
| **持續修改同一檔案** | Agent 不斷回到同一個檔案，即使任務已經轉移到其他領域 | 高 |
| **忽略 primary_source** | Agent 讀了 `knowledge/runtime/routing-registry.yaml` 的 `candidate_sources` 後，沒有檢查 `primary_source` | 高 |
| **工具偏好** | Agent 持續使用同一種工具（如 `apply_diff`），即使其他工具更適合當前任務 | 中 |

## 預防方式

1. **明確的 routing 信號** — 在 `knowledge/runtime/routing-registry.yaml` 中使用 `primary_source` 明確引導 agent 到正確路徑
2. **負面信號** — 在檔案中明確標記「不屬於這裡」的內容，幫助 agent 排除錯誤路由
3. **路由驗證** — 在關鍵 routing 點設置 validation gate，確認 agent 選擇了正確路徑
4. **限制搜尋範圍** — 在 task plan 中明確指定要操作的目錄，減少 agent 自由探索的空間
5. **使用 forbidden_routes** — 在 validation scenario 中定義不該走的路徑

## 不建議的做法

| 不建議 | 原因 |
|--------|------|
| 依賴 agent 自行發現正確路徑 | Agent 的 routing 是信號驅動的，不是邏輯驅動的 |
| 使用容易混淆的目錄名稱 | 相似名稱會產生錯誤的信號強度 |
| 在 routing 點提供過多選項 | 選項越多，agent 越可能選錯 |

## 相關 atoms

- [`cognitive-boundaries.md`](cognitive-boundaries.md) — 認知邊界
- [`attention-budgeting.md`](attention-budgeting.md) — 注意力預算管理
- [`context-collapse.md`](context-collapse.md) — 上下文崩塌

## Token Impact

錯誤路由是 token 浪費的主要來源之一。一次錯誤路由可能導致 10-20 次無效的 tool call，消耗 5K-10K token 後才被發現並修正。

---

← [回到 agent-architecture/](README.md)
