# Gen3 Workflow and Analysis Cognitive Slice Decomposition

**Status**: `in-progress`
**世代**：Gen 3 current hardening；Gen 4 prerequisite
**建立日期**：2026-05-29
**最後更新**：2026-05-29
**Pilot 決定**：`workflow/software-delivery/`（stakeholder 同意 2026-05-29）
**Glossary 決定**：`Cognitive Slice` 正式註冊刻意延後至 Phase 4 validation 之後；過渡期一律用 `loading/execution/evidence surface` 既有措辭（stakeholder 同意 2026-05-29，external review 採納）

> 本 plan 的核心決策：`workflow/` 與 `analysis/` 的 cognitive slice boundary 應在 Gen 3 完成，不等 Gen 4。Gen 3 先把 execution path（workflow）與 evidence path（analysis）切成可獨立載入、驗證與路由的 cognitive units；Gen 4 才在此基礎上做 dynamic activation / ecosystem orchestration。

---

## Decision Rationale

### Problem & Why Now

外部 skill collection（例如 `mattpocock/skills`）的設計重心是 small, composable, task-oriented skill：TDD、diagnose、triage、PRD、handoff 等。這類 skill 多數像 prompt macro 或短流程 capability，靠簡短 instruction 即可運作。

Ai-skill 目前已不是這種層級。現行 Gen 3 文件已明確定位為 **AI-native Cognitive Execution System**，具備：

- runtime state machine
- routing registry
- executable YAML contracts
- knowledge / workflow / analysis / intelligence 分離
- cognitive modes
- failure-derived validation
- memory retrieval / activation governance
- close-loop writeback transaction

因此問題不是單純「入口文件太肥」，而是 `workflow/` 與 `analysis/` 的 cognitive unit 粒度過粗。若 execution flow、artifact gates、examples、evidence acquisition、tool procedure、failure caveats 與 domain-specific triage 全混在同一 loading boundary，它會造成：

- **Execution overload**：agent 只需要流程順序時，被迫載入 examples、artifact details、edge cases 與 domain notes。
- **Evidence-path contamination**：agent 只需要證據取得方法時，被迫載入 orchestration flow、handoff、governance 或不相關 tool procedure。
- **Compression degradation**：巨型 workflow / analysis 文件被反覆摘要後，留下抽象口號，失去可執行細節。
- **Routing ambiguity**：`workflow/`、`analysis/`、`intelligence/`、summary 與 source-of-truth 邊界不清。
- **Gen 4 migration debt**：未來 dynamic activation graph 會繼承 monolith boundary，變成 monolithic cognitive ecosystem。

現有架構其實已經支持這個方向：

- Gen 3 architecture 明定 canonical navigation 入口不是 self-contained spec。
- `knowledge/summaries/` 採 summary-first loading。
- `knowledge/runtime/routing-registry.yaml` 已是 task intent → primary source 的路由層。
- `governance/document-sizing.md` 要求超過 150 行且多主題的高頻文件拆成 index + focused child files。
- `workflow/README.md` 已定義 workflow 負責 execution order，`analysis/README.md` 已定義 analysis 負責 evidence acquisition。
- Gen 4 vision 的 Watch-Out List 已警告 workflow inflation、ecosystem boundary inflation、positive-activation bias。

所以這項工作應視為 **Gen 3 runtime hardening**，不是 Gen 4 才開始的願景工作。

### Decision

在 Gen 3 先完成 **Workflow and Analysis Cognitive Slice Decomposition**：

1. 盤點 `workflow/` 與 `analysis/` 中 high-frequency、oversized、multi-topic surfaces。
2. 將巨型 workflow / analysis 文件轉為 thin index + focused child slices。
3. 把內容拆成 focused cognitive slices：
   - workflow execution slice
   - workflow artifact-gate slice
   - workflow examples slice
   - workflow handoff / closure slice
   - analysis evidence-acquisition slice
   - analysis triage / observation slice
   - analysis tool-procedure slice
   - analysis failure / caveat slice
4. 為每個 slice 定義：
   - purpose
   - load_when
   - source-of-truth owner
   - dependencies
   - summary target
   - validation signal
5. 更新 routing / summary / README links，讓 agent 先讀 index，再按 task intent 載入必要 workflow / analysis slices。
6. 補 validation scenarios，確認 execution-only 任務不載入整份 analysis，evidence-only 任務不載入整份 workflow。

Gen 4 只承接後續：

- dynamic cognitive retrieval
- activation graph
- economics-driven loading
- suppression layer
- telemetry / fitness feedback
- adaptive orchestration

本 plan 不實作 autonomous optimizer、self-modifying governance、telemetry DB 或 Gen 4 ecosystem layer。

### Alternatives Considered

- A. 等 Gen 4 再處理：reject。Gen 4 的 activation / orchestration 需要清楚 cognitive unit；若現在不拆，未來只會把 workflow / analysis monolith 帶進 ecosystem。
- B. 只把大檔縮短：reject。問題不是字數，而是 execution path、evidence path、examples、tool procedure 與 validation 混在同一 loading boundary。
- C. 直接建立完整 Gen 4 dynamic retrieval system：reject。會提前跨入 ecosystem orchestration，違反 Gen 4 vision 對 current / vision boundary 的限制。
- D. Gen 3 先做 workflow / analysis slice + index + loading rules：accept。這是 current runtime 可執行的收斂工作，也能替 Gen 4 預留乾淨 substrate。

### Why Not an ADR Yet

