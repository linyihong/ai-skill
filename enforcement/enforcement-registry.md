# Mechanical Enforcement Registry（companion）

> **本檔為 companion markdown**。Canonical executable contract 在
> [`enforcement-registry.yaml`](enforcement-registry.yaml)。修改 rule_class
> binding 的順序：先改 yaml → `ai-skill runtime compile + refresh` →
> 本檔同步說明（companion only）。

## 為什麼存在這層

2026-05-31 session 在連續 5 個 bug 之後第三方架構評審指出：
這不是 5 個獨立 bug，而是同一個 meta-pattern：

> **Knowledge Layer 有規則，Runtime Layer 沒執行器。**

每次發生時都要等使用者半年後追問才會被抓出來。本 registry 把
「Rule ↔ Executor binding」變成 framework 第一公民，讓未來的同模式
bug 在 compile time 就被擋下。

## Layer 2.5 — Meta Governance / Framework Self-Audit

Ai-skill 原本是三層架構：

```
Layer 1  Knowledge       (enforcement/, governance/, workflow/, ...)
Layer 2  Runtime         (scripts/ai-skill-cli/, hooks.go, runtime.db)
Layer 3  Governance      (constitution/, architecture/, plans/)
```

本 registry 建立缺失的 Layer 2.5：

```
Layer 2.5  Coverage Verification / Meta Governance
             Rule ←binding→ Executor ←evidence→ Verification
             的結構性驗證層
```

- Layer 1-3 管「框架的內容」（規則寫了什麼、runtime 跑什麼、治理決策）
- Layer 2.5 管「框架本身是否真的做到它說會做的」（Governance of Governance）

沒有 Layer 2.5 就沒有「rule 寫好、executor 沒接」「executor 寫了但漏覆蓋
某些 instance」「scenario 過但 production 沒人用」的結構性偵測機制。

## Rule Class 而非 Rule Instance

第七輪評審指出原始設計把 150+ rule instance 個別 binding 是維護地獄。
本 registry 用 **rule_class**（~20-30 個）抽象，每個 class 含多個 instance
但共用 executor binding。例如：

- `cognitive_mode_governance` 一條 entry 涵蓋 9 個 commit-msg validator
- `plan_governance` 一條 entry 涵蓋 plan_status_sync + checkbox_sync + archival_audit

這讓 registry 規模可控（目前 28 個 class）。

## 6-Value Coverage Enum

每個 rule_class 必須宣告自己屬於哪一種：

| Coverage | 語意 | 必填 metadata |
|---|---|---|
| `mechanical` | Executor 已存在且 enforcing | `executors[]` + `rationale` |
| `behavioral_only` | 故意不機械化，靠 agent 行為 | `rationale` + `sunset_decision.{revisit_when, success_criteria}` |
| `not_mechanizable` | 永遠不該機械化（主觀 / 無客觀 validation） | `rationale` + `objective_validation_impossible_because` |
| `pending_implementation` | 知道怎麼做、實作中或排程中 | `child_plan` + `target_promotion` |
| `research_required` | 知道應該機械化、但還不知道怎麼 | `rationale` + `research_questions` + `estimated_unblock_timeline` |
| `deprecated` | 移除中（已被取代或排程移除） | `replaced_by` OR `removal_date` |

詳細 schema / lint 行為見 [`enforcement-registry.yaml`](enforcement-registry.yaml)
§`coverage_status_spec`。

### 為什麼 `behavioral_only` 必須雙必填

`sunset_decision` 需要 `revisit_when` **且** `success_criteria` 兩個欄位
都填，原因（2026-05-31 session Q2 修正）：

| 只填 `success_criteria` | 只填 `revisit_when` | 兩者都填 |
|---|---|---|
| 有標準但永遠沒人檢查 | 沒標準但會被檢查 | 事件 trigger + 客觀判定 |
| → 假象的安全感 | → 至少會被重新看到 | → 治理閉環成立 |
| **最危險** | 可接受 | 必填 |

「永久 behavioral_only」是 Layer 2.5 最大失效模式 —— 規則寫了「暫時不
機械化」之後，沒有人回來檢查條件是否成立。雙必填讓 sunset 從「希望」
變成「事件驅動的義務」。

### 為什麼 `pending` 拆兩種

v4 把 v3 的 `pending` 拆成 `pending_implementation`（知道怎麼做）vs
`research_required`（還不知道怎麼做）。同樣是「未完成」但治理訊號完全不同：

- `pending_implementation` 訊號是「快完成」→ 看 child plan 進度
- `research_required` 訊號是「需要思考」→ 看 research_questions 有沒有解決

混在同一桶會讓 governance dashboard 失去解析度。

## 4-Level Verification Ladder

