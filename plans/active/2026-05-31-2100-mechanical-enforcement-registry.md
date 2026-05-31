# Mechanical Enforcement Registry

**Status**: `draft-v2`
**世代**：Gen 3 runtime hardening（**meta-pattern**：framework-level invariant，非個案修補）
**建立日期**：2026-05-31
**最後更新**：2026-05-31（v2 — 整合 round-7 評審：Layer 2.5 framing、Rule Class 取代 rule-level binding、新增 `not_mechanizable` coverage status、priority 反轉為 P1）
**Priority**：**P1**（v2 起，從原 P2 提升）—— 完成本 plan 後，兩個 child plan 重啟時會被本 plan 的 coverage lint 強制回答 "mechanical / behavioral_only / not_mechanizable / pending" 問題，從根本上預防未來同模式 bug

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

### Architectural Framing — Missing Layer 2.5

第七輪評審指出 Ai-skill 既有架構是三層：

```
Layer 1  Knowledge       (enforcement/, governance/, workflow/, ...)
Layer 2  Runtime         (scripts/ai-skill-cli/, hooks.go, runtime.db)
Layer 3  Governance      (constitution/, architecture/, plans/)
```

但缺一層：

```
Layer 2.5  Coverage Verification   ← NEW (本 plan 建立)
            Rule ←binding→ Executor 的結構性驗證層
```

沒有 Layer 2.5 就沒有「rule 寫好、executor 沒接」的結構性偵測機制，所有此類 bug 必須等使用者發現。本 plan 不是「補一個 executor」（child plan 的工作），而是**建立 Layer 2.5 本體**。

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
| **4-value coverage enum**（v2 採納評審 #3） | `mechanical` / `behavioral_only` / **`not_mechanizable`** / `pending`。`not_mechanizable` 是 v2 新增的關鍵狀態 —— 隔離「永遠不該機械化」（如主觀寫作品質）與「現在沒機械化但應該機械化」，避免 review queue 永遠塞著無解項目 |
| **`behavioral_only` 與 `not_mechanizable` 各自要求不同 metadata** | `behavioral_only` 需 `rationale` + `sunset_decision`（何時 revisit）；`not_mechanizable` 需 `rationale` + `objective_validation_impossible_because`（為何永遠不可機械化）。不同 metadata 強制不同思考深度 |
| **Compile-time lint，非 runtime lint** | `ai-skill runtime compile + refresh` 跑 lint，任何 registered rule_class 缺 coverage 宣告直接 compile fail。讓問題出現在「我加新 rule」當下，而非「使用者半年後追問」 |
| **CLI `ai-skill enforcement coverage` 是主要產出**（v2 評審 #4 強調） | 不只 audit 既有狀態，更重要的是**強制新規則寫作者回答 coverage 問題**。Phase 4 不是錦上添花，是 framework invariant 真正落地的地方 |
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

# ─── 4-value enum schema ────────────────────────────────────
coverage_status_spec:
  mechanical:
    requires: [executors[].symbol exists]
    lint_behavior: source 變更須同步 executors[].symbol 存在
  pending:
    requires: [child_plan, target_promotion]
    lint_behavior: child_plan 必須是 active plan 路徑
  behavioral_only:
    requires: [rationale, sunset_decision.revisit_when, sunset_decision.revisit_owner, sunset_decision.success_criteria]
    lint_behavior: 缺任一欄位 → compile fail
  not_mechanizable:
    requires: [rationale, objective_validation_impossible_because]
    lint_behavior: 缺任一欄位 → compile fail；附帶 governance review 時可挑戰是否真的 not_mechanizable
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

  ├─ Mechanical:        12  (50%)
  ├─ Behavioral only:    5  (21%)  — explicit governance choice
  ├─ Not mechanizable:   3  (12%)  — out of review queue
  └─ Pending:            4  (17%)  — implementation in progress

Pending (active implementation plans):
  workflow_activation   → plans/active/2026-05-31-1900-workflow-activation-engine.md
  sanitization          → plans/active/2026-05-31-2000-mechanical-sanitization-validator.md
  intelligence_classification → ...

