## intelligence.migration-seeder-anti-patterns

| 欄位 | 值 |
| --- | --- |
| Atom ID | `intelligence.migration-seeder-anti-patterns` |
| Source path | `intelligence/engineering/anti-patterns/migration-seeder-anti-patterns.md` |
| Lifecycle | `candidate` |
| Summary | 把大量業務資料（廠商目錄、商品/SKU、遊戲清單、權限矩陣）以巨型 `INSERT` 包進 schema migration，使資料 lifecycle 與 schema lifecycle 被強制綁定。訊號：單檔 >50KB、檔名含 dataSeeder、業務人員想改資料要工程師寫 migration。替代方案依資料性質：列舉留 migration、目錄走 application seeder/admin、大量參考資料用外部 CSV + bulk loader。 |
| When to read | 看到任一 schema migration 檔 >50KB；migration 檔名含 dataSeeder/seed/fixture；同張表反覆出現補資料 migration；業務目錄資料想新增需要 PR + deploy；不同環境需要不同資料但用同一 migration 強制執行。 |
| Do not use for | 不可取代具體 schema migration 工具的官方文件；不可作為已存在 seeder 的 rollback 操作手冊。 |
| Context cost | ~280 tokens |
| Estimated full cost | ~1800 tokens |
| Validation signal | `find <migrations_dir> -name "*.sql" -size +50k` 應為 0 或可解釋；migration 不同時做 schema change + 大量 data insert；業務目錄有文件化的 source-of-truth 與修改路徑。 |
| Last checked | 2026-05-21 |
