# Field Confidence Judgment Heuristic（欄位信心判斷啟發式）

**Status**: `candidate-intelligence`
**Source**: 從 `intelligence/engineering/analysis/heuristics/api-documentation-completeness.md` 提取的跨領域通用部分

## 問題

在 reverse engineering、API 文件化、或任何需要記錄觀測結果的任務中，如何標記每個欄位的信心水準？何時可以聲稱「已驗證」？

## 原則

- 不完整的文件比沒有文件更危險（會產生 false confidence）
- 每個欄位有三個維度：**觀測到**、**解碼/理解語意**、**可 replay/驗證**
- 高價值欄位（auth、payment、核心資料流）需要比低價值欄位更高的完成標準
- 在不確定時使用明確標記，不要 invent 不存在的資訊

## Field Confidence 判斷表

| 狀態 | 定義 | 標記 |
|------|------|------|
| confirmed | 已觀測 + 已解碼 + 可 replay | 無需標記 |
| candidate | 已觀測但語意不確定 | `candidate` |
| needs capture | 尚未觀測到 | `needs capture` |
| needs replay | 已觀測但無法重放 | `needs replay` |
| meaning unknown | 已觀測到值但無法理解用途 | `meaning unknown` |
| low confidence | 僅從 context 推測 | `low confidence` |
| out of scope | 決定不處理 | `out of scope` |
| not observed | 靜態分析推測存在但未觀測到 | `not observed` |

## 不建議的做法

- 不要對未 replay 的欄位聲稱「已驗證」
- 不要在不確定時 invent 欄位語意；使用 `candidate` 或 `meaning unknown` 標記
- 不要將 screenshots 視為欄位分析的替代品

## 相關 atoms

- `intelligence/engineering/analysis/heuristics/api-documentation-completeness.md`（領域特定版本，包含何時開始/完成 API Catalog 的決策表）

## Token 影響

低。此 atom 在任何需要記錄觀測結果的任務中 lazy-load，約 100-150 tokens。
