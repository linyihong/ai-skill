---
id: 2026-06-24-1100-adr-004-migration-completion
plan_kind: main
status: draft
execution_status: active
owner: linyihong
created: 2026-06-24
parent: null
baseline_ref: 2026-06-23-1500-adr-004-migration-drift-diagnosis
---

# ADR-004 Migration Completion Plan

> **這是 implementation plan，不是 diagnosis。** 診斷已於
> [`2026-06-23-1500-adr-004-migration-drift-diagnosis.md`](2026-06-23-1500-adr-004-migration-drift-diagnosis.md)
> 凍結（`status: frozen`）。本文件把那份 frozen diagnosis 當 **immutable baseline**，
> **不重述**證據、loop 圖或失效模式命名；只引用其錨點，並把 P0-A/B/C 與 Contract Layer
> 展開為可執行 phase。命名刻意是 *Completion*（讓 migration 真正完成）而非 *Repair*（修 bug）。

## 0. Context

### Reference Diagnosis (immutable baseline)
- Baseline：`2026-06-23-1500-adr-004-migration-drift-diagnosis`（frozen）。
- 引用而非複製：Layer 狀態表、Drift Propagation Loop、Semantic Validation Drift /
  Canonical Knowledge Propagation Failure、Canonical-Path Derivation Invariant、Contract
  Ownership 鏈，全部以 baseline 為準；本文件若與 baseline 衝突，以 baseline 為真。
- 關鍵錨點（取自 baseline「證據錨點」段，勿在此重抄行號內容）：
  `runtime.go:1213` 死 glob、`runtime.go:1876-1882` masking gate、`runtime_test.go:244,625,640`
  old-world fixtures、`routing-registry.yaml` `route.feedback.history`（Path 1 正確）。

### Out of Scope
- [ ] 重新論證 diagnosis（已 frozen）。
- [ ] 修改 `constitution/ADR-004`（憲法層正確，禁改——見 Rollback）。
- [ ] feedback lesson 內容/品質本身。
- [ ] Path 1（reference-first routing）行為變更——它正確且 load-bearing。
- [ ] <TODO: 其餘 out-of-scope 待 Phase 0 inventory 後補足>

### Success Definition
- 兩條 discovery path 一致：`Path 1 (direct read) == Path 2 (runtime query)` 對 lesson 同集合。
- canonical path 只有單一 owner（registry），無任何 consumer hardcode 路徑。
- Exit Criteria 的 derive 鏈五層全綠（見文末 Exit Criteria）。

---

## Phase 0 — Migration Inventory (Step 0)

> **排序鐵則：`inventory discovers / classification interprets`，不可反向。**
> 先定 A/B/C/D 再掃，會讓「分類先存在 → inventory 被迫塞進分類 → 新型 consumer 被忽略」——
> 這與本次 drift（fixture 先定義世界 → indexer 證明 fixture）同型。故拆三步，分類放最後。
>
> **掃描單位 = consumer，不是 string occurrence。** 且必須掃 **intent equivalence**，不只字面路徑：
> `feedback_history` / `feedback/history` / `MATCH 'feedback'`(FTS) / `route.feedback.*`(registry key) /
> README reference / generated path / helper function。教訓：沒有地方明寫舊路徑，系統仍消費舊世界觀。
>
> **Phase 0 成功條件（一句）：Every consumer must be enumerable before it becomes enforceable.**

### Step 0A — Reference Census（read-only · 只列舉不判讀）

固定欄位，**禁止**分級 / 打 P0 / 決定 owner / 判定 load-bearing：

| field | 說明 |
| --- | --- |
| file | 檔案 |
| location | line / symbol |
| reference | 命中的 path / token |
| access_mode | glob / FTS / registry / direct read / fixture / docs |
| intent | read / validate / route / describe / seed |
| execution_surface | runtime / test / docs / tooling |
| candidate_world | old / new / mixed |
| consumer_id | 暫時編號 |

**0A 種子（2026-06-24 sweep；尚未 resolve，僅 census）：**

