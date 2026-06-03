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

#### Phase 3 Round-2 Review — Schema 收緊 + 分類修正（2026-06-01）

User round-2 評審指出 F12-F22 直接 backfill 會 surface 6 個結構性問題；在 backfill 前先 land schema 修正，否則新增的 behavioral_only entries 會成為新的灰色地帶（同 registry 想預防的 failure pattern）。

**問題 → 處置 矩陣**（user priority ranking 已採納）：

| Pri | # | 問題 | 處置 |
|---|---|---|---|
| P0 | R1 | F12-F22 schema 已要求 `revisit_owner`（registry.yaml §coverage_status_spec.behavioral_only 把它列 "recommended"，但 user round-2 要求 **upgrade 為 strictly required**） | (a) 修改 `coverage_status_spec.behavioral_only.requires` 加 `sunset_decision.revisit_owner`；(b) 既有 7 個 behavioral_only entry 全部補 `revisit_owner` 欄位；(c) lint `lintBehavioralIncompleteSunset` 加第 3 個必填欄位檢查 |
| P0 | R2 | Q1 hard-block + Phase 1.3 全量 audit 組合在「成熟 repo 第一次導入」會 deadlock | (a) registry.yaml 加 `bootstrap_mode` 欄位（值：`baseline_snapshot_v1` / `strict`）；(b) `baseline_snapshot_v1` 模式：第一次 compile 時把當下所有 orphan 記入 `baseline_snapshot.{orphan_rules,orphan_executors}`，lint 對 baseline 內 entry 降為 warning；(c) 新增 orphan（不在 baseline 內）仍 hard fail；(d) baseline 必須附 `baseline_created_at` + `baseline_burndown_target_date`；(e) 過 burndown date 仍有 baseline entry → governance review trigger（不直接 fail）；(f) 本 Ai-skill repo 目前 22 finding 在「直接 backfill 一次 land」可承受範圍內，**baseline_snapshot 為 future-proofing 機制，不必本 session 用**，但 schema + lint 必須先 land |
| P1 | R3 | F19 應該是 `pending_implementation` 而非 `behavioral_only`（registry v3 已含 `coverage_evidence` schema，半機械化中） | (a) F19 改 `coverage=pending_implementation`；(b) 需要 active child_plan 路徑；**user 裁決**：A) 開 stub plan `plans/active/2026-06-01-XXXX-validation-scenario-governance-executor.md`（Phase 0/1 outline only），B) 暫接受 `child_plan: TBD` 但搭配 `pending_grace_until: <ISO date>`（schema 加新欄位允許 30 天內補齊 child_plan），C) 維持 behavioral_only 並在 sunset_decision 內明寫「promote to pending when Phase 4 coverage_evidence machinery is live」。**建議 C**（最小 schema 變動 + 誠實標明 deferred status） |
| P1 | R4 | `behavioral_only_forever`：surface-legal 的 success_criteria（「100% automation」「≥100 incidents」）lint 抓不到 | (a) registry.yaml `behavioral_only` 新增 `sunset_decision.last_reviewed_at` 必填欄位（ISO-8601）；(b) lint `behavioral_only_review_age`：>12 月未 revisit → warning；>24 月 → fail；(c) lint 算 age = now - last_reviewed_at；(d) entry 首次建立時 `last_reviewed_at = created_at`，每次手動 review 後手動更新 |
| P2 | R5 | F20 `decision_promotion_pipeline` 與 `failure_learning_system` 都是 promotion 治理，跨界 entry（incident → pattern → decision → constitution）權責不清 | (a) registry.yaml `rule_classes[]` 新增 optional `upstream_classes: []` 欄位，宣告本 class 接收哪些上游 class 的 promotion artifact；(b) F20 `decision_promotion_pipeline.upstream_classes = [failure_learning_system]`；(c) lint `upstream_chain_resolution`：引用 class 必須存在；(d) coverage report future（Phase 4）可視覺化 promotion chain |
| P2 | R6 | F22 加進 `linked_updates` 可能讓該 class 變超級桶（`linked_updates + knowledge_update_flow + markdown_yaml_sync + cli_doc_sync + runtime_yaml_projection` 全塞同 class） | (a) registry.yaml `governance_thresholds` 新增 `max_source_files_per_class: 5` + `review_when_source_files_gt: 5`；(b) lint `class_size_review_threshold`：source_files 超過 5 條 → warning（不 fail，提醒 maintainer 評估拆分）；(c) **F22 改回新獨立 class `knowledge_update_flow`**（避免立刻觸發 size warning），companion notes 在兩 class 互相 cross-link |

**F12-F22 修正稿（最終版，全部補齊 rationale + revisit_owner + 提高門檻 + 套 R3/R5/R6 修正）**：

