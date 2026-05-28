# Gen 3 Runtime Trigger Audit & Completion

**Status**: `in-progress` (Phase 1 + Phase 2 完成；Phase 3 起跑中)
**世代**：Gen 3 收斂（不是 Gen 4 forward；是「把第三代真正做完」）
**建立日期**：2026-05-28
**最後更新**：2026-05-28（Open Questions → Resolved Decisions）

> 本 plan 補完 Gen 3 的「completed」定義。當前系統 57 routes / 16 signals / 73 generated_surfaces / 15 commit-msg validators 比例顯示 ~72% routes 與 ~80% projections 沒有 runtime auto-detect 消費者 — 即多數 archived plans 雖標 `completed`，但 strengthened §define_runtime_trigger_flow 規則下實際上是 doc-only graduation。本 plan 系統性 audit 並 wire / 明確標 manual / 移除 orphan，再用機械工具（`ai-skill runtime audit` subcommand）防回流。
>
> 預計排在 [`plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md`](2026-05-27-1557-tool-runtime-signal-economics-integration.md) **之前**執行 — 理由見 §Decision Rationale §Sequencing Rationale。

---

## Decision Rationale

### Problem & Why Now

2026-05-28 governance §`define_runtime_trigger_flow` 強化新增兩條 forbidden：

- `routing_registry_entry_without_discovery_signal_or_commit_validator`
- `sqlite_projection_without_routable_or_validator_consumer`

回溯 audit 揭露既有系統有大量違規 evidence：

| Surface | 規模 | 估計 auto-detected 比例 |
|---|---|---|
| Routes in `knowledge/runtime/routing-registry.yaml` | 57 | ~28% |
| Discovery signals in `runtime/cognitive-modes-discovery.yaml` | 16 | n/a |
| Generated surfaces in `runtime/runtime.db` | 73 | ~20% |
| Commit-msg validators | 15 | n/a |
| Validation scenarios in `validation/scenarios/failure-derived/` | 44 | 估計 < 30% 被 hook 機械引用 |

具體案例 [`plans/archived/2026-05-20-1501-cognitive-state-evidence-governance.md`](../archived/2026-05-20-1501-cognitive-state-evidence-governance.md) 標 `completed` 但實際：

- `route.governance.cognitive-state-evidence` 存在但**無 discovery signal 拉它**
- `enforcement.evidence_hierarchy.contract` 已 projected (status `active`, blocking_level `blocking`) 但**無 validator 消費**
- 5 個 validation scenarios（`intent-stability-drift-v1` / `local-evidence-global-claim-v1` / etc）**無 hook 機械引用**

結果：「completed」是 doc-level 完成，runtime 完全沒觸發。同類問題遍佈系統。

**Why now**：

1. 新規則剛生效（2026-05-28），grandfather flag 必須在規則啟用初期就決定 — 越晚決定越多 plan 受波及
2. Economics integration plan（next forward plan）會新增 11 surfaces + 15 signals — 在 broken foundation 上 scale 只會放大問題
3. 沒有 audit tooling → reviewer 只能用眼睛 spot-check，無法系統性找出所有 orphan

### Decision

執行 **Gen 3 完成定義收斂**：

