# Model Checklists

本檔由 `ai-skill runtime refresh` 產生，將 routing registry 中的 model profile / compression level 轉成 agent 可直接使用的 context-loading checklist。

## Source Surfaces

| Surface | Path | Purpose |
| --- | --- | --- |
| Routing registry | [`routing-registry.yaml`](routing-registry.yaml) | 提供 route、primary source、dependencies、model profile 與 compression level。 |
| Model profiles | [`../../models/profiles/README.md`](../../models/profiles/README.md) | 定義 profile guardrails。 |
| Compression strategy | [`../../models/compression/README.md`](../../models/compression/README.md) | 定義 escalation rules。 |

## Profile Checklists

### `large`

Guardrails:

- 讀 primary source、required dependencies 與 task-relevant related sources。
- 回報 deferred sources 與 validation signal。
- 涉及 writeback、promotion、enforcement rules 或 migration 時保持 source-backed。

| Route | Checklist |
| --- | --- |
| `route.governance.durable-goal-boundary` | Primary: `enforcement/conversation-goal-ledger.md`<br>Compression: `source-backed`<br>Required: `enforcement/content-layering.md`<br>`governance/lifecycle/README.md`<br>Validation: 長期狀態已落到 durable planning 文件，且 active goal 完成後才刪除。 |
| `route.metadata.knowledge-atom-schema` | Primary: `metadata/schema.md`<br>Compression: `source-backed`<br>Required: `metadata/rules/README.md`<br>`metadata/ranking/README.md`<br>`metadata/confidence/README.md`<br>`metadata/compatibility/README.md`<br>Validation: 欄位可套用到第一批 atom candidates，且 Markdown links 可解析。 |
| `route.feedback.promotion-pipeline` | Primary: `feedback/promotion/README.md`<br>Compression: `source-backed`<br>Required: `enforcement/feedback-lessons.md`<br>`enforcement/reusable-guidance-boundary.md`<br>`enforcement/sanitization.md`<br>`enforcement/failure-learning-system.md`<br>`governance/lifecycle/README.md`<br>`governance/validation/README.md`<br>Validation: 原 lesson source 保留、promotion target 明確、runtime surfaces 與 close-loop validation 已同步。 |
| `route.models.model-aware-routing` | Primary: `models/README.md`<br>Compression: `source-backed`<br>Required: `models/routing/README.md`<br>`models/capabilities/README.md`<br>`models/compression/README.md`<br>`models/workflow-adaptation/README.md`<br>`models/governance/README.md`<br>`runtime/README.md`<br>`metadata/ranking/README.md`<br>`knowledge/summaries/README.md`<br>Validation: Strategy、capability dimensions、compression、fallback behavior、workflow shape 與 validation target 可被記錄，且 validation/scenarios/models 覆蓋主要 routing cases。 |
| `route.memory.retrieval-activation` | Primary: `memory/README.md`<br>Compression: `source-backed`<br>Required: `memory/retrieval-governance/README.md`<br>`memory/working/README.md`<br>`memory/summary/README.md`<br>`memory/episodic/README.md`<br>`memory/project/README.md`<br>`memory/failure/README.md`<br>`memory/decision/README.md`<br>`governance/ai-runtime-governance/cognitive-state-governance.md`<br>Validation: Memory activation decision 能記錄 trigger、candidate memory、qualification、replay budget、current-source revalidation 與 discard / promotion decision。 |

### `medium`

Guardrails:

- 先確認 registry record 的 model profile。
- 依 `models/profiles/README.md` 與 `models/compression/README.md` 選讀取深度。

| Route | Checklist |
| --- | --- |
| `route.governance.cognitive-state-evidence` | Primary: `governance/ai-runtime-governance/cognitive-state-governance.md`<br>Compression: `summary-first`<br>Required: `enforcement/evidence-hierarchy.md`<br>`enforcement/escalation-policy.md`<br>`metadata/evidence/README.md`<br>`metadata/evidence/domain-policies.yaml`<br>`runtime/runtime.db`<br>Validation: Claim scope、evidence owner、confidence integrity、contradiction propagation、runtime primitive boundary 與 scenario coverage 已檢查。 |

### `small`

Guardrails:

- 先讀 index、registry、summary 或 generated lookup。
- 不可跳過 required bootstrap、source-of-truth gate 或 validation signal。
- 需要修改 canonical source、遇到 conflict、缺 validation signal 時升級。

| Route | Checklist |
| --- | --- |
| `route.bootstrap.ai-skill` | Primary: `CORE_BOOTSTRAP.md`<br>Compression: `summary-first`<br>Required: `README.md`<br>`enforcement/README.md`<br>Validation: Core Bootstrap 3 rules 已讀，git status 已檢查。 |
| `route.runtime.phase-machine` | Primary: `runtime/runtime.db`<br>Compression: `source-backed`<br>Required: `CORE_BOOTSTRAP.md`<br>`README.md`<br>Validation: 目前 phase 已確認，allowed_actions 與 forbidden_actions 已載入，phase transition rules 已檢查。 |
| `route.runtime.obligation-ledger` | Primary: `runtime/runtime.db`<br>Compression: `source-backed`<br>Required: `runtime/runtime.db`<br>Validation: 本 phase 的 obligations 已確認，pending obligations 已記錄，blocking gates 已檢查。 |
| `route.runtime.blocking-gates` | Primary: `runtime/runtime.db`<br>Compression: `source-backed`<br>Required: `runtime/runtime.db`<br>`runtime/runtime.db`<br>Validation: 本 phase 的 blocking gates 已檢查，無未通過的 critical/high gates，phase transition 可進行。 |
| `route.runtime.recovery` | Primary: `runtime/runtime.db`<br>Compression: `source-backed`<br>Required: `governance/ai-runtime-governance/recovery-retry-governance.md`<br>`intelligence/engineering/agent-architecture/failure-recovery.md`<br>`intelligence/engineering/agent-architecture/cognitive-boundaries.md`<br>`runtime/runtime.db`<br>`runtime/runtime.db`<br>`metadata/recovery/escalation-levels.yaml`<br>`metadata/recovery/domain-policies.yaml`<br>Validation: 阻斷 gate 已對應到 strategy；retry limit、strategy change、domain policy、required reload set、execution graph rebuild 與 recovery validation 已檢查。 |
| `route.runtime.scheduler` | Primary: `runtime/runtime.db`<br>Compression: `source-backed`<br>Required: `runtime/runtime.db`<br>`runtime/runtime.db`<br>`runtime/runtime.db`<br>Validation: Queue 已依 priority 排序，blocking gates 優先於 obligations，dependencies 已解析。 |
| `route.runtime.transactions` | Primary: `runtime/runtime.db`<br>Compression: `source-backed`<br>Required: `runtime/runtime.db`<br>`enforcement/dependency-reading.md`<br>Validation: Transaction state 正確，所有 rules 已檢查，templates 已套用。 |
| `route.skill.discovery` | Primary: `knowledge/runtime/routing-registry.yaml`<br>Compression: `index-only`<br>Required: `CORE_BOOTSTRAP.md`<br>Validation: Task intent 已對應到 routing-registry.yaml 的 triggers，entrypoint 與 summary path 可解析。 |
| `route.runtime.context-ttl` | Primary: `governance/ai-runtime-governance/context-attention-governance.md`<br>Compression: `summary-first`<br>Required: `intelligence/engineering/agent-architecture/context-collapse.md`<br>`intelligence/engineering/agent-architecture/attention-budgeting.md`<br>`governance/lifecycle/context-ttl-philosophy.md`<br>Validation: Summary-first、attention budget、decision externalization、recap checkpoint 與 task-boundary prune 已檢查；必要時再讀 TTL policy。 |
| `route.runtime.prompt-cache-alignment` | Primary: `runtime/context/prompt-cache-playbook.md`<br>Compression: `summary-first`<br>Required: `enforcement/prompt-cache-efficiency.md`<br>`metadata/schema.md`<br>`runtime/runtime.db`<br>Validation: Prompt cache playbook、enforcement rule、metadata provider_cache 欄位與 activation rule 已同步。 |
| `route.governance.ai-runtime-five-step` | Primary: `governance/ai-runtime-governance/five-step-ai-governance.md`<br>Compression: `summary-first`<br>Required: `governance/ai-runtime-governance/README.md`<br>`intelligence/engineering/philosophy/musk-five-step-algorithm.md`<br>Validation: 新增項目已通過 necessity、deletion、simplification、cycle-time 與 automation-last 檢查，且沒有加入 always-load context。 |
| `route.runtime.context-loading` | Primary: `governance/lifecycle/routing-philosophy.md`<br>Compression: `summary-first`<br>Required: `knowledge/indexes/README.md`<br>`metadata/ranking/README.md`<br>`metadata/confidence/README.md`<br>`metadata/compatibility/README.md`<br>Validation: Primary source、deferred sources、source-of-truth gate 與 validation signal 已可被記錄。 |
| `route.governance.routing-signal` | Primary: `governance/ai-runtime-governance/routing-signal-governance.md`<br>Compression: `summary-first`<br>Required: `intelligence/engineering/agent-architecture/task-routing.md`<br>`governance/lifecycle/routing-philosophy.md`<br>Validation: Task intent 已先於 path signal 確認；primary_source、negative signals、multi-route disambiguation 與 route validation signal 已檢查。 |
| `route.governance.validation-scenario` | Primary: `governance/ai-runtime-governance/validation-scenario-governance.md`<br>Compression: `summary-first`<br>Required: `intelligence/engineering/agent-architecture/stateless-validation-necessity.md`<br>`intelligence/engineering/agent-architecture/failure-to-scenario-closure.md`<br>`validation/README.md`<br>Validation: Stateless reproduction、no answer leakage、failure class、expected/forbidden behavior 與 traceability gates 已檢查。 |
| `route.runtime.router-flow` | Primary: `runtime/README.md`<br>Compression: `index-only`<br>Required: `runtime/runtime.db`<br>`runtime/runtime.db`<br>Validation: Routing flow 已理解，activation rules 與 TTL policy 已對應到對應階段。 |
| `route.intelligence.engineering.heuristics` | Primary: `intelligence/engineering/heuristics/README.md`<br>Compression: `index-only`<br>Required: `intelligence/engineering/README.md`<br>Validation: 各 atom 有明確原則與決策表，可反查驗證。 |
| `route.runtime.context-ttl-doc` | Primary: `governance/lifecycle/context-ttl-philosophy.md`<br>Compression: `index-only`<br>Required: `runtime/runtime.db`<br>Validation: TTL policy 已套用，prune strategy 已理解。 |
| `route.workflow.documentation-ai-native` | Primary: `workflow/documentation/README.md`<br>Compression: `summary-first`<br>Required: `governance/ai-runtime-governance/documentation-context-governance.md`<br>`intelligence/engineering/agent-architecture/index-first-documentation.md`<br>`enforcement/content-layering.md`<br>Validation: 新文件或目錄具 index-first 導航、分類維度（kind/audience/stability）已標註； README-as-router、停止條件與單一真相已符合 documentation context governance； 長文已按 document-sizing 拆分；語言與工具敘述依專案自訂 policy（本 route 不預設 tool-neutral）。  |
| `route.intelligence.requirements-cognition` | Primary: `intelligence/engineering/requirements/README.md`<br>Compression: `summary-first`<br>Required: `intelligence/engineering/requirements/product-alignment/README.md`<br>`intelligence/engineering/requirements/behavior-modeling/README.md`<br>`intelligence/engineering/requirements/specification-quality/README.md`<br>`intelligence/engineering/requirements/validation-thinking/README.md`<br>`workflow/software-delivery/requirements/README.md`<br>`governance/ai-runtime-governance/software-delivery-governance.md`<br>Validation: Impact / journey evidence → requirement → behavior contract → acceptance criteria → validation target → execution artifact is traceable; no unconfirmed feature is written as requirement; high coverage is not treated as test effectiveness without targeted proof when risk is high. |
| `route.intelligence.architectural-fit` | Primary: `intelligence/engineering/architecture/architectural-fit/README.md`<br>Compression: `summary-first`<br>Required: `intelligence/engineering/architecture/README.md`<br>`intelligence/engineering/domain/README.md`<br>`intelligence/engineering/architecture/domain-modeling/README.md`<br>`workflow/software-delivery/architecture/README.md`<br>`governance/ai-runtime-governance/software-delivery-architecture-governance.md`<br>Validation: Chosen strategy、rejected lighter option、rejected heavier option、fit evidence 與 upgrade/downgrade trigger 已輸出；沒有預設套用 DDD/CQRS/event sourcing。 |
| `route.intelligence.engineering.agent-architecture` | Primary: `intelligence/engineering/agent-architecture/README.md`<br>Compression: `index-only`<br>Required: `intelligence/engineering/README.md`<br>Validation: 各 atom 有明確原則、症狀表與預防方式，可反查驗證。 |
| `route.feedback.history` | Primary: `feedback/history/README.md`<br>Compression: `index-only`<br>Required: `enforcement/feedback-lessons.md`<br>Validation: Lesson 已寫入 feedback/history/<domain>/ 對應分類，且 feedback/history/<domain>/README.md 已更新索引。 |
| `route.runtime.decision-recording` | Primary: `runtime/runtime.db`<br>Compression: `source-backed`<br>Required: `runtime/README.md`<br>`constitution/README.md`<br>Validation: 本輪若有鎖定決策，對應 tier 的檔案與 README 索引已更新。 |
| `route.constitution.adr` | Primary: `constitution/README.md`<br>Compression: `summary-first`<br>Required: `CORE_BOOTSTRAP.md`<br>`runtime/runtime.db`<br>Validation: constitution/README.md 已讀取，ADR 清單已查詢，graph edge 已確認。 |
| `route.governance.decision-promotion-pipeline` | Primary: `governance/lifecycle/decision-promotion-pipeline.yaml`<br>Compression: `source-backed`<br>Required: `governance/lifecycle/decision-promotion-pipeline.md`<br>`constitution/README.md`<br>`memory/decision/README.md`<br>`runtime/runtime.db`<br>Validation: Promotion target is selected by content type; ADR is created only when constitution criteria match; execution-affecting targets handle runtime projection.  |
| `route.architecture.permanent-docs` | Primary: `architecture/README.md`<br>Compression: `summary-first`<br>Required: `CORE_BOOTSTRAP.md`<br>Validation: architecture/README.md 已讀取，架構文件清單已查詢。 |
| `route.evaluations.scenario-results` | Primary: `evaluations/README.md`<br>Compression: `index-only`<br>Required: `validation/README.md`<br>Validation: Evaluation 記錄已依 scenario 分類存放，格式與 template 一致。 |
| `route.ai-tools.agent-onboarding` | Primary: `ai-tools/agent-onboarding.yaml`<br>Compression: `source-backed`<br>Required: `ai-tools/agent-onboarding.md`<br>`ai-tools/README.md`<br>`runtime/runtime.db`<br>Validation: Tool bootstrap entrypoints remain thin pointers; reusable rules stay in ai-tools, enforcement, governance, workflow, or runtime.db; enabled YAML contract is synced into runtime.db.generated_surfaces.  |
| `route.ai-tools.new-project-onboarding` | Primary: `ai-tools/new-project-onboarding.yaml`<br>Compression: `source-backed`<br>Required: `ai-tools/new-project-onboarding.md`<br>`ai-tools/agent-onboarding.yaml`<br>`ai-tools/README.md`<br>`scripts/ai-skill-cli/internal/app/init_project.go`<br>`runtime/runtime.db`<br>Validation: Agent adapter additions are reflected in project onboarding when they affect init-project or project-level bootstrap files; enabled YAML contract is synced into runtime.db.generated_surfaces.  |
| `route.tools.metadata-routing` | Primary: `tools/README.md`<br>Compression: `index-only`<br>Required: `ai-tools/README.md`<br>Validation: Tool metadata 已定義，compression 策略與 routing 規則已建立。 |
| `route.traces.decision-traces` | Primary: `traces/README.md`<br>Compression: `index-only`<br>Required: `validation/README.md`<br>Validation: Trace 記錄已依 scenario 分類存放，格式與 template 一致。 |
| `route.anti-patterns.runtime-patterns` | Primary: `anti-patterns/README.md`<br>Compression: `summary-first`<br>Required: `CORE_BOOTSTRAP.md`<br>Validation: Anti-pattern 已依格式記錄，症狀、預防與恢復方式已定義。 |
| `route.runtime.compiler` | Primary: `runtime/runtime.db`<br>Compression: `source-backed`<br>Required: `runtime/runtime.db`<br>`runtime/runtime.db`<br>Validation: 所有 modified sources 已編譯，runtime.db 的 generated_surfaces 表包含最新記錄，`ai-skill runtime validate` 回傳 exit 0。 |
| `route.runtime.intelligence-routing` | Primary: `runtime/runtime.db`<br>Compression: `index-only`<br>Required: `runtime/runtime.db`<br>Validation: Task intent 已對應到 intelligence-routing.yaml 的 domain，applicable_phases 已檢查，domain README 已載入。 |
| `route.runtime.output-governance` | Primary: `runtime/README.md`<br>Compression: `source-backed`<br>Required: `runtime/runtime.db`<br>`runtime/runtime.db`<br>Validation: Language consistency 已確認，sanitization 已通過，tool neutrality 已檢查， format compliance 已驗證，governance gates 全部通過。  |
| `route.runtime.distributed` | Primary: `runtime/README.md`<br>Compression: `source-backed`<br>Required: `runtime/runtime.db`<br>`runtime/runtime.db`<br>`runtime/runtime.db`<br>Validation: Distributed locks 已正確 acquire/release，multi-agent coordination rules 已遵守， async job lifecycle 狀態轉換正確，無 deadlock 或 stale state。  |
| `route.runtime.cognitive-modes` | Primary: `runtime/cognitive-modes.yaml`<br>Compression: `summary-first`<br>Required: `runtime/cognitive-modes-discovery.yaml`<br>`runtime/runtime.db`<br>Validation: cognitive_modes table 有本 task row，或 final report 含 Cognitive Mode 區塊（4 維值 + 理由）。generated_surfaces 含 runtime.cognitive_modes.contract / runtime.cognitive_modes.discovery / runtime.cognitive_modes.phase_integration / runtime.cognitive_modes.governance_integration / runtime.cognitive_modes.memory_integration。  |
| `route.governance.executable-contract-boundary` | Primary: `governance/lifecycle/executable-contract-boundary.yaml`<br>Compression: `source-backed`<br>Required: `governance/lifecycle/executable-contract-boundary.md`<br>`workflow/software-delivery/requirements/pre-build-interrogation.md`<br>`governance/lifecycle/compiler-philosophy.md`<br>`runtime/README.md`<br>`runtime/runtime.db`<br>Validation: YAML contract inventory 已載入；pre-build interrogation 已確認 canonical owner、projection boundary 與 duplication risk； runtime_projection.enabled contract 已由 compiler 投影到 runtime.db；普通 metadata / graph / validation YAML 不會自動進 runtime。  |
| `route.governance.system-upgrade` | Primary: `governance/lifecycle/system-upgrade-governance.yaml`<br>Compression: `source-backed`<br>Required: `governance/lifecycle/system-upgrade-governance.md`<br>`governance/lifecycle/executable-contract-boundary.yaml`<br>`workflow/software-delivery/requirements/pre-build-interrogation.yaml`<br>`plans/README.md`<br>`enforcement/linked-updates.md`<br>`runtime/runtime.db`<br>Validation: system-upgrade-governance.yaml 已載入並投影到 runtime.db generated_surfaces； active plan 具備 Runtime Execution Path + Trigger flow；runtime refresh/validate 與 close-loop evidence 已記錄。  |
| `route.governance.knowledge-update-flow` | Primary: `governance/lifecycle/knowledge-update-flow.yaml`<br>Compression: `source-backed`<br>Required: `governance/lifecycle/knowledge-update-flow.md`<br>`governance/ai-runtime-governance/linked-update-governance.md`<br>`intelligence/engineering/agent-architecture/linked-updates-completeness.md`<br>`enforcement/linked-updates.md`<br>`runtime/runtime.db`<br>`runtime/runtime.db`<br>`runtime/README.md`<br>Validation: 目前 knowledge update 的步驟已確認，entry conditions 已滿足，exit conditions 已檢查，reference sources 已載入；linked update completeness gates 已通過， 下一步驟已決定。  |

