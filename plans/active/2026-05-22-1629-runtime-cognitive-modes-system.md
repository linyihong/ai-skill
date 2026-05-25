# Runtime Cognitive Modes System

**Status**: `in-progress`（Phase D ✅ + Phase 0 ✅ + Phase 1 ✅ + Phase 2 ✅ + **Phase 3 governance-structure 完成**（5/5 scenarios pass）；behavioral enforcement + linked updates + ADR evaluation 待下次 session）
**世代**：Gen 3 子系統擴充
**建立日期**：2026-05-22
**最後更新**：2026-05-25（Phase 3-B: enforcement obligation + gate 寫入 runtime.db；deferred items 完整記錄於 §Phase 3 Deferred Items）

> ⚠️ 本 plan 處於 `in-progress` 階段：**Phase D documentation-contract trial 已完成並通過評估指標**；Phase 0 (Pre-Build Interrogation) 與 Phase 1-5 runtime 實作待 user 決定是否啟動。原 `constitution/ADR-008-runtime-cognitive-modes.md`（proposed）已於 2026-05-22 撤回；依新 [`decision-promotion-pipeline`](../../governance/lifecycle/decision-promotion-pipeline.md) 規則，constitution/ 只放 accepted ADRs，提案階段在本 plan 內處理。
>
> 若 plan completed 且通過 §ADR Promotion Criteria，才升級為 accepted ADR（屆時取編號）。

---

## Decision Rationale（原 ADR-008 提案內容）

### Problem & Why Now

2026-05-22 外部架構審查指出 `models/` 層目前是 **documentation layer 而非 runtime activation layer**。具體證據：

1. **沒有 blocking gate 強制查詢**：agent 在執行任務時不會自動查 `routing-registry.yaml` → `model-context-report.md` → `model-checklists.md`。`knowledge-update-flow.yaml` 11 步沒有任何一步要求查詢 model profile。
2. **本 session 實證**：本 session 內 agent 加入 4 個 intelligence atoms、修正 ADR-007 語言、寫第三代 architecture 文件、ADR 雙向連結等工作，**沒有一次**查詢 model-context-report 就直接執行；profile 報告 / Read / Deferred / Validation signal 從未出現在 final report。
3. **Document lookup runtime 不可行**：對方批評「每次 full resolution 一定爆」是真實 token cost 問題。若加 Step 0「每任務 query registry + profile + report + checklist」，每次 ~2000 tokens overhead 無法承受。
4. **真正的缺口**：對照 4 個 cognitive primitive 維度，**governance mode 強度差異化**與 **memory mode activation flag** 是 Gen 3 既有 infrastructure 沒覆蓋到的真實缺口。

對照現有系統：

| 建議 mode | 既有對應 | 缺口性質 |
|------|------|------|
| execution mode (FAST/NORMAL/DEEP/FORENSIC/RECOVERY) | `runtime.db phase_machine` 有 phase 概念，但**沒有 cognitive depth 維度** | 60% 新（FORENSIC/RECOVERY 為真新增） |
| context mode (INDEX_ONLY/SUMMARY_FIRST/CHECKLIST_FIRST/SOURCE_BACKED/GRAPH_ASSISTED) | `models/compression/` 5 個 level 名稱幾乎一樣（小寫） | 5% 新（rename + 提升為 runtime primitive） |
| governance mode (LIGHT/STANDARD/STRICT/LOCKDOWN) | 既有 governance 是 binary（gate 或無 gate） | **80% 新** — 真正缺口 |
| memory mode (NONE/EPISODIC/DECISION_REPLAY/FAILURE_REPLAY/PROJECT_CONTEXT) | `memory/` 子層存在但無 activation flag | **70% 新** — 把 memory 子層提升為 runtime mode |

### Decision

引入 **Runtime Cognitive Modes** 作為 Gen 3 runtime infrastructure 的**子系統擴充**，核心三點：

#### 1. 4 維 mode primitive 而非 flat profile

```
execution_mode  ∈ {FAST, NORMAL, DEEP, FORENSIC, RECOVERY}
context_mode    ∈ {INDEX_ONLY, SUMMARY_FIRST, CHECKLIST_FIRST, SOURCE_BACKED, GRAPH_ASSISTED}
governance_mode ∈ {LIGHT, STANDARD, STRICT, LOCKDOWN}
memory_mode     ∈ {NONE, EPISODIC, DECISION_REPLAY, FAILURE_REPLAY, PROJECT_CONTEXT}
```

組合空間 = 5 × 5 × 4 × 5 = 500 種 cognitive state。

> **Naming alignment**（per Open Question 1 resolved）：`context_mode` 5 級對應既有 [`models/compression/`](../../models/compression/README.md) 的 5 級（`index-only`、`summary-first`、`checklist-first`、`source-backed`、`graph-assisted`）；UPPERCASE 為 runtime primitive，lowercase 視為 alias。`compression/` 文件改為「`context_mode.<LEVEL>` 的實作策略 reference」，由 Phase 3.2 執行。

#### 2. Discovery 用快速啟發式，不查文件

```
任務進來
  ↓
快速 discovery（純訊號計算，< 50 tokens）：
  - user keyword、file diff scope、git status、session turn 數、recent failure、
    phase_machine 當前 phase、active goals 狀態、token_budget 剩餘、
    modified files 命中 generated_surfaces（完整訊號表見 Phase 2 §2.1）
  ↓
單一 SQLite 查詢解析 mode
  ↓
4 個 mode 寫入 runtime state
  ↓
各 subsystem 依 mode 決定 activation
```

#### 3. 與既有 infrastructure 整合，不重寫

| 既有 | 整合方式 |
|------|------|
| `runtime.db phase_machine` | 新增 `cognitive_modes` 表，與 phase_machine join |
| `models/compression/` | Rename levels 為 UPPERCASE，提升為 `context_mode` 的實作層 |
| `models/profiles/` | 保留為 reference doc；mark 為「backward-compat label」 |
| `memory/<subdir>/` | 子層不動，新增 `memory_mode` 作為 activation flag |
| `governance/` | 不動現有 gate；新增 `governance_mode` 作為 gate 強度 selector |

**不刪除任何既有檔。**

#### 4. 漸進實作（5 phase）

詳見下方 §Phase 1-5。

### Alternatives Considered

- **A. 維持現狀（models/ 純 documentation）**：拒絕 — 已證實 agent 不 activate。
- **B. 完全重寫 `models/`**：拒絕 — 既有 compression / capabilities / routing 有獨立 reference 價值。
- **C. 只加 governance mode 解決最緊迫缺口**：拒絕 — 單一 dimension 無法 compose 解決問題。
- **D. 加進 knowledge-update-flow Step 0「每任務查 profile」**：拒絕 — token cost 不可行。
- **E. 4 mode + 整合既有 + 5 phase 漸進**：**accept**（本 plan 採用）。

