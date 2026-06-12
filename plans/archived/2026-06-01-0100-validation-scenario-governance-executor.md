---
id: 2026-06-01-0100-validation-scenario-governance-executor
plan_kind: sub
status: completed
owner: linyihong
created: 2026-06-01
parent: 2026-05-31-2100-mechanical-enforcement-registry
required_for_completion: false
sub_plan_reason: >
  Phase 3 Round-4 T1 / Round-5 U2 deliverable of the parent meta-plan,
  explicitly deferred per round-2 R3 decision (option C — behavioral_only
  with promote-to-pending sunset_decision). required_for_completion: false
  because the parent plan archived without this executor being a gate;
  promotion to mechanical waits on coverage_evidence machinery.
---

# Validation Scenario Governance Executor

**Status**: `completed`（archived 2026-06-12; F19 promoted to mechanical）
**世代**：Gen 3 Layer 2.5 sub-deliverable
Owner: framework maintainer (linyihong)
**建立日期**：2026-06-01
**Source**: [`plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md`](2026-05-31-2100-mechanical-enforcement-registry.md) §Phase 3 Round-4 T1 / Round-5 U2

> **Stub Notice**：本 plan 為 stub，僅滿足 enforcement-registry `pending_implementation.child_plan_validity` 的最低門檻（Phase 0 outline + owner + Acceptance）。完整 Decision Rationale、Phase Plan 細節、Open Questions 待 implementation 啟動時擴充。

## 為什麼存在

Mechanical Enforcement Registry round-4 評審指出 F19 `validation_scenario_governance` 不應分類為 `behavioral_only`，而是 `pending_implementation`（已有 `coverage_evidence` schema，只是 executor 未實作 = 已知 implementation gap）。

Round-5 U2 要求 `pending_implementation.child_plan` 必須是合法 stub plan（路徑 resolve + Phase 0 heading + owner + Acceptance 區塊），本檔即為此 stub。

## Scope

實作 validation scenario 治理 executor，使 `enforcement/enforcement-registry.yaml §rule_classes[].coverage_evidence` 從「schema 存在但未強制」升級為「compile-time 強制驗證」：

| Source | Target | 驗證內容 |
|---|---|---|
| `validation/scenarios/<area>/*.yaml` 結構 | scenario lint | 必填 `id` / `domain` / `given` / `when` / `then` / `validation.detection_command` |
| `enforcement-registry.yaml` `coverage_evidence.validation_scenarios[]` | 路徑解析 + 存在性 lint | 引用的 scenario yaml 必須存在 |
| `coverage_evidence.coverage_target_pct` | 計算實際覆蓋 vs target | < 50% fail / < 80% warning |
| `coverage_evidence.regression_scenarios[]` | regression-pattern 反查 | 每個 historical failure_pattern 必須對應至少 1 個 regression scenario |

## Phase 0 — Preflight（done 2026-06-12）

- [x] `validation/scenarios/` 目錄結構：13 個 domain 子目錄（app-dev, bootstrap, memory,
  software-delivery, cognitive-modes, runtime, models, apk-analysis, architecture,
  failure-derived, cross-domain, enforcement, engineering）。命名 `<slug>-v<N>.yaml`。
  共 **207** 個 scenario yaml。
- [x] `scenario.schema.json` required = `[id, domain, type, priority, given, when, then]`；
  **`validation` 非 required**。`type` enum = `[routing-decision, heuristic-obedience,
  failure-recovery]` 但實際檔案大量使用 `routing-stability` / `routing-decision` 等其他值
  且 `additionalProperties: true` → schema 目前未被機械強制。
- [x] `enforcement-registry.yaml` `coverage_evidence` 使用盤點：**僅 3 個 rule_class 宣告**
  （`bootstrap_integrity` / `workflow_activation` / F19 本身尚未宣告 coverage_evidence block，
  只有 `executors_planned`）。`validation_scenarios[]` 用於 bootstrap+workflow；
  `regression_scenarios[]` 僅 workflow_activation 宣告 1 條；`coverage_target_pct`：
  bootstrap=100 / workflow=90。