| # | yaml | 最終處置 |
|---|---|---|
| F1-F11 | (見上方表) | bulk 確認，不變動 |
| F12 | authorization-scope | 新 class `authorization_scope`, **behavioral_only**. **rationale**: "authorization 範圍涉及主觀的『誰能改什麼』判斷，無單一機械 predicate". revisit_when: "≥2 失效模式累積（authorization escape incident），且至少 1 條可機械偵測". success_criteria: "authorization-detector executor 上線並覆蓋 ≥80% 歷史 incident". revisit_owner: "framework maintainer". last_reviewed_at: "2026-06-01" |
| F13 | content-layering | 加進 `document_sizing.source_files` |
| F14 | cross-skill-references | 新 class `cross_skill_references`, **behavioral_only**. **rationale**: "skill ↔ skill 引用涉及 promotion 路徑與 reusability 判斷". revisit_when: "≥3 link-rot incidents（broken cross-link / stale reference / circular promotion）累積". success_criteria: "link-resolver + promotion-target lint 上線並覆蓋三類 incident". revisit_owner: "framework maintainer". last_reviewed_at: "2026-06-01" |
| F15 | decision-efficiency | 新 class `decision_efficiency`, **not_mechanizable**. objective_validation_impossible_because: "效率判斷脈絡相依，無 absolute metric，強制機械化會獎勵 gaming" |
| F16 | document-todo-list | 新 class `document_todo_list`, **behavioral_only**. **rationale**: "文件內 TODO 收斂是寫作判斷，無 boolean predicate". revisit_when: "≥3 documented TODO leak incidents". success_criteria: "TODO-lifecycle lint 上線並覆蓋 ≥3 incident pattern". revisit_owner: "framework maintainer". last_reviewed_at: "2026-06-01" |
| F17 | goal-action-validation | 加進 `conversation_goal_ledger.source_files`（撤回原獨立 class 提案，避免 behavioral_only chain） |
| F18 | prompt-cache-efficiency | 新獨立 class `prompt_cache_efficiency`, **not_mechanizable**（拆開不與 F15 合併，避免未來 token-lint 出現時被 class 邊界擋住）. objective_validation_impossible_because: "cache 命中受多輪 context 演化影響，無 absolute metric" |
| F19 | validation-scenario-governance | **R3 處置**：採建議 C，**behavioral_only** + sunset_decision 明寫 "promote to pending_implementation when Phase 4 coverage_evidence machinery is live, then to mechanical when scenario-coverage lint executor 上線". rationale: "目前 scenario 寫作品質與覆蓋判斷無 executor；coverage_evidence schema 已預備但 enforce 邏輯未實作". revisit_when: "Phase 4 CLI `ai-skill enforcement coverage` land". success_criteria: "coverage_evidence.coverage_target_pct 在 compile time 被 enforce". revisit_owner: "framework maintainer". last_reviewed_at: "2026-06-01" |
| F20 | decision-promotion-pipeline | 新獨立 class `decision_promotion_pipeline`, **behavioral_only**, `upstream_classes: [failure_learning_system]`（R5 處置）. rationale: "ADR/governance 升級與 failure pattern 升級是兩類 promotion，目前皆 behavioral；分開以便將來各自有獨立 executor". revisit_when: "≥2 ADR promotion 失效模式累積 OR constitution 引用追蹤需求". success_criteria: "promotion-chain lint 上線並驗證 upstream_classes 鏈完整". revisit_owner: "governance maintainer". last_reviewed_at: "2026-06-01" |
| F21 | directory-structure-governance | 新 class, **behavioral_only**. rationale: "目錄結構治理涉及責任分層判斷". revisit_when: "≥3 documented directory drift incidents". success_criteria: "directory-drift lint 上線並覆蓋三類 incident". revisit_owner: "framework maintainer". last_reviewed_at: "2026-06-01" |
| F22 | knowledge-update-flow | **R6 處置**：新獨立 class `knowledge_update_flow`（不合併進 linked_updates），cross-link to linked_updates in companion notes. **behavioral_only**. rationale: "Knowledge update flow 是 linked-updates 的 step-by-step expansion，但因涉及 cross-layer (enforcement / governance / workflow) 治理，獨立成 class 避免 linked_updates 變超級桶". revisit_when: "Phase 4 後 linked_updates + knowledge_update_flow 共用 executor 是否可行". success_criteria: "兩 class 共用 single executor 而不需個別 dispatcher". revisit_owner: "framework maintainer". last_reviewed_at: "2026-06-01" |

**Schema 變更總清單（在 backfill 之前先 land）**：

1. `coverage_status_spec.behavioral_only.requires` 加 `sunset_decision.revisit_owner`、`sunset_decision.last_reviewed_at`
2. `rule_classes[]` schema 加 optional `upstream_classes: []`（R5）
3. registry 加 top-level `bootstrap_mode: strict | baseline_snapshot_v1`（R2，本 session 用 `strict` 因 finding 量可承受；schema + lint 仍須支援 baseline_snapshot 機制）
4. registry 加 top-level `baseline_snapshot:` 區塊 schema（R2，可空，但 baseline_snapshot_v1 模式需用）
5. `governance_thresholds` 加 `max_source_files_per_class: 5`、`review_when_source_files_gt: 5`（R6）
6. 既有 7 條 behavioral_only entry 全部補 `revisit_owner` + `last_reviewed_at`

**Lint 新增 check（schema 修改後同步加）**：

- `behavioral_only_missing_rationale`（補強，原 lint 漏了 rationale 檢查）
- `behavioral_only_missing_revisit_owner`（R1）
- `behavioral_only_missing_last_reviewed_at`（R4）
- `behavioral_only_review_age`（R4: >12 月 warn / >24 月 fail）
- `behavioral_only_vague_success_criteria`（黑名單 token: `TBD` / `未來` / `future` / `eventually` 等）
- `behavioral_only_revisit_chain`（revisit_when 引用其他 rule_class id 時，該 class 不得也是 behavioral_only — 避免 F17 類型的 decay chain）
- `upstream_chain_resolution`（R5: upstream_classes 引用必須解析到 active class）
- `class_size_review_threshold`（R6: source_files > 5 → warning）
- `baseline_snapshot_drift`（R2: baseline_snapshot 內 entry 若已自我修復，提醒 burndown）

**執行順序（嚴格）**：

1. Schema patch（registry.yaml 6 處變更）+ companion `enforcement-registry.md` 同步說明
2. Lint patch（新增 9 個 check）+ unit tests（每個新 check 至少 1 fail + 1 pass）
3. 既有 7 條 behavioral_only entry 補 revisit_owner + last_reviewed_at
4. Re-dry-run → 預期 surface F1-F22 + 7 個既有 entry 缺欄位（共約 29 finding）
5. Backfill F1-F22 per 修正稿
6. Re-dry-run → 0 findings
7. Wire 進 `ai-skill runtime compile`
8. Rebuild 5 platform binaries + BUILDINFO + SHA256SUMS
9. Run `ai-skill runtime compile` + 5 Phase 7 detection_commands 全 PASS
10. Owner-grouped commit + push + readback

**Open Question 補登**（這次 round-2 評審 surface 的新 OQ）：

| # | Question | 處置 |
|---|---|---|
| Q8 | `last_reviewed_at` 由誰負責更新？人工 commit-msg trailer 觸發？還是定期 governance cron？ | resolved (2026-06-01): 首版採人工更新，commit-msg 加可選 trailer `[registry-review: <class_id>]`；自動化排程列 Phase 5 後考量 |
| Q9 | `baseline_snapshot_v1` 模式如果 burndown_target_date 到期 entry 還在，governance review 由誰執行？ | open: 待 Phase 4.5 self-governance 章節決定 |
| Q10 | `upstream_classes` 形成 DAG 但若出現 cycle 怎麼處理？ | resolved (2026-06-01): lint `upstream_chain_resolution` 同時做 cycle detection，發現 cycle → fail |
| Q11 | F18 拆開的長期成本（兩個 not_mechanizable 都是 efficiency 主題）會不會反噬？ | open: 列入 Phase 4 後追蹤，coverage report 顯示「efficiency 主題 not_mechanizable 數量」作為健康指標 |

#### Phase 3 Round-3 Review — Lint 結構化 + success_criteria 可驗證性（2026-06-01）

User round-3 評審指出 round-2 land 的 9 個新 lint check 裡，**2 條 P0 lint 本質上無法可靠實作**（revisit_chain 是 NLP 問題、vague_success_criteria 黑名單脆弱），需要 schema 進一步結構化。Round-3 評分：7.5 (v3) → 8.8 (round-2 after) → 若處理本輪 3 點 → 真正的 Framework Self-Audit Layer。

**問題 → 處置 矩陣**：