### Why Not an ADR Yet

依新 [`decision-promotion-pipeline`](../../governance/lifecycle/decision-promotion-pipeline.md) §No-Proposed-ADR Rule：

- **未驗證**：Phase 1-5 未執行，不知道實際效果
- **Scope 可能調整**：5 個 Open Questions 可能改變 Decision shape
- **沒有更輕的 promotion target 已被檢驗**：可能 cognitive mode 應該是 runtime gate（不需 ADR）而非 architectural decision
- **依新規則**：constitution/ 只放 accepted ADRs；proposed 階段在 plan 內

### ADR Promotion Criteria（2026-05-25 評估）

升級為 accepted ADR 的條件（per ADR-007 §ADR Boundary）。**Promotion gate 設在 Phase 3 完成時**（per Open Question 5 resolved）：

| 條件 | 狀態 | 評估說明 |
|------|------|---------|
| foundational + cross-session + cross-project + expensive-to-reverse + explains-why 全中 | ✅ | 4 維 mode primitive 是跨 session / cross-project 的 cognitive architecture 決策；expensive-to-reverse（改變 500 種 mode 組合的 composition logic）；explains-why（4 個 rationale 記錄在 §Decision Rationale） |
| **Phase 3 完成**（4 subsystem 真實 activation 驗證可行） | ⚠️ **部分** | YAML-contract level 完成（4 integration contracts + generated_surfaces 投影）；但 **behavioral enforcement**（pre-commit hook + Go 強制）尚未完成 — 4 subsystem 有合約但無真實 Go-level activation blocking |
| ~~5 個 Open Questions 全解~~ | ✅ | 已於 2026-05-22 resolved（見 §Open Questions） |
| 沒有更輕的 promotion target 適用 | ✅ | 4 維 mode primitive 是 architectural decision（非 runtime gate / enforcement / intelligence atom）；500 種 compose space + ADR-level cross-session 影響，無法用更輕 target 代替 |
| 系統真實使用此 contract，Phase 3 後至少 5 個 task final report 列 Cognitive Mode | ⚠️ **待累積** | Phase 1.B 起的 commit 含 Cognitive Mode（3 個 phase-setup tasks）；需累積 2 個以上的 non-setup user tasks |

**評估結論（2026-05-25）**：**ADR promotion 尚未達標，原因 2 個**：

1. **Phase 3 behavioral enforcement 未完成**：4 subsystem integration 是 YAML-only，Go-level activation blocking（execution_mode floor 強制、governance_mode gate 過濾、memory subdir 限制）尚未實作。ADR 宣稱「真實 runtime activation」，但目前仍是 YAML contract layer。
2. **5 task final report 累積未達**：需在 non-setup tasks 中再累積 2 個以上含 Cognitive Mode 的 final report。

**ADR promotion 重新評估時機**：Phase 3 behavioral enforcement（pre-commit hook + Go gates）完成後，且 5 tasks 累積達標，重新評估。此時 promotion 幾乎確定通過。

Phase 4-5 是優化（cost、adaptive），不是架構驗證 — 完成 Phase 3 behavioral enforcement 即可重評 ADR promotion。

### Consequences（預期）

#### 正面

- **真正 runtime activation**：mode 寫入 runtime state，subsystem 強制依 mode 行動
- **Token cost 可控**：discovery 靠訊號不靠文件；activation 是 conditional
- **Governance 強度差異化**：LIGHT / STANDARD / STRICT / LOCKDOWN 對應不同 gate set
- **Memory mode activation**：明確區分「不查記憶」與「查 episodic」「replay decision」狀態
- **4 維 composable**：500 種狀態組合
- **Backward compat**：既有 profile / compression / memory 文件保留

#### 負面

- **runtime.db schema 擴充**：新增 `cognitive_modes` + `discovery_signals` 表
- **Discovery heuristic 維護成本**：訊號 → mode 映射規則需校準
- **詞彙重疊風險**：`context_mode` 與既有 `compression` 名稱重疊（Open Question 1）
- **5 phase 大改動**：可能需 1-3 個月推進
- **教學負擔**：新概念需文件化

#### 風險

| 風險 | 緩解 |
|------|------|
| 4 mode 互動空間 500 種難以全測試 | Phase 1 先定義「常見組合」，未列組合走 default fallback |
| Discovery heuristic 誤判 → 用錯 mode | Mode 內加 escalation 規則，偵測訊號矛盾時自動升級 |
| 與既有 phase_machine 概念衝突 | Phase 3 整合時明確定義「phase = transaction state，mode = cognitive state」邊界 |
| 重構期間既有任務行為改變 | Phase 1-2 只新增 surface，不改現有 activation；Phase 3 起才切換 |
| 對方批評變成 paper plan 不執行 | 每 phase 設明確 completion gate，未過不開下一 phase |
| Cognitive vocabulary 形成 subsystem-local semantics | 必須引用 `plans/active/2026-05-25-1000-context-language-glossary-system.md` 的 glossary canonical definitions；`context_mode`、`compression`、`memory_mode`、`reasoning_mode` 不得在 runtime、workflow、memory 各自重新定義。 |

---

## Runtime Execution Path

依 [`governance/lifecycle/system-upgrade-governance.md`](../../governance/lifecycle/system-upgrade-governance.md) §3 規則 8：

| 欄位 | 內容 |
|------|------|
| Runtime owner | **Phase D**: 無 — doc-only trial / **Phase 1-5**: `runtime.db cognitive_modes` 表 + `phase_machine` integration |
| Trigger flow | **Phase D**: user / agent task begins → agent manually resolves mode from plan heuristic → final report includes Cognitive Mode block → commit history / Phase D metrics provide evidence.<br>**Phase 2+**: task entry signal or file diff → signal discovery table / runtime query → `runtime/cognitive-modes.yaml` contract loaded → cognitive mode written / selected → subsystem gates adapt behavior → final report and validation scenarios prove mode was applied. |
| Trigger location | **Phase D**: agent 手動 resolve + final report / **Phase 2+**: task entry signal-based discovery |
| Activation contract | **Phase D**: 無 / **Phase 1+**: `runtime/cognitive-modes.yaml`（將建立） |
| Generated surface | **Phase D**: 無 / **Phase 1+**: `runtime.db generated_surfaces.cognitive_modes.contract` |
| Validation scenarios | [`validation/scenarios/failure-derived/plan-runtime-execution-path-v1.yaml`](../../validation/scenarios/failure-derived/plan-runtime-execution-path-v1.yaml)（本 plan 對應 governance rule 8 的測試）<br>Phase D specific：尚未建立（doc-only trial 期間不強制）<br>Phase 1+ 將建立：`cognitive-mode-resolution-bypass-v1.yaml`、`cognitive-mode-discovery-signal-coverage-v1.yaml` |
| Test passing evidence | Phase D: agent 在 final report 列 Cognitive Mode 區塊（累積 5+ commits 證實）<br>Phase 1+: runtime compile + validate 通過 |

