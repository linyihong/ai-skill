# Runtime Cognitive Modes System

**Status**: `in-progress`（Phase D 進行中）
**世代**：Gen 3 子系統擴充
**建立日期**：2026-05-22
**最後更新**：2026-05-22（Phase D 啟動 — documentation-contract trial）

> ⚠️ 本 plan 為 draft 階段。原 `constitution/ADR-008-runtime-cognitive-modes.md`（proposed）已於 2026-05-22 撤回；依新 [`decision-promotion-pipeline`](../../governance/lifecycle/decision-promotion-pipeline.md) 規則，constitution/ 只放 accepted ADRs，提案階段在本 plan 內處理。
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

### ADR Promotion Criteria（completed 時驗證）

升級為 accepted ADR 的條件（per ADR-007 §ADR Boundary）。**Promotion gate 設在 Phase 3 完成時**（per Open Question 5 resolved）：

- [ ] foundational + cross-session + cross-project + expensive-to-reverse + explains-why 全中
- [ ] **Phase 3 完成**（4 subsystem 真實 activation 驗證可行）
- [ ] ~~5 個 Open Questions 全解~~ ✅ 已於 2026-05-22 resolved（見 §Open Questions）
- [ ] 沒有更輕的 promotion target 適用（runtime gate / enforcement / intelligence）
- [ ] 系統真實使用此 contract，**Phase 3 完成後至少 5 個 task 在 final report 列 Cognitive Mode** 驗證

Phase 4-5 是優化（cost、adaptive），不是架構驗證 — 完成 Phase 3 即可評估 ADR promotion。

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

---

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

- [ ] 計畫書狀態：draft → in-progress（Phase 0 通過後）→ completed
- [ ] 記錄完成日期
- [ ] 記錄偏離（如有）
- [ ] 歸檔至 `plans/archived/`

### README 更新

- [ ] `README.md` OS Layout 表格新增 cognitive modes 機制簡述（不改世代名稱）
- [ ] `architecture/ai-native-cognitive-execution-system.md` 核心機制章節新增 cognitive modes
- [ ] `models/README.md` 加說明：profile 與 mode 的關係（backward-compat）

### 架構文件（per system-upgrade-governance §3 規則 6）

- [ ] Plan completed 時評估 ADR Promotion Criteria
- [ ] 若符合 → 建立 ADR（直接 accepted）並更新 `constitution/README.md` 與 `architecture/ai-native-cognitive-execution-system.md` §本世代相關 ADR
- [ ] 若不符合 → 在 plan 結尾記錄「決定不升 ADR 的理由」+ 改用 runtime gate / enforcement / intelligence 作為 promotion target

### 索引與路由

- [ ] `knowledge/runtime/routing-registry.yaml` 新增 `route.runtime.cognitive-modes`
- [ ] `knowledge/indexes/README.md` 反映新結構
- [ ] `knowledge/graphs/` 新增 cognitive-modes 與 phase/compression/memory/governance 的 edges

### Runtime Surface

- [ ] 新增 `runtime/cognitive-modes.yaml`（executable YAML contract）
- [ ] `runtime/runtime.db` 新增 `cognitive_modes` + `discovery_signals` 表
- [ ] Compiler 規則更新，投影到 `generated_surfaces`
- [ ] `runtime/runtime.db` 的 `phase_machine` / `obligation_ledger` / `blocking_gates` 對接 mode

### Linked Updates

- [ ] `governance/lifecycle/knowledge-update-flow.yaml` 評估是否加 mode resolution step（Phase 3 決定）
- [ ] `enforcement/failure-patterns/` 加 `cognitive-mode-resolution-bypass.md`（mode 未解析就執行）
- [ ] `models/compression/README.md` 標註與 `context_mode` 的對應關係
- [ ] `memory/retrieval-governance/` 標註與 `memory_mode` 的對應關係

### 閉環驗證

- [ ] Diff review
- [ ] Linked updates 完成
- [ ] `ai-skill runtime compile --native-compiler` + `ai-skill runtime validate` 通過
- [ ] Commit / push / readback

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

| 指標 | 通過條件 |
|------|------|
| 累積任務數 | ≥ 5 個 final report 含 Cognitive Mode 區塊 |
| Mode 組合多樣性 | 至少覆蓋 3 種不同 execution + 3 種不同 context + 2 種 governance |
| 誤判率 | ≤ 20% 的 mode 選擇被使用者糾正 |
| Raw signal 完整性 | 沒有「找不到適合 signal 表達」的情況 |
| 設計穩定性 | mode 定義在 5 個任務內未變動超過 2 次 |

### Phase D 完成條件

- [x] `models/cognitive-modes/README.md` 寫入
- [x] `models/README.md` 加入口
- [x] Plan 加 Phase D 段落（本段）
- [x] 至少 1 個任務（本 commit 即是）在 final report 列 Cognitive Mode
- [ ] 累積 5 個任務的 mode 報告
- [ ] 評估指標通過
- [ ] 決定下一步：進 Phase 0 / 修設計 / 撤回 plan

### Phase D Rollback