| Pri | # | 問題 | 處置 |
|---|---|---|---|
| P0 | S1 | `behavioral_only_revisit_chain` lint 無法靠 NLP 解析自由文字 revisit_when 找出引用的 rule_class（"workflow_activation Phase 6.1 land" → 是 class 引用嗎？無法可靠判斷） | (a) registry.yaml schema 加 `sunset_decision.depends_on_rule_classes: []`（optional array of rule_class id）；(b) lint `behavioral_only_revisit_chain` 改寫：只看 `depends_on_rule_classes`，不解析 revisit_when 自由文字；(c) free-text revisit_when 保留作為人類可讀觸發描述，`depends_on_rule_classes` 是機械鏈結；(d) author 若在 revisit_when 提到其他 class 名稱，rule 是「**應該**同時填 depends_on_rule_classes」——這條本身不機械強制，但 documentation 寫明 |
| P0 | S2 | `behavioral_only_vague_success_criteria` 黑名單 (`TBD`/`未來`/`future`) 容易被繞過（`"coverage substantially improved"` / `"sufficient incidents accumulated"` 都過 lint 但語意等同 TBD） | (a) 換成 **positive whitelist**：lint `behavioral_only_missing_measurable_signal` 要求 success_criteria 必須包含下列 token 至少一個：數字（`\d+`）、百分號（`%`）、`rule_class` / `executor` / `lint` / `coverage` / `validator` / `hook` / `scenario` / `gate` 等具體框架 noun；(b) 同時保留原黑名單 lint `behavioral_only_vague_success_criteria`（雙閘）；(c) 兩者都 fail (P0)，因為 round-3 用例（"coverage substantially improved"）連最弱定義都過得了，必須收緊到「至少有可量化 signal」 |
| P1 | S3 | `last_reviewed_at` 容易形式化更新（"review→什麼都沒看→改日期→過 lint"） | (a) registry.yaml schema 加 `sunset_decision.last_review_summary: string` 必填（與 last_reviewed_at 同步寫入）；(b) lint `behavioral_only_missing_review_summary`：last_reviewed_at 有但 last_review_summary 空 → fail；(c) lint `behavioral_only_review_summary_too_short`：TrimSpace(last_review_summary) < 20 chars → warning（強制留下實質內容）；(d) commit-msg trailer `[registry-review: <class_id>]` 之後可加 hook 驗證 summary 是新的（與前次不同）—— 列入 future Phase 5 增強 |
| P1 | S4 | `upstream_classes` 只解決存在性 + cycle，沒解決語意類型（`upstream_classes: [document_sizing]` 對 promotion class 在語意上不合理但 lint 過） | (a) 不立即加 schema 欄位（避免過度設計）；(b) 列入 Open Question Q12 待 Phase 4.5 self-governance 時一併處理；(c) 暫時用 companion `enforcement-registry.md` 文件化「upstream_classes 應該是 promotion-chain 上游」的寫作慣例 |
| P2 | S5 | `max_source_files_per_class: 5` 命名暗示「上限」但實際 lint 是 warning（誠實性） | (a) 改名 `max_source_files_per_class` → `source_files_review_threshold`（同義但無「上限」誤導）；(b) `review_when_source_files_gt: 5` 保留（governance_thresholds 下）；(c) lint `class_size_review_threshold` 行為不變 (warning only)；(d) companion .md 文件化「沒有 hard limit，只有 review trigger」 |

**F19 重新確認**（user round-3 留意但不阻擋）：

User 認為 F19 語意更像 `pending but blocked`，不是 `behavioral_only`。但因不開 stub child plan，本 session 維持 `behavioral_only`。**記錄 future reclassification**：列入 Open Question Q13，Phase 4 land 後第一個重新分類候選。

**F12-F22 修正稿 round-3 增量**：

所有 round-2 entries 的 sunset_decision 需要新增兩個欄位：
- `last_review_summary: "initial entry creation 2026-06-01"`（首次建立的 summary）
- `depends_on_rule_classes: []`（若 revisit_when 引用其他 class，這裡填 id；否則空陣列）

具體：

| # | depends_on_rule_classes |
|---|---|
| F12 authorization_scope | `[]`（revisit_when 引用「失效模式累積」，非特定 class） |
| F14 cross_skill_references | `[]` |
| F16 document_todo_list | `[]` |
| F19 validation_scenario_governance | `[]`（revisit_when 引用 Phase 4 CLI，非 class — 雖實際 Phase 4 對應 `coverage_governance` 之類，但目前無此 class） |
| F20 decision_promotion_pipeline | `[]`（已用 upstream_classes 表達依賴 failure_learning_system，sunset 觸發是獨立判斷） |
| F21 directory_structure_governance | `[]` |
| F22 knowledge_update_flow | `[linked_updates]`（revisit_when 提到「共用 executor 是否可行」，明確指向 linked_updates） |

**Schema 變更增量（round-3 在 round-2 基礎上）**：

7. `coverage_status_spec.behavioral_only.requires` 加 `sunset_decision.last_review_summary`
8. `rule_classes[].sunset_decision.depends_on_rule_classes: []` 新增 optional 欄位
9. `governance_thresholds.max_source_files_per_class` → 改名 `source_files_review_threshold`（語意誠實）
10. 既有 7 條 behavioral_only entry 全部補 `last_review_summary` + `depends_on_rule_classes`

**Lint 增量（round-3 在 round-2 基礎上）**：

- `behavioral_only_missing_measurable_signal`（S2 positive whitelist，含 `\d+` / `%` / 框架 noun）— **取代** round-2 提的「vague_success_criteria 加強版」
- `behavioral_only_missing_review_summary`（S3）
- `behavioral_only_review_summary_too_short`（S3 warning）
- `behavioral_only_revisit_chain` 改寫為「只看 depends_on_rule_classes」（S1）— round-2 的 NLP 版本廢棄

**最終 lint check 總數**（round-2 9 個 + round-3 增量）：

```
round-2 land (9):
  behavioral_only_missing_rationale            (P0)
  behavioral_only_missing_revisit_owner        (P0)
  behavioral_only_missing_last_reviewed_at     (P0)
  behavioral_only_review_age                   (P0 fail >24m / warn >12m)
  behavioral_only_vague_success_criteria       (P0 黑名單)
  behavioral_only_revisit_chain                (改寫: 只看 depends_on_rule_classes)
  upstream_chain_resolution                    (含 cycle detect)
  class_size_review_threshold                  (warning only)
  baseline_snapshot_drift                      (governance trigger)

round-3 增量 (3):
  behavioral_only_missing_measurable_signal    (P0 positive whitelist)
  behavioral_only_missing_review_summary       (P0)
  behavioral_only_review_summary_too_short     (warning)

合計: 12 個新 check
```

