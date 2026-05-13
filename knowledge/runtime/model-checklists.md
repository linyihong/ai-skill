# Model Checklists

本檔由 `ruby scripts/generate-model-checklists.rb --write` 產生，將 routing registry 中的 model profile / compression level 轉成 agent 可直接使用的 context-loading checklist。

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
- 涉及 writeback、promotion、shared rules 或 migration 時保持 source-backed。

| Route | Checklist |
| --- | --- |
| `route.governance.durable-goal-boundary` | Primary: `shared-rules/conversation-goal-ledger.md`<br>Compression: `source-backed`<br>Required: `shared-rules/content-layering.md`<br>`governance/lifecycle/README.md`<br>Validation: 長期狀態已落到 durable planning 文件，且 active goal 完成後才刪除。 |
| `route.metadata.knowledge-atom-schema` | Primary: `metadata/schema.md`<br>Compression: `source-backed`<br>Required: `metadata/rules/README.md`<br>`metadata/ranking/README.md`<br>`metadata/confidence/README.md`<br>`metadata/compatibility/README.md`<br>Validation: 欄位可套用到第一批 atom candidates，且 Markdown links 可解析。 |
| `route.feedback.promotion-pipeline` | Primary: `feedback/promotion/README.md`<br>Compression: `source-backed`<br>Required: `shared-rules/feedback-lessons.md`<br>`shared-rules/reusable-guidance-boundary.md`<br>`shared-rules/sanitization.md`<br>`shared-rules/failure-learning-system.md`<br>`governance/lifecycle/README.md`<br>`governance/validation/README.md`<br>Validation: 原 lesson source 保留、promotion target 明確、runtime surfaces 與 close-loop validation 已同步。 |
| `route.models.model-aware-routing` | Primary: `models/profiles/README.md`<br>Compression: `source-backed`<br>Required: `models/compression/README.md`<br>`runtime/routing/README.md`<br>`metadata/ranking/README.md`<br>`knowledge/summaries/README.md`<br>Validation: Profile、compression level、primary source、deferred sources 與 validation signal 可被記錄。 |

### `small`

Guardrails:

- 先讀 index、registry、summary 或 generated lookup。
- 不可跳過 required bootstrap、source-of-truth gate 或 validation signal。
- 需要修改 canonical source、遇到 conflict、缺 validation signal 時升級。

| Route | Checklist |
| --- | --- |
| `route.bootstrap.ai-skill` | Primary: `CORE_BOOTSTRAP.md`<br>Compression: `summary-first`<br>Required: `README.md`<br>`shared-rules/README.md`<br>Validation: Core Bootstrap 3 rules 已讀，skills-index.yaml 已查詢，git status 已檢查。 |
| `route.skill.discovery` | Primary: `skills-index.yaml`<br>Compression: `index-only`<br>Required: `CORE_BOOTSTRAP.md`<br>Validation: Task intent 已對應到 skills-index.yaml 的 triggers，entrypoint 與 summary path 可解析。 |
| `route.runtime.activation-rules` | Primary: `runtime/router/activation-rules.yaml`<br>Compression: `index-only`<br>Required: `CORE_BOOTSTRAP.md`<br>`skills-index.yaml`<br>Validation: 目前 task 已比對 activation-rules.yaml，符合條件的 rules 已載入，不符合的已 deferred。 |
| `route.runtime.context-ttl` | Primary: `runtime/context/ttl-policy.yaml`<br>Compression: `index-only`<br>Required: `runtime/router/activation-rules.yaml`<br>Validation: Context TTL 已套用，過期 context 已 prune，token 使用率在預算內。 |
| `route.runtime.context-loading` | Primary: `runtime/routing/README.md`<br>Compression: `summary-first`<br>Required: `knowledge/indexes/README.md`<br>`metadata/ranking/README.md`<br>`metadata/confidence/README.md`<br>`metadata/compatibility/README.md`<br>Validation: Primary source、deferred sources、source-of-truth gate 與 validation signal 已可被記錄。 |
| `route.runtime.router-flow` | Primary: `runtime/router/README.md`<br>Compression: `index-only`<br>Required: `runtime/router/activation-rules.yaml`<br>`runtime/context/ttl-policy.yaml`<br>Validation: Routing flow 已理解，activation rules 與 TTL policy 已對應到對應階段。 |
| `route.intelligence.engineering.heuristics` | Primary: `intelligence/engineering/heuristics/README.md`<br>Compression: `index-only`<br>Required: `intelligence/engineering/README.md`<br>Validation: 各 atom 有明確原則與決策表，可反查驗證。 |
| `route.runtime.context-ttl-doc` | Primary: `runtime/context/README.md`<br>Compression: `index-only`<br>Required: `runtime/context/ttl-policy.yaml`<br>Validation: TTL policy 已套用，prune strategy 已理解。 |
| `route.intelligence.engineering.agent-architecture` | Primary: `intelligence/engineering/agent-architecture/README.md`<br>Compression: `index-only`<br>Required: `intelligence/engineering/README.md`<br>Validation: 各 atom 有明確原則、症狀表與預防方式，可反查驗證。 |
| `route.feedback.history` | Primary: `feedback/history/README.md`<br>Compression: `index-only`<br>Required: `shared-rules/feedback-lessons.md`<br>Validation: Lesson 已寫入 feedback/history/<domain>/ 對應分類，且 feedback/history/<domain>/README.md 已更新索引。 |

