---
id: 2026-06-29-1430-preparatory-refactoring-workflow
plan_kind: main
status: draft
owner: linyihong
created: 2026-06-29
priority: P1
required_for_completion: false
---

# Software-Delivery Implementation Execution Mode — Structure Preparation

**Status**: `draft` — **觀察期**

**Execution Governance snapshot**（stakeholder 2026-06-29）：

| Field | Value |
|-------|-------|
| **Status** | Observed Across Dual Paths |
| **Confidence** | stop path → **verified** (`01`); happy path → **partial-verified** (`02`) |
| **Enforcement** | disabled（no ai-skill validator / commit-msg hook） |
| **Promotion** | pending independent audit（pointer/SHA → **Promoted**, not Verified） |

Phase 3 / 5 / validator hook **延後**。最有價值資產：evidence maturity 語言穩定，非再擴文件。

**Owner**: linyihong
**建立日期**: 2026-06-29
**最後修訂**: 2026-06-29（maturity ladder 分層：Verified = behavior proven；Promoted = independently auditable）
**Priority**: P1
**Scope**: workflow advisory + dogfood 產物（planvalidate advisory scan，`Blocking=false`）；**不**接入 commit-msg block / runtime projection / enforcement

## Executive summary

**目標（修訂後）**：引入 **behavior-preserving structure preparation** 作為 **implementation execution mode**，不是補一個 Fowler 技巧或新的 lifecycle stage。

`workflow/software-delivery/` 重心落點：

| 區段 | 強度 |
|------|------|
| 上游：需求 → 合約 → 可驗證 | 強 |
| **中段：如何安全地改 code** | **弱（本 plan 補環）** |
| 下游：驗證 → closure → evidence | 強 |

Preparatory refactoring 在此的定位是 **implementation 的 execution governance**，銜接 evidence accumulation、parity governance、plan-first ordering、surgical execution — 把 tacit judgment（「先重構再加功能」）變成 contract + workflow。

**關鍵架構決策（stakeholder 2026-06-29）**：

```text
implementation（單一 lifecycle stage，不新增 stage）
  execution_modes:
    - direct_change          # 結構已足夠，直接做 feature
    - preparatory_refactoring # structure mode → feature mode，仍在 implementation 內
```

**不採用**：獨立 `sd-preparatory-refactoring` slice 變成 intake → design → preparatory refactor → implementation → validation 的線性第五階。

---

## Decision Rationale

### Problem & Why Now

**觸發證據**（2026-06-29 workflow review + stakeholder 確認）：

- Agent 在既有 code 上加功能時兩極擺盪：過度 surgical（硬塞）vs 過度 refactor（scope creep）
- 現有分類有 `change_kind`，但 **缺 `change_kind` × `execution_mode` 雙軸路由**：feature 被 structure block 時無法選 `preparatory_refactoring` mode，易被誤吸進 replacement parity
- 最大分類衝突不是 surgical-changes，而是 agent 把 preparatory work **誤塞進 replacement parity**
- 「兩頂帽子」若只寫 commit 規則太弱 — 需升級為 **Change Intent Lock**（execution contract）
- Preparatory refactoring 最大風險是 **永遠在整理** — 需 explicit **stop condition**

**Why now**：cognitive slice Phase 2 已穩；`implementation/README.md` 仍是 stub — 正好把 retained `sd-implementation` 擴成 execution mode 契約，而非再長一個 slice（避免 replacement / migration / spike / preparatory 無限增殖）。

### Decision

#### D1 — Implementation execution mode（非新流程）

擴充 **`workflow/software-delivery/implementation/`**（canonical surface for `sd-implementation`），定義：

```yaml
implementation:
  execution_modes:
    - id: direct_change
      when: structure sufficient; feature can land locally without seam work
    - id: preparatory_refactoring
      when: change_kind is feature and blocked_by_structure is true
      sub_modes:
        - structure   # behavior_change.allowed: false
        - feature     # new acceptance / observable behavior
```

Fowler 八步 playbook **降級為 structure mode 的 reference recipe**（`recommended_recipe`，advisory checklist），不是新 lifecycle step，也不是 mandatory order。