此 plan 是文件結構、loading boundary 與 routing hardening。它可能導出新的 architecture decision，但目前尚未證明：

- workflow / analysis slice taxonomy 是否穩定；
- 哪些 slices 應成為 long-lived owner layer；
- routing registry 是否需要新增 route 類型；
- validation scenarios 是否足以防止 retrieval contamination；
- `Cognitive Slice` 是否應成為正式 framework vocabulary。

完成後若 slice taxonomy 成為跨世代、跨工具、難逆轉的基礎 decision，再評估 ADR promotion。

### ADR Promotion Criteria（completed 時驗證）

- [ ] High-frequency `workflow/` 與 `analysis/` surfaces 已完成 inventory。
- [ ] 至少一個 oversized / multi-topic workflow 或 analysis surface 已轉為 thin index + focused slices。
- [ ] 每個 slice 都有 `load_when` / owner / dependency / validation signal。
- [ ] Summary-first loading path 已更新或明確標記 not applicable。
- [ ] Routing registry / README / architecture links 已檢查並完成必要同步。
- [ ] 至少 3 個 validation scenarios 覆蓋 execution-only task / evidence-only task / mixed workflow+analysis task。
- [ ] Open Questions 全部解決或轉為明確 follow-up plan。
- [ ] 實際 agent session 能證明不再為 execution-only / evidence-only 小任務載入整份 workflow+analysis surface。

### Consequences

#### 正面

- 降低 context loading 成本與 retrieval contamination。
- 讓 Gen 3 的 runtime routing 更接近「讀對 workflow / analysis slice」而不是「讀整包流程與證據方法」。
- 替 Gen 4 dynamic activation / cognitive economics 建立乾淨邊界。
- 讓外部使用者更容易理解 Ai-skill 和一般 skill collection 的差異。

#### 負面

- 短期會增加檔案數量與索引維護成本。
- 若切片過細，可能造成 agent 多跳讀取與 route fragmentation。
- 需要補 linked updates，否則舊入口仍可能被 agent 誤讀。

#### 風險

| 風險 | 緩解 |
|---|---|
| 把「切片」誤做成 Gen 4 ecosystem layer | 本 plan 明確限定 Gen 3：index、summary、routing、validation；不做 dynamic orchestration |
| 切太細導致維護成本 > token savings（over-fragmentation） | Phase 1 granularity 原則：slice 最小單位 = 能獨立完成的 cognitive phase，非 step / concept；Phase 2 逐 slice 把關，避免 context hopping / dependency storms（external review 風險2） |
| Cross-slice dependency explosion（recursive loading / fan-out / hidden activation chains） | Phase 1 slice schema 加 dependency budget（heuristic default `max_depth:2`/`max_runtime_dependencies:4` + `override_when: task_complexity=high`，非 rigid）；Phase 4 scenario 斷言實際載入深度/廣度未超預算（external review OQ + #3 風險1，2026-05-29） |
| Example-driven loading contamination（examples 密度高、LLM 易先讀 examples override doctrine） | Phase 1：`type: examples` 預設 `default_load: false`，只在 user 要求範例或偵測 ambiguity 載入；Phase 3 suppression guidance 對齊；對應 Watch-Out Wall 5（external review #3 風險2，2026-05-29） |
| intelligence layer 吞噬 analysis（heuristic/tradeoff/anti-pattern 灰區倒成 thought dumping ground） | Phase 1 extraction direction rule + **falsifiable placement predicate**（`layer_justification` + intelligence `evidence_refs`≥2 可解析，否則退回 analysis）；Phase 4 Scenario D 負向驗證 + B/C contamination 探針抓誤放（external review #3 風險3 + placement 驗證追問，2026-05-29） |
| Placement 不可驗證（歸層淪為 honor-system 標籤，誤放無法偵測） | layer membership predicate 機械可檢查（evidence_refs gate）；Phase 4 Scenario D + contamination 間接探針；限制：無完全機械 oracle，目標是誤放可偵測/可逆，非證明每次放對（placement 驗證追問，2026-05-29） |
| routing-registry 變第二個 monolith（route inflation / flat route universe → giant cognition graph） | Phase 3 hierarchical routing 規則：route 採樹狀命名，不平鋪 leaf；新增前先確認可掛既有層級節點（external review #3 風險4，2026-05-29） |
| Taxonomy explosion / classification obsession | Phase 1 type+tags 收斂：primary `type` 只 4 種（execution/evidence/examples/failure），其餘降為 tags；新增 primary type 需回 plan 重評（external review 風險1） |
| workflow / analysis / intelligence 邊界模糊 | Phase 1 codify 三層分工：workflow=順序、analysis=證據取得+驗證、intelligence=為何長期有效/失敗；歸層用此判定（external review 風險3） |
| 過早變「理論宇宙」/ premature ecosystem abstraction | 維持節奏：small runtime hardening → measurable retrieval improvement → loading reduction → validation proof → gradual orchestration；不一次衝 full autonomous ecosystem（external review meta 警告） |
| 舊 workflow / analysis links drift | Phase 4 必須做 rg link audit + routing registry check |
| 新 taxonomy 與既有 `workflow/analysis/intelligence/knowledge` 重疊 | Phase 0 先做 owner-layer preflight；每個 slice 只導航，不重定義 canonical source |

### Glossary Impact

Glossary Impact: yes.

Candidate framework vocabulary:

- `Cognitive Slice`：可被獨立載入、驗證、路由的最小認知單元；本 plan 先落在 workflow execution slice 與 analysis evidence slice。
- `Retrieval Boundary`：agent 在某任務中應停留的載入邊界。
- `Thin Workflow/Analysis Index`：只負責 navigation / loading guidance，不承載全部 workflow 或 analysis 正文的入口。

替代命名候選（external review 2026-05-29）：`capability surface` / `cognitive surface` / `execution surface`，理由是這些更貼近「routable cognition surface」本質，而 `slice` 易被聯想成 arbitrary chunk / static partition。

Phase 1 只評估命名候選（surface vs slice）並選定**過渡期 operational wording**（`loading/execution/evidence surface`）；**正式註冊 `Cognitive Slice` 到 `knowledge/glossary/ai-skill.md` 刻意延後至 Phase 4 validation 證明 taxonomy 穩定後**，避免 premature naming lock-in / vocabulary inflation。

### Watch-Out List Citation

本 plan 對應 Gen 4 vision [`architecture/ai-native-cognitive-ecosystem-system.md`](../../architecture/ai-native-cognitive-ecosystem-system.md) §Watch-Out List：

- Wall 1. Discovery confused with Activation：本 plan 只做 Gen 3 loading boundary，不把 routing-registry 假裝成 activation graph。
- Wall 2. Workflow inflation：避免 workflow 越寫越大，並防止 analysis 變成另一個 evidence monolith，改用 focused slices。
- Wall 3. Ecosystem boundary inflation：不新增 ecosystem layer；只整理 current source surfaces。
- Wall 5. Positive-activation bias：Phase 2 必須定義 `do_not_load_when` 或 suppression note，避免「多載入就是更完整」的錯覺。

---

## Runtime Execution Path

### Doc-only Trial Statement + Graduation

目前狀態：**Gen 3 workflow / analysis documentation and routing hardening plan**。

本 plan 第一版不新增 runtime generated surface、不新增 commit-msg validator、不建立 Gen 4 activation graph。若 Phase 3 決定要更新 `knowledge/runtime/routing-registry.yaml` 或新增 validation scenario，才進入 runtime-linked execution。

Graduation condition：

| Graduation Phase | Acceptance signal |
|---|---|
| Phase 2 完成 | workflow / analysis index + slices 形狀可讀，且不重複 canonical source |
| Phase 3 完成 | loading rules / routing links 有 named consumer 或明確 manual activation |
| Phase 4 完成 | scenarios 證明 execution-only 任務不載入整份 analysis，evidence-only 任務不載入整份 workflow |

### Runtime Owner

- Current owner layer：`workflow/` + `analysis/` + affected source layer README
- Loading / summary owner：`knowledge/summaries/`
- Routing owner：`knowledge/runtime/routing-registry.yaml`（若需要）
- Validation owner：`validation/scenarios/`（若需要）
- Governance owner：`governance/document-sizing.md` + `enforcement/linked-updates.md`

### Trigger Flow

```text
event_or_signal:
  - high-frequency workflow / analysis file exceeds document-sizing threshold
  - workflow file mixes execution order / artifact gates / examples / closure / analysis method
  - analysis file mixes evidence acquisition / tool procedure / triage / caveats / workflow orchestration
  - user task only needs execution path but agent must load analysis details, or only needs evidence path but agent must load workflow details

detector:
  - Phase 0 inventory
  - rg link / route audit
  - document-sizing check

loaded source / contract:
  - workflow/README.md
  - analysis/README.md
  - architecture/ai-native-cognitive-execution-system.md
  - architecture/ai-native-cognitive-ecosystem-system.md
  - governance/document-sizing.md
  - enforcement/linked-updates.md
  - knowledge/runtime/routing-registry.yaml

runtime action:
  - keep canonical source in existing owner layer
  - create or update workflow / analysis thin index
  - create focused execution / evidence slices only when source content genuinely needs decomposition
  - update summary/routing links when needed

evidence:
  - validation scenarios
  - link audit
  - runtime refresh if routing/validation surfaces change
```

### Per-Surface Consumer 表

| Generated surface key | Named consumer(s) | Consumer 類型 |
|---|---|---|
| none in draft | manual plan execution only | `manual_activation` |

若後續新增 routing registry route、runtime projection 或 validation generated surface，Phase 3 必須更新本表，且不得留下 dead surface。

---

## Open Questions