1. **Inventory & Classify**：對所有 57 routes / 73 surfaces / 44 scenarios 做 4-way 分類：(a) auto-detected via signal, (b) consumed via validator/hook, (c) intentionally manual（workflow / discovery 性質本就需要人為觸發），(d) orphan（應移除或補 wire）
2. **Audit Tool**：`ai-skill runtime audit` subcommand 機械化此分類；輸出 JSON 報告 + 違規清單；接入 `ai-skill runtime validate` 流程作為 warning（非 block，避免阻斷正常工作）
3. **Grandfather Flag**：governance YAML 加 `pre_2026_05_28_doc_only_completion` 註記；既有 archived plans 不溯及，但須在 plans/README.md 表格標 `⚠️ doc-only completion` 而不是 `✅ completed`
4. **Wire High-Priority Orphans**：本 plan 不修全部 41 個 manual routes，只挑 5–10 個 **high-value** 補 signal / validator（criteria: blocking_level=blocking 但無 consumer / 高引用 governance contract / 經 audit tool 判定為高 drift 風險）；其他（intentionally manual）顯式標註
5. **Future-Proofing**：所有未來新 plan 必須在 commit 時通過 `validateRuntimeTriggerWiring`（commit-msg validator，第 16 個）— 若 staged 新增 `route.*` 但無新 signal / validator stage，或新增 projection target_key 但無 consumer，block
6. **Glossary Coverage Guardrail（warning-only）**：本 plan 不做全 repo blocking semantic scan，但 audit / plan-template 必須補一個 glossary impact check。當 staged diff 或 active plan 新增 framework-looking terms（snake_case / owner-layer / runtime-surface style）且不在 `glossary_terms` 或 aliases 中時，產生 candidate warning；掃描範圍包含 `plans/active/`、`architecture/`、`workflow/`、`analysis/`、`intelligence/`、`runtime/`、`ecosystem/`。此 guardrail 只提示「可能需要更新 `knowledge/glossary/ai-skill.md`」，不自動 canonicalize、不阻擋 commit。

### Alternatives Considered

- **A. 完全回流既有 41 個 manual routes 都加 signal / validator**：reject。範圍過大（可能 100+ 行 hook code、20+ 新 signal），且許多 manual route 本意就是 workflow / discovery 性質的「人選載入」，不該強加 signal
- **B. 不修，只標 grandfather flag**：reject。會讓 strengthened 規則對既有系統失效，未來新 plan 也會繼續走相同 doc-only 模式 — 規則就被掏空了
- **C. 把 audit 併進 economics plan 一起做**：reject。Audit 是 Gen 3 收斂工作，economics 是 Gen 4 forward；混在一起會讓 economics 的 12 phases 變 20+ phases，scope 失控；audit tooling 也應該先存在，作為 economics 落地 phase-by-phase 的驗證手段
- **D. 開新 plan，audit-first then wire selectively + future-proof validator**：accept。本 plan 採此方向

### Sequencing Rationale（為何先於 economics plan）

| 維度 | Audit plan（本 plan） | Economics plan |
|---|---|---|
| Scope | 5–7 phases，重點在 inventory + tool + grandfather + 少量 wire | 12 phases，新增 11 surfaces + 15 signals + 多個 new layer |
| 對 Gen 3 / Gen 4 | Gen 3 收斂 | Gen 4 forward |
| 風險方向 | 低（純 audit + 小範圍 wire） | 高（新 ecosystem layer + 新 economics contract） |
| Tooling 產出 | `ai-skill runtime audit` + `validateRuntimeTriggerWiring` | `runtime/economics/` + 多 signal |
| Foundation 依賴 | 無 | 依賴 audit 工具確保新增 surface 不變 orphan |

**結論：先做 audit plan**。理由：

1. Economics plan 會 net-add 11 surfaces — 沒有 audit tooling 防護，這 11 surfaces 可能 5/11 成 orphan，整體 orphan ratio 不降反升
2. Audit plan 產出的 `validateRuntimeTriggerWiring` 會在 economics plan 每個 phase 自動檢查 trigger chain；economics plan 完成度反而會更高
3. Audit plan 規模較小（估 5–7 phases vs 12 phases），先做不會大幅延後 Gen 4 graduation
4. Audit 揭露的 high-priority orphans（例如 `enforcement.evidence_hierarchy.contract` 無 validator）正是 economics plan 也需要消費的 surface — 先修 foundation，economics 才有 stable 對接面

### Why Not an ADR Yet

Audit 是治理工具與既有 surface 的 wiring 補完，沒到不可逆架構決策；未來若 audit tooling 成為 mandatory（每 PR 必跑），再評估 ADR promotion。

### ADR Promotion Criteria（completed 時驗證）

