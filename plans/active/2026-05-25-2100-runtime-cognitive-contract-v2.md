# Runtime Cognitive Contract v2

**Status**: `draft`
**世代**：Gen 3 子系統演進（ADR-008 amendment 或 superseding，pending Phase 7 評估）
**建立日期**：2026-05-25
**最後更新**：2026-05-25

> 本 plan 回應 2026-05-25 外部反饋：「Cognitive Mode 報告」目前是 debug telemetry，還不是穩定的 cognitive contract。Mode 是 label without capability semantics，agent self-describes 容易 inflated reporting，verbosity inflation 造成 cognitive fatigue。要把現在的 mode block 升級為 **Runtime Cognitive Contract**，加入 validation_mode + cognitive_cost、compact/full adaptive form、activation_reason 必須引 discovery signals。

---

## Decision Rationale

### Problem & Why Now

ADR-008 上線後一個 session（2026-05-25）內，agent 已在 14 commits 寫 Cognitive Mode block。實作後出現 4 個系統性問題：

1. **Mode 是 abstract label，不是 capability contract**
   - `STRICT` / `DEEP` 等值在 report 裡看不出 = 哪幾條 enforcement。
   - 雖然 capability 已在 `runtime/cognitive-modes-{phase,governance,memory}-integration.yaml` 定義，但 report 沒 surface。
   - Reader 必須回查 YAML 才知道 mode 真正 mean 什麼。

2. **Verbosity inflation**
   - 每 turn 都列 4 列表格 → 對 trivial task 是 noise；對 reader 是 cognitive fatigue。
   - 本 session 12+ turns 各帶 ~5 行 block ≈ 600 重複 tokens / session。
   - 高頻低資訊 → mode 失去 signal value，淪為儀式。

3. **Self-describing drift（inflated cognitive reporting）**
   - Agent 自己宣告 mode，沒有與 raw signals 對照 → confidence inflation 的 runtime 版本。
   - 例：實際只改了一個 typo，agent 宣告 `execution_mode=DEEP / governance_mode=STRICT`，無從反駁。
   - Discovery 已有 14 signals（`runtime/cognitive-modes-discovery.yaml`），但 commit-msg validator 沒檢查 declared mode 與 signals 的一致性。

4. **Token Estimate ≠ Cognitive Cost**
   - Token 是 observable 量；cognitive cost 是抽象負荷。
   - 例：大量 source read = 高 token 低 cognition；architecture tradeoff = 低 token 高 cognition。
   - 目前 Token Estimate 是 agent 自報，又是 self-describing 風險。

**Why now**：4 個問題同根（label-without-contract + self-describing + uniform verbosity），incremental fix 容易留 inconsistent state。趁 ADR-008 還只 1 session 經驗，整批收斂代價最低。

### Decision

引入 **Runtime Cognitive Contract v2**，6 維 + adaptive disclosure：

#### 1. 6 維 cognitive contract（v1 4 維 + 2 新增）

| 維度 | 值域 | 性質 |
|---|---|---|
| `execution_mode` | FAST / NORMAL / DEEP / FORENSIC / RECOVERY | unchanged from v1 |
| `context_mode` | INDEX_ONLY / SUMMARY_FIRST / CHECKLIST_FIRST / SOURCE_BACKED / GRAPH_ASSISTED | unchanged |
| `governance_mode` | LIGHT / STANDARD / STRICT / LOCKDOWN | unchanged |
| `memory_mode` | NONE / EPISODIC / DECISION_REPLAY / FAILURE_REPLAY / PROJECT_CONTEXT | unchanged |
| **`validation_mode`** | SKIP / CHECKLIST / SOURCE_BACKED / GRAPH_TRACED | **新增** — validation depth 與 execution mode 解耦 |
| **`cognitive_cost`** | LOW / MEDIUM / HIGH / VERY_HIGH | **新增** — derived from (execution × context) lookup，不可手寫 |

#### 2. Adaptive disclosure（compact / full / capability）