**執行順序（10 步 → 維持，schema/lint patch 內容擴充）**：步驟編號不變，但 Step 1 (schema patch) 從 6 處改 10 處變更；Step 2 (lint patch) 從 9 個 check 改 12 個 check。

**Open Question 增補**：

| # | Question | 處置 |
|---|---|---|
| Q12 | `upstream_classes` 是否需 `promotion_role: source\|intermediate\|sink` 或 `artifact_type:` 語意層？ | open，列 Phase 4.5 self-governance 時決定。當前用 companion .md 文件化寫作慣例 |
| Q13 | F19 `validation_scenario_governance` 是否在 Phase 4 land 後重新分類為 pending_implementation？ | open，Phase 4 完成後評估；user round-3 已標為「第一個重新分類候選」 |
| Q14 | `last_review_summary` 是否需 "與上次 summary 不同" 機械驗證？ | open，列 Phase 5 增強（commit-msg trailer hook 比對 git diff） |
| Q15 | `behavioral_only_missing_measurable_signal` whitelist token 是否需要 i18n（中文「個」「條」「次」也算數量）？ | resolved (2026-06-01): whitelist 包含 `\d+` 正則涵蓋所有語言數字；中文 noun (`規則類別`/`執行器`) 因低出現率暫不加入，未來 surface 真實 false negative 再補 |

#### Phase 3 Round-4 Review — 結構性退階 + 分類純度修正（2026-06-01）

User round-4 評審指出 round-2 + round-3 累積出 **2 個結構性趨勢風險**，不是「缺欄位」而是「設計軌跡本身需要回頭」：

1. **behavioral_only 過度制度化**：原本是「例外機制」，現在 7 個 metadata 欄位 + 多條 lint，治理成本可能比 mechanical 還高。
2. **F19 分類遷就流程**：user 第 3 次提出 F19 應該是 `pending_implementation`，前兩輪都用「不想開 child plan」迴避，這是讓分類遷就流程而非反映真實狀態。

加上 4 個次要結構性問題（baseline_snapshot 治理不對等 / upstream_classes 演化成 DAG / measurable_signal 假陽性 / F22 與 class_size_review 邏輯衝突）。

**問題 → 處置 矩陣**：

| Pri | # | 問題 | 處置 |
|---|---|---|---|
| P0 | T1 | F19 第 3 次被指出應為 `pending_implementation`（coverage_evidence schema 已存在，只是 executor 未實作 = 已知 implementation gap，不是 behavioral） | (a) **改採方案 A**：開 stub child plan `plans/active/2026-06-01-0100-validation-scenario-governance-executor.md`（內容僅 Phase 0/1 outline + 預估 scope，不立即實作）；(b) F19 改 `coverage=pending_implementation`，`child_plan` 指向此 stub；(c) 撤回 round-3 對 F19 的 sunset_decision 段落 |
| P0 | T2 | `behavioral_only` 已累積 7 metadata 欄位（rationale + revisit_when + revisit_owner + last_reviewed_at + last_review_summary + success_criteria + depends_on_rule_classes），治理成本逼近 mechanical | (a) **strip back to 3 hard-required**：`rationale` / `sunset_decision.revisit_when` / `sunset_decision.success_criteria`；(b) 其餘 4 個（revisit_owner / last_reviewed_at / last_review_summary / depends_on_rule_classes）**降為 recommended**；(c) lint 對 recommended 缺失 → warning 而非 fail；(d) coverage report 視覺化 recommended 完成度（governance dashboard 用），但不阻塞 compile；(e) **rationale**：behavioral_only 應該是輕量例外，governance 強度應該 < mechanical，不應變第二套制度 |
| P1 | T3 | `behavioral_only_missing_measurable_signal` 對「constitution review process formalized」這類明確但無數字的 success_criteria 假陽性 | (a) lint 等級從 P0 fail 降為 **warning**；(b) 黑名單 lint `behavioral_only_vague_success_criteria`（TBD/future/eventually）維持 P0 fail（這些確實是空話）；(c) 雙閘：黑名單擋明顯空話，白名單提示 measurable signal 但不阻塞 |
| P1 | T4 | `baseline_snapshot_v1` 治理強度反而比 behavioral_only 弱（無 owner、無 review summary） | (a) `baseline_snapshot` 區塊強制加 `baseline_owner` + `baseline_review_summary`（每筆 baseline entry 都有）；(b) lint `baseline_snapshot_missing_governance`：缺 owner/summary → fail；(c) 與 behavioral_only 治理對等，不能用「baseline 是臨時的」當藉口降低治理 |
| P1 | T5 | `upstream_classes` 演化軌跡是 Governance DAG，需要架構決策避免持續長欄位 | (a) 在 backfill 前寫 ADR `constitution/ADR-XXX-registry-as-governance-dag.md`，明確決定 commit to DAG 或 freeze；(b) **建議 freeze**：upstream_classes 維持單一向上引用 + cycle detection 即可，不加 downstream_classes / promotion_role / artifact_type；(c) Q12 改為 resolved (freeze)；(d) 若未來真需要 DAG，另開 ADR 升級 |
| P2 | T6 | F22 拆分 vs `class_size_review_threshold` 是 warning 的邏輯衝突 | (a) **F22 回歸 linked_updates**（撤回 round-2/round-3 拆獨立 class 的決定）；(b) rationale：threshold 既然是 warning，linked_updates 多收 1 個 source_file 不是問題；(c) 若 linked_updates source_files 真累積到 >7（hard threshold 候選），届時 governance review 決定是否拆 class，**不是現在 preemptively 拆**；(d) class_size 治理保持輕量；(e) F22 重新歸位後，total class 數從 round-3 的 +5 → +4（cross_skill_references / authorization_scope / decision_promotion_pipeline / directory_structure_governance）+ 加 F19 child plan 後 +1 → +5（含 F19 但不含 knowledge_update_flow） |

**Round-4 修正後 F12-F22 最終版 v3**：

| # | yaml | round-4 final |
|---|---|---|
| F12 | authorization-scope | 新 class, **behavioral_only**, 3 hard fields (rationale/revisit_when/success_criteria) + 4 recommended |
| F13 | content-layering | 加進 `document_sizing.source_files`（不變） |
| F14 | cross-skill-references | 新 class, **behavioral_only**, 3 hard + 4 recommended |
| F15 | decision-efficiency | 新 class, **not_mechanizable**（不變） |
| F16 | document-todo-list | 新 class, **behavioral_only**, 3 hard + 4 recommended |
| F17 | goal-action-validation | 加進 `conversation_goal_ledger.source_files`（不變） |
| F18 | prompt-cache-efficiency | 新獨立 class, **not_mechanizable**（不變） |
| **F19** | validation-scenario-governance | **改 `pending_implementation`**, child_plan 指向新 stub plan 0100 (T1) |
| F20 | decision-promotion-pipeline | 新 class, **behavioral_only**, `upstream_classes: [failure_learning_system]`（不變） |
| F21 | directory-structure-governance | 新 class, **behavioral_only**, 3 hard + 4 recommended |
| **F22** | knowledge-update-flow | **加進 `linked_updates.source_files`**（T6 回歸） |

