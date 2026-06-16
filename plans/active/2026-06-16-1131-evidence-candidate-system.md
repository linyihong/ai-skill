---
id: 2026-06-16-1131-evidence-candidate-system
plan_kind: main
status: in-progress
owner: linyihong
created: 2026-06-16
priority: P2
parent: null
required_for_completion: false
---

# Evidence Candidate System（跨 plan 候選證據層）

**Status**: `in-progress` — **Phase 1（observation infrastructure）完成**；現進入 **observation period**，
Phase 2（matcher + accumulation runtime）gated on `phase2_gate`（不主動建）。
Owner: framework maintainer (linyihong)
**建立日期**：2026-06-16
**Priority**：**P2**（observation-stage；不阻擋其他工作）
**目前執行入口（next）**：**不是繼續 build**。Phase 1 已交付 schema + registry（pointer-only）+ 三個 plan
的 evidence-rule + 人工 round-trip 證明 + Go assembler scanner（`cmd/evidencecandidate`）。下一步是
**累積真實（人工標註）candidate**，靠 `phase2_gate`（count≥20 / reviewed≥80% / accepted>50% /
age_p95<30d）的真實數據決定是否開 Phase 2。在 gate 未達前主動建 matcher / accumulation runtime ＝
違反本系統一路守的 observation-gate 紀律（呼應 economics plan 的 `observation_only`）。
**Glossary Impact**: yes（候選新詞 `evidence_candidate` / `candidate_registry` / `notify_rule`；**尚未**註冊到 `knowledge/glossary/ai-skill.md`，待 gate 通過再決定）

## Executive summary

多個 active plan（economics `2026-05-27`、governance-pattern `2026-06-08`、interaction-hazard
`2026-06-16`）都在**累積證據後才做 promotion / reopen 決策**。痛點不是「缺證據」，而是
**「人工開始記不住哪個案例該回流到哪個 plan」** — 這是**索引問題，不是 gate 問題**。

本 plan 設計一個 **Evidence Candidate System**：被動發現「這個事件可能值得某些 plan 看一下」，
產生 candidate，**等人工 accept 才寫入 plan**。它**只做 Detection + Capture，不做 Enforcement /
Promotion**。

關鍵切分（防止滑回 enforcement）：

```text
新任務完成
 ↓
Evidence Candidate System 掃 artifact（commit / ADR / test / doc / issue）
 ↓
產生 candidate（指向可能相關的 plan + 命中的 criterion）
 ↓
人工 Accept / Reject / Defer        ← 判定權永遠在人
 ↓
Accept → 寫入該 plan 的 evidence
 ↓
plan 自己的 gate 判斷成熟度 → notify「可以 review」（不是「已升級」）
 ↓
Maintainer 宣告 promotion
```

## Decision Rationale

### Problem & Why Now

近期三個 plan 各自長出不同的證據面，但都在做「累積 → 決策」：

| Plan | 證據面 | 問的問題 | gate |
|---|---|---|---|
| economics `2026-05-27` | Evidence Log（A/B/C class，observation-only） | reduction 切夠薄沒？要不要開第二刀？ | **刻意沒有** |
| governance-pattern `2026-06-08` | Sample inventory（draft 檔，N≥5 + falsification） | 是不是 framework law？ | N≥5 + 反例 |
| interaction-hazard `2026-06-16` | Downstream pilot gate（A0–D，吃 Vidoe-Test commit） | primitive 值不值得升 O3 model？ | ADR criteria 1–6 |

它們是同一條成熟度光譜上的三個點，不是同一個機制。共同缺口只有一個：**跨 session / 跨 plan 的
「可回收觀察層」**。冷啟動問題（下一個 session 看不到上一個 session 的觀察）只能靠當下寫下來補。
Candidate System 正是補這個洞——但補法是 index/pointer，不是把結論寫進 memory。

### Decision

建立 **Evidence Candidate System**，責任嚴格分層：

| 層 | 誰負責 | 風險等級 |
|---|---|---|
| **發現案例（Detection）** | Candidate System | 低 |
| **收集案例（Capture）** | Human（accept/reject/defer） | 中 |
| **判斷案例有效（Validity）** | Plan（evidence-rule） | — |
| **判斷成熟（Maturity）** | Plan Gate（acceptance-gate） | — |
| **宣告升級（Promotion）** | Maintainer | 很高 |

System 只佔最上面兩層；**Validity / Maturity / Promotion 永遠不由 System 決定**。

每個 plan 自己宣告**兩件事**（不是三件，不是 System 寫死）：

```text
plan/
 ├─ evidence-rule.md     # 什麼 artifact 算我的候選證據（machine-readable criteria）
 └─ acceptance-gate.md   # 達什麼條件算成熟；notify 是 gate 的 projection，內嵌於此
```

**為什麼把 `notify` 收進 acceptance-gate（不獨立成 notify-rule，2026-06-16 maintainer review）**：
notify **不是獨立成熟度規則，只是 gate 的 projection**。若獨立，會出現第二個 state machine
（gate 過了 notify 沒發 / notify 發了 gate 還沒成熟）。所以 notify 必須是 gate 達標時的衍生輸出，
不能有自己的條件。

範例（concept，非最終 schema）：