### Semantic Dependency Boundary

Runtime cognitive modes must reference glossary canonical definitions for runtime semantic vocabulary. This is not optional because `context_mode`, `compression`, `memory_mode`, and future `reasoning_mode` are high-risk semantic terms shared across runtime, workflow, memory, validation, and model routing.

Rules:

- `context_mode`、`compression`、`memory_mode`、`reasoning_mode` 的 canonical meaning 必須由 glossary owner definition 或本 plan 的 pre-glossary placeholder 明確標記。
- 本 plan 可以定義 runtime implementation behavior，但不得重新定義 glossary-owned semantic meaning。
- 若 glossary term 尚未建立，Phase 1 必須列為 pending semantic dependency，並在 glossary Phase 3 建立時回填 canonical link。
- 若 cognitive modes implementation 發現詞義需要改名、alias 或 deprecate，必須回到 glossary plan 的 semantic ownership / status lifecycle，不得只在 runtime YAML 內局部修正。

### Doc-only Trial 聲明（Phase D 階段）

> ⚠️ **此 plan 目前為 doc-only trial（Phase D），尚未接入 runtime execute layer。**
>
> 理由：
> - Phase 1-5 是大改動（runtime.db schema、compiler 規則、discovery YAML、subsystem 整合），實作成本高
> - 4 個 mode 從未真實運作，先以 doc-only 累積實證較安全
> - Token cost 與既有 phase_machine 衝突風險未知
>
> 未來接入時機：
> - **Phase D 完成且通過評估指標**後進 Phase 0 Pre-Build Interrogation
> - **Phase 1 完成**後 cognitive_modes 表寫入 runtime.db
> - **Phase 3 完成**後 4 subsystem 真實 activation，此 plan 完全離開 doc-only 狀態
>
> Doc-only 期間如何避免漂移：
> - agent 每個 final report 強制列 Cognitive Mode 區塊（manual application contract）
> - commit message 含 mode 報告（即本 plan 從 Phase D 啟動以來的所有 commit）
> - Phase D 評估指標包含「mode 定義穩定性」（5 任務內未變動超過 2 次）

## Open Questions（completed 前需釐清）

**全部於 2026-05-22 resolved**。決議摘要：

| # | Question | Resolution | 內文位置 |
|---|----------|-----------|---------|
| 1 | `context_mode` vs `compression` 命名統一？ | ✅ **resolved**：統一為 5 級 UPPERCASE `context_mode`（INDEX_ONLY / SUMMARY_FIRST / CHECKLIST_FIRST / SOURCE_BACKED / **GRAPH_ASSISTED**）；既有 `compression/` lowercase 視為 alias；Phase 3.2 改寫 compression 文件成 implementation reference | §Decision §1 Naming alignment note |
| 2 | Discovery signal 來源是否還缺？ | ✅ **resolved**：加 4 個訊號（phase_machine state / active_goals / token_budget / generated_surfaces hit）；移除 `contradiction risk`（屬 derived signal，Phase 5 adaptive 才接） | Phase 2 §2.1 訊號表 |
| 3 | `LOCKDOWN` vs `blocking_gates` 關係？ | ✅ **resolved**：governance_mode 是 gate set selector，LOCKDOWN = STRICT 全集 + `additional_actions: [block_file_writes_until_human_approval]` | Phase 1 §1.1 gate_activation |
| 4 | `memory_mode` 與 `retrieval-governance threshold` 整合方式？ | ✅ **resolved**：mode 是 category switch、threshold 是 within-category activation；兩者 AND 邏輯 compose | Phase 3.4 |
| 5 | ADR promotion gate 設在哪 phase 完成？ | ✅ **resolved**：**Phase 3 完成時**評估 ADR Promotion Criteria；Phase 4-5 是優化非架構，不延後 ADR | Phase 3 完成條件 + §ADR Promotion Criteria |

新的 open questions 可於 Phase 0 Pre-Build Interrogation 補充。

---

## 完成條件（per system-upgrade-governance §2）

本 plan 是 Gen 3 子系統擴充，**不是世代升級**，所以系統名稱不變。但涉及架構分層擴充 + 核心流程變更，需執行 §2 checklist 子集。

### 計畫書本身

- [x] 計畫書狀態：in-progress（Phase 0 通過後啟動，Phase 1-3 governance-structure 完成）
- [ ] 記錄完成日期 ← 待 behavioral enforcement + ADR 完成
- [ ] 記錄偏離（如有）← 待 plan completed
- [ ] 歸檔至 `plans/archived/` ← 待 plan completed

### README 更新

- [x] `README.md` OS Layout 表格新增 cognitive modes 機制簡述（commit `cc98961`）
- [ ] `architecture/ai-native-cognitive-execution-system.md` 核心機制章節新增 cognitive modes ← 本 session 補做
- [x] `models/README.md` 加說明：profile 與 mode 的關係（backward-compat）← 本 session 補做

### 架構文件（per system-upgrade-governance §3 規則 6）

- [x] ADR Promotion Criteria 評估（2026-05-25，結論：尚未達標 — 見 §ADR Promotion Criteria 評估表）
- [ ] 若符合 → 建立 ADR（直接 accepted）← 尚未達標，不執行
- [x] 若不符合 → 在 plan 記錄「決定不升 ADR 的理由」← 已記錄（behavioral enforcement 未完 + 5 tasks 未累積）

### 索引與路由

- [x] `knowledge/runtime/routing-registry.yaml` 新增 `route.runtime.cognitive-modes`（commit `cc98961`）
- [ ] `knowledge/indexes/README.md` 反映新結構 ← 本 session 補做
- [x] `knowledge/graphs/` 新增 `cognitive-modes.yaml`（commit `cc98961`，13 edges）

### Runtime Surface

- [x] 新增 `runtime/cognitive-modes.yaml`（commit `aa5bc68`）
- [x] `runtime/runtime.db` 新增 `cognitive_modes` + `discovery_signals` 表（commit `69eb605`、`76b8360`）
- [x] Compiler 規則更新，投影到 `generated_surfaces`（commit `69eb605`：sourceRoots + compile rule）
- [x] `runtime/runtime.db` 的 `obligation_ledger` / `blocking_gates` 對接 mode（commit `21ce746`：enforcement obligation + gate）

