# Workflow Activation Engine

**Status**: `draft-v4`
**世代**：Gen 3 runtime hardening（systemic gap remediation）
**建立日期**：2026-05-31
**最後更新**：2026-05-31（v4 — 整合第四輪評審：intelligence 不適用單一預設，改 must-declare + Phase 0.2 加 intelligence 全盤審計 gate）
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
| **Session state in-memory first，SQLite 延後** | Detector 一個 task 跑一次，結果存 Go `RuntimeContext` struct。**不**直接落 `runtime.db`。SQLite 落地等到出現實際需求（跨 session replay / 統計 / 分析）才做，避免「還沒驗證需求先固化 schema」。詳見 Phase 4。 |
| **不全部補 50 條 triggers** | 50 條無 trigger 的 route 要先分類為 `activation_mode`：`always-on` / `auto-detect` / `on-demand` / `advisory` —— 不是二元「triggered / reference-only」，因為部分 route（如 `route.intelligence.architectural-fit`）可能 multi-mode。 |
| **Discovery 與 Detector 互補，不是切開** | Discovery 不升級成 per-turn obligation（避免成本爆炸），但 detector miss 時應 fallback 到 Discovery，由 Discovery 產出 `new route candidate` 反饋給 Registry。建立 **Discovery → Detector feedback loop**，讓 Registry 隨使用自然成長。詳見 Phase 6。 |
| **Activation 不能依賴 Read File**（破環依賴） | `artifact_signals` 需要先 Read 才有內容，但 workflow activation 理論上應在 Read 之前。所以 activation 分兩階段：Phase 1 用 `user_signals` + `context_signals`（pre-Read 即可取得）→ activate workflow。Phase 2 用 `artifact_signals` 做**強化驗證**，不是 activation gate。詳見 Phase 1 schema。 |

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

#### Phase 0.2 — Route Type + Activation Mode 分類（**v3 改：self-declared，不再靠人工 classification 表**）

##### Phase 0.2a — 引入 `route_type`（self-declared，scale-friendly）

第三輪評審指出：57 routes 今天人工分類還可行，但 120 / 300 routes 後維護地獄。解法 —— 每個 route **自己宣告 `route_type`**，由 type 推導預設 `activation_mode`。新增 route 時自然就帶分類，不需中央表。

```yaml
- id: route.workflow.travel-planning
  route_type: workflow                # NEW: self-declared
  activation_mode: auto-detect        # 可省略 → 由 route_type 推導
  # 若需 override（罕見），直接寫值即可
```

**`route_type` enum 與預設 `activation_mode` 對應**：

| route_type | 預設 activation_mode | 範例 prefix |
|---|---|---|
| `bootstrap` | `always-on` | `route.bootstrap.*` |
| `runtime_core` | `always-on` | `route.runtime.{phase-machine, obligation-ledger, blocking-gates, recovery}` |
| `workflow` | `auto-detect` | `route.workflow.*` |
| `analysis` | `auto-detect` | `route.analysis.*` |
| **`intelligence`** | **`must-declare`（無預設，必須顯式宣告）** | `route.intelligence.*` —— **mixed layer，見下方特別說明** |
| `governance` | `on-demand` | `route.governance.*` |
| `constitution` | `on-demand` | `route.constitution.*` |
| `architecture` | `on-demand` | `route.architecture.*` |
| `feedback` | `advisory` | `route.feedback.*` |
| `metadata` | `on-demand` | `route.metadata.*`、`route.knowledge.*` |
| `ai_tools` | `on-demand` | `route.ai-tools.*`、`route.tools.*` |
| `models` | `advisory` | `route.models.*` |

##### Phase 0.2a-special — `intelligence` 為什麼 must-declare（v4 採納評審 #1）

第四輪評審指出 `intelligence -> advisory` 預設假設危險。具體 case：

- 使用者問「**幫我評審這個系統架構**」→ 命中 `route.intelligence.architectural-fit`
- 這是**主任務**（primary route），不是 advisory hint
- 但 `intelligence -> advisory` 預設會讓 `can_activate=false`，detector 不會鎖定，主任務變成永遠不會啟動

