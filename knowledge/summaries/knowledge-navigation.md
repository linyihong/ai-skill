## knowledge.navigation

| 欄位 | 值 |
| --- | --- |
| Atom ID | `knowledge.navigation` |
| Source path | `knowledge/README.md` |
| Lifecycle | `validated` |
| Summary | 知識導航系統：indexes（任務路由）、summaries（300-500 token 摘要）、graphs（知識圖譜邊）、runtime（routing registry、refresh policy、SQLite lookup cache）。讓 agent 用最小 token 成本找到正確知識。 |
| When to read | 需要從 task intent 路由到 knowledge source、查詢知識圖譜、或使用 runtime registry 時。 |
| Do not use for | 不可取代 enforcement/ 的可執行政策。不可用於修改 canonical source。 |
| Context cost | ~300 tokens |
| Estimated full cost | ~1500 tokens |
| Validation signal | indexes/summaries/graphs/runtime 子目錄 README 可解析，routing-registry.yaml 格式正確。 |
| Last checked | 2026-05-12 |