- [x] `enforcement/failure-patterns/` ↔ `regression_scenarios[]` 對應：failure-patterns/ 有
  38 個 md；`validation/scenarios/failure-derived/` 有平行 yaml；registry regression_scenarios
  目前只有 workflow 1 條（travel-planning-regression）。
- [x] **關鍵 preflight 發現（影響 severity 設計）**：
  - 75%（155/207）scenario **無 `validation.detection_command`** — routing scenarios 改用
    `then.validation[]` assertion array。→ `detection_command` 不是 universal invariant。
  - 38 缺 `domain`、10 缺 `given`、13 缺 `when`、15 缺 `then`（散落於舊 corpus）。
  - `bootstrap_integrity.coverage_evidence.validation_scenarios` 引用 **2 個從未建立的檔案**
    （`bootstrap-receipt-required-reads-gate-v1.yaml`、`bootstrap-bypass-on-resume-v1.yaml`）
    — Phase 2 inventory（`1b523fd`）over-promise，real evidence 存在於不同檔名。

## Open Questions（收斂 2026-06-12）

| # | Question | Resolution |
|---|---|---|
| Q1 | 結構 lint 的 scope？全部 207 還是 registry-referenced subset？ | **只 lint registry `coverage_evidence` 引用的 scenarios**（load-bearing governance evidence）。全 corpus 有大量舊格式，blanket FAIL 會打爆 compile；且未引用的 scenario 不是治理證據。 |
| Q2 | `validation.detection_command` 是否 required（FAIL）？ | **No — WARNING**。75% 合法案例沒有此欄位（routing scenarios 用 `then.validation[]`）。把它 required 等於 schema 說謊。core skeleton（id/given/when/then）FAIL；`domain` + `validation.detection_command` 為品質訊號 WARNING。**Plan 1.1 wording 已對應修正**（見下）。 |
| Q3 | `coverage_target_pct` 的「實際覆蓋」如何計算？ | compile-time yaml lint 無法 run scenario 算真實覆蓋率（需 runtime 執行）。改為**驗證宣告值對 governance floor**：mechanical class 宣告 `coverage_target_pct` < 50 → FAIL、50–79 → WARNING、≥80 → OK。對既有 entry 安全（bootstrap=100 / workflow=90 皆 pass）。 |
| Q4 | `regression_scenarios[] ↔ failure-patterns` 雙向 lint 的 reverse 方向（每個 failure pattern 都要有 regression scenario）？ | reverse 全量會要求 38 個 regression（爆量 FAIL）→ **不做全量 reverse FAIL**。forward：regression_scenarios entry 必須存在（FAIL，1.2 子集）+ 應 link 回真實 failure pattern（`failure_source` 或路徑命中 enforcement/failure-patterns 或 failure-derived）否則 WARNING。 |
| Q5 | Bootstrap 2 個 dangling ref 怎麼處理？ | **Create the 2 scenarios**（user decision，first preference B）。兩個名字對應 real designed-but-unlanded 機制（read-log gate / bypass-on-resume validated pattern），非早期草稿。建在 already-referenced 路徑 → bootstrap registry entry **零修改**即 resolve。 |
| Q6 | Validator placement 只在 runtime compile 是否重蹈 `validation-coverage-gap-executor-placement`？ | 本期 scope = wire 進 runtime compile（與既有 `LintEnforcementRegistry` 同 placement）。commit-transaction dual-placement 列為 **follow-up**（已有 precedent: runtime_index_freshness）。記入 plan §Follow-up，不阻塞 F19 promotion。 |

## Phase 1 — Implementation Design（expanded 2026-06-12）

新檔 `scripts/ai-skill-cli/internal/app/scenario_lint.go`，entry symbol **`LintValidationScenarios(repo string) ([]EnforcementRegistryLintError, error)`**
（複用既有 `EnforcementRegistryLintError` finding type + `SeverityFail/SeverityWarn`，
與 registry lint 共用 compile-time summary surface）。registry F19 `executors[].symbol`
宣告為 `LintValidationScenarios`、file=此檔，滿足 promotion R3 `symbol_exists`。

