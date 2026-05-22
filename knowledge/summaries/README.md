# Knowledge Summaries

`knowledge/summaries/` 保存 Knowledge Atoms 與 source-of-truth 文件的 compact summaries。本層協助 agent 降低 context loading cost，但不取代舊 skills、enforcement rules 或 canonical source。

## Summary 目的

Summaries 用來協助 agent：

- 在讀完整文件前判斷 source 是否相關。
- 降低 context loading cost。
- 保留 source-of-truth links。
- 支援 small-model 或 checklist-first routing。

## 目前 summaries

| Atom ID | Summary | Source path |
| --- | --- | --- |
| `root.bootstrap.ai-skill` | Root bootstrap 與 enforcement 啟動入口。 | [`root-bootstrap.md`](root-bootstrap.md) |
| `metadata.schema.knowledge-atom` | Knowledge Atom metadata schema v1。 | [`metadata-schema.md`](metadata-schema.md) |
| `architecture.apk-analysis-pilot` | `apk-analysis` 分層 pilot migration map。 | [`apk-analysis-pilot.md`](apk-analysis-pilot.md) |
| `intelligence.apk-highest-leverage-analysis` | APK 分析 highest-leverage route selection engineering intelligence。 | [`apk-highest-leverage-analysis.md`](apk-highest-leverage-analysis.md) |
| `feedback.promotion.pipeline` | Feedback lesson promotion / downgrade pipeline。 | [`feedback-promotion-pipeline.md`](feedback-promotion-pipeline.md) |
| `governance.goal-ledger-boundary` | Active conversation goal 與 durable roadmap goal 邊界。 | [`goal-ledger-boundary.md`](goal-ledger-boundary.md) |
| `skill.app-development-guidance` | App/API/Embedded 開發 guidance、控制項、檢查清單。 | [`development-guidance.md`](development-guidance.md) |
| `skill.travel-planning` | 依目的地/日期/交通規劃行程。 | [`travel-planning.md`](travel-planning.md) |
| `governance.repo-maintenance` | Repo 維護、升級、遷移與治理。 | [`repo-governance.md`](repo-governance.md) |
| `governance.executable-contract-boundary` | 流程 / gate / activation 文件何時需要 YAML contract，且如何投影到 runtime.db。 | [`executable-contract-boundary.md`](executable-contract-boundary.md) |
| `knowledge.navigation` | 知識導航：indexes、summaries、graphs、runtime。 | [`knowledge-navigation.md`](knowledge-navigation.md) |
| `runtime.operations` | Context routing、activation、TTL、pruning。 | [`runtime-operations.md`](runtime-operations.md) |
| `runtime.prompt-cache-alignment` | Provider prompt cache layout：stable prefix、semi-stable middle、volatile suffix。 | [`prompt-cache-alignment.md`](prompt-cache-alignment.md) |
| `models.routing` | 多模型協作、capability profile、compression。 | [`model-routing.md`](model-routing.md) |
| `memory.operations` | 長期記憶、episodic、project、failure memory。 | [`memory-operations.md`](memory-operations.md) |
| `architecture.context-cost-optimization` | Token 成本優化規劃（含 Phase 2.5 prompt cache alignment）。 | [`context-cost-optimization.md`](context-cost-optimization.md) |
| `intelligence.architectural-fit` | Architecture fit analysis 與 DDD adoption boundary。 | [`architectural-fit.md`](architectural-fit.md) |
| `intelligence.requirements-cognition` | Requirements cognition、BDD-lite、acceptance 與 validation target boundary。 | [`requirements-cognition.md`](requirements-cognition.md) |
| `intelligence.migration-seeder-anti-patterns` | 業務資料硬塞進 schema migration 的反模式與替代方案。 | [`migration-seeder-anti-patterns.md`](migration-seeder-anti-patterns.md) |
| `intelligence.vendor-integration-architecture` | 多廠商整合策略選型（Adapter / compile-time module / SPI / out-of-process / hybrid）。 | [`vendor-integration-architecture.md`](vendor-integration-architecture.md) |
| `analysis.dual-token-audit` | 雙簽章/雙加密 token 並存的審計方法（觀察點、五步流程、failure signals）。 | [`dual-token-audit.md`](dual-token-audit.md) |
| `intelligence.migration-feature-bundling` | 大型 migration 把搬遷與新功能綁進同一階段交付的反模式；應採 Parity-First Migration；對 stakeholder 翻譯為「失望總比絕望好」。 | [`migration-feature-bundling.md`](migration-feature-bundling.md) |
| `intelligence.module-count-discipline` | Repo 模組數量紀律；N 對應健康做法（≤5 直寫 / 5-30 模板+門檻 / 30-100 plugin/SPI/catalog / ≥100 警報）；新增 module 五道判斷；build system 中立。 | [`module-count-discipline.md`](module-count-discipline.md) |

## Summary 格式

未來新增 summary 時使用下列形狀：

| 欄位 | 必填 | 用途 |
| --- | --- | --- |
| `Atom ID` | yes | 依 `metadata/schema.md` 命名的 metadata ID。 |
| `Source path` | yes | Canonical repository-relative source path。 |
| `Lifecycle` | yes | `candidate`、`validated`、`stable` 或 `deprecated`。 |
| `Summary` | yes | 一到兩句描述 source。 |
| `When to read` | yes | 載入完整 source 的觸發條件。 |
| `Do not use for` | yes | 邊界與非目標。 |
| `Validation signal` | yes | 如何確認 summary 仍與 source 對齊。 |
| `Last checked` | optional | Summary 進入 stable 後可記錄日期或 commit。 |

## 範例

```markdown
## knowledge.indexes.task-routing

| 欄位 | 值 |
| --- | --- |
| Source path | `knowledge/indexes/README.md` |
| Lifecycle | `candidate` |
| Summary | 將 task intents 導向 canonical primary sources 與 related references。 |
| When to read | 載入深層 skill 或 shared-rule context 前使用。 |
| Do not use for | 不可取代 required dependency reading 或 old skill entrypoints。 |
| Validation signal | Links 可解析，且 routing rows 仍指向 canonical sources。 |
```

## 規則

- Summary 必須連到 source path。
- Summary 不得包含 secrets、raw evidence、private hosts、tokens 或 local absolute paths。
- Summary 不可把 candidate path 升格成 replacement path。
- Source 有實質變更時，必須 revalidate 或 downgrade summary confidence。
- Source、metadata、registry 或 graph 改動時，依 [`../runtime/refresh-policy.yaml`](../runtime/refresh-policy.yaml) 判斷是否 refresh、revalidate 或 downgrade。
