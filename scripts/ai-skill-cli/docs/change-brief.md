# Change Brief：Ai-skill CLI Runtime

> **上游計畫**：[`2026-05-21-0834-cross-platform-go-script-runtime.md`](../../../plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md)
> **對齊流程**：[`workflow/software-delivery/`](../../../workflow/software-delivery/README.md)

## Metadata

- **Change Type**：feature / tooling infrastructure
- **Priority**：p0
- **Evidence Source**：現有 `scripts/` 依賴 shell、Ruby、Python、SQLite CLI、平台權限與本機環境。
- **Date**：2026-05-21

## Evidence Summary

- `scripts/` 目前包含 shell、Ruby、Python 與 git hook，多數命令隱含 POSIX shell、Ruby、SQLite CLI、PATH、chmod / symlink 等桌面環境假設。
- `runtime/runtime.db`、knowledge runtime refresh、close-loop 與 git writeback 需要穩定驗證；不同 OS 上若缺 runtime，agent 容易只跑部分流程。
- 使用者希望以 Go 建立單一 binary，降低 deployment friction；desktop Git 維持 external dependency，不包進 binary。

## Product Impact Alignment

- **Impact / journey artifact**：not yet split into full impact map；Phase 0 先記錄核心 impact。
- **Decision**：proceed with docs-first Phase 0.
- **Blocking mismatch**：若未完成 command contract / support matrix / test fixture plan，不得開始 Go implementation。

## Scope

### In Scope

- 定義 `ai-skill` CLI 的第一版 command contract。
- 定義支援平台矩陣：Windows、macOS、Linux、iOS、Android。
- 定義 exit code、side effects、dry-run、JSON output 與 failure 行為。
- 定義 BDD-lite scenarios 與 test fixtures，先覆蓋 missing Git、dirty tree、merge / rebase state、runtime.db assertion。

### Out of Scope

- 本文件不建立 `scripts/ai-skill-cli/go.mod`、`scripts/ai-skill-cli/cmd/ai-skill/` 或 production Go code。
- 本文件不替換現有 Ruby compiler 或 shell scripts。
- 本文件不承諾 iOS native arbitrary binary。
- 本文件不決定 release binary 是否提交到 repo；此為 Phase 4 decision。

## Blocker Assessment

- [x] No blocker for Phase 0 documentation.
- [ ] Blocker for implementation：若 Phase 0 artifacts 未完成 review，不得開始 Go implementation。

## Traceability

- **Command contract**：[`command-contract.md`](command-contract.md)
- **Support matrix**：[`support-matrix.md`](support-matrix.md)
- **BDD scenarios**：[`bdd-scenarios.md`](bdd-scenarios.md)
- **Test fixtures**：[`test-fixture-plan.md`](test-fixture-plan.md)