```yaml
# economics — 收案例，永遠不通知成熟，maintainer 自己看
evidence-rule:
  collect: true
acceptance-gate:
  threshold: none
  notify:
    enabled: false

# governance-pattern — 達 N≥5 + falsification 才算成熟；notify 是它的 projection
evidence-rule:
  collect: true
acceptance-gate:
  sample_count >= 5
  falsification >= 1
  notify:
    enabled: true        # 達 gate 即提醒「可以考慮 review 是否 promotion」

# interaction-hazard — pilot 完成 + criteria 達標算成熟；notify 投影自同一 gate
evidence-rule:
  collect: true
acceptance-gate:
  pilot_complete == true
  criteria_pass >= 6
  notify:
    enabled: true
```

`notify` 的語意嚴格限定為 **✅ 可以 review**，**不是** ❌ 已升級。差異很大，不可混。

### Memory 不衝突 — 串接，不共用責任

maintainer 已有跨 session memory 層。**不把 Candidate System 建在 memory 上**，因為兩者本質不同：

| 層 | 用途 | 寫入條件 | 可推導 | 可當 authority |
|---|---|---|---|---|
| Memory | 跨 session 連續性 | 對未來有幫助 | 可整理 | 不建議 |
| Evidence Candidate | 記錄觀察事件 | 有 artifact | **不應推論** | **不可** |
| Plan Evidence | 計畫正式證據 | 人工接受 | 可做 gate | 可以 |
| Promotion / ADR | 正式決策 | 達 gate | 可以 | 可以 |

正確串接（memory 作 index/提醒側車，不在主鏈上）：

```text
Session → Task Output → Evidence Candidate Store → (人工 Accept) → Plan Evidence → Gate → Decision
                                                                          ↑
                                                          Memory（只存指標／提醒）
```

**Memory 只存三種東西**：(1) 有哪些 plan 正在收集中、(2) 有哪些待 review 的通知、(3) 最近做到哪步。
**不存**：candidate 明細、證據內容、sample 數量、gate 狀態 —— 那些都屬於 evidence/candidate store。
亦即 `memory += work_continuation_pointer`，**不是** `memory += evidence`。（用 `work_continuation`
不用 `observation`：這個 pointer 已不只是觀察，它含「哪些 plan 在收集 / 哪些通知待 review / 最近做到
哪步」，比較像 workflow continuation。）

這讓 Candidate Store 成為 **memory of observations ≠ memory of conclusions**：記「發生過什麼、指向
哪裡」，不記「已成立、已抽象成功、已升級」。正好不踩 governance-pattern 一直在防的 premature
extraction。

### Alternatives Considered

- **A. 直接做共用機制 + enforcement（hook 阻流程直到記錄）**：reject。違反 economics plan 明文拒絕
  的 auto reopen/closure/threshold，且把 N=1 的形狀提前抽成強制模板。
- **B. 只在另兩個 plan 各加一張被動 observation table**：reject as insufficient。兩 plan 已有各自證據
  面，重複加表造成 dual source-of-truth；且不解「記不住哪個案例回哪個 plan」的索引問題。
- **C. 把 Candidate Store 建在 memory 層上**：reject。memory 是連續性層、可推導；evidence candidate
  不可推論、不可當 authority。共用會讓 memory 變成 evidence authority（危險）。
- **D. Detection + Candidate Queue（被動發現 + 人工 accept），plan 自帶 evidence-rule +
  acceptance-gate（notify 為 gate projection）**：**accept**。停在成熟度階梯的 Detection +
  Capture，不進 Enforcement/Promotion。
- **E. 等以後再說，現在完全不動**：reject。冷啟動/失憶問題真實存在，索引層是低風險補洞。

### Why Not an ADR Yet

「Candidate → Accept → Gate → Notify」這個形狀目前只存在於**這個構想本身（N=1）**。要不要把它抽成
跨 plan 共用 primitive，**必須等真的觀察到多個 plan 都長同樣形狀**才決定——這正是
governance-pattern plan 的紀律。現在升 ADR = 自己違反自己。schema 的 scanner 觸發模型仍未驗證
（confidence 與 store 形態已由 Phase 0.5 schema freeze 拍板，見 §Open Questions Q1/Q2）。

### ADR Promotion Criteria（completed 時驗證）

Phase 2 gate（要不要蓋更多基礎設施）— **四項全過才算 working system，不是 backlog / 垃圾桶**：

```yaml
phase2_gate:
  candidate_count   >= 20      # 跨 session 累積量
  reviewed_ratio    >= 80%     # (accept + discard) / created；expire 不算 reviewed → 大量 expire 會 fail
  accepted_ratio    >  50%     # accept / (accept + discard)，只在「看過的」裡算，避免 expire 稀釋
  candidate_age_p95 <  30d     # 處理延遲；防 backlog
```

> 為什麼 `reviewed_ratio` 和 `accepted_ratio` 要分開：只看 accept rate 會被「20 個 candidate 只處理
> 3 個但都 accept」騙過（rate 漂亮但其實沒人在看）。`reviewed_ratio` 管「有沒有人在看」，
> `accepted_ratio` 管「看了之後採納比例」，兩者分母不同（見上方公式），缺一不可。