| consumer_id | file:location | reference | access_mode | intent | surface | world |
| --- | --- | --- | --- | --- | --- | --- |
| C-01 | `runtime.go:1212` `runtimeIndexFeedbackRecords` | `skills/*/feedback_history` | glob | read | runtime | old |
| C-02 | `runtime.go:1146-1152` index assembly | feedback records append | direct read | read | runtime | (derives C-01) |
| C-03 | `runtime.go:1876,1894` `nativeRuntimeIndexFTSCheck` | `MATCH 'feedback'` | FTS | validate | runtime | mixed |
| C-04 | `runtime.go:1379-1422` `runtime query` FTS | `MATCH ?` keyword | FTS | read | runtime | (Path 2) |
| C-05 | `runtime_test.go:244,625,640` | `skills/demo/feedback_history/...` | fixture | seed/validate | test | old |
| C-06 | `close_loop.go:468` | `HasPrefix(path,"feedback/")` | direct read | route | tooling | new |
| C-07 | `hooks.go:1861-2017` `validateFeedbackLearningReport*` | token "feedback" (Report obligation) | n/a | validate | tooling | independent (NAME COLLISION) |
| C-08 | `router_proposals.go:107` | "Discovery feedback loop" | n/a | describe | tooling | independent (NAME COLLISION) |
| C-09 | `routing-registry.yaml:1951` `route.feedback.history` | `feedback/history/README.md` | registry | route | tooling | new |
| C-10 | `routing-registry.yaml:941-958` `route.feedback.promotion-pipeline` | `feedback/promotion/...` | registry | route | tooling | new |
| C-11 | `metadata/rules/feedback-lessons.yaml` | `feedback/history/<domain>/` | direct read | describe | docs | new |
| C-12 | `enforcement/feedback-lessons.md`, `content-layering.md` | sink declaration | direct read | describe | docs | new |
| C-13 | `knowledge/summaries/feedback-promotion-pipeline.md:16` | `保留原 feedback_history/` | docs | describe | docs | old |
| C-14 | `enforcement/failure-learning-system.md:115,201` | `feedback_history/` durable loc | docs | describe | docs | old |
| C-15 | `intelligence/engineering/agent-architecture/failure-recovery.md:56` | `feedback_history/` | docs | describe | docs | old |
| C-16 | `validation/scenarios/failure-derived/feedback-history-consolidation-v1.yaml` | `feedback/history/<domain>/` | docs | validate | test | new |
| C-17 | `traces/failure-derived/...-2026-05-13.yaml` | `舊路徑...仍存在（向後相容）` | docs | describe | docs | old (now false) |

> 種子未盡：`feedback/history` 新世界 doc reference 共 308 hits（按 top-dir：feedback 58 / intelligence 19 /
> enforcement 12 / workflow 10 / knowledge 10 …），多為 docs，待 0A 補齊 consumer 收斂。
> **NAME COLLISION 提醒**：token `feedback` 橫跨 lesson / Learning-Report obligation / Discovery loop
> 三 domain（C-07/C-08）→ 直接解釋 C-03 為何在任何 `feedback` 出現時都綠。

### Step 0A-completeness — Intent Reduction（機械收斂，非全文人工閱讀）

> **目標**：eliminate hidden executable references。**非目標**：classify consumers。
> 對 census 候選先機械補三欄，再依 Exit Rule 收斂：

| field | 說明 |
| --- | --- |
| evidence_kind | read / validate / route / seed / describe / mention |
| authority_claim | canonical / derived / descriptive / none |
| executable | Y / N |

**0A Exit Rule（停止條件）**：
> **If a reference cannot execute and cannot assert authority, it does not block inventory closure.**
> 即 0A 關閉條件**不是**「所有 feedback 字串都被看過」，而是「沒有剩餘
> *executable + authority-claim* 的未分類引用」。允許 `describe` / `mention` / historical note 留在 census 外，
> 否則 inventory 永遠關不掉。

**Reduction funnel（2026-06-24 sweep）**：148 files（308 occ）→ 135 `.md`（describe/mention，
非 executable，依 Exit Rule 出局）+ 13 `.yaml` + code。executable+authority-claim 收斂後，
seed 外新發現的 **hidden executable refs**（僅 enumerate，**不** resolve）：

| consumer_id | file:location | reference | access_mode | evidence_kind | authority_claim | executable | world |
| --- | --- | --- | --- | --- | --- | --- | --- |
| C-18 | `metadata/recovery/domain-policies.yaml:36` | `source: feedback/history/apk-analysis/README.md` | direct read | route/read | derived | Y | new |
| C-19 | `validation/scenarios/failure-derived/runtime-recovery-navigation-mismatch.yaml:58` | `feedback/history/apk-analysis/README.md` | fixture/scenario | validate | derived | Y | new |
| C-20 | `validation/scenarios/failure-derived/skill-local-feedback-bypass-v1.yaml` | `feedback/history/<domain>/` 寫入期望 | fixture/scenario | validate | derived | Y | new |
| C-21 | `validation/scenarios/software-delivery/slice-load-scenario-d-placement-negative.yaml:23` | `feedback/history/manual-observation-...md` | fixture/scenario | describe | none | Y(scenario) | new |
| (note) | `hooks.go:1990,1993` `isAllowedFeedbackValue("feedback-history"…)` | Learning-Report **Target enum** | n/a | validate | independent | Y | NAME COLLISION (≠ sink) |
| (note) | `detector.go:73` `"feedback":"advisory"` | route_type weight | n/a | route | independent | Y | NAME COLLISION (≠ sink) |

