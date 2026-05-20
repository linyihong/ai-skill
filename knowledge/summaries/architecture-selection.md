## intelligence.architecture-selection

| 欄位 | 值 |
| --- | --- |
| Atom ID | `intelligence.architecture-selection` |
| Source path | `intelligence/engineering/architecture/architecture-selection/README.md` |
| Lifecycle | `candidate` |
| Summary | Architecture selection 先評估 business complexity、invariant density、integration pressure、lifecycle 與 team boundary，再選 CRUD、DDD Lite、Full DDD、event-driven 或 microservices。DDD 是 selectable architecture strategy，不是 runtime invariant。 |
| When to read | 使用者要求 architecture plan、design review、DDD、bounded context、CQRS、event sourcing、microservices，或 agent 準備提出超過 simple service layer 的架構時。 |
| Do not use for | 不可取代完整的 `workflow/software-delivery/architecture/`、`intelligence/engineering/domain/domain-driven-design/` 或 governance gate。不可把 DDD 當預設架構。 |
| Context cost | ~350 tokens |
| Estimated full cost | ~3000 tokens |
| Validation signal | Recommendation 必須包含 chosen strategy、rejected lighter option、rejected heavier option、fit evidence 與 upgrade/downgrade trigger。 |
| Last checked | 2026-05-20 |