- [x] **[resolved by Phase 0]** `workflow/` 與 `analysis/` 中哪些檔案已經是 functional monolith？是否需要優先處理 APK / software-delivery / travel-planning 等高頻 route？ → Phase 0 Inventory Record 已盤點：`development-process.md`(378, ~12 gate)、`execution-flow.md`(270)、`apk-analysis/artifact-gates.md`(575)、`travel-planning/execution-flow.md`(295) 為 multi-topic monolith；`examples/EXAMPLES.md`(528) 主題單一。Pilot 依 stakeholder 鎖定 software-delivery；APK / travel 延後。
- [ ] **[intentionally deferred pending runtime validation]** `Cognitive Slice` 是否需要正式 glossary owner？刻意延後到 Phase 4 validation 證明 slice taxonomy 穩定後再決定是否升格為 framework vocabulary，避免 premature naming lock-in。**Interim operational wording**：在正式註冊前，文件一律用既有的 `loading surface` / `execution surface` / `evidence surface` 描述，不把 `Cognitive Slice` 當已確立詞彙散播。
- [x] **[resolved]** Slice 應落在現有 domain 目錄內，還是需要每個 domain 增加 `slices/` 子目錄？ → **決定：暫不新增 generic top-level `slices/` 或 domain-local `slices/` 子目錄。優先在既有 owner layer（`workflow/<domain>/`、`analysis/<domain>/`）內用 semantic filename 切分。等 Phase 4 validation evidence 證明確有 routing / 維護需求後再重評。**
- [x] **[resolved]** 是否需要保留 public-facing tutorial？ → **決定：保留 public-facing tutorial，但 tutorial 必須維持為 non-canonical projection layer。Tutorial 可引用 workflow / analysis slices，但不得複製 execution doctrine 或 evidence procedure 正文；canonical source 永遠在 owner layer。**
- [ ] **[brought forward → Phase 1/2 test-first, Phase 4 execute]** execution-only / evidence-only task 的「不交叉載入整包對方 layer」要用哪個 scenario fixture 機械驗證？ → 改為 test-first：在 Phase 1/2 先定義 `expected_load` / `forbidden_load` fixture 形狀與 Scenario A/B/C（見 Phase 4），Phase 4 才執行。驗證必須檢查**實際載入的 cognitive surface**，不只檢查 route 是否存在。
- [ ] **[new → Phase 1 budget, Phase 4 verify]** **Cross-Slice Dependency Explosion**：切片後是否會出現 recursive loading / dependency storms / hidden activation chains / retrieval fan-out（A 依賴 B 依賴 C…一個小任務被拉進整串 slice）？緩解：Phase 1 為 slice schema 加 **dependency budget**（`max_depth: 2`、`max_runtime_dependencies: 4`），Phase 4 scenario 驗證實際載入深度與廣度未超預算。

**處置優先序（external review 2026-05-29）**：先做 validation contract + dependency control（test-first，最早）→ 再做 tutorial projection 約束 + taxonomy naming（中段）→ 最後才做 glossary 正式註冊 + filesystem 子目錄決定（validation 穩定後）。理由：先用 runtime evidence 證明 slice 邊界站得住，再做難逆轉的命名與目錄結構決定。

---

## 完成條件

- [ ] Status 從 `draft` 更新為 `in-progress` / `completed`，並記錄完成日期。
- [x] Phase 0 完成 workflow / analysis source inventory 與 architecture compatibility preflight。
- [ ] Phase 1 定義 workflow / analysis slice taxonomy 與 owner-layer decision。
- [ ] Phase 2 完成至少一個 workflow 或 analysis surface 的 thin index 化。
- [ ] Phase 3 完成 loading rules、summary links、routing links 或明確 not applicable。
- [ ] Phase 4 補足 validation scenarios 或明確說明 doc-only trial 的 validation substitute。
- [ ] Phase 5 完成 linked updates、link audit、runtime refresh（若適用）。
- [ ] 若引入新 framework vocabulary，更新 glossary 或明確拒絕並說明理由。
- [ ] Plan Completion Closure：所有 checklist 完成後，執行 `plans/README.md` 的 archival / status / commit / push 閉環。

---

## Phase 0 — Architecture Compatibility Preflight

### Phase 0.0 — Open Questions 核對（公版，必填）

逐條核對本 plan §Open Questions，標記處置並回寫：

- [x] 已讀本 plan §Open Questions 全部條目
- [x] 對每條標記 `resolved` / `still-open` / `deferred`
- [x] `resolved` 的條目已同步勾選 / 附註於 §Open Questions
- [x] 若盤點新發現問題，已加入 §Open Questions（本輪無新問題）

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| 哪些檔是 functional monolith / 是否優先 APK·software-delivery·travel | resolved | Phase 0 Inventory Record 行數 + heading 盤點；pilot 鎖定 software-delivery |
| `Cognitive Slice` 是否註冊 glossary | intentionally deferred → Phase 4 | 待 validation 證明 taxonomy 穩定；過渡期用 `loading/execution/evidence surface`（external review 採納） |
| slice 落在現有 domain 還是新增 `slices/` 子目錄 | resolved | 暫不新增 `slices/`，用既有 owner layer + semantic filename；Phase 4 後重評 |
| 是否保留 public-facing tutorial | resolved | 保留但限定為 non-canonical projection，不複製 canonical source |
| execution-only / evidence-only 不交叉載入要用哪個 scenario fixture | brought forward → Phase 1/2 test-first, Phase 4 execute | expected_load/forbidden_load + Scenario A/B/C；驗證實際載入 surface |
| Cross-slice dependency explosion（新增） | new → Phase 1 budget, Phase 4 verify | dependency budget max_depth 2 / max_runtime_dependencies 4 |

目標：確認這不是「把文件搬一搬」，而是 current Gen 3 execution path / evidence path loading boundary hardening。

- [x] 盤點 candidate workflow files（行數見下表）。
- [x] 盤點 candidate analysis files（行數見下表）。
- [x] 檢查每個 candidate 是否超過 150 行且多主題。
- [x] 判定每個 candidate 的責任（見 Inventory Record 責任欄）。
- [x] 確認不修改 generated output / mirror copy 作為 canonical source。
- [x] 列出可能受影響的 routing registry、summary、README、validation scenario。
- [x] 確認 workflow 只引用 analysis，不複製完整 evidence method；analysis 只提供 evidence acquisition，不承擔 orchestration flow。

Phase 0 exit criteria：

- [x] Candidate inventory 完成。
- [x] 每個 candidate 有 owner-layer decision。
- [x] 若發現 source-of-truth 衝突，先更新本 plan，不進 Phase 1。