> C-21 是 `describe` + authority `none` → 依 Exit Rule **不阻擋** 0A 關閉（仍登錄以利追蹤）。
> code 已全掃，無其餘 executable sink-path 引用；135 `.md` 為 doc 描述，出局。
> **0A 可關閉條件已達**：無剩餘 *executable + authority-claim* 未登錄引用。

### Step 0B 預備觀察（記錄，禁止現在 resolve）

0A census 已自然露出**三種不同的東西**，但 0A 不判讀，留給 0B：

| 暫稱 | 例 | 0B 待答 |
| --- | --- | --- |
| real consumer | `runtimeIndexFeedbackRecords`(C-01)、`runtime query`(C-04) | 真的依賴 lesson？ |
| shadow consumer | `close_loop` path prefix(C-06) | **保留問號**：讀新路徑 ≠ 一定是 discovery consumer；discovery-consumer 與 enforcement-consumer 不一定同一件事 |
| authority illusion | health-check(C-03)、docs、validator、name-collision enum | 看似相關，實則不依賴/不主張 sink authority |

> C-06 刻意維持問號，不在 0A 叫它 consumer。real / shadow / illusion 的切分是 0B 的工作。

### Step 0B — Consumer Resolution（消除假 consumer）

> **0B 鐵則：先定義失效影響，再定義 consumer 類型。** 否則「誰讀」與「誰受 drift 影響」會混淆
> （本次 drift 同型陷阱）。對每個候選先問：**「若 `feedback/history` 消失，它會讓誰做不成事？」**
> 影響答出來，分類通常自己掉出來。
> （Discovery / Enforcement / Observability / Shadow 的型別定義目前是**假說、不落文件**；本文件只落
> 判準矩陣、0B-1/2/3 staging 與 capability 結果。）
>
> **0B Invariant（由 C-06 升級）：Path access alone does not establish consumer status.**
> 讀了 canonical path ≠ 是 consumer（C-06 讀 `feedback/` 但 impact=0）。此 invariant 預期可擋掉一整類誤報。
>
> **0B 工作假設（升級自 0B-1）：目前真正 load-bearing 的不是 discovery abstraction，而是
> registry-backed direct read（Path 1）。** 故問題不是「Path 2 壞→修 Path 2→恢復系統」，而是
> 「Path 1 活、Path 2 死→修 Path 2→恢復一致性」——不同問題。

**Capability 判準矩陣（0B-1 instrument）：**

| 問題 | Yes 代表 |
| --- | --- |
| 若 sink 消失，功能失效？ | dependency |
| 若資料缺失，只降低品質？ | observer |
| 是否產生使用者可見發現結果？ | discovery |
| 是否阻止錯誤進入系統？ | enforcement |
| 是否只描述/回報現況？ | observability |

**Resolution staging（不可跳序）：**
1. **0B-1 Capability Resolution** — 只填 dependency / observer / discovery / enforcement / observability。**不填** real_consumer / load_bearing。
2. **0B-2 Authority Resolution** — **目標是 resolve authority dimensions，不是 pick one owner**。authority 拆**兩軸**：
   - `authority_of_location`（誰決定東西**應放哪**？where）
   - `authority_of_permission`（誰決定東西**允不允許存在**？whether／schema／promotion）
   先假設正交；只有當**兩元件都宣稱同一軸**才叫 overlap。再補 `path_group`（Path1 / Path2 / independent）。
   migration 欄位拆開且 **`replacement ≠ owner`**（replacement 是遷移方向，不是治理來源）：
   `authority_source`（registry / rule / contract / local）+ `replacement_target`（遷移去向）。
3. **0B-3 Criticality Resolution** — 最後才下 real_consumer / load_bearing。三因子：
   **criticality = impact × authority × replaceability**（高 impact + 高 authority 但**可完全替代** → 未必 required）。
   新增欄位 `replaceability`（none / planned / immediate）。criticality enum：
   - `required` — 高 impact + authority + **不可替代**
   - `transitional` — 現在承載，但已有接替方向（如未來 Contract→Registry 接手）
   - `latent` — 能力存在但未承載
   - `dead` — 已失效
   - `out` — 非 consumer
   **Phase 0 Closure Rule**：*Every required consumer must have a declared authority source and a
   declared replacement path.*（required 卻無 replacement 很容易長成下一輪 drift。）
   **P0 blocker 規則**：若某 consumer 同時 `required` 且 `authority_source = ambiguous`，列為 **P0 blocker**，
   **不帶進實作**。

#### 0B-1 Capability Resolution（2026-06-24；capability only）

