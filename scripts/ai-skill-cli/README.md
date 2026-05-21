# Ai-skill CLI Runtime

本目錄是跨平台 `ai-skill` CLI / runtime toolchain 的開發根目錄。目標是把現有 `scripts/` 中依賴 shell、Ruby、Python 與本機環境假設的流程，逐步升級成 Windows、macOS、Linux 可執行的單一 Go binary。

## 目錄分層

| 路徑 | 用途 |
| --- | --- |
| [`docs/`](docs/README.md) | 文件先行產物：命令契約、舊腳本 parity 盤點、支援矩陣、BDD-lite 場景、fixture 計畫 |
| `cmd/ai-skill/` | Go CLI 入口；Phase 1 已建立 `doctor` skeleton |
| `internal/` | Go internal packages；目前包含 command dispatch、`doctor` checks、JSON / plain output 與 exit code |
| `testdata/` | 未來測試 fixtures / golden 輸出（尚未建立） |

## 開發 gate

`docs/` 的 Phase 0 產物關卡已完成，Phase 1 可開始 Go implementation。後續新增命令時仍必須先對照 [`docs/command-contract.md`](docs/command-contract.md)、[`docs/script-parity-inventory.md`](docs/script-parity-inventory.md)、[`docs/bdd-scenarios.md`](docs/bdd-scenarios.md) 與 [`docs/test-fixture-plan.md`](docs/test-fixture-plan.md)。

## 開發指令

```bash
go test ./...
go run ./cmd/ai-skill doctor --json
go run ./cmd/ai-skill doctor --require-git --plain
```

Phase 1 目前不引入外部 Go dependency。Git 維持 desktop external dependency；SQLite library 選型仍以 pure Go SQLite 為目標，待後續 Phase 1 子項完成。

上游計畫：[`plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md`](../../plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md)