| 任務性質 | Form |
|---|---|
| 全 default value（NORMAL/SUMMARY_FIRST/STANDARD/NONE/CHECKLIST + LOW cost） | **Single-line compact**：`Cognitive: NORMAL·SUMMARY_FIRST·STANDARD·NONE / V:CHECKLIST / Cost:LOW / Sig:typo_fix` |
| 任一維度非 default | **Full table + activation_reason** |
| 高風險（DEEP/FORENSIC/RECOVERY/LOCKDOWN 任一） | **Full table + activation_reason + capability snippet**（從 integration YAMLs 取 1-2 行） |

升級門檻 = mode 本身的 floor。

#### 3. `activation_reason` 必填、必引 discovery signals

每 report 必含 `activation_reason` 欄位列出 ≥1 signal，且 signal 必須是 `runtime/cognitive-modes-discovery.yaml` 已定義 14 signals 之一。Validator 檢查 declared signal name ∈ known signals。Unknown signal name → block（inflated reporting）。

例：
```
activation_reason:
  - file_diff_runtime_schema      ← from discovery YAML
  - cross_layer_rule_modification ← from discovery YAML
```

#### 4. `cognitive_cost` 由 lookup table derived，不是 agent claimed

新 YAML `runtime/cognitive-modes-cost-class.yaml`：

| execution × context | cost |
|---|---|
| FAST × INDEX_ONLY | LOW |
| FAST × any other | MEDIUM |
| NORMAL × {INDEX_ONLY, SUMMARY_FIRST} | LOW |
| NORMAL × {CHECKLIST_FIRST, SOURCE_BACKED} | MEDIUM |
| DEEP × any | HIGH |
| FORENSIC × any | VERY_HIGH |
| RECOVERY × any | VERY_HIGH |

Validator 從 declared execution + context 自動算 cost，與 declared cost 對照；不一致 → block。

#### 5. Capability snippet（high-risk 才展開）

當 governance_mode ∈ {STRICT, LOCKDOWN} 或 execution_mode ∈ {DEEP, FORENSIC, RECOVERY} 時，report 末尾自動含 capability summary（由 hook 從 integration YAMLs 摘要 1-2 行）：

```
Capability summary:
  governance_mode=STRICT → source-backed required, validation_before_patch, no_global_claims
  execution_mode=DEEP → source-backed reads + dependency graph + linked updates fully resolved
```

### Alternatives Considered

- **A. 維持 v1（do nothing）**：reject — 4 問題會持續放大；每多累積 task 就多一筆 inflated reporting evidence
- **B. 只做 compact form（解 verbosity 一個問題）**：reject — 不解決 label-without-contract / self-describing；改一次只解 1/4 後又要再改
- **C. 同 batch 做 v2**（current decision）：accept — 4 問題同根，一次收斂降低 inconsistent-state 風險
- **D. 廢掉 mode block 完全靠 discovery signals 自動 attach**：reject — 目前沒有 per-tool-call interception infra；commit-msg 是現有 enforcement point，agent 還是要寫 report；自動 attach 要 Phase 5+ adaptive runtime
- **E. 直接寫 ADR-009 取代 ADR-008**：reject — v2 是 ADR-008 的演進不是廢除；先以 plan 完成驗證，Phase 7 評估是 amend ADR-008 還是 promote ADR-009

### Why Not an ADR Yet

- v2 設計仍未 implement，4 個改動互相依賴，可能在 Phase 2-5 發現需要調整 scope
- ADR-008 才 1 session 經驗，v2 是 amendment 性質；多 session 累積後才看得出 ADR-009 或 amend
- Open Questions 未全解

### ADR Promotion Criteria（completed 時驗證）

- [ ] foundational + cross-session + cross-project + expensive-to-reverse + explains-why 全中
- [ ] Plan 結果證實 6 維 + adaptive disclosure 可行
- [ ] Open Questions 全解
- [ ] 沒有更輕的 promotion target 適用
- [ ] 系統真實使用此 contract，≥5 commits 用 v2 form（含 ≥2 high-risk + ≥2 compact）
- [ ] Verbosity inflation 量化下降（v2 commits 平均 block 行數 < v1 commits）

評估時機：Phase 7 close-loop。決議 = (a) amend ADR-008 with v2 section / (b) supersede with ADR-009 / (c) keep as plan only (lighter promotion)。

### Consequences（預期）

#### 正面