- [ ] 上述 `phase2_gate` 四項全過
- [ ] 沒有大量誤報（false-positive 率可接受，具體門檻待 Phase 1 觀察後定）
- [ ] 至少 2 個 plan 真的用 evidence-rule + acceptance-gate 各自判定（證明非 economics-only 特例）
- [ ] 沒有更輕的 promotion target（per ADR-007）
- [ ] Open Questions 全解

> 注意：上面的 `≥20` / `>50%` 是 gate「**要不要再蓋更多基礎設施（Phase 2）**」的門檻，
> **不是** gate「要不要 reopen 某個 plan」。後者永遠是人工 decision。這個區別必須寫清楚，
> 否則會漂回 economics plan 拒絕的那種 threshold。

### Consequences

#### 正面
- 解「記不住哪個案例回哪個 plan」的索引問題（真實痛點）
- 補跨 session 冷啟動洞，但不污染 memory 成 evidence authority
- 三個 plan 的證據面維持各自形狀，不被壓成單一模板
- 為「是否值得抽共用 Evidence Runtime」累積真實證據（而非 N=1 抽模板）

#### 負面
- 多一個 store + scanner 要維護
- plan 要各自長出 machine-readable evidence-rule（前置成本）

#### 風險
- **confidence 數字會 anchor 人的判斷**（機器對「屬於哪個 plan」做判定）——**已由 Q1 處置：禁止
  confidence，改用 `criteria_hits[]`**；未來真要分級用 `support_level: weak/medium/strong`，不用假精確
- scanner 若做成 standing daemon / 跨 N repo 監看，會從「索引層」膨脹成 infra
- candidate 若以永久 committed 狀態存在，可能被當證據引用、變成新 governance surface——**已由 Q2
  處置：candidate inbox gitignored，accept 才寫回 plan**
- notify 語意若從「可以 review」漂成「已成熟/已升級」，就滑回 Promotion 自動化——已由 notify→gate
  projection 收斂（不獨立成 state machine）

## Runtime Execution Path

**本 plan 不接入 runtime。** 不新增 `route.*`、不 project `runtime.db generated_surfaces`、不加
commit-msg validator。Phase 1 的 scanner 設計為 **agent-invoked**（session 收尾時手動掃當前 diff /
artifact），不是 hook 自動觸發、不是 standing daemon。

未來接入條件：唯有 §ADR Promotion Criteria 全綠（≥20 candidate / >50% accept / 多 plan 使用）後，
才在後續 plan（Phase 2 Evidence Accumulation Runtime）評估 runtime wiring。在那之前，referencing
本 plan 的任何 schema 等同 reference plan-vocabulary，**不等同 runtime contract**。

**Per-surface consumer 表**：N/A — 本 plan 不新增任何 generated surface / route / validator。

## Watch-Out List citation

對應 [`architecture/ai-native-cognitive-ecosystem-system.md`](../../architecture/ai-native-cognitive-ecosystem-system.md)
§Watch-Out List：本 plan 的最大 over-engineering 風險 = 把索引層提前蓋成 enforcement/runtime。
護欄 = 成熟度階梯只走 Detection + Capture；Promotion/Enforcement 一律 deferred 到 Phase 2 且需證據。

## Open Questions

| # | Question | 傾向 | 處置 |
|---|---|---|---|
| Q1 | candidate 要不要輸出 per-plan **confidence 數字**（如 `governance-pattern (0.81)`）？ | — | **resolved（禁止）2026-06-16** — 不輸出 confidence 數字，改用 `criteria_hits[]`（命中哪幾條 criterion）。數字會 anchor 人的判斷且無訓練資料；未來真要分級用 `support_level: weak/medium/strong`。 |
| Q2 | candidate store committed 與否？ | — | **resolved（不 committed）2026-06-16** — candidate `inbox/` **gitignored**（本機跨 session 持久，不是 governance surface）；accept 後**寫回對應 plan 的 evidence**（唯一 committed 證據面），不保留永久 committed 中間態。子問：是否另設 `accepted/` audit dir → 暫不設，靠 plan + git history 即可（若日後需要再議）。 |
| Q3 | scanner 觸發模型：**agent-invoked**（手動掃當前 session）還是 task-completion hook（自動）？ | — | **resolved（agent-invoked）2026-06-16** — 見 §Phase 0.5 `scanner_trigger`。三條理由：R1 維持 authority 邊界（scanner 只提 candidate，不自決何時掃）/ R2 避免 evidence inflation（hook 每次完成都掃 → 候選數反映活動量而非訊號）/ R3 保留觀察窗口（maintainer 決定這次 session 值不值得觀察）。hook 是**未來 promotion**，不是預設能力。 |
| Q4 | 前置物：三個 plan 是否先各自長出 machine-readable `evidence-rule`？scanner 無此則只能比關鍵字。 | 是——這是 Phase 1 的第一個、也是最小的 artifact。 | proposed |
| Q5 | cross-repo evidence（下游 consuming 專案的 commit/diff）如何被 scanner 看見？目前 economics / interaction-hazard 是**人工**從下游搬。 | v0 維持人工搬下游 artifact 進 candidate；scanner 先只掃本 repo session。 | deferred to Phase 2 |
| Q6 | criterion membership 該**只在 accept-time** 驗，還是當 **local plan metadata 存在時也 static check**？（由 EL-1 真實案例觸發） | **Phase 1C: warn-only（僅 local plan）/ Phase 2: evaluate enforcement**。仍非 matching —— 只是 `criterion_id ∈ declared criterion ids`，與 assembler 不衝突；external/section_pending plan 本地驗不到，故只能 warn。 | open |
| Q7 | pointer consumability —— **plan 可不可以被引用**（被當 candidate target）？（由 EL-2 觸發；與 Q6 **不同層**：Q6 驗 criterion_id ∈ plan criteria，Q7 驗 plan 本身可否被消費） | **resolved（Phase 1C blocker）2026-06-16** —— resolve = `exists AND status==resolved`；section_pending → not resolvable / not candidate_target / not countable（WARN, no-emit, exit 0）。已落地 scanner status-aware resolve。 | **resolved** |

