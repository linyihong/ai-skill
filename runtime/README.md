## Owner-Layer Executable Contracts

`runtime/` 不接收 governance、enforcement 或 workflow source ownership。流程型 source 若需要 YAML 化，YAML contract 留在原 owner layer，例如 `governance/`、`enforcement/`、`workflow/` 或 `metadata/rules/`。

會影響 agent 執行的 YAML contract 必須設定：

```yaml
runtime_projection:
  enabled: true
```

Runtime compiler 會把這類 contract 投影到 `runtime.db` 的 `generated_surfaces`。這與 runtime internal config 不同：runtime internal config 的 canonical copy 仍直接保存在 `runtime_config_documents` 與 projection tables，不再保留 committed `runtime/**/*.yaml` mirror。

# Runtime

本層是 Ai-skill 的可執行 runtime 層：只放機器可查、可驗證、可重建的結構化狀態。

## 核心規則：不再提交 Runtime YAML

`runtime/` 目錄內的 runtime config 不再以 `*.yaml` 作為 committed source。只要內容已經是 deterministic、machine-readable，並且可以投影到 runtime tables，就必須放進 [`runtime.db`](runtime.db)：

- 完整 canonical document 放在 `runtime_config_documents`
- table 投影關係放在 `runtime_config_projections`
- 舊 source path 相容 manifest 放在 `runtime_source_files`，並標記 `source_kind='db'`
- 查詢用資料放在各專屬 projection table，例如 `phase_machine`、`obligation_ledger`、`blocking_gates`、`runtime_budget`

不要再新增或提交 `runtime/**/*.yaml` mirror。若需要人類可讀 diff，可用臨時匯出檔檢視，但不要把匯出 mirror commit 回 repo。

## 放什麼

- SQLite runtime database：`runtime.db`
- schema、registry、state definitions、transitions
- activation graph、routing rules、guard definitions
- runtime metadata：budget、health、scheduler
- navigation-only README：說明 table、來源邊界與維護方式

## 不放什麼

- 長篇解釋、哲學、推理過程
- 教學、架構討論、設計理由
- knowledge base 內容、領域知識
- 人類流程說明、操作步驟

這些內容屬於：
- `governance/`：設計哲學、生命週期、validation governance
- `workflow/`：人類可讀的 execution flow
- `intelligence/`：heuristics、分析推理、工程判斷
- `enforcement/`：可執行政策規則
- `analysis/`：分析思路與方法

## Runtime 領域表

| 領域 | 位置 | 用途 |
|--------|------|-------------|
| Activation | `activation_rules`, `activation_rules_mirror`, `core_bootstrap_rules` | lazy-load rule 與 activation condition |
| Routing | `runtime.db` activation/runtime config tables | runtime routing support；knowledge routing 仍在 `knowledge/runtime/` |
| Discovery | `discovery_checkpoints`, `capability_checkpoints` | phase-aware capability discovery checkpoint |
| Phases | `phases`, `phase_machine`, `phase_transitions` | execution phase state machine |
| Obligations | `obligations`, `obligation_ledger` | 每個 phase 的 atomic duties |
| Gates | `gates`, `blocking_gates` | phase transition 的 prerequisites |
| Compiler | `compiler_rules`, `runtime_config_documents`, `runtime_config_projections` | SQLite canonical config + deterministic prose → SQLite projections |
| Runtime DB | [`runtime.db`](runtime.db) | committed canonical runtime registry；變更時必須 commit |
| State DB | `runtime-state.db` | future mutable execution state；目前不提交 |
| Generated | [`generated/`](generated/) | legacy compiled surfaces；已遷移到 SQLite |
| Transactions | `transaction_*` | writeback transaction state machine 與 templates |
| Pipeline | `pipeline_context_flow`, `guard_chain`, `relevance_engine`, `session_lifecycle` | context flow、guard chain、relevance engine |
| Recovery | `recovery_strategies`, `state_repair`, `obligation_rebuild`, `phase_reconciliation` | recovery strategy 與 phase repair |
| Scheduler | `execution_queue`, `priority_scheduler` | execution queue 與 priority scheduler |
| Guards | `circuit_breaker`, `context_pollution` | circuit breaker、context pollution、mismatch escalation |
| Onboarding | [`onboarding/`](onboarding/) | 新專案 / 新任務 setup guidance |
| Output Governance | `language_policy`, `output_rules`, `governance_gates` | language policy、output rules、output gates |
| Prompt Artifacts | `prompt_artifact_templates`, `prompt_composition_rules` | prompt artifact templates 與 composition rules |
| Context | `context_ttl_policy` | context TTL 與 prompt cache layout |
| Budget | `runtime_budget` | token budget |
| Distributed | `distributed_locks`, `multi_agent_coordination`, `async_job_lifecycle` | multi-agent coordination、distributed locks |
| Intelligence | `intelligence_routing` | intelligence routing |
| Decision Recording | `decision_recording` | close-loop tier routing（ADR / session / project） |

## Recovery 來源分層

Runtime recovery 的 machine-readable source 已收斂到 `runtime/runtime.db` canonical documents。Agent 處理 blocking gate、phase drift、stale generated surface 或 recovery retry 時，依下列分層讀取：

