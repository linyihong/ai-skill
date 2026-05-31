# Mechanical Enforcement Registry

**Status**: `draft-v4.1`
**世代**：Gen 3 runtime hardening（**Meta Governance**：framework self-audit layer，非個案修補）
**建立日期**：2026-05-31
**最後更新**：2026-05-31（v4.1 — Phase 5.x 新增 Hook Injection Economics 作為 inaugural self-governance test case；併入而非另開 plan）
**Priority**：**P1**（v2 起）—— 此 plan 處於 architectural prevention 層級（Prevent > Detect > Repair）；child plans 仍在 Detect/Repair 層級
**Round 9 評分（user）**：8.8 / 10 —— 扣分點：runtime_observed 缺失、pending 混兩種、registry 自身治理未定義。v4 處理全部三項

**Child plans (instances of this meta-pattern)**：
- [`2026-05-31-1900-workflow-activation-engine.md`](2026-05-31-1900-workflow-activation-engine.md) — sanitizes & operationalizes `routing-registry.yaml` rules with `detector.go` executor
- [`2026-05-31-2000-mechanical-sanitization-validator.md`](2026-05-31-2000-mechanical-sanitization-validator.md) — sanitizes & operationalizes `sanitization.yaml` rules with `validateSanitizationOnWrite/Commit` executors

**Empirical trigger**：2026-05-31 session 連續暴露 5 個同模式 bug（bootstrap bypass、workflow activation、capability discovery、sanitization、intelligence classification），第三方架構評審指出這不是 5 個獨立 bug，而是同一個 meta-pattern：

> **Knowledge Layer 有規則，Runtime Layer 沒執行器。**

每次發生時都要等使用者半年後追問才會被抓出來。本 plan 把這個 meta-pattern 變成 framework 第一公民，讓未來的同模式 bug 在 compile time 就被擋下，不再依賴使用者偶然發現。

> 本 plan 不修任何個案。它**建立 invariant**：所有 enforcement rule 在 register 時必須 link 到一個 executor（或顯式宣告 `behavioral_only: true` 並附 rationale）。compile pipeline 強制此 link，違反即 fail。

---

## Decision Rationale

### Empirical Evidence（meta-pattern 體檢）

2026-05-31 session 暴露的 gap 列表，每個都符合 "rule exists, executor missing" 模式：

| 問題 | Knowledge Layer 規則 | Runtime Layer 執行器 | 發現方式 |
|---|---|---|---|
| Bootstrap bypass on resume | `core-bootstrap.yaml` `obligation.bootstrap.receipt` | ✅ 已補（read-log gate + PreToolUse） | 過去半年累積 failure pattern 才補 |
| Workflow activation | `routing-registry.yaml` `activation_triggers` | ❌ 不存在（child plan 待實作） | 使用者追問 3 次 |
| Capability discovery fallback | `capability-discovery-philosophy.md` | ❌ 不存在 | 同上 |
| Sanitization on canonical Write | `sanitization.md` + `reusable-guidance-boundary.md` | ❌ 不存在（child plan 待實作） | 使用者追問 5 次 |
| Intelligence route classification | routing-registry 已收 route id | ❌ activation semantics 缺 | 同上 |

**共同根因**：Knowledge layer 與 Runtime layer 之間**沒有結構性 binding**。新增規則時不會被強制問「executor 在哪」，所以容易出現「規則寫好、放著、忘記補執行器」。

### Why Meta-pattern, Not Just Individual Fixes

兩個 child plan（activation-engine + sanitization-validator）會解決 2 個 instance。但：
- 還有 ≥ 3 個未列入的 instance 待補（capability discovery、intelligence classification、未來其他）
- 即使全部補完，**新增的下一個 rule** 仍可能再犯同模式
- 沒有結構性檢查，每個 instance 都要等使用者追問才被發現

正確 fix：**讓 framework 自己回答**「目前有哪些規則存在，但沒有對應的 runtime executor？」這就是本 plan。

### Architectural Framing — Missing Layer 2.5 / Meta Governance

第七輪評審指出 Ai-skill 既有架構是三層：

```
Layer 1  Knowledge       (enforcement/, governance/, workflow/, ...)
Layer 2  Runtime         (scripts/ai-skill-cli/, hooks.go, runtime.db)
Layer 3  Governance      (constitution/, architecture/, plans/)
```

但缺一層：

```
Layer 2.5  Coverage Verification / Meta Governance   ← NEW (本 plan 建立)
            Rule ←binding→ Executor ←evidence→ Verification
            的結構性驗證層
```

第八輪評審把這層定位升格為 **Meta Governance / Framework Self-Audit Layer**：
- Layer 1-3 管「框架的內容」（規則寫了什麼、runtime 跑什麼、治理決策）
- Layer 2.5 管「框架本身是否真的做到它說會做的」（Governance of Governance）

沒有 Layer 2.5 就沒有「rule 寫好、executor 沒接」「executor 寫了但漏覆蓋某些 instance」的結構性偵測機制。本 plan 不是「補一個 executor」（child plan 的工作），而是**建立 Layer 2.5 本體**。

### Decision

建立 **Mechanical Enforcement Registry**（**Rule Class 級**，**非 rule instance 級**）+ **Coverage Report** + **Compile-time Lint**：

```
┌─────────────────────────────────────────────────────────────────┐
│                    Knowledge Layer (rules)                        │
│  enforcement/*.yaml, runtime/*.yaml, routing-registry.yaml, ...  │
└──────────────────────────────────┬──────────────────────────────┘
                                   │ declares
                                   ▼
┌─────────────────────────────────────────────────────────────────┐
│  enforcement-registry.yaml — Rule ↔ Executor binding table       │
│    - id: workflow_activation                                     │
│      rule_source: routing-registry.yaml                          │
│      executor: scripts/ai-skill-cli/.../detector.go              │
│      executor_symbol: validateWorkflowActivation                 │
│      enforcement_layer: PreToolUse + commit                      │
│    - id: sanitization                                            │
│      rule_source: enforcement/sanitization.yaml                  │
│      executor: scripts/ai-skill-cli/.../hooks.go                 │
│      executor_symbol: validateSanitizationOnCommit               │
│      enforcement_layer: commit (block) + PreToolUse (warn)       │
│    ...                                                            │
└──────────────────────────────────┬──────────────────────────────┘
                                   │ compile-time lint
                                   ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Runtime Layer (executors)                      │
│        hooks.go, detector.go, ai-skill CLI, runtime.db           │
└─────────────────────────────────────────────────────────────────┘
```

### Design Principles

| Decision | Rationale |
|---|---|
| **Rule Class 級 binding，非 rule instance 級**（v2 採納評審 #2） | 22 條 `obligation.commit.*` 不應產生 22 條 binding entries；都綁同一個 dispatcher 變成 registry 自我膨脹。改用 `rule_class` 抽象：sanitization / workflow_activation / bootstrap_integrity / routing / commit_governance / discovery / ... 估計 ~20-25 個 class，registry 規模可控 |
| **6-value coverage enum**（v4 採納評審 #2，從 5 升 6） | `mechanical` / `behavioral_only` / `not_mechanizable` / **`pending_implementation`** / **`research_required`** / `deprecated`。v4 把 v3 `pending` 拆兩種：**pending_implementation**（知道怎麼做、child plan 在跑）vs **research_required**（知道應該機械化、但還不知道怎麼做）。同樣是 7 個未完成，前者治理訊號是「快完成」，後者是「需要思考」 |
| **Each coverage status 不同 metadata 要求** | `behavioral_only`: rationale + sunset_decision；`not_mechanizable`: rationale + objective_validation_impossible_because；`deprecated`: rationale + `replaced_by` (or `removal_date`)；`pending_implementation`: child_plan + target_promotion；`research_required`: rationale + research_questions + estimated_unblock_timeline。Compile lint 校驗 schema |
| **`verification` 維度與 coverage 正交，含 runtime 層**（**v4 採納評審 #1**） | v3 verification 涵蓋 symbol / scenario / regression，但缺最終一層 **`runtime_observed`**。v4 新增：scenario_exists 不等於 production reality_exists（57 routes 100% scenario 涵蓋，可能實際半年都沒被 detector 觸發過 → 真實世界要嘛沒人用、要嘛 detector 壞了）。Runtime metrics 收集 last_30d 觀測，Coverage Report 顯示 Declared vs Covered vs Runtime Observed 三層 |
| **Registry Self-Governance**（**v4 採納評審 #3，框架閉環最後一塊**） | Layer 2.5 自己也需要治理：誰可以改 rule_class status？mechanical→behavioral_only 算 demotion 要 ADR；pending_implementation→mechanical 算 promotion 要 coverage_evidence pass；deprecated removal_date 屆期需 governance 決定 actually remove vs extend。詳見 Phase 5 |
| **Compile-time lint，非 runtime lint** | `ai-skill runtime compile + refresh` 跑 lint。新增 lint 條件涵蓋 verification 維度 |
| **CLI `ai-skill enforcement coverage` 是主要產出**（v2 評審 #4 強調） | 不只 audit 既有狀態，更重要的是**強制新規則寫作者回答 coverage 問題**。v3 加 verification 後，回答的問題從「mechanical 嗎」變成「mechanical 且 verified 嗎」 |
| **不重寫既有 rule**，只加 binding 層 | 既有 enforcement / runtime / governance markdown + yaml 不動，registry 是 cross-cutting 索引。降低 risk |

