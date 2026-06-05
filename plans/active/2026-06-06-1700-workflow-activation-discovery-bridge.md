---
id: 2026-06-06-1700-workflow-activation-discovery-bridge
plan_kind: sub
status: draft
owner: linyihong
created: 2026-06-06
parent: 2026-05-31-1900-workflow-activation-engine
required_for_completion: false
sub_plan_reason: >
  Workflow Activation Engine (parent) Phase 6 "Discovery → Detector feedback
  loop" 標 deferred 上線；同一失效模式（detector miss → 無 mechanical
  fallback → 靠 agent 自律 → 自律失敗）於 parent 完成隔日（2026-06-05）
  在 Travel project travel-planning 任務上原樣重演。本 sub-plan 補上 parent
  延後的 Discovery bridge，採 Light → Deep 兩階段漸進架構，避免 detector
  ontology 擴張並維持 parent §Design Principles 的 pre-Read 破環依賴原則。
  Independent sign-off：影響 per-turn cost 模型與 advisory 注入路徑，
  非單純 trigger 補洞。
---

# Workflow Activation: Discovery Bridge

**Status**: `draft`
Owner: framework maintainer (linyihong)
**世代**：Gen 3 runtime hardening — Workflow Activation Engine 第二採樣
**建立日期**：2026-06-06
**最後更新**：2026-06-06（v1 draft，rescope 自 v0 多 phase 草稿）
**Priority**：**P2**
**Parent plan**：[`2026-05-31-1900-workflow-activation-engine.md`](../archived/2026-05-31-1900-workflow-activation-engine.md)
**Empirical trigger**：2026-06-05 session — Travel project 內任務「幫我檢驗一下 `<dated-itinerary>.md`」。Detector miss（無 user keyword、無 path match），無 mechanical fallback，agent 直接憑常識做 review。文件命中 travel-planning artifact-gates 19 項中 7~10 項缺漏未被偵測。使用者三輪追問才暴露 gap，與 parent plan 2026-05-31 原 incident 為**同一結構性缺口的兩次採樣**。具體 project artifact 範例依 [`reusable-guidance-boundary.md`](../../enforcement/reusable-guidance-boundary.md) 留在原 project 文件。

> 本 plan **不擴 detector schema、不改 routing-registry**。範圍嚴格限定在 "detector miss → Discovery → advisory" 的 mechanical bridge。Filename / project metadata 等 semantic surface 議題 park 至條件式 follow-up plan，待本 plan 三週量測後依 miss rate 決定是否開。

---

## Decision Rationale

### Problem & Why Now

Parent plan v8（2026-06-04 完成）落地 detector + per-turn gate + manual-lock，**mechanical 層只完成一半**：detector hit 時機械擋住，detector miss 時 fail-open 並依賴 agent 自律 fallback 到 Discovery。Parent plan Phase 6 明寫 Discovery feedback loop「hot-hook auto-call 刻意延後」。延後代價在 parent 完成隔日就兌現（2026-06-05 Travel incident）。

失效路徑：

```
task input → detector(user_signals + context_signals(path)) miss
         → no mechanical fallback
         → agent self-route by intuition
         → workflow primary_source 未 Read
         → review 用常識做完，artifact-gates 19 項缺 7~10 項
```

**Why now**：parent plan archive 收尾 evidence 包含「detector miss 為設計接受的容忍範圍」假設。Travel incident 證明此容忍範圍對 cross-project + project-local-ontology 任務太寬。每加一個 project 就會踩同一個雷。

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

### Why Not an ADR Yet

- Phase A confidence threshold 未實測（暫定 0.5，依三週 empirical 量測再調）。
- Phase B piggyback 機制的 Read event 訂閱實作細節未確認（hook layer vs runtime layer）。
- discovery_proposals 表 lifecycle（TTL / promote / reject）governance 尚未定型。
- Cross-session proposal aggregation 是否該成為 long-term route registry 來源待 Phase 6 後決策。

### ADR Promotion Criteria（completed 時驗證）

- [ ] foundational + cross-session + cross-project + expensive-to-reverse + explains-why 全中
- [ ] Phase A + Phase B 跑過 ≥ 2 個 project 的 real task 驗證
- [ ] 三週 empirical：detector miss → Discovery hit → agent pivot 成功率 ≥ 60%
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

- [ ] 已讀本 plan §Open Questions 全部條目
- [ ] 對每條標記 `resolved` / `still-open` / `deferred`
- [ ] `resolved` 條目同步勾選於 §Open Questions
- [ ] 盤點新發現問題已加入 §Open Questions

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| Q1 confidence threshold 初值 | still-open | Phase A.3 量測後決定 |
| Q2 Phase B hijack 實作層 | still-open | Phase 0.1 + Phase B.1 設計後決定 |
| Q3 proposal TTL | still-open | Phase A.2 schema 設計時決定 |
| Q4 Cost budget 超標策略 | still-open | Phase A.4 量測後決定 |
| Q5 project overlay scan cache 邊界 | still-open | Phase A.2 設計決定 |
| Q6 cross-session proposal aggregation | deferred | 三週 empirical 後評估 |