Coverage 講「我們選擇怎麼處理這條規則」；Verification 講「實作真的做到了嗎」。
兩者正交，但 verification 自己是階梯：

```
symbol_exists → scenario_exists → regression_exists → runtime_observed
   越往右越接近 production reality
```

特別重要的是 `runtime_observed`（v4 NEW）：scenario 100% 涵蓋不等於
production 真的有跑。一個 route 半年沒被觸發，要嘛沒人用（候選 deprecate），
要嘛 detector 安靜壞了。Runtime metrics 收集會 surface 這層 reality gap。

## Executor Kind 與 Internal Helper 邊界

不是每個 Go function 都要進 registry。Q4 resolution 定義白名單：

| Kind | 是否需要 binding |
|---|---|
| `hook_dispatcher_entry`（runSessionStart 等） | ✓ 必須 |
| `commit_msg_validator`（validateXxx 系列） | ✓ 必須 |
| `runtime_state_machine_phase`（runtime.db phases） | ✓ 必須 |
| `internal_helper`（parseYaml、normalizePath 等） | ✗ 豁免 |

`internal_helper` 用 yaml 顯式 allowlist 維護
（`enforcement-registry.yaml` §`internal_helper_allowlist`），不在 Go code
散落 annotation —— 邊界放在治理可見的地方。

## Compile-time Lint（Phase 3 將實作）

第一次 lint 跑出 orphan rule / orphan executor 採 **hard block，無 grace
period**（Q1 resolution）。理由：

- Grace period 會回到「warning → 先放著 → 半年後還在」失效模式
- 違背 Prevent > Detect > Repair 哲學
- 第一次 land 預期需要密集 backfill，但這是 one-time cost

`enforcement_mode: { orphan_rule: fail, orphan_executor: fail }`。

## Registry Self-Governance（Phase 4.5）

> **狀態**：Phase 4.5 land 2026-06-02。R1/R2/R3 由 commit-msg validator
> `validateEnforcementRegistryTransition`（obligation
> `obligation.commit.enforcement_registry_transition`）機械強制；R4/R5 在
> `ai-skill enforcement coverage` 的 `## Governance Alerts` 段落 surface，
> 不阻塞 build（governance review trigger，不是 compile fail）。

Layer 2.5 自己也需要治理。沒有這層，registry 變成「一個沒人管的元數據檔」。
Self-governance 把「改 registry 一行 yaml」從 silent edit 升級為留下 ADR /
trailer / verification evidence 軌跡的治理動作。

### Status Transition Matrix

| From → To | Required action | 強制層 |
|---|---|---|
| `(new)` → `pending_implementation` | 引用 active child plan | Phase 3 lint `pending_implementation_child_plan_validity` |
| `(new)` → `research_required` | 列 `research_questions` ≥ 1 + estimated_unblock | schema lint |
| `pending_implementation` → `mechanical` | executor symbol live + coverage_evidence + verification thresholds | **Phase 4.5 R3** commit-msg |
| `research_required` → `pending_implementation` | research_questions 全 resolved + child plan | schema lint |
| `mechanical` → `behavioral_only` | **demotion，需 ADR** | **Phase 4.5 R2** commit-msg |
| `mechanical` → `not_mechanizable` | **demotion，需 ADR** | **Phase 4.5 R2** commit-msg |
| `mechanical` → `deprecated` | `replaced_by` 指向 active mechanical class | Phase 3 lint `deprecated_disposal` |
| `behavioral_only` → `not_mechanizable` | **demotion，需 ADR** | **Phase 4.5 R2** commit-msg |
| `deprecated` `removal_date` 屆期 | governance 決定 actually remove vs extend | Phase 4.5 R4 coverage alert |

### Self-Governance Lint Rules

| Code | Layer | Severity | 行為 |
|---|---|---|---|
| **R1** | commit-msg | block | status 變更 commit 必須有 `[registry-status-change]` trailer **and** `rationale: <text>` 行 |
| **R2** | commit-msg | block | demotion 必須在 rule_class entry 加 `adr_reference: constitution/ADR-NNN-*.md`，且 ADR 檔案必須存在 |
| **R3** | commit-msg | block | promotion to mechanical 觸發 Phase 3 `missing_executor_symbol` 子集 lint 於該 class；symbol 不存在於 declared file 則 reject |
| **R4** | coverage report | governance alert | deprecated 過 `removal_date` ≥ 30 天 → `## Governance Alerts` 段標紅；不阻塞 build |
| **R5** | coverage report | governance alert | research_required 過 `estimated_unblock_timeline` → `## Governance Alerts` 段標紅；不阻塞 build |

**Opt-out**：`[skip-registry-transition]` trailer 跳過 R1/R2/R3（與其他 11
個 commit-msg validator 相同的 opt-out 慣例；緊急修補時可用，但會留下
commit-msg 軌跡）。