### 為什麼這個 plan 比 child plans 更重要

| 角度 | Child plans | This plan |
|---|---|---|
| Scope | 1 個 rule + 1 套 executor | N 個 rule × M 個 executor 的 binding 系統 |
| Detection latency | 已發生才修 | Compile time 預防 |
| 未來 cost | 每個 instance 重複付 | 一次性建好 invariant，新 rule 自動 enforce |
| Governance value | 補洞 | 建 framework property |

---

## Architecture Compatibility Preflight

| 欄位 | 內容 |
|---|---|
| Candidate files | 新建 `enforcement/enforcement-registry.yaml`（canonical binding 表）、新建 `enforcement/enforcement-registry.md`（companion philosophy）、`scripts/ai-skill-cli/internal/compile/`（lint）、`scripts/ai-skill-cli/cmd/enforcement.go`（CLI subcommand）、`runtime/runtime.db`（新增 `enforcement_bindings` projection） |
| Source-of-truth | `enforcement/enforcement-registry.yaml` 為 canonical binding 表。各 rule yaml 仍是各 rule 的 canonical source；executor source 在 Go code 不變。registry 是 cross-link |
| Compiler / generated surfaces | `ai-skill runtime compile` 流程加 enforcement-registry lint 階段，fail → block compile |
| Layer responsibility | enforcement-registry 屬 enforcement layer（meta-governance）；lint 屬 runtime/compiler layer；CLI 屬 ai-skill-cli layer |
| 與現行架構衝突 | 無。本 plan 不改既有 rule semantic，只加 cross-cutting binding |
| `runtime.db` 影響 | 新增 `enforcement_bindings` table + projection；compile 失敗條件擴展 |

---

## Phase Plan

### Phase 0 — Preflight

- [ ] §Open Questions 處置
- [ ] 確認既有 `enforcement/`、`runtime/`、`governance/` 的 yaml 是否都已有 stable id field 可作為 registry key
- [ ] 確認 `hooks.go` validator dispatcher 已有可枚舉的 registry 結構（已知有 11+ commit-msg validators 註冊）

### Phase 1 — Rule Class 識別（**v2 改：不是 instance audit**）

第七輪評審指出原 v1 Phase 1 計畫掃 150+ rule instance 是 maintenance trap。v2 改為**識別 rule class**（~20-25 個），每個 class 含多個 instance 但共用 binding。

#### Phase 1.1 — 從 source 抽 rule class

逐目錄盤點，把所有 rule instance 歸類為 ~24 個 class：

| 候選 rule_class | source patterns | 預估 instance 數 |
|---|---|---|
| `bootstrap_integrity` | `runtime/core-bootstrap.yaml` per_session_obligations | 1 |
| `cognitive_mode_governance` | `runtime/cognitive-modes*.yaml`、per_turn_obligations | 1-3 |
| `commit_governance` | per_commit_obligations 全部（19 個 validator） | 19 |
| `workflow_activation` | `routing-registry.yaml` activation_triggers | 57 routes × 3 軸 |
| `capability_discovery` | `governance/lifecycle/capability-discovery-philosophy.md` | 1 |
| `sanitization` | `enforcement/sanitization.yaml`（待建） | N patterns |
| `dependency_reading` | `enforcement/dependency-reading.md` 依賴表 | 多條 mapping |
| `linked_updates` | `enforcement/linked-updates.yaml` writeback gates | 1+ |
| `rule_weight` | `enforcement/rule-weight.md` P0/P1/P2/P3 | 1（policy 級） |
| `conversation_goal_ledger` | `.agent-goals/` lifecycle | 1 |
| `document_sizing` | `governance/document-sizing.md` | 1 |
| `tool_neutral_documentation` | `enforcement/tool-neutral-documentation.md` | 1 |
| `reusable_guidance_boundary` | `enforcement/reusable-guidance-boundary.md` | 1 |
| `failure_learning_system` | `enforcement/failure-learning-system.md` | 1 |
| `glossary_governance` | glossary system | 1 |
| `plan_governance` | `plans/README.md`、archival 規則 | 1+ |
| `evidence_hierarchy` | `enforcement/evidence-hierarchy.yaml` | 1 |
| `routing_registry_evolution` | route candidate proposals | 1 |
| `ontology_consistency` | route_type 等分類規則 | 1 |
| `bootstrap_entry_thinness` | bootstrap entry 規範 | 1 |
| `runtime_yaml_projection` | yaml ↔ runtime.db projection 規則 | 1 |
| `cli_doc_sync` | CLI ↔ doc sync | 1 |
| `markdown_yaml_sync` | md ↔ yaml sibling sync | 1 |
| `intelligence_classification` | intelligence/analysis activation_mode | 1 |

**預估 24 個 rule_class**（最終視盤點而定）。每個 class 一條 registry entry，不是每個 instance 一條。

#### Phase 1.2 — Executor 盤點

掃 `scripts/ai-skill-cli/internal/app/hooks.go` 等 runtime code，列每個 executor symbol。**重要**：以 dispatcher 或 hook handler 級為粒度，不展開內部 helper。預估 ~25-30 個外部可見 executor。

#### Phase 1.3 — 第一次 Coverage Matrix

每個 rule_class 標記 4-value coverage status：
- `mechanical` — 已有 executor 且 enforcing
- `behavioral_only` — 故意不機械化，需 `rationale` + `sunset_decision`
- `not_mechanizable` — 永遠不該機械化（主觀 / 無客觀 validation），需 `rationale` + `objective_validation_impossible_because`
- `pending` — 應該機械化但尚未實作，需指向 implementation plan

預估 v1 分布（待 Phase 1.3 確認）：
- mechanical: ~10（包含已實作的 bootstrap_receipt、commit_governance 等）
- behavioral_only: ~4
- not_mechanizable: ~2-3
- pending: ~7-8（含 workflow_activation、sanitization、capability_discovery 等）

remoteness 從原 v1「掃 150+ rule + 70% orphan」變成「~24 class + 約 1/3 pending」，可管理。

### Phase 2 — 建立 `enforcement-registry.yaml`（rule_class 級）

