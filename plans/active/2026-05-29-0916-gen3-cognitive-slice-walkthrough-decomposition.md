# Gen3 Cognitive Slice Walkthrough Decomposition

**Status**: `draft`
**世代**：Gen 3 current hardening；Gen 4 prerequisite
**建立日期**：2026-05-29
**最後更新**：2026-05-29

> 本 plan 的核心決策：walkthrough / large onboarding surface 的切片化應在 Gen 3 完成，不等 Gen 4。Gen 3 先建立清楚的 cognitive unit、routing boundary、summary-first loading 與 validation；Gen 4 才在此基礎上做 dynamic activation / ecosystem orchestration。

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

因此「walkthrough 太肥」不是單純文件美觀問題，而是 runtime cognition 問題。若一個 walkthrough 同時承載 framework overview、governance、runtime rules、workflow examples、failure doctrine、philosophy 與 onboarding，它會造成：

- **Semantic overload**：agent 無法判斷當下真正需要哪一段認知。
- **Retrieval contamination**：小任務載入過多 governance / failure / orchestration context。
- **Compression degradation**：巨型文件被反覆摘要後，留下抽象口號，失去可執行細節。
- **Routing ambiguity**：entrypoint、primary_source、summary 與 source-of-truth 邊界不清。
- **Gen 4 migration debt**：未來 dynamic activation graph 會繼承 monolith boundary，變成 monolithic cognitive ecosystem。

現有架構其實已經支持這個方向：

- Gen 3 architecture 明定 canonical navigation 入口不是 self-contained spec。
- `knowledge/summaries/` 採 summary-first loading。
- `knowledge/runtime/routing-registry.yaml` 已是 task intent → primary source 的路由層。
- `governance/document-sizing.md` 要求超過 150 行且多主題的高頻文件拆成 index + focused child files。
- Gen 4 vision 的 Watch-Out List 已警告 workflow inflation、ecosystem boundary inflation、positive-activation bias。

所以這項工作應視為 **Gen 3 runtime hardening**，不是 Gen 4 才開始的願景工作。

### Decision

在 Gen 3 先完成 **Cognitive Slice Walkthrough Decomposition**：

1. 盤點現有 high-frequency onboarding / walkthrough / entry surfaces。
2. 將巨型 walkthrough 轉為 thin routing index。
3. 把內容拆成 focused cognitive slices：
   - foundation slice
   - runtime mode slice
   - workflow slice
   - failure intelligence slice
   - domain intelligence slice
   - governance slice
   - Gen 4 boundary slice
4. 為每個 slice 定義：
   - purpose
   - load_when
   - source-of-truth owner
   - dependencies
   - summary target
   - validation signal
5. 更新 routing / summary / README links，讓 agent 先讀 index，再按 task intent 載入必要 slices。
6. 補 validation scenarios，確認小任務不會載入整個 walkthrough，大任務能載入正確 slices。

Gen 4 只承接後續：

- dynamic cognitive retrieval
- activation graph
- economics-driven loading
- suppression layer
- telemetry / fitness feedback
- adaptive orchestration

本 plan 不實作 autonomous optimizer、self-modifying governance、telemetry DB 或 Gen 4 ecosystem layer。

### Alternatives Considered

- A. 等 Gen 4 再處理：reject。Gen 4 的 activation / orchestration 需要清楚 cognitive unit；若現在不拆，未來只會把 monolith 帶進 ecosystem。
- B. 只把 walkthrough 縮短：reject。問題不是字數，而是多種責任混在同一 loading boundary。
- C. 直接建立完整 Gen 4 dynamic retrieval system：reject。會提前跨入 ecosystem orchestration，違反 Gen 4 vision 對 current / vision boundary 的限制。
- D. Gen 3 先做 slice + index + loading rules：accept。這是 current runtime 可執行的收斂工作，也能替 Gen 4 預留乾淨 substrate。

### Why Not an ADR Yet

此 plan 是文件結構、loading boundary 與 routing hardening。它可能導出新的 architecture decision，但目前尚未證明：

- cognitive slice taxonomy 是否穩定；
- 哪些 slices 應成為 long-lived owner layer；
- routing registry 是否需要新增 route 類型；
- validation scenarios 是否足以防止 retrieval contamination；
- `Cognitive Slice` 是否應成為正式 framework vocabulary。

