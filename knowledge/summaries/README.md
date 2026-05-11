# Knowledge Summaries

`knowledge/summaries/` 保存 Knowledge Atoms 與 source-of-truth 文件的 compact summaries。本層協助 agent 降低 context loading cost，但不取代舊 skills、shared rules 或 canonical source。

## Summary 目的

Summaries 用來協助 agent：

- 在讀完整文件前判斷 source 是否相關。
- 降低 context loading cost。
- 保留 source-of-truth links。
- 支援 small-model 或 checklist-first routing。

## 目前 summaries

| Atom ID | Summary | Source path |
| --- | --- | --- |
| `root.bootstrap.ai-skill` | Root bootstrap 與 shared-rules 啟動入口。 | [`root-bootstrap.md`](root-bootstrap.md) |
| `metadata.schema.knowledge-atom` | Knowledge Atom metadata schema v1。 | [`metadata-schema.md`](metadata-schema.md) |
| `architecture.apk-analysis-pilot` | `apk-analysis` 分層 pilot migration map。 | [`apk-analysis-pilot.md`](apk-analysis-pilot.md) |
| `intelligence.apk-highest-leverage-analysis` | APK 分析 highest-leverage route selection engineering intelligence。 | [`apk-highest-leverage-analysis.md`](apk-highest-leverage-analysis.md) |
| `feedback.promotion.pipeline` | Feedback lesson promotion / downgrade pipeline。 | [`feedback-promotion-pipeline.md`](feedback-promotion-pipeline.md) |
| `governance.goal-ledger-boundary` | Active conversation goal 與 durable roadmap goal 邊界。 | [`goal-ledger-boundary.md`](goal-ledger-boundary.md) |

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