```yaml
# enforcement/enforcement-registry.yaml
schema_version: 2          # v2: rule_class 級而非 rule instance 級

rule_classes:
  # ─── coverage: mechanical ──────────────────────────────────────
  - id: bootstrap_integrity
    coverage: mechanical
    source_files:
      - runtime/core-bootstrap.yaml#per_session_obligations
    executors:
      - file: scripts/ai-skill-cli/internal/app/hooks.go
        symbol: validateBootstrapReceiptPresent
        hook_phase: PreToolUse
        block_or_warn: block
    rationale: |
      Receipt 是 session integrity anchor，bypass 等同跳過必讀規則。

  - id: commit_governance
    coverage: mechanical
    source_files:
      - runtime/core-bootstrap.yaml#per_commit_obligations
    executors:
      # 注意：19 個 commit-msg validators 共用 dispatcher，registry 不展開列每個
      - file: scripts/ai-skill-cli/internal/app/hooks.go
        symbol: runCommitMsgHook
        hook_phase: commit-msg
        instance_count: 19
        block_or_warn: block
    rationale: |
      所有 commit 階段 obligation 由 runCommitMsgHook 統一 dispatch。
      Rule_class 級 binding 避免 22 條 obligation.commit.* 各自一條 entry
      導致 registry 膨脹。

  # ─── coverage: pending (待 child plan 實作) ──────────────────
  - id: workflow_activation
    coverage: pending
    source_files:
      - knowledge/runtime/routing-registry.yaml#activation_triggers
    executors_planned:
      - file: scripts/ai-skill-cli/internal/app/detector.go
        symbol: DetectWorkflows
        hook_phase: PreToolUse + RuntimeContext write
        block_or_warn: block
    child_plan: plans/active/2026-05-31-1900-workflow-activation-engine.md
    target_promotion: child_plan Phase 3-5 完成後改 coverage=mechanical

  - id: sanitization
    coverage: pending
    source_files:
      - enforcement/sanitization.yaml
    executors_planned:
      preflight:
        symbol: validateSanitizationOnWrite
        block_or_warn: warn
      commit:
        symbol: validateSanitizationOnCommit
        block_or_warn: block
    child_plan: plans/active/2026-05-31-2000-mechanical-sanitization-validator.md
    target_promotion: child_plan Phase 3 完成後改 mechanical

  # ─── coverage: behavioral_only ───────────────────────────────
  - id: capability_discovery
    coverage: behavioral_only
    source_files:
      - governance/lifecycle/capability-discovery-philosophy.md
    rationale: |
      Discovery 是 detector miss 後 fallback 探索，per-turn 強制成本爆炸。
      正確做法是 detector 完成後從 hooks 觸發，不獨立成 mechanical executor。
    sunset_decision:
      revisit_when: workflow_activation child plan Phase 6.1 land
      revisit_owner: framework maintainer
      success_criteria: |
        Detector miss path 已能呼叫 Discovery 並產出 route_candidate_proposals.yaml；
        屆時本 class 從 behavioral_only 改 mechanical（與 workflow_activation 共用 executor）

  - id: rule_weight
    coverage: behavioral_only
    source_files:
      - enforcement/rule-weight.md
    rationale: |
      P0/P1/P2/P3 排序需要在規則衝突情境下 case-by-case 判斷，
      無 single mechanical rule 可實作。
    sunset_decision:
      revisit_when: 出現 ≥ 3 個可機械偵測的 P0 違反模式時
      revisit_owner: governance maintainer
      success_criteria: 為每個可偵測模式抽出獨立 rule_class 後，本 class 收斂為純策略文件

  # ─── coverage: not_mechanizable ──────────────────────────────
  # (v2 新增 enum，回應評審 #3)
  - id: tool_neutral_documentation
    coverage: not_mechanizable
    source_files:
      - enforcement/tool-neutral-documentation.md
    rationale: |
      「文件是否中立」是寫作判斷，沒有可機械驗證的 boolean。
      可寫 lint 找特定工具名稱，但「neutral writing 結構」整體無法
      規約化。
    objective_validation_impossible_because: |
      Neutrality 涉及讀者預期、語境、隱性假設等主觀面向；機械 lint
      只能抓表面 token，無法判斷 framing 是否 tool-neutral。

  - id: rule_writing_quality
    coverage: not_mechanizable
    source_files:
      - enforcement/conversation-goal-ledger.md（writing quality 部分）
      - 其他規則的寫作品質
    rationale: |
      規則本身是否寫得清楚、是否易讀、是否避免歧義，無客觀 metric。
    objective_validation_impossible_because: |
      可讀性 metric（Flesch / 字數）只測表面，不測「規則是否真的能被 agent 正確理解」。
      若強行機械化會獎勵 gaming metric 的爛規則。

  # ... (continue for all ~24 rule_classes)

# ─── 5-value enum schema (v3) ─────────────────────────────
coverage_status_spec:
  mechanical:
    requires: [executors[].symbol_exists]
    verification_required: true   # v3 NEW: 同時驗 symbol 與 coverage_evidence
    lint_behavior: source 變更須同步 executors[].symbol 存在；verification 維度另查
  pending:
    requires: [child_plan, target_promotion]
    lint_behavior: child_plan 必須是 active plan 路徑
  behavioral_only:
    requires: [rationale, sunset_decision.revisit_when, sunset_decision.revisit_owner, sunset_decision.success_criteria]
    lint_behavior: 缺任一欄位 → compile fail
  not_mechanizable:
    requires: [rationale, objective_validation_impossible_because]
    lint_behavior: 缺任一欄位 → compile fail；附帶 governance review 時可挑戰是否真的 not_mechanizable
  pending_implementation:   # v4 NEW (拆自 v3 pending)
    requires: [child_plan, target_promotion]
    semantics: "知道怎麼做、實作中或排程中"
    lint_behavior: child_plan 必須是 active plan 路徑；過 target_promotion 預期日期未 promote → warning
  research_required:        # v4 NEW (拆自 v3 pending)
    requires: [rationale, research_questions, estimated_unblock_timeline]
    semantics: "知道應該機械化、但還不知道怎麼機械化"
    lint_behavior: |
      research_questions 必須 ≥ 1 條具體未解決問題；estimated_unblock_timeline
      過期未 promote 或 demote → governance review trigger
  deprecated:
    requires: [rationale, replaced_by_or_removal_date]
    lint_behavior: |
      `replaced_by` 必須指向 active rule_class id；或 `removal_date` 必須
      是未來 ISO-8601 日期。compile time 兩者皆缺 → fail；過 removal_date
      仍存在 → fail（強制不能無限期 deprecated）

# ─── verification dimension (v3 NEW + v4 加 runtime_observed) ─
# Coverage 講「我們選擇怎麼處理這條規則」；Verification 講「實作真的做到了嗎」。
# 兩者正交，但 verification 自己有 4 層階梯：
#   symbol → scenario → regression → runtime_observed
# 每一層比前一層更接近 production reality
verification_levels:
  symbol_exists:
    check: executor symbol 在 file 中存在
    lint: 缺 → compile fail（mechanical only）
  scenario_exists:
    check: 有對應 validation/scenarios/<class>/*.yaml
    lint: 缺 → compile warning（mechanical）or fail（v5 開始 fail）
  regression_exists:
    check: validation scenario 涵蓋已知歷史失誤
    lint: 缺 → compile warning
  runtime_observed:   # v4 NEW
    check: |
      runtime_metrics 顯示過去 N 天有實際觸發紀錄。例：
      workflow_activation 的 57 routes 應至少有 K% 在 last_30d 被 detector
      觸發過；K 由 rule_class 自己宣告 expected_observation_pct（如 30%）。
    rationale: |
      Scenario exists 不等於 production reality exists。一個 route 半年沒
      被觸發，要嘛實際沒人用（可考慮 deprecate），要嘛 detector 壞了
      （symbol exists 但 silent fail）。Runtime_observed 是 v3 verification
      與 production reality 之間的最後一個 gap。
    lint: |
      runtime_observation_pct < expected_observation_pct → coverage report
      warning；< expected/2 → governance review trigger（不直接 compile fail，
      因為 runtime 數據是 observational，需人工判斷）
  coverage_evidence:
    check: |
      若 rule_class 有 instance set（如 workflow_activation 有 57 routes），
      coverage_evidence.expected_instance_count 必須宣告，且
      validation_scenarios 集合覆蓋率 >= 80%。
    lint: 覆蓋率 < 80% → compile warning；< 50% → compile fail

# ─── v4 NEW: runtime_metrics schema ───────────────────────
runtime_metrics_spec:
  collection: |
    Runtime hook 在每次 executor 觸發時寫入 runtime.db.executor_observations
    table。registry CLI 從該表 aggregate 過去 N 天指標。
  schema:
    rule_class_id: string
    observation_window_days: int (default 30)
    last_observation_at: timestamp
    activation_count: int
    instance_breakdown:        # 對 instance set 較大的 class（workflow_activation
                                # 57 routes）按 instance 統計觸發次數
      - instance_id: route.workflow.travel-planning
        count: 23
      - instance_id: route.workflow.software-delivery
        count: 0  # ← Coverage Report 會 surface 這類「scenario 過但 runtime 0」
  expected_observation_pct:    # rule_class 自宣告期望觀測比例
    description: "instance set 內期望至少 X% 在 window 內被觸發"
    default: 30  # 保守預設；無 instance set 的 class 不適用
```

