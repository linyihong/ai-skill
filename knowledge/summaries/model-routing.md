## models.routing

| 欄位 | 值 |
| --- | --- |
| Atom ID | `models.routing` |
| Source path | `models/README.md` |
| Lifecycle | `candidate` |
| Summary | 多模型協作架構：capability profile（small/large/specialized）、compression strategy（checklist/compressed/full）、model-aware context report。根據 task 複雜度選擇模型與 context 策略。 |
| When to read | 需要選擇模型、決定 compression level、或產生 model-aware context report 時。 |
| Do not use for | 不可取代單一模型的官方文件。不可用於 tool-specific 的模型設定。 |
| Context cost | ~250 tokens |
| Estimated full cost | ~1200 tokens |
| Validation signal | profiles/ 與 compression/ README 可解析，model-context-report.md 格式正確。 |
| Last checked | 2026-05-12 |