- [ ] `ai-skill runtime audit` subcommand 落地並接入 `runtime validate`
- [ ] 57 routes / 73 surfaces / 44 scenarios 全部分類完成且 JSON inventory 可查
- [ ] Grandfather flag 在 governance YAML 與 plans/README.md 同步生效
- [ ] 至少 5 個 high-priority orphans 補 wire
- [ ] `validateRuntimeTriggerWiring` commit-msg validator 上線並通過 fixture tests
- [ ] Per-surface consumer 表為新 plan 模板的 required section（plans/README.md template 更新）
- [ ] Open Questions 全解

### Consequences

#### 正面

- 「Gen 3 completed」有可驗證定義 — 不再是 doc-only
- 未來新 plan 受 `validateRuntimeTriggerWiring` 保護，自動防回流
- Economics plan 落地時有 stable foundation 與自動驗證工具
- Grandfather flag 給 archived plans 誠實標籤，避免「已完成」誤導

#### 負面

- 本 plan 落地後既有 archived plans 顯示為 `⚠️ doc-only completion`，需要使用者接受此重新分類
- `validateRuntimeTriggerWiring` 可能擋下某些 legitimate edge case（例如純文件更新不需 signal）— 需要 opt-out `[skip-runtime-trigger-wiring]`

#### 風險

| 風險 | 緩解 |
|---|---|
| Audit tool 把 intentionally manual routes 誤標 orphan | Phase 2 分類時 explicit `manual_activation: { reason: ... }` 欄位；workflow 性質 route 一律須 annotate（見 §Resolved Decisions Q4） |
| Wire 高優先 orphans 範圍失控變 41 個都做 | Phase 4 設明確 ≤ 10 上限；其餘留待 organic 補 |
| Validator 阻塞 economics plan 落地 | `validateRuntimeTriggerWiring` 預設 block；提供 `[skip-runtime-trigger-wiring]` opt-out 給 legitimate doc-only / annotation-only 變更 |
| Grandfather flag 變成永久例外 | Sunset deadline 2026-08-31 + 條件式延展至 2026-11-30（見 §Resolved Decisions Q2） |
| Audit executable YAML scenarios drift（文字與 hook 行為脫鉤） | Phase 1 強制每個 scenario 綁定一個 Go fixture test，CI 驗證 YAML ↔ hook 一致 |
| Scenario schema 蔓延（5 個 scenarios 欄位漂移） | Phase 1 鎖定 `validation/scenario.schema.json` 為 single source；新 scenario 必須通過 schema validation |

---

## Runtime Execution Path

### Doc-only Trial 聲明 + Runtime Graduation

**目前狀態（2026-05-28）**：Plan 為 **draft**，全部 Phase 0–7 為 `[ ]`，de facto doc-only。

**Graduation 階梯**：

| Graduation Phase | 達成後生效的 contract 範圍 | Acceptance signal |
|---|---|---|
| **Phase 2 完成** | `ai-skill runtime audit` subcommand 可跑；JSON inventory 輸出 4-way 分類 | `ai-skill runtime audit --json` exit 0 + 報告 routes/surfaces/scenarios 三表 |
| **Phase 5 完成** | `validateRuntimeTriggerWiring` commit-msg validator 上線；高優先 orphan 已補 wire | 第 16 個 commit-msg validator active；test fixture 證明 block + opt-out |
| **Phase 7 完成** | Plan Completion Closure；本 plan 進 archived；plans/README.md 表格使用新 4-way 狀態 enum | ADR Promotion Criteria 全綠 |

**Drift prevention during trial**：

- 本 plan 詞彙統一引用 [`knowledge/glossary/ai-skill.md`](../../knowledge/glossary/ai-skill.md) candidate entries
- Audit JSON schema 在 Phase 2 鎖定；後續變更走 governance §3 規則 8

**明文承認**：本 plan 在 Phase 2 graduation 前不構成 runtime integration；任何「已實作」claim 必須 cite 對應 phase 完成 evidence。符合 [`governance/lifecycle/system-upgrade-governance.yaml`](../../governance/lifecycle/system-upgrade-governance.yaml) §`define_runtime_trigger_flow` 的 doc-only-trial requirements。

### Runtime owner