intelligence 層本質是 **mixed layer** —— 介於 analysis / workflow / governance 之間：
- 部分 atom 是 primary route（如 architectural-fit 用於架構評審任務）
- 部分 atom 是 secondary hint（如 engineering.heuristics 暗藏在實作任務中）
- 強制單一 default 必然在某一群誤判

**v4 解法**：`intelligence` 不享有自動推導，**每條 route 必須顯式宣告 `activation_mode`**。

```yaml
# 範例：primary intelligence route
- id: route.intelligence.architectural-fit
  route_type: intelligence
  activation_mode: auto-detect   # 顯式宣告，因為這是評審任務的主路由
  activation_triggers:
    activation_any_of:
      user_signals: [架構評審, architecture review, 評估架構]

# 範例：secondary intelligence route
- id: route.intelligence.engineering.heuristics
  route_type: intelligence
  activation_mode: advisory      # 顯式宣告，作為實作任務的 hint
```

**Lint rule（commit-msg validator 機械強制）**：

```
若 route.route_type == "intelligence" AND route.activation_mode 未宣告
  → commit reject，訊息："intelligence routes must explicitly declare
    activation_mode (one of: auto-detect / on-demand / advisory) —
    no automatic default to avoid mis-categorizing primary vs
    secondary intelligence atoms"
```

**Phase 0.2 必加 Audit Gate**：所有 `route.intelligence.*`（目前 7 條）在 v4 推進到 Phase 1 之前，**每一條必須個別決策 activation_mode**。產出 audit table 列為 Phase 0.2 acceptance criteria：

| Route ID | 用於什麼任務 | 是 primary 還是 secondary | 決定的 activation_mode |
|---|---|---|---|
| route.intelligence.architectural-fit | 系統架構評審 | primary | **auto-detect**（暫定，Phase 0.2 確認） |
| route.intelligence.requirements-cognition | 需求認知拆解 | primary | **auto-detect**（暫定） |
| route.intelligence.engineering.heuristics | 實作任務輔助 | secondary | **advisory**（暫定） |
| route.intelligence.engineering.agent-architecture | agent 設計參考 | mixed | **must-declare per use case** |
| route.intelligence.apk-analysis.atoms | APK 分析輔助 | secondary | **advisory**（暫定） |
| route.intelligence.apk-highest-leverage-path | APK 高槓桿路徑分析 | primary | **auto-detect**（暫定） |

> Phase 0.2 acceptance criteria 加一條：**intelligence audit table 已 user-reviewed，每條 route 有明確 activation_mode 決策**，否則 Phase 1+ 不得啟動。

##### Phase 0.2a-extensibility — 未來其他 mixed-layer types

`must-declare` 標記可推廣到其他未來出現的 mixed-layer types。預設機制不是「給每個 type 一個答案」，而是「給每個 type 一個適當的 layer policy」：
- 單一語意明確 → 預設 mode
- 混合語意 → must-declare
- Lint rule 通用化：表格中標 `must-declare` 的 type，commit-msg 都會檢查顯式宣告

**為什麼這比人工 classification 表好**：
- 新 route 作者最知道自己屬哪類，宣告 cost 低（一行 yaml）
- type → mode 對應 table 是一次性決策，鎖定後新 route 自動繼承
- 例外覆寫機制保留彈性（極少數 route 需要與 type 預設不同）
- 沒有「中央表 vs route file」雙寫 drift 風險

##### Phase 0.2b — `activation_mode` Capability Matrix（**v3 新增**，回應評審 #1）

每個 mode 用四個 capability bit 描述行為，避免「`advisory` 是 secondary hint 但具體會發生什麼」這類模糊：

| Capability \ Mode | `always-on` | `auto-detect` | `on-demand` | `advisory` |
|---|---|---|---|---|
| `can_preload` | ✅ true | false | false | false |
| `can_activate`（rule match → 鎖定 RuntimeContext.ActiveRoute） | n/a (always loaded) | ✅ true | true（only on explicit user invocation） | ❌ **false** |
| `can_reinforce`（hit 後提升另一個 active route 的 confidence / 並列 DetectedRoutes） | n/a | true | false | ✅ **true** |
| `can_conflict`（多 hit 時參與 `workflow-routing.md` 歧義裁決） | n/a | ✅ true | false | ❌ **false** |
| `can_suggest_promotion`（如果 advisory route 持續無 auto-detect 同伴 hit → 建議升級成 auto-detect） | n/a | n/a | n/a | ✅ **true** |
| `requires_activation_triggers` | false | ✅ **true** | false | true（弱訊號 OK） |

