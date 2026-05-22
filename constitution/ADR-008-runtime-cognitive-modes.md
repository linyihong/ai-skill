# ADR-008: Runtime Cognitive Modes

## Status

**Proposed**（2026-05-22）

> ⚠️ 本 ADR 為提案階段，**尚未 accepted**。需經審查、討論與 alternatives 評估後才推進。對應實作計畫見 [`plans/active/2026-05-22-1629-runtime-cognitive-modes-system.md`](../plans/active/2026-05-22-1629-runtime-cognitive-modes-system.md)（status: draft）。

## Framework Generation

- **世代分類**：Gen 3 子系統擴充（不是 Gen 4 升級）
- **當前世代文件**：[`architecture/ai-native-cognitive-execution-system.md`](../architecture/ai-native-cognitive-execution-system.md)
- **適用狀態**：若 accepted，會擴充 Gen 3 的 runtime infrastructure，**不改變**核心架構（runtime.db、executable YAML contract、knowledge update flow 11-step 仍保留）。

## Context

外部架構審查指出 `models/` 層目前是 **documentation layer 而非 runtime activation layer**。具體證據：

1. **沒有 blocking gate 強制查詢**：agent 在執行任務時不會自動查 `routing-registry.yaml` → `model-context-report.md` → `model-checklists.md`。`knowledge-update-flow.yaml` 11 步沒有任何一步要求查詢 model profile。
2. **本 session 實證**：本 session 內 agent 加入 4 個 intelligence atoms、修正 ADR-007 語言、寫第三代 architecture 文件、ADR 雙向連結等工作，**沒有一次**查詢 model-context-report 就直接執行；profile 報告 / Read / Deferred / Validation signal 從未出現在 final report。
3. **Document lookup runtime 不可行**：對方批評「每次 full resolution 一定爆」是真實 token cost 問題。若加 Step 0「每任務 query registry + profile + report + checklist」，每次 ~2000 tokens overhead 無法承受。
4. **真正的缺口**：對照 4 個 cognitive primitive 維度，**governance mode 強度差異化**與 **memory mode activation flag** 是 Gen 3 既有 infrastructure 沒覆蓋到的真實缺口。

對照現有系統：

| 建議 mode | 既有對應 | 缺口性質 |
|------|------|------|
| execution mode (FAST/NORMAL/DEEP/FORENSIC/RECOVERY) | `runtime.db phase_machine` 有 phase 概念，但**沒有 cognitive depth 維度** | 60% 新（FORENSIC/RECOVERY 為真新增） |
| context mode (INDEX_ONLY/SUMMARY_FIRST/CHECKLIST_FIRST/SOURCE_BACKED/FULL_TRACE) | `models/compression/` 5 個 level 名稱幾乎一樣（小寫） | 5% 新（rename + 提升為 runtime primitive） |
| governance mode (LIGHT/STANDARD/STRICT/LOCKDOWN) | 既有 governance 是 binary（gate 或無 gate） | **80% 新** — 真正缺口 |
| memory mode (NONE/EPISODIC/DECISION_REPLAY/FAILURE_REPLAY/PROJECT_CONTEXT) | `memory/` 子層存在但無 activation flag | **70% 新** — 把 memory 子層提升為 runtime mode |

## Decision

引入 **Runtime Cognitive Modes** 作為 Gen 3 runtime infrastructure 的**子系統擴充**，核心三點：

### 1. 4 維 mode primitive 而非 flat profile

從 `small / large / specialized` 升級為 4 維 composable primitive：

```
execution_mode  ∈ {FAST, NORMAL, DEEP, FORENSIC, RECOVERY}
context_mode    ∈ {INDEX_ONLY, SUMMARY_FIRST, CHECKLIST_FIRST, SOURCE_BACKED, FULL_TRACE}
governance_mode ∈ {LIGHT, STANDARD, STRICT, LOCKDOWN}
memory_mode     ∈ {NONE, EPISODIC, DECISION_REPLAY, FAILURE_REPLAY, PROJECT_CONTEXT}
```