| 完整撤回 | `git revert <Phase D commit hash>` 移除 cognitive-modes/README.md、models/README 入口、本段落；無 runtime state 變更 |
| 暫停 trial | 在本段加 `paused` 標記，agent final report 不再列 Cognitive Mode |
| 修 mode 定義 | 編輯 `models/cognitive-modes/README.md`；下次任務套用新定義 |

---

## Phase 0: Pre-Build Interrogation + Architecture Compatibility Preflight

**前置條件**：Phase D 完成且通過評估指標後才開始。依 [`plans/README.md` §Architecture Compatibility Preflight](../README.md) 完成。

### Interrogation 草稿

- **Goal**：將 models/ 從 design layer 提升為 runtime activation layer，補 governance/memory 強度差異化缺口
- **Scope**：4 維 cognitive mode primitive + discovery heuristics + runtime.db 投影 + 5 個 subsystem 整合
- **Non-goals**：
  - 不重寫 `models/profiles/` / `compression/` / `capabilities/`
  - 不改世代命名（保持 Gen 3）
  - 不改 knowledge-update-flow.yaml 11 步基本結構
- **Acceptance**：Phase 1 完成後，至少一個任務能在 final report 列 `Mode: {execution, context, governance, memory}` 並由 runtime.db 強制
- **Framework discovery**：執行前須讀 runtime.db 既有 schema、所有 `models/` 文件、`memory/retrieval-governance/`、`governance/ai-runtime-governance/`
- **Duplication risk**：context_mode vs compression level 命名重複；需在 Phase 1 決定合併或保留
- **Open questions**：見上方 §Open Questions
- **Assumptions**：runtime.db 可擴充 schema 而不破壞既有 generated surfaces

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

### Phase 1 完成條件

- [ ] `runtime/cognitive-modes.yaml` 寫入並通過 schema validation
- [ ] `runtime.db cognitive_modes` 表存在
- [ ] `ai-skill runtime compile --native-compiler` 能投影
- [ ] `ai-skill runtime validate` 通過
- [ ] 至少手動寫入一個 task 的 mode resolution（POC）

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

### Phase 2 完成條件

- [ ] Discovery 訊號規則表完整
- [ ] 在 runtime.db `discovery_signals` 表中可查詢
- [ ] 規則 fallback：訊號不命中時用 default mode
- [ ] POC：5 個常見任務類型測試 discovery 輸出符合預期

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

- [ ] 4 subsystem 各有對應 mode-driven activation 邏輯
- [ ] 失敗測試：mode 未解析就執行 → 被阻擋
- [ ] Documentation 更新（含 `models/compression/` 改寫為 `context_mode` implementation reference，per Open Question 1）
- [ ] Phase 3 完成後至少 5 個任務的 final report 列 Cognitive Mode（用於 ADR promotion 評估）
- [ ] **ADR Promotion Criteria 評估**（per Open Question 5 resolved）：
  - 通過 → 建立 ADR（直接 accepted），引用本 plan 作為 evidence
  - 不通過 → 在 plan 內記錄「決定不升 ADR 的理由」+ 改用更輕 promotion target（runtime gate / enforcement / intelligence）

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

- [ ] Budget table 寫入 runtime.db
- [ ] Budget gate 在 runtime 強制
- [ ] Failure pattern `cognitive-budget-overflow.md` 建立

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

- [ ] Adaptive 規則表
- [ ] runtime 偵測訊號 → 自動切 mode
- [ ] 測試 case：模擬 contradiction / failure / budget overflow 三種情境

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
- [ ] Phase D trial 結束後評估指標
- [ ] Phase 0 Pre-Build Interrogation 草稿審閱（Phase D 通過後）
- [ ] Phase 1 實作前 candidate files 與 runtime.db schema 擴充影響範圍評估

---

## 與其他 plans 的關係

| Plan | 關係 |
|------|------|
| [`plans/archived/2026-05-20-1802-model-aware-execution-routing.md`](../archived/2026-05-20-1802-model-aware-execution-routing.md) | 前作；定義 model-aware routing 設計層；本 plan 推進到 runtime activation 層 |
| [`plans/archived/2026-05-22-0855-executable-yaml-contract-migration.md`](../archived/2026-05-22-0855-executable-yaml-contract-migration.md) | 提供 executable contract 機制基礎；本 plan 依賴此機制 |
| [`plans/archived/2026-05-20-1745-memory-retrieval-activation-governance.md`](../archived/2026-05-20-1745-memory-retrieval-activation-governance.md) | 提供 memory activation 邊界；本 plan 將 retrieve threshold 與 memory_mode 整合 |
| [`plans/archived/2026-05-20-1039-runtime-recovery-escalation-system.md`](../archived/2026-05-20-1039-runtime-recovery-escalation-system.md) | RECOVERY mode 對應此 system |
| [`plans/archived/2026-05-15-0920-runtime-execution-layer-upgrade-analysis.md`](../archived/2026-05-15-0920-runtime-execution-layer-upgrade-analysis.md) | Gen 3 升級分析；本 plan 是 Gen 3 子系統擴充 |