| consumer_id | 「sink 消失→誰做不成事」 | capability |
| --- | --- | --- |
| C-06 | 無人——`closeLoopGroupForPath` 只把 staged path top-dir 歸 commit owner-group；`feedback/` 是 ~13 prefix 之一，sink 消失即 no-op | **resolved OUT**（非 sink consumer，generic path classifier） |
| C-01 | （repoint 後）indexer 索引空 → 餵不到 Path 2 | dependency-provider（latent/目前 dead，globs 舊路徑） |
| C-04 | `runtime query` 對 lesson 回空 → 使用者可見 retrieval 能力下降 | discovery |
| C-03 | 不會失效——`MATCH 'feedback'` 仍命中 route/report/loop token；只回報現況且**目前說謊** | observability |
| C-09 | agent 經 `route.feedback.history` 直接讀的 primary_source 失效 → **Path 1 discovery 斷** | discovery（Path 1，目前真正 load-bearing） |
| C-10 | `route.feedback.promotion-pipeline` primary_source 失效 | discovery（Path 1） |
| C-11 | feedback-lessons 規則指向的寫入位置失效 → 規則 stale | enforcement（治理寫入位置） |
| C-16/19/20 | scenario 期望路由/寫入 sink → 行為錯會 fail | enforcement（validation guard） |
| C-18 | recovery reload set 少一個 pointer → 降低 recovery 品質、**不阻斷**（B-2 verified） | observer（非 dependency） |
| C-12/21 | 僅描述/單筆舉例，sink 消失不阻擋任何人 | observability / describe（候選 out） |
| C-07/C-08 | 與 sink 無關（token 同形、語意獨立） | OUT（name collision，非 consumer） |

> **0B-1 關鍵發現（capability 級，非 criticality 蓋章）**：impact lens 翻轉了 diagnosis 的重心——
> 真正「sink 消失會斷」的 discovery 是 **C-09 Path 1（registry direct-read）**；C-01→C-04（Path 2）
> 目前 dead。故 P0-B 修復的 consistency 目標是「讓已 dead 的 Path 2 重新與 load-bearing 的 Path 1 對齊」，
> 不是反向。C-06 / C-07 / C-08 三者由 impact lens 乾淨判出局，inventory 噪音下降。
>
> **未決（留 0B-2/0B-3）**：authority_source / path_group / replacement 尚未填；real_consumer /
> load_bearing 尚未蓋章。C-11/C-16/19/20 的 enforcement 是「治理寫入/驗證」而非「消費 lesson 內容」，
> 其 authority 歸屬待 0B-2 對照 Contract Ownership 鏈再定。

#### 0B-2 Authority Resolution（2026-06-24；兩軸，未蓋 criticality 章）

**C-09 vs C-11 正交假說檢驗結果**（讀 `metadata/rules/feedback-lessons.yaml` + `routing-registry.yaml:1951`）：

| 軸 | 主張者 | 判定 |
| --- | --- | --- |
| authority_of_permission（檔名/模板/schema/promotion/品質） | **C-11 rule 獨有** | ✅ 正交 |
| authority_of_location（lessons 放哪） | **C-09 registry + C-11 rule 都 textually hardcode `feedback/history/<domain>/`** | ❌ **真 overlap** |

> 結構訊號：`route.feedback.history.required_dependencies` 已指向 `enforcement/feedback-lessons.md`（registry→rule
> 有向引用），但兩者仍**各自 hardcode 路徑字串**——這就是 diagnosis「Source of Truth 未機械化」的具體形態。
> **收斂方向（Phase 4 contract）**：registry 為唯一 `authority_of_location`；rule 改為 **derive 路徑**、只留
> `authority_of_permission`。符合 Contract Ownership 鏈「path 唯一 owner = registry」。

| consumer_id | authority_of_location | authority_of_permission | path_group | authority_source | replacement_target |
| --- | --- | --- | --- | --- | --- |
| C-09 | **claims（canonical）** | — | Path1 | registry | —（已是 location owner / 遷移去向本身） |
| C-10 | claims（canonical） | — | Path1 | registry | — |
| C-11 | claims（**應撤**，改 derive） | **owns** | independent | rule（permission）/ registry（location） | location 主張 → derive from `contract.feedback.location`（待 Phase 4 命名） |
| C-01 | consumes（須由 registry 取得） | — | Path2 | registry | `route.feedback.history.primary_source` 提供的 sink glob |
| C-04 | consumes（透過 index） | — | Path2 | registry（間接） | 同 C-01（index 修好即恢復） |
| C-03 | — | — | independent | local（自證 token） | n/a（observability；P0-A 改 provenance 檢查） |
| C-16/19/20 | asserts（驗證期望） | — | independent（validation） | derived（from registry/contract） | 對齊 contract location 後重新 seed |
| C-18 | references（reload pointer，非 owner） | — | independent（recovery reload set） | reference-first（path 仍應 derive from registry） | sink README via registry（B-2: observer） |

