# Executable YAML Contract Migration Plan

> **狀態**：draft
> **建立時間**：2026-05-22 08:55 JST
> **目的**：把目前仍主要依賴 Markdown prose 的流程、gate、required reads、failure actions 與 final status requirements，整理成 owner-layer YAML executable contracts，讓 ChatGPT 以外的 agent 也能用結構化欄位遵守規則，並透過 `runtime/runtime.db` 的 generated surfaces 驗證是否同步。

## 背景

目前系統的方向是正確的：`governance/lifecycle/executable-contract-boundary.yaml` 已定義「Markdown explains, YAML executes」，`runtime/runtime.db` 也已投影部分 generated surfaces。

真正的缺口不是「沒有 YAML」，而是 YAML 分成三種，但對 agent 的強制力不同：

| 類型 | 目前例子 | 問題 |
| --- | --- | --- |
| Metadata YAML | `metadata/rules/*.yaml` | 描述 rule 的摘要、priority、activation，但多數沒有 `runtime_projection.enabled: true`，不是可執行流程契約。 |
| Prose-derived generated surface | `generated_surfaces` 中的 `enforcement/dependency-reading.md`、`workflow/*/execution-flow.md` | compiler 可抽取部分 deterministic 資料，但 agent 仍可能先讀 Markdown 而漏掉 ordered contract。 |
| Executable YAML contract | `governance/lifecycle/knowledge-update-flow.yaml`、`ai-tools/agent-onboarding.yaml` | 結構清楚、可投影 runtime，但覆蓋範圍還不夠。 |

因此其他 agent 不照規定跑，主要原因是：P0/P1 規則仍有太多「必須」藏在 Markdown 段落或表格裡，沒有統一的 `activation -> required_sources -> steps -> gates -> failure_modes -> success_criteria` contract。

## Architecture Compatibility Preflight

| 欄位 | 內容 |
| --- | --- |
| Trigger | 建立 YAML contract migration plan，尚未進入 implementation phase |
| Checked sources | `CORE_BOOTSTRAP.md`、`README.md`、`ai-tools/agent/codex.md`、`runtime/README.md`、`runtime/runtime.db`、`scripts/ai-skill-cli/internal/app/runtime_compiler.go`、`governance/lifecycle/executable-contract-boundary.md`、`governance/lifecycle/executable-contract-boundary.yaml`、`governance/lifecycle/compiler-philosophy.md`、`enforcement/dependency-reading.md`、`enforcement/linked-updates.md`、`plans/README.md`、`workflow/documentation/README.md`、`workflow/documentation/execution-flow.md`、`metadata/rules/*.yaml` P0/P1 samples |
| Current architecture | runtime internal config canonical source 是 `runtime/runtime.db`；governance / enforcement / workflow / ai-tools 的 executable contract 必須留在 owner layer；影響 agent execution 的 YAML contract 需設定 `runtime_projection.enabled: true` 並投影到 `generated_surfaces`。 |
| Conflicts | 不得新增 `runtime/**/*.yaml` mirror；不得把 owner-layer policy 搬進 runtime；metadata YAML 不能被誤當成完整 executable contract。 |
| Decision | proceed with planning only；implementation 前先補 contract schema、validator、content assertions 與 inventory drift 修正，再分批新增 YAML。現行 compiler 已會掃描 owner-layer YAML 並投影 `runtime_projection.enabled: true` 的 contract，Phase 0 不應假設 compiler coverage 從零開始。 |
| Validation | 每批 contract 新增後執行 runtime compile / refresh / validate，查詢 `generated_surfaces` target key 與 `data` 內容關鍵欄位，並新增或更新 validation scenario。 |

## Review Findings To Carry Forward

