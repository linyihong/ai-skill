# Runtime Recovery & Escalation System — System Plan

> **狀態**: draft
> **建立日期**: 2026-05-20
> **目的**: 在 agent 執行中偵測「現實已推翻原假設」的時機，強制停止局部 patch，進入 source-of-truth 重建與 execution graph replanning。

---

## 1. Problem Statement

現有系統已具備 discovery、workflow routing、dependency reading、linked updates、failure learning、validation 與 lazy-load；但這些機制主要處理「任務開始時該讀什麼」與「失敗後如何沉澱 lesson」。

缺口在於 execution 中段：

```text
Agent 已開始執行
  → 實際 UI / API / repo / workflow 證據與原始假設不一致
  → Agent 繼續沿用舊 mental model
  → 局部 patch、automation drift、context 汙染、hallucinated execution flow
```

本 plan 目標不是新增另一份 checklist，而是把「假設崩塌」變成 runtime / enforcement 可辨識、可中止、可恢復的狀態。

---

## 2. Scope

### In Scope

- 新增全庫 escalation policy，定義 mismatch trigger 與 forbidden behavior。
- 將 mismatch escalation 接入 runtime guard chain。
- 建立 recovery mode 的 required actions：suspend assumption、reload workflow source、reload source-of-truth、rebuild execution graph、explain failure。
- 建立 metadata recovery policy，讓不同 domain 可定義 required reload set。
- 將 APK analysis navigation mismatch 作為第一個 validation scenario。

### Out of Scope

- 不重寫所有 workflow。
- 不在業務專案規則中硬編單一 App route。
- 不把 project incident 的具體路徑、host、token、device 或 screenshot evidence 寫入 reusable docs。
- 不要求 agent 每次失敗都進最高級 recovery；level 由 trigger 嚴重度決定。

---

## 3. Layer Responsibility

| Layer | Responsibility | 不放什麼 |
| --- | --- | --- |
| `runtime/discovery/` | 初始 routing、context loading、workflow discovery、dependency loading | mismatch 後的認知重建策略 |
| `enforcement/escalation-policy.md` | mismatch trigger、recovery activation condition、forbidden behavior | 具體 workflow 的完整操作步驟 |
| `runtime/recovery/` | rediscovery、source-of-truth rebuild、execution graph replanning | 長篇解釋與 project evidence |
| `runtime/guards/` | execution 中段偵測 repeated failure / evidence conflict / source miss | domain-specific 操作知識 |
| `metadata/recovery/` | domain policy、recovery levels、required reload set | prose 教學 |
| `workflow/*/execution-flow.md` | domain-specific mismatch handling hook | runtime state machine 細節 |
| `validation/scenarios/failure-derived/` | 驗證 agent 是否會升級，而不是繼續 patch | project raw evidence |

---

## 4. Target State Model

現有簡化模型：

```text
DISCOVERY
→ EXECUTION
→ VALIDATION
```

目標模型：

```text
DISCOVERY
→ EXECUTION
→ VALIDATION
→ MISMATCH_DETECTION
→ ESCALATION
→ RECOVERY
→ REDISCOVERY
→ REPLAN
→ EXECUTION
```

---

## 5. Escalation Triggers

### Trigger Classes

| Class | Trigger | Default level | Required action |
| --- | --- | --- | --- |
| `repeated-failure` | 同一路徑、同一 automation、同一 checkpoint 連續失敗 2 次 | L3 | 停止 retry，reload source-of-truth |
| `user-contradiction` | 使用者指出「不是這樣」「你沒看文件」「你又在猜」 | L4 | 進入 recovery，重建 execution graph |
| `evidence-conflict` | UI / API / repo structure 與 workflow 或 contract 衝突 | L4 | reload workflow primary source 與 owner docs |
| `assumption-drift` | agent 開始靠猜 route、猜座標、沿用 stale checklist | L3 | suspend assumption，補 source reading |
| `source-of-truth-miss` | 重要操作前未讀 canonical workflow / UI map / contract | L4 | rediscovery + dependency read ledger |
| `automation-drift` | 腳本反覆 patch 但 checkpoint 截圖 / foreground / feature context 不對 | L4 | 禁止繼續自動操作，回到 navigation graph |

