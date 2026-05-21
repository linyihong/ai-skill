# Task Routing

Task routing 把工作類型映射到 execution strategy。這是 behavior contract，不是 provider model picker。

## Task Class Mapping

| Task class | 預設 strategy | Context loading | Validation target |
| --- | --- | --- | --- |
| `trivial_lookup` | `inspection-only` | `index-only` 或 `summary-first` | 回答引用 source，或明確說 not found。 |
| `bounded_doc_edit` | `source-backed` | Primary source 加上本次 edit 觸碰的 linked docs。 | Diff review、必要 lints、linked updates。 |
| `code_patch` | `execution-heavy` | Source-backed code context 加 tests / lints。 | Targeted tests，或說明未跑 tests 的原因。 |
| `architecture_planning` | `validation-heavy` | Source-backed 加 relevant architecture / governance docs。 | Plan 引用 tradeoffs、owner layers 與 blockers。 |
| `migration_promotion_deprecation` | `graph-assisted` | Primary source、dependencies、lifecycle docs、routing / metadata。 | Parity、disposition、generated surfaces 與 closure checks。 |
| `recovery_or_contradiction` | `recovery-specialized` | Source-of-truth reload 加 evidence comparison。 | Old belief downgraded、new evidence recorded、autonomy reevaluated。 |
| `long_context_handoff` | `source-backed` | Summary / goal ledger / plan，接著 reread selected source。 | Handoff 說明 assumptions、open work 與 validation gap。 |

## Escalation Rules

出現下列情況時，從較便宜的 strategy 升級到較強 strategy：

- 任務會寫 canonical source、commit、push 或更新 generated runtime artifacts。
- Evidence stale、scope 太低或被 contradiction。
- User goal、workflow success criteria 與 current action 不再一致。
- 觸碰 generated surface、runtime DB、routing registry 或 validation scenario。
- 回答會從 local evidence 推成 global claim。

## Minimality Rule

使用能安全完成 user goal 的最小 strategy。小型查詢不應啟動 architecture-level planning 或 recovery，除非出現 source-of-truth、safety 或 contradiction signal。

## Output Contract

非 trivial work 應在工作筆記或 final validation 記錄 selected strategy：

```text
Task class:
Selected strategy:
此 strategy 足夠的原因:
已檢查的 escalation trigger:
Validation target:
```
