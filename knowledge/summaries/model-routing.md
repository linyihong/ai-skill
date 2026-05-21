## models.routing

| 欄位 | 值 |
| --- | --- |
| Atom ID | `models.routing` |
| Source path | `models/README.md` |
| Lifecycle | `candidate` |
| Summary | Model-aware execution strategy：profiles、capabilities、routing、workflow adaptation、governance、runtime primitives 與 compression。用於選擇 behavior shape，不宣稱 provider model 已切換。 |
| When to read | 需要決定 task strategy、compression level、model fallback、subagent handoff、workflow adaptation 或 model-aware validation 時。 |
| Do not use for | 不可取代 provider / tool 官方模型文件。不可把 behavior-only adaptation 說成 actual model selection。不可覆蓋 workflow source-of-truth。 |
| Context cost | ~350 tokens |
| Estimated full cost | ~2200 tokens |
| Validation signal | profiles、routing、capabilities、workflow-adaptation、governance、runtime README 可解析；validation/scenarios/models 覆蓋主要 routing cases。 |
| Last checked | 2026-05-21 |