### Phase 0 Inventory Record（2026-05-29）

> 本輪 stakeholder 指示：只跑 Phase 0 盤點，不改內容；pilot 鎖定 `workflow/software-delivery/`。下列盤點皆為 read-only 觀察。

#### Preflight（依 `plans/README.md` §Architecture Compatibility Preflight）

| 欄位 | 內容 |
|---|---|
| Trigger | 開始執行本 plan 的 Phase 0（inventory only） |
| Checked sources | `plans/README.md` preflight、`enforcement/dependency-reading.md`、`workflow/software-delivery/README.md`、各 candidate heading、`knowledge/runtime/routing-registry.yaml` route.workflow.software-delivery、`workflow/software-delivery/{execution-flow,artifact-gates}.yaml` header |
| Conflicts | 無架構衝突。發現 1 個 source-of-truth 約束（見下「Source-of-truth 邊界」），須在 Phase 2 切片時遵守，但不阻擋 Phase 0。 |
| Decision | proceed — Phase 0 inventory 完成；依 stakeholder 指示在 Phase 1 前停止 |
| Validation | read-only 觀察 + `wc -l` 行數 + heading 結構；無檔案內容變更 |

#### Workflow candidates

| 檔案 | 行數 | >150 & multi-topic | 主責任 | Owner-layer decision |
|---|---|---|---|---|
| `workflow/software-delivery/development-process.md` | 378 | 是（~12 個 Gate：Default Flow / Required Contracts / Product Brief / Change Intake / Contract Governance / Traceability / BDD / Test Strategy / Embedded / Missing Info / Backfill / DoR-DoD） | execution order + 多個 artifact/governance gate 混在一檔 | **Pilot 首選**：thin index + 個別 gate slice。canonical 仍在 workflow 層 |
| `workflow/software-delivery/execution-flow.md` | 270 | 是（Step 1-9 + Surgical Changes Rules + Perf Gate + Test Strategy） | execution-order core 混入 surgical-changes caveat 與 perf gate | 切 execution-order slice + caveat/closure slice；為 routing primary_source，動它需同步 registry |
| `workflow/software-delivery/examples/EXAMPLES.md` | 528 | 行數超標但主題單一（4 組行為範例：Think Before Coding / Simplicity / Surgical / Goal-Driven） | workflow examples slice（已獨立在 `examples/` 子目錄） | 已是 de-facto example slice；可能只需瘦身索引，不必再拆 |
| `workflow/apk-analysis/artifact-gates.md` | 575 | 是 | apk artifact-gate（非 pilot） | 延後；非本輪 pilot |
| `workflow/travel-planning/execution-flow.md` | 295 | 是 | travel execution-order（非 pilot） | 延後；非本輪 pilot |

#### Analysis candidates

| 檔案 | 行數 | >150 & multi-topic | 主責任 | Owner-layer decision |
|---|---|---|---|---|
| `analysis/travel/sources-and-tools.md` | 312 | 是（~20 個來源類別） | evidence-acquisition（source catalog） | 延後；非 pilot。屬 evidence-acquisition slice |
| `analysis/apk/tools-and-failures.md` | 169 | 是（基礎工具 + 失敗判讀 + 命令模板 + 自動化安全邊界混合） | tool-procedure + failure/caveat 混合 | 延後；切 tool-procedure / failure-caveat 兩 slice |
| `analysis/apk/traffic-triage.md` | 178 | 是 | triage/observation | 延後；非 pilot |
| `analysis/apk/workflows/frida-hook-flow.md` | 190 | 是（步驟流程） | evidence-acquisition step procedure | 延後；非 pilot |
| `analysis/apk/workflows/media-hls-analysis-flow.md` | 165 | 是（步驟流程） | evidence-acquisition step procedure | 延後；非 pilot |

#### Source-of-truth 邊界（Phase 2 必須遵守）

- `workflow/software-delivery/execution-flow.yaml` 與 `artifact-gates.yaml` header 標 `source_markdown: <對應 .md>`、`owner_layer: workflow`。**`.md` 是 canonical prose，`.yaml` 是衍生 executable contract**。
- 切片若改動 `.md` 結構，必須檢查 `.yaml` 的 `source_markdown` 對映是否仍正確，並評估是否需 `ai-skill runtime compile + refresh`（Phase 3/5 處理）。
- 不得把 `.yaml` 或任何 generated surface 當 canonical prose 來切。

#### 受影響的 routing / summary / README

- `knowledge/runtime/routing-registry.yaml` → `route.workflow.software-delivery`：`primary_source = execution-flow.md`，`required_dependencies` 含 `development-process.md`、`artifact-gates.md` 等。動 pilot 檔需在 Phase 3 同步此 route。
- `knowledge/summaries/development-guidance.md`：對應 summary，Phase 3 需評估更新。
- `workflow/software-delivery/README.md`：已是 thin-ish 入口（149 行），列出各子檔導航；Phase 2 切片後需更新導航連結。

#### workflow↔analysis 邊界確認

- `workflow/software-delivery/README.md` §「與既有層的關係」明確：workflow 引用 `analysis/development-guidance/`、`analysis/repo/`、`intelligence/`，不複製其正文 → 符合 plan 對「workflow 只引用 analysis」的要求。
- analysis/apk 與 analysis/travel 提供 evidence acquisition / tool procedure，未承擔 workflow orchestration → 邊界乾淨。

#### Phase 0 結論