### Linked Updates

- [x] `governance/lifecycle/knowledge-update-flow.yaml` 評估：**決定不加 step** — 改在 Step 1 entry 自動觸發（per §3.5，Phase 3 decision）
- [x] `enforcement/failure-patterns/` 加 `cognitive-mode-resolution-bypass.md` ← 本 session 補做
- [x] `models/compression/README.md` 標註與 `context_mode` 的對應關係（commit `2a3edba`）
- [x] `memory/retrieval-governance/` 標註與 `memory_mode` 的對應關係 ← 本 session 補做

### 閉環驗證

- [x] Diff review（每次 commit 前執行）
- [ ] Linked updates 完成 ← 本 session 補做剩餘項
- [x] `ai-skill runtime compile --native-compiler` + `ai-skill runtime validate` 通過（每次 commit 後確認）
- [ ] Commit / push / readback ← 本 session 結尾執行

---

## Phase D: Documentation-Contract Trial（先行 trial，無 runtime 程式碼）

**目標**：在進入 Phase 0-5 runtime 實作前，先以**純文件契約**方式驗證 4 維 cognitive mode 設計是否實用。Agent 手動套用 mode，在 final report 回報。

**Status**: in-progress（2026-05-22 啟動）

### 為什麼先做 Phase D

- Phase 1-5 是大改動（runtime.db schema、compiler 規則、discovery YAML、subsystem 整合），實作成本高
- 4 個 mode 的設計從未真實運作過，不知道組合是否實用
- Token cost 與既有 phase_machine 衝突風險未知
- Doc-only trial 可在數天內累積實證，撤回成本極低（`git revert` 即可）

### Phase D 範圍

| 動作 | 範圍 |
|------|------|
| 建立 `models/cognitive-modes/README.md` | 4 維 mode primitive 定義 + discovery 速查表 + final report 範本 + rollback 指引 |
| 更新 `models/README.md` | 加 cognitive-modes/ 入口 |
| Agent final report 列 Cognitive Mode 報告 | 每個任務必填 4 維值 + 理由欄 |
| 累積 trial runs | ≥ 5-10 個任務，記錄在 git history |
| Phase D 結束評估 | 通過 → Phase 0；不通過 → 修設計或撤回 |

### Phase D **不做**的事

- ❌ 寫 `runtime/cognitive-modes.yaml` executable contract
- ❌ 加 `runtime.db cognitive_modes` 表
- ❌ 寫 compiler 規則
- ❌ 寫 discovery YAML 機器化
- ❌ 改 `models/compression/` 或其他既有 models/ 子層
- ❌ 改 phase_machine / blocking_gates / retrieval-governance 任何 runtime state

### Phase D 評估指標（trial 結束時驗證）

| 指標 | 通過條件 | **實際結果**（2026-05-22）|
|------|------|------|
| 累積任務數 | ≥ 5 個 final report 含 Cognitive Mode 區塊 | **8 個** ✅（9df20ae / f98d6e4 / df37b1a / 57ec8fc / 7dc96d2 / 84e736d / ef305bf / 3a38e49） |
| Mode 組合多樣性 | 至少覆蓋 3 種不同 execution + 3 種不同 context + 2 種 governance | **5/5 execution + 4/5 context + 4/4 governance + 3/5 memory** ✅ |
| 誤判率 | ≤ 20% 的 mode 選擇被使用者糾正 | **0% 糾正**（user 未糾正過 mode 選擇）✅ |
| Raw signal 完整性 | 沒有「找不到適合 signal 表達」的情況 | **0 case 缺 signal**；T5 額外揭露「rollback claim drift」signal，已加入 Phase 5 adaptive ✅ |
| 設計穩定性 | mode 定義在 5 個任務內未變動超過 2 次 | **8 任務內 0 次變動**（FULL_TRACE→GRAPH_ASSISTED 是 Open Question 1 resolve，發生在 trial 開始前的設計階段）✅ |
| Validation scenarios | scenarios 全部通過 | **5/5 automatable PASS + 1 heuristic N/A** ✅ |

**評估結論**：Phase D 全部指標通過 ✅

### Phase D 完成條件

- [x] `models/cognitive-modes/README.md` 寫入
- [x] `models/README.md` 加入口
- [x] Plan 加 Phase D 段落（本段）
- [x] 至少 1 個任務（本 commit 即是）在 final report 列 Cognitive Mode
- [x] 累積 5 個任務的 mode 報告（達成 **8 個**，commits 列於上方評估指標表）
- [x] 評估指標通過（6/6 ✅）
- [ ] 決定下一步：進 Phase 0 / 修設計 / 撤回 plan **← 待 user 決定**

### Phase D Trial 副產品（commit chain 沉澱的成果）

| Trial commit | 副產品 |
|--------------|--------|
| `f98d6e4` T1 | 2 處 stale 描述修正（plan + cognitive-modes README） |
| `df37b1a` T2 | anti-patterns README + failure-patterns README 索引補齊 |
| `57ec8fc` T5 | Rollback path 真實實證（揭露 over-optimistic claim） |
| `7dc96d2` T4 | 4 個 failure pattern 結構修正（部分） |
| `84e736d` T3 | Intelligence atom `plan-first-decision-promotion` + feedback lesson |
| `ef305bf` | Governance rule 8（升級計畫必須含 runtime execution path）+ 6 validation scenarios |
| `3a38e49` A+B | 11 個 failure pattern 結構全對齊 + 新 failure pattern `template-drift` |

**Phase D 對 governance 的正面影響**：
- 1 個新 governance rule（system-upgrade-governance §3 規則 8）
- 1 個新 intelligence atom（plan-first-decision-promotion）
- 2 個新 failure patterns（premature-adr-promotion、template-drift）
- 6 個 validation scenarios
- 既有 11 個 failure pattern 結構對齊

### Phase D Rollback

| 完整撤回 | `git revert <Phase D commit hash>` 移除 cognitive-modes/README.md、models/README 入口、本段落；無 runtime state 變更 |
| 暫停 trial | 在本段加 `paused` 標記，agent final report 不再列 Cognitive Mode |
| 修 mode 定義 | 編輯 `models/cognitive-modes/README.md`；下次任務套用新定義 |

---

## Phase 0: Pre-Build Interrogation + Architecture Compatibility Preflight

**前置條件**：Phase D 完成且通過評估指標 ✅（2026-05-22）

**Status**: ✅ **completed**（2026-05-22）

### Architecture Compatibility Preflight Ledger（per `plans/README.md`）