### `specialized`

Guardrails:

- 先讀 routing registry 與 primary source，再讀 domain workflow / technique / adapter。
- 不得讓工具能力覆蓋 enforcement rules、authorization 或 source-of-truth。
- 保留 domain-specific validation 與 project evidence boundary。

| Route | Checklist |
| --- | --- |
| `route.workflow.apk-analysis` | Primary: `workflow/apk-analysis/execution-flow.md`<br>Compression: `source-backed`<br>Required: `workflow/apk-analysis/artifact-gates.md`<br>`enforcement/README.md`<br>`enforcement/dependency-reading.md`<br>`enforcement/authorization-scope.md`<br>Validation: 新分層路徑可讀取，workflow 與 analysis 內容已分離。 |
| `route.intelligence.apk-highest-leverage-path` | Primary: `intelligence/engineering/analytical-reasoning/highest-leverage-analysis-path.md`<br>Compression: `source-backed`<br>Required: `workflow/apk-analysis/execution-flow.md`<br>`workflow/apk-analysis/artifact-gates.md`<br>`feedback/history/apk-analysis/common/2026-05-07_131000-highest-leverage-analysis-path.md`<br>Validation: 已比較可用 routes、選定 evidence-to-cost ratio 最高路線、記錄 fallback 與 attribution 回補條件。 |
| `route.workflow.software-delivery` | Primary: `workflow/software-delivery/execution-flow.md`<br>Compression: `source-backed`<br>Required: `governance/ai-runtime-governance/software-delivery-governance.md`<br>`workflow/software-delivery/artifact-gates.md`<br>`enforcement/README.md`<br>`enforcement/linked-updates.md`<br>`intelligence/engineering/requirements/README.md`<br>`intelligence/engineering/development/docs-first-bdd-closure.md`<br>Validation: Pre-build interrogation、change intake、requirements cognition、contract precedence、BDD closure、artifact completeness、test strategy、performance evidence、same-session closure 與 architecture fit analysis 已檢查；workflow、analysis、intelligence、metadata 與 governance 內容已分離。 |
| `route.workflow.greenfield` | Primary: `workflow/greenfield/execution-flow.md`<br>Compression: `source-backed`<br>Required: `workflow/greenfield/README.md`<br>`enforcement/README.md`<br>Validation: 新分層路徑可讀取，workflow 與 templates 內容完整。 |
| `route.workflow.travel-planning` | Primary: `workflow/travel-planning/execution-flow.md`<br>Compression: `source-backed`<br>Required: `workflow/travel-planning/artifact-gates.md`<br>`enforcement/README.md`<br>Validation: 新分層路徑可讀取，workflow 與 analysis 內容已分離。 |
| `route.runtime.onboarding` | Primary: `runtime/onboarding/apk-analysis-setup.md`<br>Compression: `summary-first`<br>Required: `workflow/apk-analysis/execution-flow.md`<br>`workflow/software-delivery/execution-flow.md`<br>`workflow/travel-planning/execution-flow.md`<br>Validation: 各 quickstart 的步驟可依序執行，且與對應 workflow 的內容一致。 |
| `route.analysis.apk.workflows` | Primary: `analysis/apk/workflows/README.md`<br>Compression: `summary-first`<br>Required: `analysis/apk/README.md`<br>`workflow/apk-analysis/execution-flow.md`<br>Validation: 各 workflow 有明確步驟與產出格式，可依序執行。 |
| `route.analysis.web` | Primary: `analysis/web/README.md`<br>Compression: `summary-first`<br>Required: `analysis/README.md`<br>Validation: 目標網站已評估（技術棧、JS 需求、anti-bot 保護），工具已選擇（HTTP/Dynamic/Stealth）， 提取策略已設計（selector/adaptive parsing），風險已評估（legal/technical/data quality）。  |
| `route.intelligence.apk-analysis.atoms` | Primary: `intelligence/engineering/apk-analysis/README.md`<br>Compression: `summary-first`<br>Required: `intelligence/engineering/analytical-reasoning/README.md`<br>`workflow/apk-analysis/execution-flow.md`<br>`analysis/apk/README.md`<br>Validation: 各 atom 有明確決策表或信號表，可反查驗證。 |
| `route.validation.ai-decision-contract` | Primary: `validation/README.md`<br>Compression: `summary-first`<br>Required: `validation/scenarios/apk-analysis/flutter-aot-hooking-v1.yaml`<br>`validation/scenarios/apk-analysis/local-proxy-vs-pinning-v1.yaml`<br>`validation/scenarios/apk-analysis/early-hook-prevention-v1.yaml`<br>`validation/scenarios/failure-derived/entrypoint-drift-v1.yaml`<br>`validation/scenarios/failure-derived/source-mirror-write-v1.yaml`<br>`validation/scenarios/failure-derived/shared-rules-architecture-drift-v1.yaml`<br>`validation/scenarios/failure-derived/feedback-history-consolidation-v1.yaml`<br>`validation/scenarios/failure-derived/runtime-recovery-navigation-mismatch.yaml`<br>`validation/scenarios/failure-derived/runtime-recovery-user-contradiction.yaml`<br>`validation/scenarios/failure-derived/runtime-recovery-source-miss.yaml`<br>`validation/scenarios/architecture/cargo-cult-ddd.yaml`<br>`validation/scenarios/architecture/architecture-fit-mismatch.yaml`<br>`validation/scenarios/architecture/overengineering-detection.yaml`<br>`validation/scenarios/architecture/bounded-context-collapse.yaml`<br>`validation/scenarios/architecture/aggregate-explosion.yaml`<br>`validation/scenarios/software-delivery/requirement-contradiction.yaml`<br>`validation/scenarios/software-delivery/product-impact-misalignment.yaml`<br>`validation/scenarios/software-delivery/missing-validation-target.yaml`<br>`validation/scenarios/software-delivery/stale-acceptance-criteria.yaml`<br>`validation/scenarios/software-delivery/behavior-scope-overclaim.yaml`<br>`validation/scenarios/software-delivery/mutation-testing-effectiveness.yaml`<br>Validation: Scenario 的 expected_route 與 actual trace 可比對；forbidden_routes 未被使用。 |