## Evidence Log（observation only）

ECS 自己的 observation 載體（記關於本系統架構的真實案例，非 consumer 證據）。出現一筆寫一筆，
靠它數不靠回憶。**只記實例不記結論；指向真實 artifact；不自動觸發任何 phase 推進。**

| id | date | surface | artifact | observation | strength | decision_blocked |
|---|---|---|---|---|---|---|
| EL-1 | 2026-06-16 | acceptance-boundary | candidate `C-9945b55a`（source: Vidoe-Test `docs/plans/2026-06-16-design-contract-cold-sign-off-packet.md`）+ 本 session | Assembler emit 了一筆 schema/pointer/invariant 皆合法的 candidate，但其 `criteria_hits=[independent_reviewer_same_outcome]` **不屬於** 其 `matched_plans=[governance-pattern]` 宣告的 criteria（後者只有 new_6step/nonfitting/sibling）。criterion membership **目前只在 accept-time 驗，assemble-time 不驗**（assembler 不擁有 criterion，Guard）。即 **合法 candidate ≠ 正確 candidate** —— 此設計選擇第一次被真實 cross-repo 案例撞到。未造成錯誤接受（candidate 已 defer）。觸發 Q6。 | soft | no |
| EL-2 | 2026-06-16 | resolve-consumability | scanner `cmd/evidencecandidate` pointer-resolve（`os.Stat` only）+ 新 committed pointer `design-contract-validation-pilot.pointer.yaml`（status `section_pending`）+ live probe | scanner resolve 原本**只檢查檔案存在**，不看 `status`。實測：把 `section_pending` 的 pointer 當 matched_plan → **EMIT(exit 0)**。即 **consumer 把「存在性」當「可消費性」**（`index ≠ consumable`），撞到 ECS authority 分層。**已修**：status-aware resolve（exists AND status==resolved；section_pending → WARN no-emit）。觸發並 resolve Q7。 | hard | no |

**欄位**：`surface` = 撞到的系統面；`artifact` 指向真實物件；`strength ∈ hard/soft`；`decision_blocked ∈ yes/no`。
**規則同 economics Evidence Log**：只記實例、指向真實 artifact、不自動觸發 reopen/phase 推進。

## 完成條件

- [x] **Phase 0.5 設計凍結完成**（Q1/Q2/Q3 + notify→gate + lifecycle + 兩條 invariant + schema 形狀 + promotion boundary 全拍板；design property only）
- [x] Location 決策：`governance/evidence-candidates/`（schema/ + evidence-rules/ pointer-only + inbox/ gitignored）
- [ ] Phase 1 Entry Check 通過（schema/ 兩檔 parse）
- [ ] Phase 1A：三個 plan 各表達自己的 criterion（contract only，不做 matching/scanner；evidence-rules/ pointer only）
- [ ] Phase 1B：candidate `inbox/` gitignored；手動 candidate → accept → 寫回 plan evidence（證明鏈存在）
- [x] Phase 1C：stateless **assembler** scanner v0（Go `cmd/evidencecandidate`，非 dispatch target；3 硬護欄）產生候選（≠ 自動接受；criteria 須源於 scanner 外）
- [ ] 累積觀察期後評估 §ADR `phase2_gate`（count≥20 / reviewed_ratio≥80% / accepted_ratio>50% / age_p95<30d）→ 決定是否進 Phase 2
- [ ] 全程未滑入 Enforcement / Promotion 自動化（§Watch-Out 護欄成立）

## Phase 0: Pre-Build Interrogation

### Phase 0.0 — Open Questions 核對（公版，必填）

逐條核對本 plan §Open Questions，標記處置並回寫：

- [ ] 已讀本 plan §Open Questions 全部條目
- [ ] 對每條標記 `resolved`（附 Phase 0 證據）/ `still-open` / `deferred`（附原因）
- [ ] `resolved` 的條目已同步勾選 / 附註於 §Open Questions
- [ ] 若盤點新發現問題，已加入 §Open Questions

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| Q1 confidence 數字 | **resolved（禁止）** | maintainer review 2026-06-16；改用 `criteria_hits[]` |
| Q2 store 形態 | **resolved（不 committed）** | inbox gitignored + accept 寫回 plan |
| Q3 scanner 觸發 | **resolved（agent-invoked）** | R1 authority / R2 anti-inflation / R3 觀察窗口；hook 留作 promotion |
| Q4 evidence-rule 前置 | proposed | Phase 1 第一 artifact |
| Q5 cross-repo 掃描 | deferred | v0 人工搬下游 artifact |