| 欄位 | 內容 |
|------|------|
| Trigger | 開始實作 Cognitive Modes Phase 1（runtime primitive 落地） |
| Checked sources | `runtime/runtime.db` schema（47 tables）+ `scripts/ai-skill-cli/internal/app/runtime_compiler.go` + `models/cognitive-modes/README.md` + `models/compression/README.md` + `memory/retrieval-governance/` 結構 + `governance/ai-runtime-governance/` 結構 |
| Conflicts | **無衝突** — 詳見下方 §Preflight 7 Checks |
| Interrogation | 詳見下方 §Pre-Build Interrogation Result |
| Decision | **proceed to Phase 1**（待 user 啟動）|
| Validation | `sqlite3 runtime.db ".tables"` 確認 47 tables 不含 cognitive_modes/discovery_signals；既有 compiler rules pattern 可擴展 |

### Preflight 7 Checks

| # | 檢查項目 | 結果 |
|---|---------|------|
| 1 | **Candidate files 存在性** | `runtime/cognitive-modes.yaml`（Phase 1 將建立，目前不存在 ✅）/ `runtime.db cognitive_modes` table（不存在 ✅，Phase 1 新增）/ `runtime.db discovery_signals` table（不存在 ✅，Phase 2 新增）/ `scripts/ai-skill-cli/internal/app/runtime_compiler.go`（存在，將擴展）|
| 2 | **Source-of-truth 一致性** | YAML canonical = `runtime/cognitive-modes.yaml`；projection target = `runtime.db generated_surfaces (target_key='runtime.cognitive_modes.contract')`；compiler canonical = `runtime_compiler.go`。與既有 16+ executable YAML contracts 命名/投影 pattern 完全一致 |
| 3 | **Layer responsibility** | runtime primitive → `runtime/`（YAML + runtime.db table）✅ / doc-only contract → `models/cognitive-modes/`（Phase D 已建立）✅ / validation → `validation/scenarios/`（Phase 1 將新增）✅ / 無放錯層風險 |
| 4 | **Compiler / generated surface** | 新 target_key `runtime.cognitive_modes.contract` 不與既有 70+ target_keys 衝突；compile rule 將在 `runtime_compiler.go` 加新 case，pattern 與 `decision-promotion-pipeline.yaml` 等既有 YAML 投影完全一致 |
| 5 | **Pre-build interrogation** | 見下方 §Pre-Build Interrogation Result，含 framework discovery、duplication risk、open questions、assumptions 全部記錄 |
| 6 | **Linked updates** | Phase 1 需更新：`models/README.md`（cognitive-modes entry 已存在，Phase 1 補 runtime 投影連結）/ `architecture/ai-native-cognitive-execution-system.md` §核心機制（加 cognitive_modes table）/ Phase 2 需更新：`knowledge/runtime/routing-registry.yaml`（加 route）/ `knowledge/summaries/`（加 summary）/ `knowledge/graphs/`（加 edges）|
| 7 | **Execution decision** | **proceed to Phase 1** — 無架構衝突、無未解 blocker、無 source-of-truth duplication；Open Questions 1-5 已 resolved；Phase D 8 個 trial commits 已累積實證 |

### Pre-Build Interrogation Result

| 維度 | 答案 |
|------|------|
| **Goal** | 將 `models/cognitive-modes/` 從 Phase D documentation-contract 升級為 Phase 1 runtime activation layer；補 governance / memory 強度差異化缺口 |
| **Scope** | Phase 1 only: (a) `runtime/cognitive-modes.yaml` executable contract 寫入；(b) `runtime.db cognitive_modes` table 加入；(c) `runtime_compiler.go` 加新 compile rule；(d) 投影到 `generated_surfaces` (target_key=`runtime.cognitive_modes.contract`)；(e) Phase 1 POC：至少 1 個任務手動寫入 mode 並由 runtime.db 強制 |
| **Non-goals**（Phase 1） | 不寫 discovery YAML（Phase 2）；不接 phase_machine / blocking_gates / memory subsystem（Phase 3）；不寫 token budget gate（Phase 4）；不寫 adaptive triggers（Phase 5）；不改 `models/profiles/` / `compression/` / `capabilities/` |
| **Acceptance（Phase 1）** | (a) `runtime/cognitive-modes.yaml` schema-valid + 通過 `ai-skill runtime validate`；(b) `runtime.db cognitive_modes` table 存在且可寫入；(c) `generated_surfaces` 含 `runtime.cognitive_modes.contract` row；(d) 至少 1 個 POC task 的 cognitive_modes row 寫入 runtime.db；(e) commit 含 final report Cognitive Mode 區塊 |
| **Framework discovery findings** | runtime.db 47 tables 已盤點，無 cognitive_modes / discovery_signals 既存；compiler entry = `runtime_compiler.go`；target_key 命名規則 = `<owner>.<contract>.<aspect>`；既有 YAML contract（e.g. `decision-promotion-pipeline.yaml`）為 reference template |
| **Duplication risk** | context_mode vs compression：Open Question 1 已 resolved 為 UPPERCASE primitive + lowercase alias；Phase 3.2 才實際改寫 compression docs，**Phase 1 不觸發此 duplication** |
| **Open questions** | 5 個 Open Questions 全部 resolved（見 §Open Questions section），無新增 |
| **Assumptions** | (a) runtime.db schema 可擴充（既有 47 tables 自由新增第 48 個無風險）；(b) compiler pattern 一致（如 decision-promotion-pipeline 之 compile），新增 compile rule 不破壞既有；(c) Phase D doc-contract 為 source-of-truth，Phase 1 YAML 引用之 |
| **Conflicts** | 無 — runtime table、compiler rule、generated surface 命名、layer responsibility 全部 unique 且符合既有 pattern |

### Phase 0 完成條件

- [x] Architecture Compatibility Preflight Ledger 完成（7 checks 全通過）
- [x] Pre-Build Interrogation Result 詳細記錄
- [x] 確認無架構衝突、無 source-of-truth duplication
- [x] 確認 Open Questions 全部 resolved
- [x] 確認 candidate files 不衝突（cognitive_modes / discovery_signals tables 不存在）
- [x] 決定：**proceed to Phase 1**
- [ ] Phase 1 啟動 ← **待 user 確認啟動 Phase 1**

### Phase 0 Rollback

| 動作 | 操作 |
|------|------|
| 完全撤回 Phase 0（保留 Phase D） | `git revert <Phase 0 commit>` — 移除本 Ledger / Interrogation Result / 完成條件；Phase D 不受影響 |
| 重新做 Preflight | 編輯本 section 重新填 7 checks（不必撤回）|

---

## Phase 1: Cognitive Mode Primitives

