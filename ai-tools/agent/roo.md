# Roo Code 使用說明

本檔只記錄 Roo Code 與其他 agent 工具不同的地方。通用 bootstrap 與 obligation 來源一律是 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md) 與 [`runtime/core-bootstrap.yaml`](../../runtime/core-bootstrap.yaml)。

## Thin Entry Point

Roo Code 的入口是全域 Custom Instructions 或專案 `.roomodes`。入口內容必須保持 thin pointer，只指向：

1. [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)
2. [`runtime/core-bootstrap.yaml`](../../runtime/core-bootstrap.yaml)

不要在 Custom Instructions、`.roomodes` 或本 adapter 複製 Bootstrap Receipt、Cognitive Mode、close-loop、goal ledger 或 runtime phase 的細節；需要時從 canonical source 讀取。

## Roo Code 差異

- Roo Code 是 VS Code extension，入口可能來自全域 Custom Instructions、VS Code settings 或專案 `.roomodes`。
- `.roomodes` 的 `customInstructions` 會覆蓋全域 Custom Instructions，不會合併；若使用 `.roomodes`，需確保 thin bootstrap pointer 仍存在。
- Roo Code 支援多 modes 與 file restrictions；mode 限制不得阻擋該 mode 需要讀取的 bootstrap entry。
- 修改 Roo 全域設定若碰到 VS Code `state.vscdb`，需使用 guarded flow，避免 VS Code 執行中覆寫與 WAL 未 checkpoint。

## 配置邊界

Roo-specific 的 modes、file restrictions、Custom Instructions、VS Code extension settings 與 SQLite key 留在本檔、`.roomodes` 或工具設定。跨工具規則、語言行為、goal ledger 與 validation flow 不在 adapter 中重述。

← [回到 AI 工具索引](../README.md)
