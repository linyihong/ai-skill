# Gen3 Workflow and Analysis Cognitive Slice Decomposition

**Status**: `draft`
**世代**：Gen 3 current hardening；Gen 4 prerequisite
**建立日期**：2026-05-29
**最後更新**：2026-05-29

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
| 切太細導致維護成本 > token savings | Phase 1 使用 document-sizing threshold 與 high-frequency loading evidence 判斷 |
| 舊 workflow / analysis links drift | Phase 4 必須做 rg link audit + routing registry check |
| 新 taxonomy 與既有 `workflow/analysis/intelligence/knowledge` 重疊 | Phase 0 先做 owner-layer preflight；每個 slice 只導航，不重定義 canonical source |

### Glossary Impact

Glossary Impact: yes.

Candidate framework vocabulary:

- `Cognitive Slice`：可被獨立載入、驗證、路由的最小認知單元；本 plan 先落在 workflow execution slice 與 analysis evidence slice。
- `Retrieval Boundary`：agent 在某任務中應停留的載入邊界。
- `Thin Workflow/Analysis Index`：只負責 navigation / loading guidance，不承載全部 workflow 或 analysis 正文的入口。

Phase 1 必須決定是否註冊到 `knowledge/glossary/ai-skill.md`，或改用既有詞彙避免 vocabulary inflation。

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

- [ ] `workflow/` 與 `analysis/` 中哪些檔案已經是 functional monolith？是否需要優先處理 APK / software-delivery / travel-planning 等高頻 route？
- [ ] `Cognitive Slice` 是否需要正式 glossary owner，還是先用 `summary-first loading` / `focused source` 既有詞彙即可？
- [ ] Slice 應落在現有 domain 目錄內，還是需要每個 domain 增加 `slices/` 或 `guides/` 子目錄？預設不新增 top-level layer。
- [ ] 是否需要保留 public-facing tutorial？若需要，它應只引用 workflow / analysis slices，不複製 canonical source。
- [ ] execution-only / evidence-only task 的「不交叉載入整包對方 layer」要用哪個 scenario fixture 機械驗證？

---

## 完成條件

- [ ] Status 從 `draft` 更新為 `in-progress` / `completed`，並記錄完成日期。
- [ ] Phase 0 完成 workflow / analysis source inventory 與 architecture compatibility preflight。
- [ ] Phase 1 定義 workflow / analysis slice taxonomy 與 owner-layer decision。
- [ ] Phase 2 完成至少一個 workflow 或 analysis surface 的 thin index 化。
- [ ] Phase 3 完成 loading rules、summary links、routing links 或明確 not applicable。
- [ ] Phase 4 補足 validation scenarios 或明確說明 doc-only trial 的 validation substitute。
- [ ] Phase 5 完成 linked updates、link audit、runtime refresh（若適用）。
- [ ] 若引入新 framework vocabulary，更新 glossary 或明確拒絕並說明理由。
- [ ] Plan Completion Closure：所有 checklist 完成後，執行 `plans/README.md` 的 archival / status / commit / push 閉環。

---

## Phase 0 — Architecture Compatibility Preflight

目標：確認這不是「把文件搬一搬」，而是 current Gen 3 execution path / evidence path loading boundary hardening。

- [ ] 盤點 candidate workflow files：
  - [ ] `workflow/apk-analysis/artifact-gates.md`
  - [ ] `workflow/software-delivery/examples/EXAMPLES.md`
  - [ ] `workflow/software-delivery/development-process.md`
  - [ ] `workflow/travel-planning/execution-flow.md`
  - [ ] `workflow/software-delivery/execution-flow.md`
  - [ ] 其他超過 document-sizing threshold 的 workflow surfaces
- [ ] 盤點 candidate analysis files：
  - [ ] `analysis/travel/sources-and-tools.md`
  - [ ] `analysis/apk/workflows/frida-hook-flow.md`
  - [ ] `analysis/apk/traffic-triage.md`
  - [ ] `analysis/apk/tools-and-failures.md`
  - [ ] `analysis/apk/workflows/media-hls-analysis-flow.md`
  - [ ] 其他超過 document-sizing threshold 的 analysis surfaces