### rule_class 範例：含 verification + coverage_evidence（v3 新增 schema）

```yaml
- id: workflow_activation
  coverage: pending
  source_files:
    - knowledge/runtime/routing-registry.yaml#activation_triggers
  executors_planned:
    - symbol: DetectWorkflows
  child_plan: plans/active/2026-05-31-1900-workflow-activation-engine.md
  target_promotion: child_plan Phase 3-5 完成後改 coverage=mechanical

  # v3 NEW: 預先宣告 verification 期望，promotion 時 lint 會校驗
  coverage_evidence:
    expected_instance_count: 57   # routing-registry 目前 57 routes
    expected_instance_count_query: |
      grep -c '^  - id: route\.' knowledge/runtime/routing-registry.yaml
    validation_scenarios:
      - validation/scenarios/runtime/workflow-detector-deterministic-match-v1.yaml
      - validation/scenarios/runtime/workflow-detector-conflict-resolution-v1.yaml
      - validation/scenarios/runtime/workflow-session-topic-shift-v1.yaml
      - validation/scenarios/runtime/workflow-detector-travel-planning-regression-v1.yaml
    coverage_target_pct: 90    # 期望 90%+ routes 被某個 scenario 覆蓋
    regression_scenarios:
      - validation/scenarios/runtime/workflow-detector-travel-planning-regression-v1.yaml
```

**為什麼 4-value enum 而非 3-value**：
- `behavioral_only` 隱含「應該但暫時沒做」 → review queue 有意義
- `not_mechanizable` 表達「永遠不該做」 → review queue 應排除，避免 noise
- 兩者塞同一桶會讓 review queue 永遠塞著無解項目，治理失靈

產出：
- [ ] schema v2 定稿（4-value enum + 對應 metadata requirements）
- [ ] enforcement-registry.yaml 初版（~24 rule_classes 全部填）
- [ ] companion `enforcement/enforcement-registry.md`（philosophy + Layer 2.5 framing + 寫作指南）

產出：
- [ ] schema 定稿（包含 `behavioral_only` field 與必填的 rationale / sunset_decision）
- [ ] enforcement-registry.yaml 初版（覆蓋 Phase 1 inventory 列出的所有 rule）
- [ ] companion `enforcement/enforcement-registry.md`（philosophy + 寫作指南）

### Phase 3 — Compile-time Lint

**Status (2026-05-31)**: Lint 已實作於 [`scripts/ai-skill-cli/internal/app/enforcement_registry_lint.go`](../../scripts/ai-skill-cli/internal/app/enforcement_registry_lint.go) + unit tests + scenario-shaped fail/pass coverage（全部 PASS）。**尚未** wire 進 `ai-skill runtime compile` —— 等下方 Findings backfill 與 user 確認後再 wire（避免 main branch 被破壞）。

#### Phase 3 Dry-Run Findings (2026-05-31)

Lint dry-run（registry 未 backfill 狀態）surface **22 findings**：17 個 `orphan_rule` + 5 個 `orphan_executor`。零 `behavioral_only_incomplete_sunset`、零 `missing_executor_symbol`、零 `deprecated_*`。

**Auto-classified（bulk 1-case 建議；待 user 確認後 backfill）**：

| # | Finding | 類型 | 建議處置 |
|---|---|---|---|
| F1 | `runtime/cognitive-modes-token-budget.yaml` orphan_rule | bulk | 加進 `cognitive_mode_governance.source_files`（已有 7 個 cognitive-modes-* yaml，少了這一個 token-budget 同源） |
| F2 | `governance/ai-runtime-governance/linked-update-governance.yaml` orphan_rule | bulk | 加進 `linked_updates.source_files`（同主題） |
| F3 | `governance/lifecycle/executable-contract-boundary.yaml` orphan_rule | bulk | 加進 `runtime_yaml_projection.source_files`（同主題：yaml ↔ runtime.db projection 邊界） |
| F4 | `governance/lifecycle/executable-contract-inventory.yaml` orphan_rule | bulk | 加進 `runtime_yaml_projection.source_files`（同主題） |
| F5 | `enforcement/neutral-language.yaml` orphan_rule | bulk | 新增 rule_class `neutral_language` coverage=`not_mechanizable`（與 `tool_neutral_documentation` 同類，writing judgement） |
| F6 | `enforcement/feedback-lessons.yaml` orphan_rule | bulk | 加進 `failure_learning_system.source_files`（feedback lessons 屬該系統） |
| F7 | `validateStopHookFinalTexts` orphan_executor | bulk | 加進 `internal_helper_allowlist`（plural 是 singular `validateStopHookFinalText` 的 collector helper，已綁定於 `dirty_repo_close_loop`） |
| F8 | `runHooks` orphan_executor | bulk | 加進 `internal_helper_allowlist`（是 `ai-skill hooks` CLI 入口 router，不是執行規則的 executor；類比 `runRuntime`/`runDoctor`） |
| F9 | `runPostToolUseHook` orphan_executor | bulk | 加為 `bootstrap_integrity.executors[]`（hook_dispatcher_entry, PostToolUse, warn）——目前已主動 emit Bootstrap Receipt reminder） |
| F10 | `runPreCommitHook` orphan_executor | bulk | 加為 `shell_script_governance.executors[]`（hook_dispatcher_entry, pre-commit；目前該 class 只列被 dispatch 的 validator，少了 dispatcher 自己） |
| F11 | `runUserPromptSubmitHook` orphan_executor | bulk | 加為 `bootstrap_integrity.executors[]`（hook_dispatcher_entry, UserPromptSubmit, warn）—— Phase 5.x ADR 草稿已預示這條 |

**待 user 裁決（11 條 orphan_rule 沒有明顯歸屬）**：

每條都需要 user 回覆「歸屬哪個 rule_class」或「新增 rule_class」或「mark deprecated」。

| # | yaml | 可能歸屬 | 開放問題 |
|---|---|---|---|
| F12 | `enforcement/authorization-scope.yaml` | 新 class `authorization_scope` (mechanical? behavioral_only?) | 是否有 executor 對 sanitization 旁的 authorization 範圍做檢查？若無 → behavioral_only + sunset |
| F13 | `enforcement/content-layering.yaml` | 加進 `document_sizing.source_files` OR 新 class | 是 document_sizing 的姊妹（layering vs. sizing），同 class 或拆？ |
| F14 | `enforcement/cross-skill-references.yaml` | 新 class behavioral_only | sunset trigger 寫什麼？ |
| F15 | `enforcement/decision-efficiency.yaml` | 新 class behavioral_only OR not_mechanizable | 效率判斷無客觀 metric → 傾向 not_mechanizable |
| F16 | `enforcement/document-todo-list.yaml` | 新 class behavioral_only | 是否有可機械化的「TODO 表存在」lint？目前無 → behavioral_only |
| F17 | `enforcement/goal-action-validation.yaml` | 新 class behavioral_only | 與 `conversation_goal_ledger` 是否合併？ |
| F18 | `enforcement/prompt-cache-efficiency.yaml` | 新 class behavioral_only | 同 decision-efficiency 屬於 P3 efficiency；可否合併成單一 `efficiency_governance` class？ |
| F19 | `governance/ai-runtime-governance/validation-scenario-governance.yaml` | 新 class | validation scenario 治理目前無 executor → behavioral_only 還是 pending_implementation（child plan 建立 scenario lint）？ |
| F20 | `governance/lifecycle/decision-promotion-pipeline.yaml` | 新 class behavioral_only | 與 `failure_learning_system` 是否合併？ |
| F21 | `governance/lifecycle/directory-structure-governance.yaml` | 新 class | 目錄結構治理：有 executor 嗎？若無 → behavioral_only |
| F22 | `governance/lifecycle/knowledge-update-flow.yaml` | 加進 `linked_updates.source_files` OR 新 class | 與 linked_updates 同源還是獨立流程？ |

