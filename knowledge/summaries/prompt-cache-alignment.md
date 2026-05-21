## runtime.prompt-cache-alignment

| 欄位 | 值 |
| --- | --- |
| Atom ID | `runtime.prompt-cache-alignment` |
| Source path | `runtime/context/prompt-cache-playbook.md` |
| Lifecycle | `candidate` |
| Summary | Provider prompt cache 對齊規範。定義 stable prefix、semi-stable middle、volatile suffix 的 context layout，並說明 `cacheable` 與 `provider_cache_candidate` 的邊界。 |
| When to read | 討論或修改 prompt cache、prefix cache、cached tokens、context loading、bootstrap layout 或 token cost optimization 時。 |
| Do not use for | 不可取代 required dependency reading、source-of-truth validation、safety rules 或工具/供應商專屬設定文件。 |
| Context cost | ~500 tokens |
| Estimated full cost | ~1200 tokens |
| Validation signal | `enforcement/prompt-cache-efficiency.md`、`metadata/schema.md`、`runtime/runtime.db` 與 routing registry 均指向同一 layout 規範。 |
| Last checked | 2026-05-19 |
