# Tools

`tools/` 存放 AI tool 的 metadata、routing 與 compression 策略，幫助 runtime 以最小 token 成本使用工具。

## 用途

- 定義每個 AI tool 的 metadata schema（成本、風險、適用情境）
- 實作 tool lazy activation 與 explosion detection
- 定義 tool output compression 策略

## 目錄結構

```
tools/
  compression/          # Tool output compression 策略
  metadata/             # Tool metadata schema 與成本資訊
  routing/              # Tool lazy activation 與 explosion detection
```

## 子目錄說明

| 目錄 | 說明 |
|------|------|
| [`compression/`](compression/README.md) | 定義各類型 tool output 的壓縮策略（stack trace、JSON、Git diff、log、search results、file content） |
| [`metadata/`](metadata/README.md) | 定義 tool metadata schema，記錄每個 tool 的 input/output token 成本、遞迴風險、activation 策略 |
| [`routing/`](routing/README.md) | 實作 tool lazy activation 流程與 tool explosion detection 信號 |

## 與既有層的關係

- [`ai-tools/`](../ai-tools/README.md)：工具配置與同步細節（Claude Code、Cursor 等）
- [`runtime/`](../runtime/README.md)：runtime context routing 與 guard chain
- [`models/compression/`](../models/compression/README.md)：model-aware compression strategy