**下一步**：等 user 對 F1-F11 bulk 建議 + F12-F22 個別裁決後，agent backfill registry yaml → re-dry-run 直到 PASS → wire 進 `ai-skill runtime compile` → rebuild 5 platform binaries → commit + push + readback.

#### 原 pseudo-implementation（保留作為設計參考）

在 `scripts/ai-skill-cli/internal/compile/`（或既有 compile pipeline）加 lint pass：

```go
// pseudo:
func LintEnforcementRegistry(reg EnforcementRegistry, repo Repo) []LintError {
    var errs []LintError

    // 1. 每條 binding 的 rule_source 必須存在且包含聲稱的 id
    for _, b := range reg.Bindings {
        if !repo.RuleExists(b.RuleSource) {
            errs = append(errs, LintError{...})
        }
    }

    // 2. 每條 binding 的 executor.symbol 必須在 file 中真實存在
    //    （除非 status: pending + 有 child_plan reference）
    for _, b := range reg.Bindings {
        if b.Status == "pending" { continue }
        if !repo.SymbolExists(b.Executor.File, b.Executor.Symbol) {
            errs = append(errs, LintError{...})
        }
    }

    // 3. 每條 enforcement_layer: behavioral_only 必須有 rationale + sunset_decision
    for _, b := range reg.BehavioralOnlyRules {
        if b.Rationale == "" || b.SunsetDecision == "" {
            errs = append(errs, LintError{...})
        }
    }

    // 4. 掃 repo 找 orphan rules（聲明 id 但 registry 沒 binding）
    for _, ruleID := range repo.AllDeclaredRuleIDs() {
        if !reg.HasBinding(ruleID) {
            errs = append(errs, LintError{
                Type: "orphan_rule",
                Msg:  "rule " + ruleID + " declared but no enforcement-registry binding",
            })
        }
    }

    return errs
}
```

整合：`ai-skill runtime compile` 跑 lint，任何 error 直接 exit non-zero。

產出：
- [ ] Lint 實作 + unit tests（測 4 種 lint type）
- [ ] `ai-skill runtime compile` 整合
- [ ] 第一次 compile run 預期會列出大量 orphan_rule（Phase 1 inventory 已知 ≥ 30%）；不直接 block，先記為 warning + 寫進 Q1 review 清單

### Phase 4 — CLI Coverage Report（**v2: 本 plan 主要交付**）

第七輪評審：「Registry 的最大價值不是 audit 既有，而是強制新規則作者回答 coverage 問題」。Phase 4 在 v2 升格為 primary deliverable。

新增 `ai-skill enforcement coverage`：

```
$ ai-skill enforcement coverage
Enforcement Coverage Report (2026-XX-XX)
═══════════════════════════════════════
Total Rule Classes: 24

  ├─ Mechanical:                11  (46%)
  │   ├─ Fully verified:         6  (symbol+scenario+regression+runtime)
  │   ├─ Verified, low runtime:  2  ⚠ (scenarios pass, but observed << expected)
  │   └─ Symbol only:            3  ⚠ (scenarios missing)
  ├─ Behavioral only:            5  (21%)
  ├─ Not mechanizable:           3  (12%)
  ├─ Pending implementation:     3  (12%)  — child plans active
  ├─ Research required:          1  ( 4%)  ⚠ (no clear path yet)
  └─ Deprecated:                 1  ( 4%)

Per-class detail:

  Rule Class                Status                Verification           Runtime (30d)
  ───────────────────────────────────────────────────────────────────────────────────
  bootstrap_integrity       mechanical            ✓ full                 ✓ 100% sessions
  commit_governance         mechanical            ✓ full                 ✓ 19/19 validators fired
  cognitive_mode_governance mechanical            ✓ full                 ✓
  routing                   mechanical            ⚠ symbol-only          ⚠ no metrics yet
  workflow_activation       pending_implementation —                     planned: 57 routes
                            (child: plans/.../workflow-activation-engine.md)
  sanitization              pending_implementation —                     planned: 11 patterns
                            (child: plans/.../sanitization-validator.md)
  intelligence_classification research_required   —                     —
                            (research_questions: how to disambiguate primary vs secondary
                             from a single route declaration; estimated_unblock: 2026-Q3)
  capability_discovery      behavioral_only       review queue            n/a
  rule_weight               behavioral_only       review queue            n/a
  tool_neutral_doc          not_mechanizable      n/a                    n/a
  rule_writing_quality      not_mechanizable      n/a                    n/a
  old_intelligence_route    deprecated            n/a                    replaced_by: intelligence_classification

⚠ Runtime-observed gaps (v4):
  workflow_activation       declared 57 routes / scenario 95% / runtime_observed 37%  ← potential dead-route
                            never-observed: route.analysis.security, route.intelligence.engineering.heuristics,
                            ... (36 routes 0 hits in 30d — review: deprecate or detector bug?)

Pending implementation (active child plans):
  workflow_activation   → plans/active/2026-05-31-1900-workflow-activation-engine.md (P2)
  sanitization          → plans/active/2026-05-31-2000-mechanical-sanitization-validator.md (P3)

Research required (no clear mechanization path):
  intelligence_classification  — research_questions: how to disambiguate primary vs
                                 secondary at registry declaration time;
                                 estimated_unblock: 2026-Q3

Behavioral_only awaiting sunset review:
  capability_discovery  — revisit when: workflow_activation Phase 6.1 land
  rule_weight           — revisit when: 3+ detectable P0 patterns surface

Not_mechanizable (closed, will not appear in review queue):
  tool_neutral_documentation  — subjective writing judgment
  rule_writing_quality        — would game readability metrics

Deprecated (awaiting removal):
  old_intelligence_route  — replaced_by: intelligence_classification (2026-06-30)
```

**為什麼 runtime_observed 那麼重要**：上方範例的 `workflow_activation` 行直接揭示 v3 看不到的真相 —— **57 routes declared、scenarios 95% 覆蓋，但 runtime 只觀測到 37%（21/57）**。意思是 36 個 route 半年從沒被觸發。兩種可能：
1. 真的沒人用 → 候選 deprecate
2. Detector 對這些 route 安靜壞了 → 緊急修

v3 只看 symbol + scenario 完全看不到這層 reality gap。

**強制觸發場景**（v3：含 verification 維度）：

| 場景 | Coverage report 行為 |
|---|---|
| 新增 `enforcement/<new-rule>.yaml` 但未在 registry 出現 | compile fail：`new rule class not registered` |
| Registry 加新 entry 但漏寫 `coverage` field | compile fail：`missing coverage field` |
| `coverage: behavioral_only` 但 `sunset_decision.success_criteria` 空白 | compile fail：`behavioral_only requires success_criteria` |
| `coverage: not_mechanizable` 但 `objective_validation_impossible_because` 空白 | compile fail：`not_mechanizable requires impossibility rationale` |
| **v3**：`coverage: deprecated` 但 `replaced_by` / `removal_date` 都缺 | compile fail：`deprecated requires replaced_by or removal_date` |
| **v3**：`coverage: deprecated` 且 `removal_date` 已過但 rule 仍存在 | compile fail：`deprecated rule past removal_date — actually remove or extend date` |
| 既有 mechanical class 的 executor symbol 在 hooks.go 找不到 | compile fail：`executor symbol missing` |
| **v3**：`coverage: mechanical` 但 `coverage_evidence.validation_scenarios` 為空 | compile warning：`mechanical without validation scenarios — verification level: symbol_only` |
| **v3**：mechanical 且 instance set 已知（如 workflow 57 routes）但 scenarios coverage < 50% | compile fail：`mechanical coverage_evidence under threshold` |
| **v3**：mechanical 且 scenarios coverage < 80% | compile warning：`mechanical coverage_evidence below target` |
| 新增 mechanical executor 但無對應 rule_class | compile warning：`orphan executor` |

這意味未來任何新規則寫作者，在 compile 那一刻就被強制回答：
> 「這條規則是 mechanical / behavioral_only / not_mechanizable / pending？」

無法迴避、無法 silent leak、無法等使用者半年後追問。

