## runtime.operations

| 欄位 | 值 |
| --- | --- |
| Atom ID | `runtime.operations` |
| Source path | `runtime/README.md` |
| Lifecycle | `validated` |
| Summary | Runtime 層負責 context routing、dynamic loading、context pruning、agent coordination 與 orchestration。包含 router（activation rules、cost budget）、context（TTL policy、prune strategy、prompt cache layout）。 |
| When to read | 需要理解 context 如何載入/卸載、activation rules 如何觸發、TTL policy 如何運作，或 prompt cache layout 如何排序時。 |
| Do not use for | 不可取代 enforcement/ 的可執行政策。不可用於單一 tool 的 hook 或 sync 細節。 |
| Context cost | ~300 tokens |
| Estimated full cost | ~1600 tokens |
| Validation signal | router/activation-rules.yaml、context/ttl-policy.yaml 與 context/prompt-cache-playbook.md 格式正確，README 可解析。 |
| Last checked | 2026-05-12 |