#### D2 — Change Intent Lock（升級兩頂帽子）

每個 implementation plan task step **必須**宣告 intent；validation 依 intent 分流：

| intent | behavior_change.allowed | validation 要求 |
|--------|-------------------------|-----------------|
| `structure` | `false` | observable equivalence / parity（舊行為仍受保護） |
| `feature` | `true` | new acceptance / BDD / contract proof |

範例（implementation-plan template）：

```yaml
steps:
  - id: prep-01
    intent: structure
    behavior_change:
      allowed: false
    action: extract apply_ranges()
    checkpoint:
      observable_equivalence:
        required: true

  - id: prep-02
    intent: structure
    behavior_change:
      allowed: false
    action: insert delegation seam (noop / adapter / wrapper — project-local technique)

  - id: feat-01
    intent: feature
    behavior_change:
      allowed: true
    action: add multi-range support
    validation: new acceptance criteria
```

Commit / PR 分離是 **Change Intent Lock 的副產物**，不是主要治理機制。

**Intent Transition Rule**（防 structure ↔ feature 來回切帽）：

```yaml
transition:
  - from: structure
    to: feature
    require:
      - observable_equivalence_passed   # 最近一個 structure step 的 checkpoint 已通過

  - from: feature
    to: structure
    require:
      - explicit_reopen_reason          # 記錄為何 feature 無法繼續、需重新開 structure
```

預設禁止：`feature → structure → feature → structure` 無理由 oscillation。`explicit_reopen_reason` 必須可 review（例如新發現的 seam 不足、acceptance 無法 expressible），不可是「再整理一下」。

**Intent state machine（validator 演進方向）**：未來若 advisory validator，優先驗 **`illegal transition`**（guard 未滿足即 transition），而非驗「有沒有做 structure step」。State + guard + transition condition 已足夠構成 execution protocol 骨架。

#### D3 — Observable Equivalence Checkpoint（抽象 no-op）

不綁死 `insert no-op` / `no-op test` 為唯一手法。Canonical 概念：

> 此步允許內部結構改變，但 **observable contract 必須等價**。

允許的 project-local 實作手法（非 exhaustive）：noop adapter、delegation、wrapper、branch disabled、dead-path extraction、fixture-backed parity assert、targeted mutation check。

每個 `intent: structure` step 需標記 `checkpoint.observable_equivalence.required: true` 及證據類型。

#### D4 — Intake taxonomy：`change_kind` × `execution_mode` 雙軸（Q6 resolved）

**不**把 structure preparation 塞進 `change_kind`。真正模型是 **雙軸**：

```yaml
# 軸 1 — intake：這次變更的本質是什麼？
change_kind:
  - feature
  - bugfix
  - replacement
  - migration
  - internal_refactor      # 行為不變、不為特定 feature 服務

# 軸 2 — implementation：怎麼做？
execution_mode:
  - direct_change
  - preparatory_refactoring
```

**路由規則**（避免「這次到底是 feature 還 preparatory_refactor？」）：

```yaml
if:
  change_kind: feature
  blocked_by_structure: true
then:
  execution_mode: preparatory_refactoring
  do_not_classify_as:
    - replacement
    - migration
else_if:
  change_kind: feature
  blocked_by_structure: false
then:
  execution_mode: direct_change   # 預設
```

`preparatory_refactoring` 是 **execution mode**，不是 change kind。Glossary 用 `preparatory_refactoring`（mode）與 `structure_preparation`（概念）區分，**不**新增 `change_kind.preparatory_refactor`。

與 replacement parity gate **正交**：`change_kind: feature` + `execution_mode: preparatory_refactoring` **不**觸發 parity inventory；`replacement` / `migration` 仍走 `sd-intake` parity。

#### D5 — Stop condition（防無限整理）

Structure mode **必須**在以下任一成立時 exit，進入 feature mode 或 direct_change：

```yaml
exit_when:
  - target_change_becomes_local      # 下一步 feature 已是局部改動
  - target_test_becomes_expressible  # 目標 acceptance 可寫成 failing test
  - new_abstraction_created          # 新 seam 已存在，feature 有落腳點
```