### 為什麼 demotion 必須附 ADR（R2）

| Without R2 | With R2 |
|---|---|
| 任何 dev 改一行 yaml 把 mechanical 改 behavioral_only | 必須先寫 `constitution/ADR-NNN-<slug>.md` 解釋為什麼放棄機械化 |
| Silent demotion，半年後沒人記得為何降級 | 永久 governance 軌跡，可追溯 decision context |
| 看 coverage report 以為「7 個 behavioral_only」但每個來源動機不明 | 每個降級都有 supersede 條款 + 邊界條件 |

R2 的設計意圖不是阻擋降級，而是讓「我們決定不機械化某條規則」變成
governance decision，不是 commit-time afterthought。

### 為什麼 promotion 必須過 verification（R3）

Promotion to mechanical 是 one-way ratchet：一旦 coverage=mechanical
land，下游消費者（coverage report、governance dashboards、Phase 5 bootstrap
Receipt 摘要）就會 trust 這個 claim。若實際 executor 還沒寫好，等於
registry 自己誤導下游。

R3 的最小門檻是 `verification_levels.symbol_exists`（Phase 3 lint
`missing_executor_symbol` 的 per-class scope 版本）。其他 verification
layer（scenario_exists / regression_exists / runtime_observed）由 Phase 3
compile-time lint emit WARNING 或 Phase 5 coverage report surface，
但不是 R3 的硬門檻 —— promotion 只需保證「symbol 至少真的存在」。

### 開發者快速指南

| 場景 | 必須做的 |
|---|---|
| Demote mechanical → behavioral_only | 1) 寫 ADR；2) 加 `adr_reference` 欄位；3) commit body 加 `[registry-status-change]` + `rationale:` 行 |
| Promote pending → mechanical | 1) 先 land executor symbol（hooks.go 等）；2) 改 coverage=mechanical；3) commit body 加 trailer + rationale |
| 新增 rule_class | 不觸發 R2/R3（沒有舊狀態可 demote/promote），但若 coverage 初始為 mechanical 也必須通過 symbol_exists 檢查 |
| 緊急熱修 | `[skip-registry-transition]` opt-out trailer，但需在下個 commit 補 ADR / verification |
| 想看現況 | `ai-skill enforcement coverage --format text --detail` 看 6-bucket + 每 class verification level + `## Governance Alerts` 段（R4/R5） |
| 模擬 transition | `ai-skill enforcement transition-check --old <old.yaml> --new <new.yaml> --commit-msg-file <msg.txt>` 跑與 commit-msg 相同的 R1/R2/R3 engine（scenario / CI / local debug 用） |

## 寫作指南（給 rule_class 作者）

1. **先決定 coverage**：mechanical / behavioral_only / not_mechanizable /
   pending_implementation / research_required / deprecated
2. **填對應必填 metadata**（見上方表格）
3. **source_files 列 canonical source**（yaml 優於 md；兩者都列也可）
4. **若是 mechanical**：在 hooks.go 或 runtime state machine 中指向真實
   symbol；驗證 lint 會校驗 symbol 存在
5. **若是 behavioral_only**：認真寫 `revisit_when`（事件不是日期）+
   `success_criteria`（可觀察條件）；推薦補 `revisit_owner`
6. **若是 pending_implementation**：`child_plan` 必須是
   `plans/active/*.md` 路徑
7. **若是 research_required**：`research_questions` 寫具體疑問，不是
   「未來考慮」
8. **rule_class 數量上限**：soft target 24、hard limit 40；目前 28
   在 soft-hard 區間內

### 反例

```yaml
# ❌ behavioral_only 缺 revisit_when
- id: foo
  coverage: behavioral_only
  rationale: "暫時不機械化"
  sunset_decision:
    success_criteria: "未來會處理"   # 空話 + 缺事件 trigger
# lint fail: missing sunset_decision.revisit_when
```

```yaml
# ❌ research_required 寫成 "未來考慮"
- id: bar
  coverage: research_required
  research_questions:
    - "之後再說"   # 不是具體疑問
  estimated_unblock_timeline: "未定"
# lint fail: research_questions must list ≥ 1 concrete unresolved question
```

```yaml
# ✓ 良好範例
- id: capability_discovery
  coverage: behavioral_only
  rationale: |
    Discovery 是 detector miss fallback，per-turn 強制成本爆炸。
  sunset_decision:
    revisit_when: "workflow_activation child plan Phase 6.1 lands"
    success_criteria: |
      Detector miss path can invoke Discovery and produce
      route_candidate_proposals.yaml.
    revisit_owner: "framework maintainer"
```

## 與其他 enforcement rules 的關係