### Forbidden Behaviors During Escalation

- 繼續 patch automation。
- 繼續猜 UI route、API route、repo architecture。
- 不重讀 workflow primary source。
- 不重建 execution graph。
- 沿用已被反證的 checklist。
- 把 Frida / log / partial evidence 當成 UI context 成功證據。

---

## 6. Recovery Mode Required Actions

Recovery mode 不是「修 patch」，而是「重建 mental model」。

進入 recovery 後必須完成：

1. **Suspend Assumption**
   - 明確寫出舊假設。
   - 明確寫出反證。
   - 標記該假設失效，不得繼續沿用。

2. **Reload Workflow Primary Source**
   - 依 routing registry 或 workflow activation contract 重新載入 primary source。
   - 同步載入 artifact gates、validation rules、required dependencies。

3. **Reload Source-of-Truth**
   - 依 domain policy 讀取 UI map、API catalog、architecture docs、runtime metadata、goal ledger 或 past transcript。
   - 若 source 不存在，標 `not applicable` 或 `source missing`，不得假裝已讀。

4. **Rebuild Execution Graph**
   - 重新建立：

```text
goal
→ workflow route
→ required dependencies
→ source-of-truth
→ execution checkpoint
→ validation evidence
```

5. **Explain Failure**
   - 必須輸出：

| Field | Required content |
| --- | --- |
| 原假設 | 原本以為什麼 |
| 反證 | 哪個 evidence 推翻 |
| 新 source-of-truth | 現在依據哪些文件或證據 |
| 新 execution plan | 接下來如何修正 |
| Validation | 如何確認不再沿用錯誤 mental model |

---

## 7. Recovery Levels

| Level | Meaning | Typical action |
| --- | --- | --- |
| L1 | Simple retry | 單次工具失敗，可重試一次 |
| L2 | Reload local workflow | 重新讀當前 workflow / rule |
| L3 | Reload source-of-truth | 補讀 owner docs、UI map、contract、architecture |
| L4 | Rebuild execution graph | 停止執行，重建 goal → validation chain |
| L5 | Assumption collapse + rediscovery | 原 routing 或 task framing 可能錯，回到 discovery |

---

## 8. Domain Policy: APK Analysis First

第一個落地 domain：`apk-analysis`。

### Navigation Mismatch Policy

```yaml
apk-analysis:
  repeated_navigation_mismatch:
    escalation_level: L4
    required_reload:
      - workflow/apk-analysis/execution-flow.md
      - workflow/apk-analysis/artifact-gates.md
      - project:docs/UI架構地圖/**
      - project:docs/API列表/**
    forbidden:
      - patch_coordinates_without_ui_map
      - continue_capture_without_feature_context
      - treat_target_pid_as_ui_checkpoint_success
```

### Required APK Recovery Output

- Reset baseline：`force-stop only` / `clear cache` / `clear app data` / `reinstall`。
- Current checkpoint：screen id、route id、foreground package、feature context。
- UI evidence：screenshot / hierarchy / target package validation。
- Execution graph：navigation graph before any capture.

---

## 9. Implementation Phases

### Phase 0 — Source Consistency Audit

Goal: 先確認 runtime README、runtime source、compiler、runtime.db 是否一致。

Tasks:

- [ ] 檢查 `runtime/README.md` 宣稱的 `runtime/discovery/`、`runtime/recovery/`、`runtime/phases/`、`runtime/gates/` 是否仍是 source-of-truth。
- [ ] 檢查這些 source 是否已遷移到 `runtime/compiler/embedded_data.rb` 或只存在於 `runtime.db`。
- [ ] 決定本 plan 應恢復 YAML source，或改走 compiler embedded source。
- [ ] 若 README stale，列為 blocker，不直接新增依賴不存在路徑的 rule。

Exit criteria:

- [ ] Runtime source-of-truth 路徑被確認。
- [ ] Plan 後續 phase 的檔案目標不再依賴不存在目錄。

### Phase 1 — Enforcement Policy

Goal: 建立全庫 escalation trigger 與 forbidden behavior。

Candidate files:

- `enforcement/escalation-policy.md`
- `enforcement/README.md`
- `enforcement/failure-learning-system.md`
- `enforcement/dependency-reading.md`
- `enforcement/linked-updates.md`

Tasks:

- [ ] 新增 `escalation-policy.md`，定義 trigger classes、levels、recovery required actions。
- [ ] 在 `failure-learning-system.md` 區分 real-time escalation 與 post-mortem failure learning。
- [ ] 在 `dependency-reading.md` 加入 source-of-truth miss / reload ledger 的關聯。
- [ ] 更新 `enforcement/README.md` lazy-load 表格與完整索引。
- [ ] 更新 `linked-updates.md`，加入新增 escalation policy 的連動規則。

Exit criteria:

- [ ] 使用者否定、重複失敗、evidence conflict 能從 enforcement index 找到 recovery 要求。

### Phase 2 — Runtime Guard Integration

Goal: 把 mismatch escalation 接進 execution guard chain。

Candidate files:

- `runtime/guards/mismatch-escalation.yaml`
- `runtime/guards/circuit-breaker.yaml`
- `runtime/pipeline/guard-chain.yaml`
- `runtime/pipeline/context-flow.yaml`
- `runtime/README.md`

Tasks:

- [ ] 新增或擴充 runtime guard，表示 repeated failure / user contradiction / evidence conflict / automation drift。
- [ ] 在 execution stage 的 guard chain 中插入 `mismatch_escalation`。
- [ ] 定義 action：`warn`、`suspend_execution`、`enter_recovery`、`rediscovery_required`。
- [ ] 若 runtime source YAML 已由 compiler 管理，將 source 與 generated `runtime.db` 更新納入同一工作單。

Exit criteria:

- [ ] runtime guard chain 明確知道 mismatch 不是一般 retry，而是 recovery trigger。

### Phase 3 — Recovery Runtime

Goal: 建立 recovery mode 的 machine-readable procedure。

Candidate files:

- `runtime/recovery/mismatch-recovery.yaml`
- `runtime/recovery/source-of-truth-rebuild.yaml`
- `runtime/recovery/execution-graph-rebuild.yaml`
- `runtime/README.md`

Tasks:

- [ ] 定義 recovery state transitions：`escalation → recovery → rediscovery → replan → execution`。
- [ ] 定義 required actions：suspend assumption、reload workflow、reload owner docs、rebuild execution graph、explain failure。
- [ ] 定義 recovery output schema。
- [ ] 定義 L1-L5 escalation levels。

Exit criteria:

- [ ] Recovery mode 有 machine-readable steps，不只存在 prose plan。

### Phase 4 — Metadata Recovery Policy

Goal: 讓不同 domain 可指定 required reload set。

Candidate files:

- `metadata/recovery/README.md`
- `metadata/recovery/escalation-levels.yaml`
- `metadata/recovery/domain-policies.yaml`

Tasks:

- [ ] 建立 recovery metadata 目錄與 schema。
- [ ] 加入 `apk-analysis` policy。
- [ ] 加入 `software-delivery` policy：contract / implementation / test conflict 時 reload owner contract + BDD + executable feature。
- [ ] 決定 metadata 是否需編入 `runtime.db` 或只由 routing / validation 使用。

Exit criteria:

- [ ] Domain-specific recovery 不再硬寫在單一 enforcement rule。

### Phase 5 — Workflow Integration

Goal: 讓高風險 workflow 具備 mismatch escalation hook。

Candidate files:

- `workflow/apk-analysis/execution-flow.md`
- `workflow/software-delivery/execution-flow.md`
- `workflow/workflow-routing.md`

Tasks:

- [ ] 在 APK analysis 補 `UI/navigation mismatch escalation`。
- [ ] 在 software-delivery 補 `contract / test / implementation mismatch escalation`。
- [ ] 在 workflow routing 補多 route 或 stale route 的 recovery re-entry。

Exit criteria:

- [ ] Workflow primary source 能直接告訴 agent：事情跑偏時不能繼續 patch。

### Phase 6 — Validation Scenarios

Goal: 將本失效模式變成可測 scenario。

Candidate files:

- `validation/scenarios/failure-derived/runtime-recovery-navigation-mismatch.yaml`
- `validation/scenarios/failure-derived/runtime-recovery-user-contradiction.yaml`
- `validation/scenarios/failure-derived/runtime-recovery-source-miss.yaml`
- `scripts/validate-knowledge-runtime.rb`

Tasks:

- [ ] 建立 navigation mismatch scenario。
- [ ] 建立 user contradiction scenario。
- [ ] 建立 source miss scenario。
- [ ] 驗證 expected route / forbidden behavior / required reload set 可被 validator 檢查。

Exit criteria:

- [ ] 未來 agent 若遇到類似訊號，validation 可以檢測是否應進 recovery。

### Phase 7 — Plan Completion Closure

Goal: 完成 plan 後執行閉環。

Tasks:

- [ ] 確認所有 phase 完成或明確標 blocked。
- [ ] 執行適用 validator。
- [ ] 檢查 linked updates。
- [ ] 更新 `plans/README.md` 狀態。
- [ ] 若 plan 完成，依 `plans/README.md` 搬移至 `archived/` 或標註 active 例外。
- [ ] 完成 Ai-skill writeback close-loop。

---

## 10. Validation Matrix

| Scenario | Expected behavior | Forbidden behavior |
| --- | --- | --- |
| 同一 UI navigation checkpoint 連續失敗 2 次 | enter recovery L4，reload UI map | 繼續改座標 retry |
| 使用者指出「你沒看文件」 | suspend assumption，dependency read ledger | 口頭道歉後繼續 patch |
| API response 與 contract 衝突 | reload owner contract + API docs | 只 patch test |
| repo structure 與 README 不一致 | source consistency audit | 繼續 grep 並猜 architecture |
| workflow primary source 未讀但已開始執行 | source-of-truth miss escalation | 宣稱已按 workflow 執行 |

---

## 11. Open Questions

- `runtime/recovery/` 等目錄是否應恢復為 YAML source，或保留在 compiler embedded data？
- Recovery levels 是否應屬於 `runtime/guards/`、`metadata/recovery/`，還是兩者分工？
- 是否需要 runtime-state.db 記錄 mismatch counter？
- Validation scenario 要測 prose output shape，還是只測 route / required source selection？
- Escalation trigger 是否要接入 tool adapters，讓 IDE / CLI 在連續失敗時提示 agent？

---

## 12. Dependency Read Ledger

| Field | Content |
| --- | --- |
| Trigger | 使用者要求將 Runtime Recovery & Escalation System 寫入 Ai-skill plan，且指定要按 plan 流程。 |
| Required set | `plans/README.md`、`enforcement/README.md`、`dependency-reading.md`、`linked-updates.md`、`content-layering.md`、`reusable-guidance-boundary.md`、`failure-learning-system.md`、`feedback-lessons.md`、相關 runtime guard / pipeline 文件。 |
| Read | `plans/README.md`、`enforcement/README.md`、`dependency-reading.md`、`linked-updates.md`、`content-layering.md`、`reusable-guidance-boundary.md`、`failure-learning-system.md`、`feedback-lessons.md`、`runtime/README.md`、`runtime/guards/circuit-breaker.yaml`、`runtime/guards/context-pollution.yaml`、`runtime/pipeline/guard-chain.yaml`、`runtime/pipeline/context-flow.yaml`。 |
| Not applicable | 本 plan 是 draft，尚未新增 runtime source、enforcement rule 或 validation scenario；本輪不需要更新 runtime compiler source。 |
| Deferred / blocked | 後續實作 phase 需依各 phase 再做 linked updates、runtime compiler、validator 與 readback；本 plan 只建立 active planning entry。 |
| Validation | 檔案放入 `plans/active/`；同步更新 `plans/README.md` 目前狀態；已執行 diff check、lint 與 knowledge runtime validation。 |

---

## 13. Completion Definition

本 plan 不以「新增文件」為完成。完成條件是：

- Enforcement policy 可被 indexing / dependency reading 找到。
- Runtime guard 能表達 mismatch escalation。
- Recovery procedure 有 machine-readable steps。
- APK analysis workflow 有 navigation mismatch hook。
- 至少一個 failure-derived validation scenario 能驗證 agent 不再繼續局部 patch。
- Runtime source / README / runtime.db 一致。