組合空間 = 5 × 5 × 4 × 5 = 500 種 cognitive state，遠優於既有 3 profile 的 capability fan-out。

### 2. Discovery 用快速啟發式，不查文件

```
任務進來
  ↓
快速 discovery（純訊號計算，< 50 tokens）：
  - user keyword（「改規則」「補 typo」「分析架構」「recover」「replay」）
  - file diff scope（enforcement/ vs notes/ vs ADR）
  - git status（dirty / clean / ahead-behind）
  - session 長度 / contradiction risk / recent failure
  ↓
單一 SQLite 查詢解析 mode：
  SELECT modes FROM cognitive_modes WHERE signal_pattern = ?
  ↓
4 個 mode 寫入 runtime state
  ↓
執行任務時，各 subsystem 依 mode 決定 activation：
  - execution_mode 決定 phase machine depth
  - context_mode 決定載入策略（既有 compression 邏輯）
  - governance_mode 決定哪些 gate 啟用
  - memory_mode 決定哪些 memory 子層 retrieve
```

### 3. 與既有 infrastructure 整合，不重寫

| 既有 | 新增 / 整合方式 |
|------|------|
| `runtime.db phase_machine` | 新增 `cognitive_modes` 表，與 phase_machine join |
| `models/compression/` | Rename levels 為 UPPERCASE，提升為 `context_mode` 的實作層 |
| `models/profiles/` | 保留為 reference doc；mark 為「映射到 cognitive modes 的 backward-compat label」 |
| `memory/<subdir>/` | 子層不動，新增 `memory_mode` 作為 activation flag |
| `governance/` | 不動現有 gate；新增 `governance_mode` 作為 gate 強度 selector |
| `models/capabilities/` | 保留為 cognitive primitive 的細粒度 descriptor，可用來 fine-tune mode |
| `models/routing/` | 保留為 task class / autonomy 路由設計層 |

**不刪除任何既有檔**。

### 4. 漸進實作（5 phase）

| Phase | 範圍 | 完成條件 |
|------|------|------|
| 1 | Cognitive Mode primitive 定義（YAML contract） | 4 mode 的可投影 YAML 在 runtime.db `generated_surfaces` |
| 2 | Discovery heuristics（不查文件） | 快速訊號 → mode 映射規則表寫入 runtime.db |
| 3 | Mode → activation 整合 | 4 subsystem（phase / compression / governance / memory）依 mode 自動 activate |
| 4 | Selective loading + token budget gate | Source 載入由 mode 控制，超 budget 阻擋 |
| 5 | Adaptive runtime | contradiction / recovery depth / failure repeat 動態調整 mode |

每 phase 獨立 archive 後才開下一個。Phase 1 是必要起點，Phase 5 是長期目標。

## Consequences

### 正面

- **真正 runtime activation**：mode 寫入 runtime state，subsystem 強制依 mode 行動
- **Token cost 可控**：discovery 靠訊號不靠文件；activation 是 conditional
- **Governance 強度差異化**：LIGHT 模式跳過部分 gate，STRICT 全跑，LOCKDOWN 阻擋寫入
- **Memory mode activation**：明確區分「不查記憶」「查 episodic」「replay decision」等狀態
- **4 維 composable**：500 種狀態組合，每種狀態都有明確 activation 邏輯
- **Backward compat**：既有 profile / compression / memory 文件保留為 reference

### 負面

- **runtime.db schema 擴充**：新增 `cognitive_modes` + `discovery_signals` 表，需 compiler 規則更新
- **Discovery heuristic 維護成本**：訊號 → mode 映射規則需要持續校準
- **詞彙重疊風險**：`context_mode` 與既有 `compression` 名稱幾乎相同；需明確 deprecation 或統一
- **4 phase 是大改動**：可能需要 1-3 個月推進，期間既有系統需保持穩定
- **教學負擔**：新概念（cognitive mode、discovery signal）需文件化