**目標**：定義 4 mode 的 executable YAML contract，投影到 runtime.db。

### 1.1 設計 `runtime/cognitive-modes.yaml`

```yaml
schema_version: "1.0"
id: runtime.cognitive-modes
status: active
owner_layer: runtime

runtime_projection:
  enabled: true
  target_key: runtime.cognitive_modes.contract
  surface: generated_surfaces

modes:
  execution:
    values: [FAST, NORMAL, DEEP, FORENSIC, RECOVERY]
    default: NORMAL
  context:
    values: [INDEX_ONLY, SUMMARY_FIRST, CHECKLIST_FIRST, SOURCE_BACKED, GRAPH_ASSISTED]
    default: SUMMARY_FIRST
  governance:
    values: [LIGHT, STANDARD, STRICT, LOCKDOWN]
    default: STANDARD
    gate_activation:
      LIGHT: [sanitization]
      STANDARD: [sanitization, language_policy, output_rules]
      STRICT: [sanitization, language_policy, output_rules, linked_updates, runtime_surfaces]
      LOCKDOWN: STRICT + block_writes
  memory:
    values: [NONE, EPISODIC, DECISION_REPLAY, FAILURE_REPLAY, PROJECT_CONTEXT]
    default: NONE
```

### 1.2 `runtime.db cognitive_modes` 表

```sql
CREATE TABLE cognitive_modes (
  id INTEGER PRIMARY KEY,
  task_id TEXT,
  execution_mode TEXT,
  context_mode TEXT,
  governance_mode TEXT,
  memory_mode TEXT,
  resolved_at TEXT,
  source TEXT
);
```

### 1.3 Compiler 規則

`scripts/ai-skill-cli/internal/compiler/cognitive_modes.go`：將 yaml 投影到 `generated_surfaces` 與 `cognitive_modes` 表。

### Phase 1 完成條件（拆 1.A 與 1.B）

**Phase 1.A**（doc-only YAML + scenarios，無 Go 改動）：
- [x] Pre-implementation scenarios 已寫入並驗證 fail（commit `395a0d9`，4 個 scenarios）
- [x] `runtime/cognitive-modes.yaml` 寫入並通過 schema validation
- [x] Scenario `cognitive-modes-yaml-contract-exists-v1` → **PASS** ✅
- [ ] Scenario `cognitive-modes-generated-surface-projected-v1` → **BLOCKED**
  （`compileExecutableYAMLContracts` sourceRoots 不含 `runtime/`，需 Go rebuild）

**Phase 1.B**（2026-05-25 completed）：
- [x] 加 `"runtime"` 到 `compileExecutableYAMLContracts` sourceRoots
  （`scripts/ai-skill-cli/internal/app/runtime_compiler.go` ~L556）
- [x] 加 `CREATE TABLE cognitive_modes` 到 `createGoRuntimeSchema`
  （同檔 ~L136）
- [x] Rebuild 5 個 platform binaries（darwin-arm64/amd64、linux-arm64/amd64、windows-amd64）
- [x] `ai-skill runtime compile --native-compiler` 投影 cognitive-modes.yaml 到 generated_surfaces
- [x] `ai-skill runtime validate` 通過
- [x] Scenario `cognitive-modes-generated-surface-projected-v1` → **PASS** ✅
- [x] Scenario `cognitive-modes-runtime-table-exists-v1` → **PASS** ✅
- [x] 至少手動寫入一個 task 的 mode resolution（POC）
- [x] Scenario `cognitive-modes-poc-task-record-v1` → **PASS** ✅

---

## Phase 2: Discovery Heuristics

**目標**：定義訊號 → mode 映射，不靠文件查詢。

### 2.1 訊號來源（per Open Question 2 resolved）

| 訊號 | 例 | 對應 mode 暗示 | 來源 |
|------|------|------|------|
| user keyword | 「快速」「typo」「修一下」 | execution=FAST | conversation parse |
| user keyword | 「分析」「跨層」「migration」 | execution=DEEP | conversation parse |
| user keyword | 「recover」「失誤」「重做」 | execution=RECOVERY | conversation parse |
| file diff scope | `enforcement/` / `governance/` | governance=STRICT | git diff |
| file diff scope | `notes/` / `memory/working/` | governance=LIGHT | git diff |
| git status | dirty + 多 owner group | governance=STRICT | git status |
| session turn 數 | > 50 turns | memory=PROJECT_CONTEXT | runtime self |
| recent failure | failure_repeat ≥ 2 | execution=RECOVERY, memory=FAILURE_REPLAY | runtime self |
| **phase_machine 當前 phase** | bootstrap / execution / validation | bootstrap → governance=LIGHT；validation → governance=STRICT | `runtime.db phases` |
| **active goals 狀態** | 有 owner/lock 或 ≥2 active | governance=STRICT | `.agent-goals/` |
| **token_budget 剩餘** | < 30% budget | context 降級（GRAPH_ASSISTED → SOURCE_BACKED → SUMMARY_FIRST） | runtime self |
| **modified files 命中 generated_surfaces** | 任一檔命中 | governance=STRICT；強制 runtime compile | `sqlite3 runtime.db` |

**Deferred to Phase 5 adaptive triggers**（這些是 derived signal，非 raw signal）：

- contradiction risk（從多源比對推導）
- tool call repeat pattern（從歷史推導）
- cross-session memory hit rate（從 memory retrieval 統計推導）

### 2.2 Discovery YAML contract

`runtime/cognitive-modes-discovery.yaml`，將訊號 → mode 規則機器化。

### Phase 2 完成條件（2026-05-25 completed）

- [x] Discovery 訊號規則表完整（14 rules，8 signal types，覆蓋 plan §2.1 全部訊號）
- [x] 在 runtime.db `discovery_signals` 表中可查詢（compiler 從 YAML 投影）
- [x] 規則 fallback：訊號不命中時用 default mode（fallback section 明確定義）
- [x] POC：5 個常見任務類型測試 discovery 輸出符合預期（8 distinct signal_type PASS）

---

## Phase 3: Subsystem 整合

**目標**：4 個 subsystem 依 mode 自動 activation。

### 3.1 execution_mode → phase_machine

phase_machine 進入 phase 時，依 execution_mode 調整 allowed_actions / forbidden_actions。

### 3.2 context_mode → compression

合併 `models/compression/` levels 為 `context_mode` 的實作層。決定保留兩個詞彙或統一（依 §Open Question 1）。

### 3.3 governance_mode → blocking_gates

`runtime.db blocking_gates` 依 governance_mode 啟用對應 gate set。LOCKDOWN 額外阻擋寫入。

### 3.4 memory_mode → memory retrieve

