# Knowledge Graphs

`knowledge/graphs/` 描述 Knowledge Atoms、source files、skills、shared rules 與 runtime routing surfaces 之間的關係。本目錄保存 graph record 格式與所有 atoms 的 graph records。

## 目前 graph records

### 既有 records（Phase 12）

| Graph record | 用途 | 狀態 |
| --- | --- | --- |
| [`source-boundary.yaml`](source-boundary.yaml) | 連接 active goal / durable roadmap 邊界、content layering 與 governance lifecycle。 | `candidate` |
| [`metadata-navigation.yaml`](metadata-navigation.yaml) | 連接 metadata schema、metadata 子規則、knowledge index、runtime registry 與 summaries。 | `candidate` |
| [`apk-analysis-pilot.yaml`](apk-analysis-pilot.yaml) | 連接 `skills/apk-analysis/` 舊入口與 analysis / workflow / intelligence 候選目的地。 | `candidate` |
| [`apk-highest-leverage-analysis.yaml`](apk-highest-leverage-analysis.yaml) | 連接 APK analysis 最高收益路線 intelligence、原 feedback lesson、workflow 與 runtime route。 | `candidate` |
| [`feedback-promotion-pipeline.yaml`](feedback-promotion-pipeline.yaml) | 連接 feedback promotion pipeline、shared feedback rules、failure learning 與 runtime route。 | `candidate` |

### Intelligence Atoms（Phase 13）

| Graph record | 用途 | 狀態 |
| --- | --- | --- |
| [`intelligence-heuristics.yaml`](intelligence-heuristics.yaml) | 連接 5 個 heuristics atoms（premature-optimization、abstraction-threshold、retry-smell、single-responsibility、test-driven）與 README 索引。 | `candidate` |
| [`intelligence-architecture.yaml`](intelligence-architecture.yaml) | 連接 modular-monolith-vs-microservices atom 與相關 workflow、anti-patterns、tradeoffs。 | `candidate` |
| [`intelligence-tradeoffs.yaml`](intelligence-tradeoffs.yaml) | 連接 postgres-vs-mongodb atom 與相關 architecture、domain、distributed-systems。 | `candidate` |
| [`intelligence-failure.yaml`](intelligence-failure.yaml) | 連接 connection-leak-patterns atom 與相關 heuristics、distributed-systems、failure-learning。 | `candidate` |
| [`intelligence-domain.yaml`](intelligence-domain.yaml) | 連接 aggregate-boundary-heuristics atom 與相關 architecture、tradeoffs、anti-patterns。 | `candidate` |
| [`intelligence-anti-patterns.yaml`](intelligence-anti-patterns.yaml) | 連接 generic-repository-overuse atom 與相關 architecture、domain、根目錄 anti-patterns。 | `candidate` |
| [`intelligence-distributed-systems.yaml`](intelligence-distributed-systems.yaml) | 連接 eventual-consistency-patterns atom 與相關 failure、tradeoffs、architecture。 | `candidate` |
| [`intelligence-business.yaml`](intelligence-business.yaml) | 連接 saas-pricing-heuristics atom 與 business intelligence 層。 | `candidate` |
| [`intelligence-travel.yaml`](intelligence-travel.yaml) | 連接 japan-roadtrip-fatigue、source-triage-heuristic、feasibility-build-heuristic atoms 與 travel workflow。 | `candidate` |
| [`intelligence-apk-analysis-atoms.yaml`](intelligence-apk-analysis-atoms.yaml) | 連接 evidence-first-routing、live-readiness-gates atoms 與 skill、workflow、analysis。 | `candidate` |
| [`intelligence-development-guidance.yaml`](intelligence-development-guidance.yaml) | 連接 docs-first-bdd-closure、risk-translation-heuristic、contract-governance-heuristic atoms 與 workflow、analysis 來源。 | `candidate` |
| [`intelligence-repo-analysis.yaml`](intelligence-repo-analysis.yaml) | 連接 documentation-backfill-heuristic、traceability-heuristic atoms 與 analysis/repo/、workflow 來源。 | `candidate` |
| [`intelligence-agent-architecture.yaml`](intelligence-agent-architecture.yaml) | 連接 agent-architecture 層的所有 11 個 atoms（context-collapse、rule-overload、task-routing、attention-budgeting、failure-recovery、cognitive-boundaries、pilot-first-validation、failure-to-scenario-closure、linked-updates-completeness、decomposition-strategy-selection、stateless-validation-necessity）與 shared-rules/failure-patterns/、decision-efficiency、document-sizing、validation、linked-updates、plans。 | `candidate` |

### Runtime & Pipeline（Phase 14）

