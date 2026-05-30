# APK Analysis Artifact Gates Decomposition

**Status**: `draft-skeleton`
**世代**：Gen 3 runtime hardening（second pilot applying validated framework）
**建立日期**：2026-05-30
**最後更新**：2026-05-30
**Parent plan**：[`2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md`](2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md) — provides taxonomy, three-layer rule, dependency_budget heuristic, scenario fixture template, and the placement predicate this plan applies.
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

切成 7 個 focused gate surface（Scenario F §verdict.proposed_slices）：

| Slice 檔 | 內容（原 §） | type | tags |
|---|---|---|---|
| `artifact-gates/ui-architecture-map.md` | §1 + §10 template | execution | artifact-gate, ui |
| `artifact-gates/api-catalog.md` | §2 + §11 details | execution | artifact-gate, api |
| `artifact-gates/domain-runtime-baseline.md` | §3 | execution | artifact-gate, domain |
| `artifact-gates/feature-handoff.md` | §4 | execution | artifact-gate, handoff |
| `artifact-gates/evidence-chain.md` | §5 note template + §6 chain + §7 failure log | execution | artifact-gate, evidence |
| `artifact-gates/sanitization.md` | §12 | execution | artifact-gate, sanitization |
| `artifact-gates/self-generation-audits.md` | §8 SDK + §9 Identity | failure | artifact-gate, security |

`artifact-gates.md` 變成 thin index（同 parent plan `execution-flow.md` 模式）。

### Why Not Extend Parent Plan

Parent plan §Phase 4 Extension 已記錄此 SPLIT 決定。但 parent plan scope 鎖定 software-delivery pilot；擴張會延遲 archive、模糊 close-out。本 plan 為 follow-up，繼承 parent 已驗證 framework，不重新跑 Phase 0–3 設計工作。

---

## Phase Plan（精簡，繼承 parent framework）

> 不重複 parent plan 已做的工作：taxonomy 已穩定、三層規則已驗證、scenario fixture 模板已存在。本 plan 只跑「應用」階段。

### Phase 0 — Preflight（繼承）

- [ ] 確認 parent plan 已 archive 或進入 Phase 5（避免兩 plan 同時動 routing）
- [ ] 確認 `governance/cognitive-slice-taxonomy.md` 仍是 active reference
- [ ] 確認 `workflow/apk-analysis/artifact-gates.yaml`（衍生 contract）的 `source_markdown` 對映需在 Phase 2 同步更新

### Phase 1 — Slice Schema（繼承，無需重定義）

- [ ] 在 [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) §7 新增 apk-analysis pilot section，列出 7 個 slice 與 Scenario F 已定義的欄位
- [ ] 每個 slice 套用 parent plan §1 schema（14 欄）+ `dependency_budget` heuristic

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

- [ ] 是否同時把 Scenario F 標記的 `analysis/apk/workflows/frida-hook-flow.md`（190 行）/`media-hls-analysis-flow.md`（165 行）納入本 plan 範圍？預設不納入（parent plan 標 defer 至真實單獨任務觸發再 probe）。
- [ ] `analysis/travel/sources-and-tools.md` 的第二次 probe 屬於另一個 follow-up plan，不歸本 plan。

---

## 與其他 plans 的關係

- Parent: [`2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md`](2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md)
- Sibling (potential, deferred): travel-planning 二次 probe plan
