# Ai-skill CLI Runtime 文件

本目錄保存 [`Ai-skill CLI Runtime`](../README.md) 的文件先行產物。這些文件是 Phase 0 的 source-of-truth；未完成前不得開始 Go 實作。

## 何時讀哪個文件

| 文件 | 何時讀 |
| --- | --- |
| [`change-brief.md`](change-brief.md) | 開始或調整本計畫 scope、確認為什麼要做跨平台 runtime 時 |
| [`command-contract.md`](command-contract.md) | 設計或實作任何 `ai-skill` CLI 命令前 |
| [`dependency-policy.md`](dependency-policy.md) | 新增 Go dependency、外部 binary、adapter mode 或 SQLite 行為前 |
| [`legacy-to-go-migration-map.md`](legacy-to-go-migration-map.md) | 移植或刪除 legacy surface 前，確認舊入口搬到哪個 Go command/package、source-of-truth 與驗證證據 |
| [`release-distribution.md`](release-distribution.md) | 判斷 repo-local binaries、CI artifacts、upgrade / rollback 與 package manager 策略時 |
| [`runtime-native-rewrite-assessment.md`](runtime-native-rewrite-assessment.md) | 查 runtime tooling native Go rewrite 與刪除舊 Ruby surfaces 的驗收依據時 |
| [`script-parity-inventory.md`](script-parity-inventory.md) | 檢查新 CLI 是否完整涵蓋舊腳本功能、side effects 與測試證據時 |
| [`legacy-script-disposition.md`](legacy-script-disposition.md) | 決定舊 shell / Ruby / Python scripts 的保留、刪除或 adapter 例外時 |
| [`support-matrix.md`](support-matrix.md) | 判斷 Windows、macOS、Linux、iOS、Android 支援等級與限制時 |
| [`mobile-support-decision.md`](mobile-support-decision.md) | 回答 iOS / Android 是否支援本機 binary、WASM、App-contained 或 remote runner 時 |
| [`bdd-scenarios.md`](bdd-scenarios.md) | 寫測試、驗收條件或 fixture 前 |
| [`test-fixture-plan.md`](test-fixture-plan.md) | 建立測試資料、temporary repo、缺 Git 或 runtime.db assertion fixture 前 |

## Phase 0 產物關卡

- [x] `change-brief.md` 已確認範圍、非目標與阻擋項。
- [x] `command-contract.md` 已覆蓋所有第一批 CLI 命令。
- [x] `dependency-policy.md` 已定義 pure Go 優先、Git external dependency、SQLite 選型與 adapter mode 限制。
- [x] `legacy-to-go-migration-map.md` 已補上 legacy runtime surfaces 到 Go CLI command/package 的開發者對照。
- [x] `script-parity-inventory.md` 已盤點舊腳本、未來命令、parity 狀態與最低測試證據。
- [x] `support-matrix.md` 已明確列出桌面與行動平台支援邊界。
- [x] `release-distribution.md` 已記錄 repo-local binaries、CI artifacts、upgrade / rollback 與 deferred package managers。
- [x] `mobile-support-decision.md` 已記錄 iOS / Android support decision。
- [x] `legacy-script-disposition.md` 已記錄 legacy scripts 的保留 / 刪除條件。
- [x] `bdd-scenarios.md` 已覆蓋高風險成功與失敗路徑。
- [x] `test-fixture-plan.md` 已覆蓋缺 Git、不安全 repo、Windows path、fake home、runtime.db assertion。
- [x] 文件標題、表格欄位、說明段落與 checklist 已使用繁體中文；只保留命令、flag、JSON 欄位、路徑與固定術語英文。

上述 Phase 0 gate 已完成；後續新增或修改 Go implementation 時，仍必須先對照本目錄的 command contract、parity inventory、BDD scenarios 與 fixture plan。