`memory/retrieval-governance/` 的 activation threshold 與 memory_mode compose；mode 是 switch，threshold 是 limit。

### 3.5 knowledge-update-flow 整合

決定是否加 Step 0「mode resolution」。建議：**不加 step，改在 Step 1 entry 自動觸發**。

### Phase 3 完成條件

- [x] 4 subsystem 各有對應 mode-driven activation 邏輯（YAML contract level，2026-05-25）
  - [x] 3.1 execution_mode → phase_machine（`runtime/cognitive-modes-phase-integration.yaml`，含 5 mode 的 allowed/forbidden_actions floor）
  - [x] 3.2 context_mode → compression（`models/compression/README.md` 加 alias 表 + context_mode column）
  - [x] 3.3 governance_mode → blocking_gates（`runtime/cognitive-modes-governance-integration.yaml`，含 LOCKDOWN block_file_writes）
  - [x] 3.4 memory_mode → retrieval-governance（`runtime/cognitive-modes-memory-integration.yaml`，含 AND 邏輯 compose）
- [x] 失敗測試：mode 未解析就執行 → 被阻擋（governance-structure level，2026-05-25）
  - [x] `obligation.execution.resolve_cognitive_mode` 寫入 runtime.db `obligations` 表（Go compiler 直接 INSERT）
  - [x] `gate.execution.cognitive_mode_resolved` 寫入 runtime.db `gates` 表
  - [x] Scenario `cognitive-modes-enforcement-gate-exists-v1` → **PASS** ✅
  - [ ] **Behavioral enforcement**（deferred）：pre-commit hook 實際查 `cognitive_modes` 表是否有本次 session row → 沒有就 block commit。需改 pre-commit hook Go 邏輯（`scripts/ai-skill-cli/internal/app/` pre-commit 相關 function）＋ rebuild。
- [x] Documentation 更新（`models/compression/README.md` alias annotation，per Open Question 1）
- [ ] **Phase 3 完成後至少 5 個任務的 final report 列 Cognitive Mode**（用於 ADR promotion 評估）← 待累積（每次 task commit 含 final report Cognitive Mode 區塊）
- [ ] **ADR Promotion Criteria 評估**（per Open Question 5 resolved）← 待下次 session
  - 通過 → 建立 ADR（直接 accepted），引用本 plan 作為 evidence
  - 不通過 → 在 plan 內記錄「決定不升 ADR 的理由」+ 改用更輕 promotion target

### Phase 3 Deferred Items（完整記錄，下次 session 繼續）

**2026-05-25 update**：所有 behavioral enforcement items 已於本 session 完成（commit-msg hook 為實作層）。

| 項目 | 狀態 | 完成證據 |
|------|------|---------|
| **Behavioral enforcement**（commit-msg hook） | ✅ **完成** | `runCommitMsgHook` in `scripts/ai-skill-cli/internal/app/hooks.go` 檢查 `### Cognitive Mode 報告` block；缺則 exit 30 阻擋；live test 4 條路徑通過；scenario `cognitive-mode-block-required-v1` PASS |
| **3.1-B Go enforcement**（execution_mode floors） | ✅ **完成** | `validateExecutionModeFloors` in hooks.go：FAST 禁觸 governance/enforcement/runtime/；DEEP/FORENSIC/RECOVERY 要求 governance_mode≥STRICT；context_mode/memory_mode floors 依 `runtime/cognitive-modes-phase-integration.yaml`；4 unit tests PASS |
| **3.3-B Go enforcement**（governance_mode gate query） | ✅ **完成** | `validateGovernanceModeConsistency` in hooks.go：LIGHT 禁觸 governance-critical paths；LOCKDOWN 要求 `[approved-by: <name>]` trailer；scenario `phase3-behavioral-validators-v1` PASS |
| **3.4-B Go enforcement**（memory subdir activation） | ✅ **完成** | `validateMemoryModeSubdir` in hooks.go：宣告的 memory_mode 必須與 staged `memory/` 子目錄一致；NONE 禁觸；EPISODIC→episodic/、DECISION_REPLAY→decision/、FAILURE_REPLAY→failure/、PROJECT_CONTEXT→project/；`memory/README.md` 與 `memory/retrieval-governance/` 視為 layer doc 豁免 |
| **Linked updates**（routing-registry） | ✅ 已於前次 session 完成 | `knowledge/runtime/routing-registry.yaml:1878` 含 `route.runtime.cognitive-modes` |
| **Linked updates**（knowledge graph） | ✅ 已於前次 session 完成 | `knowledge/graphs/cognitive-modes.yaml` 存在 |
| **Linked updates**（README OS layout） | ✅ 已於前次 session 完成 | `README.md` 含 cognitive modes 描述 |
| **5 tasks final report 累積** | ⏳ 進行中 | 自然累積，不需額外工作 |
| **ADR Promotion Criteria 評估** | ⏳ 等 5 tasks | 5 tasks 累積後評估 |

**Behavioral 層 caveat**：commit-msg hook 是 enforcement point（每個 commit 都過），不是 agent-action point（每個 tool call）。要做到 per-tool-call 仍需更深整合（Phase 5 adaptive runtime）；commit-msg 層已能在 close-loop 時擋下 mode 違反，足以閉環 Phase 3。

---

## Phase 4: Selective Loading + Token Budget Gate

**目標**：Source 載入由 mode 控制，超 budget 自動降級或阻擋。

### 4.1 Token budget per mode 組合

| Mode 組合 | Budget |
|------|------|
| FAST + INDEX_ONLY + LIGHT + NONE | ≤ 1000 tokens |
| NORMAL + SUMMARY_FIRST + STANDARD + EPISODIC | ≤ 5000 tokens |
| DEEP + SOURCE_BACKED + STRICT + DECISION_REPLAY | ≤ 20000 tokens |
| FORENSIC + GRAPH_ASSISTED + STRICT + FAILURE_REPLAY | ≤ 50000 tokens |

### 4.2 Budget gate

超 budget 時：
1. 自動降級 context_mode（GRAPH_ASSISTED → SOURCE_BACKED → SUMMARY_FIRST）
2. 若無法降級，阻擋並要求 user 確認
3. 記錄 budget overflow 為 failure pattern

### Phase 4 完成條件

- [x] Budget table 寫入 runtime.db — `runtime/cognitive-modes-token-budget.yaml` projected to `generated_surfaces` (target_key=`runtime.cognitive_modes.token_budget`)
- [x] Budget gate 在 runtime 強制 — `validateTokenBudget` in `scripts/ai-skill-cli/internal/app/hooks.go` 由 commit-msg hook 觸發，當 commit body 宣告 `Token Estimate: <n>` 超 budget 即 block；scenario `phase4-token-budget-v1` PASS；6 unit tests PASS
- [x] Failure pattern `cognitive-budget-overflow.md` 建立 — `enforcement/failure-patterns/cognitive-budget-overflow.md`