Behavioral_only awaiting sunset review:
  capability_discovery  — revisit when: workflow_activation Phase 6.1 land
  rule_weight           — revisit when: 3+ detectable P0 patterns surface
  ...

Not_mechanizable (closed, will not appear in review queue):
  tool_neutral_documentation  — subjective writing judgment
  rule_writing_quality        — would game readability metrics
  ...
```

**強制觸發場景**：

| 場景 | Coverage report 行為 |
|---|---|
| 新增 `enforcement/<new-rule>.yaml` 但未在 registry 出現 | compile fail：`new rule class not registered` |
| Registry 加新 entry 但漏寫 `coverage` field | compile fail：`missing coverage field` |
| `coverage: behavioral_only` 但 `sunset_decision.success_criteria` 空白 | compile fail：`behavioral_only requires success_criteria` |
| `coverage: not_mechanizable` 但 `objective_validation_impossible_because` 空白 | compile fail：`not_mechanizable requires impossibility rationale` |
| 既有 mechanical class 的 executor symbol 在 hooks.go 找不到 | compile fail：`executor symbol missing` |
| 新增 mechanical executor 但無對應 rule_class | compile warning：`orphan executor` |

這意味未來任何新規則寫作者，在 compile 那一刻就被強制回答：
> 「這條規則是 mechanical / behavioral_only / not_mechanizable / pending？」

無法迴避、無法 silent leak、無法等使用者半年後追問。

產出：
- [ ] CLI subcommand 實作
- [ ] 輸出格式（text + JSON + markdown for governance dashboards）
- [ ] 文件化（README + ai-tools/agent reference）
- [ ] CI integration：Pull Request 自動跑 coverage diff，新增規則沒填 coverage 直接 PR check 失敗

### Phase 5 — Bootstrap Integration

把 enforcement-registry 接進 bootstrap：

- [ ] `runtime/core-bootstrap.yaml` 加 contextual_activation：「修改 enforcement rule 時必讀 enforcement-registry.md」
- [ ] `ai-skill runtime receipt` 輸出加 enforcement coverage 摘要（一行）：
  `Enforcement: bound=38/152 (25%) behavioral=7 orphan=107`
- [ ] commit-msg hook 加 validator：若 commit 變更 enforcement rule yaml/md，必須同步更新 enforcement-registry.yaml；未同步 → reject

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
| Q1 | Phase 1.3 預期 orphan rate 70%；第一次 lint 直接 block compile 還是先 warn 一個 grace period？ | still-open — 建議 grace period 一個 release cycle，期間 backfill bindings + 寫 behavioral_only justifications；之後改 block |
| Q2 | `behavioral_only` 的 `sunset_decision` 強制格式？開放 free text 容易變空話 | still-open — 建議至少要列「condition 何時 revisit」+ 「revisit owner」，由 compile-time lint 校驗 schema |
| Q3 | enforcement-registry 與 `runtime/core-bootstrap.yaml` per_*_obligations 的關係？兩者都列 obligation id | resolved → core-bootstrap.yaml 是 phase-aware obligation lifecycle，enforcement-registry 是 cross-phase binding 視圖。兩者互補：bootstrap 講「何時 fire」，registry 講「何處 enforce」 |
| Q4 | Orphan executor（code 有 validator 但 registry 沒 entry）該強制 binding 還是允許 internal helper？ | still-open — 建議只強制「exported 或被 hook dispatcher 註冊」的 executor 必須有 binding；internal helper 不算 |
| Q5 | Discovery 在 future plan 是否視為「fallback rule」並進 registry？ | resolved → 是。`capability_discovery_fallback` 已列在 Phase 2 schema 範例為 behavioral_only，等 workflow-activation-engine Phase 6.1 實作後改 mechanical |
| Q6 | 已有的 11 個 commit-msg validators 是否全部要在 registry 補 binding 才算 audit 完成？ | still-open — 預期 yes，但工作量大，列入 Phase 1.3 audit |

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