禁止（即使 intent 標 structure 也不允許）：

```yaml
avoid:
  - broad_cleanup
  - style_only
  - debt_harvesting
  - opportunistic_refactor
```

原則：**當下一步 feature 已經容易加入，就停止重構；不要追求漂亮。**

**Force exit（否定條件）** — 正向 `exit_when` 未觸發但 structure 已無進展時，必須停止 structure mode：

```yaml
force_exit_when:
  - no_structure_progress_detected
  - abstraction_not_used_by_next_feature
```

**`no_structure_progress_detected`**（第一版 **不數值化** N-step，避免 gaming）— 在 plan 或 step 結束時判定，以下皆不成立即視為無進展：

- 下一步 feature 仍未變成 **local change**（`target_change_becomes_local` 仍 false）
- 目標 acceptance 仍 **不可 expressible** 為 failing test
- 新 abstraction / seam **未被下一步 feature step 引用**（與 `abstraction_not_used_by_next_feature` 重疊時以後者為準）

不在第一版規定「連續 2 步 / 3 步」— 兩步可能抽不出 seam，三步又可能太少；先以 **進展定義** 為準，dogfood 後再評估是否需 step budget。

**`abstraction_not_used_by_next_feature`**：anti-architecture-as-output — 整理出漂亮架構但下一步 feature 踩不上去，即 force exit。

觸發 `force_exit_when` 時：不得再加 structure step；必須 (a) 進入 feature mode、(b) 改 `direct_change` 硬做並縮小 scope、或 (c) 回 intake 重判 `change_kind` / `execution_mode` / 拆 plan — 不可無限 structure。

#### D6 — Surgical-changes reconciliation

[`surgical-changes.md`](../../workflow/software-delivery/surgical-changes.md) 的 scope 紀律仍適用於 **direct_change** 與 **feature intent** steps。

**例外**：`intent: structure` steps 在 implementation plan 內、有 Observable Equivalence Checkpoint、且符合 stop condition — 不算 scope creep。不得用 preparatory 名義做 `avoid` 清單內的工作。

### Alternatives Considered

- **A. 獨立 `preparatory-refactoring.md` lifecycle slice** — reject（stakeholder 2026-06-29）：易變第五階段，與 execution-flow 不一致，slice 增殖
- **B. 兩頂帽子僅 commit 分離** — reject：太弱，agent/reviewer 無法機械判斷
- **C. Fowler 八步作 canonical workflow** — reject：no-op 是實例非抽象；改為 recipe + Observable Equivalence Checkpoint
- **D. Implementation execution mode + Change Intent Lock + intake taxonomy（accept）**

### Why Not an ADR Yet

- Change Intent Lock schema 未 dogfood → **dogfood 2026-06-29**（見 dogfood-evidence；advisory illegal transition）
- `change_kind` × `execution_mode` 雙軸與 change-brief 對齊方式未實測
- Stop condition 是否需 advisory validator 未驗證

### ADR Promotion Criteria（completed 時驗證）

- [ ] foundational + cross-session + cross-project + expensive-to-reverse + explains-why 全中
- [ ] ≥1 真實任務：`preparatory_refactoring` mode + structure steps + feature step + stop condition 觸發
- [ ] Reviewer 可依 intent 機械判斷 validation 要求（無 ambiguity）
- [ ] 無 agent 將 preparatory work 誤分類為 replacement
- [ ] Open Questions 全解

### Consequences

#### 正面

- 補 **execution governance 缺環**，符合 repo 主線（tacit judgment → contract → workflow）
- 不 proliferate lifecycle slices
- Change Intent Lock 可接 plan-first ordering、implementation-plan template、未來 advisory validator

#### 負面

- `implementation/` 從 stub 擴成 canonical — 需同步 execution-flow 導航（仍 **不** 拆成獨立 top-level slice id）

#### Compatibility（既有 plan / template 遷移）

```yaml
compatibility:
  existing_plans:
    default_execution_mode: direct_change   # 未宣告 execution_mode 的 plan 視為 direct_change
  existing_implementation_plans:
    missing_intent_on_steps: allowed        # 第一版不要求回溯補 intent
    missing_execution_mode: direct_change
```

