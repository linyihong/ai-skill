# Requirements Cognition Summary

| 欄位 | 值 |
| --- | --- |
| Atom ID | `intelligence.requirements-cognition` |
| Source path | `intelligence/engineering/requirements/README.md` |
| Lifecycle | `candidate` |
| Summary | Requirements cognition 先用 pre-build interrogation 釐清 goal、scope、non-goals、acceptance、framework source-of-truth 與 duplication risk，再用 Impact Map × Customer Journey Map 對齊 product impact，並用 BDD-lite 處理 ambiguity、actor intent、behavior boundary、traceability、validation target 與 test effectiveness，而不是 Gherkin everywhere。 |
| When to read | 任務涉及 plan、framework 改動、新產品/新功能、product impact、customer journey、observable behavior、BDD、acceptance criteria、requirement mismatch、traceability、validation target、mutation testing 或 high-coverage-low-confidence 測試有效性時。 |
| Do not use for | 不可取代 project-local product brief、contract、test evidence 或 human alignment；不可把 BDD syntax promotion 成 runtime primitive。 |
| Validation signal | Pre-build interrogation → impact / journey evidence → requirement 或 framework discovery → behavior / owner contract → acceptance criteria → validation target → execution artifact 可追溯，沒有未確認 feature 或 duplicate source-of-truth 被寫成需求。 |
| Last checked | 2026-05-20 |
