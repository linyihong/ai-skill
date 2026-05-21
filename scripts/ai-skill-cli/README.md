# Ai-skill CLI Runtime

本目錄是跨平台 `ai-skill` CLI / runtime toolchain 的開發根目錄。目標是把現有 `scripts/` 中依賴 shell、Ruby、Python 與本機環境假設的流程，逐步升級成 Windows、macOS、Linux 可執行的單一 Go binary。

## 目錄分層

| 路徑 | 用途 |
| --- | --- |
| [`docs/`](docs/README.md) | 文件先行產物：命令契約、舊腳本 parity 盤點、支援矩陣、BDD-lite 場景、fixture 計畫 |
| `cmd/ai-skill/` | Go CLI 入口；Phase 1 已建立 `doctor` skeleton |
| `internal/` | Go internal packages；目前包含 command dispatch、`doctor` checks、path normalization、JSON / plain output 與 exit code |
| `testdata/` | 未來測試 fixtures / golden 輸出（尚未建立） |

## 開發 gate

`docs/` 的 Phase 0 產物關卡已完成，Phase 1 可開始 Go implementation。後續新增命令時仍必須先對照 [`docs/command-contract.md`](docs/command-contract.md)、[`docs/script-parity-inventory.md`](docs/script-parity-inventory.md)、[`docs/bdd-scenarios.md`](docs/bdd-scenarios.md) 與 [`docs/test-fixture-plan.md`](docs/test-fixture-plan.md)。

## 開發指令

```bash
./bin/ai-skill-darwin-arm64 version --json
./bin/ai-skill-darwin-arm64 doctor --json
./bin/ai-skill-darwin-arm64 doctor --check-runtime --json
./bin/ai-skill-darwin-arm64 runtime validate --repo ../.. --json
```

Phase 1 / Phase 3 採用 [`modernc.org/sqlite`](docs/dependency-policy.md) 作為 pure Go SQLite engine；`doctor --check-runtime` 已覆蓋 in-memory 與 temporary file-backed write / query / integrity proof。Git 維持 desktop external dependency；Shell / Ruby / Python 只允許作為 wrapper-mode 過渡依賴。

## Repo-local binaries

Committed binaries live in [`bin/`](bin/) so the repo can run the CLI without a local Go install:

- `bin/ai-skill-darwin-arm64`
- `bin/ai-skill-darwin-amd64`
- `bin/ai-skill-linux-amd64`
- `bin/ai-skill-linux-arm64`
- `bin/ai-skill-windows-amd64.exe`
- `bin/BUILDINFO`
- `bin/SHA256SUMS`

Use the binary matching the host OS/architecture. Rebuild these files only after CLI source changes:

```bash
go run ./cmd/releasebuild --stable-names --version "repo-$(git rev-parse --short HEAD)" --commit "$(git rev-parse --short HEAD)" --dist bin
```

`go test ./...` verifies `bin/SHA256SUMS`, checks `bin/BUILDINFO` against the latest CLI source commit, and smoke-tests the current host binary. If CLI source changes, rebuild `bin/` before committing.

Release artifacts：`go run ./cmd/releasebuild` 會輸出 Windows amd64、macOS amd64/arm64、Linux amd64/arm64 binaries 與 `SHA256SUMS`。`ai-skill version` 支援 `-ldflags` 注入 version / commit / date。

GitHub Actions：`.github/workflows/ai-skill-cli.yml` 會在 Windows、macOS、Linux 執行 `go test ./...` 與 `doctor` smoke checks。

上游計畫：[`plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md`](../../plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md)