> `replacement ≠ owner` 已遵守：`authority_source` 記治理來源，`replacement_target` 記遷移去向，兩欄不混。
> C-11 是唯一同時觸兩軸者——其 location 主張是 overlap 來源，permission 主張是它的正當核心。
> **未蓋 criticality 章**（real_consumer / load_bearing 留 0B-3：load-bearing = impact × authority）。

#### 0B-3 Criticality Resolution（2026-06-24；criticality = impact × authority × replaceability）

| consumer_id | impact | authority | replaceability | **criticality** | real_consumer |
| --- | --- | --- | --- | --- | --- |
| C-09 | high（sink 消失→Path1 斷） | registry（location，**但與 C-11 overlap**） | none（reference-first 為憲法預設，registry 穩定 owner） | **required** ⚠ | Y |
| C-10 | med-high（promotion route） | registry | none | required | Y |
| C-11 | high（治理寫入位置/schema） | **permission: owns（正當）**；location: 過度主張 | permission=none；location=planned（改 derive） | **required**（permission）⚠ | Y |
| C-01 | none（目前 globs 舊路徑） | registry（應 derive，現未） | planned（P0-B repoint） | **dead** | Y（待修） |
| C-04 | none（index 空→回空集） | registry（間接） | planned（C-01 修好即恢復） | **latent**（effectively dead until C-01 repoint） | Y（starved） |
| C-03 | none（不依賴 sink；只回報且說謊） | local（自證 token） | n/a | **out**（observability illusion；P0-A 改 provenance） | N |
| C-16/19/20 | high if enforcement active | derived（from registry/contract） | planned（對齊 contract 後 re-seed） | **required if enforcement active** | Y |
| C-18 | low（reload set 少一 pointer，品質降級不阻斷） | reference-first（無 Go loader；path 應 derive from registry） | planned | **observer**（B-2 verified；非 required） | N（非 sink dependency） |
| C-06/07/08 / C-12/21 | — | — | — | **out** | N |

**Phase 0 Closure — 觸發 2 個 blocker（不帶進實作）：**

- **B-1（P0 blocker · gates Phase 2 / P0-B）**：`authority_of_location` 目前**雙重擁有**（C-09 registry +
  C-11 rule 各 hardcode 路徑）。C-09 與 C-11 皆 `required` 且其依賴的 location authority = **ambiguous**
  → 命中 P0 blocker 規則。**若不先把 location authority 收斂成單一 owner（registry；rule 改 derive），
  P0-B 把 indexer repoint 只會再複製一份路徑字串 = 下一輪 drift。** 故 **Phase 4 contract 的 location
  收斂必須先於（或同步於）P0-B**，不可先 repoint。
- **B-2（verify · minor）— RESOLVED 2026-06-24**：C-18 無 Go loader（reference-first agent 讀）；其 feedback
  reference 只是 apk-analysis recovery reload set 的一個 pointer，sink 消失只降低 recovery 品質、不阻斷 →
  **C-18 = observer，非 required**。證明 `authority unknown + impact unknown` 收斂為 observer（不確定 ≠ 重要）。
  B-2 不新增 gate → 不改 sequencing。
  > **副帶範圍外觀察 — Observation Status: REFUTED 2026-06-24**
  > *Originally observed*：`metadata/recovery/domain-policies.yaml` 的 `source_path`
  > 與所有 inbound 引用均指向不存在的 `metadata/recovery/n.yaml`（rename 後 link 未更新）——獨立
  > reference-integrity drift，與 feedback sink 無關。
  > *Verification（current HEAD inspected）*：`source_path` 已是 `metadata/recovery/domain-policies.yaml`；
  > 全部 inbound 引用（routing-registry、model-checklists、escalation-levels、README、workflow-routing、
  > apk execution-flow、software-delivery intake、validation scenarios）皆解析到 `domain-policies.yaml`；
  > `rg '\bn\.yaml\b'` 全庫唯一命中即本段；`git log --all -- metadata/recovery/n.yaml` 為空（`n.yaml` 從未存在，
  > 無 rename 歷史）。平行的 `metadata/evidence/domain-policies.yaml` 同樣已收斂、無 drift。
  > *Disposition*：No action required。原 observation 為 **observation drift**（分析記錄脫離實際 repo 狀態，
  > 尚未進入 authority chain），非 integrity drift；標記為 refuted observation，excluded from Phase scope，不新增 blocker。

