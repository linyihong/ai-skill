# Mechanical Enforcement Registry

**Status**: `draft`
**世代**：Gen 3 runtime hardening（**meta-pattern**：framework-level invariant，非個案修補）
**建立日期**：2026-05-31
**最後更新**：2026-05-31（initial draft）

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

### Decision

建立 **Mechanical Enforcement Registry** + **Coverage Audit** + **Compile-time Lint**：

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
| **Rule 與 Executor 雙向 binding**（registry entry 必含兩端） | Rule 沒 executor → 規則只是文本；executor 沒 rule → 行為無 governance 依據。雙向都列才能 audit |
| **顯式 `behavioral_only: true` 是合法宣告，但需 rationale + sunset_decision** | 不是所有規則都該機械強制（e.g., 軟性語氣建議）。但 behavioral 必須是主動選擇，不是預設遺漏。每條 behavioral_only 要回答「為什麼不機械化」與「何時 revisit」 |
| **Compile-time lint，非 runtime lint** | `ai-skill runtime compile + refresh` 跑 lint，任何 registered rule 沒對應 executor 直接 compile fail。讓問題出現在「我加新 rule」當下，而非「使用者半年後追問」 |
| **CLI: `ai-skill enforcement coverage`** | 列當前覆蓋率 + 列舉所有 `behavioral_only: true` 條目 + 列舉所有「rule 存在但 registry 沒 entry」的孤兒。Audit 工具不只用於 lint，也用於 framework health dashboard |
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

### Phase 1 — 既存 Rule / Executor Inventory（discovery 階段）

跑一次手動 audit 列出當前狀態：

#### Phase 1.1 — 列舉所有「規則來源」

掃 `enforcement/*.yaml`、`runtime/*.yaml`、`governance/*.yaml`、`routing-registry.yaml`、`*.md` 內聲明的 `obligation.*` / `gate.*` / `rule.*` / `activation_*` id。產出：`enforcement/inventory/rules.yaml`（暫存 audit 結果）。

預估數量：
- `runtime/core-bootstrap.yaml` per_session / per_turn / per_commit obligations：~22
- `runtime/runtime.db` obligations / gates：23 + 25 = 48
- `routing-registry.yaml` route activation：57（其中 7 條有 triggers）
- `enforcement/sanitization.md` banned content categories：5
- 其他 enforcement rules：估 ~30+
- 總計：~150+ 規則項

#### Phase 1.2 — 列舉所有「執行器」

掃 `scripts/ai-skill-cli/internal/app/hooks.go` + `runtime/*`：所有 validator 函式、所有 PreToolUse hook、所有 SessionStart hook、所有 commit-msg validator。

預估：~30 個 executor symbol。

#### Phase 1.3 — 第一次 Coverage Matrix（揭發現況）

每條 rule 標記：
- `bound` — registry 有 binding 且 executor 存在
- `orphan_rule` — rule 存在但無 executor（這次 session 暴露的 5+ 個都屬此類）
- `orphan_executor` — executor 存在但無對應 rule（規則文本可能 stale）

Coverage matrix 作為 Phase 2 binding 工作的基礎。預計 orphan_rule 比例 ≥ 30%。

### Phase 2 — 建立 `enforcement-registry.yaml`

定義 binding schema：

```yaml
# enforcement/enforcement-registry.yaml
schema_version: 1
bindings:
  - id: bootstrap_receipt
    rule_source: runtime/core-bootstrap.yaml#per_session_obligations[obligation.bootstrap.receipt]
    executor:
      file: scripts/ai-skill-cli/internal/app/hooks.go
      symbol: validateBootstrapReceiptPresent
      hook_phase: PreToolUse
    enforcement_layer: mechanical
    block_or_warn: block
    rationale: |
      Receipt 是 session integrity 的 anchor。bypass 會讓 agent 跳過必讀
      規則。歷史上 SessionStart hook 注入完整 receipt 反而讓 agent 抄
      placeholder，因此採用 read-log gate 機械強制。

  - id: workflow_activation
    rule_source: knowledge/runtime/routing-registry.yaml#activation_triggers
    executor:
      file: scripts/ai-skill-cli/internal/app/detector.go
      symbol: DetectWorkflows
      hook_phase: PreToolUse + RuntimeContext write
    enforcement_layer: mechanical
    block_or_warn: block_on_canonical_workflow_paths
    status: pending  # 由 child plan workflow-activation-engine 實作
    child_plan: plans/active/2026-05-31-1900-workflow-activation-engine.md

  - id: sanitization
    rule_source: enforcement/sanitization.yaml#banned_patterns + incident_score
    executor:
      preflight:
        file: scripts/ai-skill-cli/internal/app/hooks.go
        symbol: validateSanitizationOnWrite
        hook_phase: PreToolUse
        block_or_warn: warn
      commit:
        file: scripts/ai-skill-cli/internal/app/hooks.go
        symbol: validateSanitizationOnCommit
        hook_phase: commit-msg
        block_or_warn: block
    enforcement_layer: mechanical
    status: pending
    child_plan: plans/active/2026-05-31-2000-mechanical-sanitization-validator.md

  - id: capability_discovery_fallback
    rule_source: governance/lifecycle/capability-discovery-philosophy.md
    executor:
      status: pending
    enforcement_layer: behavioral_only
    rationale: |
      Discovery 是「detector miss 後 fallback」的探索動作，不該每 turn 強制。
      預設行為強制，等實作 detector 後再 binding。
    sunset_decision: |
      child_plan workflow-activation-engine Phase 6.1 完成後重新評估是否
      mechanical。若新規模 ≥ 50 個 active route 時仍未 mechanical，重提
      sunset review。

  # ... (continue with all bindings discovered in Phase 1)

behavioral_only_rules:
  # 顯式列出「我們選擇不機械化」的規則，每條附 rationale + sunset_decision
  - id: ...
```

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

### Phase 4 — CLI Coverage Report

新增 `ai-skill enforcement coverage`：

```
$ ai-skill enforcement coverage
Enforcement Coverage Report (2026-XX-XX)
═══════════════════════════════════════
Total declared rules: 152
  ├─ Mechanically enforced (bound):   38  (25%)
  ├─ Behavioral_only (explicit):       7  ( 5%)
  └─ Orphan rules (no binding):      107  (70%)  ⚠️

Executor coverage:
  ├─ Active executors with bindings:   28
  └─ Orphan executors (no rule):        3  ⚠️

Top orphan rules by severity (P0/P1 first):
  P0  obligation.commit.sanitization_diff   (sanitization plan pending)
  P0  obligation.workflow.activation_evidence (workflow plan pending)
  ...

Behavioral_only rules pending sunset review:
  capability_discovery_fallback  — sunset: workflow-engine Phase 6.1
  ...
```

產出：
- [ ] CLI subcommand 實作
- [ ] 輸出格式（text + JSON）
- [ ] 文件化（README + ai-tools/agent reference）

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

## Source

2026-05-31 session：使用者連續 6 輪追問 / 評審，依序暴露：
1. sqlite3 vs ai-skill CLI 認知偏差
2. workflow activation gap
3. Discovery vs Detector 混淆
4. intelligence 預設 advisory 風險
5. sanitization gate 自我觸發失敗
6. **本 plan 對應**：以上 5 條共同根因是「Knowledge layer 有規則，Runtime layer 沒執行器」meta-pattern，建議建立 governance-level coverage audit

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
