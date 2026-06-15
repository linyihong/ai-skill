---
id: 2026-06-12-1600-projection-break-evidence-taxonomy-experience-runtime
plan_kind: main
status: draft
owner: linyihong
created: 2026-06-12
priority: P1
required_for_completion: false
---

# Projection Break → Evidence Taxonomy → Experience Runtime

**Status**: `draft`
Owner: framework maintainer (linyihong)
**建立日期**：2026-06-12
**Priority**：**P1**（taxonomy + gate vocabulary + authority layer；experience-runtime **cross-cutting** 先落地，升 slice 延後）
**Phase 0**：✅ **完成**（vocabulary / ownership / authority / promotion 條件已收斂）
**Phase 1**：**Ready**（開始建 `validation/evidence-types/` catalog，非探索期）
**Review**：2026-06-12 stakeholder review — L3 命名、三分 taxonomy、authority layer、OQ-5 **reject inheritance**、experience-runtime cross-cutting、`timing_gate`→`temporal_behavior`
**Pilot**：某 downstream project 的 immersive media player preview gate incident（2026-06）。具體 drama / episode / deploy host 依 [`reusable-guidance-boundary.md`](../../enforcement/reusable-guidance-boundary.md) 留在專案 plan，不寫入本 plan。

## Why this plan exists

某 downstream project 的 player preview gate 失敗不是「少寫一個測試」，而是 **L2 專案投影 → L3 驗證能力** 的投影斷裂（projection break）：

| 層 | 問題 | 當時有 |
|---|---|---|
| Contract | 試看結束要出遮罩 | ✅ `.feature` |
| Behavior | 前端真的切狀態 | ⚠️ source assert（BDD string） |
| Validation Capability | Browser 實際操作、timing poll、navigation readback | ❌ |
| Evidence | artifact + proof shape（截圖 / DOM log / URL trace） | ❌ |

`.feature` 存在 ≠ UX 被驗證。這是 AI workflow 常見洞，且現有 software-delivery 雖已有 Journey Validation、Evidence Acquisition Layer、UI Governance，仍缺：

1. **Evidence Taxonomy 邊界** — gate 仍用 `browser_review` 等抽象詞；且 `evidence_type` 與 `collection_method` 開始混用
2. **Experience Runtime cross-cutting model** — immersive player 橫切 runtime + journey + validation + ui-contract；**不應過早升成 software-delivery slice**
3. **Failure → Authority → Evolution loop** — rule → gate → test 有，但缺 authority 決策，容易每次 failure 都升 framework

本 plan 把這次 incident 提煉成 framework 演進，而不是再堆 player 特例 rule。

## Decision Rationale

### Problem & Why Now

現行治理已可描述四層，但 agent 與 project workflow 仍容易把 **BDD contract pass** 誤當 **validation complete**：

```text
L0  Runtime obligations（bootstrap, close-out）
L1  Delivery model（software-delivery slices）
L2  Project concretization（workflow YAML, project overlay）
L3  Validation capability（實際驗證能力）← player 炸在這裡
    Evidence = artifact + proof shape（Validation 的產出，不是 Validation 的同義詞）
```

垂直責任鏈（避免循環定義）：

```text
Contract
  ↓
Behavior
  ↓
Validation Capability
  ↓
Evidence（artifact + proof shape）
```

既有相關 work 已部分覆蓋 adjacent 問題，但未收斂成同一 vocabulary：

- [`2026-06-08-1544-evidence-acquisition-layer.md`](../archived/2026-06-08-1544-evidence-acquisition-layer.md) — `collection_method` 介於 domain 與 mechanism
- [`2026-06-10-0908-user-journey-validation-integration.md`](../archived/2026-06-10-0908-user-journey-validation-integration.md) — Journey Specification vs Journey Execution
- [`2026-06-09-1040-experience-validation-pipeline-evolution.md`](2026-06-09-1040-experience-validation-pipeline-evolution.md) — Coverage Model watch-list（in-progress）
- [`2026-06-10-1718-software-delivery-governance-invariants.md`](../archived/2026-06-10-1718-software-delivery-governance-invariants.md) — validation-family invariants

