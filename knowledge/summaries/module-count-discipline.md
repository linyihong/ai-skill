## intelligence.module-count-discipline

| 欄位 | 值 |
| --- | --- |
| Atom ID | `intelligence.module-count-discipline` |
| Source path | `intelligence/engineering/architecture/modularity/module-count-discipline.md` |
| Lifecycle | `candidate` |
| Summary | Repo 內 module（Maven / Gradle / npm workspace / Cargo / Go multi-module / Bazel / .NET / pnpm / Lerna / Nx / poetry workspace）數量是工程成本，N 失控時 build time、IDE 載入、升級擴散、refactor 成本線性以上增長。N 對應健康做法：≤5 直接寫、5-30 模板+門檻、30-100 plugin/SPI/catalog、≥100 結構味道警報。新增 module 必須通過 lifecycle / classpath / external use / N 警戒 / data-vs-code 五道判斷。 |
| When to read | Repo 模組數 N ≥ 30 且仍增長；CI build > 5 分鐘或本機增量 > 30 秒；IDE 載入卡頓；升級共用相依需動 ≥ 5 個模組；看到以「廠商 / SKU / 商品 / 遊戲 / 地區」為單位的編譯期模組；同 repo 內模組整合風格不一致。 |
| Do not use for | 不可作為「禁止新增 module」的教條使用；不可套用於短期 PoC / demo；不可取代特定 build system 的官方 documentation。 |
| Context cost | ~320 tokens |
| Estimated full cost | ~2900 tokens |
| Validation signal | Repo 有書面 module 新增門檻；全量 CI < 5 分鐘 或有改善 plan；升級擴散範圍 ≤ N × 10%；整合風格 ≤ 2 種；無資料分群型編譯期 module。 |
| Last checked | 2026-05-22 |