Phase 0 inventory 完成，無阻擋性架構衝突。Pilot 收斂為 **`development-process.md`（首選，gate 最多）+ `execution-flow.md`（execution-order vs caveat 分離）**，`examples/EXAMPLES.md` 可能僅需瘦身。依 stakeholder 指示在此停止，不進 Phase 1。下一步需 stakeholder 同意 Phase 1 taxonomy 與 glossary 決定後再繼續。

## Phase 1 — Cognitive Slice Taxonomy

目標：定義 workflow / analysis slice 邊界，不先改內容。

> **External review guardrails（2026-05-29 採納）**：本 phase 採 type+tags 收斂模型、granularity 原則與三層邊界規則，避免 taxonomy explosion / over-fragmentation。詳見 §Decision Rationale 風險表對應列。

- [ ] 定義最小 slice schema：
  - purpose
  - `type`（primary，**只允許 4 值**：`execution` / `evidence` / `examples` / `failure`）
  - `tags`（secondary，自由標註：artifact-gate / closure / handoff / templates / observation-triage / tool-procedure / domain-specific / extraction-to-intelligence …）
  - load_when
  - do_not_load_when
  - owner_layer
  - `layer_justification`（**歸層的 falsifiable 理由**，必須通過該層 membership predicate；見下方 placement 驗證規則）
  - `evidence_refs`（**僅 intelligence 層必填**：≥2 個獨立、已驗證的 analysis 觀察 / failure case 指標，且須可解析）
  - canonical_source
  - dependencies
  - `dependency_budget`（**heuristic default + override**，非 rigid governance；見下方規則）
  - summary_path
  - validation_signal
- [ ] **Dependency budget 規則（heuristic, not rigid）**：用 default + complexity override，避免 governance rigidity（external review 風險1，2026-05-29）。不同 surface complexity 合理上限不同（small CRUD fix 可能 depth1/deps2；deployment debugging 可能 depth3/deps6），所以**不把單一數字當硬門檻**：

  ```yaml
  dependency_budget:
    default:
      max_depth: 2
      max_runtime_dependencies: 4
    override_when:
      task_complexity: high   # 例：deployment debugging → 放寬至 depth 3 / deps 6
  ```

  超出 default 須在 slice / scenario 註記 `task_complexity: high` 理由；超出 override 上限才回 plan 重評。目的：防 recursive loading / dependency storms / retrieval fan-out（見 §Open Questions Cross-Slice Dependency Explosion），同時不讓 budget 變僵化治理。
- [ ] **type+tags 收斂規則**：primary `type` 固定 4 種，不得擴張為 first-class taxonomy；其餘責任一律降為 `tags`。新需求預設加 tag，不加 type。任何想新增第 5 個 primary type 的提議都需回到本 plan 重新評估。
- [ ] **Granularity 原則**：slice 最小單位 = **能獨立完成一個 cognitive phase**（例如 software-delivery 的 Requirement Intake / Implementation / Validation），**不是** step（Step1/Step2）也不是 concept。判準：該 slice 載入後 agent 能完成一個自足的認知階段而不需瘋狂 cross-reference。Phase 2 切片時逐個 slice 用此判準把關。
- [ ] **三層邊界規則（codify）**：`workflow` = 「要做什麼順序」；`analysis` = 「如何取得與驗證證據」；`intelligence` = 「為何這種模式長期有效 / 失敗」。slice 歸層時用此三分法判定 owner_layer，三層不得混。
- [ ] 用三層邊界規則檢查 taxonomy 是否與 `workflow/analysis/intelligence/knowledge/runtime/governance` 重疊。
- [ ] **Examples suppression bias 規則**：`type: examples` 的 slice 預設 **`default_load: false`**，只在 `user_requested_examples` 或 `ambiguity_detected` 時載入（external review 風險2，2026-05-29）。理由：examples token 密度高、pattern 明顯，LLM 易先讀 examples 而非 canonical execution，造成 **example-driven loading contamination / override doctrine**。這對應 Watch-Out Wall 5（positive-activation bias）。
- [ ] **Extraction direction rule（analysis → intelligence 單向）**：`analysis` 只產出 `observations / signals / evidence`；`intelligence` **只接受 validated repeated patterns**（external review 風險3，2026-05-29）。heuristic / tradeoff / anti-pattern / routing-heuristic 這類灰區內容，未經重複驗證前留在 analysis，不得直接倒進 intelligence，避免 intelligence 變 random thought dumping ground。slice 標 `tags: extraction-to-intelligence` 僅代表「候選」，升層需有 validation 證據。
- [ ] **Placement 可驗證規則（falsifiable membership predicate）**：歸層不是 honor-system 標籤，每個 slice 的 `layer_justification` 必須通過二選一判準，否則判定 misplacement：
  - **analysis membership test**：內容回答「如何取得 / 驗證證據」，且為 task-instance 級的 observation / signal / evidence；**不得**斷言跨實例通則。
  - **intelligence membership test**：內容是一個 generalization，**且** `evidence_refs` 含 ≥2 個獨立、已驗證、可解析的來源。少於 2 個或無法解析 → 判定 premature promotion，**強制退回 analysis**。
  - 「validated repeated pattern」即以 `evidence_refs` 數量 + 可解析性作為操作型定義（機械可檢查）。
  - **限制聲明**：歸層終究是語意分類，無完全機械 oracle；本規則目標是「讓 misplacement 可偵測、可逆、便宜修正」（falsifiable 判準 + evidence_refs gate + Phase 4 負向 scenario + contamination 探針），非「證明每次都放對」。升 / 降層走 audit trail，可逆。