避免舊 plan 被解讀成「全部必須補 `execution_mode` / intent schema 才能執行」。新欄位為 **opt-in enrichment**；僅當 task 走 `preparatory_refactoring` 或 plan 明確選 mode 時才強制 Change Intent Lock。

#### 風險

- execution mode 與 cognitive slice taxonomy 命名需對齊（implementation 內 mode vs 新 slice id）
- Stop condition 若寫太模糊，agent 仍會 over-refactor

**Glossary Impact**: yes
- 新引入：`implementation_execution_mode`、`preparatory_refactoring`（execution mode）、`change_intent_lock`、`intent_state_machine`、`illegal_transition`、`structure_intent`、`feature_intent`、`observable_equivalence_checkpoint`、`structure_preparation`（概念，非 change_kind）
- **不**引入：`change_kind.preparatory_refactor`（Q6 resolved — 雙軸）
- 需加入 [`knowledge/glossary/ai-skill.md`](../../knowledge/glossary/ai-skill.md)

---

## Runtime Execution Path

**Doc-only advisory。** 不新增 lifecycle `steps` order；擴充 `sd-implementation` loading surface 與 `execution-flow.yaml` `loading_surfaces.implementation`（若尚無則新增指向 `implementation/README.md` 或 `implementation/execution-modes.md`）。

| 項目 | 處置 |
|------|------|
| 新 lifecycle stage | **否** |
| 新 top-level cognitive slice | **否**（擴 retained `sd-implementation`） |
| commit-msg validator | **本輪否**；Change Intent Lock **不得**急於 runtime 化 |
| validation scenario | Phase 4 可選 + dogfood 必填 |

**Dogfood-before-validator（stakeholder 2026-06-29）**：Phase 4 ≥1 真實任務必須穩定跑通 `structure → stop → feature` 後，才評估 advisory validator / projection。若 dogfood 自然成立，validator 多半只是 schema 搬運 — 反之則先修 contract 文字，不加機械 block。

**Per-surface consumer 表**：N/A

---

## Open Questions

| # | 問題 | 預設 / 處置 |
|---|------|-------------|
| Q1 | 正文放 `implementation/README.md` 還是 `implementation/execution-modes.md`？ | Phase 0 依 document-sizing 決定；>150 行拆子檔 + README index |
| Q2 | `change_kind` 是否取代 change-brief `Change Type` enum？ | Phase 2 對齊表；預設擴充而非 breaking rename |
| Q3 | Skip-first acceptance test 是否 structure mode 的 recommended checkpoint？ | 是 recipe，非 mandatory technique |
| Q4 | plan-first ordering sub-plan merge 順序？ | implementation-plan template 兩 plan 都改；本 plan 定義 intent schema |
| Q5 | Change Intent Lock 第一版是否要求 YAML frontmatter vs markdown table？ | Phase 1 選最低摩擦格式（plan template 優先） |
| Q6 | `change_kind` 是否包含 preparatory 維度？ | **resolved**（第三輪）：雙軸 — `change_kind` 不含 preparatory；用 `execution_mode: preparatory_refactoring` |
| Q7 | structure progress 是否數值化 N-step？ | **resolved**（第三輪）：第一版用 `no_structure_progress_detected` 定義進展，不設 N |

---

## Phase 0 — Preflight

### Phase 0.0 — Open Questions 核對

- [ ] 已讀 §Open Questions 全部條目
- [ ] 對每條標記 `resolved` / `still-open` / `deferred`
- [ ] `resolved` 回寫

| Open Question | 處置 | 證據 |
|---|---|---|
| Q1–Q5 | still-open | Phase 0.1 |
| Q6–Q7 | resolved | 第三輪 stakeholder review |

### Phase 0.1 — 架構盤點