#### Phase 0.1 — Architecture Compatibility Preflight

- [ ] 確認 parent plan `2026-05-31-1900-workflow-activation-engine.md`（archived）狀態下增量補丁不違反 archive contract（per plans/README.md plan-tree archive_order 規則）
- [ ] 確認 [`governance/lifecycle/capability-discovery-philosophy.md`](../../governance/lifecycle/capability-discovery-philosophy.md) 對 Discovery hot-hook 啟用之立場與本 plan Light/Deep 模型相容
- [ ] 確認 `hooks.go` PreToolUse + PostToolUse pipeline 可注入 advisory + 訂閱 Read event；確認 hijack 機制應落 hook 層還是 runtime 層
- [ ] 確認 `runtime.db` 可加 `discovery_proposals` 表（schema 評估）
- [ ] cross-check `architecture/ai-native-cognitive-ecosystem-system.md` §Watch-Out List 對應的 wall 名稱

#### Phase 0.2 — 既有 Discovery code path 盤點

- [ ] 掃 `scripts/ai-skill-cli/internal/app/` 既有是否已有任何 discovery-related code（避免 double-implement）
- [ ] 確認 `knowledge/summaries/*.md` 為 Phase A primary scan target；確認 atom summary 格式穩定可解析
- [ ] 確認 `.ai-skill/project/rules/*.md` 既有 frontmatter convention（若無 metadata field，Phase A 只 scan 標題 + 第一段）

### Phase A — Light Discovery

#### Phase A.1 — Discovery scoring 模型

- [ ] 定義 confidence score 計算：term frequency match + path/ext bonus + project overlay bonus
- [ ] 明文寫入 §Decision Rationale 的 scoring vs deterministic 區隔：scoring 僅用於 advisory ranking，**從不**進入 activation gate 路徑
- [ ] 初步 threshold 預設 0.5（Phase A.3 量測後調）

#### Phase A.2 — runtime.db schema + cache

- [ ] `discovery_proposals` schema：`id` / `task_hash` / `route_candidates_json` / `best_confidence` / `status` (`awaiting_phase_b` / `advised` / `dismissed` / `expired`) / `created_at` / `updated_at`
- [ ] TTL：預設 24h；可由 `runtime.discovery.config` 調
- [ ] project overlay scan cache：per-session in-memory，cwd 改變時 invalidate

#### Phase A.3 — Discovery 實作

- [ ] `discovery.go` 新建：function `RunLightDiscovery(taskInput, openFiles, cwd) []Candidate`
- [ ] Signal extractors：user_msg tokenizer、artifact basename parser、frontmatter head reader（≤ 200B）、project overlay scanner
- [ ] Scoring：weighted sum + normalize
- [ ] 寫 proposal 到 runtime.db
- [ ] Unit tests：cross-project（≥ 3 project type）case + threshold edge case + cache invalidation

#### Phase A.4 — Advisory injector

- [ ] `hooks.go` PreToolUse pipeline：detector miss + proposal status=advised → 注入 advisory text
- [ ] Advisory format：≤ 200 token、列 top-3 candidate + 各自 primary_source 路徑、明示「non-blocking, optional Read」
- [ ] Cost 量測：p95 端到端 hook 延遲 ≤ 30ms（Phase A only）

#### Phase A.5 — Regression scenario

- [ ] `validation/scenarios/runtime/workflow-discovery-bridge-light-v1.yaml` 加 case：empirical trigger 的 task signature → expect Phase A advised travel-planning（若 confidence ≥ threshold）
- [ ] Cross-project case：fake project + non-trivial task → expect Phase A non-trivial output

**Phase A acceptance**：

- 3 個 cross-project replay 至少 2 個 Phase A hit ≥ threshold
- p95 hook 延遲 budget 達標
- Unit tests + regression scenario 綠
- Travel project empirical trigger replay → Phase A 至少寫出 candidate（即使未達 threshold，proposal 應存在）

### Phase B — Deep Discovery

#### Phase B.1 — Read hijack 機制

- [ ] 決定 hijack 落 hook 層（PostToolUse:Read 訂閱）還是 runtime 層（agent file-read wrapper）
- [ ] hook 層優點：與既有 PreToolUse pattern 對稱、實作隔離；缺點：依賴 hook delivery 時序
- [ ] runtime 層優點：時序確定；缺點：耦合度高
- [ ] Phase 0.1 + B.1 spike 後拍板

#### Phase B.2 — Content scan + 累積

- [ ] `discovery.go` 加 `RunDeepDiscovery(content, existingProposal) []Candidate`
- [ ] Content scan：keyword extract + summary match + atom signature match
- [ ] Append-only update：每次新 Read 來合併 candidate，confidence 用 max(existing, new) 而非 overwrite
- [ ] proposal status: awaiting_phase_b → advised（達 threshold 後）