**範例落地**：
- `route.workflow.travel-planning`（auto-detect）：can_activate ✅，can_conflict ✅ — 命中後鎖定 ActiveRoute，與其他 auto-detect 多 hit 時走 conflict resolver
- `route.intelligence.architectural-fit`（advisory）：can_activate ❌，can_reinforce ✅ — 命中時不單獨鎖定，但若同時有 auto-detect route hit，會被加入 DetectedRoutes 作為 secondary context；若持續單獨 hit，記 `suggest_promotion` 建議升級
- `route.runtime.phase-machine`（runtime_core → always-on）：can_preload ✅ — bootstrap 階段自動載入，detector 不參與
- `route.governance.routing-signal`（governance → on-demand）：can_activate（only user invocation） — 使用者明確問「routing signal 怎麼設」才載入

##### 產出

- [ ] `routing-registry.yaml` header 加：
  - `route_type_spec`（12 enum + 對應預設 activation_mode 表）
  - `activation_mode_spec`（4 mode + capability matrix）
- [ ] 每條 route 加 `route_type:` 欄位（required）+ optional `activation_mode:`（override 用，預設由 type 推導）
- [ ] 預估初始分布（依 route_type 自動推導）：always-on ~12、auto-detect ~12、on-demand ~25、advisory ~8

### Phase 1 — Detector Schema 定義（two-phase activation 破環依賴）

在 `routing-registry.yaml` 擴 `activation_triggers` schema，**明確分離 pre-Read / post-Read 訊號**：

```yaml
activation_triggers:
  # ─────────────────────────────────────────
  # Phase 1: Activation signals (pre-Read, 必須在 Read File 之前可取得)
  # ─────────────────────────────────────────
  activation_any_of:
    user_signals: [行程, itinerary, 旅遊]      # 對話文字（既有，重新命名）
    context_signals:                            # 檔名 / 路徑 pattern
      - "docs/*行程*.md"
      - "docs/[0-9]{8}-*.md"
  # 任一命中 → activate workflow，detector 鎖定 active_route。

  # ─────────────────────────────────────────
  # Phase 2: Reinforcement signals (post-Read, 強化驗證)
  # ─────────────────────────────────────────
  reinforcement_any_of:
    artifact_signals:                           # 已讀檔案內容 pattern
      - "Day [0-9]+"
      - "御朱印"
      - "MapCode"
  # 已 Read 的檔案命中 reinforcement_any_of → 提升 confidence、不單獨用於 activate；
  # 若 Phase 1 已 activate 則作為「方向正確」確認；若 Phase 1 miss 但 Phase 2 hit，
  # 視為「late-detected」事件，記日誌但不 retroactively rewrite history。

  task_intents: [travel-planning, itinerary]   # 既有，保留
```

**為什麼要分兩階段（破環依賴）**：第二輪評審指出 — `artifact_signals` 需要先 Read 才有內容，但 workflow activation 的目的之一是**強制 agent 在 Read 之前先讀 workflow primary_source**。若 detector 依賴 artifact_signals，就會出現「要先 Read 才知道該讀 workflow，但 workflow 又要求 Read 前先讀 workflow」的循環依賴。

Phase 1 訊號（user_signals + context_signals）來源 pre-Read 即可取得：
- `user_signals`：使用者對話文字（agent inbox 直接讀）
- `context_signals`：檔名 / 路徑（從 user 提及的檔名、open files list、cwd 等取得，不需 Read 內容）

Phase 2 訊號（artifact_signals）：在 agent 自然 Read 過程中累積，作為「我猜對了」的確認或「late-detected」訊號，不作為 activation gate。

**Deterministic rule**：任一 `activation_any_of` 子陣列內任一條 hit → activate。不加權、不算分。