**Closure Rule 檢查**：required 集合 {C-09, C-10, C-11, C-16/19/20} 已全部列舉與 resolve。
inventory closure 的定義是 **「無未解 executable consumer」**，**不是「blocker 已清」**——故 B-1（C-09/C-11 的
location authority 尚未 single-owned）作為 **implementation gate 前移**交給 Phase 2，**不**阻擋 inventory 關閉。
正式 closure 由下方 Step 0C 的 Closure Statement 宣告。

### Phase 0 → Implementation：Phase Dependency（B-2 後正式重畫）

> **重畫前置已達**：C-18 resolve 為 observer，不新增 gate → C-18 未改 sequencing。
> 證據僅支持 `contract.location → gate P0-B`，**不**支持整個 Phase 4 上移。
> 守界：是 location slice gate，不是 `Phase 4 → Phase 2`。

```
Phase 1  P0-A observability            （無 path-authority 依賴 → 可先行/並行）
   │
   ▼
GATE   Location Authority Resolution    ← B-1（Phase 4 的 location slice）
   │     registry = 唯一 authority_of_location；C-11 rule 改 derive、只留 permission
   ▼
Phase 2  P0-B repair（indexer repoint，registry-derived）
   │
   ▼
Phase 3  P0-C alignment（fixtures / negative tests）
   │
   ▼
Phase 4  remaining contract extraction（location slice 已在 B-1 gate 先行）
```

- **phase dependency ≠ consumer dependency**：本圖是從目前已 resolve 的 consumer 推導，非結構先驗。
- 唯一新增的硬邊是 **B-1 location gate**；其餘維持 diagnosis 的 observability-before-repair 序。

### Step 0C — Classification（最後，長在證據上）

> **0C 角色限制：不新增事實。** 只做 `capability + authority + criticality → class`。
> 禁止：重新判 consumer、重新判 owner、新增 blocker。分類是**結論不是前提**。

**Class = 治理動作（非技術類型）；`required ≠ Class A`（gating 才是分水嶺）：**

| Class | 定義 | 下一步 |
| --- | --- | --- |
| A | required + actionable（無 blocker） | 必須進 phase |
| B | required but gated | 先解除 blocker |
| C | latent / dead | 保留，不先修 |
| D | out / observer | 結案 |

#### Deliverable 1 — Final Consumer Table

| consumer | capability | authority | criticality | class | replacement |
| --- | --- | --- | --- | --- | --- |
| C-09 | discovery (Path1) | location: registry（與 C-11 overlap） | required | **B**（gated by B-1） | — location owner；B-1 收斂 |
| C-10 | discovery (Path1) | location: registry | required | **B**（gated by B-1） | — |
| C-11 | enforcement (permission) | permission: rule owns；location: 過度主張 | required | **B**（gated by B-1） | location → derive from `contract.feedback.location` |
| C-16/19/20 | enforcement (guard) | derived (registry/contract) | required if enforcement active | **A**（驗證現行正確 Path1，未被 B-1 擋） | P0-C re-seed/align |
| C-01 | dependency-provider | registry（應 derive，現未） | dead | **C** | P0-B repoint：`route.feedback.history` primary_source sink glob |
| C-04 | discovery (Path2) | registry（間接） | latent | **C** | C-01 repoint 後恢復 |
| C-03 | observability | local（自證 token） | out | **D** | n/a（health-check → P0-A provenance） |
| C-06 | path classifier（commit group） | — | out | **D** | n/a |
| C-07/C-08 | name collision（≠ sink） | independent | out | **D** | n/a |
| C-12/C-21 | observability / describe | none | out | **D** | n/a |
| C-18 | observer（recovery reload pointer） | reference-first（path 應 derive from registry） | observer | **D** | sink README via registry（derive） |

> 註：C-03 作為 *consumer* 是 D（結案）；其 health-check 程式在 P0-A 以 observability 修復處理，與 consumer 分類無關。
> A/B 都含 required——差別只在是否被 B-1 gated（`required ≠ Class A`）。

#### Deliverable 2 — Inventory Closure Statement

> **No unresolved executable consumer remains.**
> （不是「everything is understood」——後者永遠關不掉。）

**→ Phase 0 正式關閉。** discover → resolve → classify 完成，inventory 從「工作中假設」轉為 **baseline**。
交付 implementation 的硬邊：**B-1（location authority 收斂）gate P0-B**。下游入口：Phase 1 (P0-A) ∥ B-1 location slice 可並行。

---

## Phase 1 — Observability Restoration (P0-A)

> **原則：Truth before repair.**
> observability 不是修復的前置步驟，而是**獨立可交付能力**——故獨立成 phase。
> 對應 baseline 的「P0-A 不保證 repair」：本 phase 成功 = 量測誠實，**不**等於 lesson 已索引。

