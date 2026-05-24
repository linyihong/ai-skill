## intelligence.plan-first-decision-promotion

| 欄位 | 值 |
| --- | --- |
| Atom ID | `intelligence.plan-first-decision-promotion` |
| Source path | `intelligence/engineering/architecture/plan-first-decision-promotion.md` |
| Lifecycle | `candidate` |
| Summary | 架構決策的提案、討論、alternatives 評估在 `plans/active/<plan>.md` §Decision Rationale 完成；只有 plan completed 且通過 ADR Promotion Criteria 後才升級為 accepted ADR 寫入 constitution。憲法層只放已驗證決策，不放 proposed/draft。避免「廢棄憲法」累積、狀態爆炸、平行維護成本。Plan 模板強制 §Decision Rationale 6 子章節（Problem & Why Now / Decision / Alternatives Considered / Why Not an ADR Yet / ADR Promotion Criteria / Consequences）。 |
| When to read | 想為新架構/流程/跨層改動寫 ADR；想標 proposed/draft ADR；已有對應 plan 但仍想平行建 ADR；看到 constitution 內有非 accepted status；使用者問「這個要不要寫 ADR」。 |
| Do not use for | 不適用於業務專案 ADR 治理（各專案自行評估）；不適用於純 spike/exploratory 開發（無 plan 結構）；不取代既有 ADR-007 promotion target 路由。 |
| Context cost | ~320 tokens |
| Estimated full cost | ~2500 tokens |
| Validation signal | constitution/ 內所有 ADR Status 為 accepted/deprecated/superseded；架構/跨層 plan 含完整 §Decision Rationale 6 子章節；新 ADR 對應 plan 已 completed 且 commit message 引用為 evidence；plan completed 不全部升級為 ADR（依內容路由至更輕 layer）。 |
| Last checked | 2026-05-22 |