## Executable Contract Checklist-First Path

當任務涉及 owner-layer executable YAML contract，small / weaker agents 應先用以下 checklist，不得只讀 Markdown 或 metadata YAML：

1. 讀 [`../../metadata/executable-contract-schema.md`](../../metadata/executable-contract-schema.md)，確認 `schema_version: executable-contract/v1`、`runtime_projection.enabled`、`target_key` 與 execution-bearing fields。
2. 讀 [`../../governance/lifecycle/executable-contract-inventory.yaml`](../../governance/lifecycle/executable-contract-inventory.yaml)，確認 source 是 `contract_exists`、`contract_required`、`markdown_only` 或 `not_applicable`。
3. 若有 companion YAML，讀 YAML 的 `activation`、`required_sources`、`steps`、`gates`、`failure_modes`、`final_status_report`；Markdown 只提供背景與維護脈絡。
4. 新增或修改 executable contract 後，執行 `ai-skill runtime compile`、`ai-skill runtime refresh`、`ai-skill runtime validate`，並查 `runtime/runtime.db generated_surfaces` 的 `source_path`、`target_key`、`status`。
5. 若只看到 `metadata/rules/*.yaml`、front-matter、graph 或 routing YAML，不得當成 executable contract，除非補齊 schema 並啟用 runtime projection。

## Escalation Checklist

- Summary / registry 與 source-of-truth 可能不一致時，讀全文。
- 任務需要修改、commit、push、readback 或 promotion 時，升級到 `source-backed`。
- 涉及 safety、secrets、authorization、source/mirror 或 destructive actions 時，升級到 full source 和 enforcement rules。
- Routing registry 指向 candidate path，但 old entrypoint 仍 active 時，保留 old entrypoint gate。
- Validation signal 不足以支持結論時，停止並讀 required dependencies。

## Validation

- 產生前應先確認 `routing-registry.yaml` 可通過 `ai-skill runtime validate`。
- 產生後應重新執行 `ai-skill runtime validate`，檢查本 report links。
- 本檔是 generated view，不取代 model source docs 或 routing registry。
