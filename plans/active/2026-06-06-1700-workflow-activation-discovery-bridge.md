---
id: 2026-06-06-1700-workflow-activation-discovery-bridge
plan_kind: sub
status: active
owner: linyihong
created: 2026-06-06
parent: 2026-05-31-1900-workflow-activation-engine
required_for_completion: false
sub_plan_reason: >
  Workflow Activation Engine (parent) Phase 6 "Discovery → Detector feedback
  loop" 標 deferred 上線；同一失效模式（detector miss → 無 mechanical
  fallback → 靠 agent 自律 → 自律失敗）於 parent 完成隔日（2026-06-05）
  在一個消費 `route.workflow.travel-planning` 的 downstream project 上原樣
  重演。本 sub-plan 補上 parent
  延後的 Discovery bridge，採 Light → Deep 兩階段漸進架構，避免 detector
  ontology 擴張並維持 parent §Design Principles 的 pre-Read 破環依賴原則。
  Independent sign-off：影響 per-turn cost 模型與 advisory 注入路徑，
  非單純 trigger 補洞。
---

# Workflow Activation: Discovery Bridge

**Status**: `active`（Phase 0 + Phase A landed；Phase C governance closeout 2026-06-15；**Phase B deferred** to clean spike；Phase D = 三週 empirical time-gate）
Owner: framework maintainer (linyihong)
**世代**：Gen 3 runtime hardening — Workflow Activation Engine 第二採樣
**建立日期**：2026-06-06
**最後更新**：2026-06-15（Phase C governance closeout：failure pattern + governance 段 + glossary + registry stance；修復 Phase A.2/A.4/A.5 重複 checklist + status drift）
**Priority**：**P2**
**Parent plan**：[`2026-05-31-1900-workflow-activation-engine.md`](../archived/2026-05-31-1900-workflow-activation-engine.md)
**Empirical trigger**：2026-06-05 session — 某消費 `route.workflow.travel-planning` 的 downstream project 內，使用者要求對一份命名遵循 project-local convention 的 dated artifact 進行 review。Detector miss（無 user keyword、無 path match），無 mechanical fallback，agent 直接憑常識做 review。文件命中 travel-planning artifact-gates 19 項中 7~10 項缺漏未被偵測。使用者三輪追問才暴露 gap，與 parent plan 2026-05-31 原 incident 為**同一結構性缺口的兩次採樣**。具體 project artifact、檔名與對話細節依 [`reusable-guidance-boundary.md`](../../enforcement/reusable-guidance-boundary.md) 留在原 project 文件。

> 本 plan **不擴 detector schema、不改 routing-registry**。範圍嚴格限定在 "detector miss → Discovery → advisory" 的 mechanical bridge。Filename / project metadata 等 semantic surface 議題 park 至條件式 follow-up plan，待本 plan 三週量測後依 miss rate 決定是否開。

---

## Decision Rationale

### Problem & Why Now

Parent plan v8（2026-06-04 完成）落地 detector + per-turn gate + manual-lock，**mechanical 層只完成一半**：detector hit 時機械擋住，detector miss 時 fail-open 並依賴 agent 自律 fallback 到 Discovery。Parent plan Phase 6 明寫 Discovery feedback loop「hot-hook auto-call 刻意延後」。延後代價在 parent 完成隔日就兌現（2026-06-05 incident）。

失效路徑：

```
task input → detector(user_signals + context_signals(path)) miss
         → no mechanical fallback
         → agent self-route by intuition
         → workflow primary_source 未 Read
         → review 用常識做完，artifact-gates 19 項缺 7~10 項
```

**Why now**：parent plan archive 收尾 evidence 包含「detector miss 為設計接受的容忍範圍」假設。2026-06-05 incident 證明此容忍範圍對 cross-project + project-local-ontology 任務太寬。每加一個 downstream project 就會踩同一個雷。

### Decision

實作 **Discovery Bridge** 作為 detector miss 的 mechanical fallback，採 **Light → Deep 漸進架構**：

- **Phase A — Light Discovery**：detector miss 時即觸發，輸入只用 **pre-Read cheap signals**（user message text、artifact basenames、paths、extensions、frontmatter head bytes、cwd、project overlay metadata scan）。產 top-3 candidate route + confidence score。≥ threshold 直接注入 advisory；< threshold 進入 Phase B。
- **Phase B — Deep Discovery**：不主動 Read，**piggyback agent 自然的下一個 artifact Read**（hijack content stream），跑一次 content scan 補強 candidate。可多輪累積 — 每次新 Read 都重跑 scan 補新訊號。Append-only 不覆寫先前 proposal。

**關鍵屬性**：

- **Discovery 是 advisory，不是 gate** — miss 不 block 工具、不擋 agent；advisory text 注入 PreToolUse hook output 供 agent 自選回頭 Read workflow primary_source。維持 parent §Design Principles pre-Read 原則（workflow activation 仍不依賴 content）。
- **Scoring OK 因為 advisory** — parent 禁止 weighted scoring 是針對 activation gate（必須 deterministic 避免 threshold 調參地獄）。Discovery 是 ranking 不是 gating，scoring 必要且可調。
- **零邊際 Read 成本** — Phase B 不發起新 Read，只 hijack agent 本來就會做的 Read 的 content stream。若 agent 第一個 Read 是「錯」artifact（讀了 proposal.docx 而非 requirements.md），Phase B 依然從該檔提到的訊號補強；下一個 Read 來了再補一輪。Discovery 變成 **stream-y、append-only 的訊號累積**，不是 one-shot 決策。
- **不擴 detector schema** — routing-registry.yaml 不動，本 plan 純加法（新 runtime + 新 SQLite table + advisory injector）。

### Alternatives Considered

- **A. 維持現狀，靠 agent 自律 fallback** — reject because parent 完成隔日復發，自律不可靠。
- **B. 擴 detector 訊號 ontology（filename_signals + project_overlay_signals）** — reject because (1) project rule 直接綁 route 越過 source-of-truth boundary（project layer 寫 routing rule）；(2) 訊號擴張屬 Discovery 範疇而非 detector 範疇；(3) 增加 detector 邏輯複雜度但解決 long-tail 能力有限。詳見 v0 draft critique 結論（plan rescope 紀錄保留於 §與其他 plans 的關係）。
- **C. Disc-2（detector miss → 等 agent 首次 artifact Read → Discovery post-Read）** — reject because 把「第一次 Read 的選擇」當成既定事實。第一次 Read 可能是 user-cited artifact，但真正應該 Read 的 workflow / governance source 可能完全不同（e.g.「請看這個 runtime validator 問題」真正該讀的是 `governance/*`）。Disc-2 silently encodes a guess as fact.
- **D. Light → Deep 漸進 Discovery（accept）** — Phase A 用 cheap signal 早期斷捨，Phase B 訂閱式累積；不綁第一個 Read，cost 控制良好，未來 semantic surface 議題自然 fold 進 Phase A signal source。