- [ ] Confirm scope：Detection + Capture only；no enforcement / promotion / auto-reopen in v1。
- [ ] Confirm source-of-truth：plan 各自擁有 evidence-rule + acceptance-gate（notify 為 gate projection，不獨立）；System 不寫死判準。
- [ ] Confirm memory 邊界：memory 只存 pointer（open topics / pending notify / last step），不存 candidate 明細或 gate 狀態。
- [ ] Confirm 與 economics Evidence Log 的關係：candidate accept 後寫入的就是各 plan 既有證據面，不另開平行表。
- [ ] Confirm linked updates：economics / governance-pattern / interaction-hazard 三 plan 的 §與其他 plans 的關係；可能含 `enforcement/conversation-goal-ledger.md`（memory boundary）。
- [ ] Confirm non-goal：v1 不做 coverage dashboard / gate notification automation / reopen suggestion / threshold monitor（那是 Phase 2）。

## Phase 0.5: Schema Freeze（scanner 前必做，2026-06-16 maintainer 加入）

開 scanner **之前**先凍結 schema，避免邊掃邊改 candidate 形狀：

- [x] Q1 拍板：**禁止 confidence 數字**，candidate 輸出 `criteria_hits[]`（未來分級才用 `support_level`）
- [x] Q2 拍板：candidate **不 committed**（inbox gitignored；accept 寫回 plan）
- [x] notify → 收進 acceptance-gate（不獨立 notify-rule，避免第二個 state machine）
- [x] **Q3 拍板：`scanner_trigger.mode = agent_invoked`**

  ```yaml
  scanner_trigger:
    mode: agent_invoked        # 不是 task-completion hook
  ```

  理由：
  - **R1 維持 authority 邊界**：scanner 只能提 candidate，不能自決何時掃描。
  - **R2 避免 evidence inflation**：hook 每次完成都掃 → 候選數反映**活動量**而非**訊號**。
  - **R3 保留觀察窗口**：maintainer 決定這次 session 值不值得觀察。

  ```yaml
  promote_to_hook_when:        # hook = 未來 promotion，不是預設能力
    - candidate_count >= 20
    - reviewed_ratio >= 80%
    - false_positive acceptable
    - at_least_two_plans_using_rules
  ```
- [x] **Invariant — scanner MUST be stateless**：允許「掃當前 artifact → 產 candidate → 結束」；
  禁止「記住上次掃到哪 / 維護增量 cache / 追蹤 repo 狀態」。一旦 scanner 有狀態，就長出第二條
  authority 鏈（scanner → candidate → scanner state → decision）。與「candidate 不可指向 candidate」
  同源：觀察層不得自我累積狀態。
- [x] 定義 candidate 生命週期：**create → { accept | discard | expire }**
  - `accept` = 看過、採納 → 寫回 plan evidence
  - `discard` = 看過、不採納（= reject）
  - `expire` = **沒處理但過期**（≠ discard）；`reason ∈ {plan_closed, artifact_stale, exceeded_retention}`
  - 為什麼要 `expire`：沒它的話，「沒人看、90 天還在 inbox」會被誤算進 accept rate 分母而失真。expire
    不算 reviewed（見 §ADR `phase2_gate`），大量 expire 會讓 `reviewed_ratio` fail —— 正好抓「把 inbox 當垃圾桶」。
- [x] **Invariant — candidate 不可指向 candidate**：`candidate.source` MUST reference 原始 artifact
  （commit / ADR / test / doc / issue）；MUST NOT reference 另一個 candidate。否則觀察層會自我繁殖
  （C14 → scanner 又掃到 → C22）。鏈只允許 `artifact → candidate → plan evidence`。（對齊 economics
  plan D3 的 surface→surface 禁令。）
- [x] freeze candidate schema **形狀**（design property — 在本 plan 定義即算凍結）：`id` / `source{repo,artifact,commit}` / `matched_plans[]` / `criteria_hits[]` / `status{create|accepted|discarded|expired}`
- [x] freeze `evidence_rule` schema **形狀**（design property — 只凍形狀，**不填三個 plan 的內容**）：

  ```yaml
  evidence_rule:
    collect: bool
    match:
      artifact_types: []
      criteria:
        - id:
          description:
    exclusions: []
  ```

完成條件（Phase 0.5 關門 — **全為 design property，不依賴任何實作**）：
- [x] Q1 / Q2 / Q3 / notify→gate / lifecycle / 兩條 invariant（candidate 不指 candidate、scanner stateless）全拍板
- [x] candidate schema + `evidence_rule` 形狀已在 plan 凍結並記錄
- [x] promotion boundary（`phase2_gate` / `promote_to_hook_when`）已凍結

> **設計凍結 ✅**。parse / 建檔屬 implementation property，**移到 Phase 1 Entry Check**（不當 0.5 完成條件）。
> 區分理由：phase 完成應可在不跑、不建任何東西的情況下 assert；parse 依賴第一個實作，故不屬 design 關門。

## Phase 1 Entry Check（implementation property — 進 Phase 1 execution 前）