**Schema 變更 round-4 增量**（取代 round-3 部分定義）：

- `coverage_status_spec.behavioral_only.requires`：縮減為 3 個 hard required（rationale / sunset_decision.revisit_when / sunset_decision.success_criteria）
- `coverage_status_spec.behavioral_only.recommended`：新增區塊，列 4 個 recommended 欄位
- `baseline_snapshot` schema 加 `baseline_owner` + 每筆 entry 加 `baseline_review_summary`
- `upstream_classes` 設計凍結（不加新欄位）

**Lint 變更 round-4 增量**（取代 round-3 部分定義）：

- `behavioral_only_missing_rationale`：保留 (P0)
- `behavioral_only_missing_revisit_owner`：**降 warning**（從 fail）
- `behavioral_only_missing_last_reviewed_at`：**降 warning**
- `behavioral_only_review_age`：保留 (P0)
- `behavioral_only_vague_success_criteria`：保留 (P0 黑名單)
- `behavioral_only_missing_measurable_signal`：**降 warning** (T3)
- `behavioral_only_missing_review_summary`：**降 warning** (recommended 而非 required)
- `behavioral_only_review_summary_too_short`：**移除**（recommended 欄位不需 length lint）
- `behavioral_only_revisit_chain`：保留 (P0, 只看 depends_on_rule_classes)
- `upstream_chain_resolution` + cycle detect：保留 (P0)
- `class_size_review_threshold`：保留 (warning)
- `baseline_snapshot_drift`：保留
- `baseline_snapshot_missing_governance`：**新增 (P0)** (T4)

**最終 lint check 統計**：round-3 的 12 個 → round-4 的 12 個（替換 3 個降級 + 移除 1 個 + 新增 1 個）；P0 fail 從 round-3 的 7 個 → round-4 的 6 個（更精準、更少假陽性）。

**執行順序（10 步維持，內容更新）**：

- Step 0（新增）：寫 ADR `ADR-XXX-registry-as-governance-dag.md`（T5 freeze decision）+ 寫 F19 stub child plan
- Step 1：Schema patch（round-4 內容覆寫 round-2/3 部分定義）
- Step 2：Lint patch（12 個 check，round-4 精度版）
- Step 3-10：原 round-2 步驟不變

**Open Question 更新**：

| # | Question | 處置 |
|---|---|---|
| Q12 | upstream_classes 是否升 DAG schema？ | **resolved (round-4 T5)**: freeze，維持單一向上引用 + cycle detect；ADR 鎖定決策 |
| Q13 | F19 何時 promote 到 pending？ | **resolved (round-4 T1)**: 立即 promote 至 pending_implementation，stub child plan 0100 已開 |
| Q14 | last_review_summary 是否需 diff 驗證？ | open，但因 last_review_summary 降為 recommended，治理需求降低；列 future 增強 |
| Q15 | measurable_signal 中文 i18n？ | resolved (round-3) — round-4 後 measurable_signal 已降 warning，i18n 影響更小 |
| Q16 | behavioral_only 是否該再分層（lightweight vs strict tier）？ | open: round-4 暫採 single tier + 3 hard + 4 recommended；若未來 recommended 欄位仍持續長出，考慮拆 tier |

**為什麼 round-4 是重要的回頭路**：

Round-2/3 一直在加 lint / 加 schema 欄位，是「治 symptom」。Round-4 認識到 root cause：

- behavioral_only 想做「輕量例外」但被當「正式 coverage 類型」治理 → 解法是治理強度退階，恢復例外語意
- F19 想做「pending_implementation」但因不想開 plan 改塞 behavioral_only → 解法是開 stub plan，分類反映真實狀態
- upstream_classes 演化軌跡是 DAG → 解法是 ADR freeze 鎖定 scope，避免溫水煮青蛙

**這三條都不是欄位問題，是結構性方向錯誤的修正**。Round-4 應該是本 plan 的 final design baseline；之後若再有 round-5+，應觸發 governance review「為什麼這個 plan 需要 5 輪以上 review」（meta-failure-pattern signal）。

#### Phase 3 Round-5 Review — Final Baseline 收尾（2026-06-01）

User round-5 評審確認 **round-4 是 final design baseline**（不是因為「需要更多 round」），明確指出最大改善是「開始承認哪些東西不應該被機械化」。Round-5 不新增 lint / schema 欄位，只**收掉 round-4 留下的 3 個規格空洞**。

> **Meta clarification**：round-5 不觸發前述 meta-failure-pattern signal。Meta-pattern 指的是「持續加 lint / schema 卻仍然有 gap」的死循環；round-5 是 round-4 baseline 的規格收斂，user 已明示認可方向。Round-6+ 才是 trigger。

**問題 → 處置 矩陣**：

| Pri | # | 規格空洞 | 處置 |
|---|---|---|---|
| P0 | U1 | `behavioral_only_review_age` (P0 fail >24m) vs `last_reviewed_at` 已降 recommended（缺則 warning）—— 規格沒寫「last_reviewed_at 缺失時 age 如何計算」 | **明確規範**：lint `behavioral_only_review_age` 僅在 `last_reviewed_at` **present** 時觸發；缺失時由 `behavioral_only_missing_last_reviewed_at` (warning) 單獨處理。**不雙觸發**。Rationale：last_reviewed_at 既然 recommended，缺失只應 warning；強制 fail 違反退階意圖。Schema 文件化：「age_unknown=skip, missing=warning」 |
| P0 | U2 | F19 改 `pending_implementation` 需要 `child_plan`，但「stub plan」是否合法 round-4 沒定義 | **`child_plan` 合法性 schema**：(a) 路徑必須 resolve 到 `plans/active/*.md`；(b) plan 必須包含 `## Phase 0` heading（最低 outline）；(c) plan frontmatter 或內文必須有 owner 標示；(d) plan 必須有非空 `## Validation Plan` 或 `## Acceptance` 區塊（避免「空殼 stub」）。**stub 合法**只要滿足 (a)-(d)。新 lint `pending_implementation_child_plan_validity` (P0 fail)：違反 (a) fail；違反 (b)/(c)/(d) warning（漸進壓力，不立即阻塞 stub） |
| P0 | U3 | `upstream_classes` freeze ADR 內容若不夠精確，半年後容易被擴張 | **ADR 必填區塊**：(a) 「Scope boundary」明寫 upstream_classes **IS for**: promotion traceability；(b) 「Scope boundary」明寫 upstream_classes **IS NOT for**: execution ordering / dependency injection / runtime orchestration / DAG-based scheduling；(c) 「Supersession clause」：任何跨越此 boundary 的新欄位（downstream_classes / promotion_role / artifact_type / dependency_type 等）必須先寫新 ADR 顯式 supersede 本 ADR，不得 silent 擴張；(d) ADR 標 `status: active` 加 `revision_policy: supersede_required`（registry self-governance lint 可未來檢查此政策） |
| P2 | U4 | baseline_snapshot 治理（owner + summary 必填）反超 pending_implementation（只 child_plan）—— 治理倒掛 | **記錄為 intentional asymmetry，不立即修**。Rationale：baseline_snapshot 是「technical debt」必須有 owner + summary 才能 burndown；pending_implementation 是「tracked work」child_plan 本身已內含 owner + acceptance（U2 schema 強制）。兩者治理形式不同但實質強度相當，asymmetry 是 surface 差異不是漏洞。Companion .md 文件化此設計選擇 |

