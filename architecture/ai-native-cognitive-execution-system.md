# AI-native Cognitive Execution System

**世代**：第三代（current）
**前一代**：[`ai-native-knowledge-operating-system.md`](ai-native-knowledge-operating-system.md)（historical）

本文件是當前世代系統的 **canonical navigation 入口**，不是自包含 spec。新系統的真正 source-of-truth 已分散到 executable YAML contracts、philosophy 文件與 runtime SQLite canonical state；本檔只提供 tour guide。

> 升級治理規則見 [`governance/lifecycle/system-upgrade-governance.md`](../governance/lifecycle/system-upgrade-governance.md)。新世代系統的命名與文件規範由該文件強制。

---

## 系統世代演化

| 世代 | 名稱 | 文件 | 狀態 |
|------|------|------|------|
| 1 | Skill Repository | （未集中文件化） | historical |
| 2 | AI-native Knowledge Operating System | [`ai-native-knowledge-operating-system.md`](ai-native-knowledge-operating-system.md) | historical |
| 3 | **AI-native Cognitive Execution System** | 本檔 | **current** |

完整演化條件與規則見 [`governance/lifecycle/system-upgrade-governance.md`](../governance/lifecycle/system-upgrade-governance.md) §1。

---

## 系統定位

讓 agent 以 **runtime state machine** 驅動知識路由、phase 執行與閉環驗證的認知執行系統。相較第二代「Knowledge OS」聚焦於知識載入與閉環，第三代強調：

1. **可執行契約**：核心流程不再只靠 prose 描述，而由 YAML contracts + SQLite generated surfaces 機器強制
2. **認知執行 runtime**：phase machine / obligation ledger / blocking gates 由 `runtime/runtime.db` 統一管理
3. **多模型協作**：依任務難度與 context cost 在 small / large / specialized profile 之間路由
4. **失效學習自閉環**：agent failure → failure pattern → validation scenario → runtime guard 的完整 pipeline

---

## 當前系統 canonical 入口

