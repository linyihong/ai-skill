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
2. **0B-2 Authority Resolution** — 補 authority_source（constitution / contract / registry / local）、path_group（Path1 / Path2 / independent）、replacement。
3. **0B-3 Criticality Resolution** — 最後才下 real_consumer / load_bearing。因為 **load-bearing = impact × authority**，不是讀路徑就算。

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
| C-18 | recovery navigation source（README）缺失 → recovery 流程降級 | dependency（recovery nav） |
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

### Step 0C — Classification（最後，長在證據上）

等 0A census 穩、0B resolve 完才定 Class A/B/C/D；分類是**結論不是前提**。

**Gate（未過不得進 Phase 1）**：
- [ ] 無 Unknown Consumer（每個觸及 intent-equivalent feedback 的點都已 census + resolve）
- [ ] 每個 real_consumer 有 authority_source + owner

---

## Phase 1 — Observability Restoration (P0-A)

> **原則：Truth before repair.**
> observability 不是修復的前置步驟，而是**獨立可交付能力**——故獨立成 phase。
> 對應 baseline 的「P0-A 不保證 repair」：本 phase 成功 = 量測誠實，**不**等於 lesson 已索引。

**交付**：
- [ ] provenance health-check（取代 `runtime.go:1879` 的 `MATCH 'feedback'` token 檢查）
- [ ] `source_path` assertions（≥1 atom 且其 source_path 在 canonical lesson sink 之下）
- [ ] index coverage metrics（可量測 Path 2 的真實 lesson 覆蓋率）

**成功條件**：
- [ ] Path 2 壞 → health-check 必須**紅**（先證明它能偵測現況的空集合）
- [ ] Path 2 修 → health-check 必須**綠**

**禁止**：
- 不允許在本 phase 修改 indexer（observability 與 repair 必須可分離驗證）。

---

## Phase 2 — Runtime Repair (P0-B)

> **原則：Registry-derived indexing only.**

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

## Document TODO

| 項目 | 狀態 |
| --- | --- |
| Phase 0 分類定義（Class A/B/C/D 邊界） | pending（0C，等 census 穩） |
| 0A census 若成長超過約一屏 → 拆 companion evidence 檔（單一 census 尚不值得新檔/新 ref edge，見 N2） | trigger-watch |
| Out of Scope 補足 | pending |
| contract 命名（`<name>.yaml`） | pending |
| 各 phase 的 validator 實作細節 | pending（進 phase 時展開） |