**Q16 標為 Single Highest-Priority Open**：

> Q16 (behavioral_only tier separation) 是 round-5 後**唯一需要主動追蹤**的 open question。Trigger 條件：當第 2 個 entry 出現「目前無 executor 但理論可做」（類似原 F19 性質），不應再硬塞 behavioral_only。屆時必須回到 Q16，要嘛開 child plan（如 F19 處置），要嘛拆 tier（如 `behavioral_only_lightweight` / `behavioral_only_strict`），不得第三度遷就分類。

**Round-5 增量 lint（共 1 個，非新治理範疇，只是 U2 規格化）**：

- `pending_implementation_child_plan_validity` (P0)：U2 處置

**最終 lint check 統計**：round-4 的 12 個 → round-5 的 13 個（純增 U2 lint）。P0 fail 從 round-4 的 6 個 → round-5 的 7 個（U2 加回但只在 stub plan 不存在路徑時 fail，存在但缺 outline/owner/acceptance 只 warning）。

**Schema 變更 round-5 增量**：

- `coverage_status_spec.pending_implementation.child_plan_validity` 區塊新增 4 條規則（U2）
- `coverage_status_spec.behavioral_only.lint_behavior` 明寫 review_age 與 missing_last_reviewed_at 互斥規則（U1）
- ADR template 加 `revision_policy: supersede_required` 欄位定義（U3 配套）

**執行順序更新**（11 步 → 維持 11 步）：

- Step 0：寫 ADR + F19 stub plan
  - **Step 0a 新增子步驟**：F19 stub plan 必須通過 U2 schema (a)-(d)；如果只是「title + 一段話」會被新 lint 警告
  - **Step 0b 新增子步驟**：ADR-XXX-registry-as-governance-dag.md 必須含 U3 三區塊（Scope boundary IS / IS NOT / Supersession clause）
- Step 1-10：原 round-4 順序不變

**Round 演進總表**：

| Round | 性質 | 主要產出 |
|---|---|---|
| Round-1 | 建立 registry | 初版 schema + rule_classes |
| Round-2 | 補治理欄位（加法） | revisit_owner + bootstrap_mode + upstream_classes + class_size_review |
| Round-3 | 治理機械化（過頭） | last_review_summary + measurable_signal + depends_on_rule_classes |
| **Round-4** | **承認機械化邊界（減法）** | **strip behavioral_only to 3 hard + 4 recommended, F19 → pending, F22 回歸** |
| Round-5 | 收尾規格空洞（精修） | U1/U2/U3 規格化, Q16 標 single open, ADR scope lock |

**Final Baseline 宣告**：

Round-5 land 後，本 plan 進入 **frozen design baseline**。任何 round-6+ review **必須先回答**：

1. 這個 concern 是 round-1 ~ round-5 已 surface 但未解決的，還是真正新的？
2. 若是新的，是「實作中 surface」還是「pre-implementation review」？
3. 若是 pre-implementation review 第 6 次，觸發 meta-failure-pattern signal — 寫入 `enforcement/failure-patterns/excessive-pre-implementation-review.md`，繼續執行 round-5 baseline，不再 round-6 修改

**這條 frozen baseline 本身是 registry 治理閉環的一部分**：證明 Layer 2.5 不只治理別人，也治理自己（Q16 不會在實作前無限延伸）。

#### Phase 3 Step 1 Schema Patch v2 — 自我審查 10 findings（2026-06-01）

Schema patch v1 (commit 2a86fce) land 後，**Step 2 lint 實作開跑前**做了一輪 self-audit，surface 10 個 schema ambiguity。User round-6 評審（不觸發 meta-pattern signal — 此輪是「實作中 surface」非 pre-implementation review）分級確認後，schema patch v2 (commit c9c37b1) 處理 7 個，跳過 3 個。本節記錄審查結果作為治理軌跡。

**Round-6 性質定義**：本輪 review 在 schema patch v1 land 後 + Step 2 開跑前 surface 問題，符合 round-5 §Final Baseline 宣告的「實作中 surface」分類，**不觸發** `excessive-pre-implementation-review` failure pattern。Round-5 baseline 本身未動，只是補規格細節。

**10 findings + 處置矩陣**：