**第一個 commit 成功條件（鎖死）**：*Health-check no longer mistakes route presence for lesson presence.*
（不是「恢復功能」。）

**交付**：
- [x] provenance health-check `nativeRuntimeIndexProvenanceCheck`（runtime.go）：count `type='feedback-pattern'
  AND source_path LIKE <sink>%`；route atom（type=route）被 type filter 排除＝anti-mask 內建。**未**改既有
  `nativeRuntimeIndexFTSCheck`（masking gate 仍在，等 Gate Promotion 才替換）。
- [x] `source_path` provenance assertion（≥1 lesson atom 在 canonical sink 下）。
- [x] sink registry-derived（`feedbackCanonicalSink` 讀 `route.feedback.history.primary_source`，**不硬編**）。

**Option A — reported-red disposition（不偷升級成 enforcement）**：
- check status=failed（provenance==0）但 **visible in refresh result**、**不參與 blocked aggregation**、**不改 exit code**。
- **Gate Promotion Rule**：reported-red → blocking-red 只允許在 `P0-B complete AND authority_of_location resolved` 後。
- 限制：不得降級成 warning、不得吞掉 failed、route token 不得通過。

**驗收（四測，已鎖；不測 repair——repair 屬 P0-B）**：
- [x] no provenance atoms → check fails（refresh 續行）
- [x] route atom only → check fails（anti-mask）
- [x] valid lesson atom → check passes
- [x] sink mismatch（舊路徑）→ check fails（drift resistance）
- 全綠：`go test ./internal/app -run TestNativeRuntimeIndexProvenanceCheck`；全模組 `go test ./...` 綠。

**Live 驗證（real repo）**：`runtime_index_feedback_provenance → failed`（"Path 2 empty"），但 **overall status=success /
exit 0** → truthful failure ≠ execution blockade。

**禁止**：
- [x] 未改 indexer（observability 與 repair 分離；indexer repoint 是 P0-B）。
- [x] 未動 Path 2 修復（紅是真實狀態，非偷修）。

### B-1 Design Slice — `contract.feedback.location` (minimal shape, design-only)

> **範圍限制**：只出 minimal shape。**不出** ownership migration、registry patch、validator（那些是 Phase 2/4）。
> 草稿只回答三題，其餘留空。**design-only，不 materialize 成 runtime contract 檔。**

```yaml
# DESIGN DRAFT (not a live contract). Materialization = Phase 4.
contract.feedback.location:
  # 1. canonical location 由誰宣告（owner）
  owner: knowledge/runtime/routing-registry.yaml#route.feedback.history.primary_source
         # registry 為唯一 authority_of_location（收斂 B-1 overlap：rule 不再 hardcode 路徑）

  # 2. consumer 怎麼取得（derive）
  derive: |
    consumers MUST read the sink prefix from the registry route (primary_source 的 dir),
    never hardcode "feedback/history/". 既有 derive 範例：P0-A 的 feedbackCanonicalSink()。

  # 3. owner 變更如何失效（invalidation）
  invalidation: |
    若 route.feedback.history.primary_source 改變或缺失，所有 derived consumer 必須重新解析；
    任何 hardcoded 路徑即違反 Canonical-Path Derivation Invariant，視為 drift。

  # 其餘（permission scope / schema / promotion / migration steps）先空——非 B-1 範圍。
```

> 註：P0-A 的 `feedbackCanonicalSink()` 已是此 derive rule 的第一個 consumer（read-only），
> 但 B-1 的「rule 撤回 location 主張」收斂仍屬 Phase 2 entry gate，未在本 slice 執行。

---

## Phase 2 — Runtime Repair (P0-B)

> **原則：Registry-derived indexing only.**
>
> **Entry gate（B-1，來自 Phase 0 0B-3）**：進入 P0-B 前，`authority_of_location` 必須已收斂成單一 owner
> （registry；C-11 rule 改 derive）。**否則 repoint indexer 只是再複製一份路徑 = 下一輪 drift。**
> 即 Phase 4 contract 的 location 收斂必須先於（或同步於）本 phase。

**交付**：
- [ ] index repoint：`skills/*/feedback_history` → canonical sink（路徑由 registry 提供）
- [ ] registry consumption：indexer 從 routing-registry 取 path，不 hardcode
- [ ] compatibility notes（舊路徑消費者/快取的過渡說明）

**成功條件**：
- [ ] `Path 1 == Path 2`（同一 lesson 集合）
- [ ] Phase 1 的 health-check 由紅轉綠

**禁止**：
- 禁止任何 hardcoded canonical path（違反 Canonical-Path Derivation Invariant）。

---

## Phase 3 — Reality Alignment (P0-C)

> **原則：Fixture must model reality.**

