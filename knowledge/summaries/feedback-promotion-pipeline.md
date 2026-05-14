# feedback.promotion.pipeline

| 欄位 | 值 |
| --- | --- |
| Atom ID | `feedback.promotion.pipeline` |
| Source path | [`../../feedback/promotion/README.md`](../../feedback/promotion/README.md) |
| Lifecycle | `candidate` |
| Summary | 定義 feedback lesson 從 skill-local history 推進到 workflow、intelligence、shared-rules、memory 或 runtime surfaces 的 promotion / downgrade gate。 |
| When to read | 新 lesson 需要 promotion、feedback-derived atom 要進 runtime、或需要判斷 lesson 是否該升級到 intelligence / shared-rules 時。 |
| Do not use for | 不可取代 `enforcement/feedback-lessons.md` 的命名、模板與 agent 行為規則，也不可取代原 lesson source。 |
| Validation signal | Promotion checklist 保留 old lesson source，並要求更新 index、registry、summary、graph、reports 與 close-loop validation。 |
| Last checked | 2026-05-11 |

## Checklist

- 保留原 `feedback_history/` lesson。
- 檢查 generalized / sanitized / validation gates。
- 選擇最小 durable target。
- Runtime route 需要同步 index、registry、summary、graph 與 reports。
- Promotion 或 downgrade 後完成 close-loop validation。