完成後若 slice taxonomy 成為跨世代、跨工具、難逆轉的基礎 decision，再評估 ADR promotion。

### ADR Promotion Criteria（completed 時驗證）

- [ ] High-frequency walkthrough / onboarding surfaces 已完成 inventory。
- [ ] 至少一個 oversized / multi-topic entry surface 已轉為 thin index + focused slices。
- [ ] 每個 slice 都有 `load_when` / owner / dependency / validation signal。
- [ ] Summary-first loading path 已更新或明確標記 not applicable。
- [ ] Routing registry / README / architecture links 已檢查並完成必要同步。
- [ ] 至少 3 個 validation scenarios 覆蓋 small task / architecture task / failure-learning task。
- [ ] Open Questions 全部解決或轉為明確 follow-up plan。
- [ ] 實際 agent session 能證明不再為小任務載入整份 walkthrough。

### Consequences

#### 正面

- 降低 context loading 成本與 retrieval contamination。
- 讓 Gen 3 的 runtime routing 更接近「讀對 slice」而不是「讀一本大書」。
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
| 舊 walkthrough links drift | Phase 4 必須做 rg link audit + routing registry check |
| 新 taxonomy 與既有 `workflow/analysis/intelligence/knowledge` 重疊 | Phase 0 先做 owner-layer preflight；每個 slice 只導航，不重定義 canonical source |

### Glossary Impact

Glossary Impact: yes.

Candidate framework vocabulary:

- `Cognitive Slice`：可被獨立載入、驗證、路由的最小認知單元。
- `Retrieval Boundary`：agent 在某任務中應停留的載入邊界。
- `Thin Walkthrough Index`：只負責 navigation / loading guidance，不承載全部正文的 walkthrough 入口。

Phase 1 必須決定是否註冊到 `knowledge/glossary/ai-skill.md`，或改用既有詞彙避免 vocabulary inflation。

### Watch-Out List Citation

本 plan 對應 Gen 4 vision [`architecture/ai-native-cognitive-ecosystem-system.md`](../../architecture/ai-native-cognitive-ecosystem-system.md) §Watch-Out List：

- Wall 1. Discovery confused with Activation：本 plan 只做 Gen 3 loading boundary，不把 routing-registry 假裝成 activation graph。
- Wall 2. Workflow inflation：避免 walkthrough / workflow 越寫越大，改用 focused slices。
- Wall 3. Ecosystem boundary inflation：不新增 ecosystem layer；只整理 current source surfaces。
- Wall 5. Positive-activation bias：Phase 2 必須定義 `do_not_load_when` 或 suppression note，避免「多載入就是更完整」的錯覺。

---

## Runtime Execution Path

### Doc-only Trial Statement + Graduation

目前狀態：**Gen 3 documentation / routing hardening plan**。

本 plan 第一版不新增 runtime generated surface、不新增 commit-msg validator、不建立 Gen 4 activation graph。若 Phase 3 決定要更新 `knowledge/runtime/routing-registry.yaml` 或新增 validation scenario，才進入 runtime-linked execution。

Graduation condition：

| Graduation Phase | Acceptance signal |
|---|---|
| Phase 2 完成 | walkthrough index + slices 形狀可讀，且不重複 canonical source |
| Phase 3 完成 | loading rules / routing links 有 named consumer 或明確 manual activation |
| Phase 4 完成 | scenarios 證明小任務不載入整個 walkthrough |

### Runtime Owner

- Current owner layer：`architecture/` + `knowledge/` + affected source layer README
- Loading / summary owner：`knowledge/summaries/`
- Routing owner：`knowledge/runtime/routing-registry.yaml`（若需要）
- Validation owner：`validation/scenarios/`（若需要）
- Governance owner：`governance/document-sizing.md` + `enforcement/linked-updates.md`

### Trigger Flow

