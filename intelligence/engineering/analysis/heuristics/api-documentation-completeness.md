# API Documentation Completeness Heuristic（API 文件完整性啟發式）

> **Cross-Domain Promotion**: Field Confidence 判斷表（confirmed/candidate/needs capture/meaning unknown 等 8 種狀態）已提取到 [`intelligence/engineering/heuristics/field-confidence-judgment.md`](../../heuristics/field-confidence-judgment.md)，適用於任何需要記錄觀測結果的任務。

## 問題

HTTP API 分析中，何時該開始建立正式 API Catalog？何時可以回報 API-list 任務完成？如何判斷一個 API field 的信心水準？

## 原則

- API Catalog 的價值在於可被後續流程（BDD、fixture、contract test）直接使用
- 不完整的 API 文件比沒有文件更危險（會產生 false confidence）
- 每個 API field 有三個維度：**觀測到**、**解碼/理解語意**、**可 replay/驗證**
- 高價值 API（auth、payment、核心資料流）需要比低價值 API 更高的完成標準
- UI automation 是穩定 capture 的手段，不是文件的目的

## 決策表

### 何時開始建立 API Catalog

| 情境 | 建議做法 | 判斷信號 |
|------|---------|---------|
| 已觀測到 3+ 個 API endpoint | 開始建立 group index + per-API detail skeleton | 有 MITM/pcap/hook 輸出 |
| 已解碼 response wrapper 結構 | 建立 API entry（wrapper/decode rules） | 可辨識 `status`/`code`/`data`/`error` 欄位 |
| 已有 replay 能力 | 補齊 validation 欄位 | 可用 curl/Postman 重放請求 |
| 目標是重建 feature | 立即開始 API Catalog，即使只有 skeleton | task goal 包含 implementation |

### 何時回報完成

| 情境 | 完成標準 | 不完成就回報的風險 |
|------|---------|------------------|
| API-list 任務 | 所有觀測到的 API 在 group index 或 coverage/gap file 中 | 遺漏 endpoint 未被追蹤 |
| 高價值 API | 有 per-operation detail（不僅 method/path） | 實作時需要重新 capture |
| 每個 detail | 包含 request fields、response fields、field meaning、evidence、validation、open questions | 後續流程無法判斷欄位語意 |
| Shared behavior | headers、wrapper、auth 已文件化一次 | 每個 API detail 重複相同資訊 |
| 未驗證 API | 明確標記 confidence 等級 | 誤將 candidate 當作 confirmed |

### Field Confidence 判斷

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

- 不要在所有 API 都完成 per-API detail 前就回報完成
- 不要將 screenshots 視為 HTTP header/request/response field analysis 的替代品
- 不要對未 replay 的 API 聲稱「已驗證」
- 不要在不確定時 invent final product language；使用 `candidate` 標記

## 相關 atoms

- `intelligence/engineering/analysis/heuristics/hook-selection.md`
- `analysis/apk/workflows/http-api-documentation-flow.md`

## Token 影響

低。此 atom 在 HTTP API 分析 session 中 lazy-load，約 150-200 tokens。
