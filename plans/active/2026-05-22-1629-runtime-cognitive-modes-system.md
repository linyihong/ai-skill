# Runtime Cognitive Modes System

**Status**: `draft`
**ADR**: [ADR-008](../../constitution/ADR-008-runtime-cognitive-modes.md)（proposed）
**世代**：Gen 3 子系統擴充
**建立日期**：2026-05-22

> ⚠️ **本 plan 為 draft 階段**，依賴 [ADR-008](../../constitution/ADR-008-runtime-cognitive-modes.md) accepted 後才可進入 Phase 0 architecture compatibility preflight 與 Phase 1 實作。

---

## 緣起與動機

### 觸發來源

2026-05-22 外部架構審查指出 `models/` 層目前是 documentation layer 而非 runtime activation layer。本 session agent 從未 query model-context-report 就執行任務即為實證。詳細 context 見 [ADR-008 §Context](../../constitution/ADR-008-runtime-cognitive-modes.md#context)。

### 核心問題

| 問題 | 證據 |
|------|------|
| Model layer 沒有 runtime activation | knowledge-update-flow 11 步無任何「查 model profile」步驟 |
| Discovery cost 無法用 document lookup 解決 | 每任務 ~2000 tokens overhead 不可行 |
| Governance 無強度差異化 | 既有 governance gate 是 binary（gate / no-gate） |
| Memory 子層無 activation flag | memory/ 子層完整但 agent 不知何時 activate |

### 解決方向

引入 4 維 cognitive mode primitives（execution / context / governance / memory），與既有 runtime infrastructure 整合，**不重寫** models/。

---

## 完成條件（per system-upgrade-governance §2）

本 plan 是 Gen 3 子系統擴充，**不是世代升級**，所以系統名稱不變。但仍涉及架構分層擴充 + 核心流程變更，需執行 §2 checklist 的子集。

### 計畫書本身

- [ ] 計畫書狀態：draft → in-progress（ADR-008 accepted 後）→ completed
- [ ] 記錄完成日期
- [ ] 記錄偏離（如有）
- [ ] 歸檔至 `plans/archived/`

### README 更新

- [ ] `README.md` OS Layout 表格新增 cognitive modes 機制簡述（不改世代名稱）
- [ ] `architecture/ai-native-cognitive-execution-system.md` 核心機制章節新增 cognitive modes
- [ ] `models/README.md` 加說明：profile 與 mode 的關係（backward-compat）

### 架構文件

- [ ] `constitution/ADR-008` 由 proposed → accepted
- [ ] `architecture/ai-native-cognitive-execution-system.md` 「本世代相關 ADR」加 ADR-008 列

### 索引與路由

- [ ] `knowledge/runtime/routing-registry.yaml` 新增 `route.runtime.cognitive-modes`
- [ ] `knowledge/indexes/README.md` 反映新結構
- [ ] `knowledge/graphs/` 新增 cognitive-modes 與 phase/compression/memory/governance 的 edges

### Runtime Surface

- [ ] 新增 `runtime/cognitive-modes.yaml`（executable YAML contract）
- [ ] `runtime/runtime.db` 新增 `cognitive_modes` + `discovery_signals` 表
- [ ] Compiler 規則更新，將 yaml 投影到 `generated_surfaces`
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

## Phase 0: ADR Acceptance + Pre-Build Interrogation

**前置條件**：ADR-008 accepted 後才開始。

依 [`plans/README.md` §Architecture Compatibility Preflight](../README.md) 完成：

| 欄位 | 內容 |
|------|------|
| Trigger | 開始實作 cognitive modes system |
| Checked sources | runtime.db schema、phase_machine、compression、memory、governance 既有結構 |
| Conflicts | 待 Phase 0 執行時填寫 |
| Interrogation | Goal / Scope / Non-goals / Acceptance / Framework discovery / Duplication risk / Open questions |
| Decision | 待 Phase 0 執行時填寫 |
| Validation | runtime query、validator、link check、readback |

### Interrogation 草稿

- **Goal**：將 models/ 從 design layer 提升為 runtime activation layer，補 governance/memory 強度差異化缺口
- **Scope**：4 維 cognitive mode primitive + discovery heuristics + runtime.db 投影 + 5 個 subsystem 整合
- **Non-goals**：
  - 不重寫 `models/profiles/` / `compression/` / `capabilities/`
  - 不改世代命名（保持 Gen 3）
  - 不改 knowledge-update-flow.yaml 11 步基本結構（mode resolution 視 Phase 3 結果決定加 step）
- **Acceptance**：Phase 1 完成後，至少一個任務能在 final report 列 `Mode: {execution, context, governance, memory}` 並由 runtime.db 強制
- **Framework discovery**：執行前須讀 runtime.db 既有 schema、所有 `models/` 文件、`memory/retrieval-governance/`、`governance/ai-runtime-governance/`
- **Duplication risk**：context_mode vs compression level 命名重複；需在 Phase 1 決定合併或保留兩個詞彙
- **Open questions**：見 ADR-008 §Open Questions 1-5
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
    descriptions:
      FAST: 快速回答、不查 source；適用 typo / 簡單問答
      NORMAL: 標準流程；checklist + source 視需要
      DEEP: 跨層分析、完整 dependency；適用 migration / promotion
      FORENSIC: 高 lineage tracing；適用 incident analysis / audit
      RECOVERY: 高 validation、阻擋寫入；適用 failure recovery
  context:
    values: [INDEX_ONLY, SUMMARY_FIRST, CHECKLIST_FIRST, SOURCE_BACKED, FULL_TRACE]
    default: SUMMARY_FIRST
    # 與既有 models/compression/ 對應
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

resolution:
  required_when:
    - any task starts
  output:
    - execution_mode
    - context_mode
    - governance_mode
    - memory_mode
```

### 1.2 設計 `runtime.db cognitive_modes` 表

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

### 2.1 訊號來源

| 訊號 | 例 | 對應 mode 暗示 |
|------|------|------|
| user keyword | 「快速」「typo」「修一下」 | execution=FAST |
| user keyword | 「分析」「跨層」「migration」 | execution=DEEP |
| user keyword | 「recover」「失誤」「重做」 | execution=RECOVERY |
| file diff scope | `enforcement/` / `governance/` | governance=STRICT |
| file diff scope | `notes/` / `memory/working/` | governance=LIGHT |
| git status | dirty + 多 owner group | governance=STRICT |
| session 長度 | > 50 turns | memory=PROJECT_CONTEXT |
| recent failure | failure_repeat ≥ 2 | execution=RECOVERY, memory=FAILURE_REPLAY |
| contradiction risk | conflicting sources | context=FULL_TRACE, governance=STRICT |

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

合併 `models/compression/` levels 為 `context_mode` 的實作層。決定保留兩個詞彙或統一（依 ADR-008 §Open Question 1）。

### 3.3 governance_mode → blocking_gates

`runtime.db blocking_gates` 依 governance_mode 啟用對應 gate set。LOCKDOWN 額外阻擋寫入。

### 3.4 memory_mode → memory retrieve

`memory/retrieval-governance/` 的 activation threshold 與 memory_mode compose；mode 是 switch，threshold 是 limit。

### 3.5 knowledge-update-flow 整合

決定是否加 Step 0「mode resolution」。建議：**不加 step，改在 Step 1 entry 自動觸發**（避免膨脹 11 步）。

### Phase 3 完成條件

- [ ] 4 subsystem 各有對應 mode-driven activation 邏輯
- [ ] 失敗測試：mode 未解析就執行 → 被阻擋
- [ ] Documentation 更新

---

## Phase 4: Selective Loading + Token Budget Gate

**目標**：Source 載入由 mode 控制，超 budget 自動降級或阻擋。

### 4.1 Token budget per mode

| Mode 組合 | Budget |
|------|------|
| FAST + INDEX_ONLY + LIGHT + NONE | ≤ 1000 tokens |
| NORMAL + SUMMARY_FIRST + STANDARD + EPISODIC | ≤ 5000 tokens |
| DEEP + SOURCE_BACKED + STRICT + DECISION_REPLAY | ≤ 20000 tokens |
| FORENSIC + FULL_TRACE + STRICT + FAILURE_REPLAY | ≤ 50000 tokens |

### 4.2 Budget gate

超 budget 時：
1. 自動降級 context_mode（FULL_TRACE → SOURCE_BACKED → SUMMARY_FIRST）
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

| 訊號 | 動作 |
|------|------|
| Token budget 接近上限 | 降級 context_mode |
| 偵測到 contradiction | 升 governance_mode + 升 context_mode |
| 連續 failure ≥ 2 | 切到 RECOVERY |
| 跨 phase transition | 重新 evaluate mode |

### 5.2 Adaptive YAML contract

`runtime/cognitive-modes-adaptive.yaml`。

### Phase 5 完成條件

- [ ] Adaptive 規則表
- [ ] runtime 偵測訊號 → 自動切 mode
- [ ] 測試 case：模擬 contradiction / failure / budget overflow 三種情境

---

## 風險與緩解

| 風險 | 緩解 |
|------|------|
| 4 mode 互動空間 500 種狀態難全測試 | Phase 1 先定義「常見組合」表，未列組合走 default |
| Discovery heuristic 誤判 | Mode 內加 escalation；遇到矛盾時自動升級 |
| 與 phase_machine 概念衝突 | Phase 3 明確定義邊界：phase = transaction state, mode = cognitive state |
| Phase 1-2 重構期間既有行為改變 | Phase 1-2 只新增 surface 不改 activation；Phase 3 起才切換 |
| 對方批評變成 paper plan 不執行 | 每 phase 設明確 completion gate，未過不開下一 phase |

## Open Questions（與 ADR-008 同步）

1. `context_mode` vs `compression` 命名統一？
2. Discovery signal 來源是否還有缺？
3. `governance_mode LOCKDOWN` 與既有 `blocking_gates` 關係？
4. `memory_mode` 與 `retrieval-governance` 整合方式？
5. ADR-008 promotion gate 設在哪個 phase 完成時？

## Stakeholder 同意項目

待 review。需要 user 與後續 contributors 對下列事項 sign-off：

- [ ] 同意採 4 維 mode primitive 而非單一 enum
- [ ] 同意不重寫 `models/`，採 backward-compat
- [ ] 同意分 5 phase 漸進實作
- [ ] 同意 ADR-008 從 proposed → accepted 的 gate 設在 Phase 1 完成

---

## 與其他 plans 的關係

| Plan | 關係 |
|------|------|
| [`plans/archived/2026-05-20-1802-model-aware-execution-routing.md`](../archived/2026-05-20-1802-model-aware-execution-routing.md) | 前作；定義 model-aware routing 設計層；本 plan 把它推進到 runtime activation 層 |
| [`plans/archived/2026-05-22-0855-executable-yaml-contract-migration.md`](../archived/2026-05-22-0855-executable-yaml-contract-migration.md) | 提供 executable contract 機制基礎；本 plan 依賴此機制 |
| [`plans/archived/2026-05-20-1745-memory-retrieval-activation-governance.md`](../archived/2026-05-20-1745-memory-retrieval-activation-governance.md) | 提供 memory activation 邊界；本 plan 將 retrieve threshold 與 memory_mode 整合 |
| [`plans/archived/2026-05-20-1039-runtime-recovery-escalation-system.md`](../archived/2026-05-20-1039-runtime-recovery-escalation-system.md) | RECOVERY mode 對應此 system |
| [`plans/archived/2026-05-15-0920-runtime-execution-layer-upgrade-analysis.md`](../archived/2026-05-15-0920-runtime-execution-layer-upgrade-analysis.md) | Gen 3 升級分析；本 plan 是 Gen 3 子系統擴充 |
