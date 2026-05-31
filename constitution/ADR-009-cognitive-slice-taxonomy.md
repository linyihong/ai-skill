# ADR-009: Cognitive Slice Taxonomy

## Status

**Accepted**

## Framework Generation

- **世代分類**：Gen 3 governance promotion（既有 framework 升 ADR，無新功能）
- **當前世代文件**：[`architecture/ai-native-cognitive-execution-system.md`](../architecture/ai-native-cognitive-execution-system.md) §workflow / §analysis loading boundary
- **適用狀態**：cognitive slice taxonomy 是 Gen 3 runtime hardening 的 cross-cutting framework；本 ADR 鎖定「slice 為 routable cognition surface 的最小單位、三層邊界 enforceable、placement predicate falsifiable」這些跨 plan 決定，不取代 ADR-003 三層架構（workflow / analysis / intelligence），而是把該架構的**載入粒度**寫成可機械驗證的契約。

## Date

2026-05-31

## Source Plans

- [`plans/archived/2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md`](../plans/archived/2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md) — framework + first pilot（software-delivery 6 lifecycle slices + greenfield no-split 反向證據）+ Phase 4 5 個 scenario + Phase 4 Extension 3 個 domain probe
- [`plans/archived/2026-05-30-2200-apk-analysis-artifact-gates-decomposition.md`](../plans/archived/2026-05-30-2200-apk-analysis-artifact-gates-decomposition.md) — second pilot（apk-analysis 8 slices 含 Scheme A vs B probe override）+ Phase 4 3 個 acceptance scenario + 誠實 acceptance retrospective

完整 Decision Rationale、Open Questions、Phase 設計、Acceptance、Trial 證據與 ADR Promotion Criteria 評估皆保留於兩個 source plan。本 ADR 為 promotion plan [`2026-05-31-2200-cognitive-slice-taxonomy-adr-promotion.md`](../plans/active/2026-05-31-2200-cognitive-slice-taxonomy-adr-promotion.md) Phase 1 完成、parent plan 全 8 ADR Promotion Criteria PASS 後的 accepted promotion。

## Context

ADR-003 確立 `workflow/` / `analysis/` / `intelligence/` 三層架構但**未規範載入粒度**。實務上發現：

- 巨型 `workflow/<domain>/*.md`（例如 `workflow/apk-analysis/artifact-gates.md` 575 行 15 個 gate）即使位於正確層級，agent 為小任務也被迫載入整份 → execution overload
- `analysis/` 與 `intelligence/` 的邊界容易被 heuristic「看起來像 pattern」就升 intelligence → premature promotion contamination
- routing-registry 只能路由到「檔案」，無法路由到「檔案內某個 cognitive phase」
- 多 session 累積 summary 後，巨型 surface 的具體可執行細節被反覆壓縮成抽象口號

這不是「重寫三層架構」的問題，而是**三層架構需要更細的載入粒度單位**。

Parent plans 的 2 個切分 pilot + 1 個 no-split 反向證據 + 1 個 scheme probe override 證實：用「**cognitive slice = routable cognition surface 的最小單位**」可同時解決上述四個問題，且 falsifiable membership predicate 讓 misplacement 可偵測、可逆。

## Decision

採用 [`governance/cognitive-slice-taxonomy.md`](../governance/cognitive-slice-taxonomy.md) 為 cognitive slice framework 的 **canonical living spec**，並將下列 6 條核心規則 promote 為 architectural canon：