- Audit tooling：`scripts/ai-skill-cli/internal/audit/` 新 Go package + `ai-skill runtime audit` subcommand
- `validateRuntimeTriggerWiring` validator：`scripts/ai-skill-cli/internal/app/hooks.go`
- Grandfather flag：`governance/lifecycle/system-upgrade-governance.yaml` 新 section
- 4-way 分類規則：本 plan §Inventory & Classification（Phase 2）

### Trigger flow

```
event_or_signal:
  - file_diff matches knowledge/runtime/routing-registry.yaml (new route added)
  - OR file_diff matches runtime/*.yaml with new runtime_projection.target_key
  - OR user runs `ai-skill runtime audit`

detector:
  - commit-msg validateRuntimeTriggerWiring (catches new orphan)
  - OR runtime validate runs audit internally and surfaces warnings

loaded source / contract:
  - Audit JSON inventory (Phase 2 output)
  - Grandfather flag YAML (Phase 3)
  - Per-surface consumer table (in routing-registry / runtime YAML / generated_surfaces metadata)

runtime action / blocking gate:
  - validateRuntimeTriggerWiring block staging new orphan (or warning during trial)
  - runtime audit JSON report emitted on demand
  - audit warning included in runtime validate output

observable evidence:
  - validation/scenarios/failure-derived/orphan-routing-entry-v1.yaml
  - validation/scenarios/failure-derived/orphan-projection-target-key-v1.yaml
  - Go fixture tests in scripts/ai-skill-cli/internal/audit/
```

### Generated surfaces (with named consumer)

| Generated surface key | Named consumer(s) | Consumer 類型 |
|---|---|---|
| `runtime.audit.inventory_contract` | `ai-skill runtime audit` CLI + `ai-skill runtime validate` warning channel | CLI + runtime warning |
| `runtime.audit.classification_rules` | `validateRuntimeTriggerWiring` commit-msg validator + scenario `orphan-routing-entry-v1` | commit-msg validator + scenario |
| `runtime.audit.grandfather_flag` | plans/README.md status table renderer + scenario `pre-2026-grandfather-coverage-v1` | doc rendering + scenario |
| `runtime.audit.glossary_coverage_warning` | `ai-skill runtime audit` warning channel + future plan template glossary impact row | CLI warning + plan review |

### Validation scenarios

- `orphan-routing-entry-v1`：新增 route 無 signal / validator → audit warn / commit-msg block
- `orphan-projection-target-key-v1`：新 projection target_key 無 consumer → audit warn / commit-msg block
- `orphan-scenario-unreferenced-v1`：新增 scenario 無 hook 引用 → audit warn（不 block，scenario 本身有獨立價值）
- `pre-2026-grandfather-coverage-v1`：grandfather flag 範圍與 archived plans 一致 / sunset deadline 已設
- `framework-glossary-candidate-missing-v1`：`plans/active/`、`architecture/`、`workflow/`、`analysis/`、`intelligence/`、`runtime/`、`ecosystem/` 新增 framework-looking term，但 `glossary_terms` / aliases 查不到 → audit warn（不 block，避免短詞 false positive）

### Test passing evidence

- `ai-skill runtime audit --json` 輸出 4-way 分類報告
- `ai-skill runtime validate` 將 audit warnings 納入 checks
- `validateRuntimeTriggerWiring` fixture tests（happy / block / opt-out / warning-only mode）
- Glossary coverage warning fixture：新增一個未收錄 snake_case framework term 時產生 candidate warning；已存在於 `glossary_terms` 或 alias 時不警告
- 既有 archived plans 在 plans/README.md 表格顯示新 4-way status enum

---

## Resolved Decisions