**向後相容**：舊格式（直接 `user_signals: [...]`）仍接受，視為 `activation_any_of.user_signals`。

產出：
- [ ] `routing-registry.yaml` schema 更新 + 文件（包含 two-phase 說明）
- [ ] `governance/workflow-activation-engine.md` 新建（philosophy + schema spec + circular dependency rationale）
- [ ] 7 個既有 `activation_triggers` 路由不動（schema 向後相容）

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

### Phase 4 — RuntimeContext State（in-memory first, SQLite deferred）

#### Phase 4.0 — In-memory RuntimeContext（本 plan 唯一交付，YAGNI 原則）

**不**直接落 `runtime.db`。先在 Go process memory 維護：

```go
// scripts/ai-skill-cli/internal/app/runtime_context.go
type RuntimeContext struct {
    ActiveRoute       string         // 鎖定的 route id（單一）
    DetectedRoutes    []string       // detector 命中的全部 route（含 active 與 advisory）
    DetectionSource   DetectionSig   // 哪些 signal axis 觸發
    ActivatedAt       time.Time
    LastReinforcedAt  time.Time      // Phase 2 reinforcement 最近一次 hit
    Status            RuntimeStatus  // detected | locked | no-match | manually-overridden
}

type DetectionSig struct {
    UserSignalHits    []string
    ContextSignalHits []string
    ArtifactReinforce []string  // Phase 2 hits, 不參與 activation
}
```

**為什麼不直接落 SQLite**：第二輪評審 — `workflow_sessions` 表的需求假設「跨 session replay / 統計 / 分析」，但這些需求**還沒驗證**。先存 memory 滿足 in-task 即時讀寫，把 schema 固化延後到出現實際需求（例如使用者要求「列出我這個月做了多少次 travel-planning 任務」之類）。

Lifecycle（簡化，**移除 implicit keyword-drift invalidation**）：

1. **Task start detection**：first substantive user message 後跑 detector，結果寫 `RuntimeContext`。**"Substantive" 用 intent vocabulary 判定，不用字數**（v3 採納評審 #2）：

   ```
   substantive(msg) :=
     contains_any(msg, domain_nouns) OR contains_any(msg, action_verbs)

   domain_nouns := {旅遊, 行程, 規劃, 架構, 設計, API, APK, 分析, 評審,
                    governance, workflow, ...}  # 從 routing-registry 所有
                                                # activation_any_of.user_signals
                                                # 自動聚合，registry 改即同步
   action_verbs := {幫我, 規劃, 寫, 做, 找, 比較, 設計, 評估, 檢查, 修, ...}
   ```

   範例：
   - `幫我規劃下關行程`（8 chars）→ contains `規劃` + `行程` → ✅ substantive
   - `hi 早安`（5 chars）→ 無 domain noun / action verb → ❌ not substantive
   - 字數門檻被淘汰因為 8 字中文 message 已可表達完整 task intent，舊 ≥20 chars 規則會誤殺。
2. **Topic shift detection**（兩種，**取消 implicit drift**）：
   - ✅ **顯式 pivot**：user message 含 sentinel（`換任務 / 現在我要 / new task / switch to / 換個話題` 等）→ invalidate + 重跑 detector
   - ✅ **Manual override**：user 顯式說「用 X workflow / 跟我做 X」→ 直接覆寫 active_route
   - ❌ **取消**：連續 N turn keyword 流失 → invalidate
3. **Why no implicit drift**：第二輪評審指出 down-drill 場景會誤殺。例：旅遊規劃會 drill into「這個停車場可以嗎 / 那住宿呢 / 這神社值得去嗎 / 這邊御朱印呢」，連續多 turn 不會出現「行程 / Day1」等 trigger keyword，但仍是同一 workflow。Implicit drift 會把這類正常 drill-down 誤判成 topic shift。
4. **替代方案**：keyword 流失只記 `LastReinforcedAt`，可選擇性 warning（"已 N turn 未見此 workflow 強訊號，是否仍在此 task？"），但**不自動 invalidate**

#### Phase 4.1 — SQLite Persistence（**deferred, conditional**）