- [`dependency-reading.md`](dependency-reading.md) — 修改 enforcement rule
  時必讀；本 registry 是 cross-cutting 索引，不取代各 rule yaml
- [`linked-updates.yaml`](linked-updates.yaml) — registry 變更 commit
  須完成 writeback transaction
- [`rule-weight.md`](rule-weight.md) — registry 自己屬 P1（canonical
  repository writeback），lint 違規 = compile fail

## 與 `runtime/core-bootstrap.yaml` 的關係（Q3 resolution）

- `core-bootstrap.yaml` 是 **phase-aware obligation lifecycle**（何時 fire）
- `enforcement-registry.yaml` 是 **cross-phase binding view**（何處 enforce）
- 兩者互補，不重複

## Round-4/5 Schema 結構性退階（2026-06-01）

Round-4 評審指出 `behavioral_only` 從原本的「輕量例外」逐輪累積到 7 個
metadata 欄位，governance cost 逼近 mechanical。Round-5 收尾 3 個規格
空洞後 frozen 為 final baseline。

### `behavioral_only` 分為 hard required + recommended

**Hard required（缺即 FAIL）**：

- `rationale`
- `sunset_decision.revisit_when`
- `sunset_decision.success_criteria`

**Recommended（缺即 WARNING，不阻塞）**：

- `sunset_decision.revisit_owner`
- `sunset_decision.last_reviewed_at`
- `sunset_decision.last_review_summary`
- `sunset_decision.depends_on_rule_classes`（結構化引用，取代 NLP 解析 free-text revisit_when）

詳見 yaml `coverage_status_spec.behavioral_only.lint_behavior`。

### `pending_implementation.child_plan_validity`（U2）

`pending_implementation` 的 `child_plan` 必須滿足 4 條規則：
(a) 路徑 resolve 到 `plans/active/*.md`（**fail** 若違反）
(b) 包含 `## Phase 0` heading（**warning** 若違反）
(c) 有 owner 標示（**warning**）
(d) 有 `## Validation Plan` 或 `## Acceptance` 區塊（**warning**）

Stub plans 合法只要滿足 (a)-(d)。新 lint
`pending_implementation_child_plan_validity` 強制此 schema。

### `upstream_classes` scope freeze（ADR-010）

`rule_classes[].upstream_classes: []` 僅用於 promotion traceability /
cycle prevention / coverage report visualization。**不**用於 execution
ordering / dependency injection / runtime orchestration / DAG scheduling。

跨越此 boundary 的新欄位（`downstream_classes` / `promotion_role` /
`artifact_type` 等）必須先寫新 ADR 顯式 supersede ADR-010。詳見
[`constitution/ADR-010-registry-upstream-classes-scope-freeze.md`](../constitution/ADR-010-registry-upstream-classes-scope-freeze.md)
+ yaml `upstream_classes_scope` 區塊。

### `bootstrap_mode` 與 `baseline_snapshot`

`bootstrap_mode: strict`（預設）—— orphan_rule / orphan_executor 從
day 1 hard-fail。本 Ai-skill repo 採此模式。

`bootstrap_mode: baseline_snapshot_v1` —— 成熟 repo 第一次導入用：
記錄當下 orphan 為 baseline（lint 降 warning），新增 orphan 仍 fail。
baseline 必須含 `baseline_owner` + 每筆 entry 的 `baseline_review_summary`
（Round-4 T4 治理對等）。

## Phases（後續工作）

| Phase | 範圍 | 狀態 |
|---|---|---|
| Phase 0-2 | Open Questions + rule_class 識別 + 初版 registry | ✓ land |
| Phase 3 Round-1~5 | Compile-time lint design baseline frozen | ✓ land (round-5 = final) |
| Phase 3 Step 0-1 | ADR-010 + F19 stub plan + schema patch | ✓ land 2026-06-01 |
| Phase 3 Step 2 | Lint patch（13 個 check, P0 7 / warning 6） | in progress |
| Phase 3 Step 3-10 | Backfill / wire / rebuild / commit | pending |
| Phase 4 | CLI `ai-skill enforcement coverage` | pending |
| Phase 4.5 | Self-governance lint R1-R9 + `validateEnforcementRegistryTransition` | pending (R6-R9 added round-4/5) |
| Phase 5 | Bootstrap integration（Receipt 加 coverage 摘要） | pending |
| Phase 5.x | Hook injection economics（inaugural dogfooding case） | pending |
| Phase 6 | `enforcement/failure-patterns/rule-without-executor.md` | ✓ land |
| Phase 7 | Validation scenarios（5 個 regression case） | ✓ land |
| Phase 8 | Close-out + archive | pending |

詳見
[`plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md`](../plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md)。

← [Back to enforcement index](README.md)