缺口是 **completion claim 的 evidence requirement 仍太抽象**、**evidence_type / collection_method / artifact_shape 未三分**，以及 **client-heavy runtime 缺 cross-cutting 模型（但尚未到升 slice 的案例數）**。

### Decision

升級核心模型，並分四波落地：

```text
Contract ≠ Behavior
Behavior ≠ Validation Capability
Validation Capability produces Evidence（artifact + proof shape）
```

**C.1 — Evidence Taxonomy：三分，不混層（P1 doc + project gate vocabulary）**

新增 `validation/evidence-types/` 目錄。**只放 `evidence_type`（要證明什麼）**；`collection_method` 與 `artifact_shape` 在 README 與 project envelope 中對照，不寫進 type 檔名。

| 層 | 定義 | 歸屬 | Gate `requires` 可用？ |
|---|---|---|---|
| `evidence_type` | 要證明什麼 | `validation/evidence-types/*.md` | ✅ `evidence:user_visible` |
| `collection_method` | 怎麼取得 | 既有 Evidence Acquisition Layer | ❌ 不作 gate token |
| `artifact_shape` | 長什麼樣 | integration envelope / close-out | ❌ 不作 gate token |

範例（player preview overlay）：

```yaml
claim: preview_overlay_shown
evidence_type: user_visible
collection_method: browser_observation   # 不是 evidence_type
artifact_shape: screenshot                 # 不是 evidence_type
```

首批 **evidence_type** 檔（移除 `browser-observation`、`user-visible-proof` 等同層重疊）：

| File | evidence_type | 要證明什麼 |
|---|---|---|
| `source-contract.md` | `source_contract` | 靜態 contract / string / schema 對齊 |
| `user-visible.md` | `user_visible` | 使用者可見 UI 狀態（overlay、CTA、mask） |
| `navigation.md` | `navigation` | route / href / history / basePath 正確性 |
| `state-persistence.md` | `state_persistence` | sessionStorage / cookie / DB readback 跨導航保留 |
| `media-playback.md` | `media_playback` | video element state, HLS load, pause/resume |
| `temporal-behavior.md` | `temporal_behavior` | 時間邊界內的可觀察狀態轉換（preview 截止、debounce、poll 觸發、buffer stall 後恢復）— **不是** failure_class 名稱 |

命名原則：`evidence_type` 必須回答「證明了什麼」，不是「為什麼容易壞」。`timing_gate` 已 reject（偏 failure_class）；改用 `temporal_behavior`。

每檔必含：定義、non-goals、適用 failure class、`supported_collection_methods`、`supported_artifact_shapes`、proxy trap。**不**建立 token inheritance 樹。

**C.2 — Gate vocabulary：只要求 evidence_type（P1 project pilot, P2 framework template）**

Project workflow gate 從：

```yaml
requires:
  - integration
  - browser_review
```

改為：

```yaml
gate.<project>.validation_complete:
  requires:
    - evidence:user_visible
    - evidence:navigation
    - evidence:temporal_behavior   # 當 changed surface 含 time-bounded client transition
```

`browser_review` 降級為 **collection activity 的 umbrella label**（人類可讀摘要），不再作為 pass/fail token。`browser_observation` 是 `collection_method`，**不得**出現在 `requires:`。

**C.3 — Experience Runtime：cross-cutting concern，延後升 slice（P3 doc-only）**

Player 同時橫切 runtime + journey + validation + ui-contract。過早升 `software-delivery` slice 易與 `experience-validation`、`runtime-journey` 互吃。

Phase 3 先建 **cross-cutting** 文件，不註冊新 delivery slice：

```text
workflow/cross-cutting/experience-runtime/
  README.md
  player.yaml          # pilot template only
```

Player pilot state machine（文件模板，非 slice contract）：

```yaml
state_machine:
  - idle
  - loading
  - preview
  - gated
  - purchasing
  - unlocked
  - playing
  - suspended
  - restored

required_validation:
  - persistence
  - interruption
  - recovery
```

**Slice promotion 條件（延後）**：至少 player、editor、onboarding 三個案例收斂後，才評估升 `software-delivery` slice。

**C.4 — Failure → Authority → Evolution loop（P2）**

