# Ai-skill CLI Runtime Docs

本目錄保存 [`Ai-skill CLI Runtime`](../README.md) 的文件先行 artifacts。這些文件是 Phase 0 的 source-of-truth；未完成前不得開始 Go implementation。

## 何時讀哪個文件

| 文件 | 何時讀 |
| --- | --- |
| [`change-brief.md`](change-brief.md) | 開始或調整本計畫 scope、確認為什麼要做跨平台 runtime 時 |
| [`command-contract.md`](command-contract.md) | 設計或實作任何 `ai-skill` CLI command 前 |
| [`support-matrix.md`](support-matrix.md) | 判斷 Windows、macOS、Linux、iOS、Android 支援等級與限制時 |
| [`bdd-scenarios.md`](bdd-scenarios.md) | 寫測試、驗收條件或 fixture 前 |
| [`test-fixture-plan.md`](test-fixture-plan.md) | 建立測試資料、temporary repo、missing Git 或 runtime.db assertion fixture 前 |

## Phase 0 Artifact Gate

- [ ] `change-brief.md` 已確認 scope / non-goals / blocker。
- [ ] `command-contract.md` 已覆蓋所有第一批 CLI commands。
- [ ] `support-matrix.md` 已明確列出 desktop / mobile 支援邊界。
- [ ] `bdd-scenarios.md` 已覆蓋 high-risk success / failure paths。
- [ ] `test-fixture-plan.md` 已覆蓋 missing Git、unsafe repo、Windows path、fake home、runtime.db assertion。

完成上述項目前，不得新增 `scripts/ai-skill-cli/go.mod`、`scripts/ai-skill-cli/cmd/ai-skill/` 或 production Go implementation。