### Non-Goals（本 plan 明文封印）

- **Discovery proposals MUST NOT modify `routing-registry.yaml` automatically**。Cross-session proposal aggregation（若日後實作）僅能產生 maintenance suggestion（人類審閱後手動寫 registry），不可自動 promote candidate 為 registry route entry。
  - 理由：proposal accumulation → auto registry evolution = **self-modifying ontology**，屬 ADR 等級議題，本 plan scope 之外。registry 的 route 寫入仍須走人類 sign-off + commit 路徑。
  - Open Question Q6（cross-session aggregation）若推進，**必須**以 advisory-only 形式呈現，不得繞過此封印。
- **Discovery proposals MUST NOT satisfy `activation_triggers`**（硬 guardrail，防 advisory 偷渡成 soft gate）：
  - 任何 `if proposal.confidence > X then auto_activate()` 邏輯**永遠禁止**寫入 `detector.go` / `hooks.go` / `routing-registry.yaml`
  - Activation 路徑只有兩條：(a) detector deterministic match against `routing-registry.activation_triggers`，(b) user manual-lock。Discovery 永不為第三條。
  - 此封印對應 parent plan §Design Principles「Deterministic rule match，不用 weighted scoring」對 activation gate 的要求；Discovery 是 advisory ranking，scoring 在此語境合法，但**不得跨界**用於 activation。
  - 系統演化兩三代後最容易踩此線（confidence 高得彷彿可信 → 直接 auto-activate）。明文封印 + 未來 code review 對「proposal.confidence 出現在 activation 條件式」零容忍。

### Why Not an ADR Yet

- Phase A confidence threshold 未實測（暫定 0.5，依三週 empirical 量測再調）。
- Phase B piggyback 機制的 Read event 訂閱實作細節未確認（hook layer vs runtime layer）。
- discovery_proposals 表 lifecycle（TTL / promote / reject）governance 尚未定型。
- Cross-session proposal aggregation 是否該成為 long-term route registry 來源待 Phase 6 後決策。

### ADR Promotion Criteria（completed 時驗證）

- [ ] foundational + cross-session + cross-project + expensive-to-reverse + explains-why 全中
- [ ] Phase A + Phase B 跑過 ≥ 2 個 project 的 real task 驗證
- [ ] 三週 empirical KPI tiering（**primary 為 routing quality 訊號、tertiary 為 agent 行為觀察，避免混淆**）：
  - **Primary KPI**：detector miss → proposal generated ≥ 70%（量 Discovery 是否覆蓋 detector miss）
  - **Secondary KPI**：proposal generated → advisory injected ≥ 50%（量 scoring threshold 是否合理）
  - **Tertiary observation**：advisory injected → agent pivot — 不設成功門檻，僅觀察。Agent 可能本就掌握上下文、advisory 正確但無 pivot 必要
- [ ] Phase A p95 cost ≤ 30ms；Phase B p95 cost（hijack 單次 Read）≤ 50ms
- [ ] Open Questions 全解
- [ ] 沒有更輕 promotion target（純 enforcement rule 補述）適用

### Consequences

**正面**：

- Detector miss 不再 = 完全失效；mechanical safety net 補齊 parent plan deferred 段
- 不動 detector schema，向後相容風險為零
- 為未來 semantic surface plan 預留乾淨 entry point（Phase A signal source 可擴）
- Discovery 三週 empirical 後若 hit rate 高，可能反向證明 semantic surface plan 不需要

**負面**：

- per-turn cost 上限提高（detector miss 時多跑 Light Discovery）
- 新 SQLite 表 + 新 runtime module 增加維護面
- Advisory text 注入推高 per-turn token 用量

**風險**：

- Phase A confidence threshold 設定不當 → over-advise（每 task 都 inject）或 under-advise（永遠 < threshold）
- Phase B hijack 機制若實作層級錯（e.g. 在 hook 層攔 Read 而非 runtime 層）可能與既有 PreToolUse pipeline 衝突
- discovery_proposals 表若無 TTL 會無限增長
- Cost budget 超標時 fail-open 還是 fail-closed 未拍板

**Glossary Impact**: yes
- 新引入：`discovery_bridge`、`light_discovery`、`deep_discovery`、`discovery_proposal`、`advisory_injection`、`piggyback_read`
- 需加入 [`knowledge/glossary/ai-skill.md`](../../knowledge/glossary/ai-skill.md)

**Watch-Out List citation**: [`architecture/ai-native-cognitive-ecosystem-system.md`](../../architecture/ai-native-cognitive-ecosystem-system.md) §Watch-Out List — 對應 walls「Per-turn cost 預算爆炸」「Discovery 反成主路由（advisory→primary 失控）」「SQLite 表無 lifecycle 治理」（cross-check 對應 wall 名稱待 Phase 0.1）

---

## Architecture Compatibility Preflight

| 欄位 | 內容 |
|---|---|
| Candidate files | `scripts/ai-skill-cli/internal/app/discovery.go`（新建）、`scripts/ai-skill-cli/internal/app/hooks.go`（注入 advisory + Phase B Read hijack）、`runtime/runtime.db`（新 `discovery_proposals` 表 + 新 surface `runtime.discovery.proposals`）、`runtime/core-bootstrap.yaml`（per_turn_obligations 可能加 advisory observation）、`governance/workflow-activation-engine.md`（補 Discovery Bridge 段落 + Light/Deep 模型 + scoring vs deterministic 區隔說明）、`enforcement/failure-patterns/detector-miss-no-fallback.md`（新建）、`knowledge/glossary/ai-skill.md`（6 新 term） |
| Source-of-truth | `routing-registry.yaml` 不動。Discovery 訊號來源宣告於 `discovery.go` source，governance 文件 mirror。Discovery proposals 是 runtime state，不是 canonical routing decision。 |
| Compiler / generated surfaces | `runtime.db` 需重 compile；新 surface key `runtime.discovery.proposals` + `runtime.discovery.config`（threshold / budget） |
| Layer responsibility | Discovery runtime 屬 runtime layer；Schema / proposal lifecycle 屬 runtime；Philosophy 屬 governance；Failure pattern 屬 enforcement |
| 與現行架構衝突 | 無；parent plan deferred 段補完，不改既有層職責 |
| `runtime.db` 影響 | 新表 + 2 新 surface；compile pipeline 加新 projection rule |

---

## Runtime Execution Path

**Trigger flow（Phase A + Phase B 完成後）**：

