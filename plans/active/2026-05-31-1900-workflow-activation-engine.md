# Workflow Activation Engine

**Status**: `draft`
**世代**：Gen 3 runtime hardening（systemic gap remediation）
**建立日期**：2026-05-31
**最後更新**：2026-05-31（initial draft）
**Empirical trigger**：2026-05-31 session — agent 對 `docs/20260531-下関.md` 跑 review，doc 標題含「行程」、內容含「Day 1 / 御朱印 / MapCode / 自駕」，命中 `route.workflow.travel-planning.activation_triggers.user_signals` 全部訊號，但 workflow 從未被啟動。使用者連續三次追問才暴露此 gap。

> 本 plan 不修 travel-planning 個案，而是補齊 **Workflow Activation Engine** ——目前 framework 第一次形成「Registry ✓ + Rules ✓ + Docs ✓ + **Activation Engine ✗**」閉環的缺角。

---

## Decision Rationale

### Empirical Evidence（Registry 體檢）

`knowledge/runtime/routing-registry.yaml` 現況：

| 指標 | 數量 |
|---|---|
| Total `route.*` records | 57 |
| 有 `activation_triggers` | **7**（apk-analysis、software-delivery、greenfield、travel-planning、documentation-ai-native、governance.system-upgrade、runtime.decision-recording） |
| 無 `activation_triggers` | **50**（全部 `route.analysis.*`、全部 `route.intelligence.*`、大部分 `route.governance.*` / `route.feedback.*` / `route.memory.*` / `route.constitution.*`） |
| 命中後會自動啟動的 detector | **0** |

兩層 gap：
- **L1 data gap**：87% route 連 trigger schema 都沒填
- **L2 executor gap**：即使有 trigger 也沒人跑（`hooks.go` grep travel = 0、route.workflow = 0）

### Failure Mode Classification

這次失效非 travel-planning 個案，是 **systemic detection gap**。同樣失效模式會在 `route.analysis.web`、`route.analysis.apk.workflows`、`route.intelligence.architectural-fit`、`route.intelligence.requirements-cognition` 等任務上重現 —— 任何「規則明明寫好但沒人去觸發」的 route。

### Decision

建立 **Workflow Activation Engine** 作為 Ai-skill 第四個 runtime 層：

```
Registry (routing-registry.yaml)
        ↓
Detector (NEW: deterministic rule match)
        ↓
Workflow Session (NEW: runtime.db.workflow_sessions table)
        ↓
Execution (existing: tool calls read CurrentWorkflow)
        ↓
Enforcement (existing: per_turn_obligations + commit validators)
        ↓
Feedback (existing: feedback/history/<domain>/)
```

### Design Principles（接受第三方架構評論）

| Decision | Rationale |
|---|---|
| **Deterministic rule match，不用 weighted scoring** | 規則問題不該變成分數問題。`Day 1` + `Day 2` + `行程` + `自駕` 是 deterministic signal，不需要 confidence threshold。0.62 vs 0.58 的調參地獄是 anti-pattern。 |
| **Two-stage：rule match → conflict resolution** | Stage 1 boolean `any_of` / `artifact_any`，命中為 TRUE。Stage 2 只在多 route 同時 TRUE 時進入 `workflow/workflow-routing.md` 既有歧義裁決。 |
| **Session state，不是 per-tool PreToolUse** | Detector 一個 task 跑一次，結果寫 `workflow_sessions` 表。後續 tool calls 讀 `CurrentWorkflow`，O(1) lookup。比每次 PreToolUse 重算便宜兩個數量級。 |
| **不全部補 50 條 triggers** | 50 條無 trigger 的 route 要先分類：always-on / triggered / reference-only。只給 triggered 類補 schema。 |
| **Discovery 保留原樣** | Discovery 不是這次 bug 根因，是另一個機制（找 unknown capability）。本 plan **不**升級成 per-turn obligation，避免成本爆炸。Discovery fallback 只在 detector 全 miss 時觸發。 |

### Why Not Quick Fix Travel-Planning

幫 travel-planning 寫 special-case validator 是症狀修補：
- 不解決 `route.analysis.*` 50 條同樣問題
- 沒有可重複使用的 detection runtime
- 違反「systemic gap 需 framework patch」原則

---

## Architecture Compatibility Preflight

