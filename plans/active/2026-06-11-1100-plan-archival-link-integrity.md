---
id: 2026-06-11-1100-plan-archival-link-integrity
plan_kind: main
status: draft
owner: linyihong
created: 2026-06-11
priority: P2
required_for_completion: false
---

# Plan Archival Link Integrity

**Status**: `draft`
Owner: framework maintainer (linyihong)
**世代**：Gen 3 runtime hardening — Reference Integrity 機械化
**建立日期**：2026-06-11

> 把「plan archive 時 relative-link 斷裂」這個 2026-06-11 親身踩到的失誤模式機械化。

## Decision Rationale

### Problem & Why Now

2026-06-11 archive `2026-06-06-1800-sanitization-mechanical-enforcement`（`active/ → archived/`）時，手動發現 **兩側 relative-link 都斷**：

1. **Inbound**：8 個 active 檔案（README、2 runtime yaml `source_plan`、2 failure-pattern、2 metadata、topology-migration）指向 `plans/active/<id>`，move 後變 stale。
2. **Outbound**：被移動檔案**自身**的 relative-link（parent / sibling plan）原本假設自己在 `plans/active/`，move 到 `plans/archived/` 後 `../archived/`、same-dir、`../active/` 全部要重算。

現有 `validatePlanArchivalAudit` 只檢查 unchecked `- [ ]`，**完全不看 link**。這正是本 session 一直在處理的 **Reference Integrity** family 的一個未覆蓋 instance：move 是 surface-relocation，兩個 surface（檔案位置 ↔ 指向它/它指向的 link）之間 drift。全靠手動 grep + 逐一修，正是「rule 無 executor」風險。

### The failure this prevents

「archive 一個 plan → 連結默默斷掉 → 半年後有人點到 404 / 工具 resolve 失敗」。本 check 在 archive commit 當下就 surface。

## Scope

**In scope**
- 偵測 staged 內 plan 檔案的 `active/ ↔ archived/` rename（git `-M` rename detection）。
- **Outbound check**：解析被移動檔案的 markdown relative-link `](relpath)`，從**新位置**resolve；target 不存在 → finding。
- **Inbound check**：repo 內（active 檔案）是否仍有指向**舊路徑** `plans/active/<id>` / `active/<id>` 的 reference → finding。

**Out of scope**
- 自動修連結（只偵測 + 報告；auto-fix 列未來）。
- 非 plan 檔案的一般 link-rot（更大題目，另議）。
- 絕對 URL / 跨 repo link。

## Phase Plan

### Phase 0 — Design decision（待裁決）

- [ ] **D1 擴充 vs 新 validator**：(a) 擴充 `validatePlanArchivalAudit`（同 obligation `obligation.commit.plan_archival_audit`，less onboarding，但混兩種 concern）vs (b) 新 `validatePlanArchivalLinkIntegrity` + 新 obligation + registry entry（乾淨分離，較多 onboarding）。**傾向 (a)**：同屬 plan-archival concern，輸出分兩類 violation（unchecked / link-breakage）即可。
- [ ] **D2 severity**：outbound（target 不存在，客觀）→ **block**；inbound（舊路徑殘留，可能是歷史 prose 提及）→ 先 **warning**（避免誤殺 provenance 文字），觀察後再考慮 promote。opt-out 沿用 `[skip-plan-archival-audit]`。
- [ ] **D3 rename 偵測來源**：git `diff --cached --find-renames` 取 old→new；或比對 staged archived plan 是否在 HEAD 存在於 active/。

### Phase 1 — Implementation

- [ ] 在 hooks.go（或同檔 helper）實作 link-integrity check，接進 archival audit 路徑
- [ ] Outbound：markdown link 解析 + 從新檔位置 resolve（同 repo 相對路徑），target 不存在 → finding
- [ ] Inbound：掃 active 檔案對舊路徑的殘留 reference → warning finding

### Phase 2 — Tests

- [ ] fail/outbound：move 一個含 `../active/sibling.md`（move 後該相對路徑錯）的 plan → block
- [ ] fail/inbound：另一檔仍指 `plans/active/<moved-id>` → warning
- [ ] pass：move 且所有 inbound/outbound 連結都已更新 → 0 finding
- [ ] pass：純歷史 prose 提及 bare id（非 path link）→ 不誤報

### Phase 3 — Registry & Bootstrap Integration

- [ ] 若採 D1(a)：更新 `enforcement-registry.yaml` `plan_governance`（或 archival 所屬 class）executor/說明；若採 (b)：新 rule_class + core-bootstrap per_commit_obligation
- [ ] failure-pattern：`enforcement/failure-patterns/plan-archival-link-drift.md`（empirical: 2026-06-11 sanitization archive）
- [ ] validation scenarios
- [ ] coverage report + commit/push/readback

## Acceptance

- Archiving a plan with a stale inbound path-link emits a warning; with a broken outbound relative-link emits a block (per D2).
- Clean archive (all links retargeted) passes with zero findings.
- Bare-id provenance mentions do not false-positive.

## Validation

| 欄位 | 內容 |
|---|---|
| Trigger | 2026-06-11 sanitization plan archive 親身踩到 inbound+outbound link 斷裂 |
| Empirical evidence | commit 3f7c4b4（手動修 8 inbound + 3 outbound link） |
| Required set | `scripts/ai-skill-cli/internal/app/hooks.go`（validatePlanArchivalAudit）/ `scan_checkboxes.go` / `enforcement/enforcement-registry.yaml` |
| Deferred | auto-fix；非 plan link-rot |