- [ ] 讀 [`workflow/software-delivery/execution-flow.md`](../../workflow/software-delivery/execution-flow.md) §sd-implementation 保留決策
- [ ] 讀 [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) §3 granularity（不 over-fragment）
- [ ] 讀 [`intake.md`](../../workflow/software-delivery/intake.md) 重構 / parity 分類
- [ ] 讀 [`surgical-changes.md`](../../workflow/software-delivery/surgical-changes.md)
- [ ] 讀 [`test-strategy.md`](../../workflow/software-delivery/test-strategy.md) old vs new behavior
- [ ] 讀 [`templates/implementation-plan-template.md`](../../workflow/software-delivery/templates/implementation-plan-template.md)
- [ ] 讀 plan-first ordering sub-plan 接點

### Phase 0.2 — Stakeholder decisions 記錄

- [x] **不**新增獨立 lifecycle slice — 改 implementation execution mode（2026-06-29 stakeholder）
- [x] Change Intent Lock > commit-only 兩頂帽子
- [x] Observable Equivalence Checkpoint > no-op 綁死
- [x] `change_kind` × `execution_mode` 雙軸；preparatory 為 mode 非 kind（2026-06-29 第三輪）
- [x] Stop condition 必填
- [x] Intent Transition Rule + force_exit_when（2026-06-29 第二輪）
- [x] Compatibility default `direct_change`（2026-06-29 第二輪）
- [x] Dogfood-before-validator；不急 runtime 化 Change Intent Lock（2026-06-29 第二輪）

---

## Phase 1 — Implementation execution mode 正文

- [x] 撰寫 `workflow/software-delivery/implementation/` canonical 正文（README 或 execution-modes.md）
- [x] 定義 `direct_change` vs `preparatory_refactoring`
- [x] 定義 Change Intent Lock schema（structure / feature + behavior_change.allowed）
- [x] 定義 Intent Transition Rule（structure → feature 需 equivalence；feature → structure 需 explicit_reopen_reason）
- [x] 定義 Observable Equivalence Checkpoint
- [x] 收錄 Fowler 八步為 **structure mode recommended_recipe**（advisory checklist，非 mandatory order）
- [x] 定義 stop condition（exit_when + force_exit_when + avoid）
- [x] 定義 compatibility default（未宣告 execution_mode → direct_change）
- [x] §Failure modes：infinite refactor、intent oscillation、intent 混用、誤用 replacement parity、skip checkpoint

**完成條件**：agent 在 implementation 階段可選 mode，不需載入額外 lifecycle slice。

---

## Phase 2 — Intake taxonomy & templates

- [x] [`intake.md`](../../workflow/software-delivery/intake.md)：`change_kind` 表 + `blocked_by_structure` → `execution_mode` 路由 + 與 replacement 正交說明
- [x] 雙軸對照表（修訂）：

  | change_kind | blocked_by_structure? | parity? | execution_mode |
  |-------------|----------------------|---------|----------------|
  | `internal_refactor` | n/a | 否 | direct_change 或 structure-only batch（無 feature intent） |
  | `replacement` / `migration` | n/a | **是** | intake parity gate；通常非 preparatory mode |
  | `feature` | false | 否 | `direct_change` |
  | `feature` | true | 否 | `preparatory_refactoring` |
  | `bugfix` | 視情況 | 否 | 預設 `direct_change` |

- [x] [`templates/implementation-plan-template.md`](../../workflow/software-delivery/templates/implementation-plan-template.md)：
  - `execution_mode: direct_change | preparatory_refactoring`
  - `steps[]` with `intent` + `behavior_change` + optional `checkpoint`
  - Stop condition checklist
- [x] [`templates/change-brief-template.md`](../../workflow/software-delivery/templates/change-brief-template.md)：`change_kind` 對齊

---

## Phase 3 — Routing & cross-links（不新增 lifecycle step）

- [x] [`execution-flow.md`](../../workflow/software-delivery/execution-flow.md)：sd-implementation 列更新 — 指向 implementation execution modes（Phase 1 已落地；見 §Phase 2 進度 2026-06-29）
- [x] [`execution-flow.yaml`](../../workflow/software-delivery/execution-flow.yaml)：`loading_surfaces.implementation`（`source: implementation/execution-modes.md`；**不**加新 ordered step）
- [ ] [`README.md`](../../workflow/software-delivery/README.md) 進入方式
- [ ] [`test-strategy.md`](../../workflow/software-delivery/test-strategy.md) cross-link：intent → validation 分流
- [ ] [`surgical-changes.md`](../../workflow/software-delivery/surgical-changes.md) reconciliation 一節
- [ ] [`review-checklist.md`](../../workflow/software-delivery/review-checklist.md)：intent lock + stop condition 檢查