依 [`plans/README.md`](../../plans/README.md#plan-執行前架構相容性檢查architecture-compatibility-preflight)：

| 欄位 | 內容 |
|---|---|
| Candidate files | `knowledge/runtime/routing-registry.yaml`（擴 schema）、`scripts/ai-skill-cli/internal/app/hooks.go`（加 detector validator）、`runtime/runtime.db`（新增 `workflow_sessions` 表）、`runtime/core-bootstrap.yaml`（per_turn_obligations 加 detector check）、新建 `governance/workflow-activation-engine.md`（philosophy）、新建 `enforcement/failure-patterns/workflow-detector-missing.md` |
| Source-of-truth | `routing-registry.yaml` 仍是 trigger 唯一來源。`runtime.db.workflow_sessions` 是 runtime state，不是 canonical。 |
| Compiler / generated surfaces | `runtime.db` 需重 compile；`ai-skill runtime compile + refresh` 流程不變 |
| Layer responsibility | Detector 屬 runtime layer（scripts/ai-skill-cli）；Schema 屬 knowledge/runtime layer；Philosophy 屬 governance；Failure pattern 屬 enforcement |
| 與現行架構衝突 | 無。本 plan 補的是 missing layer，不改既有 layer 職責 |
| `runtime.db` / generated surface 影響 | 新增表 + 新增 obligation；compile pipeline 需要新 projection rule |

---

## Phase Plan

### Phase 0 — Preflight

#### Phase 0.0 — Open Questions 核對

逐條核對本 plan §Open Questions，標記處置：

- [ ] 已讀本 plan §Open Questions 全部條目
- [ ] 對每條標記 `resolved` / `still-open` / `deferred`
- [ ] 新發現問題已加入 §Open Questions

#### Phase 0.1 — Architecture Compatibility Preflight

- [ ] 確認 `governance/lifecycle/capability-discovery-philosophy.md` 與本 plan 的 Discovery vs Detector 分工不衝突（companion 章節需註記 "Detector handles known routes, Discovery handles unknown capabilities"）
- [ ] 確認 `workflow/workflow-routing.md` 既有歧義裁決可作為 Stage 2 conflict resolver
- [ ] 確認 `runtime.db` schema 可加 `workflow_sessions` 表而不破壞既有 projection
- [ ] 確認 `hooks.go` PreToolUse pipeline 可注入新 validator（非阻塞性，僅 detector miss 時 reject）

#### Phase 0.2 — Route Classification（必跑，本 plan 後續所有 phase 的基礎）

逐條檢視 57 個 route，分類為：

| 類別 | 定義 | 行動 |
|---|---|---|
| `always-on` | session bootstrap、runtime core、phase machine 等永遠該載入的 | 不需 triggers，標記 `preload: true` |
| `triggered` | workflow.\*、analysis.\*、intelligence.\* 之領域型工作流 | **必須有 activation_triggers**，本 plan Phase 2 補齊 |
| `reference-only` | governance / constitution / architecture 之描述性文件 | 等 user 明確問起，標記 `on-demand: true`，不進 detector |

產出：`routing-registry.yaml` header 加 `route_classification` schema，每條 route 加 `class:` 欄位。

預估分類（待 Phase 0.2 確認）：always-on ~10、triggered ~25、reference-only ~22。

### Phase 1 — Detector Schema 定義

在 `routing-registry.yaml` 擴 `activation_triggers` schema（向後相容）：

```yaml
activation_triggers:
  any_of:                    # NEW: deterministic boolean
    user_signals: [行程, itinerary, 旅遊]      # 對話文字（既有，重新命名）
    artifact_signals:        # NEW: 已讀檔案內容 pattern
      - "Day [0-9]+"
      - "御朱印"
      - "MapCode"
    context_signals:         # NEW: 檔名 / 路徑 pattern
      - "docs/*行程*.md"
      - "docs/[0-9]{8}-*.md"
  task_intents: [travel-planning, itinerary]   # 既有，保留
```

規則：**任一 `*_signals` 內任一條 hit → TRUE**。不加權、不算分。

向後相容：舊格式（直接 `user_signals: [...]`）仍接受，視為 `any_of.user_signals`。

產出：
- [ ] `routing-registry.yaml` schema 更新 + 文件
- [ ] `governance/workflow-activation-engine.md` 新建（philosophy + schema spec）
- [ ] 7 個既有 `activation_triggers` 路由不動（schema 已相容）

### Phase 2 — 為 triggered 類 route 補 schema

依 Phase 0.2 分類結果，為 ~25 個 `triggered` route 補 `activation_triggers`：

優先順序：
1. `route.analysis.web`、`route.analysis.apk.workflows`（最近活躍領域）
2. `route.intelligence.architectural-fit`、`route.intelligence.requirements-cognition`、`route.intelligence.engineering.agent-architecture`
3. 其餘 triggered route

每條 route 至少給 `user_signals` + `context_signals`。`artifact_signals` 可選（部分 route 沒有明顯 artifact pattern）。

產出：
- [ ] ~25 條 route 補 triggers
- [ ] 跑 `ai-skill runtime compile + refresh`
- [ ] validation：每條 route 至少 1 個 signal 來源

### Phase 3 — Detector 實作（Go）

在 `scripts/ai-skill-cli/internal/app/` 加 `detector.go`：

```go
// 簽名：
func DetectWorkflows(transcript []Message, openFiles []FileRef) []DetectedRoute

// 邏輯：
// 1. Concat transcript text (recent N user messages) + openFiles content
// 2. For each route where class == "triggered":
//      hit := any(user_signals) ∪ any(artifact_signals on content) ∪ any(context_signals on file_paths)
//      if hit { detected.append(route_id) }
// 3. Return detected (可能空、單一、多個)
```

整合點：
- **PreToolUse hook**：先查 `workflow_sessions` 表本 task 是否已 detect。已 detect → skip。未 detect → run detector，寫入表。
- **Conflict path**：detected.len > 1 → 注入 reminder 指向 `workflow/workflow-routing.md` Step 3 歧義裁決，讓 agent 自己選；不自動鎖定。
- **Miss path**：detected.len == 0 → 不阻擋，但記 `workflow_sessions.status = no-match`，未來分析這些 case 可能要加 triggers。

產出：
- [ ] `detector.go` + unit tests
- [ ] `hooks.go` 整合
- [ ] `runtime.db` 加 `workflow_sessions` 表（schema 見 Phase 4）

### Phase 4 — Workflow Session State

新增 `runtime.db.workflow_sessions` table：

```sql
CREATE TABLE workflow_sessions (
  id TEXT PRIMARY KEY,           -- uuid
  task_id TEXT NOT NULL,         -- 由 agent first substantive message hash 衍生
  session_id TEXT NOT NULL,      -- harness session id
  detected_routes TEXT,          -- JSON array of route ids
  active_route TEXT,             -- 單一鎖定的 route（conflict 解決後）
  detection_source TEXT,         -- which signals fired: user/artifact/context
  status TEXT NOT NULL,          -- detected / locked / no-match / invalidated
  activated_at TIMESTAMP,
  invalidated_at TIMESTAMP,
  invalidation_reason TEXT       -- topic-shift / explicit-pivot / keyword-drift
);
```

Lifecycle：
1. **Task start detection**：first user substantive message + first Read 後跑 detector
2. **Topic shift detection**：
   - 顯式：user message 含 `換任務 / 現在我要 / new task / switch to` 等 sentinel
   - 隱式：連續 5 turn 內 active_route 的 keywords 完全沒再出現 → invalidate，下次 user message 重跑 detector
3. **Manual override**：user 顯式說「跟我做 X」直接覆寫 active_route

產出：
- [ ] `runtime.db` migration
- [ ] `ai-skill runtime workflow-session` CLI subcommand（查當前 active）
- [ ] Session lifecycle 文件化

### Phase 5 — Obligation 整合

在 `runtime/core-bootstrap.yaml` 加 `per_turn_obligations`：

```yaml
- id: obligation.workflow.activation_evidence
  fires: first_substantive_response_after_detection
  action: |
    若 workflow_sessions.active_route != null，agent 必須在工具呼叫前
    Read 該 route 的 primary_source。validator 掃 transcript 確認。
  severity: high
  blocking_gate_id: gate.workflow.primary_source_read
```

`hooks.go` 新增 validator `validateWorkflowPrimarySourceRead`：類似 `bootstrap.receipt_present` 模式，掃 transcript 確認 Read 事件。

**這不是 Discovery、不是每 turn 跑**：只在 detector 已鎖定 active_route 後生效。沒鎖定 = 不阻擋。

### Phase 6 — Failure Pattern + Discovery 邊界澄清

- [ ] 新建 `enforcement/failure-patterns/workflow-detector-missing.md` —— 記錄這次失效（2026-05-31 session log）為 systemic gap，並把 Detector 設計指回本 plan
- [ ] 更新 `governance/lifecycle/capability-discovery-philosophy.md` —— 加章節「Discovery vs Detector 分工」，明確：
  - Detector 處理 known route 的 known trigger
  - Discovery 處理 unknown capability（detector miss 後 fallback）
  - 兩者不重疊、不取代

### Phase 7 — Validation Scenarios

新建 scenarios：
- `validation/scenarios/runtime/workflow-detector-deterministic-match-v1.yaml`
- `validation/scenarios/runtime/workflow-detector-conflict-resolution-v1.yaml`
- `validation/scenarios/runtime/workflow-session-topic-shift-v1.yaml`
- `validation/scenarios/runtime/workflow-detector-travel-planning-regression-v1.yaml`（這次 bug 的 regression test）

Acceptance：四個 scenario 全 PASS，且回放 2026-05-31 session 時 travel-planning detector 必須觸發。

### Phase 8 — Close-out

- [ ] 全部 phase done
- [ ] `git status` clean
- [ ] `git push` 完成、`git log origin/main..HEAD` empty
- [ ] 讀回更新後的 `routing-registry.yaml` / `core-bootstrap.yaml` / failure pattern
- [ ] 在本 plan 加 Phase 8 完成記錄 + archive 到 `plans/archived/`

---

## Open Questions

| # | Question | 處置 |
|---|---|---|
| Q1 | Route classification 是否需要使用者 review 才定案？50 條人工分類有主觀成分 | still-open — 建議 Phase 0.2 產出 draft 後等 user confirm |
| Q2 | Detector 的「first substantive message」定義 —— 純打招呼算嗎？ | still-open — 建議以 ≥ 20 chars + 含動詞 / 名詞為門檻 |
| Q3 | Conflict resolution 多 route 命中時自動選還是 prompt user？ | resolved → 不自動選，注入 reminder 讓 agent 走 `workflow-routing.md` |
| Q4 | `workflow_sessions` TTL？跨 session 是否保留？ | still-open — 初版建議 session-scoped，task_id 由 session_id + first message hash 構成 |
| Q5 | Detector miss（detected.len == 0）是否該自動 fallback 到 Capability Discovery？ | still-open — 建議**不**自動 fallback，只記錄 `status=no-match`，避免 silent expensive discovery 跑滿 turn |
| Q6 | 舊格式（直接 `user_signals` 不在 `any_of` 下）的 deprecation timeline？ | still-open — 建議無限期相容，Phase 2 補新 route 用新格式即可 |
| Q7 | `artifact_signals` 在 Read tool 觸發時才掃，還是每 user message 掃？ | still-open — 建議「最近 Read 的 N 個檔案 + 累積 user messages」一起掃 |

---

## Validation Plan

- [ ] Phase 0.2 route classification 經 user review
- [ ] Phase 1 schema 變更 backward compat（既有 7 條 route 不需改即可運作）
- [ ] Phase 2 新增 triggers 經抽樣 review（≥ 5 條）
- [ ] Phase 3 detector unit tests 涵蓋：single hit、multi hit、no match、舊格式相容
- [ ] Phase 4 `workflow_sessions` lifecycle 經 integration test 驗證
- [ ] Phase 5 obligation 不誤殺：當 detector miss 時 tool calls 不被擋
- [ ] Phase 7 regression scenario：2026-05-31 session 場景 replay 必須觸發 travel-planning
- [ ] Phase 8 close-loop：所有變更 commit / push / readback

---

## Dependency Read Ledger（plan drafting 階段）

| 欄位 | 內容 |
|---|---|
| Trigger | User 明確授權「可以寫入計畫」+ 要求 analysis 層體檢 |
| Required set | `CORE_BOOTSTRAP.md`、`runtime/core-bootstrap.yaml`、`enforcement/{rule-weight, dependency-reading, conversation-goal-ledger}.md`、`knowledge/runtime/routing-registry.yaml`、`workflow/travel-planning/{README, execution-flow, artifact-gates}.md`、`workflow/workflow-routing.md`、`governance/lifecycle/capability-discovery-philosophy.md`、`plans/active/*.md`（template reference） |
| Read | 全部 above |
| Not applicable | `workflow/greenfield/templates/plan-template.md` 未讀（plan 結構參考既有 plan，未直接 derive template） |
| Deferred | Implementation phase 才需要的 source（`hooks.go` 細節、`runtime.db` schema 細節）—— Phase 0 Preflight 開始時補讀 |
| Validation | 本 plan 之 Architecture Compatibility Preflight 章節已列各 candidate file；Phase 0.1 進入 implementation 前再驗證 |

---

## Source

2026-05-31 session：使用者連續三次追問才暴露 `route.workflow.travel-planning` activation gap。第三方對話建議拆解 Discovery vs Detector、放棄 weighted scoring、改用 deterministic + workflow_sessions 設計。本 plan 接受全部建議並加入 route classification + 50 條 trigger gap 的 systemic 修補。

## Companion References

- `governance/lifecycle/capability-discovery-philosophy.md` —— Discovery 機制（與本 plan 互補）
- `workflow/workflow-routing.md` —— 多 route 命中時的歧義裁決（Stage 2 conflict resolver）
- `enforcement/dependency-reading.md` §Workflow 編排 —— blocking activation 行為強制（本 plan 升級為機械強制）
- `enforcement/failure-patterns/bootstrap-bypass-on-resume.md` —— PreToolUse + transcript scan 模式範例（detector 採同模式）