#### Phase B.3 — Advisory re-injection

- [ ] Phase B 達 threshold 後，下一輪 PreToolUse 注入 advisory
- [ ] 避免 double-inject：proposal 標 advised 後不再注入（除非 user dismiss）

#### Phase B.4 — Cost 量測 + regression

- [ ] Bench：Phase B 跑 100 次（典型 markdown 1000~5000 token）p95 ≤ 50ms
- [ ] Regression：Travel incident replay → Phase A miss + Phase B hit + advisory 注入下一輪
- [ ] Edge case：Phase B 從錯 artifact（如 README 而非 itinerary）取訊號，下一輪正確 artifact Read 後 candidate 修正

**Phase B acceptance**：

- p95 cost budget 達標
- 至少 2 個 cross-project replay 證明 Phase B 累積機制有效
- Travel incident replay 在 Phase A + B 組合下達到 advised 狀態

### Phase C — Governance + Documentation

- [ ] `governance/workflow-activation-engine.md` 補 §Discovery Bridge 段落：Light/Deep 模型、scoring vs deterministic 區隔、advisory non-blocking 屬性、與 manual-lock 互動（manual-lock 時跳過 Discovery）
- [ ] `enforcement/failure-patterns/detector-miss-no-fallback.md` 新建：empirical trigger + parent plan deferred 段 + Discovery Bridge 補強
- [ ] `knowledge/glossary/ai-skill.md` 加 6 新 term（discovery_bridge / light_discovery / deep_discovery / discovery_proposal / advisory_injection / piggyback_read）
- [ ] `enforcement/enforcement-registry.yaml` 評估是否加 `discovery_bridge_advisory` rule_class（advisory 不是 mechanical gate，可能不需要 registry entry — Phase 0.1 決策）

### Phase D — Empirical Validation（三週）

- [ ] 三週 telemetry：detector miss 數 / Discovery proposal 數 / advisory inject 數 / agent pivot 成功數
- [ ] Pivot 成功定義：advisory 注入後 agent Read 了 advised route 的 primary_source
- [ ] 目標 hit rate：detector miss → Discovery proposal ≥ 70%；proposal advised → agent pivot ≥ 60%
- [ ] 若未達標：分析 miss reason，決定 (a) 調 threshold、(b) 補 signal source（觸發 semantic surface follow-up plan）、(c) 接受 baseline

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
- [ ] Travel project empirical trigger replay → Phase A or B advised → agent pivot 成功
- [ ] Open Questions 全部 resolve 或 deferred 註記
- [ ] 至少 3 個 cross-project replay 驗證

---

## Stakeholder 同意項目

- [ ] linyihong：approve Light → Deep 漸進架構（不擴 detector schema）
- [ ] linyihong：approve scoring 僅用於 advisory ranking、從不進入 gate 路徑
- [ ] linyihong：approve Phase A cost budget 30ms p95 / Phase B 50ms p95
- [ ] linyihong：approve 三週 empirical 期 + hit rate 目標（detector miss → proposal ≥ 70%；proposal advised → pivot ≥ 60%）
- [ ] linyihong：approve Phase D 結束後是否開 semantic surface follow-up plan 的決策權

---

## 與其他 plans 的關係

- **Parent**：[`2026-05-31-1900-workflow-activation-engine.md`](../archived/2026-05-31-1900-workflow-activation-engine.md) — 本 plan 補其 Phase 6 deferred 段；採 Light/Deep 漸進架構而非 parent 原 design note 所述 hot-hook auto-call。
- **Sibling meta**：[`2026-05-31-2100-mechanical-enforcement-registry.md`](../archived/2026-05-31-2100-mechanical-enforcement-registry.md) — 同樣是「rule-without-executor」meta-pattern 治理；本 plan 是該 meta 第二採樣修補，但因 Discovery 是 advisory layer 而非 mechanical gate，不直接 promote 為 `rule_classes[*].coverage = mechanical`。
- **Conditional follow-up**：`workflow-activation-semantic-surface`（尚未開）— 若本 plan Phase D 量測顯示 Phase A miss rate 過高（即 cheap signal 不足），才開此 plan 處理 filename_signals / project_metadata_signals。架構預定採「project 產 signal、registry 解釋」分層（避免 v0 draft 中 project_overlay_signals 越過 source-of-truth boundary 的問題）。
- **Drafting history**：v0 draft 含 Phase 1 (Semantic Pre-Read Surface) + Phase 2 (Discovery Bridge)；review 後 rescope 為 Discovery-only。Rationale：(a) Discovery 為 80% 根因，semantic surface 為 20%；(b) Light Discovery 自然涵蓋 filename/path/project_metadata 為 signal source，未來若需 semantic surface 是擴 Phase A 而非新 detector layer；(c) cheap 訊號層擴張會推 detector 越界。