| 需求 | 讀取位置 |
| --- | --- |
| 即時 escalation / recovery output | [`../enforcement/escalation-policy.md`](../enforcement/escalation-policy.md) |
| retry limit、strategy change、source reload、validation gate | [`../governance/ai-runtime-governance/recovery-retry-governance.md`](../governance/ai-runtime-governance/recovery-retry-governance.md) |
| domain-specific reload set / forbidden behaviors | [`../metadata/recovery/`](../metadata/recovery/) |
| machine-readable recovery strategy / phase reconciliation / state repair | `runtime.db` recovery tables |
| 修改 runtime recovery 定義 | 更新 `runtime_config_documents`，再執行 `ai-skill runtime compile` refresh projections |

`runtime/runtime.db` 是 committed runtime config 的 canonical copy。不要再保留 committed runtime YAML mirror。

## 主要入口引用

- [`route.runtime.activation-rules`](../knowledge/runtime/routing-registry.yaml:77)
- [`route.runtime.context-ttl`](../knowledge/runtime/routing-registry.yaml:102)
- [`route.runtime.context-loading`](../knowledge/runtime/routing-registry.yaml:161)
- [`route.metadata.knowledge-atom-schema`](../knowledge/runtime/routing-registry.yaml:191)
- [`route.models.model-aware-routing`](../knowledge/runtime/routing-registry.yaml:319)
- [`route.runtime.router-flow`](../knowledge/runtime/routing-registry.yaml:348)
- [`route.runtime.context-ttl-doc`](../knowledge/runtime/routing-registry.yaml:407)
- `gate.checkpoint.capability_discovery_completed` in [`runtime.db`](runtime.db)
- `obligation.checkpoint.run_capability_discovery` in [`runtime.db`](runtime.db)

## 資料庫

Runtime 使用兩種 SQLite database，生命週期不同：

### `runtime.db`（不可變：已提交的 Runtime Registry）

`runtime.db` 是 committed canonical source。Go-native `ai-skill runtime compile` 會從 DB 內的 `runtime_config_documents` 與 deterministic prose mappings refresh projection tables。

所有 committed runtime config 都在 SQLite。若 canonical copy 已在 `runtime.db`，就不能再保留 committed YAML mirror。

> **⚠️ Commit 規則：`runtime.db` 必須包含在 commit 中。**
> 當 canonical runtime config 變更時，`runtime.db` 本身就是 source-of-truth 與 projection output 的 commit artifact。
> 若手動 commit（跳過 hook），**必須**確認 `runtime validate` 通過且 `runtime.db` 已 `git add`。

| Table | Canonical Source | 用途 |
|-------|--------|---------|
| `runtime_config_documents` | `runtime.db` | 保存完整 canonical JSON documents，取代舊 runtime YAML documents |
| `runtime_config_projections` | `runtime.db` | 記錄 canonical documents 投影到哪些 runtime tables |
| `runtime_source_files` | `runtime.db` | 舊 source path 相容 manifest，`source_kind='db'` |
| `phases`, `phase_transitions` | `runtime_config_documents` | phase definitions 與 transition rules |
| `obligations`, `obligation_ledger` | `runtime_config_documents` | 每個 phase 的 duties 與 verification criteria |
| `gates`, `blocking_gates` | `runtime_config_documents` | phase transition prerequisites 與 failure actions |
| `transaction_states`, `transaction_transitions`, `transaction_rules`, `transaction_templates`, `transaction_templates_ext` | `runtime_config_documents` | transaction state、rules、templates |
| `activation_rules`, `activation_rules_mirror`, `core_bootstrap_rules` | `runtime_config_documents` | lazy-load rule definitions 與 activation conditions |
| `discovery_checkpoints`, `discovery_search_strategy`, `capability_checkpoints` | `runtime_config_documents` | phase-aware capability discovery checkpoints |
| `generated_surfaces` | prose sources | 從 workflow、enforcement、governance 文件抽取的 deterministic data |
| `compiler_metadata` | auto-generated | compiler version、compiled timestamp、schema version |
| 其他 runtime config tables | `runtime_config_documents` | budget、TTL、guards、health、routing、pipeline、prompt artifacts、recovery、scheduler、distributed runtime 的 normalized projections |

### `runtime-state.db`（可變：Execution State）

這是 future mutable execution-state surface。舊 Ruby initializer 已刪除，避免誤用。在 Go-native `ai-skill runtime state ...` command 與 validation contract 完成前，不建立也不修改這個 DB。

| Table | 用途 |
|-------|---------|
| `execution_state` | 目前 phase、status、sub-phase、metadata |
| `obligation_status` | 每個 obligation 的 completion tracking |
| `transaction_state` | active transaction lifecycle tracking |
| `execution_log` | append-only event log |

### 查詢範例

```sql
-- 目前在哪個 phase？
SELECT phase, status FROM execution_state ORDER BY id DESC LIMIT 1;

-- 哪些 obligations 還未完成？
SELECT obligation_id, phase FROM obligation_status WHERE status = 'pending';

-- 目前 phase 被哪些 gates 擋住？
SELECT g.name, g.severity FROM gates g
JOIN phases p ON g.phase = p.id
WHERE p.id = (SELECT phase FROM execution_state ORDER BY id DESC LIMIT 1);

-- 目前 state 可以轉移到哪裡？
SELECT to_state, condition FROM transaction_transitions
WHERE from_state = (SELECT state FROM transaction_state ORDER BY id DESC LIMIT 1);
```

## Source-of-Truth 邊界

Runtime 不放概念解釋。runtime design 的 source-of-truth 在：

- `governance/`：design philosophy、lifecycle、validation
- `workflow/`：human-readable execution flows
- `intelligence/`：heuristics、analytical reasoning
- `enforcement/`：executable policy rules