**交付**：
- [ ] fixture migration：`runtime_test.go` seed 改為 canonical sink
- [ ] negative tests：舊路徑**不**被索引
- [ ] deleted-world lock：機械化阻止 fixture 再教 deleted world

**成功條件**：
- [ ] 舊路徑（`skills/*/feedback_history`）若存在於 fixture/索引 → 測試 **fail**

---

## Phase 4 — Contract Extraction

把 baseline 的 Contract Ownership 鏈從文字落為可執行 artifact。

**交付**：
- [ ] `knowledge/runtime/contracts/<name>.yaml`（machine-readable contract）
- [ ] contract → routing-registry 的 single-direction materialize 連結（registry = 唯一 path owner）

**成功條件**：以下全部 **derive**（不得平行宣告 path）：
- [ ] runtime
- [ ] validator
- [ ] docs
- [ ] tests

---

## Rollback

只允許：
- [ ] disable enforcement（停用新 validator / health-check，回到觀測前狀態）

禁止：
- rollback constitution（`ADR-004` 是正確來源，回退它＝回到 drift 世界觀）

---

## Exit Criteria

**Invariant（任一條斷裂 → 不可 close）**：

```
Constitution
  ⊨ Contract
  ⊨ Registry
  ⊨ Runtime
  ⊨ Validation
```

- [ ] 五層 derive 鏈全綠
- [ ] Canonical-Path Derivation Invariant 無違反（無繞過 registry 的 path 宣告）
- [ ] `Path 1 == Path 2`
- [ ] baseline diagnosis 仍 frozen、未被本計畫修改

---

## Deferred Design Notes（記錄，禁止現在實作）

> 這些是本計畫討論中露出、但**刻意不現在做**的抽象。固定觸發條件，避免 (a) 在單一案例上過早抽象，
> (b) 兩天後忘記為何當初不做。每條都附「何時值得重啟」的 falsifiable 條件。

### N1 — Linkage 雙圖模型（已決定，僅記錄）
- `tree parent` = ownership / decomposition（execution graph）；`baseline_ref` = authority /
  derivation（reference graph）。兩者語義不同，**不可**收斂成單一 `parent`。
- 現況：`baseline_ref` 表達 *implementation derives from diagnosis*，非 *belongs to*。

### N2 — `relations:` 多型連結抽象（DEFER）
- 形狀（未來可能，**非現在**）：
  ```
  relations:
    - type: derives_from
      target: <id>
    - type: decomposes
      target: <id>
  ```
- **不現在做的理由**：目前只有一個 evidence→implementation 案例，不值得 enum surface。
- **重啟觸發（任一反覆出現 ≥2 次）**：`ADR → completion`、`spike → adoption`、
  `incident → remediation`、`RFC → rollout`。出現重複才升級成 typed relations。

### N3 — `baseline_ref` 完整性弱規則（DEFER）
- 規則（未來 completion validator，**非現在**）：
  `baseline_ref exists → referenced document status ∈ {frozen, accepted}`。
- 目的：防止 reference graph 退化成另一種 drift——completion 不得 reference 一份 mutable draft。
- 現況已部分人工滿足：本計畫的 baseline `status: frozen`。

### N4 — EL-3A（ECS plan）classified: no ADR-004 impact（append-only, 主線不動）
- 來源：`plans/active/2026-06-16-1131-evidence-candidate-system.md` Evidence Log `EL-3A`（2026-06-24，maintainer）——
  關於 Vidoe-Test perf-governance evidence system 的觀察，**與 feedback migration 無關**。
- 分類（read → classify，**未** merge）：對 ADR-004 **authority / consumer / sequencing 三者皆無影響** →
  consumer census、location/permission 兩軸、B-1 gate、phase 依賴全部不動。
- 方法論 analogy（**不 promote**）：EL-3A 的 `Observation→Registry→Executor→Validation+Rule` 4-step core 與本計畫
  Contract Ownership 鏈（`Constitution→Contract→Registry→Consumers`）形狀相近，但 `analogy ≠ same family`；
  EL-3A 自身在 ECS 仍 `status: unresolved / family 未定`。故僅記為 deferred，不抽跨域 invariant、不開 candidate。
- 處置：deferred。不重開 ADR-004 架構搜尋（P0-A 已 ship、Phase 0 已 closed）。

## Document TODO

| 項目 | 狀態 |
| --- | --- |
| Phase 0 分類定義（Class A/B/C/D 邊界） | pending（0C，等 census 穩） |
| 0A census 若成長超過約一屏 → 拆 companion evidence 檔（單一 census 尚不值得新檔/新 ref edge，見 N2） | trigger-watch |
| Out of Scope 補足 | pending |
| contract 命名（`<name>.yaml`） | pending |
| 各 phase 的 validator 實作細節 | pending（進 phase 時展開） |