**不在本 plan scope**。等以下條件之一成立才啟動 follow-up plan：
- 使用者明確要求跨 session 行為（"記住我上次做的 workflow"）
- 出現分析需求（"統計 detection accuracy"）
- 出現 multi-agent handoff 需求（"另一個 agent 接手要知道 active workflow"）

若觸發，新 plan 設計 `workflow_sessions` 表時可直接 mirror `RuntimeContext` 欄位，migration cost 低。

產出：
- [ ] `runtime_context.go` + unit tests
- [ ] `hooks.go` 整合：detector 寫入 RuntimeContext，其他 validator 從 RuntimeContext 讀
- [ ] `ai-skill runtime workflow-context` CLI subcommand（dump 當前 in-memory state）
- [ ] Lifecycle 文件化（governance/workflow-activation-engine.md §lifecycle）
- [ ] 明確記錄「SQLite 延後」決策（Phase 4.1 不在 scope）

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

### Phase 6 — Failure Pattern + Discovery Feedback Loop

#### Phase 6.0 — Failure Pattern 記錄

- [ ] 新建 `enforcement/failure-patterns/workflow-detector-missing.md` —— 記錄 2026-05-31 失效為 systemic gap，並把 Detector 設計指回本 plan

#### Phase 6.1 — Discovery → Detector Feedback Loop（**新增，採納第二輪評審**）

第二輪評審指出原版「Discovery 完全切開」過度切割。正確關係：

```
User Request
    ↓
Detector (Phase 1 signals → Phase 2 reinforcement)
    ↓
[hit]                          [miss]
    ↓                              ↓
RuntimeContext.ActiveRoute    Capability Discovery
    ↓                              ↓
Execution                     graph traversal → 找到相關 governance / workflow / intelligence atom
                                   ↓
                              suggest: new route candidate
                                   ↓
                              Registry growth proposal
                                   ↓
                              (user / future plan 決定是否新增 route + triggers)
```

**Why this matters**：
- Detector miss 不代表「沒有 workflow 可用」，可能是 Registry 還沒收錄此 task type
- Discovery 跑 graph traversal 找出可能相關的 capability atom
- 若 Discovery 一致指向同一群 capability → 暗示應新增一個 route
- 這形成 **Registry 自我成長機制**：使用越多、coverage 越廣

範例：未來 user 開始做「AI Agent Governance Audit」，現有 Registry 沒對應 route → Detector 全 miss → Discovery fallback → 找到 governance/architecture/compliance 群 atom → 提案新增 `route.workflow.governance-audit` candidate。

產出：
- [ ] 更新 `governance/lifecycle/capability-discovery-philosophy.md` —— 加章節「Discovery → Detector Feedback Loop」，明確：
  - Detector 處理 known route 的 known trigger（cheap, deterministic, per-task）
  - Discovery 處理 unknown capability（expensive, exploratory, **only fires on detector miss**）
  - Discovery 輸出可以 propose Registry growth（candidate route + suggested triggers）
- [ ] 在 detector miss path 加 fallback hook 呼叫 Discovery（但不阻擋執行流程，warn + continue）
- [ ] 新建 `runtime/router/route-candidate-proposals.yaml`（pending proposals 暫存區）—— 採 **occurrence tracking** schema 防垃圾場（v3 採納評審 #3）：

  ```yaml
  # runtime/router/route-candidate-proposals.yaml
  schema_version: 1
  proposals:
    - candidate_id: governance-audit          # slug，未 promote 前非 route.* id
      first_seen: 2026-05-31T10:00:00Z
      last_seen:  2026-06-05T14:30:00Z
      occurrence_count: 7                     # Discovery 多少次指向同樣 capability 群
      detected_capabilities:                  # Discovery 找到的相關 atom
        - governance/lifecycle/...
        - architecture/...
        - intelligence/governance-...
      suggested_user_signals: [audit, governance audit, compliance check]
      suggested_route_type: workflow
      status: accumulating                    # accumulating | ready_for_review |
                                              # promoted | rejected | stale
  ```

  **Promotion rules（防止一次性需求污染 Registry）**：

  | 狀態轉換 | 條件 |
  |---|---|
  | `accumulating` → `ready_for_review` | `occurrence_count >= 5` AND `last_seen` 在過去 30 天內 |
  | `accumulating` → `stale` | `occurrence_count < 5` AND `last_seen` 超過 60 天 → 自動 archive，不再列入活躍清單 |
  | `ready_for_review` → `promoted` | User / governance review approve，proposal 內容寫入 `routing-registry.yaml` 成為正式 route，proposal 從 yaml 移除 |
  | `ready_for_review` → `rejected` | User / governance review reject（例：太細、與既有 route 重疊），記 `rejected_reason` 後 archive |

  **CLI 輔助**：`ai-skill router proposals list --status ready_for_review` —— 只顯示真正值得 review 的，避免使用者面對 100 條一次性 proposal 的垃圾場。