- `governance/lifecycle/executable-contract-boundary.yaml` 已列 `enforcement/evidence-hierarchy.md` 為 P1 candidate，本計畫原 inventory 漏列，需補入避免 candidate drift。
- `directory-structure-governance` 在 boundary YAML 是 P2，本計畫列 P1；`sanitization` 在 metadata 是 P1，本計畫列 P0；`conversation-goal-ledger` 在 metadata 是 P0，本計畫列 P1。若保留調整，必須在 inventory 或 Phase 0 寫明 priority decision reason。
- `scripts/ai-skill-cli/internal/app/runtime_compiler.go` 已有 `compileExecutableYAMLContracts`，會掃描 `governance/`、`enforcement/`、`workflow/`、`ai-tools/`、`metadata/rules/` 並投影 enabled YAML。Phase 0 重點應改為 schema validation、contract completeness validation、content assertion 與 scenario coverage。
- 目前 `runtime/runtime.db generated_surfaces` 中 synced 的 executable YAML contract 只有 `ai-tools/agent-onboarding.yaml`、`ai-tools/new-project-onboarding.yaml`、`governance/lifecycle/decision-promotion-pipeline.yaml`、`governance/lifecycle/executable-contract-boundary.yaml`、`governance/lifecycle/knowledge-update-flow.yaml`；P0/P1 enforcement 與 workflow contracts 尚未覆蓋。
- `workflow/documentation/README.md` 已宣告本 workflow 保存「在業務專案或其他 repository 撰寫、整理、分類文件」的可執行步驟與分類表，因此 YAML 化判斷標準也必須能套用到其他專案的 docs / runbook / ADR / workflow 文件，而不是只服務 Ai-skill 本庫 migration。

## 目標

1. 建立一個統一的 executable YAML contract shape，讓不同 agent 能先讀 YAML 再讀 Markdown。
2. 將 P0/P1 blocking 規則優先 YAML 化，尤其是會導致漏讀、漏驗證、漏 close-loop、越權、錯 layer 的規則。
3. 區分 metadata YAML 與 executable contract，避免「已有 `metadata/rules/*.yaml`」被誤判為已完成流程契約。
4. 讓 runtime compiler 投影所有 execution-affecting YAML 到 `runtime.db.generated_surfaces`。
5. 為「其他 agent 不遵守流程」建立可重放 validation scenarios，而不是只靠提示詞提醒。
6. 將 YAML 化判斷標準接入 `workflow/documentation/`，讓 agent 在其他專案寫文件時，也能判斷該文件是一般 Markdown、project metadata/front-matter，還是需要 companion YAML contract。

## 非目標

- 不把所有 Markdown 都改成 YAML；哲學、背景、tradeoff、heuristic 仍留在 Markdown / intelligence。
- 不新增 `runtime/**/*.yaml` committed mirror。
- 不在第一批改變 runtime phase machine 的 canonical SQLite source。
- 不把所有規則升級成 blocking gate；只有控制安全、source-of-truth、validation、close-loop 或 routing 的項目需要阻斷。

## Contract Shape

所有新增 executable contract 建議使用同一組欄位：

```yaml
schema_version: executable-contract/v1
id: <layer.domain.contract>
title: <human title>
owner_layer: <governance|enforcement|workflow|ai-tools>
source_markdown: <path/to/source.md>
status: active
contract_type: <policy_gate|workflow_flow|onboarding_flow|promotion_gate>
blocking_level: <blocking|advisory|informational>

runtime_projection:
  enabled: true
  target_key: <stable.generated_surface.key>
  surface_type: executable_yaml_contract
  execution_scope: <agent_behavior|repo_writeback|workflow_execution>

activation:
  events: []
  required_when: []

steps: []
required_sources: []
depends_on: []
gates: []
required_evidence: []
success_criteria: []
failure_modes: []
validation:
  runtime_assertions: []
  scenario_refs: []
final_status_report: []
```

## YAMLization Decision Standard

同一套標準同時適用於 system governance、enforcement、tool adapter、`workflow/*` 文件，以及 `workflow/documentation/` 指導下產生的其他專案文件。判斷順序如下：

| 判斷 | YAML 化決策 | 理由 |
| --- | --- | --- |
| 文件要求 agent 依序執行多個步驟，且漏步會導致錯 layer、漏驗證、漏 close-loop、越權或錯誤輸出 | 必須建立 executable YAML contract | Ordered flow 是 agent 行為，不應只藏在 prose。 |
| 文件定義 activation / required reads / depends_on / blocking gates / failure action / required evidence / final report | 必須建立 executable YAML contract | 這些欄位可機讀、可驗證，且可投影到 runtime。 |
| 文件只提供哲學、背景、tradeoff、設計理由或人類導讀索引 | 保持 Markdown-only，除非後續抽出 gate | 避免把 judgment 或長篇解釋塞進 runtime 噪音。 |
| 文件已有 metadata YAML，但沒有 `contract_type`、`runtime_projection.enabled: true`、steps / gates / success criteria | 不算 executable contract | Metadata YAML 只代表 rule 摘要與 routing，不代表 agent 可直接執行。 |
| workflow 文件只描述分類智慧，但沒有固定執行順序或阻斷條件 | 保持 Markdown / intelligence，必要時由 workflow contract 引用 | 判斷智慧留在 intelligence / governance，workflow YAML 只放操作契約。 |
| workflow 文件包含「何時讀、寫入前分類、文件形狀、驗證與完成訊號」 | 建立 workflow companion YAML | `workflow/documentation/execution-flow.md` 屬此類，應作為 documentation YAML pilot。 |
| 其他專案的文件只是說明、決策背景、操作備忘或人類閱讀材料 | 使用 Markdown，必要時用專案既有 front-matter / 父 README 索引 | 不把所有 project docs 都變成 executable contract，避免製造維護負擔。 |
| 其他專案的文件定義 agent 要反覆執行的流程、檢查 gate、required evidence、failure action 或 release / runbook 步驟 | 依專案慣例建立 companion YAML 或等價 structured contract | 讓 project-local agent 也能用同一套 execution standard，不只依靠 prose。 |

