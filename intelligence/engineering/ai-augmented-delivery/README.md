# AI-Augmented Delivery（AI 輔助開發的工程判斷）

`intelligence/engineering/ai-augmented-delivery/` 收錄當開發流程被 AI codegen（Copilot、Claude、Cursor、自動化 agent 等）大幅加速後，跨工具可重用的工程取捨與設計原理。本層只放「為什麼這樣判斷」的抽象原則；具體量化資料與觀察方法在 [`analysis/ai-augmented-delivery/`](../../../analysis/ai-augmented-delivery/README.md)；具體 detection rule 在 [`enforcement/failure-patterns/`](../../../enforcement/failure-patterns/README.md)。

## Atoms

| Atom | Status | 摘要 |
|---|---|---|
| [`generation-validation-rate-parity.md`](generation-validation-rate-parity.md) | candidate-intelligence | 任何加速「產出」的工具或流程必須同步加速「驗證」；否則 net 效益會被驗證瓶頸吃掉，且風險集中流向 production。 |

## 何時引用本層

- 在 workflow 中遇到「AI 寫得快但驗證跟不上」的決策點
- 在 review 中需要解釋「為什麼仍要堅持 perf test / canary / SLO 投資」
- 在 governance 討論中需要把 AI codegen 的風險翻譯成工程語言

## 與其他層的對應

- 量化觀察與解剖 → [`analysis/ai-augmented-delivery/`](../../../analysis/ai-augmented-delivery/README.md)
- 機械化檢測規則 → [`enforcement/failure-patterns/ai-codegen-passes-ci-fails-production.md`](../../../enforcement/failure-patterns/ai-codegen-passes-ci-fails-production.md)
- 執行流程整合 → [`workflow/software-delivery/perf-risk-gate.md`](../../../workflow/software-delivery/perf-risk-gate.md)
- Validation scenario → [`validation/scenarios/software-delivery/ai-codegen-perf-risk-checklist.yaml`](../../../validation/scenarios/software-delivery/ai-codegen-perf-risk-checklist.yaml)

← [Back to intelligence/engineering](../README.md)
