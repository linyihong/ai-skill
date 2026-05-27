# ADR-008: Runtime Cognitive Modes

## Status

**Accepted**

## Framework Generation

- **世代分類**：Gen 3 子系統擴充（不是世代升級）
- **當前世代文件**：[`architecture/ai-native-cognitive-execution-system.md`](../architecture/ai-native-cognitive-execution-system.md) §12 Cognitive Modes
- **適用狀態**：4 維 cognitive mode primitive 是 Gen 3 runtime infrastructure 的擴充；既有 `runtime/runtime.db` phase_machine、blocking_gates、memory 子層與 compression 層保留，新增 mode 作為跨層的 cognitive switch。

## Date

2026-05-25

## Source Plan

[`plans/archived/2026-05-22-1629-runtime-cognitive-modes-system.md`](../plans/archived/2026-05-22-1629-runtime-cognitive-modes-system.md)

完整 Decision Rationale、Open Questions、Phase 設計、Acceptance、Trial 證據與 ADR Promotion Criteria 評估皆保留於原 plan。本 ADR 為 Phase 5 完成、5/5 criteria 通過後的 accepted promotion。

## Context

2026-05-22 外部架構審查指出 `models/` 層當時是 documentation layer 而非 runtime activation layer。具體缺口：

1. 沒有 blocking gate 強制查詢 model profile / context report / checklist
2. 實證觀察：agent 執行任務時不查 model 文件直接行動；profile / Read / Deferred / Validation signal 從未出現在 final report
3. 「每任務 full lookup」會爆 token cost，不可行
4. 真正的缺口在 **governance mode 強度差異化** 與 **memory mode activation flag**，而非單純的 model selection

對照 4 個建議 mode 維度與既有 Gen 3 infrastructure：

| 建議 mode | 既有對應 | 缺口性質 |
|---|---|---|
| `execution_mode`（FAST/NORMAL/DEEP/FORENSIC/RECOVERY） | `phase_machine` 有 phase 但無 cognitive depth 維度 | 60% 新（FORENSIC/RECOVERY 真新增） |
| `context_mode`（INDEX_ONLY/SUMMARY_FIRST/CHECKLIST_FIRST/SOURCE_BACKED/GRAPH_ASSISTED） | `models/compression/` 同名 lowercase 5 level | 5% 新（rename + 提升為 runtime primitive，UPPERCASE alias） |
| `governance_mode`（LIGHT/STANDARD/STRICT/LOCKDOWN） | 既有 governance 是 binary（gate 或無 gate） | 80% 新 — 真實缺口 |
| `memory_mode`（NONE/EPISODIC/DECISION_REPLAY/FAILURE_REPLAY/PROJECT_CONTEXT） | `memory/` 子層存在但無 activation flag | 70% 新 |

## Decision

引入 **Runtime Cognitive Modes** 作為 Gen 3 runtime infrastructure 的子系統擴充，三大核心：

### 1. 4 維 mode primitive（500 種 composable state）

```
execution_mode  ∈ {FAST, NORMAL, DEEP, FORENSIC, RECOVERY}      (5)
context_mode    ∈ {INDEX_ONLY, SUMMARY_FIRST, CHECKLIST_FIRST,
                   SOURCE_BACKED, GRAPH_ASSISTED}                (5)
governance_mode ∈ {LIGHT, STANDARD, STRICT, LOCKDOWN}            (4)
memory_mode     ∈ {NONE, EPISODIC, DECISION_REPLAY,
                   FAILURE_REPLAY, PROJECT_CONTEXT}              (5)
```

每個任務 resolve 4 個 mode，寫入 `runtime.db cognitive_modes` 表，並由各 subsystem 強制依 mode 行動。

### 2. Discovery-based resolution（不靠 file lookup）

`runtime/cognitive-modes-discovery.yaml` 定義 14 個 signals（covering 8 signal_types），從 task entry signal（file diff path、cognitive load、failure recurrence、user keyword 等）解析 mode，**不**每任務查 model 文件。Fallback 為 NORMAL / SUMMARY_FIRST / STANDARD / NONE。

### 3. Subsystem integration via integration YAMLs

每個 mode 對應一個 integration contract：

- `cognitive-modes-phase-integration.yaml` — execution_mode → phase_machine allowed/forbidden_actions
- `cognitive-modes-governance-integration.yaml` — governance_mode → gate_set_activation
- `cognitive-modes-memory-integration.yaml` — memory_mode → memory subdir activation
- `cognitive-modes-token-budget.yaml` — mode tuple → token budget ceiling
- `cognitive-modes-adaptive.yaml` — runtime adaptive trigger rules

並由 commit-msg hook（Go validators in `scripts/ai-skill-cli/internal/app/hooks.go`）行為層強制：每個 commit body 必須含 `### Cognitive Mode 報告` block；mode 之間必須符合 floor / consistency / subdir constraint；超 budget block；adaptive trigger 失配 block。

### 4. Runtime Cognitive Contract v2（2026-05-27 amendment）

ADR-008 的 reporting contract 升級為 6 維 + adaptive disclosure：

```
execution_mode / context_mode / governance_mode / memory_mode
validation_mode / cognitive_cost
```

新增規則：

