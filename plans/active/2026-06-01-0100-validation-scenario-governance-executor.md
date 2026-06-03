---
id: 2026-06-01-0100-validation-scenario-governance-executor
plan_kind: sub
status: draft
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

**Status**: `stub`
**世代**：Gen 3 Layer 2.5 sub-deliverable
Owner: framework maintainer (linyihong)
**建立日期**：2026-06-01
**Source**: [`plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md`](../archived/2026-05-31-2100-mechanical-enforcement-registry.md) §Phase 3 Round-4 T1 / Round-5 U2

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

## Phase 0 — Preflight

- [ ] 確認 `validation/scenarios/` 目錄結構（domain → scenario 命名規則）
- [ ] 確認 `enforcement-registry.yaml` 既有 `coverage_evidence` 使用情況（哪些 rule_class 已宣告、哪些未宣告）
- [ ] 確認 `enforcement/failure-patterns/` 與 `coverage_evidence.regression_scenarios[]` 的對應關係
- [ ] Open Questions: 待 implementation 啟動時收斂

## Phase 1 — Implementation Outline（待 expand）

```
Phase 1.1  scenario yaml structural lint (新檔案 internal/app/scenario_lint.go)
Phase 1.2  coverage_evidence.validation_scenarios 路徑解析 + 存在性 lint
Phase 1.3  coverage_evidence.coverage_target_pct 強制（< 50% fail / < 80% warning）
Phase 1.4  regression_scenarios ↔ failure_patterns 雙向 lint
Phase 1.5  整合進 ai-skill runtime compile
Phase 1.6  rebuild 5 platform binaries
```

詳細 phase 設計待 implementation session 撰寫；本 stub 僅滿足 child_plan validity schema。

## Promotion Criteria（registry coverage 升級條件）

F19 從 `pending_implementation` 升 `mechanical` 的條件：

- [ ] Phase 1.1-1.5 全部完成
- [ ] `enforcement-registry.yaml` 中 F19 (`validation_scenario_governance`) 改 `coverage=mechanical` + 補 `executors[]` 區塊
- [ ] 至少 1 個 regression scenario 證實能 detect 「scenario 引用不存在 yaml」歷史失效模式
- [ ] `coverage_target_pct` 機制在 ≥3 個 rule_class 上驗證有效

## Validation Plan

- [ ] Phase 0 preflight 條目逐一回答
- [ ] Phase 1.1 scenario_lint 有 unit tests（≥1 fail + 1 pass per check）
- [ ] Phase 1.5 wired 後，跑 `ai-skill runtime compile` 仍 success
- [ ] 升 mechanical 前，registry round-trip：F19 改 mechanical → re-dry-run → 0 findings

## Acceptance

- [ ] 上述 Validation Plan 全項 checked
- [ ] F19 在 `enforcement-registry.yaml` 從 `pending_implementation` 改 `mechanical`
- [ ] 本 plan 從 `plans/active/` 移至 `plans/archived/`，並更新 enforcement-registry parent plan 標示 F19 已 promote

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
