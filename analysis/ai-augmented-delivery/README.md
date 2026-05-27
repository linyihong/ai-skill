# AI-Augmented Delivery Analysis

`analysis/ai-augmented-delivery/` 收錄當 AI codegen 工具大幅進入軟體開發流程後，可重用的觀察方法、量化資料與問題解剖。本層只放「如何觀察與拆解 AI 輔助開發的衝擊」；抽象原則在 [`intelligence/engineering/ai-augmented-delivery/`](../../intelligence/engineering/ai-augmented-delivery/README.md)；具體 detection rule 在 [`enforcement/failure-patterns/`](../../enforcement/failure-patterns/README.md)。

## 何時進入此 Domain

- 評估「採用 AI codegen 工具是否會增加 production 風險」時
- 解剖「AI 生成程式碼為什麼通過 CI 卻在 production 出錯」時
- 為 reviewer / engineering leadership 準備量化說明時
- 設計新的 perf / load / observability 投資計畫時

## 目前入口

| 文件 | 角色 |
|---|---|
| [`ai-codegen-defect-distribution.md`](ai-codegen-defect-distribution.md) | 量化資料：產出倍率、production debug 比率、redeploy 次數、開發者 debug 時間佔比 |
| [`perf-test-bottleneck-anatomy.md`](perf-test-bottleneck-anatomy.md) | 為什麼 unit / integration test 抓不到 perf 缺陷；4 個常見 anti-pattern 的 detection cost 解剖 |

## 放什麼

- AI codegen 對 delivery pipeline 各環節的衝擊觀察
- AI 生成程式碼缺陷分佈的量化資料（去敏 / 引用外部研究）
- 為什麼某類缺陷會集中在某個測試階段才被發現
- 觀察方法：如何判斷某段程式碼是 AI 生成、如何抓到 perf-sensitive 路徑

## 不放什麼

- 抽象工程原則 → [`intelligence/engineering/ai-augmented-delivery/`](../../intelligence/engineering/ai-augmented-delivery/README.md)
- 具體 detection rule / commit gate → [`enforcement/failure-patterns/`](../../enforcement/failure-patterns/README.md)
- 執行流程步驟 → [`workflow/software-delivery/`](../../workflow/software-delivery/README.md)
- Raw incident logs 或某 PR 的具體 diff（屬於業務專案 evidence）
- 某工具廠商的行銷材料

## Status

整個 domain 為 `candidate`：外部研究素材已存在，但 repo 內尚無 first-party 觀察。下次 repo 內遇到「AI 寫的 code 通過 CI 卻在 production 出問題」時，補上 first-party trail，再 promote 為 active。

## Related

- [`intelligence/engineering/ai-augmented-delivery/`](../../intelligence/engineering/ai-augmented-delivery/README.md) — 抽象原則
- [`enforcement/failure-patterns/ai-codegen-passes-ci-fails-production.md`](../../enforcement/failure-patterns/ai-codegen-passes-ci-fails-production.md) — Failure pattern
- [`workflow/software-delivery/perf-risk-gate.md`](../../workflow/software-delivery/perf-risk-gate.md) — Workflow integration
- [`analysis/production/`](../production/README.md) — Production 問題分析的通用方法（本 domain 是其特化分支）

← [Back to analysis](../README.md)