產出：
- [ ] CLI subcommand 實作
- [ ] 輸出格式（text + JSON + markdown for governance dashboards）
- [ ] 文件化（README + ai-tools/agent reference）
- [ ] CI integration：Pull Request 自動跑 coverage diff，新增規則沒填 coverage 直接 PR check 失敗

### Phase 4.5 — Registry Self-Governance（**v4 NEW，採納評審 #3**）

Layer 2.5 自己也需要治理。沒有這層，registry 變成「一個沒人管的元數據檔」。

#### Status Transition Matrix

| From → To | Required action | Lint enforcement |
|---|---|---|
| (new) → `pending_implementation` | 引用 active child plan | child_plan path 必須 exist |
| (new) → `research_required` | 列 `research_questions` ≥ 1 + estimated_unblock | metadata schema 校驗 |
| `pending_implementation` → `mechanical` | child plan Phase 3+ 完成 + executor symbol live + `coverage_evidence` 填齊 | `verification_levels` 全部達門檻 |
| `research_required` → `pending_implementation` | 提出 child plan 解決所有 research_questions | research_questions 全 resolved + child_plan exists |
| `mechanical` → `behavioral_only` | **demotion，需 ADR**：列 `demotion_rationale` + `adr_reference` | adr_reference 必須指向 `constitution/ADR-*.md` |
| `mechanical` → `deprecated` | replaced_by 必須指向 active mechanical class | lint：replaced_by 解析成功 |
| `behavioral_only` sunset 屆期 | `sunset_decision.revisit_when` 觸發 → governance review queue | runtime 比對日期 / event |
| `deprecated` `removal_date` 屆期 | governance 決定 actually remove vs extend_with_rationale | lint：過期不 fail，但 coverage report 標 red |

#### 誰可以改？

- **Add new rule_class**：標準 PR review，至少一個 maintainer approval
- **Promotion**（往更 mechanical 方向）：lint pass 自動允許，無需特別 approval（因為 lint 已強制 verification 證據）
- **Demotion**（從 mechanical 退到較弱狀態）：**需 ADR 引用**。意圖是讓「我們決定不機械化某條規則」這種決策有 governance 軌跡，不能 silent 改一行 yaml 就降級
- **Status field 自身變更**：commit-msg hook 加 validator `validateEnforcementRegistryTransition`，比對 git diff 看 status 欄位是否變動、檢查 transition 合法性與是否附 ADR

#### Self-Governance Lint Rules

- **R1**：rule_class status 變更 commit 必須包含 `[registry-status-change]` trailer 與 `rationale:` body
- **R2**：demotion 必須附 ADR；無 ADR → commit reject
- **R3**：promotion 必須對應 verification_levels 達 mechanical 門檻；未達 → lint fail
- **R4**：deprecated 過 removal_date 30 天仍未實際移除 → governance alert（不阻塞 build，但 report 標 red）
- **R5**：research_required 過 estimated_unblock_timeline 仍未轉 pending_implementation 或 demote → governance review trigger

#### 為什麼 self-governance 必須是一階公民

沒有 self-governance：
- 任何 dev 可改 yaml 一行把 mechanical 改 behavioral_only，silent demotion
- deprecated 永遠 deprecated，沒人記得真的 remove
- behavioral_only 的 sunset 沒有實際觸發機制

有 self-governance：
- Demotion 留下 ADR 軌跡
- Deprecated 屆期顯示在 coverage report
- Sunset triggers 自動進 governance queue

這把 Layer 2.5 從「另一個 yaml 檔」升為「真實治理工具」。

產出：
- [ ] Phase 4.5 lint rules R1-R5 實作
- [ ] commit-msg `validateEnforcementRegistryTransition` validator
- [ ] coverage report 章節「Governance Alerts」整合 sunset + removal_date breach + research_required timeout
- [ ] `enforcement/enforcement-registry.md` companion 加 §Self-Governance

### Phase 5 — Bootstrap Integration

把 enforcement-registry 接進 bootstrap：

- [ ] `runtime/core-bootstrap.yaml` 加 contextual_activation：「修改 enforcement rule 時必讀 enforcement-registry.md」
- [ ] `ai-skill runtime receipt` 輸出加 enforcement coverage 摘要（一行）：
  `Enforcement: bound=38/152 (25%) behavioral=7 orphan=107`
- [ ] commit-msg hook 加 validator：若 commit 變更 enforcement rule yaml/md，必須同步更新 enforcement-registry.yaml；未同步 → reject

#### Phase 5.x — Hook Injection Economics（**inaugural self-governance test case**）

**Empirical trigger**：2026-05-31 session 揭露 `runUserPromptSubmitHook` 過量注入問題 —— UserPromptSubmit hook 每輪都注入完整 `CORE_BOOTSTRAP.md` + 條件式 Receipt reminder，造成 agent over-emit Bootstrap Receipt（本應第一輪 only）。

**為什麼放在 Phase 5 而不是另開 plan**：這個 fix 不只是「修 hook bug」，它**剛好同時是 registry 機制的活體測試**：

1. 修改 `runUserPromptSubmitHook` 注入策略 = 改變 `bootstrap_integrity` rule_class 的 enforcement 行為
2. 「always inject full bootstrap.md」→「transcript-aware conditional inject」是一種 enforcement 強度 demotion（從強制→條件式）
3. v4 Self-Governance lint R2「demotion 必須附 ADR」應該在這個 commit 上被觸發
4. v4 verification axis `runtime_observed` 應該在 fix land 後監測「Receipt 重複出現率」是否真的下降
5. Coverage Report 應該在 fix 前後顯示 `bootstrap_integrity` runtime_metrics 變化

換言之，這個 fix 是 registry 機制最完整的 dogfooding case —— 同時驗 self-governance lint、ADR 強制、verification axis、runtime_metrics 收集、Coverage Report 變化偵測五個 v4 機制。

**4 個 defects 明細**（commit 時必須在 ADR 中對應處理）：

| # | Defect | hooks.go 位置 | Fix |
|---|---|---|---|
| D1 | 每輪注入完整 CORE_BOOTSTRAP.md（數 K tokens × N turn） | `runUserPromptSubmitHook` L1015-1021 `combined := ... + bootstrap` | 條件式：transcript 已 ack 時跳過 bootstrap.md |
| D2 | "If bootstrap has not yet been acknowledged" 是 client-side 不可驗證 conditional | L1018 hardcoded string | 改為 hook 自己判定，agent 不再做 conditional 解讀 |
| D3 | MUST（cognitive mode）+ conditional（Receipt）黏在同一行 | L1017-1018 string concat | 拆兩個獨立 block，MUST 永遠注入、Receipt 按需 |
| D4 | Hook 無 transcript scan | 整個 `runUserPromptSubmitHook` 不讀 transcript | 加 `transcriptHasBootstrapAcknowledgment(transcript, lastN=20)` helper，模式複用既有 `transcriptHasRequiredBootstrapReads` |

**ADR 草稿**（fix commit 必須引用）：

```
constitution/ADR-XXX-conditional-bootstrap-injection.md

Title: UserPromptSubmit hook 改為 transcript-aware conditional injection
Status: proposed
Context:
  - UserPromptSubmit hook 每輪注入完整 CORE_BOOTSTRAP.md，造成 token 浪費
  - "If bootstrap has not yet been acknowledged" conditional 由 agent 解讀，
    導致 over-emit Receipt 模式（本 session 自身發現）
Decision:
  - Hook 自己掃 transcript 判定 acknowledgment，只在缺失時注入完整 bootstrap.md
  - MUST/conditional 拆為兩個獨立 block
  - Bootstrap_integrity rule_class enforcement 從 "always full-context injection"
    降級為 "conditional minimal injection + PreToolUse mechanical gate"
Consequences:
  - Token saving (~2-3K × N turn)
  - PreToolUse `gate.bootstrap.receipt_present` 不變（仍是真正 enforcement）
  - registry.yaml: bootstrap_integrity 的 enforcement_strength 標記更新
```

**Phase 5.x 步驟**（強制按順序執行）：

