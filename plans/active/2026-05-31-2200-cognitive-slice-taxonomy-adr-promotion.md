# Cognitive Slice Taxonomy ADR Promotion

**Status**: `draft`
**世代**：Gen 3 governance promotion（既有 framework 升 ADR，無新功能）
**建立日期**：2026-05-31
**最後更新**：2026-05-31
**Target ADR**：`constitution/ADR-009-cognitive-slice-taxonomy.md`
**Parent plans**：
- [`plans/archived/2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md`](../archived/2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md) — original framework + first pilot（software-delivery 6 lifecycle slice + greenfield no-split）
- [`plans/archived/2026-05-30-2200-apk-analysis-artifact-gates-decomposition.md`](../archived/2026-05-30-2200-apk-analysis-artifact-gates-decomposition.md) — second pilot（apk-analysis 8 slice）

> 本 plan 不重新設計 framework，也不跑新 pilot。它把已驗證的 cognitive slice taxonomy 從 `governance/cognitive-slice-taxonomy.md` 升為 accepted ADR，並 promote 必要的 cross-link 到 constitution layer。

---

## Decision Rationale

### Why Now

Parent plan §ADR Promotion Criteria 全 8 條 PASS（2026-05-30）；apk-analysis follow-up plan（2026-05-31 archived）追加證明：
1. **Cross-project**：3 pilot 應用——software-delivery（6 lifecycle slice）、greenfield（no-split 反向決定）、apk-analysis（8 slice 含 documentation-discipline probe override）
2. **Cross-session**：跨多個 session 持續 reference taxonomy 而未需重新設計
3. **Expensive-to-reverse**：taxonomy 已 wire 至 `knowledge/runtime/routing-registry.yaml` 的 2 個 route + `knowledge/glossary/ai-skill.md` canonical term + 2 個 archived plan 的 self-reference
4. **Foundational + explains-why**：三層邊界規則 + placement predicate + dependency_budget heuristic 是其他切分決定的 grammar，不只是 case-specific guideline
5. **System truly uses contract**：11 個 scenario fixtures + 2 個 archived plans + glossary registration + 2 個 hierarchical loading_surfaces 為 active consumers

### Decision

升 `governance/cognitive-slice-taxonomy.md` 的核心規則為 **ADR-009: Cognitive Slice Taxonomy**：

| Canonical rule | ADR-009 §內容 |
|---|---|
| §1 14-欄 slice schema（含 dependency_budget heuristic + evidence_refs gate） | Decision: 機械可檢查的 slice 定義單位 |
| §2 type+tags 收斂規則（4 primary types only） | Decision: 防 taxonomy explosion |
| §3 Granularity rule（minimum unit = 1 cognitive phase，非 step/concept） | Decision: 防 over-fragmentation |
| §4 三層邊界 + falsifiable placement predicate（intelligence 需 evidence_refs≥2，否則退回 analysis） | Decision: workflow/analysis/intelligence 邊界 enforceable |
| §5 Examples suppression（type:examples default_load:false） | Decision: 防 example-driven contamination |
| §6 Naming（cognitive_slice canonical + execution/evidence surface synonyms） | Decision: vocabulary 已 register glossary |

ADR-009 不重複 taxonomy 全文，只引用 + 標 Status: Accepted；governance 文件仍為 living spec，ADR 鎖定的是「這個 framework 已是 cross-cutting decision」這件事。

### Alternatives Considered