### Phase 1.1 — scenario 結構 lint（scoped to referenced）
- 解析 registry `coverage_evidence.validation_scenarios[]` + `regression_scenarios[]` 收集
  referenced scenario 路徑集合。
- 對每個 referenced scenario yaml：
  - **FAIL**（core BDD skeleton）：缺 `id` / `given` / `when` / `then`。
  - **WARNING**（品質訊號，依 scenario style 合法缺漏）：缺 `domain` / 缺 `validation.detection_command`。
- 檔案不存在不在此檢查報（由 1.2 報 `dangling_coverage_ref`）。

### Phase 1.2 — coverage_evidence 路徑解析 + 存在性（FAIL）
- 每個 `coverage_evidence.validation_scenarios[]` 與 `regression_scenarios[]` 路徑必須 resolve
  到既有檔案，否則 **FAIL** `dangling_coverage_ref`。
- B5-style anchor strip（`path.split('#')[0]`）。

### Phase 1.3 — coverage_target_pct governance floor（FAIL/WARNING）
- 對 `coverage: mechanical` 且宣告 `coverage_evidence.coverage_target_pct` 的 class：
  value < 50 → **FAIL** `coverage_target_below_floor`；50–79 → **WARNING**；≥80 → OK。
- 非數字 / 缺值在 mechanical+有 coverage_evidence 時 → WARNING（建議補）。

### Phase 1.4 — regression_scenarios ↔ failure-patterns 連結（WARNING）
- 每個 `regression_scenarios[]` entry：存在性已由 1.2 保證；額外驗證該 scenario 內容 link 回
  真實 failure pattern（含 `failure_source:` block 或檔內提及 `enforcement/failure-patterns/<x>.md`
  或位於 `validation/scenarios/failure-derived/`），否則 **WARNING** `regression_unlinked_pattern`。

### Phase 1.5 — wire 進 `ai-skill runtime compile`
- `runtime.go` 新增 `buildValidationScenarioLintCheck(repo) (Check, bool)`，mirror
  `buildEnforcementRegistryLintCheck`：FAIL→block compile（ExitValidationFailed）、WARNING→advisory。
- 在 `buildEnforcementRegistryLintCheck` 之後呼叫，獨立 Check name `validation_scenario_lint`。

### Phase 1.6 — tests + rebuild
- `scenario_lint_test.go`：每個 check ≥1 fail + ≥1 pass（table-driven，寫 temp repo fixture）。
- rebuild 5 platform binaries（two-stage commit：source → `releasebuild` bin/）。

## Decision Record（2026-06-12）

- **DR-1 Bootstrap refs = Create（B）**：建立 `bootstrap-receipt-required-reads-gate-v1.yaml`
  （驗證 hooks.go read-log gate 實作 + yaml read_log_requirement）與
  `bootstrap-bypass-on-resume-v1.yaml`（驗證 `resume_exempt: false` + failure pattern + gate row）。
  Both detection_command-bearing → 任何 severity 下 FAIL-clean。
- **DR-2 detection_command = WARNING**，且 Phase 1.1 wording 由「必填」改為 core skeleton
  FAIL + governance field WARNING（Q2）。
- **DR-3 coverage_target_pct = floor 驗證宣告值**（Q3），非 runtime 計算實際覆蓋。

## Decision: Close and Observe（2026-06-12）

**Decision**: Close F19 and observe. 不立即推進任何 follow-up。

**Rationale**: F19 的 **capability gap 已解決**——coverage_evidence 從「schema 存在
但無人驗證」變成「compile-time mechanical enforcement」，閉環成立，且 land 當下已有
real 收益（直接抓出 `bootstrap_integrity` 2 個不存在的 evidence）。下列三個 follow-up
本質都是 **optimization 而非 capability gap**：dual-placement 解的是「**何時**攔截」
（非「**能否**攔截」，已解）；corpus audit 解的是「歷史債務盤點」（非治理能力缺失）；
maturity ladder 解的是「如何**分級**治理」（非是否有治理）。三者收益目前 **未知**——
F19 已獲利但尚未看到下一輪收益訊號，此時加碼不如先觀察。避免「executor 完成 → 立刻
再加一層」導致治理系統自己變成永遠擴張的治理對象。