1. **先寫 ADR** 至 `constitution/ADR-XXX-conditional-bootstrap-injection.md`
2. **更新 registry**：`bootstrap_integrity` entry 加 `enforcement_changed_at` + `adr_reference: ADR-XXX`
3. **改 hooks.go**：D1-D4 patch + `transcriptHasBootstrapAcknowledgment` helper
4. **加 unit test**：兩個 case
   - first turn: transcript 無 Receipt → 注入完整 bootstrap.md
   - subsequent turn: transcript 有 Receipt → 只注入 Cognitive Mode reminder
5. **Rebuild binaries**（5 platform）+ commit code + binary 同 commit
6. **新 session 驗證**：開新 session，user message 內看 hook 注入內容應該短了
7. **30 天後檢查** `bootstrap_integrity.runtime_metrics`：Receipt 過量發送率應下降；若無數據（hook fix 也包含 runtime_metrics 收集）→ verification 維度補登

**為什麼這是最好的 inaugural test**：
- 真實 production hook（不是 toy example）
- 同時觸發 v4 五個機制（self-gov / ADR / verification / runtime_metrics / Coverage Report 變化）
- Fix 失敗 = registry 機制本身不夠成熟，有 actionable feedback
- 成功 = registry 第一個被「真實寫進 codebase 的 enforcement 變更」trigger，dogfooding 完成

**依賴**：Phase 5.x 不能比 Phase 2 (schema) + Phase 4.5 (self-governance lint) 早做。可與 Phase 5 主流程並行，不阻擋。

### Phase 6 — Failure Pattern + Cross-link to Child Plans

- [ ] 新建 `enforcement/failure-patterns/rule-without-executor.md` —— 把 2026-05-31 session 暴露的 5 個 instance 集中記錄，作為 meta-pattern 的 inaugural reference
- [ ] 兩個 child plan 各加 link 指回本 plan，明示「本 plan 是 meta-pattern 的 instance」
- [ ] `enforcement/README.md` 加章節「Mechanical Enforcement Registry」指向本 plan + registry yaml

### Phase 7 — Validation Scenarios

- [ ] `validation/scenarios/enforcement/registry-lint-orphan-rule-v1.yaml`
- [ ] `validation/scenarios/enforcement/registry-lint-missing-executor-v1.yaml`
- [ ] `validation/scenarios/enforcement/registry-lint-behavioral-without-rationale-v1.yaml`
- [ ] `validation/scenarios/enforcement/coverage-cli-output-format-v1.yaml`
- [ ] `validation/scenarios/enforcement/2026-05-31-regression-five-instances-v1.yaml`（回放 session 揭露的 5 個 gap，registry 必須 detect 全部）

### Phase 8 — Close-out

- [ ] phases done
- [ ] `git status` clean
- [ ] `git push` 完成、`git log origin/main..HEAD` empty
- [ ] 讀回更新後的 enforcement-registry.yaml / enforcement-registry.md / failure pattern
- [ ] Archive 本 plan + 確認 child plans 的 cross-link 正確

---

## Open Questions

| # | Question | 處置 |
|---|---|---|
| Q1 | Phase 1.3 預期 orphan rate 70%；第一次 lint 直接 block compile 還是先 warn 一個 grace period？ | **resolved (2026-05-31 session)** → **Hard block，無 grace period**。理由：grace period 會回到「warning → 先放著 → 半年後還在」的失效模式，違背 Prevent > Detect > Repair 哲學。`orphan_rule` 與 `orphan_executor` 都 hard fail；第一次 land 預期需要密集 backfill，但這是 one-time cost。Schema：`enforcement_mode: { orphan_rule: fail, orphan_executor: fail }`，`bootstrap_grace.enabled: false` |
| Q2 | `behavioral_only` 的 `sunset_decision` 強制格式？開放 free text 容易變空話 | **resolved (2026-05-31 session, revised)** → **`revisit_when` + `success_criteria` 雙必填**；`revisit_owner` recommended。理由：原本只選 success_criteria 會落入「有標準但永遠沒人檢查」的失效模式 —— 比「沒標準但會被檢查」更危險，因為長出虛假安全感。`revisit_when` 是「事件 trigger」，`success_criteria` 是「客觀判定」，雙鎖才能形成治理閉環。Compile lint 校驗兩個欄位都存在且非空 |
| Q3 | enforcement-registry 與 `runtime/core-bootstrap.yaml` per_*_obligations 的關係？兩者都列 obligation id | resolved → core-bootstrap.yaml 是 phase-aware obligation lifecycle，enforcement-registry 是 cross-phase binding 視圖。兩者互補：bootstrap 講「何時 fire」，registry 講「何處 enforce」 |
| Q4 | Orphan executor（code 有 validator 但 registry 沒 entry）該強制 binding 還是允許 internal helper？ | **resolved (2026-05-31 session)** → **強制 binding，但限「exported / dispatcher-registered」**。Schema 加 `executor_kind` enum: `[hook_dispatcher_entry, commit_msg_validator, runtime_state_machine_phase, internal_helper]`；`binding_required_for` 白名單只包含前三種。`internal_helper`（parseYaml / normalizePath 等 utility）在 registry 維護顯式 allowlist，避免 code-level annotation 散落 |
| Q5 | Discovery 在 future plan 是否視為「fallback rule」並進 registry？ | resolved → 是。`capability_discovery_fallback` 已列在 Phase 2 schema 範例為 behavioral_only，等 workflow-activation-engine Phase 6.1 實作後改 mechanical |
| Q6 | 已有的 11 個 commit-msg validators 是否全部要在 registry 補 binding 才算 audit 完成？ | **resolved (2026-05-31 session)** → **Phase 1.3 全部 audit**；rule_class 數量採 `soft_target=24 / hard_limit=40` 雙閾值。理由：若只 audit 子集，Coverage Report 會顯示「mechanical 80%」卻其實只覆蓋一半，是假象。盤點實際數量為 28（落在 soft-hard 中間），不需硬塞合併。Schema：`governance_thresholds: { rule_class_soft_target: 24, rule_class_hard_limit: 40 }` |
| 新 Q | rule_class 數量上限：若實際盤點遠超 24 是否該重新評估抽象粒度？ | **resolved (2026-05-31 session)** → 採雙閾值。`exceed_soft_target` → review_split_opportunities（提醒檢視是否該合併或細分）；`exceed_hard_limit` → governance_review_required（停下來重新評估粒度抽象）。實際 Phase 1.1 盤點為 28，在 soft-hard 區間內 |
| Q7 | 2026-05-31 session 揭露的 hook injection economics 問題（runUserPromptSubmitHook 過量注入），是否要另開 plan 處理？ | **resolved (v4+)** — 不另開 plan，**併入 Phase 5.x 作為 inaugural self-governance test case**。理由：fix 本身同時 dogfooding 五個 v4 機制（self-gov lint / ADR demotion / verification axis / runtime_metrics / Coverage Report 變化），是 registry 機制最完整的活體測試。詳見 Phase 5.x。 |

---

## Validation Plan

- [ ] Phase 1 inventory 數量符合預估（≥ 150 rules、≥ 30 executors）
- [ ] Phase 2 schema 經 user review（特別是 `behavioral_only` 寫作格式）
- [ ] Phase 3 lint 跑出第一次 coverage report，列出 ≥ 100 orphan rule
- [ ] Phase 4 CLI output 易讀，能直接拿去做 governance review
- [ ] Phase 5 bootstrap integration 不破壞既有 receipt format（backward compat）
- [ ] Phase 7 regression scenario 五個全 PASS

---

## Dependency Read Ledger

| 欄位 | 內容 |
|---|---|
| Trigger | 2026-05-31 session round-2 評審指出 sanitization plan v1 缺 meta-pattern 抽象，建議「Governance Rule Coverage Audit」 |
| Required set | `enforcement/sanitization.md`、`enforcement/rule-weight.md`、`enforcement/dependency-reading.md`、`runtime/core-bootstrap.yaml`、`knowledge/runtime/routing-registry.yaml`（structure ref）、`scripts/ai-skill-cli/internal/app/hooks.go`（executor enumeration ref）、`governance/lifecycle/capability-discovery-philosophy.md`、sibling plan 1900 + 2000 |
| Read | 以上 |
| Not applicable | 無 |
| Deferred | Implementation 細節 source（compile pipeline 完整 code、Go AST 解析）—— Phase 0 unlock 後補讀 |
| Validation | Architecture Compatibility Preflight 已列；Phase 0 unlock 前驗證 |