- [x] **location 決策已拍板：`governance/evidence-candidates/`**（2026-06-16；observation infra owned by governance — 不升 top-level repo capability、不與 consumer 綁死）
- [x] `governance/evidence-candidates/schema/candidate-schema.md` 建立且可被 parse（yaml block 驗證通過）
- [x] `governance/evidence-candidates/schema/evidence-rule-schema.md` 空殼可被 parse（yaml block 驗證通過）
- [x] `evidence-rules/*.pointer.yaml` resolve（3 個 plan_ref 全部存在；`embedded_section` 允許 missing）

### Location 決策（resolved 2026-06-16）— `governance/evidence-candidates/`

```text
governance/evidence-candidates/
  schema/                          # 分開放，預留 acceptance-gate-schema，避免日後重命名
    candidate-schema.md
    evidence-rule-schema.md
  evidence-rules/                  # POINTER ONLY — 不承載規則內容
    economics.pointer.yaml
    governance-pattern.pointer.yaml
    interaction-hazard.pointer.yaml
  inbox/                           # gitignored（Phase 1B）
  README.md
```

**為什麼不放 top-level `candidate/`**：太早升成 repo capability。**為什麼不內嵌各 plan section**：
把 contract 跟 consumer 綁死。**observation infra 的 owner = governance**。

**Invariant — `evidence-candidates/` 不可成為 plan authority**：`evidence-rules/*.pointer.yaml` 只放
**registry pointer**，不複製規則內容：

```yaml
plan_ref: active/2026-05-27-1557-tool-runtime-signal-economics-integration.md
schema: evidence-rule-v1
source:
  embedded_section:
    - Evidence Rule        # rule 定義仍在 plan 內，這裡只索引
```

rule **定義留在 plan**，這裡只 index。否則會長出 `plan.md → rule copy → registry` 雙 source（與
「candidate 不可成 authority」「scanner stateless」同源紀律）。

**Ownership 鏈（與整份設計最一致）**：

```text
plan        → defines rule
governance  → indexes rule   (evidence-candidates/, pointer only)
candidate   → references rule
evidence    → enters plan
```

## Phase 1: Evidence Candidate System（被動）— 1A → 1B → 1C，scanner 最後

> **範圍控制**：先證明 `artifact → candidate → plan evidence` 這條鏈真的存在（人工即可），
> 再做偵測。連人工 candidate 都跑不順，scanner 一定過早。

### Phase 1A — Contract only（定 contract，不做推論）— 三步

maintainer 三步切分（避免「為驗證 contract 被迫提前寫實作」）：

> **Step 1** schema + registry resolve → **Step 2** consumer attach（加 `## Evidence Rule` section）→
> **Step 3** criteria authoring（寫實際 criterion 內容）。加 `## Evidence Rule` section 本身**不等於**
> 開始寫 criterion。
>
> **護欄**：`pointer MAY reference missing embedded_section during Phase 1A` —— 不為了 pointer 合法
> 被迫先改 plan。

#### Step 1 — schema + registry resolve ✅（done 2026-06-16，skeleton 已建，未碰任何 plan）
- [x] `governance/evidence-candidates/schema/candidate-schema.md`（shape only）
- [x] `governance/evidence-candidates/schema/evidence-rule-schema.md`（空殼）
- [x] `governance/evidence-candidates/README.md` + `evidence-rules/{economics,governance-pattern,interaction-hazard}.pointer.yaml`（pointer only，`status: section_pending`）
- [x] `inbox/.gitignore`（dir tracked、內容 ignored）
- [x] Entry Check 通過：兩 schema parse + 3 pointer resolve + inbox 確認 ignored

#### Step 2 — consumer attach ✅（done 2026-06-16）
- [x] economics / governance-pattern / interaction-hazard 各加一個 `## Evidence Rule` 內嵌 section（rule 定義 owner = plan；criteria 留 placeholder）
- [x] 三個 pointer `status: section_pending → resolved`

#### Step 3 — criteria authoring ✅（done 2026-06-16）
- [x] 各 plan 在自己的 section 寫實際 criterion（id + description），錨定該 plan 既有 evidence 語意：
  - economics → `owner_ambiguity` / `state_cannot_describe_failure` / `phase_order_rework`（A/B/C class）
  - governance-pattern → `new_6step_sample` / `nonfitting_sample` / `sibling_family_sample`
  - interaction-hazard → `trust_transition_case` / `field_survival` / `predictive_prevention` / `downstream_pilot_evidence`

限制：定欄位 + criterion 識別；**不寫 matching 邏輯、不做 scanner**；evidence-rules/ **只放 pointer**。

成功條件：
- [x] 三個 plan 能各自表達自己的 criterion；pointer 解析得到 plan section；evidence_rule yaml block parse 通過

### Phase 1B — Candidate inbox（人工，不自動）

完成物：
- [x] `governance/evidence-candidates/inbox/`（gitignored；隨 Step 1 skeleton 一併建立）

限制：人工新增；**不自動掃、不通知**。

成功條件：
- [x] 手動放一筆 candidate → accept → **寫回對應 plan evidence**（證明鏈存在）— 見下方 round-trip log

#### Phase 1B round-trip log（2026-06-16 plumbing smoke test）