判斷結果必須落入三種狀態之一：`contract_required`、`markdown_only`、`not_applicable`。若判為 `contract_required` 但本輪不建立 YAML，需在 plan / TODO / linked update 中記錄缺口與理由。

## Candidate Inventory

| Priority | Source | Target contract | Why YAML | Current state |
| --- | --- | --- | --- | --- |
| P0 | `enforcement/dependency-reading.md` | `enforcement/dependency-reading.yaml` or upgraded `metadata/rules/dependency-reading.yaml` | Required reads、dependency ledger、source-of-truth miss、writeback transaction、final status 都是 agent 行為 gate。 | Metadata YAML exists；generated surface exists；缺完整 executable contract shape。 |
| P0 | `enforcement/authorization-scope.md` | `enforcement/authorization-scope.yaml` or upgraded metadata contract | Authorization / scope 是安全阻斷條件，不能只靠 prose。 | Metadata YAML exists；缺 runtime projection。 |
| P0 | `enforcement/sanitization.md` | `enforcement/sanitization.yaml` or upgraded metadata contract | Secrets、本機路徑、私人 evidence 是 commit blocker。 | generated surface exists；需補 explicit failure modes / evidence fields。 |
| P1 | `enforcement/linked-updates.md` | `enforcement/linked-updates.yaml` or upgraded metadata contract | 多檔同步、runtime compile、commit / push / readback 是 close-loop gate。 | Metadata YAML exists；缺完整 linked-update matrix contract。 |
| P1 | `enforcement/goal-action-validation.md` | `enforcement/goal-action-validation.yaml` or upgraded metadata contract | 重要工作要有 goal / action / validation，否則 agent 容易跳到 implementation。 | Metadata YAML exists；generated surface exists；priority metadata 目前是 P2，需重新評估 gate priority。 |
| P1 | `enforcement/evidence-hierarchy.md` | `enforcement/evidence-hierarchy.yaml` or upgraded metadata contract | Evidence precedence、confidence、escalation conditions 會影響 agent 是否能宣稱驗證完成。 | Listed as P1 candidate in executable-contract-boundary.yaml；本計畫需補齊避免 inventory drift。 |
| P1 | `enforcement/conversation-goal-ledger.md` | `enforcement/conversation-goal-ledger.yaml` or upgraded metadata contract | 多步驟任務、lock、handoff、completion deletion gate 需要結構化。 | Metadata YAML exists；routing registry references source。 |
| P1 | `governance/lifecycle/directory-structure-governance.md` | `governance/lifecycle/directory-structure-governance.yaml` | 新增 / rename 目錄會影響 layer boundary 與 routing。 | Listed as candidate in executable-contract-boundary.yaml。 |
| P1 | `governance/ai-runtime-governance/linked-update-governance.md` | `governance/ai-runtime-governance/linked-update-governance.yaml` | linked update completeness 是 governance gate，應與 enforcement linked-update table 對齊。 | Markdown governance only。 |
| P1 | `governance/ai-runtime-governance/validation-scenario-governance.md` | `governance/ai-runtime-governance/validation-scenario-governance.yaml` | failure promotion 成 scenario 的 gates 需 machine-readable。 | Routing registry references source；缺 contract。 |
| P2 | `enforcement/decision-efficiency.md` | `enforcement/decision-efficiency.yaml` or upgraded metadata contract | 決定何時少讀 / 多讀 / escalation，會影響 context loading。 | Metadata YAML exists；缺 runtime projection。 |
| P2 | `workflow/*/execution-flow.md` | `workflow/<domain>/execution-flow.yaml` | 現在多靠 prose-derived generated surface；應補 companion YAML 給 agent 直接執行。 | generated surfaces exist；缺 owner-layer YAML。 |
| P2 | `workflow/*/artifact-gates.md` | `workflow/<domain>/artifact-gates.yaml` | artifact completeness、required evidence、exit criteria 需要結構化。 | generated surfaces exist；缺 owner-layer YAML。 |
| P2 | `ai-tools/agent/*.md` | `ai-tools/agent/<tool>.yaml` | 不同 agent adapter bootstrap 應有統一 required reads / forbidden duplication contract。 | `agent-onboarding.yaml` exists；per-tool YAML missing。 |