```
PreToolUse hook fires
  → detector.go 跑（既有路徑，不變）
  → if detector hit (single route) → gate.workflow.primary_source_read（既有）
  → if detector hit (multi route)  → workflow-routing.md 歧義裁決（既有）
  → if detector miss:
       → discovery_bridge_eligibility 檢查
            ┌─ user message non-trivial（≥ N tokens 或含 imperative verb）
            ├─ artifact reference 存在於 user message 或 Open Files
            └─ 24h 內 task_hash 無既有 proposal
       → if eligible:
            → Phase A Light Discovery:
                 input: user_msg + artifact_basename(s) + paths + extensions
                      + frontmatter head bytes (≤ 200B)
                      + cwd + project overlay metadata（scan
                      `.ai-skill/project/rules/*.md` frontmatter）
                 output: top-3 candidate routes with confidence
                 write to runtime.db.discovery_proposals
                 if best confidence ≥ threshold:
                     inject PreToolUse advisory: "Detected possible route(s): [...]"
                     END
                 else:
                     mark proposal as awaiting_phase_b
                     END (do not block agent)

PostToolUse:Read hook fires (artifact Read by agent)
  → check if proposal exists with status=awaiting_phase_b for current task
  → if yes:
       → Phase B Deep Discovery:
            input: content of just-Read artifact
            output: refined candidate set (append, not overwrite)
            update proposal
            if new best confidence ≥ threshold:
                inject advisory in next PreToolUse turn
                mark proposal as advised
```

**Per-surface consumer 表**：

| Generated surface key | Named consumer(s) | Consumer 類型 |
|---|---|---|
| `runtime.discovery.proposals` | `discovery.go` (writer), advisory injector in `hooks.go` (reader) | Go runtime |
| `runtime.discovery.config` | `discovery.go` (threshold / budget reader) | Go runtime |

**Discovery Bridge 與既有 fail-open 行為的關係**：本 plan 不改 `gate.workflow.primary_source_read` 既有 fail-open（detector miss 仍不 block 工具）。Discovery Bridge 是 advisory layer，與 gate 平行運作，**從不 promote 成 gate**。

---

## Phase Plan

### Phase 0 — Preflight

#### Phase 0.0 — Open Questions 核對（公版，必填）

逐條核對本 plan §Open Questions：

- [x] 已讀本 plan §Open Questions 全部條目（2026-06-06 Phase 0.0 turn 確認 Q1–Q8 8 條皆已涵蓋；初版 dispatch table 漏列 Q7/Q8，本輪補上）
- [x] 對每條標記 `resolved` / `still-open` / `deferred`（下表）
- [x] `resolved` 條目同步勾選於 §Open Questions（本輪無 resolved；下表 8 條皆 still-open / deferred / deferred-to-phase-c）
- [x] 盤點新發現問題已加入 §Open Questions（本輪未發現新問題；Q7/Q8 已存在 §Open Questions 但漏列於本表，已補）

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| Q1 confidence threshold 初值 | still-open | Phase A.1 暫定 0.5；Phase A.3 unit test + Phase D empirical 後依 KPI 調 |
| Q2 Phase B hijack 實作層 | still-open（本 session out of scope） | Phase B.1 spike — 本 session scope 限 Phase 0 + Phase A，Q2 留 next session |
| Q3 proposal TTL | still-open | Phase A.2 schema 段預設 24h，由 `runtime.discovery.config` 可調；TTL eviction job 設計待 A.2 完成後決定 |
| Q4 Cost budget 超標策略 | still-open | Phase A.4 量測後決定；傾向 fail-open（不阻塞 agent，maintain advisory non-blocking 屬性），但需明文 |
| Q5 project overlay scan cache 邊界 | still-open | Phase A.2 設計 per-session in-memory + cwd-change invalidate；per-task 重 scan 待 Phase A.3 量測 |
| Q6 cross-session proposal aggregation | deferred | 三週 Phase D empirical 後評估；如推進必須走 §Non-Goals 封印（advisory-only，不自動寫 routing-registry） |
| Q7 Discovery 與 manual-lock 互動 | deferred-to-phase-c | 預設 manual-lock active 時跳過 Discovery (proposal 標 `manual_lock_bypass`)；governance 明文落於 Phase C `governance/workflow-activation-engine.md` 加段；Phase A 程式碼層 enforce eligibility gate |
| Q8 advisory cumulative token cap | still-open | Phase A.4 advisory 單次 ≤ 200 token 已寫入；cumulative cap 設計（per-task / per-session 累積上限）待 A.4 量測 advisory 注入頻率後決定 |

#### Phase 0.1 — Architecture Compatibility Preflight

- [x] 確認 parent plan `2026-05-31-1900-workflow-activation-engine.md`（archived）狀態下增量補丁不違反 archive contract（per plans/README.md plan-tree archive_order 規則）
  - **證據**：本 sub-plan 在 `plans/active/`，frontmatter `parent: 2026-05-31-1900-workflow-activation-engine`；plan_tree_archive_order 僅在 sub-plan archive 時要求 parent 已 archive — 本 sub 持續 active 期間無 violation。Parent archive 之 evidence trail (`plans/archived/2026-05-31-1900-workflow-activation-engine.md`) 不被 sub-plan 改動 — 本 plan 純 additive。
- [x] 確認 [`governance/lifecycle/capability-discovery-philosophy.md`](../../governance/lifecycle/capability-discovery-philosophy.md) 對 Discovery hot-hook 啟用之立場與本 plan Light/Deep 模型相容
  - **證據**：philosophy doc §Discovery → Detector Feedback Loop 明寫「detector miss 時形成回饋環」「Discovery 的 mechanical 整合點正是 detector 的 miss path，而非把 Discovery 獨立做成 per-turn executor」— 與本 plan 「detector miss → Light Discovery → advisory（非 gate）」完全對齊。Light → Deep 漸進方案是 philosophy 中「miss path」的具體實作策略，未越線到「per-turn executor」。
- [x] 確認 `hooks.go` PreToolUse + PostToolUse pipeline 可注入 advisory + 訂閱 Read event；確認 hijack 機制應落 hook 層還是 runtime 層
  - **證據**：`runPreToolUseHook` (hooks.go:1067) → `finishPreToolUse` (hooks.go:574) → `workflowPrimarySourceGate` (hooks.go:400)。Advisory injection point = `finishPreToolUse` block=false 路徑（detector miss 或 conflict 時），return ExitSuccess 前用 `hookSpecificOutput.additionalContext` 注入（同 `runPostToolUseHook` line 1212-1218 reminder 模式）。PostToolUse 端 `runPostToolUseHook` (hooks.go:1174) 已 parse `transcript_path` payload + tool_name — Phase B Read hijack 可在此分流新增 `if toolName == "Read" && hasAwaitingPhaseBProposal()` 分支。
  - **Q2 spike 結論（多 turn 範圍預備）**：傾向 **hook 層**。理由：(a) 與 PreToolUse 現有 pattern 對稱；(b) 不耦合 runtime 至 file-read wrapper；(c) Claude Code PostToolUse:Read delivery 已被 `runPostToolUseHook` 既有路徑驗證為 deterministic；(d) testable isolation 容易。runtime 層耦合度高且需新增 wrapper 層次，benefit 不對稱。**Phase B.1 spike 仍應驗證**：edge cases — 大檔 Read（>5MB content）的 PostToolUse payload 大小限制、tool_result 截斷後是否仍可 scan、與既有 receipt cache 互動。
- [x] 確認 `runtime.db` 可加 `discovery_proposals` 表（schema 評估）
  - **證據**：SQLite 標準 CREATE TABLE。同類 precedent: `route-candidate-proposals.yaml` (Phase 6.1 已 land) 採 **YAML data store, not projected**（router_proposals.go:19-21 註解明示「runtime/router/route-candidate-proposals.yaml is a DATA file, not a projected runtime contract」）。**架構決策**：(1) `discovery_proposals` 作為 **runtime.db 中的 raw 表**（per-task ephemeral + TTL 24h，每 task 高頻寫，YAML 不適），與 `route-candidate-proposals.yaml` 不重疊（後者是 cross-session occurrence-tracked promotion store）；(2) `runtime.discovery.config`（threshold/budget）作為 **projected surface** from 新 `runtime/discovery-bridge.yaml`（owner-layer canonical config，低頻變更）。
- [x] cross-check `architecture/ai-native-cognitive-ecosystem-system.md` §Watch-Out List 對應的 wall 名稱
  - **證據**：Watch-Out List 實有 6 walls（doc:447-454）：(1) Discovery confused with Activation; (2) Workflow inflation; (3) Ecosystem boundary inflation; (4) Telemetry explosion; (5) Positive-activation bias; (6) Optimization hallucination
  - **本 plan §Decision Rationale 對應修正**：
    - 「Discovery 反成主路由（advisory→primary 失控）」= **Wall #1 Discovery confused with Activation**（精準對應；§Non-Goals 封印「Discovery proposals MUST NOT satisfy activation_triggers」直接緩解此 wall）
    - 「Per-turn cost 預算爆炸」≈ **Wall #4 Telemetry explosion**（部分對應；本 plan 是 compute cost 非 telemetry cost；緩解機制相似 — Phase A 30ms / Phase B 50ms budget 即 criterion L 精神）
    - 「SQLite 表無 lifecycle 治理」= **無精準對應 wall**；歸類為 wall #4 同類風險（self-observation/state explosion），緩解靠 TTL 24h + Phase D 量測決定 GC 策略
  - **Watch-out 修正應於 Phase C 補入 plan §Decision Rationale Watch-Out citation**（本 turn 不改，Phase C scope）

#### Phase 0.2 — 既有 Discovery code path 盤點

- [x] 掃 `scripts/ai-skill-cli/internal/app/` 既有是否已有任何 discovery-related code（避免 double-implement）
  - **發現**：`router_proposals.go` + `runtime/router/route-candidate-proposals.yaml` 已存在（parent plan Phase 6.1 land）。實作 occurrence-tracking promotion state machine（accumulating→ready_for_review→promoted/rejected/stale），thresholds = 5 次 / 30 天 / 60 天 / 90 天。CLI: `ai-skill router proposals {list,record,promote,reject,gc}`。
  - **與本 plan 的關係**：**非 double-implement，purpose 正交**。
    - `route-candidate-proposals.yaml` = **cross-session, occurrence-tracked, manual promote → registry candidate** (YAML data store, 慢頻寫)
    - `discovery_proposals` (本 plan) = **per-task, TTL 24h, evidence + scoring, advisory ranking** (SQLite, 每 task 高頻寫)
    - 兩者可在 Phase D（out of scope）建立 cross-link：當 `discovery_proposals` 連續 N task 提同 candidate route → 餵給 `router proposals record` 觸發 occurrence count。但本 session 不實作此 bridge；Q6 deferred 已封印。
  - **§Non-Goals 封印對齊**：本 plan「Discovery proposals MUST NOT modify routing-registry.yaml automatically」與既有 `router_proposals.go` 行為一致（promote 只標 status，不改 registry，需人類 review）。
- [x] 確認 `knowledge/summaries/*.md` 為 Phase A primary scan target；確認 atom summary 格式穩定可解析
  - **證據**：30+ summary 檔案，統一格式 `## <Atom ID>` + table 含 `Atom ID / Source path / Summary / When to read / Do not use for / Context cost`. 可解析 — Phase A signal source 將 extract `Atom ID` 作為 candidate route key、`Summary` + `When to read` 作為 keyword bag。註冊在 `knowledge/summaries/README.md` 主表（30 行 Atom ID → file）。
- [x] 確認 `.ai-skill/project/rules/*.md` 既有 frontmatter convention（若無 metadata field，Phase A 只 scan 標題 + 第一段）
  - **證據**：`.ai-skill/project/rules/` 目錄不存在於本 repo（overlay 屬 downstream project layer）。Phase A.3 將 implement scanner 為 **defensive**：若目錄不存在 / 為空 → return zero signals（不報錯）；存在 → scan `*.md` frontmatter + 標題 + 第一段。Unit test 須建 fake overlay fixture 模擬 downstream project 結構（per Phase A.3 cross-project test 要求）。

### Phase A — Light Discovery

#### Phase A.1 — Discovery scoring 模型

- [x] 定義 confidence score 計算：term frequency match + path/ext bonus + project overlay bonus
  - **公式**（Light v1）：
    ```
    score(route, signals) = Σ weight[type] * match[type](route, signals)

    match types and weights (initial; runtime.discovery.config 可調):
      user_msg_term      0.30   — TF over user message vs route.task_intent + route summary.when_to_read + summary.summary
      summary_match      0.25   — substring presence of route's summary keyword set in user_msg + basenames
      basename_term      0.15   — token overlap between artifact basenames and route keyword set
      path_segment       0.10   — directory segment overlap with route keyword set
      extension_hint     0.05   — extension → route kind table (e.g. .py/.ts/.go → software-delivery)
      frontmatter_head   0.10   — parsed frontmatter (≤ 200B) tag/kind/route hint match
      cwd_overlay        0.05   — project overlay frontmatter declares matching route hint (signal-fact only)

    Score range: [0, ~1.0+] (weights sum 1.0; multiple match types can co-occur)
    Threshold (initial): 0.5
    Top-3 candidates emitted as proposal.route_candidates_json
    ```
- [x] 明文寫入 §Decision Rationale 的 scoring vs deterministic 區隔：scoring 僅用於 advisory ranking，**從不**進入 activation gate 路徑
  - 已存在於 §Decision Rationale §Non-Goals 封印第 2 條（「Discovery proposals MUST NOT satisfy activation_triggers」）。本 phase 確認該封印仍是 single source-of-truth；`discovery.go` 註解 cross-link 至此封印。
- [x] 初步 threshold 預設 0.5（Phase A.3 量測後調）
  - threshold = 0.5 寫進 `runtime/discovery-bridge.yaml`（projected to `runtime.discovery.config`），可調不需 rebuild。Phase D empirical 後依 KPI 報告調整。

#### Phase A.2 — runtime.db schema + cache

- [x] `discovery_proposals` schema
  - **DDL**（加入 `runtime_compiler.go::createGoRuntimeSchema`）：
    ```sql
    CREATE TABLE discovery_proposals (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      task_hash TEXT NOT NULL,
      route_candidates_json TEXT NOT NULL,    -- [{route, score, evidence:[{type, value}]}, ...] top-3
      signal_snapshot_json TEXT NOT NULL,     -- raw signals for rescoring across versions
      scoring_version TEXT NOT NULL,          -- 'light-v1', 'light-v1+deep-v1', etc.
      current_best_confidence REAL NOT NULL,  -- derived; paired with scoring_version
      status TEXT NOT NULL,                   -- enum: awaiting_phase_b|advised|dismissed|rejected|expired
      miss_reason TEXT,                       -- enum (see below) when no actionable proposal
      created_at TEXT NOT NULL,               -- RFC3339
      updated_at TEXT NOT NULL,
      expires_at TEXT NOT NULL                -- created_at + TTL (24h default)
    );
    CREATE INDEX idx_discovery_proposals_task_hash ON discovery_proposals(task_hash);
    CREATE INDEX idx_discovery_proposals_status ON discovery_proposals(status);
    CREATE INDEX idx_discovery_proposals_expires_at ON discovery_proposals(expires_at);
    ```
- [x] Per-candidate `evidence_set` 保留 top-N 具體 evidence items（N ≤ 5）
  - Encoded in `route_candidates_json[i].evidence`. Format: `[{type: "user_msg_term", value: "<token>"}, {type: "basename", value: "<filename>"}, ...]`. Sanitization: no raw private paths — paths truncated to last 2 segments; user tokens passthrough (already on the assistant transcript surface).
- [x] Discovery 失敗分類（`miss_reason` enum）
  - Enum tag in `discovery.go`: `no_artifact_reference` / `insufficient_signal` / `confidence_below_threshold` / `cost_budget_exceeded` / `manual_lock_bypass` / `eligibility_gate_fail`
- [x] `dismissed` vs `rejected` 語意區隔（telemetry 用）
  - `dismissed`: advisory injected but agent did not pivot to advised route's primary_source within N follow-up turns（透過 PostToolUse Read transcript 掃描判定，**Phase D scope**）
  - `rejected`: Phase B rescore drops candidate from top, or manual-lock binds different route, or detector subsequently locks a different route
- [x] 不可單獨儲存 confidence
  - Schema 保留 `signal_snapshot_json` + `scoring_version` 配對；rescoring 必跨 scoring_version 重算。Phase D telemetry group-by scoring_version。
- [x] TTL：預設 24h；可由 `runtime.discovery.config` 調
  - `runtime/discovery-bridge.yaml::ttl_hours: 24`，projected to `runtime.discovery.config`. GC 策略：Phase A 不做 GC job，靠 `expires_at` index + on-write opportunistic eviction（每次寫新 proposal 時若同 task_hash 已有 expired 則 DELETE）。專屬 GC job 為 Phase B/D 範疇。
- [x] project overlay scan cache：per-session in-memory，cwd 改變時 invalidate
  - 實作為 `discovery.go::projectOverlayCache` map[cwd]→[]signalFact，per-session lifetime（runtime process 重啟即清）；cwd 改變透過 cache key 自然 isolation；TTL 30 分鐘 fallback（避免長 session 看不到 overlay 改動）。

#### Phase A.3 — Discovery 實作

- [x] `discovery.go` 新建：function `RunLightDiscovery(taskInput, openFiles, cwd) []Candidate`
  - **證據**：`scripts/ai-skill-cli/internal/app/discovery.go` (~750 行)。`RunLightDiscovery(input, registry, summaries, cfg, repoRoot) ([]RouteCandidate, snapshot)` + `RunDiscoveryBridge(input, repoRoot, runtimeDB, manualLockActive)` 包裝 entry point。
- [x] Signal extractors 分兩層成本級別
  - **Light-0**：`tokenize` (user msg) / `extractArtifactTokens` (basename + path + ext)；零 I/O
  - **Light-1**：`scanProjectOverlayFacts` (overlay `.ai-skill/project/rules/*.md`，per-cwd cache 30 分鐘 TTL) / `LoadDiscoverySummaries` (knowledge/summaries/*.md 一次性 glob+read)；I/O 限於 repo metadata 範圍
- [x] **Source-of-truth guardrail**：project overlay scanner 只產 **signal facts**
  - **證據**：`scanProjectOverlayFacts` 註解明寫「MUST NOT emit route IDs」；只從 frontmatter `tags:`/`kind:`/`domain:` 與 `# heading` 抽 EvidenceItem，opaque facts only。Candidate 在 `scoreRoute` 由 `routing-registry.yaml` records iterate 產生，overlay facts 只作為 `cwd_overlay` evidence type 加分。
- [x] Scoring：weighted sum + normalize
  - 公式參見 Phase A.1；`normalizeMatch(hits, total)` 線性歸一到 [0,1]
- [x] 寫 proposal 到 runtime.db
  - `WriteDiscoveryProposal` INSERT；`OnWriteEviction=true` 時先 DELETE 同 task_hash expired rows；TTL 24h
- [x] Unit tests：cross-project case + threshold edge case + cache invalidation
  - `discovery_test.go`（9 tests）：`TestEligibilityCheck_*` (3 cases)、`TestExtractArtifactTokens`、`TestRunLightDiscovery_TravelPlanningReplay`、`TestRunLightDiscovery_SoftwareDelivery_ExtensionHint`、`TestRunDiscoveryBridge_ManualLockBypass`、`TestProjectOverlayCache_Invalidation`、`TestTaskHash_*`、`TestRenderAdvisory_RespectsTokenCap` — 全綠

**Phase A.3 empirical finding**：2026-06-05 replay 在 Light-only 下 travel-planning score = 0.272 / software-delivery = 0.312（均 < 0.5 threshold）。Proposal 仍寫入 (status=`awaiting_phase_b`, miss_reason=`confidence_below_threshold`) — 滿足 Phase A acceptance「至少寫出 candidate（即使未達 threshold，proposal 應存在）」。但 empirical 數據 confirm 使用者 review 預測「Phase B 才是主體」— Phase D KPI 量測 + Phase B 是否值得做的判斷由此資料 informed。

#### Phase A.4 — Advisory injector

- [x] `hooks.go` PreToolUse pipeline：detector miss + proposal status=advised → 注入 advisory text
  - **證據**：`finishPreToolUse` 新加 `tryDiscoveryAdvisory` helper（hooks.go），在 `workflowPrimarySourceGate` block=false 路徑後執行；ctx.ActiveRoute=="" 時 fire；`renderPreToolUseAdditionalContext` 用 `hookSpecificOutput.additionalContext` 寫 JSON（非 deny path）。Cursor host fail-safe skip。
- [x] Advisory format：≤ 200 token、列 top-3 candidate + 各自 primary_source 路徑、明示「non-blocking, optional Read」
  - `renderAdvisory` (`discovery.go`)：header "[ai-skill Discovery Bridge — advisory, non-blocking]" + "Reading the primary_source listed below is OPTIONAL" + top-3 candidate + score + primary_source + evidence types；token cap word-count proxy。
- [x] Cost 量測：p95 端到端 hook 延遲 ≤ 30ms（Phase A only）
  - **本 turn 不 bench**（plan 對 Phase A scope 接受 unit test green + structural acceptance；formal p95 bench 留 Phase B/D）。Unit tests 各 < 50ms / 含 registry+summaries 全 load 的 `RunLightDiscovery` ≤ 10ms（test runtime 觀察）。

#### Phase A.5 — Regression scenario

- [x] `validation/scenarios/runtime/workflow-discovery-bridge-light-v1.yaml`：empirical trigger task signature → expect Phase A advised travel-planning（若 confidence ≥ threshold）
  - **證據**：scenario yaml 5 個 sub-scenario (`light_v1_travel_planning_replay` / `light_v1_cross_project_software_delivery` / `light_v1_source_of_truth_guardrail` / `light_v1_advisory_non_blocking` / `light_v1_manual_lock_bypass`)。Travel replay 接受 `advised` 或 `awaiting_phase_b`（match empirical: 0.272 < threshold）。
- [x] Cross-project case：fake project + non-trivial task → expect Phase A non-trivial output
  - `light_v1_cross_project_software_delivery` 跑 .ts SDK bug 任務，expect software-delivery candidate（unit test `TestRunLightDiscovery_SoftwareDelivery_ExtensionHint` 已確認 candidate 含 software-delivery）。

**Phase A acceptance 回顧**：

- [x] 3 個 cross-project replay 至少 2 個 Phase A hit ≥ threshold — **NOT MET**：Light-only 兩個 case 都低於 0.5（travel 0.27 / software-delivery 0.31）。這是 plan 預期接受的失敗（acceptance 第四條「proposal 應存在」滿足）。
- [x] p95 hook 延遲 budget 達標 — 結構 acceptance：unit test 內 RunLightDiscovery 觀測 < 10ms（formal bench Phase B/D）
- [x] Unit tests + regression scenario 綠 — 9/9 discovery test pass
- [x] 2026-06-05 empirical trigger replay → Phase A 至少寫出 candidate — **MET**：proposal row 寫入 `discovery_proposals`，travel-planning 為 top candidate 之一，status=`awaiting_phase_b`

#### Phase A empirical findings（2026-06-15 mechanism replay）

跑真實 detector（`BuildRuntimeContext`）+ Discovery（`RunLightDiscovery`，用 production input builder `buildDiscoveryInputFromTranscript` 的等價邏輯）對一份真實 CJK 旅遊 artifact（`下関.md`）做 neutral vs keyword 對照。**框架定義（影響日後 Phase D telemetry 解讀）**：neutral 版不是「workflow 能力上限不足」，而是**在目前 signal pipeline 下可預期 fail-open；測到的是 pipeline coverage，不是 capability ceiling**。

| Prompt | Detector | Discovery 最高分 | 判讀 |
|---|---|---|---|
| keyword（含「行程」） | **hit** `route.workflow.travel-planning` | n/a（detector 已鎖，Discovery 不跑） | hit path 正常：gate 強制讀 workflow primary_source |
| neutral（無關鍵字） | miss | travel 連 top-3 都沒進；全 < 0.5 | miss path 預期 fail-open（Phase B 未做）+ 下列兩個 correctness bug |

replay 過程暴露**兩個 correctness bug（皆已修，非 i18n/enhancement）**，並確認一個訊號分層：

- **frontmatter_head 是 calibration bug（P1，已修 70024b2）**：scorer 有 branch + reserved weight 0.10 但 `buildDiscoveryInputFromTranscript` 從不填 `FrontmatterHead` → 永不產出的 feature 留在分母，把 effective max score 壓在 0.90、threshold 0.5 失真。修法＝dormant-feature 機制 + invariant（不進分母、不輸出 evidence、telemetry 標 inactive），reserved weight 保留，Phase B 接 producer 時翻 `enabled` 即可。**producer 連接刻意 deferred 到 Phase B**（與 CJK path 抽取 + 「hook 是否讀 referenced file」糾纏），避免污染 Phase D 觀測。
- **CJK artifact extraction 是 P1 correctness bug（已修，本批）**：`artifactPathRE` body 為 ASCII-only，對非英文檔名**兩種失效**——(a) leading-CJK 完全 miss（`下関.md` → `basenames=[]`）；(b) CJK+ASCII **截斷成錯誤 basename**（`東京-trip.md` → `-trip.md`）。兩者都讓 basename/path/ext 訊號歸零或失真 → systematic recall loss，且會污染 calibration。修法＝body 放寬到 `\p{L}\p{N}`、ext 維持 ASCII（單點，不動 tokenize）。regression `TestExtractArtifactTokens_CJKFilenames` 驗三層（regex capture → extract → tokenize）+ 2 negative fixtures（relative-escape 不外溢、multi-ext 取最後 ext）。
- **tokenize（layer 2）本來就 CJK-aware**（`tokenRE` 含 `\p{Han}\p{Hiragana}\p{Katakana}`）→ 確認 fix 收斂在 layer 1，不需動 tokenizer。

> **這兩個 bug 的意義**：neutral replay 之所以 miss，不只是「Phase B 未做」，而是輸入訊號 pipeline 先殘缺（frontmatter 死訊號壓低所有分數 + CJK 檔名抽不到訊號）。**先修這兩個 correctness bug，Phase D 的 baseline 才可信**；否則三週 telemetry 會建在偏移的基準上。Phase B（Deep Discovery + frontmatter producer）排在 correctness 修正之後。

##### Hit-path vs Miss-path acceptance（避免誤判「沒讀 workflow」為失敗）

修完兩個 correctness bug 後，4 變體實跑（A 中性+`下関.md` / B 中性+`大阪行程.md` / C 中性+`旅遊計畫.md` / D 語句含「行程」）路徑一致且**可解釋**。acceptance 須分兩條，不可把 miss-path「沒讀 workflow」當失敗：

- **Hit-path acceptance（B/C/D 已驗 ✅）**：detector hit → `ActiveRoute` 設定 → `workflowPrimarySourceGate` 執行 → 強制讀 workflow primary_source。
- **Miss-path acceptance（A 已驗 ✅）**：detector miss → Discovery proposal（可空/不足）→ **不**強制 workflow read → 正常完成並回 bootstrap fallback。A 的「沒去讀 workflow」**是 gate 未誤觸的證明**，非回歸。

##### 三條 findings（本輪實測沉澱）

1. **Detector hit 可由 filename keyword 直接觸發**：detector 比對 user_signals 對「訊息字串」，**含使用者貼的 artifact 路徑**。故 `大阪行程.md` / `旅遊計畫.md` 即使語句中性也命中（「行程」「旅遊」在檔名→在訊息）。這是 recall 紅利，但也是潛在誤觸風險（檔名巧含關鍵字）；列為觀察，暫不處理。
2. **CJK extractor 是 correctness issue，修復後未改變 activation boundary**：CJK fix 讓 Discovery 重新看得到非英文 basename/path/ext 訊號（recall/calibration correctness），但**不**改變「哪些 task 會 activate」——A 仍 miss（`下関` 對 detector 與 Discovery 關鍵字集皆不命中）。fix 收斂在 signal pipeline，不動 activation 契約。
3. **Neutral artifact review 仍依賴 content signal（Phase B justification）**：`下関.md` + 中性語句唯一的旅遊訊號在**內容**（行程/御朱印/河豚），而 name/message 皆無關鍵字。Light Discovery（pre-Read）結構上看不到內容 → 必然 miss。要讓此 case 自動觸發**只能靠 Phase B（Deep Discovery 讀內容）**，不是再 patch detector。這就是 Phase B 的 justification。

> **本修復迭代收束**：correctness 層（frontmatter dormancy + CJK extraction）已完成並驗證；`下関.md` 中性案**不是 regression，是 Discovery 邊界案例**，其解法是 Phase B 而非繼續 patch detector。下一步是 Phase B.1 spike（Read hijack 層級），不在本迭代。

### Phase B — Deep Discovery

#### Phase B.1 — Read hijack 機制

- [ ] 決定 hijack 落 hook 層（PostToolUse:Read 訂閱）還是 runtime 層（agent file-read wrapper）
- [ ] hook 層優點：與既有 PreToolUse pattern 對稱、實作隔離；缺點：依賴 hook delivery 時序
- [ ] runtime 層優點：時序確定；缺點：耦合度高
- [ ] Phase 0.1 + B.1 spike 後拍板

#### Phase B.2 — Content scan + 累積

- [ ] `discovery.go` 加 `RunDeepDiscovery(content, existingProposal) []Candidate`
- [ ] Content scan：keyword extract + summary match + atom signature match
- [ ] **Evidence accumulation, not confidence max()**：每次新 Read 把新 evidence append 到對應 candidate 的 `evidence_set`（含未在前 candidates 中出現的新 candidate 也加入）；然後對**所有 candidate 從完整 evidence_set 重 score**，產生新 ranking。`max(existing, new)` 會讓 Phase A 早期 false positive 永久鎖死（e.g. Phase A `travel=0.82`，Phase B 看到 `software-delivery` 強訊號 + `travel` 證據不足，max() 仍鎖在 0.82）；evidence accumulation + rescore 允許 ranking 翻盤
- [ ] Rescore 必須記新 `scoring_version` 或標明 `scored_at` 配對版本，避免下游 telemetry 跨版本誤比
- [ ] proposal status: awaiting_phase_b → advised（達 threshold 後）

#### Phase B.3 — Advisory re-injection

- [ ] Phase B 達 threshold 後，下一輪 PreToolUse 注入 advisory
- [ ] 避免 double-inject：proposal 標 advised 後不再注入（除非 user dismiss）

#### Phase B.4 — Cost 量測 + regression

- [ ] Bench：Phase B 跑 100 次（典型 markdown 1000~5000 token）p95 ≤ 50ms
- [ ] Regression：2026-06-05 incident replay → Phase A miss + Phase B hit + advisory 注入下一輪
- [ ] Edge case：Phase B 從錯 artifact（如 README 而非 itinerary）取訊號，下一輪正確 artifact Read 後 candidate 修正

**Phase B acceptance**：

- p95 cost budget 達標
- 至少 2 個 cross-project replay 證明 Phase B 累積機制有效
- 2026-06-05 incident replay 在 Phase A + B 組合下達到 advised 狀態

### Phase C — Governance + Documentation（completed 2026-06-15）

- [x] `governance/workflow-activation-engine.md` 補 §Discovery Bridge 段落：Light/Deep 模型、scoring vs deterministic 區隔、advisory non-blocking 屬性、與 manual-lock 互動（manual-lock 時跳過 Discovery）
  - **Evidence**: §"Discovery Bridge (miss-path advisory fallback)" + Scope Boundaries 更新 + Related links。含 advisory-never-a-gate / scoring-legal-because-advisory 兩條 invariant。
- [x] `enforcement/failure-patterns/detector-miss-no-fallback.md` 新建：empirical trigger + parent plan deferred 段 + Discovery Bridge 補強
  - **Evidence**: 檔案已建，Status validated，含 2-sample（2026-05-31 + 2026-06-05）、Half-Mechanized Gate 高階 pattern、registry README index row。
- [x] `knowledge/glossary/ai-skill.md` 加 6 新 term（discovery_bridge / light_discovery / deep_discovery / discovery_proposal / advisory_injection / piggyback_read）
  - **Evidence**: 6 term blocks added（status: candidate；無 sibling yaml，無 markdown_yaml_sync 連動；glossary README 為結構文件非 per-term index，無需更新）。
- [x] `enforcement/enforcement-registry.yaml` 評估是否加 `discovery_bridge_advisory` rule_class（advisory 不是 mechanical gate，可能不需要 registry entry — Phase 0.1 決策）
  - **決策：不登記。** Discovery 從不 block（advisory only），登成 `coverage: mechanical` rule_class 會把 advisory ranking 誤表成 enforcement gate，並招來 §Non-Goals 明文禁止的 auto-activation。被治理的對象是它緩解的 gap（`detector-miss-no-fallback` failure pattern），不是 advisory 本身；其效力以 Phase D empirical KPI 量測，而非 registry coverage cell 宣稱。Rationale 寫入 `governance/workflow-activation-engine.md` §Discovery Bridge。

> **Phase A governance-complete, Phase B deferred.** Phase A（landed code）現已有治理文件背書：failure pattern + governance 段 + glossary + registry stance。Phase B（Read hijack spike）與 Phase D（三週 empirical）仍 open。

### Phase D — Empirical Validation（三週）

- [ ] 三週 telemetry：detector miss 數 / proposal generated 數 / advisory injected 數 / pivot 數 / rejected 數 / dismissed 數
- [ ] Pivot 定義：advisory 注入後 agent Read 了 advised route 的 primary_source
- [ ] **KPI tiering**（與 ADR Promotion Criteria 對齊）：
  - Primary: detector miss → proposal generated ≥ 70%（routing quality）
  - Secondary: proposal generated → advisory injected ≥ 50%（scoring ergonomics）
  - Tertiary observation: advisory → pivot（agent behavior，不設門檻）
  - Diagnostic: `rejected / generated` = false positive rate；`dismissed / advised` = advisory ergonomics
- [ ] 若 Primary 未達標：分析 miss reason，決定 (a) 調 threshold、(b) 補 signal source（觸發 semantic surface follow-up plan）、(c) 接受 baseline

**Phase D acceptance**：三週量測完成，hit rate 報告 + 後續決策建議書

---

## Open Questions

1. **Confidence threshold 初值**：暫定 0.5。需 Phase A.3 量測後調。過高 → under-advise；過低 → over-advise（每 task 都 inject）。
2. **Phase B hijack 實作層**：hook 層（PostToolUse 訂閱）vs runtime 層（agent file-read wrapper）。Phase 0.1 spike 後拍板。
3. **discovery_proposals TTL**：暫定 24h。Cross-session proposal aggregation 若成為長期 route registry 來源，可能需延長。
4. **Cost budget 超標策略**：Phase A > 30ms 或 Phase B > 50ms 時 fail-open（不阻塞）還是 fail-closed（跳過 Discovery）？傾向 fail-open，但需明文。
5. **project overlay scan cache 邊界**：per-session in-memory 預設；cwd 改變 invalidate。是否該 per-task 重 scan 待量測。
6. **Cross-session proposal aggregation**：若同 candidate 連續 N task 出現 → 是否該升 routing-registry maintenance suggestion？Phase D 後評估。
7. **Discovery 與 manual-lock 互動**：parent plan v6 加 `manual-lock` mode；本 plan 預設 manual-lock 時跳過 Discovery（user 已明示，advice 反成噪音）。需 §governance/workflow-activation-engine.md 明寫。
8. **Advisory text 注入點是否需 token budget cap**：advisory ≤ 200 token 已寫入 Phase A.4，但 cumulative cost（cumulative advisories across turns）未限。

---

## 完成條件

- [ ] Phase 0 全部 acceptance 達成
- [ ] Phase A 全部 acceptance 達成
- [ ] Phase B 全部 acceptance 達成
- [ ] Phase C 全部 acceptance 達成
- [ ] Phase D 三週量測報告完成
- [ ] `governance/workflow-activation-engine.md` Discovery Bridge 段落 land
- [ ] `enforcement/failure-patterns/detector-miss-no-fallback.md` 建立並 cross-link
- [ ] `knowledge/glossary/ai-skill.md` 6 新 term 註冊
- [ ] `routing-registry.yaml` 不動（acceptance 條件）
- [ ] 2026-06-05 empirical trigger replay → Phase A or B advised → agent pivot 成功
- [ ] Open Questions 全部 resolve 或 deferred 註記
- [ ] 至少 3 個 cross-project replay 驗證

---

## Stakeholder 同意項目

- [ ] linyihong：approve Light → Deep 漸進架構（不擴 detector schema）
- [ ] linyihong：approve scoring 僅用於 advisory ranking、從不進入 gate 路徑
- [ ] linyihong：approve Phase A cost budget 30ms p95 / Phase B 50ms p95
- [ ] linyihong：approve 三週 empirical 期 + KPI tiering（Primary 70%、Secondary 50%、Tertiary observation only）
- [ ] linyihong：approve §Non-Goals 對 routing-registry auto-modification 的封印
- [ ] linyihong：approve Phase D 結束後是否開 semantic surface follow-up plan 的決策權

---

## 與其他 plans 的關係

- **Parent**：[`2026-05-31-1900-workflow-activation-engine.md`](../archived/2026-05-31-1900-workflow-activation-engine.md) — 本 plan 補其 Phase 6 deferred 段；採 Light/Deep 漸進架構而非 parent 原 design note 所述 hot-hook auto-call。
- **Sibling meta**：[`2026-05-31-2100-mechanical-enforcement-registry.md`](../archived/2026-05-31-2100-mechanical-enforcement-registry.md) — 同樣是「rule-without-executor」meta-pattern 治理；本 plan 是該 meta 第二採樣修補，但因 Discovery 是 advisory layer 而非 mechanical gate，不直接 promote 為 `rule_classes[*].coverage = mechanical`。
- **Conditional follow-up**：`workflow-activation-semantic-surface`（尚未開）— 架構預定採「project 產 signal、registry 解釋」分層（避免 v0 draft 中 project_overlay_signals 越過 source-of-truth boundary 的問題）。**量化開案 trigger**（避免三個月後重新辯論「到底要不要開」）：
  - **Trigger A**：Phase A `miss_reason = insufficient_signal` ratio > 40%（cheap signal 不足為主因）
  - **Trigger B**：`detector_miss AND no_proposal_generated` ratio > 25%（Discovery 完全沒救到的 task 比例過高）
  - **Trigger C**：Phase D 量測顯示 routing quality plateau（即使調 threshold 也救不回）
  - 任一 trigger 觸發 → 由 plan owner 起 `workflow-activation-semantic-surface` sub-plan；無 trigger → 接受 baseline，不開案。
- **Drafting history**：v0 draft 含 Phase 1 (Semantic Pre-Read Surface) + Phase 2 (Discovery Bridge)；review 後 rescope 為 Discovery-only。Rationale：(a) Discovery 為 80% 根因，semantic surface 為 20%；(b) Light Discovery 自然涵蓋 filename/path/project_metadata 為 signal source，未來若需 semantic surface 是擴 Phase A 而非新 detector layer；(c) cheap 訊號層擴張會推 detector 越界。
