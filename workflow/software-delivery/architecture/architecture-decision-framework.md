# Architecture Decision Framework

## 輸入

- Change brief 或 product brief。
- Behavior / BDD / domain contract。
- 現有架構與 constraints。
- Team ownership 與 lifecycle。
- Integration boundary 與外部系統壓力。

## Decision Record 欄位

| 欄位 | 內容 |
| --- | --- |
| Context | 需求與限制。 |
| Options | 至少列出 lighter / chosen / heavier。 |
| Decision | 選擇的架構策略。 |
| Fit evidence | 支持該策略的 complexity evidence。 |
| Rejected options | 為何不用更簡單或更複雜策略。 |
| Risks | 架構成本、測試成本、migration risk。 |
| Validation | 如何驗證決策仍適合。 |
| Revisit trigger | 何時升級、降級或重審。 |

## Gate

沒有 rejected options 的 architecture decision 不完整。