- A. **不升 ADR，繼續用 governance/**：reject。taxonomy 已是其他 plan 必引用的 framework；不升 ADR 等於把 cross-cutting decision 留在易動的 governance layer
- B. **等更多 pilot（N≥5）才升**：reject。ADR-007 lighter-target 規則已滿足；繼續累積 pilot 只延遲對未來 session/user 的價值
- C. **升 ADR 並把 taxonomy 全文搬進 ADR**：reject。違反 reference-first 原則，constitution/ 應只放 accepted decision 本體；spec 與 examples 留 governance/
- D. **升 ADR 並引用 governance（accept）**：採用

### Why ADR Now Justified

| ADR-007 promotion criterion | 證據 |
|---|---|
| foundational | 三層規則是其他切分決定的基礎 grammar |
| cross-session | 跨 session 持續被 reference 而未需重設計 |
| cross-project | 3 個 pilot domain 應用 |
| expensive-to-reverse | wire 至 routing + glossary + 2 archived plan + 11 scenarios |
| explains-why | placement predicate 把「為什麼歸這層」變 falsifiable |
| no lighter target | 已做 second pilot；no lighter promotion exists |
| system truly uses | active consumers：2 hierarchical routes + 11 scenarios + 1 glossary canonical |

### Consequences

#### 正面
- Cognitive slice framework 從 governance proposal 升為 architectural canon
- 未來其他 plan 引用 ADR-009 而非 plan-specific governance doc
- Bug-fix-mode commit-msg 的 archival audit 已就位，後續 ADR-bound plan archive 走得乾淨

#### 負面
- 維護成本：ADR-009 一旦發布，反轉成本變高（須走 supersedes 流程）
- 文件雙位點：ADR-009 是 lock-in，governance 是 living spec；需嚴守「ADR 不重複 governance 內容」的原則避免 drift

#### 風險

| 風險 | 緩解 |
|---|---|
| ADR 內容與 governance taxonomy 飄移 | ADR 只用 reference link，禁止複製 normative 內容；commit-msg `markdown_yaml_sync` 類規則覆蓋 |
| 未來 framework 改動觸發 ADR 修訂 | 預先在 ADR-009 §Future Considerations 列出已 deferred 候選（domain-local `slices/` 子目錄 / `surface` rename sweep）+ 升級路徑 |
| Pilot 證據三選都來自同一作者 / session series | 接受此限制；ADR Status 可標 `Accepted (single-author evidence)`，後續他人應用後升 `Accepted (multi-source validated)` |

### Glossary Impact

Glossary Impact: **no new vocabulary**。`cognitive_slice` 已於 parent plan Phase 5 註冊。ADR-009 只引用既有 glossary entry。

---

## Runtime Execution Path

本 plan 為 docs-only governance promotion。無 generated surface 變更、無 commit-msg validator 新增、無 runtime.db schema 變動。

| 欄位 | 值 |
|---|---|
| Runtime owner | constitution/ + governance/ |
| Trigger | manual: stakeholder 同意啟動本 plan |
| Loaded contract | constitution/ADR-007 (promotion rule) + governance/cognitive-slice-taxonomy.md + 兩個 parent archived plans |
| Runtime action | 新增 constitution/ADR-009 + cross-link governance doc |
| Validation | grep ADR-009 reference 從其他 ADR / governance / constitution README 進入 |
| Evidence | 兩個 parent archived plans + 11 scenario fixtures + glossary entry + routing-registry loading_surfaces |

### Per-Surface Consumer 表

| Generated surface key | Named consumer(s) | Consumer 類型 |
|---|---|---|
| `constitution/ADR-009-cognitive-slice-taxonomy.md` | `constitution/README.md` index + `governance/cognitive-slice-taxonomy.md` §status link back | static reference (manual_activation) |
| 無新 routing-registry 條目 | n/a | n/a |
| 無新 runtime.db generated surface | n/a | n/a |

---

## Open Questions

- [ ] ADR-009 是否同時聲明 `governance/cognitive-slice-taxonomy.md` 為 ADR-009 的 living spec（即明確 ownership boundary：ADR 鎖決定，governance 維護 schema 細節）？預設 yes。
- [ ] ADR-009 是否需要新 commit-msg validator（例如 「新 plan 引入新 slice 時 validate schema 欄位完整」）？預設 no——既有 `validatePlanArchivalAudit` + `validateMarkdownYamlSync` 已覆蓋；如有缺口再開 follow-up。
- [ ] 是否在 ADR-009 §Future Considerations 列入 deferred items（surface rename sweep、analysis/travel 二次 probe、analysis/apk frida-hook-flow probe）？預設 yes（誠實記錄而非吞掉）。

---

## 完成條件

- [ ] Phase 0 Preflight 完成（ADR 編號確認、constitution/README.md 索引位置確認、ADR-007 promotion rule re-check）
- [ ] Phase 1 ADR-009 起草（reference governance + parent plans，不複製 normative 內容）
- [ ] Phase 2 cross-link 同步（constitution/README.md / governance/cognitive-slice-taxonomy.md / 兩個 archived plan 加 ADR-009 back-reference）
- [ ] Phase 3 validation（grep ADR-009 引用點都 reachable；無 dead link；ADR 文字過 markdown lint）
- [ ] Phase 4 Plan Completion Closure（archive / commit / push / 機械 archival audit pass）

---

## Phase 0 — Preflight

- [ ] 確認下一個 ADR 編號（ADR-009——constitution/ 目前 1–8）
- [ ] 確認 ADR-008 結構為最新 template（已驗證：Status / Framework Generation / Date / Source Plan / Context / Decision / ...）
- [ ] 確認 parent plans 都已 archive（已：2026-05-29-0916 + 2026-05-30-2200）
- [ ] 確認 archival-audit validator 已生效（commit 82032c6 後生效——本 plan archive 時會被機械驗證）
- [ ] 確認 governance/cognitive-slice-taxonomy.md 為 active reference（已：被 2 plans + 11 scenarios + glossary 引用）

---

## Phase 1 — ADR-009 起草

- [ ] 建立 `constitution/ADR-009-cognitive-slice-taxonomy.md`：
  - Status: Accepted
  - Framework Generation: Gen 3 governance promotion
  - Date: 2026-05-31
  - Source Plans: 兩個 archived plan
  - Context: framework 動機 + 3 pilot evidence summary
  - Decision: 引用 governance §1–§6 為 canonical rules（不複製）
  - Consequences: 正/負/風險（從本 plan §Decision Rationale 提煉）
  - Future Considerations: deferred items（surface rename / analysis/travel 二次 probe / frida-hook-flow probe）
  - Validation Evidence: 列 11 scenario fixtures + 2 parent plans + glossary + routing surfaces

---

## Phase 2 — Cross-Link 同步

- [ ] `constitution/README.md` 加 ADR-009 entry
- [ ] `governance/cognitive-slice-taxonomy.md` 加 「Status: ADR-009 accepted 2026-05-31」標頭
- [ ] `plans/archived/2026-05-29-0916-...md` 加 「Promoted to ADR-009」尾註
- [ ] `plans/archived/2026-05-30-2200-...md` 加同樣尾註
- [ ] `knowledge/glossary/ai-skill.md` cognitive_slice entry 加 「formalized in ADR-009」

---

## Phase 3 — Validation

- [ ] `grep ADR-009` 從 constitution/README + governance taxonomy + 2 archived plans + glossary 五處可達
- [ ] 無 dead link（archived plan back-link 指 constitution/ADR-009 相對路徑正確）
- [ ] ADR 文字無 normative content 複製（只 reference）

---

## Phase 4 — Plan Completion Closure

- [ ] Link audit
- [ ] 更新本 plan 狀態與完成日期
- [ ] Plan archival audit 機械 pass（本 plan 內所有 `- [ ]` 補完）
- [ ] Archive / commit / push

---

## 與其他 plans 的關係

- Parent (framework + first pilot): [`2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md`](../archived/2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md)
- Parent (second pilot): [`2026-05-30-2200-apk-analysis-artifact-gates-decomposition.md`](../archived/2026-05-30-2200-apk-analysis-artifact-gates-decomposition.md)
- 不相依：active mechanical-enforcement / activation-engine / sanitization-validator plans（各自工作流）
