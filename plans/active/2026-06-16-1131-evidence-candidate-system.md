---
id: 2026-06-16-1131-evidence-candidate-system
plan_kind: main
status: draft
owner: linyihong
created: 2026-06-16
priority: P2
parent: null
required_for_completion: false
---

# Evidence Candidate System（跨 plan 候選證據層）

**Status**: `draft` — 設計記錄（design record），**尚未進 runtime，未寫任何 code/surface**
Owner: framework maintainer (linyihong)
**建立日期**：2026-06-16
**Priority**：**P2**（observation-stage；不阻擋其他工作）
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

## 完成條件

- [ ] Phase 0.5 schema freeze 完成（Q1 禁止 confidence ✓ / Q2 不 committed ✓ / Q3 scanner 觸發 / notify→gate 合併 ✓ / candidate 生命週期定義）
- [ ] 三個 plan 各長出 machine-readable `evidence-rule + acceptance-gate`（Phase 1 第一 artifact）
- [ ] candidate `inbox/` gitignored 建立；accept 寫回 plan evidence（無永久 committed 中間態）
- [ ] agent-invoked scanner v0 可掃當前 session diff 比對 criteria 產生 candidate（輸出 `criteria_hits[]`，無 confidence 數字）
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
- [ ] freeze candidate schema：`id` / `source{repo,artifact,commit}` / `matched_plans[]` / `criteria_hits[]` / `status{create|accepted|discarded|expired}`
- [ ] freeze `evidence_rule` schema **空殼**（只凍結形狀，**不填三個 plan 的內容** — 那是 Phase 1）：

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

  完成條件：**三個 plan 都能 parse 這個空殼**，不要求有實際 rule 內容。

完成條件（Phase 0.5 關門）：
- [ ] candidate schema + `evidence_rule` 空殼凍結並記錄
- [x] Q1 / Q2 / Q3 / notify→gate / lifecycle / 兩條 invariant（candidate 不指 candidate、scanner stateless）全拍板
- [ ] 空殼 parse 通過後才進 Phase 1（**不在 Phase 0.5 起草實際 evidence-rule 內容**）

## Phase 1: Evidence Candidate System（被動）

- [ ] 三個 plan 各加 machine-readable `evidence-rule + acceptance-gate`（collect / 命中 criterion / gate + notify projection）
- [ ] 建 candidate `inbox/`（gitignored）+ 在合適位置標明 observation-only + economics 三條規則
- [ ] 用 Phase 0.5 凍結的 candidate schema（含 `criteria_hits[]`，無 confidence）
- [ ] agent-invoked scanner v0：掃當前 session diff/artifact → 比對 criteria → 產生 candidate（不自動寫入 plan）
- [ ] 人工 Accept 流程：accept → 寫回該 plan evidence；discard → 丟棄（看過不採納）；expire → 過期（retention / plan_closed，未處理）；defer = 留在 inbox 待後續

規則（Phase 1 不可違反）：
- 不自動寫入 plan
- 不計數成「成熟」
- 不改任何 gate
- 不 reopen
- 只通知（且通知語意 = 可以 review，非已升級）

完成條件：
- [ ] candidate 能被產生、列出、人工處置；accept 後出現在對應 plan evidence
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