- [x] 是否新增 domain-local `slices/` 子目錄 → **已於 §Open Questions resolved：暫不新增 generic / domain-local `slices/`，優先在既有 owner layer 內用 semantic filename 切分；Phase 4 validation 後重評。** 本 phase 只需確認 pilot 切分落在既有 `workflow/software-delivery/` 內。
- [ ] 評估命名候選並選定**過渡期 operational wording**（`loading/execution/evidence surface`）；**正式 glossary 註冊延後至 Phase 4**（見 §Open Questions 與 §Glossary Impact）。評估 `capability surface` / `cognitive surface` / `execution surface`（review 觀點：slice 易讓人聯想 arbitrary chunk / static partition，但本質是 routable cognition surface），記錄理由但不在本 phase 鎖定 framework vocabulary。

Phase 1 exit criteria：

- [ ] Taxonomy 不重複 canonical source，且 workflow / analysis 邊界清楚。
- [ ] primary `type` 恰為 4 種，其餘為 tags（type+tags 收斂規則成立）。
- [ ] 每個 slice 有明確 `load_when` 和 `do_not_load_when`。
- [ ] Granularity 原則與三層邊界規則已寫入 taxonomy 文件。
- [ ] Examples suppression bias 規則（`type: examples` 預設 `default_load: false`）與 extraction direction rule（analysis→intelligence 單向，只收 validated repeated patterns）已寫入 taxonomy 文件。
- [ ] Placement 可驗證規則已寫入：每個 slice 有 `layer_justification` 並通過該層 membership predicate；intelligence slice 的 `evidence_refs` ≥2 且可解析（否則退回 analysis）。
- [ ] 每個 slice schema 含 `dependency_budget`（heuristic default 2 / 4 + `override_when: task_complexity=high`，非 rigid 硬門檻）。
- [ ] **Test-first validation target 已草擬**：Phase 4 fixture 形狀（`expected_load` / `forbidden_load` / `dependency_budget`）與 Scenario A/B/C 的 expected/forbidden 清單已先寫出，待 Phase 4 執行。
- [ ] Glossary / naming decision 已記錄（含是否改用 surface 命名）。

## Phase 2 — Thin Index + Focused Slices

目標：把至少一個 monolithic workflow 或 analysis surface 改成 thin index。

- [ ] 選定 pilot surface；優先從 `workflow/` 與 `analysis/` 各選一個，若 scope 過大則先選一個高頻 route。
- [ ] 將 pilot surface 的正文分為：
  - index / navigation
  - execution-order 或 evidence-acquisition core
  - artifact gates / tool procedure
  - examples / templates
  - caveats / failure notes
- [ ] 父層 index 必須說明：
  - 何時讀哪個 workflow / analysis slice
  - execution-only 任務應避免讀哪些 analysis heavy slices
  - evidence-only 任務應避免讀哪些 workflow heavy slices
  - governance / runtime / intelligence 的 owner source 在哪
- [ ] 不在 index 複製 runtime contract、governance rule、analysis method 或 workflow 正文。
- [ ] 保留舊入口兼容說明或 redirect，避免 links 斷裂。

Phase 2 exit criteria：

- [ ] Pilot surface 不再同時承擔 execution order / evidence method / examples / artifacts / caveats 多重責任。
- [ ] **每個抽出的 slice 通過 granularity 判準**（Phase 1）：是一個能獨立完成的 cognitive phase，不是 step / concept，不需瘋狂 cross-reference。
- [ ] 每個抽出的 slice 仍能回連 canonical source。
- [ ] Document-sizing check 通過。

## Phase 3 — Loading Rules, Summary, Routing

目標：讓 agent 能按 task intent 載入 workflow / analysis slices。

- [ ] 為 pilot slices 建立或更新 summary（若適用）。
- [ ] 檢查 `knowledge/runtime/routing-registry.yaml` 是否需要新增 / 修改 route。
- [ ] 若新增 route，必須同時定義 named consumer 或 `manual_activation` reason。
- [ ] **Hierarchical routing 規則**：新 route 採樹狀命名（`workflow.software-delivery.execution`），**不得**鋪成 flat route universe（`workflow.execution.api` / `workflow.execution.refactor` / `analysis.apk.network` / `analysis.apk.hls` … 平攤）。理由：防 route inflation 讓 routing-registry 本身變成第二個 monolith / giant cognition graph（external review 風險4，2026-05-29）。新增 route 前先確認能掛在既有層級節點下，不是平鋪新 leaf。
- [ ] 為常見 intent 建立 loading guidance：
  - workflow execution
  - artifact validation
  - evidence acquisition
  - tool procedure lookup
  - failure / caveat diagnosis
  - mixed workflow + analysis task
- [ ] 明確記錄 suppression guidance：哪些任務不應載入 examples / tool procedures / artifact gates / failure caveats / Gen 4 heavy slices。**examples slice 預設 suppress（`default_load: false`），只在 user 明確要求範例或偵測到 ambiguity 才載入**（對齊 Phase 1 examples suppression bias 規則）。

Phase 3 exit criteria：

- [ ] 小任務可走 index / summary，不需整份 workflow 或 analysis surface。
- [ ] 大任務能找到需要的 source。
- [ ] 新增 route 為 hierarchical（樹狀），無 flat route universe；examples 預設 suppress。
- [ ] 無 dead route / dead generated surface。

## Phase 4 — Validation Scenarios

