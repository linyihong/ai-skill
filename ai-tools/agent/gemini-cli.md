# Gemini CLI 使用說明

本檔只記錄 Gemini CLI 與其他 agent 工具不同的地方。通用 bootstrap 與 obligation 來源一律是 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md) 與 [`runtime/core-bootstrap.yaml`](../../runtime/core-bootstrap.yaml)。

## Thin Entry Point

Gemini CLI 的 repo 入口是 root `GEMINI.md`。它必須保持 thin pointer，只指向：

1. [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)
2. [`runtime/core-bootstrap.yaml`](../../runtime/core-bootstrap.yaml)

不要在 `GEMINI.md` 或本 adapter 複製 Bootstrap Receipt、Cognitive Mode、close-loop、goal ledger 或 runtime phase 的細節；需要時從 canonical source 讀取。

## Gemini CLI 差異

- `GEMINI.md` 可有全域與專案層級；專案入口應優先保持薄，只補專案特定 overlay。
- `.gemini/` 主要是 Gemini CLI 管理的工具狀態與設定；除非任務就是維護 Gemini 設定，通常不需要手動修改。
- Gemini CLI 可用外部檢索與子代理；這些能力產出的內容仍要依共用 evidence、sanitization 與 validation 規則處理。
- 執行會修改檔案或狀態的 shell command 前，先用一句話說明目的，並依工具內建確認機制處理。

## 配置邊界

Gemini-specific 入口、工具能力與設定位置留在本檔或 `GEMINI.md`。跨工具行為與輸出格式不在 adapter 中重述。

← [回到 AI 工具索引](../README.md)