| 角色 | 位置 | 說明 |
|------|------|------|
| 啟動契約 | [`CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md) | 3 條必讀核心規則（~800 tokens），其餘 lazy-load |
| OS Layout | [`README.md`](../README.md) | 根 README，列出所有層級與 Quickstart |
| Runtime State | [`runtime/runtime.db`](../runtime/runtime.db) | SQLite canonical state：phase machine、obligation ledger、blocking gates、language policy、output rules、governance gates、generated surfaces |
| 11-Step 知識更新流程 | [`governance/lifecycle/knowledge-update-flow.{md,yaml}`](../governance/lifecycle/knowledge-update-flow.md) | YAML 為 executable contract，Markdown 為 companion 說明 |
| 升級治理 | [`governance/lifecycle/system-upgrade-governance.md`](../governance/lifecycle/system-upgrade-governance.md) | 大型升級的條件、checklist、強制規則 |
| 各層哲學 | [`governance/lifecycle/*-philosophy.md`](../governance/lifecycle/) | router / routing / pipeline / compiler / context-ttl / output-governance / runtime-guards / distributed-runtime / prompt-artifacts / intelligence-routing / capability-discovery |
| Routing Registry | [`knowledge/runtime/routing-registry.yaml`](../knowledge/runtime/routing-registry.yaml) | Task intent → primary_source 的機讀路由 |
| Knowledge Summaries | [`knowledge/summaries/`](../knowledge/summaries/) | Atom 的 compact summary（300-500 tokens） |

---

## 核心機制

每一項只給簡述 + 連往 canonical source；不在本檔重複定義。

### 1. Bootstrap：3 條核心規則 + Lazy-load

舊「Default Bootstrap 12 條」拆成：
- **Core Bootstrap**（3 條）：每 session 必讀，~800 tokens
- **Lazy-load rules**（9 條）：依任務情境 activation

詳見 [`CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md) 與 [`enforcement/README.md`](../enforcement/README.md) 的 activation model。

### 2. Executable YAML Contracts

原則：**YAML 執行，Markdown 說明**。

任一流程文件若含 `ordered_steps` / `activation_conditions` / `required_reads` / `depends_on` / `exit_conditions` / `blocking_gates` / `required_evidence`，**必須**有 companion `.yaml` 並投影到 `runtime.db` 的 `generated_surfaces`。

詳見：
- [`governance/lifecycle/executable-contract-boundary.{md,yaml}`](../governance/lifecycle/executable-contract-boundary.md)
- [`governance/lifecycle/executable-contract-inventory.yaml`](../governance/lifecycle/executable-contract-inventory.yaml)
- 完成計畫：[`plans/archived/2026-05-22-0855-executable-yaml-contract-migration.md`](../plans/archived/2026-05-22-0855-executable-yaml-contract-migration.md)

### 3. Runtime State Machine

`runtime/runtime.db` 是 SQLite canonical，包含：

| 表 | 用途 |
|------|------|
| `phase_machine` / `phases` | 目前 phase、allowed_actions、forbidden_actions |
| `obligation_ledger` / `obligations` | 本 phase 的未完成義務 |
| `blocking_gates` / `gates` | 本 phase 的阻斷閘門 |
| `language_policy` | 語言強制規則（跟隨使用者語言） |
| `output_rules` | 文件輸出規則（去敏、工具中立、格式） |
| `governance_gates` | 輸出品質 blocking gates |
| `generated_surfaces` | 各 prose source 的 compiled surface |
| `runtime_config_documents` | Committed canonical config documents |

升級分析：[`plans/archived/2026-05-15-0920-runtime-execution-layer-upgrade-analysis.md`](../plans/archived/2026-05-15-0920-runtime-execution-layer-upgrade-analysis.md)

### 4. Knowledge Atom + Routing + Summary + Graph

| 層 | 內容 |
|------|------|
| `metadata/schema.md` | Knowledge Atom schema v1 |
| `knowledge/indexes/` | Task intent → canonical source 路由索引 |
| `knowledge/runtime/routing-registry.yaml` | Machine-readable routing |
| `knowledge/summaries/` | 每個 atom 的 compact summary |
| `knowledge/graphs/` | Atom 間關係 graph |

Agent 載入策略：先讀 summary（300-500 tokens），需要才展開 source。

### 5. 11-Step Knowledge Update Flow

從「學到新知識」到「commit/push 完成」的完整 master flow。Sub-pipelines（intelligence-extraction-pipeline、failure-learning-system、linked-updates）不可取代 master。

詳見 [`governance/lifecycle/knowledge-update-flow.yaml`](../governance/lifecycle/knowledge-update-flow.yaml)（executable contract）+ [`.md`](../governance/lifecycle/knowledge-update-flow.md)（rationale）。

### 6. Failure Learning System

Agent failure → classify → contain → promote → strengthen → validate 的閉環：

| 階段 | 位置 |
|------|------|
| Failure capture & taxonomy | [`enforcement/failure-learning-system.md`](../enforcement/failure-learning-system.md) |
| Cross-skill failure patterns | [`enforcement/failure-patterns/`](../enforcement/failure-patterns/) |
| Validation scenarios | [`validation/scenarios/failure-derived/`](../validation/scenarios/failure-derived/) |
| Runtime guards | `runtime/runtime.db` 的 blocking_gates |

### 7. AI Runtime Governance 5-Step

[`governance/ai-runtime-governance/`](../governance/ai-runtime-governance/README.md) 定義 routing / activation / linked-updates / validation-scenario / promotion 五個治理面。

完成計畫：[`plans/archived/2026-05-20-1307-ai-runtime-governance-five-step-integration.md`](../plans/archived/2026-05-20-1307-ai-runtime-governance-five-step-integration.md)

### 8. Cognitive Boundary System

Agent 與 system 間的認知邊界（context window、attention、memory）治理。

完成計畫：[`plans/archived/2026-05-13-0954-cognitive-boundary-system.md`](../plans/archived/2026-05-13-0954-cognitive-boundary-system.md)、[`plans/archived/2026-05-20-1501-cognitive-state-evidence-governance.md`](../plans/archived/2026-05-20-1501-cognitive-state-evidence-governance.md)

### 9. Multi-Model Routing

依任務在 small / large / specialized profile 之間路由，相關 docs：
- [`models/`](../models/README.md) — capability profile、compression、routing
- [`knowledge/runtime/model-context-report.md`](../knowledge/runtime/model-context-report.md)
- [`knowledge/runtime/model-checklists.md`](../knowledge/runtime/model-checklists.md)

完成計畫：[`plans/archived/2026-05-20-1802-model-aware-execution-routing.md`](../plans/archived/2026-05-20-1802-model-aware-execution-routing.md)

### 10. Memory Retrieval & Activation

長期記憶（episodic / project / failure / decision）的 retrieval 與 activation：[`memory/`](../memory/README.md)

完成計畫：[`plans/archived/2026-05-20-1745-memory-retrieval-activation-governance.md`](../plans/archived/2026-05-20-1745-memory-retrieval-activation-governance.md)

### 11. Recovery & Escalation

當執行卡住、源缺失、規則衝突時的恢復流程：

完成計畫：[`plans/archived/2026-05-20-1039-runtime-recovery-escalation-system.md`](../plans/archived/2026-05-20-1039-runtime-recovery-escalation-system.md)

---

## 演化里程碑（依完成順序）

| 時序 | 里程碑 | 計畫 |
|------|------|------|
| 2026-05-11 | 下一階段升級規劃啟動 | [`next-stage-upgrade-plan`](../plans/archived/2026-05-11-1112-next-stage-upgrade-plan.md) |
| 2026-05-11 | APK Analysis 作為 workflow/analysis/intelligence 分離 pilot | [`apk-analysis-pilot-migration`](../plans/archived/2026-05-11-1129-apk-analysis-pilot-migration.md) |
| 2026-05-12 | Context cost 優化（prompt cache alignment） | [`context-cost-optimization`](../plans/archived/2026-05-12-1101-context-cost-optimization.md) |
| 2026-05-12 | Technique → Intelligence pilot | [`technique-intelligence-pilot`](../plans/archived/2026-05-12-1458-technique-intelligence-pilot.md) |
| 2026-05-13 | AI Decision Contract Testing 框架 | [`ai-decision-contract-testing`](../plans/archived/2026-05-13-0837-ai-decision-contract-testing.md) |
| 2026-05-13 | Cognitive Boundary System 整合 | [`cognitive-boundary-system`](../plans/archived/2026-05-13-0954-cognitive-boundary-system.md) |
| 2026-05-13 | Knowledge Runtime Validation Gate | [`knowledge-runtime-validation-gate`](../plans/archived/2026-05-13-1331-knowledge-runtime-validation-gate.md) |
| 2026-05-14 | `shared-rules/` → `enforcement/` 搬遷 | [`shared-rules-to-enforcement-migration`](../plans/archived/2026-05-14-1028-shared-rules-to-enforcement-migration.md) |
| 2026-05-14 | Enforcement layer 強化 | [`enforcement-layer-enhancement`](../plans/archived/2026-05-14-1035-enforcement-layer-enhancement.md) |
| 2026-05-15 | **AI-native Cognitive Execution System 升級分析**（世代命名變更） | [`runtime-execution-layer-upgrade-analysis`](../plans/archived/2026-05-15-0920-runtime-execution-layer-upgrade-analysis.md) |
| 2026-05-15 | Workflow Activation Contract 搬遷 | [`workflow-activation-contract-migration`](../plans/archived/2026-05-15-0949-workflow-activation-contract-migration.md) |
| 2026-05-18 | Software Delivery output templates | [`software-delivery-output-templates`](../plans/archived/2026-05-18-0155-software-delivery-output-templates.md) |
| 2026-05-20 | Runtime Recovery & Escalation 系統 | [`runtime-recovery-escalation-system`](../plans/archived/2026-05-20-1039-runtime-recovery-escalation-system.md) |
| 2026-05-20 | AI Runtime Governance 5-Step 整合 | [`ai-runtime-governance-five-step-integration`](../plans/archived/2026-05-20-1307-ai-runtime-governance-five-step-integration.md) |
| 2026-05-20 | Cognitive State & Evidence Governance | [`cognitive-state-evidence-governance`](../plans/archived/2026-05-20-1501-cognitive-state-evidence-governance.md) |
| 2026-05-20 | DDD Intelligence + Software Delivery Governance | [`ddd-intelligence-software-delivery-governance`](../plans/archived/2026-05-20-1601-ddd-intelligence-software-delivery-governance.md) |
| 2026-05-20 | BDD + DDD Cognition-Aligned Reframe | [`bdd-ddd-cognition-aligned-reframe`](../plans/archived/2026-05-20-1635-bdd-ddd-cognition-aligned-reframe.md) |
| 2026-05-20 | Memory Retrieval & Activation Governance | [`memory-retrieval-activation-governance`](../plans/archived/2026-05-20-1745-memory-retrieval-activation-governance.md) |
| 2026-05-20 | Model-aware Execution Routing | [`model-aware-execution-routing`](../plans/archived/2026-05-20-1802-model-aware-execution-routing.md) |
| 2026-05-21 | Cross-platform Go Script Runtime | [`cross-platform-go-script-runtime`](../plans/archived/2026-05-21-0834-cross-platform-go-script-runtime.md) |
| 2026-05-22 | **Executable YAML Contract Migration**（最近一次主要升級） | [`executable-yaml-contract-migration`](../plans/archived/2026-05-22-0855-executable-yaml-contract-migration.md) |

---

## 與第二代的關係

| 第二代（Knowledge OS） | 第三代（Cognitive Execution） |
|------|------|
| `skills/` 為 capability 層 | 已 deprecated；遷移到 `workflow/` + `analysis/` + `intelligence/` |
| Default Bootstrap 12 條全讀 | Core Bootstrap 3 條 + Lazy-load 9 條 |
| `shared-rules/` 共用規則 | 搬遷到 `enforcement/` |
| Prose 為 source-of-truth | Executable YAML 為 source-of-truth，Markdown 為 companion |
| Reference-first + 工具相容層 | 仍為 reference-first，但加入 runtime state machine 與 generated surfaces |
| 閉環驗證靠 prose checklist | 由 `runtime.db` `blocking_gates` 與 validation scenarios 機器強制 |
| 失效靠 feedback lesson 累積 | Failure pattern → validation scenario → runtime guard 的完整 pipeline |

不變的原則：
- `reference-first` 仍是預設
- Canonical source 為 `<AI_SKILL_REPO>` git repository
- Path placeholder（`<AI_SKILL_REPO>` / `<PROJECT_ROOT>` / `<WORKSPACE>`）
- Tool mirror / bundle / copy snapshot 為相容層，非 source

---

## 與其他層的邊界

| 本檔 | 不是 |
|------|------|
| `architecture/` 是「**世代 canonical 入口**」 | 不是執行規則正文（在 `enforcement/`）、不是工程判斷正文（在 `intelligence/`）、不是執行計畫（在 `plans/`） |
| 本檔指向 canonical source | 不複製 canonical source 正文；source 變更時本檔可能不需動，除非世代結構變動 |

`architecture/` 與 `intelligence/engineering/architecture/` 的邊界：

| 層 | 範疇 |
|----|------|
| `architecture/`（本層） | **OS / 知識庫**架構：repo 怎麼組織、啟動、契約、世代演化 |
| `intelligence/engineering/architecture/` | **工程**架構判斷：domain modeling、modularity、coupling、選型 trade-off |

---

## 維護規則

依 [`governance/lifecycle/system-upgrade-governance.md`](../governance/lifecycle/system-upgrade-governance.md)：

- 系統升級若涉及**世代命名變更**，必須在 `architecture/` 建立新世代 canonical 文件，並把前一代標為 historical
- 本檔內容若與 canonical source（README / CORE_BOOTSTRAP / runtime.db / philosophy / archived plans）不一致，**canonical source 優先**；本檔只是 navigation
- 新增完成計畫到 `plans/archived/` 時，視重要性決定是否加入 §「演化里程碑」表