- **Capability semantics**：mode 不再只是 label，high-risk 自帶 capability summary
- **Reduced cognitive fatigue**：trivial commits 一行；只在真實高風險才 full disclosure
- **Inflated reporting blocked**：activation_reason 必須引 known signals，hook 機械驗證
- **Cognitive cost observability**：cost class derived，與 token estimate 分離但並存
- **Human ↔ runtime alignment**：user 可從 capability snippet 直接 challenge mode 選擇

#### 負面

- **Migration cost**：所有現有 commits 的 block format 與 v2 不同（不需要 retroactive 改，但未來 commits 不一致期）
- **Discovery YAML 維護負擔**：activation_reason 必須引 known signals → 新 signal 加入需更新 YAML
- **Hook 邏輯複雜度**：validator 從單一 format check → 兩種 form + signal name check + cost derive + capability surface

#### 風險

- **Compact form 被濫用**：agent 為了少寫一律報 compact → mitigation：default-only triggers compact，任一非 default → 強制 full
- **Signal vocabulary drift**：discovery YAML 與 validator 各自演進 → mitigation：validator 從 generated_surfaces 讀 signal list 而非 hardcode
- **Cost class 不準**：execution × context lookup 可能 over-simplify → mitigation：Phase 5 可加入 governance / memory 維度的影響

---

## Runtime Execution Path

| 欄位 | 內容 |
|---|---|
| Runtime owner | `runtime/cognitive-modes.yaml`（v2 amendment）+ `runtime/cognitive-modes-cost-class.yaml`（新）+ `scripts/ai-skill-cli/internal/app/hooks.go`（validator 升級） |
| Trigger flow | commit-msg hook → parse Cognitive Contract block（compact 或 full）→ verify `activation_reason` signals ∈ known discovery signals → derive `cognitive_cost` from `(execution × context)` lookup → compare with declared cost → if high-risk mode, verify capability snippet present → emit block message with rich error if any check fails |
| Trigger location | `runCommitMsgHook` in `scripts/ai-skill-cli/internal/app/hooks.go` |
| Activation contract | `runtime/cognitive-modes-cost-class.yaml`（projected to `generated_surfaces[runtime.cognitive_modes.cost_class]`）+ updated `runtime/cognitive-modes.yaml` v2 section |
| Generated surface | `runtime.cognitive_modes.cost_class`, `runtime.cognitive_modes.contract` (v2 amend) |
| Validation scenarios | `phase6-cognitive-contract-v2-{compact-form, full-form, activation-signal, cost-class, capability-snippet, inflated-rejection}-v1.yaml` (6 scenarios) |
| Test passing evidence | 全部 scenarios PASS + 6+ unit tests + ≥5 commits using v2 form + Phase 7 verbosity comparison report |

---

## Open Questions

| # | Question | Status |
|---|---|---|
| 1 | Compact form trigger boundary：「全 default」如何精確定義？是否包括 cognitive_cost=LOW 或 validation_mode=CHECKLIST? | TBD Phase 2 |
| 2 | Capability snippet 從 integration YAMLs 動態取 vs hook 內 hardcode lookup? | TBD Phase 5 |
| 3 | `activation_reason` 可否為 user-provided override（user 指定 signal name 而非 agent infer）? | TBD Phase 4 |
| 4 | Existing v1 commits 不 retroactive 改 — 是否要加 deprecation timeline? | TBD Phase 6 |
| 5 | Phase 7 ADR 決議 amend ADR-008 還是 ADR-009 supersede? | TBD Phase 7 |

---

## 完成條件

### 計畫書本身

- [ ] 計畫書狀態：`draft` → `in-progress`（Phase 0 通過後）→ `completed`
- [ ] 5 Open Questions 全部 resolved（在 Phase 0-2 補充）
- [ ] Phase 7 close-loop 完成 ADR amend / supersede / no-ADR 決議

### Behavioral evidence

- [ ] ≥5 commits 在 Phase 1-7 期間使用 v2 form（含 ≥2 high-risk full + ≥2 compact）
- [ ] Verbosity inflation 量化：v2 commits 平均 block 行數 < v1 commits 平均
- [ ] 至少 1 個 inflated-reporting attempt 被 validator 擋下並修正

### Validation

- [ ] 所有新 scenarios PASS
- [ ] 既有 21 scenarios 仍 PASS（不 regression v1 contract）