### `specialized`

Guardrails:

- 先讀 routing registry 與 primary source，再讀 domain workflow / technique / adapter。
- 不得讓工具能力覆蓋 shared rules、authorization 或 source-of-truth。
- 保留 domain-specific validation 與 project evidence boundary。

| Route | Checklist |
| --- | --- |
| `route.skill.apk-analysis` | Primary: `skills/apk-analysis/SKILL.md`<br>Compression: `source-backed`<br>Required: `skills/apk-analysis/README.md`<br>`skills/apk-analysis/WORKFLOW.md`<br>`shared-rules/README.md`<br>`shared-rules/dependency-reading.md`<br>Validation: 舊入口仍可讀，新 reference-first paths 可找到，且未 bulk migrate skill content。 |
| `route.intelligence.apk-highest-leverage-path` | Primary: `intelligence/engineering/apk-analysis/highest-leverage-analysis-path.md`<br>Compression: `source-backed`<br>Required: `skills/apk-analysis/SKILL.md`<br>`skills/apk-analysis/WORKFLOW.md`<br>`skills/apk-analysis/feedback_history/common/2026-05-07_131000-highest-leverage-analysis-path.md`<br>Validation: 已比較可用 routes、選定 evidence-to-cost ratio 最高路線、記錄 fallback 與 attribution 回補條件。 |
| `route.skill.app-development-guidance` | Primary: `skills/app-development-guidance/SKILL.md`<br>Compression: `source-backed`<br>Required: `skills/app-development-guidance/README.md`<br>`shared-rules/README.md`<br>Validation: 舊入口仍可讀，新 reference-first paths 可找到，且未 bulk migrate skill content。 |
| `route.skill.travel-planning` | Primary: `skills/travel-planning/SKILL.md`<br>Compression: `source-backed`<br>Required: `skills/travel-planning/README.md`<br>`shared-rules/README.md`<br>Validation: 舊入口仍可讀，新 reference-first paths 可找到，且未 bulk migrate skill content。 |
| `route.runtime.onboarding` | Primary: `runtime/onboarding/README.md`<br>Compression: `summary-first`<br>Required: `skills/apk-analysis/SKILL.md`<br>`skills/app-development-guidance/SKILL.md`<br>`skills/travel-planning/SKILL.md`<br>Validation: 各 quickstart 的步驟可依序執行，且與對應 SKILL.md 的 Quick Start 摘要一致。 |
| `route.analysis.apk.workflows` | Primary: `analysis/apk/workflows/README.md`<br>Compression: `summary-first`<br>Required: `analysis/apk/README.md`<br>`skills/apk-analysis/SKILL.md`<br>Validation: 各 workflow 有明確步驟與產出格式，可依序執行。 |
| `route.intelligence.apk-analysis.atoms` | Primary: `intelligence/engineering/apk-analysis/README.md`<br>Compression: `summary-first`<br>Required: `skills/apk-analysis/SKILL.md`<br>`analysis/apk/README.md`<br>Validation: 各 atom 有明確決策表或信號表，可反查驗證。 |
| `route.validation.ai-decision-contract` | Primary: `validation/README.md`<br>Compression: `summary-first`<br>Required: `validation/scenarios/apk-analysis/flutter-aot-hooking-v1.yaml`<br>`validation/scenarios/apk-analysis/local-proxy-vs-pinning-v1.yaml`<br>`validation/scenarios/apk-analysis/early-hook-prevention-v1.yaml`<br>`validation/scenarios/failure-derived/entrypoint-drift-v1.yaml`<br>`validation/scenarios/failure-derived/source-mirror-write-v1.yaml`<br>`validation/scenarios/failure-derived/shared-rules-architecture-drift-v1.yaml`<br>`validation/scenarios/failure-derived/feedback-history-consolidation-v1.yaml`<br>Validation: Scenario 的 expected_route 與 actual trace 可比對；forbidden_routes 未被使用。 |

## Escalation Checklist

- Summary / registry 與 source-of-truth 可能不一致時，讀全文。
- 任務需要修改、commit、push、readback 或 promotion 時，升級到 `source-backed`。
- 涉及 safety、secrets、authorization、source/mirror 或 destructive actions 時，升級到 full source 和 shared rules。
- Routing registry 指向 candidate path，但 old entrypoint 仍 active 時，保留 old entrypoint gate。
- Validation signal 不足以支持結論時，停止並讀 required dependencies。

## Validation

- 產生前應先確認 `routing-registry.yaml` 可通過 `ruby scripts/validate-knowledge-runtime.rb`。
- 產生後應重新執行 `ruby scripts/validate-knowledge-runtime.rb`，檢查本 report links。
- 本檔是 generated view，不取代 model source docs 或 routing registry。