| # | Question | 決定 | 落實位置 |
|---|---|---|---|
| 1 | `validateRuntimeTriggerWiring` 預設 warning 還是 block？ | **預設 block**，無 warning trial 期；legitimate edge case 走 `[skip-runtime-trigger-wiring]` opt-out | Phase 5 |
| 2 | Grandfather flag sunset deadline 設多遠？ | **2026-08-31** 為 primary deadline + 條件式延展條款（若 audit tool age < 60 天或 Phase 4 未完成，自動延至 **2026-11-30**） | Phase 3 |
| 3 | High-priority orphans 是否包含 `enforcement.evidence_hierarchy.contract`？ | **包含**，列為 Phase 4 必補 candidate | Phase 4 |
| 4 | 4-way 分類的「intentionally manual」邊界？ | workflow / discovery 性質的 route **算 intentionally manual**，但必須在 source 加 `manual_activation: { reason: <workflow_discovery|...> }` annotation；缺 annotation 仍判 orphan | Phase 2 |
| 5 | Audit tool 報告格式？ | **md 預設 + `--json` flag**；Go 端共用 inventory struct 雙渲染 | Phase 2 |
| 6 | Scenarios 走 executable YAML？ | **走 `runtime/audit/*.yaml` executable contract**；每個 scenario 必須綁定 Go fixture test；schema 鎖定為 `validation/scenario.schema.json` single source | Phase 1 |
| 7 | Glossary coverage warning term heuristic 多保守？ | 只掃 `backtick-wrapped` terms + `snake_case` ≥ 2 segments，避免自然語言短詞 false positive；自然語言詞彙明確排除 | Phase 2 / Phase 6 |

---

## 完成條件

- [ ] Phase 0–7 全部達成
- [ ] ADR Promotion Criteria 全綠
- [ ] `ai-skill runtime audit` subcommand 上線且接 `ai-skill runtime validate`
- [ ] `validateRuntimeTriggerWiring` 為第 16 個 commit-msg validator，per_commit_obligations 同步
- [ ] 至少 5 個 high-priority orphans 補 wire（含 `enforcement.evidence_hierarchy.contract`）
- [ ] Grandfather flag 在 governance YAML + plans/README.md 同步
- [ ] plans/README.md 表格新 4-way status enum（`✅ completed (auto-detected)` / `⚠️ completed (doc-only / pre-2026-strengthened)` / `🚧 in-progress` / `❌ orphan`）
- [ ] 新 plan 模板加 §Per-surface consumer 表為 required section
- [ ] 新 plan 模板加 §Glossary Impact row：有無新增 framework vocabulary；若有，更新 `knowledge/glossary/ai-skill.md` 或記錄不需新增的理由
- [ ] `ai-skill runtime audit` 能對 `plans/active/`、`architecture/`、`workflow/`、`analysis/`、`intelligence/`、`runtime/`、`ecosystem/` 的未收錄 framework-looking terms 產生 warning-only candidates

---

## Phase 0 — Pre-Build Interrogation

| 欄位 | 內容 |
|---|---|
| Trigger | 使用者要求審視整個系統是否符合新 forbidden 規則、`cognitive-state-evidence-governance` 案例揭露 doc-only completion 問題 |
| Checked sources | `governance/lifecycle/system-upgrade-governance.yaml` / `plans/README.md` / `knowledge/runtime/routing-registry.yaml` / `runtime/cognitive-modes-discovery.yaml` / `runtime/runtime.db` / `validation/scenarios/failure-derived/` / `scripts/ai-skill-cli/internal/app/hooks.go` |
| Goal | 補 Gen 3 完成定義 + 防回流 |
| Scope | Audit + classify + grandfather + wire high-priority + future-proof validator + warning-only glossary coverage guardrail |
| Non-goals | 不重做 41 個 manual routes；不重寫既有 archived plans；不延伸 Gen 4 forward 工作；不做全 repo blocking semantic scan |
| Acceptance | 4-way 分類完整、validator 上線、grandfather 明確、≥ 5 high-priority wired、new plan template 含 glossary impact row |
| Framework discovery | 既有 trigger chain 5 元素（event / detector / source / action / evidence）為驗證標的；audit JSON 為 derived projection；不重新定義 trigger chain |
| Duplication risk | 不重複 routing-registry 既有資料；audit 只 read + classify。Wire 階段補 consumer 而不重新定義 contract |
| Open questions | 全數已解，見 §Resolved Decisions |
| Decision | proceed |

---

## Phase 1 — Test-First Validation

### Tasks