---

## v4 改動摘要（Round 9 評審整合）

Round 9 評分：**8.8 / 10**。三個扣分點各對應一條 v4 修正：

| # | 評審論點 | 採納 | 對應修改 |
|---|---|---|---|
| 1 | 缺 `runtime_observed` 第 4 層 verification：scenario_exists ≠ production reality_exists | ✅ | verification_levels 從 3 升 4；新增 `runtime_metrics_spec`（observation_window_days / activation_count / instance_breakdown）；Coverage Report 加 Runtime 欄；worked example 顯示 workflow_activation 57 declared / 95% scenario / 37% runtime 的 dead-route 警告 |
| 2 | `pending` 混兩種成熟度（knows how vs needs research） | ✅ | coverage enum 從 5 升 6：`pending` 拆 `pending_implementation`（要 child_plan）+ `research_required`（要 research_questions + estimated_unblock_timeline）。Coverage Report bucket 也分開 |
| 3 | Registry 自身治理未定義（誰審核 coverage status 變更） | ✅ | Phase 4.5 新增 Registry Self-Governance：Status Transition Matrix、Demotion 必須 ADR、5 條 self-governance lint rules、commit-msg `validateEnforcementRegistryTransition` validator |

**Round 9 核心觀察（user 原話）**：
> Scenario Exists 不代表 Production Reality Exists。
>
> v3 的 Pending 其實混兩種：知道怎麼做只是還沒寫 vs 知道應該機械化但還不知道怎麼做。
>
> 任何 Rule Class 都必須顯式宣告自己的 Enforcement Strategy 與 Verification Strategy。

**v3 → v4 的核心進化**：
```
v3:  Rule  ←binding→  Executor  ←evidence→  Verification (scenarios)
v4:  Rule  ←binding→  Executor  ←evidence→  Verification (scenarios + runtime)
                                              ↑
                                     + Registry Self-Governance
                                       管理 binding 變更本身
```

從「驗證 executor 對應規則」進一步到「驗證 executor 在 production 真的有跑」+「驗證 registry 本身的演化是治理可追蹤的」。Framework Self-Audit Layer 至此真正完整。

## v3 改動摘要（Round 8 評審整合）

| # | 評審論點 | 採納 | 對應修改 |
|---|---|---|---|
| 1 | 缺 `deprecated` 5th coverage status | ✅ | enum 從 4-value 升 5-value；新增 metadata requirements（`replaced_by` 或 `removal_date`）；compile lint 校驗 removal_date 未過期 |
| 2 | `executor exists ≠ executor covers all instances`：lint 只驗 symbol 存在會錯放（DetectWorkflows 漏 intelligence routes 仍 PASS） | ✅ | 新增 `verification` 維度與 coverage 正交：symbol_exists / scenario_exists / regression_exists；新增 `coverage_evidence` schema 含 expected_instance_count + validation_scenarios + coverage_target_pct |
| 3 | Coverage Report 應顯示 declared vs covered 百分比 | ✅ | Coverage Report 新增 "Verification" + "Instance Coverage" 欄位；example 顯示 mechanical-but-symbol-only 的警告（routing class 範例） |
| 4 | 定位確認為 "Meta Governance / Framework Self-Audit Layer"（Prevent > Detect > Repair） | ✅ | Decision Rationale 章節將 Layer 2.5 升格為「Meta Governance」；Header `世代` 行同步更新 |
| 5 | 自然演進方向：Rule ↔ Executor ↔ Verification Registry（三元 binding） | ✅ | v3 schema 已包含 verification 維度 = 三元 binding。Future v4 可能擴展 verification_levels 含 runtime metrics（執行頻率、實際阻擋次數） |

**Round 8 核心觀察（user 原話）**：
> 從「找缺少 executor 的規則」變成「要求每個 Rule Class 必須宣告自己的 enforcement strategy」—— 看起來像，但是不同層級。
>
> v2 的 mechanical 只驗 `Executor Exists`，v3 進化為 `Executor Verified`。
>
> 這已經開始形成真正的 Framework Self-Audit Layer。

## v2 改動摘要（Round 7 評審整合）

| # | 評審論點 | 採納 | 對應修改 |
|---|---|---|---|
| 1 | 既有架構缺 Layer 2.5 Coverage Verification | ✅ | Decision Rationale 新增 "Architectural Framing — Missing Layer 2.5" 章節 |
| 2 | Rule instance 級 binding 會讓 registry 自我膨脹（150+ entries），改 Rule Class 級 | ✅ | Phase 1 全部重寫：從「audit 150+ rule」改「識別 ~24 rule_class」。Phase 2 schema 用 `rule_classes` 而非 `bindings`。`commit_governance` 一條 entry 涵蓋 19 個 commit-msg validator |
| 3 | 新增 `not_mechanizable` coverage status，與 `behavioral_only` 區隔 | ✅ | coverage enum 從 3-value 升 4-value。各狀態有不同 metadata requirements：behavioral_only 要 `sunset_decision`，not_mechanizable 要 `objective_validation_impossible_because` |
| 4 | 最大價值是 Coverage Report，不是 Registry 本身 | ✅ | Phase 4 升為 primary deliverable，加 6 個 compile-fail 條件確保新規則作者必須回答 coverage 問題 |
| 5 | Priority 應 reorder，本 plan 升為 P1 | ✅ | Header 加 `Priority: P1`，並說明 child plans 完成本 plan 後重啟時會被 coverage lint 強制觸發完整 binding 宣告 |

**Round 7 核心觀察（user 原話）**：
> Registry 真正的價值在 Coverage Report —— 強制新規則作者回答「mechanical / behavioral_only / not_mechanizable / pending」。這才是把 "Rule Exists, Executor Missing" 從個案修補提升成框架級不變量。

## Source

2026-05-31 session：使用者連續 7 輪追問 / 評審，依序暴露：
1. sqlite3 vs ai-skill CLI 認知偏差
2. workflow activation gap
3. Discovery vs Detector 混淆
4. intelligence 預設 advisory 風險
5. sanitization gate 自我觸發失敗
6. **本 plan 對應**：以上 5 條共同根因是「Knowledge layer 有規則，Runtime layer 沒執行器」meta-pattern，建議建立 governance-level coverage audit
7. **Round 7 v2 重構**：v1 寫成 rule instance 級 audit 是維護地獄；改 rule_class 級 + 4-value coverage enum + Coverage Report 為主要交付 + Priority 提升 P1
8. **Round 8 v3 verification 升級**：v2 的 mechanical 只驗 symbol 存在，可能漏覆蓋；v3 加 `verification` 維度（正交）+ `coverage_evidence` schema + `deprecated` 5th status + 定位升格為 Meta Governance
9. **Round 9 v4 reality + self-governance**：v3 verification 仍只到 scenario 級別；v4 加 `runtime_observed` 4th verification level + 拆 pending 為 implementation/research + 新增 Registry Self-Governance（Phase 4.5）。評分 8.8/10。

User 原話（round 6）：
> Ai-skill 現在最大的風險不是缺規則，而是「Rule Exists, Executor Missing」。
> 如果有這種 Registry，下一個 gap 會在規則新增時就被抓出來，而不是半年後被使用者追問才發現。

本 plan 是 session 累積評審的 meta-fix —— 兩個 child plan（activation-engine、sanitization-validator）解決 2 個 instance，本 plan 建立 framework property 預防所有未來 instance。

## Companion References

- `enforcement/sanitization.md` —— canonical rule（companion markdown），作為 Phase 1 inventory 樣本之一
- `enforcement/dependency-reading.md` —— writeback transaction 描述（Phase 5 bootstrap integration ref）
- `enforcement/rule-weight.md` —— rule weight priority（registry lint 排序依據）
- `runtime/core-bootstrap.yaml` —— obligation lifecycle source（與 registry binding view 互補）
- Child plan：`plans/active/2026-05-31-1900-workflow-activation-engine.md`
- Child plan：`plans/active/2026-05-31-2000-mechanical-sanitization-validator.md`