**Phase 3 進度（2026-06-29）**：routing surface 已在 execution-flow 接通；README / test-strategy / surgical-changes / review-checklist 四檔仍待補 cross-link（觀察期內可完成，不阻塞 dogfood 結論）。

---

## Phase 4 — Dogfood & optional scenario

- [x] **必填（validator 前置 gate）**：≥1 真實任務走 `change_kind: feature` + `execution_mode: preparatory_refactoring`
- [x] **有效證據路徑（dual path collected）**：
  - **Happy path**：structure → checkpoint → stop (`exit_when`) → feature — [`02`](02-vidoe-test-project-dogfood-evidence.md) **partial-verified**（structure-transition only；feature 未閉環）
  - **Stop 設計驗證 path**：structure → structure → `force_exit_when` → 縮 scope — [`01`](01-dogfood-evidence.md) **verified**
- [x] 記錄：change_kind、execution_mode、intent 序列、transition 理由（含 illegal transition 有無）、checkpoint、exit_when / force_exit_when 觸發理由
- [x] Dogfood 通過後才開 maturity ladder（優先 **illegal transition** validator）— 本 phase **不** runtime 化
- [ ] （可選）`validation/scenarios/software-delivery/implementation-mode-preparatory-refactoring.yaml`

### Evidence maturity ladder（Phase 4 收斂）

```text
Observed → Partial Verified → Verified (behavior proven) → Promoted (independently auditable)
```

| Evidence | Class | Role |
|----------|-------|------|
| [`01-dogfood-evidence.md`](01-dogfood-evidence.md) | **verified** | stop / `force_exit` 機制有效 |
| [`02-vidoe-test-project-dogfood-evidence.md`](02-vidoe-test-project-dogfood-evidence.md) | **partial-verified** | structure→transition 可運作；equivalence 未證明 |

**Phase 4 status (do not over-read)**:

- [x] dual evidence path **collected**
- [x] `force_exit` path **verified** (`01`)
- [ ] happy path **completed** — **no**
- [x] happy path **partial-verified** (`02`)
- [ ] happy path **Verified** — pending Gate A（equivalence / behavior proven）
- [ ] happy path **Promoted** — pending pointer/SHA + independent audit (+ future validator wiring)

**Upgrade gates**（collection ≠ promotion）:

| Gate | Blocks | Meaning |
|------|--------|---------|
| **Gate A** — `checkpoint_valid` / observable equivalence | **Verified** | `checkpoint_exists` ≠ `checkpoint_valid`; regression proof required |
| **Gate B** — canonical `exit_when` | **Partial Verified** (recorded) | `02` primary: `target_test_becomes_expressible` — already mapped |
| **Pointer / SHA** — external reproducibility | **Promoted** | not blocking Phase 4 collected; not blocking Verified |

> dogfood 完成 ≠ 證據完成 ≠ 治理完成。Partial = 語意清楚（transition observed, equivalence open），不是「還沒做完 vs 做完沒證明」的模糊地帶。

---

## Phase 5 — Glossary & closeout

- [ ] glossary 註冊新詞彙（**注意**：既有 `execution_mode` 詞條指 cognitive FAST/NORMAL/DEEP；implementation mode 需用 `implementation_execution_mode` 或 `preparatory_refactoring` 獨立詞條，避免 collision）
- [ ] linked-updates 檢查
- [ ] archive plan + dogfood evidence

**Phase 5 阻塞說明（2026-06-29）**：dogfood 已雙路徑（`01` ai-skill + `02` Vidoe-Test）；glossary / linked-updates / archive 仍待觀察期結束或 Phase 3 四檔 cross-link 補齊後一次收口。

---

## 完成條件

