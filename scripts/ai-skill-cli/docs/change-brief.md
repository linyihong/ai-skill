# 變更簡報：Ai-skill CLI Runtime

> **上游計畫**：[`2026-05-21-0834-cross-platform-go-script-runtime.md`](../../../plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md)
> **對齊流程**：[`workflow/software-delivery/`](../../../workflow/software-delivery/README.md)

## 中繼資料

- **變更類型**：功能 / 工具基礎設施
- **優先序**：P0
- **證據來源**：現有 `scripts/` 依賴 shell、Ruby、Python、SQLite CLI、平台權限與本機環境。
- **日期**：2026-05-21

## 證據摘要

- `scripts/` 目前包含 shell、Ruby、Python 與 git hook，多數命令隱含 POSIX shell、Ruby、SQLite CLI、PATH、chmod / symlink 等桌面環境假設。
- `runtime/runtime.db`、knowledge runtime refresh、close-loop 與 git writeback 需要穩定驗證；不同 OS 上若缺 runtime，agent 容易只跑部分流程。
- 使用者希望以 Go 建立單一 binary，降低部署摩擦；桌面 Git 維持外部依賴，不包進 binary。

## 產品影響對齊

- **影響 / 旅程產物**：尚未拆成完整 impact map；Phase 0 先記錄核心影響。
- **決策**：先推進文件先行的 Phase 0。
- **阻擋性 mismatch**：若未完成命令契約、支援矩陣、測試 fixture 計畫，不得開始 Go 實作。

## 範圍

### 範圍內

- 定義 `ai-skill` CLI 的第一版命令契約。
- 盤點既有腳本功能與未來 CLI parity，避免新功能漏掉舊能力。
- 定義支援平台矩陣：Windows、macOS、Linux、iOS、Android。
- 定義 exit code、side effect、dry-run、JSON 輸出與失敗行為。
- 定義 BDD-lite 場景與測試 fixtures，先覆蓋缺 Git、dirty tree、merge / rebase state、runtime.db assertion。

### 範圍外

- 本文件不建立 `scripts/ai-skill-cli/go.mod`、`scripts/ai-skill-cli/cmd/ai-skill/` 或 production Go 程式碼。
- 本文件不替換尚未完成 parity 的 shell scripts。
- 本文件不承諾 iOS 任意 native binary。
- 本文件不決定 release binary 是否提交到 repo；此為 Phase 4 決策。

## 阻擋項評估

- [x] Phase 0 文件沒有阻擋項。
- [ ] 實作阻擋項：若 Phase 0 產物未完成 review，不得開始 Go 實作。

## 可追溯性

- **命令契約**：[`command-contract.md`](command-contract.md)
- **舊腳本 parity 盤點**：[`script-parity-inventory.md`](script-parity-inventory.md)
- **支援矩陣**：[`support-matrix.md`](support-matrix.md)
- **BDD 場景**：[`bdd-scenarios.md`](bdd-scenarios.md)
- **測試 fixtures**：[`test-fixture-plan.md`](test-fixture-plan.md)