把 incident pattern 結構化（補 authority，避免每次都升 framework）：

```text
Failure
  → Classification（projection break layer + missing evidence_type）
  → Authority Decision（誰有資格改什麼）
  → Evolution Target（rule / gate / scenario / code / playbook）
  → Writeback
```

| Authority 類型 | 典型 writeback 目標 |
|---|---|
| framework invariant | Ai-skill validation / workflow / scenario |
| domain pattern | project overlay rule |
| implementation defect | code + integration test |
| env / deploy incident | deploy playbook + smoke checklist |

首個 catalog entry：`player-preview-gate-projection-break` — authority = **domain pattern**（project overlay）+ **implementation**（code），非 framework invariant。

Catalog 必填欄位（避免退化成 incident log）：

```yaml
failure: projection_break_at_L2_to_L3
classification: missing evidence_type user_visible + temporal_behavior
authority: [domain_pattern, implementation_defect]
evolution_target: [project overlay rule, integration test, gate requires]
writeback: [project-docs, code]
counterfactual: >
  若當時 gate 要求 evidence:user_visible 且 integration envelope
  含 claim preview_overlay_shown + artifact，應可在 merge 前發現 overlay 未出現。
```

### Alternatives Considered

- **A. 只加 player rule / player-spec（project-only）**：reject — 不解決 L2→L3 結構洞，下一個 immersive surface 會重演。
- **B. 把 BDD 升級成 Validation**：reject — BDD owns Journey **Specification**；validation owns **Execution**（見 archived user-journey plan）。
- **C. 過早升 experience-runtime 為 software-delivery slice**：reject — 案例不足，易與 journey / validation slice 互吃。
- **D. 漸進：三分 taxonomy + gate evidence_type tokens + cross-cutting experience-runtime + authority evolution loop（accept）**

### Why Not an ADR Yet

- Evidence type 三分邊界與 cross-cutting experience-runtime 仍需 downstream pilot 壓力驗證
- `browser_review` → `evidence:*` migration 需觀察 project workflow 摩擦

### ADR Promotion Criteria（completed 時驗證）

- [ ] foundational + cross-session + cross-project + expensive-to-reverse + explains-why 全中
- [ ] 至少一個 downstream project 用 evidence tokens 完成 real task closure
- [ ] Open Questions 全解
- [ ] 沒有更輕的 promotion target 適用
- [ ] validation scenarios 與 project gate 真實引用 evidence types

### Consequences（預期）

#### 正面

- Completion claim 可機械追問「缺哪種 evidence_type」，而非「有沒有做 browser review」
- Authority layer 避免每次 player bug 都升 framework invariant
- Immersive client runtime 有 cross-cutting 模型，暫不污染 software-delivery slice 命名空間

#### 負面

- Project workflow YAML 與 close-out report 欄位變多
- Agent 可能 over-claim evidence token 而未附 artifact

#### 風險

- Evidence types 膨脹成垃圾桶 taxonomy → 三分法 + document-sizing + watch-list 節制
- experience-runtime 過早升 slice → cross-cutting 先寫，三案例後再 promotion
- evidence token ontology 膨脹 → OQ-5 **reject inheritance**；用 `supported_collection_methods` / `supported_artifact_shapes` 對照

**Glossary Impact**: yes — candidate terms: `projection_break`, `validation_capability`, `evidence_type`, `artifact_shape`, `authority_decision`, `experience_runtime`（cross-cutting，非 slice 名稱）

## Runtime Execution Path

**Phase 1–2 scope = doc + project pilot vocabulary；不新增 runtime YAML gate。**

| Phase | Runtime owner | Trigger | Action | Evidence |
|---|---|---|---|---|
| P1 | `validation/evidence-types/` docs | project validation gap | agent reads taxonomy when claiming UX complete | downstream project plan cross-ref |
| P2 | `workflow/software-delivery/validation/evidence-gate-vocabulary.md` | project gate update | evidence_type tokens only in `requires:` | pilot closure |
| P3 | `workflow/cross-cutting/experience-runtime/` | immersive client feature | load cross-cutting README + player.yaml template | doc-only; no slice registration |
| P4 | `validation/scenarios/failure-derived/` + authority catalog | failure-derived | Failure→Authority→Evolution writeback | scenario pass in CI |