| ADR-009 條 | 對應 governance § | 核心約束 |
|---|---|---|
| **§1 Slice schema** | taxonomy §1（14 欄） | 每個 slice 必含 `id` / `purpose` / `type` / `tags` / `load_when` / `do_not_load_when` / `owner_layer` / `layer_justification` / `evidence_refs`（intelligence ≥2）/ `canonical_source` / `dependencies` / `dependency_budget` / `summary_path` / `validation_signal`。Schema 是機械可檢查的 slice 定義單位 |
| **§2 Type+tags 收斂** | taxonomy §2 | Primary `type` **僅** 4 種：`execution` / `evidence` / `examples` / `failure`。其餘責任（artifact-gate / closure / observation-triage / tool-procedure / extraction-to-intelligence …）一律降為 `tags`。新增第 5 個 primary type 需回 plan 重評。防 taxonomy explosion |
| **§3 Granularity rule** | taxonomy §3 | Slice 最小單位 = **能獨立完成 1 個 cognitive phase**，非 step、非 concept。判準：載入後 agent 能完成自足認知階段而無需瘋狂 cross-reference。防 over-fragmentation |
| **§4 三層邊界 + falsifiable placement predicate** | taxonomy §4 | `workflow` = 順序；`analysis` = 證據取得與驗證；`intelligence` = 為何長期有效/失敗。歸層必須通過機械 membership test：(a) **analysis_membership_test** = task-instance observation；(b) **intelligence_membership_test** = generalization **且** `evidence_refs` ≥ 2 個獨立、可解析來源。少於 2 → premature_promotion，強制 fallback 至 analysis |
| **§5 Examples suppression** | taxonomy §5 | `type: examples` 預設 `default_load: false`，只在 user 明確要求 / 偵測到 ambiguity 才載入。防 example-driven contamination（對應 Watch-Out Wall 5 positive-activation bias） |
| **§6 Naming + glossary** | taxonomy §6 + `knowledge/glossary/ai-skill.md` `cognitive_slice` entry | Canonical term：`cognitive_slice`。Operational synonyms：`execution surface` / `evidence surface` / `loading surface`。governance / 內部設計用 `slice`；對外 / runtime-oriented 文件用 `surface` |

### 約束邊界（與 ADR-003 / ADR-007 的關係）

- **不取代 ADR-003**：三層架構（workflow / analysis / intelligence）仍由 ADR-003 鎖定；本 ADR 規範該三層的**載入單位粒度**
- **遵循 ADR-007 No-Proposed-ADR Rule**：本 ADR 為 plan-completed-after-promotion accepted entry；無 proposed status
- **不引入新 runtime primitive**：sligce schema 寫進 routing-registry `loading_surfaces`（既有 hierarchical routing 機制），不新增 phase_machine state 或 commit-msg validator class

### Governance ownership boundary

- **ADR-009 鎖**：上述 6 條 canonical rule 的存在 + 機械可驗證性
- **governance/cognitive-slice-taxonomy.md 維護**：schema 細節、pilot section（§7 software-delivery / §7.5 apk-analysis）、Phase 4 fixture template、規則新增/微調

修改 6 條 canonical rule **須走 ADR supersede 流程**；修改 schema 欄位、pilot 範例、新增 tag 在 governance/ 即可。

## Consequences

### 正面

- Cognitive slice framework 從 governance proposal 升為 architectural canon；其他 plan 直接引用 ADR-009 而非 plan-specific governance doc
- 三層邊界從「honor-system 標籤」變「機械可檢查 predicate」——misplacement 可偵測、可逆
- Routing 從「routes_to file」延伸為「routes_to slice」，hierarchical loading_surfaces 為 active consumer pattern
- 為未來其他 domain split decision 提供標準 framework + probe-then-decide 方法論
- 誠實 acceptance lesson 被鎖定：per-task line savings depend on gate distribution; aggregate task-mix economy is the right metric

### 負面

- 文件雙位點：ADR-009（lock-in）+ governance/cognitive-slice-taxonomy.md（living spec）；須嚴守 ADR 不重複 normative content 的原則避免 drift
- 反轉成本上升：framework 修訂須走 supersede 流程
- 對 contributor 認知成本：理解 slice schema 14 欄 + placement predicate 比理解 monolithic file 高

### 風險