---

## Phase 0 Pre-Build Interrogation

### 目的

驗證 v2 與 ADR-008 既有實作相容，確認 6 個 integration YAMLs（discovery / phase / governance / memory / token-budget / adaptive）不需要 schema breaking change。

### Tasks

- [ ] 讀 `runtime/cognitive-modes-*.yaml` 全部 7 個 contracts 確認 v2 6 維可以兼容 layered 進去
- [ ] 確認 `cognitive-modes-discovery.yaml` 14 signals 的 name 可用作 `activation_reason` 引用（命名規範）
- [ ] 驗證 `commit-msg` hook 現有 6 個 validators 中 `validateCognitiveModeBlock` 是 v2 主要修改點，其餘 validators (`validateExecutionModeFloors` 等) 是否需配合改
- [ ] 確認 `models/cognitive-modes/README.md` 為 source-of-truth 文件，需要更新 mode 報告範本至 v2

### Phase 0 完成條件

- [ ] 上面 4 tasks all done
- [ ] No breaking-change conflict surfaced
- [ ] 若有 conflict，更新 plan §Decision 或新增 Open Question

---

## Phase 1 Test-First Validation

### 期望可觀察行為

寫失敗 scenarios 表達 v2 對外行為。預期 Phase 2-6 implementation 後逐一通過。

### Tasks

- [ ] `phase6-cognitive-contract-v2-compact-form-v1.yaml` — trivial commit 用 compact form，hook 通過
- [ ] `phase6-cognitive-contract-v2-full-form-v1.yaml` — non-default mode 用 full form，hook 通過
- [ ] `phase6-cognitive-contract-v2-activation-signal-v1.yaml` — declared signal ∈ known 14 signals → PASS；unknown signal → BLOCK
- [ ] `phase6-cognitive-contract-v2-cost-class-v1.yaml` — declared cost vs derived cost 不一致 → BLOCK
- [ ] `phase6-cognitive-contract-v2-capability-snippet-v1.yaml` — high-risk mode 缺 capability snippet → BLOCK
- [ ] `phase6-cognitive-contract-v2-inflated-rejection-v1.yaml` — typo commit 報 DEEP → cost mismatch → BLOCK

### Phase 1 完成條件

- [ ] 6 scenarios 全部寫好且 initial state = FAIL（pre-implementation）
- [ ] Scenarios commit 為一個 atomic test-first commit（[skip-cognitive-mode] / [skip-token-budget] opt-out 因為這個 commit 是 meta-scenario，不是真實 task）

---

## Phase 2 Compact / Full Form Specification

### Tasks

- [ ] 更新 `runtime/cognitive-modes.yaml` 加入 v2 section：6 維定義 + compact form spec + full form spec + capability snippet trigger rule
- [ ] 更新 `CORE_BOOTSTRAP.md` Cognitive Mode block 範本：列出 compact form + full form + 何時用哪種
- [ ] 更新 `models/cognitive-modes/README.md` 同步範本
- [ ] 確認 commit-msg hook `validateCognitiveModeBlock` 接受兩種 form

### Phase 2 完成條件

- [ ] YAML 更新通過 `runtime validate`
- [ ] CORE_BOOTSTRAP / models 文件 updated
- [ ] Scenario `phase6-cognitive-contract-v2-compact-form-v1` PASS

---

## Phase 3 Validation Mode + Cognitive Cost Dimensions

### Tasks

- [ ] 新 YAML `runtime/cognitive-modes-cost-class.yaml`（execution × context → cost lookup）
- [ ] 新 YAML or 加進 `runtime/cognitive-modes.yaml`：validation_mode 5 個值 + 與 execution_mode 的關係
- [ ] 更新 hooks.go：`deriveCognitiveCost(execution, context)` helper；`validateCognitiveCost` validator 比對 declared vs derived
- [ ] 6 unit tests cover cost lookup + validation_mode parse

### Phase 3 完成條件

- [ ] 兩個新 YAML 通過 validate
- [ ] Unit tests PASS
- [ ] Scenario `phase6-cognitive-contract-v2-cost-class-v1` PASS

---

## Phase 4 Activation Reason Enforcement

### Tasks