### Priority Normalization Notes

- `sanitization` 可從 metadata P1 升為本計畫 P0，理由是 secrets、本機路徑、私人 evidence 會直接成為 commit blocker；Phase 0 必須確認 `enforcement/sanitization.md` 的實際 blocking gate 是否足以支持 P0。
- `conversation-goal-ledger` metadata 是 P0，但本計畫可放 Phase 2，理由是它是 Core Bootstrap preload，但第一批 migration 先處理 repo writeback / safety blockers；不得把 Phase 2 誤解為 priority 降級。
- `directory-structure-governance` 在 boundary YAML 是 P2；若本計畫保留 P1，需在 Phase 4 前確認新增 / rename 目錄 gate 是否已變成 cross-agent blocking risk。
- `goal-action-validation` metadata 是 P2；若升為 P1，需明確限定為「重要工作單元或可重用文件變更」的 blocking gate，不把所有小回覆都升級為阻斷。

## Phase Plan

### Phase 0：Schema, Validator & Inventory Contract

- [ ] 定義 executable contract schema，放在 governance 或 metadata schema owner layer。
- [ ] 決定 enforcement contract 是新增 `enforcement/*.yaml`，還是升級既有 `metadata/rules/*.yaml`；預設先用 `enforcement/*.yaml`，除非 schema 能清楚區分 metadata vs executable contract。
- [ ] 確認 runtime compiler 既有 owner-layer YAML 投影行為，補 validator / tests，而不是重複建立平行 compiler path。
- [ ] 新增 schema completeness validation：enabled contract 必須有 `schema_version`、`contract_type`、`blocking_level`、`runtime_projection.target_key`、`activation`、`required_sources`、`gates` 或明確 not-applicable reason。
- [ ] 新增 content assertion validation：SQLite `generated_surfaces.data` 必須包含 contract 的關鍵 gates / failure_modes / final_status_report，不只檢查 target key synced。
- [ ] 依 `YAMLization Decision Standard` 盤點 system governance 與 workflow 文件，輸出 `contract_required` / `markdown_only` / `not_applicable`。
- [ ] 新增 validation scenario：metadata YAML 不等於 executable contract，agent 必須先讀 contract。

### Phase 1：P0 / P1 Enforcement Contracts

- [ ] YAML 化 `dependency-reading` 的 required read ledger、source-of-truth miss escalation、writeback final status gate。
- [ ] YAML 化 `authorization-scope` 的 activation、allowed / forbidden data collection、failure action。
- [ ] YAML 化 `sanitization` 的 blocker fields、redaction criteria、commit gate。
- [ ] YAML 化 `linked-updates` 的 changed path matrix、required linked checks、runtime compile / readback / dirty status gates。
- [ ] YAML 化 `goal-action-validation` 的 goal / action / validation closure，並重新評估是否應為 P1 gate。
- [ ] YAML 化 `evidence-hierarchy` 的 evidence precedence、confidence threshold、escalation condition。

### Phase 2：Conversation / Tool / Routing Contracts

- [ ] YAML 化 `conversation-goal-ledger` 的 goal lifecycle、lock decision、completion deletion gate。
- [ ] 為 `ai-tools/agent/codex.md`、`claude.md`、`cursor.md`、`roo.md` 建立 adapter contract，統一 bootstrap required reads 與 forbidden duplicated rules。
- [ ] 將 `knowledge/runtime/routing-registry.yaml` 的 route activation 與 source-of-truth gate 檢查接入 validation。

### Phase 3：Workflow Contracts