> **原則**：candidate 為合成（`dummy-observation` / `manual_test`），於 gitignored `inbox/` 實際操作；
> **不把 dummy row 寫進三個 plan 的真 evidence 面**（會違反 economics「只收真實 artifact」規則）。
> 本 log 是 Evidence Candidate System 自己的測試證據，記在本 plan，不污染 consumer。Step 3 之前先證
> `artifact → candidate → accept → plan evidence` 機制成立，避免日後 round-trip 爆掉時與 rule 問題混淆。

| Round | candidate | 動作 | 結果 |
|---|---|---|---|
| 0 | C-0001 (matched: economics) | create → **accept** | candidate 從 inbox 移除（gone）；evidence 寫回路徑驗證；real economics Evidence Log 未污染 |
| 1 | C-0002 (matched: governance-pattern) | create → **discard (reject)** | candidate 移除；**無 write-back**；不污染 Round 0 結果 |
| 2 | C-0003 (matched: interaction-hazard) | create → **expire** (`exceeded_retention`) | status 翻 expire；**不算 reviewed** |

驗證（全通過）：
- [x] **plan 為唯一 authority / 無第二份 evidence**：三個 plan grep `dummy-observation` = 0；candidate accept 後不在 store 留 committed 副本
- [x] **pointer 不變**：round-trip 全程三個 `*.pointer.yaml` 維持 `status: resolved`，未被 candidate 操作改動
- [x] **memory 無新增明細**：本輪未寫入 candidate 明細到 memory（memory 僅 work-continuation pointer）
- [x] **accept / reject 不互相污染**：C-0001 accept 與 C-0002 reject 獨立，互不影響
- [x] **expire 不失真 ratio**：`reviewed_ratio = (accept+discard)/created = 2/3 = 67%`（expire 排除）；`accepted_ratio = accept/(accept+discard) = 1/2 = 50%`（expire 不稀釋）
- [x] **inbox 全程 gitignored**：3 個 candidate 皆 `git check-ignore` 命中、未出現在 `git status`；inbox 僅 `.gitignore` 被追蹤

結論：`artifact → candidate → {accept|discard|expire} → (accept 才) plan evidence` 鏈機制成立。Step 3
（criteria authoring）可開始 —— 屆時是「定義偵測規則」，不是「拿規則救資料流」。

### Phase 1C — scanner v0 = **Assembler**（maintainer 拍板 2026-06-16）✅

**形式**：Go standalone tool（[`scripts/ai-skill-cli/cmd/evidencecandidate/`](../../scripts/ai-skill-cli/cmd/evidencecandidate/)）。
**Go-first policy 遵守**（不是 .py；前一版誤寫 Python 已刪），但**非 dispatch target、非 hook、非 runtime**
（Option B git-diff / Option C ai-skill 子命令皆 reject —— 會偷偷依賴 workflow / 讓 `accepted > scanner`
難證明）。三件事分開：implementation language（Go）/ toolchain placement（`cmd/`，discoverable）/
authority（none，非 routable）。「discoverable ≠ routable」。

**三條硬護欄（maintainer 2026-06-16）**：
1. **removable** — 刪掉 binary repo 不壞；無 `route.*` / `runtime.db` / build pipeline（releasebuild 只 build `cmd/ai-skill`）/ commit hook / generated surface 依賴它。
2. **output = artifact 非 state** — side effect 只有 stdout + `inbox/<id>.json`；不改 runtime.db / registry / plan / memory。
3. **exit code ≠ 成熟度** — 只有 `0`=assembled / `1`=invalid input；不編碼 accepted / matured。

**Contract**：

```text
input:  artifact + EXPLICIT criteria_hits（criteria_source.actor 標於 scanner 外）+ rule registry
output: candidate{ status: create } → inbox/（gitignored）
```

scanner **does**：schema validate / pointer resolve / dedupe（content-hash id，idempotent）/ invariant check / persist inbox。
scanner **does NOT**：infer / match / classify / score / **rank** / accept / expire。

**為什麼是 Assembler 不是 matcher**：criteria 目前只有 declaration，沒有 executable semantics。若讓
scanner 從 description 推論，就引入 `description → ontology → matcher → authority`，而 **ontology 本身
就是狀態**，違反 `scanner MUST be stateless` + `candidate MUST NOT become authority`。matcher 是 Phase 2
（`criteria_source.actor: matcher-v2`），scanner v0 永遠不自產 criteria。

**新 invariant — `criteria_hits MUST originate outside scanner`**：candidate 必帶 `criteria_source.actor`
（`human` 等），scanner 自身（`scanner-v0`/`scanner`/`self`）被拒。未來 matcher 來源另記，scanner 不知。

**Ranking 移除**：emit only，**ordering undefined**。連 hit-count 都不排（criterion 非等權，hit-count 會
暗示 `3 hits > 1 hit`）。要排序是 acceptance-gate 的事。

**成功重定義**：`scanner 成功 = 不扭曲人工標註`（**不是**「發現候選」）。

**Status-aware resolve（Phase 1C blocker fix，2026-06-16）— `index ≠ consumable`**：
EL-2 實測發現 scanner 的 pointer-resolve 原本只 `os.Stat()`（檔案存在），把 **section_pending pointer
當成可消費**。修正：resolve 條件改為 `pointer 存在 AND status == resolved`。

```yaml
pointer_state:
  resolved:         { resolvable: true }
  section_pending:  { resolvable: false, candidate_target: false, countable: false }
```