- [x] 新增 `validation/scenarios/failure-derived/orphan-routing-entry-v1.yaml`（commit `be87b63`）
- [x] 新增 `validation/scenarios/failure-derived/orphan-projection-target-key-v1.yaml`（commit `be87b63`）
- [x] 新增 `validation/scenarios/failure-derived/orphan-scenario-unreferenced-v1.yaml`（commit `be87b63`）
- [x] 新增 `validation/scenarios/failure-derived/pre-2026-grandfather-coverage-v1.yaml`（commit `be87b63`）
- [x] 新增 `validation/scenarios/failure-derived/framework-glossary-candidate-missing-v1.yaml`（commit `be87b63`）
- [x] 鎖定 `validation/scenario.schema.json` 為 5 個 scenarios 的 single source；governance YAML §`write_test_first_scenarios` 加 `scenario_schema_single_source` 子句（commit `be87b63`）
- [x] 每個 scenario 綁定 Go fixture test stub（`scripts/ai-skill-cli/internal/audit/scenarios_stub_test.go`，commit `be87b63`）
- [x] `ai-skill runtime refresh` + `runtime validate` 全綠

### Phase 1 完成條件

- [x] 5 個 scenarios 符合 `validation/scenario.schema.json`
- [x] Schema lock 條款寫入 governance YAML（§`write_test_first_scenarios.scenario_schema_single_source`）
- [x] 每個 scenario 對應 fixture test stub 存在
- [x] Runtime validate 通過

---

## Phase 2 — Inventory Tool（Graduation #1）

### Tasks

- [x] 新 Go package `scripts/ai-skill-cli/internal/audit/`：parser for routing-registry / cognitive-modes-discovery / runtime.db generated_surfaces / hooks.go validators / scenarios dir（commit `0f53e91`）
- [x] 定義 4-way classification rules：(a) auto-detected via signal, (b) consumed via validator/hook, (c) intentionally manual（explicit `manual_activation` annotation in source）, (d) orphan
- [x] 加入 glossary coverage warning pass：讀 `glossary_terms` / aliases，掃 7 個路徑下 backtick / snake_case ≥ 2 segments terms；warning-only（commit `642bfe2`）
- [x] 新增 `ai-skill runtime audit` subcommand：預設 markdown，`--json` flag 切換；Go 端共用 inventory struct 雙渲染
- [x] Glossary coverage term heuristic 限制：backtick-wrapped + snake_case ≥ 2 segments；含 `/` 的 path references、單一英文短詞排除
- [x] 接入 `ai-skill runtime validate`：`runtime_audit_warning` warning-only check
- [x] Update CLI docs（command-contract ✓、test-fixture-plan ✓、bdd-scenarios ✓）
- [x] Go tests（24 個全綠）

### Phase 2 完成條件（Graduation #1）

- [x] `ai-skill runtime audit --json` exit 0 + 完整分類報告
- [x] `ai-skill runtime validate` checks 含 audit warnings
- [x] Go tests cover happy / orphan-detected / classification-edge-cases
- [x] Glossary coverage warning tests cover: missing candidate, existing glossary term, alias match, workflow / intelligence / analysis path coverage

**Baseline audit 結果**：routes 55/57 orphan、surfaces 67/73 orphan、scenarios 120/125 orphan、orphan_total=242（驗證 plan §Decision Rationale 對 doc-only 比例的假設）。

---

## Phase 3 — Grandfather Flag

### Tasks

- [ ] `governance/lifecycle/system-upgrade-governance.yaml` 新增 §`pre_2026_05_28_doc_only_completion` section：列出受 grandfather 保護的 archived plans + sunset deadline **2026-08-31** + 條件式延展條款（若 `ai-skill runtime audit` tool age < 60 天 OR Phase 4 未完成，自動延至 **2026-11-30**）+ sunset 後評估規則（剩餘 doc-only items 須升 auto-detected/consumed 或降 orphan 下架）
- [ ] `plans/README.md` 模板更新：plan 狀態 enum 從 `✅ completed` / `🚧 draft` 擴成 4-way（`✅ completed (auto-detected)` / `⚠️ completed (doc-only / pre-2026-strengthened)` / `🚧 in-progress` / `❌ orphan`）
- [ ] 既有 archived plans 中受 grandfather 保護者，標 `⚠️ completed (doc-only)`；以 `cognitive-state-evidence-governance` 為首案
- [ ] Run scenario `pre-2026-grandfather-coverage-v1` 驗證