- [x] Implementation execution mode 正文落地（非獨立 lifecycle slice）
- [x] Change Intent Lock 在 implementation-plan template 可機械填寫
- [x] Intake 雙軸（`change_kind` × `execution_mode`）與 replacement 邊界清楚
- [ ] Stop condition（exit_when + force_exit_when）+ avoid + Intent Transition Rule 在正文與 review checklist
- [ ] Compatibility default 在 implementation-plan template 或正文
- [ ] Observable Equivalence Checkpoint 定義（不綁 no-op）
- [x] ≥1 dogfood evidence — dual path: `01` **verified** + `02` **partial-verified**（happy path 未 completed / 未獨立稽核）
- [ ] doc-only 宣告；無 runtime projection 本輪

---

## Stakeholder 同意項目

- [x] 定位為 **implementation execution governance**，非 Fowler 流程補丁（2026-06-29）
- [x] 不新增 lifecycle stage / top-level slice
- [x] Change Intent Lock + Intent Transition Rule 方向 sign-off（2026-06-29 第二輪）
- [x] Stop condition（含 force_exit_when）方向 sign-off（2026-06-29 第二輪）
- [x] Dogfood 任務選定 — illegal-transition advisory scan；force_exit path（2026-06-29）
- [x] 雙軸 taxonomy（Q6 resolved，2026-06-29 第三輪）
- [x] 主結構可進 Phase 0 / Phase 1，無需再大改（2026-06-29 第三輪）
- [x] **第一輪閉環完成；進入觀察期，不進 enforcement**（2026-06-29）
- [x] **Evidence maturity 收斂**（2026-06-29）— Verified = behavior proven；Promoted = independently auditable；pointer 不擋 Verified

---

## 與其他 plans 的關係

| Plan | 關係 |
|------|------|
| [`02-software-delivery-plan-first-ordering`](2026-06-22-1009-plans-system-portability-and-delivery-integration/02-software-delivery-plan-first-ordering.md) | plan artifact ⟲ preflight；本 plan 定義 implementation plan 內 intent schema |
| **Vidoe-Test landscape player** | project-layer dogfood — [`02`](02-vidoe-test-project-dogfood-evidence.md) (**partial-verified**); Gate A → Verified; pointer → Promoted |
| [`gen3-workflow-analysis-cognitive-slice-decomposition`](../archived/2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md) | 延續 sd-implementation retained；**不**新增 slice id |
| Recovery / Release 擴充 | out of scope |

---

## Watch-Out List citation

[`architecture/ai-native-cognitive-ecosystem-system.md`](../../architecture/ai-native-cognitive-ecosystem-system.md) §Watch-Out List — 防 lifecycle slice 增殖、防 infinite refactor（stop condition）、防 tacit pattern 未 contract 化。

---

## 修訂記錄

| 日期 | 修改 | 原因 | 來源 |
|------|------|------|------|
| 2026-06-29 | 初稿：獨立 preparatory-refactoring slice | workflow gap 分析 | agent draft |
| 2026-06-29 | 重寫：implementation execution mode + Change Intent Lock + intake taxonomy + stop condition | stakeholder review | 本對話 |
| 2026-06-29 | 第二輪：Intent Transition Rule、force_exit_when、compatibility default、dogfood-before-validator、Q6 命名 | stakeholder maturity review | 本對話 |
| 2026-06-29 | Phase 4 force_exit dogfood；觀察期 sign-off；enforcement 延後 | stakeholder：dogfood 驗 contract 站得住 | 本對話 |
| 2026-06-29 | Phase 1 落地 execution-modes.md + execution-flow 導航；Phase 2 intake 雙軸 + templates intent 欄位 | implementation execution governance | agent |
| 2026-06-29 | Maturity ladder refine: Verified=behavior proven; Promoted=independently auditable; pointer blocks Promoted only | stakeholder gate semantics | 本對話 |
| 2026-06-29 | Vidoe-Test project-layer dogfood **partial-verified** | landscape Phase 0 guard + structure-transition | [`02-vidoe-test-project-dogfood-evidence.md`](02-vidoe-test-project-dogfood-evidence.md) |