狀態定位：**Dormant**（非 Rejected）。保留在本 plan + failure pattern + promotion notes。

**Reopen triggers**（事件驅動，呼應 registry `sunset_decision` 哲學）：

| # | 事件 | 重開 follow-up |
|---|---|---|
| 1 | 再次出現 dangling coverage evidence incident | Commit-transaction dual-placement of scenario lint（防 `validation-coverage-gap-executor-placement` family；precedent: runtime_index_freshness） |
| 2 | Scenario corpus health 退化（domain/given/when/then 缺漏持續增加） | 全 corpus（207）結構健康度 audit（非 registry-referenced 的舊 scenario 格式正規化） |
| 3 | 出現對 differentiated scenario quality governance 的需求（"這個 scenario 算成熟嗎？" / "哪些可當 regression gate？"） | Scenario Maturity Ladder M0–M3（見 failure pattern §未來演進與觀察） |

## Promotion Criteria（registry coverage 升級條件）

F19 從 `pending_implementation` 升 `mechanical` 的條件：

- [x] Phase 1.1-1.5 全部完成
- [x] `enforcement-registry.yaml` 中 F19 (`validation_scenario_governance`) 改 `coverage=mechanical` + 補 `executors[]` 區塊
- [x] 至少 1 個 regression scenario 證實能 detect 「scenario 引用不存在 yaml」歷史失效模式
  （`scenario-lint-dangling-coverage-ref-regression-v1.yaml`；採樣 bootstrap_integrity dangling refs）
- [x] `coverage_target_pct` 機制在 ≥3 個 rule_class 上驗證有效
  （bootstrap_integrity=100 / workflow_activation=90 / F19=100 三個 class 走 floor 檢查）

## Validation Plan

- [x] Phase 0 preflight 條目逐一回答（見 §Phase 0 findings）
- [x] Phase 1.1 scenario_lint 有 unit tests（≥1 fail + 1 pass per check；`scenario_lint_test.go` 5 tests pass）
- [x] Phase 1.5 wired 後，跑 `ai-skill runtime compile` 仍 success（PASSED, 0 FAIL / 4 advisory WARNING）
- [x] 升 mechanical 前，registry round-trip：F19 改 mechanical → re-compile → 0 FAIL findings

## Acceptance

- [x] 上述 Validation Plan 全項 checked
- [x] F19 在 `enforcement-registry.yaml` 從 `pending_implementation` 改 `mechanical`
- [x] 本 plan 從 `plans/active/` 移至 `plans/archived/`，並更新 enforcement-registry parent plan 標示 F19 已 promote（archived 2026-06-12; parent note added）

## Dependency Read Ledger

| 欄位 | 內容 |
|---|---|
| Trigger | enforcement-registry round-4 T1 / round-5 U2 要求 F19 stub plan |
| Required set | 本 stub 階段僅需 `enforcement/enforcement-registry.yaml` schema 理解；implementation 啟動時補讀 `validation/scenarios/` 目錄 + `enforcement/failure-patterns/` 索引 |
| Read | enforcement-registry.yaml |
| Deferred | implementation source (Go lint code 結構、existing scenario yaml 結構分析) |
| Validation | 本 stub 通過 enforcement-registry `pending_implementation.child_plan_validity` schema (a-d) |

## Source

2026-06-01 session: enforcement-registry plan round-4 T1 + round-5 U2 觸發。

User 原話 (round-5)：
> F19 變 pending_implementation 後，要確認 schema 是否允許 stub plan。
> 如果 schema 要求 child_plan 必須 active / 必須有 milestone / 必須有 owner，
> 那「開 stub plan 就合法」這個假設要先驗證。

本 stub 即驗證 child_plan_validity 4 條規則的最小實例。