### Phase 3 完成條件

- [ ] Grandfather YAML active
- [ ] plans/README.md 4-way enum 上線
- [ ] 至少 3 個 archived plans 重新標記

---

## Phase 4 — Wire High-Priority Orphans

### Tasks

- [ ] Run `ai-skill runtime audit` 找 top-10 high-priority orphans（criteria: blocking_level=blocking + 無 consumer / 高引用 governance contract）
- [ ] 為其中 ≥ 5 個補 wire：(a) 加 discovery signal in cognitive-modes-discovery, or (b) 加 commit-msg validator (Go), or (c) 加 routable consumer
- [ ] 必補的 candidate：`enforcement.evidence_hierarchy.contract` — 加 `validateEvidenceHierarchy` 或對應 hook
- [ ] 每個 wire 都有對應 fixture test
- [ ] Re-run audit 確認 wired 的 5 個從 orphan / manual 升為 auto-detected / consumed

### Phase 4 完成條件

- [ ] ≥ 5 high-priority orphans wired
- [ ] Audit 報告中 wired 項目 status 升級
- [ ] Fixture tests 全綠

---

## Phase 5 — Future-Proof Validator（Graduation #2）

### Tasks

- [ ] **新增 `validatePlanCheckboxSync` Go validator**（sibling validator, 第 17 個 commit-msg validator）：當 staged diff 含 source code / test / scenario / generated_surface 變動且 commit message body 引用 `plans/active/<plan>.md`，檢查該 plan 同時 stage 且至少有一個 `[ ]` → `[x]` transition；缺則 emit warning（不 block，加 `[skip-plan-checkbox-sync]` opt-out 給純 hotfix）。理由：避免 agent / 開發者推進 phase 卻忘記翻 checkbox（本 plan 自身 Phase 1 + Phase 2 完成後正是這樣漏掉）
- [ ] 新增 `validateRuntimeTriggerWiring` Go validator 在 `scripts/ai-skill-cli/internal/app/hooks.go`
- [ ] 觸發條件：staged diff includes new `route.*` entry OR new `target_key` in runtime/*.yaml AND no paired discovery signal / validator / intentional manual annotation
- [ ] Opt-out: `[skip-runtime-trigger-wiring]` for legitimate doc-only / refactor / annotation-only changes
- [ ] 註冊 `obligation.commit.runtime_trigger_wiring` 在 `runtime/core-bootstrap.yaml`
- [ ] 更新 `runtime/cli-modification-policy.yaml` 加 `gate.runtime_trigger_wiring_required`
- [ ] 預設 severity: **block**（無 warning trial 期）；legitimate doc-only / annotation-only / refactor 變更走 `[skip-runtime-trigger-wiring]` opt-out trailer
- [ ] Update CLI docs
- [ ] Go fixture tests (happy / block / opt-out / warning-only mode)
- [ ] bin rebuild

### Phase 5 完成條件（Graduation #2）

- [ ] 第 16 個 commit-msg validator `validateRuntimeTriggerWiring` active
- [ ] 第 17 個 commit-msg validator `validatePlanCheckboxSync` active（warning-only）
- [ ] per_commit_obligations 含 `obligation.commit.runtime_trigger_wiring` 與 `obligation.commit.plan_checkbox_sync`
- [ ] cli-modification-policy 新 gate active
- [ ] Fixture tests green（含 block default / opt-out trailer / plan-checkbox transition 案例）

---

## Phase 6 — Plan Template Update

### Tasks

- [ ] `plans/README.md` 模板新增 §Per-surface consumer 表為 required section（仿 economics plan 的 audit-fix consumer 表）
- [ ] `plans/README.md` 模板新增 §Glossary Impact row：新框架詞彙是否引入；是否更新 `knowledge/glossary/ai-skill.md`；若沒有，明確填 `no new framework vocabulary`
- [ ] 加上 `Watch-Out List` citation requirement — 新 plan 須 cite 對應 Gen 4 vision §Watch-Out List 的 wall
- [ ] 更新 `governance/lifecycle/system-upgrade-governance.yaml` linked update

### Phase 6 完成條件

- [ ] plans/README.md 模板更新
- [ ] 後續新 plan 可直接 copy 模板過 governance check

---

## Phase 7 — Close Loop

### Tasks

- [ ] Diff review
- [ ] ReadLints
- [ ] `ai-skill runtime refresh / validate` 全綠
- [ ] `go test ./...` 全綠
- [ ] 更新 `plans/README.md` 狀態為 `✅ completed (auto-detected)`
- [ ] Move plan to `plans/archived/`
- [ ] Commit / push / readback / clean status

### Phase 7 完成條件

- [ ] ADR Promotion Criteria 全綠
- [ ] Plan archived
- [ ] Audit tool + validator + grandfather flag 三者 active 且 fixture covered

---

## Stakeholder 同意項目

- [ ] 接受 4-way plan status enum（含 `⚠️ completed (doc-only)`）
- [ ] 接受既有 archived plans（含 `cognitive-state-evidence-governance`）被重新標 `⚠️ completed (doc-only)`
- [ ] 接受 `ai-skill runtime audit` subcommand 加入 CLI surface
- [ ] 接受 `validateRuntimeTriggerWiring` 為第 16 個 commit-msg validator
- [ ] 接受 glossary coverage 先採 warning-only，不做全 repo blocking semantic scan
- [ ] 接受 grandfather flag 有 sunset deadline（不是永久例外）
- [ ] 接受 audit plan 優先於 economics integration plan 執行

---

## 與其他 plans 的關係

| Plan | 關係 |
|---|---|
| [`plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md`](2026-05-27-1557-tool-runtime-signal-economics-integration.md) | **Downstream beneficiary**。Economics plan 落地時受本 plan 的 `validateRuntimeTriggerWiring` 自動保護；其新增的 11 surfaces + 15 signals 在每個 phase commit 時自動檢查 trigger chain。Sequencing: audit plan FIRST, economics SECOND。 |
| [`plans/archived/2026-05-20-1501-cognitive-state-evidence-governance.md`](../archived/2026-05-20-1501-cognitive-state-evidence-governance.md) | **Audit case study**。`enforcement.evidence_hierarchy.contract` + 5 scenarios + `route.governance.cognitive-state-evidence` 是本 plan §Phase 4 必補 wire 對象 |
| [`plans/archived/2026-05-25-1000-context-language-glossary-system.md`](../archived/2026-05-25-1000-context-language-glossary-system.md) | **Reference plan**。其 Phase 6 「Runtime Auto-Detect Integration」是 wired-correctly 範本；本 plan 的 `validateRuntimeTriggerWiring` 模仿 `validateGlossaryRetroOwn` 結構 |
| [`architecture/ai-native-cognitive-ecosystem-system.md`](../../architecture/ai-native-cognitive-ecosystem-system.md) | **Gen 4 vision**。本 plan 完成 = Gen 3 真正 graduate；Gen 4 才能正當開始 |

---

## 預估規模

| Phase | 變動 | 大致 LOC |
|---|---|---|
| Phase 1 | 5 validation scenarios（含 glossary coverage warning） | ~180 行 YAML |
| Phase 2 | Audit Go package + CLI integration + tests + glossary coverage warning pass | ~700 行 Go |
| Phase 3 | Grandfather YAML + plans/README enum | ~50 行 |
| Phase 4 | ≥ 5 wires（signal / validator / consumer） | ~200–400 行 |
| Phase 5 | validateRuntimeTriggerWiring + policy + docs + tests | ~250 行 |
| Phase 6 | plans/README template + glossary impact row | ~40 行 |
| Phase 7 | Archive | minimal |
| **Total** | | **~1450–1650 行**，4–6 commits |

對比 economics plan 估 ~2500+ 行 / 多 phase / 多 commit — 本 plan 規模約其 50–60%，且做完後 economics plan 落地受保護，整體 ROI 較高。