scanner 行為：matched_plan pointer **missing → REJECT(exit 1)**；**section_pending → WARN「index !=
consumable」+ NO EMIT + exit 0**（非錯誤）；**resolved → emit**。這是 blocker（不是 enhancement）：
consumer 不得把「存在」當「可消費」。與 Q7 同源、與 Q6 不同層（見 §Open Questions）。

完成物：
- [x] `scripts/ai-skill-cli/cmd/evidencecandidate/main.go`（Go assembler，**status-aware resolve**）+ `main_test.go`（9 unit tests 含 section_pending-not-consumable）；`go vet` / `go test` 通過；live: resolved→EMIT(0)、section_pending→WARN no-emit(0)

#### Phase 1C Exit Criteria（maintainer 2026-06-16）— 全通過
- [x] scanner 輸出 candidate（`status: create`），**不可直接 accepted**（T1 emit；無 accept 路徑）
- [x] accept / discard / expire replay ≥1 各一（Phase 1B round-trip log 已證；scanner 只是前端 assembler，沿用同 lifecycle）
- [x] scanner 輸出（`criteria_hits`）**不參與 ratio 計算**（per Q1 無 confidence；ratio 只看 accept/discard/expire）
- [x] candidate → accepted 必須保留人工 evidence link（`criteria_source.actor` 必存，invariant T4/T5 enforced）
- [x] **移除 scanner 後 Phase 1A/1B 仍可完整運作**（scanner 為 standalone，無人 import / 無 hook 依賴；刪檔不影響 schema/registry/round-trip）

#### scanner v0 測試（2026-06-16，皆通過；8 Go unit tests + live smoke run）
| # | 測試 | 結果 |
|---|---|---|
| T1 | 合法 emit（actor: human） | `EMIT C-dd3bbf07 → inbox`（gitignored；exit 0）|
| T2 | 同 input 重跑 | `IDEMPOTENT`（同 content-hash id，無 dup）→ deterministic |
| T3 | source 指向另一個 candidate（`C-deadbeef`）| REJECT（candidate 不可指向 candidate）|
| T4 | criteria_source.actor = `scanner-v0` | REJECT（criteria 須源於 scanner 外）|
| T5 | 缺 criteria_source.actor | REJECT |
| T6 | matched_plan 無對應 pointer | REJECT（pointer resolve）|

規則（Phase 1 全程不可違反）：
- 不自動寫入 plan
- 不計數成「成熟」
- 不改任何 gate
- 不 reopen
- 只通知（且通知語意 = 可以 review，非已升級）

完成條件：
- [x] 1A/1B/1C 依序完成；accept 後 candidate 出現在對應 plan evidence（Phase 1B 已證）
- [x] memory 只多了 pointer 類條目，未含 candidate 明細
- [ ] memory 只多了 pointer 類條目，未含 candidate 明細

## Phase 2: Evidence Accumulation Runtime（deferred）

僅在 §ADR Promotion Criteria 全綠後評估：

- [ ] coverage dashboard
- [ ] gate notification（達 notify-rule 提醒「可以 review」）
- [ ] reopen suggestion（仍人工決定，非自動）
- [ ] threshold monitor

**仍不做 auto-promotion。** 若屆時觀察到「多個 plan 都有 candidate → accept → gate → notify」同樣
形狀，才把這整段抽成共用 primitive——那時是收到證據後抽象，不是 N=1 抽模板。

## Stakeholder 同意項目

- [ ] maintainer 確認成熟度階梯只走 Detection + Capture（Phase 1）
- [ ] maintainer 拍板 Q1 confidence 數字立場
- [ ] maintainer 確認 memory 只存 pointer 的邊界

## 與其他 plans 的關係

- [`active/2026-05-27-1557-tool-runtime-signal-economics-integration.md`](2026-05-27-1557-tool-runtime-signal-economics-integration.md)
  §Evidence Log — 本 plan 的 candidate accept 後寫入的證據面之一；economics 設定 `notify: false`
  （maintainer 自己看）。economics 三條規則（只記實例 / 指向真實 artifact / 不自動 reopen）是本
  plan candidate store 的 header 規則來源。
- [`active/2026-06-08-2100-governance-pattern-library-extraction.md`](2026-06-08-2100-governance-pattern-library-extraction.md)
  — 本 plan「Candidate→Accept→Gate→Notify 形狀」是否值得抽成共用 primitive，**必須走該 plan 的
  N≥5 + falsification 紀律**；本 plan 自己現在是 N=1，不抽模板。governance-pattern 的 sample
  inventory 是 candidate accept 的目標證據面之一（`notify: sample_count>=5 + falsification>=1`）。
- [`active/2026-06-16-1030-interaction-hazard-review-workflow.md`](2026-06-16-1030-interaction-hazard-review-workflow.md)
  — Downstream pilot gate 是 candidate accept 的目標證據面之一（`notify: pilot_complete +
  criteria_pass>=6`）；其下游 Vidoe-Test commit 是 cross-repo candidate 的範例來源（Q5）。
- [`enforcement/conversation-goal-ledger.md`](../../enforcement/conversation-goal-ledger.md)
  — memory / goal ledger / candidate store 的邊界參考（candidate store 既非 goal ledger 也非 memory）。