- `validation_mode` 將 validation depth 與 execution depth 解耦。
- `cognitive_cost` 由 `runtime/cognitive-modes-cost-class.yaml` 的 execution × context lookup 推導，不接受 agent self-claim。
- 全 6 維 default 時可用 compact single-line form；任一非 default 或 high-risk mode 必須用 full form。
- `activation_reason` 必須引用 `runtime/cognitive-modes-discovery.yaml` 的 known signal。
- high-risk mode（DEEP / FORENSIC / RECOVERY / STRICT / LOCKDOWN）必須附 `Capability summary`，讓 mode label 連回 capability semantics。

此 amendment 保留 ADR-008 的 4 維 runtime primitive 與 subsystem integration，只升級 reporting / validation contract；不建立 ADR-009。

## Consequences

### 正面

- **真正 runtime activation**：mode 寫入 SQLite + 由 commit-msg hook 行為層強制（不再僅是 doc layer）
- **Token cost 可控**：discovery 靠 signals 不靠 file lookup；budget table 每 mode tuple 上限明確
- **Governance 強度差異化**：LIGHT / STANDARD / STRICT / LOCKDOWN 對應不同 gate set
- **Memory mode activation 明確**：NONE 禁讀 memory；其他 mode 對應特定 subdir + AND-compose with activation thresholds
- **4 維 composable**：500 種 state，覆蓋 typo fix → forensic audit 全光譜
- **Backward compat**：既有 compression layer 保留為 lowercase alias；既有 memory subdir / governance gate 保留為 mode-activated 對象
- **Self-dogfood verified**：本 plan 自身的 Phase 3-5 commits 透過 commit-msg validators 通過或被 block + 修正，形成閉環

### 負面

- **Mode resolution overhead**：每任務 first turn 必須 resolve 4 個 mode（透過 discovery signals + Cognitive Mode block 在 commit body）— 但相比每任務查 model 文件，cost 仍低
- **500 種組合 composition complexity**：integration contracts 必須維護 floor / consistency / activation rules；若新增 mode value 需同時更新 5 個 YAML
- **Commit-msg layer 邊界**：3 個 adaptive triggers（tool_call_loop_detection、cross_session_memory_hit_rate、phase_transition_remap）需 per-tool-call 攔截，無法在 commit-msg 層 enforce；目前文件化但未啟用

## Alternatives Considered

1. **Flat profile（單一 model profile selector）** — 拒絕，無法 compose；governance + memory + execution 維度互相獨立，flat profile 表達力不足。
2. **每任務 full model lookup** — 拒絕，每次 ~2000 tokens overhead 不可行（外部審查明確指出）。
3. **保留 doc-only contract（不做 runtime activation）** — Phase D 8 commits 已證 doc-only 設計穩定，但缺少 enforcement → behavior drift 風險高。
4. **僅做 governance_mode + memory_mode（最大缺口維度）** — 拒絕，4 維 composable 才能涵蓋從 FAST 到 FORENSIC 全光譜；少做 2 維會被迫用 ad-hoc gate 補回。
5. **延後到 Gen 4 世代升級** — 拒絕，是子系統擴充非世代切換；既有 Gen 3 infrastructure（phase_machine / blocking_gates / memory / compression）皆相容，沒有 backward-incompatible 改動。

## ADR Promotion Criteria 評估結果

5/5 criteria 全部通過（2026-05-25 評估）：

| # | 條件 | 證據 |
|---|---|---|
| 1 | foundational + cross-session + cross-project + expensive-to-reverse + explains-why | 4 維 mode primitive + 500 組合 + 跨 subsystem integration |
| 2 | Phase 3 完成（4 subsystem 真實 activation） | 3.1-B/3.3-B/3.4-B Go validators 全部上線於 commit-msg hook（`scripts/ai-skill-cli/internal/app/hooks.go`） |
| 3 | 5 Open Questions 全解 | resolved 於 2026-05-22 |
| 4 | 沒有更輕的 promotion target 適用 | 4 維 primitive + 500 組合非 single gate / atom 可代替 |
| 5 | ≥5 task final report 列 Cognitive Mode | 14 commits（Phase D 8 + Plan execution 6，含 2 個 non-setup tasks） |

## Linked Validation Scenarios

`validation/scenarios/cognitive-modes/` 與 `validation/scenarios/bootstrap/` 共 26 個 cognitive scenarios，全 PASS；v2 新增 scenarios 覆蓋 compact / full / activation signal / cost class / capability snippet / inflated rejection：

- yaml-contract-exists / generated-surface-projected / runtime-table-exists / poc-task-record
- discovery-yaml-exists / discovery-signals-projected / discovery-fallback-defined / discovery-poc-coverage
- phase-integration / compression-alias / governance-integration / memory-integration / enforcement-gate-exists
- cognitive-mode-block-required / phase3-behavioral-validators / plan-status-sync-enforcement
- phase4-token-budget / phase5-adaptive-triggers / bootstrap-receipt-enforcement
- phase6-cognitive-contract-v2-compact-form / full-form / activation-signal / cost-class / capability-snippet / inflated-rejection

## Related ADRs

- **ADR-005 Memory Architecture** — `memory_mode` 把 6 子層 memory 模型提升為 runtime activation primitive
- **ADR-006 Registry-First Workflow Activation** — discovery signals 為 mode resolution 的 registry-first 入口
- **ADR-007 Constitution and Decision Promotion Boundary** — 本 ADR 依 §No-Proposed-ADR Rule 從 plan 直接 promotion 到 accepted（無 proposed 階段）