| 風險 | 緩解 |
|---|---|
| ADR 與 governance taxonomy 飄移 | ADR 全段使用 reference link，禁止複製 normative 內容；governance §6 加 ADR-009 back-reference；commit-msg `markdown_yaml_sync` 類 family 規則覆蓋 |
| 未來 framework 改動觸發 ADR 修訂 | §Future Considerations 預先列出 deferred candidates 與升級路徑 |
| Pilot 證據三選都來自 single-author / session-series | 接受此限制；Status 仍標 Accepted（已過 ADR Promotion Criteria）；未來其他 contributor 應用後可在 ADR §History 增補 multi-source validation note |
| ADR 升級後 contributor 誤把 governance 當 proposal layer | governance/cognitive-slice-taxonomy.md §6 開頭明標 `Status: ADR-009 accepted` |

## Alternatives Considered

- A. **不升 ADR，繼續用 governance/**：reject。taxonomy 已被其他 plan 必引用；不升 ADR = cross-cutting decision 留在易動的 governance layer
- B. **等更多 pilot（N ≥ 5）才升**：reject。ADR-007 lighter-target rule 已滿足（second pilot 完成）；繼續累積 pilot 延遲對未來 session 的價值
- C. **升 ADR 並把 taxonomy 全文搬進 ADR**：reject。違反 reference-first 原則；constitution/ 應只放 decision 本體，spec + examples 留 governance/
- D. **升 ADR 並引用 governance（accept）**：採用

## Future Considerations（不阻擋 Accepted；列為已 deferred 的明確項）

每個 deferred item 列 `Trigger to revisit`（觸發再評估的明確條件）+ `Owner at revisit`（誰應該推動），讓任何後續 contributor 讀到時能機械判斷「現在該不該做」，不依賴主動 review。

- **Surface rename sweep**：採用 `execution surface` / `evidence surface` 為對外詞彙的 file rename（例如 `surgical-changes.md` → `surgical-changes-surface.md`）。當前 file 名仍用 slice / 主題名。
  - **Trigger to revisit**: `knowledge/glossary/ai-skill.md` 的 `cognitive_slice` entry 經 **≥ 3 個非本系列 contributor / session-series** 引用後，且 `slice` 詞彙在新文件中造成「arbitrary chunk / static partition」誤解被觀察到 **≥ 2 次**。
  - **Owner at revisit**: 觸發 contributor，開新 rename plan（會牽動 governance taxonomy §7 / §7.5 / 多 routing-registry loading_surfaces / 多 scenario fixtures）。

- **Analysis/travel sources-and-tools 二次 probe**：Parent plan Scenario G 標 `SPLIT_CONDITIONALLY`。
  - **Trigger to revisit**: 真實 **urban walking + 公共交通** 旅遊任務出現，且在該任務中 `analysis/travel/sources-and-tools.md` 載入 **utilization < 50%**（per Scenario G `proposed_action`）。
  - **Owner at revisit**: 該真實任務的 session agent，按 Scenario G 模板寫第二份 probe fixture；如 utilization 持續 ≥ 50% 則維持 NO_SPLIT。

- **Analysis/apk frida-hook-flow / media-hls-analysis-flow probe**：Parent plan §Phase 4 Extension 標 defer；不為形式完整補測。
  - **Trigger to revisit**: 真實任務**單獨**觸發 frida-hook-flow.md 或 media-hls-analysis-flow.md（即兩者非與其他 apk surfaces 同批載入），且該 surface 載入 utilization < 50%。
  - **Owner at revisit**: 該真實任務的 session agent。

- **Slice schema 機械驗證 validator**：目前 placement predicate 為文件規則 + scenario fixture 驗證；未來可加 commit-msg validator 機械檢查新 slice header 欄位完整、`evidence_refs` ≥ 2 for intelligence。
  - **Trigger to revisit**: 觀察到 contributor **誤放案例 ≥ 月 1 次**（intelligence 升層無 evidence_refs、slice 缺 schema 欄位、type 用了非 4 種 primary）。Tracking 由 `plans/active/2026-05-31-2100-mechanical-enforcement-registry.md` 的 Coverage Report 自然覆蓋。
  - **Owner at revisit**: P1 mechanical-enforcement-registry plan owner（不在本 ADR scope）。

- **Aggregate-economy 量化追蹤**：apk-analysis Phase 4 acceptance retrospective 採 weighted task-mix（65/15/20）；未來若有 telemetry / fitness signal 可實測 task distribution。
  - **Trigger to revisit**: `plans/active/2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md` 落地後，telemetry infrastructure 提供 task-mix distribution 實測數據。
  - **Owner at revisit**: Gen 4 fitness plan owner；本 ADR 的 65/15/20 加權僅為 weighted estimate，replace 時無需 supersede ADR-009（屬 retrospective 數據 refinement）。

## Validation Evidence

**Pilots（3）**：

| Pilot | 結果 | Plan |
|---|---|---|
| software-delivery 6 lifecycle slice | done | 2026-05-29-0916 |
| greenfield no-split（反向證據） | done | 2026-05-29-0916 Phase 4 Scenario E |
| apk-analysis artifact-gates 8 slice | done | 2026-05-30-2200 |

**Scenarios（11 fixture）**：`validation/scenarios/software-delivery/slice-load-scenario-{a,b,c,d,e,f,g,h,ag-a,ag-b,ag-c,ag-acceptance-comparison,ag-schemes-a-vs-b}.yaml`（共 13 個 fixture）

**Active consumers**：
- `knowledge/runtime/routing-registry.yaml` route.workflow.software-delivery + route.workflow.apk-analysis 的 hierarchical `loading_surfaces`（共 15 個 surface）
- `knowledge/glossary/ai-skill.md` `cognitive_slice` canonical entry
- `governance/cognitive-slice-taxonomy.md` §7 + §7.5（2 pilot 盤點）
- 2 個 archived parent plans 作為 evidence trail

**Mechanical enforcement**：
- `validatePlanArchivalAudit` commit-msg validator（2026-05-31 修 rename detection bug 後生效）覆蓋 plan archive 階段的 unchecked-item gate
- Placement predicate 文件規則 + Scenario D 負向驗證機械擋下 evidence_refs < 2 的 premature promotion

## Vocabulary Evolution

`cognitive_slice` 是 canonical term；operational synonyms 已於 governance §6 + glossary entry 明示。未來文件演化中：

- 對外 / runtime-oriented：`execution surface` / `evidence surface` / `loading surface`
- governance / 內部設計：`slice`
- 機械可檢查語境（schema、predicate、scenario fixture）：`slice`

避免在 normative 文件中混用 `slice` 與 `chunk` / `partition` / `module`；後三者語意已被既有 codebase 用於其他概念。

## Related

- [ADR-003](ADR-003-three-layer-architecture.md) — 本 ADR 規範該三層架構的載入粒度，不取代三層本體
- [ADR-006](ADR-006-registry-first-workflow-activation.md) — `loading_surfaces` 為本 framework 與 registry 的接點
- [ADR-007](ADR-007-constitution-and-decision-promotion-boundary.md) — 本 ADR 依其 No-Proposed-ADR Rule 與 promotion criteria 走完整 pipeline 才 accepted
- [ADR-008](ADR-008-runtime-cognitive-modes.md) — cognitive_slice 與 cognitive_mode 為兩個獨立 primitive：mode 規範 agent 的執行偏好，slice 規範 cognition surface 的載入粒度
- [`governance/cognitive-slice-taxonomy.md`](../governance/cognitive-slice-taxonomy.md) — canonical living spec
- [`knowledge/glossary/ai-skill.md`](../knowledge/glossary/ai-skill.md) `cognitive_slice` — canonical glossary entry
