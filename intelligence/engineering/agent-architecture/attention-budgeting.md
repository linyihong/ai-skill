# Attention Budgeting（注意力預算）

**Status**: `candidate-intelligence`
**Source**: 本系統實際運作觀察

## 原則

**Agent has finite attention per session; every unnecessary file read or tool call consumes budget that could be spent on reasoning.**

Agent 每個 session 有有限的注意力；每一次不必要的檔案讀取或工具呼叫，都在消耗本可用於推理的預算。

## 為什麼

1. **注意力不是無限的** — Agent 的上下文窗口雖然大，但 attention 分佈不均。過多的輸入會稀釋對關鍵資訊的關注。
2. **每個工具呼叫都有成本** — 每次 `read_file`、`search_files`、`execute_command` 都消耗 token，而且結果需要被處理和整合。
3. **檔案讀取有機會成本** — 讀一個不相關的檔案，就少了一次讀取相關檔案的機會。
4. **Agent 不會自動優先排序** — Agent 傾向於讀取所有可讀的檔案，而不是選擇性讀取最重要的。

## 症狀

| 症狀 | 說明 | 可信度 |
|------|------|--------|
| **過度讀取** | Agent 在開始工作前讀取了 10+ 個檔案，但只使用了其中 2-3 個 | 高 |
| **重複讀取** | Agent 在同一個 session 中多次讀取同一個檔案 | 高 |
| **工具濫用** | Agent 使用 `search_files` 搜尋後又用 `read_file` 讀取每個結果 | 中 |
| **推理不足** | Agent 花費大量 token 在檔案操作上，實際推理內容很少 | 中 |

## 預防方式

1. **分層讀取** — 先讀摘要/索引，再根據需要讀取具體檔案（`knowledge/indexes/README.md` 就是為此設計）
2. **設定讀取上限** — 在 task plan 中明確限制每個階段的檔案讀取數量
3. **使用 indentation mode** — 需要讀取特定函數時使用 `indentation` 模式，而不是讀取整個檔案
4. **快取讀取結果** — 已經讀過的檔案不要重複讀取，除非內容可能已變更
5. **優先使用 search_files** — 搜尋特定內容比讀取整個檔案更節省 token

## 不建議的做法

| 不建議 | 原因 |
|--------|------|
| 在開始工作前讀取所有相關檔案 | 大部分讀取的內容不會被使用 |
| 使用 `read_file` 讀取大型檔案的全部內容 | 大型檔案（>500 行）應使用 slice 或 indentation 模式 |
| 在同一個 session 中反覆讀取同一檔案 | 除非檔案已被修改，否則內容不會改變 |

## 注意力預算分配建議

| 活動 | 建議預算佔比 | 說明 |
|------|------------|------|
| 任務理解與規劃 | 10-15% | 讀取 task plan、確定目標 |
| 檔案讀取 | 20-30% | 讀取必要的索引、規則、來源檔案 |
| 推理與決策 | 30-40% | 核心思考時間 |
| 工具執行 | 15-25% | 實際修改檔案、執行命令 |
| 驗證 | 5-10% | git status、grep 驗證、連動更新檢查 |

## 相關 atoms

- [`context-collapse.md`](context-collapse.md) — 上下文崩塌
- [`rule-overload.md`](rule-overload.md) — 規則超載
- [`task-routing.md`](task-routing.md) — 任務路由

## Token Impact

良好的注意力預算管理可以減少 30-50% 的 token 消耗。一個 session 中不必要的檔案讀取和工具呼叫是最大的 token 浪費來源。

---

← [回到 agent-architecture/](README.md)