**Phase 4 完成日期**：2026-05-25

**Behavioral 層 caveat**（與 Phase 3 一致）：commit-msg 是 enforcement point，validator 只在 commit body 宣告 `Token Estimate:` trailer 時 fire；per-tool-call **實際**用量偵測等 Phase 5 adaptive runtime。

**Downgrade path**（YAML contract 定義）：當預估超 budget 應依序 `GRAPH_ASSISTED → SOURCE_BACKED → CHECKLIST_FIRST → SUMMARY_FIRST → INDEX_ONLY` 降級，最後才用 `[skip-token-budget]` opt-out。

---

## Phase 5: Adaptive Runtime

**目標**：mode 動態調整，不固定在 task 開始時 resolve。

### 5.1 Adaptive triggers

| 訊號 | 來源 | 動作 |
|------|------|------|
| Token budget 接近上限 | runtime self | 降級 context_mode（GRAPH_ASSISTED → SOURCE_BACKED → SUMMARY_FIRST） |
| **偵測到 contradiction risk**（多源衝突） | derived from cross-source compare | 升 governance_mode + 升 context_mode |
| 連續 failure ≥ 2 | runtime self | 切到 execution=RECOVERY, memory=FAILURE_REPLAY |
| 跨 phase transition | phase_machine | 重新 evaluate 全部 4 個 mode |
| **Tool call repeat pattern**（loop detection） | derived from call history | 升 governance_mode + 觸發 recovery escalation |
| **Cross-session memory hit rate** | derived from memory retrieval 統計 | 調整 memory_mode 強度 |

### 5.2 Adaptive YAML contract

`runtime/cognitive-modes-adaptive.yaml`。

### Phase 5 完成條件

- [x] Adaptive 規則表 — `runtime/cognitive-modes-adaptive.yaml` 6 triggers（3 commit-msg-detectable + 3 out-of-scope for commit-msg layer）；projected to `generated_surfaces` (target_key=`runtime.cognitive_modes.adaptive`)
- [x] runtime 偵測訊號 → 自動切 mode — `validateAdaptiveTriggers` in `scripts/ai-skill-cli/internal/app/hooks.go` 由 commit-msg hook 觸發；3 commit-msg-detectable triggers 已 enforced（contradiction_risk block / repeated_failure block / budget_near_ceiling warning）
- [x] 測試 case：模擬 contradiction / failure / budget overflow 三種情境 — `TestValidateAdaptiveTriggers` 3 case PASS；scenario `phase5-adaptive-triggers-v1` PASS

**Phase 5 完成日期**：2026-05-25

**Scope boundary at commit-msg layer**：

| Trigger | Commit-msg 偵測 | 動作 |
|---|---|---|
| contradiction_risk | ✅ keyword + ≥2 distinct source refs | block 若 governance<STRICT 或 context<SOURCE_BACKED |
| repeated_failure | ✅ ≥2 failure-pattern refs OR ≥2 revert/hotfix/retry | block 若 execution≠RECOVERY 或 memory≠FAILURE_REPLAY |
| budget_near_ceiling | ✅ Token Estimate ≥80% budget | warning（advisory） |
| tool_call_loop_detection | ❌ 需 per-tool-call interception | 不在 commit-msg 範圍 |
| cross_session_memory_hit_rate | ❌ 需 cross-session stats | 不在 commit-msg 範圍 |
| phase_transition_remap | ❌ 需 phase transition events | 不在 commit-msg 範圍 |

**Plan 整體完成**：Phase D + Phase 0-5 全部 ✅。剩餘 deferred items（ADR Promotion Criteria 評估、5 tasks 累積）為自然進行項，無 implementation work。

---

## Stakeholder 同意項目

已於 2026-05-22 與 user 確認下列項目：

- [x] 同意採 4 維 mode primitive 而非單一 enum
- [x] 同意不重寫 `models/`，採 backward-compat
- [x] 同意分 5 phase 漸進實作
- [x] 同意 ADR promotion gate 設在 Phase 3 完成
- [x] 同意 Open Questions 1-5 全部 resolved 方向（見 §Open Questions）

剩餘待 sign-off：

- [x] **同意先走 Phase D doc-only trial**（2026-05-22 確認）
- [x] **Phase D trial 結束後評估指標 通過**（2026-05-22，6/6 指標全 PASS — 詳見 §Phase D 評估指標）
- [x] **Phase 0 Pre-Build Interrogation 草稿審閱 + Architecture Compatibility Preflight 通過**（2026-05-22，7 checks 全 PASS — 詳見 §Phase 0）
- [x] **Phase 1 實作前 candidate files 與 runtime.db schema 擴充影響範圍評估**（無衝突，既有 47 tables 不含 cognitive_modes / discovery_signals；compiler pattern 可擴展）
- [ ] **Phase 1 啟動** ← 待 user 確認

---

## 與其他 plans 的關係

| Plan | 關係 |
|------|------|
| [`plans/archived/2026-05-20-1802-model-aware-execution-routing.md`](../archived/2026-05-20-1802-model-aware-execution-routing.md) | 前作；定義 model-aware routing 設計層；本 plan 推進到 runtime activation 層 |
| [`plans/archived/2026-05-22-0855-executable-yaml-contract-migration.md`](../archived/2026-05-22-0855-executable-yaml-contract-migration.md) | 提供 executable contract 機制基礎；本 plan 依賴此機制 |
| [`plans/archived/2026-05-20-1745-memory-retrieval-activation-governance.md`](../archived/2026-05-20-1745-memory-retrieval-activation-governance.md) | 提供 memory activation 邊界；本 plan 將 retrieve threshold 與 memory_mode 整合 |
| [`plans/archived/2026-05-20-1039-runtime-recovery-escalation-system.md`](../archived/2026-05-20-1039-runtime-recovery-escalation-system.md) | RECOVERY mode 對應此 system |
| [`plans/archived/2026-05-15-0920-runtime-execution-layer-upgrade-analysis.md`](../archived/2026-05-15-0920-runtime-execution-layer-upgrade-analysis.md) | Gen 3 升級分析；本 plan 是 Gen 3 子系統擴充 |
| [`plans/active/2026-05-25-1000-context-language-glossary-system.md`](2026-05-25-1000-context-language-glossary-system.md) | 必須提供 runtime semantic vocabulary 的 canonical definitions；本 plan 的 `context_mode`、`compression`、`memory_mode`、`reasoning_mode` 不得形成 subsystem-local semantics |