```text
event_or_signal:
  - high-frequency walkthrough / onboarding file exceeds document-sizing threshold
  - file mixes framework / governance / runtime / workflow / failure / examples
  - user task only needs one topic but agent must load whole walkthrough

detector:
  - Phase 0 inventory
  - rg link / route audit
  - document-sizing check

loaded source / contract:
  - architecture/ai-native-cognitive-execution-system.md
  - architecture/ai-native-cognitive-ecosystem-system.md
  - governance/document-sizing.md
  - enforcement/linked-updates.md
  - knowledge/runtime/routing-registry.yaml

runtime action:
  - keep canonical source in existing owner layer
  - create or update thin index
  - create focused slices only when source content genuinely needs decomposition
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

- [ ] 目前 repository 中真正承擔 walkthrough 角色的是哪個檔案或哪組入口？是否是 `README.md` / `architecture/*` / `CORE_BOOTSTRAP.md` / tool adapter 之一，而不是名為 walkthrough 的檔案？
- [ ] `Cognitive Slice` 是否需要正式 glossary owner，還是先用 `summary-first loading` / `focused source` 既有詞彙即可？
- [ ] Slice 應落在現有 layer 內，還是需要 `walkthrough/` 目錄？預設偏向不新增 top-level layer，除非 Phase 0 證明現有 layer 無法承載。
- [ ] 是否需要保留一份 public-facing tutorial？若需要，它應只引用 slices，不複製 canonical source。
- [ ] 小任務「不載入整份 walkthrough」要用哪個 scenario fixture 機械驗證？

---

## 完成條件

- [ ] Status 從 `draft` 更新為 `in-progress` / `completed`，並記錄完成日期。
- [ ] Phase 0 完成 source inventory 與 architecture compatibility preflight。
- [ ] Phase 1 定義 slice taxonomy 與 owner-layer decision。
- [ ] Phase 2 完成至少一個 walkthrough / onboarding surface 的 thin index 化。
- [ ] Phase 3 完成 loading rules、summary links、routing links 或明確 not applicable。
- [ ] Phase 4 補足 validation scenarios 或明確說明 doc-only trial 的 validation substitute。
- [ ] Phase 5 完成 linked updates、link audit、runtime refresh（若適用）。
- [ ] 若引入新 framework vocabulary，更新 glossary 或明確拒絕並說明理由。
- [ ] Plan Completion Closure：所有 checklist 完成後，執行 `plans/README.md` 的 archival / status / commit / push 閉環。

---

## Phase 0 — Architecture Compatibility Preflight

目標：確認這不是「把文件搬一搬」，而是 current Gen 3 loading boundary hardening。

- [ ] 盤點 candidate files：
  - [ ] root `README.md`
  - [ ] `CORE_BOOTSTRAP.md`
  - [ ] `architecture/ai-native-cognitive-execution-system.md`
  - [ ] `architecture/ai-native-cognitive-ecosystem-system.md`
  - [ ] `ai-tools/agent/*.md`
  - [ ] high-frequency workflow / enforcement / knowledge README
- [ ] 搜尋是否存在實際 `walkthrough` 命名檔案；若無，列出 functional walkthrough surfaces。
- [ ] 檢查每個 candidate 是否超過 150 行且多主題。
- [ ] 判定每個 candidate 的責任：
  - navigation only
  - canonical source
  - executable contract companion
  - public tutorial
  - tool adapter
- [ ] 確認不修改 generated output / mirror copy 作為 canonical source。
- [ ] 列出可能受影響的 routing registry、summary、README、validation scenario。

Phase 0 exit criteria：

- [ ] Candidate inventory 完成。
- [ ] 每個 candidate 有 owner-layer decision。
- [ ] 若發現 source-of-truth 衝突，先更新本 plan，不進 Phase 1。

## Phase 1 — Cognitive Slice Taxonomy

目標：定義 slice 邊界，不先改內容。

- [ ] 定義最小 slice schema：
  - purpose
  - load_when
  - do_not_load_when
  - owner_layer
  - canonical_source
  - dependencies
  - summary_path
  - validation_signal
- [ ] 建立 initial slice taxonomy：
  - foundation
  - runtime modes
  - workflow
  - failure intelligence
  - domain intelligence
  - governance
  - Gen 4 boundary
- [ ] 檢查 taxonomy 是否與 `workflow/analysis/intelligence/knowledge/runtime/governance` 重疊。
- [ ] 決定是否需要新增 `walkthrough/` 目錄；預設不新增，優先使用現有 layer + index。
- [ ] 決定 glossary 是否註冊 `Cognitive Slice`。

Phase 1 exit criteria：

- [ ] Taxonomy 不重複 canonical source。
- [ ] 每個 slice 有明確 `load_when` 和 `do_not_load_when`。
- [ ] Glossary decision 已記錄。

## Phase 2 — Thin Index + Focused Slices

目標：把至少一個 monolithic walkthrough/onboarding surface 改成 thin index。

- [ ] 選定 pilot surface。
- [ ] 將 pilot surface 的正文分為：
  - index / navigation
  - canonical source links
  - optional examples
  - moved / extracted slice content
- [ ] 父層 index 必須說明：
  - 何時讀哪個 slice
  - 小任務應避免讀哪些 heavy slices
  - architecture / governance / runtime 的 owner source 在哪
- [ ] 不在 index 複製 runtime contract、governance rule 或 workflow 正文。
- [ ] 保留舊入口兼容說明或 redirect，避免 links 斷裂。

Phase 2 exit criteria：

- [ ] Pilot surface 不再同時承擔 framework / governance / runtime / examples / philosophy 多重責任。
- [ ] 每個抽出的 slice 仍能回連 canonical source。
- [ ] Document-sizing check 通過。

## Phase 3 — Loading Rules, Summary, Routing

目標：讓 agent 能按 task intent 載入 slices。

- [ ] 為 pilot slices 建立或更新 summary（若適用）。
- [ ] 檢查 `knowledge/runtime/routing-registry.yaml` 是否需要新增 / 修改 route。
- [ ] 若新增 route，必須同時定義 named consumer 或 `manual_activation` reason。
- [ ] 為常見 intent 建立 loading guidance：
  - debugging
  - architecture review
  - workflow execution
  - failure learning
  - Gen 4 planning
- [ ] 明確記錄 suppression guidance：哪些任務不應載入 governance / failure / Gen 4 heavy slices。

Phase 3 exit criteria：

- [ ] 小任務可走 index / summary，不需整份 walkthrough。
- [ ] 大任務能找到需要的 source。
- [ ] 無 dead route / dead generated surface。

## Phase 4 — Validation Scenarios

目標：用 scenarios 防止切片化只停在文件整理。

- [ ] 建立 small-task scenario：只問單一 task，期待不載入 full walkthrough。
- [ ] 建立 architecture-task scenario：需要 foundation + architecture slice，但不載入 unrelated domain slice。
- [ ] 建立 failure-learning scenario：需要 failure slice + enforcement source，但不載入 Gen 4 vision heavy section。
- [ ] 若 Phase 3 改 runtime/routing source，執行 `ai-skill runtime refresh` 或適用 validator。
- [ ] 記錄 scenario evidence。

Phase 4 exit criteria：

- [ ] 至少 3 個 validation scenarios PASS，或 plan 明確降級為 doc-only trial 並寫出下一階段 runtime validation plan。

## Phase 5 — Linked Updates + Closure

目標：完成全庫一致性，不留下半套入口。

- [ ] 更新受影響 README / architecture / knowledge / workflow / ai-tools links。
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

- [ ] 同意「現在 Gen 3 先做切片化，Gen 4 再做 ecosystem orchestration」。
- [ ] 同意 pilot surface 的選擇。
- [ ] 同意是否正式引入 `Cognitive Slice` vocabulary。
- [ ] 同意是否新增 `walkthrough/` 目錄；預設不新增 top-level layer。
- [ ] 同意 validation scenarios 的完成門檻。

---

## 與其他 plans 的關係

- [`2026-05-27-1557-tool-runtime-signal-economics-integration.md`](2026-05-27-1557-tool-runtime-signal-economics-integration.md)：本 plan 提供 cleaner loading units；該 plan 後續可用 economics 判斷何時載入 slices。
- [`2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md`](2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md)：本 plan 不處理 fitness / optimization memory，但避免 monolithic walkthrough 成為未來 optimization target 的污染來源。
- [`architecture/ai-native-cognitive-execution-system.md`](../../architecture/ai-native-cognitive-execution-system.md)：本 plan 屬於 Gen 3 current hardening。
- [`architecture/ai-native-cognitive-ecosystem-system.md`](../../architecture/ai-native-cognitive-ecosystem-system.md)：本 plan 是 Gen 4 prerequisite，不是 Gen 4 capability implementation。