| Graph record | 用途 | 狀態 |
| --- | --- | --- |
| [`runtime-pipeline.yaml`](runtime-pipeline.yaml) | 連接 pipeline 所有元件（session-lifecycle、context-flow、guard-chain、relevance-engine）與 budget、health、guards。 | `candidate` |
| [`runtime-prompt-artifacts.yaml`](runtime-prompt-artifacts.yaml) | 連接 prompt artifact generator（artifact-templates、composition-rules）與 workflow、intelligence、pipeline。 | `candidate` |

### Workflow / Analysis / Governance / Feedback / Decisions

| Graph record | 用途 | 狀態 |
| --- | --- | --- |
| [`workflow-layers.yaml`](workflow-layers.yaml) | 連接 workflow 所有子目錄（apk-analysis、software-delivery、repo-analysis、travel-planning）及其子檔案。 | `candidate` |
| [`workflow-software-delivery.yaml`](workflow-software-delivery.yaml) | 連接 software-delivery workflow 的 execution-flow、artifact-gates 與 skills、analysis 來源。 | `candidate` |
| [`workflow-travel-planning.yaml`](workflow-travel-planning.yaml) | 連接 travel-planning workflow 的 execution-flow、artifact-gates 與 skills、intelligence 來源。 | `candidate` |
| [`analysis-layers.yaml`](analysis-layers.yaml) | 連接 analysis 所有子目錄（apk、repo、production、issue）及其子檔案。 | `candidate` |
| [`analysis-repo-methods.yaml`](analysis-repo-methods.yaml) | 連接 analysis/repo/ 的 7 個分析方法（靜態結構、依賴、entrypoint、技術債、文件恢復、追溯性、契約治理）與 skills 來源。 | `candidate` |
| [`governance-layers.yaml`](governance-layers.yaml) | 連接 governance 所有子目錄（lifecycle、validation、cleanup、dependency）。 | `candidate` |
| [`feedback-layers.yaml`](feedback-layers.yaml) | 連接 feedback 所有子目錄（extraction、refinement、replay、promotion、pipeline）。 | `candidate` |
| [`decisions-adr.yaml`](decisions-adr.yaml) | 連接所有 ADR records（ADR-001 到 ADR-005）與 architecture、memory。 | `candidate` |

## Graph 目的

Graphs 協助 agent 理解：

- 必讀 dependencies。
- Related sources。
- Conflicts。
- Replacement 與 deprecation paths。
- 舊 skills 到新分層的 promotion flow。

## 查詢方式

低成本查詢 graph edges：

```bash
ruby scripts/query-knowledge-graph.rb --type depends_on --limit 5
ruby scripts/query-knowledge-graph.rb --query promotion --limit 5
ruby scripts/query-knowledge-graph.rb --source intelligence/engineering/analysis --limit 5
```

查詢結果只作候選 edge list。需要修改 graph、promotion、deprecation 或高信心判斷時，仍需讀回對應 graph YAML 與 canonical source。

## Edge Types

未來 graph records 使用下列 edge labels：

| Edge | 意義 |
| --- | --- |
| `depends_on` | 使用此 atom 前必須先讀 target source。 |
| `related_to` | Target source 可能有幫助，但不是必讀。 |
| `conflicts_with` | Source 可能衝突，需要 rule-weight 或 governance resolution。 |
| `replaces` | Promotion 後，新 atom 取代舊 source。 |
| `preserves_entrypoint` | 新分層 path 保留舊 source 可達性。 |
| `promotes_from` | Atom 從舊 skill / shared rule 抽取或 promotion 而來。 |
| `routes_to` | Index 或 runtime routing 指向 target source。 |

## Graph Record 格式

```yaml
id:
source:
edges:
  - type:
    target:
    reason:
    validation:
status: candidate
```

## 相容性規則

- 使用 canonical repository-relative paths 或 atom IDs。
- 不把 tool mirror paths 建模為 canonical sources。
- 若 graph 使用 `replaces`，lifecycle state 必須已是 promoted 或 deprecated。
- Candidate maps 應使用 `preserves_entrypoint`，不要使用 `replaces`。

## 新增規則

- 新 graph record 必須能解析所有 source / target path。
- Graph record 不可包含 secrets、project incident evidence、本機絕對路徑或 tool mirror source。
- 若 source 改動，graph record 需要 revalidate 或降級 confidence。
- Graph 只描述關係；可執行規則仍以 `shared-rules/` 與 active source-of-truth 文件為準。
- Source、summary、registry 或 lifecycle state 改動時，依 [`../runtime/refresh-policy.yaml`](../runtime/refresh-policy.yaml) 判斷是否 refresh、revalidate 或 downgrade。