| Pri | # | 問題 | User Grading | 處置 (c9c37b1) |
|---|---|---|---|---|
| P0 | A1 | `self_governance.lint_rules` 混 commit-msg phase (R1-R5) 與 compile-time phase (R6-R9)；命名漂移 (`R6_upstream_chain_resolution` 拼接 key) | 必修 | **拆 namespace**: `commit_msg_lint_rules: {R1..R5}` + `compile_time_lint_rules: {R1..R4}` (重編號)。不用 phase 欄位 (Go code 還要 filter)，直接拆 dict |
| P0 | A2 | `bootstrap_mode` (strict/baseline_snapshot_v1) 與 `baseline_snapshot.enabled` 雙 state source，可出現 illegal 組合 `{strict, enabled=true}` | 必修 (user 強化版) | **砍掉 `baseline_snapshot.enabled`**，`bootstrap_mode` 是 single source。Lint R3 (`baseline_snapshot_missing_governance`) trigger 改 `bootstrap_mode == baseline_snapshot_v1` |
| P0 | A3 | `upstream_classes_scope.is_for / is_not_for` 看起來像強制規則但 lint 只檢查 cycle + reference，讀者誤導 | 必修 | 新增 `mechanically_enforced: [reference_resolution, cycle_detection]` + `documentation_only: [is_for, is_not_for, supersession]` 區塊明寫邊界 |
| P0 | A4 | `executors[]` / `executors_planned[]` shape 從未形式化，只靠範例推導，未來必出現 `executor:` / `planned_executor:` / `executors_plan:` typo | 必修 | 新增 top-level `executor_schema` 區塊: required=[file, symbol, executor_kind], optional=[hook_phase, block_or_warn, notes, instance_count], used_by=[mechanical.executors[], pending.executors_planned[]] |
| P1 | B5 | `child_plan_validity.a_path_resolves` 未處理 anchor (registry 既有 source_files 用 `runtime/core-bootstrap.yaml#per_session_obligations`) | 應改 | (a) 規則加 `path.split('#')[0]` 標準化，mirror source_files anchor handling |
| P1 | B6 | `source_files_review_threshold: 5` 無 per-class override，F22 backfill 後 linked_updates 可能噪音；最初我提案 `acknowledged_size_warning: bool` 但 user 指出 bool 會被 rubber-stamp | 哲學分歧, user 強化 | 採 user 強化版: `rule_classes[].size_review_exemption_rationale: <string>`，warning **仍 emit** 不 suppress，但含 maintainer rationale 而非 bare threshold breach。強制寫 rationale 避免 rubber-stamp |
| P2 | C7 | `mechanical.requires: [executors, rationale]` 但既有 mechanical entries 不一定都有 rationale | audit 後決定 | **Audit 結果**: 14/14 mechanical entries 已有 rationale ✓，**skip**，不改 schema |
| P2 | C8 | ADR-010 與 registry.yaml 都有 `adr_revision_policy: supersede_required`，重複定義 | 可改可不改, user 偏保留 | **保留重複** — registry 是 executable contract, ADR 是歷史決策, lint 不應該解析 ADR markdown |
| P2 | C9 | `behavioral_only_missing_measurable_signal` whitelist 是 token-level，會被 `"no new lint required"` 否定句繞過 | 建議改 | 加 comment 明寫「TOKEN-LEVEL heuristic. Presence of keyword does NOT imply measurable; absence does NOT imply unmeasurable. Lint 是 documentation prompt, 非 semantic judge.」 |
| P2 | C10 | 應在 schema 註明「Step 3 之前預期會 surface ≥7 個 recommended-missing warnings」 | 不建議 | **不入 canonical** — 時間性 expectation 幾個月後過時；放 migration plan |

**處置統計**: 必修 7 個 (A1-A4 + B5-B6 + C9) → c9c37b1 全部 land。跳過 3 個 (C7 audit 結果無需改 / C8 哲學偏好保留 / C10 不入 canonical)。

**Lint check 數量更新** (round-5 13 個 → schema patch v2 後仍 13 個，但**命名與 namespace 變更**):
- commit_msg_lint_rules: R1-R5 (5 個, 不變)
- compile_time_lint_rules: R1-R4 (4 個, 從 round-5 的 R6-R9 重編號)
- behavioral_only_* lints: 4 個 hard FAIL + 5 個 WARNING (不變)

**A4 影響 Step 2 實作**: Go code 現在有明確 `executor_schema` 可 unmarshal target，避免 schema patch v1 留下的 typed-struct 模糊。Step 2 enforcement_registry_lint.go 的 `registryExecutor` struct 已有 yaml tags，與 executor_schema 對齊。

**User 哲學原則確認**:
- **不允許 rubber-stamp suppression** (B6 string vs bool 之選)
- **executable contract 與 historical decision 邊界清楚** (C8 保留重複)
- **時間性 expectation 不入 canonical** (C10)
- **lint 不做 NLP / 語意判斷** (A3 + C9 明寫邊界)

#### Phase 3 Step 6/7/9 Review — warning 分級 + wire 策略 + Step 9 scope（2026-06-01）

User round-7 評審（實作中 surface，不觸發 meta-pattern）對 Step 6 warning 處置、Step 7 wire 策略、Step 9 verify scope 給出裁決。無衝突，全部採納。

**Step 6 — warning 分級處置**：

| Warning 類型 | 處置 | 理由 |
|---|---|---|
| `class_size_review_threshold` | **保留** | registry 自身語意：> threshold ≠ error，= 提醒檢查是否該拆。cognitive_mode_governance 已記 `size_review_exemption_rationale`（cohesive cognitive-modes-* family），warning 正確且應持續顯示 |
| `child_plan_validity` (b/c/d) | **修掉** | 非設計哲學問題，是客觀 schema 缺口（缺 owner / Phase 0 / Validation）。plan 是自己控制的檔案，無必要永久背 warning（warning fatigue → 真正重要的 warning 沒人看） |
| `missing recommended field` | 補掉 | 同上（已在 Step 3 處理 7 個 behavioral_only） |
| `schema ambiguity` | 修掉 | 已在 schema patch v2 處理 |

Step 6 結果：5 findings → **1 finding (class_size only), 0 FAIL** = reviewer 定義的「理想」狀態。child_plan warning 清除方式：(1) lint regex 太嚴是 bug（只認 h2 `## Phase 0`，但 1900/2000 合法用 h3 `### Phase 0` nested under `## Phase Plan`）→ 放寬為 `^#{2,}`；(2) 1900/2000 真缺 owner → 補 `Owner:` line。

**Step 7 — wire 策略確認**：

採 severity-aware exit code，與 registry 既有 `FAIL = contract violation / WARNING = governance signal` 語意一致：

```
FAIL count > 0  → ExitValidationFailed
FAIL count = 0  → ExitSuccess
WARNING count > 0 → print warning summary, exit 0
```

關鍵原則（reviewer）：**若 compile 因 WARNING fail，severity model 就壞了**（WARNING 與 FAIL 沒差別）。compile output 必須印 FAIL/WARNING counts + 分組 summary，讓 maintainer 一眼看出「compile success 但有治理債務」，而非只 dump findings。

**Step 9 — verify scope 裁決：不在 Phase 3 補 enforcement CLI**：

Phase 7 的 5 個 scenario 的 `detection_command` 寫的是 `ai-skill enforcement lint --check <type>` —— 這驗證的是 **CLI contract**（Phase 4 public surface），不是 lint engine（Phase 3）。為了讓 scenario 跑而提前補 CLI = scope leak（argument parser / output formatter / help text / exit code contract 全部提前進來，Phase 3 被迫做半套 Phase 4）。

處置：
- **Phase 3 verify** 用 Go test / integration test 驗 `LintEnforcementRegistry()` + `buildRuntimeCompileResult()`
- **Phase 4 verify** 才驗 `ai-skill enforcement lint` CLI contract
- 5 個 Phase 7 scenario metadata 標 `phase: 4` + `requires_capability: enforcement_cli`，明示「scenario 已存在但尚未可執行（CLI 未實作）」，而非為跑 scenario 提前實作 CLI

