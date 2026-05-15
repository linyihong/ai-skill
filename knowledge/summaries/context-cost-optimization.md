## architecture.context-cost-optimization

| 欄位 | 值 |
| --- | --- |
| Atom ID | `architecture.context-cost-optimization` |
| Source path | `plans/archived/2026-05-12-context-cost-optimization.md` |
| Lifecycle | `validated` |
| Summary | Token 成本優化規劃。Phase 1（立即省錢）：Bootstrap 極小化（~800 tokens）、README 拆分、Rule lazy-load、Summary layer。Phase 2（架構升級）：Runtime Context Router、Context Cost Metadata、Skill Index、Context TTL。Phase 3（長期）：Semantic Retrieval、Episodic Memory、Multi-model Routing。 |
| When to read | 需要了解 token 成本優化策略、Bootstrap 拆分原則、或 lazy-load activation model 時。 |
| Do not use for | 不可取代 plans/archived/2026-05-11-next-stage-upgrade-plan.md 的完整架構分層規劃。不可用於日常 skill 操作。 |
| Context cost | ~350 tokens |
| Estimated full cost | ~4500 tokens |
| Validation signal | Plan 中的 migration steps 已部分完成（CORE_BOOTSTRAP、skills-index、activation-rules、ttl-policy）。 |
| Last checked | 2026-05-12 |