- [ ] 檢查每個 candidate 是否超過 150 行且多主題。
- [ ] 判定每個 candidate 的責任：
  - workflow execution order
  - workflow artifact gate / closure
  - workflow example / template
  - analysis evidence acquisition
  - analysis tool procedure
  - analysis triage / caveat
  - intelligence / governance / tool adapter spillover
- [ ] 確認不修改 generated output / mirror copy 作為 canonical source。
- [ ] 列出可能受影響的 routing registry、summary、README、validation scenario。
- [ ] 確認 workflow 只引用 analysis，不複製完整 evidence method；analysis 只提供 evidence acquisition，不承擔 orchestration flow。

Phase 0 exit criteria：

- [ ] Candidate inventory 完成。
- [ ] 每個 candidate 有 owner-layer decision。
- [ ] 若發現 source-of-truth 衝突，先更新本 plan，不進 Phase 1。

## Phase 1 — Cognitive Slice Taxonomy

目標：定義 workflow / analysis slice 邊界，不先改內容。

- [ ] 定義最小 slice schema：
  - purpose
  - load_when
  - do_not_load_when
  - owner_layer
  - canonical_source
  - dependencies
  - summary_path
  - validation_signal
- [ ] 建立 initial workflow slice taxonomy：
  - execution-order
  - artifact-gates
  - examples
  - templates
  - handoff / closure
  - validation
- [ ] 建立 initial analysis slice taxonomy：
  - evidence-acquisition
  - observation / triage
  - tool-procedure
  - failure / caveat
  - domain-specific method
  - extraction-to-intelligence
- [ ] 檢查 taxonomy 是否與 `workflow/analysis/intelligence/knowledge/runtime/governance` 重疊。
- [ ] 決定是否需要新增 domain-local `slices/` / `guides/` / `examples/` 子目錄；預設不新增 top-level layer，優先使用現有 layer + index。
- [ ] 決定 glossary 是否註冊 `Cognitive Slice`。

Phase 1 exit criteria：

- [ ] Taxonomy 不重複 canonical source，且 workflow / analysis 邊界清楚。
- [ ] 每個 slice 有明確 `load_when` 和 `do_not_load_when`。
- [ ] Glossary decision 已記錄。

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
- [ ] 每個抽出的 slice 仍能回連 canonical source。
- [ ] Document-sizing check 通過。

## Phase 3 — Loading Rules, Summary, Routing

目標：讓 agent 能按 task intent 載入 workflow / analysis slices。

- [ ] 為 pilot slices 建立或更新 summary（若適用）。
- [ ] 檢查 `knowledge/runtime/routing-registry.yaml` 是否需要新增 / 修改 route。
- [ ] 若新增 route，必須同時定義 named consumer 或 `manual_activation` reason。
- [ ] 為常見 intent 建立 loading guidance：
  - workflow execution
  - artifact validation
  - evidence acquisition
  - tool procedure lookup
  - failure / caveat diagnosis
  - mixed workflow + analysis task
- [ ] 明確記錄 suppression guidance：哪些任務不應載入 examples / tool procedures / artifact gates / failure caveats / Gen 4 heavy slices。

Phase 3 exit criteria：

- [ ] 小任務可走 index / summary，不需整份 workflow 或 analysis surface。
- [ ] 大任務能找到需要的 source。
- [ ] 無 dead route / dead generated surface。

## Phase 4 — Validation Scenarios

目標：用 scenarios 防止切片化只停在文件整理。

- [ ] 建立 execution-only scenario：只需要 workflow execution order，期待不載入 full analysis/tool-procedure surface。
- [ ] 建立 evidence-only scenario：只需要 analysis evidence acquisition，期待不載入 full workflow/artifact-gate surface。
- [ ] 建立 mixed workflow+analysis scenario：需要 workflow + specific analysis slice，但不載入 unrelated examples / caveats / Gen 4 vision heavy section。
- [ ] 若 Phase 3 改 runtime/routing source，執行 `ai-skill runtime refresh` 或適用 validator。
- [ ] 記錄 scenario evidence。

Phase 4 exit criteria：

- [ ] 至少 3 個 validation scenarios PASS，或 plan 明確降級為 doc-only trial 並寫出下一階段 runtime validation plan。

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