目標：用 scenarios 防止切片化只停在文件整理。

> **Test-first 約定（external review 2026-05-29 採納）**：本 phase 的 acceptance contract 不是 Phase 4 才設計，而是在 **Phase 1/2 就先寫好 fixture 形狀與 Scenario A/B/C 的 expected/forbidden 清單**（test-first target），Phase 4 只負責執行與蒐證。驗證必須斷言**實際載入的 cognitive surface**，不能只檢查 route 是否存在。

**Fixture 形狀（每個 scenario 一份）**：

```yaml
scenario: <id>
task_intent: <描述任務>
expected_load:        # 必須出現在載入集合
  - <surface/slice path>
forbidden_load:       # 必須不出現在載入集合
  - <surface/slice path>
dependency_budget:    # 對齊 Phase 1 slice schema（heuristic default + override）
  default: { max_depth: 2, max_runtime_dependencies: 4 }
  override_when: { task_complexity: high }   # 高複雜任務（如 Scenario C）放寬
```

- [ ] **Scenario A（execution-only）**：小型 API validation 變更。`expected_load` = software-delivery execution-order slice + 對應 artifact-gate slice；`forbidden_load` = full analysis / tool-procedure surface、examples、Gen 4 heavy section。
- [ ] **Scenario B（evidence-only）**：分析 APK 網路行為。`expected_load` = analysis evidence-acquisition / tool-procedure slice；`forbidden_load` = full workflow execution-flow / artifact-gate surface。
- [ ] **Scenario C（mixed）**：debug 失敗的 deployment pipeline。`expected_load` = workflow execution slice + 特定 analysis failure/caveat slice；`forbidden_load` = unrelated examples / 其他 domain slice / Gen 4 vision section。
- [ ] **Scenario D（placement / misplacement 負向驗證）**：故意放一條「無 evidence 或 evidence_refs < 2 的 heuristic」標成 intelligence，斷言 placement predicate **擋下並要求退回 analysis**（failure-derived validation）。同時驗證一條正確 analysis 證據 slice 的 `layer_justification` 通過 analysis membership test。
- [ ] **Contamination 作為 misplacement 間接探針**：明確記錄 Scenario B/C 的 `forbidden_load` 同時承擔 placement 驗證——若一條本該是 analysis 證據的 slice 被誤標成 intelligence doctrine，會在 evidence-only / mixed 任務的 `forbidden_load` 洩漏出來，藉此抓出歸層錯誤。
- [ ] 每個 scenario 斷言實際載入集合滿足 `expected_load` ⊆ loaded、`forbidden_load` ∩ loaded = ∅，且載入深度/廣度未超 `dependency_budget`。
- [ ] 若 Phase 3 改 runtime/routing source，執行 `ai-skill runtime refresh` 或適用 validator。
- [ ] 記錄 scenario evidence（實際 loaded surface 清單，非僅 route 宣告）。

Phase 4 exit criteria：

- [ ] Scenario A/B/C 全部 PASS（含 expected_load / forbidden_load / dependency_budget 三項斷言），或 plan 明確降級為 doc-only trial 並寫出下一階段 runtime validation plan。
- [ ] Scenario D PASS：placement predicate 擋下無證據的 intelligence 升層、放行正確 analysis slice；確認 placement 誤放可被偵測。

## Phase 5 — Linked Updates + Closure

目標：完成全庫一致性，不留下半套入口。

- [ ] 更新受影響 README / architecture / knowledge / workflow / analysis / ai-tools links。
- [ ] 執行 link / reference audit（例如搜尋舊入口 path / title）。
- [ ] 若新增 glossary terms，更新 glossary 並檢查 glossary impact。
- [ ] 若改 routing / validation / runtime source，執行 runtime compile / refresh / validate。
- [ ] 更新本 plan 狀態與完成日期。
- [ ] 執行 Plan Completion Closure，完成 archive / commit / push / readback。

Phase 5 exit criteria：

- [ ] `git status --short --branch` clean。
- [ ] `git log origin/<branch>..HEAD` 為空，或明確記錄未推送狀態與使用者授權需求。

---

## Stakeholder 同意項目

- [ ] 同意「現在 Gen 3 先做 workflow / analysis 切片化，Gen 4 再做 ecosystem orchestration」。
- [ ] 同意 workflow pilot surface 與 analysis pilot surface 的選擇。
- [ ] 同意是否正式引入 `Cognitive Slice` vocabulary。
- [ ] 同意是否新增 domain-local slice 子目錄；預設不新增 top-level layer。
- [ ] 同意 validation scenarios 的完成門檻。

---

## 與其他 plans 的關係

- [`2026-05-27-1557-tool-runtime-signal-economics-integration.md`](2026-05-27-1557-tool-runtime-signal-economics-integration.md)：本 plan 提供 cleaner workflow / analysis loading units；該 plan 後續可用 economics 判斷何時載入 slices。
- [`2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md`](2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md)：本 plan 不處理 fitness / optimization memory，但避免 monolithic workflow / analysis surfaces 成為未來 optimization target 的污染來源。
- [`architecture/ai-native-cognitive-execution-system.md`](../../architecture/ai-native-cognitive-execution-system.md)：本 plan 屬於 Gen 3 current hardening。
- [`architecture/ai-native-cognitive-ecosystem-system.md`](../../architecture/ai-native-cognitive-ecosystem-system.md)：本 plan 是 Gen 4 prerequisite，不是 Gen 4 capability implementation。