**Deferred Runtime Projection**：本 plan Phase 1 不 project 新 YAML 到 `runtime.db`。Graduation condition：downstream pilot 連續 2 次 player-class task 使用 evidence tokens 完成 closure 後，再開 follow-up plan 接 `gate.*.evidence_requirements` runtime-lite projection。

**Per-surface consumer 表**：Phase 1 N/A（doc-only）。Phase 3 新增 scenario 時再填。

## Open Questions

| # | Question | 處置 |
|---|---|---|
| OQ-1 | `evidence:*` tokens 是否進 Ai-skill runtime registry，或僅 project workflow？ | still-open — pilot 先 project-local |
| OQ-2 | cross-cutting experience-runtime 與 journey / validation execution 的 boundary checklist？ | still-open — Phase 3 README 必答 |
| OQ-3 | `temporal_behavior` 的最低 `artifact_shape` 是 poll log 還是錄影？ | still-open — player pilot 決定 |
| OQ-4 | failure evolution 的 authority 預設規則是否機械化？ | deferred — Phase 4 catalog 先 doc-only |
| OQ-5 | evidence token 是否允許 inheritance？ | **resolved — reject inheritance**。`evidence_type` 不是 subtype system；type 與 method/shape 的關係只用各 type 檔內 `supported_collection_methods` / `supported_artifact_shapes` 對照表表達，禁止 `user_visible → screenshot` 繼承樹 |

## 與其他 plans 的關係

| Plan | Relationship |
|---|---|
| [`2026-06-09-1040-experience-validation-pipeline-evolution.md`](2026-06-09-1040-experience-validation-pipeline-evolution.md) | Parent watch-list for Coverage Model; 本 plan 補 evidence **type** 維度，不取代 render_context work |
| [`2026-06-08-1544-evidence-acquisition-layer.md`](../archived/2026-06-08-1544-evidence-acquisition-layer.md) | `collection_method` = acquisition；`evidence_type` = 要證明什麼；兩者不得互換 |
| [`2026-06-10-0908-user-journey-validation-integration.md`](../archived/2026-06-10-0908-user-journey-validation-integration.md) | Journey 仍用 `validation_scope`；experience-runtime cross-cutting 管 state machine，不取代 journey execution |
| Downstream pilot project plan（專案 repo `docs/plans/`） | 第一個消費本 plan P1 taxonomy；路徑與 incident 證據不寫入 Ai-skill（見 reusable-guidance-boundary） |

## Phase 0 — Pre-Build Interrogation

### Phase 0.0 — Open Questions 核對（公版，必填）

- [x] 已讀本 plan §Open Questions 全部條目
- [x] 對每條標記 `resolved` / `still-open` / `deferred`
- [x] `resolved` 的條目已同步勾選 / 附註於 §Open Questions
- [x] 若盤點新發現問題，已加入 §Open Questions

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| OQ-1 runtime registry | still-open | 等 downstream pilot |
| OQ-2 cross-cutting boundary | still-open | Phase 3 交付 |
| OQ-3 temporal_behavior artifact_shape | still-open | preview-gate integration 實作時決定 |
| OQ-4 authority mechanization | deferred | Phase 4 catalog |
| OQ-5 token inheritance | **resolved** | 2026-06-12 review：**reject inheritance**；supported_* 對照表 |

### Phase 0.1 — Architecture Compatibility Preflight

- [x] Candidate paths：`validation/`, `workflow/software-delivery/`, `plans/active/` 存在
- [x] 不與 `collection_method` ownership 衝突 — 三分法明確分工
- [x] 不新增 `sd-browser-review` slice；不過早註冊 experience-runtime delivery slice
- [x] OQ-5 resolved — reject token inheritance
- [x] Decision: **Phase 0 complete**；**Phase 1 Ready** — 開始建 evidence-types catalog

## Phase 1 — Evidence Taxonomy Catalog（Ready）