- [ ] 新建 `ai-skill router proposals {list, promote, reject, gc}` CLI subcommands

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
| Q2 | Detector 的「first substantive message」定義 —— 純打招呼算嗎？ | **resolved (v3) → Phase 4.0 lifecycle**：採 intent vocabulary（domain_nouns ∪ action_verbs）判定，**不**用字數門檻。domain_nouns 自動從 routing-registry 的 `activation_any_of.user_signals` 聚合，registry 改即同步。 |
| Q3 | Conflict resolution 多 route 命中時自動選還是 prompt user？ | **resolved** (v1) → 不自動選，注入 reminder 讓 agent 走 `workflow-routing.md`（v2/v3 close-out 誤標 still-open，v4 修正） |
| Q8 | `route.intelligence.*` 全部 7 條 audit 結果是什麼？ | **new (v4)，still-open**：Phase 0.2a-special audit table 需 user 逐條 review，產出 acceptance criteria 之一。表中暫定值僅供討論。 |
| Q4 | `workflow_sessions` TTL？跨 session 是否保留？ | **resolved → Phase 4.0**：本 plan 不落 SQLite，TTL 等於 in-memory RuntimeContext 生命週期（process scope）。跨 session 持久化延後到 Phase 4.1 follow-up plan。 |
| Q5 | Detector miss 是否 fallback 到 Discovery？ | **resolved → Phase 6.1**：採納第二輪評審，**Yes** 但有限制 —— Discovery 只在 detector miss 時 fire（不是 per-turn），結果寫 `route-candidate-proposals.yaml` 供未來 review，不阻擋當前執行流程。 |
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

2026-05-31 session：使用者連續三次追問才暴露 `route.workflow.travel-planning` activation gap。

**Round 1 第三方對話建議**（採納於 v1）：
- 拆解 Discovery vs Detector
- 放棄 weighted scoring
- 改用 deterministic + workflow_sessions

**Round 2 第三方架構評審**（採納於 v2）：

| # | 評審論點 | 採納 | 對應修改 |
|---|---|---|---|
| 1 | `workflow_sessions` 太早進 SQLite，先 memory 即可 | ✅ | Phase 4 拆 4.0 in-memory（本 plan）+ 4.1 SQLite deferred |
| 2 | Route classification 過於二元，需 4 enum `activation_mode` | ✅ | Phase 0.2 從 3 class 改 4 mode（always-on / auto-detect / on-demand / advisory） |
| 3 | 5-turn keyword drift invalidation 會誤殺 down-drill | ✅ | Phase 4 lifecycle 移除 implicit drift，只保留 explicit pivot + manual override |
| 4 | `artifact_signals` 依賴 Read 產生循環依賴 | ✅ | Phase 1 schema 拆 two-phase：activation_any_of（pre-Read）+ reinforcement_any_of（post-Read） |
| 5 | Discovery 不該完全切開，應有 feedback loop 反餵 Registry | ✅ | Phase 6.1 新增 Discovery → Detector feedback loop + `route-candidate-proposals.yaml` |

本 plan v2 接受全部 Round 2 建議。Round 2 評分（user 給）：
- Problem Identification: A
- Root Cause Analysis: A
- System Boundary: A
- Implementation Complexity: B+ → 目標 v2 提升
- Future Maintainability: B → 目標 v2 提升（in-memory 延後固化、activation_mode 給彈性、Discovery feedback loop 讓 Registry 自我成長）

