# APK Analysis Artifact Gates Decomposition

**Status**: `in-progress`
**世代**：Gen 3 runtime hardening（second pilot applying validated framework）
**建立日期**：2026-05-30
**最後更新**：2026-05-30（Phase 0 完成）
**Parent plan**：[`2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md`](../archived/2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md) — provides taxonomy, three-layer rule, dependency_budget heuristic, scenario fixture template, and the placement predicate this plan applies.
**Pilot 決定**：`workflow/apk-analysis/artifact-gates.md`（575 行，12 個 gate）

> 本 plan 不重新設計 framework。它套用 parent plan Phase 1–4 已驗證的 cognitive slice taxonomy（[`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md)）把 `workflow/apk-analysis/artifact-gates.md` 切成 focused gate slices。

---

## Decision Rationale

### Empirical Trigger（不是想像，是 runtime evidence）

Parent plan Scenario F（[`slice-load-scenario-f-apk-analysis-probe.yaml`](../../validation/scenarios/software-delivery/slice-load-scenario-f-apk-analysis-probe.yaml)）實測結果：

- `artifact-gates.md` 575 行，12 個 distinct gate
- 真實 APK analysis 任務（capture login API + document endpoints + feature handoff）只使用 6 個 gate（~280 行）
- **inflation ratio ~1.57**（real load 575 ÷ ideal 280 ≈ 2.05x over-load 在 artifact-gates 單檔）
- 與 parent plan software-delivery 切分前症狀一致

### Decision

切成 **8 個** focused gate surface（Scenario F §verdict.proposed_slices 原列 7 + Phase 0 probe AG-schemes-A-vs-B 加 1）：

| Slice 檔 | 內容（原 §） | type | tags |
|---|---|---|---|
| `artifact-gates/ui-architecture-map.md` | §1 + §10 template | execution | artifact-gate, ui |
| `artifact-gates/api-catalog.md` | §2 + §11 details | execution | artifact-gate, api |
| `artifact-gates/domain-runtime-baseline.md` | §3 | execution | artifact-gate, domain |
| `artifact-gates/feature-handoff.md` | §4 | execution | artifact-gate, handoff |
| `artifact-gates/evidence-chain.md` | §5 note template + §6 chain + §7 failure log | execution | artifact-gate, evidence |
| `artifact-gates/sanitization.md` | §12 | execution | artifact-gate, sanitization |
| `artifact-gates/self-generation-audits.md` | §8 SDK + §9 Identity | failure | artifact-gate, security |
| `artifact-gates/documentation-discipline.md` | §13 dev notes + §14 feedback tips + §15 backfill rules | execution | artifact-gate, documentation, backfill |

`artifact-gates.md` 變成 thin index（同 parent plan `execution-flow.md` 模式）。

### Why Not Extend Parent Plan

Parent plan §Phase 4 Extension 已記錄此 SPLIT 決定。但 parent plan scope 鎖定 software-delivery pilot；擴張會延遲 archive、模糊 close-out。本 plan 為 follow-up，繼承 parent 已驗證 framework，不重新跑 Phase 0–3 設計工作。

---

## Phase Plan（精簡，繼承 parent framework）

> 不重複 parent plan 已做的工作：taxonomy 已穩定、三層規則已驗證、scenario fixture 模板已存在。本 plan 只跑「應用」階段。

### Phase 0 — Preflight（繼承）

#### Phase 0.0 — Open Questions 核對（公版，必填）

逐條核對本 plan §Open Questions，標記處置並回寫：

- [x] 已讀本 plan §Open Questions 全部條目
- [x] 對每條標記 `resolved` / `still-open` / `deferred`
- [x] `resolved` 的條目已同步勾選 / 附註於 §Open Questions
- [x] 若盤點新發現問題，已加入 §Open Questions（**Phase 0 盤點新增 2 條：gate count discrepancy 與 yaml split policy；已 append**）

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| frida-hook-flow / media-hls-analysis-flow 是否納入本 plan | deferred | parent plan 已標 defer 至真實單獨任務再 probe；維持不納入 |
| analysis/travel sources-and-tools 第二次 probe | out-of-scope | 屬另一個 follow-up plan |
| **NEW: artifact-gates.md 真實 gate 數量**（Phase 0 發現） | resolved → Phase 1 處置 | 實測 15 個 gate（§1–§15），非 Scenario F 估計的 12 個；§13/§14/§15 在 parent plan probe 時被漏算。詳見下方 §Phase 0 Inventory Record |
| **NEW: artifact-gates.yaml 是否同步切分**（Phase 0 發現） | resolved → Phase 2 處置 | 不切 yaml；保留單一 yaml contract，其 `source_markdown` 改指向新 thin-index `artifact-gates.md`，executable contract 內部 reference 在 Phase 2 同步更新 |

#### Phase 0 Preflight 核對

- [x] 確認 parent plan 已 archive 或進入 Phase 5。**Result**: parent plan 已於 2026-05-30 archive 至 `plans/archived/`。
- [x] 確認 `governance/cognitive-slice-taxonomy.md` 仍是 active reference。**Result**: file exists, status active；`cognitive_slice` 已於 2026-05-30 註冊為 canonical glossary term。
- [x] 確認 `workflow/apk-analysis/artifact-gates.yaml`（衍生 contract）的 `source_markdown` 對映需在 Phase 2 同步更新。**Result**: yaml header 確認 `source_markdown: workflow/apk-analysis/artifact-gates.md` + `runtime_projection.enabled: true`。Phase 2 須改 source_markdown 指向 thin-index 後的 `artifact-gates.md`，且若 thin index 內仍含關鍵 executable refs 則 `source_markdown` 不變；若 thin index 變成純導航則需另議 yaml 是否仍掛同檔。

#### Phase 0 Inventory Record（2026-05-30）

**Architecture Compatibility Preflight**（依 `plans/README.md` §Architecture Compatibility Preflight）：

| 欄位 | 內容 |
|---|---|
| Trigger | 啟動本 follow-up plan Phase 0 |
| Checked sources | `workflow/apk-analysis/artifact-gates.md`（575 行 heading 結構）/ `workflow/apk-analysis/artifact-gates.yaml` header / `knowledge/runtime/routing-registry.yaml` route.workflow.apk-analysis / `governance/cognitive-slice-taxonomy.md` §1/§4/§7 / parent plan Scenario F |
| Conflicts | **1 個 evidence revision**（非阻擋）：Scenario F 估 12 gate，實測為 15。Phase 1 須擴充 slice mapping。**1 個 governance question**（非阻擋）：thin index 是否仍能承擔 yaml `source_markdown` 對映；Phase 2 處置。 |
| Decision | proceed — 沒有阻擋性架構衝突；evidence revision 收進 §Phase 0 修正紀錄；slice mapping 從 7 擴為 8（提案）或維持 7（吸收 §13/§14/§15 進既有 slice）由 Phase 1 stakeholder 決 |
| Validation | read-only heading scan + line count；無檔案內容變更 |

**Gate 實測盤點（artifact-gates.md 575 行）**：

| Gate § | 標題 | 行範圍 | Scenario F 收入 | Phase 0 修正 |
|---|---|---|---|---|
| 1 | UI Architecture Map | 5–33 | ✓ → ui-architecture-map | 不變 |
| 2 | API Catalog | 34–70 | ✓ → api-catalog | 不變 |
| 3 | Domain/Runtime Baseline | 71–103 | ✓ → domain-runtime-baseline | 不變 |
| 4 | Feature Reconstruction Handoff | 104–130 | ✓ → feature-handoff | 不變 |
| 5 | 單次分析筆記模板 | 131–196 | ✓ → evidence-chain（合併 §5+§6+§7） | 不變 |
| 6 | 證據鏈要求 | 197–206 | ✓ → evidence-chain | 不變 |
| 7 | 失敗也要記錄 | 207–222 | ✓ → evidence-chain | 不變 |
| 8 | SDK Live Self-Generation Audit | 223–260 | ✓ → self-generation-audits | 不變 |
| 9 | Authorized Identity Material Self-Generation Audit | 261–287 | ✓ → self-generation-audits | 不變 |
| 10 | UI Architecture Map Template | 288–390 | ✓ → ui-architecture-map（合併 §1+§10） | 不變 |
| 11 | API Catalog Detail Requirements | 391–497 | ✓ → api-catalog（合併 §2+§11） | 不變 |
| 12 | Sanitization Rules | 498–518 | ✓ → sanitization | 不變 |
| **13** | **Developer Guidance Notes（可選）** | **519–532** | **✗ 漏** | **Phase 1 提案：併入 evidence-chain（discipline of evidence writing）** |
| **14** | **Feedback Lesson Writing Tips** | **533–560** | **✗ 漏** | **Phase 1 提案：併入 evidence-chain 或獨立 documentation-discipline slice** |
| **15** | **Backfill Rules** | **561–575** | **✗ 漏** | **Phase 1 提案：併入 evidence-chain（backfill ≈ evidence completeness）** |

**Net impact of revision**：
- 真實 gate 數 15（非 12）
- Slice 數方案經 probe 決定：
  - ~~**方案 A（Phase 0 heuristic 推薦）**：保持 7 slice，§13/§14/§15 合計 ~57 行併入 `evidence-chain.md`~~
  - **方案 B（probe 採用，2026-05-30）**：增第 8 slice `documentation-discipline.md`（§13+§14+§15 ~57 行）；evidence-chain 維持純證據紀錄。詳見 [`slice-load-scenario-ag-schemes-a-vs-b.yaml`](../../validation/scenarios/software-delivery/slice-load-scenario-ag-schemes-a-vs-b.yaml)。
- 任一方案都不需改 Scenario F 的 `~280 行 ideal load` 估算（這些 gate 在原 task 中也大多不會用到）

**Source-of-truth 邊界（Phase 2 必須遵守）**：

- `artifact-gates.yaml`（executable contract）以 `runtime_projection.enabled: true` 投影到 `runtime.db`。**該 yaml 不切**；Phase 2 改 `source_markdown` 指向新的 thin-index `artifact-gates.md`，並評估 yaml 內部步驟 reference 是否需更新指向 focused slice。
- `artifact-gates.md` 變 thin index，保留 redirect stub 兼容舊入口（parent plan 同樣模式）。

**受影響的 routing / summary / README**：

- `knowledge/runtime/routing-registry.yaml` → `route.workflow.apk-analysis`：`required_dependencies` 含 `workflow/apk-analysis/artifact-gates.md`。Phase 3 加 `loading_surfaces` hierarchical mapping。
- `knowledge/summaries/apk-analysis-pilot.md`：Phase 3 同步。
- `workflow/apk-analysis/README.md`：56 行，已是 thin-ish 入口；Phase 2 切片後需更新導航連結。

**workflow ↔ analysis 邊界確認**（沿用 parent Scenario H 結論）：

- workflow/apk-analysis 與 analysis/apk **不整併**（Scenario H decision: KEEP_SEPARATE）。本 plan 切分只影響 workflow side，不動 analysis side。

#### Phase 0 結論

- [x] Candidate inventory 完成（15 gate vs Scenario F 估的 12 gate 已修正）
- [x] Owner-layer decision 確認（全為 workflow layer，不跨 analysis）
- [x] Architecture conflict 為 0（gate count 是 evidence revision，不是 conflict）

**Phase 0 exit criteria 全達成**。下一步進 Phase 1：選 A/B slice 方案 + 在 taxonomy §7 新增 apk-analysis pilot section。

### Phase 1 — Slice Schema（繼承，無需重定義）

- [x] 在 [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) **§7.5 新增 apk-analysis pilot section**（採 Scheme B = **8 slice**，非 Scenario F 原列 7；§13/§14/§15 經 probe 證實獨立成 `apk-documentation-discipline`）
- [x] 每個 slice 套用 parent plan §1 schema（14 欄）+ `dependency_budget` heuristic（全 default 2/4，無 high override）
- [x] Granularity / placement / examples-suppression / extraction-direction 規則檢查通過（taxonomy §7.5 已記錄）

**Phase 1 exit criteria**：
- [x] Taxonomy 不重複 canonical source，artifact-gates 範圍內 8 slice 全為 workflow layer
- [x] type+tags 收斂規則成立（7 execution + 1 failure；artifact-gate 為共通 tag）
- [x] 每 slice 有明確 load_when / do_not_load_when
- [x] 三層邊界與 placement predicate 通過：8 slice 均不跨入 analysis（evidence acquisition）或 intelligence（long-term pattern）
- [x] Scheme A/B probe 完成（see [`slice-load-scenario-ag-schemes-a-vs-b.yaml`](../../validation/scenarios/software-delivery/slice-load-scenario-ag-schemes-a-vs-b.yaml)）— Scheme B 採用

### Phase 2 — Thin Index + Focused Slices

- [ ] 依 Scenario F §proposed_slices 拆 7 個檔案到 `workflow/apk-analysis/artifact-gates/`
- [ ] `artifact-gates.md` 改為 thin index（保留 redirect stub 兼容舊入口）
- [ ] 跨檔內容（§1/§10、§2/§11、§5/§6/§7、§8/§9）同批拆，避免 dual source-of-truth
- [ ] 更新 `artifact-gates.yaml` 的 `source_markdown` mapping

### Phase 3 — Routing + Summary

- [ ] `knowledge/runtime/routing-registry.yaml` route.workflow.apk-analysis 加 `loading_surfaces`（hierarchical，不平攤新 route）
- [ ] `knowledge/summaries/apk-analysis-pilot.md` 同步
- [ ] `knowledge/graphs/apk-analysis-pilot.yaml` 同步

### Phase 4 — Validation Scenarios

- [ ] **Scenario AG-A**（execution-only：產生新一輪 API catalog 但 UI map 已存在）— 預期只載入 api-catalog + evidence-chain，不載入 ui-map / self-gen audits
- [ ] **Scenario AG-B**（full analysis：對全新 APK 從零做完整 handoff）— 預期載入所有 7 slice
- [ ] **Scenario AG-C**（security focus：只做 self-generation audit）— 預期載入 self-generation-audits + evidence-chain + sanitization，不載入 ui-map / api-catalog / feature-handoff
- [ ] Acceptance：post-split 對 Scenario F 同樣任務的 load 從 575 → ~280 行（≥40% 節省）

### Phase 5 — Linked Updates + Closure

- [ ] link audit、`ai-skill runtime refresh`、plan archive

---

## Inheritance from Parent Plan

| Asset | Parent location |
|---|---|
| Slice schema | [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) §1 |
| Three-layer rule | taxonomy §4 |
| Placement predicate | taxonomy §4 + parent Scenario D |
| dependency_budget heuristic | taxonomy §1 + parent §Decision Rationale 風險表 |
| Scenario fixture shape | parent §Phase 4 Fixture |
| Hierarchical routing rule | parent §Phase 3 |
| Examples suppression rule | taxonomy §5（本 plan 無 examples slice，n/a） |
| Empirical trigger | [`slice-load-scenario-f-apk-analysis-probe.yaml`](../../validation/scenarios/software-delivery/slice-load-scenario-f-apk-analysis-probe.yaml) |

---

## Open Questions

- [x] 是否同時把 Scenario F 標記的 `analysis/apk/workflows/frida-hook-flow.md`（190 行）/`media-hls-analysis-flow.md`（165 行）納入本 plan 範圍？預設不納入（parent plan 標 defer 至真實單獨任務觸發再 probe）。**Phase 0 確認：deferred，維持不納入。**
- [x] `analysis/travel/sources-and-tools.md` 的第二次 probe 屬於另一個 follow-up plan，不歸本 plan。**Phase 0 確認：out-of-scope。**
- [x] **NEW（Phase 0 新增 2026-05-30）**：artifact-gates.md 實測有 15 個 gate（非 Scenario F 估計的 12 個）。§13 Developer Guidance Notes / §14 Feedback Lesson Writing Tips / §15 Backfill Rules 應該：(A) 併入 evidence-chain.md（總 ~57 行；保持 7 slice）／ (B) 獨立第 8 slice `documentation-discipline.md`？~~Phase 0 推薦 A（granularity 較緊；§13+§14+§15 都是 "evidence/doc discipline" 主題且行數小）。~~ **Probe override 2026-05-30**：執行 [`slice-load-scenario-ag-schemes-a-vs-b.yaml`](../../validation/scenarios/software-delivery/slice-load-scenario-ag-schemes-a-vs-b.yaml) 三任務（T1 標準分析 / T2 doc discipline 審閱 / T3 全程含 doc），weighted across 65% T1 + 20% T2 + 15% T3 顯示 **Scheme B 每 100 sessions 多省 ~7130 行**；同時 Scheme A 違反 granularity rule（evidence-chain 在 A 下混合 evidence recording + writing discipline 兩個 cognitive phase）。**結論：採 Scheme B = 8 slice**。Phase 0 heuristic 被 probe 覆蓋——這是本系列工作中**第二次** runtime probe 覆蓋直覺（parent plan Scenario H 是第一次，方向相反），證實「probe-then-decide 一致勝出純啟發」。
- [ ] **NEW（Phase 0 新增 2026-05-30）**：`artifact-gates.yaml` executable contract 是否切？**Phase 0 處置：不切**；保留單一 yaml，Phase 2 改 `source_markdown` 指向 thin-index，並評估 yaml 內部 step references 是否需要更新指向 focused slice。若 thin index 過薄無法承擔 source_markdown 對映，再回頭討論 yaml split。

---

## 與其他 plans 的關係

- Parent: [`2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md`](../archived/2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md)
- Sibling (potential, deferred): travel-planning 二次 probe plan