- [ ] 為 `workflow/software-delivery/` 建立 `execution-flow.yaml` 與 `artifact-gates.yaml`。
- [ ] 先用 `workflow/documentation/` 驗證 YAMLization Decision Standard：`execution-flow.md` 含讀者 / 生命週期、分類維度、檔案形狀、驗證與完成訊號，應建立 `workflow/documentation/execution-flow.yaml` companion contract。
- [ ] 更新 `workflow/documentation/execution-flow.md` 的文件化流程，加入「其他專案何時只寫 Markdown / front-matter，何時需要 companion YAML 或等價 structured contract」的判斷步驟。
- [ ] 更新 `workflow/documentation/README.md`，說明 YAMLization Decision Standard 是 project documentation standard 的一部分，可套用於其他 repository 的 docs / runbook / ADR / workflow 文件。
- [ ] 為 `workflow/apk-analysis/`、`workflow/travel-planning/`、`workflow/greenfield/` 建立 companion YAML；若某 workflow 只保留 philosophy / routing index，明確標 `markdown_only`。
- [ ] 讓 workflow YAML 明確列出 activation、required sources、blocking gates、success criteria、final report fields。

### Phase 4：Governance Promotion Contracts

- [ ] YAML 化 directory structure governance。
- [ ] YAML 化 linked update governance。
- [ ] YAML 化 validation scenario governance。
- [ ] 對照 `executable-contract-boundary.yaml`，將 candidates 移到 `contract_exists` 或明確標記 markdown-only。

### Phase 5：Validation & Close Loop

- [ ] 每批新增 contract 後執行 runtime compile / refresh / validate。
- [ ] 查詢 `runtime/runtime.db generated_surfaces`，確認 target key synced。
- [ ] 新增 failure-derived scenarios：agent 只讀 Markdown、不讀 YAML、把 metadata YAML 當 executable contract、漏跑 close-loop。
- [ ] 更新 `knowledge/runtime/model-checklists.md` 或 routing report，讓 small / weaker agents 有 checklist-first path。
- [ ] 完成後執行 Plan Completion Closure，更新 `plans/README.md` 並移至 `archived/`。

## Implementation Order

建議先做小而硬的順序：

1. `dependency-reading` + `linked-updates`：修掉最多 agent 漏閉環問題。
2. `authorization-scope` + `sanitization`：修掉 P0 safety / commit blocker。
3. `goal-action-validation` + `conversation-goal-ledger`：修掉做到一半、驗證不閉環、handoff 漏狀態。
4. `workflow/documentation`：先驗證 YAMLization Decision Standard 是否能判斷 workflow 需不需要 YAML 化。
5. `workflow/software-delivery`：作為 artifact gate / execution-flow 的第二個 workflow pilot。
6. 其他 workflow 與 per-tool adapter contract。

## Open Decisions

| Decision | Options | Default |
| --- | --- | --- |
| Enforcement contract location | 新增 `enforcement/*.yaml`；或升級 `metadata/rules/*.yaml` | 偏向新增 `enforcement/*.yaml`，避免 metadata summary 被弱模型誤當 executable contract；只有 schema 與 validator 清楚區分時才升級 metadata YAML。 |
| Runtime projection target key naming | `enforcement.<rule>.contract`；或沿用 metadata id | 使用 `enforcement.<rule>.contract`，避免與 metadata rule id 混淆。 |
| Workflow YAML source | 手寫 companion YAML；或 compiler 從 Markdown 產生 YAML | 第一批手寫 contract，之後再評估 deterministic generation。 |
| Blocking level | 全部 blocking；或 P0/P1 blocking、P2 advisory | P0/P1 blocking，P2 先 advisory + validation warning。 |
| YAMLization decision output | 只在 plan 記錄；或建立可驗證 registry / scenario | Phase 0 先在 plan 與 validation scenario 記錄，後續若穩定再 promotion 成 governance contract。 |
| Downstream project format | 強制所有專案使用 Ai-skill YAML shape；或允許 project-local front-matter / structured contract mapping | 不強制所有專案照抄 Ai-skill schema；`workflow/documentation/` 定義判斷標準與欄位語意，具體格式可映射到專案慣例。 |

## Completion Criteria

- P0 / P1 candidate sources 都有 executable contract 或明確 markdown-only / not-applicable decision。
- 所有 execution-affecting contracts 都有 `runtime_projection.enabled: true`。
- `runtime/runtime.db generated_surfaces` 中可查到 synced target key。
- SQLite `generated_surfaces.data` content assertion 能確認關鍵 gates、failure modes、final status report 欄位已進入投影內容。
- `workflow/documentation/` 已套用 YAMLization Decision Standard，並把「其他專案文件何時需要 companion YAML / structured contract」納入 documentation workflow。
- 至少一組 validation scenario 能重放「agent 漏讀 YAML contract」並失敗。
- `plans/README.md` 已更新狀態；完成後 plan 移至 `archived/`。