- [ ] 建立 `validation/evidence-types/README.md`（L3=Validation Capability + 三分法 + **OQ-5 reject inheritance**）
- [ ] 建立 6 個 evidence_type 檔（含 `temporal-behavior.md`；**不含** `timing_gate`）
- [ ] 每 type 檔含 `supported_collection_methods` / `supported_artifact_shapes`（禁止 inheritance 樹）
- [ ] 更新 `validation/README.md` 與 `workflow/software-delivery/validation/README.md` 交叉引用
- [ ] 新增 validation scenario stub：`validation/scenarios/software-delivery/evidence-type-projection-break-v1.yaml`
- [ ] 更新 [`2026-06-09-1040-experience-validation-pipeline-evolution.md`](2026-06-09-1040-experience-validation-pipeline-evolution.md) §與其他 plans 的關係

**Phase 1 完成條件**：6 個 evidence_type 可機械引用；README 明記 reject inheritance。

## Phase 2 — Authority + Gate Vocabulary（Framework + Project 並行）

- [ ] 建立 `workflow/software-delivery/validation/evidence-gate-vocabulary.md`（只允許 `evidence_type` 進 `requires:`）
- [ ] 建立 authority decision 對照表（framework / domain / implementation / env）
- [ ] 定義 trace chain：**gate → claim → artifact**（禁止只有 artifact 無 claim）
- [ ] 支援 downstream pilot 更新 project workflow validation gate
- [ ] Close-out：`browser_review` = activity summary；envelope 含 claim + type + method + shape

**Phase 2 完成條件**：downstream pilot 以 projection contract 消費 taxonomy；至少 1 個 integration envelope 可追蹤 **gate → claim → artifact**。

## Phase 3 — Experience Runtime Cross-Cutting（P2，延後升 slice）

- [ ] 建立 `workflow/cross-cutting/experience-runtime/README.md`（**非** software-delivery slice）
- [ ] 建立 `workflow/cross-cutting/experience-runtime/player.yaml` pilot template
- [ ] 撰寫 cross-cutting 與 journey / validation / ui-contracts 的 boundary table
- [ ] 明確記錄 slice promotion 條件（player + editor + onboarding）

**Phase 3 完成條件**：player.yaml 作為 cross-cutting 模板可用；**不**註冊 `sd-experience-runtime` slice。

## Phase 4 — Failure → Authority → Evolution（P2/P3）

- [ ] 建立 failure catalog（必填：`failure` / `classification` / `authority` / `evolution_target` / `writeback` / **`counterfactual`**）
- [ ] 首 entry: player preview gate（見 Decision C.4 範例）
- [ ] 新增 scenario：`validation/scenarios/failure-derived/projection-break-missing-browser-evidence-v1.yaml`
- [ ] 評估是否 promotion 成 enforcement advisory（非 mechanical rule_class）

**Phase 4 完成條件**：catalog entry 含 counterfactual；同類 failure 有 authority-aware evolution path。

## 完成條件（Plan-level）

- [ ] Phase 1 evidence-types catalog 落地（6 types，含 `temporal_behavior`）
- [ ] Phase 2 gate→claim→artifact trace 有 real closure example
- [ ] Phase 3 cross-cutting experience-runtime 存在且**未**升 slice
- [ ] Phase 4 至少 1 条 failure→authority→evolution catalog entry
- [ ] Glossary candidate terms 已註冊或明確 defer
- [ ] 執行 Plan Completion Closure（validator + linked-updates + archive）

## Stakeholder 同意項目

- [ ] OQ-5 reject inheritance 可接受
- [ ] `temporal_behavior` 取代 `timing_gate` 可接受
- [ ] evidence_type / collection_method / artifact_shape 三分可接受
- [ ] `browser_review` 降級為 activity label 可接受
- [ ] experience-runtime 先 cross-cutting、三案例後再升 slice 可接受
- [ ] Failure→Authority→Evolution 四步可接受
- [ ] Phase 1–2 doc-only、不接 runtime gate 可接受

## Watch-Out List citation

- 防 scope drift：不把整個 frontend 塞進 experience-runtime cross-cutting
- 防 over-engineering：Phase 1 不建 runner framework；不過早升 slice
- 防 taxonomy 混層：新檔名不得用 collection_method 或 artifact_shape 冒充 evidence_type
- 防 gate token inflation：reject inheritance；requires 只列 evidence_type