**Round 3 第三方架構評審**（採納於 v3）：

| # | 評審論點 | 採納 | 對應修改 |
|---|---|---|---|
| 1 | `advisory` 定義模糊（"secondary hint" 不夠 actionable） | ✅ | Phase 0.2b 新增 6-bit Capability Matrix（can_preload / can_activate / can_reinforce / can_conflict / can_suggest_promotion / requires_activation_triggers），每個 mode 明確標 ✅/❌ |
| 2 | "First substantive message ≥ 20 chars" 規則會誤殺短中文意圖 | ✅ | Phase 4.0 lifecycle 改用 intent vocabulary 判定（domain_nouns ∪ action_verbs），字數門檻取消。"幫我規劃下關行程"（8 chars）現在能正確識別。 |
| 3 | `route-candidate-proposals.yaml` 易變垃圾場（一次性需求污染） | ✅ | Phase 6.1 加 `occurrence_count` + `first_seen` / `last_seen` + 4 狀態機（accumulating → ready_for_review / stale → promoted / rejected）+ promotion threshold (count ≥ 5)。CLI `proposals list --status ready_for_review` 只顯示值得 review 的。 |
| 4 | 57 routes 人工分類今天可行，120/300 後維護地獄 | ✅ | Phase 0.2a 引入 **self-declared `route_type`**（12 enum）+ type → activation_mode 預設對應表。新 route 作者宣告 type 即自動帶分類，無中央維護表 drift 風險。 |

**Round 3 核心轉變**（user 點評整體結論）：
> v1 → 「發現 detector 不存在 → 補 detector」
> v3 → 「Known Capability → Detector / Unknown Capability → Discovery / Discovery → Registry Growth / Runtime → In-Memory Context / Workflow → Activation Lifecycle」

不再只是補洞，而是形成 Runtime 生態。Future maintainability 從 v2 的 B 提升目標：靠 `route_type` 自動分類 + Discovery feedback loop 讓 Registry 自我成長，使 framework 不需中央維護就能 scale。

**Round 4 第三方架構評審**（採納於 v4）：

| # | 評審論點 | 採納 | 對應修改 |
|---|---|---|---|
| 1 | `intelligence -> advisory` 預設危險，intelligence 是 mixed layer（part primary, part secondary） | ✅ | Phase 0.2a 將 intelligence 從「預設 advisory」改 `must-declare`（無預設）+ 機械 lint rule。Phase 0.2a-special 新增 intelligence audit gate，7 條 route 逐條決策成為 Phase 1 unlock 條件。Phase 0.2a-extensibility 將 must-declare 機制通用化。 |

**Round 4 評分**（user 給）：

| 項目 | v2 | v3 |
|---|---|---|
| Problem Identification | A | A |
| Root Cause Analysis | A | A |
| Runtime Design | B+ | A- |
| Scalability | B | A |
| Maintainability | B | A- |
| Registry Evolution | C+ | A |
| **Overall** | **B+** | **A-** |

User 評語：
> v3 已經從「補 travel-planning detector」進化成「建立 Capability Routing Runtime」。剩下最大架構風險不在 Detector，而是在「intelligence 是否真的能全部預設 advisory」這個分類假設。

**v4 對此風險的處置**：拒絕單一 default，引入 must-declare + Phase 0.2 audit gate，把分類決策從「framework 替你猜」轉為「每條 route 自己負責宣告」。這也是「self-declaring routes」原則的徹底化 —— 連 mode 都不靠 type 猜。

## Companion References

- `governance/lifecycle/capability-discovery-philosophy.md` —— Discovery 機制（與本 plan 互補）
- `workflow/workflow-routing.md` —— 多 route 命中時的歧義裁決（Stage 2 conflict resolver）
- `enforcement/dependency-reading.md` §Workflow 編排 —— blocking activation 行為強制（本 plan 升級為機械強制）
- `enforcement/failure-patterns/bootstrap-bypass-on-resume.md` —— PreToolUse + transcript scan 模式範例（detector 採同模式）
