## intelligence.migration-feature-bundling

| 欄位 | 值 |
| --- | --- |
| Atom ID | `intelligence.migration-feature-bundling` |
| Source path | `intelligence/engineering/anti-patterns/migration-feature-bundling.md` |
| Lifecycle | `candidate` |
| Summary | 大型 migration / rewrite / platform 升級時把搬遷（refactor）與新功能（behavior change）綁進同一階段交付的反模式。後果是驗證失去 ground truth（Verification Identity Crisis）— bug 來源無法定位、時程不可預測、回滾不可行。正確路徑是 Parity-First Migration：Phase 1 達成新版 = 舊版等價並通過舊測試套件，Phase 2 才加新功能。對 stakeholder 的有效翻譯是「失望總比絕望好」。 |
| When to read | 評估 migration / rewrite / platform 升級的 plan；plan 同時含搬遷與新功能；客戶/PM 提出「順便加功能」；新版上線後出現「不確定該不該是 bug」的差異；migration 沒有獨立 parity gate。 |
| Do not use for | 不可套用於 greenfield 專案（無舊版可比對）；不可套用於明確聲明 break compatibility 的新產品設計；不可作為「禁止任何改動」的教條使用。 |
| Context cost | ~310 tokens |
| Estimated full cost | ~2400 tokens |
| Validation signal | Migration plan 含 parity gate milestone；舊版 BDD / regression 在新版 100% 通過或明列例外；差異可在 24 小時內溯源至搬遷、環境差或已知接受差異；新功能 scope-lock 到 Phase 2。 |
| Last checked | 2026-05-22 |