### 風險

| 風險 | 緩解 |
|------|------|
| 4 個 mode 互動空間太大（500 狀態）難以全測試 | Phase 1 先定義「常見組合」，未列組合走 default fallback |
| Discovery heuristics 誤判 → 用錯 mode | Mode 內加 escalation 規則，偵測訊號矛盾時自動升級 |
| 與既有 phase_machine 概念衝突 | Phase 3 整合時明確定義「phase 是 transaction state，mode 是 cognitive state」邊界 |
| 重構期間既有任務行為改變 | Phase 1-2 只新增 surface，不改現有 activation；Phase 3 起才切換 |

## Alternatives Considered

- **A. 維持現狀（models/ 純 documentation）**：拒絕 — 已證實 agent 不 activate，違背 design intent。
- **B. 完全重寫 `models/`**：拒絕 — 既有 compression / capabilities / routing 有獨立 reference 價值；重寫會破壞 ADR-001 reference-first 原則。
- **C. 只加 governance mode 解決最緊迫缺口**：拒絕 — 單一 dimension 無法解決 compositional needs；4 維才有真正組合空間。
- **D. 加進 knowledge-update-flow Step 0「每任務查 profile」**：拒絕 — token cost 不可行（對方批評正確）。
- **E. 4 mode + 整合既有 + 5 phase 漸進**：**accept**。

## Open Questions（accepted 前需釐清）

1. `context_mode` 與 `compression level` 是否合併命名（避免雙詞彙）？建議：合併為 `context_mode`，`compression/` 文件改稱「context mode 的詳細策略文件」。
2. Discovery signal 從哪些來源讀？目前候選：user keyword、file paths、git status、session 長度、recent failure count。是否還缺？
3. `governance_mode LOCKDOWN` 與既有 `runtime.db blocking_gates` 的關係？建議：LOCKDOWN 是「全部 gates + 阻擋寫入」的 superset。
4. `memory_mode` 與 `memory/retrieval-governance/` 的 activation threshold 是否整合？建議：retrieval-governance 是 threshold，mode 是 activation switch；兩者 compose。
5. Phase 1 完成後是否要建立 ADR-008 → Accepted 的 promotion gate？建議：是 — Phase 1 視為 ADR proof-of-concept。

## Related

- [`architecture/ai-native-cognitive-execution-system.md`](../architecture/ai-native-cognitive-execution-system.md) — Gen 3 當前世代文件
- [`models/README.md`](../models/README.md) — 既有 models 層
- [`models/profiles/README.md`](../models/profiles/README.md) — 既有 small/large/specialized profile
- [`models/compression/README.md`](../models/compression/README.md) — 既有 5 級 compression
- [`memory/README.md`](../memory/README.md) — 既有 6 子層 memory
- [`runtime/runtime.db`](../runtime/runtime.db) — 將接收 cognitive_modes 表
- [`governance/lifecycle/knowledge-update-flow.yaml`](../governance/lifecycle/knowledge-update-flow.yaml) — 11 step master flow，未來可加入 mode resolution step
- [`governance/lifecycle/system-upgrade-governance.md`](../governance/lifecycle/system-upgrade-governance.md) — 本 ADR 涉及架構分層 + 核心流程變更，需依此治理
- [`plans/active/2026-05-22-1629-runtime-cognitive-modes-system.md`](../plans/active/2026-05-22-1629-runtime-cognitive-modes-system.md) — 對應實作 plan
- [ADR-002](ADR-002-intelligence-vs-knowledge-separation.md) — intelligence 與 knowledge 分離（本 ADR 不改動）
- [ADR-005](ADR-005-memory-architecture.md) — memory 6 子層（本 ADR 將其提升為 memory mode）
- [ADR-007](ADR-007-constitution-and-decision-promotion-boundary.md) — promotion target boundary（本 ADR 走 architecture-level 路徑）