架構分層確認：
```
Phase 3: registry + lint engine + compile integration  (Go test 驗)
Phase 4: enforcement lint CLI public surface           (scenario detection_command 驗)
```

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
- [x] Lint 實作 + unit tests（13 lint checks: 5 v1 + 9 round-4/5; 35 Go tests PASS）— `scripts/ai-skill-cli/internal/app/enforcement_registry_lint.go`
- [x] `ai-skill runtime compile` 整合（Step 7: severity-aware — FAIL blocks, WARNING prints + exit 0; commit 381869a + binary 5862eba）
- [x] 第一次 compile run 已處理：Q1 hard-block 下 dry-run 48 findings 全數 backfill 至 FAIL 0 / WARNING 1（class_size, acknowledged exemption）

**Phase 3 land 完成（2026-06-02）**。Step-by-step：Step 0 ADR-010 + F19 stub (ccff2c2) → Step 1 schema patch (2a86fce + c9c37b1) → Step 2 lint 9 checks (599fa05) → Step 3 behavioral backfill (43677a3) → Step 5 F1-F22 backfill (05f4b19) → Step 6 child_plan warning clear (e39ca50) → Step 7 wire compile (381869a) → Step 8 binary rebuild (5862eba) → Step 9 scenario phase-4 marking + verify. 已知 out-of-scope: Phase 4 CLI、Phase 4.5 self-gov lint、Phase 5 bootstrap integration。pre-existing unrelated: runtime-index.sqlite stale checksum (spawn task flagged)。

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
- [x] CLI subcommand 實作 — `ai-skill enforcement lint` + `ai-skill enforcement coverage` (Phase 4 land 2026-06-02, scripts/ai-skill-cli/internal/app/enforcement.go)
- [x] 輸出格式（text + JSON + markdown for governance dashboards）— 三 format + `--detail` + `--diff <ref>` + `--self-check` 全 land
- [x] 文件化（command-contract.md）— enforcement lint / coverage 兩段 + 副作用登錄表項目同步；ai-tools/agent reference 留待 Phase 5 一併處理（與 bootstrap integration 同步）
- [ ] CI integration：Pull Request 自動跑 coverage diff，新增規則沒填 coverage 直接 PR check 失敗 — deferred to Phase 4.5 / Phase 5（CI workflow 變動屬另一 owner group，本 session 範圍只到 CLI land）
- [x] 4 個 Phase 4-blocked scenario 翻 runnable — orphan-rule / missing-executor / behavioral-without-rationale / coverage-cli-output-format-v1 detection_command 改成單一 cross-platform CLI invocation（exit-code based assertion），fixture repo 提交於 `validation/scenarios/enforcement/fixtures/`

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
- [x] Phase 4.5 lint rules R1-R5 實作 — R1 (trailer + rationale) / R2 (demotion needs ADR) / R3 (promotion needs symbol_exists) 在 `scripts/ai-skill-cli/internal/app/enforcement_transition.go` `checkRegistryTransitions`；R4 (deprecated past removal_date+30d) / R5 (research_required past estimated_unblock_timeline) 在 `enforcement.go` `buildCoverageReport` + `partitionGovernanceAlerts`
- [x] commit-msg `validateEnforcementRegistryTransition` validator — 註冊為 `obligation.commit.enforcement_registry_transition`（per_commit_obligations 第 20 個），dispatch via `commitMsgValidatorRegistry`；opt-out `[skip-registry-transition]`
- [x] coverage report 章節「Governance Alerts」整合 sunset + removal_date breach + research_required timeout — text 用 `⚠ Governance Alerts (Phase 4.5 R4/R5):` 段，markdown 用 `## Governance Alerts (Phase 4.5 R4/R5)` h2，JSON 在 `alerts[]` 以 `R4_*` / `R5_*` kind 區分；與普通 alerts 拆分顯示
- [x] `enforcement/enforcement-registry.md` companion 加 §Self-Governance — 更新 Status Transition Matrix（加強制層欄）、Self-Governance Lint Rules 表（R1-R5 行為矩陣）、demotion ADR rationale、開發者快速指南
- [x] Scenario-driven verification — 新增 `registry-transition-demotion-without-adr-v1.yaml` + `registry-transition-promotion-verification-gap-v1.yaml`，搭配 self-contained fixture repo（old/new registry yaml + commit-msg.txt + stub hooks.go），detection_command 是單一 cross-platform CLI invocation
- [x] CLI standalone surface — `ai-skill enforcement transition-check` 暴露 R1/R2/R3 engine 給 scenario / CI / local debug；與 commit-msg validator 共用 `checkRegistryTransitions` engine

### Phase 5 — Bootstrap Integration

把 enforcement-registry 接進 bootstrap：

- [x] `runtime/core-bootstrap.yaml` 加 contextual_activation：「修改 enforcement rule 時必讀 enforcement-registry.md」 — `activation.enforcement_registry_editing`（land 2026-06-02），6 個 load_when trigger（editing rule yaml/md / registry yaml / 新增 rule_class / 改 coverage 值）
- [x] `ai-skill runtime receipt` 輸出加 enforcement coverage 摘要（一行） — 實際 land 格式 `Enforcement: classes=<N> mechanical=<n> behavioral=<n> not_mech=<n> pending=<n> research=<n> deprecated=<n>`（用 6-bucket enum 而非原 plan 的 bound/behavioral/orphan 三項，因為 Phase 1.3 已用 rule_class 抽象取代 rule instance binding）；registry 不可解析時降為 `skipped` check
- [x] commit-msg hook 加 validator：若 commit 變更 enforcement rule yaml/md，必須同步更新 enforcement-registry.yaml；未同步 → reject — `validateEnforcementRuleRegistrySync`（21 個 commit-msg validator，dual-gate 與 Phase 3 compile-time `orphan_rule` 並存）；opt-out `[skip-enforcement-registry-sync]`

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

- [x] 新建 `enforcement/failure-patterns/rule-without-executor.md` —— 把 2026-05-31 session 暴露的 5 個 instance 集中記錄，作為 meta-pattern 的 inaugural reference（land 2026-06-03，5 instance table + 6-value coverage mapping + cross-link 到 Phase 7 scenarios）
- [x] 兩個 child plan 各加 link 指回本 plan，明示「本 plan 是 meta-pattern 的 instance」（1900 header L9 + 2000 header L11 都已標明 parent meta-plan，2000 §Source 補了 v2 起明確標記為 meta-plan instance 的歷史脈絡）
- [x] `enforcement/README.md` 加章節「Mechanical Enforcement Registry」指向本 plan + registry yaml（§Mechanical Enforcement Registry（Layer 2.5）L78-84，cross-link 到 rule-without-executor pattern + parent plan）

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
