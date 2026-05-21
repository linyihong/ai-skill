# Ai-skill CLI Runtime

本目錄是跨平台 `ai-skill` CLI / runtime toolchain 的開發根目錄。目標是把現有 `scripts/` 中依賴 shell、Ruby、Python 與本機環境假設的流程，逐步升級成 Windows、macOS、Linux 可執行的單一 Go binary。

## 目錄分層

| 路徑 | 用途 |
| --- | --- |
| [`docs/`](docs/README.md) | 文件先行產物：命令契約、舊腳本 parity 盤點、支援矩陣、BDD-lite 場景、fixture 計畫 |
| `cmd/ai-skill/` | 未來 Go CLI 入口（尚未建立） |
| `internal/` | 未來 Go internal packages（尚未建立） |
| `testdata/` | 未來測試 fixtures / golden 輸出（尚未建立） |

## 開發 gate

在 `docs/` 的 Phase 0 產物關卡完成並審查前，不得新增 `scripts/ai-skill-cli/go.mod`、`scripts/ai-skill-cli/cmd/ai-skill/` 或 production Go 實作。

上游計畫：[`plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md`](../../plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md)