- [ ] 更新 commit-msg hook：parse `activation_reason` 區塊 → 提取 signal names → 與 `runtime/cognitive-modes-discovery.yaml` 14 signals 對照 → unknown name 列入 violation
- [ ] 從 `generated_surfaces` 讀 known signal list（而非 hardcode）→ signal vocabulary drift 自動跟上 discovery YAML
- [ ] Unit tests cover：valid signal → PASS / unknown signal → BLOCK / 空 activation_reason on non-trivial mode → BLOCK

### Phase 4 完成條件

- [ ] Hook 從 generated_surfaces 讀 signals 成功
- [ ] Scenario `phase6-cognitive-contract-v2-activation-signal-v1` + `inflated-rejection-v1` PASS

---

## Phase 5 Capability Snippet Surfacing

### Tasks

- [ ] 設計 capability snippet generator：當 mode tuple 含 high-risk value → 從 integration YAMLs 取 1-2 行 capability summary
- [ ] Open Question 2 決議：動態讀 YAML vs hardcode lookup（建議從 generated_surfaces 動態讀，跟上 contract evolution）
- [ ] Hook validator 檢查 high-risk commit 含 capability snippet section
- [ ] Unit tests

### Phase 5 完成條件

- [ ] Scenario `phase6-cognitive-contract-v2-capability-snippet-v1` PASS
- [ ] high-risk commit 自動有 capability summary 內容

---

## Phase 6 Migration And Deprecation

### Tasks

- [ ] 新增 `enforcement/failure-patterns/inflated-cognitive-mode-reporting.md` 描述 self-describing drift
- [ ] 更新 `enforcement/failure-patterns/cognitive-mode-resolution-bypass.md` cross-ref 至 v2 contract
- [ ] 既有 v1 commits 不 retroactive 改；future commits 走 v2
- [ ] Open Question 4 決議：是否要在 hook 加 deprecation window（n 個月後拒絕 v1 form?）
- [ ] 更新 `CLAUDE.md` 末尾的 modification rule 提及 v2

### Phase 6 完成條件

- [ ] Failure pattern 新增 + 既有 cross-ref 完成
- [ ] Deprecation policy decided

---

## Phase 7 Close-Loop & ADR Decision

### Tasks

- [ ] 全部 21 + 6 = 27 scenarios PASS
- [ ] 累積 evidence：≥5 v2 commits（≥2 high-risk full + ≥2 compact + ≥1 inflated-rejection 被擋）
- [ ] 量化 verbosity inflation：v1 commits 平均行數 vs v2 commits 平均行數
- [ ] 評估 ADR Promotion Criteria：(a) amend ADR-008 / (b) supersede with ADR-009 / (c) keep as plan only
- [ ] 若 (a) → 更新 ADR-008 加 v2 section
- [ ] 若 (b) → 寫 ADR-009、ADR-008 標 superseded
- [ ] Plan status → completed，移到 plans/archived/
- [ ] Plan Completion Closure 走完

### Phase 7 完成條件

- [ ] 27 scenarios PASS
- [ ] ADR decision recorded
- [ ] Plan archived

---

## Stakeholder 同意項目

- [ ] User confirms: v2 6 維（含新增 validation_mode、cognitive_cost）值得 migration cost
- [ ] User confirms: compact form 觸發條件（全 default）
- [ ] User confirms: activation_reason 必引 known signals 是 strict enforcement（不接受 free-form）
- [ ] User confirms: cognitive_cost derived (non-claimed) 設計
- [ ] User confirms: Phase 7 評估 amend vs supersede ADR

---

## 與其他 plans 的關係

| Plan | 關係 |
|---|---|
| [`archived/2026-05-22-1629-runtime-cognitive-modes-system.md`](../archived/2026-05-22-1629-runtime-cognitive-modes-system.md) | v1 的 parent plan；本 plan 是 v1 amendment / evolution |
| [`active/2026-05-25-1000-context-language-glossary-system.md`](2026-05-25-1000-context-language-glossary-system.md) | independent；ubiquitous language 是不同主題（cognitive ≠ vocabulary） |
| [`constitution/ADR-008-runtime-cognitive-modes.md`](../../constitution/ADR-008-runtime-cognitive-modes.md) | Phase 7 可能 amend 或 supersede |
